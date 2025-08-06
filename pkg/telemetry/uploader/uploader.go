package uploader

// TODO: remove this package if not required

import (
	"context"
	"sync"
	"time"
)

type upld struct {
	debounceDelay time.Duration
	baseDir       string
	fpath         string
	debouceTimer  sync.Mutex
	uploadMu      sync.Mutex
}

func New() *upld {
	return &upld{}
}

func (u *upld) WithDebounceDelay(t time.Duration) *upld {
	u.debounceDelay = t
	return u
}

func (u *upld) WithBaseDirectory(dir string) *upld {
	u.baseDir = dir
	return u
}

// Sets the file to upload. This is only the base
// file name not the complete path to the file.
//
// Use WithBaseDirectory to mark the directory for the
// file instead.
func (u *upld) WithFile(f string) *upld {
	u.fpath = f
	return u
}

func (u *upld) Do(ctx context.Context)
