# Architecture

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

## Key interfaces

These are defined first and implemented against. Concrete types are injected.

### isolation.Backend

Isolates a single tool/version environment.

```
Create(tool, version string) error
Install(spec InstallSpec) error
Run(args []string) error
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
Install(ctx context.Context, tool registry.Tool, version registry.Version) error
```

Implementations: `pip` and `gitpip` (implemented, drive the backend), plus
`binary`, `gobin`, and `cargo` (stubbed behind the interface).

### registry.Source

```
Load() (*Registry, error)
Sync() error
```

`FileSource` loads a local TOML file and refreshes it from a remote URL.

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
  tools/<name>/<ver>/  isolated environments (venvs)
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
- **Go install method** (`internal/installer/gobin.go`, TODO arsenal#3):
  `go install` into a per-tool GOBIN.
- **Cargo install method** (`internal/installer/cargo.go`, TODO arsenal#4):
  `cargo install` into a per-tool root.
- **Offline bundling** (`internal/bundle`, TODO arsenal#5): vendor wheels and
  source archives alongside a lockfile so an air-gapped host can reconstruct an
  environment with no network.
