// Package config resolves filesystem paths and loads user configuration.
//
// All persistent state lives under a single root directory, by default
// $XDG_DATA_HOME/arsenal (falling back to ~/.local/share/arsenal). The
// package never mutates global state; callers pass a Paths value explicitly.
package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Paths holds the resolved locations arsenal reads from and writes to.
type Paths struct {
	// Root is the base directory for all arsenal state.
	Root string
	// Tools holds per-tool, per-version isolated environments.
	Tools string
	// Bin holds generated PATH shims.
	Bin string
	// Ops holds engagement profiles and their lockfiles.
	Ops string
	// Bundles holds exported offline bundles.
	Bundles string
	// Cache holds downloaded registries and transient data.
	Cache string
	// ManifestFile is the installed-tools manifest (JSON).
	ManifestFile string
	// RegistryFile is the active local registry (TOML).
	RegistryFile string
}

// DefaultRoot returns the base directory derived from the environment.
func DefaultRoot() (string, error) {
	if v := os.Getenv("ARSENAL_HOME"); v != "" {
		return v, nil
	}
	if v := os.Getenv("XDG_DATA_HOME"); v != "" {
		return filepath.Join(v, "arsenal"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}
	return filepath.Join(home, ".local", "share", "arsenal"), nil
}

// NewPaths derives all paths from a root directory.
func NewPaths(root string) Paths {
	return Paths{
		Root:         root,
		Tools:        filepath.Join(root, "tools"),
		Bin:          filepath.Join(root, "bin"),
		Ops:          filepath.Join(root, "ops"),
		Bundles:      filepath.Join(root, "bundles"),
		Cache:        filepath.Join(root, "cache"),
		ManifestFile: filepath.Join(root, "manifest.json"),
		RegistryFile: filepath.Join(root, "registry.toml"),
	}
}

// ToolDir returns the environment directory for a tool at a version.
func (p Paths) ToolDir(tool, version string) string {
	return filepath.Join(p.Tools, tool, version)
}
