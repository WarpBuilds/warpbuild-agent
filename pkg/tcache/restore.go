package tcache

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	// Bytes from storage are read limitChunk at a time. With a download limit, there is a pause in between.
	limitChunk = 256 << 10
	// A piece is asked for at most pieceTries times: again when a try fails, and next to a try
	// that is slow. Some connections to storage are many times slower than the rest, and the
	// cache step gives up on a download after 5 seconds with no data.
	pieceTries = 3
	// minSlowAfter is the shortest time after which a try counts as slow.
	minSlowAfter = 1500 * time.Millisecond
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

	// The time spent waiting for the next piece from storage, and for the client to take a piece.
	// One of them is where a slow download loses its time.
	began := time.Now()
	var onStorage, onClient time.Duration
	var failed error
	for c := range pieces {
		t := time.Now()
		p := <-c
		onStorage += time.Since(t)
		<-s.slots
		if failed != nil {
			continue
		}
		if p.err != nil {
			failed = p.err
		} else {
			t = time.Now()
			_, failed = w.Write(p.data)
			onClient += time.Since(t)
		}
		if failed != nil {
			cancel()
		}
	}
	log.Logger().Infof("transparent cache: download of %d bytes ended after %s: waited %s for storage, %s for the client",
		end-start+1, time.Since(began).Round(time.Millisecond), onStorage.Round(time.Millisecond), onClient.Round(time.Millisecond))
	if failed != nil {
		// The status is already sent. Drop the connection so the client sees a short file, not a
		// complete one.
		log.Logger().Errorf("transparent cache: download stopped: %v", failed)
		panic(http.ErrAbortHandler)
	}
}

// fetchPiece gets bytes from-to of a file from storage, and the file's total size. A try that
// fails is repeated. A try that is slow gets a second one next to it on another connection, and
// the first to finish wins: the slow connection is then closed and not used again.
func (s *server) fetchPiece(ctx context.Context, storageURL string, from, to int64) ([]byte, int64, int, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	type result struct {
		data   []byte
		total  int64
		status int
		err    error
	}
	results := make(chan result, pieceTries)
	try := func() {
		go func() {
			data, total, status, err := s.fetchPieceOnce(ctx, storageURL, from, to)
			results <- result{data, total, status, err}
		}()
	}
	try()
	slow := time.NewTimer(s.slowAfter)
	defer slow.Stop()

	for tries, running := 1, 1; ; {
		select {
		case r := <-results:
			running--
			if r.err == nil || r.status == http.StatusRequestedRangeNotSatisfiable {
				return r.data, r.total, r.status, r.err
			}
			if ctx.Err() == nil && tries < pieceTries {
				log.Logger().Warnf("transparent cache: bytes %d-%d from storage failed, try %d: %v", from, to, tries, r.err)
				try()
				tries, running = tries+1, running+1
			} else if running == 0 {
				return nil, 0, r.status, r.err
			}
		case <-slow.C:
			if tries < pieceTries {
				log.Logger().Warnf("transparent cache: bytes %d-%d from storage are slow, try %d on another connection", from, to, tries+1)
				try()
				tries, running = tries+1, running+1
				slow.Reset(s.slowAfter)
			}
		}
	}
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

	// Read in chunks, no faster than the VM's download limit. The limit is shared by every piece
	// of every download, so many restores at once on a host cannot fill its network line.
	var buf bytes.Buffer
	for {
		n, err := io.CopyN(&buf, resp.Body, limitChunk)
		if err == io.EOF {
			return buf.Bytes(), total, resp.StatusCode, nil
		}
		if err != nil {
			return nil, 0, resp.StatusCode, err
		}
		if s.limit != nil {
			if err := s.limit.WaitN(ctx, int(n)); err != nil {
				return nil, 0, resp.StatusCode, err
			}
		}
	}
}
