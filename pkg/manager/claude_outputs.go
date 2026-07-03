package manager

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

const outputsUploadTimeout = 5 * time.Minute

// uploadSessionOutputs archives OutputsDir and uploads it to the backend-issued presigned URL when the
// managed-agent worker exits. Best-effort: on any error it logs and returns so a failed upload never
// fails the run. Runs synchronously before the cleanup hook so the VM isn't reaped mid-upload.
func uploadSessionOutputs(ctx context.Context, outputsDir, hostURL, pollingSecret, runnerInstanceID string) {
	if outputsDir == "" || hostURL == "" || runnerInstanceID == "" {
		return
	}
	entries, err := os.ReadDir(outputsDir)
	if err != nil {
		log.Logger().Infof("[claude_outputs] no outputs dir %s to upload (%v)", outputsDir, err)
		return
	}
	if len(entries) == 0 {
		log.Logger().Infof("[claude_outputs] outputs dir %s is empty, skipping upload", outputsDir)
		return
	}

	uploadCtx, cancel := context.WithTimeout(ctx, outputsUploadTimeout)
	defer cancel()

	archivePath, err := tarGzDir(uploadCtx, outputsDir)
	if err != nil {
		log.Logger().Errorf("[claude_outputs] failed to archive %s: %v", outputsDir, err)
		return
	}
	defer os.Remove(archivePath)

	url, err := requestOutputUploadURL(uploadCtx, hostURL, pollingSecret, runnerInstanceID)
	if err != nil {
		log.Logger().Errorf("[claude_outputs] failed to get output upload URL: %v", err)
		return
	}
	if err := putFile(uploadCtx, url, archivePath); err != nil {
		log.Logger().Errorf("[claude_outputs] failed to upload session outputs: %v", err)
		return
	}
	log.Logger().Infof("[claude_outputs] uploaded session outputs from %s", outputsDir)
}

// tarGzDir tars+gzips the contents of dir to a temp file using the system `tar` (present on the Linux/
// macOS sandbox hosts), streaming to disk so the whole archive is never held in memory. Returns the
// temp file path; the caller must remove it. `-C dir .` archives dir's contents with relative paths.
func tarGzDir(ctx context.Context, dir string) (string, error) {
	f, err := os.CreateTemp("", "warpbuild-outputs-*.tar.gz")
	if err != nil {
		return "", err
	}
	path := f.Name()
	_ = f.Close() // tar (re)creates it

	cmd := exec.CommandContext(ctx, "tar", "-czf", path, "-C", dir, ".")
	if out, cerr := cmd.CombinedOutput(); cerr != nil {
		_ = os.Remove(path)
		return "", fmt.Errorf("tar failed: %v: %s", cerr, strings.TrimSpace(string(out)))
	}
	return path, nil
}

// requestOutputUploadURL asks backend-core for a presigned PUT URL, authenticating with the runner
// polling secret (same credential the agent uses for allocation_details / cleanup_hook).
func requestOutputUploadURL(ctx context.Context, hostURL, pollingSecret, runnerInstanceID string) (string, error) {
	body, err := json.Marshal(map[string]string{"runner_instance_id": runnerInstanceID})
	if err != nil {
		return "", err
	}
	endpoint := strings.TrimRight(hostURL, "/") + "/sandboxes/outputs/upload-url"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("X-POLLING-SECRET", pollingSecret)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
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
		return "", fmt.Errorf("empty upload url in response")
	}
	return out.URL, nil
}

// putFile streams a file to a presigned S3 PUT URL (bounded memory — the archive is read from disk, not
// buffered).
func putFile(ctx context.Context, url, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, f)
	if err != nil {
		return err
	}
	req.ContentLength = info.Size()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("PUT returned %s: %s", resp.Status, string(b))
	}
	return nil
}
