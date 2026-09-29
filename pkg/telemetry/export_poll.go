package telemetry

import (
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

const telemetryPollInterval = 5 * time.Second

const maxWaitForRunningJob = 10 * time.Minute

// The backend reports "unassigned" both before allocation and once the job runs; "assigned" only in between.
const allocationStatusUnassigned = "unassigned"

type pollStep int

const (
	pollAgain       pollStep = iota // nothing decided yet
	awaitRunningJob                 // export applied for the allocated job; wait for the job that actually runs
	stopPolling                     // export is final, or none is configured
	stopTelemetry                   // collection is off and there is nothing to export
)

func (tm *TelemetryManager) monitorTelemetryStatus() {
	defer tm.wg.Done()

	log.Logger().Infof("Watching for the organization's telemetry export destination...")

	ticker := time.NewTicker(telemetryPollInterval)
	defer ticker.Stop()

	var awaitingSince time.Time
	for {
		select {
		case <-tm.ctx.Done():
			log.Logger().Infof("Context cancelled, stopping telemetry status monitoring...")
			return
		case <-ticker.C:
		}

		details, err := tm.fetchAllocationDetails()
		if err != nil {
			continue
		}

		switch tm.handleAllocationDetails(details) {
		case pollAgain:
		case awaitRunningJob:
			if awaitingSince.IsZero() {
				awaitingSince = time.Now()
			} else if time.Since(awaitingSince) > maxWaitForRunningJob {
				return
			}
		case stopPolling:
			return
		case stopTelemetry:
			tm.cancel()
			return
		}
	}
}

func (tm *TelemetryManager) fetchAllocationDetails() (*warpbuild.CommonsRunnerInstanceAllocationDetails, error) {
	details, resp, err := tm.warpbuildAPI.V1RunnerInstanceAPI.
		GetRunnerInstanceAllocationDetails(tm.ctx, tm.runnerID).
		XPOLLINGSECRET(tm.pollingSecret).
		Execute()
	if err != nil {
		log.Logger().Debugf("Failed to get runner instance allocation details: %v", err)
		if resp != nil {
			log.Logger().Debugf("Response: %+v", resp)
		}
	}
	return details, err
}

func (tm *TelemetryManager) handleAllocationDetails(details *warpbuild.CommonsRunnerInstanceAllocationDetails) pollStep {
	if details == nil {
		return pollAgain
	}
	collect := collectionEnabled(details)
	allocated := details.GetStatus() != allocationStatusUnassigned

	if export := exportConfigFrom(details.TelemetryExport); export != nil {
		tm.applyExportConfig(collect, export)
		if allocated {
			// Not started yet: GitHub may still run a sibling job on this runner.
			return awaitRunningJob
		}
		return stopPolling
	}

	switch {
	case tm.currentExportConfig() != nil:
		return awaitRunningJob
	case allocated:
		return tm.finishWithoutExport(collect)
	default:
		tm.followCollectionToggle(collect)
		return pollAgain
	}
}

// An old backend omits the flag, which must not stop collection.
func collectionEnabled(details *warpbuild.CommonsRunnerInstanceAllocationDetails) bool {
	return !details.HasTelemetryEnabled() || details.GetTelemetryEnabled()
}

func (tm *TelemetryManager) finishWithoutExport(collect bool) pollStep {
	if !collect {
		log.Logger().Infof("Collection is disabled and no export is configured. Stopping telemetry...")
		return stopTelemetry
	}
	tm.unpark()
	log.Logger().Infof("No telemetry export configured for this organization; nothing to fan out")
	return stopPolling
}

func (tm *TelemetryManager) followCollectionToggle(collect bool) {
	if collect {
		tm.unpark()
		return
	}
	tm.park()
}

func (tm *TelemetryManager) applyExportConfig(collect bool, next *exportConfig) {
	tm.mu.Lock()
	if tm.exportCfg.equal(next) && tm.collect == collect {
		tm.mu.Unlock()
		return
	}
	tm.exportCfg = next
	tm.collect = collect
	tm.parked = false
	tm.mu.Unlock()

	log.Logger().Infof("Telemetry export configured: metrics=%q logs=%q collect=%t",
		next.MetricsEndpoint, next.LogsEndpoint, collect)
	tm.restartCollector()
}

// unpark restores our own exporters when collection was turned back on while
// this runner sat idle.
func (tm *TelemetryManager) unpark() {
	tm.mu.Lock()
	if !tm.parked {
		tm.mu.Unlock()
		return
	}
	tm.parked = false
	tm.collect = true
	tm.mu.Unlock()

	log.Logger().Infof("Collection to WarpBuild is enabled again; resuming the collector")
	tm.restartCollector()
}

// park drops our own exporters while we wait to learn the export destination.
func (tm *TelemetryManager) park() {
	tm.mu.Lock()
	if tm.parked {
		tm.mu.Unlock()
		return
	}
	tm.parked = true
	tm.collect = false
	tm.mu.Unlock()

	log.Logger().Infof("Collection to WarpBuild is disabled; idling the collector until this runner is allocated")
	tm.restartCollector()
}
