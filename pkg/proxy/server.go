package proxy

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
)

type ProxyServerOptions struct {
}

func StartProxyServer(ctx context.Context, opts *ProxyServerOptions) error {
	app := fiber.New(fiber.Config{
		BodyLimit: 1024 * 1024 * 1024, // 1GB limit for body.
	})

	registerRoutes(app)

	port := os.Getenv("WARPBUILD_PROXY_PORT")
	if port == "" {
		// Use a rarely used port by default
		port = "49160"
	}

	err := app.Listen(":" + port)
	if err != nil {
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
