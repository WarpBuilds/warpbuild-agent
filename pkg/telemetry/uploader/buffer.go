package uploader

import (
	"bytes"
	"context"
	"sync"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

// Buffer manages a circular buffer of telemetry data
type Buffer struct {
	mu           sync.RWMutex
	lines        []string
	maxLines     int
	currentIndex int
	isFull       bool
	uploadChan   chan UploadRequest
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewBuffer creates a new buffer with the specified capacity
func NewBuffer(maxLines int, uploadChan chan UploadRequest) *Buffer {
	ctx, cancel := context.WithCancel(context.Background())
	return &Buffer{
		lines:      make([]string, maxLines),
		maxLines:   maxLines,
		uploadChan: uploadChan,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// AddLine adds a line to the buffer
func (b *Buffer) AddLine(line string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Add the line to the current position
	b.lines[b.currentIndex] = line
	b.currentIndex++

	// If we've reached the end, wrap around
	if b.currentIndex >= b.maxLines {
		b.currentIndex = 0
		b.isFull = true
	}

	// If buffer is full, send data to upload channel
	if b.isFull {
		go b.sendToUploadChannel()
	}
}

// sendToUploadChannel sends the current buffer content to the upload channel
func (b *Buffer) sendToUploadChannel() {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var buf bytes.Buffer

	// If buffer is full, we need to read from currentIndex to end, then from beginning to currentIndex
	if b.isFull {
		// Write from current position to end
		for i := b.currentIndex; i < b.maxLines; i++ {
			if b.lines[i] != "" {
				buf.WriteString(b.lines[i])
				buf.WriteString("\n")
			}
		}
		// Write from beginning to current position
		for i := 0; i < b.currentIndex; i++ {
			if b.lines[i] != "" {
				buf.WriteString(b.lines[i])
				buf.WriteString("\n")
			}
		}
	} else {
		// Buffer not full, just write from beginning to current position
		for i := 0; i < b.currentIndex; i++ {
			if b.lines[i] != "" {
				buf.WriteString(b.lines[i])
				buf.WriteString("\n")
			}
		}
	}

	data := buf.Bytes()
	if len(data) > 0 {
		// Create upload request with default event type (logs)
		req := UploadRequest{
			Data:      data,
			EventType: "logs", // Default to logs for buffer uploads
		}
		select {
		case b.uploadChan <- req:
			log.Logger().Debugf("Sent %d bytes to upload channel", len(data))
		case <-b.ctx.Done():
			log.Logger().Debugf("Context cancelled, not sending to upload channel")
		default:
			log.Logger().Warnf("Upload channel is full, dropping data")
		}
	}
}

// Flush sends all current data to the upload channel
func (b *Buffer) Flush() {
	b.sendToUploadChannel()
}

// Close closes the buffer and cancels the context
func (b *Buffer) Close() {
	b.cancel()
}

// GetStats returns buffer statistics
func (b *Buffer) GetStats() (int, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.isFull {
		return b.maxLines, true
	}
	return b.currentIndex, false
}
