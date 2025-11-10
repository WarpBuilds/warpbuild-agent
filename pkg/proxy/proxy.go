package proxy

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

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
	primaryKey := input.Keys[0]
	// Docker backend weirdly sends impartial key as primary key sometimes.
	restoreKeys := input.Keys

	cacheResponse, err := callCacheBackend[GetCacheResponse](ctx, CacheBackendRequest{
		Path: "/get",
		Body: GetCacheRequest{
			CacheKey:     primaryKey,
			CacheVersion: input.Version,
			RestoreKeys:  restoreKeys,
		},
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

	reserveCacheResponse, err := callCacheBackend[ReserveCacheResponse](ctx, CacheBackendRequest{
		Path: "/reserve",
		Body: ReserveCacheRequest{
			CacheKey:       input.Key,
			CacheVersion:   input.Version,
			NumberOfChunks: 1,
			ContentType:    "application/octet-stream",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to reserve cache: %w", err)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomCacheID := r.Intn(1000000)

	dockerReserveResponse := DockerGHAReserveCacheResponse{
		CacheID: randomCacheID,
	}

	// Save this cache ID for later use
	cacheEntry := &CacheEntryData{
		BackendReserveResponse: *reserveCacheResponse,
		CacheKey:               input.Key,
		CacheVersion:           input.Version,
		Chunks:                 make(map[int64]ChunkData),
	}

	cacheStore.Store(randomCacheID, cacheEntry)

	fmt.Printf("cacheEntry saved with ID: %d: %+v\n", randomCacheID, cacheEntry)

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

	switch cacheEntry.BackendReserveResponse.Provider {
	case ProviderS3:
	case ProviderR2:
		if len(cacheEntry.BackendReserveResponse.S3.PreSignedURLs) != 1 {
			return nil, fmt.Errorf("no presigned URLs found")
		}

		s3PresignedURL := cacheEntry.BackendReserveResponse.S3.PreSignedURLs[0]

		for attempt := range maxRetries {
			contentReader, totalSize := NewUnorderedChunkReader(cacheEntry.Chunks)

			req, err := http.NewRequestWithContext(ctx, http.MethodPut, s3PresignedURL, contentReader)
			if err != nil {
				return nil, fmt.Errorf("failed to create S3 request: %w", err)
			}
			req.Header.Set("Content-Type", "application/octet-stream")
			req.ContentLength = totalSize

			resp, err := http.DefaultClient.Do(req)
			if err == nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent) {
				defer resp.Body.Close()
				etag := resp.Header.Get("ETag")
				if etag == "" {
					return nil, fmt.Errorf("no ETag found in response")
				}

				partNum := int32(1)
				cacheEntry.S3Parts = []S3CompletedPart{{
					ETag:       &etag,
					PartNumber: &partNum,
				}}

				// Success, break out of retry loop
				break
			}

			// If response is not OK, log and prepare for retry
			if resp != nil {
				defer resp.Body.Close()
				if attempt < maxRetries-1 {
					fmt.Printf("Retrying upload... attempt %d/%d, error: %v\n", attempt+1, maxRetries, err)
					time.Sleep((1 << attempt) * time.Second) // Exponential backoff
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

		// Upload context
		uploadCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		wc := object.NewWriter(uploadCtx)

		for attempt := range maxRetries {
			contentReader, _ := NewUnorderedChunkReader(cacheEntry.Chunks)
			_, err = io.Copy(wc, contentReader)
			if err == nil {
				err = wc.Close()
				if err != nil {
					return nil, fmt.Errorf("failed to close GCS writer: %w", err)
				}

				break
			}

			if attempt < maxRetries-1 {
				fmt.Printf("Retrying upload... attempt %d/%d, error: %v\n", attempt+1, maxRetries, err)
				time.Sleep((1 << attempt) * time.Second)
			}
		}

	case ProviderAzureBlob:
		if cacheEntry.BackendReserveResponse.AzureBlob.PreSignedURL == "" {
			return nil, fmt.Errorf("no presigned URL found")
		}

		for attempt := range maxRetries {
			contentReader, totalSize := NewUnorderedChunkReader(cacheEntry.Chunks)
			req, err := http.NewRequestWithContext(ctx, http.MethodPut, cacheEntry.BackendReserveResponse.AzureBlob.PreSignedURL, contentReader)
			if err != nil {
				return nil, fmt.Errorf("failed to create Azure Blob request: %w", err)
			}
			req.Header.Set("Content-Type", "application/octet-stream")
			req.Header.Set("x-ms-blob-type", "BlockBlob")
			req.ContentLength = totalSize

			resp, err := http.DefaultClient.Do(req)
			if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
				resp.Body.Close()
				break
			} else {
				defer resp.Body.Close()
				if attempt < maxRetries-1 {
					fmt.Printf("Retrying upload... attempt %d/%d, error: %v\n", attempt+1, maxRetries, err)
					time.Sleep((1 << attempt) * time.Second) // Exponential backoff
				}
			}
		}

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

	fmt.Printf("cacheEntry found for ID: %d: %+v\n", input.CacheID, cacheEntry)

	payload := CommitCacheRequest{
		CacheKey:     cacheEntry.CacheKey,
		CacheVersion: cacheEntry.CacheVersion,
		VCSType:      "github",
	}

	switch cacheEntry.BackendReserveResponse.Provider {
	case ProviderS3:
	case ProviderR2:
		payload.UploadKey = cacheEntry.BackendReserveResponse.S3.UploadKey
		payload.UploadID = cacheEntry.BackendReserveResponse.S3.UploadID
		payload.Parts = cacheEntry.S3Parts
	case ProviderGCS:
	case ProviderAzureBlob:
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cacheEntry.BackendReserveResponse.Provider)
	}

	_, err = callCacheBackend[CommitCacheResponse](ctx, CacheBackendRequest{
		Path: "/commit",
		Body: payload,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to commit cache: %w", err)
	}

	cacheStore.Delete(input.CacheID)

	return &DockerGHACommitCacheResponse{}, nil
}
