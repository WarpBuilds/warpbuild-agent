package tcache

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/proxy"
)

func TestMain(m *testing.M) {
	log.Init(&log.InitOptions{})
	os.Exit(m.Run())
}

// fakeBackend stands in for backend-cache and for storage.
type fakeBackend struct {
	*httptest.Server
	mu        sync.Mutex
	parts     map[int][]byte
	committed []proxy.S3CompletedPart
	saved     bool
}

func newFakeBackend(t *testing.T) *fakeBackend {
	f := &fakeBackend{parts: map[int][]byte{}}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/cache/reserve", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer token" {
			t.Errorf("reserve: wrong Authorization header")
		}
		f.mu.Lock()
		defer f.mu.Unlock()
		if f.saved {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"code":"FVE","sub_code":"FVE_002","message":"cache entry already exists"}`)
			return
		}
		var urls []string
		for n := 1; n <= maxParts; n++ {
			urls = append(urls, fmt.Sprintf("%s/storage/%d", f.URL, n))
		}
		json.NewEncoder(w).Encode(proxy.ReserveCacheResponse{Provider: proxy.ProviderR2, S3: &proxy.S3ReserveCacheResponse{PreSignedURLs: urls, UploadKey: "k", UploadID: "u"}})
	})
	mux.HandleFunc("PUT /storage/{n}", func(w http.ResponseWriter, r *http.Request) {
		var n int
		fmt.Sscan(r.PathValue("n"), &n)
		data, _ := io.ReadAll(r.Body)
		if int64(len(data)) != r.ContentLength {
			t.Errorf("part %d: no Content-Length", n)
		}
		f.mu.Lock()
		f.parts[n] = data
		f.mu.Unlock()
		w.Header().Set("ETag", fmt.Sprintf(`"etag-%d"`, n))
	})
	mux.HandleFunc("POST /v1/cache/commit", func(w http.ResponseWriter, r *http.Request) {
		var req proxy.CommitCacheRequest
		json.NewDecoder(r.Body).Decode(&req)
		f.mu.Lock()
		f.committed, f.saved = req.Parts, true
		f.mu.Unlock()
		fmt.Fprint(w, `{}`)
	})
	f.Server = httptest.NewServer(mux)
	t.Cleanup(f.Close)
	return f
}

// file returns what storage holds after the commit: the committed parts joined in order.
func (f *fakeBackend) file(t *testing.T) []byte {
	var out []byte
	for i, p := range f.committed {
		if int(*p.PartNumber) != i+1 || *p.ETag != fmt.Sprintf(`"etag-%d"`, i+1) {
			t.Fatalf("committed part %d is wrong: number %d etag %s", i+1, *p.PartNumber, *p.ETag)
		}
		out = append(out, f.parts[i+1]...)
	}
	return out
}

func twirp(t *testing.T, base, method, body string) (int, map[string]any) {
	resp, err := http.Post(base+twirpPrefix+method, "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var out map[string]any
	json.NewDecoder(resp.Body).Decode(&out)
	return resp.StatusCode, out
}

func put(t *testing.T, url string, body []byte) {
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated || resp.Header.Get("x-ms-request-id") == "" {
		t.Errorf("PUT got status %d, x-ms-request-id %q", resp.StatusCode, resp.Header.Get("x-ms-request-id"))
	}
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

func startSave(t *testing.T) (*fakeBackend, string, string) {
	backend := newFakeBackend(t)
	srv := httptest.NewServer(newHandler(Options{CacheBackendHost: backend.URL, AuthToken: "token"}))
	t.Cleanup(srv.Close)

	status, res := twirp(t, srv.URL, "CreateCacheEntry", `{"key":"k1","version":"v1"}`)
	if status != 200 || res["ok"] != true {
		t.Fatalf("create: status %d, %v", status, res)
	}
	uploadURL := res["signed_upload_url"].(string)
	if !strings.HasPrefix(uploadURL, srv.URL+"/devstoreaccount1/upload/") || !strings.Contains(uploadURL, "?sig=") {
		t.Fatalf("unexpected upload url %s", uploadURL)
	}
	return backend, srv.URL, uploadURL
}

func finishSave(t *testing.T, backend *fakeBackend, base string, sizeBytes string, want []byte) {
	status, res := twirp(t, base, "FinalizeCacheEntryUpload", `{"key":"k1","version":"v1","size_bytes":`+sizeBytes+`}`)
	if status != 200 || res["ok"] != true || res["entry_id"] != "1" {
		t.Fatalf("finalize: status %d, %v", status, res)
	}
	if !bytes.Equal(backend.file(t), want) {
		t.Fatalf("stored file differs from what was uploaded")
	}

	// A saved key must be refused with already_exists.
	status, res = twirp(t, base, "CreateCacheEntry", `{"key":"k1","version":"v1"}`)
	if status != http.StatusConflict || res["code"] != "already_exists" {
		t.Fatalf("second create: status %d, %v", status, res)
	}
}

// actions/cache: big blocks, many at once, in any order. size_bytes is a string.
func TestSaveToolkit(t *testing.T) {
	backend, base, uploadURL := startSave(t)

	blocks := [][]byte{randomBytes(6 << 20), randomBytes(6 << 20), randomBytes(100)}
	var wg sync.WaitGroup
	for i := len(blocks) - 1; i >= 0; i-- {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "%s%012d", "4c7a2e8e-11f6-4a0c-9d3c-5b5e3a1f0b77", i))
			put(t, uploadURL+"&comp=block&blockid="+id, blocks[i])
		}()
	}
	wg.Wait()
	put(t, uploadURL+"&comp=blocklist", []byte("<BlockList></BlockList>"))

	finishSave(t, backend, base, `"12583012"`, bytes.Join(blocks, nil))
	if len(backend.committed) != 3 {
		t.Fatalf("want one part per block, got %d parts", len(backend.committed))
	}
}

// BuildKit: 1 MiB blocks one after another. size_bytes is a number.
func TestSaveBuildKit(t *testing.T) {
	backend, base, uploadURL := startSave(t)

	var want []byte
	for i := 0; i < 70; i++ {
		block := randomBytes(1 << 20)
		want = append(want, block...)
		id := make([]byte, 64)
		binary.BigEndian.PutUint32(id[16:], uint32(i))
		url := uploadURL + "&comp=block&blockid=" + base64.URLEncoding.EncodeToString(id)
		url = strings.NewReplacer("-", "%2B", "_", "%2F").Replace(url)
		put(t, url, block)
		if i == 5 {
			put(t, url, block) // a retried block must not be stored twice
		}
	}
	put(t, uploadURL+"&comp=blocklist", nil)

	finishSave(t, backend, base, "73400320", want)
	if len(backend.committed) != 2 || len(backend.parts[1]) != partSize {
		t.Fatalf("want a %d byte part and a small last part, got %d parts", partSize, len(backend.committed))
	}
}

// A small file comes as one PUT with no block id.
func TestSaveSinglePut(t *testing.T) {
	backend, base, uploadURL := startSave(t)
	want := randomBytes(1000)
	put(t, uploadURL, want)
	finishSave(t, backend, base, "1000", want)
}
