package cli

import "github.com/Okymi-X/arsenal/internal/registry"

// cmdInfo prints the registry details for a tool, including its versions.
func (a *App) cmdInfo(args []string) error {
	if len(args) != 1 {
		return usageError("info <tool>")
	}
	reg, err := a.loadRegistry()
	if err != nil {
		return err
	}
	tool, err := reg.MustFindTool(args[0])
	if err != nil {
		return err
	}
	a.printToolHeader(tool)
	a.printToolVersions(tool)
	return nil
}

func (a *App) printToolHeader(t registry.Tool) {
	a.log.Printf("%s - %s", t.Name, t.Description)
	a.log.Printf("category:       %s", t.Category)
	a.log.Printf("install method: %s", t.InstallMethod)
	a.log.Printf("python:         %s", t.PythonVersion)
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
