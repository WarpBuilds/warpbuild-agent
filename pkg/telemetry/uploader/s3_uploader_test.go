package uploader

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUploadToS3ErrorDoesNotLeakPresignedQuery(t *testing.T) {
	// The server drops the connection, so client.Do returns a *url.Error
	// that holds the full request URL.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			t.Errorf("hijack: %v", err)
			return
		}
		conn.Close()
	}))
	defer srv.Close()

	s := &S3Uploader{
		eventType:    "logs",
		presignedURL: srv.URL + "/org/file.log?X-Amz-Signature=secret-sig&X-Amz-Security-Token=secret-token",
	}

	err := s.uploadToS3([]byte("data"))
	if err == nil {
		t.Fatal("expected an upload error")
	}
	if strings.Contains(err.Error(), "secret") {
		t.Errorf("upload error contains the presigned query: %v", err)
	}
	if !strings.Contains(err.Error(), "/org/file.log?REDACTED") {
		t.Errorf("upload error does not show the redacted URL: %v", err)
	}
}
