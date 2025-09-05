package oginy

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

type Config struct {
	Servers       []ServerConfig `json:"servers"`
	ListenAddr    string         `json:"listenAddr"`
	EnableHTTP2   bool           `json:"enableHTTP2"`
	TLSMinVersion string         `json:"tlsMinVersion"`
}
type ServerConfig struct {
	ServerName     string `json:"serverName"`
	CertFile       string `json:"certFile"`
	KeyFile        string `json:"keyFile"`
	TargetURL      string `json:"targetURL"`
	TimeoutSeconds int    `json:"timeoutSeconds"`
}

type proxyEntry struct {
	cert   *tls.Certificate
	target *url.URL
	proxy  *httputil.ReverseProxy
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
		pe.proxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "no backend for host", http.StatusBadGateway)
}

func loadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	if c.ListenAddr == "" {
		c.ListenAddr = ":50052"
	}
	if c.TLSMinVersion == "" {
		c.TLSMinVersion = "1.2"
	}
	return &c, nil
}

func main() {
	cfgPath := flag.String("config", "config.json", "path to config")
	flag.Parse()

	cfg, err := loadConfig(*cfgPath)
	if err != nil {
		log.Fatalf("config: %v", err)
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

	for _, s := range cfg.Servers {
		u, err := url.Parse(s.TargetURL)
		if err != nil {
			log.Fatalf("target %s: %v", s.ServerName, err)
		}
		c, err := tls.LoadX509KeyPair(s.CertFile, s.KeyFile)
		if err != nil {
			log.Fatalf("cert %s: %v", s.ServerName, err)
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

		mp.byHost[s.ServerName] = &proxyEntry{cert: &c, target: u, proxy: rp}
		log.Printf("route: %s → %s", s.ServerName, s.TargetURL)
	}

	// TLS server config (minimal).
	var minVer uint16 = tls.VersionTLS12
	switch cfg.TLSMinVersion {
	case "1.0":
		minVer = tls.VersionTLS10
	case "1.1":
		minVer = tls.VersionTLS11
	case "1.3":
		minVer = tls.VersionTLS13
	}
	tlsCfg := &tls.Config{
		MinVersion:     minVer,
		GetCertificate: mp.GetCertificate, // SNI
	}
	if cfg.EnableHTTP2 {
		// Let Go’s default enable h2 via ALPN; no extra tuning to keep it simple.
		tlsCfg.NextProtos = []string{"h2", "http/1.1"}
	} else {
		tlsCfg.NextProtos = []string{"http/1.1"}
	}

	srv := &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      mp,
		TLSConfig:    tlsCfg,
		ReadTimeout:  0, // unlimited to avoid cutting long transfers
		WriteTimeout: 0,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("listening on %s", cfg.ListenAddr)
	log.Fatal(srv.ListenAndServeTLS("", "")) // certs come from GetCertificate
}
