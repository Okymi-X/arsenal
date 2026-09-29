# Product requirements

## Product statement

arsenal is a package and environment manager for authorized offensive-security
work. It provides a curated, reproducible catalog of known-good tool versions
and isolates installed versions so operators can switch between them without
polluting the host environment.

The curated registry is the product's core value. arsenal orchestrates existing
package managers and runtimes; it does not reimplement security tools or bundle
a Python interpreter.

## Users and primary jobs

The primary user is an authorized security practitioner who needs to:

- discover a tool and a tested version without researching it mid-engagement;
- install multiple versions without dependency conflicts;
- run or switch the active version predictably;
- reproduce a team's toolset from a reviewed lockfile;
- stage selected standalone assets in an explicit destination; and
- diagnose broken local installations without affecting unrelated tools.

## Functional requirements

### Registry

- Ship a valid embedded registry so read-only discovery works offline.
- Model source, install method, supported binaries, runtime requirements,
  versions, immutable revision data where available, test status, and notes.
- Choose the newest tested version by default and honor explicit version pins.
- Build the distributed registry deterministically from an ordered manifest and
  reviewed topical segments, and reject malformed entries.
- Replace a synced registry only after the complete download parses and
  validates successfully.
- Let operators discover GitHub tags and pin registry synchronization to a
  selected tag, branch, or commit resolved to an immutable revision.

### Isolation and installation

- Isolate every installed tool/version under arsenal's data root.
- Support multiple installed versions and one explicit active version.
- Keep install mechanisms behind a common boundary so new mechanisms do not
  leak implementation details into the CLI.
- Fail clearly when a catalogued install mechanism is not implemented.
- Never report a partial or failed install as successful.
- Report safe tested upgrades, retain prior environments for rollback, and
  never silently downgrade or replace an untracked version.
- Permit an explicit untested GitHub-ref install for supported Python, Go, and
  Cargo projects only after resolving the ref to an immutable commit.

### Execution and shims

- Pass tool arguments without reinterpretation by a shell.
- Resolve aliases to canonical registry names before consulting local state.
- Generate predictable shims for exposed binaries and update them atomically
  where practical.
- Keep normal output concise, send diagnostics to stderr, and avoid color when
  stdout is not a terminal.

### Engagement profiles

- Pin exact tool versions in named profiles.
- Export and import a deterministic, reviewable lockfile.
- Reproduce a profile without silently substituting versions.
- Surface missing or unsupported dependencies as actionable errors.

### Asset fetching

- Fetch only a registry-declared asset from its declared upstream source.
- Write only beneath the operator-selected destination, using a safe base name.
- Use an incomplete temporary file and expose the final file only after a
  successful transfer.
- Reject ambiguous, malformed, or non-successful upstream responses.
- Support future integrity metadata without changing the user-facing workflow.

### Diagnostics and operations

- Detect missing directories, runtimes, shims, and inconsistent manifest
  entries.
- Make repair explicit and limit it to arsenal-owned state.
- Collect no telemetry and do not make unrelated network requests.

## Quality attributes

- Reproducible: version selection and lockfiles produce stable results.
- Auditable: dependencies, registry changes, and filesystem effects are easy to
  review.
- Secure by default: inputs are validated at trust boundaries; paths cannot
  escape their roots; subprocesses avoid shell interpretation; network and file
  operations are bounded and failure-safe.
- Maintainable: each package, file, and function has one clear responsibility;
  shared policy and logic have one authoritative implementation.
- Portable: the distributed binary is static and supports the documented Linux
  and macOS architectures.
- Accessible: CLI messages are plain text, readable without color, and do not
  depend on symbols or animation.

## Non-goals

- Auditing or guaranteeing the behavior of third-party tools.
- Automatically running security tools against targets.
- Managing system packages or replacing language package managers.
- Hiding authorization, scope, or safety decisions from the operator.
- General-purpose project environment management.
- Silent privilege escalation or system-wide mutation.

## Release acceptance

A release is acceptable when:

- registry manifest and segment validation succeed;
- formatting, vetting, linting, race-enabled tests, and supported builds pass;
- behavior changes have tests for success and relevant failure paths;
- security-sensitive changes cover rejected inputs and cleanup after failure;
- user-visible changes update usage documentation and the changelog; and
- known incomplete features remain explicit rather than appearing functional.
