package telemetry

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

func StartTelemetryCollection(opts *TelemetryOptions) error {
	cfg := warpbuild.NewConfiguration()
	if opts.HostURL == "" {
		return fmt.Errorf("host url is required")
	}
	cfg.Servers[0].URL = opts.HostURL
	wb := warpbuild.NewAPIClient(cfg)

	ctx := context.Background()
	ctx = context.WithValue(ctx, WarpBuildAgentContextKey, wb)
	ctx = context.WithValue(ctx, WarpBuildRunnerIDContextKey, opts.ID)
	ctx = context.WithValue(ctx, WarpBuildRunnerPollingSecretContextKey, opts.PollingSecret)

	// Get the appropriate OpenTelemetry Collector Contrib binary
	collectorPath, err := getOtelCollectorPath()
	if err != nil {
		log.Logger().Errorf("Failed to get OpenTelemetry Collector binary: %v", err)
		return err
	}
	// Write the OpenTelemetry Collector configuration file
	writeOtelCollectorConfig()

	url, err := fetchPresignedURL(ctx)
	if err != nil {
		log.Logger().Errorf("failed to fetch presigned URL: %v", err)
		return err
	}
	presignedS3URL = url

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

	// Set up signal handling to catch OS kill signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Logger().Infof("Received signal %s, initiating shutdown...", sig)

		// Perform final upload before shutting down
		if err := readAndUploadFileToS3(ctx, syslogFilePath, 1000, false); err != nil {
			log.Logger().Errorf("Error during final upload: %v", err)
		}

		// Signal the OpenTelemetry Collector process to terminate
		done <- true
	}()

	<-done
	log.Logger().Infof("Shutdown complete.")

	return nil
}
