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
	Enabled       bool          `json:"enabled"`
	BaseDirectory string        `json:"base_directory"`
	PushFrequency time.Duration `json:"push_frequency"`
	RunnerID      string        `json:"id"`
	PollingSecret string        `json:"polling_secret"`
	HostURL       string        `json:"host_url"`
	Port          int           `json:"port"`
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
	if opts.BaseDirectory == "" {
		opts.BaseDirectory = "/runner/warpbuild-agent"
	}

	log.Logger().Debugf("Starting OTEL receiver-based telemetry collection...")
	log.Logger().Debugf("Telemetry configuration: port=%d, push_frequency=%v",
		opts.Port, opts.PushFrequency)

	// Initialize WarpBuild API client
	cfg := warpbuild.NewConfiguration()
	if opts.HostURL == "" {
		return fmt.Errorf("host url is required")
	}
	cfg.Servers[0].URL = opts.HostURL
	wb := warpbuild.NewAPIClient(cfg)

	log.Logger().Debugf("WarpBuild API client initialized with host URL: %s", opts.HostURL)

	// Create telemetry manager with OTEL receiver
	// Use the port from settings, default to 33931 if not specified
	port := opts.Port

	manager := NewTelemetryManager(ctx, port, opts.BaseDirectory, wb, opts.RunnerID, opts.PollingSecret, opts.HostURL)

	// Start the telemetry manager
	if err := manager.Start(); err != nil {
		log.Logger().Errorf("Failed to start telemetry manager: %v", err)
		return err
	}

	log.Logger().Debugf("OTEL receiver telemetry system started on port %d", port)
	log.Logger().Debugf("OpenTelemetry Collector logs will be displayed to stdout and stderr")

	// Wait for context cancellation
	<-ctx.Done()
	log.Logger().Debugf("Context cancelled, initiating graceful shutdown...")

	// Stop the telemetry manager
	if err := manager.Stop(); err != nil {
		log.Logger().Errorf("Error stopping telemetry manager: %v", err)
	}

	log.Logger().Debugf("OTEL receiver telemetry system shutdown complete.")
	return nil
}
