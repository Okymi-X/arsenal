# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
  interface. See docs/architecture.md.
- Offline bundle export/import, wired behind the Exporter interface.
- binary, go install, and cargo install methods, wired behind the
  InstallMethod interface.

[Unreleased]: https://github.com/Okymi-X/arsenal/compare/v0.1.2...HEAD
[0.1.2]: https://github.com/Okymi-X/arsenal/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/Okymi-X/arsenal/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/Okymi-X/arsenal/releases/tag/v0.1.0
