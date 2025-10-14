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

	// Syslog priority levels (for cross-platform compatibility)
	logInfo    = 6
	logWarning = 4
	logErr     = 3
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
	a.logToSyslog(logInfo, "Attempting to stop telemetry service")

	var err error
	switch runtime.GOOS {
	case "linux":
		err = a.stopTelemetryServiceLinux()
	case "darwin":
		err = a.stopTelemetryServiceDarwin()
	case "windows":
		err = a.stopTelemetryServiceWindows()
	default:
		err = fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	if err != nil {
		a.logToSyslog(logErr, fmt.Sprintf("Failed to stop telemetry service: %v", err))
	} else {
		a.logToSyslog(logInfo, "Successfully stopped telemetry service")
	}

	return err
}

// stopTelemetryServiceLinux stops the telemetry service on Linux using systemctl
func (a *agentImpl) stopTelemetryServiceLinux() error {
	serviceName := "warpbuild-telemetryd"
	a.logToSyslog(logInfo, fmt.Sprintf("Linux: Checking telemetry service status: %s", serviceName))

	// Check if running as root
	if os.Geteuid() != 0 {
		a.logToSyslog(logWarning, "Not running as root, will attempt with sudo")
		log.Logger().Warnf("Not running as root (euid: %d), will attempt with sudo", os.Geteuid())
	}

	// Check if service exists and is running
	checkCmd := exec.Command("systemctl", "is-active", serviceName)
	output, err := checkCmd.CombinedOutput()
	status := strings.TrimSpace(string(output))

	a.logToSyslog(logInfo, fmt.Sprintf("Service status check result: %s, error: %v", status, err))
	log.Logger().Infof("Service status: %s", status)

	if err != nil || status != "active" {
		msg := fmt.Sprintf("Telemetry service is not running (status: %s)", status)
		log.Logger().Infof(msg)
		a.logToSyslog(logInfo, msg)
		return nil
	}

	// Try to stop the service with sudo if needed
	var stopCmd *exec.Cmd
	if os.Geteuid() != 0 {
		stopCmd = exec.Command("sudo", "systemctl", "stop", serviceName)
	} else {
		stopCmd = exec.Command("systemctl", "stop", serviceName)
	}

	a.logToSyslog(logInfo, fmt.Sprintf("Stopping telemetry service with command: %v", stopCmd.Args))
	log.Logger().Infof("Stopping telemetry service: %s", serviceName)

	output, err = stopCmd.CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to stop telemetry service: %v, output: %s", err, string(output))
		log.Logger().Errorf(errMsg)
		a.logToSyslog(logErr, errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	a.logToSyslog(logInfo, "Service stopped successfully")

	// Verify the service stopped
	time.Sleep(1 * time.Second)
	verifyCmd := exec.Command("systemctl", "is-active", serviceName)
	output, _ = verifyCmd.CombinedOutput()
	newStatus := strings.TrimSpace(string(output))
	a.logToSyslog(logInfo, fmt.Sprintf("Post-stop verification status: %s", newStatus))
	log.Logger().Infof("Post-stop service status: %s", newStatus)

	// Optionally disable the service to prevent it from starting on next boot
	var disableCmd *exec.Cmd
	if os.Geteuid() != 0 {
		disableCmd = exec.Command("sudo", "systemctl", "disable", serviceName)
	} else {
		disableCmd = exec.Command("systemctl", "disable", serviceName)
	}

	output, err = disableCmd.CombinedOutput()
	if err != nil {
		warnMsg := fmt.Sprintf("Failed to disable telemetry service (non-critical): %v, output: %s", err, string(output))
		log.Logger().Warnf(warnMsg)
		a.logToSyslog(logWarning, warnMsg)
	} else {
		a.logToSyslog(logInfo, "Service disabled successfully")
	}

	log.Logger().Infof("Successfully stopped telemetry service: %s", serviceName)
	return nil
}

// stopTelemetryServiceDarwin stops the telemetry service on macOS using launchctl
func (a *agentImpl) stopTelemetryServiceDarwin() error {
	serviceName := "com.warpbuild.warpbuild-telemetryd"
	a.logToSyslog(logInfo, fmt.Sprintf("macOS: Checking telemetry service: %s", serviceName))

	// Check if running as root
	if os.Geteuid() != 0 {
		a.logToSyslog(logWarning, "Not running as root, will attempt with sudo")
		log.Logger().Warnf("Not running as root (euid: %d), will attempt with sudo", os.Geteuid())
	}

	log.Logger().Infof("Stopping telemetry service: %s", serviceName)

	// First, try to find the service using launchctl list
	listCmd := exec.Command("launchctl", "list", serviceName)
	output, err := listCmd.CombinedOutput()

	a.logToSyslog(logInfo, fmt.Sprintf("Service list check output: %s, error: %v", string(output), err))

	if err != nil {
		msg := fmt.Sprintf("Telemetry service is not loaded: %s", serviceName)
		log.Logger().Infof(msg)
		a.logToSyslog(logInfo, msg)
		return nil
	}

	// Service is running
	a.logToSyslog(logInfo, fmt.Sprintf("Service info: %s", string(output)))
	log.Logger().Infof("Service is running, details: %s", string(output))

	// Try multiple approaches to stop the service
	attempts := []struct {
		name string
		cmd  *exec.Cmd
	}{
		{
			name: "kill with SIGTERM",
			cmd:  exec.Command("sudo", "launchctl", "kill", "SIGTERM", fmt.Sprintf("system/%s", serviceName)),
		},
		{
			name: "stop service",
			cmd:  exec.Command("sudo", "launchctl", "stop", serviceName),
		},
		{
			name: "bootout system",
			cmd:  exec.Command("sudo", "launchctl", "bootout", fmt.Sprintf("system/%s", serviceName)),
		},
		{
			name: "unload plist",
			cmd:  exec.Command("sudo", "launchctl", "unload", "-w", fmt.Sprintf("/Library/LaunchDaemons/%s.plist", serviceName)),
		},
	}

	var lastErr error
	for _, attempt := range attempts {
		a.logToSyslog(logInfo, fmt.Sprintf("Attempting: %s with command: %v", attempt.name, attempt.cmd.Args))
		log.Logger().Infof("Attempting to %s", attempt.name)

		output, err = attempt.cmd.CombinedOutput()
		if err != nil {
			warnMsg := fmt.Sprintf("Failed to %s: %v, output: %s", attempt.name, err, string(output))
			log.Logger().Warnf(warnMsg)
			a.logToSyslog(logWarning, warnMsg)
			lastErr = err
		} else {
			a.logToSyslog(logInfo, fmt.Sprintf("Successfully executed: %s", attempt.name))
			log.Logger().Infof("Successfully executed: %s", attempt.name)
			lastErr = nil
			break
		}
	}

	if lastErr != nil {
		errMsg := fmt.Sprintf("All attempts to stop telemetry service failed: %v", lastErr)
		a.logToSyslog(logErr, errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// Verify the service stopped
	time.Sleep(2 * time.Second)
	verifyCmd := exec.Command("launchctl", "list", serviceName)
	output, err = verifyCmd.CombinedOutput()

	if err != nil {
		// Service not found means it's stopped
		a.logToSyslog(logInfo, "Post-stop verification: service not found (successfully stopped)")
		log.Logger().Infof("Service successfully stopped and unloaded")
	} else {
		a.logToSyslog(logWarning, fmt.Sprintf("Post-stop verification: service still listed: %s", string(output)))
		log.Logger().Warnf("Service still appears in launchctl list: %s", string(output))
	}

	log.Logger().Infof("Successfully stopped telemetry service: %s", serviceName)
	return nil
}

// stopTelemetryServiceWindows stops the telemetry service on Windows
func (a *agentImpl) stopTelemetryServiceWindows() error {
	serviceName := "warpbuild-telemetryd"

	// Note: Windows doesn't have syslog, but we log to regular logger
	log.Logger().Infof("Windows: Checking telemetry service: %s", serviceName)

	// Check if service exists
	checkCmd := exec.Command("sc", "query", serviceName)
	output, err := checkCmd.CombinedOutput()

	log.Logger().Infof("Service query output: %s, error: %v", string(output), err)

	if err != nil {
		log.Logger().Infof("Telemetry service does not exist or is not accessible: %s", serviceName)
		return nil
	}

	// Check if service is running
	if !strings.Contains(string(output), "RUNNING") {
		log.Logger().Infof("Telemetry service is not running. Status: %s", string(output))
		return nil
	}

	log.Logger().Infof("Service is running. Attempting to stop: %s", serviceName)

	// Try multiple approaches to stop the service
	attempts := []struct {
		name string
		fn   func() error
	}{
		{
			name: "sc.exe stop",
			fn: func() error {
				stopCmd := exec.Command("sc", "stop", serviceName)
				output, err := stopCmd.CombinedOutput()
				if err != nil {
					return fmt.Errorf("sc stop failed: %v, output: %s", err, string(output))
				}
				log.Logger().Infof("sc stop output: %s", string(output))
				return nil
			},
		},
		{
			name: "PowerShell Stop-Service",
			fn: func() error {
				psCmd := exec.Command("powershell", "-Command", fmt.Sprintf("Stop-Service -Name '%s' -Force", serviceName))
				output, err := psCmd.CombinedOutput()
				if err != nil {
					return fmt.Errorf("PowerShell Stop-Service failed: %v, output: %s", err, string(output))
				}
				log.Logger().Infof("PowerShell Stop-Service output: %s", string(output))
				return nil
			},
		},
		{
			name: "net stop",
			fn: func() error {
				netCmd := exec.Command("net", "stop", serviceName)
				output, err := netCmd.CombinedOutput()
				if err != nil {
					return fmt.Errorf("net stop failed: %v, output: %s", err, string(output))
				}
				log.Logger().Infof("net stop output: %s", string(output))
				return nil
			},
		},
	}

	var lastErr error
	for _, attempt := range attempts {
		log.Logger().Infof("Attempting: %s", attempt.name)
		err := attempt.fn()
		if err != nil {
			log.Logger().Warnf("Failed to %s: %v", attempt.name, err)
			lastErr = err
		} else {
			log.Logger().Infof("Successfully stopped service using: %s", attempt.name)
			lastErr = nil
			break
		}
	}

	if lastErr != nil {
		return fmt.Errorf("all attempts to stop telemetry service failed: %v", lastErr)
	}

	// Verify the service stopped
	time.Sleep(2 * time.Second)
	verifyCmd := exec.Command("sc", "query", serviceName)
	output, _ = verifyCmd.CombinedOutput()
	log.Logger().Infof("Post-stop service status: %s", string(output))

	// Optionally disable the service
	disableCmd := exec.Command("sc", "config", serviceName, "start=", "disabled")
	output, err = disableCmd.CombinedOutput()
	if err != nil {
		log.Logger().Warnf("Failed to disable telemetry service (non-critical): %v, output: %s", err, string(output))
	} else {
		log.Logger().Infof("Successfully disabled telemetry service")
	}

	log.Logger().Infof("Successfully stopped telemetry service: %s", serviceName)
	return nil
}
