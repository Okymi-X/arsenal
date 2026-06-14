// Package op manages engagement profiles ("ops"): named sets of pinned tool
// versions that make an environment reproducible and shareable across a team.
//
// An Op is the editable definition (which tools, which versions). A Lockfile
// is the resolved, fully pinned snapshot generated from an Op plus the
// registry, capturing commits and pip specs for byte-stable reproduction.
package op

// Op is an engagement profile: a named set of tool/version pins.
type Op struct {
	// Name identifies the op.
	Name string `toml:"name"`
	// Description is a short human summary.
	Description string `toml:"description"`
	// Created is the RFC-3339 creation timestamp.
	Created string `toml:"created"`
	// Pins are the requested tool/version pairs.
	Pins []Pin `toml:"pin"`
}

// Pin is a single requested tool/version in an Op.
type Pin struct {
	// Tool is the tool name.
	Tool string `toml:"tool"`
	// Version is the pinned version tag.
	Version string `toml:"version"`
}

// Lockfile is the resolved snapshot of an Op, fully pinned for reproduction.
type Lockfile struct {
	// Op is the originating op name.
	Op string `toml:"op"`
	// Generated is the RFC-3339 generation timestamp.
	Generated string `toml:"generated"`
	// RegistryVersion records the registry schema version used.
	RegistryVersion string `toml:"registry_version"`
	// Entries are the resolved, pinned tool installations.
	Entries []LockEntry `toml:"entry"`
}

// LockEntry is a fully resolved, pinned tool installation.
type LockEntry struct {
	// Tool is the tool name.
	Tool string `toml:"tool"`
	// Version is the resolved version tag.
	Version string `toml:"version"`
	// Commit pins the exact source revision, when known.
	Commit string `toml:"commit"`
	// PipSpec is the resolved pip requirement, when applicable.
	PipSpec string `toml:"pip_spec"`
	// InstallMethod records how the tool is installed.
	InstallMethod string `toml:"install_method"`
}

// SetPin adds or updates a pin for a tool.
func (o *Op) SetPin(tool, version string) {
	for i := range o.Pins {
		if o.Pins[i].Tool == tool {
			o.Pins[i].Version = version
			return
		}
	}
	o.Pins = append(o.Pins, Pin{Tool: tool, Version: version})
}
