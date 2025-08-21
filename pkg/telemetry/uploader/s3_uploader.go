package uploader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

// S3Uploader handles uploading telemetry data to S3
type S3Uploader struct {
	ctx    context.Context
	cancel context.CancelFunc
	// eventType is the otel type, can be 'metrics', 'logs', 'traces'
	eventType     string
	wg            sync.WaitGroup
	uploadChan    chan UploadRequest
	presignedURL  string // Map of eventType to presigned URL
	mu            sync.RWMutex
	warpbuildAPI  *warpbuild.APIClient
	runnerID      string
	pollingSecret string
	hostURL       string
}

// UploadRequest represents a request to upload data with event type
type UploadRequest struct {
	Data []byte
}

// S3UploaderOptions contains all the configuration options for creating a new S3Uploader
type S3UploaderOptions struct {
	WarpbuildAPI  *warpbuild.APIClient
	EventType     string
	RunnerID      string
	PollingSecret string
	HostURL       string
}

// NewS3Uploader creates a new S3 uploader
func NewS3Uploader(ctx context.Context, opts S3UploaderOptions) *S3Uploader {
	uploadCtx, cancel := context.WithCancel(ctx)
	return &S3Uploader{
		ctx:           uploadCtx,
		cancel:        cancel,
		eventType:     opts.EventType,
		uploadChan:    make(chan UploadRequest, 100), // Buffer for 100 upload requests
		presignedURL:  "",
		warpbuildAPI:  opts.WarpbuildAPI,
		runnerID:      opts.RunnerID,
		pollingSecret: opts.PollingSecret,
		hostURL:       opts.HostURL,
	}
}

// Start starts the S3 uploader
func (s *S3Uploader) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Logger().Infof("[S3Uploader] Starting upload")

	log.Logger().Infof("[S3Uploader] Fetching initial presigned urls")

	// Get initial presigned URLs for logs, metrics, and traces
	if err := s.refreshPresignedURL(); err != nil {
		return fmt.Errorf("failed to get initial logs presigned URL: %w", err)
	}

	// Start the upload worker
	s.wg.Add(1)
	go s.uploadWorker()

	log.Logger().Infof("S3 uploader started")
	return nil
}

// uploadWorker processes upload requests
func (s *S3Uploader) uploadWorker() {
	defer s.wg.Done()

	for {
		select {
		case req := <-s.uploadChan:
			if err := s.uploadToS3(req.Data); err != nil {
				log.Logger().Errorf("Failed to upload data to S3: %v", err)
				// Continue processing other uploads
			}
		case <-s.ctx.Done():
			log.Logger().Debugf("S3 upload worker shutting down")
			return
		}
	}
}

// uploadToS3 uploads data to S3 using the presigned URL
func (s *S3Uploader) uploadToS3(data []byte) error {
	s.mu.RLock()
	presignedURL := s.presignedURL
	eventType := s.eventType
	s.mu.RUnlock()

	log.Logger().Infof("Uploading %d bytes of %s data to S3", len(data), eventType)

	if presignedURL == "" {
		return fmt.Errorf("no presigned URL available for event type: %s", eventType)
	}

	req, err := http.NewRequest("PUT", presignedURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload data, status: %v", resp.Status)
	}

	log.Logger().Debugf("Successfully uploaded %d bytes of %s data to S3", len(data), eventType)

	// Refresh presigned URL for next upload
	if err := s.refreshPresignedURL(); err != nil {
		log.Logger().Errorf("Failed to refresh presigned URL: %v", err)
	}

	return nil
}

// refreshPresignedURL fetches a new presigned URL from the API
func (s *S3Uploader) refreshPresignedURL() error {
	eventType := s.eventType
	// Create context with timeout for the API call
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	logFileName := fmt.Sprintf("%s.%s.log", time.Now().Format("20060102-150405"), eventType)
	out, resp, err := s.warpbuildAPI.V1RunnerInstanceAPI.
		GetRunnerInstancePresignedLogUploadURL(ctx, s.runnerID).
		XPOLLINGSECRET(s.pollingSecret).
		LogFileName(logFileName).
		Execute()
	if err != nil {
		return fmt.Errorf("failed to call presigned URL API: %w", err)
	}

	if resp == nil || resp.Body == nil {
		return fmt.Errorf("empty response from presigned URL API")
	}

	// Read and log response body for debugging
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Logger().Errorf("Error reading response body: %v", err)
	}
	log.Logger().Debugf("Response body: %s", string(body))

	if out == nil || out.Url == nil {
		return fmt.Errorf("no URL received in response")
	}

	s.presignedURL = *out.Url

	log.Logger().Debugf("Refreshed presigned URL for event type: %s", eventType)
	return nil
}

// GetUploadChannel returns the upload channel for external use
func (s *S3Uploader) GetUploadChannel() chan UploadRequest {
	return s.uploadChan
}
