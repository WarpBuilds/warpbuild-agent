package oginy

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Global logger and mutex for thread-safe file writing
var (
	proxyLogger *log.Logger
	logMutex    sync.Mutex
	logFile     *os.File
)

// initLogger initializes the file logger for proxy traffic
func initLogger(logDir string) error {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("create log dir: %v", err)
	}

	// Create log file with timestamp
	logFileName := fmt.Sprintf("oginy-proxy-%s.log", time.Now().Format("2006-01-02-15-04-05"))
	logPath := filepath.Join(logDir, logFileName)

	var err error
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("open log file: %v", err)
	}

	proxyLogger = log.New(logFile, "", log.LstdFlags|log.Lmicroseconds)
	proxyLogger.Printf("=== OGINY Proxy Log Started ===")
	log.Printf("Proxy traffic log file: %s", logPath)

	return nil
}

// logRequest logs detailed request information
func logRequest(r *http.Request, destination string) {
	if proxyLogger == nil {
		return
	}

	logMutex.Lock()
	defer logMutex.Unlock()

	// Capture request body if present
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore body
	}

	// Use httputil.DumpRequest for comprehensive request capture
	requestDump, err := httputil.DumpRequest(r, false)
	if err != nil {
		proxyLogger.Printf("[ERROR] Failed to dump request: %v", err)
	}

	proxyLogger.Printf("\n========== REQUEST ==========")
	proxyLogger.Printf("Time: %s", time.Now().Format(time.RFC3339))
	proxyLogger.Printf("Destination: %s", destination)
	proxyLogger.Printf("Client IP: %s", r.RemoteAddr)
	proxyLogger.Printf("\n%s", string(requestDump))

	if len(bodyBytes) > 0 {
		proxyLogger.Printf("Body (%d bytes):", len(bodyBytes))
		// Log body as string if it's printable, otherwise as hex
		if isPrintable(bodyBytes) {
			proxyLogger.Printf("%s", string(bodyBytes))
		} else {
			proxyLogger.Printf("[Binary data - %d bytes]", len(bodyBytes))
		}
	}
	proxyLogger.Printf("========== END REQUEST ==========\n")
}

// isPrintable checks if the byte slice contains printable characters
func isPrintable(data []byte) bool {
	for _, b := range data {
		if b < 32 && b != '\n' && b != '\r' && b != '\t' {
			return false
		}
	}
	return true
}

// loggingRoundTripper wraps an http.RoundTripper to log responses
type loggingRoundTripper struct {
	transport http.RoundTripper
	name      string
}

func (l *loggingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	// Log the request
	logRequest(r, l.name)

	// Execute the request
	resp, err := l.transport.RoundTrip(r)
	if err != nil {
		if proxyLogger != nil {
			logMutex.Lock()
			proxyLogger.Printf("[ERROR] %s request failed: %v", l.name, err)
			logMutex.Unlock()
		}
		return nil, err
	}

	// Log the response
	if proxyLogger != nil && resp != nil {
		logMutex.Lock()
		defer logMutex.Unlock()

		// Capture response for logging
		responseDump, _ := httputil.DumpResponse(resp, false)

		proxyLogger.Printf("\n========== RESPONSE ==========")
		proxyLogger.Printf("Time: %s", time.Now().Format(time.RFC3339))
		proxyLogger.Printf("From: %s", l.name)
		proxyLogger.Printf("Status: %s", resp.Status)
		proxyLogger.Printf("\n%s", string(responseDump))

		// Always log response body up to 1MB
		maxBodySize := int64(1024 * 1024) // 1MB max

		// Log response body if not too large
		if resp.Body != nil && resp.ContentLength > 0 && resp.ContentLength <= maxBodySize {
			// Read the response body
			bodyBytes, err := io.ReadAll(resp.Body)
			if err == nil {
				// Restore the body for the actual consumer
				resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				proxyLogger.Printf("Response Body (%d bytes):", len(bodyBytes))
				if isPrintable(bodyBytes) {
					proxyLogger.Printf("%s", string(bodyBytes))
				} else {
					proxyLogger.Printf("[Binary data - %d bytes]", len(bodyBytes))
				}
			}
		} else if resp.ContentLength > maxBodySize {
			proxyLogger.Printf("Response Body: [Skipped - too large: %d bytes > %d bytes limit]", resp.ContentLength, maxBodySize)
		} else if resp.ContentLength >= 0 {
			proxyLogger.Printf("Content-Length: %d bytes", resp.ContentLength)
		}
		proxyLogger.Printf("========== END RESPONSE ==========\n")
	}

	return resp, nil
}

type proxyEntry struct {
	cert             *tls.Certificate
	target           *url.URL
	proxy            *httputil.ReverseProxy
	remoteProxy      *httputil.ReverseProxy // For proxying to the actual domain
	pathBasedRouting bool                   // Whether to use path-based routing
}

type muxProxy struct{ byHost map[string]*proxyEntry }

func (m *muxProxy) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if pe := m.byHost[hello.ServerName]; pe != nil {
		return pe.cert, nil
	}
	return nil, fmt.Errorf("no certificate for %s", hello.ServerName)
}
func (m *muxProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	if pe := m.byHost[host]; pe != nil {
		// If path-based routing is enabled, check the path
		if pe.pathBasedRouting {
			// Forward Twirp cache service requests to local service
			if strings.HasPrefix(r.URL.Path, "/twirp/github.actions.results.api.v1.CacheService/") {
				log.Printf("[OGINY ROUTING] %s %s → LOCAL SERVICE (port %s)", r.Method, r.URL.Path, pe.target.Host)
				pe.proxy.ServeHTTP(w, r)
			} else if pe.remoteProxy != nil {
				// Forward all other requests to the actual domain
				log.Printf("[OGINY ROUTING] %s %s → REMOTE DOMAIN", r.Method, r.URL.Path)
				pe.remoteProxy.ServeHTTP(w, r)
			} else {
				log.Printf("[OGINY ROUTING] %s %s → ERROR: no remote proxy configured", r.Method, r.URL.Path)
				if proxyLogger != nil {
					logMutex.Lock()
					proxyLogger.Printf("[ERROR] No remote proxy configured for %s %s", r.Method, r.URL.Path)
					logMutex.Unlock()
				}
				http.Error(w, "no remote proxy configured", http.StatusBadGateway)
			}
		} else {
			// Standard routing - forward all requests to local service
			log.Printf("[OGINY ROUTING] %s %s → LOCAL SERVICE (port %s) [no path routing]", r.Method, r.URL.Path, pe.target.Host)
			pe.proxy.ServeHTTP(w, r)
		}
		return
	}
	log.Printf("[OGINY ROUTING] %s %s → ERROR: no backend for host %s", r.Method, r.URL.Path, host)
	if proxyLogger != nil {
		logMutex.Lock()
		proxyLogger.Printf("[ERROR] No backend for host %s", host)
		logMutex.Unlock()
	}
	http.Error(w, "no backend for host", http.StatusBadGateway)
}

// generateCA generates a CA certificate and key
func generateCA(certPath, keyPath string) error {
	// Generate CA private key
	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generate CA key: %v", err)
	}

	// Create CA certificate template
	caTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "WarpBuild Local CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(825 * 24 * time.Hour), // 825 days
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		MaxPathLen:            0,
		MaxPathLenZero:        true,
	}

	// Generate CA certificate
	caCert, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		return fmt.Errorf("create CA cert: %v", err)
	}

	// Save CA certificate
	certOut, err := os.Create(certPath)
	if err != nil {
		return fmt.Errorf("create CA cert file: %v", err)
	}
	defer certOut.Close()
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: caCert}); err != nil {
		return fmt.Errorf("encode CA cert: %v", err)
	}

	// Save CA private key
	keyOut, err := os.Create(keyPath)
	if err != nil {
		return fmt.Errorf("create CA key file: %v", err)
	}
	defer keyOut.Close()
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(caKey)}); err != nil {
		return fmt.Errorf("encode CA key: %v", err)
	}

	return nil
}

// generateLeafCert generates a leaf certificate signed by the CA
func generateLeafCert(hostname, certPath, keyPath, caCertPath, caKeyPath string) error {
	// Load CA certificate and key
	caCertPEM, err := os.ReadFile(caCertPath)
	if err != nil {
		return fmt.Errorf("read CA cert: %v", err)
	}
	caCertBlock, _ := pem.Decode(caCertPEM)
	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		return fmt.Errorf("parse CA cert: %v", err)
	}

	caKeyPEM, err := os.ReadFile(caKeyPath)
	if err != nil {
		return fmt.Errorf("read CA key: %v", err)
	}
	caKeyBlock, _ := pem.Decode(caKeyPEM)
	caKey, err := x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("parse CA key: %v", err)
	}

	// Generate leaf private key
	leafKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generate leaf key: %v", err)
	}

	// Create leaf certificate template
	leafTemplate := x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			CommonName: hostname,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour), // 1 year
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{hostname},
	}

	// Generate leaf certificate
	leafCert, err := x509.CreateCertificate(rand.Reader, &leafTemplate, caCert, &leafKey.PublicKey, caKey)
	if err != nil {
		return fmt.Errorf("create leaf cert: %v", err)
	}

	// Save leaf certificate
	certOut, err := os.Create(certPath)
	if err != nil {
		return fmt.Errorf("create leaf cert file: %v", err)
	}
	defer certOut.Close()
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: leafCert}); err != nil {
		return fmt.Errorf("encode leaf cert: %v", err)
	}

	// Save leaf private key
	keyOut, err := os.Create(keyPath)
	if err != nil {
		return fmt.Errorf("create leaf key file: %v", err)
	}
	defer keyOut.Close()
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(leafKey)}); err != nil {
		return fmt.Errorf("encode leaf key: %v", err)
	}

	return nil
}

// resolveRealIP uses DNS over HTTPS to get the real IP, bypassing /etc/hosts
func resolveRealIP(hostname string) (string, error) {
	// Use Cloudflare's DNS over HTTPS
	url := fmt.Sprintf("https://1.1.1.1/dns-query?name=%s&type=A", hostname)
	log.Printf("[DNS RESOLUTION] Resolving real IP for %s using DNS over HTTPS...", hostname)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/dns-json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Answer []struct {
			Data string `json:"data"`
			Type int    `json:"type"`
		} `json:"Answer"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	// Find the first A record (type 1)
	for _, answer := range result.Answer {
		if answer.Type == 1 && answer.Data != "" {
			log.Printf("[DNS RESOLUTION] Found IP %s for %s", answer.Data, hostname)
			return answer.Data, nil
		}
	}

	return "", fmt.Errorf("no A record found for %s", hostname)
}

// Start starts the OGINY TLS reverse proxy service
// If port is > 0, it uses that port, otherwise defaults to 443
// If loggingEnabled is true, all proxy traffic will be logged to files
func Start(port int, loggingEnabled bool) error {
	// Ignore cfgPath since we're inlining the config
	listenAddr := ":443"
	if port > 0 {
		listenAddr = fmt.Sprintf(":%d", port)
	}

	// Initialize logger only if logging is enabled
	if loggingEnabled {
		// Always use ~/oginy-logs as the log directory
		homeBase := os.Getenv("HOME")
		if homeBase == "" {
			homeBase = "/home"
		}
		logDir := filepath.Join(homeBase, "oginy-logs")

		if err := initLogger(logDir); err != nil {
			log.Printf("Warning: Failed to initialize file logger: %v", err)
			log.Printf("Continuing without file logging...")
		} else {
			// Ensure log file is closed on exit
			defer func() {
				if logFile != nil {
					proxyLogger.Printf("=== OGINY Proxy Log Ended ===")
					logFile.Close()
				}
			}()
			log.Printf("Proxy traffic logging enabled")
		}
	} else {
		log.Printf("Proxy traffic logging disabled")
	}

	// Get results-receiver hostname from env var or use default
	resultsReceiverHost := "results-receiver.actions.githubusercontent.com"
	if actionsURL := os.Getenv("ACTIONS_RESULTS_URL"); actionsURL != "" {
		if u, err := url.Parse(actionsURL); err == nil && u.Host != "" {
			resultsReceiverHost = u.Host
		}
	}

	// Set up certificate directory
	var certDir string
	if dir := os.Getenv("OGINY_CERT_DIR"); dir != "" {
		certDir = dir
	} else {
		// Use $HOME env var with fallback to /home
		homeBase := os.Getenv("HOME")
		if homeBase == "" {
			homeBase = "/home"
		}
		certDir = filepath.Join(homeBase, "runner", "certs")
	}

	// Create certificate directory if it doesn't exist
	if err := os.MkdirAll(certDir, 0755); err != nil {
		return fmt.Errorf("create cert dir: %v", err)
	}
	log.Printf("Using certificate directory: %s", certDir)

	// Cleanup any existing certificates
	os.RemoveAll(filepath.Join(certDir, "*.crt"))
	os.RemoveAll(filepath.Join(certDir, "*.key"))

	// Generate CA if it doesn't exist
	caCertPath := filepath.Join(certDir, "localCA.crt")
	caKeyPath := filepath.Join(certDir, "localCA.key")
	if _, err := os.Stat(caCertPath); os.IsNotExist(err) {
		log.Printf("Generating CA certificate...")
		if err := generateCA(caCertPath, caKeyPath); err != nil {
			return fmt.Errorf("generate CA: %v", err)
		}
		log.Printf("CA certificate generated at %s", caCertPath)
	}

	// Set environment variables for current process and children
	os.Setenv("NODE_OPTIONS", "--use-openssl-ca")
	os.Setenv("NODE_EXTRA_CA_CERTS", caCertPath)
	os.Setenv("SSL_CERT_FILE", caCertPath)

	// If running in GitHub Actions, write to GITHUB_ENV
	if githubEnv := os.Getenv("GITHUB_ENV"); githubEnv != "" {
		log.Printf("Detected GitHub Actions environment, writing to GITHUB_ENV")
		if err := appendToFile(githubEnv, fmt.Sprintf("NODE_OPTIONS=--use-openssl-ca\nNODE_EXTRA_CA_CERTS=%s\n", caCertPath)); err != nil {
			log.Printf("Warning: failed to write to GITHUB_ENV: %v", err)
		}
		// Write SSL_CERT_FILE in a separate call to avoid duplicate-check skipping
		if err := appendToFile(githubEnv, fmt.Sprintf("SSL_CERT_FILE=%s\n", caCertPath)); err != nil {
			log.Printf("Warning: failed to write SSL_CERT_FILE to GITHUB_ENV: %v", err)
		}
	}

	// Also write to /etc/environment if we have permissions (for system-wide)
	if err := appendToFile("/etc/environment", fmt.Sprintf("NODE_OPTIONS=\"--use-openssl-ca\"\nNODE_EXTRA_CA_CERTS=\"%s\"\n", caCertPath)); err != nil {
		// This is expected to fail if not running as root
		log.Printf("Note: Could not write to /etc/environment (need root): %v", err)
		log.Printf("To set system-wide, run as root or manually add to /etc/environment:")
		log.Printf("  NODE_OPTIONS=\"--use-openssl-ca\"")
		log.Printf("  NODE_EXTRA_CA_CERTS=\"%s\"", caCertPath)
		log.Printf("  SSL_CERT_FILE=\"%s\"", caCertPath)
	} else {
		log.Printf("Successfully updated /etc/environment")
		// Write SSL_CERT_FILE in a separate call to avoid duplicate-check skipping
		if err := appendToFile("/etc/environment", fmt.Sprintf("SSL_CERT_FILE=\"%s\"\n", caCertPath)); err != nil {
			log.Printf("Warning: failed to append SSL_CERT_FILE to /etc/environment: %v", err)
		}
	}

	// Generate certificates for each domain
	servers := []struct {
		serverName string
		certFile   string
		keyFile    string
		targetURL  string
	}{
		{
			serverName: "warpbuild.blob.core.windows.net",
			certFile:   filepath.Join(certDir, "warpbuild.crt"),
			keyFile:    filepath.Join(certDir, "warpbuild.key"),
			targetURL:  "http://127.0.0.1:50053",
		},
		{
			serverName: resultsReceiverHost,
			certFile:   filepath.Join(certDir, "results-receiver.crt"),
			keyFile:    filepath.Join(certDir, "results-receiver.key"),
			targetURL:  "http://127.0.0.1:50052",
		},
	}

	// Generate leaf certificates if they don't exist
	for _, s := range servers {
		if _, err := os.Stat(s.certFile); os.IsNotExist(err) {
			log.Printf("Generating certificate for %s...", s.serverName)
			if err := generateLeafCert(s.serverName, s.certFile, s.keyFile, caCertPath, caKeyPath); err != nil {
				return fmt.Errorf("generate cert for %s: %v", s.serverName, err)
			}
		}
	}

	// Shared, fast transport to backends.
	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           (&net.Dialer{Timeout: 10 * time.Second, KeepAlive: 60 * time.Second}).DialContext,
		ForceAttemptHTTP2:     false, // backend is http://
		MaxIdleConns:          1024,
		MaxIdleConnsPerHost:   512,
		MaxConnsPerHost:       0, // unlimited
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 0,
		DisableCompression:    true, // preserve encodings; avoid cpu
	}

	mp := &muxProxy{byHost: make(map[string]*proxyEntry)}

	for _, s := range servers {
		u, err := url.Parse(s.targetURL)
		if err != nil {
			return fmt.Errorf("target %s: %v", s.serverName, err)
		}
		c, err := tls.LoadX509KeyPair(s.certFile, s.keyFile)
		if err != nil {
			return fmt.Errorf("cert %s: %v", s.serverName, err)
		}

		// Create proxy with logging transport
		rp := httputil.NewSingleHostReverseProxy(u)
		loggingTransport := &loggingRoundTripper{
			transport: tr,
			name:      fmt.Sprintf("LOCAL:%s→%s", s.serverName, s.targetURL),
		}
		rp.Transport = loggingTransport
		origDirector := rp.Director
		rp.Director = func(r *http.Request) {
			origHost := r.Host
			origDirector(r)   // sets scheme/host/path to target
			r.Host = origHost // preserve original Host
			r.Header.Set("X-Forwarded-Proto", "https")
		}
		// No ModifyResponse / no custom flusher -> minimal overhead.

		entry := &proxyEntry{cert: &c, target: u, proxy: rp}

		// Special handling for results-receiver - enable path-based routing
		if s.serverName == resultsReceiverHost {
			entry.pathBasedRouting = true

			// Resolve the real IP to bypass /etc/hosts
			realIP, err := resolveRealIP(resultsReceiverHost)
			if err != nil {
				return fmt.Errorf("failed to resolve real IP for %s: %v", resultsReceiverHost, err)
			}
			log.Printf("Resolved real IP for %s: %s", resultsReceiverHost, realIP)

			// Create remote proxy for non-artifactcache requests using real IP
			remoteURL, err := url.Parse(fmt.Sprintf("https://%s", realIP))
			if err != nil {
				return fmt.Errorf("parse remote URL for %s: %v", resultsReceiverHost, err)
			}

			// Create custom transport that sets the proper SNI
			remoteTransport := &http.Transport{
				Proxy:       http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{Timeout: 10 * time.Second, KeepAlive: 60 * time.Second}).DialContext,
				TLSClientConfig: &tls.Config{
					ServerName: resultsReceiverHost, // Set SNI to the original hostname
					NextProtos: []string{"http/1.1"},
				},
				ForceAttemptHTTP2:     false,
				TLSNextProto:          map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
				MaxIdleConns:          1024,
				MaxIdleConnsPerHost:   512,
				MaxConnsPerHost:       0, // unlimited
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 0,
				DisableCompression:    true,
			}

			// Wrap remote transport with logging
			loggingRemoteTransport := &loggingRoundTripper{
				transport: remoteTransport,
				name:      fmt.Sprintf("REMOTE:%s→%s", resultsReceiverHost, realIP),
			}

			remoteProxy := httputil.NewSingleHostReverseProxy(remoteURL)
			remoteProxy.Transport = loggingRemoteTransport
			remoteProxy.Director = func(r *http.Request) {
				// Log the request details before modification
				origURL := r.URL.String()
				origHost := r.Host

				r.URL.Scheme = remoteURL.Scheme
				r.URL.Host = remoteURL.Host
				r.Host = resultsReceiverHost // Keep the original Host header
				// Don't set X-Forwarded-Proto for remote requests as they're already HTTPS

				// Log where the request is being sent
				log.Printf("[REMOTE PROXY] Forwarding request: %s %s (orig host: %s) → %s (IP: %s, Host header: %s)",
					r.Method, origURL, origHost, r.URL.String(), realIP, r.Host)
			}
			entry.remoteProxy = remoteProxy

			log.Printf("route: %s → %s (twirp only), other paths → %s (IP: %s)", s.serverName, s.targetURL, resultsReceiverHost, realIP)
		} else {
			log.Printf("route: %s → %s", s.serverName, s.targetURL)
		}

		mp.byHost[s.serverName] = entry
	}

	// TLS server config (minimal) - using TLS 1.2 as minimum and enabling HTTP/2
	tlsCfg := &tls.Config{
		MinVersion:     tls.VersionTLS12,
		GetCertificate: mp.GetCertificate, // SNI
		NextProtos:     []string{"http/1.1"},
	}

	srv := &http.Server{
		Addr:         listenAddr,
		Handler:      mp,
		TLSConfig:    tlsCfg,
		ReadTimeout:  0, // unlimited to avoid cutting long transfers
		WriteTimeout: 0,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("listening on %s", listenAddr)
	log.Printf("CA certificate location: %s", caCertPath)
	log.Printf("Set NODE_EXTRA_CA_CERTS=%s to trust the CA", caCertPath)
	return srv.ListenAndServeTLS("", "") // certs come from GetCertificate
}

// appendToFile appends content to a file if it doesn't already exist in the file
func appendToFile(filepath, content string) error {
	// Check if file exists and is writable
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read existing content to check for duplicates
	existing, err := os.ReadFile(filepath)
	if err == nil && strings.Contains(string(existing), content) {
		// Already configured with the same content, skip
		return nil
	}

	_, err = file.WriteString(content)
	return err
}
