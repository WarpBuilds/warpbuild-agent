package manager

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

type ghcriManager struct {
	*GithubCRIOptions
}

type GithubCRIOptions struct {
	StdoutFile string      `json:"stdout_file"`
	StderrFile string      `json:"stderr_file"`
	RunnerDir  string      `json:"runner_dir"`
	CMDOptions *CMDOptions `json:"cmd_options"`
}

type CMDOptions struct {
	CMD  string               `json:"cmd"`
	Args []string             `json:"args"`
	Dir  string               `json:"dir"`
	Envs EnvironmentVariables `json:"envs"`
}

var _ IManager = &ghcriManager{}

func newGithubCRIManager(opts *ManagerOptions) IManager {
	return &ghcriManager{
		GithubCRIOptions: opts.GithubCRI,
	}
}

func (m *ghcriManager) StartRunner(ctx context.Context, opts *StartRunnerOptions) (*StartRunnerOutput, error) {
	err := m.init(ctx)
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, m.CMDOptions.CMD, m.CMDOptions.Args...)
	cmd.Env = append(cmd.Env, "WARPBUILD_GH_JIT_TOKEN="+opts.JitToken)
	for _, env := range m.CMDOptions.Envs {
		log.Logger().Infof("setting env %s=%s", env.Key, env.Value)
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", env.Key, env.Value))
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
				Provider:  ProviderGithubCRI,
				GithubCRI: m.GithubCRIOptions,
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
			fmt.Fprintln(stdoutFile, out)
			fmt.Fprintln(os.Stdout, out)
		case err := <-stderrChan:
			fmt.Fprintln(stderrFile, err)
			fmt.Fprintln(os.Stderr, err)
		case <-doneChan:

			wg.Wait()

			// Exit the loop when command completes
			// Run all the post-end hooks
			for _, hook := range GetHooks[IPostEndHook]() {
				err := hook.PostEndHook(ctx, &PostEndHookOptions{
					StartRunnerOptions: opts,
					ManagerOptions: &ManagerOptions{
						Provider:  ProviderGithubCRI,
						GithubCRI: m.GithubCRIOptions,
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

func (m *ghcriManager) init(ctx context.Context) error {

	err := m.createFiles(ctx)
	if err != nil {
		log.Logger().Errorf("error creating files: %v", err)
		return err
	}

	return nil

}

func (m *ghcriManager) createFiles(ctx context.Context) error {

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
