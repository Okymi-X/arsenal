package cli

import "github.com/Okymi-X/arsenal/internal/doctor"

// cmdDoctor runs health checks and, with --fix, repairs what it can.
func (a *App) cmdDoctor(args []string) error {
	fix := len(args) > 0 && args[0] == "--fix"
	doc := doctor.New(a.doctorChecks())

	for _, r := range doc.Run() {
		marker := "[ok]"
		if !r.OK {
			marker = "[fail]"
		}
		a.log.Printf("%-6s %-12s %s", marker, r.Name, r.Detail)
	}
	if !fix {
		return nil
	}
	fixed, err := doc.Repair()
	for _, name := range fixed {
		a.log.Printf("-> repaired %s", name)
	}
	return err
}

// doctorChecks builds the standard set of health checks.
func (a *App) doctorChecks() []doctor.Check {
	return []doctor.Check{
		doctor.NewDirsCheck(a.paths),
		doctor.NewPythonCheck(a.cfg.PythonBin),
		doctor.NewPathCheck(a.paths.Bin),
		doctor.NewManifestCheck(a.store),
	}
}
