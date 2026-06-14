// Package cli routes command-line arguments to commands. Each command lives in
// its own file and receives the wired App explicitly. The router owns global
// flag parsing (verbosity) and usage; commands own their own behavior.
package cli

import (
	"fmt"
	"sort"
)

// command is a single CLI subcommand.
type command struct {
	name    string
	summary string
	run     func(a *App, args []string) error
}

// commands returns the full command table, built fresh per call to avoid
// shared mutable state.
func (a *App) commands() []command {
	return []command{
		{"install", "Install a tool at a tested version", (*App).cmdInstall},
		{"remove", "Remove an installed tool", (*App).cmdRemove},
		{"switch", "Switch the active version of a tool", (*App).cmdSwitch},
		{"list", "List installed tools", (*App).cmdList},
		{"search", "Search the registry", (*App).cmdSearch},
		{"info", "Show registry details for a tool", (*App).cmdInfo},
		{"run", "Run an installed tool", (*App).cmdRun},
		{"op", "Manage engagement profiles", (*App).cmdOp},
		{"sync", "Sync the registry from upstream", (*App).cmdSync},
		{"doctor", "Diagnose and repair installs", (*App).cmdDoctor},
		{"bundle", "Export an offline bundle", (*App).cmdBundle},
		{"version", "Print the build version", (*App).cmdVersion},
	}
}

// Main parses global flags, dispatches to a command, and returns an exit code.
func (a *App) Main(args []string) int {
	verbose, rest := extractVerbose(args)
	a.setVerbose(verbose)

	if len(rest) == 0 || isHelp(rest[0]) {
		a.printUsage()
		return 0
	}
	name, cmdArgs := rest[0], rest[1:]
	for _, c := range a.commands() {
		if c.name == name {
			if err := c.run(a, cmdArgs); err != nil {
				a.log.Errorf("%v", err)
				return 1
			}
			return 0
		}
	}
	a.log.Errorf("unknown command %q", name)
	a.printUsage()
	return 2
}

func (a *App) printUsage() {
	a.log.Printf("arsenal - package manager for offensive-security tooling")
	a.log.Printf("")
	a.log.Printf("Usage: arsenal [-v|--verbose] <command> [args]")
	a.log.Printf("")
	a.log.Printf("Commands:")
	cmds := a.commands()
	sort.Slice(cmds, func(i, j int) bool { return cmds[i].name < cmds[j].name })
	for _, c := range cmds {
		a.log.Printf("  %-9s %s", c.name, c.summary)
	}
}

func isHelp(arg string) bool {
	return arg == "-h" || arg == "--help" || arg == "help"
}

// extractVerbose removes -v/--verbose from args and reports whether it was set.
func extractVerbose(args []string) (bool, []string) {
	verbose := false
	rest := make([]string, 0, len(args))
	for _, a := range args {
		if a == "-v" || a == "--verbose" {
			verbose = true
			continue
		}
		rest = append(rest, a)
	}
	return verbose, rest
}

// usageError returns an error describing correct command usage.
func usageError(usage string) error {
	return fmt.Errorf("usage: arsenal %s", usage)
}
