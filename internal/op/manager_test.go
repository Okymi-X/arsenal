package op

import (
	"os"
	"path/filepath"
	"testing"
)

func TestManagerRejectsUnsafeNames(t *testing.T) {
	manager := NewManager(filepath.Join(t.TempDir(), "ops"))
	for _, name := range []string{"", ".", "..", "../escape", "nested/name", "name with space"} {
		t.Run(name, func(t *testing.T) {
			if _, err := manager.Path(name); err == nil {
				t.Fatalf("Path(%q) unexpectedly succeeded", name)
			}
			if _, err := manager.LockPath(name); err == nil {
				t.Fatalf("LockPath(%q) unexpectedly succeeded", name)
			}
			if err := manager.Create(&Op{Name: name}); err == nil {
				t.Fatalf("Create(%q) unexpectedly succeeded", name)
			}
		})
	}
}

func TestManagerSaveDoesNotFollowPredictableTemporarySymlink(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "ops")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(t.TempDir(), "outside")
	if err := os.WriteFile(outside, []byte("unchanged"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(dir, "safe.toml.tmp")); err != nil {
		t.Fatal(err)
	}
	manager := NewManager(dir)
	if err := manager.Save(&Op{Name: "safe"}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	data, err := os.ReadFile(outside)
	if err != nil || string(data) != "unchanged" {
		t.Fatalf("symlink target changed: %q, %v", data, err)
	}
}

func TestManagerLoadRejectsMismatchedDeclaredName(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "expected.toml"), []byte("name = \"other\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := NewManager(dir).Load("expected"); err == nil {
		t.Fatal("expected mismatched op name to be rejected")
	}
}

func TestManagerLoadRejectsUnknownFields(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "expected.toml"), []byte("name = \"expected\"\nunknown = true\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := NewManager(dir).Load("expected"); err == nil {
		t.Fatal("expected unknown op field to be rejected")
	}
}
