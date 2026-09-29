package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Okymi-X/arsenal/internal/registry"
)

func readManifest(path string) (registry.Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return registry.Manifest{}, fmt.Errorf("read manifest: %w", err)
	}
	return registry.ParseManifest(data)
}

func writeManifest(path string, manifest registry.Manifest) error {
	manifest.SegmentSHA256 = make(map[string]string, len(manifest.Segments))
	for _, name := range manifest.Segments {
		data, err := os.ReadFile(filepath.Join(filepath.Dir(path), filepath.FromSlash(name)))
		if err != nil {
			return fmt.Errorf("read segment %s: %w", name, err)
		}
		digest := sha256.Sum256(data)
		manifest.SegmentSHA256[name] = hex.EncodeToString(digest[:])
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), ".registry-manifest-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer func() { _ = os.Remove(tmpPath) }()
	w := bufio.NewWriter(tmp)
	_, err = fmt.Fprintf(w, "# Registry metadata and generated segment integrity pins.\n# Edit entries in registry/segments/*.toml, then run make registry.\n\nversion = %q\nupdated = %q\nsegments = [\n", manifest.Version, manifest.Updated)
	for _, name := range manifest.Segments {
		if err == nil {
			_, err = fmt.Fprintf(w, "  %q,\n", name)
		}
	}
	if err == nil {
		_, err = fmt.Fprintln(w, "]\n\n[segment_sha256]")
	}
	for _, name := range manifest.Segments {
		if err == nil {
			_, err = fmt.Fprintf(w, "%q = %q\n", name, manifest.SegmentSHA256[name])
		}
	}
	if err != nil {
		_ = tmp.Close()
		return err
	}
	if err := w.Flush(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Chmod(0o644); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}

func checkCoverage(manifestPath, dir string, segments []string) error {
	manifestDir := filepath.Dir(manifestPath)
	segmentDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("resolve segment directory: %w", err)
	}
	referenced := make(map[string]struct{}, len(segments))
	for _, name := range segments {
		full, err := filepath.Abs(filepath.Join(manifestDir, filepath.FromSlash(name)))
		if err != nil {
			return fmt.Errorf("resolve segment %s: %w", name, err)
		}
		rel, err := filepath.Rel(segmentDir, full)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return fmt.Errorf("segment %q is outside %s", name, dir)
		}
		referenced[full] = struct{}{}
	}

	files, err := filepath.Glob(filepath.Join(segmentDir, "*.toml"))
	if err != nil {
		return fmt.Errorf("list segments: %w", err)
	}
	for _, file := range files {
		if _, ok := referenced[file]; !ok {
			return fmt.Errorf("segment %s is not listed in the manifest", file)
		}
	}
	return nil
}
