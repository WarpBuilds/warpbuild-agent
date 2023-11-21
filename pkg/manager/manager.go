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
	Provider Provider       `json:"provider"`
	Github   *GithubOptions `json:"github"`
}

func NewManager(opts *ManagerOptions) IManager {
	switch opts.Provider {
	case ProviderGithub:
		return newGithubManager(opts)
	default:
		panic("unknown provider")
	}
}
