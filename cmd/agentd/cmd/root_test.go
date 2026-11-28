package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderTelemetryTemplateFlags(t *testing.T) {
	dir := t.TempDir()
	vars := filepath.Join(dir, "vars.json")
	out := filepath.Join(dir, "config.yaml")
	require.NoError(t, os.WriteFile(vars, []byte(`{
	  "os": "darwin", "arch": "arm64", "port": 33931, "runner_id": "wr_1", "collect": true,
	  "export": {"metrics_endpoint": "https://otlp.example.com/v1/metrics", "headers": {"api-key": "do-not-write-me"}},
	  "data": {"export": {"metrics": [{"name": "system.cpu.utilization", "attributes": ["aggregation"]}]}}
	}`), 0o644))

	rootCmd.SetArgs([]string{
		"--render-telemetry-template", filepath.Join("..", "..", "..", "pkg", "telemetry", "otel-collector-config.tmpl"),
		"--render-telemetry-vars", vars,
		"--render-telemetry-out", out,
	})
	require.NoError(t, rootCmd.Execute())

	got, err := os.ReadFile(out)
	require.NoError(t, err)
	for _, s := range []string{
		"include: ['/var/log/system.log']",
		`host.name`,
		`value: "wr_1"`,
		`metrics_endpoint: "http://localhost:33931/v1/metrics"`,
		`metrics_endpoint: "https://otlp.example.com/v1/metrics"`,
		"exporters: [otlphttp, forward/customer_metrics]",
		`- "keep_keys(attributes, [\"aggregation\"]) where metric.name == \"system.cpu.utilization\""`,
	} {
		assert.Contains(t, string(got), s)
	}
	assert.NotContains(t, string(got), "do-not-write-me")
}
