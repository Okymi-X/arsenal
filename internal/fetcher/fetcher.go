// Package fetcher pulls precompiled binaries from upstream and writes them to a
// destination directory for staging onto a target.
//
// Unlike the installer, the fetcher performs no isolation and creates no shims:
// it locates the latest version of a release asset or a raw repository file and
// downloads it where the operator asks.
package fetcher

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Okymi-X/arsenal/internal/registry"
)

// Fetcher downloads assets from GitHub.
type Fetcher struct {
	client *http.Client
	token  string
}

// New returns a Fetcher. An optional GitHub token raises API rate limits and is
// used only for api.github.com listing calls, never for downloads.
func New(token string) *Fetcher {
	return &Fetcher{
		client: &http.Client{Timeout: 60 * time.Second},
		token:  token,
	}
}

// Selection narrows what to fetch and where to write it.
type Selection struct {
	// Binary names the file to fetch, overriding the asset's default pattern.
	// It is required for a collection asset.
	Binary string
	// Build overrides a collection asset's directory, such as a SharpCollection
	// .NET build folder.
	Build string
	// DestDir is the directory the file is written into.
	DestDir string
}

// Result describes a completed fetch.
type Result struct {
	// File is the written file's base name.
	File string
	// Path is the path written.
	Path string
	// Size is the number of bytes written.
	Size int64
	// Version is the upstream release tag, or "latest" for a branch file.
	Version string
}

// Fetch resolves the asset to a single upstream file and downloads it into
// sel.DestDir, returning what was written.
func (f *Fetcher) Fetch(ctx context.Context, asset registry.Asset, sel Selection) (Result, error) {
	name, url, version, err := f.resolve(ctx, asset, sel)
	if err != nil {
		return Result{}, err
	}
	if err := os.MkdirAll(sel.DestDir, 0o755); err != nil {
		return Result{}, fmt.Errorf("create dest %s: %w", sel.DestDir, err)
	}
	// Take the base name only: the file name comes from upstream, so this keeps
	// a crafted name from escaping the destination directory.
	name = filepath.Base(name)
	dest := filepath.Join(sel.DestDir, name)
	size, err := f.download(ctx, url, dest)
	if err != nil {
		return Result{}, err
	}
	return Result{File: name, Path: dest, Size: size, Version: version}, nil
}

// List returns the candidate file names available for an asset: the names in
// its latest release, or the entries in a collection's build directory.
func (f *Fetcher) List(ctx context.Context, asset registry.Asset, build string) ([]string, error) {
	switch asset.Source {
	case registry.AssetGitHubRelease:
		rel, err := f.latestRelease(ctx, asset.Repo)
		if err != nil {
			return nil, err
		}
		return assetNames(rel.Assets), nil
	case registry.AssetGitHubRaw:
		entries, err := f.repoDir(ctx, asset.Repo, branchOf(asset), dirOf(asset, build))
		if err != nil {
			return nil, err
		}
		return fileNames(entries), nil
	default:
		return nil, fmt.Errorf("unknown asset source %q", asset.Source)
	}
}
