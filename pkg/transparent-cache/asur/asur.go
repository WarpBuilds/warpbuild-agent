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
	"github.com/warpbuilds/warpbuild-agent/pkg/transparent-cache/derp"
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

// getS3ClientForURL retrieves or creates S3 clients based on credentials from the URL
func (s *server) getS3ClientForURL(urlStr string) (*s3ClientPair, error) {
	// Get credentials from derp based on the URL
	provider, derpCreds, found := derp.GetCredentialsFromURL(urlStr)
	if !found {
		return nil, fmt.Errorf("no credentials found for URL: %s", urlStr)
	}

	// Create a cache key based on provider and credentials
	cacheKey := fmt.Sprintf("%s:%v", provider, derpCreds)

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if we already have clients for these credentials
	if clientPair, exists := s.s3ClientCache[cacheKey]; exists {
		clientPair.lastUsed = time.Now()
		return clientPair, nil
	}

	// Create new S3 clients based on the credentials
	var s3Client, s3NoRetryClient *s3.Client
	var uploader Uploader

	switch provider {
	case derp.ProviderS3:
		// Handle S3 credentials - use access credentials only
		var accessKeyID, secretAccessKey, sessionToken string
		var region string = "auto" // Default region

		// Extract credentials based on response type
		switch creds := derpCreds.(type) {
		case *derp.S3GetCacheResponse:
			// Use access grant if available
			if creds.AccessGrant != nil {
				accessKeyID = creds.AccessGrant.AccessKeyID
				secretAccessKey = creds.AccessGrant.SecretAccessKey
				sessionToken = creds.AccessGrant.SessionToken
			}
		case *derp.S3ReserveCacheResponse:
			// Use access grant if available
			if creds.AccessGrant != nil {
				accessKeyID = creds.AccessGrant.AccessKeyID
				secretAccessKey = creds.AccessGrant.SecretAccessKey
				sessionToken = creds.AccessGrant.SessionToken
			}
		}

		// We require explicit credentials for S3
		if accessKeyID == "" || secretAccessKey == "" {
			return nil, fmt.Errorf("no AWS credentials available for S3 provider")
		}

		// Extract bucket name for endpoint determination
		var bucketName string
		switch creds := derpCreds.(type) {
		case *derp.S3GetCacheResponse:
			if creds.AccessGrant != nil {
				bucketName = creds.AccessGrant.BucketName
			}
		case *derp.S3ReserveCacheResponse:
			if creds.AccessGrant != nil {
				bucketName = creds.AccessGrant.BucketName
			}
		}

		// Determine endpoint and region based on bucket
		var endpoint string
		if bucketName != "" && strings.Contains(bucketName, ".") {
			// For S3, check if bucket name contains endpoint information
			// Some S3-compatible services encode the endpoint in the bucket name
			// e.g., "bucket.s3.us-west-2.amazonaws.com"
			parts := strings.SplitN(bucketName, ".", 2)
			if len(parts) > 1 && strings.Contains(parts[1], "amazonaws.com") {
				// Extract region from bucket name if it follows AWS pattern
				if strings.HasPrefix(parts[1], "s3.") {
					regionParts := strings.Split(parts[1], ".")
					if len(regionParts) >= 3 {
						region = regionParts[1]
					}
				}
			}
		}

		if debugMode {
			log.Printf("S3 configuration: provider=%s, bucket=%s, region=%s, endpoint=%s",
				provider, bucketName, region, endpoint)
		}

		// Create AWS config
		cfg, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				accessKeyID,
				secretAccessKey,
				sessionToken,
			)),
			config.WithHTTPClient(s.httpClient),
			config.WithRetryer(func() aws.Retryer {
				return retry.NewStandard(func(o *retry.StandardOptions) {
					o.MaxAttempts = 3
				})
			}),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create AWS config: %w", err)
		}

		// Create S3 client options
		s3Config := func(o *s3.Options) {
			o.UsePathStyle = true
			o.UseARNRegion = false
			if endpoint != "" {
				o.BaseEndpoint = aws.String(endpoint)
			}
			o.DisableS3ExpressSessionAuth = aws.Bool(true)
		}

		// S3 client with no retries - for streaming operations
		s3NoRetryConfig := func(o *s3.Options) {
			o.UsePathStyle = true
			o.Retryer = retry.AddWithMaxAttempts(retry.NewStandard(), 1)
			o.UseARNRegion = false
			if endpoint != "" {
				o.BaseEndpoint = aws.String(endpoint)
			}
			o.DisableS3ExpressSessionAuth = aws.Bool(true)
		}

		// Create S3 clients
		s3Client = s3.NewFromConfig(cfg, s3Config)
		s3NoRetryClient = s3.NewFromConfig(cfg, s3NoRetryConfig)

		// Create uploader
		uploader = NewS3Uploader(s3Client, s3NoRetryClient, s.defaultConcurrency)

	case derp.ProviderR2:
		// Handle R2 credentials - use HTTPUploader for R2
		var accessKeyID, secretAccessKey, sessionToken string
		var accountID string

		// Extract credentials and account ID based on response type
		switch creds := derpCreds.(type) {
		case *derp.S3GetCacheResponse:
			// Use access grant if available
			if creds.AccessGrant != nil {
				accessKeyID = creds.AccessGrant.AccessKeyID
				secretAccessKey = creds.AccessGrant.SecretAccessKey
				sessionToken = creds.AccessGrant.SessionToken
				accountID = creds.AccessGrant.AccountID
			}
		case *derp.S3ReserveCacheResponse:
			// Use access grant if available
			if creds.AccessGrant != nil {
				accessKeyID = creds.AccessGrant.AccessKeyID
				secretAccessKey = creds.AccessGrant.SecretAccessKey
				sessionToken = creds.AccessGrant.SessionToken
				accountID = creds.AccessGrant.AccountID
			}
		}

		// We require explicit credentials for R2
		if accessKeyID == "" || secretAccessKey == "" {
			return nil, fmt.Errorf("no AWS credentials available for R2 provider")
		}

		// We require account ID for R2
		if accountID == "" {
			return nil, fmt.Errorf("no account ID available for R2 provider")
		}

		// R2 endpoint format: https://<account-id>.r2.cloudflarestorage.com
		endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)

		if debugMode {
			log.Printf("R2 configuration: accountID=%s, endpoint=%s", accountID, endpoint)
		}

		// Create AWS config for R2
		cfg, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion("auto"), // R2 uses "auto" region
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				accessKeyID,
				secretAccessKey,
				sessionToken,
			)),
			config.WithHTTPClient(s.httpClient),
			config.WithRetryer(func() aws.Retryer {
				return retry.NewStandard(func(o *retry.StandardOptions) {
					o.MaxAttempts = 3
				})
			}),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create AWS config for R2: %w", err)
		}

		// Create S3 client options for R2
		s3Config := func(o *s3.Options) {
			o.UsePathStyle = true
			o.BaseEndpoint = aws.String(endpoint)
			o.UseARNRegion = false
			o.DisableS3ExpressSessionAuth = aws.Bool(true)
		}

		// S3 client with no retries - for streaming operations
		s3NoRetryConfig := func(o *s3.Options) {
			o.UsePathStyle = true
			o.Retryer = retry.AddWithMaxAttempts(retry.NewStandard(), 1)
			o.BaseEndpoint = aws.String(endpoint)
			o.UseARNRegion = false
			o.DisableS3ExpressSessionAuth = aws.Bool(true)
		}

		// Create S3 clients for R2 (needed for GET/HEAD operations)
		s3Client = s3.NewFromConfig(cfg, s3Config)
		s3NoRetryClient = s3.NewFromConfig(cfg, s3NoRetryConfig)

		// Create HTTPUploader for R2 (used for upload operations)
		uploader = NewHTTPUploader(endpoint, s.httpClient, accessKeyID, secretAccessKey, sessionToken, s.defaultConcurrency)

	case derp.ProviderAzureBlob:
		// For Azure, we'd typically use Azure SDK, but for now we'll return an error
		// In a full implementation, you'd create an Azure Blob client here
		return nil, fmt.Errorf("Azure Blob provider not yet implemented for dynamic clients")

	case derp.ProviderGCS:
		// For GCS, we'd use the GCS SDK with the short-lived token
		// In a full implementation, you'd create a GCS client here
		return nil, fmt.Errorf("GCS provider not yet implemented for dynamic clients")

	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	// Create and cache the client pair
	clientPair := &s3ClientPair{
		client:        s3Client,
		noRetryClient: s3NoRetryClient,
		uploader:      uploader,
		lastUsed:      time.Now(),
	}

	s.s3ClientCache[cacheKey] = clientPair
	return clientPair, nil
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

	// Validate block number is within reasonable range
	// Azure allows up to 50,000 blocks per blob
	if blockNum < 0 || blockNum >= 50000 {
		return 0, fmt.Errorf("block number %d out of valid range (0-49999)", blockNum)
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
	// Default S3 client for fallback (if any)
	s3 *s3.Client
	// s3NoRetry disables automatic retries to avoid rewind requirement on streaming bodies
	s3NoRetry *s3.Client
	// mutex for thread-safe operations
	mu sync.Mutex
	// in-memory upload states: "bucket/key" -> uploadState
	uploadStates map[string]*uploadState
	// Dynamic S3 client cache: cacheKey -> *s3ClientPair
	s3ClientCache map[string]*s3ClientPair
	// Default HTTP client for creating S3 clients
	httpClient *http.Client
	// Default upload concurrency
	defaultConcurrency int
}

// s3ClientPair holds both regular and no-retry S3 clients
type s3ClientPair struct {
	client        *s3.Client
	noRetryClient *s3.Client
	uploader      Uploader
	lastUsed      time.Time
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Build the full URL from the request
	fullURL := fmt.Sprintf("https://%s%s", r.Host, r.URL.Path)

	if debugMode {
		log.Printf("ServeHTTP: Processing URL: %s", fullURL)
	}

	// Get credentials from derp
	provider, credentials, found := derp.GetCredentialsFromURL(fullURL)
	if !found {
		http.Error(w, "no credentials found for URL", http.StatusBadRequest)
		return
	}

	if debugMode {
		log.Printf("ServeHTTP: Found provider=%s for URL=%s", provider, fullURL)
	}

	// Extract bucket name from credentials
	bucket, err := getBucketFromCredentials(provider, credentials)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get bucket: %v", err), http.StatusBadRequest)
		return
	}

	// Extract the blob key from the URL path
	// Expected format: /{provider}/{cacheKey}--{version}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		http.Error(w, "invalid path format", http.StatusBadRequest)
		return
	}
	blobKey := strings.Join(parts[1:], "/") // Everything after the provider prefix

	// For S3/R2, check if we have a specific key from the backend
	actualKey := blobKey
	if provider == derp.ProviderS3 || provider == derp.ProviderR2 {
		// Extract cache identifier from blobKey (format: cacheKey--version)
		if r.Method == http.MethodPut {
			// For uploads, check for UploadKey
			if uploadKey, found := derp.GetS3UploadKey(blobKey); found {
				actualKey = uploadKey
				if debugMode {
					log.Printf("ServeHTTP: Using backend UploadKey=%s instead of blobKey=%s", uploadKey, blobKey)
				}
			}
		} else if r.Method == http.MethodGet || r.Method == http.MethodHead {
			// For downloads, check for CacheKey
			if cacheKey, found := derp.GetS3CacheKey(blobKey); found {
				actualKey = cacheKey
				if debugMode {
					log.Printf("ServeHTTP: Using backend CacheKey=%s instead of blobKey=%s", cacheKey, blobKey)
				}
			}
		}
	}

	if debugMode {
		log.Printf("ServeHTTP: bucket=%s, blobKey=%s, actualKey=%s", bucket, blobKey, actualKey)
	}

	switch r.Method {
	case http.MethodPut:
		q := r.URL.Query()
		switch strings.ToLower(q.Get("comp")) {
		case "block":
			s.handlePutBlock(ctx, w, r, bucket, actualKey, blobKey, q.Get("blockid"))
		case "blocklist":
			s.handlePutBlockList(ctx, w, r, bucket, actualKey, blobKey)
		default:
			s.handlePutBlob(ctx, w, r, bucket, actualKey)
		}
	case http.MethodGet:
		s.handleGet(ctx, w, r, bucket, actualKey)
	case http.MethodHead:
		s.handleHead(ctx, w, r, bucket, actualKey)
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
		// Build query string with ? prefix if present
		queryString := ""
		if r.URL.RawQuery != "" {
			queryString = "?" + r.URL.RawQuery
		}

		fullURL := fmt.Sprintf("https://%s%s%s", r.Host, r.URL.Path, queryString)

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
	// Build the full Azure-style URL from the request
	fullURL := fmt.Sprintf("https://%s%s", r.Host, r.URL.Path)

	// Get S3 client for this URL
	clientPair, err := s.getS3ClientForURL(fullURL)
	if err != nil {
		log.Printf("Failed to get S3 client for URL %s: %v", fullURL, err)
		http.Error(w, fmt.Sprintf("Failed to get storage client: %v", err), http.StatusBadGateway)
		return
	}

	contentLength := r.ContentLength
	log.Printf("handlePutBlob: key=%s size=%d MB", key, contentLength/(1024*1024))

	// Extract cache identifier from the URL path (format: /provider/cacheKey--version)
	var cacheIdentifier string
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(parts) >= 2 {
		cacheIdentifier = parts[1] // This should be cacheKey--version
	}

	// For small files (< 5MB), use direct upload
	if contentLength > 0 && contentLength < 5*1024*1024 {
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := clientPair.uploader.UploadSmallFile(ctx, bucket, key, buf)
		if err != nil {
			log.Printf("UploadSmallFile failed: %v", err)
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		// Track small file upload completion if we have a cache identifier
		if cacheIdentifier != "" {
			// Store a special marker indicating direct upload completion
			derp.AddS3UploadedBlock(cacheIdentifier, "DIRECT_UPLOAD_SMALL", 1, result.ETag)
			derp.SetS3BlockOrder(cacheIdentifier, []string{"DIRECT_UPLOAD_SMALL"})
		}

		azureOKHeaders(w)
		w.Header().Set("ETag", result.ETag)
		w.Header().Set("Last-Modified", result.LastModified.Format(http.TimeFormat))
		w.WriteHeader(http.StatusCreated)
		return
	}

	// For larger files, use concurrent multipart upload with tracking
	log.Printf("Using multipart upload for large file: %s", key)

	// Get the backend's upload ID (required for multipart uploads)
	if cacheIdentifier == "" {
		log.Printf("ERROR: No cache identifier found for multipart upload")
		http.Error(w, "Cache identifier required for multipart upload", http.StatusBadRequest)
		return
	}

	uploadID, ok := derp.GetS3UploadID(cacheIdentifier)
	if !ok || uploadID == "" {
		log.Printf("ERROR: No backend upload ID found for cacheIdentifier=%s", cacheIdentifier)
		http.Error(w, "Backend upload ID not found", http.StatusBadRequest)
		return
	}

	log.Printf("Using backend-provided upload ID: %s for cacheIdentifier=%s", uploadID, cacheIdentifier)

	// Use the backend's upload ID
	result, err := clientPair.uploader.UploadLargeFile(ctx, bucket, key, r.Body, contentLength, cacheIdentifier, uploadID)

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

func (s *server) handlePutBlock(ctx context.Context, w http.ResponseWriter, r *http.Request, bucket, key, cacheIdentifier, blockIDB64 string) {
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

	// Build the full Azure-style URL from the request
	fullURL := fmt.Sprintf("https://%s%s", r.Host, r.URL.Path)

	// Get S3 client for this URL
	clientPair, err := s.getS3ClientForURL(fullURL)
	if err != nil {
		log.Printf("Failed to get S3 client for URL %s: %v", fullURL, err)
		http.Error(w, fmt.Sprintf("Failed to get storage client: %v", err), http.StatusBadGateway)
		return
	}

	// Use UploadID provided by backend (stored by derp)
	uploadID, ok := derp.GetS3UploadID(cacheIdentifier)
	if !ok || uploadID == "" {
		http.Error(w, "missing upload id for multipart", http.StatusBadRequest)
		return
	}
	log.Printf("Using provided multipart upload ID: %s", uploadID)

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
		s.handlePutBlockBuffered(ctx, w, r, bucket, key, cacheIdentifier, blockIDB64, uploadID, partNum, clientPair.uploader)
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
		etag, err := clientPair.uploader.UploadPartStream(ctx, bucket, key, uploadID, partNum, pr, contentLength)

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

	// Record the part in derp for finalize
	derp.AddS3UploadedBlock(cacheIdentifier, blockIDB64, partNum, result.etag)

	// Return success to client
	azureOKHeaders(w)
	w.Header().Set("Content-MD5", md5Sum)
	w.Header().Set("x-ms-request-server-encrypted", "false")
	w.WriteHeader(http.StatusCreated)
}

// handlePutBlockBuffered is the fallback for when Content-Length is not known
func (s *server) handlePutBlockBuffered(ctx context.Context, w http.ResponseWriter, r *http.Request, bucket, key, cacheIdentifier, blockIDB64, uploadID string, partNum int32, uploader Uploader) {
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
	etag, err := uploader.UploadPartStream(ctx, bucket, key, uploadID, partNum, bytes.NewReader(blockBytes), int64(blockSize))
	if err != nil {
		http.Error(w, fmt.Sprintf("upload part %d failed: %v", partNum, err), http.StatusBadGateway)
		return
	}

	// Record the part in derp for finalize
	derp.AddS3UploadedBlock(cacheIdentifier, blockIDB64, partNum, etag)

	// Return success to client
	azureOKHeaders(w)
	w.Header().Set("Content-MD5", md5Sum)
	w.Header().Set("x-ms-request-server-encrypted", "false")
	w.WriteHeader(http.StatusCreated)
}

func (s *server) handlePutBlockList(ctx context.Context, w http.ResponseWriter, r *http.Request, bucket, key, cacheIdentifier string) {
	// Build the full Azure-style URL from the request
	fullURL := fmt.Sprintf("https://%s%s", r.Host, r.URL.Path)

	// Get S3 client for this URL
	clientPair, err := s.getS3ClientForURL(fullURL)
	if err != nil {
		log.Printf("Failed to get S3 client for URL %s: %v", fullURL, err)
		http.Error(w, fmt.Sprintf("Failed to get storage client: %v", err), http.StatusBadGateway)
		return
	}

	_ = clientPair // clientPair retained for parity; not used here

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	var bl blockListXML
	if err := xml.Unmarshal(body, &bl); err != nil {
		http.Error(w, "bad blocklist XML", http.StatusBadRequest)
		return
	}

	// Azure enforces the order in the XML; build order in that sequence.
	order := make([]string, 0, len(bl.Latest)+len(bl.Uncommitted)+len(bl.Committed))
	order = append(order, bl.Latest...)
	order = append(order, bl.Uncommitted...)
	order = append(order, bl.Committed...)

	log.Printf("Received block list with %d blocks (deferring completion to finalize)", len(order))

	// Store order for finalize stage
	derp.SetS3BlockOrder(cacheIdentifier, order)

	// Do NOT complete multipart here. Only acknowledge.
	azureOKHeaders(w)
	w.Header().Set("ETag", "")
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	w.WriteHeader(http.StatusCreated) // Azure Put Block List returns 201
}

func (s *server) handleGet(ctx context.Context, w http.ResponseWriter, r *http.Request, bucket, key string) {
	// Build the full Azure-style URL from the request
	fullURL := fmt.Sprintf("https://%s%s", r.Host, r.URL.Path)

	// Get S3 client for this URL
	clientPair, err := s.getS3ClientForURL(fullURL)
	if err != nil {
		log.Printf("Failed to get S3 client for URL %s: %v", fullURL, err)
		http.Error(w, fmt.Sprintf("Failed to get storage client: %v", err), http.StatusBadGateway)
		return
	}

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
	out, err := clientPair.client.GetObject(ctx, rangeOpt)
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
	// Build the full Azure-style URL from the request
	fullURL := fmt.Sprintf("https://%s%s", r.Host, r.URL.Path)

	// Get S3 client for this URL
	clientPair, err := s.getS3ClientForURL(fullURL)
	if err != nil {
		log.Printf("Failed to get S3 client for URL %s: %v", fullURL, err)
		http.Error(w, fmt.Sprintf("Failed to get storage client: %v", err), http.StatusBadGateway)
		return
	}

	// Use HeadObject for proper HEAD request
	out, err := clientPair.client.HeadObject(ctx, &s3.HeadObjectInput{
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

// Start starts the ASUR Azure→S3 proxy service
func Start(port int) error {
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

	// Create custom HTTP client with optimized settings for large transfers
	httpClient := &http.Client{
		Timeout: 30 * time.Minute, // Extended timeout for large uploads
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
			// Add response header timeout - increased for large uploads
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

	s := &server{
		s3:           nil, // No default S3 client - will be created dynamically
		s3NoRetry:    nil, // No default S3 client - will be created dynamically
		uploadStates: make(map[string]*uploadState),

		s3ClientCache:      make(map[string]*s3ClientPair),
		httpClient:         httpClient,
		defaultConcurrency: defaultConcurrency,
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
	log.Printf("Dynamic credential loading enabled - no environment variables required")
	log.Printf("Credentials will be retrieved from derp based on request URLs")

	if debugMode {
		log.Printf("Health check: http://localhost:%d/_debug/health", port)
		log.Printf("Stats endpoint: http://localhost:%d/_debug/stats", port)
	}

	// Start a goroutine to clean up old cached clients
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			s.cleanupOldClients(context.Background())
		}
	}()

	return srv.ListenAndServe()
}

// getBucketFromCredentials extracts the bucket name from derp credentials
func getBucketFromCredentials(provider derp.Provider, credentials interface{}) (string, error) {
	if debugMode {
		log.Printf("getBucketFromCredentials: provider=%s, credType=%T", provider, credentials)
	}

	switch provider {
	case derp.ProviderS3, derp.ProviderR2:
		// For S3/R2, check both GetCacheResponse and ReserveCacheResponse
		switch creds := credentials.(type) {
		case *derp.S3GetCacheResponse:
			if creds.AccessGrant != nil && creds.AccessGrant.BucketName != "" {
				if debugMode {
					log.Printf("getBucketFromCredentials: Found S3 bucket from GetCacheResponse: %s", creds.AccessGrant.BucketName)
				}
				return creds.AccessGrant.BucketName, nil
			}
		case *derp.S3ReserveCacheResponse:
			if creds.AccessGrant != nil && creds.AccessGrant.BucketName != "" {
				if debugMode {
					log.Printf("getBucketFromCredentials: Found S3 bucket from ReserveCacheResponse: %s", creds.AccessGrant.BucketName)
				}
				return creds.AccessGrant.BucketName, nil
			}
		}
	case derp.ProviderAzureBlob:
		switch creds := credentials.(type) {
		case *derp.AzureBlobGetCacheResponse:
			if creds.BucketName != "" {
				if debugMode {
					log.Printf("getBucketFromCredentials: Found Azure bucket from GetCacheResponse: %s", creds.BucketName)
				}
				return creds.BucketName, nil
			}
		case *derp.AzureBlobReserveCacheResponse:
			if creds.ContainerName != "" {
				if debugMode {
					log.Printf("getBucketFromCredentials: Found Azure container from ReserveCacheResponse: %s", creds.ContainerName)
				}
				return creds.ContainerName, nil
			}
		}
	case derp.ProviderGCS:
		switch creds := credentials.(type) {
		case *derp.GCSGetCacheResponse:
			if creds.BucketName != "" {
				if debugMode {
					log.Printf("getBucketFromCredentials: Found GCS bucket from GetCacheResponse: %s", creds.BucketName)
				}
				return creds.BucketName, nil
			}
		case *derp.GCSReserveCacheResponse:
			if creds.BucketName != "" {
				if debugMode {
					log.Printf("getBucketFromCredentials: Found GCS bucket from ReserveCacheResponse: %s", creds.BucketName)
				}
				return creds.BucketName, nil
			}
		}
	}

	// Fallback to default bucket
	defaultBucket := "prajtestcacheproxyeuauto"
	if debugMode {
		log.Printf("getBucketFromCredentials: Using fallback bucket: %s", defaultBucket)
	}
	return defaultBucket, nil
}

// cleanupOldClients removes cached S3 clients that haven't been used recently
func (s *server) cleanupOldClients(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	expireTime := 1 * time.Hour // Remove clients not used for 1 hour
	cleaned := 0

	for key, clientPair := range s.s3ClientCache {
		if now.Sub(clientPair.lastUsed) > expireTime {
			delete(s.s3ClientCache, key)
			cleaned++
		}
	}

	if cleaned > 0 {
		log.Printf("Cleaned up %d old S3 clients", cleaned)
	}
}
