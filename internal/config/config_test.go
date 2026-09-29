package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveConfigRoundTripAndPermissions(t *testing.T) {
	paths := Paths{Root: t.TempDir()}
	want := DefaultConfig()
	want.RegistryRef = "v1.0.0"
	want.RegistryCommit = "0123456789abcdef0123456789abcdef01234567"
	if err := paths.SaveConfig(want); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}
	got, err := paths.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if got != want {
		t.Fatalf("LoadConfig = %#v, want %#v", got, want)
	}
	info, err := os.Stat(filepath.Join(paths.Root, "config.json"))
	if err != nil {
		t.Fatalf("stat config: %v", err)
	}
	if gotMode := info.Mode().Perm(); gotMode != 0o600 {
		t.Fatalf("config permissions = %o, want 600", gotMode)
	}
}
