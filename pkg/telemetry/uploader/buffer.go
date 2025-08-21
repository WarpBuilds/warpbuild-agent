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
	log.Logger().Infof("Creating new buffer for event type: %s with max lines: %d", opts.EventType, opts.MaxLines)

	// Create S3 uploader
	s3Uploader := NewS3Uploader(ctx, S3UploaderOptions{
		WarpbuildAPI:  opts.WarpbuildAPI,
		EventType:     opts.EventType,
		RunnerID:      opts.RunnerID,
		PollingSecret: opts.PollingSecret,
		HostURL:       opts.HostURL,
	})

	log.Logger().Infof("Starting S3 Uploader for event type: %s", opts.EventType)

	// Start S3 uploader
	if err := s3Uploader.Start(); err != nil {
		log.Logger().Errorf("Failed to start S3 uploader for %v: %w", opts.EventType, err)
		return nil, fmt.Errorf("failed to start S3 uploader for %v: %w", opts.EventType, err)
	}

	log.Logger().Infof("S3 Uploader started successfully for event type: %s", opts.EventType)

	ctx, cancel := context.WithCancel(ctx)
	buffer := &Buffer{
		buf:          bytes.NewBufferString(""),
		api:          opts.WarpbuildAPI,
		eventType:    opts.EventType,
		maxLines:     opts.MaxLines,
		currentIndex: 0,
		uploadChan:   s3Uploader.uploadChan,
		ctx:          ctx,
		cancel:       cancel,
		s3Uploader:   s3Uploader,
	}

	log.Logger().Infof("Buffer created successfully for event type: %s, upload channel capacity: %d",
		opts.EventType, cap(s3Uploader.uploadChan))

	return buffer, nil
}

// AddLineWithType adds a line to the buffer with specified event type
func (b *Buffer) AddLineWithType(line []byte) {
	// Safety check to prevent nil pointer dereference
	if b.buf == nil {
		log.Logger().Errorf("Buffer buf is nil for event type %s, initializing", b.eventType)
		b.buf = bytes.NewBufferString("")
	}

	log.Logger().Debugf("Buffer[%s]: Adding line of %d bytes, current index: %d/%d",
		b.eventType, len(line), b.currentIndex, b.maxLines)

	b.buf.Write(line)
	b.buf.WriteString("\n")
	b.currentIndex++

	log.Logger().Debugf("Buffer[%s]: After adding line, current index: %d/%d, buffer size: %d bytes",
		b.eventType, b.currentIndex, b.maxLines, b.buf.Len())

	if b.currentIndex == b.maxLines {
		log.Logger().Infof("Buffer[%s]: Buffer is full (%d/%d), triggering upload",
			b.eventType, b.currentIndex, b.maxLines)
		b.sendToUploadChannel()
		b.buf.Reset()
		b.currentIndex = 0
		log.Logger().Debugf("Buffer[%s]: Buffer reset, current index: %d", b.eventType, b.currentIndex)
	}
}

// sendToUploadChannel sends the current buffer content to the upload channel
func (b *Buffer) sendToUploadChannel() {
	log.Logger().Debugf("Buffer[%s]: sendToUploadChannel called, buffer size: %d bytes",
		b.eventType, b.buf.Len())

	// Log the upload process
	log.Logger().Debugf("Sending %d event types to upload channel", b.eventType)

	data := b.buf.Bytes()
	if len(data) > 0 {
		log.Logger().Debugf("Buffer[%s]: Preparing to send %d bytes to upload channel",
			b.eventType, len(data))

		req := UploadRequest{
			Data: data,
		}

		// Try to send to upload channel
		select {
		case b.uploadChan <- req:
			log.Logger().Infof("Buffer[%s]: Successfully sent %d bytes to upload channel",
				b.eventType, len(data))
		case <-b.ctx.Done():
			log.Logger().Debugf("Buffer[%s]: Context cancelled, not sending to upload channel", b.eventType)
			return
		default:
			log.Logger().Warnf("Buffer[%s]: Upload channel is full, dropping %d bytes of data",
				b.eventType, len(data))
			// Don't mark as sent if we couldn't send it
		}
	} else {
		log.Logger().Warnf("Buffer[%s]: No data to send (buffer is empty)", b.eventType)
	}
}
