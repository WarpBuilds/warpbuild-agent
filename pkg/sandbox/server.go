//go:build darwin

package sandbox

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os/user"
	"sync"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/sandboxspec/filesystem/filesystemconnect"
	"github.com/warpbuilds/warpbuild-agent/pkg/sandboxspec/process/processconnect"
)

const (
	idleTimeout       = 60 * time.Second
	readHeaderTimeout = 10 * time.Second
)

type Server struct {
	opts  Options
	users *userCache
	procs *processService
	fs    *filesystemService

	pauseMu sync.Mutex
}

func newServer(opts Options) (*Server, error) {
	u, err := user.Lookup(opts.GuestUser)
	if err != nil {
		u, err = user.Current()
		if err != nil {
			return nil, fmt.Errorf("resolve guest user: %w", err)
		}
	}

	users := newUserCache(u)

	return &Server{
		opts:  opts,
		users: users,
		procs: newProcessService(users),
		fs:    newFilesystemService(users),
	}, nil
}

func (s *Server) handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc(healthPath, s.handleHealth)
	mux.HandleFunc("/metrics", s.handleMetrics)
	mux.HandleFunc("/envs", s.handleEnvs)
	mux.HandleFunc(filesPath, s.handleFiles)
	mux.HandleFunc("/pause-prepare", s.handlePausePrepare)

	procPath, procHandler := processconnect.NewProcessHandler(s.procs)
	mux.Handle(procPath, procHandler)

	fsPath, fsHandler := filesystemconnect.NewFilesystemHandler(s.fs)
	mux.Handle(fsPath, fsHandler)

	return &tokenAuth{token: s.opts.ControlToken, next: mux}
}

func tlsConfig(opts Options) (*tls.Config, error) {
	cert, err := tls.X509KeyPair(opts.TLSCert, opts.TLSKey)
	if err != nil {
		return nil, fmt.Errorf("sandbox TLS certificate: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
		NextProtos:   []string{"http/1.1"},
	}, nil
}

func newHTTPServer(h http.Handler) *http.Server {
	return &http.Server{
		Handler:           h,
		IdleTimeout:       idleTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
	}
}

func Serve(ctx context.Context, opts Options) error {
	opts.applyDefaults()

	tlsCfg, err := tlsConfig(opts)
	if err != nil {
		return err
	}

	srv, err := newServer(opts)
	if err != nil {
		return err
	}

	ln, err := srv.listen()
	if err != nil {
		return err
	}

	httpSrv := newHTTPServer(srv.handler())

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = httpSrv.Shutdown(shutdownCtx)
	}()

	if opts.ControlToken == "" {
		log.Logger().Errorf("sandbox: no control token configured; every request is rejected")
	}
	log.Logger().Infof("sandbox: serving TLS on %s", ln.Addr())

	if err := httpSrv.Serve(tls.NewListener(ln, tlsCfg)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) listen() (net.Listener, error) {
	if s.opts.ListenAddr != "" {
		return net.Listen("tcp", s.opts.ListenAddr)
	}

	return listenVsock(uint32(s.opts.VsockPort))
}
