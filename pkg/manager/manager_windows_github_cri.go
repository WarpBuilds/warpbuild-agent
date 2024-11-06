//go:build windows
// +build windows

package manager

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"unsafe"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"golang.org/x/sys/windows"
)

type GithubWindowsCRIOptions struct {
	PassAllEnvs bool        `json:"pass_all_envs"`
	Username    string      `json:"username"`
	Password    string      `json:"password"`
	Domain      string      `json:"domain"`
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
	err := m.init(ctx)
	if err != nil {
		return nil, err
	}

	// Logon and impersonate the user
	token, err := logonUser(m.Username, m.Password, m.Domain)
	if err != nil {
		log.Logger().Errorf("LogonUser failed: %v", err)
		return nil, err
	}
	defer windows.CloseHandle(token)

	err = impersonateLoggedOnUser(token)
	if err != nil {
		log.Logger().Errorf("ImpersonateLoggedOnUser failed: %v", err)
		return nil, err
	}

	log.Logger().Infof("Logged in as %s", m.Username)

	// log all the envs
	log.Logger().Infof("Environment variables available to the warpbuild agent:")
	for _, env := range os.Environ() {
		log.Logger().Infof("env: %s", env)
	}

	cmd := exec.CommandContext(ctx, m.CMDOptions.CMD, m.CMDOptions.Args...)
	cmd.Env = append(cmd.Env, "WARPBUILD_GH_JIT_TOKEN="+opts.JitToken)
	for _, env := range m.CMDOptions.Envs {
		log.Logger().Infof("setting env %s=%s", env.Key, env.Value)
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", env.Key, env.Value))
	}

	if m.PassAllEnvs {
		log.Logger().Infof("Adding all available envs to command...")
		for _, env := range os.Environ() {
			log.Logger().Infof("env: %s", env)
			cmd.Env = append(cmd.Env, env)
		}
	}

	cmd.Dir = m.CMDOptions.Dir

	log.Logger().Infof("starting runner with command: %s", cmd.String())
	log.Logger().Infof("JIT Token: %s", opts.JitToken)

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

	for _, hook := range GetHooks[IPreStartHook]() {
		err := hook.PreStartHook(ctx, &PreStartHookOptions{
			StartRunnerOptions: opts,
			ManagerOptions: &ManagerOptions{
				Provider:         ProviderGithubWindowsCRI,
				GithubWindowsCRI: m.GithubWindowsCRIOptions,
			},
		})
		if err != nil {
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

	// Open files to write stdout and stderr
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

	for {
		select {
		case out := <-stdoutChan:
			log.Logger().Infof(out)
			fmt.Fprintln(stdoutFile, out)
		case err := <-stderrChan:
			log.Logger().Errorf(err)
			fmt.Fprintln(stderrFile, err)
		case <-doneChan:

			wg.Wait()

			// Exit the loop when command completes
			// Run all the post-end hooks
			for _, hook := range GetHooks[IPostEndHook]() {
				err := hook.PostEndHook(ctx, &PostEndHookOptions{
					StartRunnerOptions: opts,
					ManagerOptions: &ManagerOptions{
						Provider:         ProviderGithubWindowsCRI,
						GithubWindowsCRI: m.GithubWindowsCRIOptions,
					},
				})
				if err != nil {
					log.Logger().Errorf("error running post-end hook %s: %v", hook.HookID(), err)
					return nil, err
				}
			}

			return &StartRunnerOutput{
				RunCompletedSuccessfully: true,
			}, nil
		}
	}

}

func (m *ghWindowsCriManager) init(ctx context.Context) error {

	err := m.createFiles(ctx)
	if err != nil {
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

	// Iterate over the full paths
	for _, fullPath := range fullPaths {
		// Get the base directory from the full path
		baseDir := filepath.Dir(fullPath)

		// Ensure the base directory exists
		if _, err := os.Stat(baseDir); os.IsNotExist(err) {
			err := os.MkdirAll(baseDir, 0755)
			if err != nil {
				log.Logger().Errorf("Failed to create base directory %s: %v", baseDir, err)
				return err
			}
		}

		// Create the file if it doesn't exist
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

const (
	LOGON32_LOGON_INTERACTIVE = 2
	LOGON32_PROVIDER_DEFAULT  = 0
)

// Load the necessary Windows API DLLs and functions
var (
	modAdvapi32                 = windows.NewLazySystemDLL("advapi32.dll")
	procLogonUserW              = modAdvapi32.NewProc("LogonUserW")
	procImpersonateLoggedOnUser = modAdvapi32.NewProc("ImpersonateLoggedOnUser")
)

func logonUser(username, password, domain string) (windows.Handle, error) {
	var token windows.Handle
	usernamePtr, _ := windows.UTF16PtrFromString(username)
	domainPtr, _ := windows.UTF16PtrFromString(domain)
	passwordPtr, _ := windows.UTF16PtrFromString(password)

	// Call LogonUserW from advapi32.dll
	ret, _, err := procLogonUserW.Call(
		uintptr(unsafe.Pointer(usernamePtr)),
		uintptr(unsafe.Pointer(domainPtr)),
		uintptr(unsafe.Pointer(passwordPtr)),
		uintptr(LOGON32_LOGON_INTERACTIVE),
		uintptr(LOGON32_PROVIDER_DEFAULT),
		uintptr(unsafe.Pointer(&token)),
	)
	if ret == 0 {
		return 0, err
	}
	return token, nil
}

func impersonateLoggedOnUser(token windows.Handle) error {
	// Call ImpersonateLoggedOnUser from advapi32.dll
	ret, _, err := procImpersonateLoggedOnUser.Call(uintptr(token))
	if ret == 0 {
		return err
	}
	return nil
}
