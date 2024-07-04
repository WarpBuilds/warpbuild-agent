package telemetry

import (
	"context"
	"fmt"

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
	// ID is the warpbuild assigned id.
	ID            string `json:"id"`
	PollingSecret string `json:"polling_secret"`
	HostURL       string `json:"host_url"`
}

func StartTelemetryCollection(ctx context.Context, opts *TelemetryOptions) error {
	log.Logger().Infof("Starting telemetry collection...")

	cfg := warpbuild.NewConfiguration()
	if opts.HostURL == "" {
		return fmt.Errorf("host url is required")
	}
	cfg.Servers[0].URL = opts.HostURL
	wb := warpbuild.NewAPIClient(cfg)

	ctx = context.WithValue(ctx, WarpBuildAgentContextKey, wb)
	ctx = context.WithValue(ctx, WarpBuildRunnerIDContextKey, opts.ID)
	ctx = context.WithValue(ctx, WarpBuildRunnerPollingSecretContextKey, opts.PollingSecret)

	log.Logger().Infof("WarpBuild API client initialized with values: [%+v]", wb)

	// Get the appropriate OpenTelemetry Collector Contrib binary
	collectorPath, err := getOtelCollectorPath()
	if err != nil {
		log.Logger().Errorf("Failed to get OpenTelemetry Collector binary: %v", err)
		return err
	}

	log.Logger().Infof("OpenTelemetry Collector binary path: %s", collectorPath)

	// Write the OpenTelemetry Collector configuration file
	writeOtelCollectorConfig()

	log.Logger().Infof("OpenTelemetry Collector configuration file written to: %s", configFilePath)

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
		startOtelCollector(collectorPath, done)
	}()

	// Setup a filewatcher to monitor changes in the otel output file
	if err := enableOtelOutputFileWatcher(ctx); err != nil {
		log.Logger().Errorf("Failed to enable file watcher: %v", err)
	}

	<-ctx.Done()
	log.Logger().Infof("Context cancelled, initiating shutdown...")

	// Perform final upload before shutting down
	if err := readAndUploadFileToS3(ctx, syslogFilePath, 1000, false); err != nil {
		log.Logger().Errorf("Error during final upload: %v", err)
	}

	// Signal the OpenTelemetry Collector process to terminate
	done <- true
	log.Logger().Infof("Shutdown complete.")

	return nil
}
