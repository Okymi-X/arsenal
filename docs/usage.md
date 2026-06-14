# Usage

Output is quiet by default. Add `-v` or `--verbose` before the command for
detail. There is no color unless stdout is a TTY, and no emoji anywhere. Status
markers are plain ASCII: `[ok]`, `[fail]`, `[warn]`, `->`.

## Global form

```
arsenal [-v|--verbose] <command> [args]
```

## Commands

### install

```
arsenal install <tool>[@version]
```

Resolves the tool against the registry, creates an isolated environment,
installs the pinned version, writes shims, and marks it active. Without a
version, the newest tested version is selected. A warning is printed if the
selected version is not flagged tested.

### remove

```
arsenal remove <tool>
```

Removes every installed version of the tool, its environments, and its shims,
then prunes the manifest.

### switch

```
arsenal switch <tool> <version>
```

Repoints the tool's shims to an already-installed version and marks it active.

### list

```
arsenal list
```

Lists installed tools and versions. `[*]` marks the active version of each tool.

### search

```
arsenal search <query>
```

Lists registry tools whose name, description, tags, or aliases match the query.
An empty query lists every tool.

### info

```
arsenal info <tool>
```

Prints the registry details for a tool, including each catalogued version.
`[ok]` marks tested versions.

### run

```
arsenal run <tool> -- <args...>
```

Runs the active version of the tool inside its isolated environment, forwarding
everything after `--`.

### op

```
arsenal op create <name> [description]
arsenal op pin <name> <tool>[@version]
arsenal op list
arsenal op export <name>
arsenal op use <name>
arsenal op import <lockfile>
```

Manages engagement profiles. `create` makes an empty op. `pin` resolves and
records a tool/version. `export` writes a fully pinned TOML lockfile. `use`
resolves the op and installs every pinned tool. `import` installs from a shared
lockfile.

### sync

```
arsenal sync
```

Downloads the registry from the configured upstream URL, validates it, and
replaces the local copy atomically.

### doctor

```
arsenal doctor [--fix]
```

Runs health checks: directory tree, Python interpreter, shim directory on PATH,
and manifest integrity. With `--fix`, repairs what it safely can (recreating
directories, pruning manifest entries whose environments are gone).

### bundle

```
arsenal bundle --offline [op]
```

Exports a self-contained offline bundle. Stubbed; see docs/architecture.md.

### version

```
arsenal version
```

Prints the build version, injected at link time.

## Environment variables

- `ARSENAL_HOME` - override the state root directory.
- `XDG_DATA_HOME` - used to derive the default root when `ARSENAL_HOME` is unset.

## Configuration

`config.json` under the root holds:

```json
{
  "registry_url": "https://raw.githubusercontent.com/Okymi-X/arsenal/main/registry/registry.toml",
  "default_backend": "venv",
  "python_bin": "python3"
}
```
