package app

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/manager"
)

type ApplicationOptions struct {
	SettingsFile string `json:"settings_file"`
}

func (opts *ApplicationOptions) Default() {
	if opts.SettingsFile == "" {
		opts.SettingsFile = "/var/lib/warpbuild-agent/settings.json"
	}
}

type Settings struct {
	Agent  *AgentSettings  `json:"agent"`
	Runner *RunnerSettings `json:"runner"`
}

type AgentSettings struct {
	ID string `json:"id"`
}

type RunnerSettings struct {
	Provider Provider        `json:"provider"`
	Github   *GithubSettings `json:"github"`
}

type GithubSettings struct {
	RunnerDir  string `json:"runner_dir"`
	Script     string `json:"script"`
	StdoutFile string `json:"stdout_file"`
	StderrFile string `json:"stderr_file"`
}

type Provider string

const (
	ProviderGithub Provider = "github"
)

func NewApp(ctx context.Context, opts *ApplicationOptions) error {

	err := log.Init()
	if err != nil {
		return err
	}

	log.Logger().Infof("starting warpbuild agent")

	var settings Settings
	var foundSettings bool
	// read the settings file every 200ms
	// the settings file might not be present at startup
	ticker := time.NewTicker(200 * time.Millisecond)
	timeout := time.After(120 * time.Second)
	for {
		select {
		case <-ticker.C:

			// read the settings file
			settingsData, err := os.ReadFile(opts.SettingsFile)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				log.Logger().Errorf("failed to read settings file: %v", err)
				return err
			}

			log.Logger().Infof("found settings file at %s", opts.SettingsFile)

			// found the settings file, parse it
			if err := json.Unmarshal(settingsData, &settings); err != nil {
				log.Logger().Errorf("failed to parse settings file: %v", err)
				return err
			}

			log.Logger().Debugf("settings: %v", settings)

			foundSettings = true

		case <-timeout:
			log.Logger().Errorf("timed out waiting for settings file")
			return nil

		case <-ctx.Done():
			log.Logger().Infof("context cancelled")
			return nil
		}

		if foundSettings {
			break
		}
	}

	agent, err := manager.NewAgent(&manager.AgentOptions{
		ID: settings.Agent.ID,
	})
	if err != nil {
		log.Logger().Errorf("failed to create agent: %v", err)
		return err
	}

	err = agent.StartAgent(ctx, &manager.StartAgentOptions{
		Manager: &manager.ManagerOptions{
			Provider: manager.Provider(string(settings.Runner.Provider)),
			Github: &manager.GithubOptions{
				RunnerDir:  settings.Runner.Github.RunnerDir,
				Script:     settings.Runner.Github.Script,
				StdoutFile: settings.Runner.Github.StdoutFile,
				StderrFile: settings.Runner.Github.StderrFile,
			},
		},
	})
	if err != nil {
		log.Logger().Errorf("failed to start agent: %v", err)
		return err
	}

	return nil
}
