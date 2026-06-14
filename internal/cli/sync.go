package cli

// cmdSync refreshes the local registry from its configured upstream URL.
func (a *App) cmdSync(args []string) error {
	if err := a.paths.EnsureDirs(); err != nil {
		return err
	}
	a.log.Printf("-> syncing registry from %s", a.cfg.RegistryURL)
	if err := a.source.Sync(); err != nil {
		return err
	}
	reg, err := a.source.Load()
	if err != nil {
		return err
	}
	a.log.Printf("[ok] registry updated: %d tools (schema %s)", len(reg.Tools), reg.Version)
	return nil
}
