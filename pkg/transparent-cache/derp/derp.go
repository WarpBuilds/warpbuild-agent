package derp

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	cachepb "github.com/warpbuilds/warpbuild-agent/pkg/transparent-cache/derp/generated_go/results/api/v1"
	entitiespb "github.com/warpbuilds/warpbuild-agent/pkg/transparent-cache/derp/generated_go/results/entities/v1"
)

// cacheServiceImpl implements the generated CacheService interface
type cacheServiceImpl struct {
	// In-memory storage for demonstration
	cacheEntries map[string]*cacheEntry
}

type cacheEntry struct {
	key       string
	version   string
	size      int64
	uploadURL string
	data      []byte
	metadata  *entitiespb.CacheMetadata
}

// NewCacheService creates a new cache service implementation
func NewCacheService() cachepb.CacheService {
	return &cacheServiceImpl{
		cacheEntries: make(map[string]*cacheEntry),
	}
}

// CreateCacheEntry implements CacheService.CreateCacheEntry
func (s *cacheServiceImpl) CreateCacheEntry(ctx context.Context, req *cachepb.CreateCacheEntryRequest) (*cachepb.CreateCacheEntryResponse, error) {
	log.Printf("CreateCacheEntry: key=%s, version=%s", req.Key, req.Version)

	// Validate request
	if req.Key == "" {
		return nil, fmt.Errorf("cache key is required")
	}

	// Check if cache entry already exists
	cacheKey := fmt.Sprintf("%s:%s", req.Key, req.Version)
	if _, exists := s.cacheEntries[cacheKey]; exists {
		return &cachepb.CreateCacheEntryResponse{
			Ok:              false,
			SignedUploadUrl: "",
		}, nil
	}

	// Generate a mock signed upload URL
	uploadURL := fmt.Sprintf("https://storage.example.com/upload/%s?token=%d", cacheKey, time.Now().Unix())

	// Store the pending cache entry
	s.cacheEntries[cacheKey] = &cacheEntry{
		key:       req.Key,
		version:   req.Version,
		uploadURL: uploadURL,
		metadata:  req.Metadata,
	}

	return &cachepb.CreateCacheEntryResponse{
		Ok:              true,
		SignedUploadUrl: uploadURL,
	}, nil
}

// FinalizeCacheEntryUpload implements CacheService.FinalizeCacheEntryUpload
func (s *cacheServiceImpl) FinalizeCacheEntryUpload(ctx context.Context, req *cachepb.FinalizeCacheEntryUploadRequest) (*cachepb.FinalizeCacheEntryUploadResponse, error) {
	log.Printf("FinalizeCacheEntryUpload: key=%s, size=%d bytes", req.Key, req.SizeBytes)

	// Validate request
	if req.Key == "" {
		return nil, fmt.Errorf("cache key is required")
	}

	// Find the cache entry
	cacheKey := fmt.Sprintf("%s:%s", req.Key, req.Version)
	entry, exists := s.cacheEntries[cacheKey]
	if !exists {
		return &cachepb.FinalizeCacheEntryUploadResponse{
			Ok:      false,
			EntryId: 0,
		}, nil
	}

	// Update the cache entry with size
	entry.size = req.SizeBytes

	// Generate a mock entry ID
	entryID := int64(time.Now().UnixNano())

	return &cachepb.FinalizeCacheEntryUploadResponse{
		Ok:      true,
		EntryId: entryID,
	}, nil
}

// GetCacheEntryDownloadURL implements CacheService.GetCacheEntryDownloadURL
func (s *cacheServiceImpl) GetCacheEntryDownloadURL(ctx context.Context, req *cachepb.GetCacheEntryDownloadURLRequest) (*cachepb.GetCacheEntryDownloadURLResponse, error) {
	log.Printf("GetCacheEntryDownloadURL: key=%s, restore_keys=%v", req.Key, req.RestoreKeys)

	// Validate request
	if req.Key == "" && len(req.RestoreKeys) == 0 {
		return nil, fmt.Errorf("cache key or restore keys are required")
	}

	// Get download URL from environment or use a placeholder
	downloadURL := os.Getenv("TEST_DOWNLOAD_URL")
	if downloadURL == "" {
		// Return a placeholder URL that indicates configuration is needed
		downloadURL = "https://your-endpoint.example.com/your-bucket/your-file?configure-TEST_DOWNLOAD_URL-env-var"
	}

	// Use the requested key as the matched key
	matchedKey := req.Key
	if matchedKey == "" && len(req.RestoreKeys) > 0 {
		matchedKey = req.RestoreKeys[0]
	}

	return &cachepb.GetCacheEntryDownloadURLResponse{
		Ok:                true,
		SignedDownloadUrl: downloadURL,
		MatchedKey:        matchedKey,
	}, nil
}

// Start starts the DERP cache Twirp service
func Start(port int) error {
	// Create cache service implementation
	service := NewCacheService()

	// Create Twirp server
	twirpHandler := cachepb.NewCacheServiceServer(service)

	// Create HTTP mux and add handlers
	mux := http.NewServeMux()

	// Mount the Twirp handler - it handles its own routing
	mux.Handle(cachepb.CacheServicePathPrefix, twirpHandler)

	// Add health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy"}`))
	})

	log.Printf("Cache Twirp server starting on port %d...", port)
	log.Printf("Twirp service available at: %s", cachepb.CacheServicePathPrefix)

	// Start HTTP server
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
