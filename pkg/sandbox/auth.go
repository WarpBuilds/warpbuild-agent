//go:build darwin

package sandbox

import (
	"crypto/subtle"
	"net/http"
	"time"
)

func basicAuthUsername(h http.Header) string {
	user, _, ok := (&http.Request{Header: h}).BasicAuth()
	if !ok {
		return ""
	}

	return user
}

type tokenAuth struct {
	token string
	next  http.Handler
}

func (a *tokenAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	presented := r.Header.Get("X-Access-Token")
	if a.token == "" || subtle.ConstantTimeCompare([]byte(presented), []byte(a.token)) != 1 {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Connection", "close")
		_ = http.NewResponseController(w).SetReadDeadline(time.Now())
		http.Error(w, "unauthorized", http.StatusUnauthorized)

		return
	}

	a.next.ServeHTTP(w, r)
}
