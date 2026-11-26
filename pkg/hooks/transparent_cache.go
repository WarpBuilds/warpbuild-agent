package hooks

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/manager"
	"github.com/warpbuilds/warpbuild-agent/pkg/tcache"
)

const TRANSPARENT_CACHE_HOOK = "TRANSPARENT_CACHE_HOOK"

// The marker file sits next to the job hook (prerun.sh), which reads it. It exists only when the
// agent has patched the runner.
const transparentCacheMarker = ".warpbuild-transparent-cache"

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
	if tc == nil {
		return nil
	}
	var runnerDir string
	var envs manager.EnvironmentVariables
	switch mo := opts.ManagerOptions; {
	case mo.Github != nil:
		runnerDir, envs = mo.Github.RunnerDir, mo.Github.Envs
	case mo.GithubCRI != nil && mo.GithubCRI.CMDOptions != nil:
		runnerDir, envs = mo.GithubCRI.RunnerDir, mo.GithubCRI.CMDOptions.Envs
	default:
		return nil
	}
	var jobHook string
	for _, env := range envs {
		if env.Key == "ACTIONS_RUNNER_HOOK_JOB_STARTED" {
			jobHook = env.Value
		}
	}
	if jobHook == "" {
		log.Logger().Errorf("transparent cache: no job hook configured, using GitHub's cache")
		return nil
	}

	// Marker first, patch second. A patched runner with no marker breaks artifact calls, because
	// nobody sets ACTIONS_RESULTS_URL. A marker with a stock runner is harmless: the runner
	// overwrites what the job hook sets.
	marker := filepath.Join(filepath.Dir(jobHook), transparentCacheMarker)
	content := fmt.Sprintf("port=%d\naddress=%s\nrunner_dir=%s\ngithub_results_url=%s\n", tc.Port, tcache.VMAddress(), runnerDir, tc.GitHubResultsURL)
	if err := os.WriteFile(marker, []byte(content), 0644); err != nil {
		log.Logger().Errorf("transparent cache: marker not written, using GitHub's cache: %v", err)
		return nil
	}
	if err := tcache.PatchRunner(runnerDir); err != nil {
		log.Logger().Errorf("transparent cache: runner not patched, using GitHub's cache: %v", err)
		os.Remove(marker)
		return nil
	}

	// The server is a child of the agent, not of the runner, so the runner's end-of-job clean-up
	// cannot kill it. Every artifact call of the job goes through it, so it is started again if it
	// ever stops. While it is down at job start, the job hook falls back to GitHub.
	if exe, err := os.Executable(); err != nil {
		log.Logger().Errorf("transparent cache: failed to start the cache server: %v", err)
	} else {
		go func() {
			for {
				cmd := exec.Command(exe, "--launch-transparent-cache", "--settings", tc.SettingsFile)
				cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
				log.Logger().Errorf("transparent cache: the cache server stopped, starting it again: %v", cmd.Run())
				time.Sleep(time.Second)
			}
		}()
	}
	log.Logger().Infof("transparent cache: runner patched, marker at %s", marker)
	return nil
}
