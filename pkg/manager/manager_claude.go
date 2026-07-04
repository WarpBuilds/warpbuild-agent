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
	home, err := os.UserHomeDir()
	if err != nil || home == "" || home == "/" {
		home = "/tmp"
	}
	workdir := "/workspace"
	outputsDir := "/mnt/session/outputs"
	stdout := filepath.Join(home, ".warpbuild", "warpbuild-agentd", "runner.claude.stdout.log")
	stderr := filepath.Join(home, ".warpbuild", "warpbuild-agentd", "runner.claude.stderr.log")
	if runtime.GOOS == "windows" {
		workdir = `C:\workspace`
		outputsDir = ""
		stdout = `C:\ProgramData\warpbuild\logs\runner.claude.stdout.log`
		stderr = `C:\ProgramData\warpbuild\logs\runner.claude.stderr.log`
	} else if runtime.GOOS == "darwin" {
		// macOS runners boot with a sealed, read-only system volume (SIP), so /workspace and
		// /mnt/session/outputs can't be created even via sudo. Anchor everything under the runner's
		// home on the writable data volume instead.
		workdir = filepath.Join(home, ".warpbuild", "workspace")
		outputsDir = filepath.Join(home, ".warpbuild", "session-outputs")
		stdout = filepath.Join(home, ".warpbuild", "agent", "log", "runner.claude.stdout.log")
		stderr = filepath.Join(home, ".warpbuild", "agent", "log", "runner.claude.stderr.log")
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

func (m *claudeManager) createFiles() error {
	// /workspace and /mnt/session/outputs (the worker's workdir + session deliverables) are Anthropic-
	// documented paths pre-created by the VM cloud-init — owned by the agentd user and world-writable —
	// because agentd may run non-root (e.g. x64 images) and can't create dirs under the root-owned
	// filesystem root. ensureWritableDir just ensures each exists and is 0777; stdout/stderr live under a
	// log dir agentd can create itself. All idempotent.
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

// ensureWritableDir makes dir exist and world-writable (0777) so the worker's (possibly non-root) tool
// user can write to it. The session dirs (/workspace, /mnt/session/outputs) are pre-created by the VM
// cloud-init owned by the agentd user, so this finds them and re-applies 0777; the log dir, which agentd
// can create itself, is created here. Std-lib only, no privilege escalation. Idempotent.
func ensureWritableDir(dir string) error {
	if dir == "" {
		return nil
	}
	// MkdirAll is a no-op on the cloud-init-pre-created session dirs (and creates the log dir where agentd
	// has permission); its mode is umask-masked, so chmod explicitly to guarantee 0777 for the non-root
	// worker. 0777 is safe on a single-use, single-tenant sandbox VM.
	if err := os.MkdirAll(dir, 0o777); err != nil {
		return err
	}
	if runtime.GOOS == "windows" {
		return nil
	}
	return os.Chmod(dir, 0o777)
}
