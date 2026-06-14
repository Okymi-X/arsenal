package cli

import "testing"

func TestParseFetchArgs(t *testing.T) {
	tests := []struct {
		name                string
		args                []string
		wantName, wantBin   string
		wantBuild, wantDest string
		wantList            bool
		wantErr             bool
	}{
		{name: "name only", args: []string{"pspy"}, wantName: "pspy", wantDest: "."},
		{name: "name and binary", args: []string{"winpeas", "winPEASany.exe"}, wantName: "winpeas", wantBin: "winPEASany.exe", wantDest: "."},
		{
			name: "build and dest", args: []string{"sharpcollection", "Rubeus", "--build", "NetFramework_4.7_x64", "--dest", "/tmp/x"},
			wantName: "sharpcollection", wantBin: "Rubeus", wantBuild: "NetFramework_4.7_x64", wantDest: "/tmp/x",
		},
		{name: "inline dest", args: []string{"pspy", "--dest=/loot"}, wantName: "pspy", wantDest: "/loot"},
		{name: "list", args: []string{"sharp", "--list"}, wantName: "sharp", wantDest: ".", wantList: true},
		{name: "no positional", args: []string{"--list"}, wantErr: true},
		{name: "too many positionals", args: []string{"a", "b", "c"}, wantErr: true},
		{name: "unknown flag", args: []string{"pspy", "--bogus"}, wantErr: true},
		{name: "dest missing value", args: []string{"pspy", "--dest"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, err := parseFetchArgs(tt.args)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %+v", o)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if o.name != tt.wantName || o.binary != tt.wantBin || o.build != tt.wantBuild || o.dest != tt.wantDest || o.list != tt.wantList {
				t.Fatalf("parseFetchArgs(%v) = %+v", tt.args, o)
			}
		})
	}
}
