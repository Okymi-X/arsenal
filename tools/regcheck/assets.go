package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Okymi-X/arsenal/internal/registry"
)

// verifyAsset checks that a fetchable asset resolves upstream: that its latest
// release contains a matching file, or that its raw build directory exists.
func (c *checker) verifyAsset(a registry.Asset) error {
	owner, repo, err := parseGitHubRepo(a.Repo)
	if err != nil {
		return err
	}
	switch a.Source {
	case registry.AssetGitHubRelease:
		return c.checkReleaseAsset(owner, repo, a.Pattern)
	case registry.AssetGitHubRaw:
		return c.checkRawDir(owner, repo, assetBranch(a), a.Dir)
	default:
		return fmt.Errorf("unknown asset source %q", a.Source)
	}
}

func assetBranch(a registry.Asset) string {
	if a.Branch != "" {
		return a.Branch
	}
	return "master"
}

func (c *checker) checkReleaseAsset(owner, repo, pattern string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	var payload struct {
		Assets []struct {
			Name string `json:"name"`
		} `json:"assets"`
	}
	if err := c.githubJSON(url, &payload); err != nil {
		return err
	}
	if pattern == "" {
		if len(payload.Assets) == 0 {
			return fmt.Errorf("latest release of %s/%s has no assets", owner, repo)
		}
		return nil
	}
	want := strings.ToLower(pattern)
	for _, as := range payload.Assets {
		if strings.Contains(strings.ToLower(as.Name), want) {
			return nil
		}
	}
	return fmt.Errorf("no asset matches %q in latest release of %s/%s", pattern, owner, repo)
}

func (c *checker) checkRawDir(owner, repo, branch, dir string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", owner, repo, dir, branch)
	var entries []struct {
		Name string `json:"name"`
	}
	if err := c.githubJSON(url, &entries); err != nil {
		return err
	}
	if len(entries) == 0 {
		return fmt.Errorf("directory %q is empty in %s/%s", dir, owner, repo)
	}
	return nil
}

// githubJSON performs a GET and decodes a JSON body, mapping common GitHub
// status codes to clear errors.
func (c *checker) githubJSON(url string, v any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.githubToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.githubToken)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("query %s: %w", url, err)
	}
	defer func() { _ = resp.Body.Close() }()
	switch resp.StatusCode {
	case http.StatusOK:
		return json.NewDecoder(resp.Body).Decode(v)
	case http.StatusNotFound:
		return fmt.Errorf("not found: %s", url)
	case http.StatusForbidden:
		return fmt.Errorf("GitHub rate limit hit; set GITHUB_TOKEN")
	default:
		return fmt.Errorf("GitHub returned status %d for %s", resp.StatusCode, url)
	}
}
