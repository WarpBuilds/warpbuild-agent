package hooks

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/manager"
)

func TestMain(m *testing.M) {
	if _, err := log.Init(&log.InitOptions{LogLevel: "error"}); err != nil {
		panic(err)
	}
	m.Run()
}

func drainOpts(port int) *manager.PostEndHookOptions {
	return &manager.PostEndHookOptions{
		StartRunnerOptions: &manager.StartRunnerOptions{
			AgentOptions: &manager.AgentOptions{TelemetryPort: port},
		},
	}
}

// stubTelemetryd stands in for the sibling process on loopback.
func stubTelemetryd(t *testing.T, handler http.HandlerFunc) int {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	u, err := url.Parse(srv.URL)
	require.NoError(t, err, "parse stub url")
	port, err := strconv.Atoi(u.Port())
	require.NoError(t, err, "parse stub port")
	return port
}

func TestTelemetryDrainHook_CallsDrain(t *testing.T) {
	var gotPath, gotMethod string
	port := stubTelemetryd(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		w.WriteHeader(http.StatusOK)
	})

	require.NoError(t, (&TelemetryDrainHook{}).PostEndHook(context.Background(), drainOpts(port)))
	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, "/internal/drain", gotPath)
}

// Telemetry must never fail a customer's run, so every failure mode here
// is a no-op rather than an error.
func TestTelemetryDrainHook_NeverFailsTheRun(t *testing.T) {
	hook := &TelemetryDrainHook{}

	cases := map[string]*manager.PostEndHookOptions{
		"nil opts":           nil,
		"nil start options":  {},
		"nil agent options":  {StartRunnerOptions: &manager.StartRunnerOptions{}},
		"telemetry disabled": drainOpts(0),
		"nothing listening":  drainOpts(1),
	}
	for name, opts := range cases {
		t.Run(name, func(t *testing.T) {
			assert.NoError(t, hook.PostEndHook(context.Background(), opts))
		})
	}

	t.Run("telemetryd errors", func(t *testing.T) {
		port := stubTelemetryd(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		assert.NoError(t, hook.PostEndHook(context.Background(), drainOpts(port)))
	})
}

// The cleanup callback tells the backend it may reap the VM, so draining
// after it is a race we would lose.
func TestTelemetryDrainRunsBeforeCleanup(t *testing.T) {
	postEnd := manager.GetHooks[manager.IPostEndHook]()

	drainAt, cleanupAt := -1, -1
	for i, h := range postEnd {
		switch h.HookID() {
		case TELEMETRY_DRAIN_HOOK:
			drainAt = i
		case CLEANUP_CALLBACK_HOOK:
			cleanupAt = i
		}
	}

	require.GreaterOrEqual(t, drainAt, 0, "telemetry drain hook is not registered")
	require.GreaterOrEqual(t, cleanupAt, 0, "cleanup callback hook is not registered")
	assert.Less(t, drainAt, cleanupAt, "drain must run before cleanup")
}
