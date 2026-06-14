package installer

import (
	"context"
	"fmt"

	"github.com/Okymi-X/arsenal/internal/registry"
)

// CargoMethod installs a Rust tool via "cargo install" pinned to a version.
//
// TODO(arsenal#4): Implement cargo install into a per-tool root.
// Tracking: docs/architecture.md "Cargo install method".
type CargoMethod struct{}

// NewCargoMethod returns a CargoMethod.
func NewCargoMethod() *CargoMethod { return &CargoMethod{} }

// Supports reports whether the tool declares the cargo install method.
func (m *CargoMethod) Supports(tool registry.Tool) bool {
	return tool.InstallMethod == MethodCargo
}

// Install is not implemented yet.
func (m *CargoMethod) Install(ctx context.Context, tool registry.Tool, version registry.Version) error {
	return fmt.Errorf("cargo install method not implemented for %q yet", tool.Name)
}

var _ InstallMethod = (*CargoMethod)(nil)
