package venv

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/Okymi-X/arsenal/internal/isolation"
)

func TestInstallVerifiesPythonDependencies(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test helper uses POSIX scripts")
	}
	dir := t.TempDir()
	binDir := filepath.Join(dir, "bin")
	if err := os.Mkdir(binDir, 0o755); err != nil {
		t.Fatal(err)
	}
	logPath := filepath.Join(dir, "commands.log")
	writeCommandStub(t, filepath.Join(binDir, "pip"), logPath)
	writeCommandStub(t, filepath.Join(binDir, "python"), logPath)

	backend := &Backend{dir: dir}
	if err := backend.Install(context.Background(), isolation.InstallSpec{PipSpecs: []string{"example==1.0.0"}}); err != nil {
		t.Fatalf("Install: %v", err)
	}
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	if !strings.Contains(got, "install --no-input example==1.0.0\n") || !strings.Contains(got, "-m pip check\n") {
		t.Fatalf("commands = %q", got)
	}
}

func writeCommandStub(t *testing.T, path, logPath string) {
	t.Helper()
	script := "#!/bin/sh\nprintf '%s\\n' \"$*\" >> " + logPath + "\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
}
