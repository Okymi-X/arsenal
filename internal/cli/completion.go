package cli

import "fmt"

// cmdCompletion prints a shell completion script for the requested shell.
//
//	arsenal completion bash > /etc/bash_completion.d/arsenal
//	arsenal completion zsh  > "${fpath[1]}/_arsenal"
//	arsenal completion fish > ~/.config/fish/completions/arsenal.fish
//
// The scripts complete subcommands and, dynamically, registry tool names (for
// install and info) and installed tool names (for run, remove, and switch).
func (a *App) cmdCompletion(args []string) error {
	if len(args) != 1 {
		return usageError("completion <bash|zsh|fish>")
	}
	switch args[0] {
	case "bash":
		a.log.Printf("%s", bashCompletion)
	case "zsh":
		a.log.Printf("%s", zshCompletion)
	case "fish":
		a.log.Printf("%s", fishCompletion)
	default:
		return fmt.Errorf("unsupported shell %q (want bash, zsh, or fish)", args[0])
	}
	return nil
}
