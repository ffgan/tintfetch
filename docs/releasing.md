# Release and distribution preparation

Current changes are development work. Do not tag or advertise a stable release
as full upstream parity. The repository URL remains `ffgan/gf`; coordinate its
rename to `ffgan/tintfetch` before publishing the new Go module path. Check
repository, package and executable name availability with Debian/Fedora
maintainers before submission. No trademark or package-name clearance is implied.

Before publishing:

1. Resolve the remaining compatibility milestones or explicitly scope the release
   as a native subset with optional upstream adapters; update documentation accordingly.
2. Review license/provenance, `THIRD_PARTY_NOTICES` and `debian/copyright`.
   Re-run `scripts/check-provenance.py` against the pinned source when changing assets.
3. Run `make check`, `go test -race ./...`, smoke and reproducibility checks.
   Run the CI matrix and inspect actual results, not just workflow presence.
4. Choose a release version; update Debian/RPM versions, changelog and the binary
   default together. The supplied recipes use development versions intentionally.
5. Build source archives with `make dist VERSION=... SOURCE_DATE_EPOCH=...` using
   the release commit timestamp. Build from a clean, reviewed checkout; the archive
   includes the working tree, so do not keep secrets or unrelated files there.
6. Verify install/uninstall using DESTDIR; compare the source archive and binary
   licenses/docs. Publish checksums and release notes through the normal maintainer
   review process. Publication is not automatic in this repository.

## Debian and Ubuntu

The root `debian/` directory is a development source-package recipe. Build with
`dpkg-buildpackage -us -uc -b` in a clean Debian/Ubuntu chroot containing Go 1.24+
and debhelper 13+. For source packages, first create the corresponding original
source tarball with Debian's version naming. Use `sbuild`/`pbuilder`, inspect
`lintian`, install/remove the package in a disposable environment and then seek
sponsorship/review. Older Ubuntu releases without a suitable Go compiler are
not supported by the current recipe. Set the correct maintainer/contact data
before submission; the current noreply address is a repository identity only.

## Fedora

Use the current Fedora Go packaging guidelines and package review process; the
spec is an initial standalone application recipe, not an approved Fedora package.
The Go module has no external dependencies to unbundle. In a clean Fedora build
root with Go 1.24+, make and RPM development tools:

```sh
make dist VERSION=0.1.0
rpmbuild -ba packaging/rpm/tintfetch.spec --define "_sourcedir $PWD/build"
```

`Source0` describes the intended release-tag archive. For development, the command
above supplies a local archive; that release tag has not been created. Use `mock`
and `rpmlint`, review generated dependencies and debuginfo, and validate optional
upstream engine availability. Update architecture exclusions only from concrete
build failures. Do not assume local RPM build success implies policy acceptance.

Both recipes build PIE executables with RELRO using the system C toolchain,
retain build IDs/debug information, and stage into `/usr`, ship notices and a manual, avoid downloading
modules/toolchains during build, and install examples under documentation rather
than creating mutable user configuration. CI validates Debian unstable and current
Fedora; it does not establish support for every Debian/Ubuntu/Fedora release.
