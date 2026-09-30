package cli

import "testing"

func TestParseSyncArgs(t *testing.T) {
	tests := []struct {
		args    []string
		want    syncOptions
		wantErr bool
	}{
		{args: nil, want: syncOptions{}},
		{args: []string{"--list-refs"}, want: syncOptions{listRefs: true}},
		{args: []string{"--select-ref"}, want: syncOptions{selectRef: true}},
		{args: []string{"--ref", "v1.0.0"}, want: syncOptions{reference: "v1.0.0"}},
		{args: []string{"--repo=owner/repo", "--ref=main"}, want: syncOptions{repo: "owner/repo", reference: "main"}},
		{args: []string{"--ref"}, wantErr: true},
		{args: []string{"--list-refs", "--ref", "main"}, wantErr: true},
		{args: []string{"--select-ref", "--ref", "main"}, wantErr: true},
		{args: []string{"--select-ref", "--list-refs"}, wantErr: true},
		{args: []string{"unexpected"}, wantErr: true},
	}
	for _, test := range tests {
		got, err := parseSyncArgs(test.args)
		if test.wantErr {
			if err == nil {
				t.Fatalf("parseSyncArgs(%v) unexpectedly succeeded", test.args)
			}
			continue
		}
		if err != nil || got != test.want {
			t.Fatalf("parseSyncArgs(%v) = %#v, %v; want %#v", test.args, got, err, test.want)
		}
	}
}

func TestParseInstallArgs(t *testing.T) {
	tests := []struct {
		args    []string
		want    installOptions
		wantErr bool
	}{
		{args: []string{"nxc"}, want: installOptions{spec: "nxc"}},
		{args: []string{"nxc@1.0.0"}, want: installOptions{spec: "nxc@1.0.0"}},
		{args: []string{"--select"}, want: installOptions{selectTool: true, selectVersion: true}},
		{args: []string{"--select", "ad"}, want: installOptions{query: "ad", selectTool: true, selectVersion: true}},
		{args: []string{"nxc", "--select"}, want: installOptions{spec: "nxc", selectVersion: true}},
		{args: []string{"nxc", "--github-select"}, want: installOptions{spec: "nxc", selectGitHub: true}},
		{args: []string{"nxc", "--github-ref", "v2.0.0"}, want: installOptions{spec: "nxc", githubRef: "v2.0.0"}},
		{args: []string{"nxc@1.0.0", "--github-ref", "main"}, wantErr: true},
		{args: []string{"nxc@1.0.0", "--select"}, wantErr: true},
		{args: []string{"nxc@1.0.0", "--github-select"}, wantErr: true},
		{args: []string{"nxc", "--github-ref"}, wantErr: true},
	}
	for _, test := range tests {
		got, err := parseInstallArgs(test.args)
		if test.wantErr {
			if err == nil {
				t.Fatalf("parseInstallArgs(%v) unexpectedly succeeded", test.args)
			}
			continue
		}
		if err != nil || got != test.want {
			t.Fatalf("parseInstallArgs(%v) = %#v, %v; want %#v", test.args, got, err, test.want)
		}
	}
}
