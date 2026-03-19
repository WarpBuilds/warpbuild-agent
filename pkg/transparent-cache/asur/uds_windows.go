//go:build windows

package asur

import (
	"net/http"
	"os"
)

// UDSSocketPath returns the ASUR Unix Domain Socket path.
// Set ASUR_UDS_PATH environment variable to override the default.
func UDSSocketPath() string {
	if p := os.Getenv("ASUR_UDS_PATH"); p != "" {
		return p
	}
	return "/tmp/warpbuild-asur.sock"
}

// startUDSListener is a no-op on Windows.
func startUDSListener(handler http.Handler) {}
