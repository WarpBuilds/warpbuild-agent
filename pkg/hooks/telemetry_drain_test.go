package hooks

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

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
	if err != nil {
		t.Fatalf("parse stub url: %v", err)
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		t.Fatalf("parse stub port: %v", err)
	}
	return port
}

func TestTelemetryDrainHook_CallsDrain(t *testing.T) {
	var gotPath, gotMethod string
	port := stubTelemetryd(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		w.WriteHeader(http.StatusOK)
	})

	if err := (&TelemetryDrainHook{}).PostEndHook(context.Background(), drainOpts(port)); err != nil {
		t.Fatalf("PostEndHook: %v", err)
	}
	if gotPath != "/internal/drain" || gotMethod != http.MethodPost {
		t.Errorf("called %s %s, want POST /internal/drain", gotMethod, gotPath)
	}
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
			if err := hook.PostEndHook(context.Background(), opts); err != nil {
				t.Errorf("expected nil error, got %v", err)
			}
		})
	}

	t.Run("telemetryd errors", func(t *testing.T) {
		port := stubTelemetryd(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		if err := hook.PostEndHook(context.Background(), drainOpts(port)); err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
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

	if drainAt < 0 {
		t.Fatal("telemetry drain hook is not registered")
	}
	if cleanupAt < 0 {
		t.Fatal("cleanup callback hook is not registered")
	}
	if drainAt > cleanupAt {
		t.Errorf("drain runs at %d, after cleanup at %d", drainAt, cleanupAt)
	}
}
