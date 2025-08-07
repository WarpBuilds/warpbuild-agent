package uploader

import (
	"bytes"
	"context"
	"sync"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

// BufferLine represents a line with its event type
type BufferLine struct {
	Data      string
	EventType string
}

// Buffer manages a circular buffer of telemetry data
type Buffer struct {
	mu           sync.RWMutex
	lines        []BufferLine
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
		lines:      make([]BufferLine, maxLines),
		maxLines:   maxLines,
		uploadChan: uploadChan,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// AddLine adds a line to the buffer with default event type (logs)
func (b *Buffer) AddLine(line string) {
	b.AddLineWithType(line, "logs")
}

// AddLineWithType adds a line to the buffer with specified event type
func (b *Buffer) AddLineWithType(line string, eventType string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Add the line to the current position
	b.lines[b.currentIndex] = BufferLine{
		Data:      line,
		EventType: eventType,
	}
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

	// Group data by event type
	dataByType := make(map[string]*bytes.Buffer)

	// If buffer is full, we need to read from currentIndex to end, then from beginning to currentIndex
	if b.isFull {
		// Read from current position to end
		for i := b.currentIndex; i < b.maxLines; i++ {
			if b.lines[i].Data != "" {
				if dataByType[b.lines[i].EventType] == nil {
					dataByType[b.lines[i].EventType] = &bytes.Buffer{}
				}
				dataByType[b.lines[i].EventType].WriteString(b.lines[i].Data)
				dataByType[b.lines[i].EventType].WriteString("\n")
			}
		}
		// Read from beginning to current position
		for i := 0; i < b.currentIndex; i++ {
			if b.lines[i].Data != "" {
				if dataByType[b.lines[i].EventType] == nil {
					dataByType[b.lines[i].EventType] = &bytes.Buffer{}
				}
				dataByType[b.lines[i].EventType].WriteString(b.lines[i].Data)
				dataByType[b.lines[i].EventType].WriteString("\n")
			}
		}
	} else {
		// Buffer not full, just read from beginning to current position
		for i := 0; i < b.currentIndex; i++ {
			if b.lines[i].Data != "" {
				if dataByType[b.lines[i].EventType] == nil {
					dataByType[b.lines[i].EventType] = &bytes.Buffer{}
				}
				dataByType[b.lines[i].EventType].WriteString(b.lines[i].Data)
				dataByType[b.lines[i].EventType].WriteString("\n")
			}
		}
	}

	log.Logger().Debugf("Sending %d event types to upload channel", len(dataByType))

	// Send each event type separately
	for eventType, buf := range dataByType {
		data := buf.Bytes()
		if len(data) > 0 {
			req := UploadRequest{
				Data:      data,
				EventType: eventType,
			}
			select {
			case b.uploadChan <- req:
				log.Logger().Debugf("Sent %d bytes of %s data to upload channel", len(data), eventType)
			case <-b.ctx.Done():
				log.Logger().Debugf("Context cancelled, not sending to upload channel")
			default:
				log.Logger().Warnf("Upload channel is full, dropping %s data", eventType)
			}
		}
	}
}

// Flush sends all current data to the upload channel
func (b *Buffer) Flush() {
	log.Logger().Debugf("Flushing buffer...")
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
