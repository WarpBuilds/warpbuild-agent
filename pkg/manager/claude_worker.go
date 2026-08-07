package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/tools/agenttoolset"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

// runClaudeWorker runs the managed-agent session worker in-process, replacing the
// `ant beta:worker run` exec. It drives the anthropic-sdk-go SessionToolRunner
// (kept open for the whole session, so it rides every turn) alongside our own
// work-item heartbeat. Unlike the SDK's own EnvironmentWorker, our heartbeat does
// NOT cancel the session runner when the work item reports state=stopping/stopped
// — the SSE session stream is independent of the work-item lease, so the worker
// keeps serving turns. It exits only on:
//   - session.status_terminated / session.deleted (ErrSessionTerminated),
//   - MaxIdle of continuous end_turn idle (ErrIdleTimeout),
//   - ctx cancellation (SIGTERM / teardown), or
//   - a heartbeat staleness ceiling (lease presumed lost).
//
// Creds come from exactly the ANTHROPIC_* vars agentd already sets on the worker
// (nothing extra); no org API key ever touches the VM.
func runClaudeWorker(ctx context.Context, c *ClaudeOptions) error {
	creds := claudeEnvMap(c.Envs)
	envKey := creds["ANTHROPIC_ENVIRONMENT_KEY"]
	envID := creds["ANTHROPIC_ENVIRONMENT_ID"]
	sessionID := creds["ANTHROPIC_SESSION_ID"]
	workID := creds["ANTHROPIC_WORK_ID"]

	logger := claudeSlog(c.StderrFile)
	logger.Info("claude in-process worker start",
		slog.String("env_id", envID), slog.String("session_id", sessionID),
		slog.String("work_id", workID), slog.Int("env_key_len", len(envKey)),
		slog.String("workdir", c.Workdir), slog.String("max_idle", c.MaxIdle))
	if envKey == "" || envID == "" || sessionID == "" || workID == "" {
		return fmt.Errorf("claude worker missing creds (need ANTHROPIC_ENVIRONMENT_KEY/ENVIRONMENT_ID/SESSION_ID/WORK_ID)")
	}

	client := anthropic.NewClient() // no API key on the VM; env key is attached per-request
	// The environment key is a bearer; it must be paired with an X-Api-Key delete
	// or the parent client's default api-key rides along and the events stream
	// rejects the dual auth.
	bearer := []option.RequestOption{
		option.WithHeaderDel("X-Api-Key"),
		option.WithAuthToken(envKey),
	}

	env := &agenttoolset.AgentToolContext{Workdir: c.Workdir}
	if err := env.SetupSkills(ctx, client, sessionID, bearer...); err != nil {
		logger.Warn("skill setup failed (continuing without skills)", slog.Any("error", err))
	}
	defer func() { _ = env.Cleanup() }()
	tools := wrapNonEmptyResult(agenttoolset.BetaAgentToolset20260401(env))
	defer agenttoolset.CloseAll(tools)

	// The runner's context is ours: the heartbeat cancels it only on a genuine
	// staleness/lease loss, NOT on a work-item stop.
	runnerCtx, runnerCancel := context.WithCancel(ctx)
	defer runnerCancel()

	hbDone := make(chan struct{})
	go func() {
		defer close(hbDone)
		claudeHeartbeat(runnerCtx, runnerCancel, client, workID, envID, bearer, logger)
	}()

	mi := claudeMaxIdle(c.MaxIdle)
	runner := client.Beta.Sessions.Events.NewToolRunner(runnerCtx, sessionID, anthropic.SessionToolRunnerOptions{
		Tools:          tools,
		MaxIdle:        mi,
		Logger:         logger,
		RequestOptions: bearer,
	})
	start := time.Now()
	for runner.Next() {
		call := runner.Current()
		logger.Info("tool dispatched",
			slog.String("name", call.Name), slog.Bool("custom", call.Custom),
			slog.Bool("is_error", call.IsError), slog.Bool("posted", call.Posted))
	}
	err := runner.Err()
	_ = runner.Close()
	runnerCancel()
	<-hbDone

	logger.Info("claude worker runner exited",
		slog.String("elapsed", time.Since(start).Round(time.Second).String()),
		slog.Any("error", err),
		slog.Bool("session_terminated", errors.Is(err, anthropic.ErrSessionTerminated)),
		slog.Bool("idle_timeout", errors.Is(err, anthropic.ErrIdleTimeout)))

	// Session-terminated and idle-timeout are benign, expected terminations.
	if err != nil && !errors.Is(err, anthropic.ErrSessionTerminated) && !errors.Is(err, anthropic.ErrIdleTimeout) {
		return err
	}
	return nil
}

// claudeHeartbeat keeps the work-item lease alive and logs each beat's state.
// Crucially — unlike the SDK's runHeartbeat — on a work-item state=stopping/stopped
// (or a lease-not-extended) it does NOT cancel the session runner: it just stops
// beating. The session stream is what governs the session's lifetime. It cancels
// the runner only on a genuine staleness ceiling (lease presumed lost server-side).
func claudeHeartbeat(ctx context.Context, cancel context.CancelFunc, client anthropic.Client, workID, envID string, reqOpts []option.RequestOption, logger *slog.Logger) {
	const floor = time.Second
	interval := 30 * time.Second
	ttl := 30 * time.Second
	last := "NO_HEARTBEAT"
	lastSuccess := time.Now()
	reclaims := 0

	beat := func() bool { // returns: keep beating?
		opts := append(reqOpts[:len(reqOpts):len(reqOpts)], option.WithRequestTimeout(interval))
		resp, err := client.Beta.Environments.Work.Heartbeat(ctx, workID, anthropic.BetaEnvironmentWorkHeartbeatParams{
			EnvironmentID:         envID,
			ExpectedLastHeartbeat: param.NewOpt(last),
		}, opts...)
		if err != nil {
			if ctx.Err() != nil {
				return false
			}
			// A 412 means the work item already carries a live lease — typically a
			// prior worker instance from before an agentd restart (crash or the
			// cloud-init binary swap). Since ONE VM owns this session, reclaim the
			// lease by adopting the server's current last_heartbeat instead of dying.
			// Bounded so two genuinely-live workers can't ping-pong forever.
			var apierr *anthropic.Error
			if errors.As(err, &apierr) && apierr.StatusCode == 412 && reclaims < 3 {
				if lh := reclaimLastHeartbeat(apierr); lh != "" && lh != last {
					reclaims++
					logger.Warn("heartbeat 412; reclaiming lease from current_state",
						slog.String("last_heartbeat", lh), slog.Int("reclaim", reclaims))
					last = lh
					lastSuccess = time.Now()
					return true
				}
			}
			if stale := time.Since(lastSuccess); stale > ttl {
				logger.Error("heartbeat staleness ceiling; cancelling session runner",
					slog.String("since_success", stale.String()), slog.String("ttl", ttl.String()), slog.Any("error", err))
				cancel()
				return false
			}
			logger.Warn("transient heartbeat error", slog.Any("error", err))
			return true
		}
		last = resp.LastHeartbeat
		lastSuccess = time.Now()
		reclaims = 0
		if resp.TTLSeconds > 0 {
			ttl = max(time.Duration(resp.TTLSeconds)*time.Second, floor)
			interval = max(floor, min(30*time.Second, ttl/2))
		}
		logger.Info("heartbeat", slog.String("state", string(resp.State)),
			slog.Int64("ttl_s", resp.TTLSeconds), slog.Bool("lease_extended", resp.LeaseExtended))
		switch resp.State {
		case anthropic.BetaSelfHostedWorkHeartbeatResponseStateStopping,
			anthropic.BetaSelfHostedWorkHeartbeatResponseStateStopped:
			logger.Warn("work item stopped; keeping session runner alive (policy), stopping heartbeat only")
			return false
		}
		if !resp.LeaseExtended {
			logger.Warn("lease not extended; keeping session runner alive (policy), stopping heartbeat only")
			return false
		}
		return true
	}

	if !beat() {
		return
	}
	for {
		t := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
		}
		if !beat() {
			return
		}
	}
}

// nonEmptyResultTool wraps a BetaTool so an empty tool result never reaches the
// events API as an empty text block (which the API rejects with
// "content.0.text: value is required", stalling the turn). Empty output → "(no output)".
type nonEmptyResultTool struct{ inner anthropic.BetaTool }

func (t nonEmptyResultTool) Name() string                                { return t.inner.Name() }
func (t nonEmptyResultTool) Description() string                         { return t.inner.Description() }
func (t nonEmptyResultTool) InputSchema() anthropic.BetaToolInputSchemaParam { return t.inner.InputSchema() }

func (t nonEmptyResultTool) Execute(ctx context.Context, input json.RawMessage) ([]anthropic.BetaToolResultBlockParamContentUnion, error) {
	res, err := t.inner.Execute(ctx, input)
	if err != nil {
		return res, err
	}
	if len(res) == 0 {
		return []anthropic.BetaToolResultBlockParamContentUnion{{OfText: &anthropic.BetaTextBlockParam{Text: "(no output)"}}}, nil
	}
	for i := range res {
		if res[i].OfText != nil && strings.TrimSpace(res[i].OfText.Text) == "" {
			res[i].OfText.Text = "(no output)"
		}
	}
	return res, nil
}

// Close forwards to the wrapped tool so agenttoolset.CloseAll still tears down the
// persistent bash session.
func (t nonEmptyResultTool) Close() error {
	if c, ok := t.inner.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

func wrapNonEmptyResult(tools []anthropic.BetaTool) []anthropic.BetaTool {
	out := make([]anthropic.BetaTool, len(tools))
	for i, t := range tools {
		out[i] = nonEmptyResultTool{inner: t}
	}
	return out
}

// reclaimLastHeartbeat extracts error.details.current_state.last_heartbeat from a
// 412 heartbeat error body so the worker can adopt (reclaim) the existing lease.
func reclaimLastHeartbeat(apierr *anthropic.Error) string {
	var body struct {
		Error struct {
			Details struct {
				CurrentState struct {
					LastHeartbeat string `json:"last_heartbeat"`
				} `json:"current_state"`
			} `json:"details"`
		} `json:"error"`
	}
	if json.Unmarshal([]byte(apierr.RawJSON()), &body) != nil {
		return ""
	}
	return body.Error.Details.CurrentState.LastHeartbeat
}

func claudeEnvMap(envs EnvironmentVariables) map[string]string {
	m := make(map[string]string, len(envs))
	for _, e := range envs {
		m[e.Key] = e.Value
	}
	return m
}

func claudeMaxIdle(s string) *time.Duration {
	d, err := time.ParseDuration(strings.TrimSpace(s))
	if err != nil || d < 0 {
		def, _ := time.ParseDuration(anthropicWorkerMaxIdle)
		return &def
	}
	return &d
}

// claudeSlog writes the worker's structured logs to the same stderr log file the
// ant exec used (so telemetry that reads runner.claude.stderr.log keeps working);
// falls back to os.Stderr.
func claudeSlog(path string) *slog.Logger {
	var w io.Writer = os.Stderr
	if path != "" {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err == nil {
			if f, ferr := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); ferr == nil {
				w = f
			} else {
				log.Logger().Warnf("claude worker: could not open %s, logging to stderr: %v", path, ferr)
			}
		}
	}
	return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
