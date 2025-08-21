package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

func handlePanic() {
	if r := recover(); r != nil {
		log.Logger().Errorf("Recovered from panic: %v", r)
	}
}

type contextKey string

const (
	WarpBuildAgentContextKey               contextKey = "WarpBuildAgentContextKey"
	WarpBuildRunnerIDContextKey            contextKey = "WarpBuildRunnerIDContextKey"
	WarpBuildRunnerPollingSecretContextKey contextKey = "WarpBuildRunnerPollingSecretContextKey"
)

type TelemetryOptions struct {
	Enabled                   bool          `json:"enabled"`
	BaseDirectory             string        `json:"base_directory"`
	SysLogNumberOfLinesToRead int           `json:"syslog_number_of_lines_to_read"`
	PushFrequency             time.Duration `json:"push_frequency"`
	RunnerID                  string        `json:"id"`
	PollingSecret             string        `json:"polling_secret"`
	HostURL                   string        `json:"host_url"`
	Port                      int           `json:"port"`
}

func StartTelemetryCollection(ctx context.Context, opts *TelemetryOptions) error {
	if !opts.Enabled {
		log.Logger().Infof("Telemetry collection is disabled.")
		return nil
	}

	// Fallback to defaults
	if opts.PushFrequency == 0 {
		opts.PushFrequency = 60 * time.Second
	}
	if opts.SysLogNumberOfLinesToRead == 0 {
		opts.SysLogNumberOfLinesToRead = 100
	}
	if opts.BaseDirectory == "" {
		opts.BaseDirectory = "/runner/warpbuild-agent"
	}

	log.Logger().Infof("Starting OTEL receiver-based telemetry collection...")
	log.Logger().Infof("Telemetry configuration: port=%d, buffer_size=%d, push_frequency=%v",
		opts.Port, opts.SysLogNumberOfLinesToRead, opts.PushFrequency)

	// Initialize WarpBuild API client
	cfg := warpbuild.NewConfiguration()
	if opts.HostURL == "" {
		return fmt.Errorf("host url is required")
	}
	cfg.Servers[0].URL = opts.HostURL
	wb := warpbuild.NewAPIClient(cfg)

	log.Logger().Infof("WarpBuild API client initialized with host URL: %s", opts.HostURL)

	// Create telemetry manager with OTEL receiver
	// Use the port from settings, default to 33931 if not specified
	port := opts.Port
	maxBufferSize := opts.SysLogNumberOfLinesToRead

	manager := NewTelemetryManager(ctx, port, maxBufferSize, opts.BaseDirectory, wb, opts.RunnerID, opts.PollingSecret, opts.HostURL)

	// Start the telemetry manager
	if err := manager.Start(); err != nil {
		log.Logger().Errorf("Failed to start telemetry manager: %v", err)
		return err
	}

	log.Logger().Infof("OTEL receiver telemetry system started on port %d with buffer size %d", port, maxBufferSize)
	log.Logger().Infof("OpenTelemetry Collector logs will be displayed to stdout and stderr")

	// Wait for context cancellation
	<-ctx.Done()
	log.Logger().Infof("Context cancelled, initiating graceful shutdown...")

	// Stop the telemetry manager
	if err := manager.Stop(); err != nil {
		log.Logger().Errorf("Error stopping telemetry manager: %v", err)
	}

	log.Logger().Infof("OTEL receiver telemetry system shutdown complete.")
	return nil
}
