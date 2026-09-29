# Engineering rules

This document is the single source of truth for implementation, review, and
verification rules in this repository. `MUST` and `MUST NOT` are mandatory.
`SHOULD` identifies the default; deviations need a concrete reason visible in
the code or review.

## Responsibilities and size

- A package MUST own one cohesive domain concern.
- A production source file MUST have one primary reason to change.
- A function MUST perform one job at one level of abstraction.
- Production Go files SHOULD remain below 180 lines and functions SHOULD remain
  below 40 lines. These are review triggers, not targets: split only along a
  real responsibility boundary, and do not fragment cohesive code to satisfy a
  count.
- Cyclomatic complexity MUST remain below the configured lint threshold of 15.
- Dependency wiring belongs at the application boundary. Domain packages MUST
  NOT construct unrelated concrete dependencies or depend on the CLI package.
- Introduce an interface only for a real boundary, multiple implementations, or
  test isolation. Keep it as small as its consumer requires.

## Duplication and reuse

- Business rules, validation, path-safety logic, protocol parsing, and security
  decisions MUST have one authoritative implementation.
- Search for an existing implementation before adding one. Extend or extract
  the owning helper when the same policy would otherwise exist twice.
- Do not abstract coincidentally similar code with different reasons to change.
  Duplication of syntax is cheaper than coupling unrelated concepts.
- Generated files MUST be produced from their source and MUST NOT become a
  second hand-maintained source of truth.

## Go implementation

- Preserve public interfaces unless a breaking change and migration are part of
  the approved task.
- Pass dependencies explicitly. Package-level mutable state is prohibited.
- Accept `context.Context` for operations that can block or cross a process,
  filesystem, or network boundary, and propagate cancellation where supported.
- Wrap errors with an actionable operation and `%w`. Error text begins with a
  lowercase word and does not end with punctuation.
- Handle every meaningful error. Cleanup errors may be ignored only when an
  earlier error already determines the result and cleanup is best effort.
- Use typed data and the standard library before adding a dependency. A new
  dependency requires a maintained source, a clear ownership benefit, and a
  review of transitive and supply-chain risk.
- Exported packages and symbols MUST have useful godoc comments.
- Comments explain invariants, intent, or non-obvious tradeoffs. They MUST NOT
  narrate obvious code or contain generated filler.
- Code, comments, docs, logs, UI copy, and commits MUST NOT contain emoji.

## Secure design and logic

Treat registry data, remote responses, configuration, command arguments,
lockfiles, manifests, filenames, URLs, archives, and environment variables as
untrusted inputs.

- Validate once at each trust boundary, normalize into a safe internal form,
  and keep invalid states out of business logic.
- Authorization and engagement scope remain operator responsibilities. arsenal
  MUST NOT automatically discover, select, or attack targets.
- Subprocesses MUST use an executable plus an argument slice. Never construct a
  shell command from untrusted text or invoke a shell for convenience.
- File destinations MUST be resolved beneath the intended root. Reject absolute
  paths, traversal, unsafe base names, and symlink-based escapes where an
  attacker can influence the path.
- Files MUST use the least permissions required. Secret-bearing files, if ever
  introduced, MUST default to owner-only access and MUST NOT be logged.
- Network clients MUST use timeouts or caller cancellation, require expected
  status codes, close response bodies, and apply a justified size limit before
  buffering untrusted responses.
- Remote content MUST be written to a temporary file, validated completely, and
  atomically promoted. Failure MUST leave the previous good state intact and
  remove incomplete artifacts where practical.
- Immutable revisions and checksums MUST be verified when registry metadata
  provides them. Never label unverified content as verified or tested.
- Parsers MUST reject malformed, ambiguous, duplicate, or unsupported data
  rather than guessing. Do not silently downgrade security behavior.
- Logs and errors MUST NOT expose credentials, tokens, private keys, sensitive
  environment values, or unnecessary target data.
- Cryptography and token handling MUST use maintained libraries and established
  formats; custom cryptography is prohibited.
- Security controls MUST NOT be weakened to make a test pass.

For any security-sensitive change, identify the asset, trust boundary, attacker
controlled input, abuse case, and safe failure mode before implementation. The
result can be brief, but it must be reflected in code and tests.

## CLI and data behavior

- Normal output goes to stdout; diagnostics and failures go to stderr.
- Output MUST remain useful without color and stable enough for documented
  automation. Breaking output changes require explicit documentation.
- Destructive actions MUST target an explicit arsenal-owned path. Broad,
  unresolved, or empty paths are invalid.
- Persistent state updates SHOULD be atomic and deterministic. Readers MUST
  receive either the old valid state or the new valid state, never a partial
  write.
- Configuration and wire-format additions SHOULD be backward compatible and
  have deterministic defaults.

## Tests and verification

- Every bug fix MUST add a regression test that fails before the fix.
- New logic MUST cover the normal path, relevant boundary conditions, and
  failure behavior. Security controls require both allowed and denied cases.
- Prefer table-driven unit tests and injected fakes for process, network, and
  filesystem boundaries. Tests MUST NOT depend on public services unless they
  are explicitly integration tests.
- Tests MUST be deterministic, isolated, and safe to run in parallel when
  marked parallel.
- Run the narrowest relevant test during development. Before completion, run:

```text
gofmt -l .
go vet ./...
go test -race ./...
golangci-lint run ./...
make build
```

If a command cannot run, report the exact command, reason, and closest check
performed. Do not claim unrun verification.

## Completion criteria

A change is complete only when code, tests, documentation, and generated files
agree; relevant checks pass or blockers are disclosed; no unrelated changes
were overwritten; and temporary files, debug output, dead code, and secrets are
absent from the diff.
