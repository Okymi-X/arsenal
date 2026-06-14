# arsenal

A package and environment manager specialized for offensive-security tooling.
It is a domain-aware alternative to `pipx`/`uv`.

The point of arsenal is not the packaging mechanics. It is the curated registry:
a hand-maintained, tested mapping from tool name to known-good versions
(NetExec/nxc, Impacket, Certipy, and more), pinned and annotated so you never
have to open a repo mid-engagement to find out which version actually works.

arsenal ships as a single statically linked Go binary with no runtime
dependencies on the host. It orchestrates the tools already on the system
(`python -m venv`, `pip`, and optionally a container runtime); it does not
reimplement them or bundle a Python interpreter.

## Features

- Curated TOML registry of tested versions, pinned by commit, with the required
  Python version and operational notes.
- Per tool/version isolation using Python virtualenvs (default), with a
  container backend planned behind the same interface.
- Engagement profiles ("ops"): pin a set of tool versions, produce a lockfile,
  and make an environment reproducible and shareable across a team.
- PATH shims so multiple versions coexist and the active one is switchable.
- A `doctor` command that reports and repairs broken installs.
- Works offline out of the box: the curated registry is embedded in the binary.

## Install

Build from source (Go latest stable required):

```
make build
sudo make install
```

Then add the shim directory to your PATH (printed by `arsenal doctor`):

```
export PATH="$HOME/.local/share/arsenal/bin:$PATH"
```

## Quick start

```
arsenal search ad              # browse the registry by keyword
arsenal info nxc               # see tested versions of NetExec
arsenal install nxc            # install the newest tested version
arsenal install impacket@0.12.0
arsenal list                   # show installed tools; [*] marks active
arsenal run nxc -- smb 10.0.0.1
arsenal switch nxc 1.3.0       # repoint shims to another installed version
arsenal remove nxc
```

Output is quiet by default. Add `-v`/`--verbose` for detail. There is no color
unless stdout is a TTY, and no emoji anywhere.

## The registry concept

Each tool entry records its repo, category, install method, required Python
version, exposed binaries, and a list of versions. Each version carries a tag, a
pinned commit, a `tested` flag, a pip spec, a date, and notes. The newest tested
version is selected by default; you can always pin an explicit one with
`tool@version`.

Refresh the registry from upstream at any time:

```
arsenal sync
```

The full schema is documented in [docs/registry-format.md](docs/registry-format.md).

## The op workflow

An "op" is an engagement profile: a named set of pinned tool versions that you
can reproduce and share.

```
arsenal op create redteam-q3 "Q3 internal"
arsenal op pin redteam-q3 nxc            # pins the newest tested version
arsenal op pin redteam-q3 impacket@0.12.0
arsenal op export redteam-q3             # writes a TOML lockfile
arsenal op use redteam-q3                # installs everything in the lockfile
```

Hand the lockfile to a teammate and they reproduce the exact set:

```
arsenal op import redteam-q3.lock.toml
```

## Documentation

- [docs/architecture.md](docs/architecture.md) - package layout and interfaces
- [docs/registry-format.md](docs/registry-format.md) - the registry schema
- [docs/usage.md](docs/usage.md) - full command reference

## Status

The venv backend, the pip and git+pip install methods, the registry, the shim
system, the op manager with lockfiles, and the core commands are implemented and
tested end to end. The container backend, the offline bundle, and the binary/go/
cargo install methods are stubbed behind their interfaces with tracking notes.

## License

MIT. See [LICENSE](LICENSE).
