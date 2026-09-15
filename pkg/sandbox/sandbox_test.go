//go:build darwin

package sandbox

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"connectrpc.com/connect"

	fsrpc "github.com/warpbuilds/warpbuild-agent/pkg/sandboxspec/filesystem"
	"github.com/warpbuilds/warpbuild-agent/pkg/sandboxspec/filesystem/filesystemconnect"
	procrpc "github.com/warpbuilds/warpbuild-agent/pkg/sandboxspec/process"
	"github.com/warpbuilds/warpbuild-agent/pkg/sandboxspec/process/processconnect"
)

func testServer(t *testing.T, opts Options) *httptest.Server {
	t.Helper()
	opts.applyDefaults()
	// The tests run as whoever invoked them, not as the guest user.
	opts.GuestUser = ""
	srv, err := newServer(opts)
	if err != nil {
		t.Fatalf("newServer: %v", err)
	}
	ts := httptest.NewServer(srv.handler())
	t.Cleanup(ts.Close)

	return ts
}

// tokenRoundTripper stamps the token on every request, streams included, which
// a unary interceptor would miss.
type tokenRoundTripper struct {
	token string
	next  http.RoundTripper
}

func (t *tokenRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if t.token != "" {
		r.Header.Set("X-Access-Token", t.token)
	}

	return t.next.RoundTrip(r)
}

func procClient(t *testing.T, ts *httptest.Server, token string) processconnect.ProcessClient {
	t.Helper()
	c := ts.Client()
	c.Transport = &tokenRoundTripper{token: token, next: c.Transport}

	return processconnect.NewProcessClient(c, ts.URL)
}

// collect drains a Start stream into the ordered list of event kinds plus the
// terminal event, which is what every ordering guarantee is expressed in terms of.
func collectStart(t *testing.T, stream *connect.ServerStreamForClient[procrpc.StartResponse]) ([]string, string, *procrpc.ProcessEvent_EndEvent) {
	t.Helper()
	var kinds []string
	var stdout strings.Builder
	var end *procrpc.ProcessEvent_EndEvent

	for stream.Receive() {
		ev := stream.Msg().GetEvent()
		switch e := ev.GetEvent().(type) {
		case *procrpc.ProcessEvent_Start:
			_ = e
			kinds = append(kinds, "start")
		case *procrpc.ProcessEvent_Data:
			kinds = append(kinds, "data")
			stdout.Write(e.Data.GetStdout())
		case *procrpc.ProcessEvent_Keepalive:
			kinds = append(kinds, "keepalive")
		case *procrpc.ProcessEvent_End:
			kinds = append(kinds, "end")
			end = e.End
		}
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}

	return kinds, stdout.String(), end
}

func TestStartStreamsOutputAndOrdersEvents(t *testing.T) {
	ts := testServer(t, Options{})
	c := procClient(t, ts, "")

	stream, err := c.Start(context.Background(), connect.NewRequest(&procrpc.StartRequest{
		Process: &procrpc.ProcessConfig{Cmd: "/bin/sh", Args: []string{"-c", "echo hello"}},
	}))
	if err != nil {
		t.Fatal(err)
	}

	kinds, stdout, end := collectStart(t, stream)

	if len(kinds) < 2 || kinds[0] != "start" || kinds[len(kinds)-1] != "end" {
		t.Errorf("expected start first and end last, got %v", kinds)
	}
	for _, k := range kinds[1 : len(kinds)-1] {
		if k == "start" || k == "end" {
			t.Errorf("start/end must appear exactly once, got %v", kinds)
		}
	}
	if strings.TrimSpace(stdout) != "hello" {
		t.Errorf("stdout = %q, want hello", stdout)
	}
	if end == nil || end.GetExitCode() != 0 || !end.GetExited() {
		t.Errorf("end = %+v, want exit 0 exited", end)
	}
}

// quiesceGuest in backend-core matches the literal string "exit status 3" out of
// the JSON stream, so this format is load-bearing beyond being informative.
func TestEndEventStatusMatchesGoProcessStateFormat(t *testing.T) {
	ts := testServer(t, Options{})
	c := procClient(t, ts, "")

	stream, err := c.Start(context.Background(), connect.NewRequest(&procrpc.StartRequest{
		Process: &procrpc.ProcessConfig{Cmd: "/bin/sh", Args: []string{"-c", "exit 3"}},
	}))
	if err != nil {
		t.Fatal(err)
	}
	_, _, end := collectStart(t, stream)

	if end == nil {
		t.Fatal("no end event")
	}
	if end.GetStatus() != "exit status 3" {
		t.Errorf("status = %q, want %q", end.GetStatus(), "exit status 3")
	}
	if end.GetExitCode() != 3 {
		t.Errorf("exit code = %d, want 3", end.GetExitCode())
	}
	if !end.GetExited() {
		t.Error("exited = false, want true")
	}
}

func TestStartRejectsMissingWorkingDirectory(t *testing.T) {
	ts := testServer(t, Options{})
	c := procClient(t, ts, "")

	missing := "/definitely/not/here"
	stream, err := c.Start(context.Background(), connect.NewRequest(&procrpc.StartRequest{
		Process: &procrpc.ProcessConfig{Cmd: "/bin/sh", Args: []string{"-c", "true"}, Cwd: &missing},
	}))
	if err == nil {
		stream.Receive()
		err = stream.Err()
	}
	if err == nil {
		t.Fatal("expected an error for a non-existent cwd")
	}
	if got := connect.CodeOf(err); got != connect.CodeInvalidArgument {
		t.Errorf("code = %v, want InvalidArgument", got)
	}
}

// The child must not inherit the agent's environment; only PATH carries over.
func TestEnvironmentIsBuiltFromScratch(t *testing.T) {
	t.Setenv("WARP_LEAKY_SECRET", "should-not-appear")
	ts := testServer(t, Options{})
	c := procClient(t, ts, "")

	stream, err := c.Start(context.Background(), connect.NewRequest(&procrpc.StartRequest{
		Process: &procrpc.ProcessConfig{
			Cmd:  "/bin/sh",
			Args: []string{"-c", "echo \"${WARP_LEAKY_SECRET:-absent}|${FROM_REQUEST:-unset}|${HOME:+home-set}\""},
			Envs: map[string]string{"FROM_REQUEST": "present"},
		},
	}))
	if err != nil {
		t.Fatal(err)
	}
	_, stdout, _ := collectStart(t, stream)

	got := strings.TrimSpace(stdout)
	if got != "absent|present|home-set" {
		t.Errorf("env = %q, want %q", got, "absent|present|home-set")
	}
}

// Rule 3: a detached child must not outlive the command that spawned it. This is
// the nohup shape — output redirected, so the leader's pipes close and the
// command returns immediately, leaving the child to be reaped with the group.
func TestDetachedChildIsReapedWithItsProcessGroup(t *testing.T) {
	ts := testServer(t, Options{})
	c := procClient(t, ts, "")

	marker := filepath.Join(t.TempDir(), "survivor")
	stream, err := c.Start(context.Background(), connect.NewRequest(&procrpc.StartRequest{
		Process: &procrpc.ProcessConfig{
			Cmd:  "/bin/sh",
			Args: []string{"-c", "(sleep 3; touch " + marker + ") >/dev/null 2>&1 </dev/null & echo spawned"},
		},
	}))
	if err != nil {
		t.Fatal(err)
	}
	_, _, end := collectStart(t, stream)
	if end == nil {
		t.Fatal("no end event")
	}

	time.Sleep(4 * time.Second)
	if _, err := os.Stat(marker); err == nil {
		t.Error("detached child survived its process group being reaped")
	}
}

// The complement: a background child that still holds the command's stdout keeps
// the stream open until it finishes, rather than being cut off mid-write. This is
// what makes "$(long-running) &" still deliver its output.
func TestBackgroundChildHoldingStdoutKeepsStreamOpen(t *testing.T) {
	ts := testServer(t, Options{})
	c := procClient(t, ts, "")

	stream, err := c.Start(context.Background(), connect.NewRequest(&procrpc.StartRequest{
		Process: &procrpc.ProcessConfig{
			Cmd:  "/bin/sh",
			Args: []string{"-c", "(sleep 1; echo late) & echo early"},
		},
	}))
	if err != nil {
		t.Fatal(err)
	}
	_, stdout, end := collectStart(t, stream)
	if end == nil {
		t.Fatal("no end event")
	}
	if !strings.Contains(stdout, "late") {
		t.Errorf("stdout = %q, want it to include the background child's output", stdout)
	}
}

func TestListReportsRunningProcesses(t *testing.T) {
	ts := testServer(t, Options{})
	c := procClient(t, ts, "")

	stream, err := c.Start(context.Background(), connect.NewRequest(&procrpc.StartRequest{
		Process: &procrpc.ProcessConfig{Cmd: "/bin/sh", Args: []string{"-c", "sleep 2"}},
		Tag:     strPtr("sleeper"),
	}))
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()

	// Wait for the StartEvent so the process is registered before listing.
	if !stream.Receive() {
		t.Fatalf("no start event: %v", stream.Err())
	}

	resp, err := c.List(context.Background(), connect.NewRequest(&procrpc.ListRequest{}))
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, p := range resp.Msg.GetProcesses() {
		if p.GetTag() == "sleeper" {
			found = true
		}
	}
	if !found {
		t.Errorf("running process not listed, got %+v", resp.Msg.GetProcesses())
	}
}

// A reattach arriving just after exit must still learn how the process ended,
// rather than getting NotFound.
func TestConnectReplaysRetainedExit(t *testing.T) {
	ts := testServer(t, Options{})
	c := procClient(t, ts, "")

	stream, err := c.Start(context.Background(), connect.NewRequest(&procrpc.StartRequest{
		Process: &procrpc.ProcessConfig{Cmd: "/bin/sh", Args: []string{"-c", "exit 7"}},
		Tag:     strPtr("shortlived"),
	}))
	if err != nil {
		t.Fatal(err)
	}
	_, _, end := collectStart(t, stream)
	if end == nil {
		t.Fatal("no end event")
	}

	cs, err := c.Connect(context.Background(), connect.NewRequest(&procrpc.ConnectRequest{
		Process: &procrpc.ProcessSelector{Selector: &procrpc.ProcessSelector_Tag{Tag: "shortlived"}},
	}))
	if err != nil {
		t.Fatal(err)
	}
	var sawStart, sawEnd bool
	for cs.Receive() {
		switch e := cs.Msg().GetEvent().GetEvent().(type) {
		case *procrpc.ProcessEvent_Start:
			sawStart = true
		case *procrpc.ProcessEvent_End:
			sawEnd = true
			if e.End.GetExitCode() != 7 {
				t.Errorf("replayed exit code = %d, want 7", e.End.GetExitCode())
			}
		}
	}
	if !sawStart || !sawEnd {
		t.Errorf("reattach did not replay start+end (start=%v end=%v)", sawStart, sawEnd)
	}
}

func TestHealthIsExemptFromAuthAndReturns204(t *testing.T) {
	ts := testServer(t, Options{ControlToken: "sbxt_secret"})

	resp, err := ts.Client().Get(ts.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("status = %d, want 204", resp.StatusCode)
	}
}

func TestTokenIsRequiredEverywhereElse(t *testing.T) {
	ts := testServer(t, Options{ControlToken: "sbxt_secret"})

	for _, path := range []string{"/envs", "/metrics", "/files?path=/tmp/x"} {
		resp, err := ts.Client().Get(ts.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("%s without token: status = %d, want 401", path, resp.StatusCode)
		}
	}

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/envs", nil)
	req.Header.Set("X-Access-Token", "sbxt_secret")
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("with token: status = %d, want 200", resp.StatusCode)
	}
}

func TestRPCRequiresToken(t *testing.T) {
	ts := testServer(t, Options{ControlToken: "sbxt_secret"})

	unauth := procClient(t, ts, "")
	_, err := unauth.List(context.Background(), connect.NewRequest(&procrpc.ListRequest{}))
	if err == nil {
		t.Fatal("expected unauthenticated List to fail")
	}

	auth := procClient(t, ts, "sbxt_secret")
	if _, err := auth.List(context.Background(), connect.NewRequest(&procrpc.ListRequest{})); err != nil {
		t.Fatalf("authenticated List failed: %v", err)
	}
}

func TestUploadRejectsUnsupportedContentType(t *testing.T) {
	ts := testServer(t, Options{})
	target := filepath.Join(t.TempDir(), "f.txt")

	cases := map[string]int{
		"application/octet-stream": http.StatusOK,
		"text/plain":               http.StatusBadRequest,
		"":                         http.StatusBadRequest,
	}

	for ct, want := range cases {
		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/files?path="+target, strings.NewReader("body"))
		if ct != "" {
			req.Header.Set("Content-Type", ct)
		} else {
			req.Header.Del("Content-Type")
		}
		resp, err := ts.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != want {
			t.Errorf("content-type %q: status = %d, want %d", ct, resp.StatusCode, want)
		}
	}
}

func TestUploadThenDownloadRoundTrips(t *testing.T) {
	ts := testServer(t, Options{})
	target := filepath.Join(t.TempDir(), "nested", "f.txt")

	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/files?path="+target, strings.NewReader("payload"))
	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("upload status = %d", resp.StatusCode)
	}

	// Parent directories are created on the way.
	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("uploaded file missing: %v", err)
	}
	if string(got) != "payload" {
		t.Errorf("content = %q", got)
	}

	dl, err := ts.Client().Get(ts.URL + "/files?path=" + target)
	if err != nil {
		t.Fatal(err)
	}
	defer dl.Body.Close()
	body, _ := io.ReadAll(dl.Body)
	if string(body) != "payload" {
		t.Errorf("download = %q, want payload", body)
	}
}

func TestListDirDepth(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "a", "b"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "a", "b", "deep.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "top.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	ts := testServer(t, Options{})
	c := filesystemconnect.NewFilesystemClient(ts.Client(), ts.URL)

	// Depth is unset, which must behave as depth 1 rather than listing nothing.
	resp, err := c.ListDir(context.Background(), connect.NewRequest(&fsrpc.ListDirRequest{Path: dir}))
	if err != nil {
		t.Fatal(err)
	}
	names := map[string]bool{}
	for _, e := range resp.Msg.GetEntries() {
		names[e.GetName()] = true
	}
	if !names["top.txt"] || !names["a"] {
		t.Errorf("depth 1 should list top.txt and a, got %v", names)
	}
	if names["deep.txt"] {
		t.Error("depth 1 must not descend to deep.txt")
	}

	deep, err := c.ListDir(context.Background(), connect.NewRequest(&fsrpc.ListDirRequest{Path: dir, Depth: 3}))
	if err != nil {
		t.Fatal(err)
	}
	deepNames := map[string]bool{}
	for _, e := range deep.Msg.GetEntries() {
		deepNames[e.GetName()] = true
	}
	if !deepNames["deep.txt"] {
		t.Errorf("depth 3 should reach deep.txt, got %v", deepNames)
	}
}

func TestStatReportsTypeAndSymlinkTarget(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "real.txt")
	link := filepath.Join(dir, "link.txt")
	if err := os.WriteFile(file, []byte("1234"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(file, link); err != nil {
		t.Fatal(err)
	}

	ts := testServer(t, Options{})
	c := filesystemconnect.NewFilesystemClient(ts.Client(), ts.URL)

	st, err := c.Stat(context.Background(), connect.NewRequest(&fsrpc.StatRequest{Path: file}))
	if err != nil {
		t.Fatal(err)
	}
	if st.Msg.GetEntry().GetType() != fsrpc.FileType_FILE_TYPE_FILE {
		t.Errorf("file type = %v", st.Msg.GetEntry().GetType())
	}
	if st.Msg.GetEntry().GetSize() != 4 {
		t.Errorf("size = %d, want 4", st.Msg.GetEntry().GetSize())
	}

	ls, err := c.Stat(context.Background(), connect.NewRequest(&fsrpc.StatRequest{Path: link}))
	if err != nil {
		t.Fatal(err)
	}
	if ls.Msg.GetEntry().GetType() != fsrpc.FileType_FILE_TYPE_SYMLINK {
		t.Errorf("symlink type = %v", ls.Msg.GetEntry().GetType())
	}
	if ls.Msg.GetEntry().GetSymlinkTarget() != file {
		t.Errorf("symlink target = %q, want %q", ls.Msg.GetEntry().GetSymlinkTarget(), file)
	}
}

func TestWatchDirSendsStartFrameThenEvents(t *testing.T) {
	dir := t.TempDir()
	ts := testServer(t, Options{})
	c := filesystemconnect.NewFilesystemClient(ts.Client(), ts.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := c.WatchDir(ctx, connect.NewRequest(&fsrpc.WatchDirRequest{Path: dir}))
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()

	if !stream.Receive() {
		t.Fatalf("no first frame: %v", stream.Err())
	}
	if _, ok := stream.Msg().GetEvent().(*fsrpc.WatchDirResponse_Start); !ok {
		t.Fatalf("first frame = %T, want StartEvent", stream.Msg().GetEvent())
	}

	// Only safe to write after the start frame says the watch is armed.
	if err := os.WriteFile(filepath.Join(dir, "new.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	for stream.Receive() {
		if e, ok := stream.Msg().GetEvent().(*fsrpc.WatchDirResponse_Filesystem); ok {
			if e.Filesystem.GetName() != "new.txt" {
				t.Errorf("event name = %q, want new.txt (relative to the watched dir)", e.Filesystem.GetName())
			}

			return
		}
	}
	t.Fatalf("no filesystem event received: %v", stream.Err())
}

func TestPausePrepareIsIdempotentWithoutAVolume(t *testing.T) {
	ts := testServer(t, Options{DataVolume: filepath.Join(t.TempDir(), "not-mounted")})

	for i := 0; i < 2; i++ {
		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/pause-prepare", nil)
		resp, err := ts.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("call %d: status = %d, body=%s", i, resp.StatusCode, body)
		}
		if !strings.Contains(string(body), `"unmounted":false`) {
			t.Errorf("call %d: body = %s, want unmounted false", i, body)
		}
	}
}

func strPtr(s string) *string { return &s }
