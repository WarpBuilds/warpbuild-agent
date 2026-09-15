//go:build darwin

package sandbox

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"syscall"
	"time"

	"connectrpc.com/connect"

	rpc "github.com/warpbuilds/warpbuild-agent/pkg/sandboxspec/process"
	"github.com/warpbuilds/warpbuild-agent/pkg/sandboxspec/process/processconnect"
)

const terminatedRetention = 30 * time.Second

type terminatedProc struct {
	end *rpc.ProcessEvent_EndEvent
	at  time.Time
	tag string
}

type processService struct {
	processconnect.UnimplementedProcessHandler

	users *userCache

	mu         sync.Mutex
	procs      map[int]*procHandler
	terminated map[int]*terminatedProc
}

func newProcessService(users *userCache) *processService {
	return &processService{
		users:      users,
		procs:      make(map[int]*procHandler),
		terminated: make(map[int]*terminatedProc),
	}
}

// processTimeout bounds the spawned process, not the request. Connect emits this
// header from any client-side deadline, and e2b's SDKs rely on it as the command
// timeout, so it must outlive the HTTP call.
func processTimeout(h http.Header) (time.Duration, error) {
	v := h.Get("Connect-Timeout-Ms")
	if v == "" {
		return 0, nil
	}
	ms, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}

	return time.Duration(ms) * time.Millisecond, nil
}

func (s *processService) register(h *procHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// A reused pid, or a reused tag, must not resolve to a stale exit.
	delete(s.terminated, h.pid)
	if h.tag != "" {
		for pid, t := range s.terminated {
			if t.tag == h.tag {
				delete(s.terminated, pid)
			}
		}
	}
	s.procs[h.pid] = h
}

func (s *processService) retire(h *procHandler, end *rpc.ProcessEvent_EndEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if cur, ok := s.procs[h.pid]; ok && cur == h {
		delete(s.procs, h.pid)
	}
	s.terminated[h.pid] = &terminatedProc{end: end, at: time.Now(), tag: h.tag}
	for pid, t := range s.terminated {
		if time.Since(t.at) > terminatedRetention {
			delete(s.terminated, pid)
		}
	}
}

func (s *processService) resolve(sel *rpc.ProcessSelector) (*procHandler, *terminatedProc, int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch v := sel.GetSelector().(type) {
	case *rpc.ProcessSelector_Pid:
		pid := int(v.Pid)
		if h, ok := s.procs[pid]; ok {
			return h, nil, pid
		}
		if t, ok := s.terminated[pid]; ok && time.Since(t.at) <= terminatedRetention {
			return nil, t, pid
		}
	case *rpc.ProcessSelector_Tag:
		for _, h := range s.procs {
			if h.tag == v.Tag {
				return h, nil, h.pid
			}
		}
		for pid, t := range s.terminated {
			if t.tag == v.Tag && time.Since(t.at) <= terminatedRetention {
				return nil, t, pid
			}
		}
	}

	return nil, nil, 0
}

// liveProc resolves a selector to a running process, ignoring retained exits.
func (s *processService) liveProc(sel *rpc.ProcessSelector) *procHandler {
	h, _, _ := s.resolve(sel)

	return h
}

func (s *processService) Start(
	ctx context.Context,
	req *connect.Request[rpc.StartRequest],
	stream *connect.ServerStream[rpc.StartResponse],
) error {
	timeout, err := processTimeout(req.Header())
	if err != nil {
		return connect.NewError(connect.CodeInvalidArgument, err)
	}

	u := s.users.lookup(basicAuthUsername(req.Header()))

	// Unset means true, for clients that predate the field.
	stdin := true
	if req.Msg.Stdin != nil {
		stdin = req.Msg.GetStdin()
	}

	h, err := startProcess(spawnOptions{
		Config:  req.Msg.GetProcess(),
		PTY:     req.Msg.GetPty(),
		Tag:     req.Msg.GetTag(),
		Stdin:   stdin,
		User:    u,
		Timeout: timeout,
	})
	if err != nil {
		return connect.NewError(connect.CodeInvalidArgument, err)
	}

	s.register(h)
	go h.wait(func(end *rpc.ProcessEvent_EndEvent) { s.retire(h, end) })

	id, events, ok := h.events.subscribe()
	if !ok {
		return connect.NewError(connect.CodeInternal, errProcessGone)
	}
	defer h.events.unsubscribe(id)

	if err := stream.Send(&rpc.StartResponse{Event: &rpc.ProcessEvent{
		Event: &rpc.ProcessEvent_Start{Start: &rpc.ProcessEvent_StartEvent{Pid: uint32(h.pid)}},
	}}); err != nil {
		return err
	}

	return s.pump(ctx, req.Header(), events, func(ev *rpc.ProcessEvent) error {
		return stream.Send(&rpc.StartResponse{Event: ev})
	})
}

func (s *processService) Connect(
	ctx context.Context,
	req *connect.Request[rpc.ConnectRequest],
	stream *connect.ServerStream[rpc.ConnectResponse],
) error {
	h, term, pid := s.resolve(req.Msg.GetProcess())
	if h == nil && term == nil {
		return connect.NewError(connect.CodeNotFound, errProcessNotFound)
	}

	send := func(ev *rpc.ProcessEvent) error {
		return stream.Send(&rpc.ConnectResponse{Event: ev})
	}

	if err := send(&rpc.ProcessEvent{
		Event: &rpc.ProcessEvent_Start{Start: &rpc.ProcessEvent_StartEvent{Pid: uint32(pid)}},
	}); err != nil {
		return err
	}

	// The process already exited; replay the retained terminal event so a late
	// reattach still learns how it ended.
	if h == nil {
		return send(&rpc.ProcessEvent{Event: &rpc.ProcessEvent_End{End: term.end}})
	}

	id, events, ok := h.events.subscribe()
	if !ok {
		if end := h.endEvent.Load(); end != nil {
			return send(&rpc.ProcessEvent{Event: &rpc.ProcessEvent_End{End: end}})
		}

		return nil
	}
	defer h.events.unsubscribe(id)

	return s.pump(ctx, req.Header(), events, send)
}

// pump forwards process events, ending the stream on the terminal event.
func (s *processService) pump(
	ctx context.Context,
	hdr http.Header,
	events <-chan *rpc.ProcessEvent,
	send func(*rpc.ProcessEvent) error,
) error {
	return pumpWithKeepalive(ctx, events, keepAliveInterval(hdr), send,
		func() error {
			return send(&rpc.ProcessEvent{
				Event: &rpc.ProcessEvent_Keepalive{Keepalive: &rpc.ProcessEvent_KeepAlive{}},
			})
		},
		func(ev *rpc.ProcessEvent) bool {
			_, done := ev.GetEvent().(*rpc.ProcessEvent_End)

			return done
		})
}

func (s *processService) List(
	ctx context.Context,
	req *connect.Request[rpc.ListRequest],
) (*connect.Response[rpc.ListResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]*rpc.ProcessInfo, 0, len(s.procs))
	for _, h := range s.procs {
		out = append(out, h.info())
	}

	return connect.NewResponse(&rpc.ListResponse{Processes: out}), nil
}

func (s *processService) Update(
	ctx context.Context,
	req *connect.Request[rpc.UpdateRequest],
) (*connect.Response[rpc.UpdateResponse], error) {
	h := s.liveProc(req.Msg.GetProcess())
	if h == nil {
		return nil, connect.NewError(connect.CodeNotFound, errProcessNotFound)
	}
	if p := req.Msg.GetPty(); p != nil {
		if err := h.resize(p.GetSize()); err != nil {
			return nil, connect.NewError(connect.CodeFailedPrecondition, err)
		}
	}

	return connect.NewResponse(&rpc.UpdateResponse{}), nil
}

func (s *processService) SendSignal(
	ctx context.Context,
	req *connect.Request[rpc.SendSignalRequest],
) (*connect.Response[rpc.SendSignalResponse], error) {
	h := s.liveProc(req.Msg.GetProcess())
	if h == nil {
		return nil, connect.NewError(connect.CodeNotFound, errProcessNotFound)
	}

	var sig syscall.Signal
	switch req.Msg.GetSignal() {
	case rpc.Signal_SIGNAL_SIGTERM:
		sig = syscall.SIGTERM
	case rpc.Signal_SIGNAL_SIGKILL:
		sig = syscall.SIGKILL
	default:
		return nil, connect.NewError(connect.CodeUnimplemented, errUnsupportedSignal)
	}

	if err := h.signal(sig); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&rpc.SendSignalResponse{}), nil
}

// running reports the pids the agent still owns, for the pause path.
func (s *processService) running() []*procHandler {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]*procHandler, 0, len(s.procs))
	for _, h := range s.procs {
		out = append(out, h)
	}

	return out
}
