package hooks

import (
	"context"
	_ "embed"
	"os"
	"os/exec"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/manager"
)

//go:embed upload_outputs.js
var uploadOutputsScript []byte

const CLAUDE_OUTPUTS_UPLOAD_HOOK string = "CLAUDE_OUTPUTS_UPLOAD_HOOK"

const outputsUploadTimeout = 5 * time.Minute

type ClaudeOutputsUploadHook struct{}

var _ manager.IPostEndHook = &ClaudeOutputsUploadHook{}

func init() {
	manager.RegisterHook[manager.IPostEndHook](&ClaudeOutputsUploadHook{})
}

func (*ClaudeOutputsUploadHook) HookID() string {
	return CLAUDE_OUTPUTS_UPLOAD_HOOK
}

// PostEndHook uploads the claude session deliverables to backend-cache via the node warp-cache client,
// keyed by the anthropic_session_id. No-op for non-claude runs. Best-effort: logs and returns on error
// so a failed upload never fails the run. Runs before the cleanup hook reaps this single-use VM.
func (*ClaudeOutputsUploadHook) PostEndHook(ctx context.Context, opts *manager.PostEndHookOptions) error {
	if opts.ManagerOptions == nil || opts.ManagerOptions.Provider != manager.ProviderClaudeAgent {
		return nil
	}
	c := opts.ManagerOptions.Claude
	if c == nil || c.OutputsDir == "" || c.SessionID == "" || c.CacheBackendHost == "" {
		return nil
	}

	f, err := os.CreateTemp("", "wb-upload-outputs-*.js")
	if err != nil {
		log.Logger().Errorf("[claude_outputs] failed to stage upload script: %v", err)
		return nil
	}
	scriptPath := f.Name()
	defer os.Remove(scriptPath)
	if _, err := f.Write(uploadOutputsScript); err != nil {
		_ = f.Close()
		log.Logger().Errorf("[claude_outputs] failed to write upload script: %v", err)
		return nil
	}
	_ = f.Close()

	runCtx, cancel := context.WithTimeout(ctx, outputsUploadTimeout)
	defer cancel()

	cmd := exec.CommandContext(runCtx, "node", scriptPath)
	env := append(os.Environ(),
		"WARPBUILD_CACHE_URL="+c.CacheBackendHost,
		"WARPBUILD_RUNNER_VERIFICATION_TOKEN="+c.RunnerVerificationToken,
		"SANDBOX_OUTPUTS_DIR="+c.OutputsDir,
		"SANDBOX_SESSION_ID="+c.SessionID,
		"RUNNER_TEMP="+os.TempDir(),
	)
	if os.Getenv("NODE_PATH") == "" {
		env = append(env, "NODE_PATH="+os.Getenv("HOME")+"/.warpbuild/cache-client/node_modules")
	}
	cmd.Env = env
	out, cerr := cmd.CombinedOutput()
	if cerr != nil {
		log.Logger().Errorf("[claude_outputs] upload failed: %v: %s", cerr, string(out))
		return nil
	}
	log.Logger().Infof("[claude_outputs] %s", string(out))
	return nil
}
