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
	b.mu.Lock()
	defer b.mu.Unlock()

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

	// Log the upload process
	log.Logger().Debugf("Sending %d event types to upload channel", len(dataByType))
	log.Logger().Debugf("Event types to send: %v", getKeys(dataByType))

	// Track which event types were successfully sent
	sentEventTypes := make(map[string]bool)

	// Send each event type separately
	for eventType, buf := range dataByType {
		data := buf.Bytes()
		if len(data) > 0 {
			req := UploadRequest{
				Data:      data,
				EventType: eventType,
			}

			// Try to send to upload channel
			select {
			case b.uploadChan <- req:
				log.Logger().Debugf("Sent %d bytes of %s data to upload channel", len(data), eventType)
				sentEventTypes[eventType] = true
			case <-b.ctx.Done():
				log.Logger().Debugf("Context cancelled, not sending to upload channel")
				return
			default:
				log.Logger().Warnf("Upload channel is full, dropping %s data", eventType)
				// Don't mark as sent if we couldn't send it
				continue
			}
		}
	}

	log.Logger().Debugf("Successfully sent event types: %v", getBoolKeys(sentEventTypes))

	// Only drop events of successfully sent event types to prevent duplicates
	for eventType := range sentEventTypes {
		b.dropEventsOfType(eventType)
		log.Logger().Debugf("Dropped all %s events after successful upload", eventType)
	}
}

// getKeys returns the keys of a map as a slice
func getKeys(m map[string]*bytes.Buffer) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// getBoolKeys returns the keys of a map[string]bool as a slice
func getBoolKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// dropEventsOfType removes all events of a specific type from the buffer
func (b *Buffer) dropEventsOfType(eventType string) {
	droppedCount := 0

	// Clear all lines of the specified event type
	for i := 0; i < b.maxLines; i++ {
		if b.lines[i].EventType == eventType {
			b.lines[i] = BufferLine{}
			droppedCount++
		}
	}

	// If we dropped all events, reset the buffer state
	if droppedCount > 0 {
		// Check if buffer is now empty
		remainingEvents := 0
		for i := 0; i < b.maxLines; i++ {
			if b.lines[i].Data != "" {
				remainingEvents++
			}
		}

		if remainingEvents == 0 {
			// Buffer is completely empty, reset state
			b.currentIndex = 0
			b.isFull = false
		} else {
			// Compact the buffer by removing empty lines
			b.compactBuffer()
		}

		log.Logger().Debugf("Dropped %d %s events, %d events remaining", droppedCount, eventType, remainingEvents)
	}
}

// compactBuffer removes empty lines and compacts the buffer
func (b *Buffer) compactBuffer() {
	// Create a new compacted buffer
	compacted := make([]BufferLine, b.maxLines)
	compactedIndex := 0

	// Copy non-empty lines to the beginning
	for i := 0; i < b.maxLines; i++ {
		if b.lines[i].Data != "" {
			compacted[compactedIndex] = b.lines[i]
			compactedIndex++
		}
	}

	// Update the buffer with compacted data
	b.lines = compacted
	b.currentIndex = compactedIndex
	b.isFull = false

	log.Logger().Debugf("Buffer compacted, new size: %d", compactedIndex)
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

// GetDetailedStats returns detailed buffer statistics including event type counts
func (b *Buffer) GetDetailedStats() map[string]interface{} {
	b.mu.RLock()
	defer b.mu.RUnlock()

	stats := map[string]interface{}{
		"max_lines":     b.maxLines,
		"current_index": b.currentIndex,
		"is_full":       b.isFull,
	}

	// Count events by type
	eventTypeCounts := make(map[string]int)
	totalLines := 0

	if b.isFull {
		// Count all lines in the buffer
		for i := 0; i < b.maxLines; i++ {
			if b.lines[i].Data != "" {
				eventTypeCounts[b.lines[i].EventType]++
				totalLines++
			}
		}
	} else {
		// Count lines from beginning to current position
		for i := 0; i < b.currentIndex; i++ {
			if b.lines[i].Data != "" {
				eventTypeCounts[b.lines[i].EventType]++
				totalLines++
			}
		}
	}

	stats["total_lines"] = totalLines
	stats["event_type_counts"] = eventTypeCounts

	return stats
}

// Clear clears the entire buffer and resets its state
func (b *Buffer) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Clear all lines in the buffer
	for i := 0; i < b.maxLines; i++ {
		b.lines[i] = BufferLine{}
	}

	// Reset buffer state
	b.currentIndex = 0
	b.isFull = false

	log.Logger().Debugf("Buffer manually cleared")
}
