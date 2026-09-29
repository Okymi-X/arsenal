package installer

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Okymi-X/arsenal/internal/registry"
)

type runnerFunc func(context.Context, Command) error

func (f runnerFunc) Run(ctx context.Context, command Command) error { return f(ctx, command) }

func TestGoBinInstallUsesPinnedIsolatedCommand(t *testing.T) {
	root := t.TempDir()
	var got Command
	runner := runnerFunc(func(_ context.Context, command Command) error {
		got = command
		return writeExecutable(command.Env["GOBIN"], "ffuf")
	})
	method := NewGoBinMethod(root, runner)
	tool := registry.Tool{Name: "ffuf", InstallMethod: MethodGoBin, Binary: "ffuf"}
	version := registry.Version{
		Tag: "2.1.0", Commit: "v2.1.0",
		InstallTargets: map[string]string{"ffuf": "github.com/ffuf/ffuf/v2"},
	}

	result, err := method.Install(context.Background(), tool, version)
	if err != nil {
		t.Fatalf("Install: %v", err)
	}
	wantPath := filepath.Join(root, "ffuf", "2.1.0")
	if result.Path != wantPath || result.Backend != MethodGoBin {
		t.Fatalf("result = %#v", result)
	}
	if got.Name != "go" || strings.Join(got.Args, " ") != "install github.com/ffuf/ffuf/v2@v2.1.0" {
		t.Fatalf("command = %#v", got)
	}
	for _, key := range []string{"GOBIN", "GOCACHE", "GOMODCACHE", "GOPATH"} {
		if !strings.HasPrefix(got.Env[key], filepath.Dir(wantPath)) {
			t.Fatalf("%s is not scoped to staging: %q", key, got.Env[key])
		}
	}
	assertExecutable(t, filepath.Join(wantPath, "bin", "ffuf"))
	if _, err := os.Stat(filepath.Join(wantPath, "cache")); !os.IsNotExist(err) {
		t.Fatalf("cache retained after install: %v", err)
	}
}

func TestGoBinInstallCleansFailureAndKeepsExistingInstall(t *testing.T) {
	root := t.TempDir()
	final := filepath.Join(root, "ffuf", "2.1.0")
	if err := os.MkdirAll(final, 0o755); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(final, "existing")
	if err := os.WriteFile(marker, []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}
	method := NewGoBinMethod(root, runnerFunc(func(context.Context, Command) error {
		return errors.New("download failed")
	}))
	tool := registry.Tool{Name: "ffuf", Binary: "ffuf", InstallMethod: MethodGoBin}
	version := registry.Version{Tag: "2.1.0", Commit: "v2.1.0", InstallTargets: map[string]string{
		"ffuf": "github.com/ffuf/ffuf/v2",
	}}

	if _, err := method.Install(context.Background(), tool, version); err == nil {
		t.Fatal("expected install failure")
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("existing install changed: %v", err)
	}
	entries, err := os.ReadDir(filepath.Dir(final))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != filepath.Base(final) {
		t.Fatalf("staging directory was not cleaned: %v", entries)
	}
}

func TestRemoveInstallTreeHandlesReadOnlyCache(t *testing.T) {
	root := filepath.Join(t.TempDir(), "cache")
	nested := filepath.Join(root, "module", "source")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, "file"), []byte("data"), 0o444); err != nil {
		t.Fatal(err)
	}
	for _, dir := range []string{nested, filepath.Join(root, "module"), root} {
		if err := os.Chmod(dir, 0o555); err != nil {
			t.Fatal(err)
		}
	}
	if err := removeInstallTree(root); err != nil {
		t.Fatalf("removeInstallTree: %v", err)
	}
	if _, err := os.Stat(root); !os.IsNotExist(err) {
		t.Fatalf("read-only tree still exists: %v", err)
	}
}

func TestGoBinInstallRejectsUnsafeInputs(t *testing.T) {
	called := false
	method := NewGoBinMethod(t.TempDir(), runnerFunc(func(context.Context, Command) error {
		called = true
		return nil
	}))
	tool := registry.Tool{Name: "../escape", Binary: "tool", InstallMethod: MethodGoBin}
	version := registry.Version{Tag: "1.0.0", Commit: "v1.0.0", InstallTargets: map[string]string{
		"tool": "github.com/owner/tool",
	}}
	if _, err := method.Install(context.Background(), tool, version); err == nil {
		t.Fatal("expected unsafe path to be rejected")
	}
	if called {
		t.Fatal("runner called for rejected input")
	}
}

func TestGoBinInstallRejectsMovingRef(t *testing.T) {
	method := NewGoBinMethod(t.TempDir(), runnerFunc(func(context.Context, Command) error {
		t.Fatal("runner must not be called")
		return nil
	}))
	tool := registry.Tool{Name: "ffuf", Binary: "ffuf", InstallMethod: MethodGoBin}
	version := registry.Version{Tag: "main", Commit: "main", InstallTargets: map[string]string{
		"ffuf": "github.com/ffuf/ffuf/v2",
	}}
	if _, err := method.Install(context.Background(), tool, version); err == nil {
		t.Fatal("expected moving ref to be rejected")
	}
}

func TestGoBinaryNameHandlesModuleMajorSuffix(t *testing.T) {
	if got := goBinaryName("github.com/OJ/gobuster/v3"); got != "gobuster" {
		t.Fatalf("goBinaryName = %q, want gobuster", got)
	}
	if got := goBinaryName("github.com/example/tool/cmd/agent"); got != "agent" {
		t.Fatalf("goBinaryName = %q, want agent", got)
	}
}

func writeExecutable(dir, name string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, name), []byte("test"), 0o755)
}

func assertExecutable(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm()&0o111 == 0 {
		t.Fatalf("expected executable %s: %v", path, err)
	}
}
