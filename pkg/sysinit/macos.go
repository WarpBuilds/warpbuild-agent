package sysinit

import (
	"os/exec"
	"strings"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

// runMacOSDiagnostics runs macOS-specific system initialization diagnostics
func runMacOSDiagnostics() error {
	// whoami
	whoamiCmd := exec.Command("sh", "-c", "whoami")
	if output, err := whoamiCmd.CombinedOutput(); err == nil {
		log.Logger().Infof("whoami: %s", strings.TrimSpace(string(output)))
	} else {
		log.Logger().Errorf("whoami command failed: %v", err)
	}

	// console user
	consoleUserCmd := exec.Command("sh", "-c", "stat -f '%Su' /dev/console")
	if output, err := consoleUserCmd.CombinedOutput(); err == nil {
		log.Logger().Infof("console user: %s", strings.TrimSpace(string(output)))
	} else {
		log.Logger().Errorf("console user command failed: %v", err)
	}

	// pgrep Finder
	pgrepCmd := exec.Command("sh", "-c", "pgrep -lf Finder || true")
	if output, err := pgrepCmd.CombinedOutput(); err == nil {
		if len(output) > 0 {
			log.Logger().Infof("Finder processes: %s", strings.TrimSpace(string(output)))
		} else {
			log.Logger().Infof("Finder processes: none found")
		}
	} else {
		log.Logger().Errorf("pgrep Finder command failed: %v", err)
	}

	// osascript startup disk
	// This is what triggers the popup for 'Automations' on warpbuild-agent
	osascriptCmd := exec.Command("osascript", "-e", "tell application \"Finder\" to get name of startup disk")
	if output, err := osascriptCmd.CombinedOutput(); err == nil {
		log.Logger().Infof("startup disk: %s", strings.TrimSpace(string(output)))
	} else {
		log.Logger().Errorf("osascript startup disk command failed: %v", err)
	}

	return nil
}
