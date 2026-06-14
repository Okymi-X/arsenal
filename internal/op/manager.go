package op

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

// Manager persists ops and their lockfiles under a single ops directory.
type Manager struct {
	dir string
}

// NewManager returns a Manager rooted at the given ops directory.
func NewManager(dir string) *Manager { return &Manager{dir: dir} }

// Path returns the op definition file path for a name.
func (m *Manager) Path(name string) string {
	return filepath.Join(m.dir, name+".toml")
}

// LockPath returns the lockfile path for an op name.
func (m *Manager) LockPath(name string) string {
	return filepath.Join(m.dir, name+".lock.toml")
}

// Create writes a new, empty op, refusing to overwrite an existing one.
func (m *Manager) Create(o *Op) error {
	if o.Name == "" {
		return fmt.Errorf("op name is required")
	}
	if _, err := os.Stat(m.Path(o.Name)); err == nil {
		return fmt.Errorf("op %q already exists", o.Name)
	}
	return m.Save(o)
}

// Save writes an op definition to disk as TOML, atomically.
func (m *Manager) Save(o *Op) error {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(o); err != nil {
		return fmt.Errorf("encode op: %w", err)
	}
	if err := os.MkdirAll(m.dir, 0o755); err != nil {
		return fmt.Errorf("create ops dir: %w", err)
	}
	tmp := m.Path(o.Name) + ".tmp"
	if err := os.WriteFile(tmp, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write op: %w", err)
	}
	if err := os.Rename(tmp, m.Path(o.Name)); err != nil {
		return fmt.Errorf("replace op: %w", err)
	}
	return nil
}

// Load reads an op definition by name.
func (m *Manager) Load(name string) (*Op, error) {
	data, err := os.ReadFile(m.Path(name))
	if err != nil {
		return nil, fmt.Errorf("read op %q: %w", name, err)
	}
	var o Op
	if err := toml.Unmarshal(data, &o); err != nil {
		return nil, fmt.Errorf("parse op %q: %w", name, err)
	}
	return &o, nil
}

// List returns the names of all defined ops, sorted.
func (m *Manager) List() ([]string, error) {
	entries, err := os.ReadDir(m.dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read ops dir: %w", err)
	}
	var names []string
	for _, e := range entries {
		n := e.Name()
		if e.IsDir() || !strings.HasSuffix(n, ".toml") || strings.HasSuffix(n, ".lock.toml") {
			continue
		}
		names = append(names, strings.TrimSuffix(n, ".toml"))
	}
	sort.Strings(names)
	return names, nil
}
