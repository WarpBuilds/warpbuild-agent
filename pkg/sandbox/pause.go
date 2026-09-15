//go:build darwin

package sandbox

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

const unmountTimeout = 20 * time.Second

type pausePrepareResult struct {
	// Unmounted reports the data volume is gone. Past this point the guest has
	// no $HOME and cannot start a process, so the sandbox is no longer servable.
	Unmounted bool   `json:"unmounted"`
	HadVolume bool   `json:"had_volume"`
	Detail    string `json:"detail,omitempty"`
}

// handlePausePrepare flushes the guest and releases the data volume so the host
// can snapshot it. Two-phase on purpose: writers are stopped before the flush,
// because a process that writes between the sync and the unmount would otherwise
// lose those writes with nothing to show for it.
func (s *Server) handlePausePrepare(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		writeError(w, http.StatusMethodNotAllowed, errMethodNotAllowed)

		return
	}

	s.pauseMu.Lock()
	defer s.pauseMu.Unlock()

	if s.opts.DataVolume == "" {
		writeJSON(w, http.StatusOK, pausePrepareResult{Detail: "no data volume configured"})

		return
	}

	mounted, err := isMountPoint(s.opts.DataVolume)
	if err != nil || !mounted {
		// Idempotent: a second call, or a sandbox that never had a volume, is a
		// success with nothing to do.
		writeJSON(w, http.StatusOK, pausePrepareResult{Detail: "data volume is not mounted"})

		return
	}

	stopped := s.stopWriters()
	syscall.Sync()

	// Our own cwd must not be the thing holding the volume open.
	if err := os.Chdir("/"); err != nil {
		s.resumeWriters(stopped)
		writeJSON(w, http.StatusOK, pausePrepareResult{
			HadVolume: true,
			Detail:    "could not leave the data volume: " + err.Error(),
		})

		return
	}

	if err := unmount(r.Context(), s.opts.DataVolume); err != nil {
		// Deliberately no force unmount: forcing here is what silently discards
		// the writes this endpoint exists to preserve. Let the writers continue
		// and report failure; the caller suspends anyway.
		s.resumeWriters(stopped)
		writeJSON(w, http.StatusOK, pausePrepareResult{
			HadVolume: true,
			Detail:    "unmount failed: " + err.Error(),
		})

		return
	}

	syscall.Sync()
	writeJSON(w, http.StatusOK, pausePrepareResult{Unmounted: true, HadVolume: true})
}

// stopWriters SIGSTOPs every process group the agent spawned and returns the
// pids that were actually signalled, so they can be resumed if the pause aborts.
func (s *Server) stopWriters() []int {
	var stopped []int
	for _, h := range s.procs.running() {
		if h.pid <= 0 {
			continue
		}
		if err := syscall.Kill(-h.pid, syscall.SIGSTOP); err == nil {
			stopped = append(stopped, h.pid)
		}
	}

	return stopped
}

func (s *Server) resumeWriters(pids []int) {
	for _, pid := range pids {
		_ = syscall.Kill(-pid, syscall.SIGCONT)
	}
}

func unmount(ctx context.Context, path string) error {
	ctx, cancel := context.WithTimeout(ctx, unmountTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sudo", "diskutil", "unmount", path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}

	return nil
}
