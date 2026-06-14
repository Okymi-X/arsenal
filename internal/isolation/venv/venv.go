// Package venv implements the isolation.Backend interface using Python
// virtual environments. Each tool/version gets its own venv directory; pip
// installs land inside it, and binaries are executed from its bin directory.
package venv

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Okymi-X/arsenal/internal/isolation"
)

// Backend provisions and drives a single Python virtualenv.
//
// It satisfies isolation.Backend. A Backend is bound to one tool/version via
// the directory derived from the configured tools root.
type Backend struct {
	pythonBin string
	toolsRoot string
	dir       string
}

// New returns a venv Backend rooted at toolsRoot, using pythonBin to create
// environments. The returned Backend is not yet bound to a tool; call Create.
func New(pythonBin, toolsRoot string) *Backend {
	return &Backend{pythonBin: pythonBin, toolsRoot: toolsRoot}
}

// Path returns the environment directory, empty until Create is called.
func (b *Backend) Path() string { return b.dir }

// Exists reports whether the virtualenv has been provisioned.
func (b *Backend) Exists() bool {
	if b.dir == "" {
		return false
	}
	_, err := os.Stat(b.pythonExe())
	return err == nil
}

// Create provisions a virtualenv for tool at version.
func (b *Backend) Create(tool, version string) error {
	b.dir = filepath.Join(b.toolsRoot, tool, version)
	if b.Exists() {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(b.dir), 0o755); err != nil {
		return fmt.Errorf("create tool parent dir: %w", err)
	}
	if err := runCommand(b.pythonBin, "-m", "venv", b.dir); err != nil {
		return fmt.Errorf("create virtualenv for %s@%s: %w", tool, version, err)
	}
	return nil
}

// Remove tears down the virtualenv directory.
func (b *Backend) Remove() error {
	if b.dir == "" {
		return fmt.Errorf("backend not bound to a tool")
	}
	if err := os.RemoveAll(b.dir); err != nil {
		return fmt.Errorf("remove virtualenv %s: %w", b.dir, err)
	}
	return nil
}

func (b *Backend) binDir() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(b.dir, "Scripts")
	}
	return filepath.Join(b.dir, "bin")
}

func (b *Backend) pythonExe() string { return filepath.Join(b.binDir(), "python") }

func (b *Backend) pipExe() string { return filepath.Join(b.binDir(), "pip") }

// Static assertion that Backend satisfies the interface.
var _ isolation.Backend = (*Backend)(nil)
