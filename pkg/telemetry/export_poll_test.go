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

// Until the runner is allocated, the response carries no job details — so an
// absent export block says nothing yet and discovery keeps waiting.
func TestPollWaitsWhileUnassigned(t *testing.T) {
	tm := &TelemetryManager{}

	for i := range 3 {
		got := tm.applyAllocationDetails(allocationDetails(allocationStatusUnassigned, enabled(), nil))
		assert.Equalf(t, pollWait, got, "poll %d", i)
		assert.Nil(t, tm.exportCfg, "an unassigned response must not configure an export")
	}
}

// Once the job details arrive the answer is final: a destination is applied
// and discovery stops. A later change takes effect on the next run.
func TestPollAppliesExportOnceAllocated(t *testing.T) {
	tm := &TelemetryManager{}

	got := tm.applyAllocationDetails(allocationDetails("assigned", enabled(), apiExport("https://otlp.example.com")))

	require.Equal(t, pollApplied, got)
	require.NotNil(t, tm.exportCfg)
	assert.Equal(t, "https://otlp.example.com", tm.exportCfg.Endpoint)
}

// The job details arriving without an export block is a definitive "this org
// has none" — not something to keep polling for.
func TestPollStopsWhenAllocatedWithoutExport(t *testing.T) {
	tm := &TelemetryManager{}

	got := tm.applyAllocationDetails(allocationDetails("assigned", enabled(), nil))

	assert.Equal(t, pollNoExport, got)
	assert.Nil(t, tm.exportCfg)
}

// telemetry_enabled is populated on every response path, including while the
// runner is still unassigned, so it stops collection outright.
func TestPollStopsWhenTelemetryDisabled(t *testing.T) {
	tm := &TelemetryManager{}
	disabled := false

	got := tm.applyAllocationDetails(allocationDetails(allocationStatusUnassigned, &disabled, nil))

	assert.Equal(t, pollDisabled, got)
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
