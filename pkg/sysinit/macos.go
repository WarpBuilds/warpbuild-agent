package sysinit

import (
	"fmt"
	"os/exec"
)

// runMacOSDiagnostics runs macOS-specific system initialization diagnostics
func runMacOSDiagnostics() error {
	// whoami
	whoamiCmd := exec.Command("sh", "-c", "echo \"whoami: $(whoami)\"")
	if output, err := whoamiCmd.CombinedOutput(); err == nil {
		fmt.Print(string(output))
	} else {
		fmt.Printf("whoami command failed: %v\n", err)
	}

	// console user
	consoleUserCmd := exec.Command("sh", "-c", "echo \"console user: $(stat -f '%Su' /dev/console)\"")
	if output, err := consoleUserCmd.CombinedOutput(); err == nil {
		fmt.Print(string(output))
	} else {
		fmt.Printf("console user command failed: %v\n", err)
	}

	// pgrep Finder
	pgrepCmd := exec.Command("sh", "-c", "pgrep -lf Finder || true")
	if output, err := pgrepCmd.CombinedOutput(); err == nil {
		fmt.Print(string(output))
	} else {
		fmt.Printf("pgrep Finder command failed: %v\n", err)
	}

	// osascript startup disk
	osascriptCmd := exec.Command("osascript", "-e", "tell application \"Finder\" to get name of startup disk")
	if output, err := osascriptCmd.CombinedOutput(); err == nil {
		fmt.Print(string(output))
	} else {
		fmt.Printf("osascript startup disk command failed: %v\n", err)
	}

	return nil
}
