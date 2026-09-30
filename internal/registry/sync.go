package registry

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Okymi-X/arsenal/internal/fsutil"
)

// FileSource loads a registry from a local TOML file and refreshes it from a
// remote manifest or legacy monolithic URL. It is the default Source used by
// the CLI.
type FileSource struct {
	path string
	url  string
	http *http.Client
}

// NewFileSource builds a Source backed by a local file and a remote URL.
func NewFileSource(path, url string) *FileSource {
	return &FileSource{
		path: path,
		url:  url,
		http: &http.Client{Timeout: 30 * time.Second},
	}
}

// Load reads and parses the local registry file.
func (s *FileSource) Load() (*Registry, error) {
	return Load(s.path)
}

// Sync downloads the remote registry, validates it, and replaces the local
// copy atomically. The previous file is left untouched on any failure.
func (s *FileSource) Sync(ctx context.Context) error {
	if s.url == "" {
		return fmt.Errorf("no registry URL configured")
	}
	data, err := s.fetchRegistry(ctx)
	if err != nil {
		return err
	}
	if _, err := Parse(data); err != nil {
		return fmt.Errorf("remote registry invalid: %w", err)
	}
	return s.writeAtomic(data)
}

func (s *FileSource) writeAtomic(data []byte) error {
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create registry dir: %w", err)
	}
	if err := fsutil.WriteFileAtomic(s.path, data, 0o644); err != nil {
		return fmt.Errorf("write registry: %w", err)
	}
	return nil
}
