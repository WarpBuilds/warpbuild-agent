package uploader

import (
	"bytes"
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
func createTestBuffer(maxLines int) *Buffer {
	ctx, cancel := context.WithCancel(context.Background())
	return &Buffer{
		buf:          bytes.NewBufferString(""),
		maxLines:     maxLines,
		currentIndex: 0,
		uploadChan:   make(chan UploadRequest, 100),
		ctx:          ctx,
		cancel:       cancel,
	}
}

func TestBufferAddLine(t *testing.T) {
	buffer := createTestBuffer(5)

	// Add a line
	buffer.AddLineWithType([]byte("test line"))

	// Check that line was added
	if buffer.currentIndex != 1 {
		t.Errorf("Expected current index 1, got %d", buffer.currentIndex)
	}

	// Check buffer content
	content := buffer.buf.String()
	expected := "test line\n"
	if content != expected {
		t.Errorf("Expected buffer content '%s', got '%s'", expected, content)
	}
}

func TestBufferFull(t *testing.T) {
	buffer := createTestBuffer(3)

	// Add lines until buffer is full
	buffer.AddLineWithType([]byte("line1"))
	buffer.AddLineWithType([]byte("line2"))
	buffer.AddLineWithType([]byte("line3"))

	// Buffer should be full and reset
	if buffer.currentIndex != 0 {
		t.Errorf("Expected current index 0 after buffer full, got %d", buffer.currentIndex)
	}

	// Buffer should be empty after reset
	content := buffer.buf.String()
	if content != "" {
		t.Errorf("Expected empty buffer after reset, got '%s'", content)
	}
}

func TestBufferSendToUploadChannel(t *testing.T) {
	buffer := createTestBuffer(2)

	// Add some data
	buffer.AddLineWithType([]byte("line1"))
	buffer.AddLineWithType([]byte("line2"))

	// Start a consumer goroutine
	received := make(chan []byte, 1)
	go func() {
		req := <-buffer.uploadChan
		received <- req.Data
	}()

	// Manually trigger send to upload channel
	buffer.sendToUploadChannel()

	// Check that data was sent
	select {
	case data := <-received:
		expected := "line1\nline2\n"
		if string(data) != expected {
			t.Errorf("Expected sent data '%s', got '%s'", expected, string(data))
		}
	case <-buffer.ctx.Done():
		t.Error("Context was cancelled unexpectedly")
	}
}

func TestBufferContextCancellation(t *testing.T) {
	buffer := createTestBuffer(5)

	// Add some data
	buffer.AddLineWithType([]byte("line1"))

	// Cancel the context
	buffer.cancel()

	// Try to send to upload channel (should not block due to context cancellation)
	buffer.sendToUploadChannel()

	// Test should complete without hanging
}

func TestBufferMultipleAdds(t *testing.T) {
	buffer := createTestBuffer(10)

	// Add multiple lines
	lines := []string{"line1", "line2", "line3", "line4", "line5"}
	for _, line := range lines {
		buffer.AddLineWithType([]byte(line))
	}

	// Check current index
	if buffer.currentIndex != 5 {
		t.Errorf("Expected current index 5, got %d", buffer.currentIndex)
	}

	// Check buffer content
	content := buffer.buf.String()
	expected := "line1\nline2\nline3\nline4\nline5\n"
	if content != expected {
		t.Errorf("Expected buffer content '%s', got '%s'", expected, content)
	}
}

func TestBufferEmptySend(t *testing.T) {
	buffer := createTestBuffer(5)

	// Try to send empty buffer
	buffer.sendToUploadChannel()

	// Should not panic or hang
	// Buffer should remain empty
	if buffer.currentIndex != 0 {
		t.Errorf("Expected current index 0, got %d", buffer.currentIndex)
	}
}

func TestBufferMaxLines(t *testing.T) {
	buffer := createTestBuffer(1)

	// Add one line
	buffer.AddLineWithType([]byte("single line"))

	// Buffer should be full and reset
	if buffer.currentIndex != 0 {
		t.Errorf("Expected current index 0 after single line in 1-line buffer, got %d", buffer.currentIndex)
	}

	// Buffer should be empty
	content := buffer.buf.String()
	if content != "" {
		t.Errorf("Expected empty buffer after reset, got '%s'", content)
	}
}

func TestBufferNilBufRecovery(t *testing.T) {
	// Create a buffer with nil buf to test recovery
	ctx, cancel := context.WithCancel(context.Background())
	buffer := &Buffer{
		buf:          nil, // Intentionally nil to test recovery
		maxLines:     5,
		currentIndex: 0,
		uploadChan:   make(chan UploadRequest, 100),
		ctx:          ctx,
		cancel:       cancel,
		eventType:    "test",
	}

	// This should not panic and should recover by creating a new buf
	buffer.AddLineWithType([]byte("recovery test"))

	// Check that buf was created and line was added
	if buffer.buf == nil {
		t.Error("Buffer buf should not be nil after recovery")
	}

	if buffer.currentIndex != 1 {
		t.Errorf("Expected current index 1 after recovery, got %d", buffer.currentIndex)
	}

	// Check buffer content
	content := buffer.buf.String()
	expected := "recovery test\n"
	if content != expected {
		t.Errorf("Expected buffer content '%s', got '%s'", expected, content)
	}

	cancel()
}
