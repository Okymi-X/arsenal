// Package shim generates and switches PATH shims so multiple tool versions
// coexist and the active one is selectable.
//
// A shim is a tiny shell script in the user's bin directory that execs the
// real binary from the active version's isolated environment. Switching a
// version rewrites the shim to point at the new environment; no other state
// changes.
package shim

import (
	"fmt"
	"os"
	"path/filepath"
)

// Manager creates and removes shims in a single bin directory.
type Manager struct {
	binDir string
}

// NewManager returns a Manager writing shims into binDir.
func NewManager(binDir string) *Manager { return &Manager{binDir: binDir} }

// BinDir returns the directory holding generated shims.
func (m *Manager) BinDir() string { return m.binDir }

// Write creates or replaces a shim named binary that execs targetBin.
func (m *Manager) Write(binary, targetBin string) error {
	if err := os.MkdirAll(m.binDir, 0o755); err != nil {
		return fmt.Errorf("create bin dir: %w", err)
	}
	path := filepath.Join(m.binDir, binary)
	content := script(binary, targetBin)
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		return fmt.Errorf("write shim %s: %w", binary, err)
	}
	return nil
}

// Remove deletes a shim if present.
func (m *Manager) Remove(binary string) error {
	path := filepath.Join(m.binDir, binary)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove shim %s: %w", binary, err)
	}
	return nil
}

// Path returns the absolute path of a shim by binary name.
func (m *Manager) Path(binary string) string {
	return filepath.Join(m.binDir, binary)
}

// Exists reports whether a shim for binary is present.
func (m *Manager) Exists(binary string) bool {
	_, err := os.Stat(filepath.Join(m.binDir, binary))
	return err == nil
}
