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
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	uploadChan    chan UploadRequest
	presignedURL  string
	mu            sync.RWMutex
	warpbuildAPI  *warpbuild.APIClient
	runnerID      string
	pollingSecret string
	hostURL       string
	isRunning     bool
}

// UploadRequest represents a request to upload data with event type
type UploadRequest struct {
	Data      []byte
	EventType string // "logs", "metrics", "traces"
}

// NewS3Uploader creates a new S3 uploader
func NewS3Uploader(ctx context.Context, warpbuildAPI *warpbuild.APIClient, runnerID, pollingSecret, hostURL string) *S3Uploader {
	uploadCtx, cancel := context.WithCancel(ctx)
	return &S3Uploader{
		ctx:           uploadCtx,
		cancel:        cancel,
		uploadChan:    make(chan UploadRequest, 100), // Buffer for 100 upload requests
		warpbuildAPI:  warpbuildAPI,
		runnerID:      runnerID,
		pollingSecret: pollingSecret,
		hostURL:       hostURL,
	}
}

// Start starts the S3 uploader
func (s *S3Uploader) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return fmt.Errorf("S3 uploader is already running")
	}

	// Get initial presigned URL (default to logs)
	if err := s.refreshPresignedURL("logs"); err != nil {
		return fmt.Errorf("failed to get initial presigned URL: %w", err)
	}

	s.isRunning = true

	// Start the upload worker
	s.wg.Add(1)
	go s.uploadWorker()

	log.Logger().Infof("S3 uploader started")
	return nil
}

// Stop stops the S3 uploader
func (s *S3Uploader) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return nil
	}

	log.Logger().Infof("Stopping S3 uploader...")

	s.cancel()
	s.wg.Wait()
	s.isRunning = false

	log.Logger().Infof("S3 uploader stopped")
	return nil
}

// Upload uploads data to S3
func (s *S3Uploader) Upload(data []byte, eventType string) error {
	req := UploadRequest{
		Data:      data,
		EventType: eventType,
	}
	select {
	case s.uploadChan <- req:
		return nil
	case <-s.ctx.Done():
		return fmt.Errorf("uploader is stopped")
	default:
		return fmt.Errorf("upload channel is full")
	}
}

// uploadWorker processes upload requests
func (s *S3Uploader) uploadWorker() {
	defer s.wg.Done()

	for {
		select {
		case req := <-s.uploadChan:
			if err := s.uploadToS3(req.Data, req.EventType); err != nil {
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
func (s *S3Uploader) uploadToS3(data []byte, eventType string) error {
	s.mu.RLock()
	presignedURL := s.presignedURL
	s.mu.RUnlock()

	if presignedURL == "" {
		return fmt.Errorf("no presigned URL available")
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

	log.Logger().Debugf("Successfully uploaded %d bytes to S3", len(data))

	// Refresh presigned URL for next upload
	if err := s.refreshPresignedURL(eventType); err != nil {
		log.Logger().Errorf("Failed to refresh presigned URL: %v", err)
	}

	return nil
}

// refreshPresignedURL fetches a new presigned URL from the API
func (s *S3Uploader) refreshPresignedURL(eventType string) error {
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

	s.mu.Lock()
	s.presignedURL = *out.Url
	s.mu.Unlock()

	log.Logger().Debugf("Refreshed presigned URL")
	return nil
}

// GetUploadChannel returns the upload channel for external use
func (s *S3Uploader) GetUploadChannel() chan UploadRequest {
	return s.uploadChan
}

// IsRunning returns whether the uploader is currently running
func (s *S3Uploader) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isRunning
}
