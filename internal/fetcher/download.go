package fetcher

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func (f *Fetcher) download(ctx context.Context, rawURL, destPath string, expectedSize int64) (int64, error) {
	if f.maxDownloadBytes <= 0 {
		return 0, fmt.Errorf("asset download limit is not configured")
	}
	if expectedSize < 0 {
		return 0, fmt.Errorf("asset metadata has invalid negative size %d", expectedSize)
	}
	if expectedSize > f.maxDownloadBytes {
		return 0, fmt.Errorf("asset metadata size %d exceeds %d-byte limit", expectedSize, f.maxDownloadBytes)
	}
	if f.validateDownloadURL != nil {
		if err := f.validateDownloadURL(rawURL); err != nil {
			return 0, err
		}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return 0, fmt.Errorf("build download request: %w", err)
	}
	req.Header.Set("User-Agent", "arsenal")
	resp, err := f.client.Do(req)
	if err != nil {
		var urlError *url.Error
		if errors.As(err, &urlError) {
			err = urlError.Err
		}
		return 0, fmt.Errorf("download request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("download request: %s", resp.Status)
	}
	if resp.ContentLength > f.maxDownloadBytes {
		return 0, fmt.Errorf("asset response size %d exceeds %d-byte limit", resp.ContentLength, f.maxDownloadBytes)
	}
	return f.writeDownload(resp.Body, destPath, expectedSize)
}

func (f *Fetcher) writeDownload(body io.Reader, destPath string, expectedSize int64) (int64, error) {
	out, err := os.CreateTemp(filepath.Dir(destPath), "."+filepath.Base(destPath)+".part-*")
	if err != nil {
		return 0, fmt.Errorf("create temporary asset: %w", err)
	}
	tmp := out.Name()
	defer func() {
		_ = out.Close()
		_ = os.Remove(tmp)
	}()
	if err := out.Chmod(0o755); err != nil {
		return 0, fmt.Errorf("set temporary asset permissions: %w", err)
	}
	n, err := io.Copy(out, io.LimitReader(body, f.maxDownloadBytes+1))
	if err != nil {
		return 0, fmt.Errorf("write temporary asset: %w", err)
	}
	if n > f.maxDownloadBytes {
		return 0, fmt.Errorf("asset exceeds %d-byte limit", f.maxDownloadBytes)
	}
	if n != expectedSize {
		return 0, fmt.Errorf("asset size mismatch: got %d bytes, expected %d", n, expectedSize)
	}
	if err := out.Sync(); err != nil {
		return 0, fmt.Errorf("sync temporary asset: %w", err)
	}
	if err := out.Close(); err != nil {
		return 0, fmt.Errorf("close temporary asset: %w", err)
	}
	if err := os.Rename(tmp, destPath); err != nil {
		return 0, fmt.Errorf("finalize asset: %w", err)
	}
	return n, nil
}
