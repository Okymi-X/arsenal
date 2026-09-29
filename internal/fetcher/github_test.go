package fetcher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestContentsURL(t *testing.T) {
	tests := []struct {
		name             string
		ownerRepo, dir   string
		branch, expected string
	}{
		{
			"spaces are percent-encoded",
			"swisskyrepo/PayloadsAllTheThings", "SQL Injection/Intruder", "master",
			"https://api.github.com/repos/swisskyrepo/PayloadsAllTheThings/contents/SQL%20Injection/Intruder?ref=master",
		},
		{
			"plain path with branch",
			"fortra/nanodump", "dist", "main",
			"https://api.github.com/repos/fortra/nanodump/contents/dist?ref=main",
		},
		{
			"empty dir and branch",
			"o/r", "", "",
			"https://api.github.com/repos/o/r/contents/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contentsURL(tt.ownerRepo, tt.dir, tt.branch); got != tt.expected {
				t.Fatalf("contentsURL = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestValidateGitHubURL(t *testing.T) {
	tests := []struct {
		url     string
		wantErr bool
	}{
		{"https://github.com/owner/repo/releases/file", false},
		{"https://raw.githubusercontent.com/owner/repo/ref/file", false},
		{"https://release-assets.githubusercontent.com/file", false},
		{"http://github.com/owner/repo", true},
		{"https://github.com:8443/owner/repo", true},
		{"https://github.com.evil.test/file", true},
		{"https://user@github.com/file", true},
	}
	for _, tt := range tests {
		if err := validateGitHubURL(tt.url); (err != nil) != tt.wantErr {
			t.Errorf("validateGitHubURL(%q) error = %v", tt.url, err)
		}
	}
}

func TestGetJSONRejectsOversizedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(strings.Repeat("x", maxGitHubAPIBytes+1)))
	}))
	defer server.Close()
	fetcher := &Fetcher{client: server.Client()}
	var result any
	if err := fetcher.getJSON(context.Background(), server.URL, &result); err == nil {
		t.Fatal("expected oversized GitHub response to be rejected")
	}
}
