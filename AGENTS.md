# Repository instructions for coding agents

These instructions apply to every automated contributor working in this
repository, including Codex. They apply recursively from the repository root.

## Required context

Before changing code, read the following documents in order:

1. `docs/RULES.md` for the mandatory engineering and security rules.
2. `docs/PRD.md` for product scope and acceptance criteria.
3. `docs/ARCHITECTURE.md` for boundaries, dependencies, and trust zones.
4. `docs/DESIGN.md` for CLI and implementation decisions.
5. `docs/TASKS.md` for the current work plan.
6. `docs/MEMORY.md` for durable, verified project facts.

Then inspect the files and tests directly related to the requested change. The
code and executable tests are authoritative for current behavior; do not rely
on memory when they disagree.

## Working agreement

- Follow `docs/RULES.md`. It is the single source of truth for code quality,
  security, testing, and completion criteria; do not restate it here.
- Preserve unrelated work and public behavior unless the requested change
  explicitly requires a migration or breaking change.
- Make the smallest complete change. Keep responsibilities narrow and reuse
  existing logic instead of copying it.
- Add or update tools in `registry/segments/*.toml`. Edit the small
  `registry/registry.toml` manifest only when adding, removing, or reordering a
  segment, then run `make registry` to validate the complete catalog.
- Update the relevant document when product behavior, architecture, design,
  operational guidance, or durable project facts change.
- Record only verified, durable facts in `docs/MEMORY.md`. Do not store secrets,
  credentials, private target data, or temporary session notes there.
- Run the narrowest relevant tests first, followed by the required checks in
  `docs/RULES.md` when the environment permits.

More specific instructions in a subdirectory override this file for that
subtree, but they may not weaken security requirements.
