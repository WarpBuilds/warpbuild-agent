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

// Hook IDs. Declared here rather than alongside each implementation in
// pkg/hooks so hookRunOrder below can name them: pkg/hooks imports this
// package, not the other way round.
const (
	CLAUDE_OUTPUTS_UPLOAD_HOOK = "CLAUDE_OUTPUTS_UPLOAD_HOOK"
	CLEANUP_CALLBACK_HOOK      = "CLEANUP_CALLBACK_HOOK"
)

// hookRunOrder is the order hooks run in. Anything not listed runs after
// these, in registration order.
//
// CLEANUP_CALLBACK_HOOK tells the backend it may reap the VM, so anything
// that needs the VM alive belongs above it.
var hookRunOrder = []string{
	CLAUDE_OUTPUTS_UPLOAD_HOOK,
	CLEANUP_CALLBACK_HOOK,
}

// hookIDer is what both hook interfaces embed; used to order them.
type hookIDer interface {
	HookID() string
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
	// Stable, so unlisted hooks keep registration order behind the listed ones.
	sort.SliceStable(result, func(i, j int) bool {
		return hookRunRank(result[i]) < hookRunRank(result[j])
	})
	return result
}

// hookRunRank is the hook's index in hookRunOrder, or len(hookRunOrder)
// for anything unlisted.
func hookRunRank(hook any) int {
	h, ok := hook.(hookIDer)
	if !ok {
		return len(hookRunOrder)
	}
	for i, name := range hookRunOrder {
		if name == h.HookID() {
			return i
		}
	}
	return len(hookRunOrder)
}
