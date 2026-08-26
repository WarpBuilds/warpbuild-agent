package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

func apiExport(metricsEndpoint string) *warpbuild.CommonsTelemetryExportConfig {
	out := warpbuild.NewCommonsTelemetryExportConfig()
	if metricsEndpoint != "" {
		out.SetMetricsEndpoint(metricsEndpoint)
	}
	out.SetHeaders(map[string]string{"dd-api-key": "secret"})
	out.SetResourceAttrs(map[string]string{"service.name": "ci"})
	return out
}

func TestExportConfigFrom(t *testing.T) {
	cases := []struct {
		in      *warpbuild.CommonsTelemetryExportConfig
		name    string
		wantNil bool
	}{
		{name: "nil payload", in: nil, wantNil: true},
		{name: "no endpoint", in: apiExport(""), wantNil: true},
		{name: "metrics endpoint set", in: apiExport("https://x.example.com/v1/metrics")},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := exportConfigFrom(tc.in)
			if tc.wantNil {
				assert.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			assert.Equal(t, "https://x.example.com/v1/metrics", got.MetricsEndpoint)
			assert.True(t, got.exportsMetrics())
			assert.False(t, got.exportsLogs())
		})
	}
}

func TestExportConfigHeaderEnv(t *testing.T) {
	cfg := testExportConfig()
	env := cfg.headerEnv()

	require.Len(t, env, len(cfg.Headers))
	assert.Equal(t, exportHeaderEnvPrefix+"0", env["dd-api-key"])
	assert.Equal(t, exportHeaderEnvPrefix+"1", env["dd-otel-metric-config"])

	pairs := cfg.envPairs()
	require.Len(t, pairs, len(cfg.Headers))
	assert.Contains(t, pairs, exportHeaderEnvPrefix+"0=super-secret-value")
	assert.Contains(t, pairs, exportHeaderEnvPrefix+`1={"resource_attributes_as_tags": true}`)
}

func TestExportConfigNilSafety(t *testing.T) {
	var cfg *exportConfig
	assert.Nil(t, cfg.headerEnv())
	assert.Nil(t, cfg.envPairs())
}
