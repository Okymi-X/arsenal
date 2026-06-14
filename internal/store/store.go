// Package store persists the manifest of installed tools.
//
// The manifest records which tool versions are installed, which one is active,
// and metadata needed by list, switch, doctor, and remove. The Store interface
// is swappable so persistence (JSON file today) can change without affecting
// callers.
package store

// InstalledTool records a single installed tool/version environment.
type InstalledTool struct {
	// Name is the canonical tool name.
	Name string `json:"name"`
	// Version is the installed version tag.
	Version string `json:"version"`
	// Backend names the isolation backend used ("venv" or "container").
	Backend string `json:"backend"`
	// Path is the environment directory.
	Path string `json:"path"`
	// Binaries are the executables exposed via shims.
	Binaries []string `json:"binaries"`
	// Active marks the version currently selected for its tool.
	Active bool `json:"active"`
	// InstalledAt is the RFC-3339 timestamp of installation.
	InstalledAt string `json:"installed_at"`
}

// Manifest is the complete set of installed tools.
type Manifest struct {
	// Tools is every installed tool/version.
	Tools []InstalledTool `json:"tools"`
}

// Store reads and writes the installed-tools manifest.
type Store interface {
	// Load returns the manifest, or an empty manifest when none exists.
	Load() (*Manifest, error)
	// Save persists the manifest in full.
	Save(m *Manifest) error
}
