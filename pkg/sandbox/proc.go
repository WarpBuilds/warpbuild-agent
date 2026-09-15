//go:build darwin

package sandbox

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/creack/pty"

	rpc "github.com/warpbuilds/warpbuild-agent/pkg/sandboxspec/process"
)

const (
	outputBufferSize = 64
	ptyReadChunk     = 16 * 1024
	pipeReadChunk    = 32 * 1024
)

// broadcaster fans one process's events out to every attached stream. Events are
// dropped when nobody is attached: output is live, never replayed, so a client
// that reconnects mid-run sees only what follows.
type broadcaster struct {
	mu     sync.Mutex
	subs   map[int]chan *rpc.ProcessEvent
	nextID int
	closed bool
}

func newBroadcaster() *broadcaster {
	return &broadcaster{subs: make(map[int]chan *rpc.ProcessEvent)}
}

func (b *broadcaster) subscribe() (int, <-chan *rpc.ProcessEvent, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return 0, nil, false
	}
	id := b.nextID
	b.nextID++
	ch := make(chan *rpc.ProcessEvent, outputBufferSize)
	b.subs[id] = ch

	return id, ch, true
}

func (b *broadcaster) unsubscribe(id int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if ch, ok := b.subs[id]; ok {
		delete(b.subs, id)
		close(ch)
	}
}

func (b *broadcaster) send(ev *rpc.ProcessEvent) {
	b.sendFunc(func() *rpc.ProcessEvent { return ev })
}

// sendFunc builds the event only when somebody is attached, so a process
// nobody is watching does not pay for copying its own output.
func (b *broadcaster) sendFunc(build func() *rpc.ProcessEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.subs) == 0 {
		return
	}

	ev := build()
	for _, ch := range b.subs {
		select {
		case ch <- ev:
		default:
		}
	}
}

func (b *broadcaster) close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	b.closed = true
	for id, ch := range b.subs {
		delete(b.subs, id)
		close(ch)
	}
}

type procHandler struct {
	pid    int
	tag    string
	config *rpc.ProcessConfig

	cmd *exec.Cmd
	tty *os.File

	stdinMu sync.Mutex
	stdin   io.WriteCloser

	events *broadcaster

	endEvent atomic.Pointer[rpc.ProcessEvent_EndEvent]

	cancelProc context.CancelFunc
	outWG      sync.WaitGroup
}

type spawnOptions struct {
	Config *rpc.ProcessConfig
	PTY    *rpc.PTY
	Tag    string
	Stdin  bool
	User   *user.User
	// Timeout bounds the process, not the request: a dropped connection must
	// leave the process reattachable, so the context is rooted at Background.
	Timeout time.Duration
}

// expandPath resolves ~ against the given user's home and makes relative paths
// relative to it. The ~user form is not supported, matching envd.
func expandPath(p string, u *user.User) string {
	if p == "" {
		return u.HomeDir
	}
	if p == "~" {
		return u.HomeDir
	}
	if strings.HasPrefix(p, "~/") {
		return filepath.Join(u.HomeDir, p[2:])
	}
	if filepath.IsAbs(p) {
		return filepath.Clean(p)
	}

	return filepath.Join(u.HomeDir, p)
}

// buildEnv constructs the child environment from scratch rather than inheriting
// ours: only PATH carries over, then the identity of the resolved user, then the
// request's own vars last so they win.
func buildEnv(cfg *rpc.ProcessConfig, u *user.User) []string {
	env := []string{
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + u.HomeDir,
		"USER=" + u.Username,
		"LOGNAME=" + u.Username,
	}
	for k, v := range cfg.GetEnvs() {
		env = append(env, k+"="+v)
	}

	return env
}

func startProcess(opts spawnOptions) (*procHandler, error) {
	cfg := opts.Config
	if cfg.GetCmd() == "" {
		return nil, fmt.Errorf("cmd is required")
	}

	cwd := expandPath(cfg.GetCwd(), opts.User)
	if fi, err := os.Stat(cwd); err != nil || !fi.IsDir() {
		return nil, fmt.Errorf("working directory %q does not exist", cwd)
	}

	procCtx, cancelProc := context.Background(), context.CancelFunc(func() {})
	if opts.Timeout > 0 {
		procCtx, cancelProc = context.WithTimeout(context.Background(), opts.Timeout)
	}

	cmd := exec.CommandContext(procCtx, cfg.GetCmd(), cfg.GetArgs()...)
	cmd.Dir = cwd
	cmd.Env = buildEnv(cfg, opts.User)

	h := &procHandler{
		tag:        opts.Tag,
		config:     cfg,
		cmd:        cmd,
		events:     newBroadcaster(),
		cancelProc: cancelProc,
	}

	fail := func(err error) (*procHandler, error) {
		cancelProc()
		cmd.Cancel = nil
		closePipes(cmd)

		return nil, err
	}

	if opts.PTY != nil {
		size := &pty.Winsize{Cols: uint16(opts.PTY.GetSize().GetCols()), Rows: uint16(opts.PTY.GetSize().GetRows())}
		tty, err := pty.StartWithSize(cmd, size)
		if err != nil {
			return fail(fmt.Errorf("start pty: %w", err))
		}
		h.tty = tty
		h.outWG.Add(1)
		go h.pump(tty, ptyReadChunk, func(b []byte) *rpc.ProcessEvent {
			return dataEvent(&rpc.ProcessEvent_DataEvent{Output: &rpc.ProcessEvent_DataEvent_Pty{Pty: b}})
		})
	} else {
		// New process group so the whole tree can be reaped once the leader exits;
		// this is what stops a backgrounded child from outliving its command.
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return fail(err)
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return fail(err)
		}
		if opts.Stdin {
			in, err := cmd.StdinPipe()
			if err != nil {
				return fail(err)
			}
			h.stdin = in
		}
		if err := cmd.Start(); err != nil {
			return fail(fmt.Errorf("start process: %w", err))
		}
		h.outWG.Add(2)
		go h.pump(stdout, pipeReadChunk, func(b []byte) *rpc.ProcessEvent {
			return dataEvent(&rpc.ProcessEvent_DataEvent{Output: &rpc.ProcessEvent_DataEvent_Stdout{Stdout: b}})
		})
		go h.pump(stderr, pipeReadChunk, func(b []byte) *rpc.ProcessEvent {
			return dataEvent(&rpc.ProcessEvent_DataEvent{Output: &rpc.ProcessEvent_DataEvent_Stderr{Stderr: b}})
		})
	}

	h.pid = cmd.Process.Pid

	return h, nil
}

func dataEvent(d *rpc.ProcessEvent_DataEvent) *rpc.ProcessEvent {
	return &rpc.ProcessEvent{Event: &rpc.ProcessEvent_Data{Data: d}}
}

func (h *procHandler) pump(r io.Reader, chunk int, wrap func([]byte) *rpc.ProcessEvent) {
	defer h.outWG.Done()
	buf := make([]byte, chunk)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			h.events.sendFunc(func() *rpc.ProcessEvent {
				cp := make([]byte, n)
				copy(cp, buf[:n])

				return wrap(cp)
			})
		}
		if err != nil {
			return
		}
	}
}

// wait blocks until the process exits, then publishes the terminal event and
// reaps anything left in its process group.
func (h *procHandler) wait(onExit func(*rpc.ProcessEvent_EndEvent)) {
	h.outWG.Wait()
	err := h.cmd.Wait()

	if h.tty != nil {
		h.tty.Close()
	}

	var errMsg *string
	if err != nil {
		msg := err.Error()
		errMsg = &msg
	}

	end := &rpc.ProcessEvent_EndEvent{
		Error:    errMsg,
		ExitCode: int32(h.cmd.ProcessState.ExitCode()),
		Exited:   h.cmd.ProcessState.Exited(),
		Status:   h.cmd.ProcessState.String(),
	}
	h.endEvent.Store(end)

	h.events.send(&rpc.ProcessEvent{Event: &rpc.ProcessEvent_End{End: end}})

	// Retain the exit before closing, so a Connect that arrives just after the
	// close falls back to the retention cache instead of blocking.
	if onExit != nil {
		onExit(end)
	}
	h.events.close()

	h.reapGroup()
	h.cancelProc()
}

// reapGroup kills anything still running in the leader's process group. Without
// this a nohup'd child keeps the sandbox busy after its command returned.
func (h *procHandler) reapGroup() {
	if h.pid <= 0 {
		return
	}
	_ = syscall.Kill(-h.pid, syscall.SIGKILL)
}

func (h *procHandler) signal(sig syscall.Signal) error {
	if h.cmd.Process == nil {
		return fmt.Errorf("process is not running")
	}

	return h.cmd.Process.Signal(sig)
}

func (h *procHandler) resize(size *rpc.PTY_Size) error {
	if h.tty == nil {
		return fmt.Errorf("process has no pty")
	}

	return pty.Setsize(h.tty, &pty.Winsize{Cols: uint16(size.GetCols()), Rows: uint16(size.GetRows())})
}

func (h *procHandler) writeStdin(b []byte) error {
	h.stdinMu.Lock()
	defer h.stdinMu.Unlock()
	if h.tty != nil {
		return fmt.Errorf("process has a pty; write to the pty instead")
	}
	if h.stdin == nil {
		return fmt.Errorf("process has no stdin")
	}
	_, err := h.stdin.Write(b)

	return err
}

func (h *procHandler) writeTTY(b []byte) error {
	if h.tty == nil {
		return fmt.Errorf("process has no pty")
	}
	_, err := h.tty.Write(b)

	return err
}

func (h *procHandler) closeStdin() error {
	h.stdinMu.Lock()
	defer h.stdinMu.Unlock()
	if h.tty != nil {
		return fmt.Errorf("cannot close stdin on a pty process; send Ctrl+D (0x04) instead")
	}
	if h.stdin == nil {
		return nil
	}
	err := h.stdin.Close()
	h.stdin = nil

	return err
}

func (h *procHandler) info() *rpc.ProcessInfo {
	info := &rpc.ProcessInfo{Config: h.config, Pid: uint32(h.pid)}
	if h.tag != "" {
		tag := h.tag
		info.Tag = &tag
	}

	return info
}

// closePipes releases descriptors opened by StdoutPipe/StderrPipe/StdinPipe on a
// spawn that never reached Start. exec.Cmd only closes them from Start or Wait,
// so without this they leak for the life of the daemon.
func closePipes(cmd *exec.Cmd) {
	for _, c := range []io.Closer{asCloser(cmd.Stdout), asCloser(cmd.Stderr), asCloser(cmd.Stdin)} {
		if c != nil {
			_ = c.Close()
		}
	}
}

func asCloser(v any) io.Closer {
	c, _ := v.(io.Closer)

	return c
}
