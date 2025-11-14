package uploader

import (
	"context"
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

// createTestBuffer creates a minimal buffer for testing without S3Uploader
func createTestBuffer() *Buffer {
	ctx, cancel := context.WithCancel(context.Background())
	return &Buffer{
		eventType:  "test",
		uploadChan: make(chan UploadRequest, 100),
		ctx:        ctx,
		cancel:     cancel,
	}
}

func TestBufferAddLine(t *testing.T) {
	buffer := createTestBuffer()

	// Start a consumer goroutine to receive data
	received := make(chan []byte, 1)
	go func() {
		req := <-buffer.uploadChan
		received <- req.Data
	}()

	// Add a line - should immediately send to uploader
	buffer.AddLineWithType([]byte("test line"))

	// Check that data was sent immediately
	select {
	case data := <-received:
		expected := "test line"
		if string(data) != expected {
			t.Errorf("Expected sent data '%s', got '%s'", expected, string(data))
		}
	case <-buffer.ctx.Done():
		t.Error("Context was cancelled unexpectedly")
	}
}

func TestBufferEmptyData(t *testing.T) {
	buffer := createTestBuffer()

	// Try to add empty data - should be skipped
	buffer.AddLineWithType([]byte(""))

	// Should not panic or hang
	// No data should be sent to upload channel
}

func TestBufferMultipleAdds(t *testing.T) {
	buffer := createTestBuffer()

	// Start a consumer goroutine to receive data
	received := make(chan []byte, 3)
	go func() {
		for i := 0; i < 3; i++ {
			req := <-buffer.uploadChan
			received <- req.Data
		}
	}()

	// Add multiple lines - each should be sent immediately
	lines := []string{"line1", "line2", "line3"}
	for _, line := range lines {
		buffer.AddLineWithType([]byte(line))
	}

	// Check that all data was sent immediately
	for i, expected := range lines {
		select {
		case data := <-received:
			if string(data) != expected {
				t.Errorf("Expected sent data '%s' at index %d, got '%s'", expected, i, string(data))
			}
		case <-buffer.ctx.Done():
			t.Error("Context was cancelled unexpectedly")
		}
	}
}

func TestBufferContextCancellation(t *testing.T) {
	buffer := createTestBuffer()

	// Cancel the context
	buffer.cancel()

	// Try to add data - should not send due to context cancellation
	buffer.AddLineWithType([]byte("test line"))

	// Test should complete without hanging
}

func TestBufferUploadChannelFull(t *testing.T) {
	buffer := createTestBuffer()

	// Fill the upload channel
	for i := 0; i < 100; i++ {
		select {
		case buffer.uploadChan <- UploadRequest{Data: []byte("test")}:
		default:
			break
		}
	}

	// Try to add data - should not block due to full channel
	buffer.AddLineWithType([]byte("test line"))

	// Test should complete without hanging
}

func TestBufferNilData(t *testing.T) {
	buffer := createTestBuffer()

	// Try to add nil data - should not panic
	buffer.AddLineWithType(nil)

	// Test should complete without panic
}
