// Package builtin embeds the curated registry that ships inside the binary so
// arsenal works out of the box and offline, before any `arsenal sync`.
//
// registry.toml in this directory is the canonical, in-repo source of the
// curated catalog; it is embedded at build time and written to the user's
// data directory on first run.
package builtin

import _ "embed"

//go:embed registry.toml
var data []byte

// Bytes returns the embedded registry TOML.
func Bytes() []byte { return data }
