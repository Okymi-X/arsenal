package cli

// cmdSearch lists registry tools matching a query, or all tools when empty.
func (a *App) cmdSearch(args []string) error {
	query := ""
	if len(args) > 0 {
		query = args[0]
	}
	reg, err := a.loadRegistry()
	if err != nil {
		return err
	}
	matches := reg.Search(query)
	if len(matches) == 0 {
		a.log.Printf("no tools match %q", query)
		return nil
	}
	for _, t := range matches {
		a.log.Printf("%-16s [%s] %s", t.Name, t.Category, t.Description)
	}
	return nil
}
