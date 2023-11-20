package manager

import (
	"context"
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
}

func NewAgent(opts *AgentOptions) (IAgent, error) {
	wb := warpbuild.NewAPIClient(&warpbuild.Configuration{})
	return &agentImpl{
		client:        wb,
		id:            opts.ID,
		pollingSecret: opts.PollingSecret,
	}, nil
}

type agentImpl struct {
	client        *warpbuild.APIClient
	id            string
	pollingSecret string
}

func (a *agentImpl) StartAgent(ctx context.Context, opts *StartAgentOptions) error {

	ticker := time.NewTicker(200 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			log.Logger().Infof("checking for runner instance allocation details for %s", a.id)

			allocationDetails, _, err := a.client.V1RunnerInstanceAPI.
				GetRunnerInstanceAllocationDetails(ctx, a.id).
				XPOLLINGSECRET(a.pollingSecret).
				Execute()
			if err != nil {
				log.Logger().Errorf("failed to get runner instance allocation details: %v", err)
				return err
			}

			// TODO: verify the correct status
			if *allocationDetails.RunnerInstance.Status == "allocated" {
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
