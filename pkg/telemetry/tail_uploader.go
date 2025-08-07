package telemetry

// TODO: remove this package if not required

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

type upld struct {
	interval    time.Duration
	baseDir     string
	fpath       string
	processChan chan<- string
}

func New() *upld {
	return &upld{}
}

func (u *upld) WithInterval(t time.Duration) *upld {
	u.interval = t
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

func (u *upld) Do(ctx context.Context) error {

	handlePanic := func() {
		if r := recover(); r != nil {
			log.Logger().Errorf("Recovered from panic: %v", r)
		}
	}

	defer handlePanic()

	if u.interval == 0 {
		return errors.New("invalid interval '0' passed")
	}

	u.processChan = make(chan<- string)

	go u.tailLines(ctx)
	go u.upload(ctx)

	<-ctx.Done()

	return nil
}

func (u *upld) tailLines(ctx context.Context) error {
	interval := u.interval
	path := filepath.Join(u.baseDir, u.fpath)
	processChan := u.processChan

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var offset int64
	remainder := make([]byte, 0, 4096)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			fi, err := f.Stat()
			if err != nil {
				return err
			}
			// nothing new?
			if fi.Size() <= offset {
				continue
			}

			// read just the new bytes
			toRead := fi.Size() - offset
			chunk := make([]byte, toRead)
			n, err := f.ReadAt(chunk, offset)
			if err != nil && err != io.EOF {
				return err
			}
			offset += int64(n)

			// prepend any leftover from last time
			data := append(remainder, chunk[:n]...)

			// split on newline
			parts := bytes.Split(data, []byte("\n"))

			fullLines := ""
			// all except last are full lines
			for i := 0; i < len(parts)-1; i++ {
				fullLines += string(parts[i]) + "\n"
			}

			processChan <- fullLines

			// last element may be a partial line—keep for next round
			remainder = parts[len(parts)-1]
		}
	}
}

func (u *upld) upload(ctx context.Context) {

}
