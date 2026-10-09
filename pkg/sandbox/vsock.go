//go:build darwin

package sandbox

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
	"unsafe"
)

const (
	afVSOCK = 40

	vmaddrCIDAny = ^uint32(0)
)

type sockaddrVM struct {
	Len       uint8
	Family    uint8
	Reserved1 uint16
	Port      uint32
	CID       uint32
}

type vsockAddr struct{ port uint32 }

func (a *vsockAddr) Network() string { return "vsock" }
func (a *vsockAddr) String() string  { return fmt.Sprintf("vsock:%d", a.port) }

type vsockListener struct {
	fd   int
	addr *vsockAddr
}

func (l *vsockListener) Accept() (net.Conn, error) {
	for {
		r, _, errno := syscall.Syscall(syscall.SYS_ACCEPT, uintptr(l.fd), 0, 0)
		if errno != 0 {
			if errno == syscall.EINTR || errno == syscall.ECONNABORTED {
				continue
			}

			return nil, &net.OpError{Op: "accept", Net: "vsock", Addr: l.addr, Err: errno}
		}
		nfd := int(r)
		syscall.CloseOnExec(nfd)

		return &vsockConn{f: os.NewFile(uintptr(nfd), l.addr.String()), addr: l.addr}, nil
	}
}

func (l *vsockListener) Close() error   { return syscall.Close(l.fd) }
func (l *vsockListener) Addr() net.Addr { return l.addr }

type vsockConn struct {
	f    *os.File
	addr *vsockAddr
}

func (c *vsockConn) Read(b []byte) (int, error)         { return c.f.Read(b) }
func (c *vsockConn) Write(b []byte) (int, error)        { return c.f.Write(b) }
func (c *vsockConn) Close() error                       { return c.f.Close() }
func (c *vsockConn) LocalAddr() net.Addr                { return c.addr }
func (c *vsockConn) RemoteAddr() net.Addr               { return c.addr }
func (c *vsockConn) SetDeadline(t time.Time) error      { return c.f.SetDeadline(t) }
func (c *vsockConn) SetReadDeadline(t time.Time) error  { return c.f.SetReadDeadline(t) }
func (c *vsockConn) SetWriteDeadline(t time.Time) error { return c.f.SetWriteDeadline(t) }

func listenVsock(port uint32) (net.Listener, error) {
	fd, err := syscall.Socket(afVSOCK, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, fmt.Errorf("vsock socket: %w", err)
	}
	syscall.CloseOnExec(fd)

	sa := sockaddrVM{
		Len:    uint8(unsafe.Sizeof(sockaddrVM{})),
		Family: afVSOCK,
		Port:   port,
		CID:    vmaddrCIDAny,
	}
	if _, _, errno := syscall.Syscall(
		syscall.SYS_BIND,
		uintptr(fd),
		uintptr(unsafe.Pointer(&sa)),
		uintptr(unsafe.Sizeof(sa)),
	); errno != 0 {
		syscall.Close(fd)

		return nil, fmt.Errorf("vsock bind port %d: %w", port, errno)
	}

	if err := syscall.Listen(fd, 64); err != nil {
		syscall.Close(fd)

		return nil, fmt.Errorf("vsock listen port %d: %w", port, err)
	}

	return &vsockListener{fd: fd, addr: &vsockAddr{port: port}}, nil
}
