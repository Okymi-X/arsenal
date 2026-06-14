package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

// download streams url into destPath through a temporary ".part" file, renaming
// on success so an interrupted download never leaves a usable-looking file.
// Files land mode 0755 since fetched binaries are typically run on the target.
func (f *Fetcher) download(ctx context.Context, url, destPath string) (int64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", "arsenal")
	resp, err := f.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("download %s: %w", url, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("download %s: %s", url, resp.Status)
	}

	tmp := destPath + ".part"
	out, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o755)
	if err != nil {
		return 0, fmt.Errorf("create %s: %w", tmp, err)
	}
	n, copyErr := io.Copy(out, resp.Body)
	closeErr := out.Close()
	if copyErr != nil {
		_ = os.Remove(tmp)
		return 0, fmt.Errorf("write %s: %w", tmp, copyErr)
	}
	if closeErr != nil {
		_ = os.Remove(tmp)
		return 0, fmt.Errorf("close %s: %w", tmp, closeErr)
	}
	if err := os.Rename(tmp, destPath); err != nil {
		_ = os.Remove(tmp)
		return 0, fmt.Errorf("finalize %s: %w", destPath, err)
	}
	return n, nil
}
