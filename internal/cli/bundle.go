package cli

import "github.com/Okymi-X/arsenal/internal/bundle"

// cmdBundle exports a self-contained offline bundle.
//
// The exporter is stubbed; the command is wired against its interface so the
// surface is stable once the implementation lands.
func (a *App) cmdBundle(args []string) error {
	if len(args) == 0 || args[0] != "--offline" {
		return usageError("bundle --offline [op]")
	}
	opName := ""
	if len(args) > 1 {
		opName = args[1]
	}
	exporter := bundle.NewStubExporter()
	a.log.Printf("-> exporting offline bundle for op %q", opName)
	return exporter.Export(opName, a.paths.Bundles)
}
