package telemetry

import (
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
