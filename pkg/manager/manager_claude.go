package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

// ClaudeOptions configures the Anthropic managed-agent worker a claude_agent sandbox VM runs
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
}

const anthropicWorkerMaxIdle = "300s"

// Claude worker filesystem layout. Workdir + session-deliverables dir vary per OS; the log leaf
// names are shared and joined under an OS-specific dir at runtime.
const (
	claudeWorkspaceDir  = "/workspace"           // linux worker workdir
	claudeOutputsDir    = "/mnt/session/outputs" // session deliverables (linux + darwin; firmlinked on macOS)
	claudeStdoutLogName = "runner.claude.stdout.log"
	claudeStderrLogName = "runner.claude.stderr.log"
)

// DefaultClaudeOptions runs `ant beta:worker run --workdir <workdir> --max-idle <maxIdle>`
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
		outputsDir = claudeOutputsDir // intentionally the linux path — macOS firmlinks /mnt/session/outputs
		stdout = filepath.Join(home, ".warpbuild", "agent", "log", claudeStdoutLogName)
		stderr = filepath.Join(home, ".warpbuild", "agent", "log", claudeStderrLogName)
	default: // linux
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

// resolvedCommand is the absolute path the worker is invoked by. The `ant` binary is resolved to the
// path cloud-init installs it to (never via PATH — the generic runner image ships a different `ant`).
func (c *ClaudeOptions) resolvedCommand() string {
	if c.Command == "" || c.Command == anthropicWorkerBinary {
		return anthropicWorkerPath()
	}
	return c.Command
}

// newClaudeManager runs the Anthropic worker through the shared command runner (ghcriManager): the
// worker is just a command (CMDOptions, env via CMDOptions.Envs), and the outputs upload + cleanup run
// as post-end hooks. provisionClaudeWorker must have run first (the agent calls it before NewManager).
func newClaudeManager(opts *ManagerOptions) IManager {
	c := opts.Claude
	return &ghcriManager{
		GithubCRIOptions: &GithubCRIOptions{
			StdoutFile:       c.StdoutFile,
			StderrFile:       c.StderrFile,
			InheritParentEnv: true, // the ant worker needs PATH/HOME from the agent process
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

// provisionClaudeWorker verifies the ant CLI is installed and creates the worker's writable dirs
// (workspace + session outputs) before the worker starts. The shared runner creates the agentd log
// files; these two dirs are what the non-root worker itself writes to.
func provisionClaudeWorker(c *ClaudeOptions) error {
	antPath := c.resolvedCommand()
	if _, err := os.Stat(antPath); err != nil {
		return fmt.Errorf("anthropic worker CLI not found at %s (cloud-init installs it on claude_agent sandbox VMs): %w", antPath, err)
	}
	for _, dir := range []string{c.Workdir, c.OutputsDir} {
		if err := ensureWritableDir(dir); err != nil {
			log.Logger().Errorf("Failed to provision claude worker dir %s: %v", dir, err)
			return err
		}
	}
	return nil
}

// ensureWritableDir ensures a worker dir (/workspace, /mnt/session/outputs) exists and is 0777.
func ensureWritableDir(dir string) error {
	// MkdirAll is a no-op on the cloud-init-pre-created session dirs; its mode is umask-masked, so chmod
	// explicitly to guarantee 0777 for the non-root worker. 0777 is safe on a single-use sandbox VM.
	if err := os.MkdirAll(dir, 0o777); err != nil {
		return err
	}
	return os.Chmod(dir, 0o777)
}
