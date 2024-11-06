//go:build windows

package cmd

import (
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

type flagsStruct struct {
	stdoutFile        string
	stderrFile        string
	agentdServiceName string
	restartInterval   time.Duration
}

var flags flagsStruct

var rootCmd = &cobra.Command{
	Use:   "agentd-restarter",
	Short: "Restarts the agentd service",
	Long:  `Restarts the agentd service`,
	RunE: func(cmd *cobra.Command, args []string) error {

		for {
			time.Sleep(flags.restartInterval)

			serviceName := flags.agentdServiceName

			// Connect to the service manager
			m, err := mgr.Connect()
			if err != nil {
				log.Fatalf("Failed to connect to service manager: %v", err)
			}
			defer m.Disconnect()

			// Open the specified service
			service, err := m.OpenService(serviceName)
			if err != nil {
				log.Fatalf("Could not access service %s: %v", serviceName, err)
			}
			defer service.Close()

			// Query the current status of the service
			status, err := service.Query()
			if err != nil {
				log.Fatalf("Could not query service %s: %v", serviceName, err)
			}

			if status.State == svc.Stopped {
				// restart the service
				err := service.Start()
				if err != nil {
					log.Fatalf("Could not start service %s: %v", serviceName, err)
				}
			}

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
	rootCmd.PersistentFlags().StringVar(&flags.agentdServiceName, "agentd-service-name", "warpbuild-agentd", "agentd service name")
	rootCmd.PersistentFlags().DurationVar(&flags.restartInterval, "restart-interval", 1*time.Minute, "restart interval")
}
