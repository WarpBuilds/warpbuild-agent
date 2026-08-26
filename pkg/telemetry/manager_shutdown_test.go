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

func waitForFile(t *testing.T, path string) {
	t.Helper()

	require.Eventuallyf(t, func() bool {
		_, err := os.Stat(path)
		return err == nil
	}, 5*time.Second, 10*time.Millisecond, "child never created %s", path)
}

func startWaitable(t *testing.T, name string, args ...string) (*exec.Cmd, chan error) {
	t.Helper()

	cmd := exec.Command(name, args...)
	require.NoErrorf(t, cmd.Start(), "start %s", name)

	waitDone := make(chan error, 1)
	go func() { waitDone <- cmd.Wait() }()
	return cmd, waitDone
}

func exitSignal(t *testing.T, cmd *exec.Cmd) syscall.Signal {
	t.Helper()

	require.NotNil(t, cmd.ProcessState, "terminateCollector should have reaped the process")
	status, ok := cmd.ProcessState.Sys().(syscall.WaitStatus)
	require.True(t, ok, "expected syscall.WaitStatus")
	require.True(t, status.Signaled(), "process was not signalled: %v", status)
	return status.Signal()
}

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

func TestTerminateCollector_SignalsRatherThanKills(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("no graceful signal on windows")
	}

	tm := newTestManager(t, nil)
	cmd, waitDone := startWaitable(t, "sleep", "60")

	runTerminate(t, tm, cmd, waitDone)

	assert.Equal(t, syscall.SIGTERM, exitSignal(t, cmd))
}

func TestTerminateCollector_KillsWhenDrainStalls(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("no graceful signal on windows")
	}

	ready := filepath.Join(t.TempDir(), "ready")
	cmd, waitDone := startWaitable(t, "sh", "-c",
		"trap '' TERM; touch "+ready+"; while :; do sleep 1; done")
	t.Cleanup(func() { _ = cmd.Process.Kill() })

	waitForFile(t, ready)

	tm := newTestManager(t, nil)
	tm.drainTimeout = 300 * time.Millisecond

	start := time.Now()
	runTerminate(t, tm, cmd, waitDone)

	assert.GreaterOrEqual(t, time.Since(start), tm.drainTimeout,
		"killed before the drain window elapsed")
	assert.Equal(t, syscall.SIGKILL, exitSignal(t, cmd))
}

func TestDrainStopsRatherThanRestarts(t *testing.T) {
	tm := newTestManager(t, nil)

	tm.Drain()

	assert.Len(t, tm.drainCh, 1, "drain should be pending")
	assert.Empty(t, tm.restartCh, "drain must not schedule a restart")
}

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
