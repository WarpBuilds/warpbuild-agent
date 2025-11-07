package proxy

import (
	"encoding/xml"
	"time"
)

// Docker GHA backend dto

type CacheBackendInfo struct {
	HostURL   string `json:"hostURL"`
	AuthToken string `json:"authToken"`
}

type DockerGHAGetCacheRequest struct {
	Keys    []string `json:"keys"`
	Version string   `json:"version"`
}

type DockerGHAGetCacheResponse struct {
	CacheKey        string `json:"cacheKey"`
	ArchiveLocation string `json:"archiveLocation"`
}

type DockerGHAReserveCacheRequest struct {
	Key     string `json:"key"`
	Version string `json:"version"`
}

type DockerGHAReserveCacheResponse struct {
	CacheID int `json:"cacheID"`
}

type DockerGHAUploadCacheRequest struct {
	CacheID      int    `json:"cacheID"`
	Content      []byte `json:"content"`
	ContentRange string `json:"contentRange"`
}

type DockerGHAUploadCacheResponse struct{}

type DockerGHACommitCacheRequest struct {
	CacheID int   `json:"cacheID"`
	Size    int64 `json:"size"`
}

type DockerGHACommitCacheResponse struct{}

// WarpCache DTO
type Provider string

const (
	ProviderGCS       Provider = "gcs"
	ProviderS3        Provider = "s3"
	ProviderAzureBlob Provider = "azure_blob"
	ProviderR2        Provider = "r2"
)

type AuthMethod string

const (
	MethodShortLivedToken AuthMethod = "short_lived_token"
	MethodPresignedURL    AuthMethod = "presigned_url"
)

type CacheEntry struct {
	ID                     string             `json:"id"`
	CreatedAt              time.Time          `json:"created_at"`
	UpdatedAt              time.Time          `json:"updated_at"`
	StorageBackendId       string             `json:"storage_backend_id"`
	StorageBackendLocation string             `json:"storage_backend_location"`
	CacheKey               string             `json:"cache_key"`
	CacheUserGivenKey      string             `json:"cache_user_given_key"`
	CacheVersion           string             `json:"cache_version"`
	VCSOrganizationName    string             `json:"vcs_organization_name"`
	VCSRepositoryName      string             `json:"vcs_repository_name"`
	VCSRef                 string             `json:"vcs_ref"`
	OrganizationID         string             `json:"organization_id"`
	Provider               Provider           `json:"provider" enum:"gcs,s3,azure_blob"`
	Metadata               CacheEntryMetadata `json:"metadata"`
}

type GetCacheRequest struct {
	CacheKey     string   `json:"cache_key" validate:"required"`
	CacheVersion string   `json:"cache_version" validate:"required"`
	RestoreKeys  []string `json:"restore_keys"`
}

type GetCacheResponse struct {
	Provider   Provider                   `json:"provider" enum:"gcs,s3"`
	GCS        *GCSGetCacheResponse       `json:"gcs,omitempty"`
	S3         *S3GetCacheResponse        `json:"s3,omitempty"`
	AzureBlob  *AzureBlobGetCacheResponse `json:"azure_blob,omitempty"`
	CacheEntry *CacheEntry                `json:"cache_entry"`
}

type CacheEntryMetadata struct {
	StackId           string `json:"stack_id"`
	StackName         string `json:"stack_name"`
	CloudConnectionId string `json:"cloud_connection_id"`
}

type S3GetCacheResponse struct {
	CacheKey     string `json:"cache_key"`
	CacheVersion string `json:"cache_version"`
	PreSignedURL string `json:"pre_signed_url"`
}

type GCSGetCacheResponse struct {
	Method          AuthMethod       `json:"method" enum:"short_lived_token"`
	ShortLivedToken *ShortLivedToken `json:"short_lived_token,omitempty"`
	PreSignedURL    string           `json:"pre_signed_url"`
	BucketName      string           `json:"bucket_name"`
	ProjectID       string           `json:"project_id"`
	CacheKey        string           `json:"cache_key"`
	CacheVersion    string           `json:"cache_version"`
}

type AzureBlobGetCacheResponse struct {
	PreSignedURL string `json:"pre_signed_url"`
	CacheKey     string `json:"cache_key"`
	CacheVersion string `json:"cache_version"`
	BucketName   string `json:"bucket_name"`
}

type ShortLivedToken struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type ReserveCacheRequest struct {
	CacheKey       string `json:"cache_key" validate:"required"`
	CacheVersion   string `json:"cache_version" validate:"required"`
	NumberOfChunks uint   `json:"number_of_chunks"`
	ContentType    string `json:"content_type"`
}

type ReserveCacheResponse struct {
	Provider  Provider                       `json:"provider" enum:"gcs,s3"`
	GCS       *GCSReserveCacheResponse       `json:"gcs,omitempty"`
	S3        *S3ReserveCacheResponse        `json:"s3,omitempty"`
	AzureBlob *AzureBlobReserveCacheResponse `json:"azure_blob,omitempty"`
}

type GCSReserveCacheResponse struct {
	Method          AuthMethod       `json:"method" enum:"short_lived_token"`
	ShortLivedToken *ShortLivedToken `json:"short_lived_token,omitempty"`
	BucketName      string           `json:"bucket_name"`
	ProjectID       string           `json:"project_id"`
	CacheKey        string           `json:"cache_key" validate:"required"`
}

type S3ReserveCacheResponse struct {
	PreSignedURLs []string `json:"pre_signed_urls"`
	UploadKey     string   `json:"upload_key"`
	UploadID      string   `json:"upload_id"`
}

type AzureBlobReserveCacheResponse struct {
	PreSignedURL  string `json:"pre_signed_url"`
	ContainerName string `json:"container_name"`
	BlobName      string `json:"blob_name"`
}

type CommitCacheRequest struct {
	CacheKey     string            `json:"cache_key" validate:"required"`
	CacheVersion string            `json:"cache_version" validate:"required"`
	UploadKey    string            `json:"upload_key"`
	UploadID     string            `json:"upload_id"`
	Parts        []S3CompletedPart `json:"parts" validate:"required"`
	VCSType      string            `json:"vcs_type" validate:"required"`
	Provider     Provider          `json:"provider" enum:"gcs,s3"`
}

type CommitCacheResponse struct {
	CacheEntry *CacheEntry                   `json:"cache_entry"`
	Provider   Provider                      `json:"provider" enum:"gcs,s3"`
	GCS        *GCSCommitCacheResponse       `json:"gcs,omitempty"`
	S3         *S3CommitCacheResponse        `json:"s3,omitempty"`
	AzureBlob  *AzureBlobCommitCacheResponse `json:"azure_blob,omitempty"`
}

type GCSCommitCacheResponse struct {
	Method          AuthMethod       `json:"method" enum:"short_lived_token"`
	ShortLivedToken *ShortLivedToken `json:"short_lived_token,omitempty"`
	BucketName      string           `json:"bucket_name"`
	ProjectID       string           `json:"project_id"`
	CacheKey        string           `json:"cache_key" validate:"required"`
}

type S3CommitCacheResponse struct {
	CacheKey     string `json:"cache_key" validate:"required"`
	CacheVersion string `json:"cache_version" validate:"required"`
}

type AzureBlobCommitCacheResponse struct {
	CacheKey     string `json:"cache_key" validate:"required"`
	CacheVersion string `json:"cache_version" validate:"required"`
}

type DeleteCacheRequest struct {
	CacheKey     string `json:"cache_key" validate:"required"`
	CacheVersion string `json:"cache_version" validate:"required"`
}

type DeleteCacheResponse struct {
	CacheEntry *CacheEntry                   `json:"cache_entry"`
	Provider   Provider                      `json:"provider" enum:"gcs,s3,azure_blob"`
	GCS        *GCSDeleteCacheResponse       `json:"gcs,omitempty"`
	S3         *S3DeleteCacheResponse        `json:"s3,omitempty"`
	AzureBlob  *AzureBlobDeleteCacheResponse `json:"azure_blob,omitempty"`
}

type GCSDeleteCacheResponse struct {
	CacheKey     string `json:"cache_key" validate:"required"`
	CacheVersion string `json:"cache_version" validate:"required"`
}

type S3DeleteCacheResponse struct {
	CacheKey     string `json:"cache_key" validate:"required"`
	CacheVersion string `json:"cache_version" validate:"required"`
}

type AzureBlobDeleteCacheResponse struct {
	CacheKey     string `json:"cache_key" validate:"required"`
	CacheVersion string `json:"cache_version" validate:"required"`
}

type S3PartNumber uint16

// Taken from s3 v2 sdk
type S3CompletedPart struct {

	// Entity tag returned when the part was uploaded.
	ETag *string

	// Part number that identifies the part. This is a positive integer between 1 and
	// 10,000.
	//
	//   - General purpose buckets - In CompleteMultipartUpload , when a additional
	//   checksum (including x-amz-checksum-crc32 , x-amz-checksum-crc32c ,
	//   x-amz-checksum-sha1 , or x-amz-checksum-sha256 ) is applied to each part, the
	//   PartNumber must start at 1 and the part numbers must be consecutive.
	//   Otherwise, Amazon S3 generates an HTTP 400 Bad Request status code and an
	//   InvalidPartOrder error code.
	//
	//   - Directory buckets - In CompleteMultipartUpload , the PartNumber must start
	//   at 1 and the part numbers must be consecutive.
	PartNumber S3PartNumber
}

// ExtendReserveCacheRequest is used to request additional presigned URLs for an existing upload
type ExtendReserveCacheRequest struct {
	CacheKey     string `json:"cache_key"`
	CacheVersion string `json:"cache_version"`
	UploadID     string `json:"upload_id"`
	FromPart     int32  `json:"from_part"`
	Count        int32  `json:"count"`
}

// ExtendReserveCacheResponse contains additional presigned URLs
type ExtendReserveCacheResponse struct {
	PreSignedURLs []string `json:"pre_signed_urls"`
}

// blockListXML represents Azure Block Blob block list for commit
type blockListXML struct {
	XMLName xml.Name `xml:"BlockList"`
	Latest  []string `xml:"Latest"`
}
