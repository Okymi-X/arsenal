package cli

import "fmt"

// cmdRun executes the active version of a tool, forwarding arguments after an
// optional "--" separator.
func (a *App) cmdRun(args []string) error {
	if len(args) == 0 {
		return usageError("run <tool> -- <args...>")
	}
	name := args[0]
	forwarded := forwardedArgs(args[1:])

	m, err := a.store.Load()
	if err != nil {
		return err
	}
	active, ok := m.Active(name)
	if !ok {
		return fmt.Errorf("%s has no active version; run 'arsenal install %s'", name, name)
	}
	if len(active.Binaries) == 0 {
		return fmt.Errorf("%s has no runnable binary", name)
	}
	backend := a.newBackend()
	if err := backend.Create(active.Name, active.Version); err != nil {
		return err
	}
	return backend.Run(append([]string{active.Binaries[0]}, forwarded...))
}

// forwardedArgs strips a leading "--" separator if present.
func forwardedArgs(args []string) []string {
	if len(args) > 0 && args[0] == "--" {
		return args[1:]
	}
	return args
}
