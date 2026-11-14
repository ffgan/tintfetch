# Compatibility contract

## Engine choice

`native` is the default and runs the built-in Go implementation. `--engine hyfetch`
runs a separately installed `hyfetch`; `--engine neofetch` prefers `neowofetch`,
then uses `neofetch`. Arguments other than the engine selector, standard input,
standard output, standard error and positive exit codes are forwarded. These
adapters provide access to upstream functionality, **not a native port**. Features
and platform support depend on the installed upstream version and prerequisites.
No upstream binaries, Python/Rust runtime or full upstream scripts are bundled.

For the original project's full options (music players, advanced image backends,
looped rendering and interactive HyFetch setup), select the corresponding
installed upstream engine. That program uses its own configuration locations;
`--shell-config`, `--plain`, and native JSON/preset semantics do not apply there.

## Native support

Core modules include title, underline, distro, OS, model, kernel, uptime, packages,
shell, editor, CPU, GPU, memory, resolution, desktop/window-manager and theme
information, fonts, icons, cursor, terminal, network, Bluetooth and BIOS.
Probe availability and accuracy depend on the OS and optional commands.
Additional modules include disk, battery, power adapter, local/public IP, users,
locale, Java/Python/Node versions and GPU driver. Battery, power adapter and GPU
driver probes currently use Linux sysfs; disk uses GNU `df`; these additions
are not portable implementations. `cols` is a recognized layout marker and emits
no native color-block row yet.

Working native settings: `title_fqdn`, `os_arch`, `distro_shorthand`,
`kernel_shorthand`, `uptime_shorthand`, shell/editor path/version switches,
`cpu_cores`, `cpu_speed`, `cpu_temp`, `cpu_brand`, `speed_type`, `speed_shorthand`,
package manager naming/separation/minimal switches, memory
unit/precision/percent/display, bar length/characters, separator, underline
switch/character, image source (`auto`, `ascii`, file or literal ASCII), image
backend (`ascii`, `off`), ASCII distro/colors/bold, gap, stdout mode, disk
selection/subtitle/percent and local/public-IP settings. Public-IP requests use
a timeout of 1–60 seconds and require an IP-only response.

Flags retain historical underscore spelling. Multiword/array arguments must be
quoted as a single shell argument. `--print-config` and `--no-config` are the
preferred hyphenated forms. Six built-in `--preset` stripe palettes are available;
these are not the full upstream HyFetch preset catalog or interactive setup.

Other historical registered flags are retained with an explicit unsupported
error in the native engine. Changing unsupported configuration settings from
their defaults also raises an error. They are not silently treated as implemented.
Use `--help` for the complete registered option list and an upstream engine for
its own help. In particular, GPU selection, refresh-rate
control, GTK toggles, music, most text/color-block/bar styling controls and image
processing options are not implemented natively.

## Configuration subset and precedence

Native precedence: defaults, selected config, explicit flags. User configuration
replaces the system fallback; they are not merged. Missing default config files
are allowed; an explicitly selected missing file is an error. The metadata
commands (`--help`, `--version`, `--print-config`, `--list-presets`) do not read
configuration or probe the machine.

Assignments and arrays are single-line. Quotes, escaped characters and trailing
comments are supported. Only `info [label] module` is accepted inside one optional
`print_info()` block. Unknown keys, command substitution, other functions and
shell statements require migration or the explicit Bash mode. Unknown/unsupported
modules produce an error rather than reflection panics.

`--shell-config --config PATH` sources a trusted file with Bash startup hooks
disabled, captures known setting values, then calls its `print_info` function
with a replacement `info` collector. It is limited to five seconds. Arbitrary
`prin` functions, dynamically added fields and upstream internal globals are
not emulated; select an upstream engine for full executable-config compatibility.

## Platform commitments

Linux is the primary runtime target. CI checks compilation, argument/configuration
handling and JSON/metadata behavior on Linux, macOS and Windows. Go
cross-compilation alone does not establish probe correctness. Native runtime
parity on BSD, macOS, Windows, Android, Haiku and other upstream platforms is
not established. Optional desktop utilities are not hard runtime dependencies.

## Remaining milestones

- Complete the full native flag and info-module implementation with fixture tests.
- Port and verify upstream image protocols, crop/resize/cache/loop behavior.
- Port HyFetch's full preset catalog, interactive configuration and color modes.
- Establish a pinned behavioral parity suite and real platform runtime coverage.
- Have distribution maintainers review licensing, names, packaging and policy.

These milestones are required before describing Tintfetch as a full pure-Go
replacement for both upstream projects. Adapter availability is not completion
of those milestones.
