//go:build darwin

package sandbox

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

const defaultKeepAliveInterval = 90 * time.Second

// keepAliveInterval honours the per-request override, in seconds, and otherwise
// pings every 90s of silence so an idle stream is not mistaken for a dead one.
func keepAliveInterval(h http.Header) time.Duration {
	if v := h.Get("Keepalive-Ping-Interval"); v != "" {
		if secs, err := strconv.Atoi(v); err == nil && secs > 0 {
			return time.Duration(secs) * time.Second
		}
	}

	return defaultKeepAliveInterval
}

// pumpWithKeepalive forwards events until the source closes, injecting a
// keepalive whenever the source has been quiet for the ping interval. Both
// streaming planes share it so the "only genuine silence produces a keepalive"
// rule cannot drift between them.
//
// done reports that an event is terminal, so a stream can close on it rather
// than waiting for the source; a nil done means the source closing is the only
// end.
func pumpWithKeepalive[T any](
	ctx context.Context,
	events <-chan T,
	ping time.Duration,
	send func(T) error,
	keepalive func() error,
	done func(T) bool,
) error {
	ticker := time.NewTicker(ping)
	defer ticker.Stop()

	last := time.Now()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev, ok := <-events:
			if !ok {
				return nil
			}
			if err := send(ev); err != nil {
				return err
			}
			// Deliberately not ticker.Reset: resetting re-sifts the runtime timer
			// heap on every chunk of a noisy process, where all we need is to
			// suppress a keepalive that would have been redundant.
			last = time.Now()
			if done != nil && done(ev) {
				return nil
			}
		case <-ticker.C:
			if time.Since(last) < ping {
				continue
			}
			if err := keepalive(); err != nil {
				return err
			}
			last = time.Now()
		}
	}
}
