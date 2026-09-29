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

// TestEmbeddedRegistryIncludesVerifiedExpansion guards the installable tools
// added from the Exegol reference set and their expected shim names.
func TestEmbeddedRegistryIncludesVerifiedExpansion(t *testing.T) {
	reg, err := registry.Parse(builtin.Bytes())
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	want := map[string]string{
		"bbot":              "bbot",
		"bloodhound-import": "bloodhound-import",
		"censys":            "censys",
		"fierce":            "fierce",
		"holehe":            "holehe",
		"ldeep":             "ldeep",
		"maigret":           "maigret",
		"name-that-hash":    "nth",
		"sherlock-project":  "sherlock",
		"ssh-audit":         "ssh-audit",
	}
	for name, binary := range want {
		tool, ok := reg.FindTool(name)
		if !ok {
			t.Errorf("expected tool %q in registry", name)
			continue
		}
		if tool.InstallMethod != "pip" || tool.Binary != binary {
			t.Errorf("%s: install=%q binary=%q", name, tool.InstallMethod, tool.Binary)
		}
		if len(tool.Versions) == 0 || !tool.Versions[0].Tested {
			t.Errorf("%s: newest version is not tested", name)
		}
	}
}
