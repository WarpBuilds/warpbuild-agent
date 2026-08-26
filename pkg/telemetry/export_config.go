package telemetry

import (
	"fmt"
	"sort"

	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

const (
	exportHeaderEnvPrefix = "WARPBUILD_OTLP_HEADER_"
)

type exportConfig struct {
	MetricsEndpoint string
	LogsEndpoint    string
	Headers         map[string]string
	ResourceAttrs   map[string]string
}

func (e *exportConfig) exportsMetrics() bool { return e != nil && e.MetricsEndpoint != "" }
func (e *exportConfig) exportsLogs() bool    { return e != nil && e.LogsEndpoint != "" }

func exportConfigFrom(in *warpbuild.CommonsTelemetryExportConfig) *exportConfig {
	if in == nil {
		return nil
	}
	out := &exportConfig{
		MetricsEndpoint: in.GetMetricsEndpoint(),
		LogsEndpoint:    in.GetLogsEndpoint(),
		Headers:         in.GetHeaders(),
		ResourceAttrs:   in.GetResourceAttrs(),
	}
	if !out.exportsMetrics() && !out.exportsLogs() {
		return nil
	}
	return out
}

func (e *exportConfig) headerEnv() map[string]string {
	if e == nil {
		return nil
	}
	names := make([]string, 0, len(e.Headers))
	for name := range e.Headers {
		names = append(names, name)
	}
	sort.Strings(names)

	env := make(map[string]string, len(names))
	for i, name := range names {
		env[name] = fmt.Sprintf("%s%d", exportHeaderEnvPrefix, i)
	}
	return env
}

func (e *exportConfig) envPairs() []string {
	if e == nil {
		return nil
	}
	env := e.headerEnv()
	pairs := make([]string, 0, len(env))
	for name, key := range env {
		pairs = append(pairs, fmt.Sprintf("%s=%s", key, e.Headers[name]))
	}
	sort.Strings(pairs)
	return pairs
}
