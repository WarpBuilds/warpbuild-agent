package telemetry

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

const (
	versionA = "aaaaaaaaaaaa"
	versionB = "bbbbbbbbbbbb"
)

type recordedRequest struct {
	path   string
	query  url.Values
	header http.Header
}

type fakeBackend struct {
	*httptest.Server

	mu       sync.Mutex
	status   int
	template remoteTemplate
	details  string
	requests []recordedRequest
}

func newFakeBackend(t *testing.T) *fakeBackend {
	t.Helper()
	b := &fakeBackend{
		status:   http.StatusOK,
		template: remoteTemplate{Contract: telemetryTemplateContract, Version: versionA, Template: "receivers: {{ .OS }}\n"},
		details:  `{"status":"unassigned","telemetry_enabled":true}`,
	}
	b.Server = httptest.NewServer(http.HandlerFunc(b.serve))
	t.Cleanup(b.Close)
	return b
}

func (b *fakeBackend) serve(w http.ResponseWriter, r *http.Request) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.requests = append(b.requests, recordedRequest{r.URL.Path, r.URL.Query(), r.Header.Clone()})

	w.Header().Set("Content-Type", "application/json")
	switch {
	case strings.HasSuffix(r.URL.Path, "/telemetry_template"):
		if b.status != http.StatusOK {
			w.WriteHeader(b.status)
			return
		}
		_ = json.NewEncoder(w).Encode(b.template)
	case strings.HasSuffix(r.URL.Path, "/allocation_details"):
		_, _ = io.WriteString(w, b.details)
	default:
		http.NotFound(w, r)
	}
}

func (b *fakeBackend) set(status int, tmpl remoteTemplate) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.status, b.template = status, tmpl
}

func (b *fakeBackend) templateRequests() []recordedRequest {
	b.mu.Lock()
	defer b.mu.Unlock()
	var out []recordedRequest
	for _, r := range b.requests {
		if strings.HasSuffix(r.path, "/telemetry_template") {
			out = append(out, r)
		}
	}
	return out
}

func (b *fakeBackend) lastRequest() recordedRequest {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.requests[len(b.requests)-1]
}

func newBackendManager(t *testing.T, b *fakeBackend, base string) *TelemetryManager {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Join(base, "pkg", "telemetry"), 0o755))
	hostURL := b.URL + "/api/v1"
	tm := NewTelemetryManager(t.Context(), 33931, base, newAPIClient(hostURL, "1.2.3"), "wr_test_runner", "polling-secret", hostURL, false, "", "")
	tm.restartCh = make(chan struct{}, 4)
	return tm
}

func templateBlock(version string, data map[string]any) *warpbuild.CommonsTelemetryTemplate {
	out := warpbuild.NewCommonsTelemetryTemplate()
	out.SetContract(telemetryTemplateContract)
	if version != "" {
		out.SetVersion(version)
	}
	if data != nil {
		out.SetData(data)
	}
	return out
}

func withTemplate(d *warpbuild.CommonsRunnerInstanceAllocationDetails, block *warpbuild.CommonsTelemetryTemplate) *warpbuild.CommonsRunnerInstanceAllocationDetails {
	d.TelemetryTemplate = block
	return d
}

func TestAPIClientIdentifiesTheAgentOnEveryRequest(t *testing.T) {
	b := newFakeBackend(t)
	tm := newBackendManager(t, b, t.TempDir())

	_, err := tm.fetchAllocationDetails()
	require.NoError(t, err)
	poll := b.lastRequest()

	_, err = tm.fetchRemoteTemplate()
	require.NoError(t, err)
	fetch := b.lastRequest()

	for name, r := range map[string]recordedRequest{"allocation details": poll, "template": fetch} {
		assert.Equalf(t, "1.2.3", r.header.Get(agentVersionHeader), "%s request", name)
		assert.Equalf(t, "2", r.header.Get(telemetryContractHeader), "%s request", name)
		assert.Equalf(t, "polling-secret", r.header.Get(pollingSecretHeader), "%s request", name)
	}
	assert.Equal(t, "/api/v1/runners_instance/wr_test_runner/allocation_details", poll.path)
	assert.Equal(t, "/api/v1/runners_instance/wr_test_runner/telemetry_template", fetch.path)
	assert.Equal(t, "2", fetch.query.Get("contract"))
	assert.Equal(t, poll.header.Get("User-Agent"), fetch.header.Get("User-Agent"))
}

func TestLoadRemoteTemplate_DownloadsAndCaches(t *testing.T) {
	b := newFakeBackend(t)
	base := t.TempDir()
	tm := newBackendManager(t, b, base)

	tm.loadRemoteTemplate()

	require.NotNil(t, tm.remoteTemplate)
	assert.Equal(t, versionA, tm.remoteTemplate.Version)

	text, err := os.ReadFile(filepath.Join(base, "pkg", "telemetry", "otel-collector-config.remote.tmpl"))
	require.NoError(t, err)
	assert.Equal(t, "receivers: {{ .OS }}\n", string(text))
	meta, err := os.ReadFile(filepath.Join(base, "pkg", "telemetry", "otel-collector-config.remote.json"))
	require.NoError(t, err)
	assert.JSONEq(t, `{"contract":2,"version":"aaaaaaaaaaaa"}`, string(meta))
}

func TestLoadRemoteTemplate_FallsBackToTheCache(t *testing.T) {
	b := newFakeBackend(t)
	base := t.TempDir()
	newBackendManager(t, b, base).loadRemoteTemplate()

	b.set(http.StatusInternalServerError, remoteTemplate{})
	tm := newBackendManager(t, b, base)
	tm.loadRemoteTemplate()

	require.NotNil(t, tm.remoteTemplate, "an unreachable backend falls back to the last good download")
	assert.Equal(t, versionA, tm.remoteTemplate.Version)
	assert.Equal(t, "receivers: {{ .OS }}\n", tm.remoteTemplate.Template)
}

func TestLoadRemoteTemplate_UsesTheImageTemplate(t *testing.T) {
	t.Run("no cache", func(t *testing.T) {
		b := newFakeBackend(t)
		b.set(http.StatusBadGateway, remoteTemplate{})
		tm := newBackendManager(t, b, t.TempDir())

		tm.loadRemoteTemplate()

		assert.Nil(t, tm.remoteTemplate)
	})

	t.Run("unsupported contract ignores the cache", func(t *testing.T) {
		b := newFakeBackend(t)
		base := t.TempDir()
		newBackendManager(t, b, base).loadRemoteTemplate()

		b.set(http.StatusNotFound, remoteTemplate{})
		tm := newBackendManager(t, b, base)
		tm.loadRemoteTemplate()

		assert.Nil(t, tm.remoteTemplate, "404 means no remote template for this contract")
	})
}

func TestLoadRemoteTemplate_RejectsUnusableDownloads(t *testing.T) {
	cases := map[string]remoteTemplate{
		"other contract": {Contract: 3, Version: versionA, Template: "receivers: {}\n"},
		"no version":     {Contract: telemetryTemplateContract, Template: "receivers: {}\n"},
		"no text":        {Contract: telemetryTemplateContract, Version: versionA},
		"does not parse": {Contract: telemetryTemplateContract, Version: versionA, Template: "{{ if }}"},
		"not hermetic":   {Contract: telemetryTemplateContract, Version: versionA, Template: `{{ env "HOME" }}`},
	}
	for name, tmpl := range cases {
		t.Run(name, func(t *testing.T) {
			b := newFakeBackend(t)
			b.set(http.StatusOK, tmpl)
			base := t.TempDir()
			tm := newBackendManager(t, b, base)

			tm.loadRemoteTemplate()

			assert.Nil(t, tm.remoteTemplate)
			assert.NoFileExists(t, tm.getRemoteTemplatePath())
		})
	}
}

func TestLoadCachedTemplate_IgnoresAMismatchedCache(t *testing.T) {
	base := t.TempDir()
	tm := newTestManager(t, nil)
	tm.baseDirectory = base
	require.NoError(t, os.MkdirAll(filepath.Join(base, "pkg", "telemetry"), 0o755))

	require.NoError(t, os.WriteFile(tm.getRemoteTemplatePath(), []byte("receivers: {}\n"), 0o644))
	require.NoError(t, os.WriteFile(tm.getRemoteTemplateMetaPath(), []byte(`{"contract":1,"version":"aaaaaaaaaaaa"}`), 0o644))

	assert.Nil(t, tm.loadCachedTemplate(), "a cache written for another contract must not be rendered")
}

func TestPollFetchesAnAnnouncedTemplate(t *testing.T) {
	b := newFakeBackend(t)
	b.details = `{"status":"unassigned","telemetry_enabled":true,
	  "telemetry_template":{"contract":2,"version":"aaaaaaaaaaaa","data":{"export":{"metrics":"on"}}}}`
	tm := newBackendManager(t, b, t.TempDir())

	details, err := tm.fetchAllocationDetails()
	require.NoError(t, err)
	require.Equal(t, pollAgain, tm.handleAllocationDetails(details))

	require.NotNil(t, tm.remoteTemplate)
	assert.Equal(t, versionA, tm.remoteTemplate.Version)
	assert.Equal(t, map[string]any{"export": map[string]any{"metrics": "on"}}, tm.templateData)
	assert.Len(t, tm.restartCh, 1, "an idle runner applies the new template right away")
	assert.FileExists(t, tm.getRemoteTemplatePath())
}

func TestPollFetchesOnlyWhenTheVersionChanges(t *testing.T) {
	b := newFakeBackend(t)
	tm := newBackendManager(t, b, t.TempDir())
	tm.remoteTemplate = &remoteTemplate{Contract: telemetryTemplateContract, Version: versionA, Template: "a"}

	idle := func(version string) pollStep {
		return tm.handleAllocationDetails(withTemplate(allocationDetails(allocationStatusUnassigned, enabled(), nil), templateBlock(version, nil)))
	}

	require.Equal(t, pollAgain, idle(versionA))
	assert.Empty(t, b.templateRequests(), "the active version is not fetched again")

	b.set(http.StatusOK, remoteTemplate{Contract: telemetryTemplateContract, Version: versionB, Template: "b"})
	require.Equal(t, pollAgain, idle(versionB))
	require.Len(t, b.templateRequests(), 1)
	assert.Equal(t, versionB, tm.remoteTemplateVersion())

	require.Equal(t, pollAgain, idle(versionB))
	assert.Len(t, b.templateRequests(), 1)
	assert.Len(t, tm.restartCh, 1)
}

func TestPollRetriesAFailingVersionAtMostOncePerMinute(t *testing.T) {
	b := newFakeBackend(t)
	b.set(http.StatusInternalServerError, remoteTemplate{})
	tm := newBackendManager(t, b, t.TempDir())

	poll := func() {
		tm.handleAllocationDetails(withTemplate(allocationDetails(allocationStatusUnassigned, enabled(), nil), templateBlock(versionB, nil)))
	}

	for range 3 {
		poll()
	}
	require.Len(t, b.templateRequests(), 1)
	assert.Empty(t, tm.restartCh, "a failed download changes nothing")

	tm.templateAttemptAt = time.Now().Add(-templateRetryInterval - time.Second)
	b.set(http.StatusOK, remoteTemplate{Contract: telemetryTemplateContract, Version: versionB, Template: "b"})
	poll()

	assert.Len(t, b.templateRequests(), 2)
	assert.Equal(t, versionB, tm.remoteTemplateVersion())
}

func TestPollWithoutATemplateBlockKeepsTheActiveTemplate(t *testing.T) {
	active := &remoteTemplate{Contract: telemetryTemplateContract, Version: versionA, Template: "a"}
	data := map[string]any{"k": "v"}
	tm := &TelemetryManager{collect: true, restartCh: make(chan struct{}, 4), remoteTemplate: active, templateData: data}

	got := tm.handleAllocationDetails(mainBackendWire(t, allocationStatusUnassigned, true))

	assert.Equal(t, pollAgain, got)
	assert.Same(t, active, tm.remoteTemplate)
	assert.Equal(t, data, tm.templateData)
	assert.Empty(t, tm.restartCh)
}

func TestTemplateChangeRestartsAnIdleCollector(t *testing.T) {
	tm := &TelemetryManager{collect: true, restartCh: make(chan struct{}, 4)}
	idle := func(data map[string]any) {
		require.Equal(t, pollAgain,
			tm.handleAllocationDetails(withTemplate(allocationDetails(allocationStatusUnassigned, enabled(), nil), templateBlock("", data))))
	}

	idle(map[string]any{"a": 1.0})
	assert.Len(t, tm.restartCh, 1)

	idle(map[string]any{"a": 1.0})
	assert.Len(t, tm.restartCh, 1, "unchanged data is not a change")

	idle(map[string]any{"a": 2.0})
	assert.Len(t, tm.restartCh, 2)

	idle(nil)
	idle(map[string]any{})
	assert.Len(t, tm.restartCh, 3, "dropping the data is a change; absent and empty are the same")
}

func TestTemplateChangeAppliesAtAllocation(t *testing.T) {
	tm := &TelemetryManager{collect: true, restartCh: make(chan struct{}, 4)}
	assigned := func(data map[string]any) pollStep {
		return tm.handleAllocationDetails(withTemplate(allocationDetails("assigned", enabled(), exportForJob("build")), templateBlock("", data)))
	}

	require.Equal(t, awaitRunningJob, assigned(map[string]any{"a": 1.0}))
	require.Len(t, tm.restartCh, 1, "the export and the data arrive together: one restart")

	require.Equal(t, awaitRunningJob, assigned(map[string]any{"a": 2.0}))
	assert.Len(t, tm.restartCh, 2, "same export, new data, job not started yet: restart")

	require.Equal(t, awaitRunningJob, assigned(map[string]any{"a": 2.0}))
	assert.Len(t, tm.restartCh, 2)
}

func TestTemplateChangeWaitsForTheRunningJob(t *testing.T) {
	tm := &TelemetryManager{collect: true, restartCh: make(chan struct{}, 4)}

	require.Equal(t, awaitRunningJob, tm.handleAllocationDetails(withTemplate(
		allocationDetails("assigned", enabled(), exportForJob("build")), templateBlock("", map[string]any{"a": 1.0}))))
	require.Len(t, tm.restartCh, 1)

	require.Equal(t, awaitRunningJob, tm.handleAllocationDetails(withTemplate(
		allocationDetails(allocationStatusUnassigned, enabled(), nil), templateBlock("", map[string]any{"a": 2.0}))))
	assert.Len(t, tm.restartCh, 1, "the job may be running while the backend has not named it")

	require.Equal(t, stopPolling, tm.handleAllocationDetails(withTemplate(
		allocationDetails(allocationStatusUnassigned, enabled(), exportForJob("build")), templateBlock("", map[string]any{"a": 3.0}))))
	assert.Len(t, tm.restartCh, 1, "a running job keeps its collector")
}

func TestTemplateChangeRidesAlongWithASiblingJobRestart(t *testing.T) {
	tm := &TelemetryManager{collect: true, restartCh: make(chan struct{}, 4)}

	require.Equal(t, awaitRunningJob, tm.handleAllocationDetails(withTemplate(
		allocationDetails("assigned", enabled(), exportForJob("build (0)")), templateBlock("", map[string]any{"a": 1.0}))))
	require.Equal(t, stopPolling, tm.handleAllocationDetails(withTemplate(
		allocationDetails(allocationStatusUnassigned, enabled(), exportForJob("build (5)")), templateBlock("", map[string]any{"a": 2.0}))))

	assert.Len(t, tm.restartCh, 2)
	assert.False(t, tm.renderDirty, "the sibling-job restart renders the new data too")
}

func TestTemplateChangeAppliesAtAllocationWithoutExport(t *testing.T) {
	tm := &TelemetryManager{collect: true, restartCh: make(chan struct{}, 4)}

	require.Equal(t, pollAgain, tm.handleAllocationDetails(withTemplate(
		allocationDetails(allocationStatusUnassigned, enabled(), nil), templateBlock("", map[string]any{"a": 1.0}))))
	require.Len(t, tm.restartCh, 1)

	require.Equal(t, stopPolling, tm.handleAllocationDetails(withTemplate(
		allocationDetails("assigned", enabled(), nil), templateBlock("", map[string]any{"org": "options"}))))
	assert.Len(t, tm.restartCh, 2, "per-org data arrives at allocation, before the job starts")
}

func TestParkingRestartCoversATemplateChange(t *testing.T) {
	tm := &TelemetryManager{collect: true, restartCh: make(chan struct{}, 4)}
	disabled := false

	require.Equal(t, pollAgain, tm.handleAllocationDetails(withTemplate(
		allocationDetails(allocationStatusUnassigned, &disabled, nil), templateBlock("", map[string]any{"a": 1.0}))))

	assert.True(t, tm.parked)
	assert.Len(t, tm.restartCh, 1, "park already restarts with the new data; no second restart")
}

const backendMarker = "# rendered from the backend template"

func backendTemplate(t *testing.T) *remoteTemplate {
	t.Helper()
	image, err := os.ReadFile(imageTemplatePath(t))
	require.NoError(t, err)
	return &remoteTemplate{
		Contract: telemetryTemplateContract,
		Version:  versionB,
		Template: string(image) + "\n" + backendMarker + ` ({{ dig "flavour" "plain" .Data }})` + "\n",
	}
}

type validatorCall struct {
	config string
	env    []string
}

func rejecting(calls *[]validatorCall, needles ...string) configValidator {
	return func(_ context.Context, path string, env []string) error {
		config, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		*calls = append(*calls, validatorCall{string(config), env})
		for _, needle := range needles {
			if strings.Contains(string(config), needle) {
				return errors.New("rejected: " + needle)
			}
		}
		return nil
	}
}

func TestRender_UsesAValidBackendTemplate(t *testing.T) {
	var calls []validatorCall
	tm := newTestManager(t, testExportConfig())
	tm.remoteTemplate = backendTemplate(t)
	tm.templateData = map[string]any{"flavour": "spicy"}
	tm.validator = rejecting(&calls)

	got := renderConfig(t, tm)

	assert.Contains(t, got, backendMarker+" (spicy)")
	require.Len(t, calls, 1)
	assert.Contains(t, calls[0].env, exportHeaderEnvPrefix+"0=super-secret-value",
		"the validator sees the same header env the collector gets")
}

func TestRender_FallsBackToTheImageTemplate(t *testing.T) {
	var calls []validatorCall
	tm := newTestManager(t, testExportConfig())
	tm.remoteTemplate = backendTemplate(t)
	tm.validator = rejecting(&calls, backendMarker)

	got := renderConfig(t, tm)

	assert.Equal(t, renderConfig(t, newTestManager(t, testExportConfig())), got,
		"an invalid backend template falls back to the image's template with the same data")
	assert.Len(t, calls, 2)
}

func TestRender_DropsTheExportAsALastResort(t *testing.T) {
	var calls []validatorCall
	tm := newTestManager(t, testExportConfig())
	tm.remoteTemplate = backendTemplate(t)
	tm.validator = rejecting(&calls, "otlphttp/customer")

	got := renderConfig(t, tm)

	assert.NotContains(t, got, "customer")
	assert.NotContains(t, got, backendMarker)
	for _, line := range []string{"      exporters: [otlphttp]", "      exporters: [otlphttp/gha_logs]"} {
		assert.Contains(t, got, line, "our own collection keeps running")
	}
	assert.Len(t, calls, 2, "the last resort is used as is")
}

func TestRender_ParksWhenOnlyTheExportWasLeft(t *testing.T) {
	var calls []validatorCall
	tm := newCollectOffManager(t, testExportConfig())
	tm.validator = rejecting(&calls, "otlphttp/customer")

	got := renderConfig(t, tm)

	assert.Contains(t, got, "    metrics/idle:", "collection is off and the export is unusable: idle rather than exit")
	assert.NotContains(t, got, "customer")
}

func TestRender_LeavesNoTempFilesBehind(t *testing.T) {
	var calls []validatorCall
	tm := newTestManager(t, testExportConfig())
	tm.remoteTemplate = backendTemplate(t)
	tm.validator = rejecting(&calls, "otlphttp/customer")

	renderConfig(t, tm)

	entries, err := os.ReadDir(filepath.Join(tm.baseDirectory, "pkg", "telemetry"))
	require.NoError(t, err)
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	assert.ElementsMatch(t, []string{"binaries", "otel-collector-config.tmpl", "otel-collector-config.yaml"}, names)
}

func TestAllocationDetailsDecodeTheTemplateBlock(t *testing.T) {
	out := warpbuild.NewCommonsRunnerInstanceAllocationDetails()
	require.NoError(t, json.Unmarshal([]byte(`{"status":"unassigned",
	  "telemetry_template":{"contract":2,"version":"0123456789ab","data":{"export":{"keep":["host.name"]}}}}`), out))

	require.NotNil(t, out.TelemetryTemplate)
	assert.EqualValues(t, 2, out.TelemetryTemplate.GetContract())
	assert.Equal(t, "0123456789ab", out.TelemetryTemplate.GetVersion())
	assert.Equal(t, map[string]any{"export": map[string]any{"keep": []any{"host.name"}}}, out.TelemetryTemplate.GetData())
	assert.NotContains(t, out.AdditionalProperties, "telemetry_template")

	assert.Nil(t, mainBackendWire(t, "assigned", true).TelemetryTemplate, "an older backend sends no block")
}
