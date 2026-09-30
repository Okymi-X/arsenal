package cli

import (
	"fmt"
)

// cmdSwitch makes a specific installed version the active one by repointing
// its shims and updating the manifest.
func (a *App) cmdSwitch(args []string) error {
	if len(args) != 2 || args[0] == "" || args[1] == "" {
		return usageError("switch <tool> <version|--select>")
	}
	name, version := a.canonicalName(args[0]), args[1]
	m, err := a.store.Load()
	if err != nil {
		return err
	}
	if version == "--select" {
		versions := m.Versions(name)
		options := make([]selectOption, 0, len(versions))
		for _, installed := range versions {
			status := "installed"
			if installed.Active {
				status = "active"
			}
			options = append(options, selectOption{
				label: fmt.Sprintf("%-24s %s", installed.Version, status),
				value: installed.Version,
			})
		}
		version, err = selectOne(a.in, a.out, "Select an installed version for "+name+":", options)
		if err != nil {
			return err
		}
	}
	target, ok := m.Find(name, version)
	if !ok {
		return fmt.Errorf("%s@%s is not installed", name, version)
	}
	if err := a.activateInstalled(m, target); err != nil {
		return err
	}
	a.log.Printf("[ok] switched %s to %s", name, version)
	return nil
}
