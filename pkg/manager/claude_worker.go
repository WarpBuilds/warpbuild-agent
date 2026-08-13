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

// runClaudeWorker runs the managed-agent session worker in-process.
// It exits only on:
//   - session.status_terminated / session.deleted (ErrSessionTerminated),
//   - MaxIdle of continuous end_turn idle (ErrIdleTimeout),
//   - ctx cancellation (SIGTERM / teardown), or
//   - a heartbeat staleness ceiling (lease presumed lost).
func runClaudeWorker(ctx context.Context, c *ClaudeOptions) error {
	envKey := c.EnvKey
	envID := c.EnvID
	sessionID := c.SessionID
	workID := c.WorkID

	logger := claudeSlog(c.StderrFile)
	logger.Info("claude in-process worker start",
		slog.String("env_id", envID), slog.String("session_id", sessionID),
		slog.String("work_id", workID), slog.Int("env_key_len", len(envKey)),
		slog.String("workdir", c.Workdir), slog.String("max_idle", c.MaxIdle))
	if envKey == "" || envID == "" || sessionID == "" || workID == "" {
		return fmt.Errorf("claude worker missing creds (need ANTHROPIC_ENVIRONMENT_KEY/ENVIRONMENT_ID/SESSION_ID/WORK_ID)")
	}

	client := anthropic.NewClient() // no API key on the VM; env key is attached per-request
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

	hb := &claudeHeartbeater{
		client:      client,
		workID:      workID,
		envID:       envID,
		reqOpts:     bearer,
		logger:      logger,
		cancel:      runnerCancel,
		interval:    heartbeatInterval,
		ttl:         heartbeatTTL,
		last:        "NO_HEARTBEAT",
		lastSuccess: time.Now(),
	}
	hbDone := make(chan struct{})
	go func() {
		defer close(hbDone)
		hb.run(runnerCtx)
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

const (
	heartbeatFloor       = time.Second
	heartbeatInterval    = 30 * time.Second
	heartbeatTTL         = 30 * time.Second
	heartbeatMaxReclaims = 3
)

// claudeHeartbeater keeps the work-item lease alive and logs each beat's state.
// Crucially — unlike the SDK's runHeartbeat — on a work-item state=stopping/stopped
// (or a lease-not-extended) it does NOT cancel the session runner: it just stops
// beating. The session stream is what governs the session's lifetime. It cancels
// the runner only on a genuine staleness ceiling (lease presumed lost server-side).
type claudeHeartbeater struct {
	client  anthropic.Client
	workID  string
	envID   string
	reqOpts []option.RequestOption
	logger  *slog.Logger
	cancel  context.CancelFunc

	interval    time.Duration
	ttl         time.Duration
	last        string
	lastSuccess time.Time
	reclaims    int
}

// run beats once immediately, then every h.interval until a beat says stop or ctx ends.
func (h *claudeHeartbeater) run(ctx context.Context) {
	for h.beat(ctx) {
		t := time.NewTimer(h.interval)
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
		}
	}
}

// beat sends one heartbeat and reports whether to keep beating.
func (h *claudeHeartbeater) beat(ctx context.Context) bool {
	opts := append(h.reqOpts[:len(h.reqOpts):len(h.reqOpts)], option.WithRequestTimeout(h.interval))
	resp, err := h.client.Beta.Environments.Work.Heartbeat(ctx, h.workID, anthropic.BetaEnvironmentWorkHeartbeatParams{
		EnvironmentID:         h.envID,
		ExpectedLastHeartbeat: param.NewOpt(h.last),
	}, opts...)
	if err != nil {
		return h.onError(ctx, err)
	}

	h.last = resp.LastHeartbeat
	h.lastSuccess = time.Now()
	h.reclaims = 0
	if resp.TTLSeconds > 0 {
		h.ttl = max(time.Duration(resp.TTLSeconds)*time.Second, heartbeatFloor)
		h.interval = max(heartbeatFloor, min(heartbeatInterval, h.ttl/2))
	}
	h.logger.Info("heartbeat", slog.String("state", string(resp.State)),
		slog.Int64("ttl_s", resp.TTLSeconds), slog.Bool("lease_extended", resp.LeaseExtended))

	switch {
	case resp.State == anthropic.BetaSelfHostedWorkHeartbeatResponseStateStopping,
		resp.State == anthropic.BetaSelfHostedWorkHeartbeatResponseStateStopped:
		h.logger.Warn("work item stopped; keeping session runner alive (policy), stopping heartbeat only")
		return false
	case !resp.LeaseExtended:
		h.logger.Warn("lease not extended; keeping session runner alive (policy), stopping heartbeat only")
		return false
	default:
		return true
	}
}

// onError decides whether to keep beating after a failed heartbeat: reclaim a 412
// lease, tolerate transient errors until the staleness ceiling, else cancel the runner.
func (h *claudeHeartbeater) onError(ctx context.Context, err error) bool {
	if ctx.Err() != nil {
		return false
	}
	if h.tryReclaimLease(err) {
		return true
	}
	if stale := time.Since(h.lastSuccess); stale > h.ttl {
		h.logger.Error("heartbeat staleness ceiling; cancelling session runner",
			slog.String("since_success", stale.String()), slog.String("ttl", h.ttl.String()), slog.Any("error", err))
		h.cancel()
		return false
	}
	h.logger.Warn("transient heartbeat error", slog.Any("error", err))
	return true
}

// tryReclaimLease adopts the server's current last_heartbeat on a 412 (the work item
// already carries a live lease — typically a prior worker). One VM owns this session,
// so we reclaim instead of dying; bounded so two live workers can't ping-pong forever.
func (h *claudeHeartbeater) tryReclaimLease(err error) bool {
	var apierr *anthropic.Error
	if !errors.As(err, &apierr) || apierr.StatusCode != 412 || h.reclaims >= heartbeatMaxReclaims {
		return false
	}
	lh := reclaimLastHeartbeat(apierr)
	if lh == "" || lh == h.last {
		return false
	}
	h.reclaims++
	h.logger.Warn("heartbeat 412; reclaiming lease from current_state",
		slog.String("last_heartbeat", lh), slog.Int("reclaim", h.reclaims))
	h.last = lh
	h.lastSuccess = time.Now()
	return true
}

// nonEmptyResultTool wraps a BetaTool so an empty tool result never reaches the
// events API as an empty text block (which the API rejects with
// "content.0.text: value is required", stalling the turn). Empty output → "(no output)".
type nonEmptyResultTool struct{ inner anthropic.BetaTool }

func (t nonEmptyResultTool) Name() string        { return t.inner.Name() }
func (t nonEmptyResultTool) Description() string { return t.inner.Description() }
func (t nonEmptyResultTool) InputSchema() anthropic.BetaToolInputSchemaParam {
	return t.inner.InputSchema()
}

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

// heartbeatErrorBody mirrors the 412 heartbeat error body just far enough to reach
// error.details.current_state.last_heartbeat — the server's authoritative lease token.
// The SDK's *anthropic.Error doesn't type this detail, so we parse it from raw JSON.
type heartbeatErrorBody struct {
	Error heartbeatErrorContent `json:"error"`
}

type heartbeatErrorContent struct {
	Details heartbeatErrorDetails `json:"details"`
}

type heartbeatErrorDetails struct {
	CurrentState heartbeatWorkCurrentState `json:"current_state"`
}

type heartbeatWorkCurrentState struct {
	LastHeartbeat string `json:"last_heartbeat"`
}

// reclaimLastHeartbeat extracts error.details.current_state.last_heartbeat from a
// 412 heartbeat error body so the worker can adopt (reclaim) the existing lease.
func reclaimLastHeartbeat(apierr *anthropic.Error) string {
	var body heartbeatErrorBody
	if json.Unmarshal([]byte(apierr.RawJSON()), &body) != nil {
		return ""
	}
	return body.Error.Details.CurrentState.LastHeartbeat
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
	return slog.New(slog.NewTextHandler(claudeLogWriter(path), &slog.HandlerOptions{Level: slog.LevelDebug}))
}

// claudeLogWriter opens the log file, falling back to os.Stderr on any failure.
func claudeLogWriter(path string) io.Writer {
	if path == "" {
		return os.Stderr
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return os.Stderr
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Logger().Warnf("claude worker: could not open %s, logging to stderr: %v", path, err)
		return os.Stderr
	}
	return f
}
