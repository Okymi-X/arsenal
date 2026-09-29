# Contributing to arsenal

Thanks for helping improve arsenal. This project aims to be clean, minimal, and
auditable enough to ship in distribution repositories. Please keep changes in
that spirit.

## Ground rules

All changes must follow [docs/RULES.md](docs/RULES.md), the single source of
truth for responsibility boundaries, file and function size, reuse, secure
design, testing, and completion criteria. Read [docs/PRD.md](docs/PRD.md) and
[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) before changing product behavior or
package boundaries.

## Development

```
make build     # compile into bin/
make test      # go test ./...
make lint      # golangci-lint (config in .golangci.yml)
make fmt       # gofumpt -w .
```

Before opening a pull request, run the completion checks in
[docs/RULES.md](docs/RULES.md). CI repeats the applicable checks.

## Tests

Add table-driven tests for behavior changes. Mock the `isolation.Backend` and
`installer.InstallMethod` interfaces rather than touching the network or the
real filesystem where avoidable.

## Adding tools to the registry

Edit the appropriate topical file under `registry/segments/`, then run `make
registry` to validate the manifest and `make verify-registry` to confirm every
entry resolves upstream. Add a path to the small `registry/registry.toml`
manifest only when creating a segment. Follow `docs/registry-format.md`. Only
set `tested = true` on a version you have actually installed and invoked in a
clean environment. Include the pinned `commit` for `gitpip` tools and a
`pip_spec` for `pip` tools. Precompiled upload-binaries go in an asset segment
as `[[asset]]` blocks.

## Commits and versioning

- Use [Conventional Commits](https://www.conventionalcommits.org/) for messages
  (`feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `build:`, `ci:`, `chore:`).
- The project follows [Semantic Versioning](https://semver.org/). Releases are
  tagged `vMAJOR.MINOR.PATCH` and update `CHANGELOG.md`
  ([Keep a Changelog](https://keepachangelog.com/)).
