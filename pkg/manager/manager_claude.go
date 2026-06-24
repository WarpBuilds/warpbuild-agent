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

// ClaudeOptions configures the Anthropic managed-agent worker a claude_agent sandbox VM
// runs. The VM boots from the generic runner image, so there is no claude block in
// settings.json — DefaultClaudeOptions supplies sensible defaults. The Anthropic identity
// (ANTHROPIC_ENVIRONMENT_ID/KEY/SESSION_ID) is exported into the process environment by the
// agent before StartRunner, so the worker inherits it.
type ClaudeOptions struct {
	Command    string   `json:"command"`
	Args       []string `json:"args"`
	Workdir    string   `json:"workdir"`
	StdoutFile string   `json:"stdout_file"`
	StderrFile string   `json:"stderr_file"`
}

// anthropicWorkerMaxIdle is the FALLBACK --max-idle for `ant beta:worker run`, used only when the
// backend's allocation_details doesn't supply one (e.g. an older backend). Normally the value comes
// from the backend (sandbox.idle_ttl_seconds) via DefaultClaudeOptions, so the worker self-terminates
// on the same threshold the idle-TTL reaper backstops. ant's own default is 60s.
const anthropicWorkerMaxIdle = "300s"

// DefaultClaudeOptions runs `ant beta:worker run --workdir /workspace --max-idle <maxIdle>`. maxIdle
// is supplied by the backend's allocation_details (sandbox.idle_ttl_seconds), which the agent polls
// before starting the worker; empty falls back to anthropicWorkerMaxIdle. `ant` (the Anthropic CLI)
// is installed on demand by StartRunner (see ensureAnthropicWorkerInstalled), so it does not need to
// be baked into the runner image.
func DefaultClaudeOptions(maxIdle string) *ClaudeOptions {
	if maxIdle == "" {
		maxIdle = anthropicWorkerMaxIdle
	}
	workdir := "/workspace"
	stdout := "/var/log/warpbuild-agentd/runner.claude.stdout.log"
	stderr := "/var/log/warpbuild-agentd/runner.claude.stderr.log"
	if runtime.GOOS == "windows" {
		workdir = `C:\workspace`
		stdout = `C:\ProgramData\warpbuild\logs\runner.claude.stdout.log`
		stderr = `C:\ProgramData\warpbuild\logs\runner.claude.stderr.log`
	}
	return &ClaudeOptions{
		Command:    anthropicWorkerBinary,
		Args:       []string{"beta:worker", "run", "--workdir", workdir, "--max-idle", maxIdle},
		Workdir:    workdir,
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
	if err := m.createFiles(); err != nil {
		return nil, err
	}

	// claude_agent sandbox VMs boot the generic runner image, which does not ship the
	// Anthropic CLI. Install it on demand (a single cross-OS path) and run it by absolute
	// path. A non-default Command is treated as an explicit override and left untouched.
	command := m.Command
	if command == anthropicWorkerBinary {
		antPath, err := ensureAnthropicWorkerInstalled(ctx)
		if err != nil {
			return nil, err
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

			// Run post-end hooks (the cleanup hook tears the VM down) — same lifecycle
			// as the github manager.
			for _, hook := range GetHooks[IPostEndHook]() {
				if err := hook.PostEndHook(ctx, &PostEndHookOptions{
					StartRunnerOptions: opts,
					ManagerOptions:     managerOpts,
				}); err != nil {
					log.Logger().Errorf("error running post-end hook %s: %v", hook.HookID(), err)
				}
			}

			return &StartRunnerOutput{
				RunCompletedSuccessfully: true,
			}, nil
		}
	}
}

func (m *claudeManager) createFiles() error {
	if m.Workdir != "" {
		if err := os.MkdirAll(m.Workdir, 0755); err != nil {
			log.Logger().Errorf("Failed to create workdir %s: %v", m.Workdir, err)
			return err
		}
	}
	for _, fullPath := range []string{m.StderrFile, m.StdoutFile} {
		baseDir := filepath.Dir(fullPath)
		if _, err := os.Stat(baseDir); os.IsNotExist(err) {
			if err := os.MkdirAll(baseDir, 0755); err != nil {
				log.Logger().Errorf("Failed to create base directory %s: %v", baseDir, err)
				return err
			}
		}
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
