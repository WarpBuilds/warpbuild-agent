package tcache

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/proxy"
)

// startRestore runs the server against a fake backend-cache that knows one key, "hit", and a fake
// storage that serves file with Range support.
func startRestore(t *testing.T, file []byte, limitMBPerSec ...int) string {
	mux := http.NewServeMux()
	var backendURL string
	mux.HandleFunc("POST /v1/cache/get", func(w http.ResponseWriter, r *http.Request) {
		var req proxy.GetCacheRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.CacheKey != "hit" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"sub_code":"NFE_001","message":"cache miss"}`)
			return
		}
		json.NewEncoder(w).Encode(proxy.GetCacheResponse{
			Provider:   proxy.ProviderR2,
			S3:         &proxy.S3GetCacheResponse{PreSignedURL: backendURL + "/storage/file"},
			CacheEntry: &proxy.CacheEntry{CacheUserGivenKey: "hit-matched"},
		})
	})
	mux.HandleFunc("GET /storage/file", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Range") == "" {
			t.Errorf("storage was asked for the whole file in one request")
		}
		http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(file))
	})
	backend := httptest.NewServer(mux)
	t.Cleanup(backend.Close)
	backendURL = backend.URL

	opts := Options{CacheBackendHost: backend.URL, AuthToken: "token"}
	if len(limitMBPerSec) > 0 {
		opts.DownloadLimitMBPerSec = limitMBPerSec[0]
	}
	srv := httptest.NewServer(newHandler(opts))
	t.Cleanup(srv.Close)
	return srv.URL
}

func get(t *testing.T, url string, header ...string) (*http.Response, []byte) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	for i := 0; i+1 < len(header); i += 2 {
		req.Header.Set(header[i], header[i+1])
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return resp, body
}

func TestRestore(t *testing.T) {
	// Big enough for many pieces, with a last piece that is not full.
	file := randomBytes(40*pieceSize + 12345)
	base := startRestore(t, file)

	status, res := twirp(t, base, "GetCacheEntryDownloadURL", `{"key":"miss","restore_keys":["miss"],"version":"v1"}`)
	if status != 200 || res["ok"] != false || res["signed_download_url"] != "" {
		t.Fatalf("miss: status %d, %v", status, res)
	}

	status, res = twirp(t, base, "GetCacheEntryDownloadURL", `{"key":"hit","restore_keys":null,"version":"v1"}`)
	if status != 200 || res["ok"] != true || res["matched_key"] != "hit-matched" {
		t.Fatalf("hit: status %d, %v", status, res)
	}
	url := res["signed_download_url"].(string)
	if !strings.HasPrefix(url, base+"/download/") {
		t.Fatalf("unexpected download url %s", url)
	}

	// The whole file: one reply with Content-Length.
	resp, body := get(t, url)
	if resp.StatusCode != 200 || resp.ContentLength != int64(len(file)) || !bytes.Equal(body, file) {
		t.Fatalf("full download: status %d, Content-Length %d, %d bytes", resp.StatusCode, resp.ContentLength, len(body))
	}

	// A resume, with either header.
	from := 3*pieceSize + 7
	resp, body = get(t, url, "Range", fmt.Sprintf("bytes=%d-", from))
	if resp.StatusCode != 206 || !bytes.Equal(body, file[from:]) ||
		resp.Header.Get("Content-Range") != fmt.Sprintf("bytes %d-%d/%d", from, len(file)-1, len(file)) {
		t.Fatalf("resume: status %d, Content-Range %q, %d bytes", resp.StatusCode, resp.Header.Get("Content-Range"), len(body))
	}
	resp, body = get(t, url, "x-ms-range", "bytes=10-19")
	if resp.StatusCode != 206 || !bytes.Equal(body, file[10:20]) {
		t.Fatalf("x-ms-range: status %d, %d bytes", resp.StatusCode, len(body))
	}
}

func TestRestoreEmptyFile(t *testing.T) {
	base := startRestore(t, nil)
	_, res := twirp(t, base, "GetCacheEntryDownloadURL", `{"key":"hit","version":"v1"}`)
	resp, body := get(t, res["signed_download_url"].(string))
	if resp.StatusCode != 200 || len(body) != 0 {
		t.Fatalf("empty file: status %d, %d bytes", resp.StatusCode, len(body))
	}
}

func TestReady(t *testing.T) {
	base := startRestore(t, nil)
	resp, err := http.Post(base+"/_warpbuild/ready", "text/plain", strings.NewReader("https://results.example/"))
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("ready with backend up: %v, %v", resp, err)
	}

	// A backend that cannot be reached means not ready.
	srv := httptest.NewServer(newHandler(Options{CacheBackendHost: "http://127.0.0.1:1", AuthToken: "token"}))
	defer srv.Close()
	resp, err = http.Post(srv.URL+"/_warpbuild/ready", "text/plain", strings.NewReader("https://results.example/"))
	if err != nil || resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("ready with backend down: %v, %v", resp, err)
	}
}

func TestLocalCaller(t *testing.T) {
	vm := netip.MustParseAddr("10.90.16.146")
	local := []netip.Prefix{
		netip.MustParsePrefix("127.0.0.1/8"),
		netip.MustParsePrefix("10.90.16.146/20"), // the VM's network, shared with other VMs
		netip.MustParsePrefix("172.17.0.1/16"),   // Docker's default network
		netip.MustParsePrefix("172.18.0.1/16"),   // the network of a container job
		netip.MustParsePrefix("fe80::1/64"),
	}
	for caller, want := range map[string]bool{
		"127.0.0.1":           true, // the job hook
		"::1":                 true,
		"10.90.16.146":        true, // a step on the VM calling the VM's address
		"::ffff:10.90.16.146": true,
		"172.17.0.2":          true,  // BuildKit's container
		"172.18.0.3":          true,  // a step of a container job
		"10.90.16.147":        false, // a neighbour VM
		"10.90.25.130":        false,
		"8.8.8.8":             false,
	} {
		if got := localCaller(netip.MustParseAddr(caller), vm, local); got != want {
			t.Errorf("localCaller(%s) = %v, want %v", caller, got, want)
		}
	}
	if localCaller(netip.Addr{}, vm, local) {
		t.Error("an unparsed caller address must be refused")
	}
	// The VM's own IP is not known: Docker's network cannot be told apart from the VM's, so both are refused.
	if localCaller(netip.MustParseAddr("172.17.0.2"), netip.MustParseAddr("127.0.0.1"), local) {
		t.Error("no network may pass when the VM's address is not known")
	}
}

func TestGuardRefusesOtherAddresses(t *testing.T) {
	h := guard(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	for addr, want := range map[string]int{"127.0.0.1:5000": 200, "203.0.113.9:5000": 403, "": 403} {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = addr
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != want {
			t.Errorf("caller %q: status %d, want %d", addr, rec.Code, want)
		}
	}
}

func TestRestoreDownloadLimit(t *testing.T) {
	// 60 MB at 20 MB/s must take about 3 s, shared by all the pieces fetched at once.
	file := randomBytes(60_000_000)
	base := startRestore(t, file, 20)
	_, res := twirp(t, base, "GetCacheEntryDownloadURL", `{"key":"hit","version":"v1"}`)

	began := time.Now()
	_, body := get(t, res["signed_download_url"].(string))
	took := time.Since(began)
	if !bytes.Equal(body, file) {
		t.Fatalf("got %d bytes, want %d", len(body), len(file))
	}
	if took < 2500*time.Millisecond || took > 5*time.Second {
		t.Errorf("60 MB at 20 MB/s took %v, want about 3 s", took)
	}
}

func TestForward(t *testing.T) {
	github := func(name string) *httptest.Server {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, name+" "+r.URL.Path)
		}))
		t.Cleanup(s.Close)
		return s
	}
	backup, real := github("backup"), github("real")
	call := "/twirp/github.actions.results.api.v1.ArtifactService/CreateArtifact"

	// No address at all: the call cannot be forwarded.
	none := httptest.NewServer(newHandler(Options{}))
	defer none.Close()
	if resp, _ := get(t, none.URL+call); resp.StatusCode != http.StatusBadGateway {
		t.Fatalf("no address: status %d", resp.StatusCode)
	}

	// Before the job hook has called, as after a restart in the middle of a job: the backup.
	backend := httptest.NewServer(http.NotFoundHandler()) // a miss is all the ready check needs
	defer backend.Close()
	srv := httptest.NewServer(newHandler(Options{CacheBackendHost: backend.URL, GitHubResultsURL: backup.URL + "/"}))
	defer srv.Close()
	if _, body := get(t, srv.URL+call); string(body) != "backup "+call {
		t.Fatalf("before ready: %q", body)
	}
	// After the job hook has called: the address it gave.
	http.Post(srv.URL+"/_warpbuild/ready", "text/plain", strings.NewReader(real.URL+"/"))
	if _, body := get(t, srv.URL+call); string(body) != "real "+call {
		t.Fatalf("after ready: %q", body)
	}
}

func TestRestoreStalledPiece(t *testing.T) {
	// Storage hangs on one request: headers, then nothing. The piece must be fetched again, and
	// the client must never wait 5 s for data, which is when the cache step gives up.
	file := randomBytes(6 * pieceSize)
	var requests atomic.Int32
	mux := http.NewServeMux()
	var backendURL string
	mux.HandleFunc("POST /v1/cache/get", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(proxy.GetCacheResponse{
			S3:         &proxy.S3GetCacheResponse{PreSignedURL: backendURL + "/storage/file"},
			CacheEntry: &proxy.CacheEntry{CacheUserGivenKey: "hit"},
		})
	})
	mux.HandleFunc("GET /storage/file", func(w http.ResponseWriter, r *http.Request) {
		if requests.Add(1) == 3 {
			w.Header().Set("Content-Range", fmt.Sprintf("bytes 0-0/%d", len(file)))
			w.WriteHeader(http.StatusPartialContent)
			w.(http.Flusher).Flush()
			select {
			case <-r.Context().Done():
			case <-time.After(20 * time.Second):
			}
			return
		}
		http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(file))
	})
	backend := httptest.NewServer(mux)
	defer backend.Close()
	backendURL = backend.URL
	srv := httptest.NewServer(newHandler(Options{CacheBackendHost: backend.URL}))
	defer srv.Close()

	_, res := twirp(t, srv.URL, "GetCacheEntryDownloadURL", `{"key":"hit","version":"v1"}`)
	began := time.Now()
	_, body := get(t, res["signed_download_url"].(string))
	if !bytes.Equal(body, file) {
		t.Fatalf("got %d bytes, want %d", len(body), len(file))
	}
	if took := time.Since(began); took > 4*time.Second {
		t.Errorf("the download took %v with one stalled piece, want under 4 s", took)
	}
	if n := requests.Load(); n != 7 {
		t.Errorf("storage got %d requests, want 7: six pieces and one retry", n)
	}
}
