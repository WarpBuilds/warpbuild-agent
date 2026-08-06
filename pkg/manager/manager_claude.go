package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

type ClaudeOptions struct {
	Command          string               `json:"command"`
	Args             []string             `json:"args"`
	Workdir          string               `json:"workdir"`
	OutputsDir       string               `json:"outputs_dir"`
	StdoutFile       string               `json:"stdout_file"`
	StderrFile       string               `json:"stderr_file"`
	Envs             EnvironmentVariables `json:"envs"`
	HostURL          string               `json:"host_url"`
	PollingSecret    string               `json:"polling_secret"`
	RunnerInstanceID string               `json:"runner_instance_id"`
	SessionID        string               `json:"session_id"`
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
		Command:    anthropicWorkerBinary,
		Args:       []string{"beta:worker", "run", "--workdir", workdir, "--max-idle", maxIdle},
		Workdir:    workdir,
		OutputsDir: outputsDir,
		StdoutFile: stdout,
		StderrFile: stderr,
	}
}

func (c *ClaudeOptions) resolvedCommand() string {
	if c.Command == "" || c.Command == anthropicWorkerBinary {
		return anthropicWorkerPath()
	}
	return c.Command
}

func newClaudeManager(opts *ManagerOptions) IManager {
	c := opts.Claude
	return &ghcriManager{
		GithubCRIOptions: &GithubCRIOptions{
			StdoutFile:       c.StdoutFile,
			StderrFile:       c.StderrFile,
			InheritParentEnv: true,
			CMDOptions: &CMDOptions{
				CMD:  c.resolvedCommand(),
				Args: c.Args,
				Dir:  c.Workdir,
				Envs: c.Envs,
			},
		},
		provider:    ProviderClaudeAgent,
		managerOpts: opts,
	}
}

func provisionClaudeWorker(c *ClaudeOptions) error {
	antPath := c.resolvedCommand()
	if _, err := os.Stat(antPath); err != nil {
		return fmt.Errorf("anthropic worker CLI not found at %s (cloud-init installs it on claude_agent VMs): %w", antPath, err)
	}
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
