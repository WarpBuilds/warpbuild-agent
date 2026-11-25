// Package tcache is the transparent cache server. It runs on the runner VM and the job's
// ACTIONS_RESULTS_URL points at it.
package tcache

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

type Options struct {
	Port int
	// CacheBackendHost and AuthToken are used to call backend-cache.
	CacheBackendHost string
	AuthToken        string
}

type server struct {
	opts   Options
	client *http.Client

	entries  sync.Map // upload id -> *entry
	entryIDs sync.Map // key + version -> upload id
	saved    atomic.Int64

	// githubURL is GitHub's real results address for this job. The job hook sets it once.
	githubURL atomic.Pointer[url.URL]
}

// ListenHost is the address the server listens on: every address on Linux, so steps inside
// containers can reach it, and loopback on Mac, where jobs never run in containers.
func ListenHost() string {
	if runtime.GOOS == "darwin" {
		return "127.0.0.1"
	}
	return ""
}

func newHandler(opts Options) http.Handler {
	s := &server{opts: opts, client: &http.Client{}}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /_warpbuild/ready", s.ready)
	mux.HandleFunc("POST "+twirpPrefix+"CreateCacheEntry", s.createEntry)
	mux.HandleFunc("POST "+twirpPrefix+"FinalizeCacheEntryUpload", s.finalizeEntry)
	mux.HandleFunc("PUT /devstoreaccount1/upload/{id}", s.upload)
	mux.HandleFunc("/", s.forward)
	return mux
}

func Start(ctx context.Context, opts Options) error {
	srv := &http.Server{Addr: net.JoinHostPort(ListenHost(), strconv.Itoa(opts.Port)), Handler: newHandler(opts)}
	go func() {
		<-ctx.Done()
		srv.Close()
	}()

	log.Logger().Infof("transparent cache server listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

// ready is called by the job hook with GitHub's results address as the body.
func (s *server) ready(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 4096))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u, err := url.Parse(strings.TrimSpace(string(body)))
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" {
		http.Error(w, "body must be GitHub's results URL", http.StatusBadRequest)
		return
	}
	// The first caller wins: the hook runs before any step of the job.
	s.githubURL.CompareAndSwap(nil, u)
	fmt.Fprintln(w, "ok")
}

// forward sends every call that is not a cache call to GitHub unchanged, with the caller's headers.
func (s *server) forward(w http.ResponseWriter, r *http.Request) {
	target := s.githubURL.Load()
	if target == nil {
		http.Error(w, "GitHub's results URL is not set", http.StatusBadGateway)
		return
	}
	log.Logger().Infof("forwarding to GitHub: %s %s", r.Method, r.URL.Path)
	(&httputil.ReverseProxy{Rewrite: func(pr *httputil.ProxyRequest) { pr.SetURL(target) }}).ServeHTTP(w, r)
}
