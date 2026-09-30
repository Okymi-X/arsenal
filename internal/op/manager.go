package op

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/Okymi-X/arsenal/internal/fsutil"
	"github.com/Okymi-X/arsenal/internal/safepath"
	"github.com/Okymi-X/arsenal/internal/strictdecode"
)

// Manager persists ops and their lockfiles under a single ops directory.
type Manager struct {
	dir string
}

// NewManager returns a Manager rooted at the given ops directory.
func NewManager(dir string) *Manager { return &Manager{dir: dir} }

// Path returns the op definition file path for a validated name.
func (m *Manager) Path(name string) (string, error) {
	if err := safepath.ValidateComponent("op name", name); err != nil {
		return "", err
	}
	return filepath.Join(m.dir, name+".toml"), nil
}

// LockPath returns the lockfile path for a validated op name.
func (m *Manager) LockPath(name string) (string, error) {
	if err := safepath.ValidateComponent("op name", name); err != nil {
		return "", err
	}
	return filepath.Join(m.dir, name+".lock.toml"), nil
}

// Create writes a new, empty op, refusing to overwrite an existing one.
func (m *Manager) Create(o *Op) error {
	if err := validateOp(o); err != nil {
		return err
	}
	path, err := m.Path(o.Name)
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("op %q already exists", o.Name)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("inspect op %q: %w", o.Name, err)
	}
	return m.Save(o)
}

// Save writes an op definition to disk as TOML, atomically.
func (m *Manager) Save(o *Op) error {
	if err := validateOp(o); err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(o); err != nil {
		return fmt.Errorf("encode op: %w", err)
	}
	if err := os.MkdirAll(m.dir, 0o755); err != nil {
		return fmt.Errorf("create ops dir: %w", err)
	}
	path, err := m.Path(o.Name)
	if err != nil {
		return err
	}
	if err := fsutil.WriteFileAtomic(path, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write op: %w", err)
	}
	return nil
}

// Load reads an op definition by name.
func (m *Manager) Load(name string) (*Op, error) {
	path, err := m.Path(name)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read op %q: %w", name, err)
	}
	var o Op
	if err := strictdecode.TOML(data, &o); err != nil {
		return nil, fmt.Errorf("parse op %q: %w", name, err)
	}
	if err := validateOp(&o); err != nil {
		return nil, fmt.Errorf("validate op %q: %w", name, err)
	}
	if o.Name != name {
		return nil, fmt.Errorf("op file %q declares name %q", name, o.Name)
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
		name := strings.TrimSuffix(n, ".toml")
		if safepath.ValidateComponent("op name", name) != nil {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}

func validateOp(o *Op) error {
	if o == nil {
		return fmt.Errorf("op is required")
	}
	if err := safepath.ValidateComponent("op name", o.Name); err != nil {
		return err
	}
	seen := make(map[string]struct{}, len(o.Pins))
	for _, pin := range o.Pins {
		if err := safepath.ValidateComponent("tool name", pin.Tool); err != nil {
			return err
		}
		if err := safepath.ValidateComponent("version", pin.Version); err != nil {
			return err
		}
		if _, duplicate := seen[pin.Tool]; duplicate {
			return fmt.Errorf("duplicate pin for tool %q", pin.Tool)
		}
		seen[pin.Tool] = struct{}{}
	}
	return nil
}
