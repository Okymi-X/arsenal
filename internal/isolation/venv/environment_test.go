package venv

import (
	"os"
	"slices"
	"testing"
)

func TestVirtualEnvironmentRemovesHostPythonOverrides(t *testing.T) {
	base := []string{
		"PATH=/host/bin",
		"VIRTUAL_ENV=/host/venv",
		"PYTHONHOME=/host/python",
		"PYTHONPATH=/host/modules",
		"PYTHONNOUSERSITE=0",
		"KEEP=value",
	}
	want := []string{
		"KEEP=value",
		"VIRTUAL_ENV=/arsenal/tool",
		"PYTHONNOUSERSITE=1",
		"PATH=/arsenal/tool/bin" + string(os.PathListSeparator) + "/host/bin",
	}
	if got := virtualEnvironment(base, "/arsenal/tool/bin", "/arsenal/tool"); !slices.Equal(got, want) {
		t.Fatalf("virtualEnvironment = %#v, want %#v", got, want)
	}
}

func TestCleanPythonEnvironment(t *testing.T) {
	base := []string{"PATH=/host/bin", "PYTHONPATH=/host/modules", "KEEP=value"}
	want := []string{"PATH=/host/bin", "KEEP=value", "PYTHONNOUSERSITE=1"}
	if got := cleanPythonEnvironment(base); !slices.Equal(got, want) {
		t.Fatalf("cleanPythonEnvironment = %#v, want %#v", got, want)
	}
}
