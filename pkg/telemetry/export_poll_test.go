package telemetry

import (
	"testing"

	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

func allocationDetails(telemetryEnabled *bool, export *warpbuild.CommonsObservabilityExportConfig) *warpbuild.CommonsRunnerInstanceAllocationDetails {
	out := warpbuild.NewCommonsRunnerInstanceAllocationDetails()
	if telemetryEnabled != nil {
		out.SetTelemetryEnabled(*telemetryEnabled)
	}
	if export != nil {
		out.ObservabilityExport = export
	}
	return out
}

// The backend only populates observability_export on the freshly-ALLOCATED
// response path. Every later poll -- the whole time the job actually runs --
// omits it. Reading that as "export removed" tore the customer exporter out of
// the collector seconds into every job while the internal pipeline, which is
// rendered unconditionally, carried on none the wiser.
func TestPollWithoutExportBlockKeepsExisting(t *testing.T) {
	tm := &TelemetryManager{}
	enabled := true

	if !tm.applyAllocationDetails(allocationDetails(&enabled, apiExport("https://otlp.example.com", []string{"metrics"}))) {
		t.Fatal("telemetry should stay enabled")
	}

	applied, fingerprint := tm.exportCfg, tm.exportFingerprint
	if applied == nil {
		t.Fatal("the allocated response should have applied an export config")
	}

	// Every subsequent poll for the life of the job looks like this.
	for i := 0; i < 3; i++ {
		if !tm.applyAllocationDetails(allocationDetails(&enabled, nil)) {
			t.Fatal("telemetry should stay enabled")
		}
	}

	if tm.exportCfg == nil {
		t.Fatal("export config was cleared by a poll that simply did not carry one")
	}
	if tm.exportCfg != applied || tm.exportFingerprint != fingerprint {
		t.Errorf("export config changed: %+v (fingerprint %s), want %+v (fingerprint %s)",
			tm.exportCfg, tm.exportFingerprint, applied, fingerprint)
	}
}

func TestPollAppliesChangedExport(t *testing.T) {
	tm := &TelemetryManager{}
	enabled := true

	tm.applyAllocationDetails(allocationDetails(&enabled, apiExport("https://one.example.com", []string{"metrics"})))
	first := tm.exportFingerprint

	tm.applyAllocationDetails(allocationDetails(&enabled, apiExport("https://two.example.com", []string{"metrics", "logs"})))

	if tm.exportFingerprint == first {
		t.Error("a delivered export block with new values must still be applied")
	}
	if tm.exportCfg == nil || tm.exportCfg.Endpoint != "https://two.example.com" {
		t.Errorf("expected the second endpoint, got %+v", tm.exportCfg)
	}
}

// telemetry_enabled is populated on every response path, so it stays the kill
// switch even though a missing export block no longer is one.
func TestPollStopsWhenTelemetryDisabled(t *testing.T) {
	tm := &TelemetryManager{}
	disabled := false

	if tm.applyAllocationDetails(allocationDetails(&disabled, nil)) {
		t.Error("telemetry_enabled=false must stop telemetry")
	}
}

func TestPollWithAbsentTelemetryFlagKeepsRunning(t *testing.T) {
	tm := &TelemetryManager{}

	if !tm.applyAllocationDetails(allocationDetails(nil, nil)) {
		t.Error("an absent telemetry_enabled key defaults to enabled")
	}
}
