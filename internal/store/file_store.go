package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// FileStore persists the manifest as a JSON file. It satisfies Store.
type FileStore struct {
	path string
}

// NewFileStore returns a FileStore backed by the file at path.
func NewFileStore(path string) *FileStore { return &FileStore{path: path} }

// Load reads the manifest, returning an empty one if the file is absent.
func (s *FileStore) Load() (*Manifest, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, fs.ErrNotExist) {
		return &Manifest{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read manifest: %w", err)
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}
	return &m, nil
}

// Save writes the manifest atomically as indented JSON.
func (s *FileStore) Save(m *Manifest) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("create manifest dir: %w", err)
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("encode manifest: %w", err)
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write manifest: %w", err)
	}
	if err := os.Rename(tmp, s.path); err != nil {
		return fmt.Errorf("replace manifest: %w", err)
	}
	return nil
}

var _ Store = (*FileStore)(nil)
