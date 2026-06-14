package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Config holds user-tunable settings persisted as JSON under the root.
type Config struct {
	// RegistryURL is the remote registry fetched by `arsenal sync`.
	RegistryURL string `json:"registry_url"`
	// DefaultBackend selects the isolation backend ("venv" or "container").
	DefaultBackend string `json:"default_backend"`
	// PythonBin is the interpreter used to create virtualenvs.
	PythonBin string `json:"python_bin"`
}

// DefaultConfig returns the built-in configuration.
func DefaultConfig() Config {
	return Config{
		RegistryURL:    "https://raw.githubusercontent.com/Okymi-X/arsenal/main/registry/registry.toml",
		DefaultBackend: "venv",
		PythonBin:      "python3",
	}
}

func (p Paths) configFile() string { return filepath.Join(p.Root, "config.json") }

// LoadConfig reads the config file, returning defaults when it is absent.
func (p Paths) LoadConfig() (Config, error) {
	data, err := os.ReadFile(p.configFile())
	if errors.Is(err, fs.ErrNotExist) {
		return DefaultConfig(), nil
	}
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}
	cfg := DefaultConfig()
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

// SaveConfig writes the config file, creating the root if needed.
func (p Paths) SaveConfig(cfg Config) error {
	if err := os.MkdirAll(p.Root, 0o755); err != nil {
		return fmt.Errorf("create root: %w", err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	if err := os.WriteFile(p.configFile(), data, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}
