package fetcher

import "testing"

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

func TestOwnerRepo(t *testing.T) {
	tests := []struct {
		in, want string
		wantErr  bool
	}{
		{"https://github.com/lc/gau", "lc/gau", false},
		{"https://github.com/lc/gau.git", "lc/gau", false},
		{"https://github.com/lc/gau/", "lc/gau", false},
		{"https://gitlab.com/x/y", "", true},
		{"not a url", "", true},
	}
	for _, tt := range tests {
		got, err := ownerRepo(tt.in)
		if tt.wantErr {
			if err == nil {
				t.Errorf("ownerRepo(%q) = %q, want error", tt.in, got)
			}
			continue
		}
		if err != nil || got != tt.want {
			t.Errorf("ownerRepo(%q) = (%q, %v), want %q", tt.in, got, err, tt.want)
		}
	}
}
