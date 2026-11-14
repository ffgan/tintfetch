# Release and distribution preparation

This merge retains the native features and release history through gf v0.3.6 and
starts Tintfetch development at 0.3.7-dev. The GitHub remote remains `ffgan/gf`;
coordinate its rename before publishing the `github.com/ffgan/tintfetch` Go module
path. Existing release binaries are still named gf; newly generated archives use
tintfetch. No release tag or distribution submission is created by this merge.

Before publication:

1. Review pending artwork provenance in assets-pending-provenance.txt and the
   native HyFetch-derived data/code notices. Preserve LICENSE and
   THIRD_PARTY_NOTICES in all archives and binary packages.
2. Run `make check`, `go test -race ./...`, scripts/smoke.sh and
   scripts/reproducible.sh. Inspect actual cross-platform CI/e2e results.
3. Verify package/executable name availability and maintainer contact information.
   The current noreply maintainer address is a repository identity only.
4. Update the binary, release-script, DEB and RPM versions together. Build source
   archives from a clean checkout using `make dist VERSION=...` and the release
   timestamp as SOURCE_DATE_EPOCH. The archive includes the working tree.
5. Validate installation/removal, metadata, notices, debuginfo and checksums.
   Publish through maintainer review; presence of packaging is not archive approval.

Go 1.26+ is required by the integrated upstream implementation. Debian/Ubuntu
build roots without that compiler are unsupported. Root debian/ contains a
source-package recipe: run `dpkg-buildpackage -us -uc -b`, inspect Lintian, and
use sbuild/pbuilder before seeking sponsorship. For source packages, provide an
appropriately named original tarball. dwz is skipped for Go's compressed DWARF;
normal stripping and debug-package generation remain enabled.

For Fedora, use a clean Go 1.26+ build root with make, gcc, glibc-devel and RPM
build tools:

```sh
make dist VERSION=0.3.7
rpmbuild -ba packaging/rpm/tintfetch.spec --define "_sourcedir $PWD/build"
```

The spec is an initial application recipe, not an approved Fedora package.
Source0 refers to the intended future release tag; local development supplies the
archive explicitly. Run mock/rpmlint and complete Fedora package review.

Both recipes build PIE/RELRO executables with system C linking, retain build IDs
and debug information, stage into /usr, install a manual and copyright notices,
and keep examples in documentation. There are no third-party Go modules and
no module/toolchain downloads during builds. The retained release script also
supports CGO-disabled builds where the target allows them. Current packaging CI
uses Debian unstable and Fedora; it does not promise support for all older releases.
