package uploader

import (
	"context"
	"fmt"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

// BufferOptions contains all the configuration options for creating a new Buffer
type BufferOptions struct {
	WarpbuildAPI  *warpbuild.APIClient
	RunnerID      string
	PollingSecret string
	HostURL       string
	EventType     string
}

// Buffer manages a simple buffer that immediately pushes data to the uploader
//
// Buffer is no longer required and can be cleaned up.
type Buffer struct {
	api        *warpbuild.APIClient
	s3Uploader *S3Uploader
	eventType  string
	uploadChan chan UploadRequest
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewBuffer creates a new buffer with the specified options
func NewBuffer(ctx context.Context, opts BufferOptions) (*Buffer, error) {
	log.Logger().Debugf("Creating new buffer for event type: %s", opts.EventType)

	// Create S3 uploader
	s3Uploader := NewS3Uploader(ctx, S3UploaderOptions{
		WarpbuildAPI:  opts.WarpbuildAPI,
		EventType:     opts.EventType,
		RunnerID:      opts.RunnerID,
		PollingSecret: opts.PollingSecret,
		HostURL:       opts.HostURL,
	})

	log.Logger().Debugf("Starting S3 Uploader for event type: %s", opts.EventType)

	// Start S3 uploader
	if err := s3Uploader.Start(); err != nil {
		log.Logger().Errorf("Failed to start S3 uploader for %v: %w", opts.EventType, err)
		return nil, fmt.Errorf("failed to start S3 uploader for %v: %w", opts.EventType, err)
	}

	log.Logger().Debugf("S3 Uploader started successfully for event type: %s", opts.EventType)

	ctx, cancel := context.WithCancel(ctx)
	buffer := &Buffer{
		api:        opts.WarpbuildAPI,
		eventType:  opts.EventType,
		uploadChan: s3Uploader.uploadChan,
		ctx:        ctx,
		cancel:     cancel,
		s3Uploader: s3Uploader,
	}

	log.Logger().Debugf("Buffer created successfully for event type: %s, upload channel capacity: %d",
		opts.EventType, cap(s3Uploader.uploadChan))

	return buffer, nil
}

// AddLineWithType immediately sends data to the uploader (buffer size 1)
func (b *Buffer) AddLineWithType(line []byte) {
	log.Logger().Debugf("Buffer[%s]: Immediately sending %d bytes to uploader", b.eventType, len(line))

	// Validate input
	if len(line) == 0 {
		log.Logger().Warnf("Buffer[%s]: Empty data received, skipping", b.eventType)
		return
	}

	// Immediately send to upload channel
	b.sendToUploadChannel(line)
}

// sendToUploadChannel sends the data directly to the upload channel
func (b *Buffer) sendToUploadChannel(data []byte) {
	log.Logger().Debugf("Buffer[%s]: Sending %d bytes to upload channel", b.eventType, len(data))

	req := UploadRequest{
		Data: data,
	}

	// Try to send to upload channel
	select {
	case b.uploadChan <- req:
		log.Logger().Debugf("Buffer[%s]: Successfully sent %d bytes to upload channel", b.eventType, len(data))
	case <-b.ctx.Done():
		log.Logger().Debugf("Buffer[%s]: Context cancelled, not sending to upload channel", b.eventType)
		return
	default:
		log.Logger().Warnf("Buffer[%s]: Upload channel is full, dropping %d bytes of data", b.eventType, len(data))
	}
}
