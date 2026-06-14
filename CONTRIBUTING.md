# Contributing to arsenal

Thanks for helping improve arsenal. This project aims to be clean, minimal, and
auditable enough to ship in distribution repositories. Please keep changes in
that spirit.

## Ground rules

- Single Responsibility Principle: one concern per package, one primary
  responsibility per file, one job per function.
- Small files (soft cap ~180 lines) and small functions. Split by
  responsibility rather than growing a file.
- No emoji anywhere: code, comments, docs, commit messages, or CLI output. Use
  ASCII markers: `[ok]`, `[fail]`, `[warn]`, `->`.
- Program to interfaces; inject concrete implementations. No global mutable
  state. Wrap errors with context (`fmt.Errorf("...: %w", err)`).
- Every exported symbol and every package has a godoc comment.

## Development

```
make build     # compile into bin/
make test      # go test ./...
make lint      # golangci-lint (config in .golangci.yml)
make fmt       # gofumpt -w .
```

Before opening a pull request, ensure `go vet ./...`, `gofmt -l .` (empty),
`make test`, and `make lint` all pass. CI runs the same checks.

## Tests

Add table-driven tests for behavior changes. Mock the `isolation.Backend` and
`installer.InstallMethod` interfaces rather than touching the network or the
real filesystem where avoidable.

## Adding tools to the registry

Edit the per-category segment file under `registry/segments/` (not the generated
`registry/registry.toml`), then run `make registry` to reassemble it and `make
verify-registry` to confirm the entry resolves upstream. Follow
`docs/registry-format.md`. Only set `tested = true` on a version you have
actually verified. Include the pinned `commit` for `gitpip` tools and a
`pip_spec` for `pip` tools. Precompiled upload-binaries go in
`registry/segments/assets.toml` as `[[asset]]` blocks.

## Commits and versioning

- Use [Conventional Commits](https://www.conventionalcommits.org/) for messages
  (`feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `build:`, `ci:`, `chore:`).
- The project follows [Semantic Versioning](https://semver.org/). Releases are
  tagged `vMAJOR.MINOR.PATCH` and update `CHANGELOG.md`
  ([Keep a Changelog](https://keepachangelog.com/)).
