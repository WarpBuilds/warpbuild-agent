package uploader

import (
	"context"
	"fmt"
	"sync"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

// TelemetryProcessor defines the interface for processing telemetry data
type TelemetryProcessor interface {
	ProcessLogs(ctx context.Context, data []byte) error
	ProcessMetrics(ctx context.Context, data []byte) error
	ProcessTraces(ctx context.Context, data []byte) error
	GetStats() map[string]interface{}
}

// TelemetryService implements the TelemetryProcessor interface
type TelemetryService struct {
	buffer     *Buffer
	s3Uploader *S3Uploader
	mu         sync.RWMutex
	stats      map[string]interface{}
}

// NewTelemetryService creates a new telemetry service
func NewTelemetryService(buffer *Buffer, s3Uploader *S3Uploader) *TelemetryService {
	return &TelemetryService{
		buffer:     buffer,
		s3Uploader: s3Uploader,
		stats: map[string]interface{}{
			"logs_processed":    0,
			"metrics_processed": 0,
			"traces_processed":  0,
			"errors":            0,
		},
	}
}

// ProcessLogs processes log data and adds it to the buffer
func (s *TelemetryService) ProcessLogs(ctx context.Context, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.buffer == nil {
		return fmt.Errorf("buffer is not initialized")
	}

	// Validate input
	if len(data) == 0 {
		return fmt.Errorf("empty log data")
	}

	// Process the logs and add to buffer
	s.buffer.AddLine(string(data))

	// Update statistics
	if count, ok := s.stats["logs_processed"].(int); ok {
		s.stats["logs_processed"] = count + 1
	}

	// log.Logger().Debugf("Processed %d bytes of log data", len(data))
	return nil
}

// ProcessMetrics processes metrics data and adds it to the buffer
func (s *TelemetryService) ProcessMetrics(ctx context.Context, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.buffer == nil {
		return fmt.Errorf("buffer is not initialized")
	}

	// Validate input
	if len(data) == 0 {
		return fmt.Errorf("empty metrics data")
	}

	log.Logger().Debugf("Processing %d bytes of metrics data", len(data))

	// Process the metrics and add to buffer with metrics event type
	s.buffer.AddLineWithType(string(data), "metrics")

	// Update statistics
	if count, ok := s.stats["metrics_processed"].(int); ok {
		s.stats["metrics_processed"] = count + 1
	}

	log.Logger().Debugf("Processed %d bytes of metrics data", len(data))
	return nil
}

// ProcessTraces processes trace data and adds it to the buffer
func (s *TelemetryService) ProcessTraces(ctx context.Context, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.buffer == nil {
		return fmt.Errorf("buffer is not initialized")
	}

	// Validate input
	if len(data) == 0 {
		return fmt.Errorf("empty trace data")
	}

	// Process the traces and add to buffer
	s.buffer.AddLine(string(data))

	// Update statistics
	if count, ok := s.stats["traces_processed"].(int); ok {
		s.stats["traces_processed"] = count + 1
	}

	// log.Logger().Debugf("Processed %d bytes of trace data", len(data))
	return nil
}

// GetStats returns the current statistics
func (s *TelemetryService) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create a copy of stats to avoid race conditions
	stats := make(map[string]interface{})
	for k, v := range s.stats {
		stats[k] = v
	}

	// Add buffer stats if available
	if s.buffer != nil {
		bufferSize, isFull := s.buffer.GetStats()
		stats["buffer_size"] = bufferSize
		stats["buffer_is_full"] = isFull
	}

	return stats
}

// RecordError records an error in the statistics
func (s *TelemetryService) RecordError() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if count, ok := s.stats["errors"].(int); ok {
		s.stats["errors"] = count + 1
	}
}

// ResetStats resets all statistics
func (s *TelemetryService) ResetStats() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stats = map[string]interface{}{
		"logs_processed":    0,
		"metrics_processed": 0,
		"traces_processed":  0,
		"errors":            0,
	}
}
