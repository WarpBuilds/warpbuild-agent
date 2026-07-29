package manager

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

// ClaudeOptions configures the Anthropic managed-agent worker a claude_agent sandbox VM runs
type ClaudeOptions struct {
	Command          string   `json:"command"`
	Args             []string `json:"args"`
	Workdir          string   `json:"workdir"`
	OutputsDir       string   `json:"outputs_dir"`
	StdoutFile       string   `json:"stdout_file"`
	StderrFile       string   `json:"stderr_file"`
	HostURL          string   `json:"host_url"`
	PollingSecret    string   `json:"polling_secret"`
	RunnerInstanceID string   `json:"runner_instance_id"`
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

type claudeManager struct {
	*ClaudeOptions
}

var _ IManager = (*claudeManager)(nil)

func NewClaudeManager(opts *ClaudeOptions) IManager {
	return &claudeManager{ClaudeOptions: opts}
}

func (m *claudeManager) StartRunner(ctx context.Context, opts *StartRunnerOptions) (*StartRunnerOutput, error) {
	if err := m.provisionWorkerFiles(); err != nil {
		return nil, err
	}

	command := m.Command
	if command == anthropicWorkerBinary {
		antPath := anthropicWorkerPath()
		if _, err := os.Stat(antPath); err != nil {
			return nil, fmt.Errorf("anthropic worker CLI not found at %s (cloud-init installs it on claude_agent sandbox VMs): %w", antPath, err)
		}
		command = antPath
	}

	cmd := exec.CommandContext(ctx, command, m.Args...)
	if m.Workdir != "" {
		cmd.Dir = m.Workdir
	}

	log.Logger().Infof("starting claude managed-agent worker with command: %s", cmd.String())

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Logger().Errorf("error creating stdout pipe: %v", err)
		return nil, err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Logger().Errorf("error creating stderr pipe: %v", err)
		return nil, err
	}

	managerOpts := &ManagerOptions{Provider: ProviderClaudeAgent}

	for _, hook := range GetHooks[IPreStartHook]() {
		if err := hook.PreStartHook(ctx, &PreStartHookOptions{
			StartRunnerOptions: opts,
			ManagerOptions:     managerOpts,
		}); err != nil {
			log.Logger().Errorf("error running pre-start hook %s: %v", hook.HookID(), err)
			return nil, err
		}
	}

	if err := cmd.Start(); err != nil {
		log.Logger().Errorf("error starting command: %v", err)
		return nil, err
	}

	stdoutChan := make(chan string)
	stderrChan := make(chan string)
	doneChan := make(chan bool)

	go captureOutput(stdoutPipe, stdoutChan)
	go captureOutput(stderrPipe, stderrChan)

	stdoutFile, err := openFile(m.StdoutFile)
	if err != nil {
		log.Logger().Errorf("error opening stdout file: %v", err)
		return nil, err
	}
	defer stdoutFile.Close()

	stderrFile, err := openFile(m.StderrFile)
	if err != nil {
		log.Logger().Errorf("error opening stderr file: %v", err)
		return nil, err
	}
	defer stderrFile.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		cmd.Wait()
		doneChan <- true
	}()

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case out := <-stdoutChan:
			fmt.Fprintln(stdoutFile, out)
			fmt.Fprintln(os.Stdout, out)
		case errLine := <-stderrChan:
			fmt.Fprintln(stderrFile, errLine)
			fmt.Fprintln(os.Stderr, errLine)
		case <-ticker.C:
		case <-doneChan:
			wg.Wait()

			// Upload the session deliverables (/mnt/session/outputs) to S3 BEFORE the post-end
			// hooks run: the cleanup hook reaps this single-use VM. Synchronous + best-effort.
			uploadSessionOutputs(ctx, m.OutputsDir, m.HostURL, m.PollingSecret, m.RunnerInstanceID)
			// The worker (`ant beta:worker run`) exited — the session ended (end_turn + --max-idle) or
			// the lease was lost. Run the post-end hooks: the cleanup hook calls the backend cleanup_hook
			// → RemoveRunner, which reaps this single-use VM. Then return.
			for _, hook := range GetHooks[IPostEndHook]() {
				if err := hook.PostEndHook(ctx, &PostEndHookOptions{
					StartRunnerOptions: opts,
					ManagerOptions:     managerOpts,
				}); err != nil {
					log.Logger().Errorf("error running post-end hook %s: %v", hook.HookID(), err)
				}
			}

			return &StartRunnerOutput{RunCompletedSuccessfully: true}, nil
		}
	}
}

// provisionWorkerFiles ensures the worker's dirs (workspace, session outputs, log dir) and log files exist.
func (m *claudeManager) provisionWorkerFiles() error {
	// ensureWritableDir is idempotent, so overlapping dirs are fine.
	for _, dir := range []string{m.Workdir, m.OutputsDir, filepath.Dir(m.StdoutFile), filepath.Dir(m.StderrFile)} {
		if err := ensureWritableDir(dir); err != nil {
			log.Logger().Errorf("Failed to provision claude worker dir %s: %v", dir, err)
			return err
		}
	}
	for _, fullPath := range []string{m.StderrFile, m.StdoutFile} {
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			f, err := os.Create(fullPath)
			if err != nil {
				log.Logger().Errorf("Failed to create file %s: %v", fullPath, err)
				return err
			}
			f.Close()
		}
	}
	return nil
}

// ensureWritableDir ensures /workspace and /mnt/session/outputs (the worker's workdir + session deliverables) exists and is 0777
func ensureWritableDir(dir string) error {
	// MkdirAll is a no-op on the cloud-init-pre-created session dirs (and creates the log dir where agentd
	// has permission); its mode is umask-masked, so chmod explicitly to guarantee 0777 for the non-root
	// worker. 0777 is safe on a single-use, single-tenant sandbox VM.
	if err := os.MkdirAll(dir, 0o777); err != nil {
		return err
	}
	return os.Chmod(dir, 0o777)
}
