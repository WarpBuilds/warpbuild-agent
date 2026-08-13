package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	transparentcache "github.com/warpbuilds/warpbuild-agent/pkg/transparent-cache"
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
	// TransparentCacheOginyPort is the port for the transparent cache oginy server.
	TransparentCacheOginyPort int `json:"transparent_cache_oginy_port"`
	// CacheBackendHost + RunnerVerificationToken let a claude agent upload its deliverables to
	// backend-cache (via the node warp-cache client).
	CacheBackendHost        string `json:"cache_backend_host"`
	RunnerVerificationToken string `json:"runner_verification_token"`
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
	cfg.UserAgent = "warpbuild-agent"

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

			// TODO: verify the correct status
			if allocationDetails.Status == nil || *allocationDetails.Status != "assigned" {
				status := "<nil>"
				if allocationDetails.Status != nil {
					status = *allocationDetails.Status
				}
				log.Logger().Infof("runner instance allocation details status: %s", status)
				log.Logger().Infof("Retrying in %s", Interval)
				continue
			}

			switch runnerApplication(allocationDetails) {
			case ProviderClaudeAgent:
				if err := a.handleClaudeAgentAllocation(ctx, allocationDetails); err != nil {
					return err
				}
			default:
				if err := a.handleGithubAllocation(ctx, opts, allocationDetails); err != nil {
					return err
				}
			}

		case <-ctx.Done():
			log.Logger().Infof("Context cancelled. Agent is exiting...")
			return nil
		}
	}

}

func runnerApplication(details *warpbuild.CommonsRunnerInstanceAllocationDetails) Provider {
	if details.RunnerApplication != nil {
		return Provider(*details.RunnerApplication)
	}
	return ProviderGithub
}

func (a *agentImpl) handleGithubAllocation(ctx context.Context, opts *StartAgentOptions, allocationDetails *warpbuild.CommonsRunnerInstanceAllocationDetails) error {
	if allocationDetails.GhRunnerApplicationDetails == nil || allocationDetails.GhRunnerApplicationDetails.Variables == nil {
		log.Logger().Warnf("assigned allocation missing GitHub runner application details; retrying in %s", Interval)
		return nil
	}

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

	// Read WARPBUILD_TRANSPARENT_CACHE_ENABLED and if it is true, then start the transparent cache server
	if (*allocationDetails.GhRunnerApplicationDetails.Variables)["WARPBUILD_TRANSPARENT_CACHE_ENABLED"] == "true" {
		log.Logger().Infof("Starting transparent cache server")
		if err := transparentcache.SetupNetworking(a.opts.TransparentCacheOginyPort); err != nil {
			log.Logger().Errorf("failed to configure networking for transparent cache: %v", err)
			return err
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
		if err := a.writeExitFile(ctx, startRunnerOutput); err != nil {
			log.Logger().Errorf("failed to write exit file: %v", err)
			return err
		}
	}
	return nil
}

func (a *agentImpl) handleClaudeAgentAllocation(ctx context.Context, allocationDetails *warpbuild.CommonsRunnerInstanceAllocationDetails) error {
	startRunnerOutput, err := a.startClaudeAgent(ctx, allocationDetails)
	if err != nil {
		log.Logger().Errorf("failed to start claude agent worker: %v", err)
		return err
	}
	if startRunnerOutput.RunCompletedSuccessfully {
		if err := a.writeExitFile(ctx, startRunnerOutput); err != nil {
			log.Logger().Errorf("failed to write exit file: %v", err)
			return err
		}
	}
	return nil
}

func (a *agentImpl) startClaudeAgent(ctx context.Context, allocationDetails *warpbuild.CommonsRunnerInstanceAllocationDetails) (*StartRunnerOutput, error) {
	details := allocationDetails.ClaudeAgentApplicationDetails
	if details == nil {
		return nil, fmt.Errorf("claude_agent allocation is missing claude_agent_application_details")
	}

	log.Logger().Infof("Starting Claude managed-agent worker for session %s", details.GetSessionId())

	copts := DefaultClaudeOptions(details.GetMaxIdle())
	copts.HostURL = a.hostURL
	copts.PollingSecret = a.pollingSecret
	copts.RunnerInstanceID = a.id
	copts.EnvID = details.GetEnvId()
	copts.EnvKey = details.GetEnvKey()
	copts.SessionID = details.GetSessionId()
	if details.WorkId != nil {
		copts.WorkID = *details.WorkId
	}
	copts.CacheBackendHost = a.opts.CacheBackendHost
	copts.RunnerVerificationToken = a.opts.RunnerVerificationToken
	if err := provisionClaudeWorker(copts); err != nil {
		return nil, err
	}
	m := NewManager(&ManagerOptions{Provider: ProviderClaudeAgent, Claude: copts})
	return m.StartRunner(ctx, &StartRunnerOptions{AgentOptions: a.opts})
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
