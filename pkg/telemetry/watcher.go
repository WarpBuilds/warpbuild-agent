package telemetry

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/fsnotify/fsnotify"
	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

var watcher *fsnotify.Watcher

func disableOtelOutputFileWatcher() {
	watcher.Close()
}

func enableOtelOutputFileWatcher(ctx context.Context, baseDirectory string) error {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %w", err)
	}

	go func() {
		defer handlePanic()
		watchOtelOutputFile(ctx, baseDirectory)
	}()
	return nil
}

func watchOtelOutputFile(ctx context.Context, baseDirectory string) {
	defer watcher.Close()

	// Ensure the log file exists
	if _, err := os.Stat(getOtelCollectorOutputFilePath(baseDirectory, runtime.GOOS)); os.IsNotExist(err) {
		file, err := os.Create(getOtelCollectorOutputFilePath(baseDirectory, runtime.GOOS))
		if err != nil {
			log.Logger().Errorf("failed to create log file: %v", err)
		}
		file.Close()
	}

	watchPath := getOtelCollectorOutputFilePath(baseDirectory, runtime.GOOS)
	if runtime.GOOS == "windows" {
		// ? on windows we try to watch the path for the directory
		// ? since apparently the events are only transmitted on dir level
		// ? not on path level.
		//
		// TODO: validate that the above statement is true.
		watchPath = filepath.Dir(watchPath)
	}

	err := watcher.Add(watchPath)
	if err != nil {
		log.Logger().Errorf("failed to watch file: %v", err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Logger().Infof("Unknown event")
				return
			}
			if event.Op == fsnotify.Write {
				log.Logger().Infof("Watching the following paths: %+v", watcher.WatchList())
				log.Logger().Infof("Modified file:", event.Name)
				debouncedOtelUpload(ctx, baseDirectory)
				// TODO: remove below log
				log.Logger().Infof("Completed upload for event:", event.Name)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Logger().Errorf("Error watching file:", err)
		}
	}
}
