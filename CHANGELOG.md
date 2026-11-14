# Changelog

## Unreleased

- Rename gf to Tintfetch; establish separate command and configuration namespaces.
- Add native declarative configuration parsing, explicit trusted Bash configuration,
  configuration precedence, JSON/plain output and six color presets.
- Replace go-radix with an in-tree prefix lookup; support offline builds.
- Bound optional subprocess execution, skip unavailable CPU/memory probes and
  replace predictable macOS temporary filenames.
- Add explicit adapters for separately installed HyFetch/Neofetch engines.
- Restore upstream MIT attribution and verify the 384 embedded artwork files.
- Add staged installation, man page, development DEB/RPM recipes, CI and regression tests.

Full native upstream parity and distribution acceptance remain future milestones.
