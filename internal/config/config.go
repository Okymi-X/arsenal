package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Okymi-X/arsenal/internal/fsutil"
	"github.com/Okymi-X/arsenal/internal/strictdecode"
)

// Config holds user-tunable settings persisted as JSON under the root.
type Config struct {
	// RegistryURL is the remote manifest or legacy registry fetched by sync.
	RegistryURL string `json:"registry_url"`
	// RegistryRepo is the GitHub owner/repository used by sync --ref.
	RegistryRepo string `json:"registry_repo"`
	// RegistryPath is the manifest path inside RegistryRepo.
	RegistryPath string `json:"registry_path"`
	// RegistryRef is the last user-selected GitHub tag, branch, or SHA.
	RegistryRef string `json:"registry_ref,omitempty"`
	// RegistryCommit is the immutable commit resolved from RegistryRef.
	RegistryCommit string `json:"registry_commit,omitempty"`
	// DefaultBackend selects the isolation backend ("venv" or "container").
	DefaultBackend string `json:"default_backend"`
	// PythonBin is the interpreter used to create virtualenvs.
	PythonBin string `json:"python_bin"`
}

// DefaultConfig returns the built-in configuration.
func DefaultConfig() Config {
	return Config{
		RegistryURL:    "https://raw.githubusercontent.com/Okymi-X/arsenal/main/registry/registry.toml",
		RegistryRepo:   "Okymi-X/arsenal",
		RegistryPath:   "registry/registry.toml",
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
	if err := strictdecode.JSON(data, &cfg); err != nil {
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
	if err := fsutil.WriteFileAtomic(p.configFile(), data, 0o600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}
