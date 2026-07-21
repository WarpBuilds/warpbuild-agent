package manager

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

const (
	// anthropicWorkerBinary is the Anthropic CLI that runs the self-hosted worker
	// (`ant beta:worker run`). workerAssetForPlatform appends .exe on Windows.
	anthropicWorkerBinary = "ant"

	// defaultAnthropicWorkerVersion is the pinned `ant` release installed on sandbox VMs.
	// Bump it (and re-verify the asset names) to take a newer worker. Overridable per-VM via
	// WARPBUILD_ANTHROPIC_WORKER_VERSION without rebuilding the agent.
	defaultAnthropicWorkerVersion = "1.16.0"
	anthropicWorkerVersionEnv     = "WARPBUILD_ANTHROPIC_WORKER_VERSION"

	anthropicCLIDownloadBase = "https://github.com/anthropics/anthropic-cli/releases/download"

	extTarGz = "tar.gz"
	extZip   = "zip"
)

// workerAsset describes the release asset for the running platform. Anthropic publishes the
// CLI as ant_<version>_<os>_<arch>.<ext> — note the OS token is `macos` (not Go's `darwin`),
// and only Linux ships a tarball; macOS and Windows ship zips (with an ant.exe on Windows).
type workerAsset struct {
	osToken string // linux | macos | windows
	arch    string // amd64 | arm64
	ext     string // extTarGz | extZip
	binary  string // ant | ant.exe
}

// ensureAnthropicWorkerInstalled installs Anthropic's `ant` CLI on demand and returns its path.
// Neither the generic runner image nor cloud-init ships it — cloud-init only pre-creates the worker
// dirs and swaps agentd — so the agent installs it here, idempotently (reused if a prior run did).
func ensureAnthropicWorkerInstalled(ctx context.Context) (string, error) {
	asset, err := workerAssetForPlatform()
	if err != nil {
		return "", err
	}
	dst := filepath.Join(anthropicWorkerInstallDir(), asset.binary)

	// Idempotent: reuse the binary WE installed on a prior run/retry. We deliberately do NOT use
	// exec.LookPath("ant") — on the generic runner image /usr/bin/ant is Apache Ant (the Java build
	// tool). Invoking that with `beta:worker run …` just prints Ant's usage and exits, so the worker
	// never connects. Always resolve the Anthropic CLI by our own install path instead.
	if _, statErr := os.Stat(dst); statErr == nil {
		log.Logger().Infof("anthropic worker CLI already installed at %s", dst)
		return dst, nil
	}

	version := resolveAnthropicWorkerVersion()
	url := fmt.Sprintf("%s/v%s/ant_%s_%s_%s.%s", anthropicCLIDownloadBase, version, version, asset.osToken, asset.arch, asset.ext)

	log.Logger().Infof("installing anthropic worker CLI v%s from %s", version, url)

	archivePath, err := downloadToTempFile(ctx, url)
	if err != nil {
		return "", fmt.Errorf("download anthropic worker CLI: %w", err)
	}
	defer os.Remove(archivePath)

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return "", err
	}

	switch asset.ext {
	case extTarGz:
		err = extractFromTarGz(archivePath, asset.binary, dst)
	case extZip:
		err = extractFromZip(archivePath, asset.binary, dst)
	default:
		err = fmt.Errorf("unknown archive format %q", asset.ext)
	}
	if err != nil {
		return "", fmt.Errorf("install anthropic worker CLI: %w", err)
	}

	log.Logger().Infof("anthropic worker CLI v%s installed to %s", version, dst)
	return dst, nil
}

// workerAssetForPlatform maps the running OS/arch to its release asset shape.
func workerAssetForPlatform() (workerAsset, error) {
	a := workerAsset{arch: runtime.GOARCH}
	switch runtime.GOARCH {
	case "amd64", "arm64":
	default:
		return a, fmt.Errorf("anthropic worker auto-install is not supported on %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	switch runtime.GOOS {
	case "linux":
		a.osToken, a.ext, a.binary = "linux", extTarGz, anthropicWorkerBinary
	case "darwin":
		a.osToken, a.ext, a.binary = "macos", extZip, anthropicWorkerBinary
	case "windows":
		a.osToken, a.ext, a.binary = "windows", extZip, anthropicWorkerBinary+".exe"
	default:
		return a, fmt.Errorf("anthropic worker auto-install is not supported on %s", runtime.GOOS)
	}
	return a, nil
}

// resolveAnthropicWorkerVersion returns the pinned version, or an environment override
// (a leading "v" is tolerated).
func resolveAnthropicWorkerVersion() string {
	if v := strings.TrimSpace(os.Getenv(anthropicWorkerVersionEnv)); v != "" {
		return strings.TrimPrefix(v, "v")
	}
	return defaultAnthropicWorkerVersion
}

// anthropicWorkerInstallDir is a writable, conventional bin dir per OS. The worker is invoked by its
// absolute install path, so this dir does not need to be on PATH. It must be writable by the non-root
// `runner` the worker runs as, so on unix we use a dir under the user's home rather than /usr/local/bin
// (which only root can write). Falls back to TempDir if HOME is unavailable.
func anthropicWorkerInstallDir() string {
	if runtime.GOOS == "windows" {
		if pd := os.Getenv("ProgramData"); pd != "" {
			return filepath.Join(pd, "warpbuild", "bin")
		}
		return filepath.Join(os.TempDir(), "warpbuild", "bin")
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".warpbuild", "bin")
	}
	return filepath.Join(os.TempDir(), "warpbuild", "bin")
}

func downloadToTempFile(ctx context.Context, url string) (string, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("downloading %s returned %s", url, resp.Status)
	}

	f, err := os.CreateTemp("", "ant-archive-*")
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(f, resp.Body); err != nil {
		os.Remove(f.Name())
		return "", err
	}
	return f.Name(), nil
}

// extractFromTarGz writes the entry whose base name is binary into dst, executable.
func extractFromTarGz(archivePath, binary, dst string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		// Archives may be flat (ant) or wrapped in a dir (ant_x/ant); match on the base name.
		if hdr.Typeflag == tar.TypeReg && filepath.Base(hdr.Name) == binary {
			return writeExecutable(dst, tr)
		}
	}
	return fmt.Errorf("archive %s did not contain %q", archivePath, binary)
}

// extractFromZip writes the entry whose base name is binary into dst, executable.
func extractFromZip(archivePath, binary, dst string) error {
	zr, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer zr.Close()

	for _, zf := range zr.File {
		if zf.FileInfo().IsDir() || filepath.Base(zf.Name) != binary {
			continue
		}
		rc, err := zf.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		return writeExecutable(dst, rc)
	}
	return fmt.Errorf("archive %s did not contain %q", archivePath, binary)
}

// writeExecutable writes r to a temp file in dst's directory, then atomically renames it into
// place with an executable mode so a partial/concurrent write never leaves a broken binary.
func writeExecutable(dst string, r io.Reader) error {
	tmp, err := os.CreateTemp(filepath.Dir(dst), ".ant-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if _, err := io.Copy(tmp, r); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpName, 0o755); err != nil {
		return err
	}
	return os.Rename(tmpName, dst)
}
