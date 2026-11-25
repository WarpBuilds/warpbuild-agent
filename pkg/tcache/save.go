package tcache

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/proxy"
	"github.com/warpbuilds/warpbuild-agent/pkg/redact"
)

const (
	twirpPrefix = "/twirp/github.actions.results.api.v1.CacheService/"

	// One entry is at most maxParts parts of partSize: 10 GB, the same limit as GitHub's cache.
	maxParts = 160
	partSize = 64 << 20

	// backend-cache's code for "cache entry already exists".
	duplicateRequestCode = "FVE_002"
)

// entry is one save in progress. It lives in memory from CreateCacheEntry to FinalizeCacheEntryUpload.
type entry struct {
	key, version string
	upload       *proxy.S3ReserveCacheResponse

	mu    sync.Mutex
	parts map[int]string // part number -> ETag

	// BuildKit sends 1 MiB blocks one after another. Storage needs parts of at least 5 MiB, so
	// they are collected here and sent as one part when partSize is reached.
	buf       bytes.Buffer
	nextBlock int
	nextPart  int
}

func entryKey(key, version string) string { return key + "\x00" + version }

// callBackend calls backend-cache the same way pkg/proxy does.
func (s *server) callBackend(ctx context.Context, path string, in, out any) error {
	body, err := json.Marshal(in)
	if err != nil {
		return err
	}
	u, err := url.JoinPath(s.opts.CacheBackendHost, "/v1/cache", path)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.opts.AuthToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("backend-cache %s failed with status %d: %s", path, resp.StatusCode, data)
	}
	return json.Unmarshal(data, out)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func twirpError(w http.ResponseWriter, status int, code string, err error) {
	log.Logger().Errorf("transparent cache: %s: %v", code, err)
	writeJSON(w, status, map[string]string{"code": code, "msg": err.Error()})
}

type twirpRequest struct {
	Key     string `json:"key"`
	Version string `json:"version"`
}

// createEntry handles CreateCacheEntry: reserve the upload at backend-cache and hand the client
// an upload address on this server.
func (s *server) createEntry(w http.ResponseWriter, r *http.Request) {
	var req twirpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		twirpError(w, http.StatusBadRequest, "malformed", err)
		return
	}

	var res proxy.ReserveCacheResponse
	err := s.callBackend(r.Context(), "/reserve", proxy.ReserveCacheRequest{
		CacheKey:       req.Key,
		CacheVersion:   req.Version,
		NumberOfChunks: maxParts,
		ContentType:    "application/octet-stream",
	}, &res)
	if err != nil {
		// A saved key must be already_exists, not ok:false: BuildKit loops on a refused save.
		if strings.Contains(err.Error(), duplicateRequestCode) {
			twirpError(w, http.StatusConflict, "already_exists", err)
			return
		}
		twirpError(w, http.StatusInternalServerError, "internal", err)
		return
	}
	if res.S3 == nil || len(res.S3.PreSignedURLs) == 0 {
		twirpError(w, http.StatusInternalServerError, "internal", fmt.Errorf("storage provider %q is not supported", res.Provider))
		return
	}

	idBytes := make([]byte, 16)
	rand.Read(idBytes)
	id := hex.EncodeToString(idBytes)

	s.entries.Store(id, &entry{key: req.Key, version: req.Version, upload: res.S3, parts: map[int]string{}, nextPart: 1})
	s.entryIDs.Store(entryKey(req.Key, req.Version), id)

	// The Azure SDK in both clients accepts this shape: an Azurite-style path and a sig.
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":                true,
		"signed_upload_url": "http://" + r.Host + "/devstoreaccount1/upload/" + id + "?sig=" + id,
		"message":           "",
	})
}

// upload handles the file PUTs from the Azure SDK.
func (s *server) upload(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	v, ok := s.entries.Load(id)
	if !ok {
		http.Error(w, "unknown upload", http.StatusNotFound)
		return
	}
	e := v.(*entry)

	var err error
	q := r.URL.Query()
	switch q.Get("comp") {
	case "blocklist":
		// Nothing to do: the parts are committed in FinalizeCacheEntryUpload.
	case "block":
		// A "+" in the block id arrives as a space when the client does not escape it.
		var blockID []byte
		blockID, err = base64.StdEncoding.DecodeString(strings.ReplaceAll(q.Get("blockid"), " ", "+"))
		switch {
		case err != nil:
		case len(blockID) == 48:
			// actions/cache: 36-char uuid + 12-digit index. Big blocks, many at once: each block
			// is one part and is streamed straight through.
			var index int
			if index, err = strconv.Atoi(string(blockID[36:])); err == nil {
				err = s.putPart(r.Context(), e, index+1, r.Body, r.ContentLength)
			}
		case len(blockID) == 64:
			// BuildKit: index is big-endian at bytes 16-19.
			err = s.addBlock(r.Context(), e, int(binary.BigEndian.Uint32(blockID[16:20])), r.Body)
		default:
			err = fmt.Errorf("unknown block id of %d bytes", len(blockID))
		}
	case "":
		// The whole file in one PUT.
		err = s.putPart(r.Context(), e, 1, r.Body, r.ContentLength)
	default:
		err = fmt.Errorf("unsupported comp=%s", q.Get("comp"))
	}
	if err != nil {
		log.Logger().Errorf("transparent cache: upload failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("x-ms-request-id", id)
	w.Header().Set("Content-Length", "0")
	w.WriteHeader(http.StatusCreated)
}

// putPart streams one part to storage and records its ETag.
func (s *server) putPart(ctx context.Context, e *entry, part int, body io.Reader, size int64) error {
	etag, err := s.sendPart(ctx, e, part, body, size)
	if err != nil {
		return err
	}
	e.mu.Lock()
	e.parts[part] = etag
	e.mu.Unlock()
	return nil
}

func (s *server) sendPart(ctx context.Context, e *entry, part int, body io.Reader, size int64) (string, error) {
	if part < 1 || part > len(e.upload.PreSignedURLs) {
		return "", fmt.Errorf("part %d is over the limit of %d parts", part, len(e.upload.PreSignedURLs))
	}
	if size == 0 {
		body = http.NoBody
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, e.upload.PreSignedURLs[part-1], body)
	if err != nil {
		return "", redact.URLError(err)
	}
	req.ContentLength = size

	resp, err := s.client.Do(req)
	if err != nil {
		return "", redact.URLError(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", fmt.Errorf("storage refused part %d with status %d: %s", part, resp.StatusCode, data)
	}
	return resp.Header.Get("ETag"), nil
}

// addBlock collects one small block. Blocks must come in order.
func (s *server) addBlock(ctx context.Context, e *entry, index int, body io.Reader) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if index < e.nextBlock {
		// A retry of a block we already have.
		return nil
	}
	if index > e.nextBlock {
		return fmt.Errorf("block %d came before block %d", index, e.nextBlock)
	}
	before := e.buf.Len()
	if _, err := io.Copy(&e.buf, body); err != nil {
		e.buf.Truncate(before)
		return err
	}
	e.nextBlock++

	if e.buf.Len() >= partSize {
		return s.flush(ctx, e)
	}
	return nil
}

// flush sends the collected blocks as the next part. The caller holds e.mu.
func (s *server) flush(ctx context.Context, e *entry) error {
	etag, err := s.sendPart(ctx, e, e.nextPart, bytes.NewReader(e.buf.Bytes()), int64(e.buf.Len()))
	if err != nil {
		return err
	}
	e.parts[e.nextPart] = etag
	e.nextPart++
	e.buf.Reset()
	return nil
}

// finalizeEntry handles FinalizeCacheEntryUpload: commit the parts at backend-cache.
func (s *server) finalizeEntry(w http.ResponseWriter, r *http.Request) {
	var req twirpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		twirpError(w, http.StatusBadRequest, "malformed", err)
		return
	}
	id, ok := s.entryIDs.Load(entryKey(req.Key, req.Version))
	if !ok {
		twirpError(w, http.StatusNotFound, "not_found", fmt.Errorf("no upload in progress for key %q", req.Key))
		return
	}
	v, _ := s.entries.Load(id)
	e := v.(*entry)

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.buf.Len() > 0 {
		if err := s.flush(r.Context(), e); err != nil {
			twirpError(w, http.StatusInternalServerError, "internal", err)
			return
		}
	}

	var parts []proxy.S3CompletedPart
	for n := 1; n <= maxParts; n++ {
		if etag, ok := e.parts[n]; ok {
			partNumber := int32(n)
			parts = append(parts, proxy.S3CompletedPart{ETag: &etag, PartNumber: &partNumber})
		}
	}

	var res proxy.CommitCacheResponse
	if err := s.callBackend(r.Context(), "/commit", proxy.CommitCacheRequest{
		CacheKey:     e.key,
		CacheVersion: e.version,
		UploadKey:    e.upload.UploadKey,
		UploadID:     e.upload.UploadID,
		Parts:        parts,
		VCSType:      "github",
	}, &res); err != nil {
		twirpError(w, http.StatusInternalServerError, "internal", err)
		return
	}

	s.entries.Delete(id)
	s.entryIDs.Delete(entryKey(req.Key, req.Version))

	// The clients want entry_id as digits in a string.
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":       true,
		"entry_id": strconv.FormatInt(s.saved.Add(1), 10),
		"message":  "",
	})
}
