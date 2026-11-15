# Changelog

## 0.3.7 — Unreleased

- Standardize Tintfetch naming across binaries, archives, module path, repository,
  tests and native configuration namespace; retain Neofetch config fallbacks.
- Preserve native HyFetch coloring/configuration, fastfetch-style modules, platform
  fixes, full release matrix and end-to-end tests merged through de3e5aa (v0.3.6).
- Add explicit trusted Bash configuration and clearer plain/config CLI aliases.
- Replace go-radix with an in-tree prefix registry; retain offline Go builds.
- Bound subprocess pipe waits and replace predictable macOS temporary files.
- Restore upstream/runtime license notices and verify all 401 artwork files
  against pinned upstream references.
- Add staged installation, manual, contribution/security docs and DEB/RPM recipes
  with package CI, PIE/RELRO and debug metadata.

Distribution acceptance and complete upstream image-protocol parity remain
future milestones. No release tag is created by this preparation work.

### Distribution preparation

- Complete pinned-source matching for all 401 logos and both HyFetch JSON files.
- Preserve the MIT attribution for Björn Ottosson's OKHSL conversion.
- Enforce offline provenance checks during checks/source archive creation; compare
  pinned upstream references in CI.
- Add installed Debian smoke tests, maintainer contact, RPM system cgo flags,
  and install the audit manifests with documentation.

### Fedora review and license remediation

- Regenerate the ARM ID registry from pinned util-linux GPL-2.0-or-later and
  Fastfetch MIT source inputs; retain all 209 earlier mappings.
- Include GPL terms, complete generator inputs and corresponding-source release
  gates while preserving MIT source and BSD runtime notices.
- Disclose substantial AI-assisted development and require human responsibility.
- Use Fedora Go macros, a pinned source snapshot and go-vendor-tools license
  checks; fix the RPM description spelling error and check source/binary RPMs.
- Use HTTPS for the default optional public-IP endpoint.
