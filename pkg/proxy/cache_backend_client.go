package proxy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type CacheBackendRequest struct {
	Path string
	Body any
}

func getCacheBackendInfo(ctx context.Context) CacheBackendInfo {
	opts := ctx.Value(PROXY_SERVER_OPTIONS_CONTEXT_KEY).(*ProxyServerOptions)
	return opts.CacheBackendInfo
}

func callCacheBackend[T any](ctx context.Context, req CacheBackendRequest) (*T, error) {
	info := getCacheBackendInfo(ctx)
	requestURL := fmt.Sprintf("%s%s", info.HostURL, req.Path)

	var bodyBytes []byte
	var err error
	if req.Body != nil {
		bodyBytes, err = json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	statusCode, body, errs := fiber.Post(requestURL).
		Body(bodyBytes).
		Add("Content-Type", "application/json").
		Add("Accept", "application/json").
		Add("Authorization", fmt.Sprintf("Bearer %s", info.AuthToken)).
		Bytes()

	if len(errs) > 0 {
		return nil, fmt.Errorf("request failed: %v", errs)
	}

	if statusCode < 200 || statusCode >= 300 {
		return nil, fmt.Errorf("request failed with status %d: %s", statusCode, string(body))
	}

	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}
