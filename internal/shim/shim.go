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
	"strings"
	"unicode"
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
	if err := validateBinaryName(binary); err != nil {
		return err
	}
	if err := os.MkdirAll(m.binDir, 0o755); err != nil {
		return fmt.Errorf("create bin dir: %w", err)
	}
	path := filepath.Join(m.binDir, binary)
	content := script(binary, targetBin)
	if err := writeAtomic(path, []byte(content)); err != nil {
		return fmt.Errorf("write shim %s: %w", binary, err)
	}
	return nil
}

func writeAtomic(path string, content []byte) error {
	temporary, err := os.CreateTemp(filepath.Dir(path), ".shim-*")
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(temporary.Name()) }()
	if err := temporary.Chmod(0o755); err != nil {
		_ = temporary.Close()
		return err
	}
	if _, err := temporary.Write(content); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Sync(); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	return os.Rename(temporary.Name(), path)
}

// Remove deletes a shim if present.
func (m *Manager) Remove(binary string) error {
	if err := validateBinaryName(binary); err != nil {
		return err
	}
	path := filepath.Join(m.binDir, binary)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove shim %s: %w", binary, err)
	}
	return nil
}

// Path returns the absolute path of a shim by binary name.
func (m *Manager) Path(binary string) string {
	if validateBinaryName(binary) != nil {
		return ""
	}
	return filepath.Join(m.binDir, binary)
}

// Exists reports whether a shim for binary is present.
func (m *Manager) Exists(binary string) bool {
	if validateBinaryName(binary) != nil {
		return false
	}
	_, err := os.Stat(filepath.Join(m.binDir, binary))
	return err == nil
}

func validateBinaryName(binary string) error {
	if binary == "" || binary == "." || binary == ".." || filepath.Base(binary) != binary {
		return fmt.Errorf("invalid shim name %q", binary)
	}
	for _, char := range binary {
		if unicode.IsLetter(char) || unicode.IsDigit(char) || strings.ContainsRune("._+-", char) {
			continue
		}
		return fmt.Errorf("invalid shim name %q", binary)
	}
	return nil
}
