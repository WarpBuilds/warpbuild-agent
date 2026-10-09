//go:build darwin

package sandbox

import (
	"crypto/subtle"
	"net/http"
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
	if a.token == "" || r.URL.Path == healthPath {
		a.next.ServeHTTP(w, r)

		return
	}

	presented := r.Header.Get("X-Access-Token")
	if presented == "" {
		presented = r.URL.Query().Get("access_token")
	}

	if subtle.ConstantTimeCompare([]byte(presented), []byte(a.token)) != 1 {
		w.Header().Set("Cache-Control", "no-store")
		http.Error(w, "unauthorized", http.StatusUnauthorized)

		return
	}

	a.next.ServeHTTP(w, r)
}
