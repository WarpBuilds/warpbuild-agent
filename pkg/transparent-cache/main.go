package transparentcache

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/warpbuilds/warpbuild-agent/pkg/transparent-cache/asur"
	"github.com/warpbuilds/warpbuild-agent/pkg/transparent-cache/derp"
	"github.com/warpbuilds/warpbuild-agent/pkg/transparent-cache/oginy"
)

func main() {
	log.Println("========================================")
	log.Println("Starting Transparent Cache Services")
	log.Println("========================================")

	// Default ports (can be overridden by environment variables)
	derpPort := 50051
	oginyPort := 50052
	asurPort := 50053

	// Allow port configuration via environment variables
	if p := os.Getenv("DERP_PORT"); p != "" {
		if port, err := strconv.Atoi(p); err == nil {
			derpPort = port
		}
	}
	if p := os.Getenv("OGINY_PORT"); p != "" {
		if port, err := strconv.Atoi(p); err == nil {
			oginyPort = port
		}
	}
	if p := os.Getenv("ASUR_PORT"); p != "" {
		if port, err := strconv.Atoi(p); err == nil {
			asurPort = port
		}
	}

	// Create a WaitGroup to manage all services
	var wg sync.WaitGroup

	// Channel to collect errors from services
	errChan := make(chan error, 3)

	// Start DERP service (Cache Twirp server)
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("Starting DERP Cache Service on port %d...", derpPort)
		if err := derp.Start(derpPort); err != nil {
			errChan <- fmt.Errorf("DERP service error: %v", err)
		}
	}()

	// Start ASUR service (Azure→S3 Proxy)
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("Starting ASUR Azure→S3 Proxy on port %d...", asurPort)
		if err := asur.Start(asurPort); err != nil {
			errChan <- fmt.Errorf("ASUR service error: %v", err)
		}
	}()

	// Start OGINY service (TLS reverse proxy)
	// Only start if config file exists or is specified
	oginyConfig := os.Getenv("OGINY_CONFIG_PATH")
	if oginyConfig == "" {
		oginyConfig = "config.json"
	}

	if _, err := os.Stat(oginyConfig); err == nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("Starting OGINY TLS Reverse Proxy on port %d...", oginyPort)
			if err := oginy.Start(oginyConfig, oginyPort); err != nil {
				errChan <- fmt.Errorf("OGINY service error: %v", err)
			}
		}()
	} else {
		log.Printf("OGINY: Skipping service - config file not found: %s", oginyConfig)
	}

	// Monitor for errors in a separate goroutine
	go func() {
		for err := range errChan {
			log.Printf("Service error: %v", err)
			// Optionally exit on first error
			// os.Exit(1)
		}
	}()

	// Wait for all services to complete (they shouldn't unless there's an error)
	wg.Wait()
	close(errChan)

	log.Println("All services have stopped")
}
