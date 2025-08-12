package uploader

import (
	"testing"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

func init() {
	// Initialize logger for tests
	_, err := log.Init(&log.InitOptions{})
	if err != nil {
		panic("Failed to initialize logger for tests")
	}
}

func TestBufferCleanupAfterUpload(t *testing.T) {
	// Create a small buffer for testing with a large upload channel
	uploadChan := make(chan UploadRequest, 100)
	buffer := NewBuffer(5, uploadChan)

	// Start a consumer goroutine to consume from the upload channel
	go func() {
		for range uploadChan {
			// Consume all upload requests
		}
	}()

	// Add some lines of different event types
	buffer.AddLineWithType("log1", "logs")
	buffer.AddLineWithType("metric1", "metrics")
	buffer.AddLineWithType("log2", "logs")
	buffer.AddLineWithType("trace1", "traces")
	buffer.AddLineWithType("log3", "logs")

	// Buffer should be full now
	if !buffer.isFull {
		t.Error("Buffer should be full after adding 5 lines to a 5-line buffer")
	}

	// Get initial stats
	initialSize, initialFull := buffer.GetStats()
	if initialSize != 5 || !initialFull {
		t.Errorf("Expected buffer size 5 and full=true, got size=%d, full=%v", initialSize, initialFull)
	}

	// Check initial event type counts
	initialStats := buffer.GetDetailedStats()
	initialEventTypeCounts := initialStats["event_type_counts"].(map[string]int)
	t.Logf("Initial event counts: %v", initialEventTypeCounts)

	// Simulate sending data to upload channel (this will trigger cleanup)
	buffer.sendToUploadChannel()

	// Check that all events were dropped (since they were all sent to upload channel)
	stats := buffer.GetDetailedStats()
	eventTypeCounts := stats["event_type_counts"].(map[string]int)
	t.Logf("After upload event counts: %v", eventTypeCounts)

	// All event types should be dropped (0 remaining) since they were all sent
	if eventTypeCounts["logs"] != 0 {
		t.Errorf("Expected 0 logs after upload, got %d", eventTypeCounts["logs"])
	}

	if eventTypeCounts["metrics"] != 0 {
		t.Errorf("Expected 0 metrics after upload, got %d", eventTypeCounts["metrics"])
	}

	if eventTypeCounts["traces"] != 0 {
		t.Errorf("Expected 0 traces after upload, got %d", eventTypeCounts["traces"])
	}

	// Buffer should not be full anymore since all events were dropped
	if buffer.isFull {
		t.Error("Buffer should not be full after dropping all events")
	}

	// Buffer should be empty
	if buffer.currentIndex != 0 {
		t.Error("Buffer current index should be 0 after dropping all events")
	}
}

func TestBufferManualClear(t *testing.T) {
	buffer := NewBuffer(5, make(chan UploadRequest, 10))

	// Add some lines
	buffer.AddLineWithType("line1", "logs")
	buffer.AddLineWithType("line2", "metrics")

	// Manually clear the buffer
	buffer.Clear()

	// Buffer should be cleared
	size, _ := buffer.GetStats()
	if size != 0 {
		t.Errorf("Expected 0 lines after manual clear, got %d", size)
	}

	if buffer.currentIndex != 0 {
		t.Error("Buffer current index should be 0 after manual clear")
	}
}

func TestDetailedStats(t *testing.T) {
	buffer := NewBuffer(5, make(chan UploadRequest, 10))

	// Add lines with different event types
	buffer.AddLineWithType("log1", "logs")
	buffer.AddLineWithType("log2", "logs")
	buffer.AddLineWithType("metric1", "metrics")

	// Get detailed stats
	stats := buffer.GetDetailedStats()

	// Check basic stats
	if stats["max_lines"] != 5 {
		t.Errorf("Expected max_lines=5, got %v", stats["max_lines"])
	}

	if stats["total_lines"] != 3 {
		t.Errorf("Expected total_lines=3, got %v", stats["total_lines"])
	}

	// Check event type counts
	eventTypeCounts := stats["event_type_counts"].(map[string]int)
	if eventTypeCounts["logs"] != 2 {
		t.Errorf("Expected 2 logs, got %d", eventTypeCounts["logs"])
	}

	if eventTypeCounts["metrics"] != 1 {
		t.Errorf("Expected 1 metric, got %d", eventTypeCounts["metrics"])
	}
}

func TestDropEventsOfType(t *testing.T) {
	buffer := NewBuffer(5, make(chan UploadRequest, 10))

	// Add mixed event types
	buffer.AddLineWithType("log1", "logs")
	buffer.AddLineWithType("metric1", "metrics")
	buffer.AddLineWithType("log2", "logs")
	buffer.AddLineWithType("trace1", "traces")

	// Drop all logs events
	buffer.dropEventsOfType("logs")

	// Check that logs are gone but others remain
	stats := buffer.GetDetailedStats()
	eventTypeCounts := stats["event_type_counts"].(map[string]int)

	if eventTypeCounts["logs"] != 0 {
		t.Errorf("Expected 0 logs after dropping, got %d", eventTypeCounts["logs"])
	}

	if eventTypeCounts["metrics"] != 1 {
		t.Errorf("Expected 1 metric after dropping logs, got %d", eventTypeCounts["metrics"])
	}

	if eventTypeCounts["traces"] != 1 {
		t.Errorf("Expected 1 trace after dropping logs, got %d", eventTypeCounts["traces"])
	}
}

func TestCompactBuffer(t *testing.T) {
	buffer := NewBuffer(5, make(chan UploadRequest, 10))

	// Add some events
	buffer.AddLineWithType("log1", "logs")
	buffer.AddLineWithType("metric1", "metrics")
	buffer.AddLineWithType("log2", "logs")

	// Drop logs events (this will create gaps)
	buffer.dropEventsOfType("logs")

	// Compact the buffer
	buffer.compactBuffer()

	// Check that buffer is compacted
	stats := buffer.GetDetailedStats()
	eventTypeCounts := stats["event_type_counts"].(map[string]int)

	if eventTypeCounts["logs"] != 0 {
		t.Errorf("Expected 0 logs after compacting, got %d", eventTypeCounts["logs"])
	}

	if eventTypeCounts["metrics"] != 1 {
		t.Errorf("Expected 1 metric after compacting, got %d", eventTypeCounts["metrics"])
	}

	// Buffer should not be full and current index should be correct
	if buffer.isFull {
		t.Error("Buffer should not be full after compacting")
	}

	if buffer.currentIndex != 1 {
		t.Errorf("Expected current index 1 after compacting, got %d", buffer.currentIndex)
	}
}
