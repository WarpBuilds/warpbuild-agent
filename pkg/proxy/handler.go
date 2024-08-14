package proxy

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type GithubAPIError struct {
	Message   string `json:"message"`
	TypeName  string `json:"typeName"`
	TypeKey   string `json:"typeKey"`
	ErrorCode int    `json:"errorCode"`
}

func sendErrorResponse(c *fiber.Ctx, status int, message, typeName, typeKey string, errorCode int) error {
	errorResponse := GithubAPIError{
		Message:   message,
		TypeName:  typeName,
		TypeKey:   typeKey,
		ErrorCode: errorCode,
	}
	return c.Status(status).JSON(errorResponse)
}

func getAuthorizationToken() string {
	return os.Getenv("WARPBUILD_RUNNER_VERIFICATION_TOKEN")
}

func getAuthorizationTokenFromHeader(c *fiber.Ctx) string {
	// Extract the Authorization header
	// This is removed for the time being as the token passing mechanism in GitHub is very finnicky
	// authHeader := c.Get("Authorization")
	// if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
	// 	return strings.TrimPrefix(authHeader, "Bearer ")
	// }

	// Fall back to environment variable if header is not present
	return getAuthorizationToken()
}

func GetCacheEntryHandler(c *fiber.Ctx) error {
	queryKeys := c.Query("keys")
	version := c.Query("version")

	if queryKeys == "" || version == "" {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Keys and version are required.", "InvalidRequest", "InvalidRequest", 1001)
	}

	keys := strings.Split(queryKeys, ",")
	if len(keys) == 0 {
		return sendErrorResponse(c, fiber.StatusBadRequest, "No keys provided.", "InvalidRequest", "InvalidRequest", 1002)
	}

	resp, err := GetCache(c.Context(), DockerGHAGetCacheRequest{Keys: keys, Version: version, AuthToken: getAuthorizationTokenFromHeader(c)})
	if err != nil {
		fmt.Printf("Error getting cache: %v\n", err)
		// GHA backend expects a 200 response even if the cache is not found. It checks if the cache key is empty.
		return c.JSON(DockerGHAGetCacheResponse{CacheKey: "", ArchiveLocation: ""})
	}

	return c.JSON(resp)
}

func ReserveCacheHandler(c *fiber.Ctx) error {
	var req DockerGHAReserveCacheRequest
	if err := c.BodyParser(&req); err != nil {
		fmt.Printf("Error parsing request: %v\n", err)
		return sendErrorResponse(c, fiber.StatusBadRequest, "Failed to parse request body.", "InvalidRequest", "InvalidRequest", 2001)
	}

	resp, err := ReserveCache(c.Context(), DockerGHAReserveCacheRequest{Key: req.Key, Version: req.Version, AuthToken: getAuthorizationTokenFromHeader(c)})
	if err != nil {
		fmt.Printf("Error reserving cache: %v\n", err)
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Failed to reserve cache.", "CacheReserveFailed", "CacheReserveFailed", 2002)
	}

	return c.JSON(resp)
}

func UploadCacheHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		fmt.Printf("Invalid cache ID: %v\n", err)
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid cache ID.", "InvalidCacheID", "InvalidCacheID", 3001)
	}

	fmt.Printf("Received UploadCache request for Cache ID: %d\n", id)
	contentRange := c.Get("Content-Range")
	if contentRange == "" {
		return sendErrorResponse(c, fiber.StatusBadRequest, "Content-Range header is missing.", "MissingHeader", "MissingHeader", 3002)
	}

	resp, err := UploadCache(c.Context(), DockerGHAUploadCacheRequest{CacheID: id, Content: c.Body(), ContentRange: contentRange, AuthToken: getAuthorizationTokenFromHeader(c)})
	if err != nil {
		fmt.Printf("Error uploading cache: %v\n", err)
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Failed to upload cache.", "CacheUploadFailed", "CacheUploadFailed", 3003)
	}

	return c.JSON(resp)
}

func CommitCacheHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		fmt.Printf("Invalid cache ID: %v\n", err)
		return sendErrorResponse(c, fiber.StatusBadRequest, "Invalid cache ID.", "InvalidCacheID", "InvalidCacheID", 4001)
	}

	resp, err := CommitCache(c.Context(), DockerGHACommitCacheRequest{CacheID: id, AuthToken: getAuthorizationTokenFromHeader(c)})
	if err != nil {
		fmt.Printf("Error committing cache: %v\n", err)
		return sendErrorResponse(c, fiber.StatusInternalServerError, "Failed to commit cache.", "CacheCommitFailed", "CacheCommitFailed", 4003)
	}

	return c.JSON(resp)

}
