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
		opts.SysLogNumberOfLinesToRead = 1000
	}
	if opts.BaseDirectory == "" {
		opts.BaseDirectory = "/runner/warpbuild-agent"
	}

	log.Logger().Infof("Starting telemetry collection...")

	cfg := warpbuild.NewConfiguration()
	if opts.HostURL == "" {
		return fmt.Errorf("host url is required")
	}
	cfg.Servers[0].URL = opts.HostURL
	wb := warpbuild.NewAPIClient(cfg)

	ctx = context.WithValue(ctx, WarpBuildAgentContextKey, wb)
	ctx = context.WithValue(ctx, WarpBuildRunnerIDContextKey, opts.RunnerID)
	ctx = context.WithValue(ctx, WarpBuildRunnerPollingSecretContextKey, opts.PollingSecret)

	log.Logger().Infof("WarpBuild API client initialized with values: [%+v]", wb)

	// Get the appropriate OpenTelemetry Collector Contrib binary
	collectorPath, err := getOtelCollectorPath(opts.BaseDirectory)
	if err != nil {
		log.Logger().Errorf("Failed to get OpenTelemetry Collector binary: %v", err)
		return err
	}

	log.Logger().Infof("OpenTelemetry Collector binary path: %s", collectorPath)

	// Write the OpenTelemetry Collector configuration file
	writeOtelCollectorConfig(opts.BaseDirectory, opts.PushFrequency)

	log.Logger().Infof("OpenTelemetry Collector configuration file written to: %s", getConfigFilePath(opts.BaseDirectory))

	url, err := fetchPresignedURL(ctx)
	if err != nil {
		log.Logger().Errorf("failed to fetch presigned URL: %v", err)
		return err
	}
	presignedS3URL = url

	log.Logger().Infof("Fetched initial Presigned S3 URL: %s", presignedS3URL)

	// Channel to signal when the application should terminate
	done := make(chan bool, 1)

	// Start OpenTelemetry Collector Contrib
	go func() {
		defer handlePanic()
		startOtelCollector(opts.BaseDirectory, collectorPath, done)
	}()

	// Setup a filewatcher to monitor changes in the otel output file
	if err := enableOtelOutputFileWatcher(ctx, opts.BaseDirectory); err != nil {
		log.Logger().Errorf("Failed to enable file watcher: %v", err)
	}

	<-ctx.Done()
	log.Logger().Infof("Context cancelled, initiating shutdown...")

	// Perform final upload before shutting down
	if err := readAndUploadFileToS3(ctx, opts.BaseDirectory, syslogFilePath, opts.SysLogNumberOfLinesToRead, false); err != nil {
		log.Logger().Errorf("Error during final upload: %v", err)
	}

	// Signal the OpenTelemetry Collector process to terminate
	done <- true
	log.Logger().Infof("Shutdown complete.")

	return nil
}
