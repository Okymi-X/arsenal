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
arsenal run <tool> <binary> -- <args...>
```

Runs the active version of the tool inside its isolated environment, forwarding
everything after `--`. With no binary selector, the tool's primary binary runs.

For multi-binary tools (such as Impacket) name the binary as the first argument.
The selector is matched loosely, ignoring case, a `.py` suffix, and a `<tool>-`
prefix, so these are equivalent:

```
arsenal run impacket getTGT -- -dc-ip 10.0.0.1 domain/user
arsenal run impacket getTGT.py -- -dc-ip 10.0.0.1 domain/user
arsenal run impacket impacket-getTGT -- -dc-ip 10.0.0.1 domain/user
```

Alternatively, add the shim directory to your PATH (see `doctor`) and call the
binaries directly: `getTGT.py`, `secretsdump.py`, `nxc`, and so on.

### fetch

```
arsenal fetch <asset> [binary] [--dest <dir>] [--build <build>] [--list]
```

Pulls the latest version of a precompiled binary and writes it into a
directory, for staging onto a target. Unlike `install`, `fetch` creates no
isolated environment and no shim: it just downloads the file.

`<asset>` is an entry from the registry's asset catalog (`SharpCollection`,
`winpeas`, `linpeas`, `pspy`, ...); `arsenal search` lists assets with an
`(asset)` marker. Without `--dest`, the file lands in the current directory.
Downloaded files are made executable.

For a single-binary asset the default file is used unless you name another as
`[binary]`:

```
arsenal fetch linpeas --dest ./www          # -> ./www/linpeas.sh
arsenal fetch winpeas --dest ./www          # -> ./www/winPEASx64.exe
arsenal fetch winpeas winPEASany.exe        # pick a different release asset
arsenal fetch pspy pspy32                    # pick a variant
```

For a collection such as SharpCollection, name the binary and optionally pick a
build with `--build` (default `NetFramework_4.7_x64`); `--list` shows what is
available:

```
arsenal fetch sharpcollection --list
arsenal fetch sharpcollection Rubeus --dest ./www
arsenal fetch sharpcollection Certify --build NetFramework_4.5_x64
```

The same model serves payload lists from PayloadsAllTheThings; fetch a list and
feed it to a fuzzer:

```
arsenal fetch payloads-lfi --list
arsenal fetch payloads-lfi JHADDIX_LFI.txt --dest ./wl
ffuf -u 'http://target/?page=FUZZ' -w ./wl/JHADDIX_LFI.txt
```

The binary selector is matched loosely, ignoring case and a `.exe`/`.sh`
suffix, so `rubeus` resolves `Rubeus.exe`. A `GITHUB_TOKEN` or `GH_TOKEN` in the
environment raises GitHub API rate limits for `--list` and release lookups.

### completion

```
arsenal completion bash > /etc/bash_completion.d/arsenal
arsenal completion zsh  > "${fpath[1]}/_arsenal"
arsenal completion fish > ~/.config/fish/completions/arsenal.fish
```

Prints a shell completion script for bash, zsh, or fish. Completion covers the
subcommands and, dynamically, registry tool names (for `install`/`info`) and
installed tool names (for `run`/`remove`/`switch`). arsenal must be on PATH for
the dynamic tool-name completion to work.

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
