package installer

import (
	"context"
	"fmt"

	"github.com/Okymi-X/arsenal/internal/isolation"
	"github.com/Okymi-X/arsenal/internal/registry"
)

// PipMethod installs a tool from a PyPI package into an isolation backend.
type PipMethod struct {
	backend isolation.Backend
}

// NewPipMethod returns a PipMethod bound to the given backend.
func NewPipMethod(backend isolation.Backend) *PipMethod {
	return &PipMethod{backend: backend}
}

// Supports reports whether the tool declares the pip install method.
func (m *PipMethod) Supports(tool registry.Tool) bool {
	return tool.InstallMethod == MethodPip
}

// Install creates the environment and installs the pinned pip spec.
func (m *PipMethod) Install(ctx context.Context, tool registry.Tool, version registry.Version) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	spec := version.PipSpec
	if spec == "" {
		spec = fmt.Sprintf("%s==%s", tool.Name, version.Tag)
	}
	if err := m.backend.Create(tool.Name, version.Tag); err != nil {
		return fmt.Errorf("provision environment: %w", err)
	}
	install := isolation.InstallSpec{PipSpecs: []string{spec}}
	if err := m.backend.Install(install); err != nil {
		return fmt.Errorf("install %s via pip: %w", tool.Name, err)
	}
	return nil
}

var _ InstallMethod = (*PipMethod)(nil)
