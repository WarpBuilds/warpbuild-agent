package manager

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

type ClaudeOptions struct {
	MaxIdle          string `json:"max_idle"`
	Workdir          string `json:"workdir"`
	OutputsDir       string `json:"outputs_dir"`
	StdoutFile       string `json:"stdout_file"`
	StderrFile       string `json:"stderr_file"`
	HostURL          string `json:"host_url"`
	PollingSecret    string `json:"polling_secret"`
	RunnerInstanceID string `json:"runner_instance_id"`
	EnvID            string `json:"env_id"`
	EnvKey           string `json:"env_key"`
	SessionID        string `json:"session_id"`
	WorkID           string `json:"work_id"`
	// CacheBackendHost + RunnerVerificationToken drive the deliverables upload to backend-cache.
	CacheBackendHost        string `json:"cache_backend_host"`
	RunnerVerificationToken string `json:"runner_verification_token"`
}

const anthropicWorkerMaxIdle = "300s"

const (
	claudeWorkspaceDir  = "/workspace"
	claudeOutputsDir    = "/mnt/session/outputs"
	claudeStdoutLogName = "runner.claude.stdout.log"
	claudeStderrLogName = "runner.claude.stderr.log"
)

func DefaultClaudeOptions(maxIdle string) *ClaudeOptions {
	if maxIdle == "" {
		maxIdle = anthropicWorkerMaxIdle
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" || home == "/" {
		home = "/tmp"
	}

	var workdir, outputsDir, stdout, stderr string
	switch runtime.GOOS {
	case "darwin":
		workdir = filepath.Join(home, ".warpbuild", "workspace")
		outputsDir = claudeOutputsDir
		stdout = filepath.Join(home, ".warpbuild", "agent", "log", claudeStdoutLogName)
		stderr = filepath.Join(home, ".warpbuild", "agent", "log", claudeStderrLogName)
	default:
		workdir = claudeWorkspaceDir
		outputsDir = claudeOutputsDir
		stdout = filepath.Join(home, ".warpbuild", "warpbuild-agentd", claudeStdoutLogName)
		stderr = filepath.Join(home, ".warpbuild", "warpbuild-agentd", claudeStderrLogName)
	}

	return &ClaudeOptions{
		MaxIdle:    maxIdle,
		Workdir:    workdir,
		OutputsDir: outputsDir,
		StdoutFile: stdout,
		StderrFile: stderr,
	}
}

func newClaudeManager(opts *ManagerOptions) IManager {
	return &claudeInprocManager{opts: opts}
}

// claudeInprocManager runs the managed-agent session worker in-process (via the
// anthropic-sdk-go SessionToolRunner) instead of exec'ing the `ant` binary, so we
// own the exit/lifecycle policy end to end. On worker exit it fires the same
// post-end hooks (outputs upload + cleanup callback) the ghcri-based path fired.
type claudeInprocManager struct {
	opts *ManagerOptions
}

var _ IManager = &claudeInprocManager{}

func (m *claudeInprocManager) StartRunner(ctx context.Context, opts *StartRunnerOptions) (*StartRunnerOutput, error) {
	c := m.opts.Claude
	if c == nil {
		return nil, fmt.Errorf("claude manager: missing claude options")
	}
	if c.StderrFile != "" {
		_ = os.MkdirAll(filepath.Dir(c.StderrFile), 0o755)
	}
	log.Logger().Infof("Starting in-process Claude worker for session %s", c.SessionID)

	for _, hook := range GetHooks[IPreStartHook]() {
		if err := hook.PreStartHook(ctx, &PreStartHookOptions{StartRunnerOptions: opts, ManagerOptions: m.opts}); err != nil {
			log.Logger().Errorf("error running pre-start hook %s: %v", hook.HookID(), err)
			return nil, err
		}
	}

	if err := runClaudeWorker(ctx, c); err != nil {
		log.Logger().Errorf("claude in-process worker exited with error: %v", err)
	}

	// Post-end hooks (outputs upload + cleanup callback → VM reap), exactly as
	// ghcriManager.StartRunner fired them after the ant exec returned.
	for _, hook := range GetHooks[IPostEndHook]() {
		if err := hook.PostEndHook(ctx, &PostEndHookOptions{StartRunnerOptions: opts, ManagerOptions: m.opts}); err != nil {
			log.Logger().Errorf("error running post-end hook %s: %v", hook.HookID(), err)
		}
	}

	return &StartRunnerOutput{RunCompletedSuccessfully: true}, nil
}

func provisionClaudeWorker(c *ClaudeOptions) error {
	for _, dir := range []string{c.Workdir, c.OutputsDir} {
		if err := ensureWritableDir(dir); err != nil {
			log.Logger().Errorf("Failed to provision claude worker dir %s: %v", dir, err)
			return err
		}
	}
	return nil
}

func ensureWritableDir(dir string) error {
	if err := os.MkdirAll(dir, 0o777); err != nil {
		return err
	}
	return os.Chmod(dir, 0o777)
}
