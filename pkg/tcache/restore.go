package tcache

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/proxy"
	"github.com/warpbuilds/warpbuild-agent/pkg/redact"
)

const (
	// A download is fetched from storage in pieces of pieceSize, many at once.
	pieceSize = 8 << 20
	// downloadSlots is how many pieces the whole VM fetches or holds in memory at once, across
	// every download. It caps the VM's connections to storage and the server's memory.
	downloadSlots = 16
)

// getDownloadURL handles GetCacheEntryDownloadURL: look the key up at backend-cache and hand the
// client a download address on this server.
func (s *server) getDownloadURL(w http.ResponseWriter, r *http.Request) {
	var req struct {
		twirpRequest
		RestoreKeys []string `json:"restore_keys"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		twirpError(w, http.StatusBadRequest, "malformed", err)
		return
	}

	var res proxy.GetCacheResponse
	err := s.callBackend(r.Context(), "/get", proxy.GetCacheRequest{
		CacheKey:     req.Key,
		CacheVersion: req.Version,
		RestoreKeys:  req.RestoreKeys,
	}, &res)
	if err != nil || res.CacheEntry == nil || res.S3 == nil {
		// backend-cache answers a miss with an error. Any other failure is also a miss for the job.
		log.Logger().Infof("transparent cache: miss for key %q: %v", req.Key, err)
		writeJSON(w, http.StatusOK, map[string]any{"ok": false, "signed_download_url": "", "matched_key": ""})
		return
	}

	idBytes := make([]byte, 16)
	rand.Read(idBytes)
	id := hex.EncodeToString(idBytes)
	s.downloads.Store(id, res.S3.PreSignedURL)

	writeJSON(w, http.StatusOK, map[string]any{
		"ok":                  true,
		"signed_download_url": "http://" + r.Host + "/download/" + id,
		"matched_key":         res.CacheEntry.CacheUserGivenKey,
	})
}

type piece struct {
	data []byte
	err  error
}

// download streams one file to the client as a single reply. Behind it, the file is fetched from
// storage in many pieces at once and written out in order.
func (s *server) download(w http.ResponseWriter, r *http.Request) {
	v, ok := s.downloads.Load(r.PathValue("id"))
	if !ok {
		http.Error(w, "unknown download", http.StatusNotFound)
		return
	}
	storageURL := v.(string)
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Clients resume with Range or x-ms-range: "bytes=start-" or "bytes=start-end".
	start, end, ranged := int64(0), int64(-1), false
	rangeHeader := r.Header.Get("Range")
	if rangeHeader == "" {
		rangeHeader = r.Header.Get("x-ms-range")
	}
	if from, to, ok := strings.Cut(strings.TrimPrefix(rangeHeader, "bytes="), "-"); ok && from != "" {
		if n, err := strconv.ParseInt(from, 10, 64); err == nil {
			start, ranged = n, true
			if n, err := strconv.ParseInt(to, 10, 64); err == nil {
				end = n
			}
		}
	}

	// The first piece also tells us the size of the file.
	firstEnd := start + pieceSize - 1
	if end >= 0 && end < firstEnd {
		firstEnd = end
	}
	first, total, status, err := s.fetchPiece(ctx, storageURL, start, firstEnd)
	if status == http.StatusRequestedRangeNotSatisfiable && !ranged {
		// An empty file.
		w.Header().Set("Content-Length", "0")
		return
	}
	if err != nil {
		log.Logger().Errorf("transparent cache: download failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	if end < 0 || end >= total {
		end = total - 1
	}

	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))
	if ranged {
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, total))
		w.WriteHeader(http.StatusPartialContent)
	}
	if _, err := w.Write(first); err != nil {
		return
	}

	// Start the pieces in order, each holding a slot until it is written to the client.
	pieces := make(chan chan piece, downloadSlots)
	go func() {
		defer close(pieces)
		for from := start + int64(len(first)); from <= end; from += pieceSize {
			select {
			case s.slots <- struct{}{}:
			case <-ctx.Done():
				return
			}
			to := min(from+pieceSize-1, end)
			c := make(chan piece, 1)
			pieces <- c
			go func() {
				data, _, _, err := s.fetchPiece(ctx, storageURL, from, to)
				c <- piece{data, err}
			}()
		}
	}()

	var failed error
	for c := range pieces {
		p := <-c
		<-s.slots
		if failed != nil {
			continue
		}
		if p.err != nil {
			failed = p.err
		} else if _, err := w.Write(p.data); err != nil {
			failed = err
		}
		if failed != nil {
			cancel()
		}
	}
	if failed != nil {
		// The status is already sent. Drop the connection so the client sees a short file, not a
		// complete one.
		log.Logger().Errorf("transparent cache: download stopped: %v", failed)
		panic(http.ErrAbortHandler)
	}
}

// fetchPiece gets bytes from-to of a file from storage, and the file's total size. It tries 3 times.
func (s *server) fetchPiece(ctx context.Context, storageURL string, from, to int64) (data []byte, total int64, status int, err error) {
	for attempt := 0; attempt < 3 && ctx.Err() == nil; attempt++ {
		data, total, status, err = s.fetchPieceOnce(ctx, storageURL, from, to)
		if err == nil || status == http.StatusRequestedRangeNotSatisfiable {
			break
		}
	}
	return data, total, status, err
}

func (s *server) fetchPieceOnce(ctx context.Context, storageURL string, from, to int64) ([]byte, int64, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, storageURL, nil)
	if err != nil {
		return nil, 0, 0, redact.URLError(err)
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", from, to))

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, 0, 0, redact.URLError(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK && resp.ContentLength == 0 {
		// An empty file: some storage answers a Range on it with 200, some with 416.
		return nil, 0, resp.StatusCode, nil
	}
	if resp.StatusCode != http.StatusPartialContent {
		return nil, 0, resp.StatusCode, fmt.Errorf("storage answered bytes %d-%d with status %d", from, to, resp.StatusCode)
	}
	// Content-Range: bytes 0-8388607/2097152000
	_, size, _ := strings.Cut(resp.Header.Get("Content-Range"), "/")
	total, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		return nil, 0, resp.StatusCode, fmt.Errorf("storage sent no file size: %q", resp.Header.Get("Content-Range"))
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, resp.StatusCode, err
	}
	return data, total, resp.StatusCode, nil
}
