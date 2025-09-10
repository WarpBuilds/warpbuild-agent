package asur

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// UploadResult represents the result of an upload operation
type UploadResult struct {
	ETag         string
	LastModified time.Time
}

// Uploader interface defines methods for different upload strategies
type Uploader interface {
	// UploadSmallFile uploads a file directly (for files < 5MB)
	UploadSmallFile(ctx context.Context, bucket, key string, data []byte) (*UploadResult, error)

	// UploadLargeFile uploads a large file using multipart/concurrent upload
	UploadLargeFile(ctx context.Context, bucket, key string, reader io.Reader, contentLength int64) (*UploadResult, error)

	// Block-based upload methods for Azure compatibility

	// EnsureMultipartUpload ensures a multipart upload exists, creating if necessary
	EnsureMultipartUpload(ctx context.Context, bucket, key string) (uploadID string, err error)

	// UploadPart uploads a single part with streaming support
	UploadPartStream(ctx context.Context, bucket, key, uploadID string, partNumber int32, reader io.Reader, contentLength int64) (etag string, err error)

	// CompleteMultipartUpload completes a multipart upload
	CompleteMultipartUpload(ctx context.Context, bucket, key, uploadID string, parts []CompletedPart) (*UploadResult, error)

	// AbortMultipartUpload cancels a multipart upload
	AbortMultipartUpload(ctx context.Context, bucket, key, uploadID string) error
}

// CompletedPart represents a completed part in multipart upload
type CompletedPart struct {
	PartNumber int32
	ETag       string
}

// S3Uploader implements Uploader using AWS S3 SDK
type S3Uploader struct {
	s3Client       *s3.Client
	s3NoRetry      *s3.Client
	stats          *performanceStats
	maxConcurrency int
}

// NewS3Uploader creates a new S3Uploader
func NewS3Uploader(s3Client, s3NoRetry *s3.Client, stats *performanceStats, maxConcurrency int) *S3Uploader {
	return &S3Uploader{
		s3Client:       s3Client,
		s3NoRetry:      s3NoRetry,
		stats:          stats,
		maxConcurrency: maxConcurrency,
	}
}

// UploadSmallFile implements direct upload for small files
func (u *S3Uploader) UploadSmallFile(ctx context.Context, bucket, key string, data []byte) (*UploadResult, error) {
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	}

	out, err := u.s3Client.PutObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("PutObject failed: %w", err)
	}

	return &UploadResult{
		ETag:         aws.ToString(out.ETag),
		LastModified: time.Now().UTC(),
	}, nil
}

// UploadLargeFile implements concurrent multipart upload for large files
func (u *S3Uploader) UploadLargeFile(ctx context.Context, bucket, key string, reader io.Reader, contentLength int64) (*UploadResult, error) {
	// Track performance
	uploadStart := time.Now()
	u.stats.incrementActive()
	defer u.stats.decrementActive()

	// Create multipart upload
	createInput := &s3.CreateMultipartUploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	createOut, err := u.s3Client.CreateMultipartUpload(ctx, createInput)
	if err != nil {
		return nil, fmt.Errorf("create multipart upload failed: %w", err)
	}
	uploadID := aws.ToString(createOut.UploadId)
	log.Printf("Created multipart upload: %s", uploadID)

	// Determine optimal chunk size and concurrency
	const minChunkSize = 5 * 1024 * 1024   // 5MB minimum
	const maxChunkSize = 100 * 1024 * 1024 // 100MB maximum
	const targetChunks = 100               // Target number of chunks for optimal concurrency

	chunkSize := contentLength / targetChunks
	if chunkSize < minChunkSize {
		chunkSize = minChunkSize
	} else if chunkSize > maxChunkSize {
		chunkSize = maxChunkSize
	}

	// Use channels to coordinate concurrent uploads
	type uploadJob struct {
		partNum int32
		data    []byte
	}

	type uploadResult struct {
		partNum int32
		part    s3types.CompletedPart
		err     error
	}

	jobChan := make(chan uploadJob, u.maxConcurrency)
	resultChan := make(chan uploadResult, u.maxConcurrency)

	// Start worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < u.maxConcurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobChan {
				// Upload with retries
				var out *s3.UploadPartOutput
				var uploadErr error

				for attempt := 1; attempt <= 3; attempt++ {
					timeoutDuration := max(time.Duration(len(job.data)/(10*1024*1024)+6)*time.Minute, 8*time.Minute)
					uploadCtx, cancel := context.WithTimeout(ctx, timeoutDuration)

					input := &s3.UploadPartInput{
						Bucket:        aws.String(bucket),
						Key:           aws.String(key),
						UploadId:      aws.String(uploadID),
						PartNumber:    aws.Int32(job.partNum),
						Body:          bytes.NewReader(job.data),
						ContentLength: aws.Int64(int64(len(job.data))),
					}

					log.Printf("[Worker %d] Uploading part %d (size=%d MB) - attempt %d/3",
						workerID, job.partNum, len(job.data)/(1024*1024), attempt)
					start := time.Now()

					out, uploadErr = u.s3Client.UploadPart(uploadCtx, input)
					cancel()

					if uploadErr == nil {
						duration := time.Since(start)
						speed := float64(len(job.data)) / (1024 * 1024) / duration.Seconds()
						log.Printf("[Worker %d] Uploaded part %d in %v (%.1f MB/s)",
							workerID, job.partNum, duration, speed)
						break
					}

					logError("UploadPart", uploadErr, map[string]interface{}{
						"worker_id":     workerID,
						"part_number":   job.partNum,
						"attempt":       attempt,
						"chunk_size_mb": len(job.data) / (1024 * 1024),
					})

					if attempt < 3 {
						time.Sleep(time.Duration(attempt) * time.Second)
					}
				}

				if uploadErr != nil {
					resultChan <- uploadResult{
						partNum: job.partNum,
						err:     uploadErr,
					}
				} else {
					resultChan <- uploadResult{
						partNum: job.partNum,
						part: s3types.CompletedPart{
							ETag:       out.ETag,
							PartNumber: aws.Int32(job.partNum),
						},
					}
				}
			}
		}(i)
	}

	// Result collector goroutine
	parts := make([]s3types.CompletedPart, 0)
	partsMu := sync.Mutex{}
	uploadErrors := make([]error, 0)

	go func() {
		for result := range resultChan {
			if result.err != nil {
				partsMu.Lock()
				uploadErrors = append(uploadErrors, result.err)
				partsMu.Unlock()
			} else {
				partsMu.Lock()
				parts = append(parts, result.part)
				partsMu.Unlock()
			}
		}
	}()

	// Read and distribute chunks to workers
	partNum := int32(1)
	totalRead := int64(0)

	for totalRead < contentLength {
		// Calculate chunk size for this iteration
		remainingBytes := contentLength - totalRead
		currentChunkSize := int64(chunkSize)
		if remainingBytes < currentChunkSize {
			currentChunkSize = remainingBytes
		}

		// Read chunk
		chunk := make([]byte, currentChunkSize)
		n, err := io.ReadFull(reader, chunk)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Printf("Read error: %v", err)
			close(jobChan)
			wg.Wait()
			close(resultChan)

			// Abort the upload
			u.s3Client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
				Bucket:   aws.String(bucket),
				Key:      aws.String(key),
				UploadId: aws.String(uploadID),
			})
			return nil, fmt.Errorf("read error: %w", err)
		}

		// Trim to actual size read
		chunk = chunk[:n]
		totalRead += int64(n)

		// Send to worker
		jobChan <- uploadJob{
			partNum: partNum,
			data:    chunk,
		}

		partNum++

		// Break if we've read everything
		if n == 0 || err == io.EOF {
			break
		}
	}

	// Close job channel and wait for workers to finish
	close(jobChan)
	wg.Wait()
	close(resultChan)

	// Check for errors
	if len(uploadErrors) > 0 {
		// Abort the upload
		u.s3Client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
			Bucket:   aws.String(bucket),
			Key:      aws.String(key),
			UploadId: aws.String(uploadID),
		})
		return nil, fmt.Errorf("upload failed: %w", uploadErrors[0])
	}

	// Sort parts by part number (required by S3)
	sort.Slice(parts, func(i, j int) bool {
		return *parts[i].PartNumber < *parts[j].PartNumber
	})

	// Complete multipart upload
	log.Printf("Completing multipart upload with %d parts", len(parts))
	completeOut, err := u.s3Client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadID),
		MultipartUpload: &s3types.CompletedMultipartUpload{
			Parts: parts,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("complete multipart upload failed: %w", err)
	}

	// Record performance statistics
	uploadDuration := time.Since(uploadStart)
	u.stats.recordUpload(contentLength, uploadDuration, len(parts), u.maxConcurrency)
	throughput := float64(contentLength) / (1024 * 1024) / uploadDuration.Seconds()
	log.Printf("Upload completed: %s, size=%d MB, duration=%v, throughput=%.1f MB/s",
		key, contentLength/(1024*1024), uploadDuration, throughput)

	return &UploadResult{
		ETag:         aws.ToString(completeOut.ETag),
		LastModified: time.Now().UTC(),
	}, nil
}

// EnsureMultipartUpload ensures a multipart upload exists, creating if necessary
func (u *S3Uploader) EnsureMultipartUpload(ctx context.Context, bucket, key string) (string, error) {
	input := &s3.CreateMultipartUploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	out, err := u.s3Client.CreateMultipartUpload(ctx, input)
	if err != nil {
		return "", fmt.Errorf("create multipart upload failed: %w", err)
	}
	return aws.ToString(out.UploadId), nil
}

// UploadPartStream uploads a single part with streaming support
func (u *S3Uploader) UploadPartStream(ctx context.Context, bucket, key, uploadID string, partNumber int32, reader io.Reader, contentLength int64) (string, error) {
	// Create timeout context for streaming upload
	timeoutDuration := max(time.Duration(contentLength/(10*1024*1024)+6)*time.Minute, 8*time.Minute)
	uploadCtx, cancel := context.WithTimeout(ctx, timeoutDuration)
	defer cancel()

	input := &s3.UploadPartInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		UploadId:      aws.String(uploadID),
		PartNumber:    aws.Int32(partNumber),
		Body:          reader,
		ContentLength: aws.Int64(contentLength),
	}

	log.Printf("Uploading part %d (streaming, size=%d MB)",
		partNumber, contentLength/(1024*1024))
	uploadStart := time.Now()

	// Use s3NoRetry client since we're streaming and can't retry
	out, err := u.s3NoRetry.UploadPart(uploadCtx, input)

	if err == nil {
		uploadDuration := time.Since(uploadStart)
		uploadSpeed := float64(contentLength) / (1024 * 1024) / uploadDuration.Seconds()
		log.Printf("Uploaded part %d in %v (%.1f MB/s)", partNumber, uploadDuration, uploadSpeed)
		return aws.ToString(out.ETag), nil
	}

	// Log error - streaming uploads cannot be retried
	logError("UploadPart (streaming)", err, map[string]interface{}{
		"bucket":         bucket,
		"key":            key,
		"part_number":    partNumber,
		"content_length": contentLength,
	})

	return "", fmt.Errorf("upload part %d failed: %w", partNumber, err)
}

// CompleteMultipartUpload completes a multipart upload
func (u *S3Uploader) CompleteMultipartUpload(ctx context.Context, bucket, key, uploadID string, parts []CompletedPart) (*UploadResult, error) {
	// Convert to S3 types
	s3Parts := make([]s3types.CompletedPart, len(parts))
	for i, part := range parts {
		s3Parts[i] = s3types.CompletedPart{
			PartNumber: aws.Int32(part.PartNumber),
			ETag:       aws.String(part.ETag),
		}
	}

	out, err := u.s3Client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadID),
		MultipartUpload: &s3types.CompletedMultipartUpload{
			Parts: s3Parts,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("complete multipart upload failed: %w", err)
	}

	return &UploadResult{
		ETag:         aws.ToString(out.ETag),
		LastModified: time.Now().UTC(),
	}, nil
}

// AbortMultipartUpload cancels a multipart upload
func (u *S3Uploader) AbortMultipartUpload(ctx context.Context, bucket, key, uploadID string) error {
	_, err := u.s3Client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadID),
	})
	return err
}

// HTTPUploader implements Uploader using direct HTTP client for R2
type HTTPUploader struct {
	endpoint       string
	httpClient     *http.Client
	stats          *performanceStats
	accessKey      string
	secretKey      string
	maxConcurrency int
}

// partInfo represents a completed part in a multipart upload
type partInfo struct {
	PartNumber int    `xml:"PartNumber"`
	ETag       string `xml:"ETag"`
}

// NewHTTPUploader creates a new HTTPUploader for direct R2 uploads
func NewHTTPUploader(endpoint string, httpClient *http.Client, accessKey, secretKey string, stats *performanceStats, maxConcurrency int) *HTTPUploader {
	return &HTTPUploader{
		endpoint:       endpoint,
		httpClient:     httpClient,
		stats:          stats,
		accessKey:      accessKey,
		secretKey:      secretKey,
		maxConcurrency: maxConcurrency,
	}
}

// UploadSmallFile implements direct upload for small files using HTTP
func (u *HTTPUploader) UploadSmallFile(ctx context.Context, bucket, key string, data []byte) (*UploadResult, error) {
	// URL encode the key to handle special characters
	encodedKey := url.PathEscape(key)
	uploadURL := fmt.Sprintf("%s/%s/%s", u.endpoint, bucket, encodedKey)

	req, err := http.NewRequestWithContext(ctx, "PUT", uploadURL, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))
	u.signRequest(req, data)

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upload failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, body)
	}

	etag := resp.Header.Get("ETag")
	return &UploadResult{
		ETag:         etag,
		LastModified: time.Now().UTC(),
	}, nil
}

// UploadLargeFile implements multipart upload for large files using HTTP
func (u *HTTPUploader) UploadLargeFile(ctx context.Context, bucket, key string, reader io.Reader, contentLength int64) (*UploadResult, error) {
	// Track performance
	uploadStart := time.Now()
	u.stats.incrementActive()
	defer u.stats.decrementActive()

	// Step 1: Initiate multipart upload
	uploadID, err := u.initiateMultipartUpload(ctx, bucket, key)
	if err != nil {
		return nil, err
	}

	log.Printf("HTTP multipart upload initiated: %s", uploadID)

	// Determine optimal chunk size
	const minChunkSize = 5 * 1024 * 1024   // 5MB minimum
	const maxChunkSize = 100 * 1024 * 1024 // 100MB maximum
	const targetChunks = 100               // Target number of chunks

	chunkSize := contentLength / targetChunks
	if chunkSize < minChunkSize {
		chunkSize = minChunkSize
	} else if chunkSize > maxChunkSize {
		chunkSize = maxChunkSize
	}

	// Use channels to coordinate concurrent uploads
	type uploadJob struct {
		partNum int
		data    []byte
	}

	type uploadResult struct {
		partNum int
		etag    string
		err     error
	}

	jobChan := make(chan uploadJob, u.maxConcurrency)
	resultChan := make(chan uploadResult, u.maxConcurrency)

	// Start worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < u.maxConcurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobChan {
				etag, err := u.uploadPart(ctx, bucket, key, uploadID, job.partNum, job.data, workerID)
				resultChan <- uploadResult{
					partNum: job.partNum,
					etag:    etag,
					err:     err,
				}
			}
		}(i)
	}

	// Result collector goroutine
	parts := make([]partInfo, 0)
	partsMu := sync.Mutex{}
	uploadErrors := make([]error, 0)

	go func() {
		for result := range resultChan {
			if result.err != nil {
				partsMu.Lock()
				uploadErrors = append(uploadErrors, result.err)
				partsMu.Unlock()
			} else {
				partsMu.Lock()
				parts = append(parts, partInfo{
					PartNumber: result.partNum,
					ETag:       result.etag,
				})
				partsMu.Unlock()
			}
		}
	}()

	// Read and distribute chunks to workers
	partNum := 1
	totalRead := int64(0)

	for totalRead < contentLength {
		// Calculate chunk size for this iteration
		remainingBytes := contentLength - totalRead
		currentChunkSize := int64(chunkSize)
		if remainingBytes < currentChunkSize {
			currentChunkSize = remainingBytes
		}

		// Read chunk
		chunk := make([]byte, currentChunkSize)
		n, err := io.ReadFull(reader, chunk)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Printf("Read error: %v", err)
			close(jobChan)
			wg.Wait()
			close(resultChan)

			// Abort the upload
			u.abortMultipartUpload(ctx, bucket, key, uploadID)
			return nil, fmt.Errorf("read error: %w", err)
		}

		// Trim to actual size read
		chunk = chunk[:n]
		totalRead += int64(n)

		// Send to worker
		jobChan <- uploadJob{
			partNum: partNum,
			data:    chunk,
		}

		partNum++

		// Break if we've read everything
		if n == 0 || err == io.EOF {
			break
		}
	}

	// Close job channel and wait for workers to finish
	close(jobChan)
	wg.Wait()
	close(resultChan)

	// Check for errors
	if len(uploadErrors) > 0 {
		u.abortMultipartUpload(ctx, bucket, key, uploadID)
		return nil, fmt.Errorf("upload failed: %w", uploadErrors[0])
	}

	// Sort parts by part number
	sort.Slice(parts, func(i, j int) bool {
		return parts[i].PartNumber < parts[j].PartNumber
	})

	// Complete multipart upload
	etag, err := u.completeMultipartUpload(ctx, bucket, key, uploadID, parts)
	if err != nil {
		return nil, err
	}

	// Record performance statistics
	uploadDuration := time.Since(uploadStart)
	u.stats.recordUpload(contentLength, uploadDuration, len(parts), u.maxConcurrency)
	throughput := float64(contentLength) / (1024 * 1024) / uploadDuration.Seconds()
	log.Printf("HTTP upload completed: %s, size=%d MB, duration=%v, throughput=%.1f MB/s",
		key, contentLength/(1024*1024), uploadDuration, throughput)

	return &UploadResult{
		ETag:         etag,
		LastModified: time.Now().UTC(),
	}, nil
}

// EnsureMultipartUpload ensures a multipart upload exists, creating if necessary
func (u *HTTPUploader) EnsureMultipartUpload(ctx context.Context, bucket, key string) (string, error) {
	return u.initiateMultipartUpload(ctx, bucket, key)
}

// UploadPartStream uploads a single part with streaming support
func (u *HTTPUploader) UploadPartStream(ctx context.Context, bucket, key, uploadID string, partNumber int32, reader io.Reader, contentLength int64) (string, error) {
	// URL encode the key to handle special characters
	// Use PathEscape for S3-compatible URL encoding
	encodedKey := url.PathEscape(key)
	uploadURL := fmt.Sprintf("%s/%s/%s?partNumber=%d&uploadId=%s",
		u.endpoint, bucket, encodedKey, partNumber, uploadID)

	log.Printf("HTTP Upload: part %d, uploadID=%s, url=%s, size=%d MB",
		partNumber, uploadID, uploadURL, contentLength/(1024*1024))

	// For streaming uploads, we can only try once since we can't rewind the reader
	req, err := http.NewRequestWithContext(ctx, "PUT", uploadURL, reader)
	if err != nil {
		return "", fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Length", fmt.Sprintf("%d", contentLength))
	req.ContentLength = contentLength // Add this line to fix 411 error

	// For streaming, we need to handle the signing differently
	// since we can't read the full body to calculate the hash
	u.signRequestStreaming(req, contentLength)

	start := time.Now()
	resp, err := u.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("upload part %d failed: %w", partNumber, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload part %d failed with status %d: %s", partNumber, resp.StatusCode, body)
	}

	duration := time.Since(start)
	speed := float64(contentLength) / (1024 * 1024) / duration.Seconds()
	log.Printf("Uploaded part %d in %v (%.1f MB/s) via HTTP", partNumber, duration, speed)

	return resp.Header.Get("ETag"), nil
}

// CompleteMultipartUpload completes a multipart upload
func (u *HTTPUploader) CompleteMultipartUpload(ctx context.Context, bucket, key, uploadID string, parts []CompletedPart) (*UploadResult, error) {
	// Convert to partInfo for XML marshaling
	xmlParts := make([]partInfo, len(parts))
	for i, part := range parts {
		xmlParts[i] = partInfo{
			PartNumber: int(part.PartNumber),
			ETag:       part.ETag,
		}
	}

	etag, err := u.completeMultipartUpload(ctx, bucket, key, uploadID, xmlParts)
	if err != nil {
		return nil, err
	}

	return &UploadResult{
		ETag:         etag,
		LastModified: time.Now().UTC(),
	}, nil
}

// AbortMultipartUpload cancels a multipart upload
func (u *HTTPUploader) AbortMultipartUpload(ctx context.Context, bucket, key, uploadID string) error {
	u.abortMultipartUpload(ctx, bucket, key, uploadID)
	return nil
}

// signRequestStreaming adds AWS v4 signature for streaming uploads
func (u *HTTPUploader) signRequestStreaming(req *http.Request, contentLength int64) {
	now := time.Now().UTC()
	dateStr := now.Format("20060102T150405Z")
	dateShort := now.Format("20060102")
	region := "auto" // R2 uses "auto" region

	// For streaming, use UNSIGNED-PAYLOAD
	payloadHash := "UNSIGNED-PAYLOAD"
	req.Header.Set("X-Amz-Date", dateStr)
	req.Header.Set("X-Amz-Content-Sha256", payloadHash)
	req.Header.Set("Host", req.Host)

	// Create signing key
	key := []byte("AWS4" + u.secretKey)
	dateKey := hmacSHA256(key, []byte(dateShort))
	regionKey := hmacSHA256(dateKey, []byte(region))
	serviceKey := hmacSHA256(regionKey, []byte("s3"))
	signingKey := hmacSHA256(serviceKey, []byte("aws4_request"))

	// Build canonical request
	canonicalHeaders := fmt.Sprintf("host:%s\nx-amz-content-sha256:%s\nx-amz-date:%s\n",
		req.Host, payloadHash, dateStr)
	signedHeaders := "host;x-amz-content-sha256;x-amz-date"

	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		req.Method,
		req.URL.Path,
		req.URL.RawQuery,
		canonicalHeaders,
		signedHeaders,
		payloadHash)

	// String to sign
	scope := fmt.Sprintf("%s/%s/s3/aws4_request", dateShort, region)
	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s",
		dateStr, scope, sha256sum([]byte(canonicalRequest)))

	// Calculate signature
	signature := hex.EncodeToString(hmacSHA256(signingKey, []byte(stringToSign)))

	// Set authorization header
	req.Header.Set("Authorization", fmt.Sprintf(
		"AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		u.accessKey, scope, signedHeaders, signature))
}

// Helper methods for HTTPUploader

// signRequest adds AWS v4 signature to the request
func (u *HTTPUploader) signRequest(req *http.Request, payload []byte) {
	now := time.Now().UTC()
	dateStr := now.Format("20060102T150405Z")
	dateShort := now.Format("20060102")
	region := "auto" // R2 uses "auto" region

	// Headers required by AWS v4
	payloadHash := sha256sum(payload)
	req.Header.Set("X-Amz-Date", dateStr)
	req.Header.Set("X-Amz-Content-Sha256", payloadHash)
	req.Header.Set("Host", req.Host)

	// Create signing key
	key := []byte("AWS4" + u.secretKey)
	dateKey := hmacSHA256(key, []byte(dateShort))
	regionKey := hmacSHA256(dateKey, []byte(region))
	serviceKey := hmacSHA256(regionKey, []byte("s3"))
	signingKey := hmacSHA256(serviceKey, []byte("aws4_request"))

	// Build canonical request
	canonicalHeaders := fmt.Sprintf("host:%s\nx-amz-content-sha256:%s\nx-amz-date:%s\n",
		req.Host, payloadHash, dateStr)
	signedHeaders := "host;x-amz-content-sha256;x-amz-date"

	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		req.Method,
		req.URL.Path,
		req.URL.RawQuery,
		canonicalHeaders,
		signedHeaders,
		payloadHash)

	// String to sign
	scope := fmt.Sprintf("%s/%s/s3/aws4_request", dateShort, region)
	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s",
		dateStr, scope, sha256sum([]byte(canonicalRequest)))

	// Calculate signature
	signature := hex.EncodeToString(hmacSHA256(signingKey, []byte(stringToSign)))

	// Set authorization header
	req.Header.Set("Authorization", fmt.Sprintf(
		"AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		u.accessKey, scope, signedHeaders, signature))
}

// initiateMultipartUpload starts a new multipart upload
func (u *HTTPUploader) initiateMultipartUpload(ctx context.Context, bucket, key string) (string, error) {
	// URL encode the key to handle special characters
	// Use PathEscape for S3-compatible URL encoding
	encodedKey := url.PathEscape(key)
	uploadURL := fmt.Sprintf("%s/%s/%s?uploads=", u.endpoint, bucket, encodedKey)

	log.Printf("HTTP: Initiating multipart upload for %s/%s at %s", bucket, key, uploadURL)

	req, err := http.NewRequestWithContext(ctx, "POST", uploadURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request failed: %w", err)
	}

	u.signRequest(req, nil)

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("initiate multipart upload failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("initiate failed with status %d: %s", resp.StatusCode, body)
	}

	// Parse upload ID from response
	var result struct {
		UploadID string `xml:"UploadId"`
	}
	if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	log.Printf("HTTP: Created multipart upload ID: %s", result.UploadID)
	return result.UploadID, nil
}

// uploadPart uploads a single part
func (u *HTTPUploader) uploadPart(ctx context.Context, bucket, key, uploadID string, partNum int, data []byte, workerID int) (string, error) {
	// URL encode the key to handle special characters
	encodedKey := url.PathEscape(key)
	uploadURL := fmt.Sprintf("%s/%s/%s?partNumber=%d&uploadId=%s",
		u.endpoint, bucket, encodedKey, partNum, uploadID)

	log.Printf("[Worker %d] Uploading part %d (size=%d MB)...",
		workerID, partNum, len(data)/(1024*1024))

	// Retry logic
	for attempt := 1; attempt <= 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, "PUT", uploadURL, bytes.NewReader(data))
		if err != nil {
			return "", fmt.Errorf("create request failed: %w", err)
		}

		req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))
		u.signRequest(req, data)

		start := time.Now()
		resp, err := u.httpClient.Do(req)
		if err != nil {
			if attempt < 3 {
				log.Printf("[Worker %d] Part %d upload attempt %d failed: %v, retrying...",
					workerID, partNum, attempt, err)
				time.Sleep(time.Duration(attempt) * time.Second)
				continue
			}
			return "", fmt.Errorf("upload part %d failed after %d attempts: %w", partNum, attempt, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			if attempt < 3 {
				log.Printf("[Worker %d] Part %d upload attempt %d failed with status %d: %s, retrying...",
					workerID, partNum, attempt, resp.StatusCode, body)
				time.Sleep(time.Duration(attempt) * time.Second)
				continue
			}
			return "", fmt.Errorf("upload part %d failed with status %d: %s", partNum, resp.StatusCode, body)
		}

		duration := time.Since(start)
		speed := float64(len(data)) / (1024 * 1024) / duration.Seconds()
		log.Printf("[Worker %d] Uploaded part %d in %v (%.1f MB/s)",
			workerID, partNum, duration, speed)

		return resp.Header.Get("ETag"), nil
	}

	return "", fmt.Errorf("upload part %d failed after all retries", partNum)
}

// completeMultipartUpload finishes the multipart upload
func (u *HTTPUploader) completeMultipartUpload(ctx context.Context, bucket, key, uploadID string, parts []partInfo) (string, error) {
	// Build XML payload
	completeXML, err := xml.Marshal(struct {
		XMLName xml.Name   `xml:"CompleteMultipartUpload"`
		Parts   []partInfo `xml:"Part"`
	}{Parts: parts})
	if err != nil {
		return "", fmt.Errorf("failed to marshal XML: %w", err)
	}

	url := fmt.Sprintf("%s/%s/%s?uploadId=%s", u.endpoint, bucket, key, uploadID)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(completeXML))
	if err != nil {
		return "", fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(completeXML)))
	u.signRequest(req, completeXML)

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("complete multipart upload failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("complete failed with status %d: %s", resp.StatusCode, body)
	}

	return resp.Header.Get("ETag"), nil
}

// abortMultipartUpload cancels a multipart upload
func (u *HTTPUploader) abortMultipartUpload(ctx context.Context, bucket, key, uploadID string) {
	url := fmt.Sprintf("%s/%s/%s?uploadId=%s", u.endpoint, bucket, key, uploadID)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		log.Printf("Failed to create abort request: %v", err)
		return
	}

	u.signRequest(req, nil)

	resp, err := u.httpClient.Do(req)
	if err != nil {
		log.Printf("Failed to abort multipart upload: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 && resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Abort multipart upload returned status %d: %s", resp.StatusCode, body)
	}
}

// Helper functions for AWS v4 signing
func sha256sum(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func hmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}
