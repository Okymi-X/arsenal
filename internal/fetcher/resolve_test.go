package fetcher

import (
	"testing"

	"github.com/Okymi-X/arsenal/internal/registry"
)

func TestSelectedDirAllowsOnlyCataloguedBuilds(t *testing.T) {
	asset := registry.Asset{
		Name:   "collection",
		Dir:    "default",
		Builds: []string{"allowed"},
	}
	tests := []struct {
		name     string
		override string
		want     string
		wantErr  bool
	}{
		{name: "default", want: "default"},
		{name: "catalogued", override: "allowed", want: "allowed"},
		{name: "uncatalogued", override: "other", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := selectedDir(asset, test.override)
			if (err != nil) != test.wantErr {
				t.Fatalf("selectedDir error = %v, wantErr %v", err, test.wantErr)
			}
			if got != test.want {
				t.Fatalf("selectedDir = %q, want %q", got, test.want)
			}
		})
	}
}

func TestSelectedDirRejectsOverrideWithoutBuildCatalog(t *testing.T) {
	asset := registry.Asset{Name: "single", Dir: "dist"}
	if _, err := selectedDir(asset, "arbitrary"); err == nil {
		t.Fatal("expected override without a build catalog to be rejected")
	}
}
