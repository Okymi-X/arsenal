package registry

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"
)

func TestExpandSegmentedRegistry(t *testing.T) {
	segment := []byte(`
[[tool]]
name = "example"
category = "recon"
install_method = "pip"
binary = "example"

  [[tool.version]]
  tag = "1.0.0"
  tested = true
  pip_spec = "example==1.0.0"
`)
	manifest := testManifest("segments/tools.toml", segment)

	data, err := Expand(manifest, func(name string) ([]byte, error) {
		if name != "segments/tools.toml" {
			return nil, fmt.Errorf("unexpected segment %q", name)
		}
		return segment, nil
	})
	if err != nil {
		t.Fatalf("expand: %v", err)
	}
	reg, err := Parse(data)
	if err != nil {
		t.Fatalf("parse expanded registry: %v", err)
	}
	if len(reg.Tools) != 1 || reg.Tools[0].Name != "example" {
		t.Fatalf("unexpected tools: %#v", reg.Tools)
	}
}

func TestExpandRejectsMissingOrMismatchedChecksum(t *testing.T) {
	segment := []byte("[[tool]]\nname = \"example\"\n")
	missing := []byte("version = \"1\"\nsegments = [\"segments/tools.toml\"]\n")
	if _, err := Expand(missing, func(string) ([]byte, error) { return segment, nil }); err == nil {
		t.Fatal("expected missing checksum to be rejected")
	}
	tampered := testManifest("segments/tools.toml", []byte("original"))
	if _, err := Expand(tampered, func(string) ([]byte, error) { return segment, nil }); err == nil {
		t.Fatal("expected mismatched checksum to be rejected")
	}
}

func TestParseManifestRejectsUnsafeSegments(t *testing.T) {
	tests := []struct {
		name     string
		segments string
	}{
		{name: "parent traversal", segments: `"../outside.toml"`},
		{name: "absolute", segments: `"/tmp/outside.toml"`},
		{name: "backslash", segments: `"segments\\outside.toml"`},
		{name: "encoded traversal", segments: `"segments/%2e%2e/outside.toml"`},
		{name: "wrong extension", segments: `"segments/tools.txt"`},
		{name: "duplicate", segments: `"segments/tools.toml", "segments/tools.toml"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := fmt.Appendf(nil, "version = %q\nsegments = [%s]\n", "1", tt.segments)
			if _, err := ParseManifest(data); err == nil {
				t.Fatal("expected manifest to be rejected")
			}
		})
	}
}

func TestExpandRejectsOversizedRegistry(t *testing.T) {
	manifest := []byte("version = \"1\"\nsegments = [\"segments/large.toml\"]\n")
	_, err := Expand(manifest, func(string) ([]byte, error) {
		return []byte(strings.Repeat("x", maxAssembledSize+1)), nil
	})
	if err == nil {
		t.Fatal("expected oversized registry to be rejected")
	}
}

func testManifest(name string, segment []byte) []byte {
	digest := sha256.Sum256(segment)
	return fmt.Appendf(nil, "version = %q\nsegments = [%q]\n[segment_sha256]\n%q = %q\n",
		"1", name, name, fmt.Sprintf("%x", digest))
}
