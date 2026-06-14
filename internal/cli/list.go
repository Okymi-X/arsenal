package cli

import (
	"sort"

	"github.com/Okymi-X/arsenal/internal/store"
)

// cmdList prints installed tools and their versions, marking the active one.
func (a *App) cmdList(args []string) error {
	m, err := a.store.Load()
	if err != nil {
		return err
	}
	if len(m.Tools) == 0 {
		a.log.Printf("no tools installed")
		return nil
	}
	for _, t := range sortInstalls(m.Tools) {
		marker := "   "
		if t.Active {
			marker = "[*]"
		}
		a.log.Printf("%s %s@%s -> %s", marker, t.Name, t.Version, t.Path)
	}
	return nil
}

// sortInstalls returns the installs ordered by name then version.
func sortInstalls(tools []store.InstalledTool) []store.InstalledTool {
	out := make([]store.InstalledTool, len(tools))
	copy(out, tools)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Name != out[j].Name {
			return out[i].Name < out[j].Name
		}
		return out[i].Version < out[j].Version
	})
	return out
}
