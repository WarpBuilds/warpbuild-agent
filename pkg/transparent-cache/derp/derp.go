package derp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	cachepb "github.com/warpbuilds/warpbuild-agent/pkg/transparent-cache/derp/generated_go/results/api/v1"
)

// Global stores for caching backend responses and credentials
var (
	cacheStore       = sync.Map{} // Stores cache entries by key
	credentialsStore = sync.Map{} // Stores credentials for different providers
)

// CacheEntryInfo stores the backend response and metadata
type CacheEntryInfo struct {
	BackendResponse interface{} // Can be GetCacheResponse, ReserveCacheResponse, etc.
	Provider        Provider
	CacheKey        string
	CacheVersion    string
	CreatedAt       time.Time
}

// cacheServiceImpl implements the generated CacheService interface
type cacheServiceImpl struct {
	backendURL string
	authToken  string
}

func NewCacheService(backendURL, authToken string) cachepb.CacheService {
	if backendURL == "" {
		backendURL = "https://api.warpbuild.com" // Default backend URL
	}
	return &cacheServiceImpl{
		backendURL: backendURL,
		authToken:  authToken,
	}
}

func (s *cacheServiceImpl) CreateCacheEntry(ctx context.Context, req *cachepb.CreateCacheEntryRequest) (*cachepb.CreateCacheEntryResponse, error) {
	log.Printf("CreateCacheEntry: key=%s, version=%s", req.Key, req.Version)

	if req.Key == "" {
		return nil, fmt.Errorf("cache key is required")
	}

	// Call backend to reserve cache
	reserveReq := ReserveCacheRequest{
		CacheKey:       req.Key,
		CacheVersion:   req.Version,
		NumberOfChunks: 1,
		ContentType:    "application/octet-stream",
	}

	reserveResp, err := s.callBackendReserve(ctx, reserveReq)
	if err != nil {
		return nil, fmt.Errorf("failed to reserve cache: %w", err)
	}

	// Validate that S3/R2 responses have AccessGrant
	if (reserveResp.Provider == ProviderS3 || reserveResp.Provider == ProviderR2) &&
		(reserveResp.S3 == nil || reserveResp.S3.AccessGrant == nil) {
		return nil, fmt.Errorf("S3/R2 response missing AccessGrant for cache key %s", req.Key)
	}

	// Store the response in memory
	cacheKey := fmt.Sprintf("%s:%s", req.Key, req.Version)
	cacheInfo := &CacheEntryInfo{
		BackendResponse: reserveResp,
		Provider:        reserveResp.Provider,
		CacheKey:        req.Key,
		CacheVersion:    req.Version,
		CreatedAt:       time.Now(),
	}
	cacheStore.Store(cacheKey, cacheInfo)

	// Store credentials if available
	s.storeCredentials(reserveResp, req.Key, req.Version)

	// Generate appropriate upload URL based on provider
	uploadURL := s.generateUploadURL(reserveResp, req.Key, req.Version)

	return &cachepb.CreateCacheEntryResponse{
		Ok:              true,
		SignedUploadUrl: uploadURL,
	}, nil
}

func (s *cacheServiceImpl) FinalizeCacheEntryUpload(ctx context.Context, req *cachepb.FinalizeCacheEntryUploadRequest) (*cachepb.FinalizeCacheEntryUploadResponse, error) {
	log.Printf("FinalizeCacheEntryUpload: key=%s, version=%s, size=%d bytes", req.Key, req.Version, req.SizeBytes)

	if req.Key == "" {
		return nil, fmt.Errorf("cache key is required")
	}

	// Call backend to commit cache
	commitReq := CommitCacheRequest{
		CacheKey:     req.Key,
		CacheVersion: req.Version,
		Parts:        []S3CompletedPart{}, // Will be populated for S3/R2
		VCSType:      "github",
	}

	// Get the cached info to determine provider
	cacheKey := fmt.Sprintf("%s:%s", req.Key, req.Version)
	if cacheInfo, ok := cacheStore.Load(cacheKey); ok {
		info := cacheInfo.(*CacheEntryInfo)
		commitReq.Provider = info.Provider
	}

	// For S3/R2, include UploadKey, UploadID and completed parts captured by asur
	if commitReq.Provider == ProviderS3 || commitReq.Provider == ProviderR2 {
		cacheIdentifier := fmt.Sprintf("%s--%s", req.Key, req.Version)

		if uploadKey, ok := GetS3UploadKey(cacheIdentifier); ok {
			commitReq.UploadKey = uploadKey
		}
		if uploadID, ok := GetS3UploadID(cacheIdentifier); ok {
			commitReq.UploadID = uploadID
		}
		if parts, ok := GetS3CompletedParts(cacheIdentifier); ok {
			commitReq.Parts = parts
		} else {
			log.Printf("FinalizeCacheEntryUpload: no parts found for cacheIdentifier=%s", cacheIdentifier)

			// Try alternative cache identifiers if no parts found
			alternativeIdentifiers := []string{}

			// If version is empty, try with just key--
			if req.Version == "" {
				alternativeIdentifiers = append(alternativeIdentifiers, fmt.Sprintf("%s--", req.Key))
			}

			// Try with some common default versions
			commonDefaultVersions := []string{"0", "1", "default"}
			for _, defaultVer := range commonDefaultVersions {
				if req.Version != defaultVer {
					alternativeIdentifiers = append(alternativeIdentifiers, fmt.Sprintf("%s--%s", req.Key, defaultVer))
				}
			}

			// Try each alternative
			for _, altId := range alternativeIdentifiers {
				if parts, ok := GetS3CompletedParts(altId); ok {
					commitReq.Parts = parts
					break
				}
			}

			// If still no parts found, try finding by cache key prefix
			if len(commitReq.Parts) == 0 {
				if parts, foundId, ok := findPartsForCacheKey(req.Key); ok {
					commitReq.Parts = parts

					// Also try to get the upload key and ID with this identifier
					if uploadKey, ok := GetS3UploadKey(foundId); ok {
						commitReq.UploadKey = uploadKey
					}
					if uploadID, ok := GetS3UploadID(foundId); ok {
						commitReq.UploadID = uploadID
					}
				}
			}
		}
	}

	commitResp, err := s.callBackendCommit(ctx, commitReq)
	if err != nil {
		return nil, fmt.Errorf("failed to commit cache: %w", err)
	}

	entryID := int64(0)
	if commitResp.CacheEntry != nil {
		// Use a hash of the cache entry ID as the entry ID
		entryID = int64(hash(commitResp.CacheEntry.ID))
	}

	return &cachepb.FinalizeCacheEntryUploadResponse{
		Ok:      true,
		EntryId: entryID,
	}, nil
}

func (s *cacheServiceImpl) GetCacheEntryDownloadURL(ctx context.Context, req *cachepb.GetCacheEntryDownloadURLRequest) (*cachepb.GetCacheEntryDownloadURLResponse, error) {
	log.Printf("GetCacheEntryDownloadURL: key=%s, restore_keys=%v", req.Key, req.RestoreKeys)

	if req.Key == "" && len(req.RestoreKeys) == 0 {
		return nil, fmt.Errorf("cache key or restore keys are required")
	}

	// Call backend to get cache
	getReq := GetCacheRequest{
		CacheKey:     req.Key,
		CacheVersion: req.Version,
		RestoreKeys:  req.RestoreKeys,
	}

	getResp, err := s.callBackendGet(ctx, getReq)
	if err != nil {
		return &cachepb.GetCacheEntryDownloadURLResponse{
			Ok: false,
		}, nil
	}

	if getResp.CacheEntry == nil {
		return &cachepb.GetCacheEntryDownloadURLResponse{
			Ok: false,
		}, nil
	}

	// Validate that S3/R2 responses have AccessGrant
	if (getResp.Provider == ProviderS3 || getResp.Provider == ProviderR2) &&
		(getResp.S3 == nil || getResp.S3.AccessGrant == nil) {
		log.Printf("WARNING: S3/R2 response missing AccessGrant for cache key %s", req.Key)
		return &cachepb.GetCacheEntryDownloadURLResponse{
			Ok: false,
		}, nil
	}

	// Store the response in memory
	cacheKey := fmt.Sprintf("%s:%s", getResp.CacheEntry.CacheUserGivenKey, getResp.CacheEntry.CacheVersion)
	cacheInfo := &CacheEntryInfo{
		BackendResponse: getResp,
		Provider:        getResp.Provider,
		CacheKey:        getResp.CacheEntry.CacheUserGivenKey,
		CacheVersion:    getResp.CacheEntry.CacheVersion,
		CreatedAt:       time.Now(),
	}
	cacheStore.Store(cacheKey, cacheInfo)

	s.storeCredentials(getResp, getResp.CacheEntry.CacheUserGivenKey, getResp.CacheEntry.CacheVersion)

	// Generate appropriate download URL based on provider
	downloadURL := s.generateDownloadURL(getResp)

	return &cachepb.GetCacheEntryDownloadURLResponse{
		Ok:                true,
		SignedDownloadUrl: downloadURL,
		MatchedKey:        getResp.CacheEntry.CacheUserGivenKey,
	}, nil
}

// Backend communication methods
func (s *cacheServiceImpl) callBackendGet(ctx context.Context, req GetCacheRequest) (*GetCacheResponse, error) {
	requestURL := fmt.Sprintf("%s/v1/cache/get", s.backendURL)

	payloadBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	serviceURL, err := url.Parse(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service URL: %w", err)
	}

	agent := fiber.Post(serviceURL.String())
	agent.Body(payloadBytes)
	agent.Add("Content-Type", "application/json")
	agent.Add("Accept", "application/json")
	if s.authToken != "" {
		agent.Add("Authorization", fmt.Sprintf("Bearer %s", s.authToken))
	}

	statusCode, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to send request: %v", errs)
	}

	if statusCode < 200 || statusCode >= 300 {
		return nil, fmt.Errorf("backend returned error: %s", string(body))
	}

	var resp GetCacheResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

func (s *cacheServiceImpl) callBackendReserve(ctx context.Context, req ReserveCacheRequest) (*ReserveCacheResponse, error) {
	requestURL := fmt.Sprintf("%s/v1/cache/reserve", s.backendURL)

	payloadBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	serviceURL, err := url.Parse(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service URL: %w", err)
	}

	agent := fiber.Post(serviceURL.String())
	agent.Body(payloadBytes)
	agent.Add("Content-Type", "application/json")
	agent.Add("Accept", "application/json")
	if s.authToken != "" {
		agent.Add("Authorization", fmt.Sprintf("Bearer %s", s.authToken))
	}

	statusCode, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to send request: %v", errs)
	}

	if statusCode < 200 || statusCode >= 300 {
		return nil, fmt.Errorf("backend returned error: %s", string(body))
	}

	var resp ReserveCacheResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

func (s *cacheServiceImpl) callBackendCommit(ctx context.Context, req CommitCacheRequest) (*CommitCacheResponse, error) {
	requestURL := fmt.Sprintf("%s/v1/cache/commit", s.backendURL)

	payloadBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	serviceURL, err := url.Parse(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service URL: %w", err)
	}

	agent := fiber.Post(serviceURL.String())
	agent.Body(payloadBytes)
	agent.Add("Content-Type", "application/json")
	agent.Add("Accept", "application/json")
	if s.authToken != "" {
		agent.Add("Authorization", fmt.Sprintf("Bearer %s", s.authToken))
	}

	statusCode, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to send request: %v", errs)
	}

	if statusCode < 200 || statusCode >= 300 {
		return nil, fmt.Errorf("backend returned error: %s", string(body))
	}

	var resp CommitCacheResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// Credential storage methods
func (s *cacheServiceImpl) storeCredentials(response interface{}, cacheKey, cacheVersion string) {
	switch resp := response.(type) {
	case *GetCacheResponse:
		switch resp.Provider {
		case ProviderGCS:
			if resp.GCS != nil && resp.GCS.ShortLivedToken != nil {
				credentialsStore.Store("gcs_token", resp.GCS.ShortLivedToken)
			}
		case ProviderS3, ProviderR2:
			if resp.S3 != nil {
				credentialsStore.Store("s3_access_grant", resp.S3.AccessGrant)
				// Store the CacheKey for S3/R2 downloads
				if resp.S3.CacheKey != "" {
					// Store with cache identifier for lookup
					cacheIdentifier := fmt.Sprintf("%s--%s", cacheKey, cacheVersion)
					credentialsStore.Store(fmt.Sprintf("s3_cache_key_%s", cacheIdentifier), resp.S3.CacheKey)
				}
			}
		case ProviderAzureBlob:
			if resp.AzureBlob != nil {
				credentialsStore.Store("azure_url", resp.AzureBlob.PreSignedURL)
				credentialsStore.Store("azure_bucket", resp.AzureBlob.BucketName)
			}
		}
	case *ReserveCacheResponse:
		switch resp.Provider {
		case ProviderGCS:
			if resp.GCS != nil && resp.GCS.ShortLivedToken != nil {
				credentialsStore.Store("gcs_token", resp.GCS.ShortLivedToken)
			}
		case ProviderS3, ProviderR2:
			if resp.S3 != nil && resp.S3.AccessGrant != nil {
				credentialsStore.Store("s3_access_grant", resp.S3.AccessGrant)
				// Store the UploadKey and UploadID for S3/R2 uploads
				if resp.S3.UploadKey != "" || resp.S3.UploadID != "" {
					cacheIdentifier := fmt.Sprintf("%s--%s", cacheKey, cacheVersion)
					if resp.S3.UploadKey != "" {
						credentialsStore.Store(fmt.Sprintf("s3_upload_key_%s", cacheIdentifier), resp.S3.UploadKey)
					}
					if resp.S3.UploadID != "" {
						credentialsStore.Store(fmt.Sprintf("s3_upload_id_%s", cacheIdentifier), resp.S3.UploadID)
					}
				}
			}
		case ProviderAzureBlob:
			if resp.AzureBlob != nil {
				credentialsStore.Store("azure_url", resp.AzureBlob.PreSignedURL)
				credentialsStore.Store("azure_container", resp.AzureBlob.ContainerName)
				credentialsStore.Store("azure_blob", resp.AzureBlob.BlobName)
			}
		}
	}
}

// URL generation methods
func (s *cacheServiceImpl) generateUploadURL(resp *ReserveCacheResponse, key, version string) string {
	// Always generate Azure-style URLs for asur to intercept
	// Encode the cache key and version in the URL so asur can retrieve credentials

	// Create a unique cache identifier that asur can use for lookup
	cacheIdentifier := fmt.Sprintf("%s--%s", key, version)

	// Store provider-specific information for asur to use later
	// The container encodes the provider type, and blob encodes the cache key/version
	var container string

	switch resp.Provider {
	case ProviderAzureBlob:
		container = "azure"
		// Store additional Azure-specific info if needed
		if resp.AzureBlob != nil {
			// Could store container/blob mapping if different from default
			credentialsStore.Store(fmt.Sprintf("azure_container_%s", cacheIdentifier), resp.AzureBlob.ContainerName)
			credentialsStore.Store(fmt.Sprintf("azure_blob_%s", cacheIdentifier), resp.AzureBlob.BlobName)
		}
	case ProviderS3, ProviderR2:
		if resp.Provider == ProviderS3 {
			container = "s3"
		} else {
			container = "r2"
		}
	case ProviderGCS:
		container = "gcs"
		if resp.GCS != nil && resp.GCS.BucketName != "" {
			// Store the actual bucket name for GCS
			credentialsStore.Store(fmt.Sprintf("gcs_bucket_%s", cacheIdentifier), resp.GCS.BucketName)
		}
	default:
		container = "cache"
	}

	// Always return Azure-style URL with encoded cache identifier
	// Format: https://warpbuild.blob.core.windows.net/{provider}/{cacheKey}--{version}
	return fmt.Sprintf("https://warpbuild.blob.core.windows.net/%s/%s", container, cacheIdentifier)
}

func (s *cacheServiceImpl) generateDownloadURL(resp *GetCacheResponse) string {
	// Always generate Azure-style URLs for asur to intercept
	if resp.CacheEntry == nil {
		return ""
	}

	// Create a unique cache identifier that asur can use for lookup
	cacheIdentifier := fmt.Sprintf("%s--%s", resp.CacheEntry.CacheUserGivenKey, resp.CacheEntry.CacheVersion)

	// Store provider-specific information for asur to use later
	var container string

	switch resp.Provider {
	case ProviderAzureBlob:
		container = "azure"
		if resp.AzureBlob != nil && resp.AzureBlob.BucketName != "" {
			credentialsStore.Store(fmt.Sprintf("azure_container_%s", cacheIdentifier), resp.AzureBlob.BucketName)
		}
	case ProviderS3, ProviderR2:
		if resp.Provider == ProviderS3 {
			container = "s3"
		} else {
			container = "r2"
		}
	case ProviderGCS:
		container = "gcs"
		if resp.GCS != nil && resp.GCS.BucketName != "" {
			credentialsStore.Store(fmt.Sprintf("gcs_bucket_%s", cacheIdentifier), resp.GCS.BucketName)
		}
	default:
		container = "cache"
	}

	// Always return Azure-style URL with encoded cache identifier
	return fmt.Sprintf("https://warpbuild.blob.core.windows.net/%s/%s", container, cacheIdentifier)
}

// Helper function to hash string to int64
func hash(s string) uint32 {
	h := uint32(0)
	for _, c := range s {
		h = h*31 + uint32(c)
	}
	return h
}

// GetStoredCredentials returns stored credentials for a given provider
func GetStoredCredentials(provider string) interface{} {
	if val, ok := credentialsStore.Load(provider); ok {
		return val
	}
	return nil
}

// GetCacheInfo returns stored cache information for a given cache key
func GetCacheInfo(cacheKey string) *CacheEntryInfo {
	if val, ok := cacheStore.Load(cacheKey); ok {
		if info, ok := val.(*CacheEntryInfo); ok {
			return info
		}
	}
	return nil
}

func Start(port int) error {
	// Get backend URL and auth token from environment
	backendURL := os.Getenv("WARPBUILD_CACHE_URL")
	if backendURL == "" {
		backendURL = "https://cache.warpbuild.com"
	}
	authToken := os.Getenv("WARPBUILD_RUNNER_VERIFICATION_TOKEN")

	service := NewCacheService(backendURL, authToken)

	twirpHandler := cachepb.NewCacheServiceServer(service)

	mux := http.NewServeMux()

	mux.Handle(cachepb.CacheServicePathPrefix, twirpHandler)

	log.Printf("Cache Twirp server starting on port %d...", port)
	log.Printf("Twirp service available at: %s", cachepb.CacheServicePathPrefix)
	log.Printf("Backend URL: %s", backendURL)

	// Start cleanup routine to remove old cache entries every 30 minutes
	// Remove entries older than 2 hours
	StartCleanupRoutine(30*time.Minute, 2*time.Hour)
	log.Printf("Started cleanup routine for cache entries")

	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

// GetAzureCredentials returns Azure-specific credentials if available
func GetAzureCredentials() (presignedURL, containerName, blobName string, found bool) {
	urlVal, urlOk := credentialsStore.Load("azure_url")
	containerVal, containerOk := credentialsStore.Load("azure_container")
	blobVal, blobOk := credentialsStore.Load("azure_blob")

	if urlOk {
		if url, ok := urlVal.(string); ok {
			presignedURL = url
			found = true
		}
	}
	if containerOk {
		if container, ok := containerVal.(string); ok {
			containerName = container
		}
	}
	if blobOk {
		if blob, ok := blobVal.(string); ok {
			blobName = blob
		}
	}

	return
}

// GetS3Credentials returns S3-specific credentials if available
func GetS3Credentials() (accessGrant *AccessGrant, found bool) {
	if val, ok := credentialsStore.Load("s3_access_grant"); ok {
		if grant, ok := val.(*AccessGrant); ok {
			accessGrant = grant
			found = true
		}
		return
	}

	return
}

// GetS3UploadKey returns the S3 upload key for a given cache identifier
func GetS3UploadKey(cacheIdentifier string) (uploadKey string, found bool) {
	key := fmt.Sprintf("s3_upload_key_%s", cacheIdentifier)
	if val, ok := credentialsStore.Load(key); ok {
		if key, ok := val.(string); ok {
			return key, true
		}
	}
	return "", false
}

// GetS3CacheKey returns the S3 cache key for a given cache identifier
func GetS3CacheKey(cacheIdentifier string) (cacheKey string, found bool) {
	key := fmt.Sprintf("s3_cache_key_%s", cacheIdentifier)
	if val, ok := credentialsStore.Load(key); ok {
		if key, ok := val.(string); ok {
			return key, true
		}
	}
	return "", false
}

// GetS3UploadID returns the S3 upload ID for a given cache identifier
func GetS3UploadID(cacheIdentifier string) (uploadID string, found bool) {
	key := fmt.Sprintf("s3_upload_id_%s", cacheIdentifier)
	if val, ok := credentialsStore.Load(key); ok {
		if id, ok := val.(string); ok && id != "" {
			return id, true
		}
	}
	return "", false
}

// internal struct for uploaded part info
type uploadedPartInfo struct {
	PartNumber int32
	ETag       string
}

// AddS3UploadedBlock records an uploaded block (by Azure block ID) for a cache identifier
func AddS3UploadedBlock(cacheIdentifier, blockID string, partNumber int32, etag string) {
	key := fmt.Sprintf("s3_parts_map_%s", cacheIdentifier)
	val, _ := credentialsStore.LoadOrStore(key, &sync.Map{})
	partsMap := val.(*sync.Map)
	partsMap.Store(blockID, uploadedPartInfo{PartNumber: partNumber, ETag: etag})
}

// SetS3BlockOrder records the order of Azure block IDs used when committing
func SetS3BlockOrder(cacheIdentifier string, order []string) {
	key := fmt.Sprintf("s3_block_order_%s", cacheIdentifier)
	credentialsStore.Store(key, order)
}

// GetS3CompletedParts returns the completed parts in the correct order for commit
func GetS3CompletedParts(cacheIdentifier string) ([]S3CompletedPart, bool) {
	// Load parts map
	partsKey := fmt.Sprintf("s3_parts_map_%s", cacheIdentifier)
	val, ok := credentialsStore.Load(partsKey)
	if !ok {
		return nil, false
	}
	partsMap, ok := val.(*sync.Map)
	if !ok {
		return nil, false
	}

	// Try to get block order
	orderKey := fmt.Sprintf("s3_block_order_%s", cacheIdentifier)
	var orderedIDs []string
	if v, ok := credentialsStore.Load(orderKey); ok {
		if ids, ok := v.([]string); ok {
			orderedIDs = ids
		}
	}

	buildPart := func(info uploadedPartInfo) S3CompletedPart {
		pn := info.PartNumber
		etag := info.ETag
		return S3CompletedPart{PartNumber: &pn, ETag: &etag}
	}

	result := make([]S3CompletedPart, 0)
	if len(orderedIDs) > 0 {
		for _, bid := range orderedIDs {
			if v, ok := partsMap.Load(bid); ok {
				info := v.(uploadedPartInfo)
				result = append(result, buildPart(info))
			}
		}
		if len(result) == 0 {
			return nil, false
		}
		return result, true
	}

	// Fallback: collect all and sort by PartNumber
	tmp := make([]uploadedPartInfo, 0)
	partsMap.Range(func(_, v interface{}) bool {
		tmp = append(tmp, v.(uploadedPartInfo))
		return true
	})
	if len(tmp) == 0 {
		return nil, false
	}
	sort.Slice(tmp, func(i, j int) bool { return tmp[i].PartNumber < tmp[j].PartNumber })
	for _, info := range tmp {
		result = append(result, buildPart(info))
	}
	return result, true
}

// GetGCSCredentials returns GCS-specific credentials if available
func GetGCSCredentials() (*ShortLivedToken, bool) {
	if val, ok := credentialsStore.Load("gcs_token"); ok {
		if token, ok := val.(*ShortLivedToken); ok {
			return token, true
		}
	}
	return nil, false
}

// GetCredentialsFromURL extracts cache identifier from URL and returns appropriate credentials
// This is the main function asur will use to get credentials from just the URL
func GetCredentialsFromURL(urlStr string) (provider Provider, credentials interface{}, found bool) {
	// Parse the URL to extract container and blob
	container, blob, err := ParseAzureURL(urlStr)
	if err != nil {
		return "", nil, false
	}

	// Determine provider from container name
	switch container {
	case "azure":
		provider = ProviderAzureBlob
	case "s3":
		provider = ProviderS3
	case "r2":
		provider = ProviderR2
	case "gcs":
		provider = ProviderGCS
	default:
		provider = ProviderAzureBlob // Default fallback
	}

	// The blob contains the cache identifier (key--version)
	cacheIdentifier := blob

	// Try to get cache info using the identifier
	// Convert identifier back to cache key format for lookup
	parts := strings.Split(cacheIdentifier, "--")
	if len(parts) >= 2 {
		cacheKey := fmt.Sprintf("%s:%s", parts[0], parts[1])
		if cacheInfo := GetCacheInfo(cacheKey); cacheInfo != nil {
			// Extract credentials based on provider
			switch provider {
			case ProviderAzureBlob:
				if resp, ok := cacheInfo.BackendResponse.(*GetCacheResponse); ok && resp.AzureBlob != nil {
					return provider, resp.AzureBlob, true
				} else if resp, ok := cacheInfo.BackendResponse.(*ReserveCacheResponse); ok && resp.AzureBlob != nil {
					return provider, resp.AzureBlob, true
				}
			case ProviderS3, ProviderR2:
				if resp, ok := cacheInfo.BackendResponse.(*GetCacheResponse); ok && resp.S3 != nil {
					return provider, resp.S3, true
				} else if resp, ok := cacheInfo.BackendResponse.(*ReserveCacheResponse); ok && resp.S3 != nil {
					return provider, resp.S3, true
				}
			case ProviderGCS:
				if resp, ok := cacheInfo.BackendResponse.(*GetCacheResponse); ok && resp.GCS != nil {
					return provider, resp.GCS, true
				} else if resp, ok := cacheInfo.BackendResponse.(*ReserveCacheResponse); ok && resp.GCS != nil {
					return provider, resp.GCS, true
				}
			}
		}
	}

	// Fallback: try to get credentials from the general stores
	switch provider {
	case ProviderAzureBlob:
		if url, _, _, ok := GetAzureCredentials(); ok {
			return provider, url, true
		}
	case ProviderS3, ProviderR2:
		if accessGrant, ok := GetS3Credentials(); ok {
			return provider, accessGrant, true
		}
	case ProviderGCS:
		if token, ok := GetGCSCredentials(); ok {
			return provider, token, true
		}
	}

	return provider, nil, false
}

// ParseAzureURL extracts container and blob name from a warpbuild.blob.core.windows.net URL
func ParseAzureURL(url string) (container, blob string, err error) {
	// Expected format: https://warpbuild.blob.core.windows.net/container/blob
	prefix := "https://warpbuild.blob.core.windows.net/"
	if !strings.HasPrefix(url, prefix) {
		// Also try without https://
		prefix = "warpbuild.blob.core.windows.net/"
		if !strings.HasPrefix(url, prefix) {
			return "", "", fmt.Errorf("invalid Azure URL format")
		}
	}

	path := strings.TrimPrefix(url, prefix)
	parts := strings.SplitN(path, "/", 2)

	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid Azure URL path")
	}

	return parts[0], parts[1], nil
}

// cleanupOldCacheEntries removes cache entries older than the specified duration
func cleanupOldCacheEntries(maxAge time.Duration) {
	now := time.Now()
	keysToDelete := []string{}

	// Find old entries
	cacheStore.Range(func(key, value interface{}) bool {
		if cacheKey, ok := key.(string); ok {
			if info, ok := value.(*CacheEntryInfo); ok {
				if now.Sub(info.CreatedAt) > maxAge {
					keysToDelete = append(keysToDelete, cacheKey)
				}
			}
		}
		return true
	})

	// Delete old entries
	for _, key := range keysToDelete {
		cacheStore.Delete(key)
	}

	if len(keysToDelete) > 0 {
		log.Printf("Cleaned up %d old cache entries", len(keysToDelete))
	}
}

// StartCleanupRoutine starts a background goroutine to periodically clean up old cache entries
func StartCleanupRoutine(cleanupInterval, maxAge time.Duration) {
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()
		for range ticker.C {
			cleanupOldCacheEntries(maxAge)
		}
	}()
}

// findPartsForCacheKey tries to find parts for a cache key by checking various possible cache identifiers
func findPartsForCacheKey(cacheKey string) ([]S3CompletedPart, string, bool) {
	var foundParts []S3CompletedPart
	var foundIdentifier string

	// Look for any parts map that starts with the cache key
	credentialsStore.Range(func(key, value interface{}) bool {
		if keyStr, ok := key.(string); ok {
			if strings.HasPrefix(keyStr, "s3_parts_map_") {
				// Extract cache identifier from the key
				cacheId := strings.TrimPrefix(keyStr, "s3_parts_map_")

				// Check if this cache identifier starts with our cache key
				if strings.HasPrefix(cacheId, cacheKey+"--") {
					if parts, ok := GetS3CompletedParts(cacheId); ok {
						foundParts = parts
						foundIdentifier = cacheId
						return false // Stop iteration
					}
				}
			}
		}
		return true
	})

	return foundParts, foundIdentifier, len(foundParts) > 0
}
