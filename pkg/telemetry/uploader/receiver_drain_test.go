package uploader

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleDrain(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		remoteAddr string
		onDrain    bool
		wantStatus int
		wantCalled bool
	}{
		{"loopback post", http.MethodPost, "127.0.0.1:5555", true, http.StatusOK, true},
		{"ipv6 loopback", http.MethodPost, "[::1]:5555", true, http.StatusOK, true},
		// The port is local to the box, but so is the customer's job code.
		{"remote rejected", http.MethodPost, "10.1.2.3:5555", true, http.StatusForbidden, false},
		{"get rejected", http.MethodGet, "127.0.0.1:5555", true, http.StatusMethodNotAllowed, false},
		{"no callback", http.MethodPost, "127.0.0.1:5555", false, http.StatusNotImplemented, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			r := NewReceiver(0, nil)
			if tc.onDrain {
				r.SetOnDrain(func() { called = true })
			}

			req := httptest.NewRequest(tc.method, "/internal/drain", nil)
			req.RemoteAddr = tc.remoteAddr
			rec := httptest.NewRecorder()

			r.handleDrain(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status %d, want %d", rec.Code, tc.wantStatus)
			}
			if called != tc.wantCalled {
				t.Errorf("callback fired=%t, want %t", called, tc.wantCalled)
			}
		})
	}
}

func TestIsLoopback(t *testing.T) {
	loopback := []string{"127.0.0.1:1", "[::1]:1", "127.0.0.1"}
	remote := []string{"10.0.0.1:1", "192.168.1.5:1", "8.8.8.8:1", "", "not-an-ip:1"}

	for _, addr := range loopback {
		if !isLoopback(addr) {
			t.Errorf("%q should be loopback", addr)
		}
	}
	for _, addr := range remote {
		if isLoopback(addr) {
			t.Errorf("%q should not be loopback", addr)
		}
	}
}
