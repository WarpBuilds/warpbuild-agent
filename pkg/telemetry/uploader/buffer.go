package uploader

import (
	"bytes"
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
	MaxLines      int
	EventType     string
}

// Buffer manages a circular buffer of telemetry data
type Buffer struct {
	buf          *bytes.Buffer
	api          *warpbuild.APIClient
	s3Uploader   *S3Uploader
	eventType    string
	maxLines     int
	currentIndex int
	uploadChan   chan UploadRequest
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewBuffer creates a new buffer with the specified options
func NewBuffer(ctx context.Context, opts BufferOptions) (*Buffer, error) {
	// Create S3 uploader
	s3Uploader := NewS3Uploader(ctx, S3UploaderOptions{
		WarpbuildAPI:  opts.WarpbuildAPI,
		EventType:     opts.EventType,
		RunnerID:      opts.RunnerID,
		PollingSecret: opts.PollingSecret,
		HostURL:       opts.HostURL,
	})

	log.Logger().Infof("Starting S3 Uploader")

	// Start S3 uploader
	if err := s3Uploader.Start(); err != nil {
		return nil, fmt.Errorf("failed to start S3 uploader for %v: %w", opts.EventType, err)
	}

	ctx, cancel := context.WithCancel(ctx)
	return &Buffer{
		buf:          bytes.NewBufferString(""),
		api:          opts.WarpbuildAPI,
		eventType:    opts.EventType,
		maxLines:     opts.MaxLines,
		currentIndex: 0,
		uploadChan:   s3Uploader.uploadChan,
		ctx:          ctx,
		cancel:       cancel,
		s3Uploader:   s3Uploader,
	}, nil
}

// AddLineWithType adds a line to the buffer with specified event type
func (b *Buffer) AddLineWithType(line []byte) {

	b.buf.Write(line)
	b.buf.WriteString("\n")
	b.currentIndex++

	if b.currentIndex == b.maxLines {
		b.sendToUploadChannel()
		b.buf.Reset()
		b.currentIndex = 0
	}
}

// sendToUploadChannel sends the current buffer content to the upload channel
func (b *Buffer) sendToUploadChannel() {

	// Log the upload process
	log.Logger().Debugf("Sending %d event types to upload channel", b.eventType)

	data := b.buf.Bytes()
	if len(data) > 0 {
		req := UploadRequest{
			Data: data,
		}

		// Try to send to upload channel
		select {
		case b.uploadChan <- req:
			log.Logger().Debugf("Sent %d bytes of %s data to upload channel", len(data), b.eventType)
		case <-b.ctx.Done():
			log.Logger().Debugf("Context cancelled, not sending to upload channel")
			return
		default:
			log.Logger().Warnf("Upload channel is full, dropping %s data", b.eventType)
			// Don't mark as sent if we couldn't send it
		}
	}

}
