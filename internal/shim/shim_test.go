package shim

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteRejectsPathTraversal(t *testing.T) {
	manager := NewManager(filepath.Join(t.TempDir(), "bin"))
	for _, name := range []string{"../outside", "sub/tool", "bad\nname"} {
		if err := manager.Write(name, "/safe/tool"); err == nil {
			t.Fatalf("shim name %q unexpectedly accepted", name)
		}
	}
}

func TestWriteQuotesTargetAsShellData(t *testing.T) {
	dir := t.TempDir()
	manager := NewManager(dir)
	target := `/tmp/a'$(touch /tmp/not-run)`
	if err := manager.Write("safe-tool", target); err != nil {
		t.Fatalf("Write: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(dir, "safe-tool"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), `'/tmp/a'"'"'$(touch /tmp/not-run)'`) {
		t.Fatalf("target was not safely quoted: %s", content)
	}
}
