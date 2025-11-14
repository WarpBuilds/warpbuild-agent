package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	cachepb "github.com/warpbuilds/warpbuild-agent/pkg/transparent-cache/derp/generated_go/results/api/v1"
	entitiespb "github.com/warpbuilds/warpbuild-agent/pkg/transparent-cache/derp/generated_go/results/entities/v1"
)

// TwirpClient implements a Twirp protocol client
type TwirpClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewTwirpClient creates a new Twirp client
func NewTwirpClient(baseURL string) *TwirpClient {
	return &TwirpClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// doTwirpRequest performs a Twirp RPC call
func (c *TwirpClient) doTwirpRequest(ctx context.Context, service, method string, request, response interface{}) error {
	// Serialize request to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Construct Twirp URL
	url := fmt.Sprintf("%s/twirp/%s/%s", c.baseURL, service, method)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set Twirp headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check for Twirp errors
	if resp.StatusCode != http.StatusOK {
		var twirpErr struct {
			Code string `json:"code"`
			Msg  string `json:"msg"`
		}
		if err := json.Unmarshal(respBody, &twirpErr); err == nil && twirpErr.Code != "" {
			return fmt.Errorf("twirp error %s: %s", twirpErr.Code, twirpErr.Msg)
		}
		return fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Deserialize response
	if response != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, response); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// GetCacheEntryDownloadURL calls the Twirp GetCacheEntryDownloadURL method
func (c *TwirpClient) GetCacheEntryDownloadURL(ctx context.Context, req *cachepb.GetCacheEntryDownloadURLRequest) (*cachepb.GetCacheEntryDownloadURLResponse, error) {
	var resp cachepb.GetCacheEntryDownloadURLResponse
	err := c.doTwirpRequest(ctx, "github.actions.results.api.v1.CacheService", "GetCacheEntryDownloadURL", req, &resp)
	return &resp, err
}

// CreateCacheEntry calls the Twirp CreateCacheEntry method
func (c *TwirpClient) CreateCacheEntry(ctx context.Context, req *cachepb.CreateCacheEntryRequest) (*cachepb.CreateCacheEntryResponse, error) {
	var resp cachepb.CreateCacheEntryResponse
	err := c.doTwirpRequest(ctx, "github.actions.results.api.v1.CacheService", "CreateCacheEntry", req, &resp)
	return &resp, err
}

// FinalizeCacheEntryUpload calls the Twirp FinalizeCacheEntryUpload method
func (c *TwirpClient) FinalizeCacheEntryUpload(ctx context.Context, req *cachepb.FinalizeCacheEntryUploadRequest) (*cachepb.FinalizeCacheEntryUploadResponse, error) {
	var resp cachepb.FinalizeCacheEntryUploadResponse
	err := c.doTwirpRequest(ctx, "github.actions.results.api.v1.CacheService", "FinalizeCacheEntryUpload", req, &resp)
	return &resp, err
}

func main() {
	// Create Twirp client
	client := NewTwirpClient("http://localhost:59991")

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Example: GitHub Actions v2 cache restore (similar to the code you showed)
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
		log.Fatalf("Failed to get cache entry: %v", err)
	}

	if !response.Ok {
		log.Printf("Cache not found for key: %s", request.Key)
		return
	}

	log.Printf("Cache hit! Matched key: %s", response.MatchedKey)
	log.Printf("Download URL: %s", response.SignedDownloadUrl)

	// In real GitHub Actions, they would:
	// 1. Download the cache from response.SignedDownloadUrl
	// 2. Extract the archive
	// 3. Restore the cached files

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
		log.Fatalf("Failed to create cache entry: %v", err)
	}

	if createResp.Ok {
		log.Printf("Cache entry created! Upload URL: %s", createResp.SignedUploadUrl)
		// In real GitHub Actions, they would upload the cache archive to this URL
	} else {
		log.Println("Failed to create cache entry (maybe it already exists)")
	}
}
