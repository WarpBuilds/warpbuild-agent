package manager

import "context"

type IManager interface {
	StartRunner(ctx context.Context, opts *StartRunnerOptions) error
}

type StartRunnerOptions struct {
	JitToken string `json:"jit_token"`
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
