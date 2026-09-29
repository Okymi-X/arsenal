// Command regbuild validates the segmented registry manifest and its files.
//
// The historical command name is retained for compatibility with existing
// contributor workflows. The registry is no longer assembled into a tracked
// monolithic file; runtime consumers expand the manifest directly.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Okymi-X/arsenal/internal/registry"
)

func main() {
	dir := flag.String("dir", "registry/segments", "directory of segment files")
	manifestPath := flag.String("out", "registry/registry.toml", "registry manifest path")
	write := flag.Bool("write", false, "regenerate segment checksums in the manifest")
	_ = flag.Bool("verify", false, "validate without modifying files")
	flag.Parse()

	manifest, err := readManifest(*manifestPath)
	if err != nil {
		fail("%v", err)
	}
	if len(manifest.Segments) == 0 {
		fail("manifest has no segments")
	}
	if err := checkCoverage(*manifestPath, *dir, manifest.Segments); err != nil {
		fail("%v", err)
	}
	if *write {
		if err := writeManifest(*manifestPath, manifest); err != nil {
			fail("write manifest: %v", err)
		}
	}
	reg, err := registry.Load(*manifestPath)
	if err != nil {
		fail("load segmented registry: %v", err)
	}
	fmt.Printf("[ok] registry manifest references %d segments, %d tools, %d assets\n",
		len(manifest.Segments), len(reg.Tools), len(reg.Assets))
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "[fail] "+format+"\n", args...)
	os.Exit(1)
}
