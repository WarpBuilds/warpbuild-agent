package hooks

import (
	"context"
	"net/http"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/manager"
	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

const CLEANUP_CALLBACK_HOOK string = "CLEANUP_CALLBACK_HOOK"

type CleanupCallbackHook struct{}

var _ manager.IPostEndHook = &CleanupCallbackHook{}

func init() {
	manager.RegisterHook[manager.IPostEndHook](&CleanupCallbackHook{})
}

func NewWarpBuildClient(opts *manager.AgentOptions) *warpbuild.APIClient {
	cfg := warpbuild.NewConfiguration()
	cfg.Servers[0].URL = opts.HostURL
	return warpbuild.NewAPIClient(cfg)
}

// HookID implements manager.IPostEndHook.
func (*CleanupCallbackHook) HookID() string {
	return CLEANUP_CALLBACK_HOOK
}

// PostEndHook implements manager.IPostEndHook.
func (*CleanupCallbackHook) PostEndHook(ctx context.Context, opts *manager.PostEndHookOptions) error {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			_, res, err := NewWarpBuildClient(opts.StartRunnerOptions.AgentOptions).V1RunnerInstanceApi.
				RunnerInstanceCleanupHook(ctx, opts.StartRunnerOptions.AgentOptions.ID).
				XPOLLINGSECRET(opts.StartRunnerOptions.AgentOptions.PollingSecret).
				Execute()
			if err != nil {
				log.Logger().Errorf("failed to call cleanup hook runner instance: %v", err)
				continue
			}
			if res.StatusCode == http.StatusOK {
				return nil
			}
		case <-ctx.Done():
			log.Logger().Infof("Context cancelled. Agent is exiting...")
			return nil
		}
	}
}
