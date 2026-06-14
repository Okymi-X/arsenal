package cli

import (
	"github.com/Okymi-X/arsenal/internal/store"
)

// cmdRemove removes every installed version of a tool, its environments, and
// its shims, then prunes the manifest in a single save.
func (a *App) cmdRemove(args []string) error {
	if len(args) != 1 {
		return usageError("remove <tool>")
	}
	name := args[0]
	m, err := a.store.Load()
	if err != nil {
		return err
	}
	installs := m.Versions(name)
	if len(installs) == 0 {
		a.log.Warnf("%s is not installed", name)
		return nil
	}
	for _, t := range installs {
		if err := a.tearDown(t); err != nil {
			return err
		}
		m.Delete(t.Name, t.Version)
	}
	if err := a.store.Save(m); err != nil {
		return err
	}
	a.log.Printf("[ok] removed %s", name)
	return nil
}

// tearDown removes an installation's environment and shims, leaving the
// manifest to the caller.
func (a *App) tearDown(t store.InstalledTool) error {
	backend := a.newBackend()
	if err := backend.Create(t.Name, t.Version); err == nil && backend.Exists() {
		if err := backend.Remove(); err != nil {
			return err
		}
	}
	for _, bin := range t.Binaries {
		if err := a.shims.Remove(bin); err != nil {
			return err
		}
	}
	return nil
}
