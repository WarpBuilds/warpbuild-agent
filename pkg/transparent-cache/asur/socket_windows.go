//go:build windows
// +build windows

package asur

import (
	"syscall"
)

// setSocketOptions sets platform-specific socket options for Windows systems
func setSocketOptions(fd uintptr) error {
	// Enable TCP_NODELAY for low latency
	if err := syscall.SetsockoptInt(syscall.Handle(fd), syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 1); err != nil {
		return err
	}

	// Increase socket buffer sizes for better throughput
	if err := syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_SNDBUF, 4*1024*1024); err != nil {
		return err
	}

	if err := syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUF, 4*1024*1024); err != nil {
		return err
	}

	return nil
}
