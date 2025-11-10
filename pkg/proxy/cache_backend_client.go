package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

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
	requestURL, err := url.JoinPath(info.HostURL, "/v1/cache", req.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	f := fiber.Post(requestURL).
		ContentType("application/json").
		Add("Accept", "application/json").
		Add("Authorization", fmt.Sprintf("Bearer %s", info.AuthToken))

	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		f = f.Body(bodyBytes)
		fmt.Printf("\tPayload: %s\n", string(bodyBytes))
	}

	statusCode, body, errs := f.Bytes()
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
