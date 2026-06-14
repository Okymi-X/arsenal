package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Okymi-X/arsenal/internal/registry"
)

// checkPyPI verifies that a package version is published on PyPI.
func (c *checker) checkPyPI(pkg, version string) error {
	url := fmt.Sprintf("https://pypi.org/pypi/%s/json", pkg)
	resp, err := c.http.Get(url)
	if err != nil {
		return fmt.Errorf("query PyPI for %s: %w", pkg, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("package %q not found on PyPI", pkg)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("PyPI returned status %d for %s", resp.StatusCode, pkg)
	}
	var payload struct {
		Releases map[string]json.RawMessage `json:"releases"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return fmt.Errorf("decode PyPI response for %s: %w", pkg, err)
	}
	if _, ok := payload.Releases[version]; !ok {
		return fmt.Errorf("version %q not published on PyPI", version)
	}
	return nil
}

// checkGitRef verifies that at least one candidate ref resolves in the tool's
// GitHub repository.
func (c *checker) checkGitRef(tool registry.Tool, v registry.Version) error {
	owner, repo, err := parseGitHubRepo(tool.Repo)
	if err != nil {
		return err
	}
	refs := refCandidates(v)
	if len(refs) == 0 {
		return fmt.Errorf("no commit or tag to verify")
	}
	for _, ref := range refs {
		ok, err := c.gitHubRefExists(owner, repo, ref)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
	}
	return fmt.Errorf("none of the refs %v resolve in %s/%s", refs, owner, repo)
}

// gitHubRefExists reports whether a ref (commit, tag, or branch) resolves.
func (c *checker) gitHubRefExists(owner, repo, ref string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", owner, repo, ref)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.githubToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.githubToken)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return false, fmt.Errorf("query GitHub for %s/%s@%s: %w", owner, repo, ref, err)
	}
	defer func() { _ = resp.Body.Close() }()
	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound, http.StatusUnprocessableEntity:
		return false, nil
	case http.StatusForbidden:
		return false, fmt.Errorf("GitHub rate limit hit; set GITHUB_TOKEN")
	default:
		return false, fmt.Errorf("GitHub returned status %d for %s", resp.StatusCode, ref)
	}
}

// parseGitHubRepo extracts owner and repo from a GitHub URL.
func parseGitHubRepo(repoURL string) (owner, repo string, err error) {
	const host = "github.com/"
	i := strings.Index(repoURL, host)
	if i < 0 {
		return "", "", fmt.Errorf("not a GitHub repo URL: %q", repoURL)
	}
	parts := strings.Split(strings.Trim(repoURL[i+len(host):], "/"), "/")
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("cannot parse owner/repo from %q", repoURL)
	}
	return parts[0], strings.TrimSuffix(parts[1], ".git"), nil
}
