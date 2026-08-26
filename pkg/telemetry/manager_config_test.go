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

// The package logs through a process-global logger that the daemon
// initialises at startup; tests have to stand it up themselves.
func TestMain(m *testing.M) {
	if _, err := log.Init(&log.InitOptions{LogLevel: "error"}); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

// repoRoot resolves the module root, which is what baseDirectory points
// at on a runner — the template and collector binaries hang off it.
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
		Endpoint: "https://otlp.datadoghq.com",
		Headers: map[string]string{
			"dd-api-key":            "super-secret-value",
			"dd-otel-metric-config": `{"resource_attributes_as_tags": true}`,
		},
		ResourceAttrs: map[string]string{
			"service.name":                 "ci",
			"host.name":                    "warp-ubuntu-latest-x64-4x",
			"warpbuild.runner.instance_id": "wr_test_runner",
			"vcs.repository.name":          "widgets",
			// Deliberately awkward: must survive YAML quoting.
			"warpbuild.note": `he said "hi"\and left`,
		},
	}
}

// renderConfig writes the collector config into a temp base dir that
// symlinks the real template, and returns the rendered text.
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

// validateWithCollector runs the bundled collector's own config
// validator. Skips where the binary for this platform isn't present.
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
	// Our own pipelines must still be intact.
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
		`endpoint: "https://otlp.datadoghq.com"`,
	} {
		assert.Containsf(t, got, present, "expected %q in rendered config:\n%s", present, got)
	}
	validateWithCollector(t, tm)
}

// The credential must reach the collector through the environment, never
// the config file — that file is long-lived and ends up in bug reports.
func TestRenderConfig_SecretNotInFile(t *testing.T) {
	tm := newTestManager(t, testExportConfig())
	got := renderConfig(t, tm)

	require.NotContainsf(t, got, "super-secret-value", "credential leaked into rendered config:\n%s", got)
	assert.Containsf(t, got, "${env:"+exportHeaderEnvPrefix, "expected header to be an env reference:\n%s", got)

	assert.Contains(t, tm.currentExportConfig().envPairs(), exportHeaderEnvPrefix+"0=super-secret-value",
		"credential missing from collector environment")
}

// Our own pipelines must keep exporting to the local receiver even when
// an export is configured — the customer path is additive.
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

	// resource/customer must never appear in our own metrics pipeline:
	// it rewrites host.name, which our ClickHouse view keys on.
	internalMetrics := section(got, "    metrics:\n", "\n    ")
	assert.NotContainsf(t, internalMetrics, "resource/customer",
		"resource/customer leaked into the internal metrics pipeline:\n%s", internalMetrics)
}

// Each signal gets its own exporter instance, so a destination that
// rejects one keeps accepting the other rather than sharing a queue.
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

// section returns the slice of s starting at start up to the next
// occurrence of sep, for coarse pipeline assertions.
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
