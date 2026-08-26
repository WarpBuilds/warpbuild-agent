package telemetry

import (
	"fmt"
	"sort"

	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

const (
	// exportHeaderEnvPrefix keeps the org's ingest credential out of the
	// rendered collector config, which is long-lived and gets attached to
	// bug reports. Values reach the collector as process env instead.
	exportHeaderEnvPrefix = "WARPBUILD_OTLP_HEADER_"
)

// exportConfig is the runner-local view of the org's observability
// export, delivered on the allocation-details poll.
// Both metrics and logs are always exported; the two signals get their
// own exporter instances in the collector config so a destination that
// rejects one keeps accepting the other.
type exportConfig struct {
	Endpoint      string
	Headers       map[string]string
	ResourceAttrs map[string]string
}

// exportConfigFrom converts an allocation-details payload, returning nil
// when the org has not configured an export or the config is unusable.
func exportConfigFrom(in *warpbuild.CommonsObservabilityExportConfig) *exportConfig {
	if in == nil || in.GetEndpoint() == "" {
		return nil
	}

	return &exportConfig{
		Endpoint:      in.GetEndpoint(),
		Headers:       in.GetHeaders(),
		ResourceAttrs: in.GetResourceAttrs(),
	}
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
