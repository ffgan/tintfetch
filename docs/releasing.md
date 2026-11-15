# Release and distribution preparation

The current development version is 0.3.7-dev. The module and repository are
`github.com/ffgan/tintfetch`; executables and archives are named tintfetch.
A development snapshot can be reviewed by a distribution without a stable tag.
Preparing a candidate does not create a Fedora approval or a release tag.

## License and corresponding source

On 2026-11-15 the maintainer authorized observing GPL-2.0-or-later for the ARM
ID registry and its resulting binary-distribution obligations. The generated
registry combines util-linux GPL data and Fastfetch MIT additions. Its pinned
upstream inputs and generator are in the source tree. Other MIT source files
retain their license; the compiled Go runtime retains its BSD-3-Clause notice.
The aggregate RPM License expression is
`GPL-2.0-or-later AND MIT AND BSD-3-Clause`.

COPYING contains the full GPL text; LICENSE explains the MIT scope, and
THIRD_PARTY_NOTICES preserves source/runtime attribution. `make install`, Debian,
RPM and binary release archives install/include all three notice files.

The binary release script automatically creates and verifies complete
corresponding source on every invocation, including partial target builds.
The source archive uses a portable deterministic Python writer, and `make dist`
uses the same writer. Before publication CI independently recreates the source
from the exact checkout used for the builds:

```sh
make dist VERSION="$VERSION"
cp "build/tintfetch-$VERSION.tar.gz" "dist/tintfetch_${VERSION}_source.tar.gz"
python3 scripts/check-release-source.py dist "$VERSION"
```

The verifier compares the source files, generator inputs and build scripts with
the build checkout and checks COPYING, LICENSE, THIRD_PARTY_NOTICES and SOURCE.txt
in every binary archive. CI runs this before publication. The source archive is
published in the same release, included in checksums, and must remain available
with the binaries. A repository link alone is not a replacement for the GPL
corresponding-source obligation. Downstream distributors must observe the GPL
when redistributing the resulting binaries; Debian/Fedora source packages include
the exact upstream archive and packaging needed to rebuild them.

## Source and behavioral checks

```sh
make check
make build
sh scripts/smoke.sh
go test -race ./...
sh scripts/reproducible.sh
```

`make check-provenance` validates the artwork/JSON hashes and regenerates the
ARM table offline. Follow provenance.md for an independent comparison with the
pinned artwork references. The source distribution includes the unmodified ARM
inputs and the generator; all 209 earlier mappings retain their display names.
Go 1.26+ is required; Python 3 is a build/test dependency, not a runtime dependency.
No third-party Go module dependencies or toolchain/module downloads are required.
External detection/rendering commands remain optional and probes run as the user.
Bash configuration and command modules are explicit features; selected weather
and public-IP modules contact external services, using HTTPS defaults.

Source lineage, specific copyright notices, exact data references and the limits
of historical reconstruction are in provenance.md. Known AI usage and contributor
responsibilities are disclosed in CONTRIBUTING.md. Scanners and tests support
review; their results do not establish ownership or replace human responsibility.

## Fedora review candidate

The SPEC uses an immutable upstream commit for Source0 and the matching
`packaging/rpm/go-vendor-tools.toml` for Source1. Both must refer to a published
source revision containing the current notices and source changes. When making
a new source snapshot, update commit, snapdate, Release and the changelog together.
The SPEC builds with Fedora Go macros and system C flags, keeps PIE/RELRO/build
IDs and debuginfo, and checks tests, provenance, CLI behavior and license metadata.

On a Fedora build host, install rpm-build, rpmdevtools, go-rpm-macros,
go-vendor-tools, askalono-cli, make, gcc, glibc-devel, Python 3 and Go 1.26+.
Then download the fixed sources and build:

```sh
mkdir -p build/fedora-review/SOURCES
spectool -g -C build/fedora-review/SOURCES packaging/rpm/tintfetch.spec
rpmbuild -ba packaging/rpm/tintfetch.spec \
  --define "_topdir $PWD/build/fedora-review"
rpmlint packaging/rpm/tintfetch.spec \
  build/fedora-review/SRPMS/*.rpm build/fedora-review/RPMS/*/*.rpm
sudo mock -r fedora-rawhide-x86_64 \
  --rebuild build/fedora-review/SRPMS/*.src.rpm \
  --resultdir "$PWD/build/fedora-review/mock-rawhide"
```

There are no Go modules to vendor. go-vendor-tools scans the license files and
checks the aggregate expression. The config explicitly classifies COPYING's
later-version grant and the nonstandard THIRD_PARTY_NOTICES filename by hash;
these classifications are based on the supplied license texts and source grants.
It does not suppress missing/unknown licenses. The license files are installed
explicitly with `%license`, with no separate Go library packages.

Before filing a review, check existing packages, retirement history and review
requests, confirm the maintainer's Fedora account/group status, and obtain
packager sponsorship if needed. A Bugzilla search on 2026-11-15 returned no
Package Review ticket with tintfetch in its summary; that search alone does
not establish name availability, account membership or sponsorship.
Contact identity: ffgan <ffgandev@gmail.com>; FAS username supplied: ffgan.

Upload a publicly accessible SPEC and SRPM and file a Fedora Package Review
request. Identify upstream maintainership, disclose the AI-assisted packaging and
review text, and include real mock/rpmlint results for that exact source version.
A human maintainer must review and authorize the actual submission. No historical
human sign-offs, sponsorship or Fedora acceptance are fabricated by these tools.

## Debian

The debian/ recipe targets build roots with Go 1.26+. Run
`dpkg-buildpackage -us -uc -b`, Lintian and the debian/tests installed-package
smoke test; use sbuild/pbuilder before seeking sponsorship. A source upload needs
an appropriately named original tarball. dwz is skipped for Go compressed DWARF;
normal stripping and debug-package generation remain enabled. Examples stay in
documentation, and no maintainer scripts execute user configuration.

## Policy references

Inspected on 2026-11-15:

- Fedora AI policy v1.0 (2025-10-24), council/docs at
  8eeb42f8007c49118856d78c38b9be4cef721108:
  https://forge.fedoraproject.org/council/docs/src/commit/8eeb42f8007c49118856d78c38b9be4cef721108/council/modules/ROOT/pages/policy/ai-contribution-policy.adoc
- Fedora packaging guidelines at 7f8ef2ff5494acff0c316ea036741443eb330a84,
  including Golang, SourceURL, LicensingGuidelines and ReviewGuidelines:
  https://forge.fedoraproject.org/packaging/guidelines/src/commit/7f8ef2ff5494acff0c316ea036741443eb330a84/guidelines/modules/ROOT/pages/Golang.adoc
- New contributor process at dd2a0b21f9dbca82ee68a8850971a367f3c3c2a9:
  https://pagure.io/fedora-docs/package-maintainer-docs/blob/dd2a0b21f9dbca82ee68a8850971a367f3c3c2a9/f/modules/ROOT/pages/New_Package_Process_for_New_Contributors.adoc
- Debian source policy:
  https://www.debian.org/doc/debian-policy/ch-source.html

Fedora permits AI assistance and requires disclosure of significant unchanged
AI output, with other material assistance recommended for disclosure. The human
contributor is accountable for quality, licensing and usefulness. AI must not be
the sole/final subjective arbiter. The policy does not require reconstructing
all historical prompts/models, rewriting Git history, or banning AI-assisted
upstream projects. This AI-assisted preparation is evidence for human review,
not the final subjective decision to accept a Fedora contribution.
