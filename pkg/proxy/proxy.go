package proxy

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
)

// Optimized HTTP client for blob storage uploads (based on ASUR)
var blobStorageHTTPClient = &http.Client{
	Timeout: 30 * time.Minute,
	Transport: &http.Transport{
		MaxIdleConns:          1000,
		MaxIdleConnsPerHost:   500,
		MaxConnsPerHost:       100,
		IdleConnTimeout:       90 * time.Second,
		DisableCompression:    true,
		WriteBufferSize:       1024 * 1024,
		ReadBufferSize:        1024 * 1024,
		ResponseHeaderTimeout: 5 * time.Minute,
		ExpectContinueTimeout: 30 * time.Second,
	},
}

var cacheStore = sync.Map{}

type CacheEntryData struct {
	CacheKey               string
	CacheVersion           string
	Mutex                  sync.Mutex
	Reserved               bool
	BackendReserveResponse ReserveCacheResponse

	// S3/R2 streaming state
	UploadedParts     map[S3PartNumber]S3CompletedPart // for ordering
	InferredChunkSize int64                            // Inferred from first chunk

	// GCS streaming state
	GCSClient *storage.Client
	GCSWriter *storage.Writer

	// Azure streaming state
	AzureBlocks map[int64]string // offset → blockID for ordering
}

func GetCache(ctx context.Context, input DockerGHAGetCacheRequest) (*DockerGHAGetCacheResponse, error) {
	primaryKey := input.Keys[0]
	// Docker backend weirdly sends impartial key as primary key sometimes.
	restoreKeys := input.Keys

	payload := GetCacheRequest{
		CacheKey:     primaryKey,
		CacheVersion: input.Version,
		RestoreKeys:  restoreKeys,
	}

	cacheResponse, err := callCacheBackend[GetCacheResponse](ctx, CacheBackendRequest{
		Path: "/v1/cache/get",
		Body: payload,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get cache: %w", err)
	}

	dockerGetResponse := DockerGHAGetCacheResponse{
		CacheKey: cacheResponse.CacheEntry.CacheUserGivenKey,
	}

	presignedURL := ""
	switch cacheResponse.Provider {
	case ProviderS3:
	case ProviderR2:
		presignedURL = cacheResponse.S3.PreSignedURL
	case ProviderGCS:
		presignedURL = cacheResponse.GCS.PreSignedURL
	case ProviderAzureBlob:
		presignedURL = cacheResponse.AzureBlob.PreSignedURL
	}

	if cacheResponse.CacheEntry != nil {
		dockerGetResponse.ArchiveLocation = presignedURL
	}

	return &dockerGetResponse, nil
}

func ReserveCache(ctx context.Context, input DockerGHAReserveCacheRequest) (*DockerGHAReserveCacheResponse, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomCacheID := r.Intn(1000000)

	cacheStore.Store(randomCacheID, &CacheEntryData{
		CacheKey:      input.Key,
		CacheVersion:  input.Version,
		Reserved:      false,
		UploadedParts: make(map[S3PartNumber]S3CompletedPart),
		AzureBlocks:   make(map[int64]string),
	})

	return &DockerGHAReserveCacheResponse{CacheID: randomCacheID}, nil
}

func UploadCache(ctx context.Context, input DockerGHAUploadCacheRequest) (*DockerGHAUploadCacheResponse, error) {
	cacheData, ok := cacheStore.Load(input.CacheID)
	if !ok {
		return nil, fmt.Errorf("cache ID not found")
	}
	cacheEntry := cacheData.(*CacheEntryData)

	cacheEntry.Mutex.Lock()
	defer cacheEntry.Mutex.Unlock()

	start, end, err := parseContentRange(input.ContentRange)
	if err != nil {
		return nil, fmt.Errorf("invalid Content-Range header: %w", err)
	}

	// Infer chunk size and reserve backend on first UploadCache call
	if !cacheEntry.Reserved {
		cacheEntry.InferredChunkSize = end - start + 1

		if err := ensureBackendReserved(ctx, cacheEntry, cacheEntry.InferredChunkSize); err != nil {
			return nil, fmt.Errorf("failed to reserve backend: %w", err)
		}
	}

	switch cacheEntry.BackendReserveResponse.Provider {
	case ProviderS3, ProviderR2:
		// Calculate logical part number from offset (handles out-of-order)
		partNum := S3PartNumber((start / cacheEntry.InferredChunkSize) + 1)

		if _, exists := cacheEntry.UploadedParts[partNum]; exists {
			fmt.Printf("Part %d (offset %d) already uploaded, skipping\n", partNum, start)
			return &DockerGHAUploadCacheResponse{}, nil
		}

		// Since we only generate 100 URLs initially, we need to check if there are more parts to upload.
		// If we need more, we get more presigned URLs.
		if err := ensureURLForPart(ctx, cacheEntry, partNum); err != nil {
			return nil, fmt.Errorf("failed to ensure presigned URL: %w", err)
		}

		dataCopy := append([]byte{}, input.Content...)

		if err := uploadPartToS3(ctx, cacheEntry, partNum, dataCopy); err != nil {
			return nil, fmt.Errorf("failed to upload part to S3: %w", err)
		}

		return &DockerGHAUploadCacheResponse{}, nil

	case ProviderGCS:
		dataCopy := append([]byte{}, input.Content...)
		if err := uploadChunkGCS(ctx, cacheEntry, dataCopy); err != nil {
			return nil, fmt.Errorf("failed to upload chunk to GCS: %w", err)
		}
		return &DockerGHAUploadCacheResponse{}, nil

	case ProviderAzureBlob:
		dataCopy := append([]byte{}, input.Content...)
		if err := uploadChunkAzure(ctx, cacheEntry, start, dataCopy); err != nil {
			return nil, fmt.Errorf("failed to upload chunk to Azure: %w", err)
		}
		return &DockerGHAUploadCacheResponse{}, nil

	default:
		return nil, fmt.Errorf("unsupported provider: %s", cacheEntry.BackendReserveResponse.Provider)
	}
}

// Content-Range is of the form "bytes 1234-5678/*"
var contentRangeRegex = regexp.MustCompile(`^bytes (\d+)-(\d+)`)

func parseContentRange(contentRange string) (int64, int64, error) {
	matches := contentRangeRegex.FindStringSubmatch(contentRange)
	if len(matches) != 3 {
		return 0, 0, fmt.Errorf("invalid Content-Range format: %s", contentRange)
	}

	start, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid start offset: %v", err)
	}

	end, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid end offset: %v", err)
	}

	return start, end, nil
}

func ensureBackendReserved(ctx context.Context, entry *CacheEntryData, firstChunkSize int64) error {
	if entry.Reserved {
		return nil
	}

	chunkSize := firstChunkSize
	if chunkSize < 5*1024*1024 {
		chunkSize = 32 * 1024 * 1024 // Assume standard for small final chunks
	}
	entry.InferredChunkSize = chunkSize

	const initialURLs = 100 // Covers ~3.2GB with 32MB chunks

	payload := ReserveCacheRequest{
		CacheKey:       entry.CacheKey,
		CacheVersion:   entry.CacheVersion,
		NumberOfChunks: initialURLs,
		ContentType:    "application/octet-stream",
	}

	reserveResp, err := callCacheBackend[ReserveCacheResponse](ctx, CacheBackendRequest{
		Path: "/v1/cache/reserve",
		Body: payload,
	})
	if err != nil {
		return fmt.Errorf("backend reserve failed: %w", err)
	}

	entry.BackendReserveResponse = *reserveResp
	entry.Reserved = true
	return nil
}

func ensureURLForPart(ctx context.Context, entry *CacheEntryData, partNum S3PartNumber) error {
	currentURLs := len(entry.BackendReserveResponse.S3.PreSignedURLs)

	if int(partNum) <= currentURLs {
		return nil
	}

	// Need more URLs - request another batch
	extendResp, err := callCacheBackend[ExtendReserveCacheResponse](ctx, CacheBackendRequest{
		Path: "/v1/cache/reserve/extend",
		Body: ExtendReserveCacheRequest{
			CacheKey:     entry.CacheKey,
			CacheVersion: entry.CacheVersion,
			UploadID:     entry.BackendReserveResponse.S3.UploadID,
			FromPart:     int32(currentURLs + 1),
			Count:        100,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to extend presigned URLs: %w", err)
	}

	entry.BackendReserveResponse.S3.PreSignedURLs = append(
		entry.BackendReserveResponse.S3.PreSignedURLs,
		extendResp.PreSignedURLs...,
	)

	fmt.Printf("Extended presigned URLs: now have %d URLs (requested from part %d)\n",
		len(entry.BackendReserveResponse.S3.PreSignedURLs), currentURLs+1)

	return nil
}

func uploadPartToS3(ctx context.Context, entry *CacheEntryData, partNum S3PartNumber, data []byte) error {
	// Get presigned URL for this part
	if int(partNum-1) >= len(entry.BackendReserveResponse.S3.PreSignedURLs) {
		return fmt.Errorf("part number %d exceeds available presigned URLs", partNum)
	}
	presignedURL := entry.BackendReserveResponse.S3.PreSignedURLs[partNum-1]

	// Calculate timeout based on data size (ASUR pattern)
	sizeMB := len(data) / (1024 * 1024)
	timeoutMin := max((sizeMB/10)+6, 8)
	uploadCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutMin)*time.Minute)
	defer cancel()

	uploadStart := time.Now()
	var etag string
	var lastErr error

	for attempt := range maxRetries {
		req, err := http.NewRequestWithContext(uploadCtx, http.MethodPut, presignedURL,
			bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("Content-Length", strconv.Itoa(len(data)))

		fmt.Printf("Uploading part %d (size=%d MB) - attempt %d/%d\n", partNum, sizeMB, attempt+1, maxRetries)

		resp, err := blobStorageHTTPClient.Do(req)
		if err == nil && (resp.StatusCode == 200 || resp.StatusCode == 204) {
			etag = resp.Header.Get("ETag")
			resp.Body.Close()

			duration := time.Since(uploadStart)
			speed := float64(len(data)) / (1024 * 1024) / duration.Seconds()
			fmt.Printf("Uploaded part %d in %v (%.1f MB/s)\n", partNum, duration, speed)
			break
		}

		lastErr = err
		if resp != nil {
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d", resp.StatusCode)
		}
		fmt.Printf("Part %d upload failed (attempt %d/%d): %v\n", partNum, attempt+1, maxRetries, lastErr)

		if attempt < maxRetries-1 {
			backoff := min(initialBackoff*(1<<attempt), maxBackoff)
			time.Sleep(backoff)
		}
	}

	if etag == "" {
		return fmt.Errorf("failed to upload part %d after 3 attempts: %v", partNum, lastErr)
	}

	entry.UploadedParts[partNum] = S3CompletedPart{
		ETag:       &etag,
		PartNumber: partNum,
	}

	return nil
}

func uploadChunkGCS(ctx context.Context, entry *CacheEntryData, data []byte) error {
	if entry.GCSWriter == nil {
		if entry.BackendReserveResponse.GCS.ShortLivedToken == nil {
			return fmt.Errorf("no GCS short-lived token found")
		}

		creds := option.WithTokenSource(oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: entry.BackendReserveResponse.GCS.ShortLivedToken.AccessToken,
		}))
		client, err := storage.NewClient(ctx, creds)
		if err != nil {
			return fmt.Errorf("failed to create GCS client: %w", err)
		}

		entry.GCSClient = client

		bucket := client.Bucket(entry.BackendReserveResponse.GCS.BucketName)
		object := bucket.Object(entry.BackendReserveResponse.GCS.CacheKey)
		entry.GCSWriter = object.NewWriter(ctx)
	}

	_, err := entry.GCSWriter.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write chunk to GCS: %w", err)
	}

	return nil
}

func uploadChunkAzure(ctx context.Context, entry *CacheEntryData, offset int64, data []byte) error {
	blockNum := (offset / entry.InferredChunkSize) + 1
	blockID := base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "block-%06d", blockNum))

	if _, exists := entry.AzureBlocks[offset]; exists {
		fmt.Printf("Azure block at offset %d already uploaded, skipping\n", offset)
		return nil
	}

	baseURL := entry.BackendReserveResponse.AzureBlob.PreSignedURL
	blockURL := fmt.Sprintf("%s&comp=block&blockid=%s", baseURL, url.QueryEscape(blockID))

	var lastErr error
	for attempt := range maxRetries {
		req, err := http.NewRequestWithContext(ctx, http.MethodPut, blockURL, bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("failed to create Azure block request: %w", err)
		}
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("Content-Length", strconv.Itoa(len(data)))

		resp, err := blobStorageHTTPClient.Do(req)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			resp.Body.Close()
			entry.AzureBlocks[offset] = blockID
			fmt.Printf("Uploaded Azure block %s (offset %d)\n", blockID, offset)
			return nil
		}

		lastErr = err
		if resp != nil {
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d", resp.StatusCode)
		}

		if attempt < maxRetries-1 {
			backoff := min(initialBackoff*(1<<attempt), maxBackoff)
			time.Sleep(backoff)
		}
	}

	return fmt.Errorf("failed to upload Azure block after %d attempts: %v", maxRetries, lastErr)
}

func commitAzureBlockList(ctx context.Context, entry *CacheEntryData) error {
	// Sort offsets to get correct block order
	offsets := make([]int64, 0, len(entry.AzureBlocks))
	for offset := range entry.AzureBlocks {
		offsets = append(offsets, offset)
	}
	sort.Slice(offsets, func(i, j int) bool {
		return offsets[i] < offsets[j]
	})

	// Build ordered block ID list
	blockIDs := make([]string, len(offsets))
	for i, offset := range offsets {
		blockIDs[i] = entry.AzureBlocks[offset]
	}

	blockList := blockListXML{Latest: blockIDs}
	xmlBytes, err := xml.Marshal(blockList)
	if err != nil {
		return fmt.Errorf("failed to marshal block list XML: %w", err)
	}

	commitURL := entry.BackendReserveResponse.AzureBlob.PreSignedURL + "&comp=blocklist"
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, commitURL, bytes.NewReader(xmlBytes))
	if err != nil {
		return fmt.Errorf("failed to create block list request: %w", err)
	}
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Content-Length", strconv.Itoa(len(xmlBytes)))

	resp, err := blobStorageHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to commit Azure block list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("azure block list commit failed with status %d", resp.StatusCode)
	}

	fmt.Printf("Committed Azure block list with %d blocks\n", len(blockIDs))
	return nil
}

const (
	maxRetries     = 3
	initialBackoff = 1 * time.Second
	maxBackoff     = 16 * time.Second
)

func CommitCache(ctx context.Context, input DockerGHACommitCacheRequest) (*DockerGHACommitCacheResponse, error) {
	cacheData, ok := cacheStore.Load(input.CacheID)
	if !ok {
		return nil, fmt.Errorf("cache ID not found")
	}
	cacheEntry := cacheData.(*CacheEntryData)

	cacheEntry.Mutex.Lock()
	defer cacheEntry.Mutex.Unlock()

	var payload CommitCacheRequest
	switch cacheEntry.BackendReserveResponse.Provider {
	case ProviderS3, ProviderR2:
		parts := make([]S3CompletedPart, 0, len(cacheEntry.UploadedParts))
		for _, part := range cacheEntry.UploadedParts {
			parts = append(parts, part)
		}
		slices.SortFunc(parts, func(a, b S3CompletedPart) int {
			return int(a.PartNumber - b.PartNumber)
		})

		payload = CommitCacheRequest{
			CacheKey:     cacheEntry.CacheKey,
			CacheVersion: cacheEntry.CacheVersion,
			UploadKey:    cacheEntry.BackendReserveResponse.S3.UploadKey,
			UploadID:     cacheEntry.BackendReserveResponse.S3.UploadID,
			Parts:        parts,
			VCSType:      "github",
		}

	case ProviderGCS:
		if cacheEntry.GCSWriter != nil {
			if err := cacheEntry.GCSWriter.Close(); err != nil {
				return nil, fmt.Errorf("failed to close GCS writer: %w", err)
			}
		}
		if cacheEntry.GCSClient != nil {
			if err := cacheEntry.GCSClient.Close(); err != nil {
				return nil, fmt.Errorf("failed to close GCS client: %w", err)
			}
		}

		payload = CommitCacheRequest{
			CacheKey:     cacheEntry.CacheKey,
			CacheVersion: cacheEntry.CacheVersion,
			Parts:        []S3CompletedPart{},
			VCSType:      "github",
		}

	case ProviderAzureBlob:
		if err := commitAzureBlockList(ctx, cacheEntry); err != nil {
			return nil, fmt.Errorf("failed to commit Azure block list: %w", err)
		}

		payload = CommitCacheRequest{
			CacheKey:     cacheEntry.CacheKey,
			CacheVersion: cacheEntry.CacheVersion,
			Parts:        []S3CompletedPart{},
			VCSType:      "github",
		}

	default:
		return nil, fmt.Errorf("unsupported provider: %s", cacheEntry.BackendReserveResponse.Provider)
	}

	_, err := callCacheBackend[CommitCacheResponse](ctx, CacheBackendRequest{
		Path: "/v1/cache/commit",
		Body: payload,
	})
	if err != nil {
		return nil, fmt.Errorf("backend commit failed: %w", err)
	}

	cacheStore.Delete(input.CacheID)

	return &DockerGHACommitCacheResponse{}, nil
}
