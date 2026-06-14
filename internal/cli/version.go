package cli

// cmdVersion prints the injected build version.
func (a *App) cmdVersion(args []string) error {
	a.log.Printf("arsenal %s", a.version)
	return nil
}
