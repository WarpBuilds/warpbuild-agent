package telemetry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"

	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

const telemetryTemplateContract = 2

func yamlStr(v string) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

var ottlEscaper = strings.NewReplacer(`\`, `\\`, `"`, `\"`)

func ottlStr(s string) string {
	return `"` + ottlEscaper.Replace(s) + `"`
}

func templateFuncs() template.FuncMap {
	funcs := sprig.HermeticTxtFuncMap()
	funcs["yamlStr"] = yamlStr
	funcs["ottlStr"] = ottlStr
	return funcs
}

func parseCollectorTemplate(name, text string) (*template.Template, error) {
	tmpl, err := template.New(name).Funcs(templateFuncs()).Parse(text)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}
	return tmpl, nil
}

type renderInput struct {
	OS                    string
	Arch                  string
	Port                  int
	RunnerID              string
	LogExportFilePath     string
	MetricsExportFilePath string
	SigNozEnable          bool
	SigNozEndpoint        string
	SigNozAPIKey          string
	Collect               bool
	Parked                bool
	Export                *exportConfig
	Data                  map[string]any
}

type collectorTemplateData struct {
	LogExportFilePath     string
	MetricsExportFilePath string
	PushFrequency         time.Duration
	OS                    string
	Arch                  string
	Port                  int
	RunnerID              string
	SigNozEndpoint        string
	SigNozAPIKey          string
	EnableSigNoz          bool
	Collect               bool
	Parked                bool
	LogsExporters         string
	GHALogsExporters      string
	MetricsExporters      string
	ExportMetrics         bool
	ExportLogs            bool
	ExportMetricsEndpoint string
	ExportLogsEndpoint    string
	ExportHeaderEnv       map[string]string
	ExportResourceAttrs   map[string]string
	ExportDatapointAttrs  []string
	ExportMetricsInterval int
	Data                  map[string]any
}

func newCollectorTemplateData(in renderInput) collectorTemplateData {
	data := collectorTemplateData{
		LogExportFilePath:     in.LogExportFilePath,
		MetricsExportFilePath: in.MetricsExportFilePath,
		PushFrequency:         60 * time.Second, // Default push frequency
		OS:                    in.OS,
		Arch:                  in.Arch,
		Port:                  in.Port,
		RunnerID:              in.RunnerID,
		SigNozEndpoint:        in.SigNozEndpoint,
		SigNozAPIKey:          in.SigNozAPIKey,
		EnableSigNoz:          in.SigNozEnable && in.SigNozEndpoint != "" && in.SigNozAPIKey != "",
		Collect:               in.Collect,
		Parked:                in.Parked,
		Data:                  cloneData(in.Data),
	}
	if export := in.Export; export != nil {
		data.ExportMetrics = export.exportsMetrics()
		data.ExportLogs = export.exportsLogs()
		data.ExportMetricsEndpoint = escapeExpansion(export.MetricsEndpoint)
		data.ExportLogsEndpoint = escapeExpansion(export.LogsEndpoint)
		data.ExportHeaderEnv = escapeExpansionKeys(export.headerEnv())
		data.ExportResourceAttrs = escapeExpansionMap(export.ResourceAttrs)
		data.ExportDatapointAttrs = datapointLabelAttrs(export.ResourceAttrs)
		data.ExportMetricsInterval = export.MetricsIntervalSeconds
	}

	data.LogsExporters = exporterList(
		exporterEntry{"otlphttp", in.Collect},
		exporterEntry{"forward/customer_logs", data.ExportLogs},
	)
	data.GHALogsExporters = exporterList(
		exporterEntry{"otlphttp/gha_logs", in.Collect},
		exporterEntry{"forward/customer_logs", data.ExportLogs},
	)
	data.MetricsExporters = exporterList(
		exporterEntry{"otlphttp", in.Collect},
		exporterEntry{"otlp/signoz", data.EnableSigNoz},
		exporterEntry{"debug", data.EnableSigNoz},
		exporterEntry{"forward/customer_metrics", data.ExportMetrics},
	)
	return data
}

func cloneData(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = cloneValue(v)
	}
	return out
}

func cloneValue(v any) any {
	switch t := v.(type) {
	case map[string]any:
		return cloneData(t)
	case []any:
		out := make([]any, len(t))
		for i, e := range t {
			out[i] = cloneValue(e)
		}
		return out
	default:
		return v
	}
}

func renderCollectorConfig(name, text string, data collectorTemplateData) ([]byte, error) {
	tmpl, err := parseCollectorTemplate(name, text)
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	if err := tmpl.Execute(&out, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}
	return out.Bytes(), nil
}

type renderVars struct {
	OS             string                                  `json:"os"`
	Arch           string                                  `json:"arch"`
	Port           int                                     `json:"port"`
	RunnerID       string                                  `json:"runner_id"`
	Collect        *bool                                   `json:"collect"`
	Parked         bool                                    `json:"parked"`
	EnableSigNoz   bool                                    `json:"enable_signoz"`
	SigNozEndpoint string                                  `json:"signoz_endpoint"`
	SigNozAPIKey   string                                  `json:"signoz_api_key"`
	Export         *warpbuild.CommonsTelemetryExportConfig `json:"export"`
	Data           map[string]any                          `json:"data"`
}

func (v renderVars) renderInput() renderInput {
	in := renderInput{
		OS:             v.OS,
		Arch:           v.Arch,
		Port:           v.Port,
		RunnerID:       v.RunnerID,
		SigNozEnable:   v.EnableSigNoz,
		SigNozEndpoint: v.SigNozEndpoint,
		SigNozAPIKey:   v.SigNozAPIKey,
		Collect:        v.Collect == nil || *v.Collect,
		Parked:         v.Parked,
		Export:         exportConfigFrom(v.Export),
		Data:           v.Data,
	}
	if in.OS == "" {
		in.OS = runtime.GOOS
	}
	if in.Arch == "" {
		in.Arch = runtime.GOARCH
	}
	return in
}

func RenderTemplateFile(templatePath, varsPath, outPath string) error {
	text, err := os.ReadFile(templatePath)
	if err != nil {
		return err
	}

	var vars renderVars
	if varsPath != "" {
		raw, err := os.ReadFile(varsPath)
		if err != nil {
			return err
		}
		dec := json.NewDecoder(bytes.NewReader(raw))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&vars); err != nil {
			return fmt.Errorf("failed to parse %s: %w", varsPath, err)
		}
	}

	out, err := renderCollectorConfig(filepath.Base(templatePath), string(text), newCollectorTemplateData(vars.renderInput()))
	if err != nil {
		return err
	}
	if outPath == "" {
		_, err = os.Stdout.Write(out)
		return err
	}
	return os.WriteFile(outPath, out, 0o644)
}
