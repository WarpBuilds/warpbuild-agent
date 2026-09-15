//go:build darwin

package sandbox

import (
	"crypto/subtle"
	"net/http"
)

// basicAuthUsername reads the username half of a Basic header. It selects which
// user's home resolves ~ and relative paths; the password is not consulted,
// because the agent runs as a single unprivileged user and never setuids.
func basicAuthUsername(h http.Header) string {
	// Through the stdlib parser so the case-insensitive scheme prefix RFC 7617
	// allows is handled the same way every other Go server handles it.
	user, _, ok := (&http.Request{Header: h}).BasicAuth()
	if !ok {
		return ""
	}

	return user
}

// tokenAuth gates every route but /health on a shared per-sandbox token. An
// empty configured token leaves the agent open, which is only for local runs.
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
