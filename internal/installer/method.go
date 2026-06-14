// Package installer selects and drives the right installation method for a
// tool: pip, git+pip, prebuilt binary, go install, or cargo install.
//
// Each method is a small implementation of InstallMethod. The orchestrator
// owns method selection; methods themselves only know how to install. Methods
// are constructed with their dependencies (the isolation backend) injected,
// honoring dependency inversion.
package installer

import (
	"context"

	"github.com/Okymi-X/arsenal/internal/registry"
)

// InstallMethod knows how to install one class of tool.
type InstallMethod interface {
	// Supports reports whether this method handles the given tool.
	Supports(tool registry.Tool) bool
	// Install provisions the tool at the given version.
	Install(ctx context.Context, tool registry.Tool, version registry.Version) error
}

// Method name constants matching the registry install_method field.
const (
	MethodPip    = "pip"
	MethodGitPip = "gitpip"
	MethodBinary = "binary"
	MethodGoBin  = "gobin"
	MethodCargo  = "cargo"
)
