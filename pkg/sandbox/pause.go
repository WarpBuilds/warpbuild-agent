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

const (
	unmountTimeout = 20 * time.Second
	unmountRetries = 3
	unmountBackoff = 2 * time.Second
)

type pausePrepareResult struct {
	Unmounted bool   `json:"unmounted"`
	HadVolume bool   `json:"had_volume"`
	Forced    bool   `json:"forced,omitempty"`
	Detail    string `json:"detail,omitempty"`
}

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
		writeJSON(w, http.StatusOK, pausePrepareResult{Detail: "data volume is not mounted"})

		return
	}

	stopped := s.stopWriters()
	syscall.Sync()

	if err := os.Chdir("/"); err != nil {
		s.resumeWriters(stopped)
		writeJSON(w, http.StatusOK, pausePrepareResult{
			HadVolume: true,
			Detail:    "could not leave the data volume: " + err.Error(),
		})

		return
	}

	forced, err := unmount(r.Context(), s.opts.DataVolume)
	if err != nil {
		s.resumeWriters(stopped)
		writeJSON(w, http.StatusOK, pausePrepareResult{
			HadVolume: true,
			Detail:    "unmount failed: " + err.Error(),
		})

		return
	}

	syscall.Sync()
	writeJSON(w, http.StatusOK, pausePrepareResult{
		Unmounted: true, HadVolume: true, Forced: forced,
	})
}

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

func unmount(ctx context.Context, path string) (bool, error) {
	var last error
	for attempt := range unmountRetries {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return false, ctx.Err()
			case <-time.After(unmountBackoff):
			}
		}
		if err := diskutilUnmount(ctx, path); err == nil {
			return false, nil
		} else {
			last = err
		}
	}

	if err := diskutilUnmount(ctx, path, "force"); err != nil {
		return false, fmt.Errorf("%w (polite unmount: %v)", err, last)
	}

	return true, nil
}

func diskutilUnmount(ctx context.Context, path string, extra ...string) error {
	ctx, cancel := context.WithTimeout(ctx, unmountTimeout)
	defer cancel()

	args := append(append([]string{"diskutil", "unmount"}, extra...), path)
	out, err := exec.CommandContext(ctx, "sudo", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}

	return nil
}
