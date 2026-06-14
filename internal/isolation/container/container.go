// Package container is the container-based isolation backend for tools with
// heavy system dependencies, using podman or docker.
//
// TODO(arsenal#1): Implement the container backend. The interface is wired so
// the orchestrator can select it once ready. Tracking: docs/architecture.md
// "Container backend". Until implemented, every method returns a clear
// not-implemented error so misconfiguration fails loudly rather than silently.
package container

import (
	"fmt"

	"github.com/Okymi-X/arsenal/internal/isolation"
)

// Backend is the not-yet-implemented container isolation backend.
//
// It satisfies isolation.Backend so it can be injected today, but every
// operation reports that the backend is unavailable.
type Backend struct {
	// Runtime is the container runtime binary, "podman" or "docker".
	Runtime string
}

// New returns a container Backend using the given runtime binary.
func New(runtime string) *Backend { return &Backend{Runtime: runtime} }

var errNotImplemented = fmt.Errorf("container backend not implemented yet (see docs/architecture.md)")

// Create is not implemented.
func (b *Backend) Create(tool, version string) error { return errNotImplemented }

// Install is not implemented.
func (b *Backend) Install(spec isolation.InstallSpec) error { return errNotImplemented }

// Run is not implemented.
func (b *Backend) Run(args []string) error { return errNotImplemented }

// Remove is not implemented.
func (b *Backend) Remove() error { return errNotImplemented }

// Path returns an empty string; no environment exists.
func (b *Backend) Path() string { return "" }

// Exists always reports false.
func (b *Backend) Exists() bool { return false }

var _ isolation.Backend = (*Backend)(nil)
