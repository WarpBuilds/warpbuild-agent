package proxy

import (
	"context"
	"fmt"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"

	"github.com/gofiber/fiber/v2"
)

const PROXY_SERVER_OPTIONS_CONTEXT_KEY = "proxy_server_options"

type ProxyServerOptions struct {
	CacheProxyPort                   string
	CacheBackendHost                 string
	WarpBuildRunnerVerificationToken string
}

func StartProxyServer(ctx context.Context, opts *ProxyServerOptions) error {
	if opts.CacheProxyPort == "" {
		opts.CacheProxyPort = "49160"
	}

	if opts.CacheBackendHost == "" {
		opts.CacheBackendHost = "https://cache.warpbuild.com"
	}

	if opts.WarpBuildRunnerVerificationToken == "" {
		log.Logger().Errorf("WARPBUILD_RUNNER_VERIFICATION_TOKEN is required")
		return fmt.Errorf("WARPBUILD_RUNNER_VERIFICATION_TOKEN is required")
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 1024 * 1024 * 1024, // 1GB limit for body.
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Locals(PROXY_SERVER_OPTIONS_CONTEXT_KEY, opts)
		return c.Next()
	})

	registerRoutes(app)

	log.Logger().Infof("Starting cache proxy server on port %s", opts.CacheProxyPort)
	err := app.Listen(":" + opts.CacheProxyPort)
	if err != nil {
		log.Logger().Errorf("Failed to start cache proxy server: %v", err)
		return err
	}

	return nil
}

func registerRoutes(app *fiber.App) {
	api := app.Group("/_apis/artifactcache")
	{
		api.Get("/cache", GetCacheEntryHandler)
		api.Post("/caches", ReserveCacheHandler)
		api.Patch("/caches/:id", UploadCacheHandler)
		api.Post("/caches/:id", CommitCacheHandler)
	}

}
