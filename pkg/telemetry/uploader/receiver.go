package uploader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

// Receiver handles incoming OTEL telemetry data
type Receiver struct {
	port    int
	service TelemetryProcessor
	server  *http.Server
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	mu      sync.RWMutex
}

// NewReceiver creates a new OTEL receiver
func NewReceiver(port int, service TelemetryProcessor) *Receiver {
	ctx, cancel := context.WithCancel(context.Background())
	return &Receiver{
		port:    port,
		service: service,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// loggingMiddleware logs request details and body content
func (r *Receiver) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Log request details
		log.Logger().Infof("Request received: %s %s from %s", req.Method, req.URL.Path, req.RemoteAddr)

		// Only log body for POST requests to telemetry endpoints
		if req.Method == http.MethodPost && (req.URL.Path == "/v1/logs" || req.URL.Path == "/v1/metrics" || req.URL.Path == "/v1/traces") {
			// Read the body
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				log.Logger().Errorf("Failed to read request body for logging: %v", err)
				http.Error(w, "Failed to read request body", http.StatusBadRequest)
				return
			}

			// Log body content (truncate if too long)
			bodyStr := string(bodyBytes)
			if len(bodyStr) > 1000 {
				bodyStr = bodyStr[:1000] + "... [truncated]"
			}
			log.Logger().Infof("Request body (%d bytes): %s", len(bodyBytes), bodyStr)

			// Restore the body for the actual handler
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Call the next handler
		next(w, req)
	}
}

// Start starts the receiver on the specified port
func (r *Receiver) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	mux := http.NewServeMux()

	// Handle OTEL logs endpoint with logging middleware
	mux.HandleFunc("/v1/logs", r.loggingMiddleware(r.handleLogs))

	// Handle OTEL metrics endpoint with logging middleware
	mux.HandleFunc("/v1/metrics", r.loggingMiddleware(r.handleMetrics))

	// Handle OTEL traces endpoint with logging middleware
	mux.HandleFunc("/v1/traces", r.loggingMiddleware(r.handleTraces))

	// Health check endpoint (no middleware needed)
	mux.HandleFunc("/health", r.handleHealth)

	r.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", r.port),
		Handler: mux,
	}

	// Start the server in a goroutine
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		log.Logger().Infof("Starting OTEL receiver on port %d", r.port)

		if err := r.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Logger().Errorf("OTEL receiver error: %v", err)
		}
	}()

	// Wait a moment to ensure server is starting
	time.Sleep(100 * time.Millisecond)

	// Check if server is actually listening
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", r.port), time.Second)
	if err != nil {
		return fmt.Errorf("failed to start receiver: %w", err)
	}
	conn.Close()

	log.Logger().Infof("OTEL receiver started successfully on port %d", r.port)
	return nil
}

// Stop stops the receiver
func (r *Receiver) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	log.Logger().Infof("Stopping OTEL receiver...")

	r.cancel()

	if r.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := r.server.Shutdown(ctx); err != nil {
			log.Logger().Errorf("Error shutting down server: %v", err)
		}
	}

	r.wg.Wait()

	log.Logger().Infof("OTEL receiver stopped")
	return nil
}

// handleLogs handles incoming log data
func (r *Receiver) handleLogs(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := r.readRequestBody(req)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Process the logs and add to buffer
	r.processLogs(body)

	w.WriteHeader(http.StatusOK)
}

// handleMetrics handles incoming metrics data
func (r *Receiver) handleMetrics(w http.ResponseWriter, req *http.Request) {
	log.Logger().Debugf("handleMetrics called from %s with method %s", req.RemoteAddr, req.Method)

	if req.Method != http.MethodPost {
		log.Logger().Warnf("Invalid method %s for metrics endpoint", req.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := r.readRequestBody(req)
	if err != nil {
		log.Logger().Errorf("Failed to read metrics request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	log.Logger().Debugf("Received metrics request with %d bytes from %s", len(body), req.RemoteAddr)

	// Process the metrics and add to buffer
	r.processMetrics(body)

	log.Logger().Debugf("Metrics processing completed for %d bytes", len(body))
	w.WriteHeader(http.StatusOK)
}

// handleTraces handles incoming trace data
func (r *Receiver) handleTraces(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := r.readRequestBody(req)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Process the traces and add to buffer
	r.processTraces(body)

	w.WriteHeader(http.StatusOK)
}

// handleHealth handles health check requests
func (r *Receiver) handleHealth(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

// readRequestBody reads and returns the request body
func (r *Receiver) readRequestBody(req *http.Request) ([]byte, error) {
	// Limit body size to prevent memory issues
	req.Body = http.MaxBytesReader(nil, req.Body, 10*1024*1024) // 10MB limit

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	return body, nil
}

// processLogs processes log data using the service layer
func (r *Receiver) processLogs(data []byte) {
	if err := r.service.ProcessLogs(r.ctx, data); err != nil {
		log.Logger().Errorf("Failed to process logs: %v", err)
	}
}

// processMetrics processes metrics data using the service layer
func (r *Receiver) processMetrics(data []byte) {
	if err := r.service.ProcessMetrics(r.ctx, data); err != nil {
		log.Logger().Errorf("Failed to process metrics: %v", err)
	}
}

// processTraces processes trace data using the service layer
func (r *Receiver) processTraces(data []byte) {
	if err := r.service.ProcessTraces(r.ctx, data); err != nil {
		log.Logger().Errorf("Failed to process traces: %v", err)
	}
}
