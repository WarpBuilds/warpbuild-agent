package main

import (
	"context"
	"log"
	"net/http"
	"time"

	cachepb "github.com/warpbuilds/warpbuild-agent/pkg/transparent-cache/derp/generated_go/results/api/v1"
	entitiespb "github.com/warpbuilds/warpbuild-agent/pkg/transparent-cache/derp/generated_go/results/entities/v1"
)

func main() {
	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create Twirp client using the generated client
	client := cachepb.NewCacheServiceJSONClient("http://localhost:50051", httpClient)

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Example: GitHub Actions v2 cache restore
	log.Println("=== Simulating GitHub Actions v2 cache restore ===")

	// Prepare request similar to GitHub's v2 cache
	request := &cachepb.GetCacheEntryDownloadURLRequest{
		Metadata: &entitiespb.CacheMetadata{
			RepositoryId: 12345,
			Scope: []*entitiespb.CacheScope{
				{
					Scope:      "refs/heads/main",
					Permission: 1, // Read permission
				},
			},
		},
		Key:         "v2-cache-key",
		RestoreKeys: []string{"v2-cache", "v2"},
		Version:     "abc123", // This would be computed from paths, compression method, etc.
	}

	// Call GetCacheEntryDownloadURL using Twirp
	response, err := client.GetCacheEntryDownloadURL(ctx, request)
	if err != nil {
		log.Printf("Failed to get cache entry: %v", err)
		// Try to create a new cache entry instead
	} else {
		if !response.Ok {
			log.Printf("Cache not found for key: %s", request.Key)
		} else {
			log.Printf("Cache hit! Matched key: %s", response.MatchedKey)
			log.Printf("Download URL: %s", response.SignedDownloadUrl)
		}
	}

	// Example 2: Create a new cache entry (for cache save operation)
	log.Println("\n=== Creating new v2 cache entry ===")

	createReq := &cachepb.CreateCacheEntryRequest{
		Metadata: &entitiespb.CacheMetadata{
			RepositoryId: 12345,
			Scope: []*entitiespb.CacheScope{
				{
					Scope:      "refs/heads/main",
					Permission: 3, // Read + Write
				},
			},
		},
		Key:     "v2-new-cache",
		Version: "xyz789",
	}

	createResp, err := client.CreateCacheEntry(ctx, createReq)
	if err != nil {
		log.Printf("Failed to create cache entry: %v", err)
	} else {
		if createResp.Ok {
			log.Printf("Cache entry created! Upload URL: %s", createResp.SignedUploadUrl)

			// Simulate upload completion
			log.Println("\n=== Finalizing cache upload ===")
			finalizeReq := &cachepb.FinalizeCacheEntryUploadRequest{
				Metadata:  createReq.Metadata,
				Key:       createReq.Key,
				Version:   createReq.Version,
				SizeBytes: 1024 * 1024, // 1MB
			}

			finalizeResp, err := client.FinalizeCacheEntryUpload(ctx, finalizeReq)
			if err != nil {
				log.Printf("Failed to finalize upload: %v", err)
			} else {
				if finalizeResp.Ok {
					log.Printf("Upload finalized! Entry ID: %d", finalizeResp.EntryId)
				} else {
					log.Println("Failed to finalize upload")
				}
			}
		} else {
			log.Println("Failed to create cache entry (maybe it already exists)")
		}
	}

	// Example 3: Test cache retrieval
	log.Println("\n=== Testing cache retrieval ===")
	getReq := &cachepb.GetCacheEntryDownloadURLRequest{
		Metadata: createReq.Metadata,
		Key:      "v2-new-cache",
		Version:  "xyz789",
	}

	getResp, err := client.GetCacheEntryDownloadURL(ctx, getReq)
	if err != nil {
		log.Printf("Failed to get cache: %v", err)
	} else {
		if getResp.Ok {
			log.Printf("Cache found! Download URL: %s", getResp.SignedDownloadUrl)
		} else {
			log.Println("Cache not found")
		}
	}
}
