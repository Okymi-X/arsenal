package op

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/Okymi-X/arsenal/internal/fsutil"
	"github.com/Okymi-X/arsenal/internal/safepath"
	"github.com/Okymi-X/arsenal/internal/strictdecode"
)

// WriteLockfile encodes a lockfile to TOML at path, atomically.
func WriteLockfile(path string, lf *Lockfile) error {
	if err := validateLockfile(lf); err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(lf); err != nil {
		return fmt.Errorf("encode lockfile: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create lockfile dir: %w", err)
	}
	if err := fsutil.WriteFileAtomic(path, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write lockfile: %w", err)
	}
	return nil
}

// ReadLockfile decodes and validates a lockfile from TOML at path.
func ReadLockfile(path string) (*Lockfile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read lockfile: %w", err)
	}
	return ParseLockfile(data)
}

// ParseLockfile decodes and validates a lockfile from raw TOML.
func ParseLockfile(data []byte) (*Lockfile, error) {
	var lf Lockfile
	if err := strictdecode.TOML(data, &lf); err != nil {
		return nil, fmt.Errorf("parse lockfile: %w", err)
	}
	if err := validateLockfile(&lf); err != nil {
		return nil, err
	}
	return &lf, nil
}

func validateLockfile(lf *Lockfile) error {
	if lf == nil {
		return fmt.Errorf("lockfile is required")
	}
	if err := safepath.ValidateComponent("op name", lf.Op); err != nil {
		return fmt.Errorf("lockfile: %w", err)
	}
	if lf.RegistryVersion == "" {
		return fmt.Errorf("lockfile has no registry version")
	}
	seen := make(map[string]struct{}, len(lf.Entries))
	for _, e := range lf.Entries {
		if err := safepath.ValidateComponent("tool name", e.Tool); err != nil {
			return fmt.Errorf("lockfile: %w", err)
		}
		if err := safepath.ValidateComponent("version", e.Version); err != nil {
			return fmt.Errorf("lockfile entry for %q: %w", e.Tool, err)
		}
		if e.InstallMethod == "" {
			return fmt.Errorf("lockfile entry for %q has no install method", e.Tool)
		}
		if _, dup := seen[e.Tool]; dup {
			return fmt.Errorf("duplicate lockfile entry for %q", e.Tool)
		}
		seen[e.Tool] = struct{}{}
	}
	return nil
}
