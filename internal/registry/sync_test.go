package registry

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestFileSourceSyncHonorsCancellation(t *testing.T) {
	started := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		close(started)
		<-request.Context().Done()
	}))
	defer server.Close()
	path := filepath.Join(t.TempDir(), "registry.toml")
	previous := []byte("previous registry")
	if err := os.WriteFile(path, previous, 0o644); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- NewFileSource(path, server.URL).Sync(ctx) }()
	<-started
	cancel()
	if err := <-done; !errors.Is(err, context.Canceled) {
		t.Fatalf("Sync error = %v, want context cancellation", err)
	}
	got, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(got, previous) {
		t.Fatalf("previous registry changed: %q, %v", got, err)
	}
}

func TestFileSourceSyncsSegmentedRegistry(t *testing.T) {
	segment := `[[tool]]
name = "example"
category = "recon"
install_method = "pip"
binary = "example"

  [[tool.version]]
  tag = "1.0.0"
  tested = true
  pip_spec = "example==1.0.0"
`
	manifest := testToolsManifest([]byte(segment))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/registry.toml":
			_, _ = w.Write(manifest)
		case "/segments/tools.toml":
			_, _ = w.Write([]byte(segment))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	path := filepath.Join(t.TempDir(), "registry.toml")
	source := NewFileSource(path, server.URL+"/registry.toml")
	if err := source.Sync(context.Background()); err != nil {
		t.Fatalf("sync: %v", err)
	}
	reg, err := source.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(reg.Tools) != 1 || reg.Tools[0].Name != "example" {
		t.Fatalf("unexpected tools: %#v", reg.Tools)
	}
}

func TestFileSourceRejectsOversizedResponseWithoutReplacingRegistry(t *testing.T) {
	previous := []byte(`version = "1"
[[tool]]
name = "existing"
install_method = "pip"
  [[tool.version]]
  tag = "1.0.0"
`)
	path := filepath.Join(t.TempDir(), "registry.toml")
	if err := os.WriteFile(path, previous, 0o644); err != nil {
		t.Fatalf("seed registry: %v", err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(strings.Repeat("x", maxRemoteRegistryFile+1)))
	}))
	defer server.Close()

	source := NewFileSource(path, server.URL+"/registry.toml")
	if err := source.Sync(context.Background()); err == nil {
		t.Fatal("expected oversized response to be rejected")
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read registry: %v", err)
	}
	if !bytes.Equal(got, previous) {
		t.Fatal("failed sync replaced the previous registry")
	}
}

func TestFileSourceRejectsChecksumMismatchWithoutReplacingRegistry(t *testing.T) {
	previous := []byte(`version = "1"
[[tool]]
name = "existing"
install_method = "pip"
  [[tool.version]]
  tag = "1.0.0"
`)
	path := filepath.Join(t.TempDir(), "registry.toml")
	if err := os.WriteFile(path, previous, 0o644); err != nil {
		t.Fatalf("seed registry: %v", err)
	}
	manifest := testToolsManifest([]byte("expected"))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/registry.toml" {
			_, _ = w.Write(manifest)
			return
		}
		_, _ = w.Write([]byte("tampered"))
	}))
	defer server.Close()

	source := NewFileSource(path, server.URL+"/registry.toml")
	if err := source.Sync(context.Background()); err == nil {
		t.Fatal("expected checksum mismatch to be rejected")
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read registry: %v", err)
	}
	if !bytes.Equal(got, previous) {
		t.Fatal("failed sync replaced the previous registry")
	}
}

func TestFileSourceRejectsEmptyManifest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("version = \"1\"\nsegments = []\n"))
	}))
	defer server.Close()

	source := NewFileSource(filepath.Join(t.TempDir(), "registry.toml"), server.URL)
	if err := source.Sync(context.Background()); err == nil {
		t.Fatal("expected empty manifest to be rejected")
	}
}

func TestFileSourceBoundsConcurrentSegmentPulls(t *testing.T) {
	segments := make(map[string][]byte)
	names := make([]string, 0, 8)
	for i := range 8 {
		name := fmt.Sprintf("segments/tool-%d.toml", i)
		names = append(names, name)
		segments[name] = fmt.Appendf(nil, "[[tool]]\nname = %q\ninstall_method = \"pip\"\nbinary = %q\n[[tool.version]]\ntag = \"1.0.0\"\ntested = true\n", fmt.Sprintf("tool-%d", i), fmt.Sprintf("tool-%d", i))
	}
	manifest := multiSegmentManifest(names, segments)
	var active atomic.Int32
	var maximum atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/registry.toml" {
			_, _ = w.Write(manifest)
			return
		}
		name := strings.TrimPrefix(r.URL.Path, "/")
		segment, ok := segments[name]
		if !ok {
			http.NotFound(w, r)
			return
		}
		current := active.Add(1)
		defer active.Add(-1)
		for {
			previous := maximum.Load()
			if current <= previous || maximum.CompareAndSwap(previous, current) {
				break
			}
		}
		time.Sleep(25 * time.Millisecond)
		_, _ = w.Write(segment)
	}))
	defer server.Close()

	source := NewFileSource(filepath.Join(t.TempDir(), "registry.toml"), server.URL+"/registry.toml")
	if err := source.Sync(context.Background()); err != nil {
		t.Fatalf("sync: %v", err)
	}
	if got := maximum.Load(); got < 2 || got > pullWorkers {
		t.Fatalf("maximum concurrent pulls = %d, want 2..%d", got, pullWorkers)
	}
}

func TestFileSourceRejectsCrossOriginRedirect(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("version = \"1\"\n"))
	}))
	defer target.Close()
	sourceServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL+"/registry.toml", http.StatusFound)
	}))
	defer sourceServer.Close()

	source := NewFileSource(filepath.Join(t.TempDir(), "registry.toml"), sourceServer.URL+"/registry.toml")
	if err := source.Sync(context.Background()); err == nil || !strings.Contains(err.Error(), "redirect changed remote origin") {
		t.Fatalf("Sync() error = %v", err)
	}
}

func multiSegmentManifest(names []string, segments map[string][]byte) []byte {
	var builder strings.Builder
	builder.WriteString("version = \"1\"\nsegments = [\n")
	for _, name := range names {
		fmt.Fprintf(&builder, "%q,\n", name)
	}
	builder.WriteString("]\n[segment_sha256]\n")
	for _, name := range names {
		digest := sha256.Sum256(segments[name])
		fmt.Fprintf(&builder, "%q = %q\n", name, fmt.Sprintf("%x", digest))
	}
	return []byte(builder.String())
}
