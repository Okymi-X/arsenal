package venv

import (
	"os"
	"strings"
)

func (b *Backend) environment() []string {
	return virtualEnvironment(os.Environ(), b.binDir(), b.dir)
}

func virtualEnvironment(base []string, binDir, virtualEnv string) []string {
	out := make([]string, 0, len(base)+3)
	hostPath := ""
	for _, entry := range base {
		name, value, _ := strings.Cut(entry, "=")
		switch name {
		case "PATH":
			hostPath = value
		case "VIRTUAL_ENV", "PYTHONHOME", "PYTHONPATH", "PYTHONNOUSERSITE":
		default:
			out = append(out, entry)
		}
	}
	path := binDir
	if hostPath != "" {
		path += string(os.PathListSeparator) + hostPath
	}
	return append(out, "VIRTUAL_ENV="+virtualEnv, "PYTHONNOUSERSITE=1", "PATH="+path)
}

func cleanPythonEnvironment(base []string) []string {
	out := make([]string, 0, len(base)+1)
	for _, entry := range base {
		name, _, _ := strings.Cut(entry, "=")
		switch name {
		case "VIRTUAL_ENV", "PYTHONHOME", "PYTHONPATH", "PYTHONNOUSERSITE":
		default:
			out = append(out, entry)
		}
	}
	return append(out, "PYTHONNOUSERSITE=1")
}
