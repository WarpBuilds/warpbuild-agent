//go:build windows
// +build windows

package manager

import (
	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

// logToSyslog logs a message using Windows Event Log or falls back to regular logging
func (a *agentImpl) logToSyslog(priority int, message string) {
	// Windows doesn't have syslog, just use regular logging
	// In a production system, you might want to use Windows Event Log here
	switch priority {
	case logInfo:
		log.Logger().Infof("[SYSLOG] %s", message)
	case logWarning:
		log.Logger().Warnf("[SYSLOG] %s", message)
	case logErr:
		log.Logger().Errorf("[SYSLOG] %s", message)
	}
}
