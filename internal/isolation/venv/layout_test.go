package venv

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateRejectsUnsafeComponents(t *testing.T) {
	backend := New("python", t.TempDir())
	for _, test := range []struct {
		tool    string
		version string
	}{
		{tool: "../outside", version: "1.0.0"},
		{tool: "example", version: "../outside"},
	} {
		if err := backend.Create(context.Background(), test.tool, test.version); err == nil {
			t.Fatalf("Create(%q, %q) succeeded", test.tool, test.version)
		}
	}
}

func TestCreateRejectsSymlinkedToolDirectory(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, "example")); err != nil {
		t.Fatal(err)
	}
	backend := New("python", root)
	if err := backend.Create(context.Background(), "example", "1.0.0"); err == nil {
		t.Fatal("Create succeeded through symlinked tool directory")
	}
}

func TestCreateRejectsSymlinkedVersionDirectory(t *testing.T) {
	root := t.TempDir()
	toolDir := filepath.Join(root, "example")
	if err := os.Mkdir(toolDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(t.TempDir(), filepath.Join(toolDir, "1.0.0")); err != nil {
		t.Fatal(err)
	}
	backend := New("python", root)
	if err := backend.Create(context.Background(), "example", "1.0.0"); err == nil {
		t.Fatal("Create succeeded through symlinked version directory")
	}
}
