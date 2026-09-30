package installer

import (
	"context"
	"fmt"

	"github.com/Okymi-X/arsenal/internal/isolation"
	"github.com/Okymi-X/arsenal/internal/registry"
)

// GitPipMethod installs a tool from a Git repository, pinned by commit, into
// an isolation backend using pip's VCS support.
type GitPipMethod struct {
	backend isolation.Backend
}

// NewGitPipMethod returns a GitPipMethod bound to the given backend.
func NewGitPipMethod(backend isolation.Backend) *GitPipMethod {
	return &GitPipMethod{backend: backend}
}

// Supports reports whether the tool declares the git+pip install method.
func (m *GitPipMethod) Supports(tool registry.Tool) bool {
	return tool.InstallMethod == MethodGitPip
}

// Install creates the environment and installs from the pinned Git revision.
func (m *GitPipMethod) Install(ctx context.Context, tool registry.Tool, version registry.Version) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	repo := tool.RepoFor(version)
	if repo == "" {
		return Result{}, fmt.Errorf("tool %q has no repo for git+pip install", tool.Name)
	}
	if err := m.backend.Create(ctx, tool.Name, version.Tag); err != nil {
		return Result{}, fmt.Errorf("provision environment: %w", err)
	}
	install := isolation.InstallSpec{
		GitURL:   repo,
		Commit:   version.Commit,
		PipSpecs: append([]string(nil), version.PipDependencies...),
	}
	if err := m.backend.Install(ctx, install); err != nil {
		return Result{}, fmt.Errorf("install %s via git+pip: %w", tool.Name, err)
	}
	return Result{Path: m.backend.Path()}, nil
}

var _ InstallMethod = (*GitPipMethod)(nil)
