package cli

import (
	"testing"

	"github.com/Okymi-X/arsenal/internal/op"
	"github.com/Okymi-X/arsenal/internal/registry"
)

func TestVerifyLockfileAcceptsMatchingRegistry(t *testing.T) {
	reg := lockTestRegistry()
	lockfile := lockTestFile()
	resolved, err := verifyLockfile(reg, lockfile)
	if err != nil {
		t.Fatalf("verifyLockfile: %v", err)
	}
	if len(resolved) != 1 || resolved[0].Tool.Name != "example" {
		t.Fatalf("resolved = %#v", resolved)
	}
}

func TestVerifyLockfileRejectsDriftBeforeInstall(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*op.Lockfile)
	}{
		{name: "registry version", mutate: func(lf *op.Lockfile) { lf.RegistryVersion = "2" }},
		{name: "install method", mutate: func(lf *op.Lockfile) { lf.Entries[0].InstallMethod = "gitpip" }},
		{name: "commit", mutate: func(lf *op.Lockfile) { lf.Entries[0].Commit = "other" }},
		{name: "package spec", mutate: func(lf *op.Lockfile) { lf.Entries[0].PipSpec = "example==2.0.0" }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lockfile := lockTestFile()
			test.mutate(lockfile)
			if _, err := verifyLockfile(lockTestRegistry(), lockfile); err == nil {
				t.Fatal("expected metadata drift to be rejected")
			}
		})
	}
}

func lockTestRegistry() *registry.Registry {
	return &registry.Registry{
		Version: "1",
		Tools: []registry.Tool{{
			Name:          "example",
			InstallMethod: "pip",
			Versions: []registry.Version{{
				Tag:     "1.0.0",
				Commit:  "v1.0.0",
				PipSpec: "example==1.0.0",
			}},
		}},
	}
}

func lockTestFile() *op.Lockfile {
	return &op.Lockfile{
		Op:              "eng",
		RegistryVersion: "1",
		Entries: []op.LockEntry{{
			Tool:          "example",
			Version:       "1.0.0",
			Commit:        "v1.0.0",
			PipSpec:       "example==1.0.0",
			InstallMethod: "pip",
		}},
	}
}
