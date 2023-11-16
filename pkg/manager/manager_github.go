package manager

import (
	"context"
	"fmt"
	"os/exec"
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
}

var _ IManager = (*ghManager)(nil)

func newGithubManager(opts *ManagerOptions) IManager {
	return &ghManager{
		GithubOptions: opts.Github,
	}
}

func (m *ghManager) StartRunner(ctx context.Context, opts *StartRunnerOptions) error {
	cmd := exec.CommandContext(ctx, "runner.sh", "--jitToken", opts.JitToken)
	cmd.Dir = m.RunnerDir

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Logger().Errorf("error creating stdout pipe: %v", err)
		return err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Logger().Errorf("error creating stderr pipe: %v", err)
		return err
	}

	if err := cmd.Start(); err != nil {
		log.Logger().Errorf("error starting command: %v", err)
		return err
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
		return err
	}
	defer stdoutFile.Close()

	stderrFile, err := openFile(m.StderrFile)
	if err != nil {
		log.Logger().Errorf("error opening stderr file: %v", err)
		return err
	}
	defer stderrFile.Close()

	// Goroutine to wait for the command to finish
	go func() {
		cmd.Wait()
		doneChan <- true
	}()

	// Ticker to capture output every second
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case out := <-stdoutChan:
			fmt.Fprintln(stdoutFile, out)
		case err := <-stderrChan:
			fmt.Fprintln(stderrFile, err)
		case <-ticker.C:
			// Handle output every second
		case <-doneChan:
			// Exit the loop when command completes
			return nil
		}
	}

}
