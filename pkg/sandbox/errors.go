//go:build darwin

package sandbox

import "errors"

var (
	errProcessNotFound   = errors.New("process not found")
	errProcessGone       = errors.New("process exited before the stream was attached")
	errUnsupportedSignal = errors.New("only SIGTERM and SIGKILL are supported")
	errNotADirectory     = errors.New("path is not a directory")
	errIsADirectory      = errors.New("path is a directory")
	errMethodNotAllowed  = errors.New("method not allowed")
)
