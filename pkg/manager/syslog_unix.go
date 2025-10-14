//go:build !windows
// +build !windows

package manager

import (
	"log/syslog"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

// logToSyslog logs a message to syslog for debugging (Unix systems only)
func (a *agentImpl) logToSyslog(priority int, message string) {
	syslogger, err := syslog.New(syslog.Priority(priority), "warpbuild-agent")
	if err != nil {
		log.Logger().Warnf("Failed to connect to syslog: %v", err)
		return
	}
	defer syslogger.Close()

	switch priority {
	case logInfo:
		syslogger.Info(message)
	case logWarning:
		syslogger.Warning(message)
	case logErr:
		syslogger.Err(message)
	}
}
