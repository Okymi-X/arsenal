package op

import (
	"path/filepath"
	"testing"
)

func TestLockfileRoundTrip(t *testing.T) {
	lf := &Lockfile{
		Op:              "redteam-q3",
		Generated:       "2026-06-14T00:00:00Z",
		RegistryVersion: "1",
		Entries: []LockEntry{
			{Tool: "netexec", Version: "1.4.0", PipSpec: "netexec==1.4.0", InstallMethod: "pip"},
			{Tool: "impacket", Version: "0.12.0", PipSpec: "impacket==0.12.0", InstallMethod: "pip"},
		},
	}
	path := filepath.Join(t.TempDir(), "op.lock.toml")
	if err := WriteLockfile(path, lf); err != nil {
		t.Fatalf("WriteLockfile: %v", err)
	}
	got, err := ReadLockfile(path)
	if err != nil {
		t.Fatalf("ReadLockfile: %v", err)
	}
	if got.Op != lf.Op || len(got.Entries) != 2 {
		t.Fatalf("round trip mismatch: %+v", got)
	}
	if got.Entries[0].Tool != "netexec" || got.Entries[0].Version != "1.4.0" {
		t.Fatalf("unexpected first entry: %+v", got.Entries[0])
	}
}

func TestValidateLockfile(t *testing.T) {
	tests := []struct {
		name    string
		lf      *Lockfile
		wantErr bool
	}{
		{"valid", &Lockfile{Op: "x", Entries: []LockEntry{{Tool: "a", Version: "1"}}}, false},
		{"no op", &Lockfile{Entries: []LockEntry{{Tool: "a", Version: "1"}}}, true},
		{"entry no tool", &Lockfile{Op: "x", Entries: []LockEntry{{Version: "1"}}}, true},
		{"entry no version", &Lockfile{Op: "x", Entries: []LockEntry{{Tool: "a"}}}, true},
		{"duplicate tool", &Lockfile{Op: "x", Entries: []LockEntry{{Tool: "a", Version: "1"}, {Tool: "a", Version: "2"}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLockfile(tt.lf)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateLockfile(t *testing.T) {
	o := &Op{Name: "eng", Pins: []Pin{{Tool: "netexec", Version: "1.4.0"}}}
	resolve := func(tool, version string) (LockEntry, error) {
		return LockEntry{Tool: tool, Version: version, InstallMethod: "pip"}, nil
	}
	lf, err := GenerateLockfile(o, "1", "now", resolve)
	if err != nil {
		t.Fatalf("GenerateLockfile: %v", err)
	}
	if len(lf.Entries) != 1 || lf.Entries[0].Tool != "netexec" {
		t.Fatalf("unexpected lockfile: %+v", lf)
	}
}
