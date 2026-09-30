package registry

import (
	"fmt"

	"github.com/Okymi-X/arsenal/internal/pipspec"
)

func validatePipMetadata(tool Tool, version Version) error {
	if version.PipSpec != "" {
		if tool.InstallMethod != "pip" {
			return fmt.Errorf("tool %q version %q has pip_spec for %s install", tool.Name, version.Tag, tool.InstallMethod)
		}
		if _, err := pipspec.ParseExact(version.PipSpec); err != nil {
			return fmt.Errorf("tool %q version %q: %w", tool.Name, version.Tag, err)
		}
	}
	if len(version.PipDependencies) != 0 && tool.InstallMethod != "pip" && tool.InstallMethod != "gitpip" {
		return fmt.Errorf("tool %q version %q has pip dependencies for %s install", tool.Name, version.Tag, tool.InstallMethod)
	}
	seen := make(map[string]struct{}, len(version.PipDependencies))
	for _, requirement := range version.PipDependencies {
		pin, err := pipspec.ParseExact(requirement)
		if err != nil {
			return fmt.Errorf("tool %q version %q: %w", tool.Name, version.Tag, err)
		}
		name := pin.NormalizedName()
		if _, duplicate := seen[name]; duplicate {
			return fmt.Errorf("tool %q version %q has duplicate pip dependency %q", tool.Name, version.Tag, pin.Name)
		}
		seen[name] = struct{}{}
	}
	return nil
}
