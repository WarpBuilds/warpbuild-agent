package telemetry

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"
	"time"
)

// waitForFile blocks until path exists, so a test can wait for a child
// to reach a known state instead of sleeping and hoping.
func waitForFile(t *testing.T, path string) {
	t.Helper()

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("child never created %s", path)
}

// startWaitable spawns a child and the wait goroutine that
// terminateCollector expects.
func startWaitable(t *testing.T, name string, args ...string) (*exec.Cmd, chan error) {
	t.Helper()

	cmd := exec.Command(name, args...)
	if err := cmd.Start(); err != nil {
		t.Fatalf("start %s: %v", name, err)
	}
	waitDone := make(chan error, 1)
	go func() { waitDone <- cmd.Wait() }()
	return cmd, waitDone
}

func signalOf(t *testing.T, err error) syscall.Signal {
	t.Helper()

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected *exec.ExitError, got %T (%v)", err, err)
	}
	status, ok := exitErr.Sys().(syscall.WaitStatus)
	if !ok {
		t.Fatalf("expected syscall.WaitStatus, got %T", exitErr.Sys())
	}
	if !status.Signaled() {
		t.Fatalf("process was not signalled: %v", status)
	}
	return status.Signal()
}

// A SIGKILL here would drop up to a full batch interval — the tail of the
// job, which is exactly the part people care about.
func TestTerminateCollector_SignalsRatherThanKills(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("no graceful signal on windows")
	}

	tm := newTestManager(t, nil)
	cmd, waitDone := startWaitable(t, "sleep", "60")

	done := make(chan struct{})
	go func() {
		tm.terminateCollector(cmd, waitDone)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("terminateCollector did not return; it should not have waited out the drain timeout")
	}

	if sig := signalOf(t, <-collectExit(t, cmd)); sig != syscall.SIGTERM {
		t.Errorf("expected SIGTERM, got %v", sig)
	}
}

// A collector that ignores SIGTERM must still be reaped, or a stuck
// child would hold up the whole shutdown.
func TestTerminateCollector_KillsWhenDrainStalls(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("no graceful signal on windows")
	}

	// Ignores TERM, so only SIGKILL ends it. Two subtleties: the loop
	// stops `sh -c` from exec'ing away the trap on a trailing simple
	// command, and the ready file closes the race where the signal lands
	// before the trap is installed.
	ready := filepath.Join(t.TempDir(), "ready")
	cmd, waitDone := startWaitable(t, "sh", "-c",
		"trap '' TERM; touch "+ready+"; while :; do sleep 1; done")
	t.Cleanup(func() { _ = cmd.Process.Kill() })

	waitForFile(t, ready)

	tm := newTestManager(t, nil)
	// The behaviour under test is the fallback, not the real window.
	tm.drainTimeout = 300 * time.Millisecond

	start := time.Now()
	done := make(chan struct{})
	go func() {
		tm.terminateCollector(cmd, waitDone)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("terminateCollector never fell back to kill")
	}

	if elapsed := time.Since(start); elapsed < tm.drainTimeout {
		t.Errorf("killed after %s, before the drain window elapsed", elapsed)
	}
	if sig := signalOf(t, <-collectExit(t, cmd)); sig != syscall.SIGKILL {
		t.Errorf("expected SIGKILL fallback, got %v", sig)
	}
}

// collectExit re-reads the process state after terminateCollector has
// already consumed waitDone.
func collectExit(t *testing.T, cmd *exec.Cmd) chan error {
	t.Helper()

	out := make(chan error, 1)
	if cmd.ProcessState == nil {
		t.Fatal("process state not populated; terminateCollector should have waited")
	}
	out <- &exec.ExitError{ProcessState: cmd.ProcessState}
	return out
}
