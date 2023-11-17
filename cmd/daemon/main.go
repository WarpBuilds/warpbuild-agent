package main

import (
	"context"
	"log"
	"os"

	"github.com/warpbuilds/warpbuild-agent/pkg/app"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	settingsFile := os.Getenv("WARPBUILD_AGENTD_SETTINGS_FILE")

	err := app.NewApp(ctx, &app.ApplicationOptions{
		SettingsFile: settingsFile,
	})
	if err != nil {
		return err
	}

	return nil
}
