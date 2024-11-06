//go:build !windows
// +build !windows

package manager

import "context"

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
	return nil, nil
}
