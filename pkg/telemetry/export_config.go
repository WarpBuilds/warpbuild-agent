package telemetry

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

const (
	signalMetrics = "metrics"
	signalLogs    = "logs"

	// exportHeaderEnvPrefix keeps the org's ingest credential out of the
	// rendered collector config, which is long-lived and gets attached to
	// bug reports. Values reach the collector as process env instead.
	exportHeaderEnvPrefix = "WARPBUILD_OTLP_HEADER_"
)

// exportConfig is the runner-local view of the org's observability
// export, delivered on the allocation-details poll.
type exportConfig struct {
	Endpoint      string
	Metrics       bool
	Logs          bool
	Headers       map[string]string
	ResourceAttrs map[string]string
}

// exportConfigFrom converts an allocation-details payload, returning nil
// when the org has not configured an export or the config is unusable.
func exportConfigFrom(in *warpbuild.CommonsObservabilityExportConfig) *exportConfig {
	if in == nil || in.GetEndpoint() == "" {
		return nil
	}

	out := &exportConfig{
		Endpoint:      in.GetEndpoint(),
		Headers:       in.GetHeaders(),
		ResourceAttrs: in.GetResourceAttrs(),
	}
	for _, s := range in.GetSignals() {
		switch s {
		case signalMetrics:
			out.Metrics = true
		case signalLogs:
			out.Logs = true
		}
	}

	if !out.Metrics && !out.Logs {
		return nil
	}
	return out
}

// headerEnv maps each header name to the env var carrying its value.
// Header names are not valid env var identifiers, so they are indexed by
// sorted position — which also keeps the rendered config stable.
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

// envPairs renders the header values as KEY=value strings for the
// collector's environment.
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

// fingerprint identifies a config for change detection. Hashed rather
// than compared field-wise so the credential never reaches a log line.
func (e *exportConfig) fingerprint() string {
	if e == nil {
		return ""
	}
	payload, err := json.Marshal(struct {
		Endpoint      string            `json:"endpoint"`
		Metrics       bool              `json:"metrics"`
		Logs          bool              `json:"logs"`
		Headers       map[string]string `json:"headers"`
		ResourceAttrs map[string]string `json:"resource_attrs"`
	}{e.Endpoint, e.Metrics, e.Logs, e.Headers, e.ResourceAttrs})
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:8])
}
