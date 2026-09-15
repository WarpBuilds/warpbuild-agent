//go:build darwin

package sandbox

import (
	"context"
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

// idleTimeout is generous because exec and watch streams are long-lived and
// legitimately silent; the keepalive frames are what prove liveness.
const idleTimeout = 640 * time.Second

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

// Serve runs the sandbox data plane until ctx is cancelled. It listens on
// AF_VSOCK so the host can reach it without guest networking, which is what lets
// a sandbox answer before it has an IP.
func Serve(ctx context.Context, opts Options) error {
	opts.applyDefaults()

	srv, err := newServer(opts)
	if err != nil {
		return err
	}

	ln, err := srv.listen()
	if err != nil {
		return err
	}

	httpSrv := &http.Server{
		Handler:     srv.handler(),
		IdleTimeout: idleTimeout,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = httpSrv.Shutdown(shutdownCtx)
	}()

	if opts.ControlToken == "" {
		log.Logger().Warnf("sandbox: no control token configured; the data plane is unauthenticated")
	}
	log.Logger().Infof("sandbox: serving on %s", ln.Addr())

	if err := httpSrv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
