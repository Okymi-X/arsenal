package registry

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
)

const (
	maxSegments      = 64
	maxAssembledSize = 16 << 20
)

// Manifest describes a segmented registry and the order in which its files
// are assembled.
type Manifest struct {
	Version       string            `toml:"version"`
	Updated       string            `toml:"updated"`
	Segments      []string          `toml:"segments"`
	SegmentSHA256 map[string]string `toml:"segment_sha256"`
}

// ParseManifest decodes and validates segmented-registry metadata.
func ParseManifest(data []byte) (Manifest, error) {
	var manifest Manifest
	if err := toml.Unmarshal(data, &manifest); err != nil {
		return Manifest{}, fmt.Errorf("parse registry manifest: %w", err)
	}
	if len(manifest.Segments) == 0 {
		return manifest, nil
	}
	if manifest.Version == "" {
		return Manifest{}, fmt.Errorf("registry manifest has no version")
	}
	if len(manifest.Segments) > maxSegments {
		return Manifest{}, fmt.Errorf("registry manifest exceeds %d segments", maxSegments)
	}
	seen := make(map[string]struct{}, len(manifest.Segments))
	for _, name := range manifest.Segments {
		if err := validateSegmentPath(name); err != nil {
			return Manifest{}, err
		}
		if _, exists := seen[name]; exists {
			return Manifest{}, fmt.Errorf("duplicate registry segment %q", name)
		}
		seen[name] = struct{}{}
	}
	for name, digest := range manifest.SegmentSHA256 {
		if _, exists := seen[name]; !exists {
			return Manifest{}, fmt.Errorf("checksum references unknown segment %q", name)
		}
		decoded, err := hex.DecodeString(digest)
		if err != nil || len(decoded) != sha256.Size {
			return Manifest{}, fmt.Errorf("invalid SHA-256 for registry segment %q", name)
		}
	}
	return manifest, nil
}

// Expand assembles a segmented registry. A legacy monolithic registry is
// returned unchanged when it has no segments field.
func Expand(data []byte, load func(string) ([]byte, error)) ([]byte, error) {
	manifest, err := ParseManifest(data)
	if err != nil {
		return nil, err
	}
	if len(manifest.Segments) == 0 {
		return data, nil
	}
	if load == nil {
		return nil, fmt.Errorf("registry segment loader is nil")
	}

	var out bytes.Buffer
	fmt.Fprintf(&out, "version = %q\nupdated = %q\n", manifest.Version, manifest.Updated)
	for _, name := range manifest.Segments {
		segment, err := load(name)
		if err != nil {
			return nil, fmt.Errorf("load registry segment %s: %w", name, err)
		}
		if out.Len()+len(segment) > maxAssembledSize {
			return nil, fmt.Errorf("assembled registry exceeds %d bytes", maxAssembledSize)
		}
		if err := verifySegment(manifest, name, segment); err != nil {
			return nil, err
		}
		out.WriteByte('\n')
		out.Write(bytes.TrimSpace(segment))
		out.WriteByte('\n')
	}
	return out.Bytes(), nil
}

func verifySegment(manifest Manifest, name string, data []byte) error {
	expected, ok := manifest.SegmentSHA256[name]
	if !ok {
		return fmt.Errorf("registry segment %q has no SHA-256", name)
	}
	digest := sha256.Sum256(data)
	if !strings.EqualFold(hex.EncodeToString(digest[:]), expected) {
		return fmt.Errorf("registry segment %q failed SHA-256 verification", name)
	}
	return nil
}

func validateSegmentPath(name string) error {
	if name == "" || strings.Contains(name, "\\") || path.IsAbs(name) ||
		path.Clean(name) != name || strings.HasPrefix(name, "../") ||
		path.Ext(name) != ".toml" || !isSafeSegmentName(name) {
		return fmt.Errorf("invalid registry segment path %q", name)
	}
	return nil
}

func isSafeSegmentName(name string) bool {
	for _, char := range name {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') || char == '/' || char == '-' ||
			char == '_' || char == '.' {
			continue
		}
		return false
	}
	return true
}
