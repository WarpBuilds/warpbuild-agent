package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

const (
	Interval = 1 * time.Second
)

type StartAgentOptions struct {
	Manager *ManagerOptions `json:"manager"`
}

type IAgent interface {
	StartAgent(ctx context.Context, opts *StartAgentOptions) error
}

type AgentOptions struct {
	// ID is the warpbuild assigned id.
	ID               string `json:"id"`
	PollingSecret    string `json:"polling_secret"`
	HostURL          string `json:"host_url"`
	ExitFileLocation string `json:"exit_file_location"`
	// WindowsOptions are options for the windows agent.
	WindowsOptions *WindowsOptions `json:"windows_options"`
}

type WindowsOptions struct {
	ServiceName string `json:"service_name"`
}

func NewAgent(opts *AgentOptions) (IAgent, error) {
	cfg := warpbuild.NewConfiguration()

	if opts.HostURL == "" {
		return nil, fmt.Errorf("host url is required")
	}

	cfg.Servers[0].URL = opts.HostURL

	wb := warpbuild.NewAPIClient(cfg)
	return &agentImpl{
		client:           wb,
		id:               opts.ID,
		pollingSecret:    opts.PollingSecret,
		hostURL:          opts.HostURL,
		exitFileLocation: opts.ExitFileLocation,
		opts:             opts,
	}, nil
}

type agentImpl struct {
	client               *warpbuild.APIClient
	id                   string
	pollingSecret        string
	hostURL              string
	exitFileLocation     string
	opts                 *AgentOptions
	telemetryKilled      bool
	lastTelemetryEnabled *bool
}

type ExitFile struct {
	ExitCode     int                `json:"exit_code"`
	MachineState RunnerMachineState `json:"machine_state"`
}

type RunnerMachineState string

const (
	RunnerMachineStateDirty RunnerMachineState = "dirty"
)

func (a *agentImpl) StartAgent(ctx context.Context, opts *StartAgentOptions) error {

	if a.exitFileLocation == "" {
		return fmt.Errorf("exit file location is required")
	}

	ticker := time.NewTicker(Interval)
	for {
		select {
		case <-ticker.C:

			if err := a.verifyExitFile(); err != nil {
				log.Logger().Errorf("exit file verification failed: %v", err)
				log.Logger().Infof("Runner will not be started and polling will not happen.")
				continue
			}

			log.Logger().Infof("host url: %s", a.hostURL)
			log.Logger().Infof("checking for runner instance allocation details for %s", a.id)
			log.Logger().Infof("polling secret: %s", a.pollingSecret)

			allocationDetails, resp, err := a.client.V1RunnerInstanceAPI.
				GetRunnerInstanceAllocationDetails(ctx, a.id).
				XPOLLINGSECRET(a.pollingSecret).
				Execute()
			if err != nil {
				// get url from resp
				log.Logger().Errorf("failed to get runner instance allocation details: %v", err)
				log.Logger().Errorf("Response: %+v", resp)
				continue
			}

			if allocationDetails == nil {
				log.Logger().Infof("No runner instance allocation details found. Retrying in %s", Interval)
				continue
			}

			// Check telemetry status and kill telemetry agent if disabled
			if allocationDetails.HasTelemetryEnabled() {
				telemetryEnabled := allocationDetails.GetTelemetryEnabled()

				// Check if telemetry status has changed or if we haven't checked before
				if a.lastTelemetryEnabled == nil || *a.lastTelemetryEnabled != telemetryEnabled {
					log.Logger().Infof("Telemetry enabled status: %v", telemetryEnabled)
					a.lastTelemetryEnabled = &telemetryEnabled

					if !telemetryEnabled && !a.telemetryKilled {
						log.Logger().Infof("Telemetry is disabled. Killing telemetry agent asynchronously...")
						a.telemetryKilled = true // Set flag immediately to prevent duplicate attempts
						go func() {
							if err := a.killTelemetryProcess(); err != nil {
								log.Logger().Errorf("Failed to kill telemetry process: %v", err)
							} else {
								log.Logger().Infof("Telemetry agent successfully killed")
							}
						}()
					}
				}
			}

			// TODO: verify the correct status
			if *allocationDetails.Status == "assigned" {

				log.Logger().Infof("Setting additonal environment variables")
				for key, val := range *allocationDetails.GhRunnerApplicationDetails.Variables {
					os.Setenv(key, val)
				}

				if opts.Manager.Provider == ProviderGithubCRI {
					for key, val := range *allocationDetails.GhRunnerApplicationDetails.Variables {
						opts.Manager.GithubCRI.CMDOptions.Envs = append(opts.Manager.GithubCRI.CMDOptions.Envs, EnvironmentVariable{
							Key:   key,
							Value: val,
						})
					}
				}

				log.Logger().Infof("Starting runner")
				m := NewManager(opts.Manager)
				startRunnerOutput, err := m.StartRunner(ctx, &StartRunnerOptions{
					JitToken:     *allocationDetails.GhRunnerApplicationDetails.Jit,
					AgentOptions: a.opts,
				})
				if err != nil {
					log.Logger().Errorf("failed to start runner: %v", err)
					return err
				}

				if startRunnerOutput.RunCompletedSuccessfully {
					err := a.writeExitFile(ctx, startRunnerOutput)
					if err != nil {
						log.Logger().Errorf("failed to write exit file: %v", err)
						return err
					}
				}

			} else {
				log.Logger().Infof("runner instance allocation details status: %s", *allocationDetails.Status)
				log.Logger().Infof("Retrying in %s", Interval)
			}

		case <-ctx.Done():
			log.Logger().Infof("Context cancelled. Agent is exiting...")
			return nil
		}
	}

}

func (a *agentImpl) writeExitFile(ctx context.Context, opts *StartRunnerOutput) error {
	log.Logger().Infof("Runner completed successfully. Marking vm as dirty")

	ef := &ExitFile{
		ExitCode:     0,
		MachineState: RunnerMachineStateDirty,
	}

	data, err := json.Marshal(ef)
	if err != nil {
		log.Logger().Errorf("failed to marshal exit file: %v", err)
		return err
	}

	f, err := os.Create(a.exitFileLocation)
	if err != nil && !os.IsExist(err) {
		log.Logger().Errorf("failed to create exit file: %v", err)
		return err
	}

	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		log.Logger().Errorf("failed to write exit file: %v", err)
		return err
	}

	log.Logger().Infof("Exiting...")
	return nil

}

func (a *agentImpl) verifyExitFile() error {
	// read the exit file
	log.Logger().Infof("Verifying exit file at %s", a.exitFileLocation)

	f, err := os.Open(a.exitFileLocation)
	if err != nil && os.IsNotExist(err) {
		log.Logger().Infof("exit file does not exist. VM is clean. Continuing with agent startup...")
		// the exit file does not exist which means the vm is clean
		return nil
	} else if err != nil {
		log.Logger().Errorf("failed to open exit file: %v", err)
		return err
	}
	defer f.Close()

	// read the exit file
	var ef ExitFile
	err = json.NewDecoder(f).Decode(&ef)
	if err != nil {
		log.Logger().Errorf("failed to decode exit file: %v", err)
		return err
	}

	if ef.MachineState == RunnerMachineStateDirty {
		log.Logger().Errorf("exit file exists and machine state is dirty. Exiting...")
		return fmt.Errorf("exit file exists and machine state is dirty. VMs must be clean before starting the agent")
	} else {
		log.Logger().Infof("exit file exists and machine state is '%s'. Continuing with agent startup...", ef.MachineState)
	}

	return nil

}

// killTelemetryProcess stops the telemetry service using system service managers
// This is better than killing the process directly because:
// 1. The telemetry agent runs as a system service with auto-restart enabled
// 2. Killing the process would cause the service manager to restart it
// 3. Stopping the service is cleaner and respects the service lifecycle
func (a *agentImpl) killTelemetryProcess() error {
	log.Logger().Infof("Attempting to stop telemetry service...")

	switch runtime.GOOS {
	case "linux":
		return a.stopTelemetryServiceLinux()
	case "darwin":
		return a.stopTelemetryServiceDarwin()
	case "windows":
		return a.stopTelemetryServiceWindows()
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

// stopTelemetryServiceLinux stops the telemetry service on Linux using systemctl
func (a *agentImpl) stopTelemetryServiceLinux() error {
	serviceName := "warpbuild-telemetryd"

	// Check if service exists and is running
	checkCmd := exec.Command("systemctl", "is-active", serviceName)
	output, err := checkCmd.Output()
	status := strings.TrimSpace(string(output))

	if err != nil || status != "active" {
		log.Logger().Infof("Telemetry service is not running (status: %s)", status)
		return nil
	}

	// Stop the service
	log.Logger().Infof("Stopping telemetry service: %s", serviceName)
	stopCmd := exec.Command("systemctl", "stop", serviceName)
	if err := stopCmd.Run(); err != nil {
		log.Logger().Errorf("Failed to stop telemetry service: %v", err)
		return err
	}

	// Optionally disable the service to prevent it from starting on next boot
	disableCmd := exec.Command("systemctl", "disable", serviceName)
	if err := disableCmd.Run(); err != nil {
		log.Logger().Warnf("Failed to disable telemetry service (non-critical): %v", err)
	}

	log.Logger().Infof("Successfully stopped telemetry service: %s", serviceName)
	return nil
}

// stopTelemetryServiceDarwin stops the telemetry service on macOS using launchctl
func (a *agentImpl) stopTelemetryServiceDarwin() error {
	serviceName := "warpbuild-telemetryd-launcher"

	// Try to unload the service (this stops it and prevents it from restarting)
	log.Logger().Infof("Stopping telemetry service: %s", serviceName)

	// First, try to find the service using launchctl list
	listCmd := exec.Command("launchctl", "list", serviceName)
	if err := listCmd.Run(); err != nil {
		log.Logger().Infof("Telemetry service is not loaded: %s", serviceName)
		return nil
	}

	// Stop the service using launchctl bootout or stop
	// Try modern approach first (macOS 10.11+)
	stopCmd := exec.Command("launchctl", "stop", serviceName)
	if err := stopCmd.Run(); err != nil {
		log.Logger().Warnf("Failed to stop telemetry service with 'stop' command: %v", err)
	}

	// Also try to bootout (remove from launchd)
	bootoutCmd := exec.Command("launchctl", "bootout", "gui/"+os.Getenv("UID"), serviceName)
	if err := bootoutCmd.Run(); err != nil {
		log.Logger().Warnf("Failed to bootout telemetry service (trying unload): %v", err)

		// Fallback to older unload command
		unloadCmd := exec.Command("launchctl", "unload", "-w", fmt.Sprintf("/Library/LaunchDaemons/%s.plist", serviceName))
		if err := unloadCmd.Run(); err != nil {
			log.Logger().Errorf("Failed to unload telemetry service: %v", err)
			return err
		}
	}

	log.Logger().Infof("Successfully stopped telemetry service: %s", serviceName)
	return nil
}

// stopTelemetryServiceWindows stops the telemetry service on Windows
func (a *agentImpl) stopTelemetryServiceWindows() error {
	serviceName := "warpbuild-telemetryd"

	// Check if service exists
	checkCmd := exec.Command("sc", "query", serviceName)
	output, err := checkCmd.Output()

	if err != nil {
		log.Logger().Infof("Telemetry service does not exist or is not accessible: %s", serviceName)
		return nil
	}

	// Check if service is running
	if !strings.Contains(string(output), "RUNNING") {
		log.Logger().Infof("Telemetry service is not running")
		return nil
	}

	// Stop the service
	log.Logger().Infof("Stopping telemetry service: %s", serviceName)
	stopCmd := exec.Command("sc", "stop", serviceName)
	if err := stopCmd.Run(); err != nil {
		// Try using PowerShell as fallback
		log.Logger().Warnf("Failed to stop service with sc.exe, trying PowerShell: %v", err)
		psCmd := exec.Command("powershell", "-Command", fmt.Sprintf("Stop-Service -Name '%s' -Force", serviceName))
		if err := psCmd.Run(); err != nil {
			log.Logger().Errorf("Failed to stop telemetry service with PowerShell: %v", err)
			return err
		}

	}

	// Optionally disable the service
	disableCmd := exec.Command("sc", "config", serviceName, "start=", "disabled")
	if err := disableCmd.Run(); err != nil {
		log.Logger().Warnf("Failed to disable telemetry service (non-critical): %v", err)
	}

	log.Logger().Infof("Successfully stopped telemetry service: %s", serviceName)
	return nil
}
