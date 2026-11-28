package telemetry

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/telemetry/uploader"
	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

const defaultCollectorDrainTimeout = 20 * time.Second

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

	exportCfg        *exportConfig
	collect          bool
	parked           bool
	remoteTemplate   *remoteTemplate
	templateData     map[string]any
	renderDirty      bool
	validator        configValidator
	restartCh        chan struct{}
	drainCh          chan struct{}
	collectorDone    chan struct{}
	collectorStarted atomic.Bool

	drainTimeout time.Duration

	templateAttemptVersion string
	templateAttemptAt      time.Time
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
		collect:        true,
		restartCh:      make(chan struct{}, 1),
		drainCh:        make(chan struct{}, 1),
		collectorDone:  make(chan struct{}),
		drainTimeout:   defaultCollectorDrainTimeout,
	}
}

// Start starts the telemetry manager and all its components
func (tm *TelemetryManager) Start() error {
	tm.loadRemoteTemplate()

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
	tm.collectorStarted.Store(true)
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

	log.Logger().Debugf("Stopping telemetry manager...")

	// Cancel context to stop all goroutines
	tm.cancel()

	// Stop receiver
	if tm.receiver != nil {
		if err := tm.receiver.Stop(); err != nil {
			log.Logger().Errorf("Error stopping receiver: %v", err)
		}
	}

	// Released before Wait: the collector goroutine takes tm.mu on its way out.
	tm.mu.Unlock()

	// Wait for all goroutines to finish
	tm.wg.Wait()

	log.Logger().Infof("Telemetry manager stopped")
	return nil
}

// startOtelCollector starts the OTEL collector process
func (tm *TelemetryManager) startOtelCollector() {
	defer tm.wg.Done()
	defer close(tm.collectorDone)

	log.Logger().Infof("Starting OpenTelemetry Collector process...")

	// Get the appropriate OpenTelemetry Collector Contrib binary
	collectorPath, err := tm.getOtelCollectorPath()
	if err != nil {
		log.Logger().Errorf("Failed to get OpenTelemetry Collector binary: %v", err)
		return
	}

	log.Logger().Infof("OpenTelemetry Collector binary path: %s", collectorPath)

	if tm.validator == nil {
		tm.validator = collectorValidator(collectorPath)
	}

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

func (tm *TelemetryManager) runOtelCollector(collectorPath string) bool {
	configPath := tm.getConfigFilePath()
	log.Logger().Infof("Starting OpenTelemetry Collector with config: %s", configPath)

	cmd := exec.Command(collectorPath, "--config", configPath)

	// Ensure OpenTelemetry collector logs are captured and displayed
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

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

	restart := false
	select {
	case <-tm.ctx.Done():
		log.Logger().Infof("Context cancelled, stopping OTEL collector (PID: %d)...", cmd.Process.Pid)
		tm.terminateCollector(cmd, waitDone)

	case <-tm.restartCh:
		log.Logger().Infof("Reloading OTEL collector (PID: %d) for new export configuration...", cmd.Process.Pid)
		tm.terminateCollector(cmd, waitDone)
		restart = true

	case <-tm.drainCh:
		log.Logger().Infof("Draining OTEL collector (PID: %d) at end of job...", cmd.Process.Pid)
		tm.terminateCollector(cmd, waitDone)

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

const collectorKillTimeout = 5 * time.Second

func (tm *TelemetryManager) terminateCollector(cmd *exec.Cmd, waitDone <-chan error) {
	if tm.drainCollector(cmd, waitDone) {
		return
	}
	killCollector(cmd, waitDone)
}

func (tm *TelemetryManager) drainCollector(cmd *exec.Cmd, waitDone <-chan error) bool {
	if err := requestCollectorShutdown(cmd); err != nil {
		log.Logger().Warnf("Cannot signal OpenTelemetry Collector (%v); killing it", err)
		return false
	}

	if awaitCollectorExit(waitDone, tm.drainTimeout) {
		return true
	}

	log.Logger().Warnf("OpenTelemetry Collector did not drain within %s; killing it", tm.drainTimeout)
	return false
}

func killCollector(cmd *exec.Cmd, waitDone <-chan error) {
	if err := cmd.Process.Kill(); err != nil {
		log.Logger().Errorf("Failed to kill OpenTelemetry Collector: %v", err)
	}

	if !awaitCollectorExit(waitDone, collectorKillTimeout) {
		log.Logger().Warnf("OpenTelemetry Collector still running %s after kill", collectorKillTimeout)
	}
}

func requestCollectorShutdown(cmd *exec.Cmd) error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("graceful signals are unsupported on windows")
	}
	log.Logger().Infof("Draining OpenTelemetry Collector (PID: %d)...", cmd.Process.Pid)
	return cmd.Process.Signal(syscall.SIGTERM)
}

func awaitCollectorExit(waitDone <-chan error, timeout time.Duration) bool {
	select {
	case err := <-waitDone:
		log.Logger().Infof("OpenTelemetry Collector exited: %v", err)
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

// writeOtelCollectorConfig writes the OTEL collector configuration
func (tm *TelemetryManager) writeOtelCollectorConfig() error {
	imagePath := tm.getConfigTemplatePath()
	image, err := os.ReadFile(imagePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	in, remote := tm.takeRenderInput()
	data := newCollectorTemplateData(in)
	log.Logger().Infof("Rendering the collector config: collect=%t parked=%t signoz=%t export_metrics=%t export_logs=%t",
		data.Collect, data.Parked, data.EnableSigNoz, data.ExportMetrics, data.ExportLogs)

	type candidate struct{ source, name, text string }
	var candidates []candidate
	if remote != nil {
		candidates = append(candidates, candidate{"telemetry template " + remote.Version, remote.Version, remote.Template})
	}
	candidates = append(candidates, candidate{"the image's template", filepath.Base(imagePath), string(image)})

	for _, c := range candidates {
		err := tm.installConfig(c.name, c.text, data, in.Export)
		if err == nil {
			log.Logger().Infof("Rendered the collector config from %s", c.source)
			return nil
		}
		log.Logger().Warnf("Rejected the collector config rendered from %s: %v", c.source, err)
	}

	fallback := in
	fallback.Export = nil
	fallback.Parked = in.Parked || !in.Collect
	config, err := renderCollectorConfig(filepath.Base(imagePath), string(image), newCollectorTemplateData(fallback))
	if err != nil {
		return err
	}
	log.Logger().Warnf("Running the collector from the image's template without the customer export")
	return writeFileAtomic(tm.getConfigFilePath(), config)
}

func (tm *TelemetryManager) installConfig(name, text string, data collectorTemplateData, export *exportConfig) error {
	config, err := renderCollectorConfig(name, text, data)
	if err != nil {
		return err
	}

	path := tm.getConfigFilePath()
	tmp, err := writeTempFile(filepath.Dir(path), "otel-collector-config.*.yaml", config)
	if err != nil {
		return err
	}
	if tm.validator != nil {
		if err := tm.validator(tm.ctx, tmp, export.envPairs()); err != nil {
			os.Remove(tmp)
			return err
		}
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}

func (tm *TelemetryManager) currentExportConfig() *exportConfig {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.exportCfg
}

func (tm *TelemetryManager) takeRenderInput() (renderInput, *remoteTemplate) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.renderDirty = false
	return renderInput{
		OS:                    runtime.GOOS,
		Arch:                  runtime.GOARCH,
		Port:                  tm.port,
		RunnerID:              tm.runnerID,
		LogExportFilePath:     tm.getOtelCollectorOutputFilePath(false),
		MetricsExportFilePath: tm.getOtelCollectorOutputFilePath(true),
		SigNozEnable:          tm.sigNozEnable,
		SigNozEndpoint:        tm.sigNozEndpoint,
		SigNozAPIKey:          tm.sigNozAPIKey,
		Collect:               tm.collect,
		Parked:                tm.parked,
		Export:                tm.exportCfg,
		Data:                  tm.templateData,
	}, tm.remoteTemplate
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

// Drain blocks until the collector has flushed and exited, so the caller can
// sequence it before VM teardown.
func (tm *TelemetryManager) Drain() {
	select {
	case tm.drainCh <- struct{}{}:
	default:
	}

	if !tm.collectorStarted.Load() {
		return
	}

	select {
	case <-tm.collectorDone:
	case <-time.After(tm.drainTimeout + collectorKillTimeout):
		log.Logger().Warnf("Telemetry drain did not finish within %s", tm.drainTimeout+collectorKillTimeout)
	}
}

func (tm *TelemetryManager) restartCollector() {
	select {
	case tm.restartCh <- struct{}{}:
	default:
	}
}
