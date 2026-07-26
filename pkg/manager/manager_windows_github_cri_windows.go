//go:build windows

package manager

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"golang.org/x/sys/windows"
)

type GithubWindowsCRIOptions struct {
	PassAllEnvs bool        `json:"pass_all_envs"`
	StdoutFile  string      `json:"stdout_file"`
	StderrFile  string      `json:"stderr_file"`
	RunnerDir   string      `json:"runner_dir"`
	CMDOptions  *CMDOptions `json:"cmd_options"`
}

type ghWindowsCriManager struct {
	*GithubWindowsCRIOptions
}

var _ IManager = &ghWindowsCriManager{}

func newGithubWindowsCRIManager(opts *ManagerOptions) IManager {
	return &ghWindowsCriManager{
		GithubWindowsCRIOptions: opts.GithubWindowsCRI,
	}
}

func (m *ghWindowsCriManager) StartRunner(ctx context.Context, opts *StartRunnerOptions) (*StartRunnerOutput, error) {
	if err := m.init(ctx); err != nil {
		return nil, err
	}

	// log all the envs
	log.Logger().Infof("Environment variables available to the warpbuild agent:")
	for _, env := range os.Environ() {
		log.Logger().Infof("env: %s", env)
	}

	env := []string{"WARPBUILD_GH_JIT_TOKEN=" + opts.JitToken}
	for _, kv := range m.CMDOptions.Envs {
		log.Logger().Infof("setting env %s=%s", kv.Key, kv.Value)
		env = append(env, fmt.Sprintf("%s=%s", kv.Key, kv.Value))
	}
	if m.PassAllEnvs {
		log.Logger().Infof("Adding all available envs to command...")
		for _, e := range os.Environ() {
			log.Logger().Infof("env: %s", e)
			env = append(env, e)
		}
	}

	for _, hook := range GetHooks[IPreStartHook]() {
		if err := hook.PreStartHook(ctx, &PreStartHookOptions{
			StartRunnerOptions: opts,
			ManagerOptions: &ManagerOptions{
				Provider:         ProviderGithubWindowsCRI,
				GithubWindowsCRI: m.GithubWindowsCRIOptions,
			},
		}); err != nil {
			log.Logger().Errorf("error running pre-start hook %s: %v", hook.HookID(), err)
			return nil, err
		}
	}

	log.Logger().Infof("starting runner with command: %s %v", m.CMDOptions.CMD, m.CMDOptions.Args)
	log.Logger().Infof("JIT Token: %s", opts.JitToken)

	proc, err := startInteractiveProcess(ctx, m.CMDOptions, env)
	if err != nil {
		log.Logger().Errorf("error starting command: %v", err)
		return nil, err
	}

	stdoutChan := make(chan string)
	stderrChan := make(chan string)
	doneChan := make(chan bool)

	go captureOutput(proc.stdout, stdoutChan)
	go captureOutput(proc.stderr, stderrChan)

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
		_ = proc.wait()
		doneChan <- true
	}()

	for {
		select {
		case out := <-stdoutChan:
			log.Logger().Infof(out)
			fmt.Fprintln(stdoutFile, out)
		case errLine := <-stderrChan:
			log.Logger().Errorf(errLine)
			fmt.Fprintln(stderrFile, errLine)
		case <-doneChan:
			wg.Wait()

			for _, hook := range GetHooks[IPostEndHook]() {
				if err := hook.PostEndHook(ctx, &PostEndHookOptions{
					StartRunnerOptions: opts,
					ManagerOptions: &ManagerOptions{
						Provider:         ProviderGithubWindowsCRI,
						GithubWindowsCRI: m.GithubWindowsCRIOptions,
					},
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

func (m *ghWindowsCriManager) init(ctx context.Context) error {
	if err := m.createFiles(ctx); err != nil {
		log.Logger().Errorf("error creating files: %v", err)
		return err
	}
	return nil
}

func (m *ghWindowsCriManager) createFiles(ctx context.Context) error {
	fullPaths := []string{
		m.StderrFile,
		m.StdoutFile,
	}
	for _, fullPath := range fullPaths {
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

// startInteractiveProcess launches the runner inside the active interactive
// desktop session (Session 1, WinSta0\Default) instead of inheriting the
// Session 0 service window station. This is required for UI/GUI workloads
// (Selenium, WinAppDriver, MAUI/WinUI, anything that touches winsta0\default).
// See WARP-849.
//
// Flow: WTSEnumerateSessions -> pick WTSActive -> WTSQueryUserToken to obtain
// the logged-on user's primary token -> DuplicateTokenEx for a fresh primary
// token -> CreateEnvironmentBlock for that user -> CreateProcessAsUser with
// STARTUPINFOW.lpDesktop = "winsta0\\default" and stdout/stderr wired to
// anonymous pipes for log forwarding.
func startInteractiveProcess(ctx context.Context, cmdOpts *CMDOptions, env []string) (*interactiveProcess, error) {
	sessionID, err := pickActiveSessionID()
	if err != nil {
		return nil, fmt.Errorf("locating active interactive session: %w", err)
	}
	log.Logger().Infof("launching runner into interactive session id=%d", sessionID)

	var userToken windows.Token
	if err := windows.WTSQueryUserToken(sessionID, &userToken); err != nil {
		return nil, fmt.Errorf("WTSQueryUserToken(session=%d): %w", sessionID, err)
	}
	defer userToken.Close()

	var primaryToken windows.Token
	if err := windows.DuplicateTokenEx(
		userToken,
		windows.MAXIMUM_ALLOWED,
		nil,
		windows.SecurityImpersonation,
		windows.TokenPrimary,
		&primaryToken,
	); err != nil {
		return nil, fmt.Errorf("DuplicateTokenEx: %w", err)
	}

	envBlockPtr, err := buildEnvironmentBlock(primaryToken, env)
	if err != nil {
		primaryToken.Close()
		return nil, fmt.Errorf("building environment block: %w", err)
	}

	stdoutR, stdoutW, err := makeInheritablePipe(false)
	if err != nil {
		primaryToken.Close()
		return nil, err
	}
	stderrR, stderrW, err := makeInheritablePipe(false)
	if err != nil {
		stdoutR.Close()
		stdoutW.Close()
		primaryToken.Close()
		return nil, err
	}
	stdinR, stdinW, err := makeInheritablePipe(true)
	if err != nil {
		stdoutR.Close()
		stdoutW.Close()
		stderrR.Close()
		stderrW.Close()
		primaryToken.Close()
		return nil, err
	}
	stdinW.Close()

	// lpApplicationName = NULL so Windows PATH-resolves the first token of
	// lpCommandLine against the *user's* environment (e.g. "powershell.exe"
	// resolves without us pre-resolving it).
	cmdLine, err := windows.UTF16PtrFromString(buildCommandLine(cmdOpts.CMD, cmdOpts.Args))
	if err != nil {
		closeAll(primaryToken, stdoutR, stdoutW, stderrR, stderrW, stdinR)
		return nil, err
	}
	desktop, err := windows.UTF16PtrFromString(`winsta0\default`)
	if err != nil {
		closeAll(primaryToken, stdoutR, stdoutW, stderrR, stderrW, stdinR)
		return nil, err
	}
	var dirPtr *uint16
	if cmdOpts.Dir != "" {
		dirPtr, err = windows.UTF16PtrFromString(cmdOpts.Dir)
		if err != nil {
			closeAll(primaryToken, stdoutR, stdoutW, stderrR, stderrW, stdinR)
			return nil, err
		}
	}

	si := windows.StartupInfo{
		Cb:         uint32(unsafe.Sizeof(windows.StartupInfo{})),
		Desktop:    desktop,
		Flags:      windows.STARTF_USESTDHANDLES,
		StdInput:   windows.Handle(stdinR.Fd()),
		StdOutput:  windows.Handle(stdoutW.Fd()),
		StdErr:     windows.Handle(stderrW.Fd()),
		ShowWindow: windows.SW_HIDE,
	}
	var pi windows.ProcessInformation
	flags := uint32(windows.CREATE_UNICODE_ENVIRONMENT | windows.CREATE_NO_WINDOW)

	if err := windows.CreateProcessAsUser(
		primaryToken,
		nil,
		cmdLine,
		nil,
		nil,
		true,
		flags,
		envBlockPtr,
		dirPtr,
		&si,
		&pi,
	); err != nil {
		closeAll(primaryToken, stdoutR, stdoutW, stderrR, stderrW, stdinR)
		_ = windows.DestroyEnvironmentBlock(envBlockPtr)
		return nil, fmt.Errorf("CreateProcessAsUser: %w", err)
	}

	stdoutW.Close()
	stderrW.Close()
	stdinR.Close()
	_ = windows.DestroyEnvironmentBlock(envBlockPtr)

	proc := &interactiveProcess{
		process:      pi.Process,
		thread:       pi.Thread,
		primaryToken: primaryToken,
		stdout:       stdoutR,
		stderr:       stderrR,
		done:         make(chan struct{}),
	}

	go func() {
		select {
		case <-ctx.Done():
			_ = windows.TerminateProcess(proc.process, 1)
		case <-proc.done:
		}
	}()

	return proc, nil
}

// pickActiveSessionID returns the session id of the active interactive console
// session. Prefers any WTSActive session over the bare console session id, so
// that RDP-attached sessions also work.
func pickActiveSessionID() (uint32, error) {
	var sessions *windows.WTS_SESSION_INFO
	var count uint32
	if err := windows.WTSEnumerateSessions(0, 0, 1, &sessions, &count); err != nil {
		return 0, fmt.Errorf("WTSEnumerateSessions: %w", err)
	}
	defer windows.WTSFreeMemory(uintptr(unsafe.Pointer(sessions)))

	size := unsafe.Sizeof(*sessions)
	base := unsafe.Pointer(sessions)
	for i := uint32(0); i < count; i++ {
		info := (*windows.WTS_SESSION_INFO)(unsafe.Add(base, uintptr(i)*size))
		if info.State == windows.WTSActive {
			return info.SessionID, nil
		}
	}
	if id := windows.WTSGetActiveConsoleSessionId(); id != 0xFFFFFFFF {
		return id, nil
	}
	return 0, errors.New("no active interactive session found (autologon may have failed)")
}

// buildEnvironmentBlock returns a pointer to a UTF-16 environment block,
// merging the user's profile env with the caller-supplied overlay (overlay
// wins). Caller must DestroyEnvironmentBlock the result.
func buildEnvironmentBlock(token windows.Token, overlay []string) (*uint16, error) {
	var userBlock *uint16
	if err := windows.CreateEnvironmentBlock(&userBlock, token, false); err != nil {
		return nil, fmt.Errorf("CreateEnvironmentBlock: %w", err)
	}
	if len(overlay) == 0 {
		return userBlock, nil
	}

	envMap := make(map[string]string)
	for _, kv := range readEnvironmentBlock(userBlock) {
		if idx := strings.Index(kv, "="); idx > 0 {
			envMap[strings.ToUpper(kv[:idx])] = kv[idx+1:]
		}
	}
	for _, kv := range overlay {
		if idx := strings.Index(kv, "="); idx > 0 {
			envMap[strings.ToUpper(kv[:idx])] = kv[idx+1:]
		}
	}
	_ = windows.DestroyEnvironmentBlock(userBlock)

	var buf []uint16
	for k, v := range envMap {
		entry, err := windows.UTF16FromString(k + "=" + v)
		if err != nil {
			return nil, err
		}
		buf = append(buf, entry...)
	}
	buf = append(buf, 0)
	return &buf[0], nil
}

func readEnvironmentBlock(block *uint16) []string {
	if block == nil {
		return nil
	}
	var out []string
	p := unsafe.Pointer(block)
	for {
		n := 0
		for *(*uint16)(unsafe.Add(p, uintptr(n)*2)) != 0 {
			n++
		}
		if n == 0 {
			break
		}
		s := unsafe.Slice((*uint16)(p), n)
		out = append(out, windows.UTF16ToString(s))
		p = unsafe.Add(p, uintptr(n+1)*2)
	}
	return out
}

// makeInheritablePipe returns a (parent, child) file pair. childReads=true
// marks the read end inheritable (for stdin); childReads=false marks the
// write end inheritable (for stdout/stderr).
func makeInheritablePipe(childReads bool) (parentEnd, childEnd *os.File, err error) {
	var r, w windows.Handle
	sa := &windows.SecurityAttributes{
		Length:        uint32(unsafe.Sizeof(windows.SecurityAttributes{})),
		InheritHandle: 1,
	}
	if err := windows.CreatePipe(&r, &w, sa, 0); err != nil {
		return nil, nil, fmt.Errorf("CreatePipe: %w", err)
	}
	if childReads {
		if err := windows.SetHandleInformation(w, windows.HANDLE_FLAG_INHERIT, 0); err != nil {
			windows.CloseHandle(r)
			windows.CloseHandle(w)
			return nil, nil, fmt.Errorf("SetHandleInformation: %w", err)
		}
		return os.NewFile(uintptr(w), "stdin-w"), os.NewFile(uintptr(r), "stdin-r"), nil
	}
	if err := windows.SetHandleInformation(r, windows.HANDLE_FLAG_INHERIT, 0); err != nil {
		windows.CloseHandle(r)
		windows.CloseHandle(w)
		return nil, nil, fmt.Errorf("SetHandleInformation: %w", err)
	}
	return os.NewFile(uintptr(r), "child-out-r"), os.NewFile(uintptr(w), "child-out-w"), nil
}

func closeAll(token windows.Token, files ...*os.File) {
	for _, f := range files {
		if f != nil {
			f.Close()
		}
	}
	if token != 0 {
		token.Close()
	}
}

// buildCommandLine joins exe and args into a single Windows command line,
// using the standard CommandLineToArgvW quoting rules.
func buildCommandLine(cmd string, args []string) string {
	var b strings.Builder
	b.WriteString(syscall.EscapeArg(cmd))
	for _, a := range args {
		b.WriteByte(' ')
		b.WriteString(syscall.EscapeArg(a))
	}
	return b.String()
}

type interactiveProcess struct {
	process      windows.Handle
	thread       windows.Handle
	primaryToken windows.Token
	stdout       *os.File
	stderr       *os.File
	done         chan struct{}
}

func (p *interactiveProcess) wait() error {
	_, waitErr := windows.WaitForSingleObject(p.process, windows.INFINITE)

	// Signal the ctx-watcher goroutine before closing the handle, so a racing
	// TerminateProcess can't fire against a closed handle.
	close(p.done)

	defer func() {
		windows.CloseHandle(p.thread)
		windows.CloseHandle(p.process)
		p.primaryToken.Close()
	}()

	if waitErr != nil {
		return fmt.Errorf("WaitForSingleObject: %w", waitErr)
	}
	var code uint32
	if err := windows.GetExitCodeProcess(p.process, &code); err != nil {
		return fmt.Errorf("GetExitCodeProcess: %w", err)
	}
	if code != 0 {
		return fmt.Errorf("runner exited with code %d", code)
	}
	return nil
}
