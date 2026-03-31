package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var sysInitCmd = &cobra.Command{
	Use:    "sys-init",
	Short:  "System initialization commands",
	Hidden: true,
}

var sysInitMacosCmd = &cobra.Command{
	Use:   "macos",
	Short: "macOS system initialization",
	RunE: func(cmd *cobra.Command, args []string) error {
		// whoami
		whoamiCmd := exec.Command("sh", "-c", "echo \"whoami: $(whoami)\"")
		if output, err := whoamiCmd.CombinedOutput(); err == nil {
			fmt.Print(string(output))
		}

		// console user
		consoleUserCmd := exec.Command("sh", "-c", "echo \"console user: $(stat -f '%Su' /dev/console)\"")
		if output, err := consoleUserCmd.CombinedOutput(); err == nil {
			fmt.Print(string(output))
		}

		// pgrep Finder
		pgrepCmd := exec.Command("sh", "-c", "pgrep -lf Finder || true")
		if output, err := pgrepCmd.CombinedOutput(); err == nil {
			fmt.Print(string(output))
		}

		// osascript startup disk
		osascriptCmd := exec.Command("osascript", "-e", "tell application \"Finder\" to get name of startup disk")
		if output, err := osascriptCmd.CombinedOutput(); err == nil {
			fmt.Print(string(output))
		}

		return nil
	},
}

func init() {
	sysInitCmd.AddCommand(sysInitMacosCmd)
	rootCmd.AddCommand(sysInitCmd)
}
