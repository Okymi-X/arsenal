package op

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
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
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write lockfile: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("replace lockfile: %w", err)
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
	if err := toml.Unmarshal(data, &lf); err != nil {
		return nil, fmt.Errorf("parse lockfile: %w", err)
	}
	if err := validateLockfile(&lf); err != nil {
		return nil, err
	}
	return &lf, nil
}

func validateLockfile(lf *Lockfile) error {
	if lf.Op == "" {
		return fmt.Errorf("lockfile has no op name")
	}
	seen := make(map[string]struct{}, len(lf.Entries))
	for _, e := range lf.Entries {
		if e.Tool == "" {
			return fmt.Errorf("lockfile entry has no tool name")
		}
		if e.Version == "" {
			return fmt.Errorf("lockfile entry for %q has no version", e.Tool)
		}
		if _, dup := seen[e.Tool]; dup {
			return fmt.Errorf("duplicate lockfile entry for %q", e.Tool)
		}
		seen[e.Tool] = struct{}{}
	}
	return nil
}
