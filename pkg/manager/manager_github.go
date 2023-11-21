package manager

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

type ghManager struct {
	*GithubOptions
}

type GithubOptions struct {
	StdoutFile string `json:"stdout_file"`
	StderrFile string `json:"stderr_file"`
	RunnerDir  string `json:"runner_dir"`
	Script     string `json:"script"`
}

var _ IManager = (*ghManager)(nil)

func newGithubManager(opts *ManagerOptions) IManager {
	return &ghManager{
		GithubOptions: opts.Github,
	}
}

func (m *ghManager) StartRunner(ctx context.Context, opts *StartRunnerOptions) (*StartRunnerOutput, error) {
	err := m.init(ctx)
	if err != nil {
		return nil, err
	}

	os.Setenv("RUNNER_ALLOW_RUNASROOT", "1")

	cmd := exec.CommandContext(ctx, m.Script, "--jitconfig", opts.JitToken)
	cmd.Dir = m.RunnerDir

	log.Logger().Infof("starting runner with command: %s", cmd.String())

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
				Provider: ProviderGithub,
				Github:   m.GithubOptions,
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

	// Ticker to capture output every second
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case out := <-stdoutChan:
			fmt.Fprintln(stdoutFile, out)
			fmt.Fprintln(os.Stdout, out)
		case err := <-stderrChan:
			fmt.Fprintln(stderrFile, err)
			fmt.Fprintln(os.Stderr, err)
		case <-ticker.C:
			// Handle output every second
		case <-doneChan:

			wg.Wait()

			// Exit the loop when command completes
			// Run all the post-end hooks
			for _, hook := range GetHooks[IPostEndHook]() {
				err := hook.PostEndHook(ctx, &PostEndHookOptions{
					StartRunnerOptions: opts,
					ManagerOptions: &ManagerOptions{
						Provider: ProviderGithub,
						Github:   m.GithubOptions,
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

func (m *ghManager) init(ctx context.Context) error {

	err := m.createFiles(ctx)
	if err != nil {
		log.Logger().Errorf("error creating files: %v", err)
		return err
	}

	return nil

}

func (m *ghManager) createFiles(ctx context.Context) error {

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
