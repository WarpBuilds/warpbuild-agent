package hooks

import (
	"context"

	"github.com/warpbuilds/warpbuild-agent/pkg/manager"
)

const CLAUDE_OUTPUTS_UPLOAD_HOOK string = "CLAUDE_OUTPUTS_UPLOAD_HOOK"

type ClaudeOutputsUploadHook struct{}

var _ manager.IPostEndHook = &ClaudeOutputsUploadHook{}

// This file sorts before cleanup.go, so its init registers first and GetHooks returns it first: the
// deliverables must upload BEFORE the cleanup hook, which reaps this single-use VM.
func init() {
	manager.RegisterHook[manager.IPostEndHook](&ClaudeOutputsUploadHook{})
}

// HookID implements manager.IPostEndHook.
func (*ClaudeOutputsUploadHook) HookID() string {
	return CLAUDE_OUTPUTS_UPLOAD_HOOK
}

// PostEndHook uploads the claude session deliverables (/mnt/session/outputs) after the worker exits.
// No-op for non-claude runs. Best-effort (UploadSessionOutputs never fails the run).
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
