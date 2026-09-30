# Project memory

This file contains durable, verified facts that help a new maintainer or coding
agent resume work. It is not a transcript, scratchpad, task tracker, or place
for assumptions. Update it only when a fact is confirmed in the repository.

## Stable facts

- The module path is `github.com/Okymi-X/arsenal` and the implementation is Go.
- The shipped artifact is a statically linked CLI binary named `arsenal`.
- Users can install a tagged release with
  `go install github.com/Okymi-X/arsenal/cmd/arsenal@latest`; the CLI reads the
  module build metadata when no link-time version was injected.
- The canonical project icon is `assets/branding/arsenal-icon.png`.
- `registry/registry.toml` is a small manifest. The embedded catalog is
  assembled and validated directly from its ordered `registry/segments/*.toml`
  files at build time.
- Python virtual environments are the implemented default isolation backend.
- `arsenal run` resolves managed binaries inside their installation root. The
  private shim directory stays off the host `PATH` by default so existing
  system, pipx, and user-installed commands are not shadowed.
- Pip, git-plus-pip, Go, and Cargo are implemented install methods. The
  prebuilt-binary method remains an explicit stub.
- Standalone fetched assets intentionally bypass installation, isolation,
  shims, and the installed-tool manifest.
- Engagement profiles are called ops and use TOML lockfiles.
- Lockfile application validates all entries against the active registry before
  installation and rejects registry-version, method, commit, or package-spec
  drift.
- Local state is rooted at `ARSENAL_HOME`, then `XDG_DATA_HOME/arsenal`, then
  the platform default documented in `ARCHITECTURE.md`.
- The CLI has no telemetry and normal output is quiet by default.
- Segmented registries are SHA-256 pinned. GitHub registry refs and explicit
  Python, Go, and Cargo tool refs are resolved to immutable commits before
  pulling.
- Interactive selection is explicit: tool and curated-version discovery use
  `install --select`, upstream tool tags use `--github-select`, installed
  versions use `switch --select`, and registry tags use `sync --select-ref`.
  Remote choices are numeric, terminal-sanitized, and rejected if their tag
  moves between display and final resolution.
- Native Go and Cargo installations use per-tool staging roots and caches,
  verify declared executables, and promote completed installs atomically.
- Asset downloads accept at most 512 MiB from HTTPS GitHub-controlled hosts,
  verify declared and actual sizes, and atomically promote temporary files.
- The CLI propagates interrupt and termination cancellation through installers,
  GitHub operations, registry synchronization, and asset downloads.
- Safe upgrades target the newest tested registry entry, retain the previous
  installed environment, and never auto-change ahead or untracked versions.
- Registry versions may declare exact `pip_dependencies`. Python installers
  resolve those pins with the primary package in one transaction, lockfiles
  preserve them, upstream verification confirms every pin exists on PyPI, and
  every completed Python installation must pass `pip check`.
- Python virtualenv operations reject unsafe or symlinked installation paths,
  remove inherited `PYTHONHOME` and `PYTHONPATH`, and disable the user site.
- Production Go source currently follows a soft 180-line file limit and the
  linter rejects cyclomatic complexity of 15 or greater.
- Persisted configuration, registry records, installed manifests, ops, and
  lockfiles reject unknown schema fields. Shared path-component validation and
  atomic file replacement live in `internal/safepath` and `internal/fsutil`.
- CI verifies module checksums and scans reachable Go code using a pinned
  `govulncheck` version. Release tags must be documented and point to a commit
  on `main` before artifacts are published.

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
