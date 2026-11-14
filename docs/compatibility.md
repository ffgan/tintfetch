# Compatibility after the master merge

Tintfetch is a native Go system information application. The integration preserves origin/master
at de3e5aa, including its native Go HyFetch color profiles, complete preset/font
catalogs, interactive configuration, fastfetch-style module catalog, and platform
probes. Go 1.26+ is required. Those features are implemented in this repository;
a separately installed HyFetch program is not required.

Use `--help` for the actual CLI and `--list-presets` for the native preset catalog. Native HyFetch flags include `-p`/`--preset`,
`-c` for interactive setup, `-C` for its JSON configuration, color modes and
alignment. `hyfetch.json` retains its upstream namespace for compatibility.
Optional external rendering backends (`-b fastfetch`, `-b neofetch`, `-b macchina`)
require the corresponding programs. Image rendering may use chafa/catimg/viu.
Some crop, image-loop and image-size flags remain compatibility placeholders;
their presence is not proof of complete upstream image-protocol parity.

The native configuration namespace is now `tintfetch/config.conf`. User config
uses the OS user config directory, including XDG_CONFIG_HOME on Linux;
Neofetch paths remain fallback candidates. Explicit `--config` selects a
file and CLI options override it. `--no_config` disables info configuration loading.
HyFetch JSON configuration is handled independently. Existing underscore flags
remain supported; `--no-config`, `--print-config`, and `--plain` are aliases for
`--no_config`, `--print_config`, and `--stdout`.

Ordinary info configuration is parsed as data rather than sourced by a shell.
The remote parser's existing `info` and `prin` statements are retained.
`--shell-config --config PATH` explicitly executes a trusted file with Bash,
captures known settings and invokes its print_info function through an info
collector. It has a five-second execution limit. That opt-in mode requires an
explicit filename and is not a sandbox. Its print_info body must use info calls;
arbitrary upstream rendering functions are not emulated. The separate command
information module intentionally executes its configured command; public-IP
and weather modules intentionally make network requests when selected.

`--stdout`/`--plain` prints without logo or ANSI styling. `--json` preserves the
remote interface: an object with `version` and an ordered `result` array; rows
contain `key`, optional `label`, and `value`. The native JSON wire format is preserved.

Linux is the primary distribution-packaging target. Existing Linux/macOS/Windows
end-to-end tests, native platform tests, and full release-target matrix are
retained. Cross-compilation does not establish runtime correctness for every
upstream platform or module. See README.md for module and release usage, and
releasing.md for distribution review requirements. Newly inherited artwork that
still requires source review is listed in assets-pending-provenance.txt.
