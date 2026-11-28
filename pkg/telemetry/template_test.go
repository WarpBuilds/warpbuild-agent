package telemetry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func renderText(t *testing.T, text string, in renderInput) string {
	t.Helper()
	out, err := renderCollectorConfig("test", text, newCollectorTemplateData(in))
	require.NoError(t, err)
	return string(out)
}

func TestOttlStr(t *testing.T) {
	cases := map[string]string{
		"host.name":       `"host.name"`,
		`he said "hi"`:    `"he said \"hi\""`,
		`C:\runner`:       `"C:\\runner"`,
		`trailing\`:       `"trailing\\"`,
		"":                `""`,
		`mixed\"both"\\`:  `"mixed\\\"both\"\\\\"`,
		"unicode ✓ stays": `"unicode ✓ stays"`,
	}
	for in, want := range cases {
		assert.Equalf(t, want, ottlStr(in), "ottlStr(%q)", in)
	}
}

func TestRenderTemplate_DataAndHelpers(t *testing.T) {
	const text = `metrics: {{ dig "export" "metrics" "none" .Data }}
logs: {{ dig "export" "logs" "none" .Data }}
{{- $keys := list }}
{{- range (dig "export" "keep" (list) .Data) }}{{ $keys = append $keys (ottlStr .) }}{{ end }}
statement: {{ yamlStr (printf "keep_keys(attributes, [%s])" (join ", " $keys)) }}
{{- range $name := (dig "export" "drop" (list) .Data) }}
drop: {{ yamlStr $name }}
{{- end }}
`
	data := map[string]any{
		"export": map[string]any{
			"metrics": "on",
			"keep":    []any{"host.name", `odd"key\`},
			"drop":    []any{"system.disk.io", "system.network.io"},
		},
	}

	got := renderText(t, text, renderInput{Data: data})

	assert.Equal(t, `metrics: on
logs: none
statement: "keep_keys(attributes, [\"host.name\", \"odd\\\"key\\\\\"])"
drop: "system.disk.io"
drop: "system.network.io"
`, got)
}

func TestRenderTemplate_DataIsNeverNil(t *testing.T) {
	got := renderText(t, `{{ len .Data }} {{ dig "export" "metrics" "unset" .Data }} {{ hasKey .Data "x" }}`, renderInput{})

	assert.Equal(t, "0 unset false", got)
}

func TestRenderTemplate_KeepsExistingFields(t *testing.T) {
	got := renderText(t, `{{.OS}}/{{.Arch}} {{.Port}} {{.RunnerID}} {{.Collect}} {{.Parked}} [{{.MetricsExporters}}] {{.ExportMetrics}} {{.ExportMetricsInterval}}`,
		renderInput{OS: "darwin", Arch: "arm64", Port: 33931, RunnerID: "wr_1", Collect: true, Export: testExportConfig()})

	assert.Equal(t, "darwin/arm64 33931 wr_1 true false [otlphttp, forward/customer_metrics] true 0", got)
}

func TestRenderTemplate_FunctionsAreHermetic(t *testing.T) {
	for _, call := range []string{`env "HOME"`, `expandenv "$HOME"`, `getHostByName "example.com"`, `now`, `randAlpha 4`, `uuidv4`} {
		_, err := renderCollectorConfig("test", "{{ "+call+" }}", newCollectorTemplateData(renderInput{}))
		assert.ErrorContainsf(t, err, "not defined", "%s must not be available to templates", call)
	}
}

func TestRenderTemplate_MutationsStayInTheRender(t *testing.T) {
	data := map[string]any{"nested": map[string]any{"keep": "me"}, "list": []any{"a"}}

	renderText(t, `{{ $_ := set .Data "added" 1 }}{{ $_ := set .Data.nested "keep" "changed" }}{{ $_ := unset .Data "list" }}`,
		renderInput{Data: data})

	assert.Equal(t, map[string]any{"nested": map[string]any{"keep": "me"}, "list": []any{"a"}}, data,
		"the poller compares against this map; a template must not be able to change it")
}

func imageTemplatePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(repoRoot(t), "pkg", "telemetry", "otel-collector-config.tmpl")
}

func renderModeFiles(t *testing.T, vars any) (varsPath, outPath string) {
	t.Helper()
	dir := t.TempDir()
	raw, err := json.Marshal(vars)
	require.NoError(t, err)
	varsPath = filepath.Join(dir, "vars.json")
	require.NoError(t, os.WriteFile(varsPath, raw, 0o644))
	return varsPath, filepath.Join(dir, "config.yaml")
}

func wireExport() map[string]any {
	export := testExportConfig()
	return map[string]any{
		"metrics_endpoint": export.MetricsEndpoint,
		"logs_endpoint":    export.LogsEndpoint,
		"headers":          export.Headers,
		"resource_attrs":   export.ResourceAttrs,
	}
}

func TestRenderMode_MatchesLiveRender(t *testing.T) {
	live := renderConfig(t, newTestManager(t, testExportConfig()))

	varsPath, outPath := renderModeFiles(t, map[string]any{
		"os":        runtime.GOOS,
		"arch":      runtime.GOARCH,
		"port":      33931,
		"runner_id": "wr_test_runner",
		"collect":   true,
		"export":    wireExport(),
	})
	require.NoError(t, RenderTemplateFile(imageTemplatePath(t), varsPath, outPath))

	got, err := os.ReadFile(outPath)
	require.NoError(t, err)
	assert.Equal(t, live, string(got), "render mode must build the template data exactly as telemetryd does")
}

func TestRenderMode_OverridesOS(t *testing.T) {
	cases := map[string][]string{
		"windows": {"windowseventlog/application:", "windowsperfcounters/processor:", `'C:\warpbuilds\runner\_diag\*.log'`},
		"darwin":  {"include: ['/var/log/system.log']", "filter/drop_free_memory:", "metricstransform/mem_aggregate_darwin"},
		"linux":   {"include: ['/var/log/syslog']", "filter/keep_only_used_memory:"},
	}
	for goos, want := range cases {
		t.Run(goos, func(t *testing.T) {
			varsPath, outPath := renderModeFiles(t, map[string]any{
				"os": goos, "arch": "amd64", "port": 33931, "runner_id": "wr_1", "collect": true, "export": wireExport(),
			})
			require.NoError(t, RenderTemplateFile(imageTemplatePath(t), varsPath, outPath))

			got, err := os.ReadFile(outPath)
			require.NoError(t, err)
			for _, s := range append(want, "otlphttp/customer_metrics:", "otlphttp/customer_logs:", "${env:"+exportHeaderEnvPrefix+"0}") {
				assert.Containsf(t, string(got), s, "%s render", goos)
			}
			assert.NotContains(t, string(got), "super-secret-value", "header values stay out of the file")
		})
	}
}

func TestRenderMode_CollectDefaultsOnAndDataReachesTheTemplate(t *testing.T) {
	tmpl := filepath.Join(t.TempDir(), "t.tmpl")
	require.NoError(t, os.WriteFile(tmpl, []byte(`{{.Collect}} {{.OS}} {{ dig "a" "b" "" .Data }}`), 0o644))
	varsPath, outPath := renderModeFiles(t, map[string]any{"os": "windows", "data": map[string]any{"a": map[string]any{"b": "c"}}})

	require.NoError(t, RenderTemplateFile(tmpl, varsPath, outPath))

	got, err := os.ReadFile(outPath)
	require.NoError(t, err)
	assert.Equal(t, "true windows c", string(got))
}

func TestRenderMode_RejectsUnknownVars(t *testing.T) {
	varsPath, outPath := renderModeFiles(t, map[string]any{"os": "linux", "enable_sig_noz": true})

	err := RenderTemplateFile(imageTemplatePath(t), varsPath, outPath)

	assert.ErrorContains(t, err, "enable_sig_noz")
	assert.NoFileExists(t, outPath)
}
