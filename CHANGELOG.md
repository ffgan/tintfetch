# Changelog

## 0.3.7 — Unreleased

- Rename gf to Tintfetch, including binary/archive names, module path and native
  configuration namespace; retain legacy gf/Neofetch config fallbacks.
- Preserve native HyFetch coloring/configuration, fastfetch-style modules, platform
  fixes, full release matrix and end-to-end tests merged through de3e5aa (gf v0.3.6).
- Add explicit trusted Bash configuration and clearer plain/config CLI aliases.
- Replace go-radix with an in-tree prefix registry; retain offline Go builds.
- Bound subprocess pipe waits and replace predictable macOS temporary files.
- Restore upstream/runtime license notices; verify 384 original artwork files and
  record 17 inherited artwork files requiring further provenance review.
- Add staged installation, manual, contribution/security docs and DEB/RPM recipes
  with package CI, PIE/RELRO and debug metadata.

Distribution acceptance and complete upstream image-protocol parity remain
future milestones. No release tag is created by this preparation work.
