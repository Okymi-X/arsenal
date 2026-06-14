# Registry format

The registry is a single TOML file. It is the canonical, in-repo catalog of
tested, known-good tool versions. It is embedded in the binary (so arsenal works
offline) and written to the user's data directory on first run. `arsenal sync`
refreshes it from the configured upstream URL.

## Top level

```toml
version = "1"          # registry schema version
updated = "2026-06-14" # ISO-8601 date of the last revision

[[tool]]
# ... one block per tool
```

| Field     | Type   | Required | Description                          |
|-----------|--------|----------|--------------------------------------|
| `version` | string | yes      | Registry schema version.             |
| `updated` | string | no       | ISO-8601 date the registry changed.  |
| `tool`    | array  | yes      | Array of tool tables (`[[tool]]`).   |

## Tool entry

```toml
[[tool]]
name = "netexec"
aliases = ["nxc", "cme"]
description = "Network execution tool, successor to CrackMapExec"
repo = "https://github.com/Pennyw0rth/NetExec"
homepage = "https://www.netexec.wiki"
category = "ad"
install_method = "pip"
python_version = ">=3.9"
binary = "nxc"
binaries = ["nxcdb"]
dependencies = []
tags = ["smb", "ldap", "active-directory"]
requires_root = false
notes = "Use the tested 1.x line for stable engagements."
```

| Field            | Type     | Required | Description                                                |
|------------------|----------|----------|------------------------------------------------------------|
| `name`           | string   | yes      | Canonical identifier used on the command line.             |
| `aliases`        | [string] | no       | Alternative names that resolve to this tool.               |
| `description`    | string   | no       | One-line summary.                                          |
| `repo`           | string   | no       | Source repository URL (required for `gitpip`).             |
| `homepage`       | string   | no       | Project's primary URL.                                     |
| `category`       | string   | yes      | One of `ad`, `web`, `recon`, `password`, `exploit`, `c2`, `misc`. |
| `install_method` | string   | yes      | One of `pip`, `gitpip`, `binary`, `gobin`, `cargo`.        |
| `python_version` | string   | no       | Minimum interpreter, e.g. `>=3.9`.                         |
| `binary`         | string   | no       | Primary executable exposed via a shim.                     |
| `binaries`       | [string] | no       | Additional executables exposed via shims.                  |
| `dependencies`   | [string] | no       | Extra system or pip requirements, as notes.                |
| `tags`           | [string] | no       | Free-form search keywords.                                 |
| `requires_root`  | bool     | no       | Whether the tool typically needs elevated privileges.      |
| `notes`          | string   | no       | Operational guidance.                                      |
| `version`        | array    | yes      | Array of version tables (`[[tool.version]]`).              |

## Version entry

Versions are listed newest first. The default selection picks the newest
version whose `tested` flag is true; if none is tested, the newest version is
used.

```toml
  [[tool.version]]
  tag = "1.4.0"
  commit = ""
  tested = true
  pip_spec = "netexec==1.4.0"
  date = "2025-12-01"
  notes = "Recommended stable release."
```

| Field      | Type   | Required | Description                                        |
|------------|--------|----------|----------------------------------------------------|
| `tag`      | string | yes      | Human-facing version label, e.g. `1.4.0`.          |
| `commit`   | string | no       | Exact source revision (used by `gitpip`).          |
| `tested`   | bool   | no       | Marks a version the maintainers verified.          |
| `pip_spec` | string | no       | pip requirement string, e.g. `netexec==1.4.0`.     |
| `date`     | string | no       | ISO-8601 release date.                             |
| `notes`    | string | no       | Version-specific guidance.                         |

## Install method semantics

- `pip` - installs `pip_spec` (or `name==tag` if absent) into a venv.
- `gitpip` - installs `git+<repo>@<commit>` into a venv.
- `binary`, `gobin`, `cargo` - reserved; implementations are stubbed and will
  fail loudly until completed (see docs/architecture.md).

## Validation rules

The loader rejects a registry where a tool has no `name`, no `install_method`,
or no versions, or where two tools share a `name`.

## Upstream verification

Every version is verified against its official source by the `registry-check` CI
workflow (and locally via `make verify-registry`):

- `pip` tools: the version from `pip_spec` (or `name==tag`) must be published on
  PyPI.
- `gitpip`, `gobin`, `cargo`, `binary` tools: the `commit` or `tag` must resolve
  as a ref in the GitHub repository (the checker also tries a leading `v`).

A pull request that adds a version which does not exist upstream fails CI, so a
non-existent version cannot be merged. The workflow also runs weekly to catch
versions that were later yanked or moved.
