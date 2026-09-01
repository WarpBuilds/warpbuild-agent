package telemetry

import (
	"fmt"
	"sort"
	"strings"

	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

// The collector expands ${...} in config keys and values, and expands again on
// what it pulls from env, so customer-controlled strings are doubled to stay literal.
func escapeExpansion(v string) string {
	return strings.ReplaceAll(v, "$", "$$")
}

func escapeExpansionMap(in map[string]string) map[string]string {
	if in == nil {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[escapeExpansion(k)] = escapeExpansion(v)
	}
	return out
}

func escapeExpansionKeys(in map[string]string) map[string]string {
	if in == nil {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[escapeExpansion(k)] = v
	}
	return out
}

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
		pairs = append(pairs, fmt.Sprintf("%s=%s", key, escapeExpansion(e.Headers[name])))
	}
	sort.Strings(pairs)
	return pairs
}
