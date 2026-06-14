// Command regbuild assembles the canonical registry.toml from the per-category
// segment files in registry/segments.
//
// The registry is authored as small segment files (one per category) so it
// stays maintainable; this tool concatenates them, under the version/updated
// metadata, into the single registry.toml that is embedded, synced, and
// verified. Run it via `make registry`. With -check it reports drift without
// writing, for CI.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/BurntSushi/toml"
	"github.com/Okymi-X/arsenal/internal/registry"
)

// segmentOrder is the preferred ordering of segments in the assembled file.
// Any segment not listed here is appended in alphabetical order.
var segmentOrder = []string{"ad", "web", "recon", "password", "misc", "assets"}

const header = `# arsenal curated registry - GENERATED FILE, DO NOT EDIT.
#
# Assembled from registry/segments/*.toml by tools/regbuild (run ` + "`make registry`" + `).
# Edit the segment files, not this file. Every version and asset here is
# verified against its official source by the registry-check CI workflow.
`

func main() {
	dir := flag.String("dir", "registry/segments", "directory of segment files")
	out := flag.String("out", "registry/registry.toml", "assembled registry path")
	verify := flag.Bool("verify", false, "report drift without writing (for CI)")
	flag.Parse()

	assembled, err := build(*dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[fail] %v\n", err)
		os.Exit(1)
	}

	if *verify {
		current, err := os.ReadFile(*out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[fail] read %s: %v\n", *out, err)
			os.Exit(1)
		}
		if !bytes.Equal(current, assembled) {
			fmt.Fprintf(os.Stderr, "[fail] %s is out of date; run `make registry`\n", *out)
			os.Exit(1)
		}
		fmt.Printf("[ok] %s is up to date with %s\n", *out, *dir)
		return
	}

	if err := os.WriteFile(*out, assembled, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "[fail] write %s: %v\n", *out, err)
		os.Exit(1)
	}
	fmt.Printf("[ok] wrote %s from %s\n", *out, *dir)
}

// build assembles and validates the registry bytes from a segments directory.
func build(dir string) ([]byte, error) {
	meta, err := os.ReadFile(filepath.Join(dir, "_meta.toml"))
	if err != nil {
		return nil, fmt.Errorf("read metadata: %w", err)
	}
	var m struct {
		Version string `toml:"version"`
		Updated string `toml:"updated"`
	}
	if err := toml.Unmarshal(meta, &m); err != nil {
		return nil, fmt.Errorf("parse metadata: %w", err)
	}

	var buf bytes.Buffer
	buf.WriteString(header)
	buf.WriteString("\n")
	fmt.Fprintf(&buf, "version = %q\n", m.Version)
	fmt.Fprintf(&buf, "updated = %q\n", m.Updated)

	for _, name := range orderedSegments(dir) {
		seg, err := os.ReadFile(filepath.Join(dir, name+".toml"))
		if err != nil {
			return nil, fmt.Errorf("read segment %s: %w", name, err)
		}
		buf.WriteString("\n")
		buf.Write(bytes.TrimRight(seg, "\n"))
		buf.WriteString("\n")
	}

	assembled := buf.Bytes()
	if _, err := registry.Parse(assembled); err != nil {
		return nil, fmt.Errorf("assembled registry is invalid: %w", err)
	}
	return assembled, nil
}

// orderedSegments lists segment names (without extension) in preferred order,
// appending any unknown segments alphabetically. _meta is excluded.
func orderedSegments(dir string) []string {
	entries, _ := filepath.Glob(filepath.Join(dir, "*.toml"))
	present := make(map[string]bool)
	for _, e := range entries {
		base := filepath.Base(e)
		name := base[:len(base)-len(".toml")]
		if name != "_meta" {
			present[name] = true
		}
	}
	var ordered []string
	for _, name := range segmentOrder {
		if present[name] {
			ordered = append(ordered, name)
			delete(present, name)
		}
	}
	rest := make([]string, 0, len(present))
	for name := range present {
		rest = append(rest, name)
	}
	sort.Strings(rest)
	return append(ordered, rest...)
}
