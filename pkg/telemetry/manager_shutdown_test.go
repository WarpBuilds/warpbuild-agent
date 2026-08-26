package telemetry

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// waitForFile blocks until path exists, so a test can wait for a child to
// reach a known state instead of sleeping and hoping.
func waitForFile(t *testing.T, path string) {
	t.Helper()

	require.Eventuallyf(t, func() bool {
		_, err := os.Stat(path)
		return err == nil
	}, 5*time.Second, 10*time.Millisecond, "child never created %s", path)
}

// startWaitable spawns a child and the wait goroutine terminateCollector expects.
func startWaitable(t *testing.T, name string, args ...string) (*exec.Cmd, chan error) {
	t.Helper()

	cmd := exec.Command(name, args...)
	require.NoErrorf(t, cmd.Start(), "start %s", name)

	waitDone := make(chan error, 1)
	go func() { waitDone <- cmd.Wait() }()
	return cmd, waitDone
}

// exitSignal reports the signal a finished process died from.
func exitSignal(t *testing.T, cmd *exec.Cmd) syscall.Signal {
	t.Helper()

	require.NotNil(t, cmd.ProcessState, "terminateCollector should have reaped the process")
	status, ok := cmd.ProcessState.Sys().(syscall.WaitStatus)
	require.True(t, ok, "expected syscall.WaitStatus")
	require.True(t, status.Signaled(), "process was not signalled: %v", status)
	return status.Signal()
}

// runTerminate calls terminateCollector and fails if it hangs.
func runTerminate(t *testing.T, tm *TelemetryManager, cmd *exec.Cmd, waitDone chan error) {
	t.Helper()

	done := make(chan struct{})
	go func() {
		tm.terminateCollector(cmd, waitDone)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("terminateCollector did not return")
	}
}

// A SIGKILL here would drop up to a full batch interval — the tail of the
// job, which is exactly the part people care about.
func TestTerminateCollector_SignalsRatherThanKills(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("no graceful signal on windows")
	}

	tm := newTestManager(t, nil)
	cmd, waitDone := startWaitable(t, "sleep", "60")

	runTerminate(t, tm, cmd, waitDone)

	assert.Equal(t, syscall.SIGTERM, exitSignal(t, cmd))
}

// A collector that ignores SIGTERM must still be reaped, or a stuck child
// would hold up the whole shutdown.
func TestTerminateCollector_KillsWhenDrainStalls(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("no graceful signal on windows")
	}

	// Ignores TERM, so only SIGKILL ends it. Two subtleties: the loop stops
	// `sh -c` from exec'ing away the trap on a trailing simple command, and
	// the ready file closes the race where the signal lands before the trap
	// is installed.
	ready := filepath.Join(t.TempDir(), "ready")
	cmd, waitDone := startWaitable(t, "sh", "-c",
		"trap '' TERM; touch "+ready+"; while :; do sleep 1; done")
	t.Cleanup(func() { _ = cmd.Process.Kill() })

	waitForFile(t, ready)

	tm := newTestManager(t, nil)
	// The behaviour under test is the fallback, not the real window.
	tm.drainTimeout = 300 * time.Millisecond

	start := time.Now()
	runTerminate(t, tm, cmd, waitDone)

	assert.GreaterOrEqual(t, time.Since(start), tm.drainTimeout,
		"killed before the drain window elapsed")
	assert.Equal(t, syscall.SIGKILL, exitSignal(t, cmd))
}

// Drain stops the collector rather than cycling it: the job is over, so
// there is nothing to come back up for, and a restart would only churn the
// process and briefly drop our own collection too.
func TestDrainStopsRatherThanRestarts(t *testing.T) {
	tm := newTestManager(t, nil)

	tm.Drain()

	assert.Len(t, tm.drainCh, 1, "drain should be pending")
	assert.Empty(t, tm.restartCh, "drain must not schedule a restart")
}

// The post-end hook can fire more than once; a second drain must not block.
func TestDrainIsIdempotent(t *testing.T) {
	tm := newTestManager(t, nil)

	done := make(chan struct{})
	go func() {
		tm.Drain()
		tm.Drain()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Drain blocked on a full channel")
	}
	assert.Len(t, tm.drainCh, 1)
}
