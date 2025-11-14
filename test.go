package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	transparentcache "github.com/warpbuilds/warpbuild-agent/pkg/transparent-cache"
)

func main() {
	// Define command line flags
	var (
		derpPort       = flag.Int("derp-port", 59992, "Port for DERP (Cache Twirp) service")
		asurPort       = flag.Int("asur-port", 59993, "Port for ASUR (Azure→S3 Proxy) service")
		oginyPort      = flag.Int("oginy-port", 59991, "Port for OGINY (TLS Reverse Proxy) service. Set to 0 to disable")
		backendURL     = flag.String("backend-url", "http://localhost:8000", "WarpCache backend URL")
		authToken      = flag.String("auth-token", "", "WarpCache authentication token")
		debug          = flag.Bool("debug", true, "Enable debug mode")
		skipNetworking = flag.Bool("skip-networking", false, "Skip networking setup (useful for local testing)")
	)

	flag.Parse()

	// Setup environment variables
	setupEnvironment(*backendURL, *authToken, *debug, *skipNetworking)
	// If no auth token was provided via flag, fallback to env set by setupEnvironment
	if *authToken == "" {
		if v := os.Getenv("WARPBUILD_RUNNER_VERIFICATION_TOKEN"); v != "" {
			*authToken = v
		}
	}

	log.Println("========================================")
	log.Println("Transparent Cache Test Runner")
	log.Println("========================================")
	log.Printf("DERP Port: %d", *derpPort)
	log.Printf("ASUR Port: %d", *asurPort)
	log.Printf("OGINY Port: %d", *oginyPort)
	log.Printf("Backend URL: %s", *backendURL)
	log.Printf("Debug Mode: %v", *debug)
	log.Printf("Skip Networking: %v", *skipNetworking)
	log.Println("========================================")

	// If skip networking is enabled, disable OGINY
	if *skipNetworking {
		log.Println("Networking setup disabled, setting OGINY port to 0")
		*oginyPort = 0
	}

	// Create a context that can be cancelled
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Channel to receive error from Start function
	errChan := make(chan error, 1)

	// Start the transparent cache services in a goroutine
	go func() {
		// Use empty string for certDir to use default behavior
		err := transparentcache.Start(*derpPort, *oginyPort, *asurPort, *backendURL, *authToken, "", *debug)
		errChan <- err
	}()

	// Wait for either an error or interrupt signal
	select {
	case err := <-errChan:
		if err != nil {
			log.Fatalf("Transparent cache failed: %v", err)
		}
		log.Println("Transparent cache stopped successfully")
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
		log.Println("Shutting down transparent cache...")
		cancel()
		// Give services time to shutdown gracefully
		// In a real implementation, we'd need to modify the Start function
		// to accept a context for proper shutdown
		os.Exit(0)
	}
}

func setupEnvironment(backendURL, authToken string, debug, skipNetworking bool) {
	os.Setenv("WARPBUILD_RUNNER_VERIFICATION_TOKEN", "<Add warpbuild token here>")
	// Set up environment variables for DERP
	if backendURL != "" {
		os.Setenv("WARPCACHE_BACKEND_URL", backendURL)
	}
	if authToken != "" {
		os.Setenv("WARPCACHE_AUTH_TOKEN", authToken)
	}

	// Set up debug mode for ASUR
	if debug {
		os.Setenv("AZPROXY_DEBUG", "true")
	}

	// Example R2 credentials for ASUR (these would be real in production)
	// You can override these with actual values
	if os.Getenv("R2_ACCESS_KEY_ID") == "" {
		log.Println("WARNING: R2_ACCESS_KEY_ID not set. ASUR may not work properly for R2 storage.")
		// Set dummy values for testing
		os.Setenv("R2_ACCESS_KEY_ID", "test_access_key")
		os.Setenv("R2_SECRET_ACCESS_KEY", "test_secret_key")
		os.Setenv("R2_ENDPOINT", "https://test.r2.cloudflarestorage.com")
	}

	// Set default upload method and concurrency
	if os.Getenv("AZPROXY_UPLOAD_METHOD") == "" {
		os.Setenv("AZPROXY_UPLOAD_METHOD", "http")
	}
	if os.Getenv("AZPROXY_UPLOAD_CONCURRENCY") == "" {
		os.Setenv("AZPROXY_UPLOAD_CONCURRENCY", "10")
	}

	// If we're in test mode and have a custom ACTIONS_RESULTS_URL, set it
	if os.Getenv("ACTIONS_RESULTS_URL") == "" && !skipNetworking {
		// Use the default GitHub Actions results URL
		os.Setenv("ACTIONS_RESULTS_URL", "https://results-receiver.actions.githubusercontent.com")
	}

	log.Println("Environment setup completed")
	log.Printf("WARPCACHE_BACKEND_URL: %s", os.Getenv("WARPCACHE_BACKEND_URL"))
	log.Printf("AZPROXY_DEBUG: %s", os.Getenv("AZPROXY_DEBUG"))
	log.Printf("AZPROXY_UPLOAD_METHOD: %s", os.Getenv("AZPROXY_UPLOAD_METHOD"))
}

func init() {
	// Set up logging format
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Print usage information
	fmt.Println("Transparent Cache Test Runner")
	fmt.Println("=============================")
	fmt.Println("This tool starts all transparent cache services for testing.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run test.go [flags]")
	fmt.Println()
	fmt.Println("For local testing without root privileges:")
	fmt.Println("  go run test.go --skip-networking")
	fmt.Println()
	fmt.Println("For production-like setup (requires root):")
	fmt.Println("  sudo go run test.go --auth-token=<your-token>")
	fmt.Println()
}
