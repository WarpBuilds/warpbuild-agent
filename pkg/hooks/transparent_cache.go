package hooks

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/manager"
	"github.com/warpbuilds/warpbuild-agent/pkg/tcache"
)

const TRANSPARENT_CACHE_HOOK = "TRANSPARENT_CACHE_HOOK"

// TransparentCacheMarkerEnv tells the job hook (prerun.sh) where the marker file is.
const TransparentCacheMarkerEnv = "WARPBUILD_TRANSPARENT_CACHE_MARKER"

type TransparentCacheHook struct{}

var _ manager.IPreStartHook = &TransparentCacheHook{}

func init() {
	manager.RegisterHook[manager.IPreStartHook](&TransparentCacheHook{})
}

// HookID implements manager.IPreStartHook.
func (*TransparentCacheHook) HookID() string {
	return TRANSPARENT_CACHE_HOOK
}

// PreStartHook patches the runner, starts the cache server and writes the marker file.
// It never returns an error: if anything fails the runner starts stock and uses GitHub's cache.
func (*TransparentCacheHook) PreStartHook(ctx context.Context, opts *manager.PreStartHookOptions) error {
	tc := opts.StartRunnerOptions.AgentOptions.TransparentCache
	if tc == nil || opts.ManagerOptions.Github == nil {
		return nil
	}
	runnerDir := opts.ManagerOptions.Github.RunnerDir

	if err := tcache.PatchRunner(runnerDir); err != nil {
		log.Logger().Errorf("transparent cache: runner not patched, using GitHub's cache: %v", err)
		return nil
	}

	// The runner is patched, so from here the marker must be written: without it nobody sets
	// ACTIONS_RESULTS_URL and artifact calls break. If the server is down the job hook falls
	// back to GitHub's address.
	exe, err := os.Executable()
	if err == nil {
		// A child of the agent, not of the runner, so the runner's end-of-job clean-up cannot kill it.
		cmd := exec.Command(exe, "--launch-transparent-cache", "--settings", tc.SettingsFile)
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		err = cmd.Start()
	}
	if err != nil {
		log.Logger().Errorf("transparent cache: failed to start the cache server: %v", err)
	}

	marker := filepath.Join(runnerDir, ".warpbuild-transparent-cache")
	content := fmt.Sprintf("port=%d\naddress=%s\ngithub_results_url=%s\n", tc.Port, vmAddress(), tc.GitHubResultsURL)
	if err := os.WriteFile(marker, []byte(content), 0644); err != nil {
		log.Logger().Errorf("transparent cache: failed to write the marker file: %v", err)
		return nil
	}
	// The runner inherits the agent's environment, and the job hook inherits the runner's.
	os.Setenv(TransparentCacheMarkerEnv, marker)
	return nil
}

// vmAddress is the address jobs use to reach the cache server: the VM's own IP where the server
// listens on every address, so steps inside containers can reach it.
func vmAddress() string {
	if host := tcache.ListenHost(); host != "" {
		return host
	}
	// No packet is sent: this only asks the kernel which local address the default route uses.
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
