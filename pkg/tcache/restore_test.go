package tcache

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/proxy"
)

// startRestore runs the server against a fake backend-cache that knows one key, "hit", and a fake
// storage that serves file with Range support.
func startRestore(t *testing.T, file []byte) string {
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

	srv := httptest.NewServer(newHandler(Options{CacheBackendHost: backend.URL, AuthToken: "token"}))
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
