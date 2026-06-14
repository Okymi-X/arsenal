package installer

import (
	"context"
	"testing"

	"github.com/Okymi-X/arsenal/internal/isolation"
	"github.com/Okymi-X/arsenal/internal/registry"
)

// fakeBackend records the calls made against it.
type fakeBackend struct {
	created   bool
	installed []isolation.InstallSpec
}

func (f *fakeBackend) Create(tool, version string) error { f.created = true; return nil }
func (f *fakeBackend) Install(spec isolation.InstallSpec) error {
	f.installed = append(f.installed, spec)
	return nil
}
func (f *fakeBackend) Run(args []string) error { return nil }
func (f *fakeBackend) Remove() error           { return nil }
func (f *fakeBackend) Path() string            { return "/fake" }
func (f *fakeBackend) Exists() bool            { return true }

func TestSelect(t *testing.T) {
	backend := &fakeBackend{}
	orch := NewOrchestrator(DefaultMethods(backend))
	tests := []struct {
		name   string
		method string
		want   string
	}{
		{"pip", MethodPip, "*installer.PipMethod"},
		{"gitpip", MethodGitPip, "*installer.GitPipMethod"},
		{"binary", MethodBinary, "*installer.BinaryMethod"},
		{"gobin", MethodGoBin, "*installer.GoBinMethod"},
		{"cargo", MethodCargo, "*installer.CargoMethod"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := registry.Tool{Name: "t", InstallMethod: tt.method}
			m, err := orch.Select(tool)
			if err != nil {
				t.Fatalf("Select: %v", err)
			}
			if got := typeName(m); got != tt.want {
				t.Fatalf("Select method = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestSelectUnknown(t *testing.T) {
	orch := NewOrchestrator(DefaultMethods(&fakeBackend{}))
	if _, err := orch.Select(registry.Tool{Name: "t", InstallMethod: "nope"}); err == nil {
		t.Fatal("expected error for unknown method")
	}
}

func TestPipInstallUsesBackend(t *testing.T) {
	backend := &fakeBackend{}
	orch := NewOrchestrator(DefaultMethods(backend))
	tool := registry.Tool{Name: "netexec", InstallMethod: MethodPip}
	ver := registry.Version{Tag: "1.4.0", PipSpec: "netexec==1.4.0"}
	if err := orch.Install(context.Background(), tool, ver); err != nil {
		t.Fatalf("Install: %v", err)
	}
	if !backend.created {
		t.Fatal("expected backend.Create to be called")
	}
	if len(backend.installed) != 1 || backend.installed[0].PipSpecs[0] != "netexec==1.4.0" {
		t.Fatalf("unexpected install specs: %+v", backend.installed)
	}
}

func typeName(v any) string {
	switch v.(type) {
	case *PipMethod:
		return "*installer.PipMethod"
	case *GitPipMethod:
		return "*installer.GitPipMethod"
	case *BinaryMethod:
		return "*installer.BinaryMethod"
	case *GoBinMethod:
		return "*installer.GoBinMethod"
	case *CargoMethod:
		return "*installer.CargoMethod"
	default:
		return "unknown"
	}
}
