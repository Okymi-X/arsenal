# Registry format

The registry is the canonical, in-repo catalog of tested, known-good tool
versions. It is embedded in the binary (so arsenal works offline) and written to
the user's data directory on first run. `arsenal sync` refreshes it from the
configured upstream URL.

## Authoring: manifest and segments

`registry/registry.toml` is a small ordered manifest. Tool and asset entries
live in topical files under `registry/segments/`; each segment contains only
`[[tool]]` or `[[asset]]` blocks. New topics can be added without growing an
existing category file.

Builds embed the manifest and every segment. Local loads and remote sync verify
each segment's SHA-256, expand the files in manifest order, validate the
complete catalog, and expose one in-memory registry. Remote segments download
concurrently with a fixed worker bound. Legacy monolithic registries remain
readable for backward compatibility.

To change the catalog:

```
# 1. edit the relevant registry/segments/*.toml
# 2. bump `updated` in registry/registry.toml
# 3. add the path to `segments` only when creating a new segment
make registry
# 4. verify every entry still resolves upstream
make verify-registry
```

`make registry` and CI reject missing, unlisted, duplicate, unsafe, or invalid
segments.

## Top level

```toml
version = "1"          # registry schema version
updated = "2026-09-29" # ISO-8601 date of the last revision
segments = [
  "segments/ad.toml",
  "segments/ad-recon.toml",
  "segments/web.toml",
]

[segment_sha256]
"segments/ad.toml" = "<64 lowercase hexadecimal characters>"
"segments/ad-recon.toml" = "<64 lowercase hexadecimal characters>"
"segments/web.toml" = "<64 lowercase hexadecimal characters>"
```

| Field      | Type     | Required | Description                         |
|------------|----------|----------|-------------------------------------|
| `version`  | string   | yes      | Registry schema version.            |
| `updated`  | string   | no       | ISO-8601 date the registry changed. |
| `segments` | [string] | yes      | Ordered relative segment paths.     |
| `segment_sha256` | table | yes | Generated SHA-256 pin for each segment. |

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
  pip_dependencies = ["dploot==3.1.2"]
  date = "2025-12-01"
  notes = "Recommended stable release."
```

| Field      | Type   | Required | Description                                              |
|------------|--------|----------|----------------------------------------------------------|
| `tag`      | string | yes      | Human-facing version label, e.g. `1.4.0`.                |
| `commit`   | string | no       | Source revision for `gitpip`: commit SHA, tag, or branch.|
| `repo`     | string | no       | Overrides the tool's repo for this version (fork/branch).|
| `tested`   | bool   | no       | Marks a version the maintainers verified.                |
| `pip_spec` | string | no       | pip requirement string, e.g. `impacket==0.12.0`.         |
| `pip_dependencies` | [string] | no | Additional exact `name==version` requirements installed with this version. |
| `install_targets` | table | for Go/Cargo | Exposed binary to Go package or Cargo crate mapping. |
| `date`     | string | no       | ISO-8601 release date.                                   |
| `notes`    | string | no       | Version-specific guidance.                               |

A version may pin a fork or a branch by combining `repo` and `commit`. For
example, to install NetExec from a fork's feature branch:

```toml
  [[tool.version]]
  tag = "badsuccessor"
  commit = "feat/refactor-badsuccessor"
  repo = "https://github.com/azoxlpf/NetExec"
  tested = false
  notes = "Fork branch; tracks a moving branch, not reproducible."
```

Install it with `arsenal install nxc@badsuccessor`. Branch refs move over time,
so leave such versions `tested = false`.

## Install method semantics

- `pip` - installs `pip_spec` (or `name==tag` if absent) and any exact
  `pip_dependencies` into a venv in one resolver transaction.
- `gitpip` - installs `git+<repo>@<commit>` and any exact
  `pip_dependencies` into a venv in one resolver transaction.
- `gobin` - runs `go install <install_target>@<commit>` with isolated `GOBIN`,
  module, and build caches, then verifies and promotes the declared binaries.
- `cargo` - runs `cargo install --version <tag> --locked` with an isolated root
  and cache, then verifies and promotes the declared binaries.
- `binary` - reserved; the implementation remains a tracked stub.

Every `gobin` and `cargo` version must map all exposed binaries. Keys are the
public shim names and values are package-manager targets. This supports safe
renaming without parsing notes, for example:

```toml
install_targets = { "ligolo-proxy" = "github.com/nicocha30/ligolo-ng/cmd/proxy", "ligolo-agent" = "github.com/nicocha30/ligolo-ng/cmd/agent" }
```

## Asset entry

Assets are precompiled binaries that `arsenal fetch` downloads to a directory
for staging onto a target. They are not installed, isolated, versioned in the
manifest, or shimmed: arsenal always pulls the latest upstream version.

```toml
[[asset]]
name = "sharpcollection"
description = "Precompiled .NET offensive binaries"
repo = "https://github.com/Flangvik/SharpCollection"
category = "ad"
source = "github-raw"          # or "github-release"
branch = "master"              # github-raw: ref to read
dir = "NetFramework_4.7_x64"   # github-raw: default build directory
builds = ["NetFramework_4.7_x64", "NetFramework_4.5_x64"]  # allowed --build values
tags = ["dotnet", "upload"]
notes = "..."
```

| Field         | Type     | Required | Description                                                        |
|---------------|----------|----------|--------------------------------------------------------------------|
| `name`        | string   | yes      | Canonical identifier used on the command line.                     |
| `aliases`     | array    | no       | Alternative names.                                                 |
| `description` | string   | no       | One-line summary.                                                  |
| `repo`        | string   | yes      | GitHub repository URL.                                             |
| `category`    | string   | no       | Offensive-security domain.                                         |
| `source`      | string   | yes      | `github-release` or `github-raw`.                                  |
| `pattern`     | string   | no       | github-release: default asset filename to match.                  |
| `branch`      | string   | no       | github-raw: ref to read (default `master`).                       |
| `dir`         | string   | no       | github-raw: default in-repo directory of selectable binaries.     |
| `builds`      | array    | no       | github-raw: allowed `--build` overrides for `dir`.                |
| `tags`        | array    | no       | Search keywords.                                                   |
| `notes`       | string   | no       | Operational guidance.                                              |

A `github-release` asset fetches the file matching `pattern` (or a binary named
on the command line) from the repository's latest release. A `github-raw` asset
fetches a file from `dir` (overridable with `--build`) on `branch`; a binary
name is required and `--list` enumerates the directory.

## Validation rules

Manifest paths must be relative `.toml` paths with no traversal or duplicates.
Every segment must have exactly one valid SHA-256 pin. Remote segments and
redirects must remain on the manifest's origin, and each response plus the
assembled catalog has a size bound. The loader also rejects a registry where a
tool has unsafe path components, no `install_method`, or no versions; duplicate
version tags or tool names; non-exact or option-like Python dependency pins;
and incomplete Go/Cargo target mappings.

## Upstream verification

Every version and asset is verified against its official source by the
`registry-check` CI workflow (and locally via `make verify-registry`):

- `pip` tools: the version from `pip_spec` (or `name==tag`) must be published on
  PyPI.
- Every `pip_dependencies` entry must use exact `name==version` syntax and that
  release must be published on PyPI.
- `gitpip`, `gobin`, `cargo`, `binary` tools: the `commit` or `tag` must resolve
  as a ref in the GitHub repository (the checker also tries a leading `v`).
- `github-release` assets: the latest release must contain a file matching
  `pattern`. `github-raw` assets: the `dir` directory must exist and be
  non-empty.

A pull request that adds a version which does not exist upstream fails CI, so a
non-existent version cannot be merged. The workflow also runs weekly to catch
versions that were later yanked or moved.
