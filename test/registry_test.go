// Package test holds integration tests that exercise multiple packages
// together, including the curated registry data shipped in the binary.
package test

import (
	"testing"

	"github.com/Okymi-X/arsenal/internal/registry"
	"github.com/Okymi-X/arsenal/internal/resolver"
	builtin "github.com/Okymi-X/arsenal/registry"
)

// TestEmbeddedRegistryParses guards the curated data file: it must parse,
// validate, and expose the flagship tools used across engagements.
func TestEmbeddedRegistryParses(t *testing.T) {
	reg, err := registry.Parse(builtin.Bytes())
	if err != nil {
		t.Fatalf("embedded registry failed to parse: %v", err)
	}
	if len(reg.Tools) == 0 {
		t.Fatal("embedded registry has no tools")
	}
	for _, name := range []string{"netexec", "impacket", "certipy"} {
		if _, ok := reg.FindTool(name); !ok {
			t.Errorf("expected tool %q in registry", name)
		}
	}
}

// TestEmbeddedRegistryResolves checks that aliases and default selection work
// against the real data end to end.
func TestEmbeddedRegistryResolves(t *testing.T) {
	reg, err := registry.Parse(builtin.Bytes())
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	r := resolver.New(reg)
	res, err := r.Resolve(resolver.Request{Tool: "nxc"})
	if err != nil {
		t.Fatalf("resolve nxc: %v", err)
	}
	if res.Tool.Name != "netexec" {
		t.Fatalf("alias nxc resolved to %q", res.Tool.Name)
	}
	if !res.Version.Tested {
		t.Fatalf("default version %q should be tested", res.Version.Tag)
	}
}
