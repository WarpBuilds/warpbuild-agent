package telemetry

import (
	"context"
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

// TelemetryManager coordinates all telemetry components
type TelemetryManager struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	mu     sync.RWMutex

	// Components
	buffer     *uploader.Buffer
	receiver   *uploader.Receiver
	s3Uploader *uploader.S3Uploader

	// Configuration
	port          int
	maxBufferSize int
	baseDirectory string
	warpbuildAPI  *warpbuild.APIClient
	runnerID      string
	pollingSecret string
	hostURL       string

	// State
	isRunning bool
}

// NewTelemetryManager creates a new telemetry manager
func NewTelemetryManager(ctx context.Context, port, maxBufferSize int, baseDirectory string, warpbuildAPI *warpbuild.APIClient, runnerID, pollingSecret, hostURL string) *TelemetryManager {
	managerCtx, cancel := context.WithCancel(ctx)
	return &TelemetryManager{
		ctx:           managerCtx,
		cancel:        cancel,
		port:          port,
		maxBufferSize: maxBufferSize,
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

	if tm.isRunning {
		return fmt.Errorf("telemetry manager is already running")
	}

	log.Logger().Infof("Starting telemetry manager on port %d with buffer size %d", tm.port, tm.maxBufferSize)

	// Create S3 uploader
	tm.s3Uploader = uploader.NewS3Uploader(tm.ctx, tm.warpbuildAPI, tm.runnerID, tm.pollingSecret, tm.hostURL)

	log.Logger().Infof("Starting S3 Uploader")

	// Start S3 uploader
	if err := tm.s3Uploader.Start(); err != nil {
		return fmt.Errorf("failed to start S3 uploader: %w", err)
	}

	log.Logger().Infof("Started S3 Uploader")

	// Create buffer with upload channel
	tm.buffer = uploader.NewBuffer(tm.maxBufferSize, tm.s3Uploader.GetUploadChannel())

	// Create telemetry service
	service := uploader.NewTelemetryService(tm.buffer, tm.s3Uploader)

	// Create receiver
	tm.receiver = uploader.NewReceiver(tm.port, service)

	// Start receiver
	if err := tm.receiver.Start(); err != nil {
		// Clean up S3 uploader
		tm.s3Uploader.Stop()
		return fmt.Errorf("failed to start receiver: %w", err)
	}

	log.Logger().Infof("Started receiver")

	// Start periodic buffer flush
	tm.wg.Add(1)
	go tm.periodicFlush()

	// Start OTEL collector
	tm.wg.Add(1)
	go tm.startOtelCollector()

	tm.isRunning = true
	log.Logger().Infof("Telemetry manager started successfully")
	return nil
}

// Stop stops the telemetry manager and all its components
func (tm *TelemetryManager) Stop() error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if !tm.isRunning {
		return nil
	}

	log.Logger().Infof("Stopping telemetry manager...")

	// Cancel context to stop all goroutines
	tm.cancel()

	// Stop receiver
	if tm.receiver != nil {
		if err := tm.receiver.Stop(); err != nil {
			log.Logger().Errorf("Error stopping receiver: %v", err)
		}
	}

	// Flush buffer before stopping
	if tm.buffer != nil {
		tm.buffer.Flush()
		tm.buffer.Close()
	}

	// Stop S3 uploader
	if tm.s3Uploader != nil {
		if err := tm.s3Uploader.Stop(); err != nil {
			log.Logger().Errorf("Error stopping S3 uploader: %v", err)
		}
	}

	// Wait for all goroutines to finish
	tm.wg.Wait()

	tm.isRunning = false
	log.Logger().Infof("Telemetry manager stopped")
	return nil
}

// periodicFlush periodically flushes the buffer to ensure data is uploaded even if buffer doesn't fill up
func (tm *TelemetryManager) periodicFlush() {
	defer tm.wg.Done()

	ticker := time.NewTicker(60 * time.Second) // Flush every minute
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if tm.buffer != nil {
				tm.buffer.Flush()
				log.Logger().Debugf("Periodic buffer flush completed")
			}
		case <-tm.ctx.Done():
			log.Logger().Debugf("Periodic flush worker shutting down")
			return
		}
	}
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

	// Set environment variables for better logging
	cmd.Env = append(os.Environ(),
		"OTEL_LOG_LEVEL=info",
		"OTEL_SERVICE_NAME=warpbuild-agent",
	)

	log.Logger().Infof("OpenTelemetry Collector command: %s --config %s", collectorPath, configPath)

	err := cmd.Start()
	if err != nil {
		log.Logger().Errorf("Failed to start OpenTelemetry Collector: %v", err)
		return
	}

	log.Logger().Infof("OpenTelemetry Collector started with PID: %d", cmd.Process.Pid)

	go func() {
		<-done
		log.Logger().Infof("Signaling OpenTelemetry Collector to terminate...")
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
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
		SyslogFilePath        string
		LogExportFilePath     string
		MetricsExportFilePath string
		PushFrequency         time.Duration
		OS                    string
		Arch                  string
		Port                  int
	}{
		SyslogFilePath:        tm.getSyslogFilePath(),
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

// getSyslogFilePath gets the syslog file path
func (tm *TelemetryManager) getSyslogFilePath() string {
	switch runtime.GOOS {
	case "darwin":
		return "/var/log/system.log"
	case "linux":
		return "/var/log/syslog"
	case "windows":
		return `C:\Windows\System32\winevt\Logs\System.evtx`
	default:
		log.Logger().Errorf("Unsupported OS: %s", runtime.GOOS)
		return ""
	}
}

// getOtelCollectorOutputFilePath gets the OTEL collector output file path
func (tm *TelemetryManager) getOtelCollectorOutputFilePath(isMetrics bool) string {
	if isMetrics {
		return filepath.Join(tm.baseDirectory, "otel-metrics-out.log")
	}
	return filepath.Join(tm.baseDirectory, "otel-out.log")
}

// GetStats returns statistics about the telemetry manager
func (tm *TelemetryManager) GetStats() map[string]interface{} {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	stats := map[string]interface{}{
		"is_running":  tm.isRunning,
		"port":        tm.port,
		"buffer_size": tm.maxBufferSize,
	}

	if tm.buffer != nil {
		bufferSize, isFull := tm.buffer.GetStats()
		stats["current_buffer_size"] = bufferSize
		stats["buffer_is_full"] = isFull
	}

	if tm.receiver != nil {
		stats["receiver_running"] = tm.receiver.IsRunning()
	}

	if tm.s3Uploader != nil {
		stats["s3_uploader_running"] = tm.s3Uploader.IsRunning()
	}

	return stats
}

// IsRunning returns whether the telemetry manager is currently running
func (tm *TelemetryManager) IsRunning() bool {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.isRunning
}
