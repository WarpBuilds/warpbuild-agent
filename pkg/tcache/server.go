// Package tcache is the transparent cache server. It runs on the runner VM and the job's
// ACTIONS_RESULTS_URL points at it.
package tcache

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/netip"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/proxy"
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

	downloads sync.Map      // download id -> signed storage URL
	slots     chan struct{} // see downloadSlots

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

// VMAddress is the address jobs use to reach the server: the VM's own IP where the server listens
// on every address, so steps inside containers can reach it.
func VMAddress() string {
	if host := ListenHost(); host != "" {
		return host
	}
	// No packet is sent: this only asks the kernel which local address the default route uses.
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}

// guard answers only callers on this VM. The server has no login of its own and holds the
// runner's cache token, so a neighbour VM must not be able to use it.
func guard(next http.Handler) http.Handler {
	vm, _ := netip.ParseAddr(VMAddress())
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var local []netip.Prefix
		addrs, _ := net.InterfaceAddrs()
		for _, a := range addrs {
			if p, err := netip.ParsePrefix(a.String()); err == nil {
				local = append(local, p)
			}
		}
		caller, _ := netip.ParseAddrPort(r.RemoteAddr)
		if !localCaller(caller.Addr(), vm, local) {
			log.Logger().Warnf("transparent cache: refused caller %s", r.RemoteAddr)
			http.Error(w, "callers must be on this VM", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// localCaller reports whether caller is on this VM: loopback, one of the VM's own addresses, or
// a network that lives on the VM, such as Docker's. The network of the VM's own IP holds other
// VMs, so on that one only the VM itself passes. local is the VM's addresses with their networks.
// When the VM's own IP is not known, no network passes: the two kinds cannot be told apart.
func localCaller(caller, vm netip.Addr, local []netip.Prefix) bool {
	caller, vm = caller.Unmap(), vm.Unmap()
	if caller.IsLoopback() {
		return true
	}
	vmKnown := vm.IsValid() && !vm.IsLoopback()
	for _, p := range local {
		own := p.Addr().Unmap()
		if caller == own {
			return true
		}
		// Unmap can leave a prefix length that does not fit the address; Contains is then false.
		if vmKnown && own != vm && netip.PrefixFrom(own, p.Bits()).Contains(caller) {
			return true
		}
	}
	return false
}

func newHandler(opts Options) http.Handler {
	// Pieces of a file go to storage over many connections at once. Keep them open and resume
	// TLS sessions, because new handshakes burn CPU. HTTP/2 is off: it would put every piece on
	// one connection.
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConnsPerHost = downloadSlots
	transport.TLSClientConfig = &tls.Config{ClientSessionCache: tls.NewLRUClientSessionCache(0)}
	transport.ForceAttemptHTTP2 = false
	transport.TLSNextProto = map[string]func(string, *tls.Conn) http.RoundTripper{}

	s := &server{opts: opts, client: &http.Client{Transport: transport}, slots: make(chan struct{}, downloadSlots)}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /_warpbuild/ready", s.ready)
	mux.HandleFunc("POST "+twirpPrefix+"CreateCacheEntry", s.createEntry)
	mux.HandleFunc("POST "+twirpPrefix+"FinalizeCacheEntryUpload", s.finalizeEntry)
	mux.HandleFunc("PUT /devstoreaccount1/upload/{id}", s.upload)
	mux.HandleFunc("POST "+twirpPrefix+"GetCacheEntryDownloadURL", s.getDownloadURL)
	mux.HandleFunc("GET /download/{id}", s.download)
	mux.HandleFunc("/", s.forward)
	return mux
}

func Start(ctx context.Context, opts Options) error {
	srv := &http.Server{Addr: net.JoinHostPort(ListenHost(), strconv.Itoa(opts.Port)), Handler: guard(newHandler(opts))}
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

// ready is called by the job hook with GitHub's results address as the body. It answers 200 only
// when backend-cache can be reached with this runner's token, so the hook can fall back to GitHub.
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
	// A lookup for a key that does not exist: a miss means backend-cache is up and knows this runner.
	var res proxy.GetCacheResponse
	err = s.callBackend(r.Context(), "/get", proxy.GetCacheRequest{CacheKey: "warpbuild-ready-check", CacheVersion: "0"}, &res)
	var backendErr *backendError
	if err != nil && !(errors.As(err, &backendErr) && (backendErr.status == http.StatusNotFound || backendErr.status == http.StatusBadRequest)) {
		log.Logger().Errorf("transparent cache: not ready: %v", err)
		http.Error(w, "backend-cache is not reachable", http.StatusServiceUnavailable)
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
