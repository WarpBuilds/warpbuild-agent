package uploader

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRedactURL(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "presigned query is replaced",
			in:   "https://bucket.s3.us-east-1.amazonaws.com/org/run/file.log?X-Amz-Signature=sig&X-Amz-Security-Token=tok",
			want: "https://bucket.s3.us-east-1.amazonaws.com/org/run/file.log?REDACTED",
		},
		{
			name: "URL without query is unchanged",
			in:   "https://bucket.s3.amazonaws.com/org/run/file.log",
			want: "https://bucket.s3.amazonaws.com/org/run/file.log",
		},
		{
			name: "user info and fragment are removed",
			in:   "https://user:pass@host/path?a=b#frag",
			want: "https://host/path?REDACTED",
		},
		{
			name: "unparseable URL is fully replaced",
			in:   "http://[::1]:port/path?X-Amz-Signature=sig",
			want: "REDACTED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := redactURL(tt.in); got != tt.want {
				t.Errorf("redactURL(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestRedactURLErrorKeepsWrappedError(t *testing.T) {
	err := &url.Error{
		Op:  "Put",
		URL: "https://bucket.s3.amazonaws.com/file.log?X-Amz-Signature=secret",
		Err: context.DeadlineExceeded,
	}

	got := redactURLError(err)

	if strings.Contains(got.Error(), "secret") {
		t.Errorf("error still contains the presigned query: %v", got)
	}
	if !errors.Is(got, context.DeadlineExceeded) {
		t.Errorf("redactURLError lost the wrapped error: %v", got)
	}
}

func TestRedactURLErrorIgnoresOtherErrors(t *testing.T) {
	err := errors.New("plain error")
	if got := redactURLError(err); got != err {
		t.Errorf("redactURLError changed a non-URL error: %v", got)
	}
}

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

func TestUploadToS3BadURLDoesNotLeakPresignedQuery(t *testing.T) {
	s := &S3Uploader{
		eventType:    "logs",
		presignedURL: "http://[::1]:port/org/file.log?X-Amz-Signature=secret-sig",
	}

	err := s.uploadToS3([]byte("data"))
	if err == nil {
		t.Fatal("expected a request error")
	}
	if strings.Contains(err.Error(), "secret") {
		t.Errorf("request error contains the presigned query: %v", err)
	}
}
