package telemetry

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/warpbuild"
)

func readLastNLines(filePath string, n int) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > n {
			lines = lines[1:]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return []byte(strings.Join(lines, "\n")), nil
}

func fetchPresignedURL(ctx context.Context) (string, error) {
	client := ctx.Value(WarpBuildAgentContextKey).(*warpbuild.APIClient)
	runnerId := ctx.Value(WarpBuildRunnerIDContextKey).(string)
	pollingSecret := ctx.Value(WarpBuildRunnerPollingSecretContextKey).(string)

	logFileName := fmt.Sprintf("%s.log", time.Now().Format("20060102-150405"))
	_, resp, err := client.V1RunnerInstanceApi.
		GetRunnerInstancePresignedLogUploadURL(context.Background(), runnerId).
		XPOLLINGSECRET(pollingSecret).
		LogFileName(logFileName).
		Execute()
	if err != nil {
		return "", fmt.Errorf("failed to fetch presigned URL: %w", err)
	}

	var url string
	err = json.NewDecoder(resp.Body).Decode(&url)
	if err != nil {
		return "", fmt.Errorf("failed to decode presigned URL: %w", err)
	}

	return url, nil
}
