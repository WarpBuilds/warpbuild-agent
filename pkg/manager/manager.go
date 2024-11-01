package manager

import (
	"context"
)

type StartRunnerOutput struct {
	RunCompletedSuccessfully bool `json:"run_completed_successfully"`
}

type IManager interface {
	StartRunner(ctx context.Context, opts *StartRunnerOptions) (*StartRunnerOutput, error)
}

type StartRunnerOptions struct {
	JitToken     string        `json:"jit_token"`
	AgentOptions *AgentOptions `json:"agent_options"`
}

type ManagerOptions struct {
	Provider         Provider                 `json:"provider"`
	Github           *GithubOptions           `json:"github"`
	GithubCRI        *GithubCRIOptions        `json:"github_cri"`
	GithubWindowsCRI *GithubWindowsCRIOptions `json:"github_windows_cri"`
}

func NewManager(opts *ManagerOptions) IManager {
	switch opts.Provider {
	case ProviderGithub:
		return newGithubManager(opts)
	case ProviderGithubCRI:
		return newGithubCRIManager(opts)
	case ProviderGithubWindowsCRI:
		return newGithubWindowsCRIManager(opts)
	default:
		panic("unknown provider")
	}
}
