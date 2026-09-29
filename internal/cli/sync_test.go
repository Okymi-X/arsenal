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
		{args: []string{"--ref", "v1.0.0"}, want: syncOptions{reference: "v1.0.0"}},
		{args: []string{"--repo=owner/repo", "--ref=main"}, want: syncOptions{repo: "owner/repo", reference: "main"}},
		{args: []string{"--ref"}, wantErr: true},
		{args: []string{"--list-refs", "--ref", "main"}, wantErr: true},
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
		args                    []string
		wantSpec, wantGitHubRef string
		wantErr                 bool
	}{
		{args: []string{"nxc"}, wantSpec: "nxc"},
		{args: []string{"nxc@1.0.0"}, wantSpec: "nxc@1.0.0"},
		{args: []string{"nxc", "--github-ref", "v2.0.0"}, wantSpec: "nxc", wantGitHubRef: "v2.0.0"},
		{args: []string{"nxc@1.0.0", "--github-ref", "main"}, wantErr: true},
		{args: []string{"nxc", "--github-ref"}, wantErr: true},
	}
	for _, test := range tests {
		spec, ref, err := parseInstallArgs(test.args)
		if test.wantErr {
			if err == nil {
				t.Fatalf("parseInstallArgs(%v) unexpectedly succeeded", test.args)
			}
			continue
		}
		if err != nil || spec != test.wantSpec || ref != test.wantGitHubRef {
			t.Fatalf("parseInstallArgs(%v) = %q, %q, %v", test.args, spec, ref, err)
		}
	}
}
