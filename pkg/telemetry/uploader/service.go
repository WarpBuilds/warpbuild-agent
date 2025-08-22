package uploader

import (
	"context"
	"fmt"
	"sync"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
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
	// Store API client and other required fields for creating buffers
	maxLines      int
	warpbuildAPI  *warpbuild.APIClient
	runnerID      string
	pollingSecret string
	hostURL       string
}

// NewTelemetryService creates a new telemetry service
func NewTelemetryService(warpbuildAPI *warpbuild.APIClient, runnerID, pollingSecret, hostURL string, maxLines int) *TelemetryService {
	return &TelemetryService{
		buffer:        map[string]*Buffer{},
		maxLines:      maxLines,
		warpbuildAPI:  warpbuildAPI,
		runnerID:      runnerID,
		pollingSecret: pollingSecret,
		hostURL:       hostURL,
	}
}

// createBufferIfNil creates a properly initialized buffer using NewBuffer if it doesn't exist
func (s *TelemetryService) createBufferIfNil(eventType string) error {
	if s.buffer[eventType] == nil {
		log.Logger().Debugf("Creating new buffer for event type: %s", eventType)
		ctx := context.Background()
		buffer, err := NewBuffer(ctx, BufferOptions{
			WarpbuildAPI:  s.warpbuildAPI,
			EventType:     eventType,
			RunnerID:      s.runnerID,
			PollingSecret: s.pollingSecret,
			HostURL:       s.hostURL,
			MaxLines:      s.maxLines, // Default max lines
		})
		if err != nil {
			log.Logger().Errorf("Failed to create buffer for %s: %v", eventType, err)
			return fmt.Errorf("failed to create buffer for %s: %w", eventType, err)
		}
		s.buffer[eventType] = buffer
		log.Logger().Debugf("Successfully created buffer for event type: %s", eventType)
	} else {
		log.Logger().Debugf("Buffer already exists for event type: %s", eventType)
	}
	return nil
}

// ProcessLogs processes log data and adds it to the buffer
func (s *TelemetryService) ProcessLogs(ctx context.Context, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.createBufferIfNil("logs"); err != nil {
		return fmt.Errorf("failed to create logs buffer: %w", err)
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

	log.Logger().Debugf("ProcessMetrics called with %d bytes of data", len(data))

	if err := s.createBufferIfNil("metrics"); err != nil {
		log.Logger().Errorf("Failed to create metrics buffer: %v", err)
		return fmt.Errorf("failed to create metrics buffer: %w", err)
	}

	// Validate input
	if len(data) == 0 {
		log.Logger().Warnf("Empty metrics data received")
		return fmt.Errorf("empty metrics data")
	}

	log.Logger().Debugf("Processing %d bytes of metrics data", len(data))

	// Process the metrics and add to buffer with metrics event type
	s.buffer["metrics"].AddLineWithType(data)

	log.Logger().Debugf("Successfully processed %d bytes of metrics data", len(data))
	return nil
}

// ProcessTraces processes trace data and adds it to the buffer
func (s *TelemetryService) ProcessTraces(ctx context.Context, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.createBufferIfNil("traces"); err != nil {
		return fmt.Errorf("failed to create traces buffer: %w", err)
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
