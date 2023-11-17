package manager

import "context"

type PreStartHookOptions struct {
	StartRunnerOptions *StartRunnerOptions `json:"start_runner_options"`
	ManagerOptions     *ManagerOptions     `json:"manager_options"`
}

type IPreStartHook interface {
	// HookID is a unique name for the hook. This will be used in logs.
	HookID() string
	PreStartHook(ctx context.Context, opts *PreStartHookOptions) error
}

type PostEndHookOptions struct {
	StartRunnerOptions *StartRunnerOptions `json:"start_runner_options"`
	ManagerOptions     *ManagerOptions     `json:"manager_options"`
}

type IPostEndHook interface {
	// HookID is a unique name for the hook. This will be used in logs.
	HookID() string
	PostEndHook(ctx context.Context, opts *PostEndHookOptions) error
}

var hooks []any

func RegisterHook[T any](hook T) {
	hooks = append(hooks, hook)
}

func GetHooks[T any]() []T {
	var result []T
	for _, hook := range hooks {
		if h, ok := hook.(T); ok {
			result = append(result, h)
		}
	}
	return result
}
