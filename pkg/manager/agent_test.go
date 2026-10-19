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
	"time"

	"github.com/stretchr/testify/require"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

type scriptedReply func(*http.Request) (*http.Response, error)

// scriptedTransport answers each allocation poll from a script, records when
// the request arrived, and calls done once the script is used up.
type scriptedTransport struct {
	mu     sync.Mutex
	script []scriptedReply
	times  []time.Time
	done   func()
}

func (t *scriptedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	i := len(t.times)
	t.times = append(t.times, time.Now())
	if i >= len(t.script) {
		return nil, errors.New("script exhausted")
	}
	if i == len(t.script)-1 {
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

func TestStartAgentPollsAtOnceAndRetriesFastUntilAnswered(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tr := &scriptedTransport{done: cancel, script: []scriptedReply{
		connectionRefused,                  // address not configured yet
		connectionRefused,                  // still not
		reply(500, `{}`),                   // an answer, even though it failed
		reply(200, `{"status":"pending"}`), // next ticker interval
	}}
	a := newTestAgent(t, tr)

	start := time.Now()
	err := a.StartAgent(ctx, &StartAgentOptions{Manager: &ManagerOptions{Provider: ProviderGithub}})
	require.NoError(t, err)

	times := tr.requestTimes()
	require.Len(t, times, 4)
	since := func(i int) time.Duration { return times[i].Sub(start) }

	// The first poll goes out at once, not after the first Interval tick.
	require.Less(t, since(0), 500*time.Millisecond)
	// Polls that got no response are retried every firstPollRetryInterval.
	require.Less(t, since(2), 2*firstPollRetryInterval+400*time.Millisecond)
	// After the first answer (a 500 counts) polling follows the ticker.
	require.GreaterOrEqual(t, since(3), Interval-50*time.Millisecond)
}
