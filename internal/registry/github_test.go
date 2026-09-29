package registry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGitHubClientListsAndResolvesRefs(t *testing.T) {
	const commit = "0123456789abcdef0123456789abcdef01234567"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != "arsenal" {
			t.Errorf("User-Agent = %q", r.Header.Get("User-Agent"))
		}
		switch r.URL.Path {
		case "/repos/owner/repo/tags":
			_, _ = w.Write([]byte(`[{"name":"v1.0.0","commit":{"sha":"` + commit + `"}}]`))
		case "/repos/owner/repo/commits/v1.0.0":
			_, _ = w.Write([]byte(`{"sha":"` + commit + `"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client, err := NewGitHubClient("owner/repo", "registry/registry.toml", "")
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	client.apiBase = server.URL
	refs, err := client.ListTags(context.Background(), 10)
	if err != nil || len(refs) != 1 || refs[0].Name != "v1.0.0" || refs[0].Commit != commit {
		t.Fatalf("ListTags() = %#v, %v", refs, err)
	}
	gotCommit, err := client.ResolveCommit(context.Background(), "v1.0.0")
	if err != nil || gotCommit != commit {
		t.Fatalf("ResolveCommit() = %q, %v", gotCommit, err)
	}
	gotURL, err := client.RegistryURL(commit)
	if err != nil {
		t.Fatalf("RegistryURL: %v", err)
	}
	wantURL := "https://raw.githubusercontent.com/owner/repo/" + commit + "/registry/registry.toml"
	if gotURL != wantURL {
		t.Fatalf("RegistryURL = %q, want %q", gotURL, wantURL)
	}
}

func TestNewGitHubClientRejectsUnsafeInput(t *testing.T) {
	for _, test := range []struct{ repo, file string }{
		{repo: "owner", file: "registry/registry.toml"},
		{repo: "owner/repo/extra", file: "registry/registry.toml"},
		{repo: "owner/repo", file: "../registry.toml"},
		{repo: "owner/repo", file: "/registry.toml"},
	} {
		if _, err := NewGitHubClient(test.repo, test.file, ""); err == nil {
			t.Fatalf("expected repo=%q file=%q to fail", test.repo, test.file)
		}
	}
}

func TestParseGitHubRepo(t *testing.T) {
	for _, raw := range []string{"https://github.com/owner/repo", "https://github.com/owner/repo.git"} {
		got, err := ParseGitHubRepo(raw)
		if err != nil || got != "owner/repo" {
			t.Fatalf("ParseGitHubRepo(%q) = %q, %v", raw, got, err)
		}
	}
	for _, raw := range []string{"http://github.com/owner/repo", "https://example.com/owner/repo", "git@github.com:owner/repo.git"} {
		if _, err := ParseGitHubRepo(raw); err == nil {
			t.Fatalf("ParseGitHubRepo(%q) unexpectedly succeeded", raw)
		}
	}
}
