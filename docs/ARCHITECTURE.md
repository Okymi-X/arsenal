# Architecture

This document defines the runtime structure, ownership boundaries, and trust
boundaries of arsenal. Mandatory implementation constraints live in
`RULES.md`.

arsenal is structured around the Single Responsibility Principle: one concern
per package, one primary responsibility per file, one job per function. Files
are kept small (soft cap ~180 lines). Backends and install methods are defined
as interfaces and injected; no package hard-references another's concrete
implementation across a boundary, and there is no global mutable state.

## Flow of control

```
cmd/arsenal/main.go        wire dependencies, call cli, exit
        |
internal/cli               parse args, dispatch to a command
        |
internal/registry          load/sync/query the curated catalog
internal/resolver          turn "tool[@version]" into a concrete version
internal/installer         select an install method and drive it
internal/fetcher           download upload-binaries (assets) to a directory
internal/isolation         isolate the environment (venv or container)
internal/shim              write PATH shims to the active version
internal/store             persist the installed-tools manifest
internal/op                ops and lockfiles for reproducibility
internal/doctor            health checks and repair
internal/safepath          validate filesystem path components
internal/strictdecode      reject unsupported persisted-schema fields
internal/fsutil            atomically replace complete files
```

`main` does nothing but construct `config.Paths`, load `config.Config`, and hand
them to `cli.New`. The CLI owns wiring of the registry source, the store, the
shim manager, and the op manager, and constructs an isolation backend and
installer orchestrator per command invocation.

`internal/fetcher` is a separate path used by `arsenal fetch`. It bypasses the
installer, isolation, shim, and store layers entirely: given a registry
`Asset`, it resolves the latest upstream file (a GitHub release asset or a raw
repository file) and downloads it to an operator-chosen directory for staging
onto a target. Assets are never isolated, versioned in the manifest, or shimmed
onto the operator's PATH.

Managed tool commands are resolved from the manifest and executed from their
absolute installation path by `arsenal run`. The private shim directory is not
required on the host `PATH`; this prevents Arsenal tools from shadowing system,
pipx, or user-managed commands with the same name.

## Key interfaces

These are defined first and implemented against. Concrete types are injected.

### isolation.Backend

Isolates a single tool/version environment.

```
Create(ctx context.Context, tool, version string) error
Install(ctx context.Context, spec InstallSpec) error
Run(ctx context.Context, args []string) error
Remove() error
Path() string
Exists() bool
```

Implementations:

- `isolation/venv` - one Python virtualenv per tool/version (default).
- `isolation/container` - podman/docker backend. Stubbed; see below.

### installer.InstallMethod

Installs one class of tool. The orchestrator owns selection; methods only know
how to install.

```
Supports(tool registry.Tool) bool
Install(ctx context.Context, tool registry.Tool, version registry.Version) (Result, error)
```

Implementations: `pip` and `gitpip` drive an isolation backend; `gobin` and
`cargo` stage exact package versions in private roots and atomically promote
verified executables. `binary` remains stubbed behind the interface.

### registry.Source

```
Load() (*Registry, error)
Sync(ctx context.Context) error
```

`FileSource` loads either an assembled local TOML catalog or a segmented
manifest. Remote sync fetches checksum-pinned same-origin segments with size,
path, redirect, and concurrency bounds, validates the complete catalog, then
atomically replaces the local assembled copy. `GitHubClient` lists tags and
resolves user-selected refs to immutable commits before a pull or explicit
GitHub tool install. Legacy single-file registry URLs remain supported.

The process root context is canceled on interrupt or termination and is passed
through CLI commands to installers, registry synchronization, GitHub queries,
and asset transfers so partial work can clean up before exit.

### store.Store

```
Load() (*Manifest, error)
Save(m *Manifest) error
```

`FileStore` persists the manifest as JSON. The interface is swappable.

## On-disk layout

All state lives under one root (`$ARSENAL_HOME`, else `$XDG_DATA_HOME/arsenal`,
else `~/.local/share/arsenal`):

```
root/
  config.json          user configuration
  registry.toml        active registry (seeded from the embedded copy)
  manifest.json        installed-tools manifest
  tools/<name>/<ver>/  isolated venv, Go, or Cargo installation root
  bin/                 generated PATH shims
  ops/<name>.toml      op definitions
  ops/<name>.lock.toml op lockfiles
  bundles/             exported offline bundles
  cache/               transient data
```

## Stubbed work (tracking)

The following are wired behind their interfaces with clear, intentional
not-implemented errors so misconfiguration fails loudly:

- **Container backend** (`internal/isolation/container`, TODO arsenal#1):
  provision tools with heavy system dependencies in podman/docker, mapping
  `InstallSpec` to a build, and `Run` to a `run --rm` invocation.
- **Binary install method** (`internal/installer/binary.go`, TODO arsenal#2):
  download a release asset and verify its checksum.
- **Offline bundling** (`internal/bundle`, TODO arsenal#5): vendor wheels and
  source archives alongside a lockfile so an air-gapped host can reconstruct an
  environment with no network.

## Trust boundaries

arsenal crosses four important boundaries:

1. CLI arguments, environment variables, configuration, manifests, and
   lockfiles enter from the local operator environment.
2. Registry sync, package managers, source repositories, and asset downloads
   return remote, untrusted content.
3. Installers start external runtimes and third-party tools with inherited host
   access unless an isolation backend constrains them.
4. Stores, shims, environments, assets, and lockfiles write to the local
   filesystem.

Validation belongs at the package that first crosses each boundary. Safe values
then move inward as typed data. Path containment, subprocess argument handling,
network limits, integrity checks, and atomic writes follow `RULES.md` and
`DESIGN.md`. Persisted JSON and TOML schemas reject unknown fields. Filesystem
components share one validator, and complete-file replacement uses unpredictable
sibling temporary files through the shared atomic writer.

Python virtualenv execution removes inherited `PYTHONHOME` and `PYTHONPATH`,
disables the user site, and prepends only the selected environment's bin
directory. Other environment values remain available to tools that need
operator-supplied configuration.

arsenal manages tools used for authorized testing, but it does not authorize or
scope their execution. The operator selects targets and remains responsible for
permission. The CLI must not infer targets or automatically execute a fetched
tool.
