package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/Okymi-X/arsenal/internal/registry"
)

type installTarget struct {
	Binary string
	Source string
}

func orderedTargets(tool registry.Tool, version registry.Version) ([]installTarget, error) {
	binaries := tool.AllBinaries()
	if len(version.InstallTargets) != len(binaries) {
		return nil, fmt.Errorf("%s@%s install targets do not match exposed binaries", tool.Name, version.Tag)
	}
	wanted := make(map[string]struct{}, len(binaries))
	for _, binary := range binaries {
		if err := validatePathComponent("binary", binary); err != nil {
			return nil, err
		}
		wanted[binary] = struct{}{}
	}
	keys := make([]string, 0, len(version.InstallTargets))
	for binary := range version.InstallTargets {
		if _, ok := wanted[binary]; !ok {
			return nil, fmt.Errorf("install target has unknown binary %q", binary)
		}
		keys = append(keys, binary)
	}
	sort.Strings(keys)
	targets := make([]installTarget, 0, len(keys))
	for _, binary := range keys {
		targets = append(targets, installTarget{Binary: binary, Source: version.InstallTargets[binary]})
	}
	return targets, nil
}

func renameInstalledBinary(root, installed, exposed string) error {
	if err := validatePathComponent("installed binary", installed); err != nil {
		return err
	}
	if err := validatePathComponent("exposed binary", exposed); err != nil {
		return err
	}
	if installed == exposed {
		return nil
	}
	binDir := filepath.Join(root, "bin")
	if err := os.Rename(filepath.Join(binDir, installed), filepath.Join(binDir, exposed)); err != nil {
		return fmt.Errorf("name installed binary %q as %q: %w", installed, exposed, err)
	}
	return nil
}
