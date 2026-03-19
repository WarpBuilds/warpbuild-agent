//go:build !windows

package asur

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

// UDSSocketPath returns the path to the ASUR Unix Domain Socket.
// Set ASUR_UDS_PATH environment variable to override the default.
func UDSSocketPath() string {
	if p := os.Getenv("ASUR_UDS_PATH"); p != "" {
		return p
	}
	return "/tmp/warpbuild-asur.sock"
}

// startUDSListener starts an HTTP server on a Unix Domain Socket for small-payload requests.
// This runs alongside the main TCP listener to give latency/throughput benefits for
// payloads under 32KB (see benchmark results in pkg/transparent-cache/benchmark/uds_vs_tcp.go).
func startUDSListener(handler http.Handler) {
	path := UDSSocketPath()
	os.Remove(path)

	ln, err := net.Listen("unix", path)
	if err != nil {
		log.Printf("Warning: Failed to start ASUR UDS listener at %s: %v. Using TCP only.", path, err)
		return
	}

	if err := os.Chmod(path, 0666); err != nil {
		log.Printf("Warning: Failed to chmod ASUR UDS socket: %v", err)
	}

	log.Printf("ASUR UDS listener started at %s", path)

	udsSrv := &http.Server{
		Handler:        handler,
		ReadTimeout:    30 * time.Minute,
		WriteTimeout:   30 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := udsSrv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Printf("ASUR UDS server error: %v", err)
		}
	}()
}
