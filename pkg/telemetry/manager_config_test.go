package telemetry

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

func TestMain(m *testing.M) {
	if _, err := log.Init(&log.InitOptions{LogLevel: "error"}); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	require.NoError(t, err, "getwd")
	return filepath.Dir(filepath.Dir(wd))
}

func newTestManager(t *testing.T, export *exportConfig) *TelemetryManager {
	t.Helper()
	tm := NewTelemetryManager(t.Context(), 33931, repoRoot(t), nil, "wr_test_runner", "secret", "https://api.warpbuild.com/api/v1", false, "", "")
	tm.exportCfg = export
	return tm
}

func testExportConfig() *exportConfig {
	return &exportConfig{
		MetricsEndpoint: "https://otlp.datadoghq.com/v1/metrics",
		LogsEndpoint:    "https://http-intake.logs.datadoghq.com/v1/logs",
		Headers: map[string]string{
			"dd-api-key":            "super-secret-value",
			"dd-otel-metric-config": `{"resource_attributes_as_tags": true}`,
		},
		ResourceAttrs: map[string]string{
			"service.name":                 "ci",
			"host.name":                    "warp-ubuntu-latest-x64-4x",
			"warpbuild.runner.instance_id": "wr_test_runner",
			"vcs.repository.name":          "widgets",
			"warpbuild.note":               `he said "hi"\and left`,
		},
	}
}

func renderConfig(t *testing.T, tm *TelemetryManager) string {
	t.Helper()

	base := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(base, "pkg", "telemetry"), 0o755))
	for _, name := range []string{"otel-collector-config.tmpl", "binaries"} {
		src := filepath.Join(repoRoot(t), "pkg", "telemetry", name)
		dst := filepath.Join(base, "pkg", "telemetry", name)
		require.NoErrorf(t, os.Symlink(src, dst), "symlink %s", name)
	}
	tm.baseDirectory = base

	require.NoError(t, tm.writeOtelCollectorConfig())
	out, err := os.ReadFile(tm.getConfigFilePath())
	require.NoError(t, err, "read rendered config")
	return string(out)
}

func validateWithCollector(t *testing.T, tm *TelemetryManager) {
	t.Helper()

	binary, err := tm.getOtelCollectorPath()
	if err != nil {
		t.Skipf("collector binary unavailable: %v", err)
	}

	cmd := exec.Command(binary, "validate", "--config", tm.getConfigFilePath())
	cmd.Env = append(os.Environ(), tm.currentExportConfig().envPairs()...)
	out, runErr := cmd.CombinedOutput()
	require.NoErrorf(t, runErr, "collector rejected config:\n%s", out)
}

func TestRenderConfig_WithoutExport(t *testing.T) {
	tm := newTestManager(t, nil)
	got := renderConfig(t, tm)

	for _, absent := range []string{"otlphttp/customer", "resource/customer", "connectors:", "forward/customer"} {
		assert.NotContainsf(t, got, absent, "unconfigured export should not emit %q", absent)
	}
	for _, present := range []string{"otlphttp:", "otlphttp/gha_logs:", "hostmetrics:"} {
		assert.Containsf(t, got, present, "expected %q in rendered config", present)
	}
	validateWithCollector(t, tm)
}

func TestRenderConfig_WithExport(t *testing.T) {
	tm := newTestManager(t, testExportConfig())
	got := renderConfig(t, tm)

	for _, present := range []string{
		"otlphttp/customer_metrics:",
		"otlphttp/customer_logs:",
		"resource/customer:",
		"forward/customer_metrics:",
		"forward/customer_logs:",
		"metrics/customer:",
		"logs/customer:",
		`metrics_endpoint: "https://otlp.datadoghq.com/v1/metrics"`,
	} {
		assert.Containsf(t, got, present, "expected %q in rendered config:\n%s", present, got)
	}
	validateWithCollector(t, tm)
}

func TestRenderConfig_SecretNotInFile(t *testing.T) {
	tm := newTestManager(t, testExportConfig())
	got := renderConfig(t, tm)

	require.NotContainsf(t, got, "super-secret-value", "credential leaked into rendered config:\n%s", got)
	assert.Containsf(t, got, "${env:"+exportHeaderEnvPrefix, "expected header to be an env reference:\n%s", got)

	assert.Contains(t, tm.currentExportConfig().envPairs(), exportHeaderEnvPrefix+"0=super-secret-value",
		"credential missing from collector environment")
}

func TestRenderConfig_InternalPipelinesUntouched(t *testing.T) {
	tm := newTestManager(t, testExportConfig())
	got := renderConfig(t, tm)

	for _, present := range []string{
		"exporters: [otlphttp, forward/customer_metrics]",
		"exporters: [otlphttp, forward/customer_logs]",
		"exporters: [otlphttp/gha_logs, forward/customer_logs]",
	} {
		assert.Containsf(t, got, present, "expected %q in rendered config:\n%s", present, got)
	}

	internalMetrics := section(got, "    metrics:\n", "\n    ")
	assert.NotContainsf(t, internalMetrics, "resource/customer",
		"resource/customer leaked into the internal metrics pipeline:\n%s", internalMetrics)
}

func TestRenderConfig_SignalsExportIndependently(t *testing.T) {
	tm := newTestManager(t, testExportConfig())
	got := renderConfig(t, tm)

	for _, present := range []string{
		"exporters: [otlphttp/customer_metrics]",
		"exporters: [otlphttp/customer_logs]",
	} {
		assert.Containsf(t, got, present, "expected %q in rendered config:\n%s", present, got)
	}
	validateWithCollector(t, tm)
}

func TestRenderConfig_MetricsOnly(t *testing.T) {
	cfg := testExportConfig()
	cfg.LogsEndpoint = ""
	tm := newTestManager(t, cfg)
	got := renderConfig(t, tm)

	for _, present := range []string{"otlphttp/customer_metrics:", "forward/customer_metrics:", "metrics/customer:"} {
		assert.Containsf(t, got, present, "expected %q", present)
	}
	for _, absent := range []string{"otlphttp/customer_logs", "forward/customer_logs", "logs/customer:"} {
		assert.NotContainsf(t, got, absent, "did not expect %q", absent)
	}
	validateWithCollector(t, tm)
}

func TestRenderConfig_LogsOnly(t *testing.T) {
	cfg := testExportConfig()
	cfg.MetricsEndpoint = ""
	tm := newTestManager(t, cfg)
	got := renderConfig(t, tm)

	for _, present := range []string{"otlphttp/customer_logs:", "forward/customer_logs:", "logs/customer:"} {
		assert.Containsf(t, got, present, "expected %q", present)
	}
	for _, absent := range []string{"otlphttp/customer_metrics", "forward/customer_metrics", "metrics/customer:"} {
		assert.NotContainsf(t, got, absent, "did not expect %q", absent)
	}
	validateWithCollector(t, tm)
}

func section(s, start, sep string) string {
	i := strings.Index(s, start)
	if i < 0 {
		return ""
	}
	rest := s[i+len(start):]
	if j := strings.Index(rest, sep); j >= 0 {
		return rest[:j]
	}
	return rest
}

func newCollectOffManager(t *testing.T, export *exportConfig) *TelemetryManager {
	t.Helper()
	tm := newTestManager(t, export)
	tm.collect = false
	return tm
}

// The exporter lists are the whole of what collection-off changes, so pin them
// verbatim for the states that were already shipped.
func TestRenderConfig_CollectOnExporterListsUnchanged(t *testing.T) {
	t.Run("no export", func(t *testing.T) {
		got := renderConfig(t, newTestManager(t, nil))
		for _, line := range []string{
			"      exporters: [otlphttp]",
			"      exporters: [otlphttp/gha_logs]",
		} {
			assert.Containsf(t, got, line, "expected %q", line)
		}
	})

	t.Run("with export", func(t *testing.T) {
		got := renderConfig(t, newTestManager(t, testExportConfig()))
		for _, line := range []string{
			"      exporters: [otlphttp, forward/customer_logs]",
			"      exporters: [otlphttp/gha_logs, forward/customer_logs]",
			"      exporters: [otlphttp, forward/customer_metrics]",
		} {
			assert.Containsf(t, got, line, "expected %q", line)
		}
	})
}

func TestRenderConfig_CollectOffExportsOnly(t *testing.T) {
	tm := newCollectOffManager(t, testExportConfig())
	got := renderConfig(t, tm)

	for _, present := range []string{
		"      exporters: [forward/customer_logs]",
		"      exporters: [forward/customer_metrics]",
		"    logs/customer:",
		"    metrics/customer:",
	} {
		assert.Containsf(t, got, present, "expected %q:\n%s", present, got)
	}
	for _, absent := range []string{"[otlphttp]", "[otlphttp,", "[otlphttp/gha_logs"} {
		assert.NotContainsf(t, got, absent, "our own exporter %q must be off every pipeline", absent)
	}
	validateWithCollector(t, tm)
}

// A pipeline left with an empty exporter list is a fatal config error, so the
// single-signal exports are where collection-off would take the collector down.
func TestRenderConfig_CollectOffSingleSignal(t *testing.T) {
	t.Run("metrics only", func(t *testing.T) {
		export := testExportConfig()
		export.LogsEndpoint = ""
		tm := newCollectOffManager(t, export)
		got := renderConfig(t, tm)

		assert.Contains(t, got, "      exporters: [forward/customer_metrics]")
		assert.NotContains(t, got, "      processors: [batch/logs]", "log pipelines have no destination left")
		assert.NotContains(t, got, "forward/customer_logs")
		validateWithCollector(t, tm)
	})

	t.Run("logs only", func(t *testing.T) {
		export := testExportConfig()
		export.MetricsEndpoint = ""
		tm := newCollectOffManager(t, export)
		got := renderConfig(t, tm)

		assert.Contains(t, got, "      exporters: [forward/customer_logs]")
		assert.NotContains(t, got, "      receivers: [hostmetrics", "the metrics pipeline has no destination left")
		assert.NotContains(t, got, "forward/customer_metrics")
		validateWithCollector(t, tm)
	})
}

func TestRenderConfig_Parked(t *testing.T) {
	tm := newCollectOffManager(t, nil)
	tm.parked = true
	got := renderConfig(t, tm)

	for _, present := range []string{"  nop: {}", "    metrics/idle:", "      receivers: [nop]", "      exporters: [nop]"} {
		assert.Containsf(t, got, present, "expected %q:\n%s", present, got)
	}
	for _, absent := range []string{"      processors: [batch/logs]", "      receivers: [hostmetrics"} {
		assert.NotContainsf(t, got, absent, "a parked collector runs no real pipeline (%q)", absent)
	}
	validateWithCollector(t, tm)
}
