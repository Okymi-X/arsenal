# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Added opt-in numbered selection for registry tool discovery, curated install
  versions, upstream GitHub tool tags, installed-version switching, and
  registry release tags. Existing commands remain non-interactive unless a
  selection flag is passed.

### Security

- Interactive GitHub choices are terminal-sanitized, resolved to immutable
  commits, and rejected if a selected tag moves after it is displayed.

## [0.5.4] - 2026-09-30

### Changed

- `arsenal run` is now the documented conflict-free execution path. Private
  shims remain outside the host `PATH`, and `doctor` reports global shim
  exposure as a conflict instead of requiring it.

### Security

- Python environment creation and execution now reject unsafe or symlinked
  tool paths, remove inherited `PYTHONHOME` and `PYTHONPATH`, and disable the
  user site so host Python packages cannot leak into managed environments.

## [0.5.3] - 2026-09-30

### Fixed

- Added version-specific exact Python dependency pins and propagated them
  through installation, lockfiles, schema validation, and upstream checks.
  NetExec 1.5.1 now installs dploot 3.1.3, matching NetExec's upstream lockfile,
  instead of accepting incompatible dploot 4.x releases that break SMB and
  several modules at import time. NetExec 1.4.0 similarly pins dploot 3.1.2.
  Every Python-backed tool now runs `pip check` before installation succeeds.

## [0.5.2] - 2026-09-30

### Changed

- Upgraded the TOML parser to 1.6.0 and made registry, configuration, manifest,
  op, and lockfile decoding reject unsupported fields instead of silently
  ignoring misspelled metadata.
- Added a pinned `govulncheck` target and CI job, dependency verification in CI
  and releases, and release-tag checks for semantic versioning, changelog
  coverage, and ancestry on `main`.

### Fixed

- Op names, lockfile entries, install methods, asset names, GitHub repositories,
  and raw asset build paths now fail validation before they can reach filesystem
  or network boundaries. Raw asset `--build` overrides are restricted to the
  catalogued build list for both listing and fetching.
- Applying a lockfile now verifies its registry version, install method, commit,
  and package spec against the active registry before installing any entry,
  preventing silent metadata drift.

### Security

- Op file paths can no longer escape Arsenal's state directory through crafted
  names or mismatched names embedded in an op file.
- Persistent state, lockfiles, registries, and shims now share one atomic writer
  that uses unpredictable sibling temporary files, preventing predictable-temp
  symlink overwrites while preserving the previous destination on failure.

## [0.5.1] - 2026-09-30

### Added

- Official `go install github.com/Okymi-X/arsenal/cmd/arsenal@latest`
  installation path with module-aware version reporting.

## [0.5.0] - 2026-09-29

### Added

- Fully isolated `gobin` and `cargo` installers with exact-version package
  targets, per-install build caches, executable verification, atomic promotion,
  native `run` support, and manifest-aware removal.
- Direct `--github-ref` installation for Go and Cargo tools. Human-readable
  refs are resolved through GitHub and package managers receive only the
  resulting full commit SHA.

### Changed

- Long-running installs, GitHub queries, registry synchronization, and asset
  downloads now share the process cancellation context and stop cleanly on
  interrupt or termination signals.
- Native installer metadata now maps each exposed binary to its Go package or
  Cargo crate, including multi-binary tools whose upstream command names differ
  from Arsenal's public shims.
- Updated GitHub workflow actions to their Node 24-based releases, pinned each
  action to an immutable commit, and pinned runners to Ubuntu 24.04 for
  reproducible CI behavior.

### Security

- Asset downloads are capped at 512 MiB, checked against upstream metadata and
  actual bytes, restricted to HTTPS GitHub-controlled hosts, staged under
  unpredictable filenames, synced, and atomically promoted. Oversized,
  truncated, interrupted, or redirected-to-untrusted transfers preserve the
  previous destination and remove partial files.
- GitHub asset API responses are bounded to 4 MiB before JSON decoding.
- Native installs use private staging roots, scoped caches, shell-free process
  execution, strict path and target validation, failure cleanup, and atomic
  replacement of incomplete installations.
- Registry names, versions, binaries, and shim filenames are validated as safe
  path components. Shim targets are POSIX-quoted so configured paths remain
  data rather than executable shell syntax.

## [0.4.0] - 2026-09-29

### Added

- Official Arsenal toolbox icon under `assets/branding/` and README branding.
- `versions <tool> --github` and explicit `install <tool> --github-ref <ref>`
  workflows. GitHub refs are resolved to immutable commits before installation;
  direct upstream installs are clearly marked untested.
- `outdated` and `upgrade` commands that target the newest tested registry
  version, retain prior environments for rollback, and refuse silent downgrades
  or updates of untracked versions.
- GitHub-aware registry synchronization with `sync --list-refs`, `--ref`, and
  `--repo`, including persistence of the selected ref and resolved commit.
- Ten independently install-verified Python tools selected from Exegol's
  official image build set: `ldeep`, `bloodhound-import`, `bbot`, `fierce`,
  `ssh-audit`, `holehe`, `sherlock-project`, `maigret`, `censys`, and
  `name-that-hash`. Exact versions and CLI entrypoints were verified in clean
  Python 3.14 virtual environments.
- Topical `ad-recon`, `recon-audit`, and `osint` registry segments so future
  additions do not grow already broad category files.

### Changed

- Replaced the 40 KB generated `registry/registry.toml` catalog with a small,
  ordered manifest. Builds, local validation, upstream verification, and
  `arsenal sync` now assemble the referenced segment files directly while
  retaining support for legacy monolithic registry URLs.
- Registry segments now download with bounded concurrency, and reinstalling an
  already healthy version simply activates it instead of contacting upstream.

### Security

- Registry sync now rejects unsafe or duplicate segment paths, cross-origin
  segment resolution, oversized responses, and oversized assembled catalogs.
  Temporary files use unpredictable names and a failed sync preserves the
  previous valid registry.
- Every segmented registry file is SHA-256 pinned in the manifest. A mixed or
  tampered remote snapshot fails before the active registry is replaced.

## [0.3.1] - 2026-06-14

### Added

- Eight payload-list assets from PayloadsAllTheThings, fetchable into a
  directory and fed to a fuzzer (ffuf, wfuzz, ...): `payloads-sqli`,
  `payloads-xss`, `payloads-lfi`, `payloads-traversal`, `payloads-ssrf`,
  `payloads-xxe`, `payloads-cmdi`, and `payloads-nosqli`. Each is a `github-raw`
  collection of a category's payload folder; use `--list`, then name a file.

### Fixed

- `fetch` and the registry-check bot now percent-encode `github-raw` directory
  paths, so asset folders containing spaces (such as the PayloadsAllTheThings
  category folders) resolve correctly.
- `fetch` writes the downloaded file under its base name only, so a crafted
  upstream file name cannot escape the destination directory.
- `info` shows accurate `directory:`/`default file:` lines for assets and omits
  the empty `python:` line for tools without a Python requirement.

## [0.3.0] - 2026-06-14

### Added

- Modular registry authoring. The catalog is now split into per-category segment
  files under `registry/segments/` (`ad`, `web`, `recon`, `password`, `misc`,
  `assets`, plus `_meta`). `tools/regbuild` assembles them into the canonical
  `registry/registry.toml` (`make registry`); CI fails if that file is out of
  date with its segments. The single assembled file is still what gets embedded,
  synced, and verified.
- A large, upstream-verified batch of catalog entries:
  - AD and credentials (pip): `smbmap`, `lsassy`, `pypykatz`, `adidnsdump`,
    `donpapi`, `masky`, `pywerview`, `dploot`, `certsync`.
  - Web (pip): `sqlmap`, `dirsearch`, `wfuzz`, `arjun`, `wafw00f`, `sslyze`.
  - Recon (pip): `autorecon`, `dnsrecon`, `sublist3r`.
  - Password (pip): `hashid`. Post-exploitation (pip): `pwncat-cs`.
  - Go tools (gobin, behind the pending method): `httpx`, `subfinder`, `naabu`,
    `katana`, `waybackurls`, `gau`.
  - Upload assets (`fetch`): `mimikatz`, `sharphound`, `nanodump`,
    `printspoofer`, `godpotato`, `juicypotatong`, `lse`.
- Asset validation in the registry loader (name and source are checked).

### Changed

- `fetch` asset matching skips checksum and signature siblings (`.sha256`,
  `.asc`, ...) and prefers the shortest match, so a plain release file wins over
  a longer `+debug` variant. Single-file `github-raw` assets (such as `lse`) now
  fetch without naming a binary, via the asset's default pattern.

## [0.2.0] - 2026-06-14

### Added

- `arsenal fetch <asset> [binary]` pulls the latest version of a precompiled
  upload-binary into a directory (`--dest`, default the current directory),
  with no isolated environment and no shim - a distinct workflow from `install`
  for binaries you stage onto a target. A new `[[asset]]` registry section backs
  it, with two source kinds: `github-release` (latest release asset, matched by
  pattern or an explicit binary argument) and `github-raw` (a file from a repo
  branch, used for collections). Seeded assets: `sharpcollection` (with `--list`
  and `--build` to choose a .NET build), `winpeas`, `linpeas`, and `pspy`
  (moved here from the tool catalog). `arsenal search` now also lists matching
  assets, tagged `(asset)`, and shell completion completes asset names for
  `fetch`. The registry-check workflow verifies every asset resolves upstream.

- Six tools from 0xdf's offensive-Python toolkit, each verified against its
  official source by the registry-check workflow:
  - `bloodyad` (AD privilege-escalation swiss army knife, PyPI `bloodyad`).
  - `pywhisker` (Shadow Credentials attack, gitpip from `ShutdownRepo/pywhisker`).
  - `ldapdomaindump` (LDAP domain dumper, with the `ldd2bloodhound` and
    `ldd2pretty` helper binaries).
  - `flask-unsign` (crack and forge Flask session cookies).
  - `git-dumper` (reconstruct a source tree from an exposed `.git`).
  - `sshuttle` (transparent SSH proxy VPN for pivoting).
- Eight more 0xdf staples in Go and Rust, catalogued (with their pinned upstream
  tags and exact `go install`/`cargo install` paths in the notes) behind the
  still-pending gobin/cargo/binary install methods:
  - `kerbrute` and `pretender` (gobin) - Kerberos user enumeration/spraying and
    LLMNR/mDNS/DHCPv6 spoofing.
  - `gobuster` and `nuclei` (gobin) - content/DNS brute-forcing and templated
    vulnerability scanning.
  - `chisel` and `ligolo-ng` (gobin) - HTTP and TUN-based pivoting.
  - `rustscan` (cargo) - fast port scanner that hands off to nmap.
  - `pspy` (binary) - rootless Linux process and cron snooping.

## [0.1.2] - 2026-06-14

### Added

- Per-version `repo` override in the registry so a single version can pin a fork
  or a branch without changing the tool's canonical repo. Used to add a NetExec
  `badsuccessor` version installed from the `azoxlpf/NetExec`
  `feat/refactor-badsuccessor` branch (`arsenal install nxc@badsuccessor`).
- `arsenal completion bash|zsh|fish` prints a shell completion script that
  completes subcommands and, dynamically, registry and installed tool names.
- `arsenal run <tool> <binary>` selects a specific binary of a multi-binary
  tool (for example `arsenal run impacket getTGT`). The selector is matched
  loosely, ignoring case, a `.py` suffix, and a `<tool>-` prefix.

### Fixed

- Impacket binary names corrected to the actual installed script names
  (`secretsdump.py`, `getTGT.py`, ...); the previous `impacket-*` names did not
  exist, so `run impacket` and its shims were broken.
- `arsenal run` now rejects a flag-like first argument instead of treating it as
  a tool name.

## [0.1.1] - 2026-06-14

### Added

- GitHub Actions CI (format, vet, race tests with coverage, golangci-lint, and a
  cross-platform build matrix) and a tag-triggered release workflow that
  cross-compiles static binaries and publishes them with checksums.
- Registry verification bot: `tools/regcheck` and a `registry-check` workflow
  (on registry changes and weekly) that confirm every catalogued version exists
  at its official source - PyPI for pip tools, the upstream Git repository for
  gitpip/gobin/cargo/binary tools - so a non-existent version can never ship.
  Also available locally via `make verify-registry`.
- Dependabot for Go modules and GitHub Actions.
- CONTRIBUTING, SECURITY policy, editorconfig, issue templates, and a PR template.

### Changed

- Registry search now also matches a tool's category.
- All registry versions corrected against their official sources.

### Fixed

- NetExec is installed via `gitpip` from its pinned Git tag; it is not published
  on PyPI, so the previous `pip` entry could never install.
- `run`, `switch`, and `remove` now resolve a tool alias (for example `nxc`) to
  its canonical name before consulting the manifest.

## [0.1.0] - 2026-06-14

### Added

- Curated TOML registry of tested, known-good pentest tool versions
  (NetExec, Impacket, Certipy, BloodHound.py, mitm6, Coercer, and more).
- Virtualenv isolation backend: one Python venv per tool/version.
- pip and git+pip install methods driven by an orchestrator that selects the
  right method from the registry entry.
- PATH shim system so multiple versions coexist and the active one is
  switchable.
- Engagement profiles ("ops") with reproducible TOML lockfiles, plus
  create, pin, use, list, export, and import subcommands.
- Commands: install, remove, switch, list, search, info, run, op, sync,
  doctor, version.
- doctor command with directory, Python, PATH, and manifest checks and a
  `--fix` repair mode.
- Embedded registry so the tool works offline out of the box.
- Build version injected at link time via -ldflags.

### Stubbed

- Container isolation backend (podman/docker), wired behind the Backend
  interface. See docs/ARCHITECTURE.md.
- Offline bundle export/import, wired behind the Exporter interface.
- binary, go install, and cargo install methods, wired behind the
  InstallMethod interface.

[Unreleased]: https://github.com/Okymi-X/arsenal/compare/v0.5.4...HEAD
[0.5.4]: https://github.com/Okymi-X/arsenal/compare/v0.5.3...v0.5.4
[0.5.3]: https://github.com/Okymi-X/arsenal/compare/v0.5.2...v0.5.3
[0.5.2]: https://github.com/Okymi-X/arsenal/compare/v0.5.1...v0.5.2
[0.5.1]: https://github.com/Okymi-X/arsenal/compare/v0.5.0...v0.5.1
[0.5.0]: https://github.com/Okymi-X/arsenal/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/Okymi-X/arsenal/compare/v0.3.1...v0.4.0
[0.3.1]: https://github.com/Okymi-X/arsenal/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/Okymi-X/arsenal/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/Okymi-X/arsenal/compare/v0.1.2...v0.2.0
[0.1.2]: https://github.com/Okymi-X/arsenal/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/Okymi-X/arsenal/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/Okymi-X/arsenal/releases/tag/v0.1.0
