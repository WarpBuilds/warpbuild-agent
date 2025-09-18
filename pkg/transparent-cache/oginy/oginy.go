package oginy

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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
			// Only forward /_api/artifactcache requests to local service
			if strings.HasPrefix(r.URL.Path, "/_api/artifactcache") {
				pe.proxy.ServeHTTP(w, r)
			} else if pe.remoteProxy != nil {
				// Forward all other requests to the actual domain
				pe.remoteProxy.ServeHTTP(w, r)
			} else {
				http.Error(w, "no remote proxy configured", http.StatusBadGateway)
			}
		} else {
			// Standard routing - forward all requests to local service
			pe.proxy.ServeHTTP(w, r)
		}
		return
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

// Start starts the OGINY TLS reverse proxy service
// If port is > 0, it uses that port, otherwise defaults to 443
func Start(port int) error {
	// Ignore cfgPath since we're inlining the config
	listenAddr := ":443"
	if port > 0 {
		listenAddr = fmt.Sprintf(":%d", port)
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

	// If running in GitHub Actions, write to GITHUB_ENV
	if githubEnv := os.Getenv("GITHUB_ENV"); githubEnv != "" {
		log.Printf("Detected GitHub Actions environment, writing to GITHUB_ENV")
		if err := appendToFile(githubEnv, fmt.Sprintf("NODE_OPTIONS=--use-openssl-ca\nNODE_EXTRA_CA_CERTS=%s\n", caCertPath)); err != nil {
			log.Printf("Warning: failed to write to GITHUB_ENV: %v", err)
		}
	}

	// Also write to /etc/environment if we have permissions (for system-wide)
	if err := appendToFile("/etc/environment", fmt.Sprintf("NODE_OPTIONS=\"--use-openssl-ca\"\nNODE_EXTRA_CA_CERTS=\"%s\"\n", caCertPath)); err != nil {
		// This is expected to fail if not running as root
		log.Printf("Note: Could not write to /etc/environment (need root): %v", err)
		log.Printf("To set system-wide, run as root or manually add to /etc/environment:")
		log.Printf("  NODE_OPTIONS=\"--use-openssl-ca\"")
		log.Printf("  NODE_EXTRA_CA_CERTS=\"%s\"", caCertPath)
	} else {
		log.Printf("Successfully updated /etc/environment")
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

	// Transport for remote HTTPS connections
	remoteTr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           (&net.Dialer{Timeout: 10 * time.Second, KeepAlive: 60 * time.Second}).DialContext,
		ForceAttemptHTTP2:     true,
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

		rp := httputil.NewSingleHostReverseProxy(u)
		rp.Transport = tr
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

			// Create remote proxy for non-artifactcache requests
			remoteURL, err := url.Parse("https://" + resultsReceiverHost)
			if err != nil {
				return fmt.Errorf("parse remote URL for %s: %v", resultsReceiverHost, err)
			}

			remoteProxy := httputil.NewSingleHostReverseProxy(remoteURL)
			remoteProxy.Transport = remoteTr
			remoteProxy.Director = func(r *http.Request) {
				r.URL.Scheme = remoteURL.Scheme
				r.URL.Host = remoteURL.Host
				r.Host = resultsReceiverHost
				// Don't set X-Forwarded-Proto for remote requests as they're already HTTPS
			}
			entry.remoteProxy = remoteProxy

			log.Printf("route: %s → %s (artifactcache only), other paths → %s", s.serverName, s.targetURL, remoteURL.String())
		} else {
			log.Printf("route: %s → %s", s.serverName, s.targetURL)
		}

		mp.byHost[s.serverName] = entry
	}

	// TLS server config (minimal) - using TLS 1.2 as minimum and enabling HTTP/2
	tlsCfg := &tls.Config{
		MinVersion:     tls.VersionTLS12,
		GetCertificate: mp.GetCertificate,          // SNI
		NextProtos:     []string{"h2", "http/1.1"}, // Enable HTTP/2
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
	if err == nil && strings.Contains(string(existing), "NODE_EXTRA_CA_CERTS") {
		// Already configured, skip
		return nil
	}

	_, err = file.WriteString(content)
	return err
}
