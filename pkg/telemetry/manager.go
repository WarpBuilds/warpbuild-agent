package telemetry

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"text/template"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/telemetry/uploader"
	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

// TelemetryManager coordinates all telemetry components
type TelemetryManager struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	mu     sync.RWMutex

	// Components
	receiver          *uploader.Receiver
	s3Uploader        *uploader.S3Uploader
	otelCollectorCmd  *exec.Cmd
	otelCollectorDone chan bool

	// Configuration
	port          int
	baseDirectory string
	warpbuildAPI  *warpbuild.APIClient
	runnerID      string
	pollingSecret string
	hostURL       string
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

	// Initialize the done channel for OTEL collector
	tm.otelCollectorDone = make(chan bool, 1)

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

	// Write the OpenTelemetry Collector configuration file
	if err := tm.writeOtelCollectorConfig(); err != nil {
		log.Logger().Errorf("Failed to write OTEL collector config: %v", err)
		return
	}

	log.Logger().Infof("OpenTelemetry Collector configuration written successfully")

	// Channel to signal when the application should terminate
	done := make(chan bool, 1)

	// Start OpenTelemetry Collector Contrib
	go func() {
		defer tm.handlePanic()
		log.Logger().Infof("Launching OpenTelemetry Collector in background...")
		tm.runOtelCollector(collectorPath, done)
	}()

	// Wait for context cancellation
	<-tm.ctx.Done()
	log.Logger().Infof("Context cancelled, stopping OTEL collector...")

	// Signal the OpenTelemetry Collector process to terminate
	done <- true
}

// runOtelCollector runs the OTEL collector process
func (tm *TelemetryManager) runOtelCollector(collectorPath string, done chan bool) {
	configPath := tm.getConfigFilePath()
	log.Logger().Infof("Starting OpenTelemetry Collector with config: %s", configPath)

	cmd := exec.Command(collectorPath, "--config", configPath)

	// Ensure OpenTelemetry collector logs are captured and displayed
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Logger().Infof("OpenTelemetry Collector command: %s --config %s", collectorPath, configPath)

	err := cmd.Start()
	if err != nil {
		log.Logger().Errorf("Failed to start OpenTelemetry Collector: %v", err)
		return
	}

	// Store the command reference so we can stop it later
	tm.mu.Lock()
	tm.otelCollectorCmd = cmd
	tm.mu.Unlock()

	log.Logger().Infof("OpenTelemetry Collector started with PID: %d", cmd.Process.Pid)

	go func() {
		<-done
		log.Logger().Infof("Signaling OpenTelemetry Collector to terminate...")

		// Kill the process - this works across all platforms (Linux, macOS, Windows)
		if err := cmd.Process.Kill(); err != nil {
			log.Logger().Errorf("Failed to terminate OpenTelemetry Collector: %v", err)
		}
	}()

	go func() {
		if err := cmd.Wait(); err != nil {
			log.Logger().Errorf("OpenTelemetry Collector exited with error: %v", err)
		} else {
			log.Logger().Infof("OpenTelemetry Collector exited successfully")
		}
	}()
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
	tmpl, err := template.ParseFiles(tm.getConfigTemplatePath())
	if err != nil {
		return fmt.Errorf("failed to parse template file: %w", err)
	}

	file, err := os.Create(tm.getConfigFilePath())
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	data := struct {
		LogExportFilePath     string
		MetricsExportFilePath string
		PushFrequency         time.Duration
		OS                    string
		Arch                  string
		Port                  int
	}{
		LogExportFilePath:     tm.getOtelCollectorOutputFilePath(false),
		MetricsExportFilePath: tm.getOtelCollectorOutputFilePath(true),
		PushFrequency:         60 * time.Second, // Default push frequency
		OS:                    runtime.GOOS,
		Arch:                  runtime.GOARCH,
		Port:                  tm.port,
	}

	log.Logger().Infof("Parsing template with vars: %+v", data)

	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
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

// monitorTelemetryStatus monitors the telemetry enabled status via API polling
func (tm *TelemetryManager) monitorTelemetryStatus() {
	defer tm.wg.Done()

	log.Logger().Infof("Starting telemetry status monitoring...")

	ticker := time.NewTicker(1 * time.Second)
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

			// Check if telemetry is disabled
			if allocationDetails.HasTelemetryEnabled() {
				telemetryEnabled := allocationDetails.GetTelemetryEnabled()

				if !telemetryEnabled {
					log.Logger().Infof("Telemetry has been disabled via API. Stopping telemetry collection...")
					tm.stopOtelCollector()

					// Cancel the context to stop the entire telemetry manager
					tm.cancel()
					return
				}
			}

		case <-tm.ctx.Done():
			log.Logger().Infof("Context cancelled, stopping telemetry status monitoring...")
			return
		}
	}
}

// stopOtelCollector stops the OTEL collector process
func (tm *TelemetryManager) stopOtelCollector() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tm.otelCollectorCmd == nil || tm.otelCollectorCmd.Process == nil {
		log.Logger().Infof("OTEL collector process not running")
		return
	}

	log.Logger().Infof("Stopping OTEL collector process (PID: %d)...", tm.otelCollectorCmd.Process.Pid)

	// Send signal to the done channel to trigger graceful shutdown
	select {
	case tm.otelCollectorDone <- true:
		log.Logger().Infof("Sent termination signal to OTEL collector")
	default:
		// Channel already has a signal or is closed
		log.Logger().Debugf("OTEL collector already signaled for termination")
	}
}
