package registry

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Load parses a registry from TOML at the given path.
func Load(path string) (*Registry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read registry %s: %w", path, err)
	}
	return Parse(data)
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
	seen := make(map[string]struct{}, len(reg.Tools))
	for i := range reg.Tools {
		t := &reg.Tools[i]
		if t.Name == "" {
			return fmt.Errorf("tool at index %d has no name", i)
		}
		if _, dup := seen[t.Name]; dup {
			return fmt.Errorf("duplicate tool name %q", t.Name)
		}
		seen[t.Name] = struct{}{}
		if t.InstallMethod == "" {
			return fmt.Errorf("tool %q has no install_method", t.Name)
		}
		if len(t.Versions) == 0 {
			return fmt.Errorf("tool %q has no versions", t.Name)
		}
	}
	return nil
}
