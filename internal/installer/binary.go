package installer

import (
	"context"
	"fmt"

	"github.com/Okymi-X/arsenal/internal/registry"
)

// BinaryMethod installs a prebuilt release binary for a tool.
//
// TODO(arsenal#2): Implement release-asset download and checksum verification.
// Tracking: docs/architecture.md "Binary install method".
type BinaryMethod struct{}

// NewBinaryMethod returns a BinaryMethod.
func NewBinaryMethod() *BinaryMethod { return &BinaryMethod{} }

// Supports reports whether the tool declares the binary install method.
func (m *BinaryMethod) Supports(tool registry.Tool) bool {
	return tool.InstallMethod == MethodBinary
}

// Install is not implemented yet.
func (m *BinaryMethod) Install(ctx context.Context, tool registry.Tool, version registry.Version) error {
	return fmt.Errorf("binary install method not implemented for %q yet", tool.Name)
}

var _ InstallMethod = (*BinaryMethod)(nil)
