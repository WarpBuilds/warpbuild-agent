package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/warpbuilds/warpbuild-agent/pkg/app"
)

type flagsStruct struct {
	stdoutFile        string
	stderrFile        string
	settingsFile      string
	launchTelemetry   bool
	launchProxyServer bool
}

var flags flagsStruct

var rootCmd = &cobra.Command{
	Use:   "agentd",
	Short: "Manages runner lifecycle",
	Long: `Manages runner lifecycle and synchronizes runner state with the WarpBuild.
	This is run as a daemon on the runner host.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		err := app.NewApp(cmd.Context(), &app.ApplicationOptions{
			SettingsFile:      flags.settingsFile,
			StdoutFile:        flags.stdoutFile,
			StderrFile:        flags.stderrFile,
			LaunchTelemetry:   flags.launchTelemetry,
			LaunchProxyServer: flags.launchProxyServer,
		})
		if err != nil {
			return err
		}

		return nil

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func ExecuteWithErr() error {
	return rootCmd.Execute()
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.agentd.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVar(&flags.stdoutFile, "stdout", "", "stdout file")
	rootCmd.PersistentFlags().StringVar(&flags.stderrFile, "stderr", "", "stderr file")
	rootCmd.PersistentFlags().StringVar(&flags.settingsFile, "settings", "", "settings file")
	rootCmd.PersistentFlags().BoolVar(&flags.launchTelemetry, "launch-telemetry", false, "launch telemetry")
	rootCmd.PersistentFlags().BoolVar(&flags.launchProxyServer, "launch-proxy-server", false, "launch proxy server")
}
