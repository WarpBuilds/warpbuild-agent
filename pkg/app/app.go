package app

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/manager"
	"github.com/warpbuilds/warpbuild-agent/pkg/proxy"
	"github.com/warpbuilds/warpbuild-agent/pkg/telemetry"
	transparentcache "github.com/warpbuilds/warpbuild-agent/pkg/transparent-cache"
)

type ApplicationOptions struct {
	SettingsFile           string `json:"settings_file"`
	StdoutFile             string `json:"stdout_file"`
	StderrFile             string `json:"stderr_file"`
	LaunchTelemetry        bool   `json:"launch_telemetry"`
	LaunchProxyServer      bool   `json:"launch_cache_proxy_server"`
	LaunchTransparentCache bool   `json:"launch_transparent_cache"`
	LogLevel               string `json:"log_level"`
}

func (opts *ApplicationOptions) ApplyDefaults() {
	if opts.SettingsFile == "" {
		opts.SettingsFile = "/var/lib/warpbuild-agentd/settings.json"
	}
	if opts.LogLevel == "" {
		opts.LogLevel = "info"
	}
}

type Settings struct {
	Agent            *AgentSettings            `json:"agent"`
	Runner           *RunnerSettings           `json:"runner"`
	Telemetry        *TelemetrySettings        `json:"telemetry"`
	Proxy            *ProxySettings            `json:"proxy"`
	TransparentCache *TransparentCacheSettings `json:"transparent_cache"`
}

func (s *Settings) ApplyDefaults() {
	if s.Telemetry != nil {
		s.Telemetry.ApplyDefaults()
	}
	if s.TransparentCache != nil {
		s.TransparentCache.ApplyDefaults()
	}
}

type AgentSettings struct {
	ID                      string `json:"id"`
	PollingSecret           string `json:"polling_secret"`
	RunnerVerificationToken string `json:"runner_verification_token"`
	HostURL                 string `json:"host_url"`
	ExitFileLocation        string `json:"exit_file_location"`
}

type TelemetrySettings struct {
	BaseDirectory string `json:"base_directory"`
	Enabled       bool   `json:"enabled"`
	// At what frequency to push the telemetry data to the server. This is in seconds.
	PushFrequency string `json:"push_frequency"`
	// Port is the port on which the otel receiver is exposed.
	//
	// Default: 33931
	Port int `json:"port"`
}

func (t *TelemetrySettings) ApplyDefaults() {

	if t.Port == 0 {
		t.Port = 33931
	}

}

type ProxySettings struct {
	CacheProxyPort   string `json:"cache_proxy_port"`
	CacheBackendHost string `json:"cache_backend_host"`
}

type TransparentCacheSettings struct {
	Enabled        bool `json:"enabled"`
	LoggingEnabled bool `json:"logging_enabled"`
	DerpPort       int  `json:"derp_port"`
	OginyPort      int  `json:"oginy_port"`
	AsurPort       int  `json:"asur_port"`
}

func (t *TransparentCacheSettings) ApplyDefaults() {
	if t.DerpPort == 0 {
		t.DerpPort = 50052
	}
	if t.OginyPort == 0 {
		t.OginyPort = 50051
	}
	if t.AsurPort == 0 {
		t.AsurPort = 50053
	}
}

type RunnerSettings struct {
	Provider         Provider                         `json:"provider"`
	Github           *GithubSettings                  `json:"github"`
	GithubCRI        *manager.GithubCRIOptions        `json:"github_cri"`
	GithubWindowsCRI *manager.GithubWindowsCRIOptions `json:"github_windows_cri"`
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

	// Initialize logger with default level first
	lm, err := log.Init(&log.InitOptions{
		StdoutFile: opts.StdoutFile,
		StderrFile: opts.StderrFile,
		LogLevel:   opts.LogLevel, // Use application-level log level
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
				if strings.Contains(err.Error(), "unexpected end of JSON input") {
					log.Logger().Infof("unexpected end of JSON input. We'll retry.")
					continue
				}
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

	settings.ApplyDefaults()

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
			BaseDirectory: settings.Telemetry.BaseDirectory,
			RunnerID:      settings.Agent.ID,
			PollingSecret: settings.Agent.PollingSecret,
			HostURL:       settings.Agent.HostURL,
			Enabled:       settings.Telemetry.Enabled,
			PushFrequency: pushFrequency,
			Port:          settings.Telemetry.Port,
		}); err != nil {
			log.Logger().Errorf("failed to start telemetry: %v", err)
		}

	} else if opts.LaunchProxyServer {
		proxy.StartProxyServer(ctx, &proxy.ProxyServerOptions{
			CacheBackendHost:                 settings.Proxy.CacheBackendHost,
			CacheProxyPort:                   settings.Proxy.CacheProxyPort,
			WarpBuildRunnerVerificationToken: settings.Agent.RunnerVerificationToken,
		})
	} else if opts.LaunchTransparentCache {
		// Start the transparent cache server with configured ports
		if settings.TransparentCache == nil {
			log.Logger().Errorf("transparent cache settings not configured")
			return errors.New("transparent cache settings not configured")
		}

		// Check if transparent cache is enabled
		if !settings.TransparentCache.Enabled {
			log.Logger().Infof("transparent cache is disabled in settings")
			return nil
		}

		if err := transparentcache.Start(
			settings.TransparentCache.DerpPort,
			settings.TransparentCache.OginyPort,
			settings.TransparentCache.AsurPort,
			settings.Proxy.CacheBackendHost,
			settings.Agent.RunnerVerificationToken,
			settings.TransparentCache.LoggingEnabled,
		); err != nil {
			log.Logger().Errorf("failed to start transparent cache: %v", err)
			return err
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
		case manager.ProviderGithubWindowsCRI:
			startAgentOpts.Manager.GithubWindowsCRI = settings.Runner.GithubWindowsCRI
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
