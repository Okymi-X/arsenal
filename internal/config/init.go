package config

import (
	"fmt"
	"os"
)

// EnsureDirs creates every directory arsenal relies on, idempotently.
func (p Paths) EnsureDirs() error {
	dirs := []string{p.Root, p.Tools, p.Bin, p.Ops, p.Bundles, p.Cache}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return fmt.Errorf("create directory %s: %w", d, err)
		}
	}
	return nil
}
