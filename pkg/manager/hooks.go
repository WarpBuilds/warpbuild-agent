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

const (
	TELEMETRY_DRAIN_HOOK       = "TELEMETRY_DRAIN_HOOK"
	CLAUDE_OUTPUTS_UPLOAD_HOOK = "CLAUDE_OUTPUTS_UPLOAD_HOOK"
	CLEANUP_CALLBACK_HOOK      = "CLEANUP_CALLBACK_HOOK"
)

var hookRunOrder = []string{
	TELEMETRY_DRAIN_HOOK,
	CLAUDE_OUTPUTS_UPLOAD_HOOK,
	CLEANUP_CALLBACK_HOOK,
}

type hookIDer interface {
	HookID() string
}

var hooks []any

func RegisterHook[T any](hook T) {
	hooks = append(hooks, hook)
}

func GetHooks[T any]() []T {
	var matching []T
	for _, hook := range hooks {
		if h, ok := hook.(T); ok {
			matching = append(matching, h)
		}
	}

	result := make([]T, 0, len(matching))
	placed := make([]bool, len(matching))

	for _, name := range hookRunOrder {
		for i, hook := range matching {
			if placed[i] {
				continue
			}
			if h, ok := any(hook).(hookIDer); ok && h.HookID() == name {
				result = append(result, hook)
				placed[i] = true
			}
		}
	}

	for i, hook := range matching {
		if !placed[i] {
			result = append(result, hook)
		}
	}

	return result
}
