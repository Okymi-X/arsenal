package registry

import (
	"fmt"
	"sort"
	"strings"
)

// AssetSource classifies where a fetchable binary is pulled from.
type AssetSource string

// Recognized asset sources.
const (
	// AssetGitHubRelease pulls a named asset from a repository's latest
	// GitHub release.
	AssetGitHubRelease AssetSource = "github-release"
	// AssetGitHubRaw pulls a raw file from a repository branch. The branch
	// HEAD is always the latest version, so no pin is recorded.
	AssetGitHubRaw AssetSource = "github-raw"
)

// Asset is a precompiled binary that arsenal fetches to a destination
// directory for staging onto a target.
//
// Unlike a Tool, an Asset is never isolated in a virtualenv or exposed via a
// shim on the operator's PATH. arsenal pulls its latest version and writes the
// file where the operator asks.
type Asset struct {
	// Name is the canonical identifier used on the command line.
	Name string `toml:"name"`
	// Aliases are alternative names that resolve to this asset.
	Aliases []string `toml:"aliases"`
	// Description is a one-line summary.
	Description string `toml:"description"`
	// Repo is the source repository URL.
	Repo string `toml:"repo"`
	// Category places the asset in an offensive-security domain.
	Category Category `toml:"category"`
	// Source selects how the binary is located: github-release or github-raw.
	Source AssetSource `toml:"source"`
	// Pattern selects the file for a single-binary asset. For github-release
	// it matches a release asset filename (for example "winPEASx64.exe"); for
	// github-raw it is unused. A binary argument on the command line overrides
	// it.
	Pattern string `toml:"pattern"`
	// Branch is the git ref for github-raw fetches; defaults to "master".
	Branch string `toml:"branch"`
	// Dir is the in-repo directory holding the selectable binaries of a
	// github-raw collection (for example a SharpCollection build folder).
	Dir string `toml:"dir"`
	// Builds are the selectable Dir values for a collection, used to validate
	// a --build override. Empty means any value is accepted.
	Builds []string `toml:"builds"`
	// Tags are free-form search keywords.
	Tags []string `toml:"tags"`
	// Notes holds operational guidance.
	Notes string `toml:"notes"`
}

// Collection reports whether the asset is a github-raw directory of many
// selectable binaries (such as SharpCollection) rather than a single file.
func (a Asset) Collection() bool {
	return a.Source == AssetGitHubRaw && a.Dir != ""
}

// FindAsset returns the asset matching name or one of its aliases.
//
// Matching is case-insensitive. The returned bool is false when no asset
// matches, mirroring the comma-ok idiom for lookups.
func (r *Registry) FindAsset(name string) (Asset, bool) {
	want := strings.ToLower(strings.TrimSpace(name))
	for _, a := range r.Assets {
		if strings.ToLower(a.Name) == want {
			return a, true
		}
		for _, alias := range a.Aliases {
			if strings.ToLower(alias) == want {
				return a, true
			}
		}
	}
	return Asset{}, false
}

// SearchAssets returns assets whose name, description, category, tags, or
// aliases match the query, ordered by name. An empty query returns all assets.
func (r *Registry) SearchAssets(query string) []Asset {
	q := strings.ToLower(strings.TrimSpace(query))
	var out []Asset
	for _, a := range r.Assets {
		if q == "" || assetMatches(a, q) {
			out = append(out, a)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func assetMatches(a Asset, q string) bool {
	if strings.Contains(strings.ToLower(a.Name), q) ||
		strings.Contains(strings.ToLower(a.Description), q) ||
		strings.EqualFold(string(a.Category), q) {
		return true
	}
	for _, tag := range a.Tags {
		if strings.Contains(strings.ToLower(tag), q) {
			return true
		}
	}
	for _, alias := range a.Aliases {
		if strings.Contains(strings.ToLower(alias), q) {
			return true
		}
	}
	return false
}

// MustFindAsset is FindAsset that returns an error instead of a bool.
func (r *Registry) MustFindAsset(name string) (Asset, error) {
	a, ok := r.FindAsset(name)
	if !ok {
		return Asset{}, fmt.Errorf("asset %q not found in registry", name)
	}
	return a, nil
}
