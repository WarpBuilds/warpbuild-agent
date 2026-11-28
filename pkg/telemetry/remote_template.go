package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

const (
	telemetryContractHeader = "X-Warpbuild-Telemetry-Contract"
	pollingSecretHeader     = "X-POLLING-SECRET"

	templateFetchTimeout  = 5 * time.Second
	templateRetryInterval = time.Minute
	maxTemplateBytes      = 1 << 20

	configValidateTimeout = 30 * time.Second
)

var errNoRemoteTemplate = errors.New("the backend has no telemetry template for this contract")

type remoteTemplate struct {
	Contract int    `json:"contract"`
	Version  string `json:"version"`
	Template string `json:"template,omitempty"`
}

func (t *remoteTemplate) check() error {
	if t.Contract != telemetryTemplateContract {
		return fmt.Errorf("template is for contract %d, not %d", t.Contract, telemetryTemplateContract)
	}
	if t.Version == "" || t.Template == "" {
		return errors.New("template is missing its version or text")
	}
	_, err := parseCollectorTemplate(t.Version, t.Template)
	return err
}

type configValidator func(ctx context.Context, configPath string, env []string) error

func collectorValidator(collectorPath string) configValidator {
	return func(ctx context.Context, configPath string, env []string) error {
		ctx, cancel := context.WithTimeout(ctx, configValidateTimeout)
		defer cancel()

		cmd := exec.CommandContext(ctx, collectorPath, "validate", "--config", configPath)
		cmd.Env = append(os.Environ(), env...)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("%w: %s", err, bytes.TrimSpace(out))
		}
		return nil
	}
}

func newAPIClient(hostURL string) *warpbuild.APIClient {
	cfg := warpbuild.NewConfiguration()
	cfg.Servers[0].URL = hostURL
	cfg.AddDefaultHeader(telemetryContractHeader, strconv.Itoa(telemetryTemplateContract))
	return warpbuild.NewAPIClient(cfg)
}

func (tm *TelemetryManager) fetchRemoteTemplate() (*remoteTemplate, error) {
	if tm.warpbuildAPI == nil {
		return nil, errors.New("no API client")
	}
	cfg := tm.warpbuildAPI.GetConfig()
	base, err := cfg.ServerURL(0, nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(tm.ctx, templateFetchTimeout)
	defer cancel()

	endpoint := fmt.Sprintf("%s/runners_instance/%s/telemetry_template?contract=%d",
		base, url.PathEscape(tm.runnerID), telemetryTemplateContract)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", cfg.UserAgent)
	for name, value := range cfg.DefaultHeader {
		req.Header.Set(name, value)
	}
	req.Header.Set(pollingSecretHeader, tm.pollingSecret)

	resp, err := cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, errNoRemoteTemplate
	default:
		return nil, fmt.Errorf("unexpected status %s", resp.Status)
	}

	var out remoteTemplate
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxTemplateBytes)).Decode(&out); err != nil {
		return nil, fmt.Errorf("failed to decode template: %w", err)
	}
	if err := out.check(); err != nil {
		return nil, err
	}
	return &out, nil
}

func (tm *TelemetryManager) getRemoteTemplatePath() string {
	return filepath.Join(tm.baseDirectory, "pkg/telemetry/otel-collector-config.remote.tmpl")
}

func (tm *TelemetryManager) getRemoteTemplateMetaPath() string {
	return filepath.Join(tm.baseDirectory, "pkg/telemetry/otel-collector-config.remote.json")
}

func (tm *TelemetryManager) cacheRemoteTemplate(t *remoteTemplate) {
	meta, err := json.Marshal(remoteTemplate{Contract: t.Contract, Version: t.Version})
	if err == nil {
		err = writeFileAtomic(tm.getRemoteTemplatePath(), []byte(t.Template))
	}
	if err == nil {
		err = writeFileAtomic(tm.getRemoteTemplateMetaPath(), meta)
	}
	if err != nil {
		log.Logger().Warnf("Failed to cache telemetry template %s: %v", t.Version, err)
	}
}

func (tm *TelemetryManager) loadCachedTemplate() *remoteTemplate {
	meta, err := os.ReadFile(tm.getRemoteTemplateMetaPath())
	if err != nil {
		return nil
	}
	text, err := os.ReadFile(tm.getRemoteTemplatePath())
	if err != nil {
		return nil
	}

	var t remoteTemplate
	if err := json.Unmarshal(meta, &t); err != nil {
		log.Logger().Warnf("Ignoring the cached telemetry template: %v", err)
		return nil
	}
	t.Template = string(text)
	if err := t.check(); err != nil {
		log.Logger().Warnf("Ignoring the cached telemetry template: %v", err)
		return nil
	}
	return &t
}

func (tm *TelemetryManager) loadRemoteTemplate() {
	t, err := tm.fetchRemoteTemplate()
	switch {
	case err == nil:
		log.Logger().Infof("Using telemetry template %s from the backend", t.Version)
		tm.cacheRemoteTemplate(t)
	case errors.Is(err, errNoRemoteTemplate):
		log.Logger().Infof("Using the image's telemetry template: %v", err)
		return
	default:
		log.Logger().Warnf("Failed to fetch the telemetry template: %v", err)
		if t = tm.loadCachedTemplate(); t == nil {
			log.Logger().Infof("Using the image's telemetry template")
			return
		}
		log.Logger().Infof("Using cached telemetry template %s", t.Version)
	}

	tm.mu.Lock()
	tm.remoteTemplate = t
	tm.mu.Unlock()
}

func (tm *TelemetryManager) remoteTemplateVersion() string {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	if tm.remoteTemplate == nil {
		return ""
	}
	return tm.remoteTemplate.Version
}

func (tm *TelemetryManager) syncTemplate(in *warpbuild.CommonsTelemetryTemplate) {
	if in == nil || (in.HasContract() && in.GetContract() != telemetryTemplateContract) {
		return
	}
	if version := in.GetVersion(); version != "" && version != tm.remoteTemplateVersion() {
		tm.refreshRemoteTemplate(version)
	}
	tm.setTemplateData(in.GetData())
}

func (tm *TelemetryManager) refreshRemoteTemplate(want string) {
	if want == tm.templateAttemptVersion && time.Since(tm.templateAttemptAt) < templateRetryInterval {
		return
	}
	tm.templateAttemptVersion, tm.templateAttemptAt = want, time.Now()

	t, err := tm.fetchRemoteTemplate()
	if err != nil {
		log.Logger().Warnf("Failed to fetch telemetry template %s: %v", want, err)
		return
	}

	tm.mu.Lock()
	changed := tm.remoteTemplate == nil || tm.remoteTemplate.Version != t.Version
	if changed {
		tm.remoteTemplate = t
		tm.renderDirty = true
	}
	tm.mu.Unlock()

	if changed {
		log.Logger().Infof("Fetched telemetry template %s", t.Version)
		tm.cacheRemoteTemplate(t)
	}
}

func (tm *TelemetryManager) setTemplateData(data map[string]any) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if (len(data) == 0 && len(tm.templateData) == 0) || reflect.DeepEqual(data, tm.templateData) {
		return
	}
	tm.templateData = data
	tm.renderDirty = true
}

func (tm *TelemetryManager) applyTemplateChange() {
	tm.mu.Lock()
	if !tm.renderDirty {
		tm.mu.Unlock()
		return
	}
	tm.renderDirty = false
	tm.mu.Unlock()

	log.Logger().Infof("Telemetry template changed; restarting the collector")
	tm.restartCollector()
}

func writeTempFile(dir, pattern string, data []byte) (string, error) {
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return "", err
	}
	_, writeErr := f.Write(data)
	if err := errors.Join(writeErr, f.Close(), os.Chmod(f.Name(), 0o644)); err != nil {
		os.Remove(f.Name())
		return "", err
	}
	return f.Name(), nil
}

func writeFileAtomic(path string, data []byte) error {
	tmp, err := writeTempFile(filepath.Dir(path), filepath.Base(path)+".*", data)
	if err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}
