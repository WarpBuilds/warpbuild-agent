package manager

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

const gitSnapshotUploadTimeout = 10 * time.Minute

// snapshotManifest is dropped by the warpbuilds/checkout action on a cache miss. The
// agent uploads the snapshot after the runner exits, so the tar+upload is off the
// customer's billed job time. git_dir is absolute.
type snapshotManifest struct {
	RepoKey string `json:"repo_key"`
	Sha     string `json:"sha"`
	GitDir  string `json:"git_dir"`
}

// uploadGitSnapshots uploads any checkout snapshots the finished job recorded. Manifests
// live under <runnerDir>/_work/_temp/wb-snapshots (RUNNER_TEMP as the action sees it).
// Best-effort: every failure logs and is skipped so it never affects the (already-done)
// job or blocks VM teardown. Reuses the runner's own WARPBUILD_HOST_URL + verification
// token — the same endpoint + credential the in-job action would have used.
func uploadGitSnapshots(ctx context.Context, runnerDir string) {
	dir := filepath.Join(runnerDir, "_work", "_temp", "wb-snapshots")
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) == 0 {
		return // no misses recorded
	}
	hostURL := os.Getenv("WARPBUILD_HOST_URL")
	token := os.Getenv("WARPBUILD_RUNNER_VERIFICATION_TOKEN")
	if hostURL == "" || token == "" {
		log.Logger().Infof("[git-snapshot] WARPBUILD host/token not in env; skipping %d manifest(s)", len(entries))
		return
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		m, err := readSnapshotManifest(filepath.Join(dir, e.Name()))
		if err != nil {
			log.Logger().Errorf("[git-snapshot] bad manifest %s: %v", e.Name(), err)
			continue
		}
		if err := uploadOneGitSnapshot(ctx, hostURL, token, m); err != nil {
			log.Logger().Errorf("[git-snapshot] upload failed for %s: %v", m.Sha, err)
			continue
		}
		log.Logger().Infof("[git-snapshot] uploaded snapshot for %s", m.Sha)
	}
}

func readSnapshotManifest(path string) (snapshotManifest, error) {
	var m snapshotManifest
	b, err := os.ReadFile(path)
	if err != nil {
		return m, err
	}
	if err := json.Unmarshal(b, &m); err != nil {
		return m, err
	}
	if m.RepoKey == "" || m.Sha == "" || m.GitDir == "" {
		return m, fmt.Errorf("incomplete manifest")
	}
	return m, nil
}

func uploadOneGitSnapshot(ctx context.Context, hostURL, token string, m snapshotManifest) error {
	uploadCtx, cancel := context.WithTimeout(ctx, gitSnapshotUploadTimeout)
	defer cancel()

	archivePath, err := tarGitObjects(m.GitDir)
	if err != nil {
		return fmt.Errorf("archive: %w", err)
	}
	defer os.Remove(archivePath)

	url, err := requestGitMirrorUploadURL(uploadCtx, hostURL, token, m.RepoKey, m.Sha)
	if err != nil {
		return err
	}
	if url == "" {
		return nil // disabled org or another job holds the lock — skip quietly
	}
	return putFile(uploadCtx, url, archivePath)
}

// tarGitObjects writes an uncompressed tar of gitDir/objects (+ shallow) with Go's
// archive/tar — no `tar` binary, works on every OS. Member names are objects/... and
// shallow, exactly what the checkout action's restore (`tar -xf` into .git) expects.
func tarGitObjects(gitDir string) (string, error) {
	f, err := os.CreateTemp("", "wb-git-snapshot-*.tar")
	if err != nil {
		return "", err
	}
	defer f.Close()
	tw := tar.NewWriter(f)

	if err := addTreeToTar(tw, filepath.Join(gitDir, "objects"), "objects"); err != nil {
		_ = os.Remove(f.Name())
		return "", err
	}
	if shallow := filepath.Join(gitDir, "shallow"); fileExists(shallow) {
		if err := addFileToTar(tw, shallow, "shallow"); err != nil {
			_ = os.Remove(f.Name())
			return "", err
		}
	}
	if err := tw.Close(); err != nil {
		_ = os.Remove(f.Name())
		return "", err
	}
	return f.Name(), nil
}

func addTreeToTar(tw *tar.Writer, root, prefix string) error {
	return filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		name := prefix
		if rel != "." {
			name = prefix + "/" + filepath.ToSlash(rel)
		}
		if info.IsDir() {
			return tw.WriteHeader(&tar.Header{Name: name + "/", Mode: 0o755, Typeflag: tar.TypeDir})
		}
		if !info.Mode().IsRegular() {
			return nil // git object stores hold only regular files; skip anything else
		}
		return writeRegularFile(tw, p, name, info.Size())
	})
}

func addFileToTar(tw *tar.Writer, path, name string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	return writeRegularFile(tw, path, name, info.Size())
}

func writeRegularFile(tw *tar.Writer, path, name string, size int64) error {
	if err := tw.WriteHeader(&tar.Header{Name: name, Mode: 0o644, Size: size, Typeflag: tar.TypeReg}); err != nil {
		return err
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(tw, f)
	return err
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// requestGitMirrorUploadURL asks backend-core for a presigned PUT URL, hitting the same
// endpoint + verification token the in-job action used. "" (skip) on 403 disabled / 409
// locked — the org has no cache stack or another job is already uploading this sha.
func requestGitMirrorUploadURL(ctx context.Context, hostURL, token, repoKey, sha string) (string, error) {
	body, err := json.Marshal(map[string]string{"repo_key": repoKey, "sha": sha})
	if err != nil {
		return "", err
	}
	endpoint := strings.TrimRight(hostURL, "/") + "/api/v1/git-mirrors/upload-url"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusConflict {
		return "", nil
	}
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload-url returned %s: %s", resp.Status, string(b))
	}
	var out struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if out.URL == "" {
		return "", fmt.Errorf("empty upload url")
	}
	return out.URL, nil
}
