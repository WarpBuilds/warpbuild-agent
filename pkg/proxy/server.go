package proxy

import (
	"context"
	"fmt"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"

	"github.com/gofiber/fiber/v2"
)

const PROXY_SERVER_OPTIONS_CONTEXT_KEY = "proxy_server_options"

type ProxyServerOptions struct {
	Port string
	CacheBackendInfo
}

func StartProxyServer(ctx context.Context, opts *ProxyServerOptions) error {
	if opts.Port == "" {
		opts.Port = "49160"
	}

	if opts.CacheBackendInfo.HostURL == "" {
		opts.CacheBackendInfo.HostURL = "https://cache.warpbuild.com"
	}

	if opts.CacheBackendInfo.AuthToken == "" {
		log.Logger().Errorf("WARPBUILD_RUNNER_VERIFICATION_TOKEN is required")
		return fmt.Errorf("WARPBUILD_RUNNER_VERIFICATION_TOKEN is required")
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 1024 * 1024 * 1024 * 10, // 10GB limit for body.
		// Increase header size limits for BuildKit requests
		// ReadBufferSize:  16 * 1024, // 16KB read buffer (default is 4KB)
		// WriteBufferSize: 16 * 1024, // 16KB write buffer (default is 4KB)
		ReadBufferSize:  32 * 1024, // 32KB
		WriteBufferSize: 32 * 1024, // 32KB
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Locals(PROXY_SERVER_OPTIONS_CONTEXT_KEY, opts)
		return c.Next()
	})

	registerRoutes(app)

	log.Logger().Infof("Starting cache proxy server on port %s", opts.Port)
	err := app.Listen(":" + opts.Port)
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
