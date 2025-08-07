# OTEL Receiver Telemetry System

This package implements a new telemetry system that uses an OpenTelemetry (OTEL) receiver instead of file-based reading. The system consists of several components working together to collect, buffer, and upload telemetry data.

## Architecture

### Components

1. **Buffer** (`buffer.go`)
   - Maintains a circular buffer of ~1000 lines (configurable)
   - Automatically sends data to upload channel when buffer is full
   - Provides periodic flush functionality

2. **Receiver** (`receiver.go`)
   - HTTP server listening on specified port (default: 33931)
   - Handles OTEL protocol endpoints:
     - `/v1/logs` - for log data
     - `/v1/metrics` - for metrics data
     - `/v1/traces` - for trace data
     - `/health` - health check endpoint

3. **S3Uploader** (`s3_uploader.go`)
   - Handles uploading telemetry data to S3 using presigned URLs
   - Manages presigned URL refresh
   - Processes upload requests asynchronously
   - Organizes uploads by event type with appropriate suffixes

4. **TelemetryManager** (`pkg/telemetry/manager.go`)
   - Coordinates all components
   - Manages lifecycle (start/stop)
   - Provides statistics and monitoring
   - Starts and manages OTEL collector process

5. **OTEL Collector** (integrated)
   - Runs the OpenTelemetry Collector binary
   - Processes system logs and metrics
   - Sends data to the HTTP receiver endpoints

## Configuration

The system is configured through the `TelemetrySettings` in the main application:

```json
{
  "telemetry": {
    "enabled": true,
    "port": 33931,
    "syslog_number_of_lines_to_read": 1000,
    "push_frequency": "60s",
    "base_directory": "/runner/warpbuild-agent"
  }
}
```

## Usage

### Starting the System

```go
// Create telemetry manager
manager := telemetry.NewTelemetryManager(ctx, port, maxBufferSize, baseDirectory, warpbuildAPI, runnerID, pollingSecret, hostURL)

// Start the system
if err := manager.Start(); err != nil {
    log.Error("Failed to start telemetry manager:", err)
}

// ... system runs ...

// Stop the system
if err := manager.Stop(); err != nil {
    log.Error("Failed to stop telemetry manager:", err)
}
```

### Sending Data

The system accepts OTEL protocol data via HTTP POST requests. Data is uploaded without any prefixes or formatting:

```bash
# Send logs
curl -X POST http://localhost:33931/v1/logs \
  -H "Content-Type: application/json" \
  -d '{"logs": [{"body": "test log message"}]}'

# Send metrics
curl -X POST http://localhost:33931/v1/metrics \
  -H "Content-Type: application/json" \
  -d '{"metrics": [{"name": "cpu_usage", "value": 75.5}]}'

# Send traces
curl -X POST http://localhost:33931/v1/traces \
  -H "Content-Type: application/json" \
  -d '{"traces": [{"span_id": "123", "trace_id": "456"}]}'
```

### File Organization

Uploaded files are organized by event type with suffixes:
- Logs: `YYYYMMDD-HHMMSS.logs.log`
- Metrics: `YYYYMMDD-HHMMSS.metrics.log`
- Traces: `YYYYMMDD-HHMMSS.traces.log`

## Data Flow

1. **System Collection**: OTEL collector processes system logs and metrics
2. **HTTP Ingestion**: OTEL data arrives via HTTP endpoints
3. **Processing**: Data is processed without any prefixes
4. **Buffering**: Data is added to circular buffer
5. **Upload**: When buffer is full or periodically, data is sent to S3
6. **Storage**: Data is uploaded to S3 with event-specific filenames (logs, metrics, traces)

## Monitoring

The system provides statistics through the `GetStats()` method:

```go
stats := manager.GetStats()
// Returns map with:
// - is_running: bool
// - port: int
// - buffer_size: int
// - current_buffer_size: int
// - buffer_is_full: bool
// - receiver_running: bool
// - s3_uploader_running: bool
```

## Benefits

1. **Real-time**: Data is processed as it arrives, not from file reads
2. **Scalable**: Buffer system prevents memory issues
3. **Reliable**: Automatic retry and error handling
4. **Standards-compliant**: Uses OTEL protocol
5. **Configurable**: Port, buffer size, and other parameters are configurable
6. **Event-specific**: Uploads are organized by event type (logs, metrics, traces)
7. **Clean data**: No prefixes or formatting added to the original data

## Migration from File-based System

The new system replaces the previous file-based telemetry system that:
- Read from syslog files
- Used OTEL collector binaries
- Watched for file changes
- Required complex file management

The new system is simpler, more efficient, and provides better real-time capabilities. 