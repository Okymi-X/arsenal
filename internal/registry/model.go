// Package registry defines the curated catalog of offensive-security tools
// and the operations to load, query, and synchronize it.
//
// The registry is the product's core value: a hand-maintained, tested mapping
// from tool name to known-good versions, pinned by commit, annotated with the
// required Python version and operational notes.
package registry

// Category classifies a tool by its offensive-security domain.
type Category string

// Recognized tool categories.
const (
	CategoryAD       Category = "ad"
	CategoryWeb      Category = "web"
	CategoryRecon    Category = "recon"
	CategoryPassword Category = "password"
	CategoryExploit  Category = "exploit"
	CategoryC2       Category = "c2"
	CategoryMisc     Category = "misc"
)

// Registry is the top-level catalog parsed from registry.toml.
type Registry struct {
	// Version is the registry schema version.
	Version string `toml:"version"`
	// Updated is the ISO-8601 date the registry was last revised.
	Updated string `toml:"updated"`
	// Tools is the curated list of catalog entries.
	Tools []Tool `toml:"tool"`
}

// Tool is a single catalog entry describing an installable program.
type Tool struct {
	// Name is the canonical identifier used on the command line.
	Name string `toml:"name"`
	// Aliases are alternative names that resolve to this tool.
	Aliases []string `toml:"aliases"`
	// Description is a one-line summary.
	Description string `toml:"description"`
	// Repo is the source repository URL.
	Repo string `toml:"repo"`
	// Homepage is the project's primary URL.
	Homepage string `toml:"homepage"`
	// Category places the tool in an offensive-security domain.
	Category Category `toml:"category"`
	// InstallMethod selects the installer: pip, gitpip, binary, gobin, cargo.
	InstallMethod string `toml:"install_method"`
	// PythonVersion is the minimum interpreter, e.g. ">=3.9".
	PythonVersion string `toml:"python_version"`
	// Binary is the primary executable exposed via a shim.
	Binary string `toml:"binary"`
	// Binaries are additional executables exposed via shims.
	Binaries []string `toml:"binaries"`
	// Dependencies are extra system or pip requirements, as notes.
	Dependencies []string `toml:"dependencies"`
	// Tags are free-form search keywords.
	Tags []string `toml:"tags"`
	// RequiresRoot indicates the tool typically needs elevated privileges.
	RequiresRoot bool `toml:"requires_root"`
	// Notes holds operational guidance.
	Notes string `toml:"notes"`
	// Versions is the list of catalogued releases, newest first.
	Versions []Version `toml:"version"`
}

// Version is a single catalogued release of a tool.
type Version struct {
	// Tag is the human-facing version label, e.g. "1.1.0".
	Tag string `toml:"tag"`
	// Commit pins the exact source revision (commit SHA, tag, or branch).
	Commit string `toml:"commit"`
	// Repo optionally overrides the tool's repo for this version, e.g. to pin
	// a fork or a branch. Empty means use the tool's repo.
	Repo string `toml:"repo"`
	// Tested marks a version the maintainers verified as known-good.
	Tested bool `toml:"tested"`
	// PipSpec is the pip requirement string, e.g. "netexec==1.1.0".
	PipSpec string `toml:"pip_spec"`
	// Date is the ISO-8601 release date.
	Date string `toml:"date"`
	// Notes holds version-specific guidance.
	Notes string `toml:"notes"`
}

// RepoFor returns the effective source repository for a version: the version's
// repo override when set, otherwise the tool's repo.
func (t Tool) RepoFor(v Version) string {
	if v.Repo != "" {
		return v.Repo
	}
	return t.Repo
}

// AllBinaries returns the primary binary plus any additional binaries.
func (t Tool) AllBinaries() []string {
	out := make([]string, 0, len(t.Binaries)+1)
	if t.Binary != "" {
		out = append(out, t.Binary)
	}
	out = append(out, t.Binaries...)
	return out
}
