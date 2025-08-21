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
}

// TelemetryService implements the TelemetryProcessor interface
type TelemetryService struct {
	buffer map[string]*Buffer
	mu     sync.RWMutex
}

// NewTelemetryService creates a new telemetry service
func NewTelemetryService(buffer map[string]*Buffer) *TelemetryService {
	return &TelemetryService{
		buffer: buffer,
	}
}

// ProcessLogs processes log data and adds it to the buffer
func (s *TelemetryService) ProcessLogs(ctx context.Context, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.buffer["logs"] == nil {
		s.buffer["logs"] = &Buffer{}
	}

	// Validate input
	if len(data) == 0 {
		return fmt.Errorf("empty log data")
	}

	// Process the logs and add to buffer
	s.buffer["logs"].AddLineWithType(data)

	// log.Logger().Debugf("Processed %d bytes of log data", len(data))
	return nil
}

// ProcessMetrics processes metrics data and adds it to the buffer
func (s *TelemetryService) ProcessMetrics(ctx context.Context, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.buffer["metrics"] == nil {
		s.buffer["metrics"] = &Buffer{}
	}

	// Validate input
	if len(data) == 0 {
		return fmt.Errorf("empty metrics data")
	}

	log.Logger().Debugf("Processing %d bytes of metrics data", len(data))

	// Process the metrics and add to buffer with metrics event type
	s.buffer["metrics"].AddLineWithType(data)

	log.Logger().Debugf("Processed %d bytes of metrics data", len(data))
	return nil
}

// ProcessTraces processes trace data and adds it to the buffer
func (s *TelemetryService) ProcessTraces(ctx context.Context, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.buffer["traces"] == nil {
		s.buffer["traces"] = &Buffer{}
	}

	// Validate input
	if len(data) == 0 {
		return fmt.Errorf("empty trace data")
	}

	// Process the traces and add to buffer with traces event type
	s.buffer["traces"].AddLineWithType(data)

	// log.Logger().Debugf("Processed %d bytes of trace data", len(data))
	return nil
}
