//go:build darwin

package sandbox

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

const (
	healthPath = "/health"
	filesPath  = "/files"
)

// restEntryInfo is the upload response shape. It is narrower than the RPC
// EntryInfo and the type is a lowercase literal, not an enum name.
type restEntryInfo struct {
	Path     string            `json:"path"`
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type metricsPayload struct {
	TS          int64   `json:"ts"`
	CPUCount    int     `json:"cpu_count"`
	CPUUsedPct  float64 `json:"cpu_used_pct"`
	MemTotal    uint64  `json:"mem_total"`
	MemUsed     uint64  `json:"mem_used"`
	MemCache    uint64  `json:"mem_cache"`
	MemTotalMiB uint64  `json:"mem_total_mib"`
	MemUsedMiB  uint64  `json:"mem_used_mib"`
	DiskUsed    uint64  `json:"disk_used"`
	DiskTotal   uint64  `json:"disk_total"`
}

func writeJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, code int, err error) {
	writeJSON(w, code, map[string]any{"code": code, "message": err.Error()})
}

// handleHealth answers 204, and is the one route left unauthenticated: the
// control plane probes it to decide the sandbox is up, before it holds a token.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleEnvs(w http.ResponseWriter, r *http.Request) {
	out := map[string]string{}
	for _, kv := range os.Environ() {
		if k, v, ok := strings.Cut(kv, "="); ok {
			out[k] = v
		}
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	m := metricsPayload{TS: time.Now().Unix(), CPUCount: runtime.NumCPU()}

	if pcts, err := cpu.Percent(0, false); err == nil && len(pcts) > 0 {
		m.CPUUsedPct = pcts[0]
	}
	if vm, err := mem.VirtualMemory(); err == nil {
		m.MemTotal = vm.Total
		m.MemUsed = vm.Used
		m.MemCache = vm.Cached
		m.MemTotalMiB = vm.Total / (1 << 20)
		m.MemUsedMiB = vm.Used / (1 << 20)
	}
	var st syscall.Statfs_t
	if err := syscall.Statfs(s.users.fallback.HomeDir, &st); err == nil {
		bsize := uint64(st.Bsize)
		m.DiskTotal = st.Blocks * bsize
		m.DiskUsed = (st.Blocks - st.Bavail) * bsize
	}

	writeJSON(w, http.StatusOK, m)
}

func (s *Server) handleFiles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.downloadFile(w, r)
	case http.MethodPost:
		s.uploadFile(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(w, http.StatusMethodNotAllowed, errMethodNotAllowed)
	}
}

func (s *Server) resolveRequestPath(r *http.Request) string {
	u := s.users.lookup(r.URL.Query().Get("username"))

	return expandPath(r.URL.Query().Get("path"), u)
}

func (s *Server) downloadFile(w http.ResponseWriter, r *http.Request) {
	path := s.resolveRequestPath(r)

	fi, err := os.Stat(path)
	if err != nil {
		writeError(w, http.StatusNotFound, err)

		return
	}
	if fi.IsDir() {
		writeError(w, http.StatusBadRequest, errIsADirectory)

		return
	}

	f, err := os.Open(path)
	if err != nil {
		writeError(w, http.StatusNotFound, err)

		return
	}
	defer f.Close()

	w.Header().Set("Content-Disposition", "inline; filename="+filepath.Base(path))
	// ServeContent so Range and conditional requests keep working.
	http.ServeContent(w, r, filepath.Base(path), fi.ModTime(), f)
}

func (s *Server) uploadFile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	contentType := r.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New("unsupported content type"))

		return
	}

	switch {
	case mediaType == "application/octet-stream":
		if r.URL.Query().Get("path") == "" {
			writeError(w, http.StatusBadRequest, errors.New("path query parameter is required"))

			return
		}
		entry, err := writeUpload(s.resolveRequestPath(r), r.Body)
		if err != nil {
			writeUploadError(w, err)

			return
		}
		writeJSON(w, http.StatusOK, []restEntryInfo{*entry})
	case strings.HasPrefix(mediaType, "multipart/"):
		s.uploadMultipart(w, r, params["boundary"])
	default:
		writeError(w, http.StatusBadRequest, errors.New("unsupported content type"))
	}
}

func (s *Server) uploadMultipart(w http.ResponseWriter, r *http.Request, boundary string) {
	if boundary == "" {
		writeError(w, http.StatusBadRequest, errors.New("multipart boundary is missing"))

		return
	}

	mr, err := r.MultipartReader()
	if err != nil {
		writeError(w, http.StatusBadRequest, err)

		return
	}

	u := s.users.lookup(r.URL.Query().Get("username"))
	base := r.URL.Query().Get("path")

	var out []restEntryInfo
	for {
		part, err := mr.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			writeError(w, http.StatusBadRequest, err)

			return
		}
		if part.FormName() != "file" {
			part.Close()

			continue
		}

		target := base
		if target == "" {
			target = part.FileName()
		}
		entry, err := writeUpload(expandPath(target, u), part)
		part.Close()
		if err != nil {
			writeUploadError(w, err)

			return
		}
		out = append(out, *entry)
	}

	writeJSON(w, http.StatusOK, out)
}

func writeUpload(path string, body io.Reader) (*restEntryInfo, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := io.Copy(f, body); err != nil {
		return nil, err
	}

	return &restEntryInfo{Path: path, Name: filepath.Base(path), Type: "file"}, nil
}

// writeUploadError maps a full disk to 507 so a client can tell "no space" from
// a generic failure; everything else is the caller's problem or ours.
func writeUploadError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, syscall.ENOSPC), errors.Is(err, syscall.EDQUOT):
		writeError(w, http.StatusInsufficientStorage, err)
	// OpenFile already reports a directory target; no pre-flight stat needed.
	case errors.Is(err, syscall.EISDIR), errors.Is(err, errIsADirectory):
		writeError(w, http.StatusBadRequest, err)
	default:
		writeError(w, http.StatusInternalServerError, err)
	}
}
