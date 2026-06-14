package doctor

import (
	"fmt"
	"os"

	"github.com/Okymi-X/arsenal/internal/store"
)

// ManifestCheck verifies every installed environment recorded in the manifest
// still exists on disk, and can prune entries whose directories are gone.
type ManifestCheck struct {
	store store.Store
}

// NewManifestCheck returns a ManifestCheck backed by the given store.
func NewManifestCheck(s store.Store) *ManifestCheck { return &ManifestCheck{store: s} }

// Name identifies the check.
func (c *ManifestCheck) Name() string { return "manifest" }

// Run reports any installed environments whose directories are missing.
func (c *ManifestCheck) Run() Result {
	m, err := c.store.Load()
	if err != nil {
		return Result{Name: c.Name(), OK: false, Detail: err.Error()}
	}
	missing := c.missing(m)
	if len(missing) == 0 {
		return Result{Name: c.Name(), OK: true, Detail: fmt.Sprintf("%d installs tracked", len(m.Tools))}
	}
	return Result{
		Name:    c.Name(),
		OK:      false,
		Detail:  fmt.Sprintf("%d installs missing on disk", len(missing)),
		Fixable: true,
	}
}

// Fix removes manifest entries whose environment directories are missing.
func (c *ManifestCheck) Fix() error {
	m, err := c.store.Load()
	if err != nil {
		return err
	}
	for _, t := range c.missing(m) {
		m.Delete(t.Name, t.Version)
	}
	return c.store.Save(m)
}

func (c *ManifestCheck) missing(m *store.Manifest) []store.InstalledTool {
	var out []store.InstalledTool
	for _, t := range m.Tools {
		if _, err := os.Stat(t.Path); err != nil {
			out = append(out, t)
		}
	}
	return out
}
