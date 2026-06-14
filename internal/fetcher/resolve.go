package fetcher

import (
	"context"
	"fmt"

	"github.com/Okymi-X/arsenal/internal/registry"
)

// resolve locates the single upstream file an asset and selection refer to,
// returning its file name, download URL, and version label.
func (f *Fetcher) resolve(ctx context.Context, asset registry.Asset, sel Selection) (name, url, version string, err error) {
	switch asset.Source {
	case registry.AssetGitHubRelease:
		return f.resolveRelease(ctx, asset, sel)
	case registry.AssetGitHubRaw:
		return f.resolveRaw(ctx, asset, sel)
	default:
		return "", "", "", fmt.Errorf("unknown asset source %q", asset.Source)
	}
}

func (f *Fetcher) resolveRelease(ctx context.Context, asset registry.Asset, sel Selection) (string, string, string, error) {
	rel, err := f.latestRelease(ctx, asset.Repo)
	if err != nil {
		return "", "", "", err
	}
	want := sel.Binary
	if want == "" {
		want = asset.Pattern
	}
	if want == "" {
		return "", "", "", fmt.Errorf("%s: no binary given and asset has no default pattern", asset.Name)
	}
	chosen, ok := matchName(want, assetNames(rel.Assets))
	if !ok {
		return "", "", "", fmt.Errorf("%s: no release asset matches %q (try --list)", asset.Name, want)
	}
	for _, a := range rel.Assets {
		if a.Name == chosen {
			return a.Name, a.URL, rel.Tag, nil
		}
	}
	return "", "", "", fmt.Errorf("%s: release asset %q vanished", asset.Name, chosen)
}

func (f *Fetcher) resolveRaw(ctx context.Context, asset registry.Asset, sel Selection) (string, string, string, error) {
	if sel.Binary == "" {
		return "", "", "", fmt.Errorf("%s: a binary name is required (try --list)", asset.Name)
	}
	if !buildAllowed(asset, sel.Build) {
		return "", "", "", fmt.Errorf("%s: unknown build %q (try --list of builds in notes)", asset.Name, sel.Build)
	}
	dir := dirOf(asset, sel.Build)
	entries, err := f.repoDir(ctx, asset.Repo, branchOf(asset), dir)
	if err != nil {
		return "", "", "", err
	}
	chosen, ok := matchName(sel.Binary, fileNames(entries))
	if !ok {
		return "", "", "", fmt.Errorf("%s: no binary matches %q in %s (try --list)", asset.Name, sel.Binary, dir)
	}
	for _, e := range entries {
		if e.Name == chosen {
			return e.Name, e.URL, "latest", nil
		}
	}
	return "", "", "", fmt.Errorf("%s: entry %q vanished", asset.Name, chosen)
}

func branchOf(a registry.Asset) string {
	if a.Branch != "" {
		return a.Branch
	}
	return "master"
}

func dirOf(a registry.Asset, override string) string {
	if override != "" {
		return override
	}
	return a.Dir
}

func buildAllowed(a registry.Asset, override string) bool {
	if override == "" || len(a.Builds) == 0 {
		return true
	}
	for _, b := range a.Builds {
		if b == override {
			return true
		}
	}
	return false
}
