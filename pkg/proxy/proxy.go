package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

var cacheStore = sync.Map{}

type CacheEntryData struct {
	BackendReserveResponse ReserveCacheResponse
	S3Parts                []S3CompletedPart
	CacheKey               string
	CacheVersion           string
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
		ContentType:    "application/zstd",
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
	cacheStore.Store(randomCacheID, CacheEntryData{
		BackendReserveResponse: reserveCacheResponse,
		CacheKey:               input.Key,
		CacheVersion:           input.Version,
	})

	return &dockerReserveResponse, nil
}

type BufferData struct {
	Content []byte
	mu      sync.Mutex
}

var bufferStore = sync.Map{}

func UploadCache(ctx context.Context, input DockerGHAUploadCacheRequest) (*DockerGHAUploadCacheResponse, error) {
	bufferData, _ := bufferStore.LoadOrStore(input.CacheID, &BufferData{})
	buffer := bufferData.(*BufferData)

	// Lock the buffer for writing the incoming data
	buffer.mu.Lock()
	buffer.Content = append(buffer.Content, input.Content...)
	buffer.mu.Unlock()

	return &DockerGHAUploadCacheResponse{}, nil
}

func uploadToBlobStorage(ctx context.Context, cacheID int) (*DockerGHAUploadCacheResponse, error) {
	bufferData, ok := bufferStore.Load(cacheID)
	if !ok {
		return nil, fmt.Errorf("buffer data not found for cache ID %d", cacheID)
	}

	buffer, ok := bufferData.(*BufferData)
	if !ok {
		return nil, fmt.Errorf("buffer data is not of type BufferData")
	}

	cacheEntryData, ok := cacheStore.Load(cacheID)
	if !ok {
		return nil, fmt.Errorf("cache ID not found")
	}

	cacheEntry := cacheEntryData.(CacheEntryData)

	if cacheEntry.BackendReserveResponse.Provider == ProviderS3 {
		if len(cacheEntry.BackendReserveResponse.S3.PreSignedURLs) != 1 {
			return nil, fmt.Errorf("no presigned URLs found")
		}

		s3PresignedURL := cacheEntry.BackendReserveResponse.S3.PreSignedURLs[0]

		contentReader := bytes.NewReader(buffer.Content)
		req, err := http.NewRequestWithContext(ctx, http.MethodPut, s3PresignedURL, contentReader)
		if err != nil {
			return nil, fmt.Errorf("failed to create S3 request: %w", err)
		}
		req.Header.Set("Content-Type", "application/octet-stream")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to upload to S3: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("S3 upload failed with status %d: %s", resp.StatusCode, string(bodyBytes))
		}

		etag := resp.Header.Get("ETag")
		if etag == "" {
			return nil, fmt.Errorf("no ETag found in response")
		}

		partNumberPtr := int32(1)
		cacheEntry.S3Parts = []S3CompletedPart{{
			ETag:       &etag,
			PartNumber: &partNumberPtr,
		}}

		cacheStore.Store(cacheID, cacheEntry)

	} else {
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

	cacheEntry := cacheEntryData.(CacheEntryData)

	defer bufferStore.Delete(input.CacheID)

	if cacheEntry.BackendReserveResponse.Provider == ProviderS3 {
		requestURL := fmt.Sprintf("%s/v1/cache/commit", input.HostURL)

		payload := CommitCacheRequest{
			CacheKey:     cacheEntry.CacheKey,
			CacheVersion: cacheEntry.CacheVersion,
			UploadKey:    cacheEntry.BackendReserveResponse.S3.UploadKey,
			UploadID:     cacheEntry.BackendReserveResponse.S3.UploadID,
			Parts:        cacheEntry.S3Parts,
			VCSType:      "github",
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

	} else {
		return nil, fmt.Errorf("unsupported provider: %s", cacheEntry.BackendReserveResponse.Provider)
	}

	return &DockerGHACommitCacheResponse{}, nil
}
