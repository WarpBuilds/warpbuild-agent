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

// telemetryPollInterval is how often allocation details are polled while
// waiting for the runner to be allocated to a job.
const telemetryPollInterval = 5 * time.Second

// allocationStatusUnassigned is the status the backend reports until the
// runner is allocated to a job. Until then the response carries no job
// details, and therefore no export destination.
const allocationStatusUnassigned = "unassigned"

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

	sigNozEnable   bool
	sigNozEndpoint string
	sigNozAPIKey   string

	// Org-owned OTLP destination, delivered on the allocation poll. Nil
	// until the runner is allocated to a job — warm runners sitting idle
	// export nothing.
	exportCfg *exportConfig
	restartCh chan struct{}

	// drainTimeout is a field rather than a constant so tests can stop
	// waiting out the real window.
	drainTimeout time.Duration
}

// NewTelemetryManager creates a new telemetry manager
func NewTelemetryManager(ctx context.Context, port int, baseDirectory string, warpbuildAPI *warpbuild.APIClient, runnerID, pollingSecret, hostURL string, sigNozEnable bool, sigNozEndpoint, sigNozAPIKey string) *TelemetryManager {
	managerCtx, cancel := context.WithCancel(ctx)
	return &TelemetryManager{
		ctx:            managerCtx,
		cancel:         cancel,
		port:           port,
		baseDirectory:  baseDirectory,
		warpbuildAPI:   warpbuildAPI,
		runnerID:       runnerID,
		pollingSecret:  pollingSecret,
		hostURL:        hostURL,
		sigNozEnable:   sigNozEnable,
		sigNozEndpoint: sigNozEndpoint,
		sigNozAPIKey:   sigNozAPIKey,
		restartCh:      make(chan struct{}, 1),
		drainTimeout:   defaultCollectorDrainTimeout,
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

// collectorKillTimeout bounds the wait after SIGKILL, which the kernel
// should honour immediately; this only guards against an unreapable process.
const collectorKillTimeout = 5 * time.Second

// terminateCollector stops the collector and returns once it is gone.
//
// It asks nicely first: a graceful shutdown is what makes the collector
// flush its batch processors and drain its sending queue. Killing outright
// would discard up to a full batch interval, which on an ephemeral runner
// is the tail of the job — the part people care about most.
func (tm *TelemetryManager) terminateCollector(cmd *exec.Cmd, waitDone <-chan error) {
	if err := requestCollectorShutdown(cmd); err != nil {
		log.Logger().Warnf("Cannot signal OpenTelemetry Collector (%v); killing it", err)
	} else if awaitCollectorExit(waitDone, tm.drainTimeout) {
		return
	} else {
		log.Logger().Warnf("OpenTelemetry Collector did not drain within %s; killing it", tm.drainTimeout)
	}

	if err := cmd.Process.Kill(); err != nil {
		log.Logger().Errorf("Failed to kill OpenTelemetry Collector: %v", err)
	}
	if !awaitCollectorExit(waitDone, collectorKillTimeout) {
		log.Logger().Warnf("OpenTelemetry Collector still running %s after kill", collectorKillTimeout)
	}
}

// requestCollectorShutdown asks the collector to exit gracefully. Windows
// has no equivalent signal for another process, so it reports an error and
// the caller falls back to a kill.
func requestCollectorShutdown(cmd *exec.Cmd) error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("graceful signals are unsupported on windows")
	}
	log.Logger().Infof("Draining OpenTelemetry Collector (PID: %d)...", cmd.Process.Pid)
	return cmd.Process.Signal(syscall.SIGTERM)
}

// awaitCollectorExit reports whether the process exited within timeout.
func awaitCollectorExit(waitDone <-chan error, timeout time.Duration) bool {
	select {
	case err := <-waitDone:
		if err != nil {
			log.Logger().Infof("OpenTelemetry Collector exited: %v", err)
		} else {
			log.Logger().Infof("OpenTelemetry Collector exited cleanly")
		}
		return true
	case <-time.After(timeout):
		return false
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

	file, err := os.Create(tm.getConfigFilePath())
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	export := tm.currentExportConfig()

	data := struct {
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
		// Export* stays zero until the runner is allocated to a job.
		ExportEnabled       bool
		ExportEndpoint      string
		ExportHeaderEnv     map[string]string
		ExportResourceAttrs map[string]string
	}{
		LogExportFilePath:     tm.getOtelCollectorOutputFilePath(false),
		MetricsExportFilePath: tm.getOtelCollectorOutputFilePath(true),
		PushFrequency:         60 * time.Second, // Default push frequency
		OS:                    runtime.GOOS,
		Arch:                  runtime.GOARCH,
		Port:                  tm.port,
		RunnerID:              tm.runnerID,
		SigNozEndpoint:        tm.sigNozEndpoint,
		SigNozAPIKey:          tm.sigNozAPIKey,
		EnableSigNoz:          tm.sigNozEnable && tm.sigNozEndpoint != "" && tm.sigNozAPIKey != "",
	}
	if export != nil {
		data.ExportEnabled = true
		data.ExportEndpoint = export.Endpoint
		data.ExportHeaderEnv = export.headerEnv()
		data.ExportResourceAttrs = export.ResourceAttrs
	}

	// ExportHeaderEnv holds env var *names*, not credentials; the values
	// only ever live in the collector's process environment.
	log.Logger().Infof("Parsing template with vars: %+v", data)

	err = tmpl.Execute(file, data)
	if err != nil {
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

// allocationPollResult is what one allocation-details response tells us.
type allocationPollResult int

const (
	// pollWait means the runner is not allocated yet, so the response
	// carries no job details and no verdict on the export.
	pollWait allocationPollResult = iota
	// pollApplied means the export destination was found and applied.
	pollApplied
	// pollNoExport means the job details arrived without an export block:
	// this org has none configured, so this run exports nothing.
	pollNoExport
	// pollDisabled means the org has telemetry switched off.
	pollDisabled
)

// monitorTelemetryStatus polls allocation details until the runner is
// allocated to a job, reads the export destination off that response, and
// stops.
//
// It has to keep polling past the first responses because at boot the
// runner is unassigned and carries no job details. Once those arrive the
// answer is final either way — a destination, or none configured for this
// org — so there is nothing further to wait for. A change to the org's
// export config mid-job takes effect on the next run, not this one.
func (tm *TelemetryManager) monitorTelemetryStatus() {
	defer tm.wg.Done()

	log.Logger().Infof("Watching for the organization's observability export destination...")

	ticker := time.NewTicker(telemetryPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
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

			switch tm.applyAllocationDetails(allocationDetails) {
			case pollDisabled:
				// Cancel the context to stop the entire telemetry manager
				tm.cancel()
				return
			case pollApplied:
				return
			case pollNoExport:
				log.Logger().Infof("No observability export configured for this organization; nothing to fan out")
				return
			case pollWait:
			}

		case <-tm.ctx.Done():
			log.Logger().Infof("Context cancelled, stopping telemetry status monitoring...")
			return
		}
	}
}

// applyAllocationDetails folds one poll response into the manager's state.
func (tm *TelemetryManager) applyAllocationDetails(details *warpbuild.CommonsRunnerInstanceAllocationDetails) allocationPollResult {
	if details == nil {
		return pollWait
	}

	// An absent key defaults to enabled.
	if details.HasTelemetryEnabled() && !details.GetTelemetryEnabled() {
		log.Logger().Infof("Telemetry has been disabled via API. Stopping telemetry collection...")
		return pollDisabled
	}

	// Until the runner is allocated the response carries no job details, so
	// an absent export block says nothing yet.
	if details.GetStatus() == allocationStatusUnassigned {
		return pollWait
	}

	// The job details are here. Whatever the export block says now is the
	// answer for this run.
	next := exportConfigFrom(details.ObservabilityExport)
	if next == nil {
		return pollNoExport
	}

	tm.applyExportConfig(next)
	return pollApplied
}

// applyExportConfig stores the export config and restarts the collector
// onto it. Called once per run.
func (tm *TelemetryManager) applyExportConfig(next *exportConfig) {
	tm.mu.Lock()
	tm.exportCfg = next
	tm.mu.Unlock()

	log.Logger().Infof("Observability export configured: endpoint=%s", next.Endpoint)
	tm.restartCollector()
}

// restartCollector asks the supervisor to re-render the config and start a
// fresh collector. Non-blocking: a restart already pending does the same job.
func (tm *TelemetryManager) restartCollector() {
	select {
	case tm.restartCh <- struct{}{}:
	default:
	}
}
