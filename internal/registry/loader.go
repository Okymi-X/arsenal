package registry

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Okymi-X/arsenal/internal/safepath"
	"github.com/Okymi-X/arsenal/internal/strictdecode"
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
	if err := strictdecode.TOML(data, &reg); err != nil {
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
	if err := safepath.ValidateComponent("tool name", tool.Name); err != nil {
		return err
	}
	if _, duplicate := seen[tool.Name]; duplicate {
		return fmt.Errorf("duplicate tool name %q", tool.Name)
	}
	seen[tool.Name] = struct{}{}
	if tool.InstallMethod == "" {
		return fmt.Errorf("tool %q has no install_method", tool.Name)
	}
	if !validInstallMethod(tool.InstallMethod) {
		return fmt.Errorf("tool %q has invalid install_method %q", tool.Name, tool.InstallMethod)
	}
	if len(tool.Versions) == 0 {
		return fmt.Errorf("tool %q has no versions", tool.Name)
	}
	for _, binary := range tool.AllBinaries() {
		if err := safepath.ValidateComponent("binary", binary); err != nil {
			return fmt.Errorf("tool %q: %w", tool.Name, err)
		}
	}
	return validateVersions(*tool)
}

func validateVersions(tool Tool) error {
	seen := make(map[string]struct{}, len(tool.Versions))
	for _, version := range tool.Versions {
		if err := safepath.ValidateComponent("version tag", version.Tag); err != nil {
			return fmt.Errorf("tool %q: %w", tool.Name, err)
		}
		if _, duplicate := seen[version.Tag]; duplicate {
			return fmt.Errorf("tool %q has duplicate version %q", tool.Name, version.Tag)
		}
		seen[version.Tag] = struct{}{}
		if err := validatePipMetadata(tool, version); err != nil {
			return err
		}
		if err := validateInstallTargets(tool, version); err != nil {
			return err
		}
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
		if err := safepath.ValidateComponent("asset name", a.Name); err != nil {
			return fmt.Errorf("asset at index %d: %w", i, err)
		}
		if _, dup := seen[a.Name]; dup {
			return fmt.Errorf("duplicate asset name %q", a.Name)
		}
		seen[a.Name] = struct{}{}
		if a.Source != AssetGitHubRelease && a.Source != AssetGitHubRaw {
			return fmt.Errorf("asset %q has invalid source %q", a.Name, a.Source)
		}
		if _, err := ParseGitHubRepo(a.Repo); err != nil {
			return fmt.Errorf("asset %q: %w", a.Name, err)
		}
		if a.Source == AssetGitHubRelease && len(a.Builds) != 0 {
			return fmt.Errorf("asset %q: builds require github-raw source", a.Name)
		}
		if err := validateAssetRepoPath("directory", a.Dir, true); err != nil {
			return fmt.Errorf("asset %q: %w", a.Name, err)
		}
		for _, build := range a.Builds {
			if err := validateAssetRepoPath("build", build, false); err != nil {
				return fmt.Errorf("asset %q: %w", a.Name, err)
			}
		}
	}
	return nil
}

func validInstallMethod(method string) bool {
	switch method {
	case "pip", "gitpip", "binary", "gobin", "cargo":
		return true
	default:
		return false
	}
}

func validateAssetRepoPath(kind, value string, allowEmpty bool) error {
	if value == "" && allowEmpty {
		return nil
	}
	if value == "" || strings.Contains(value, "\\") || path.IsAbs(value) ||
		path.Clean(value) != value || strings.HasPrefix(value, "../") {
		return fmt.Errorf("invalid asset %s %q", kind, value)
	}
	return nil
}
