package telemetry

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

var (
	debounceDelay = 10 * time.Second
	debounceTimer *time.Timer
	debounceMu    sync.Mutex
	uploadMu      sync.Mutex
)

func debouncedOtelUpload(ctx context.Context, baseDirectory string) {
	debounceMu.Lock()
	defer debounceMu.Unlock()

	if debounceTimer != nil {
		debounceTimer.Stop()
	}
	debounceTimer = time.AfterFunc(debounceDelay, func() {
		defer handlePanic()
		if err := readAndUploadFileToS3(ctx, baseDirectory, getOtelCollectorOutputFilePath(baseDirectory), 0, true); err != nil {
			log.Logger().Errorf("Error during upload: %v", err)
		}
	})
}

func readAndUploadFileToS3(ctx context.Context, baseDirectory string, filePath string, linesToRead int, shouldTruncateAfterRead bool) error {
	uploadMu.Lock()
	defer uploadMu.Unlock()

	var data []byte
	var err error
	if linesToRead > 0 {
		data, err = readLastNLines(filePath, 1000)
	} else {
		data, err = os.ReadFile(filePath)
	}
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Remove null bytes from the data. OTEL collector prepends null bytes to the log file
	data = bytes.ReplaceAll(data, []byte{0}, []byte{})
	if len(data) == 0 {
		return nil
	}

	err = uploadToPresignedURL(presignedS3URL, data)
	if err != nil {
		return fmt.Errorf("failed to upload data to S3 using presigned URL: %w", err)
	}

	// Refresh the presigned url
	url, err := fetchPresignedURL(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch presigned URL: %w", err)
	}
	presignedS3URL = url

	if shouldTruncateAfterRead {
		// Disable watcher before truncating the file. Otherwise, it will go into infinite loop
		disableOtelOutputFileWatcher()

		// Truncate the log file after successful upload
		err = os.WriteFile(filePath, []byte{}, 0)
		if err != nil {
			return fmt.Errorf("failed to truncate otel collector output file: %w", err)
		}

		// Re-enable watcher after truncating
		err = enableOtelOutputFileWatcher(ctx, baseDirectory)
		if err != nil {
			return fmt.Errorf("failed to re-enable watcher: %w", err)
		}
	}

	return nil
}

func uploadToPresignedURL(presignedURL string, data []byte) error {
	req, err := http.NewRequest("PUT", presignedURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload data, status: %v", resp.Status)
	}

	return nil
}
