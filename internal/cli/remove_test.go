package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Okymi-X/arsenal/internal/config"
	"github.com/Okymi-X/arsenal/internal/shim"
	"github.com/Okymi-X/arsenal/internal/store"
)

func TestTearDownRemovesOnlyRecordedOwnedPath(t *testing.T) {
	paths := config.NewPaths(t.TempDir())
	app := &App{paths: paths, shims: shim.NewManager(paths.Bin)}
	install := paths.ToolDir("ffuf", "2.1.0")
	if err := os.MkdirAll(filepath.Join(install, "bin"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := app.shims.Write("ffuf", filepath.Join(install, "bin", "ffuf")); err != nil {
		t.Fatal(err)
	}
	record := store.InstalledTool{
		Name: "ffuf", Version: "2.1.0", Path: install, Binaries: []string{"ffuf"},
	}
	if err := app.tearDown(record); err != nil {
		t.Fatalf("tearDown: %v", err)
	}
	if _, err := os.Stat(install); !os.IsNotExist(err) {
		t.Fatalf("installation still exists: %v", err)
	}
	if app.shims.Exists("ffuf") {
		t.Fatal("shim still exists")
	}
}

func TestTearDownRejectsManifestPathOutsideToolsRoot(t *testing.T) {
	paths := config.NewPaths(t.TempDir())
	app := &App{paths: paths, shims: shim.NewManager(paths.Bin)}
	outside := filepath.Join(t.TempDir(), "outside")
	if err := os.Mkdir(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	record := store.InstalledTool{Name: "ffuf", Version: "2.1.0", Path: outside}
	if err := app.tearDown(record); err == nil {
		t.Fatal("expected untrusted path to be rejected")
	}
	if _, err := os.Stat(outside); err != nil {
		t.Fatalf("outside path was changed: %v", err)
	}
}
