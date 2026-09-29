package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Okymi-X/arsenal/internal/registry"
)

const (
	apiBase           = "https://api.github.com"
	maxGitHubAPIBytes = 4 << 20
)

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
	or, err := registry.ParseGitHubRepo(repo)
	if err != nil {
		return release{}, err
	}
	var rel release
	err = f.getJSON(ctx, apiBase+"/repos/"+or+"/releases/latest", &rel)
	return rel, err
}

// repoDir returns the entries of a directory in a repository at a branch.
func (f *Fetcher) repoDir(ctx context.Context, repo, branch, dir string) ([]contentEntry, error) {
	or, err := registry.ParseGitHubRepo(repo)
	if err != nil {
		return nil, err
	}
	var entries []contentEntry
	err = f.getJSON(ctx, contentsURL(or, dir, branch), &entries)
	return entries, err
}

// contentsURL builds the GitHub contents API URL for a directory. It uses
// url.URL so directory segments with spaces (such as the PayloadsAllTheThings
// category folders) are percent-encoded.
func contentsURL(ownerRepo, dir, branch string) string {
	u := &url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   "/repos/" + ownerRepo + "/contents/" + dir,
	}
	if branch != "" {
		u.RawQuery = url.Values{"ref": {branch}}.Encode()
	}
	return u.String()
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
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxGitHubAPIBytes+1))
	if err != nil {
		return fmt.Errorf("read github response: %w", err)
	}
	if len(body) > maxGitHubAPIBytes {
		return fmt.Errorf("github response exceeds %d-byte limit", maxGitHubAPIBytes)
	}
	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("decode github response: %w", err)
	}
	return nil
}

func newGitHubHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many download redirects")
			}
			return validateGitHubURL(req.URL.String())
		},
	}
}

func validateGitHubURL(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme != "https" || parsed.User != nil || parsed.Port() != "" {
		return fmt.Errorf("download URL is not a valid HTTPS URL")
	}
	host := strings.ToLower(parsed.Hostname())
	if host == "github.com" || host == "api.github.com" || host == "raw.githubusercontent.com" ||
		strings.HasSuffix(host, ".githubusercontent.com") {
		return nil
	}
	return fmt.Errorf("download URL host %q is not trusted", host)
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
