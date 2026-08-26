package hooks

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/manager"
)

const TELEMETRY_DRAIN_HOOK = manager.TELEMETRY_DRAIN_HOOK

const telemetryDrainTimeout = 10 * time.Second

type TelemetryDrainHook struct{}

var _ manager.IPostEndHook = &TelemetryDrainHook{}

func init() {
	manager.RegisterHook[manager.IPostEndHook](&TelemetryDrainHook{})
}

func (*TelemetryDrainHook) HookID() string {
	return TELEMETRY_DRAIN_HOOK
}

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
