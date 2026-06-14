// Command regcheck verifies that every version in the curated registry exists
// at its official source: PyPI for pip tools, and the upstream Git repository
// for gitpip, gobin, cargo, and binary tools.
//
// It is run in CI (the registry-check workflow) on changes to the registry and
// on a schedule, so a version that does not exist upstream can never ship.
// Exit code 0 means every version resolved; non-zero means at least one did
// not, and the offending entries are printed.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Okymi-X/arsenal/internal/registry"
)

func main() {
	path := flag.String("registry", "registry/registry.toml", "path to registry.toml")
	flag.Parse()

	reg, err := registry.Load(*path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[fail] load registry: %v\n", err)
		os.Exit(2)
	}

	c := newChecker(20*time.Second, os.Getenv("GITHUB_TOKEN"))
	failures := run(c, reg)

	fmt.Printf("\nchecked %d tools\n", len(reg.Tools))
	if failures > 0 {
		fmt.Printf("[fail] %d version(s) did not resolve upstream\n", failures)
		os.Exit(1)
	}
	fmt.Println("[ok] every registry version resolved upstream")
}

// run verifies every version of every tool, printing one line per version and
// returning the number of failures.
func run(c *checker, reg *registry.Registry) int {
	failures := 0
	for _, tool := range reg.Tools {
		for _, v := range tool.Versions {
			if err := c.verify(tool, v); err != nil {
				failures++
				fmt.Printf("[fail] %s@%s (%s): %v\n", tool.Name, v.Tag, tool.InstallMethod, err)
				continue
			}
			fmt.Printf("[ok]   %s@%s (%s)\n", tool.Name, v.Tag, tool.InstallMethod)
		}
	}
	return failures
}
