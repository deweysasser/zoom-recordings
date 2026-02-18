package zoom

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadFile downloads a recording file to the specified destination path.
// It uses the access_token query parameter as required by the Zoom download API.
// The download is atomic: it writes to a temp file first, then renames.
func (c *Client) DownloadFile(ctx context.Context, downloadURL, accessToken, destPath string) error {
	reqURL := fmt.Sprintf("%s?access_token=%s", downloadURL, accessToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return fmt.Errorf("creating download request: %w", err)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("downloading file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	tmpPath := destPath + ".tmp"
	f, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("writing file: %w", err)
	}

	if err := f.Close(); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("closing temp file: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("creating destination directory: %w", err)
	}

	if err := os.Rename(tmpPath, destPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("moving file to destination: %w", err)
	}

	return nil
}
