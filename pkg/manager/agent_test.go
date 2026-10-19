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

func repeat(n int, r scriptedReply) []scriptedReply {
	s := make([]scriptedReply, n)
	for i := range s {
		s[i] = r
	}
	return s
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

// The tests run StartAgent inside a synctest bubble, so time.After, the
// ticker and time.Now are the bubble's fake clock and request times can be
// asserted exactly. Nothing in the loop blocks outside the bubble: the
// transport is in-process and the logger writes synchronously.

func TestStartAgentPollsAtOnceAndRetriesFastUntilAnswered(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tr := &scriptedTransport{done: cancel, script: []scriptedReply{
		connectionRefused, // address not configured yet
		connectionRefused, // still not
		taking(100*time.Millisecond, reply(500, `{}`)), // an answer, even though it failed
		reply(200, `{"status":"pending"}`),             // one Interval after that answer
	}}
	a := newTestAgent(t, tr)

	synctest.Test(t, func(t *testing.T) {
		start := time.Now()
		require.NoError(t, a.StartAgent(ctx, &StartAgentOptions{Manager: &ManagerOptions{Provider: ProviderGithub}}))

		var got []time.Duration
		for _, at := range tr.requestTimes() {
			got = append(got, at.Sub(start))
		}
		require.Equal(t, []time.Duration{
			0,                          // the first poll goes out at once
			firstPollRetryInterval,     // retried firstPollRetryInterval after the failure
			2 * firstPollRetryInterval, // and again
			// The 500 arrived at 500 ms + 100 ms and counts as an answer; the
			// ticker was reset then, so the next poll is one full Interval
			// later rather than at the tick that was already due at 1 s.
			2*firstPollRetryInterval + 100*time.Millisecond + Interval,
		}, got)
	})
}

func TestStartAgentFallsBackToIntervalAfterRetryLimit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// The first poll plus firstPollRetryLimit fast retries all get no
	// response; the one after them must come from the ticker.
	tr := &scriptedTransport{done: cancel, script: repeat(firstPollRetryLimit+2, connectionRefused)}
	a := newTestAgent(t, tr)

	synctest.Test(t, func(t *testing.T) {
		start := time.Now()
		require.NoError(t, a.StartAgent(ctx, &StartAgentOptions{Manager: &ManagerOptions{Provider: ProviderGithub}}))

		times := tr.requestTimes()
		require.Len(t, times, firstPollRetryLimit+2)
		for i := 0; i <= firstPollRetryLimit; i++ {
			require.Equal(t, time.Duration(i)*firstPollRetryInterval, times[i].Sub(start), "request %d", i)
		}
		// The last fast retry failed at limit * 250 ms with no answer. The
		// ticker was reset then, so the fallback poll is one Interval later,
		// not the tick that was already due at that instant.
		last := time.Duration(firstPollRetryLimit)*firstPollRetryInterval + Interval
		require.Equal(t, last, times[firstPollRetryLimit+1].Sub(start))
	})
}

func TestStartAgentReturnsWhenCancelledDuringRetryWait(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tr := &scriptedTransport{script: []scriptedReply{connectionRefused}}
	a := newTestAgent(t, tr)

	synctest.Test(t, func(t *testing.T) {
		start := time.Now()
		// Cancel 100 ms into the 250 ms wait before the first retry.
		time.AfterFunc(100*time.Millisecond, cancel)
		require.NoError(t, a.StartAgent(ctx, &StartAgentOptions{Manager: &ManagerOptions{Provider: ProviderGithub}}))
		require.Equal(t, 100*time.Millisecond, time.Since(start))
		require.Len(t, tr.requestTimes(), 1)
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
		// Cancel 1 s into a request that would take 10 s.
		time.AfterFunc(time.Second, cancel)
		require.NoError(t, a.StartAgent(ctx, &StartAgentOptions{Manager: &ManagerOptions{Provider: ProviderGithub}}))
		require.Equal(t, time.Second, time.Since(start))
		require.Len(t, tr.requestTimes(), 1)
	})
}
