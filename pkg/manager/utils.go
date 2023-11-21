package manager

import (
	"bufio"
	"io"
	"os"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

func openFile(filename string) (*os.File, error) {
	if _, err := os.Stat(filename); err == nil {
		// Truncate file if it exists
		return os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0666)
	}
	// Create file if it doesn't exist
	return os.Create(filename)
}

// captureOutput reads from a pipe and sends the output to a channel
func captureOutput(pipe io.ReadCloser, ch chan<- string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		ch <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Logger().Errorf("error reading from pipe: %v", err)
	}
}
