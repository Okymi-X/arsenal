package registry

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// FileSource loads a registry from a local TOML file and refreshes it from a
// remote URL. It is the default Source used by the CLI.
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
func (s *FileSource) Sync() error {
	if s.url == "" {
		return fmt.Errorf("no registry URL configured")
	}
	data, err := s.fetch()
	if err != nil {
		return err
	}
	if _, err := Parse(data); err != nil {
		return fmt.Errorf("remote registry invalid: %w", err)
	}
	return s.writeAtomic(data)
}

func (s *FileSource) fetch() ([]byte, error) {
	resp, err := s.http.Get(s.url)
	if err != nil {
		return nil, fmt.Errorf("fetch registry: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch registry: status %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read registry response: %w", err)
	}
	return data, nil
}

func (s *FileSource) writeAtomic(data []byte) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("create registry dir: %w", err)
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write registry: %w", err)
	}
	if err := os.Rename(tmp, s.path); err != nil {
		return fmt.Errorf("replace registry: %w", err)
	}
	return nil
}
