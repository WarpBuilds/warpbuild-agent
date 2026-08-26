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

func TestGetHooksAppliesRunOrder(t *testing.T) {
	withHooks(t, CLEANUP_CALLBACK_HOOK, CLAUDE_OUTPUTS_UPLOAD_HOOK, TELEMETRY_DRAIN_HOOK)

	assert.Equal(t, []string{
		TELEMETRY_DRAIN_HOOK,
		CLAUDE_OUTPUTS_UPLOAD_HOOK,
		CLEANUP_CALLBACK_HOOK,
	}, hookIDs(t))
}

func TestGetHooksPutsUnlistedHooksLast(t *testing.T) {
	withHooks(t, "ZEBRA_HOOK", CLEANUP_CALLBACK_HOOK, "ALPHA_HOOK", TELEMETRY_DRAIN_HOOK)

	assert.Equal(t, []string{
		TELEMETRY_DRAIN_HOOK,
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
