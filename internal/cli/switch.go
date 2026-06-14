package cli

import (
	"fmt"
	"path/filepath"
)

// cmdSwitch makes a specific installed version the active one by repointing
// its shims and updating the manifest.
func (a *App) cmdSwitch(args []string) error {
	if len(args) != 2 {
		return usageError("switch <tool> <version>")
	}
	name, version := a.canonicalName(args[0]), args[1]
	m, err := a.store.Load()
	if err != nil {
		return err
	}
	target, ok := m.Find(name, version)
	if !ok {
		return fmt.Errorf("%s@%s is not installed", name, version)
	}
	for _, bin := range target.Binaries {
		dest := filepath.Join(target.Path, "bin", bin)
		if err := a.shims.Write(bin, dest); err != nil {
			return err
		}
	}
	if !m.SetActive(name, version) {
		return fmt.Errorf("failed to activate %s@%s", name, version)
	}
	if err := a.store.Save(m); err != nil {
		return err
	}
	a.log.Printf("[ok] switched %s to %s", name, version)
	return nil
}
