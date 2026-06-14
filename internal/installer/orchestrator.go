package installer

import (
	"context"
	"fmt"

	"github.com/Okymi-X/arsenal/internal/isolation"
	"github.com/Okymi-X/arsenal/internal/registry"
)

// Orchestrator selects the appropriate InstallMethod for a tool and drives it.
//
// It holds an ordered list of methods and delegates to the first one whose
// Supports returns true. Methods are injected, so tests can supply fakes.
type Orchestrator struct {
	methods []InstallMethod
}

// NewOrchestrator returns an Orchestrator over the given methods.
func NewOrchestrator(methods []InstallMethod) *Orchestrator {
	return &Orchestrator{methods: methods}
}

// DefaultMethods builds the standard set of methods bound to a backend.
//
// pip and git+pip use the backend; binary, go, and cargo are independent.
func DefaultMethods(backend isolation.Backend) []InstallMethod {
	return []InstallMethod{
		NewPipMethod(backend),
		NewGitPipMethod(backend),
		NewBinaryMethod(),
		NewGoBinMethod(),
		NewCargoMethod(),
	}
}

// Select returns the first method that supports the tool.
func (o *Orchestrator) Select(tool registry.Tool) (InstallMethod, error) {
	for _, m := range o.methods {
		if m.Supports(tool) {
			return m, nil
		}
	}
	return nil, fmt.Errorf("no install method supports %q (method %q)", tool.Name, tool.InstallMethod)
}

// Install selects a method for the tool and installs the version.
func (o *Orchestrator) Install(ctx context.Context, tool registry.Tool, version registry.Version) error {
	m, err := o.Select(tool)
	if err != nil {
		return err
	}
	return m.Install(ctx, tool, version)
}
