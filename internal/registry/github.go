package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"
)

const maxGitHubResponse = 1 << 20

var githubRepoPattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$`)

// GitHubRef is a repository tag and the commit it currently identifies.
type GitHubRef struct {
	Name   string
	Commit string
}

// GitHubClient resolves human-readable refs before registry files are pulled.
type GitHubClient struct {
	repo    string
	file    string
	token   string
	apiBase string
	rawBase string
	http    *http.Client
}

// NewGitHubClient constructs a read-only GitHub registry client.
func NewGitHubClient(repo, file, token string) (*GitHubClient, error) {
	if !githubRepoPattern.MatchString(repo) {
		return nil, fmt.Errorf("invalid GitHub repository %q (want owner/repo)", repo)
	}
	cleanFile := path.Clean(file)
	if cleanFile == "." || strings.HasPrefix(cleanFile, "../") || path.IsAbs(cleanFile) || !isSafeSegmentName(cleanFile) {
		return nil, fmt.Errorf("invalid registry path %q", file)
	}
	return &GitHubClient{
		repo:    repo,
		file:    cleanFile,
		token:   token,
		apiBase: "https://api.github.com",
		rawBase: "https://raw.githubusercontent.com",
		http:    &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// ParseGitHubRepo extracts owner/repo from a canonical HTTPS GitHub URL.
func ParseGitHubRepo(repoURL string) (string, error) {
	parsed, err := url.Parse(repoURL)
	if err != nil || parsed.Scheme != "https" || !strings.EqualFold(parsed.Host, "github.com") ||
		parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", fmt.Errorf("repository %q is not an HTTPS GitHub URL", repoURL)
	}
	repo := strings.Trim(strings.TrimSuffix(parsed.Path, ".git"), "/")
	if !githubRepoPattern.MatchString(repo) {
		return "", fmt.Errorf("repository %q does not identify owner/repo", repoURL)
	}
	return repo, nil
}

// ListTags returns up to the requested number of repository tags.
func (c *GitHubClient) ListTags(ctx context.Context, limit int) ([]GitHubRef, error) {
	if limit < 1 || limit > 100 {
		return nil, fmt.Errorf("GitHub tag limit must be between 1 and 100")
	}
	endpoint := fmt.Sprintf("%s/repos/%s/tags?per_page=%d", c.apiBase, c.repo, limit)
	var response []struct {
		Name   string `json:"name"`
		Commit struct {
			SHA string `json:"sha"`
		} `json:"commit"`
	}
	if err := c.getJSON(ctx, endpoint, &response); err != nil {
		return nil, err
	}
	refs := make([]GitHubRef, 0, len(response))
	for _, tag := range response {
		if tag.Name == "" || !isCommitSHA(tag.Commit.SHA) {
			return nil, fmt.Errorf("GitHub returned an invalid tag response")
		}
		refs = append(refs, GitHubRef{Name: tag.Name, Commit: tag.Commit.SHA})
	}
	return refs, nil
}

// ResolveCommit resolves a tag, branch, or SHA to an immutable commit SHA.
func (c *GitHubClient) ResolveCommit(ctx context.Context, ref string) (string, error) {
	if strings.TrimSpace(ref) == "" {
		return "", fmt.Errorf("GitHub ref is empty")
	}
	endpoint := fmt.Sprintf("%s/repos/%s/commits/%s", c.apiBase, c.repo, url.PathEscape(ref))
	var response struct {
		SHA string `json:"sha"`
	}
	if err := c.getJSON(ctx, endpoint, &response); err != nil {
		return "", err
	}
	if !isCommitSHA(response.SHA) {
		return "", fmt.Errorf("GitHub returned an invalid commit SHA")
	}
	return response.SHA, nil
}

// RegistryURL returns the raw registry manifest URL pinned to commit.
func (c *GitHubClient) RegistryURL(commit string) (string, error) {
	if !isCommitSHA(commit) {
		return "", fmt.Errorf("invalid GitHub commit SHA %q", commit)
	}
	return fmt.Sprintf("%s/%s/%s/%s", c.rawBase, c.repo, commit, c.file), nil
}

func (c *GitHubClient) getJSON(ctx context.Context, endpoint string, destination any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("build GitHub request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "arsenal")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("GitHub request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if !sameOrigin(req.URL, resp.Request.URL) {
		return fmt.Errorf("GitHub request redirect changed remote origin")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub request failed with status %d", resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxGitHubResponse+1))
	if err != nil {
		return fmt.Errorf("read GitHub response: %w", err)
	}
	if len(data) > maxGitHubResponse {
		return fmt.Errorf("GitHub response exceeds %d bytes", maxGitHubResponse)
	}
	if err := json.Unmarshal(data, destination); err != nil {
		return fmt.Errorf("decode GitHub response: %w", err)
	}
	return nil
}

func isCommitSHA(value string) bool {
	if len(value) != 40 {
		return false
	}
	for _, character := range value {
		if (character < '0' || character > '9') && (character < 'a' || character > 'f') && (character < 'A' || character > 'F') {
			return false
		}
	}
	return true
}
