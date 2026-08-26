package manager

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakePostEndHook struct{ id string }

func (h fakePostEndHook) HookID() string { return h.id }
func (h fakePostEndHook) PostEndHook(context.Context, *PostEndHookOptions) error {
	return nil
}

// withHooks swaps the package-level registry for the duration of a test.
func withHooks(t *testing.T, ids ...string) {
	t.Helper()

	original := hooks
	t.Cleanup(func() { hooks = original })

	hooks = nil
	for _, id := range ids {
		RegisterHook[IPostEndHook](fakePostEndHook{id: id})
	}
}

func hookIDs(t *testing.T) []string {
	t.Helper()

	var got []string
	for _, h := range GetHooks[IPostEndHook]() {
		got = append(got, h.HookID())
	}
	return got
}

// The cleanup callback tells the backend it may reap the VM, so anything
// that needs the VM alive has to run before it — regardless of the order
// package init happened to register them in.
func TestGetHooksAppliesRunOrder(t *testing.T) {
	withHooks(t, CLEANUP_CALLBACK_HOOK, CLAUDE_OUTPUTS_UPLOAD_HOOK)

	assert.Equal(t, []string{CLAUDE_OUTPUTS_UPLOAD_HOOK, CLEANUP_CALLBACK_HOOK}, hookIDs(t))
}

// Hooks nobody ordered run after the ordered ones, keeping registration order.
func TestGetHooksPutsUnlistedHooksLast(t *testing.T) {
	withHooks(t, "ZEBRA_HOOK", CLEANUP_CALLBACK_HOOK, "ALPHA_HOOK", CLAUDE_OUTPUTS_UPLOAD_HOOK)

	assert.Equal(t, []string{
		CLAUDE_OUTPUTS_UPLOAD_HOOK,
		CLEANUP_CALLBACK_HOOK,
		"ZEBRA_HOOK",
		"ALPHA_HOOK",
	}, hookIDs(t))
}

func TestGetHooksFiltersByType(t *testing.T) {
	withHooks(t, CLEANUP_CALLBACK_HOOK)

	require.Len(t, GetHooks[IPostEndHook](), 1)
	assert.Empty(t, GetHooks[IPreStartHook]())
}

// A hook named in the run order ranks at its position; anything else sorts
// after the whole list.
func TestHookRunRank(t *testing.T) {
	for i, id := range hookRunOrder {
		assert.Equalf(t, i, hookRunRank(fakePostEndHook{id: id}), "rank for %s", id)
	}
	assert.Equal(t, len(hookRunOrder), hookRunRank(fakePostEndHook{id: "UNLISTED_HOOK"}))
	assert.Equal(t, len(hookRunOrder), hookRunRank("not a hook"))
}
