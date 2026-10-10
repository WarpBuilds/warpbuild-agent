//go:build darwin

package sandbox

import "errors"

var (
	errProcessNotFound   = errors.New("process not found")
	errUnsupportedSignal = errors.New("only SIGTERM and SIGKILL are supported")
	errNotADirectory     = errors.New("path is not a directory")
	errIsADirectory      = errors.New("path is a directory")
	errMethodNotAllowed  = errors.New("method not allowed")
	errEmptyInput        = errors.New("input carries neither stdin nor pty bytes")
	errInputBeforeStart  = errors.New("input data arrived before the start event")
	errInputClosed       = errors.New("process input closed")
)
