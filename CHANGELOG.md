# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/Okymi-X/arsenal/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/Okymi-X/arsenal/releases/tag/v0.1.0
