# Derp - Cache Service with Backend Integration

Derp is a Twirp-based cache service that connects with the WarpCache backend and provides credential management for transparent cache operations.

## Overview

Derp serves as an intermediary between clients and the WarpCache backend, handling:

- Communication with the WarpCache backend API
- Storage of credentials and cache metadata in memory
- Generation of proxy-friendly URLs for all storage providers using Azure-style URLs

## Key Features

1. **Backend Integration**: Communicates with WarpCache backend using the DTOs defined in `sdk.go`
2. **Credential Storage**: Stores credentials received from the backend in memory for use by other services
3. **Unified URL Format**: Generates Azure-style URLs (`warpbuild.blob.core.windows.net/container/...`) for ALL storage providers to ensure all requests go through asur
4. **URL-based Credential Retrieval**: Asur can retrieve all necessary credentials using only the URL

## Configuration

The service is configured via environment variables:

- `WARPCACHE_BACKEND_URL`: URL of the WarpCache backend (default: `https://api.warpbuild.com`)
- `WARPCACHE_AUTH_TOKEN`: Authentication token for the backend

## Architecture

### Request Flow

1. **Client → Derp**: Client requests cache operations via Twirp API
2. **Derp → Backend**: Derp forwards requests to WarpCache backend
3. **Backend → Derp**: Backend returns credentials and storage URLs
4. **Derp Storage**: Derp stores credentials in memory
5. **Derp → Client**: Returns Azure-style proxy URLs to client (regardless of actual storage provider)
6. **Client → Asur**: Client uses returned URLs for actual storage operations
7. **Asur → Storage**: Asur intercepts requests, retrieves credentials from URL, and proxies to the correct storage backend

### Integration with Asur

Asur (Azure proxy) intercepts ALL storage requests and can retrieve credentials using just the URL:

```go
// Main function to get credentials from URL
provider, credentials, found := derp.GetCredentialsFromURL(url)

// The function returns:
// - provider: The storage provider type (azure, s3, r2, gcs)
// - credentials: Provider-specific credential object
// - found: Whether credentials were found

// Example usage in asur:
fullURL := "https://warpbuild.blob.core.windows.net/s3/mycache--v1"
provider, creds, found := derp.GetCredentialsFromURL(fullURL)
if found {
    switch provider {
    case derp.ProviderS3:
        s3Creds := creds.(*derp.S3GetCacheResponse)
        // Use s3Creds.PreSignedURL
    }
}
```

### URL Format

ALL storage providers return Azure-style URLs with embedded information:

```
https://warpbuild.blob.core.windows.net/{provider}/{cacheKey}--{version}
```

Components:

- **provider**: Identifies the storage backend (`azure`, `s3`, `r2`, `gcs`)
- **cacheKey--version**: Unique identifier for credential lookup (uses `--` as separator)

Examples:

- Azure: `https://warpbuild.blob.core.windows.net/azure/build-cache--v1`
- S3: `https://warpbuild.blob.core.windows.net/s3/artifacts--abc123`
- GCS: `https://warpbuild.blob.core.windows.net/gcs/models--latest`

## API Endpoints

The service implements the following Twirp endpoints:

- `CreateCacheEntry`: Reserves a cache entry and returns an upload URL
- `FinalizeCacheEntryUpload`: Commits a cache entry after upload
- `GetCacheEntryDownloadURL`: Retrieves a download URL for a cache entry

## Storage Providers

Derp supports multiple storage providers, but ALL return Azure-style URLs:

- **Azure Blob Storage**: Returns URLs with `azure` provider prefix
- **AWS S3**: Returns URLs with `s3` provider prefix
- **Cloudflare R2**: Returns URLs with `r2` provider prefix
- **Google Cloud Storage**: Returns URLs with `gcs` provider prefix

## Example Usage

See `example_integration.go` for a complete example of how asur can integrate with derp and handle different storage providers.

## Why Azure-style URLs for All Providers?

Using a consistent URL format for all storage providers ensures:

1. All storage requests go through asur for proper authentication
2. The toolkit doesn't need to know about different storage backends
3. Credentials are never exposed to the client
4. A single proxy (asur) can handle all storage operations
5. Asur can retrieve all necessary information from just the URL

## Development

To run the service:

```bash
# Set environment variables
export WARPCACHE_BACKEND_URL="https://api.warpbuild.com"
export WARPCACHE_AUTH_TOKEN="your-token"

# Start the service (default port: 8080)
go run cmd/derp/main.go
```

## Testing Credential Retrieval

You can test the URL-based credential retrieval:

```go
// Example URL
url := "https://warpbuild.blob.core.windows.net/s3/test-cache--v1"

// Parse and retrieve credentials
provider, creds, found := derp.GetCredentialsFromURL(url)
fmt.Printf("Provider: %s, Found: %v\n", provider, found)
```
