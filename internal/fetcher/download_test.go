package fetcher

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDownloadAtomicallyWritesExpectedContent(t *testing.T) {
	const content = "verified-content"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(content))
	}))
	defer server.Close()
	destination := filepath.Join(t.TempDir(), "asset")
	fetcher := testFetcher(server.Client(), 1024)

	size, err := fetcher.download(context.Background(), server.URL, destination, int64(len(content)))
	if err != nil {
		t.Fatalf("download: %v", err)
	}
	if size != int64(len(content)) {
		t.Fatalf("size = %d", size)
	}
	data, err := os.ReadFile(destination)
	if err != nil || string(data) != content {
		t.Fatalf("downloaded content = %q, %v", data, err)
	}
	info, err := os.Stat(destination)
	if err != nil {
		t.Fatalf("stat downloaded file: %v", err)
	}
	if info.Mode().Perm() != 0o755 {
		t.Fatalf("downloaded mode = %v", info.Mode().Perm())
	}
	assertNoPartFiles(t, filepath.Dir(destination))
}

func TestDownloadRejectsOversizedMetadataBeforeRequest(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		called = true
	}))
	defer server.Close()
	fetcher := testFetcher(server.Client(), 4)
	if _, err := fetcher.download(context.Background(), server.URL, filepath.Join(t.TempDir(), "asset"), 5); err == nil {
		t.Fatal("expected oversized metadata to be rejected")
	}
	if called {
		t.Fatal("server called for rejected metadata")
	}
}

func TestDownloadRejectsInvalidMetadataSize(t *testing.T) {
	fetcher := testFetcher(http.DefaultClient, 4)
	for _, size := range []int64{-1, 5} {
		if _, err := fetcher.download(context.Background(), "https://unused.invalid", filepath.Join(t.TempDir(), "asset"), size); err == nil {
			t.Fatalf("metadata size %d unexpectedly accepted", size)
		}
	}
}

func TestDownloadLimitPreservesExistingDestination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Transfer-Encoding", "chunked")
		_, _ = w.Write([]byte("oversized"))
	}))
	defer server.Close()
	dir := t.TempDir()
	destination := filepath.Join(dir, "asset")
	if err := os.WriteFile(destination, []byte("previous"), 0o755); err != nil {
		t.Fatal(err)
	}
	fetcher := testFetcher(server.Client(), 4)
	if _, err := fetcher.download(context.Background(), server.URL, destination, 0); err == nil {
		t.Fatal("expected oversized body to be rejected")
	}
	assertPreviousFile(t, destination)
	assertNoPartFiles(t, dir)
}

func TestDownloadSizeMismatchPreservesExistingDestination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("short"))
	}))
	defer server.Close()
	dir := t.TempDir()
	destination := filepath.Join(dir, "asset")
	if err := os.WriteFile(destination, []byte("previous"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := testFetcher(server.Client(), 1024).download(context.Background(), server.URL, destination, 10); err == nil {
		t.Fatal("expected mismatched size to be rejected")
	}
	assertPreviousFile(t, destination)
	assertNoPartFiles(t, dir)
}

func TestDownloadCancellationCleansPartialFile(t *testing.T) {
	started := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Transfer-Encoding", "chunked")
		_, _ = w.Write([]byte("partial"))
		w.(http.Flusher).Flush()
		close(started)
		<-r.Context().Done()
	}))
	defer server.Close()
	dir := t.TempDir()
	destination := filepath.Join(dir, "asset")
	if err := os.WriteFile(destination, []byte("previous"), 0o755); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := testFetcher(server.Client(), 1024).download(ctx, server.URL, destination, 0)
		done <- err
	}()
	<-started
	cancel()
	if err := <-done; !errors.Is(err, context.Canceled) && !strings.Contains(err.Error(), "context canceled") {
		t.Fatalf("download error = %v, want cancellation", err)
	}
	assertPreviousFile(t, destination)
	assertNoPartFiles(t, dir)
}

func TestWriteDownloadTimeoutCleansPartialFile(t *testing.T) {
	dir := t.TempDir()
	destination := filepath.Join(dir, "asset")
	if err := os.WriteFile(destination, []byte("previous"), 0o755); err != nil {
		t.Fatal(err)
	}
	body := io.MultiReader(strings.NewReader("partial"), errorReader{err: context.DeadlineExceeded})
	if _, err := testFetcher(http.DefaultClient, 1024).writeDownload(body, destination, 100); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("writeDownload error = %v, want deadline exceeded", err)
	}
	assertPreviousFile(t, destination)
	assertNoPartFiles(t, dir)
}

func testFetcher(client *http.Client, limit int64) *Fetcher {
	return &Fetcher{client: client, maxDownloadBytes: limit}
}

func assertNoPartFiles(t *testing.T, dir string) {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(dir, ".*.part-*"))
	if err != nil || len(matches) != 0 {
		t.Fatalf("partial files = %v, %v", matches, err)
	}
}

func assertPreviousFile(t *testing.T, path string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "previous" {
		t.Fatalf("file content = %q, %v; want previous", data, err)
	}
}

type errorReader struct {
	err error
}

func (r errorReader) Read([]byte) (int, error) {
	return 0, r.err
}
