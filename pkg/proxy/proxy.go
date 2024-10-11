package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var cacheStore = sync.Map{}

type ChunkData struct {
	StartOffset int64
	EndOffset   int64
	Content     []byte
}

type CacheEntryData struct {
	BackendReserveResponse ReserveCacheResponse
	S3Parts                []S3CompletedPart
	CacheKey               string
	CacheVersion           string
	Chunks                 map[int64]ChunkData
	Mutex                  sync.Mutex // Mutex to protect access to chunks
}

func GetCache(ctx context.Context, input DockerGHAGetCacheRequest) (*DockerGHAGetCacheResponse, error) {
	requestURL := fmt.Sprintf("%s/v1/cache/get", input.HostURL)

	primaryKey := input.Keys[0]
	// Docker backend weirdly sends impartial key as primary key sometimes.
	restoreKeys := input.Keys

	payload := GetCacheRequest{
		CacheKey:     primaryKey,
		CacheVersion: input.Version,
		RestoreKeys:  restoreKeys,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	serviceURL, err := url.Parse(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service URL: %w", err)
	}

	agent := fiber.Post(serviceURL.String())

	agent.Body(payloadBytes)

	agent.Add("Content-Type", "application/json")
	agent.Add("Accept", "application/json")
	agent.Add("Authorization", fmt.Sprintf("Bearer %s", input.AuthToken))

	statusCode, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to send request to cache backend: %v", errs)
	}

	if statusCode < 200 || statusCode >= 300 {
		return nil, fmt.Errorf("failed to get cache: %s", string(body))
	}

	var cacheResponse GetCacheResponse
	if err := json.Unmarshal(body, &cacheResponse); err != nil {
		return nil, fmt.Errorf("failed to parse backend response: %w", err)
	}

	dockerGetResponse := DockerGHAGetCacheResponse{
		CacheKey: cacheResponse.CacheEntry.CacheUserGivenKey,
	}

	presignedURL := ""
	switch cacheResponse.Provider {
	case ProviderS3:
		presignedURL = cacheResponse.S3.PreSignedURL
	case ProviderGCS:
		presignedURL = cacheResponse.GCS.PreSignedURL
	}

	if cacheResponse.CacheEntry != nil {
		dockerGetResponse.ArchiveLocation = presignedURL
	}

	return &dockerGetResponse, nil
}

func ReserveCache(ctx context.Context, input DockerGHAReserveCacheRequest) (*DockerGHAReserveCacheResponse, error) {
	requestURL := fmt.Sprintf("%s/v1/cache/reserve", input.HostURL)

	payload := ReserveCacheRequest{
		CacheKey:       input.Key,
		CacheVersion:   input.Version,
		NumberOfChunks: 1,
		ContentType:    "application/octet-stream",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	serviceURL, err := url.Parse(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service URL: %w", err)
	}

	agent := fiber.Post(serviceURL.String())

	agent.Body(payloadBytes)

	agent.Add("Content-Type", "application/json")
	agent.Add("Accept", "application/json")
	agent.Add("Authorization", fmt.Sprintf("Bearer %s", input.AuthToken))

	statusCode, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to send request to cache backend: %v", errs)
	}

	if statusCode < 200 || statusCode >= 300 {
		return nil, fmt.Errorf("failed to reserve cache: %s", string(body))
	}

	var reserveCacheResponse ReserveCacheResponse
	if err := json.Unmarshal(body, &reserveCacheResponse); err != nil {
		return nil, fmt.Errorf("failed to parse backend response: %w", err)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomCacheID := r.Intn(1000000)

	dockerReserveResponse := DockerGHAReserveCacheResponse{
		CacheID: randomCacheID,
	}

	// Save this cache ID for later use
	cacheStore.Store(randomCacheID, &CacheEntryData{
		BackendReserveResponse: reserveCacheResponse,
		CacheKey:               input.Key,
		CacheVersion:           input.Version,
		Chunks:                 make(map[int64]ChunkData),
	})

	return &dockerReserveResponse, nil
}

func UploadCache(ctx context.Context, input DockerGHAUploadCacheRequest) (*DockerGHAUploadCacheResponse, error) {
	cacheData, _ := cacheStore.LoadOrStore(input.CacheID, &CacheEntryData{
		Chunks: make(map[int64]ChunkData),
	})
	cacheEntry := cacheData.(*CacheEntryData)

	start, end, err := parseContentRange(input.ContentRange)
	if err != nil {
		return nil, fmt.Errorf("invalid Content-Range header: %w", err)
	}

	chunk := ChunkData{
		StartOffset: start,
		EndOffset:   end,
		Content:     append([]byte{}, input.Content...), // Make a copy of the content otherwise it gets corrupted
	}

	cacheEntry.Mutex.Lock()
	cacheEntry.Chunks[start] = chunk
	cacheEntry.Mutex.Unlock()

	return &DockerGHAUploadCacheResponse{}, nil
}

func parseContentRange(contentRange string) (int64, int64, error) {
	parts := strings.Split(contentRange, " ")
	if len(parts) != 2 || parts[0] != "bytes" {
		return 0, 0, fmt.Errorf("invalid Content-Range format: %s", contentRange)
	}

	// Extract the range part: "0-1023" or "0-1023/*"
	rangeParts := strings.Split(parts[1], "-")
	if len(rangeParts) != 2 {
		return 0, 0, fmt.Errorf("invalid Content-Range range section: %s", parts[1])
	}

	start, err := strconv.ParseInt(rangeParts[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid start offset: %s, error: %v", rangeParts[0], err)
	}

	// Handle the end offset part, accounting for missing total size
	endParts := strings.Split(rangeParts[1], "/")
	end, err := strconv.ParseInt(endParts[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid end offset: %s, error: %v", endParts[0], err)
	}

	return start, end, nil
}

const (
	maxRetries     = 5
	initialBackoff = 1 * time.Second
	maxBackoff     = 16 * time.Second
)

func uploadToBlobStorage(ctx context.Context, cacheID int) (*DockerGHAUploadCacheResponse, error) {
	cacheData, ok := cacheStore.Load(cacheID)
	if !ok {
		return nil, fmt.Errorf("cache ID not found")
	}

	cacheEntry := cacheData.(*CacheEntryData)

	cacheEntry.Mutex.Lock()
	defer cacheEntry.Mutex.Unlock()

	// Reassemble chunks in the correct order
	var finalBuffer bytes.Buffer
	offsets := make([]int64, 0, len(cacheEntry.Chunks))
	for offset := range cacheEntry.Chunks {
		offsets = append(offsets, offset)
	}

	sort.Slice(offsets, func(i, j int) bool {
		return offsets[i] < offsets[j]
	})

	for _, offset := range offsets {
		finalBuffer.Write(cacheEntry.Chunks[offset].Content)
	}

	switch cacheEntry.BackendReserveResponse.Provider {
	case ProviderS3:
		if len(cacheEntry.BackendReserveResponse.S3.PreSignedURLs) != 1 {
			return nil, fmt.Errorf("no presigned URLs found")
		}

		s3PresignedURL := cacheEntry.BackendReserveResponse.S3.PreSignedURLs[0]

		for attempt := 0; attempt < maxRetries; attempt++ {
			contentReader := bytes.NewReader(finalBuffer.Bytes())
			req, err := http.NewRequestWithContext(ctx, http.MethodPut, s3PresignedURL, contentReader)
			if err != nil {
				return nil, fmt.Errorf("failed to create S3 request: %w", err)
			}
			req.Header.Set("Content-Type", "application/octet-stream")

			resp, err := http.DefaultClient.Do(req)
			if err == nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent) {
				defer resp.Body.Close()
				etag := resp.Header.Get("ETag")
				if etag == "" {
					return nil, fmt.Errorf("no ETag found in response")
				}

				partNumberPtr := int32(1)
				cacheEntry.S3Parts = []S3CompletedPart{{
					ETag:       &etag,
					PartNumber: &partNumberPtr,
				}}

				// Success, break out of retry loop
				return &DockerGHAUploadCacheResponse{}, nil
			}

			// If response is not OK, log and prepare for retry
			if resp != nil {
				defer resp.Body.Close()
				if attempt < maxRetries-1 {
					fmt.Printf("Retrying upload... attempt %d/%d, error: %v\n", attempt+1, maxRetries, err)
					time.Sleep(1 << attempt * time.Second) // Exponential backoff
				}
			}
		}

	case ProviderGCS:
		if cacheEntry.BackendReserveResponse.GCS.ShortLivedToken == nil {
			return nil, fmt.Errorf("no short lived token found")
		}

		creds := option.WithTokenSource(oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: cacheEntry.BackendReserveResponse.GCS.ShortLivedToken.AccessToken,
		}))
		client, err := storage.NewClient(ctx, creds)
		if err != nil {
			return nil, fmt.Errorf("failed to create GCS client: %w", err)
		}

		defer client.Close()

		bucket := client.Bucket(cacheEntry.BackendReserveResponse.GCS.BucketName)
		object := bucket.Object(cacheEntry.BackendReserveResponse.GCS.CacheKey)

		wc := object.NewWriter(ctx)

		for attempt := 0; attempt < maxRetries; attempt++ {
			_, err = wc.Write(finalBuffer.Bytes())
			if err == nil {
				err = wc.Close()
				if err == nil {
					return &DockerGHAUploadCacheResponse{}, nil
				}
			}

			if attempt < maxRetries-1 {
				fmt.Printf("Retrying upload... attempt %d/%d, error: %v\n", attempt+1, maxRetries, err)
				time.Sleep(1 << attempt * time.Second)
			}
		}

		return nil, fmt.Errorf("failed to upload to GCS: %w", err)

	default:
		return nil, fmt.Errorf("unsupported provider: %s", cacheEntry.BackendReserveResponse.Provider)

	}

	return &DockerGHAUploadCacheResponse{}, nil
}

func CommitCache(ctx context.Context, input DockerGHACommitCacheRequest) (*DockerGHACommitCacheResponse, error) {
	// Trigger upload to S3 now that we are sure all chunks have been received
	_, err := uploadToBlobStorage(ctx, input.CacheID)
	if err != nil {
		return nil, err
	}

	cacheEntryData, ok := cacheStore.Load(input.CacheID)
	if !ok {
		return nil, fmt.Errorf("cache ID not found")
	}

	cacheEntry := cacheEntryData.(*CacheEntryData)

	requestURL := fmt.Sprintf("%s/v1/cache/commit", input.HostURL)

	var payload CommitCacheRequest

	switch cacheEntry.BackendReserveResponse.Provider {
	case ProviderS3:
		payload = CommitCacheRequest{
			CacheKey:     cacheEntry.CacheKey,
			CacheVersion: cacheEntry.CacheVersion,
			UploadKey:    cacheEntry.BackendReserveResponse.S3.UploadKey,
			UploadID:     cacheEntry.BackendReserveResponse.S3.UploadID,
			Parts:        cacheEntry.S3Parts,
			VCSType:      "github",
		}

	case ProviderGCS:
		payload = CommitCacheRequest{
			CacheKey:     cacheEntry.CacheKey,
			CacheVersion: cacheEntry.CacheVersion,
			VCSType:      "github",
		}

	default:
		return nil, fmt.Errorf("unsupported provider: %s", cacheEntry.BackendReserveResponse.Provider)
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	serviceURL, err := url.Parse(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service URL: %w", err)
	}

	agent := fiber.Post(serviceURL.String())

	agent.Body(payloadBytes)

	agent.Add("Content-Type", "application/json")
	agent.Add("Accept", "application/json")
	agent.Add("Authorization", fmt.Sprintf("Bearer %s", input.AuthToken))

	statusCode, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to send request to cache backend: %v", errs)
	}

	if statusCode < 200 || statusCode >= 300 {
		return nil, fmt.Errorf("failed to commit cache: %s", string(body))
	}

	var commitCacheResponse CommitCacheResponse
	if err := json.Unmarshal(body, &commitCacheResponse); err != nil {
		return nil, fmt.Errorf("failed to parse backend response: %w", err)
	}

	cacheStore.Delete(input.CacheID)

	return &DockerGHACommitCacheResponse{}, nil
}
