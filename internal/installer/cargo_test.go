package installer

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Okymi-X/arsenal/internal/registry"
)

func TestCargoInstallUsesExactVersionAndScopedState(t *testing.T) {
	root := t.TempDir()
	var got Command
	runner := runnerFunc(func(_ context.Context, command Command) error {
		got = command
		return writeExecutable(filepath.Join(command.Args[2], "bin"), "rustscan")
	})
	method := NewCargoMethod(root, runner)
	tool := registry.Tool{Name: "rustscan", InstallMethod: MethodCargo, Binary: "rustscan"}
	version := registry.Version{Tag: "2.4.1", InstallTargets: map[string]string{"rustscan": "rustscan"}}

	result, err := method.Install(context.Background(), tool, version)
	if err != nil {
		t.Fatalf("Install: %v", err)
	}
	wantPath := filepath.Join(root, "rustscan", "2.4.1")
	if result.Path != wantPath || result.Backend != MethodCargo {
		t.Fatalf("result = %#v", result)
	}
	wantArgs := "install --root " + got.Args[2] + " --version 2.4.1 --locked rustscan"
	if got.Name != "cargo" || strings.Join(got.Args, " ") != wantArgs {
		t.Fatalf("command = %#v", got)
	}
	for _, key := range []string{"CARGO_HOME", "CARGO_TARGET_DIR"} {
		if !strings.Contains(got.Env[key], ".install-") {
			t.Fatalf("%s is not staging-scoped: %q", key, got.Env[key])
		}
	}
	assertExecutable(t, filepath.Join(wantPath, "bin", "rustscan"))
	if _, err := os.Stat(filepath.Join(wantPath, "cache")); !os.IsNotExist(err) {
		t.Fatalf("cache retained after install: %v", err)
	}
}

func TestCargoInstallRejectsMovingOrUnsafeVersions(t *testing.T) {
	method := NewCargoMethod(t.TempDir(), runnerFunc(func(context.Context, Command) error {
		t.Fatal("runner must not be called")
		return nil
	}))
	tool := registry.Tool{Name: "rustscan", InstallMethod: MethodCargo, Binary: "rustscan"}
	for _, tag := range []string{"latest", "2.4", "../2.4.1"} {
		version := registry.Version{Tag: tag, InstallTargets: map[string]string{"rustscan": "rustscan"}}
		if _, err := method.Install(context.Background(), tool, version); err == nil {
			t.Fatalf("version %q unexpectedly accepted", tag)
		}
	}
}

func TestCargoGitHubCommandUsesResolvedCommit(t *testing.T) {
	commit := strings.Repeat("a", 40)
	command, err := cargoInstallCommand("/stage", registry.Tool{
		Repo: "https://github.com/RustScan/RustScan",
	}, registry.Version{Tag: "github-" + commit[:12], Commit: commit}, installTarget{
		Binary: "rustscan", Source: "rustscan",
	}, true)
	if err != nil {
		t.Fatalf("cargoInstallCommand: %v", err)
	}
	want := "install --root /stage --git https://github.com/RustScan/RustScan --rev " + commit + " --locked rustscan"
	if got := strings.Join(command.Args, " "); got != want {
		t.Fatalf("args = %q, want %q", got, want)
	}
}
