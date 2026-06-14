// Package isolation defines the Backend interface that isolates a tool's
// runtime from the host and from other tools.
//
// Two implementations exist: a Python virtualenv backend (the default,
// covering most tools) and a container backend for tools with heavy system
// dependencies. Concrete backends are injected, never referenced across
// package boundaries by name.
package isolation

// InstallSpec describes what to install into an isolated environment.
//
// It is intentionally backend-agnostic: a venv backend interprets it as pip
// arguments, while a container backend interprets it as build instructions.
type InstallSpec struct {
	// PipSpecs are pip requirement strings, e.g. "netexec==1.1.0".
	PipSpecs []string
	// GitURL is a source repository to install from, when set.
	GitURL string
	// Commit pins the revision for a Git install.
	Commit string
	// Editable installs the Git checkout in editable mode.
	Editable bool
	// Extras are optional pip arguments appended verbatim.
	Extras []string
}

// Backend isolates a single tool/version environment.
//
// A Backend instance is bound to one tool and version. Create provisions the
// environment, Install populates it, Run executes a binary inside it, Remove
// tears it down, Path reports its location, and Exists reports provisioning.
type Backend interface {
	// Create provisions the isolated environment for tool at version.
	Create(tool, version string) error
	// Install populates the environment according to spec.
	Install(spec InstallSpec) error
	// Run executes args inside the environment, wiring through stdio.
	Run(args []string) error
	// Remove tears down the environment.
	Remove() error
	// Path returns the environment's root directory.
	Path() string
	// Exists reports whether the environment is provisioned.
	Exists() bool
}
