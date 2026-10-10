//go:build darwin

package sandbox

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

const defaultKeepAliveInterval = 90 * time.Second

func keepAliveInterval(h http.Header) time.Duration {
	if v := h.Get("Keepalive-Ping-Interval"); v != "" {
		if secs, err := strconv.Atoi(v); err == nil && secs > 0 {
			return time.Duration(secs) * time.Second
		}
	}

	return defaultKeepAliveInterval
}

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
