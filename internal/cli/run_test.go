package cli

import (
	"reflect"
	"testing"
)

func TestSelectBinary(t *testing.T) {
	bins := []string{"secretsdump.py", "getTGT.py", "GetUserSPNs.py"}
	tests := []struct {
		name    string
		rest    []string
		wantBin string
		wantFwd []string
	}{
		{"default primary", nil, "secretsdump.py", []string{}},
		{"select by exact", []string{"getTGT.py"}, "getTGT.py", []string{}},
		{"select loose case and suffix", []string{"gettgt"}, "getTGT.py", []string{}},
		{"select with tool prefix", []string{"impacket-getTGT"}, "getTGT.py", []string{}},
		{"select then forward", []string{"getTGT", "-dc-ip", "1.2.3.4"}, "getTGT.py", []string{"-dc-ip", "1.2.3.4"}},
		{"select with explicit separator", []string{"getTGT", "--", "-key", "v"}, "getTGT.py", []string{"-key", "v"}},
		{"no match forwards all", []string{"ssh", "127.0.0.1"}, "secretsdump.py", []string{"ssh", "127.0.0.1"}},
		{"separator only forwards", []string{"--", "-h"}, "secretsdump.py", []string{"-h"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bin, fwd := selectBinary("impacket", bins, tt.rest)
			if bin != tt.wantBin {
				t.Fatalf("bin = %q, want %q", bin, tt.wantBin)
			}
			if !reflect.DeepEqual(fwd, tt.wantFwd) {
				t.Fatalf("forwarded = %#v, want %#v", fwd, tt.wantFwd)
			}
		})
	}
}

func TestNormalizeBinary(t *testing.T) {
	tests := []struct {
		tool, name, want string
	}{
		{"impacket", "getTGT.py", "gettgt"},
		{"impacket", "impacket-GetUserSPNs", "getuserspns"},
		{"netexec", "nxc", "nxc"},
	}
	for _, tt := range tests {
		if got := normalizeBinary(tt.tool, tt.name); got != tt.want {
			t.Fatalf("normalizeBinary(%q, %q) = %q, want %q", tt.tool, tt.name, got, tt.want)
		}
	}
}
