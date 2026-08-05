package hooks

import (
	"context"

	"github.com/warpbuilds/warpbuild-agent/pkg/manager"
)

const CLAUDE_OUTPUTS_UPLOAD_HOOK string = "CLAUDE_OUTPUTS_UPLOAD_HOOK"

type ClaudeOutputsUploadHook struct{}

var _ manager.IPostEndHook = &ClaudeOutputsUploadHook{}

func init() {
	manager.RegisterHook[manager.IPostEndHook](&ClaudeOutputsUploadHook{})
}

func (*ClaudeOutputsUploadHook) HookID() string {
	return CLAUDE_OUTPUTS_UPLOAD_HOOK
}

func (*ClaudeOutputsUploadHook) PostEndHook(ctx context.Context, opts *manager.PostEndHookOptions) error {
	if opts.ManagerOptions == nil || opts.ManagerOptions.Provider != manager.ProviderClaudeAgent {
		return nil
	}
	c := opts.ManagerOptions.Claude
	if c == nil {
		return nil
	}
	a := opts.StartRunnerOptions.AgentOptions
	manager.UploadSessionOutputs(ctx, c.OutputsDir, a.HostURL, a.PollingSecret, a.ID)
	return nil
}
