package manager

import (
	"context"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

type scriptedReply func(*http.Request) (*http.Response, error)

// scriptedTransport answers each allocation poll from a script and records
// when the request arrived. It calls done when the last reply is handed out.
type scriptedTransport struct {
	mu     sync.Mutex
	script []scriptedReply
	times  []time.Time
	done   func()
}

func (t *scriptedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.mu.Lock()
	i := len(t.times)
	t.times = append(t.times, time.Now())
	t.mu.Unlock()
	if i >= len(t.script) {
		return nil, errors.New("script exhausted")
	}
	if i == len(t.script)-1 && t.done != nil {
		t.done()
	}
	return t.script[i](req)
}

func (t *scriptedTransport) requestTimes() []time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()
	return append([]time.Time(nil), t.times...)
}

func connectionRefused(*http.Request) (*http.Response, error) {
	return nil, errors.New("dial tcp: connect: connection refused")
}

func reply(code int, body string) scriptedReply {
	return func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: code,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body)),
			Request:    req,
		}, nil
	}
}

// taking makes a reply take d, or fail with the request's context error if
// that ends first.
func taking(d time.Duration, r scriptedReply) scriptedReply {
	return func(req *http.Request) (*http.Response, error) {
		select {
		case <-time.After(d):
			return r(req)
		case <-req.Context().Done():
			return nil, req.Context().Err()
		}
	}
}

func newTestAgent(t *testing.T, rt http.RoundTripper) *agentImpl {
	t.Helper()
	_, err := log.Init(&log.InitOptions{})
	require.NoError(t, err)

	cfg := warpbuild.NewConfiguration()
	cfg.Servers[0].URL = "http://backend.test/api/v1"
	cfg.HTTPClient = &http.Client{Transport: rt}
	return &agentImpl{
		client:           warpbuild.NewAPIClient(cfg),
		id:               "ri-test",
		pollingSecret:    "secret",
		hostURL:          cfg.Servers[0].URL,
		exitFileLocation: filepath.Join(t.TempDir(), "exit.json"),
		opts:             &AgentOptions{},
	}
}

// StartAgent runs inside a synctest bubble, so the timers are the bubble's fake clock and request
// times can be asserted exactly; the transport is in-process and nothing blocks outside the bubble.

func TestStartAgentPollsAtOnceThenOnTheTicker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tr := &scriptedTransport{done: cancel, script: []scriptedReply{
		connectionRefused,
		connectionRefused,
		reply(200, `{"status":"pending"}`),
	}}
	a := newTestAgent(t, tr)

	synctest.Test(t, func(t *testing.T) {
		start := time.Now()
		require.NoError(t, a.StartAgent(ctx, &StartAgentOptions{Manager: &ManagerOptions{Provider: ProviderGithub}}))

		var got []time.Duration
		for _, at := range tr.requestTimes() {
			got = append(got, at.Sub(start))
		}
		require.Equal(t, []time.Duration{0, Interval, 2 * Interval}, got)
	})
}

func TestStartAgentReturnsWhenCancelledDuringRequest(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tr := &scriptedTransport{script: []scriptedReply{
		taking(10*time.Second, reply(200, `{"status":"pending"}`)),
	}}
	a := newTestAgent(t, tr)

	synctest.Test(t, func(t *testing.T) {
		start := time.Now()
		time.AfterFunc(Interval/2, cancel) // not at a tick, where select may take the tick first
		require.NoError(t, a.StartAgent(ctx, &StartAgentOptions{Manager: &ManagerOptions{Provider: ProviderGithub}}))
		require.Equal(t, Interval/2, time.Since(start))
		require.Len(t, tr.requestTimes(), 1)
	})
}
