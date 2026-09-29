package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/BurntSushi/toml"
)

// Load parses a registry from TOML at the given path.
func Load(path string) (*Registry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read registry %s: %w", path, err)
	}
	expanded, err := Expand(data, func(name string) ([]byte, error) {
		return os.ReadFile(filepath.Join(filepath.Dir(path), filepath.FromSlash(name)))
	})
	if err != nil {
		return nil, err
	}
	return Parse(expanded)
}

// Parse decodes a registry from raw TOML bytes.
func Parse(data []byte) (*Registry, error) {
	var reg Registry
	if err := toml.Unmarshal(data, &reg); err != nil {
		return nil, fmt.Errorf("parse registry: %w", err)
	}
	if err := validate(&reg); err != nil {
		return nil, err
	}
	return &reg, nil
}

func validate(reg *Registry) error {
	if len(reg.Tools) == 0 && len(reg.Assets) == 0 {
		return fmt.Errorf("registry has no tools or assets")
	}
	seen := make(map[string]struct{}, len(reg.Tools))
	for i := range reg.Tools {
		if err := validateTool(&reg.Tools[i], i, seen); err != nil {
			return err
		}
	}
	return validateAssets(reg)
}

func validateTool(tool *Tool, index int, seen map[string]struct{}) error {
	if tool.Name == "" {
		return fmt.Errorf("tool at index %d has no name", index)
	}
	if err := validateRegistryComponent("tool name", tool.Name); err != nil {
		return err
	}
	if _, duplicate := seen[tool.Name]; duplicate {
		return fmt.Errorf("duplicate tool name %q", tool.Name)
	}
	seen[tool.Name] = struct{}{}
	if tool.InstallMethod == "" {
		return fmt.Errorf("tool %q has no install_method", tool.Name)
	}
	if len(tool.Versions) == 0 {
		return fmt.Errorf("tool %q has no versions", tool.Name)
	}
	for _, binary := range tool.AllBinaries() {
		if err := validateRegistryComponent("binary", binary); err != nil {
			return fmt.Errorf("tool %q: %w", tool.Name, err)
		}
	}
	return validateVersions(*tool)
}

func validateVersions(tool Tool) error {
	seen := make(map[string]struct{}, len(tool.Versions))
	for _, version := range tool.Versions {
		if err := validateRegistryComponent("version tag", version.Tag); err != nil {
			return fmt.Errorf("tool %q: %w", tool.Name, err)
		}
		if _, duplicate := seen[version.Tag]; duplicate {
			return fmt.Errorf("tool %q has duplicate version %q", tool.Name, version.Tag)
		}
		seen[version.Tag] = struct{}{}
		if err := validateInstallTargets(tool, version); err != nil {
			return err
		}
	}
	return nil
}

func validateRegistryComponent(kind, value string) error {
	if value == "" || value == "." || value == ".." || filepath.Base(value) != value {
		return fmt.Errorf("invalid %s %q", kind, value)
	}
	for _, char := range value {
		if unicode.IsLetter(char) || unicode.IsDigit(char) || strings.ContainsRune("._+-", char) {
			continue
		}
		return fmt.Errorf("invalid %s %q", kind, value)
	}
	return nil
}

func validateInstallTargets(tool Tool, version Version) error {
	if tool.InstallMethod != "gobin" && tool.InstallMethod != "cargo" {
		return nil
	}
	binaries := tool.AllBinaries()
	if len(binaries) == 0 {
		return fmt.Errorf("tool %q has no exposed binary", tool.Name)
	}
	if len(version.InstallTargets) != len(binaries) {
		return fmt.Errorf("tool %q version %q install_targets must map every binary", tool.Name, version.Tag)
	}
	seen := make(map[string]struct{}, len(binaries))
	for _, binary := range binaries {
		if _, duplicate := seen[binary]; duplicate {
			return fmt.Errorf("tool %q has duplicate binary %q", tool.Name, binary)
		}
		seen[binary] = struct{}{}
		if version.InstallTargets[binary] == "" {
			return fmt.Errorf("tool %q version %q has no install target for %q", tool.Name, version.Tag, binary)
		}
	}
	return nil
}

func validateAssets(reg *Registry) error {
	seen := make(map[string]struct{}, len(reg.Assets))
	for i := range reg.Assets {
		a := &reg.Assets[i]
		if a.Name == "" {
			return fmt.Errorf("asset at index %d has no name", i)
		}
		if _, dup := seen[a.Name]; dup {
			return fmt.Errorf("duplicate asset name %q", a.Name)
		}
		seen[a.Name] = struct{}{}
		if a.Source != AssetGitHubRelease && a.Source != AssetGitHubRaw {
			return fmt.Errorf("asset %q has invalid source %q", a.Name, a.Source)
		}
	}
	return nil
}
