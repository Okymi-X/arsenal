package resolver

import (
	"testing"

	"github.com/Okymi-X/arsenal/internal/registry"
)

func sampleRegistry() *registry.Registry {
	return &registry.Registry{
		Version: "1",
		Tools: []registry.Tool{
			{
				Name:    "netexec",
				Aliases: []string{"nxc"},
				Versions: []registry.Version{
					{Tag: "2.0.0", Tested: false},
					{Tag: "1.4.0", Tested: true},
					{Tag: "1.3.0", Tested: true},
				},
			},
			{
				Name: "untested",
				Versions: []registry.Version{
					{Tag: "0.9.0", Tested: false},
				},
			},
		},
	}
}

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name    string
		spec    string
		want    Request
		wantErr bool
	}{
		{"name only", "netexec", Request{Tool: "netexec"}, false},
		{"name and version", "netexec@1.4.0", Request{Tool: "netexec", Version: "1.4.0"}, false},
		{"alias", "nxc@1.3.0", Request{Tool: "nxc", Version: "1.3.0"}, false},
		{"empty", "", Request{}, true},
		{"trailing at", "netexec@", Request{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRequest(tt.spec)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestResolve(t *testing.T) {
	r := New(sampleRegistry())
	tests := []struct {
		name        string
		req         Request
		wantTool    string
		wantVersion string
		wantErr     bool
	}{
		{"explicit version", Request{Tool: "netexec", Version: "1.3.0"}, "netexec", "1.3.0", false},
		{"default picks newest tested", Request{Tool: "netexec"}, "netexec", "1.4.0", false},
		{"alias resolves", Request{Tool: "nxc"}, "netexec", "1.4.0", false},
		{"default falls back to newest when none tested", Request{Tool: "untested"}, "untested", "0.9.0", false},
		{"unknown tool", Request{Tool: "ghost"}, "", "", true},
		{"unknown version", Request{Tool: "netexec", Version: "9.9.9"}, "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.Resolve(tt.req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got.Tool.Name != tt.wantTool || got.Version.Tag != tt.wantVersion {
				t.Fatalf("got %s@%s, want %s@%s", got.Tool.Name, got.Version.Tag, tt.wantTool, tt.wantVersion)
			}
		})
	}
}
