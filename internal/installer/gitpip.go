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
func (m *GitPipMethod) Install(ctx context.Context, tool registry.Tool, version registry.Version) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if tool.Repo == "" {
		return fmt.Errorf("tool %q has no repo for git+pip install", tool.Name)
	}
	if err := m.backend.Create(tool.Name, version.Tag); err != nil {
		return fmt.Errorf("provision environment: %w", err)
	}
	install := isolation.InstallSpec{
		GitURL: tool.Repo,
		Commit: version.Commit,
	}
	if err := m.backend.Install(install); err != nil {
		return fmt.Errorf("install %s via git+pip: %w", tool.Name, err)
	}
	return nil
}

var _ InstallMethod = (*GitPipMethod)(nil)
