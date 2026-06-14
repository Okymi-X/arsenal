package installer

import (
	"context"
	"fmt"

	"github.com/Okymi-X/arsenal/internal/registry"
)

// GoBinMethod installs a tool via "go install" pinned to a version.
//
// TODO(arsenal#3): Implement go install into a per-tool GOBIN.
// Tracking: docs/architecture.md "Go install method".
type GoBinMethod struct{}

// NewGoBinMethod returns a GoBinMethod.
func NewGoBinMethod() *GoBinMethod { return &GoBinMethod{} }

// Supports reports whether the tool declares the gobin install method.
func (m *GoBinMethod) Supports(tool registry.Tool) bool {
	return tool.InstallMethod == MethodGoBin
}

// Install is not implemented yet.
func (m *GoBinMethod) Install(ctx context.Context, tool registry.Tool, version registry.Version) error {
	return fmt.Errorf("go install method not implemented for %q yet", tool.Name)
}

var _ InstallMethod = (*GoBinMethod)(nil)
