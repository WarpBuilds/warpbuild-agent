package uploader

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

const dropConnection = -1

type fakeS3 struct {
	*httptest.Server
	mu      sync.Mutex
	script  []int
	puts    []string
	objects map[string]string
}

func newFakeS3(t *testing.T, script ...int) *fakeS3 {
	f := &fakeS3{script: script, objects: map[string]string{}}
	f.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		f.mu.Lock()
		status := http.StatusOK
		if len(f.script) > 0 {
			status, f.script = f.script[0], f.script[1:]
		}
		f.puts = append(f.puts, r.URL.Path)
		if status == http.StatusOK {
			f.objects[r.URL.Path] = string(body)
		}
		f.mu.Unlock()

		if status == dropConnection {
			conn, _, err := w.(http.Hijacker).Hijack()
			if err != nil {
				t.Errorf("hijack: %v", err)
				return
			}
			conn.Close()
			return
		}
		w.WriteHeader(status)
	}))
	t.Cleanup(f.Close)
	return f
}

type fakePresigner struct {
	*httptest.Server
	mu     sync.Mutex
	script []int
}

func newFakePresigner(t *testing.T, s3 *fakeS3, script ...int) *fakePresigner {
	f := &fakePresigner{script: script}
	f.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/runners_instance/wr_test_runner/presigned_log_upload_url") ||
			r.Header.Get("X-POLLING-SECRET") != "polling-secret" {
			http.NotFound(w, r)
			return
		}
		f.mu.Lock()
		status := http.StatusOK
		if len(f.script) > 0 {
			status, f.script = f.script[0], f.script[1:]
		}
		f.mu.Unlock()

		if status != http.StatusOK {
			w.WriteHeader(status)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		url := s3.URL + "/" + r.URL.Query().Get("log-file-name") + "?X-Amz-Signature=sig"
		if err := json.NewEncoder(w).Encode(map[string]string{"url": url}); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	t.Cleanup(f.Close)
	return f
}

func newTestUploader(t *testing.T, backendURL string) (*S3Uploader, *[]time.Duration) {
	cfg := warpbuild.NewConfiguration()
	cfg.Servers[0].URL = backendURL + "/api/v1"
	u := NewS3Uploader(t.Context(), S3UploaderOptions{
		WarpbuildAPI:  warpbuild.NewAPIClient(cfg),
		EventType:     "metrics",
		RunnerID:      "wr_test_runner",
		PollingSecret: "polling-secret",
	})
	waits := &[]time.Duration{}
	u.wait = func(ctx context.Context, d time.Duration) bool {
		*waits = append(*waits, d)
		return ctx.Err() == nil
	}
	return u, waits
}

func TestUploadRetriesTransientFailuresOnTheSameKey(t *testing.T) {
	s3 := newFakeS3(t, dropConnection, http.StatusServiceUnavailable)
	backend := newFakePresigner(t, s3)
	u, waits := newTestUploader(t, backend.URL)

	require.NoError(t, u.uploadToS3([]byte("batch")))

	require.Len(t, s3.puts, 3)
	assert.Equal(t, s3.puts[0], s3.puts[1])
	assert.Equal(t, s3.puts[0], s3.puts[2])
	assert.Equal(t, map[string]string{s3.puts[0]: "batch"}, s3.objects)
	require.Len(t, *waits, 2)
	assert.InDelta(t, float64(time.Second), float64((*waits)[0]), float64(200*time.Millisecond))
	assert.InDelta(t, float64(2*time.Second), float64((*waits)[1]), float64(400*time.Millisecond))
}

func TestUploadRetriesTooManyRequests(t *testing.T) {
	s3 := newFakeS3(t, http.StatusTooManyRequests)
	backend := newFakePresigner(t, s3)
	u, waits := newTestUploader(t, backend.URL)

	require.NoError(t, u.uploadToS3([]byte("batch")))

	assert.Len(t, s3.puts, 2)
	assert.Len(t, *waits, 1)
}

func TestUploadGetsANewURLAfterAForbidden(t *testing.T) {
	s3 := newFakeS3(t, http.StatusForbidden)
	backend := newFakePresigner(t, s3)
	u, _ := newTestUploader(t, backend.URL)

	require.NoError(t, u.uploadToS3([]byte("batch")))

	require.Len(t, s3.puts, 2)
	assert.NotEqual(t, s3.puts[0], s3.puts[1])
	assert.Equal(t, map[string]string{s3.puts[1]: "batch"}, s3.objects)
}

func TestUploadGivesUpOnAClientError(t *testing.T) {
	s3 := newFakeS3(t, http.StatusBadRequest)
	backend := newFakePresigner(t, s3)
	u, waits := newTestUploader(t, backend.URL)

	err := u.uploadToS3([]byte("batch"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "attempt 1/3")
	assert.Contains(t, err.Error(), "400")
	assert.Len(t, s3.puts, 1)
	assert.Empty(t, *waits)
}

func TestUploadGivesUpAfterThreeAttempts(t *testing.T) {
	s3 := newFakeS3(t, http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable)
	backend := newFakePresigner(t, s3)
	u, waits := newTestUploader(t, backend.URL)

	err := u.uploadToS3([]byte("batch"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "attempt 3/3")
	assert.Len(t, s3.puts, 3)
	assert.Len(t, *waits, 2)
	assert.Empty(t, s3.objects)
}

func TestUploadRetriesAFailedURLRequest(t *testing.T) {
	s3 := newFakeS3(t)
	backend := newFakePresigner(t, s3, http.StatusBadGateway)
	u, waits := newTestUploader(t, backend.URL)

	require.NoError(t, u.uploadToS3([]byte("batch")))

	assert.Len(t, s3.puts, 1)
	assert.Len(t, *waits, 1)
}

func TestUploadGivesUpWhenTheBackendRejectsTheURLRequest(t *testing.T) {
	s3 := newFakeS3(t)
	backend := newFakePresigner(t, s3, http.StatusBadRequest)
	u, waits := newTestUploader(t, backend.URL)

	err := u.uploadToS3([]byte("batch"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "attempt 1/3")
	assert.Empty(t, s3.puts)
	assert.Empty(t, *waits)
}

func TestAFailedNextURLRequestDoesNotReuseTheKey(t *testing.T) {
	s3 := newFakeS3(t)
	backend := newFakePresigner(t, s3, http.StatusOK, http.StatusInternalServerError)
	u, _ := newTestUploader(t, backend.URL)

	require.NoError(t, u.uploadToS3([]byte("first")))
	require.NoError(t, u.uploadToS3([]byte("second")))

	require.Len(t, s3.puts, 2)
	assert.NotEqual(t, s3.puts[0], s3.puts[1])
	assert.Equal(t, "first", s3.objects[s3.puts[0]])
	assert.Equal(t, "second", s3.objects[s3.puts[1]])
}

func TestUploadStopsWaitingWhenTheUploaderStops(t *testing.T) {
	s3 := newFakeS3(t, http.StatusServiceUnavailable)
	backend := newFakePresigner(t, s3)
	u, _ := newTestUploader(t, backend.URL)
	u.wait = func(ctx context.Context, d time.Duration) bool {
		u.cancel()
		return waitFor(ctx, d)
	}

	start := time.Now()
	err := u.uploadToS3([]byte("batch"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "stopped after attempt 1/3")
	assert.Len(t, s3.puts, 1)
	assert.Less(t, time.Since(start), 500*time.Millisecond)
}

func TestEachURLRequestNamesADistinctFile(t *testing.T) {
	s3 := newFakeS3(t)
	backend := newFakePresigner(t, s3)
	u, _ := newTestUploader(t, backend.URL)

	seen := map[string]bool{}
	for range 20 {
		url, err := u.fetchPresignedURL()
		require.NoError(t, err)
		assert.Falsef(t, seen[url], "URL handed out twice: %s", url)
		seen[url] = true
		assert.Contains(t, url, ".metrics.log?")
	}
}

func TestUploadToS3ErrorDoesNotLeakPresignedQuery(t *testing.T) {
	// The server drops the connection, so client.Do returns a *url.Error
	// that holds the full request URL.
	srv := newFakeS3(t, dropConnection, dropConnection, dropConnection)
	u, _ := newTestUploader(t, "http://127.0.0.1:1")
	u.setURL(srv.URL + "/org/file.log?X-Amz-Signature=secret-sig&X-Amz-Security-Token=secret-token")

	err := u.uploadToS3([]byte("data"))
	if err == nil {
		t.Fatal("expected an upload error")
	}
	if strings.Contains(err.Error(), "secret") {
		t.Errorf("upload error contains the presigned query: %v", err)
	}
	if !strings.Contains(err.Error(), "/org/file.log?REDACTED") {
		t.Errorf("upload error does not show the redacted URL: %v", err)
	}
}
