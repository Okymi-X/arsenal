package store

import (
	"path/filepath"
	"testing"
)

func TestManifestUpsertAndDelete(t *testing.T) {
	m := &Manifest{}
	m.Upsert(InstalledTool{Name: "nxc", Version: "1.4.0"})
	m.Upsert(InstalledTool{Name: "nxc", Version: "1.4.0", Path: "/p"})
	if len(m.Tools) != 1 {
		t.Fatalf("upsert should not duplicate: %d", len(m.Tools))
	}
	if m.Tools[0].Path != "/p" {
		t.Fatalf("upsert should replace: %+v", m.Tools[0])
	}
	if !m.Delete("nxc", "1.4.0") {
		t.Fatal("delete should report success")
	}
	if len(m.Tools) != 0 {
		t.Fatalf("delete should remove entry: %d", len(m.Tools))
	}
}

func TestManifestSetActive(t *testing.T) {
	m := &Manifest{Tools: []InstalledTool{
		{Name: "nxc", Version: "1.4.0", Active: true},
		{Name: "nxc", Version: "1.3.0"},
	}}
	if !m.SetActive("nxc", "1.3.0") {
		t.Fatal("SetActive should find the target")
	}
	a, ok := m.Active("nxc")
	if !ok || a.Version != "1.3.0" {
		t.Fatalf("active = %+v ok=%v", a, ok)
	}
	if m.SetActive("nxc", "9.9.9") {
		t.Fatal("SetActive should fail for missing version")
	}
}

func TestFileStoreRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "manifest.json")
	s := NewFileStore(path)
	loaded, err := s.Load()
	if err != nil || len(loaded.Tools) != 0 {
		t.Fatalf("empty load: %v %+v", err, loaded)
	}
	loaded.Upsert(InstalledTool{Name: "impacket", Version: "0.12.0"})
	if err := s.Save(loaded); err != nil {
		t.Fatalf("save: %v", err)
	}
	again, err := s.Load()
	if err != nil || len(again.Tools) != 1 {
		t.Fatalf("reload: %v %+v", err, again)
	}
}
