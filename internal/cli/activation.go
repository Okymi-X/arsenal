package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Okymi-X/arsenal/internal/store"
)

func (a *App) activateInstalled(manifest *store.Manifest, target store.InstalledTool) error {
	if !installationHealthy(target) {
		return fmt.Errorf("%s@%s installation is incomplete; run doctor", target.Name, target.Version)
	}
	for _, binary := range target.Binaries {
		destination := filepath.Join(target.Path, "bin", binary)
		if err := a.shims.Write(binary, destination); err != nil {
			return err
		}
	}
	if !manifest.SetActive(target.Name, target.Version) {
		return fmt.Errorf("failed to activate %s@%s", target.Name, target.Version)
	}
	return a.store.Save(manifest)
}

func installationHealthy(tool store.InstalledTool) bool {
	if info, err := os.Stat(tool.Path); err != nil || !info.IsDir() {
		return false
	}
	for _, binary := range tool.Binaries {
		info, err := os.Stat(filepath.Join(tool.Path, "bin", binary))
		if err != nil || info.IsDir() {
			return false
		}
	}
	return len(tool.Binaries) > 0
}
