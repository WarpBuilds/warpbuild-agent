package proxy

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUploadToBlobStorageErrorDoesNotLeakPresignedQuery(t *testing.T) {
	// The server drops the connection, so http.DefaultClient.Do returns a
	// *url.Error that holds the full request URL.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			t.Errorf("hijack: %v", err)
			return
		}
		conn.Close()
	}))
	defer srv.Close()

	presignedURL := srv.URL + "/cache?X-Amz-Signature=secret-sig&sig=secret-sas"
	responses := map[string]ReserveCacheResponse{
		"s3": {
			Provider: ProviderS3,
			S3:       &S3ReserveCacheResponse{PreSignedURLs: []string{presignedURL}},
		},
		"azure_blob": {
			Provider:  ProviderAzureBlob,
			AzureBlob: &AzureBlobReserveCacheResponse{PreSignedURL: presignedURL},
		},
	}

	cacheID := -1
	for name, resp := range responses {
		t.Run(name, func(t *testing.T) {
			cacheStore.Store(cacheID, &CacheEntryData{
				BackendReserveResponse: resp,
				Chunks:                 map[int64]ChunkData{0: {StartOffset: 0, EndOffset: 3, Content: []byte("data")}},
			})
			defer cacheStore.Delete(cacheID)

			err := uploadToBlobStorage(t.Context(), cacheID)
			if err == nil {
				t.Fatal("expected an upload error")
			}
			if strings.Contains(err.Error(), "secret") {
				t.Errorf("upload error contains the presigned query: %v", err)
			}
		})
	}
}
