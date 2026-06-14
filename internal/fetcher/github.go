package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const apiBase = "https://api.github.com"

type releaseAsset struct {
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
	Size int64  `json:"size"`
}

type release struct {
	Tag    string         `json:"tag_name"`
	Assets []releaseAsset `json:"assets"`
}

type contentEntry struct {
	Name string `json:"name"`
	Type string `json:"type"`
	URL  string `json:"download_url"`
	Size int64  `json:"size"`
}

// latestRelease returns the repository's latest GitHub release.
func (f *Fetcher) latestRelease(ctx context.Context, repo string) (release, error) {
	or, err := ownerRepo(repo)
	if err != nil {
		return release{}, err
	}
	var rel release
	err = f.getJSON(ctx, apiBase+"/repos/"+or+"/releases/latest", &rel)
	return rel, err
}

// repoDir returns the entries of a directory in a repository at a branch.
func (f *Fetcher) repoDir(ctx context.Context, repo, branch, dir string) ([]contentEntry, error) {
	or, err := ownerRepo(repo)
	if err != nil {
		return nil, err
	}
	url := apiBase + "/repos/" + or + "/contents/" + dir
	if branch != "" {
		url += "?ref=" + branch
	}
	var entries []contentEntry
	err = f.getJSON(ctx, url, &entries)
	return entries, err
}

func (f *Fetcher) getJSON(ctx context.Context, url string, v any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "arsenal")
	if f.token != "" {
		req.Header.Set("Authorization", "Bearer "+f.token)
	}
	resp, err := f.client.Do(req)
	if err != nil {
		return fmt.Errorf("github request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("github %s: %s: %s", url, resp.Status, strings.TrimSpace(string(body)))
	}
	return json.NewDecoder(resp.Body).Decode(v)
}

// ownerRepo extracts the "owner/name" slug from a GitHub repository URL.
func ownerRepo(repoURL string) (string, error) {
	s := strings.TrimSuffix(strings.TrimSpace(repoURL), "/")
	s = strings.TrimSuffix(s, ".git")
	s = strings.TrimPrefix(s, "https://github.com/")
	s = strings.TrimPrefix(s, "http://github.com/")
	if s == "" || strings.Count(s, "/") != 1 {
		return "", fmt.Errorf("not a github repository url: %q", repoURL)
	}
	return s, nil
}

// checksumSuffixes are release-asset extensions that are never the binary
// itself, so they are hidden from listing and matching.
var checksumSuffixes = []string{".sha256", ".sha1", ".sha512", ".md5", ".asc", ".sig"}

func assetNames(assets []releaseAsset) []string {
	out := make([]string, 0, len(assets))
	for _, a := range assets {
		if isChecksum(a.Name) {
			continue
		}
		out = append(out, a.Name)
	}
	return out
}

func isChecksum(name string) bool {
	l := strings.ToLower(name)
	for _, s := range checksumSuffixes {
		if strings.HasSuffix(l, s) {
			return true
		}
	}
	return false
}

func fileNames(entries []contentEntry) []string {
	var out []string
	for _, e := range entries {
		if e.Type == "file" && e.Name != "README.md" {
			out = append(out, e.Name)
		}
	}
	return out
}
