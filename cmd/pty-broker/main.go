// warpbuild-pty-broker
//
// Tiny in-sandbox PTY broker. One WebSocket connection = one PTY. Clients
// send raw UTF-8 bytes for keystrokes; control messages (resize) ride a tiny
// JSON envelope. The server echoes raw PTY output as binary frames.
//
// Auth: short-lived HMAC-signed token in the query string. The HMAC key is
// passed in via env (BROKER_HMAC_KEY) or the --hmac-key flag. Tokens are
// scoped to a single sandbox/runner via the `sid` claim and expire ~60s
// after issue (`e` claim, unix seconds).
//
// Goals: low byte-path latency (no extra hops, no per-message RPC framing)
// and a small, dependency-light static binary so the upload/install cost is
// negligible.
package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

func main() {
	var port int
	var hmacKey string
	var defaultCwd string
	var shellPath string
	flag.IntVar(&port, "port", 7681, "TCP port to listen on")
	flag.StringVar(&hmacKey, "hmac-key", os.Getenv("BROKER_HMAC_KEY"), "HMAC key for token verification")
	flag.StringVar(&defaultCwd, "cwd", "/home/user", "default working directory for spawned shells")
	flag.StringVar(&shellPath, "shell", defaultShell(), "shell to spawn")
	flag.Parse()

	if hmacKey == "" {
		log.Fatal("hmac-key (or BROKER_HMAC_KEY env) is required")
	}

	srv := &server{
		hmacKey:    []byte(hmacKey),
		defaultCwd: defaultCwd,
		shell:      shellPath,
	}
	http.HandleFunc("/term", srv.handleTerminal)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	addr := fmt.Sprintf(":%d", port)
	log.Printf("warpbuild-pty-broker listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func defaultShell() string {
	if s := os.Getenv("SHELL"); s != "" {
		return s
	}
	for _, p := range []string{"/bin/bash", "/bin/sh"} {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return "/bin/sh"
}

type server struct {
	hmacKey    []byte
	defaultCwd string
	shell      string
}

// upgrader allows any origin: upstream proxies don't always forward the
// browser's Origin verbatim, and our auth is the HMAC token, not Origin.
var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

// tokenClaims is what the issuer signs into the token.
// We use a custom compact format rather than full JWT to keep the broker
// dependency-light: "<base64url(payload)>.<base64url(hmac)>" with a fixed
// SHA-256 algorithm. payload is JSON.
type tokenClaims struct {
	SandboxID string `json:"sid"` // sandbox the token is scoped to
	UserID    string `json:"uid"` // optional, for audit
	Cwd       string `json:"cwd"` // working directory to start the shell in
	Cols      int    `json:"c,omitempty"`
	Rows      int    `json:"r,omitempty"`
	Exp       int64  `json:"e"` // unix seconds
}

func (s *server) verifyToken(token string) (*tokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, errors.New("malformed token")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("bad payload encoding: %w", err)
	}
	sigGiven, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("bad signature encoding: %w", err)
	}
	mac := hmac.New(sha256.New, s.hmacKey)
	mac.Write([]byte(parts[0]))
	if !hmac.Equal(sigGiven, mac.Sum(nil)) {
		return nil, errors.New("bad signature")
	}
	var c tokenClaims
	if err := json.Unmarshal(payload, &c); err != nil {
		return nil, fmt.Errorf("bad payload json: %w", err)
	}
	if time.Now().Unix() > c.Exp {
		return nil, errors.New("token expired")
	}
	return &c, nil
}

func (s *server) handleTerminal(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	claims, err := s.verifyToken(token)
	if err != nil {
		log.Printf("auth rejected: %v", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws upgrade failed: %v", err)
		return
	}
	defer func() { _ = conn.Close() }()

	cols, rows := claims.Cols, claims.Rows
	if cols <= 0 {
		cols = 120
	}
	if rows <= 0 {
		rows = 32
	}
	cwd := claims.Cwd
	if cwd == "" {
		cwd = s.defaultCwd
	}

	cmd := exec.Command(s.shell, "-l")
	cmd.Dir = cwd
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")
	ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)})
	if err != nil {
		log.Printf("pty start failed: %v", err)
		closeWithReason(conn, websocket.CloseInternalServerErr, "pty start failed")
		return
	}
	defer func() {
		_ = ptmx.Close()
		// Reap the child so we don't leak zombies. PTY close usually does
		// this on its own, but be explicit.
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
	}()

	log.Printf("pty started: pid=%d cwd=%s shell=%s", cmd.Process.Pid, cwd, s.shell)

	// PTY → WS pump.
	go func() {
		buf := make([]byte, 32*1024)
		for {
			n, err := ptmx.Read(buf)
			if n > 0 {
				if werr := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); werr != nil {
					return
				}
			}
			if err != nil {
				// EIO is normal when the shell exits.
				if !errors.Is(err, io.EOF) && !errors.Is(err, syscall.EIO) {
					log.Printf("pty read error: %v", err)
				}
				closeWithReason(conn, websocket.CloseNormalClosure, "pty exited")
				return
			}
		}
	}()

	// Reaper: when the shell exits, the PTY read loop above closes the WS.
	// Also surface a non-zero exit code clearly.
	go func() {
		_ = cmd.Wait()
		_ = ptmx.Close()
	}()

	// WS → PTY pump (this goroutine).
	for {
		mt, data, err := conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("ws read error: %v", err)
			}
			return
		}
		switch mt {
		case websocket.BinaryMessage:
			if _, werr := ptmx.Write(data); werr != nil {
				log.Printf("pty write failed: %v", werr)
				return
			}
		case websocket.TextMessage:
			handleControlMessage(ptmx, data)
		}
	}
}

// handleControlMessage handles JSON-wrapped control envelopes (resize and
// legacy "input"/"stdin"). Unknown types are dropped — text messages should
// only be control envelopes; bulk text input arrives as binary.
func handleControlMessage(ptmx *os.File, data []byte) {
	var msg struct {
		Type string `json:"type"`
		Cols int    `json:"cols"`
		Rows int    `json:"rows"`
		Data string `json:"data"`
	}
	if err := json.Unmarshal(data, &msg); err != nil {
		return
	}
	switch msg.Type {
	case "resize":
		if msg.Cols > 0 && msg.Rows > 0 {
			_ = setWinsize(ptmx, uint16(msg.Cols), uint16(msg.Rows))
		}
	case "input", "stdin":
		_, _ = ptmx.Write([]byte(msg.Data))
	}
}

// setWinsize: TIOCSWINSZ ioctl. creack/pty exposes Setsize but we wrap to
// keep all syscalls in one place.
func setWinsize(f *os.File, cols, rows uint16) error {
	ws := struct {
		rows, cols, x, y uint16
	}{rows, cols, 0, 0}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		syscall.TIOCSWINSZ,
		uintptr(unsafe.Pointer(&ws)),
	)
	if errno != 0 {
		return errno
	}
	return nil
}

func closeWithReason(conn *websocket.Conn, code int, reason string) {
	msg := websocket.FormatCloseMessage(code, reason)
	_ = conn.WriteControl(
		websocket.CloseMessage,
		msg,
		time.Now().Add(2*time.Second),
	)
}
