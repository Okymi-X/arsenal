# Project memory

This file contains durable, verified facts that help a new maintainer or coding
agent resume work. It is not a transcript, scratchpad, task tracker, or place
for assumptions. Update it only when a fact is confirmed in the repository.

## Stable facts

- The module path is `github.com/Okymi-X/arsenal` and the implementation is Go.
- The shipped artifact is a statically linked CLI binary named `arsenal`.
- The canonical project icon is `assets/branding/arsenal-icon.png`.
- `registry/registry.toml` is a small manifest. The embedded catalog is
  assembled and validated directly from its ordered `registry/segments/*.toml`
  files at build time.
- Python virtual environments are the implemented default isolation backend.
- Pip, git-plus-pip, Go, and Cargo are implemented install methods. The
  prebuilt-binary method remains an explicit stub.
- Standalone fetched assets intentionally bypass installation, isolation,
  shims, and the installed-tool manifest.
- Engagement profiles are called ops and use TOML lockfiles.
- Local state is rooted at `ARSENAL_HOME`, then `XDG_DATA_HOME/arsenal`, then
  the platform default documented in `ARCHITECTURE.md`.
- The CLI has no telemetry and normal output is quiet by default.
- Segmented registries are SHA-256 pinned. GitHub registry refs and explicit
  Python, Go, and Cargo tool refs are resolved to immutable commits before
  pulling.
- Native Go and Cargo installations use per-tool staging roots and caches,
  verify declared executables, and promote completed installs atomically.
- Safe upgrades target the newest tested registry entry, retain the previous
  installed environment, and never auto-change ahead or untracked versions.
- Production Go source currently follows a soft 180-line file limit and the
  linter rejects cyclomatic complexity of 15 or greater.

## Authoritative sources

- Product scope and acceptance: `PRD.md`
- Mandatory engineering and security policy: `RULES.md`
- Package boundaries and runtime flow: `ARCHITECTURE.md`
- CLI and implementation decisions: `DESIGN.md`
- Incomplete repository-level work: `TASKS.md`
- User commands: `usage.md`
- Registry schema and authoring: `registry-format.md`
- Release history: `../CHANGELOG.md`

When memory disagrees with code or tests, inspect the implementation and update
this file. Never store credentials, tokens, private infrastructure, customer or
target data, unpublished vulnerabilities, or temporary debugging notes here.
