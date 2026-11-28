package uploader

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/redact"
	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

const (
	uploadAttempts       = 3
	uploadRequestTimeout = 30 * time.Second
	uploadRetryJitter    = 0.2
)

var uploadRetryWaits = [uploadAttempts - 1]time.Duration{time.Second, 2 * time.Second}

// S3Uploader handles uploading telemetry data to S3
type S3Uploader struct {
	ctx    context.Context
	cancel context.CancelFunc
	wait   func(ctx context.Context, d time.Duration) bool
	// eventType is the otel type, can be 'metrics', 'logs', 'traces'
	eventType     string
	wg            sync.WaitGroup
	uploadChan    chan UploadRequest
	presignedURL  string
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

type retryableError struct {
	err error
}

func (e retryableError) Error() string { return e.err.Error() }

func (e retryableError) Unwrap() error { return e.err }

func retryable(err error) error {
	return retryableError{err: err}
}

func isRetryable(err error) bool {
	var r retryableError
	return errors.As(err, &r)
}

func retryableStatus(code int) bool {
	return code == http.StatusTooManyRequests || code >= http.StatusInternalServerError
}

func jittered(d time.Duration) time.Duration {
	return time.Duration(float64(d) * (1 + uploadRetryJitter*(2*rand.Float64()-1)))
}

func waitFor(ctx context.Context, d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
		return true
	case <-ctx.Done():
		return false
	}
}

// NewS3Uploader creates a new S3 uploader
func NewS3Uploader(ctx context.Context, opts S3UploaderOptions) *S3Uploader {
	uploadCtx, cancel := context.WithCancel(ctx)
	return &S3Uploader{
		ctx:           uploadCtx,
		cancel:        cancel,
		wait:          waitFor,
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
	log.Logger().Debugf("[S3Uploader] Starting upload")

	// Start the upload worker
	s.wg.Add(1)
	go s.uploadWorker()

	log.Logger().Debugf("S3 uploader started")
	return nil
}

// uploadWorker processes upload requests
func (s *S3Uploader) uploadWorker() {
	defer s.wg.Done()

	log.Logger().Debugf("S3Uploader[%s]: Upload worker started", s.eventType)

	for {
		select {
		case req := <-s.uploadChan:
			log.Logger().Debugf("S3Uploader[%s]: Received upload request with %d bytes", s.eventType, len(req.Data))
			if err := s.uploadToS3(req.Data); err != nil {
				log.Logger().Errorf("S3Uploader[%s]: Failed to upload data to S3: %v", s.eventType, err)
				// Continue processing other uploads
			} else {
				log.Logger().Debugf("S3Uploader[%s]: Successfully processed upload request", s.eventType)
			}
		case <-s.ctx.Done():
			log.Logger().Debugf("S3Uploader[%s]: Upload worker shutting down", s.eventType)
			return
		}
	}
}

func (s *S3Uploader) uploadToS3(data []byte) error {
	for attempt := 1; ; attempt++ {
		err := s.attemptUpload(data)
		if err == nil {
			return nil
		}
		if !isRetryable(err) || attempt == uploadAttempts {
			return fmt.Errorf("attempt %d/%d: %w", attempt, uploadAttempts, err)
		}
		log.Logger().Warnf("S3Uploader[%s]: Upload attempt %d/%d failed, retrying: %v", s.eventType, attempt, uploadAttempts, err)
		if !s.wait(s.ctx, jittered(uploadRetryWaits[attempt-1])) {
			return fmt.Errorf("stopped after attempt %d/%d: %w", attempt, uploadAttempts, err)
		}
	}
}

func (s *S3Uploader) attemptUpload(data []byte) error {
	url, err := s.uploadURL()
	if err != nil {
		return err
	}

	log.Logger().Debugf("Uploading %d bytes of %s data to S3", len(data), s.eventType)
	if err := s.put(url, data); err != nil {
		return err
	}
	log.Logger().Debugf("Successfully uploaded %d bytes of %s data to S3", len(data), s.eventType)

	s.setURL("")
	next, err := s.fetchPresignedURL()
	if err != nil {
		log.Logger().Warnf("S3Uploader[%s]: Failed to fetch the next presigned URL; the next upload fetches its own: %v", s.eventType, err)
		return nil
	}
	s.setURL(next)
	return nil
}

func (s *S3Uploader) uploadURL() (string, error) {
	s.mu.RLock()
	url := s.presignedURL
	s.mu.RUnlock()
	if url != "" {
		return url, nil
	}

	url, err := s.fetchPresignedURL()
	if err != nil {
		return "", err
	}
	s.setURL(url)
	return url, nil
}

func (s *S3Uploader) setURL(url string) {
	s.mu.Lock()
	s.presignedURL = url
	s.mu.Unlock()
}

func (s *S3Uploader) put(url string, data []byte) error {
	req, err := http.NewRequestWithContext(s.ctx, http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", redact.URLError(err))
	}

	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{
		Timeout: uploadRequestTimeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return retryable(fmt.Errorf("failed to execute HTTP request: %w", redact.URLError(err)))
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	err = fmt.Errorf("failed to upload data, status: %v", resp.Status)
	switch {
	case resp.StatusCode == http.StatusForbidden:
		s.setURL("")
		return retryable(err)
	case retryableStatus(resp.StatusCode):
		return retryable(err)
	}
	return err
}

func (s *S3Uploader) fetchPresignedURL() (string, error) {
	ctx, cancel := context.WithTimeout(s.ctx, uploadRequestTimeout)
	defer cancel()

	logFileName := fmt.Sprintf("%s.%s.log", time.Now().Format("20060102-150405.000000000"), s.eventType)
	out, resp, err := s.warpbuildAPI.V1RunnerInstanceAPI.
		GetRunnerInstancePresignedLogUploadURL(ctx, s.runnerID).
		XPOLLINGSECRET(s.pollingSecret).
		LogFileName(logFileName).
		Execute()
	if err != nil {
		err = fmt.Errorf("failed to call presigned URL API: %w", err)
		if resp == nil || retryableStatus(resp.StatusCode) {
			return "", retryable(err)
		}
		return "", err
	}

	if out == nil || out.Url == nil || *out.Url == "" {
		return "", errors.New("no URL received in response")
	}

	log.Logger().Debugf("Refreshed presigned URL for event type: %s", s.eventType)
	return *out.Url, nil
}

// GetUploadChannel returns the upload channel for external use
func (s *S3Uploader) GetUploadChannel() chan UploadRequest {
	return s.uploadChan
}
