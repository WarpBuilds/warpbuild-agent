package telemetry

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

const (
	configFilePath              = "/runner/warpbuild-agent/pkg/telemetry/otel-collector-config.yaml"
	configTemplatePath          = "/runner/warpbuild-agent/pkg/telemetry/otel-collector-config.tmpl"
	otelCollectorOutputFilePath = "/runner/warpbuild-agent/pkg/telemetry/otel-out.log"
	binariesDir                 = "/runner/warpbuild-agent/pkg/telemetry/binaries"
)

var syslogFilePath = getSyslogFilePath()
var presignedS3URL = ""

func getSyslogFilePath() string {
	switch runtime.GOOS {
	case "darwin":
		return "/var/log/system.log"
	case "linux":
		return "/var/log/syslog"
	case "windows":
		return `C:\Windows\System32\winevt\Logs\System.evtx`
	default:
		log.Logger().Errorf("Unsupported OS: %s", runtime.GOOS)
		return ""
	}
}

func getOtelCollectorPath() (string, error) {
	var collectorPath string
	systemArch := runtime.GOARCH
	systemOS := runtime.GOOS

	switch systemOS {
	case "linux":
		switch systemArch {
		case "amd64":
			collectorPath = filepath.Join(binariesDir, "linux", "amd64", "otelcol-contrib")
		case "arm64":
			collectorPath = filepath.Join(binariesDir, "linux", "arm64", "otelcol-contrib")
		default:
			return "", fmt.Errorf("unsupported architecture: %s", systemArch)
		}
	case "darwin":
		switch systemArch {
		case "amd64":
			collectorPath = filepath.Join(binariesDir, "darwin", "amd64", "otelcol-contrib")
		case "arm64":
			collectorPath = filepath.Join(binariesDir, "darwin", "arm64", "otelcol-contrib")
		default:
			return "", fmt.Errorf("unsupported architecture: %s", systemArch)
		}
	case "windows":
		if systemArch == "amd64" {
			collectorPath = filepath.Join(binariesDir, "windows", "amd64", "otelcol-contrib.exe")
		} else {
			return "", fmt.Errorf("unsupported architecture: %s", systemArch)
		}
	default:
		return "", fmt.Errorf("unsupported OS: %s", systemOS)
	}

	// Ensure the binary exists
	if _, err := os.Stat(collectorPath); os.IsNotExist(err) {
		return "", fmt.Errorf("collector binary not found at %s", collectorPath)
	}

	// Make the binary executable
	if systemOS != "windows" {
		if err := os.Chmod(collectorPath, 0755); err != nil {
			return "", fmt.Errorf("failed to make the OpenTelemetry Collector binary executable: %w", err)
		}
	}

	if systemOS == "darwin" {
		if err := exec.Command("xattr", "-rd", "com.apple.quarantine", collectorPath).Run(); err != nil {
			return "", fmt.Errorf("failed to remove quarantine attribute from binary: %w", err)
		}
	}

	return collectorPath, nil
}

func writeOtelCollectorConfig() error {
	tmpl, err := template.ParseFiles(configTemplatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template file: %w", err)
	}

	file, err := os.Create(configFilePath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	data := struct {
		SyslogFilePath string
		ExportFilePath string
	}{
		SyslogFilePath: syslogFilePath,
		ExportFilePath: otelCollectorOutputFilePath,
	}

	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
