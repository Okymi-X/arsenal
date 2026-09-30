package main

import "testing"

func TestSelectBuildVersion(t *testing.T) {
	tests := []struct {
		name          string
		injected      string
		moduleVersion string
		want          string
	}{
		{name: "injected release", injected: "v1.2.3", moduleVersion: "v1.2.2", want: "v1.2.3"},
		{name: "go install module", injected: "dev", moduleVersion: "v1.2.3", want: "v1.2.3"},
		{name: "development build", injected: "dev", moduleVersion: "(devel)", want: "dev"},
		{name: "missing metadata", want: "dev"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := selectBuildVersion(test.injected, test.moduleVersion)
			if got != test.want {
				t.Fatalf("selectBuildVersion() = %q, want %q", got, test.want)
			}
		})
	}
}
