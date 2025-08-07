package telemetry

// var watcher *fsnotify.Watcher

// func disableOtelOutputFileWatcher() {
// 	watcher.Close()
// }

// func enableOtelOutputFileWatcher(ctx context.Context, baseDirectory string) error {
// 	var err error
// 	watcher, err = fsnotify.NewWatcher()
// 	if err != nil {
// 		return fmt.Errorf("failed to create file watcher: %w", err)
// 	}

// 	go func() {
// 		defer handlePanic()
// 		watchOtelOutputFile(ctx, baseDirectory)
// 	}()
// 	return nil
// }

// func watchOtelOutputFile(ctx context.Context, baseDirectory string) {
// 	defer watcher.Close()

// 	// Ensure the log file exists
// 	if _, err := os.Stat(getOtelCollectorOutputFilePath(baseDirectory, runtime.GOOS, false)); os.IsNotExist(err) {
// 		file, err := os.Create(getOtelCollectorOutputFilePath(baseDirectory, runtime.GOOS, false))
// 		if err != nil {
// 			log.Logger().Errorf("failed to create log file: %v", err)
// 		}
// 		file.Close()
// 	}

// 	// Ensure the metrics file exists
// 	if _, err := os.Stat(getOtelCollectorOutputFilePath(baseDirectory, runtime.GOOS, true)); os.IsNotExist(err) {
// 		file, err := os.Create(getOtelCollectorOutputFilePath(baseDirectory, runtime.GOOS, true))
// 		if err != nil {
// 			log.Logger().Errorf("failed to create log file: %v", err)
// 		}
// 		file.Close()
// 	}

// 	// We watch the directory for fs notification events
// 	watchPath := getOtelCollectorOutputFilePath(baseDirectory, runtime.GOOS, false)
// 	watchPath = filepath.Dir(watchPath)

// 	err := watcher.Add(watchPath)
// 	if err != nil {
// 		log.Logger().Errorf("failed to watch file: %v", err)
// 	}

// 	for {
// 		select {
// 		case event, ok := <-watcher.Events:
// 			if !ok {
// 				log.Logger().Infof("Unknown event")
// 				return
// 			}
// 			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename|fsnotify.Chmod) != 0 {
// 				log.Logger().Infof("Watching the following paths: %+v", watcher.WatchList())
// 				log.Logger().Infof("Modified file:", event.Name)

// 				metricsPath := filepath.Base(getOtelCollectorOutputFilePath(baseDirectory, runtime.GOOS, true))
// 				if strings.Contains(event.Name, metricsPath) {
// 					log.Logger().Infof("Path is metrics path: %v", metricsPath)
// 					debouncedOtelUpload(ctx, baseDirectory, true)
// 				}

// 				logsPath := filepath.Base(getOtelCollectorOutputFilePath(baseDirectory, runtime.GOOS, false))
// 				if strings.Contains(event.Name, logsPath) {
// 					log.Logger().Infof("Path is logs path: %v", logsPath)
// 					debouncedOtelUpload(ctx, baseDirectory, false)
// 				}

// 				// TODO: remove below log
// 				log.Logger().Infof("Completed upload for event:", event.Name)
// 			}
// 		case err, ok := <-watcher.Errors:
// 			if !ok {
// 				return
// 			}
// 			log.Logger().Errorf("Error watching file:", err)
// 		}
// 	}
// }
