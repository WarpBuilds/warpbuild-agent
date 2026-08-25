package telemetry

import (
	"strings"
	"testing"

	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

func apiExport(endpoint string, signals []string) *warpbuild.CommonsObservabilityExportConfig {
	out := warpbuild.NewCommonsObservabilityExportConfig()
	if endpoint != "" {
		out.SetEndpoint(endpoint)
	}
	if signals != nil {
		out.SetSignals(signals)
	}
	out.SetHeaders(map[string]string{"dd-api-key": "secret"})
	out.SetResourceAttrs(map[string]string{"service.name": "ci"})
	return out
}

func TestExportConfigFrom(t *testing.T) {
	cases := []struct {
		name     string
		in       *warpbuild.CommonsObservabilityExportConfig
		wantNil  bool
		wantMetr bool
		wantLogs bool
	}{
		{name: "nil payload", in: nil, wantNil: true},
		{name: "no endpoint", in: apiExport("", []string{"metrics"}), wantNil: true},
		{name: "no signals", in: apiExport("https://x.example.com", nil), wantNil: true},
		{name: "unknown signal only", in: apiExport("https://x.example.com", []string{"traces"}), wantNil: true},
		{name: "metrics only", in: apiExport("https://x.example.com", []string{"metrics"}), wantMetr: true},
		{name: "logs only", in: apiExport("https://x.example.com", []string{"logs"}), wantLogs: true},
		{name: "both", in: apiExport("https://x.example.com", []string{"metrics", "logs"}), wantMetr: true, wantLogs: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := exportConfigFrom(tc.in)
			if tc.wantNil {
				if got != nil {
					t.Fatalf("expected nil, got %+v", got)
				}
				return
			}
			if got == nil {
				t.Fatal("expected a config, got nil")
			}
			if got.Metrics != tc.wantMetr || got.Logs != tc.wantLogs {
				t.Errorf("signals: metrics=%t logs=%t, want metrics=%t logs=%t",
					got.Metrics, got.Logs, tc.wantMetr, tc.wantLogs)
			}
		})
	}
}

func TestExportConfigFingerprint(t *testing.T) {
	var nilCfg *exportConfig
	if nilCfg.fingerprint() != "" {
		t.Error("nil config should fingerprint empty, so an unconfigured runner never restarts")
	}

	base := testExportConfig()
	if base.fingerprint() != testExportConfig().fingerprint() {
		t.Error("identical configs must fingerprint the same, or the collector restarts on every poll")
	}

	// A rotated credential has to trigger a restart.
	rotated := testExportConfig()
	rotated.Headers["dd-api-key"] = "rotated"
	if rotated.fingerprint() == base.fingerprint() {
		t.Error("credential change must change the fingerprint")
	}

	// So does a new job's attributes.
	reattr := testExportConfig()
	reattr.ResourceAttrs["cicd.pipeline.task.run.id"] = "999"
	if reattr.fingerprint() == base.fingerprint() {
		t.Error("resource attribute change must change the fingerprint")
	}

	// The fingerprint is logged, so it must not carry the secret.
	if strings.Contains(base.fingerprint(), "super-secret-value") {
		t.Error("fingerprint must not embed the credential")
	}
}

func TestExportConfigHeaderEnv(t *testing.T) {
	cfg := testExportConfig()
	env := cfg.headerEnv()

	if len(env) != len(cfg.Headers) {
		t.Fatalf("expected %d env mappings, got %d", len(cfg.Headers), len(env))
	}
	// Sorted-position indexing keeps the rendered config stable across
	// restarts that change nothing.
	if env["dd-api-key"] != exportHeaderEnvPrefix+"0" {
		t.Errorf("dd-api-key mapped to %q", env["dd-api-key"])
	}
	if env["dd-otel-metric-config"] != exportHeaderEnvPrefix+"1" {
		t.Errorf("dd-otel-metric-config mapped to %q", env["dd-otel-metric-config"])
	}

	pairs := cfg.envPairs()
	if len(pairs) != len(cfg.Headers) {
		t.Fatalf("expected %d env pairs, got %d", len(cfg.Headers), len(pairs))
	}
	for _, want := range []string{
		exportHeaderEnvPrefix + "0=super-secret-value",
		exportHeaderEnvPrefix + `1={"resource_attributes_as_tags": true}`,
	} {
		found := false
		for _, p := range pairs {
			if p == want {
				found = true
			}
		}
		if !found {
			t.Errorf("missing env pair %q in %v", want, pairs)
		}
	}
}

func TestExportConfigNilSafety(t *testing.T) {
	var cfg *exportConfig
	if cfg.headerEnv() != nil {
		t.Error("headerEnv on nil should be nil")
	}
	// runOtelCollector appends this unconditionally.
	if pairs := cfg.envPairs(); pairs != nil {
		t.Errorf("envPairs on nil should be nil, got %v", pairs)
	}
}
