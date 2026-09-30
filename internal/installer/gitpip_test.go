package installer

import (
	"context"
	"testing"

	"github.com/Okymi-X/arsenal/internal/registry"
)

func TestGitPipInstallIncludesPinnedDependencies(t *testing.T) {
	backend := &fakeBackend{}
	method := NewGitPipMethod(backend)
	tool := registry.Tool{
		Name:          "netexec",
		Repo:          "https://github.com/Pennyw0rth/NetExec",
		InstallMethod: MethodGitPip,
	}
	version := registry.Version{
		Tag:             "1.5.1",
		Commit:          "v1.5.1",
		PipDependencies: []string{"dploot==3.1.3"},
	}
	if _, err := method.Install(context.Background(), tool, version); err != nil {
		t.Fatalf("Install: %v", err)
	}
	if len(backend.installed) != 1 {
		t.Fatalf("install calls = %d, want 1", len(backend.installed))
	}
	got := backend.installed[0]
	if got.GitURL != tool.Repo || got.Commit != version.Commit || len(got.PipSpecs) != 1 || got.PipSpecs[0] != "dploot==3.1.3" {
		t.Fatalf("install spec = %#v", got)
	}
}
