package asur

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

const azureVersion = "2020-04-08" // cosmetic response header

var debugMode = os.Getenv("AZPROXY_DEBUG") == "true"

// -------- Azure blocklist XML --------
type blockListXML struct {
	XMLName     xml.Name `xml:"BlockList"`
	Latest      []string `xml:"Latest"`
	Uncommitted []string `xml:"Uncommitted"`
	Committed   []string `xml:"Committed"`
}

// -------- in-memory multipart state --------
type uploadedPart struct {
	PartNumber int32
	ETag       string
}

type uploadState struct {
	UploadID string
	Parts    map[string]uploadedPart // blockID (base64) -> part info
}

// stateKey generates a unique key for the upload state map
func stateKey(bucket, key string) string {
	return bucket + "/" + key
}

// getUploadState retrieves upload state from memory
func (s *server) getUploadState(bucket, key string) *uploadState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.uploadStates[stateKey(bucket, key)]
}

// ensureMultipart ensures a multipart upload exists, creating if necessary
func (s *server) ensureMultipart(ctx context.Context, bucket, key string, r *http.Request) (*uploadState, error) {
	stateMapKey := stateKey(bucket, key)

	// Hold the mutex while checking and possibly creating to avoid races where
	// multiple concurrent requests all create independent multipart uploads.
	s.mu.Lock()
	us := s.uploadStates[stateMapKey]
	if us != nil {
		s.mu.Unlock()
		return us, nil
	}

	// Create new multipart upload while still holding the lock so only one is created.
	uploadID, err := s.uploader.EnsureMultipartUpload(ctx, bucket, key)
	if err != nil {
		s.mu.Unlock()
		return nil, err
	}

	us = &uploadState{
		UploadID: uploadID,
		Parts:    map[string]uploadedPart{},
	}

	s.uploadStates[stateMapKey] = us
	s.mu.Unlock()

	return us, nil
}

// deleteUploadState removes upload state from memory
func (s *server) deleteUploadState(bucket, key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.uploadStates, stateKey(bucket, key))
}

// -------- utility --------
func azureOKHeaders(w http.ResponseWriter) {
	w.Header().Set("x-ms-version", azureVersion)
	w.Header().Set("x-ms-request-id", randomReqID())
	w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
}

func randomReqID() string {
	var b [16]byte
	rand.Read(b[:])
	return fmt.Sprintf("%x", b[:])
}

// extractBlockNumber extracts the block index from an Azure block ID
// Azure block IDs are base64-encoded strings containing a UUID prefix and padded block number
func extractBlockNumber(blockIDB64 string) (int, error) {
	decoded, err := base64.StdEncoding.DecodeString(blockIDB64)
	if err != nil {
		return 0, err
	}

	decodedStr := string(decoded)
	// Azure uses 6 digits for block number at the end
	if len(decodedStr) < 6 {
		return 0, fmt.Errorf("block ID too short")
	}

	// Extract last 6 characters as block number
	blockNumStr := decodedStr[len(decodedStr)-6:]
	blockNum, err := strconv.Atoi(blockNumStr)
	if err != nil {
		return 0, fmt.Errorf("invalid block number: %v", err)
	}

	return blockNum, nil
}

// parseTotalFromContentRange extracts the total length from a Content-Range header like "bytes start-end/total"
func parseTotalFromContentRange(cr string) (int64, bool) {
	_, after, ok := strings.Cut(cr, "/")
	if !ok || after == "*" || after == "" {
		return 0, false
	}
	// Trim any trailing non-digit chars just in case
	i := 0
	for i < len(after) && after[i] >= '0' && after[i] <= '9' {
		i++
	}
	after = after[:i]
	total, err := strconv.ParseInt(after, 10, 64)
	if err != nil {
		return 0, false
	}
	return total, true
}

// logError logs errors with context
func logError(operation string, err error, context map[string]interface{}) {
	if err == nil {
		return
	}

	log.Printf("ERROR in %s: %v", operation, err)

	if debugMode && len(context) > 0 {
		if jsonBytes, err := json.MarshalIndent(context, "", "  "); err == nil {
			log.Printf("Context:\n%s", string(jsonBytes))
		}
	}
}

// -------- server --------
type server struct {
	s3 *s3.Client
	// s3NoRetry disables automatic retries to avoid rewind requirement on streaming bodies
	s3NoRetry *s3.Client
	// mutex for thread-safe operations
	mu sync.Mutex
	// in-memory upload states: "bucket/key" -> uploadState
	uploadStates map[string]*uploadState
	// performance statistics
	stats *performanceStats
	// uploader strategy
	uploader Uploader
}

// performanceStats tracks upload performance metrics
type performanceStats struct {
	mu              sync.RWMutex
	totalUploads    int64
	totalBytes      int64
	totalDuration   time.Duration
	activeUploads   int32
	peakConcurrency int32
	uploadHistory   []uploadStat // Ring buffer of recent uploads
	historyIndex    int
}

type uploadStat struct {
	timestamp   time.Time
	size        int64
	duration    time.Duration
	throughput  float64 // MB/s
	partCount   int
	concurrency int
}

func newPerformanceStats() *performanceStats {
	return &performanceStats{
		uploadHistory: make([]uploadStat, 100), // Keep last 100 uploads
	}
}

func (ps *performanceStats) recordUpload(size int64, duration time.Duration, partCount, concurrency int) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.totalUploads++
	ps.totalBytes += size
	ps.totalDuration += duration

	throughput := float64(size) / (1024 * 1024) / duration.Seconds()

	// Add to ring buffer
	ps.uploadHistory[ps.historyIndex] = uploadStat{
		timestamp:   time.Now(),
		size:        size,
		duration:    duration,
		throughput:  throughput,
		partCount:   partCount,
		concurrency: concurrency,
	}
	ps.historyIndex = (ps.historyIndex + 1) % len(ps.uploadHistory)
}

func (ps *performanceStats) incrementActive() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.activeUploads++
	if ps.activeUploads > ps.peakConcurrency {
		ps.peakConcurrency = ps.activeUploads
	}
}

func (ps *performanceStats) decrementActive() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.activeUploads--
}

func (ps *performanceStats) getStats() map[string]interface{} {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	avgThroughput := float64(0)
	if ps.totalDuration > 0 {
		avgThroughput = float64(ps.totalBytes) / (1024 * 1024) / ps.totalDuration.Seconds()
	}

	// Calculate recent performance (last 10 uploads)
	recentStats := make([]map[string]interface{}, 0)
	count := 0
	for i := 0; i < len(ps.uploadHistory) && count < 10; i++ {
		idx := (ps.historyIndex - 1 - i + len(ps.uploadHistory)) % len(ps.uploadHistory)
		stat := ps.uploadHistory[idx]
		if stat.timestamp.IsZero() {
			continue
		}
		recentStats = append(recentStats, map[string]interface{}{
			"timestamp":       stat.timestamp.Format(time.RFC3339),
			"size_mb":         float64(stat.size) / (1024 * 1024),
			"duration_s":      stat.duration.Seconds(),
			"throughput_mbps": stat.throughput,
			"parts":           stat.partCount,
			"concurrency":     stat.concurrency,
		})
		count++
	}

	return map[string]interface{}{
		"total_uploads":       ps.totalUploads,
		"total_bytes_mb":      float64(ps.totalBytes) / (1024 * 1024),
		"avg_throughput_mbps": avgThroughput,
		"active_uploads":      ps.activeUploads,
		"peak_concurrency":    ps.peakConcurrency,
		"recent_uploads":      recentStats,
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Handle debug endpoints
	if r.URL.Path == "/_debug/health" {
		s.handleDebugHealth(w, r)
		return
	}
	if r.URL.Path == "/_debug/stats" {
		s.handleDebugStats(w, r)
		return
	}

	ctx := r.Context()
	account, container, blobKey, err := parseAzuriteStyle(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid azure path", http.StatusBadRequest)
		return
	}
	bucket := bucketFor(account, container)

	switch r.Method {
	case http.MethodPut:
		q := r.URL.Query()
		switch strings.ToLower(q.Get("comp")) {
		case "block":
			s.handlePutBlock(ctx, w, r, bucket, blobKey, q.Get("blockid"))
		case "blocklist":
			s.handlePutBlockList(ctx, w, r, bucket, blobKey)
		default:
			s.handlePutBlob(ctx, w, r, bucket, blobKey)
		}
	case http.MethodGet:
		s.handleGet(ctx, w, r, bucket, blobKey)
	case http.MethodHead:
		s.handleHead(ctx, w, r, bucket, blobKey)
	default:
		http.Error(w, "unsupported", http.StatusNotImplemented)
	}
}

// loggingMiddleware wraps an http.Handler to log all requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Construct the full URL
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}

		// Build query string with ? prefix if present
		queryString := ""
		if r.URL.RawQuery != "" {
			queryString = "?" + r.URL.RawQuery
		}

		fullURL := fmt.Sprintf("%s://%s%s%s", scheme, r.Host, r.URL.Path, queryString)

		// Log request
		if debugMode {
			// In debug mode, format headers nicely on separate lines
			log.Printf("[%s] Exact URL: %s", r.Method, fullURL)
			log.Printf("  Headers:")
			for key, values := range r.Header {
				log.Printf("    %s: %s", key, strings.Join(values, ", "))
			}
		} else {
			// Format headers for single-line logging
			headers := make([]string, 0, len(r.Header))
			for key, values := range r.Header {
				headers = append(headers, fmt.Sprintf("%s: %s", key, strings.Join(values, ", ")))
			}
			headerStr := strings.Join(headers, " | ")

			// Log request with exact URL and headers
			if r.Header.Get("Range") != "" || r.Header.Get("x-ms-range") != "" {
				// Log range requests specially
				log.Printf("[%s] Exact URL: %s | Range: %s | Headers: {%s}",
					r.Method, fullURL,
					r.Header.Get("Range")+r.Header.Get("x-ms-range"),
					headerStr)
			} else {
				log.Printf("[%s] Exact URL: %s | Headers: {%s}", r.Method, fullURL, headerStr)
			}
		}

		// Call the next handler
		next.ServeHTTP(wrapped, r)

		// Log response with duration
		duration := time.Since(start)
		if debugMode {
			log.Printf("[%s] Response: %s - Status: %d - Duration: %v",
				r.Method, fullURL,
				wrapped.statusCode, duration)
			log.Printf("  Response Headers:")
			for key, values := range wrapped.Header() {
				log.Printf("    %s: %s", key, strings.Join(values, ", "))
			}
		} else {
			log.Printf("[%s] Exact URL: %s - Status: %d - Duration: %v",
				r.Method, fullURL,
				wrapped.statusCode, duration)
		}
	})
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (s *server) handlePutBlob(ctx context.Context, w http.ResponseWriter, r *http.Request, bucket, key string) {
	contentLength := r.ContentLength
	log.Printf("handlePutBlob: key=%s size=%d MB", key, contentLength/(1024*1024))

	// For small files (< 5MB), use direct upload
	if contentLength > 0 && contentLength < 5*1024*1024 {
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := s.uploader.UploadSmallFile(ctx, bucket, key, buf)
		if err != nil {
			log.Printf("UploadSmallFile failed: %v", err)
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		azureOKHeaders(w)
		w.Header().Set("ETag", result.ETag)
		w.Header().Set("Last-Modified", result.LastModified.Format(http.TimeFormat))
		w.WriteHeader(http.StatusCreated)
		return
	}

	// For larger files, use concurrent multipart upload
	log.Printf("Using multipart upload for large file: %s", key)

	result, err := s.uploader.UploadLargeFile(ctx, bucket, key, r.Body, contentLength)
	if err != nil {
		log.Printf("UploadLargeFile failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	azureOKHeaders(w)
	w.Header().Set("ETag", result.ETag)
	w.Header().Set("Last-Modified", result.LastModified.Format(http.TimeFormat))
	w.WriteHeader(http.StatusCreated) // Azure Put Blob returns 201
}

func (s *server) handlePutBlock(ctx context.Context, w http.ResponseWriter, r *http.Request, bucket, key, blockIDB64 string) {
	if blockIDB64 == "" {
		http.Error(w, "missing blockid", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	// Azure block IDs are base64-encoded opaque IDs. We store them as-is (base64).
	if _, err := base64.StdEncoding.DecodeString(blockIDB64); err != nil {
		http.Error(w, "blockid must be base64", http.StatusBadRequest)
		return
	}

	// Ensure multipart context and get state
	log.Printf("Ensuring multipart upload for %s/%s", bucket, key)
	us, err := s.ensureMultipart(ctx, bucket, key, r)
	if err != nil || us == nil {
		http.Error(w, "create multipart: "+err.Error(), http.StatusBadGateway)
		return
	}
	log.Printf("Got upload state with ID: %s", us.UploadID)

	// Extract block number from Azure block ID
	blockNum, err := extractBlockNumber(blockIDB64)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid block ID: %v", err), http.StatusBadRequest)
		return
	}

	// S3 part numbers start at 1, so block 0 becomes part 1
	partNum := int32(blockNum + 1)

	// Get content length if available
	contentLength := r.ContentLength
	if contentLength < 0 {
		// If content length is unknown, we need to buffer
		// Fall back to the original implementation for this edge case
		s.handlePutBlockBuffered(ctx, w, r, bucket, key, blockIDB64, us, partNum)
		return
	}

	log.Printf("Streaming block: key=%s blockid=%s blocknum=%d size=%d MB",
		key, blockIDB64, blockNum, contentLength/(1024*1024))

	// Create a pipe to stream data from request to S3
	pr, pw := io.Pipe()

	// Create MD5 hasher for content verification
	hasher := md5.New()

	// Channel to communicate the upload result
	type uploadResult struct {
		etag string
		err  error
	}
	resultChan := make(chan uploadResult, 1)

	// Start upload goroutine
	go func() {
		defer pr.Close()

		// Upload part using the uploader interface
		etag, err := s.uploader.UploadPartStream(ctx, bucket, key, us.UploadID, partNum, pr, contentLength)

		resultChan <- uploadResult{
			etag: etag,
			err:  err,
		}
	}()

	// Copy from request body to pipe writer and calculate MD5
	multiWriter := io.MultiWriter(pw, hasher)
	_, copyErr := io.Copy(multiWriter, r.Body)
	pw.Close() // Signal EOF to the reader

	if copyErr != nil {
		log.Printf("Failed to stream request body: %v", copyErr)
		http.Error(w, fmt.Sprintf("failed to stream request body: %v", copyErr), http.StatusBadRequest)
		return
	}

	// Calculate MD5 after copy completes
	md5Sum := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	// Wait for upload result
	result := <-resultChan

	if result.err != nil {
		log.Printf("Upload part %d failed: %v", partNum, result.err)
		http.Error(w, fmt.Sprintf("upload part %d failed: %v", partNum, result.err), http.StatusBadGateway)
		return
	}

	// Store the part information for this block
	s.mu.Lock()
	us.Parts[blockIDB64] = uploadedPart{
		PartNumber: partNum,
		ETag:       result.etag,
	}
	s.mu.Unlock()

	// Return success to client
	azureOKHeaders(w)
	w.Header().Set("Content-MD5", md5Sum)
	w.Header().Set("x-ms-request-server-encrypted", "false")
	w.WriteHeader(http.StatusCreated)
}

// handlePutBlockBuffered is the fallback for when Content-Length is not known
func (s *server) handlePutBlockBuffered(ctx context.Context, w http.ResponseWriter, r *http.Request, bucket, key, blockIDB64 string, us *uploadState, partNum int32) {
	// Read the entire block data into memory
	hasher := md5.New()
	teeReader := io.TeeReader(r.Body, hasher)

	blockBytes, err := io.ReadAll(teeReader)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		http.Error(w, fmt.Sprintf("failed to read request body: %v", err), http.StatusBadRequest)
		return
	}

	md5Sum := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	blockSize := len(blockBytes)
	log.Printf("Buffered block: key=%s blockid=%s size=%d MB md5=%s",
		key, blockIDB64, blockSize/(1024*1024), md5Sum)

	// Upload part using the uploader interface (with buffer reader)
	etag, err := s.uploader.UploadPartStream(ctx, bucket, key, us.UploadID, partNum, bytes.NewReader(blockBytes), int64(blockSize))
	if err != nil {
		http.Error(w, fmt.Sprintf("upload part %d failed: %v", partNum, err), http.StatusBadGateway)
		return
	}

	// Store the part information for this block
	s.mu.Lock()
	us.Parts[blockIDB64] = uploadedPart{
		PartNumber: partNum,
		ETag:       etag,
	}
	s.mu.Unlock()

	// Return success to client
	azureOKHeaders(w)
	w.Header().Set("Content-MD5", md5Sum)
	w.Header().Set("x-ms-request-server-encrypted", "false")
	w.WriteHeader(http.StatusCreated)
}

func (s *server) handlePutBlockList(ctx context.Context, w http.ResponseWriter, r *http.Request, bucket, key string) {
	// Get existing upload state
	us := s.getUploadState(bucket, key)
	if us == nil {
		http.Error(w, "no multipart context", http.StatusConflict)
		return
	}

	body, _ := io.ReadAll(r.Body)
	var bl blockListXML
	if err := xml.Unmarshal(body, &bl); err != nil {
		http.Error(w, "bad blocklist XML", http.StatusBadRequest)
		return
	}

	// Azure enforces the order in the XML; build S3 CompletedParts in that order.
	order := make([]string, 0, len(bl.Latest)+len(bl.Uncommitted)+len(bl.Committed))
	order = append(order, bl.Latest...)
	order = append(order, bl.Uncommitted...)
	order = append(order, bl.Committed...)

	log.Printf("Committing %d blocks...", len(order))
	startTime := time.Now()

	// Gather all parts that were uploaded for these blocks
	var parts []CompletedPart

	for _, blockID := range order {
		// Get the part info for this block
		part, ok := us.Parts[blockID]
		if !ok {
			http.Error(w, "missing block "+blockID, http.StatusBadRequest)
			return
		}

		parts = append(parts, CompletedPart(part))
	}

	log.Printf("Completing multipart upload with %d parts", len(parts))

	// Complete the multipart upload
	result, err := s.uploader.CompleteMultipartUpload(ctx, bucket, key, us.UploadID, parts)
	if err != nil {
		http.Error(w, "complete multipart: "+err.Error(), http.StatusBadGateway)
		return
	}

	elapsed := time.Since(startTime)
	log.Printf("Multipart upload completed in %v", elapsed)

	// Cleanup state object
	s.deleteUploadState(bucket, key)

	azureOKHeaders(w)
	w.Header().Set("ETag", result.ETag)
	w.Header().Set("Last-Modified", result.LastModified.Format(http.TimeFormat))
	w.WriteHeader(http.StatusCreated) // Azure Put Block List returns 201
}

func (s *server) handleGet(ctx context.Context, w http.ResponseWriter, r *http.Request, bucket, key string) {
	// Forward Range header if present
	rangeOpt := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	// Support both standard Range and Azure's x-ms-range
	if rangeHeader := r.Header.Get("Range"); rangeHeader != "" {
		rangeOpt.Range = aws.String(rangeHeader)
	} else if msRange := r.Header.Get("x-ms-range"); msRange != "" {
		// Azure sends ranges via x-ms-range: bytes=start-end
		rangeOpt.Range = aws.String(msRange)
	}

	// Get object from S3
	out, err := s.s3.GetObject(ctx, rangeOpt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer out.Body.Close()

	// Set Azure response headers
	azureOKHeaders(w)

	// Copy relevant headers from S3 response
	if out.ETag != nil {
		w.Header().Set("ETag", *out.ETag)
	}
	if out.ContentType != nil {
		w.Header().Set("Content-Type", *out.ContentType)
		w.Header().Set("x-ms-blob-content-type", *out.ContentType)
	}
	if out.LastModified != nil {
		w.Header().Set("Last-Modified", out.LastModified.Format(http.TimeFormat))
	}
	if out.ContentLength != nil {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", *out.ContentLength))
		w.Header().Set("x-ms-blob-content-length", fmt.Sprintf("%d", *out.ContentLength))
	}
	if out.ContentRange != nil {
		w.Header().Set("Content-Range", *out.ContentRange)
		if _, ok := w.Header()["X-Ms-Blob-Content-Length"]; !ok {
			if total, ok := parseTotalFromContentRange(*out.ContentRange); ok {
				w.Header().Set("x-ms-blob-content-length", fmt.Sprintf("%d", total))
			}
		}
	}
	if out.AcceptRanges != nil {
		w.Header().Set("Accept-Ranges", *out.AcceptRanges)
	}

	// Add Azure-specific headers
	w.Header().Set("x-ms-blob-type", "BlockBlob")
	w.Header().Set("x-ms-server-encrypted", "false")

	// Set appropriate status code
	status := http.StatusOK
	if out.ContentRange != nil {
		status = http.StatusPartialContent
	}

	w.WriteHeader(status)

	// Stream using io.Copy with no intermediate buffering; let the runtime handle it
	if _, err := io.Copy(w, out.Body); err != nil {
		log.Printf("stream copy error: %v", err)
		return
	}
}

func (s *server) handleHead(ctx context.Context, w http.ResponseWriter, r *http.Request, bucket, key string) {
	// Use HeadObject for proper HEAD request
	out, err := s.s3.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	// Set Azure response headers
	azureOKHeaders(w)

	// Copy relevant headers from S3 response
	if out.ETag != nil {
		w.Header().Set("ETag", *out.ETag)
	}
	if out.ContentType != nil {
		w.Header().Set("Content-Type", *out.ContentType)
		w.Header().Set("x-ms-blob-content-type", *out.ContentType)
	}
	if out.LastModified != nil {
		w.Header().Set("Last-Modified", out.LastModified.Format(http.TimeFormat))
	}
	if out.ContentLength != nil {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", *out.ContentLength))
		w.Header().Set("x-ms-blob-content-length", fmt.Sprintf("%d", *out.ContentLength))
	}
	if out.AcceptRanges != nil {
		w.Header().Set("Accept-Ranges", *out.AcceptRanges)
	}
	w.Header().Set("x-ms-blob-type", "BlockBlob")
	w.Header().Set("x-ms-server-encrypted", "false")

	w.WriteHeader(http.StatusOK)
}

// handleDebugHealth provides a health check endpoint
func (s *server) handleDebugHealth(w http.ResponseWriter, r *http.Request) {
	// Test S3 connectivity
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	health := map[string]interface{}{
		"status":     "ok",
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
		"debug_mode": debugMode,
	}

	// Try to list buckets as a health check
	_, err := s.s3.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		health["status"] = "unhealthy"
		health["s3_error"] = err.Error()
	} else {
		health["s3_status"] = "connected"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

// handleDebugStats provides runtime statistics
func (s *server) handleDebugStats(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	activeUploads := len(s.uploadStates)
	uploadKeys := make([]string, 0, activeUploads)
	for k := range s.uploadStates {
		uploadKeys = append(uploadKeys, k)
	}
	s.mu.Unlock()

	stats := map[string]interface{}{
		"timestamp":      time.Now().UTC().Format(time.RFC3339),
		"active_uploads": activeUploads,
		"upload_keys":    uploadKeys,
		"debug_mode":     debugMode,
		"goroutines":     runtime.NumGoroutine(),
		"performance":    s.stats.getStats(),
	}

	// Memory statistics
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	stats["memory"] = map[string]interface{}{
		"alloc_mb":       float64(m.Alloc) / (1024 * 1024),
		"total_alloc_mb": float64(m.TotalAlloc) / (1024 * 1024),
		"sys_mb":         float64(m.Sys) / (1024 * 1024),
		"num_gc":         m.NumGC,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// Start starts the ASUR Azure→S3 proxy service
func Start(port int) error {
	// Try to load .env file if it exists (optional)
	loadEnvFile()

	// Load credentials from environment variables
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	endpoint := os.Getenv("R2_ENDPOINT")

	// Fallback to S3 credentials if R2 not set
	if accessKeyID == "" {
		accessKeyID = os.Getenv("S3_ACCESS_KEY_ID")
	}
	if secretAccessKey == "" {
		secretAccessKey = os.Getenv("S3_SECRET_ACCESS_KEY")
	}

	// Validate required credentials
	if accessKeyID == "" || secretAccessKey == "" {
		return fmt.Errorf("missing required credentials. Please set R2_ACCESS_KEY_ID and R2_SECRET_ACCESS_KEY environment variables")
	}

	// Validate endpoint is provided
	if endpoint == "" {
		return fmt.Errorf("missing R2_ENDPOINT environment variable")
	}

	// Check if HTTP/2 should be enabled (default: disabled for stability)
	forceHTTP2 := os.Getenv("AZPROXY_ENABLE_HTTP2") == "true"
	if forceHTTP2 {
		log.Printf("HTTP/2 enabled via AZPROXY_ENABLE_HTTP2=true")
	} else {
		log.Printf("HTTP/2 disabled (default). Set AZPROXY_ENABLE_HTTP2=true to enable")
	}

	// Get max connections per host (default: 100)
	maxConnsPerHost := 100
	if v := os.Getenv("AZPROXY_MAX_CONNS_PER_HOST"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			maxConnsPerHost = n
		}
	}
	log.Printf("Max connections per host: %d", maxConnsPerHost)

	// Get default upload concurrency
	defaultConcurrency := 10
	if envConcurrency := os.Getenv("AZPROXY_UPLOAD_CONCURRENCY"); envConcurrency != "" {
		if v, err := strconv.Atoi(envConcurrency); err == nil && v > 0 {
			defaultConcurrency = v
		}
	}

	// Determine upload method from environment
	uploadMethod := os.Getenv("AZPROXY_UPLOAD_METHOD")
	if uploadMethod == "" {
		uploadMethod = "http" // Default to HTTP method for direct R2 uploads
	}

	// Create custom HTTP client with optimized settings for large transfers
	httpClient := &http.Client{
		Timeout: 30 * time.Minute, // Extended timeout for large R2 uploads
		Transport: &http.Transport{
			MaxIdleConns:        1000,            // Increased from 512
			MaxIdleConnsPerHost: 500,             // Increased from 256
			MaxConnsPerHost:     maxConnsPerHost, // Configurable limit
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  true,
			// Increase buffer sizes for better throughput
			WriteBufferSize: 1024 * 1024, // 1MB write buffer (increased from 256KB)
			ReadBufferSize:  1024 * 1024, // 1MB read buffer (increased from 256KB)
			// Control HTTP/2 based on environment variable
			ForceAttemptHTTP2: forceHTTP2,
			// Add response header timeout - increased for R2 large uploads
			ResponseHeaderTimeout: 5 * time.Minute,
			// Expect continue timeout
			ExpectContinueTimeout: 30 * time.Second,
			// Disable keep-alives to force new connections (can help with load balancing)
			DisableKeepAlives: os.Getenv("AZPROXY_DISABLE_KEEPALIVE") == "true",
			// Set aggressive TCP keep-alive
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				d := &net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					Control: func(network, address string, c syscall.RawConn) error {
						return c.Control(func(fd uintptr) {
							// Enable TCP_NODELAY for low latency
							syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 1)
							// Increase socket buffer sizes for better throughput
							syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_SNDBUF, 4*1024*1024) // 4MB send buffer
							syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUF, 4*1024*1024) // 4MB receive buffer
						})
					},
				}
				return d.DialContext(ctx, network, addr)
			},
		},
	}

	// Configure AWS SDK logging based on debug mode
	var logMode aws.ClientLogMode
	if debugMode {
		logMode = aws.LogRetries | aws.LogRequest | aws.LogResponse | aws.LogRequestEventMessage | aws.LogResponseEventMessage
		log.Printf("AWS SDK verbose logging enabled")
	} else {
		logMode = aws.LogRetries
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("auto"), // R2 uses "auto" region
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKeyID,
			secretAccessKey,
			"",
		)),
		config.WithHTTPClient(httpClient),
		config.WithRetryer(func() aws.Retryer {
			return retry.NewStandard(func(o *retry.StandardOptions) {
				o.MaxAttempts = 3 // Reduced from 8 for faster failure feedback
			})
		}),
		config.WithClientLogMode(logMode),
	)
	if err != nil {
		return fmt.Errorf("aws cfg: %v", err)
	}

	// Create S3 client
	s3Config := func(o *s3.Options) {
		o.UsePathStyle = true
		// Disable ARN region for R2 compatibility
		o.UseARNRegion = false
		// Set R2 endpoint
		o.BaseEndpoint = aws.String(endpoint)
		// Disable S3 Express session auth which R2 doesn't support
		o.DisableS3ExpressSessionAuth = aws.Bool(true)
	}

	// S3 client with no retries - for streaming operations
	s3NoRetryConfig := func(o *s3.Options) {
		o.UsePathStyle = true
		o.Retryer = retry.AddWithMaxAttempts(retry.NewStandard(), 1)
		// Disable ARN region for R2 compatibility
		o.UseARNRegion = false
		// Set R2 endpoint
		o.BaseEndpoint = aws.String(endpoint)
		// Disable S3 Express session auth which R2 doesn't support
		o.DisableS3ExpressSessionAuth = aws.Bool(true)
	}

	// Create performance stats
	stats := newPerformanceStats()

	// Create S3 clients
	s3Client := s3.NewFromConfig(cfg, s3Config)
	s3NoRetryClient := s3.NewFromConfig(cfg, s3NoRetryConfig)

	// Create the appropriate uploader based on configuration
	var uploader Uploader
	switch uploadMethod {
	case "http":
		log.Printf("Using HTTP uploader for direct R2 uploads")
		uploader = NewHTTPUploader(endpoint, httpClient, accessKeyID, secretAccessKey, stats, defaultConcurrency)
	default:
		log.Printf("Using S3 SDK uploader")
		uploader = NewS3Uploader(s3Client, s3NoRetryClient, stats, defaultConcurrency)
	}

	s := &server{
		s3:           s3Client,
		s3NoRetry:    s3NoRetryClient,
		uploadStates: make(map[string]*uploadState),
		stats:        stats,
		uploader:     uploader,
	}

	// Configure server with optimized settings
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      loggingMiddleware(s),
		ReadTimeout:  30 * time.Minute,
		WriteTimeout: 30 * time.Minute,
		// Increase max header size for large requests
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	log.Printf("========================================")
	log.Printf("Azure→S3 Proxy Starting")
	log.Printf("========================================")
	log.Printf("Listening on: %s", srv.Addr)
	log.Printf("Debug mode: %v", debugMode)
	log.Printf("HTTP/2: %v", forceHTTP2)
	log.Printf("Max connections per host: %d", maxConnsPerHost)
	log.Printf("Upload concurrency: %d (set AZPROXY_UPLOAD_CONCURRENCY to change)", defaultConcurrency)
	log.Printf("Upload method: %s (set AZPROXY_UPLOAD_METHOD to 's3' for AWS SDK uploads)", uploadMethod)
	log.Printf("Using endpoint: %s", endpoint)
	log.Printf("Credentials loaded from environment variables")

	if debugMode {
		log.Printf("Health check: http://localhost:%d/_debug/health", port)
		log.Printf("Stats endpoint: http://localhost:%d/_debug/stats", port)
	}

	// Start a goroutine to clean up old incomplete uploads
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			s.cleanupIncompleteUploads(context.Background())
		}
	}()

	return srv.ListenAndServe()
}

// cleanupIncompleteUploads removes incomplete multipart uploads
func (s *server) cleanupIncompleteUploads(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cleaned := 0
	for key, state := range s.uploadStates {
		parts := strings.SplitN(key, "/", 2)
		if len(parts) != 2 {
			continue
		}
		bucket, objectKey := parts[0], parts[1]

		// Try to abort the multipart upload
		err := s.uploader.AbortMultipartUpload(ctx, bucket, objectKey, state.UploadID)
		if err == nil {
			delete(s.uploadStates, key)
			cleaned++
		}
	}

	if cleaned > 0 {
		log.Printf("Cleaned up %d incomplete uploads", cleaned)
	}
}

func parseAzuriteStyle(p string) (account, container, key string, _ error) {
	parts := strings.Split(strings.TrimPrefix(p, "/"), "/")
	if len(parts) < 3 {
		return "", "", "", fmt.Errorf("bad path")
	}
	account, container = parts[0], parts[1]
	key = strings.Join(parts[2:], "/")
	return
}

func bucketFor(_account, _container string) string {
	return "prajtestcacheproxyeuauto"
}

// loadEnvFile tries to load .env file if it exists
func loadEnvFile() {
	// Simple .env file loader - reads key=value pairs
	data, err := os.ReadFile(".env")
	if err != nil {
		// .env file not found is okay, just use system env vars
		return
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Only set if not already set in environment
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}
	log.Printf("Loaded configuration from .env file")
}
