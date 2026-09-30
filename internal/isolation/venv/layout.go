package venv

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Okymi-X/arsenal/internal/safepath"
)

func (b *Backend) bind(tool, version string) error {
	if err := safepath.ValidateComponent("tool", tool); err != nil {
		return err
	}
	if err := safepath.ValidateComponent("version", version); err != nil {
		return err
	}
	if b.toolsRoot == "" {
		return fmt.Errorf("tools root is empty")
	}
	toolDir := filepath.Join(b.toolsRoot, tool)
	if err := rejectUnsafeExistingDirectory(toolDir); err != nil {
		return err
	}
	b.dir = filepath.Join(toolDir, version)
	return rejectUnsafeExistingDirectory(b.dir)
}

func rejectUnsafeExistingDirectory(path string) error {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect environment path %s: %w", path, err)
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return fmt.Errorf("environment path is not a regular directory: %s", path)
	}
	return nil
}
