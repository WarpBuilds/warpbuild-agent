package app

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/manager"
	"github.com/warpbuilds/warpbuild-agent/pkg/telemetry"
)

type ApplicationOptions struct {
	SettingsFile    string `json:"settings_file"`
	StdoutFile      string `json:"stdout_file"`
	StderrFile      string `json:"stderr_file"`
	LaunchTelemetry bool   `json:"launch_telemetry"`
}

func (opts *ApplicationOptions) ApplyDefaults() {
	if opts.SettingsFile == "" {
		opts.SettingsFile = "/var/lib/warpbuild-agentd/settings.json"
	}
}

type Settings struct {
	Agent     *AgentSettings     `json:"agent"`
	Runner    *RunnerSettings    `json:"runner"`
	Telemetry *TelemetrySettings `json:"telemetry"`
}

type AgentSettings struct {
	ID               string `json:"id"`
	PollingSecret    string `json:"polling_secret"`
	HostURL          string `json:"host_url"`
	ExitFileLocation string `json:"exit_file_location"`
}

type TelemetrySettings struct {
	BaseDirectory string `json:"base_directory"`
	Enabled       bool   `json:"enabled"`
	// The telemetry agent reads the defined number of lines from syslog file of the system and pushes to the server
	SysLogNumberOfLinesToRead int `json:"syslog_number_of_lines_to_read"`
	// At what frequency to push the telemetry data to the server. This is in seconds.
	PushFrequency string `json:"push_frequency"`
}

type RunnerSettings struct {
	Provider  Provider                  `json:"provider"`
	Github    *GithubSettings           `json:"github"`
	GithubCRI *manager.GithubCRIOptions `json:"github_cri"`
}

type ContainerOptionsVolume struct {
	HostPath      string `json:"host_path"`
	ContainerPath string `json:"container_path"`
	AccessMode    string `json:"access_mode"`
}

type GithubSettings struct {
	RunnerDir  string                       `json:"runner_dir"`
	Script     string                       `json:"script"`
	StdoutFile string                       `json:"stdout_file"`
	StderrFile string                       `json:"stderr_file"`
	Envs       manager.EnvironmentVariables `json:"envs"`
}

type Provider string

const (
	ProviderGithub Provider = "github"
	// ProviderGithubCRI is the provider for the github custom runner image.
	ProviderGithubCRI Provider = "github_cri"
)

func NewApp(ctx context.Context, opts *ApplicationOptions) error {

	opts.ApplyDefaults()

	lm, err := log.Init(&log.InitOptions{
		StdoutFile: opts.StdoutFile,
		StderrFile: opts.StderrFile,
	})
	if err != nil {
		return err
	}

	defer lm.Sync()

	log.Logger().Infof("starting warpbuild agent")
	log.Logger().Infof("settings file: %s", opts.SettingsFile)

	var elapsedTime time.Duration
	var settings Settings
	var foundSettings bool
	// read the settings file every 200ms
	// the settings file might not be present at startup
	ticker := time.NewTicker(200 * time.Millisecond)
	for {
		select {
		case <-ticker.C:

			log.Logger().Infof("checking for settings file at %s", opts.SettingsFile)
			log.Logger().Infof("elapsed time: %v", elapsedTime)

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

		case <-ctx.Done():
			log.Logger().Infof("context cancelled")
			return nil
		}

		if foundSettings {
			break
		}
	}

	if opts.LaunchTelemetry {
		telemetryCtx, telemetryCtxCancel := context.WithCancel(context.Background())
		defer telemetryCtxCancel()

		// Set up signal handling to catch OS kill signals
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			sig := <-sigs
			log.Logger().Infof("Received signal %s, initiating shutdown...", sig)
			telemetryCtxCancel()
		}()

		pushFrequency, _ := time.ParseDuration(settings.Telemetry.PushFrequency)
		if err := telemetry.StartTelemetryCollection(telemetryCtx, &telemetry.TelemetryOptions{
			BaseDirectory:             settings.Telemetry.BaseDirectory,
			RunnerID:                  settings.Agent.ID,
			PollingSecret:             settings.Agent.PollingSecret,
			HostURL:                   settings.Agent.HostURL,
			Enabled:                   settings.Telemetry.Enabled,
			PushFrequency:             pushFrequency,
			SysLogNumberOfLinesToRead: settings.Telemetry.SysLogNumberOfLinesToRead,
		}); err != nil {
			log.Logger().Errorf("failed to start telemetry: %v", err)
		}

	} else {
		agent, err := manager.NewAgent(&manager.AgentOptions{
			ID:               settings.Agent.ID,
			PollingSecret:    settings.Agent.PollingSecret,
			HostURL:          settings.Agent.HostURL,
			ExitFileLocation: settings.Agent.ExitFileLocation,
		})
		if err != nil {
			log.Logger().Errorf("failed to create agent: %v", err)
			return err
		}

		startAgentOpts := manager.StartAgentOptions{
			Manager: &manager.ManagerOptions{
				Provider: manager.Provider(string(settings.Runner.Provider)),
			},
		}
		switch startAgentOpts.Manager.Provider {
		case manager.ProviderGithub:
			startAgentOpts.Manager.Github = &manager.GithubOptions{
				RunnerDir:  settings.Runner.Github.RunnerDir,
				Script:     settings.Runner.Github.Script,
				StdoutFile: settings.Runner.Github.StdoutFile,
				StderrFile: settings.Runner.Github.StderrFile,
				Envs:       settings.Runner.Github.Envs,
			}
		case manager.ProviderGithubCRI:
			startAgentOpts.Manager.GithubCRI = settings.Runner.GithubCRI
		default:
			log.Logger().Errorf("unknown provider: %s", startAgentOpts.Manager.Provider)
			return errors.New("unknown provider")
		}

		err = agent.StartAgent(ctx, &startAgentOpts)
		if err != nil {
			log.Logger().Errorf("failed to start agent: %v", err)
			return err
		}
	}

	log.Logger().Infof("Shutdown complete.")
	return nil
}
