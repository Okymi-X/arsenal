# Design guide

This document captures product-facing and implementation design decisions. The
package dependency map lives in `ARCHITECTURE.md`; mandatory constraints live in
`RULES.md`.

## Design priorities

When choices conflict, prefer the option that is safer, more reproducible,
easier to audit, and simpler to operate. Convenience does not justify hidden
network access, ambiguous version selection, shell interpretation, or writes
outside arsenal-owned or explicitly selected paths.

## CLI model

Commands use a predictable verb-first shape:

```text
arsenal <command> [subject] [flags]
```

- Keep global flags scarce. Put command-specific behavior on the owning
  command.
- Require explicit names for state-changing operations.
- Use `tool@version` consistently wherever a version can be selected.
- Pass arguments following `--` to the selected tool unchanged.
- Print primary results to stdout and actionable diagnostics to stderr.
- Keep quiet mode suitable for interactive use and automation. Verbose mode may
  expose operational detail but never secrets.
- Use plain text and ASCII status markers. Do not require color to communicate
  meaning.
- Error messages state the failed operation and, where possible, the corrective
  action. A command must not print success before durable work succeeds.

## Domain boundaries

- `cli` parses intent and renders results; it does not own registry,
  installation, or persistence rules.
- `registry` owns catalog parsing, validation, querying, and synchronization.
- `resolver` maps a user selector to one concrete catalog version.
- `installer` selects an installation strategy and coordinates installation.
- `isolation` owns lifecycle operations for an isolated environment.
- `fetcher` resolves and downloads standalone assets; it does not update the
  installed-tool manifest or create shims.
- `store`, `op`, and `shim` each own one persistence format or filesystem side
  effect.
- `doctor` observes consistency and performs only explicit, bounded repairs.

A new feature belongs in the package that owns its invariant. If no package has
that responsibility, add a focused package rather than turning `cli` into a
domain layer.

## Data and state

- `registry/registry.toml` is the ordered catalog manifest. Tool and asset
  records live in topical segment files, which are assembled and validated at
  build, local-load, and remote-sync boundaries.
- Manifests describe local installed state; lockfiles describe reproducible
  requested state. Do not merge their responsibilities.
- Persisted formats should be reviewable text with stable ordering.
- Write complete state to a sibling temporary file, sync when durability
  matters, set intended permissions, and rename into place.
- Schema changes require backward-compatible reading or an explicit migration.

## External processes

- Build commands from a fixed executable and validated argument list.
- Inherit only the environment needed by the child process; override variables
  deliberately.
- Connect stdout and stderr according to the CLI contract.
- Propagate cancellation and preserve the child exit failure with context.
- Never infer elevated privileges or execute through a shell.

## Network operations

- Centralize upstream-specific URL construction and response parsing.
- Use clients with timeouts and contexts; do not use an unbounded default client
  for runtime operations.
- Bound response sizes before buffering them in memory or on disk.
- Validate metadata before selecting a download and validate content before
  replacing durable state.
- Resolve user-selected GitHub refs to immutable commits and persist the
  human-readable ref separately from the resolved commit.
- Bound concurrency as well as individual response size for segmented pulls.
- Keep network behavior injectable so tests can use local servers.

## Failure design

Failures should be local and recoverable:

- A failed registry sync retains the previous valid registry.
- A failed download leaves no final-looking asset.
- A failed install does not become an active version.
- A failed shim update does not destroy an existing working shim.
- Repair commands report each action and do not modify unrelated files.

Retries are appropriate only for transient, idempotent work and must be bounded.
Invalid input, authentication failure, and integrity failure are not retryable.

## Adding a feature

1. Confirm it is within `PRD.md` or update the product decision explicitly.
2. Identify the owning package and trust boundaries.
3. Define the smallest contract and failure behavior.
4. Implement the domain behavior independently of CLI rendering.
5. Add allowed, denied, boundary, and cleanup tests as relevant.
6. Wire the feature into the CLI and update usage documentation.
7. Run the completion checks from `RULES.md`.
