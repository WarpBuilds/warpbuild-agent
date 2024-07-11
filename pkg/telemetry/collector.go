package telemetry

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

func startOtelCollector(baseDirectory string, collectorPath string, done chan bool) {
	cmd := exec.Command(collectorPath, "--config", getConfigFilePath(baseDirectory))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Logger().Errorf("failed to start OpenTelemetry Collector: %v", err)
		return
	}

	go func() {
		defer handlePanic()
		<-done
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			log.Logger().Errorf("Failed to terminate OpenTelemetry Collector: %v", err)
		}
	}()

	go func() {
		defer handlePanic()
		if err := cmd.Wait(); err != nil {
			log.Logger().Errorf("OpenTelemetry Collector exited with error: %v", err)
		}
	}()
}
