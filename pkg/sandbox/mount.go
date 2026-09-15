//go:build darwin

package sandbox

import (
	"os"
	"path/filepath"
	"syscall"
)

// isMountPoint reports whether path is the root of its own filesystem, by
// comparing its device with its parent's.
func isMountPoint(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return false, nil
	}

	parent := filepath.Dir(path)
	pfi, err := os.Stat(parent)
	if err != nil {
		return false, err
	}
	pst, ok := pfi.Sys().(*syscall.Stat_t)
	if !ok {
		return false, nil
	}

	return st.Dev != pst.Dev, nil
}
