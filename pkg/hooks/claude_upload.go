package hooks

import (
	"context"
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"
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

	nodeBin := resolveNodeBinary()
	cmd := exec.CommandContext(runCtx, nodeBin, scriptPath)
	env := append(os.Environ(),
		"WARPBUILD_CACHE_URL="+c.CacheBackendHost,
		"WARPBUILD_RUNNER_VERIFICATION_TOKEN="+c.RunnerVerificationToken,
		"AGENT_RUNNERS_OUTPUTS_DIR="+c.OutputsDir,
		"AGENT_RUNNERS_SESSION_ID="+c.SessionID,
		"RUNNER_TEMP="+os.TempDir(),
	)
	env = append(env, "NODE_PATH="+resolveNodeModulePath())
	// Put node's own dir first on PATH so node — and anything the cache client shells out to (tar/zstd),
	// which on macOS is also under Homebrew — resolves even under a minimal launchd/systemd service PATH.
	// Last PATH= wins per exec.Cmd.Env semantics.
	if filepath.IsAbs(nodeBin) {
		env = append(env, "PATH="+filepath.Dir(nodeBin)+string(os.PathListSeparator)+os.Getenv("PATH"))
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

// resolveNodeModulePath returns a directory containing @warpbuilds/cache for the upload hook's NODE_PATH.
// The cloud-init installs the module to <runner-home>/.warpbuild/cache-client/node_modules but no longer
// exports NODE_PATH; the agentd's own $HOME may differ from the runner home (e.g. a root systemd service),
// so probe an explicit NODE_PATH first, then $HOME, then the known runner-home locations, and pick the one
// that actually has the module. Falls back to the first candidate so upload_outputs.js surfaces a clear
// not-found error rather than silently mis-resolving.
func resolveNodeModulePath() string {
	var candidates []string
	if np := os.Getenv("NODE_PATH"); np != "" {
		candidates = append(candidates, np)
	}
	if home := os.Getenv("HOME"); home != "" {
		candidates = append(candidates, filepath.Join(home, ".warpbuild", "cache-client", "node_modules"))
	}
	candidates = append(candidates,
		"/home/runner/.warpbuild/cache-client/node_modules",
		"/Users/runner/.warpbuild/cache-client/node_modules",
		"/root/.warpbuild/cache-client/node_modules",
	)
	for _, c := range candidates {
		if fi, err := os.Stat(filepath.Join(c, "@warpbuilds", "cache")); err == nil && fi.IsDir() {
			return c
		}
	}
	if len(candidates) > 0 {
		return candidates[0]
	}
	return ""
}

// resolveNodeBinary returns an absolute path to the node executable. A launchd (macOS) or systemd service
// is spawned with a minimal PATH (e.g. /usr/bin:/bin:/usr/sbin:/sbin) that omits Homebrew (/opt/homebrew/bin)
// and nvm, so a bare "node" lookup fails on mac even when node is installed. Fall back to well-known
// locations; last resort "node" surfaces a clear not-found error rather than silently mis-resolving.
func resolveNodeBinary() string {
	if p, err := exec.LookPath("node"); err == nil {
		return p
	}
	for _, p := range []string{"/opt/homebrew/bin/node", "/usr/local/bin/node", "/usr/bin/node"} {
		if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
			return p
		}
	}
	return "node"
}
