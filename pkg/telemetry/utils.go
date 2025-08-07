package telemetry

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
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

	// Increase buffer size to handle very long lines (up to 1MB)
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 1024*1024)

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

	log.Logger().Infof("Fetching presigned URL for runner ID: %s and polling secret %s from url %v and client [%+v]", runnerId, pollingSecret, client.GetConfig().Host, client)

	logFileName := fmt.Sprintf("%s.log", time.Now().Format("20060102-150405"))
	out, resp, err := client.V1RunnerInstanceAPI.
		GetRunnerInstancePresignedLogUploadURL(context.Background(), runnerId).
		XPOLLINGSECRET(pollingSecret).
		LogFileName(logFileName).
		Execute()
	if err != nil {
		return "", err
	}

	// Print the response body
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Logger().Errorf("Error reading response body: %v", err)
	}

	fmt.Printf("Response body: %s\n", body)

	// You can also print other details like status code and headers if needed
	fmt.Printf("Status code: %d\n", resp.StatusCode)
	fmt.Printf("Response headers: %v\n", resp.Header)

	if out == nil || out.Url == nil {
		return "", fmt.Errorf("no url received in response")
	}

	return *out.Url, nil
}
