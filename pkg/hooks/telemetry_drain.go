package hooks

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/manager"
)

const TELEMETRY_DRAIN_HOOK string = "TELEMETRY_DRAIN_HOOK"

const telemetryDrainTimeout = 10 * time.Second

type TelemetryDrainHook struct{}

var (
	_ manager.IPostEndHook  = &TelemetryDrainHook{}
	_ manager.IHookPriority = &TelemetryDrainHook{}
)

func init() {
	manager.RegisterHook[manager.IPostEndHook](&TelemetryDrainHook{})
}

func (*TelemetryDrainHook) HookID() string {
	return TELEMETRY_DRAIN_HOOK
}

// HookPriority runs this ahead of the cleanup callback, which tells the
// backend it may reap the VM. Draining after that is a race.
func (*TelemetryDrainHook) HookPriority() int {
	return -100
}

// PostEndHook asks the telemetryd sibling process to flush.
//
// Batches are on a 30s timer, so the tail of a job — usually the
// interesting part — would otherwise still be buffered when this VM is
// destroyed. Best-effort: telemetry must never fail a run, and a runner
// with telemetry disabled has nothing listening at all.
func (*TelemetryDrainHook) PostEndHook(ctx context.Context, opts *manager.PostEndHookOptions) error {
	if opts == nil || opts.StartRunnerOptions == nil || opts.StartRunnerOptions.AgentOptions == nil {
		return nil
	}
	port := opts.StartRunnerOptions.AgentOptions.TelemetryPort
	if port == 0 {
		return nil
	}

	reqCtx, cancel := context.WithTimeout(ctx, telemetryDrainTimeout)
	defer cancel()

	url := fmt.Sprintf("http://127.0.0.1:%d/internal/drain", port)
	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, url, nil)
	if err != nil {
		log.Logger().Warnf("failed to build telemetry drain request: %v", err)
		return nil
	}

	resp, err := (&http.Client{Timeout: telemetryDrainTimeout}).Do(req)
	if err != nil {
		log.Logger().Infof("telemetry drain skipped (telemetryd unreachable): %v", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Logger().Warnf("telemetry drain returned %s", resp.Status)
		return nil
	}

	log.Logger().Infof("Telemetry flushed before VM teardown")
	return nil
}
