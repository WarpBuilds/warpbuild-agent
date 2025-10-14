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
	client           *warpbuild.APIClient
	id               string
	pollingSecret    string
	hostURL          string
	exitFileLocation string
	opts             *AgentOptions
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

			// Check telemetry status and terminate process if disabled
			if allocationDetails.HasTelemetryEnabled() {
				telemetryEnabled := allocationDetails.GetTelemetryEnabled()
				log.Logger().Infof("Telemetry enabled status: %v", telemetryEnabled)

				if !telemetryEnabled {
					log.Logger().Infof("Telemetry is disabled. Terminating telemetry process...")
					if err := a.killTelemetryProcess(); err != nil {
						log.Logger().Errorf("Failed to terminate telemetry process: %v", err)
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
		log.Logger().Errorf("Failed to stop telemetry service: %v", err)
	} else {
		log.Logger().Infof("Successfully stopped telemetry service")
	}

	return err
}

// stopTelemetryServiceLinux stops the telemetry service on Linux using systemctl
func (a *agentImpl) stopTelemetryServiceLinux() error {
	serviceName := "warpbuild-telemetryd"

	// Check if service is running
	checkCmd := exec.Command("systemctl", "is-active", serviceName)
	output, err := checkCmd.CombinedOutput()
	status := strings.TrimSpace(string(output))

	log.Logger().Infof("Service status: %s", status)

	if status != "active" {
		log.Logger().Infof("Telemetry service is not active, nothing to stop")
		return nil
	}

	// Stop the service
	euid := os.Geteuid()
	var stopCmd *exec.Cmd
	if euid != 0 {
		stopCmd = exec.Command("sudo", "-n", "systemctl", "stop", serviceName)
		log.Logger().Infof("Stopping service with sudo (euid: %d)", euid)
	} else {
		stopCmd = exec.Command("systemctl", "stop", serviceName)
		log.Logger().Infof("Stopping service as root")
	}

	log.Logger().Infof("Executing: %v", stopCmd.Args)
	output, err = stopCmd.CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to stop service: %v, output: %s", err, string(output))
		log.Logger().Errorf(errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	log.Logger().Infof("Successfully stopped telemetry service")

	// Verify the service stopped
	time.Sleep(1 * time.Second)
	verifyCmd := exec.Command("systemctl", "is-active", serviceName)
	output, _ = verifyCmd.CombinedOutput()
	newStatus := strings.TrimSpace(string(output))
	log.Logger().Infof("Post-stop status: %s", newStatus)

	return nil
}

// stopTelemetryServiceDarwin stops the telemetry service on macOS using launchctl
func (a *agentImpl) stopTelemetryServiceDarwin() error {
	serviceName := "com.warpbuild.warpbuild-telemetryd"

	// Check if service is loaded
	listCmd := exec.Command("launchctl", "list", serviceName)
	output, err := listCmd.CombinedOutput()

	if err != nil {
		log.Logger().Infof("Telemetry service is not loaded")
		return nil
	}

	log.Logger().Infof("Service is running: %s", strings.TrimSpace(string(output)))

	// Stop the service with unload
	euid := os.Geteuid()
	plistPath := fmt.Sprintf("/Library/LaunchDaemons/%s.plist", serviceName)

	var unloadCmd *exec.Cmd
	if euid != 0 {
		unloadCmd = exec.Command("sudo", "launchctl", "unload", "-w", plistPath)
		log.Logger().Infof("Unloading service with sudo")
	} else {
		unloadCmd = exec.Command("launchctl", "unload", "-w", plistPath)
		log.Logger().Infof("Unloading service as root")
	}

	log.Logger().Infof("Executing: %v", unloadCmd.Args)
	output, err = unloadCmd.CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to unload service: %v, output: %s", err, string(output))
		log.Logger().Errorf(errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	log.Logger().Infof("Successfully stopped telemetry service")

	// Verify the service stopped
	time.Sleep(1 * time.Second)
	verifyCmd := exec.Command("launchctl", "list", serviceName)
	_, err = verifyCmd.CombinedOutput()

	if err != nil {
		log.Logger().Infof("Post-stop: service not found (success)")
	} else {
		log.Logger().Warnf("Post-stop: service still listed")
	}

	return nil
}

// stopTelemetryServiceWindows stops the telemetry service on Windows
func (a *agentImpl) stopTelemetryServiceWindows() error {
	serviceName := "warpbuild-telemetryd"

	// Check if service exists and is running
	checkCmd := exec.Command("sc", "query", serviceName)
	output, err := checkCmd.CombinedOutput()

	if err != nil {
		log.Logger().Infof("Telemetry service does not exist: %s", serviceName)
		return nil
	}

	if !strings.Contains(string(output), "RUNNING") {
		log.Logger().Infof("Telemetry service is not running")
		return nil
	}

	log.Logger().Infof("Stopping service: %s", serviceName)

	// Stop the service
	stopCmd := exec.Command("sc", "stop", serviceName)
	output, err = stopCmd.CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to stop service: %v, output: %s", err, string(output))
		log.Logger().Errorf(errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	log.Logger().Infof("Successfully stopped telemetry service")

	// Verify the service stopped
	time.Sleep(1 * time.Second)
	verifyCmd := exec.Command("sc", "query", serviceName)
	output, _ = verifyCmd.CombinedOutput()

	if strings.Contains(string(output), "STOPPED") {
		log.Logger().Infof("Post-stop: service stopped")
	} else {
		log.Logger().Warnf("Post-stop: service status unclear: %s", string(output))
	}

	return nil
}
