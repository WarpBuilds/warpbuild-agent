package manager

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/warpbuilds/warpbuild-agent/pkg/log"
)

const outputsUploadTimeout = 5 * time.Minute

func UploadSessionOutputs(ctx context.Context, outputsDir, hostURL, pollingSecret, runnerInstanceID string) {
	if outputsDir == "" || hostURL == "" || runnerInstanceID == "" {
		return
	}
	entries, err := os.ReadDir(outputsDir)
	if err != nil {
		log.Logger().Infof("[claude_outputs] no outputs dir %s to upload (%v)", outputsDir, err)
		return
	}
	if len(entries) == 0 {
		log.Logger().Infof("[claude_outputs] outputs dir %s is empty, skipping upload", outputsDir)
		return
	}

	uploadCtx, cancel := context.WithTimeout(ctx, outputsUploadTimeout)
	defer cancel()

	archivePath, err := tarGzDir(uploadCtx, outputsDir)
	if err != nil {
		log.Logger().Errorf("[claude_outputs] failed to archive %s: %v", outputsDir, err)
		return
	}
	defer os.Remove(archivePath)

	url, err := requestOutputUploadURL(uploadCtx, hostURL, pollingSecret, runnerInstanceID)
	if err != nil {
		log.Logger().Errorf("[claude_outputs] failed to get output upload URL: %v", err)
		return
	}
	if err := putFile(uploadCtx, url, archivePath); err != nil {
		log.Logger().Errorf("[claude_outputs] failed to upload session outputs: %v", err)
		return
	}
	log.Logger().Infof("[claude_outputs] uploaded session outputs from %s", outputsDir)
}

func tarGzDir(ctx context.Context, dir string) (string, error) {
	f, err := os.CreateTemp("", "warpbuild-outputs-*.tar.gz")
	if err != nil {
		return "", err
	}
	path := f.Name()

	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)

	walkErr := filepath.Walk(dir, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		rel, err := filepath.Rel(dir, file)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		link := ""
		if fi.Mode()&os.ModeSymlink != 0 {
			if link, err = os.Readlink(file); err != nil {
				return err
			}
		}
		hdr, err := tar.FileInfoHeader(fi, link)
		if err != nil {
			return err
		}
		hdr.Name = filepath.ToSlash(rel)
		if fi.IsDir() {
			hdr.Name += "/"
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if !fi.Mode().IsRegular() {
			return nil
		}
		src, err := os.Open(file)
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, src)
		_ = src.Close()
		return err
	})

	if err := errors.Join(walkErr, tw.Close(), gz.Close(), f.Close()); err != nil {
		_ = os.Remove(path)
		return "", fmt.Errorf("archive %s: %w", dir, err)
	}
	return path, nil
}

func requestOutputUploadURL(ctx context.Context, hostURL, pollingSecret, runnerInstanceID string) (string, error) {
	body, err := json.Marshal(map[string]string{"runner_instance_id": runnerInstanceID})
	if err != nil {
		return "", err
	}
	endpoint := strings.TrimRight(hostURL, "/") + "/sandboxes/outputs/upload-url"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("X-POLLING-SECRET", pollingSecret)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload-url returned %s: %s", resp.Status, string(b))
	}
	var out struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if out.URL == "" {
		return "", fmt.Errorf("empty upload url in response")
	}
	return out.URL, nil
}

func putFile(ctx context.Context, url, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, f)
	if err != nil {
		return err
	}
	req.ContentLength = info.Size()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("PUT returned %s: %s", resp.Status, string(b))
	}
	return nil
}
