// Package builtin embeds the segmented curated registry so arsenal works
// offline before any `arsenal sync`.
package builtin

import (
	"embed"

	internalregistry "github.com/Okymi-X/arsenal/internal/registry"
)

//go:embed registry.toml segments/*.toml
var files embed.FS

var data = mustAssemble()

// Bytes returns the embedded registry TOML.
func Bytes() []byte { return data }

func mustAssemble() []byte {
	manifest, err := files.ReadFile("registry.toml")
	if err != nil {
		panic(err)
	}
	data, err := internalregistry.Expand(manifest, files.ReadFile)
	if err != nil {
		panic(err)
	}
	if _, err := internalregistry.Parse(data); err != nil {
		panic(err)
	}
	return data
}
