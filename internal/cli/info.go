package cli

import (
	"fmt"

	"github.com/Okymi-X/arsenal/internal/registry"
)

// cmdInfo prints the registry details for a tool (with its versions) or an
// asset.
func (a *App) cmdInfo(args []string) error {
	if len(args) != 1 {
		return usageError("info <tool|asset>")
	}
	reg, err := a.loadRegistry()
	if err != nil {
		return err
	}
	if tool, ok := reg.FindTool(args[0]); ok {
		a.printToolHeader(tool)
		a.printToolVersions(tool)
		return nil
	}
	if asset, ok := reg.FindAsset(args[0]); ok {
		a.printAssetInfo(asset)
		return nil
	}
	return fmt.Errorf("%q not found in registry", args[0])
}

func (a *App) printAssetInfo(as registry.Asset) {
	a.log.Printf("%s - %s (asset)", as.Name, as.Description)
	a.log.Printf("category:       %s", as.Category)
	a.log.Printf("source:         %s", as.Source)
	a.log.Printf("repo:           %s", as.Repo)
	if len(as.Aliases) > 0 {
		a.log.Printf("aliases:        %v", as.Aliases)
	}
	if as.Dir != "" {
		a.log.Printf("directory:      %s", as.Dir)
	}
	if as.Pattern != "" {
		a.log.Printf("default file:   %s", as.Pattern)
	}
	if as.Notes != "" {
		a.log.Printf("notes:          %s", as.Notes)
	}
	a.log.Printf("fetch:          arsenal fetch %s [binary] --dest <dir>", as.Name)
}

func (a *App) printToolHeader(t registry.Tool) {
	a.log.Printf("%s - %s", t.Name, t.Description)
	a.log.Printf("category:       %s", t.Category)
	a.log.Printf("install method: %s", t.InstallMethod)
	if t.PythonVersion != "" {
		a.log.Printf("python:         %s", t.PythonVersion)
	}
	a.log.Printf("repo:           %s", t.Repo)
	if len(t.Aliases) > 0 {
		a.log.Printf("aliases:        %v", t.Aliases)
	}
	if t.RequiresRoot {
		a.log.Printf("requires root:  yes")
	}
	if t.Notes != "" {
		a.log.Printf("notes:          %s", t.Notes)
	}
}

func (a *App) printToolVersions(t registry.Tool) {
	a.log.Printf("versions:")
	for _, v := range t.Versions {
		marker := "   "
		if v.Tested {
			marker = "[ok]"
		}
		a.log.Printf("  %s %-10s %s %s", marker, v.Tag, v.Date, v.Notes)
	}
}
