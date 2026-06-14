// Command arsenal is a package and environment manager specialized for
// offensive-security tooling. This entry point only wires dependencies and
// delegates to the cli package.
package main

import (
	"fmt"
	"os"

	"github.com/Okymi-X/arsenal/internal/cli"
	"github.com/Okymi-X/arsenal/internal/config"
)

// version is injected at build time via -ldflags. It is not hardcoded.
var version = "dev"

func main() {
	os.Exit(run())
}

// run wires dependencies and dispatches, returning a process exit code.
func run() int {
	root, err := config.DefaultRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[fail] %v\n", err)
		return 1
	}
	paths := config.NewPaths(root)
	cfg, err := paths.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[fail] %v\n", err)
		return 1
	}
	app := cli.New(cli.Options{
		Paths:   paths,
		Cfg:     cfg,
		Version: version,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	})
	return app.Main(os.Args[1:])
}
