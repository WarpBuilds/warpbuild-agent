package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

type StartAgentOptions struct {
	Manager *ManagerOptions `json:"manager"`
}

type IAgent interface {
	StartAgent(ctx context.Context, opts *StartAgentOptions) error
}

type AgentOptions struct {
	// ID is the warpbuild assigned id.
	ID            string `json:"id"`
	PollingSecret string `json:"pollingSecret"`
	HostURL       string `json:"hostURL"`
}

func NewAgent(opts *AgentOptions) (IAgent, error) {
	cfg := warpbuild.NewConfiguration()

	if opts.HostURL == "" {
		return nil, fmt.Errorf("host url is required")
	}

	cfg.Servers[0].URL = opts.HostURL

	wb := warpbuild.NewAPIClient(cfg)
	return &agentImpl{
		client:        wb,
		id:            opts.ID,
		pollingSecret: opts.PollingSecret,
		hostURL:       opts.HostURL,
	}, nil
}

type agentImpl struct {
	client        *warpbuild.APIClient
	id            string
	pollingSecret string
	hostURL       string
}

func (a *agentImpl) StartAgent(ctx context.Context, opts *StartAgentOptions) error {

	ticker := time.NewTicker(200 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			log.Logger().Errorf("host url: %s", a.hostURL)
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

			// TODO: verify the correct status
			if *allocationDetails.Status == "assigned" {
				m := NewManager(opts.Manager)
				err := m.StartRunner(ctx, &StartRunnerOptions{
					JitToken: *allocationDetails.GhRunnerApplicationDetails.Jit,
				})
				if err != nil {
					log.Logger().Errorf("failed to start runner: %v", err)
					return err
				}
			}

		case <-ctx.Done():
			log.Logger().Infof("Context cancelled. Agent is exiting...")
			return nil
		}
	}

	return nil
}
