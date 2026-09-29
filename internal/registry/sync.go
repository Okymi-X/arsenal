package registry

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
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
func (s *FileSource) Sync() error {
	if s.url == "" {
		return fmt.Errorf("no registry URL configured")
	}
	data, err := s.fetchRegistry()
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
	tmp, err := os.CreateTemp(dir, ".registry-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary registry: %w", err)
	}
	tmpPath := tmp.Name()
	defer func() { _ = os.Remove(tmpPath) }()
	if err := writeRegistryTemp(tmp, data); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, s.path); err != nil {
		return fmt.Errorf("replace registry: %w", err)
	}
	return nil
}

func writeRegistryTemp(file *os.File, data []byte) error {
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		return fmt.Errorf("write temporary registry: %w", err)
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return fmt.Errorf("sync temporary registry: %w", err)
	}
	if err := file.Chmod(0o644); err != nil {
		_ = file.Close()
		return fmt.Errorf("set registry permissions: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close temporary registry: %w", err)
	}
	return nil
}
