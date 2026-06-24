package manager

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
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
	OutputsDir string   `json:"outputs_dir"`
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
	// /mnt/session/outputs is where the worker harness has Claude write final deliverables (the docs'
	// system default for self-hosted sandbox mode). Linux-pathed; left empty on Windows.
	outputsDir := "/mnt/session/outputs"
	stdout := "/var/log/warpbuild-agentd/runner.claude.stdout.log"
	stderr := "/var/log/warpbuild-agentd/runner.claude.stderr.log"
	if runtime.GOOS == "windows" {
		workdir = `C:\workspace`
		outputsDir = ""
		stdout = `C:\ProgramData\warpbuild\logs\runner.claude.stdout.log`
		stderr = `C:\ProgramData\warpbuild\logs\runner.claude.stderr.log`
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
	// The worker runs as the non-root `runner`. /workspace and /mnt/session/outputs are Anthropic-
	// documented paths it must use (they can't be relocated), and its stdout/stderr live under a log
	// dir — all at locations `runner` can't create under the filesystem root. ensureWritableDir creates
	// each and hands ownership to the current user, escalating via the sandbox VM's passwordless sudo on
	// a unix permission error. All idempotent.
	seen := map[string]bool{}
	for _, dir := range []string{m.Workdir, m.OutputsDir, filepath.Dir(m.StdoutFile), filepath.Dir(m.StderrFile)} {
		if dir == "" || dir == "." || seen[dir] {
			continue
		}
		seen[dir] = true
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

// ensureWritableDir makes dir exist and writable by the current (non-root worker) user. It tries a
// direct MkdirAll first — which succeeds on Windows, when the user already has permission, or when the
// dir already exists — and on a unix permission error escalates via the sandbox VM's passwordless sudo,
// creating the dir as root and chowning it back to the current user. Idempotent.
func ensureWritableDir(dir string) error {
	if dir == "" {
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err == nil {
		return nil
	} else if runtime.GOOS == "windows" || !os.IsPermission(err) {
		return err
	}
	u, err := user.Current()
	if err != nil {
		return err
	}
	if out, err := exec.Command("sudo", "-n", "mkdir", "-p", dir).CombinedOutput(); err != nil {
		return fmt.Errorf("sudo mkdir -p %s: %w (%s)", dir, err, strings.TrimSpace(string(out)))
	}
	if out, err := exec.Command("sudo", "-n", "chown", u.Username, dir).CombinedOutput(); err != nil {
		return fmt.Errorf("sudo chown %s %s: %w (%s)", u.Username, dir, err, strings.TrimSpace(string(out)))
	}
	return nil
}
