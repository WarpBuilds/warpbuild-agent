package transparentcache

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/warpbuilds/warpbuild-agent/pkg/transparent-cache/asur"
	"github.com/warpbuilds/warpbuild-agent/pkg/transparent-cache/derp"
	"github.com/warpbuilds/warpbuild-agent/pkg/transparent-cache/oginy"
)

// setupNetworking configures loopback IPs and nftables rules for transparent caching
func setupNetworking(oginyPort int) error {
	log.Println("Setting up networking for transparent cache...")

	// Skip networking setup if OGINY is not enabled
	if oginyPort <= 0 {
		log.Println("OGINY port not configured; skipping networking setup")
		return nil
	}

	// Ensure we have root privileges for system networking changes
	if os.Geteuid() != 0 {
		return fmt.Errorf("networking setup requires root privileges; please run as root or with sudo")
	}

	// Get results-receiver hostname from environment variable or use default
	resultsReceiverHost := "results-receiver.actions.githubusercontent.com"
	if actionsResultsURL := os.Getenv("ACTIONS_RESULTS_URL"); actionsResultsURL != "" {
		if parsedURL, err := url.Parse(actionsResultsURL); err == nil && parsedURL.Host != "" {
			resultsReceiverHost = parsedURL.Hostname()
			log.Printf("Using results-receiver hostname from ACTIONS_RESULTS_URL: %s", resultsReceiverHost)
		} else {
			log.Printf("Failed to parse ACTIONS_RESULTS_URL (%s), using default: %s", actionsResultsURL, resultsReceiverHost)
		}
	} else {
		log.Printf("ACTIONS_RESULTS_URL not set, using default results-receiver hostname: %s", resultsReceiverHost)
	}

	// Define the hostname to IP mappings
	hostMappings := map[string]string{
		"warpbuild.blob.core.windows.net": "127.77.77.77",
		resultsReceiverHost:               "127.77.77.78",
	}

	// Add loopback aliases
	for hostname, ip := range hostMappings {
		log.Printf("Adding loopback alias %s for %s", ip, hostname)

		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			// On macOS, use ifconfig to add loopback alias
			cmd = exec.Command("ifconfig", "lo0", "alias", ip, "netmask", "255.255.255.255")
		case "linux":
			// On Linux, use ip command
			cmd = exec.Command("ip", "addr", "add", fmt.Sprintf("%s/32", ip), "dev", "lo")
		default:
			return fmt.Errorf("unsupported OS for loopback setup: %s", runtime.GOOS)
		}

		if output, err := cmd.CombinedOutput(); err != nil {
			// Check if the address already exists (not an error in this case)
			errMsg := string(output)
			if !strings.Contains(errMsg, "File exists") && !strings.Contains(errMsg, "already exists") {
				return fmt.Errorf("failed to add loopback alias %s: %v - %s", ip, err, errMsg)
			}
			log.Printf("Loopback alias %s already exists", ip)
		}
	}

	// Update /etc/hosts
	log.Println("Updating /etc/hosts...")
	hostsFile, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open /etc/hosts: %v", err)
	}
	defer hostsFile.Close()

	for hostname, ip := range hostMappings {
		// Check if entry already exists
		cmd := exec.Command("grep", "-q", fmt.Sprintf("%s %s", ip, hostname), "/etc/hosts")
		if err := cmd.Run(); err != nil {
			// Entry doesn't exist, add it
			entry := fmt.Sprintf("%s %s\n", ip, hostname)
			if _, err := hostsFile.WriteString(entry); err != nil {
				return fmt.Errorf("failed to write to /etc/hosts: %v", err)
			}
			log.Printf("Added /etc/hosts entry: %s", entry)
		} else {
			log.Printf("/etc/hosts entry already exists for %s", hostname)
		}
	}

	// Setup firewall rules based on OS
	switch runtime.GOOS {
	case "darwin":
		if err := setupPfctl(hostMappings, oginyPort); err != nil {
			return fmt.Errorf("failed to setup pfctl: %v", err)
		}
	case "linux":
		if err := setupNftables(hostMappings, oginyPort); err != nil {
			return fmt.Errorf("failed to setup nftables: %v", err)
		}
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	log.Println("Networking setup completed successfully")
	return nil
}

// setupPfctl configures packet filter rules on macOS
func setupPfctl(hostMappings map[string]string, oginyPort int) error {
	log.Println("Setting up pfctl rules for macOS...")

	// Create pf.conf content
	var pfRules strings.Builder
	pfRules.WriteString("# Transparent cache redirect rules\n")

	// First, ensure oginy's own egress to the loopback-mapped IPs is NOT redirected
	for hostname, ip := range hostMappings {
		noRdr := fmt.Sprintf("no rdr on lo0 inet proto tcp from 127.0.0.1 to %s port 443\n", ip)
		pfRules.WriteString(noRdr)
		log.Printf("Adding pfctl no-rdr rule for oginy egress to %s (%s)", hostname, ip)
	}

	for hostname, ip := range hostMappings {
		rule := fmt.Sprintf("rdr pass on lo0 inet proto tcp from any to %s port 443 -> 127.0.0.1 port %d\n", ip, oginyPort)
		pfRules.WriteString(rule)
		log.Printf("Adding pfctl redirect rule for %s (%s) -> port %d", hostname, ip, oginyPort)
	}

	// Write rules to a temporary file
	tmpFile := "/tmp/transparent-cache-pf.conf"
	if err := os.WriteFile(tmpFile, []byte(pfRules.String()), 0644); err != nil {
		return fmt.Errorf("failed to write pf rules: %v", err)
	}
	defer os.Remove(tmpFile)

	// Load the rules
	cmd := exec.Command("pfctl", "-f", tmpFile)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to load pf rules: %v - %s", err, string(output))
	}

	// Enable pfctl
	cmd = exec.Command("pfctl", "-e")
	if output, err := cmd.CombinedOutput(); err != nil {
		// pfctl might already be enabled, check the error message
		if !strings.Contains(string(output), "already enabled") {
			return fmt.Errorf("failed to enable pfctl: %v - %s", err, string(output))
		}
		log.Println("pfctl already enabled")
	}

	return nil
}

// setupNftables configures nftables rules on Linux
func setupNftables(hostMappings map[string]string, oginyPort int) error {
	log.Println("Setting up nftables rules for Linux...")

	// Create table and chain if they don't exist
	nftCommands := [][]string{
		{"nft", "add", "table", "ip", "nat"},
		{"nft", "add", "chain", "ip", "nat", "output", "{", "type", "nat", "hook", "output", "priority", "0;", "}"},
	}

	for _, cmdArgs := range nftCommands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		if output, err := cmd.CombinedOutput(); err != nil {
			// Ignore errors if table/chain already exists
			log.Printf("nft command output: %s", string(output))
		}
	}

	// Add return (skip) rules for oginy egress first, so we don't redirect its own outbound connections
	for hostname, ip := range hostMappings {
		log.Printf("Adding nftables return rule for oginy egress to %s (%s)", hostname, ip)
		cmd := exec.Command("nft", "add", "rule", "ip", "nat", "output",
			"ip", "saddr", "127.0.0.1", "ip", "daddr", ip, "tcp", "dport", "443", "return")
		if output, err := cmd.CombinedOutput(); err != nil {
			outStr := string(output)
			if strings.Contains(outStr, "exist") {
				log.Printf("nftables return rule for %s already exists", ip)
				continue
			}
			log.Printf("nft command output: %s", outStr)
		}
	}

	// Add redirect rules for each IP
	for hostname, ip := range hostMappings {
		log.Printf("Adding nftables redirect rule for %s (%s) -> port %d", hostname, ip, oginyPort)
		cmd := exec.Command("nft", "add", "rule", "ip", "nat", "output",
			"ip", "daddr", ip, "tcp", "dport", "443", "redirect", "to", fmt.Sprintf(":%d", oginyPort))
		if output, err := cmd.CombinedOutput(); err != nil {
			outStr := string(output)
			if strings.Contains(outStr, "exist") {
				log.Printf("nftables rule for %s already exists", ip)
				continue
			}
			return fmt.Errorf("failed to add nftables rule for %s: %v - %s", ip, err, outStr)
		}
	}

	return nil
}

func Start(derpPort, oginyPort, asurPort int) error {
	log.Println("========================================")
	log.Println("Starting Transparent Cache Services")
	log.Println("========================================")

	// Setup networking before starting services
	if oginyPort > 0 {
		if err := setupNetworking(oginyPort); err != nil {
			return fmt.Errorf("failed to setup networking: %v", err)
		}
	} else {
		log.Println("OGINY port not set; skipping networking setup")
	}

	// Create a WaitGroup to manage all services
	var wg sync.WaitGroup

	// Channel to collect errors from services
	errChan := make(chan error, 3)
	// Channel to signal first error
	firstErr := make(chan error, 1)

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
		log.Printf("Starting ASUR Azure→R2 Proxy on port %d...", asurPort)
		if err := asur.Start(asurPort); err != nil {
			errChan <- fmt.Errorf("ASUR service error: %v", err)
		}
	}()

	// Start OGINY service (TLS reverse proxy)
	if oginyPort > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("Starting OGINY TLS Reverse Proxy on port %d...", oginyPort)
			if err := oginy.Start(oginyPort); err != nil {
				errChan <- fmt.Errorf("OGINY service error: %v", err)
			}
		}()
	}

	// Monitor for errors in a separate goroutine
	go func() {
		for err := range errChan {
			log.Printf("Service error: %v", err)
			// Send first error to the main function
			select {
			case firstErr <- err:
			default:
			}
		}
	}()

	// Wait for all services to complete (they shouldn't unless there's an error)
	wg.Wait()
	close(errChan)

	// Check if any error occurred
	select {
	case err := <-firstErr:
		return err
	default:
		log.Println("All services have stopped")
		return nil
	}
}
