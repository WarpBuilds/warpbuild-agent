//go:build !windows

package manager

import (
	"context"
	"errors"
)

// The github_windows_cri provider is Windows-only. This file exists so the
// package compiles on linux/darwin (for tests, cross-builds, IDE tooling).
// The real implementation lives in manager_windows_github_cri_windows.go.
// In production this provider is only ever selected on Windows runners.

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
	return nil, errors.New("github_windows_cri provider is not implemented on this OS (Windows only)")
}
