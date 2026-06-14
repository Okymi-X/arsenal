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
	tools := reg.Search(query)
	assets := reg.SearchAssets(query)
	if len(tools) == 0 && len(assets) == 0 {
		a.log.Printf("no tools or assets match %q", query)
		return nil
	}
	for _, t := range tools {
		a.log.Printf("%-16s [%s] %s", t.Name, t.Category, t.Description)
	}
	for _, as := range assets {
		a.log.Printf("%-16s [%s] %s (asset)", as.Name, as.Category, as.Description)
	}
	return nil
}
