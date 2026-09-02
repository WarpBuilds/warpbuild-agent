package telemetry

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

func allocationDetails(status string, telemetryEnabled *bool, export *warpbuild.CommonsTelemetryExportConfig) *warpbuild.CommonsRunnerInstanceAllocationDetails {
	out := warpbuild.NewCommonsRunnerInstanceAllocationDetails()
	out.SetStatus(status)
	if telemetryEnabled != nil {
		out.SetTelemetryEnabled(*telemetryEnabled)
	}
	if export != nil {
		out.TelemetryExport = export
	}
	return out
}

func enabled() *bool { b := true; return &b }

func TestPollWaitsWhileUnassigned(t *testing.T) {
	tm := &TelemetryManager{}

	for i := range 3 {
		got := tm.applyAllocationDetails(allocationDetails(allocationStatusUnassigned, enabled(), nil))
		assert.Equalf(t, pollWait, got, "poll %d", i)
		assert.Nil(t, tm.exportCfg, "an unassigned response must not configure an export")
	}
}

func TestPollAppliesExportOnceAllocated(t *testing.T) {
	tm := &TelemetryManager{}

	got := tm.applyAllocationDetails(allocationDetails("assigned", enabled(), apiExport("https://otlp.example.com/v1/metrics")))

	require.Equal(t, pollApplied, got)
	require.NotNil(t, tm.exportCfg)
	assert.Equal(t, "https://otlp.example.com/v1/metrics", tm.exportCfg.MetricsEndpoint)
}

func TestPollStopsWhenAllocatedWithoutExport(t *testing.T) {
	tm := &TelemetryManager{}

	got := tm.applyAllocationDetails(allocationDetails("assigned", enabled(), nil))

	assert.Equal(t, pollNoExport, got)
	assert.Nil(t, tm.exportCfg)
}

func TestPollParksWhileUnassignedWithCollectionOff(t *testing.T) {
	tm := &TelemetryManager{collect: true}
	disabled := false

	for i := range 3 {
		got := tm.applyAllocationDetails(allocationDetails(allocationStatusUnassigned, &disabled, nil))
		assert.Equalf(t, pollWait, got, "poll %d must keep waiting for the export destination", i)
	}

	assert.True(t, tm.parked)
	assert.False(t, tm.collect)
}

func TestPollStopsWhenAllocatedWithCollectionOffAndNoExport(t *testing.T) {
	tm := &TelemetryManager{collect: true}
	disabled := false

	got := tm.applyAllocationDetails(allocationDetails("assigned", &disabled, nil))

	assert.Equal(t, pollDisabled, got)
}

func TestPollExportsWithCollectionOff(t *testing.T) {
	tm := &TelemetryManager{collect: true}
	disabled := false

	got := tm.applyAllocationDetails(
		allocationDetails("assigned", &disabled, apiExport("https://otlp.example.com/v1/metrics")))

	require.Equal(t, pollApplied, got)
	require.NotNil(t, tm.exportCfg)
	assert.False(t, tm.collect, "our own exporters must be dropped")
	assert.False(t, tm.parked)
}

func TestPollKeepsCollectingWhenFlagAbsent(t *testing.T) {
	tm := &TelemetryManager{collect: true}

	got := tm.applyAllocationDetails(
		allocationDetails("assigned", nil, apiExport("https://otlp.example.com/v1/metrics")))

	require.Equal(t, pollApplied, got)
	assert.True(t, tm.collect, "an old backend omits the flag; that must not stop collection")
}

func TestPollWithAbsentTelemetryFlagKeepsRunning(t *testing.T) {
	tm := &TelemetryManager{}

	got := tm.applyAllocationDetails(allocationDetails(allocationStatusUnassigned, nil, nil))

	assert.Equal(t, pollWait, got)
}

func TestPollWithNilDetailsWaits(t *testing.T) {
	tm := &TelemetryManager{}

	assert.Equal(t, pollWait, tm.applyAllocationDetails(nil))
}

func TestPollResumesCollectionAfterParking(t *testing.T) {
	tm := &TelemetryManager{collect: true}
	disabled, reenabled := false, true

	require.Equal(t, pollWait,
		tm.applyAllocationDetails(allocationDetails(allocationStatusUnassigned, &disabled, nil)))
	require.True(t, tm.parked)

	got := tm.applyAllocationDetails(allocationDetails("assigned", &reenabled, nil))

	assert.Equal(t, pollNoExport, got)
	assert.False(t, tm.parked, "a re-enabled runner must not sit idle for the whole job")
	assert.True(t, tm.collect)
}

// backend-core main has no export feature at all: its allocation details carry
// telemetry_enabled and no telemetry_export key whatsoever. Decode that exact
// wire shape rather than constructing the struct, so the generated client's
// absent-vs-null handling is part of what is under test.
func mainBackendWire(t *testing.T, status string, telemetryEnabled bool) *warpbuild.CommonsRunnerInstanceAllocationDetails {
	t.Helper()
	raw := fmt.Sprintf(`{"status":%q,"runner_application":"github","telemetry_enabled":%t,
	  "gh_runner_application_details":{"runner_name":"wr_1","labels":["warp-ubuntu-latest-x64-2x"]}}`,
		status, telemetryEnabled)
	out := warpbuild.NewCommonsRunnerInstanceAllocationDetails()
	require.NoError(t, json.Unmarshal([]byte(raw), out))
	require.Nil(t, out.TelemetryExport, "main never sends telemetry_export")
	return out
}

func TestBackwardCompat_MainBackend_TelemetryOn(t *testing.T) {
	tm := &TelemetryManager{collect: true}

	assert.Equal(t, pollWait, tm.applyAllocationDetails(mainBackendWire(t, allocationStatusUnassigned, true)))
	assert.False(t, tm.parked, "collection is on; nothing to park")

	assert.Equal(t, pollNoExport, tm.applyAllocationDetails(mainBackendWire(t, "assigned", true)))
	assert.True(t, tm.collect, "our own exporters must stay wired against an old backend")
	assert.False(t, tm.parked)
}

func TestBackwardCompat_MainBackend_TelemetryOff(t *testing.T) {
	tm := &TelemetryManager{collect: true}

	assert.Equal(t, pollWait, tm.applyAllocationDetails(mainBackendWire(t, allocationStatusUnassigned, false)))
	assert.True(t, tm.parked, "collection off with no export: idle until allocation")
	assert.False(t, tm.collect)

	assert.Equal(t, pollDisabled, tm.applyAllocationDetails(mainBackendWire(t, "assigned", false)))
}
