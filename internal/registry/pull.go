package registry

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

const (
	maxRemoteRegistryFile = 4 << 20
	pullWorkers           = 4
)

func (s *FileSource) fetchRegistry(parent context.Context) ([]byte, error) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	data, err := s.fetchURL(ctx, s.url)
	if err != nil {
		return nil, err
	}
	manifest, err := ParseManifest(data)
	if err != nil || len(manifest.Segments) == 0 {
		return Expand(data, nil)
	}

	segments := make(map[string][]byte, len(manifest.Segments))
	jobs := make(chan string)
	errCh := make(chan error, 1)
	var mu sync.Mutex
	totalSize := 0
	var workers sync.WaitGroup
	for range min(pullWorkers, len(manifest.Segments)) {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for name := range jobs {
				segmentURL, resolveErr := resolveSegmentURL(s.url, name)
				if resolveErr != nil {
					reportPullError(errCh, cancel, resolveErr)
					return
				}
				segment, fetchErr := s.fetchURL(ctx, segmentURL)
				if fetchErr != nil {
					reportPullError(errCh, cancel, fetchErr)
					return
				}
				if verifyErr := verifySegment(manifest, name, segment); verifyErr != nil {
					reportPullError(errCh, cancel, verifyErr)
					return
				}
				mu.Lock()
				if totalSize+len(segment) > maxAssembledSize {
					mu.Unlock()
					reportPullError(errCh, cancel, fmt.Errorf("remote registry segments exceed %d bytes", maxAssembledSize))
					return
				}
				totalSize += len(segment)
				segments[name] = segment
				mu.Unlock()
			}
		}()
	}
dispatch:
	for _, name := range manifest.Segments {
		select {
		case jobs <- name:
		case <-ctx.Done():
			break dispatch
		}
	}
	close(jobs)
	workers.Wait()
	select {
	case err := <-errCh:
		return nil, err
	default:
	}
	return Expand(data, func(name string) ([]byte, error) {
		segment, ok := segments[name]
		if !ok {
			return nil, fmt.Errorf("registry segment %q was not downloaded", name)
		}
		return segment, nil
	})
}

func reportPullError(errCh chan<- error, cancel context.CancelFunc, err error) {
	select {
	case errCh <- err:
		cancel()
	default:
	}
}

func (s *FileSource) fetchURL(ctx context.Context, rawURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build registry request: %w", err)
	}
	req.Header.Set("User-Agent", "arsenal")
	resp, err := s.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", rawURL, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if !sameOrigin(req.URL, resp.Request.URL) {
		return nil, fmt.Errorf("fetch %s: redirect changed remote origin", rawURL)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch %s: status %d", rawURL, resp.StatusCode)
	}
	if resp.ContentLength > maxRemoteRegistryFile {
		return nil, fmt.Errorf("fetch %s: response exceeds %d bytes", rawURL, maxRemoteRegistryFile)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxRemoteRegistryFile+1))
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", rawURL, err)
	}
	if len(data) > maxRemoteRegistryFile {
		return nil, fmt.Errorf("fetch %s: response exceeds %d bytes", rawURL, maxRemoteRegistryFile)
	}
	return data, nil
}

func resolveSegmentURL(manifestURL, name string) (string, error) {
	base, err := url.Parse(manifestURL)
	if err != nil {
		return "", fmt.Errorf("parse registry URL: %w", err)
	}
	reference, err := url.Parse(name)
	if err != nil {
		return "", fmt.Errorf("parse registry segment URL: %w", err)
	}
	resolved := base.ResolveReference(reference)
	if !sameOrigin(base, resolved) {
		return "", fmt.Errorf("registry segment %q changes remote origin", name)
	}
	return resolved.String(), nil
}

func sameOrigin(a, b *url.URL) bool {
	return a != nil && b != nil && a.Scheme == b.Scheme && a.Host == b.Host
}
