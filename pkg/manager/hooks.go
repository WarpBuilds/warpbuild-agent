package manager

import (
	"context"
	"sort"
)

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

// IHookPriority lets a hook state where it belongs in the running order
// rather than inheriting whatever order package init happened to produce.
// Lower runs first; a hook that does not implement this is 0.
//
// This matters for post-end hooks: the cleanup callback tells the backend
// it may reap the VM, so anything that needs the VM alive has to run
// before it.
type IHookPriority interface {
	HookPriority() int
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
	// Stable, so hooks that share a priority keep registration order.
	sort.SliceStable(result, func(i, j int) bool {
		return hookPriority(result[i]) < hookPriority(result[j])
	})
	return result
}

func hookPriority(hook any) int {
	if p, ok := hook.(IHookPriority); ok {
		return p.HookPriority()
	}
	return 0
}
