package telemetry

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"text/template"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/telemetry/uploader"
	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

// telemetryPollInterval is how often allocation details are polled for
// the telemetry kill switch and the org's export destination.
const telemetryPollInterval = 5 * time.Second

// defaultCollectorDrainTimeout bounds how long we wait for the collector
// to flush on shutdown. Generous enough for a 30s batch plus a queued
// export, short enough not to hold up a VM that is being reaped.
const defaultCollectorDrainTimeout = 25 * time.Second

// TelemetryManager coordinates all telemetry components
type TelemetryManager struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	mu     sync.RWMutex

	// Components
	receiver         *uploader.Receiver
	s3Uploader       *uploader.S3Uploader
	otelCollectorCmd *exec.Cmd

	// Configuration
	port          int
	baseDirectory string
	warpbuildAPI  *warpbuild.APIClient
	runnerID      string
	pollingSecret string
	hostURL       string

	// Org-owned OTLP destination, delivered on the allocation poll. Nil
	// until the runner is allocated to a job — warm runners sitting idle
	// export nothing.
	exportCfg         *exportConfig
	exportFingerprint string
	restartCh         chan struct{}

	// drainTimeout is a field rather than a constant so tests can stop
	// waiting out the real window.
	drainTimeout time.Duration
}

// NewTelemetryManager creates a new telemetry manager
func NewTelemetryManager(ctx context.Context, port int, baseDirectory string, warpbuildAPI *warpbuild.APIClient, runnerID, pollingSecret, hostURL string) *TelemetryManager {
	managerCtx, cancel := context.WithCancel(ctx)
	return &TelemetryManager{
		ctx:           managerCtx,
		cancel:        cancel,
		port:          port,
		baseDirectory: baseDirectory,
		warpbuildAPI:  warpbuildAPI,
		runnerID:      runnerID,
		pollingSecret: pollingSecret,
		hostURL:       hostURL,
		restartCh:     make(chan struct{}, 1),
		drainTimeout:  defaultCollectorDrainTimeout,
	}
}

// Start starts the telemetry manager and all its components
func (tm *TelemetryManager) Start() error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	log.Logger().Debugf("Starting telemetry manager on port %d", tm.port)

	log.Logger().Debugf("Started S3 Uploader")

	// Create telemetry service with required parameters
	service := uploader.NewTelemetryService(tm.warpbuildAPI, tm.runnerID, tm.pollingSecret, tm.hostURL)

	// Create receiver
	tm.receiver = uploader.NewReceiver(tm.port, service)
	tm.receiver.SetOnDrain(tm.Drain)

	// Start receiver
	if err := tm.receiver.Start(); err != nil {
		return fmt.Errorf("failed to start receiver: %w", err)
	}

	log.Logger().Debugf("Started receiver")

	// Start OTEL collector
	tm.wg.Add(1)
	go tm.startOtelCollector()

	// Start telemetry status monitoring
	tm.wg.Add(1)
	go tm.monitorTelemetryStatus()

	log.Logger().Infof("Telemetry manager started successfully")
	return nil
}

// Stop stops the telemetry manager and all its components
func (tm *TelemetryManager) Stop() error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	log.Logger().Debugf("Stopping telemetry manager...")

	// Cancel context to stop all goroutines
	tm.cancel()

	// Stop receiver
	if tm.receiver != nil {
		if err := tm.receiver.Stop(); err != nil {
			log.Logger().Errorf("Error stopping receiver: %v", err)
		}
	}

	// Wait for all goroutines to finish
	tm.wg.Wait()

	log.Logger().Infof("Telemetry manager stopped")
	return nil
}

// startOtelCollector starts the OTEL collector process
func (tm *TelemetryManager) startOtelCollector() {
	defer tm.wg.Done()

	log.Logger().Infof("Starting OpenTelemetry Collector process...")

	// Get the appropriate OpenTelemetry Collector Contrib binary
	collectorPath, err := tm.getOtelCollectorPath()
	if err != nil {
		log.Logger().Errorf("Failed to get OpenTelemetry Collector binary: %v", err)
		return
	}

	log.Logger().Infof("OpenTelemetry Collector binary path: %s", collectorPath)

	// Supervise: the collector is re-rendered and restarted whenever the
	// org's export config changes, which is how a runner picks up its
	// destination on allocation. The collector has no config hot-reload,
	// so a restart is the mechanism.
	for {
		if err := tm.writeOtelCollectorConfig(); err != nil {
			log.Logger().Errorf("Failed to write OTEL collector config: %v", err)
			return
		}

		log.Logger().Infof("OpenTelemetry Collector configuration written successfully")
		log.Logger().Infof("Launching OpenTelemetry Collector in background...")

		if restart := tm.runOtelCollector(collectorPath); !restart {
			break
		}

		log.Logger().Infof("Export configuration changed, restarting OTEL collector")
	}

	log.Logger().Infof("OTEL collector goroutine exited")
}

// runOtelCollector runs one OTEL collector process. It returns true when
// the caller should re-render the config and start a new one, false when
// the manager is shutting down or the process is gone for good.
func (tm *TelemetryManager) runOtelCollector(collectorPath string) bool {
	configPath := tm.getConfigFilePath()
	log.Logger().Infof("Starting OpenTelemetry Collector with config: %s", configPath)

	cmd := exec.Command(collectorPath, "--config", configPath)

	// Ensure OpenTelemetry collector logs are captured and displayed
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Ingest credentials are passed as env rather than written into the
	// config file, which is long-lived and ends up in bug reports.
	cmd.Env = append(os.Environ(), tm.currentExportConfig().envPairs()...)

	log.Logger().Infof("OpenTelemetry Collector command: %s --config %s", collectorPath, configPath)

	if err := cmd.Start(); err != nil {
		log.Logger().Errorf("Failed to start OpenTelemetry Collector: %v", err)
		return false
	}

	// Store the command reference so we can stop it later
	tm.mu.Lock()
	tm.otelCollectorCmd = cmd
	tm.mu.Unlock()

	log.Logger().Infof("OpenTelemetry Collector started with PID: %d", cmd.Process.Pid)

	// Channel to track when cmd.Wait() completes
	waitDone := make(chan error, 1)

	// Wait for the process to exit in a separate goroutine
	go func() {
		waitDone <- cmd.Wait()
	}()

	// Wait for context cancellation, a config change, or process exit
	restart := false
	select {
	case <-tm.ctx.Done():
		log.Logger().Infof("Context cancelled, stopping OTEL collector (PID: %d)...", cmd.Process.Pid)
		tm.terminateCollector(cmd, waitDone)

	case <-tm.restartCh:
		log.Logger().Infof("Reloading OTEL collector (PID: %d) for new export configuration...", cmd.Process.Pid)
		tm.terminateCollector(cmd, waitDone)
		restart = true

	case err := <-waitDone:
		if err != nil {
			log.Logger().Errorf("OpenTelemetry Collector exited with error: %v", err)
		} else {
			log.Logger().Infof("OpenTelemetry Collector exited successfully")
		}
	}

	log.Logger().Infof("OpenTelemetry Collector process handler completed")
	return restart
}

// terminateCollector stops the process and waits, bounded, for it to go.
//
// SIGTERM first: the collector drains its batch processors and sending
// queue on a graceful shutdown, and a SIGKILL here would drop up to a
// full batch interval of data — which on an ephemeral runner is the tail
// of the job, the part people most want. SIGKILL remains the fallback.
func (tm *TelemetryManager) terminateCollector(cmd *exec.Cmd, waitDone <-chan error) {
	if signalErr := signalCollectorShutdown(cmd); signalErr != nil {
		log.Logger().Warnf("Graceful stop unavailable (%v), killing OpenTelemetry Collector", signalErr)
	} else {
		select {
		case err := <-waitDone:
			logCollectorExit(err)
			return
		case <-time.After(tm.drainTimeout):
			log.Logger().Warnf("OpenTelemetry Collector did not drain within %s, killing it", tm.drainTimeout)
		}
	}

	if err := cmd.Process.Kill(); err != nil {
		log.Logger().Errorf("Failed to kill OpenTelemetry Collector process: %v", err)
	}

	select {
	case err := <-waitDone:
		logCollectorExit(err)
	case <-time.After(5 * time.Second):
		log.Logger().Warnf("Timeout waiting for OpenTelemetry Collector to exit after 5 seconds")
	}
}

// signalCollectorShutdown asks the collector to shut down gracefully.
// Windows has no SIGTERM equivalent for a foreign process, so the caller
// falls back to a kill there.
func signalCollectorShutdown(cmd *exec.Cmd) error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("graceful signals are not supported on windows")
	}
	log.Logger().Infof("Draining OpenTelemetry Collector (PID: %d)...", cmd.Process.Pid)
	return cmd.Process.Signal(syscall.SIGTERM)
}

func logCollectorExit(err error) {
	if err != nil {
		log.Logger().Infof("OpenTelemetry Collector terminated with error: %v", err)
	} else {
		log.Logger().Infof("OpenTelemetry Collector terminated successfully")
	}
}

// handlePanic handles panics in goroutines
func (tm *TelemetryManager) handlePanic() {
	if r := recover(); r != nil {
		log.Logger().Errorf("Recovered from panic: %v", r)
	}
}

// getOtelCollectorPath gets the path to the OTEL collector binary
func (tm *TelemetryManager) getOtelCollectorPath() (string, error) {
	var collectorPath string
	systemArch := runtime.GOARCH
	systemOS := runtime.GOOS

	binariesDir := tm.getBinariesDir()

	switch systemOS {
	case "linux":
		switch systemArch {
		case "amd64":
			collectorPath = filepath.Join(binariesDir, "linux", "amd64", "otelcol-contrib")
		case "arm64":
			collectorPath = filepath.Join(binariesDir, "linux", "arm64", "otelcol-contrib")
		default:
			return "", fmt.Errorf("unsupported architecture: %s", systemArch)
		}
	case "darwin":
		switch systemArch {
		case "amd64":
			collectorPath = filepath.Join(binariesDir, "darwin", "amd64", "otelcol-contrib")
		case "arm64":
			collectorPath = filepath.Join(binariesDir, "darwin", "arm64", "otelcol-contrib")
		default:
			return "", fmt.Errorf("unsupported architecture: %s", systemArch)
		}
	case "windows":
		if systemArch == "amd64" {
			collectorPath = filepath.Join(binariesDir, "windows", "amd64", "otelcol-contrib.exe")
		} else {
			return "", fmt.Errorf("unsupported architecture: %s", systemArch)
		}
	default:
		return "", fmt.Errorf("unsupported OS: %s", systemOS)
	}

	// Ensure the binary exists
	if _, err := os.Stat(collectorPath); os.IsNotExist(err) {
		return "", fmt.Errorf("collector binary not found at %s", collectorPath)
	}

	// Make the binary executable
	if systemOS != "windows" {
		if err := os.Chmod(collectorPath, 0755); err != nil {
			return "", fmt.Errorf("failed to make the OpenTelemetry Collector binary executable: %w", err)
		}
	}

	if systemOS == "darwin" {
		if err := exec.Command("xattr", "-rd", "com.apple.quarantine", collectorPath).Run(); err != nil {
			return "", fmt.Errorf("failed to remove quarantine attribute from binary: %w", err)
		}
	}

	return collectorPath, nil
}

// collectorTemplateData is everything otel-collector-config.tmpl reads.
// Export* is empty until the runner is allocated to a job.
type collectorTemplateData struct {
	LogExportFilePath     string
	MetricsExportFilePath string
	PushFrequency         time.Duration
	OS                    string
	Arch                  string
	Port                  int
	RunnerID              string

	ExportEnabled       bool
	ExportEndpoint      string
	ExportMetrics       bool
	ExportLogs          bool
	ExportHeaderEnv     map[string]string
	ExportResourceAttrs map[string]string
}

// yamlStr renders a value as a quoted YAML scalar. Resource attributes
// and endpoints are partly org-supplied, so they cannot be interpolated
// raw. YAML's double-quoted escapes match JSON's.
func yamlStr(v string) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// writeOtelCollectorConfig writes the OTEL collector configuration
func (tm *TelemetryManager) writeOtelCollectorConfig() error {
	tmplPath := tm.getConfigTemplatePath()
	tmpl, err := template.New(filepath.Base(tmplPath)).
		Funcs(template.FuncMap{"yamlStr": yamlStr}).
		ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse template file: %w", err)
	}

	export := tm.currentExportConfig()

	data := collectorTemplateData{
		LogExportFilePath:     tm.getOtelCollectorOutputFilePath(false),
		MetricsExportFilePath: tm.getOtelCollectorOutputFilePath(true),
		PushFrequency:         60 * time.Second, // Default push frequency
		OS:                    runtime.GOOS,
		Arch:                  runtime.GOARCH,
		Port:                  tm.port,
		RunnerID:              tm.runnerID,
	}
	if export != nil {
		data.ExportEnabled = true
		data.ExportEndpoint = export.Endpoint
		data.ExportMetrics = export.Metrics
		data.ExportLogs = export.Logs
		data.ExportHeaderEnv = export.headerEnv()
		data.ExportResourceAttrs = export.ResourceAttrs
	}

	// Logged without the header env map's values — those live in the
	// process environment, never in a log line.
	log.Logger().Infof("Rendering collector config: os=%s arch=%s port=%d export_enabled=%t export_metrics=%t export_logs=%t",
		data.OS, data.Arch, data.Port, data.ExportEnabled, data.ExportMetrics, data.ExportLogs)

	file, err := os.Create(tm.getConfigFilePath())
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// currentExportConfig returns a snapshot so rendering never holds the lock.
func (tm *TelemetryManager) currentExportConfig() *exportConfig {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.exportCfg
}

// getConfigFilePath gets the path to the OTEL collector config file
func (tm *TelemetryManager) getConfigFilePath() string {
	return filepath.Join(tm.baseDirectory, "pkg/telemetry/otel-collector-config.yaml")
}

// getConfigTemplatePath gets the path to the OTEL collector config template
func (tm *TelemetryManager) getConfigTemplatePath() string {
	return filepath.Join(tm.baseDirectory, "pkg/telemetry/otel-collector-config.tmpl")
}

// getBinariesDir gets the binaries directory
func (tm *TelemetryManager) getBinariesDir() string {
	return filepath.Join(tm.baseDirectory, "pkg/telemetry/binaries")
}

// getOtelCollectorOutputFilePath gets the OTEL collector output file path
func (tm *TelemetryManager) getOtelCollectorOutputFilePath(isMetrics bool) string {
	if isMetrics {
		return filepath.Join(tm.baseDirectory, "otel-metrics-out.log")
	}
	return filepath.Join(tm.baseDirectory, "otel-out.log")
}

// monitorTelemetryStatus polls allocation details for the telemetry kill
// switch and for the org's export destination.
//
// This must keep polling rather than settling after the first response:
// at boot the runner is unassigned and carries no export config, and the
// destination only appears once it is allocated to a job. Exiting early
// would mean never seeing that transition.
func (tm *TelemetryManager) monitorTelemetryStatus() {
	defer tm.wg.Done()

	log.Logger().Infof("Starting telemetry status monitoring...")

	ticker := time.NewTicker(telemetryPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Poll the API to check telemetry status
			allocationDetails, resp, err := tm.warpbuildAPI.V1RunnerInstanceAPI.
				GetRunnerInstanceAllocationDetails(tm.ctx, tm.runnerID).
				XPOLLINGSECRET(tm.pollingSecret).
				Execute()

			if err != nil {
				log.Logger().Debugf("Failed to get runner instance allocation details: %v", err)
				if resp != nil {
					log.Logger().Debugf("Response: %+v", resp)
				}
				continue
			}

			if allocationDetails == nil {
				log.Logger().Debugf("No runner instance allocation details found")
				continue
			}

			// An absent key defaults to enabled.
			if allocationDetails.HasTelemetryEnabled() && !allocationDetails.GetTelemetryEnabled() {
				log.Logger().Infof("Telemetry has been disabled via API. Stopping telemetry collection...")

				// Cancel the context to stop the entire telemetry manager
				tm.cancel()
				return
			}

			tm.applyExportConfig(exportConfigFrom(allocationDetails.ObservabilityExport))

		case <-tm.ctx.Done():
			log.Logger().Infof("Context cancelled, stopping telemetry status monitoring...")
			return
		}
	}
}

// Drain flushes whatever the collector is holding.
//
// Implemented as a restart rather than a stop: shutdown is what makes the
// collector drain its batch processors and sending queue, and the
// supervisor brings it straight back up, so a drain on a VM that turns
// out to live longer costs a second of collection rather than ending it.
func (tm *TelemetryManager) Drain() {
	select {
	case tm.restartCh <- struct{}{}:
	default:
		// A restart is already pending, which drains just the same.
	}
}

// applyExportConfig stores a new export config and asks the supervisor to
// restart the collector, if anything actually changed.
func (tm *TelemetryManager) applyExportConfig(next *exportConfig) {
	fingerprint := next.fingerprint()

	tm.mu.Lock()
	if fingerprint == tm.exportFingerprint {
		tm.mu.Unlock()
		return
	}
	tm.exportCfg = next
	tm.exportFingerprint = fingerprint
	tm.mu.Unlock()

	if next == nil {
		log.Logger().Infof("Observability export cleared")
	} else {
		log.Logger().Infof("Observability export configured: endpoint=%s metrics=%t logs=%t config=%s",
			next.Endpoint, next.Metrics, next.Logs, fingerprint)
	}

	// Non-blocking: a restart already pending will pick up this config
	// anyway, since the supervisor re-reads it when it re-renders.
	select {
	case tm.restartCh <- struct{}{}:
	default:
	}
}
