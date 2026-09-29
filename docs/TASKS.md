# Work plan

This file tracks verified, repository-level work that is not yet complete. It
is not a session log. Keep entries small, testable, and linked to a concrete
package or behavior. Remove an entry when it is delivered and record
user-visible completion in `CHANGELOG.md`.

## Priority order

1. Prevent unsafe or corrupt behavior.
2. Complete already-exposed interfaces and documented workflows.
3. Improve reproducibility and verification.
4. Add new catalog breadth or convenience.

## Active backlog

### Integrity verification for fetched and installed content

Status: registry segment integrity completed; package and asset integrity
remain planned

- Extend package and asset metadata with checksums or signatures.
- Verify package and asset content before atomic promotion or installation.
- Define a safe compatibility path for existing entries without integrity
  metadata; never present those entries as cryptographically verified.
- Cover mismatch, missing metadata, cleanup, and valid-content cases.

### Binary install method

Status: stubbed in `internal/installer/binary.go`

- Resolve a platform-specific immutable release asset.
- Download with bounds, verify integrity, and install atomically.
- Reject unsupported platforms, unsafe archive paths, and ambiguous assets.

### Container isolation backend

Status: stubbed in `internal/isolation/container`

- Support a configured rootless Podman or Docker runtime behind
  `isolation.Backend`.
- Pin base images by digest, drop unnecessary capabilities, avoid privileged
  mode, and constrain mounts to explicit paths.
- Specify lifecycle, networking, cache, and cleanup behavior before enabling it
  in configuration.

### Offline bundles

Status: stubbed in `internal/bundle`

- Export a lockfile with required source artifacts and integrity metadata.
- Import only after validating the complete bundle, including archive paths and
  size limits.
- Reproduce an environment without network access or silent version changes.

## Definition of ready

A task is ready for implementation when its user outcome, owning package,
trust boundaries, compatibility constraints, and acceptance tests are clear.
If any of those would materially change product scope, update `PRD.md` or
`DESIGN.md` before writing code.
