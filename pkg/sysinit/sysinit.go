package sysinit

import (
	"fmt"
	"runtime"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

// SysInit runs system initialization diagnostics based on the current OS and architecture.
// It automatically detects the platform and runs the appropriate initialization routine.
// If the platform is not supported, it silently skips without error.
func SysInit() error {
	// Only run on macOS (darwin)
	if runtime.GOOS == "darwin" {
		log.Logger().Infof("=== Running macOS System Initialization Diagnostics ===")
		if err := runMacOSDiagnostics(); err != nil {
			return fmt.Errorf("macOS diagnostics failed: %w", err)
		}
		log.Logger().Infof("=== System Initialization Diagnostics Complete ===")
	}

	return nil
}
