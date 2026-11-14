# Tintfetch

Tintfetch (formerly **gf**) displays system information and colorful distribution
artwork. The built-in Go engine builds as one executable without third-party Go
modules. Optional adapters run installed Neofetch/HyFetch programs for their
additional options, image backends, interactive configuration and platforms.

**Development status:** the native engine is not a complete Go port of HyFetch
or Neofetch. DEB/RPM packaging is included for validation, not evidence of
acceptance into any distribution. See [compatibility](docs/compatibility.md) and
[release preparation](docs/releasing.md).

## Build and run

Requires Go 1.24 or newer. GNU Make is optional for building the executable and
required by the supplied Linux packaging. Builds do not download dependencies
or automatically download a newer Go toolchain.

```sh
git clone https://github.com/ffgan/gf.git
cd gf
make check
make build
./build/tintfetch --no-config
./build/tintfetch --no-config --json
./build/tintfetch --no-config --preset trans
```

The product, executable, module and configuration namespace are now `tintfetch`.
The GitHub repository remains at `ffgan/gf` until its owner renames it; the Go
module path `github.com/ffgan/tintfetch` is reserved for that intended destination.
Build from this checkout until that migration is completed. No remote repository
has been renamed by these changes.

## Install

```sh
make install PREFIX=/usr/local
# For distribution staging, without writing to the host's /usr:
make install PREFIX=/usr DESTDIR="$PWD/build/stage"
```

Installation includes `tintfetch`, its manual page, MIT license and third-party
notices, documentation, and an example configuration. It does not overwrite a
user's configuration. `make uninstall` removes the explicitly installed files.
Packagers can override `BINDIR`, `MANDIR`, `DOCDIR`, and `LICENSEDIR` independently.

## Configuration

```sh
mkdir -p "$HOME/.config/tintfetch"
./build/tintfetch --print-config > "$HOME/.config/tintfetch/config.conf"
./build/tintfetch --config ./my-config.conf --plain
```

Search order is an explicit `--config` path, otherwise the OS user configuration
directory's `tintfetch/config.conf`, otherwise `/etc/tintfetch/config.conf`.
On Linux, `XDG_CONFIG_HOME` is honored and otherwise `~/.config` is used.
Defaults apply when no configuration exists; explicitly supplied command options
override file values. `--no-config` ignores configuration files. No files are
created on ordinary invocations. Neofetch configuration is not loaded implicitly.

The default parser accepts quoted assignments, single-line arrays, comments,
and `print_info() { ... }` containing `info [label] module` statements. This
resembles the historical configuration syntax but **does not execute Bash**.
Unknown keys and unsupported syntax are errors. To intentionally execute a
trusted Bash configuration, use both `--config PATH` and `--shell-config`.
The latter requires Bash and executes code with your user permissions. Its
`print_info` body must emit only `info` calls; use an upstream engine for arbitrary
Neofetch rendering functions. See [compatibility](docs/compatibility.md).

## Output and upstream engines

```sh
./build/tintfetch --no-config --plain
./build/tintfetch --no-config --json
./build/tintfetch --list-presets
# Requires the corresponding separately installed program:
./build/tintfetch --engine neofetch --help
./build/tintfetch --engine hyfetch --help
```

Native redirected output is plain text. `--stdout=off` forces artwork/color;
`--plain`, `--stdout=on`, `TERM=dumb` in automatic mode, and nonempty `NO_COLOR`
disable it. JSON contains `version` and an ordered `entries` list of objects
with `label`, `module`, and `value` strings, without artwork or ANSI styling.
Public-IP lookup is opt-in by adding the `public_ip` information module.
Ordinary native startup does not make network requests.

Upstream engines receive their arguments and streams directly. They use their
own configuration and dependencies; native flags do not configure them. The
Neofetch adapter prefers `neowofetch` over `neofetch`. Missing engines produce a
clear error and are never installed automatically.

## Contributing and licensing

See [CONTRIBUTING.md](CONTRIBUTING.md), [SECURITY.md](SECURITY.md), and
[THIRD_PARTY_NOTICES](THIRD_PARTY_NOTICES). The project is MIT licensed and
includes Neofetch-derived code and artwork. All 384 embedded artwork files were
matched against a pinned upstream script; see [provenance](docs/provenance.md).

中文：Tintfetch 是原 gf 项目的新名称，用于在终端显示系统信息与发行版图案。
原生 Go 引擎尚未完整移植 HyFetch/Neofetch；需要上游完整功能时，可显式选择已安装的上游引擎。
