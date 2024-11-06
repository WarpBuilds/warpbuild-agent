//go:build windows

package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/warpbuilds/warpbuild-agent/pkg/log"
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

		lm, err := log.Init(&log.InitOptions{
			StdoutFile: flags.stdoutFile,
			StderrFile: flags.stderrFile,
		})
		if err != nil {
			return err
		}
		defer lm.Sync()

		for {
			log.Logger().Infof("sleeping for %v", flags.restartInterval)

			time.Sleep(flags.restartInterval)

			serviceName := flags.agentdServiceName

			// Connect to the service manager
			m, err := mgr.Connect()
			if err != nil {
				log.Logger().Errorf("Failed to connect to service manager: %v", err)
				continue
			}
			defer m.Disconnect()

			// Open the specified service
			service, err := m.OpenService(serviceName)
			if err != nil {
				log.Logger().Errorf("Could not access service %s: %v", serviceName, err)
				continue
			}
			defer service.Close()

			// Query the current status of the service
			status, err := service.Query()
			if err != nil {
				log.Logger().Errorf("Could not query service %s: %v", serviceName, err)
				continue
			}

			if status.State == svc.Stopped {
				log.Logger().Infof("service %s is stopped, restarting", serviceName)
				// restart the service
				err := service.Start()
				if err != nil {
					log.Logger().Errorf("Could not start service %s: %v", serviceName, err)
					continue
				}

				log.Logger().Infof("service %s restarted", serviceName)
			}

			if status.State == svc.StopPending {
				// kill the service process
				err := killServiceProcess(service)
				if err != nil {
					log.Logger().Errorf("Could not kill service %s: %v", serviceName, err)
					continue
				}
			}

			if status.State == svc.Running {
				log.Logger().Infof("service %s is running. Sleeping for another %v", serviceName, flags.restartInterval)
				time.Sleep(flags.restartInterval)
				continue
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

func killServiceProcess(service *mgr.Service) error {
	// Get the process ID of the service
	status, err := service.Query()
	if err != nil {
		log.Logger().Errorf("Could not query service %s: %v", service.Name, err)
		return err
	}

	log.Logger().Infof("service %s has process id %d", service.Name, status.ProcessId)

	if status.ProcessId == 0 {
		log.Logger().Infof("service %s has no process id", service.Name)
		return nil // Process is already gone
	}

	log.Logger().Infof("finding process %d", status.ProcessId)

	// Open the process
	process, err := os.FindProcess(int(status.ProcessId))
	if err != nil {
		log.Logger().Errorf("Could not find process %d: %v", status.ProcessId, err)
		return err
	}

	log.Logger().Infof("killing process %d", status.ProcessId)

	// Kill the process
	err = process.Kill()
	if err != nil {
		log.Logger().Errorf("Could not kill process %d: %v", status.ProcessId, err)
		return err
	}

	log.Logger().Infof("killed process %d", status.ProcessId)

	return nil
}
