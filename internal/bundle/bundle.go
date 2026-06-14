// Package bundle exports a self-contained, offline bundle of installed tools
// for transfer to isolated or air-gapped networks.
//
// TODO(arsenal#5): Implement bundle export/import. The export must capture the
// resolved lockfile plus vendored wheels and source archives so a target host
// can reconstruct environments with no network. Tracking: docs/architecture.md
// "Offline bundling". The Exporter interface is defined now so the CLI can be
// wired ahead of the implementation.
package bundle

import "fmt"

// Exporter produces and consumes offline bundles.
type Exporter interface {
	// Export writes a self-contained bundle for the named op to destDir.
	Export(opName, destDir string) error
	// Import reconstructs environments from a bundle at srcPath.
	Import(srcPath string) error
}

// StubExporter is the not-yet-implemented bundle exporter.
type StubExporter struct{}

// NewStubExporter returns a StubExporter.
func NewStubExporter() *StubExporter { return &StubExporter{} }

// Export is not implemented yet.
func (e *StubExporter) Export(opName, destDir string) error {
	return fmt.Errorf("offline bundle export not implemented yet (see docs/architecture.md)")
}

// Import is not implemented yet.
func (e *StubExporter) Import(srcPath string) error {
	return fmt.Errorf("offline bundle import not implemented yet (see docs/architecture.md)")
}

var _ Exporter = (*StubExporter)(nil)
