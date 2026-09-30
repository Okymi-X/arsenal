// Command arsenal is a package and environment manager specialized for
// offensive-security tooling. This entry point only wires dependencies and
// delegates to the cli package.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Okymi-X/arsenal/internal/cli"
	"github.com/Okymi-X/arsenal/internal/config"
)

// version may be injected at build time via -ldflags.
var version = "dev"

func main() {
	os.Exit(run())
}

// run wires dependencies and dispatches, returning a process exit code.
func run() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
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
		Context: ctx,
		Paths:   paths,
		Cfg:     cfg,
		Version: buildVersion(),
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	})
	return app.Main(os.Args[1:])
}
