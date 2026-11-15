# Source and license provenance

The maintainer identified `ffgan/hyfetch` as the migration source. The references
below were inspected on 2026-11-15. They establish matching source text, not the
exact revisions used for the historical migration or independent authorship of
each distribution's artwork. Distribution names and marks belong to their
respective owners; no endorsement is implied.

## Embedded artwork and data

`provenance.json` records pinned URLs, upstream SHA-256 hashes, local hashes and
comparison methods for every file in `assets/` and `internal/hyfetch/data/`.
`assets.sha256` is the corresponding artwork checksum list.

- 384 original logos match `ffgan/hyfetch`'s `neofetch` at
  `2252fc36e76eebc42030c383784f63f9489f5630` (script version 8.0.4).
- 16 additional logos match `hykilpikonna/hyfetch`'s `neofetch` at
  `a54daf37166d7ec69a4315a005fe38071caa3a6c`.
- `AOSC_Retro` matches `ffgan/hyfetch`'s `neofetch` at
  `ace0245097b3c9736975ada66cb45bb684011f01`, before the Afterglow rename.
- `font_logos.json` is byte-identical to the original reference's
  `hyfetch/data/font_logos.json`.
- Every entry in `presets.json` matches the newer reference's
  `hyfetch/data/presets.json`. The local file omits `demiromantic`; formatting
  differs. Embedded source comments are preserved.

Artwork comparison removes only surrounding newlines, preserving spaces,
color markers and all other bytes. All three reference repositories/revisions
carry Dylan Araps's MIT notice (2015–2023) in `LICENSE.md`. The script also
contains his earlier notice (2015–2021). LICENSE and THIRD_PARTY_NOTICES retain
upstream attribution. `cmd/README` and `configs/config` retain their legacy
notices. No artwork remains on the former pending-provenance list.

To check the recorded local bytes without network access:

```sh
make check-provenance
```

For an independent upstream comparison, download the pinned files under their
reference IDs, then run:

```sh
mkdir -p build/upstream
python3 scripts/check-provenance.py --fetch | while read -r name url; do
    curl --fail --location "$url" -o "build/upstream/$name" || exit 1
done
python3 scripts/check-provenance.py --upstream-dir build/upstream
```

Builds never download these references. CI runs both checks. A changed or new
embedded file requires an explicit provenance manifest update and source review.
The local checksum check alone does not independently prove upstream matching.

## Translated implementation and color conversion

The Go implementation includes manually translated and modified Neofetch and
HyFetch behavior. This audit establishes lineage and data matches; it does not
reconstruct every translated function's historical revision. Preserve upstream
notices and identify additional sources in new contributions.

`internal/hyfetch/okhsl.go` adapts Björn Ottosson's MIT-licensed Oklab/OKHSL
conversion, with Copyright (c) 2021 Björn Ottosson. Reference:
https://github.com/bottosson/bottosson.github.io/blob/7561fbab5c8b982020ed212aebb0b8620c44b228/misc/ok_color.h
The header includes the MIT permission and warranty notice. This is a Go
adaptation, not a byte-identical copy. Its notice is retained in the source,
THIRD_PARTY_NOTICES and debian/copyright.

Compiled executables contain the Go runtime and standard library. Their BSD
redistribution notice is included in THIRD_PARTY_NOTICES and installed with
binary packages. The combined executable is GPL-2.0-or-later; MIT source
portions retain their original license. There are no
third-party Go modules. A distribution reviewer may still request individual
artwork attribution or further translation history before acceptance.

## ARM registry: GPL lineage resolved by observing its terms

The maintainer chose to observe GPL-2.0-or-later for the ARM registry on
2026-11-15. The registry is now generated from two checked-in, unmodified inputs:

- `internal/CLI/data/lscpu-arm.c`, util-linux revision
  `d6abbbab2cf284597720b8426fdd67c5cef7c9f8`, GPL-2.0-or-later,
  Copyright (C) 2018 Riku Voipio.
- `internal/CLI/data/fastfetch_cpu_arm.h`, Fastfetch revision
  `73368d36d99c5e647782346b7c075aac12599e45`, with the upstream MIT
  attribution for Linus Dierheimer and Carter Li. This file identifies the
  util-linux lineage; its additions and spellings retain the MIT terms.

`scripts/generate-arm-ids.py` checks both input hashes and emits
`internal/CLI/armids_data.go` with `GPL-2.0-or-later AND MIT` metadata.
It starts from util-linux's tables and adds Fastfetch's entries and display
spellings. All 209 previously supported implementer/part mappings were retained;
the resulting registry contains 212 mappings. No util-linux or Fastfetch C/C++
library is built or linked; the original files are retained as generator source.
`make check-provenance` reproduces and compares the generated file offline.
`internal/CLI/data/README` records full URLs and notices.

COPYING contains the complete current FSF GPL version 2 text from
https://www.gnu.org/licenses/old-licenses/gpl-2.0.txt (including its current
contact URL rather than an obsolete postal address). LICENSE defines the scope of
the MIT portions; THIRD_PARTY_NOTICES retains all upstream notices. Binaries
incorporating the registry are distributed under GPL-2.0-or-later. The release
workflow publishes complete corresponding source and refuses binary archives
with absent/mismatched notices or source. Source distributions include the
original generator inputs and the generator, as well as the Go sources.

## Implementation-source review scope

The 403-file artwork/JSON inventory and its pinned-reference comparison cover
embedded data. They do not establish legal ownership of every translated line.
The implementation has these known source/attribution groups:

- The original Go application, CLI, configuration parser, utility wrappers,
  platform probes, tests and build scripts were declared MIT by this project.
  Its initial implementation is in `8716e49`, with the MIT grant added in
  `1754779`. Later changes retain the project license unless specifically marked.
- Neofetch-compatible probing, formatting, distribution matching, configuration
  text and artwork follow the maintainer-identified Neofetch/HyFetch lineage.
  Dylan Araps's MIT notice is retained. `36603f7` is the large neofetch-style
  completion change; this is historical evidence, not a claim that it was the
  first use of any specific upstream function.
- `internal/hyfetch` was introduced in `69a5973` as native Go compatibility
  code and data for HyFetch. The pinned HyFetch references above supply the
  data and behavior; MIT notices are retained. The Ottosson conversion is
  separately attributed above.
- Native module implementations under `internal/modules`, hardware helpers and
  newer platform probes are project Go code following local OS interfaces and
  Fastfetch's module behavior. Major introduction commits include `e41f664`,
  `99e1d07` and `cc106c5`. Fastfetch's MIT notice is retained conservatively for
  reference-derived material, without claiming a particular line-by-line C port.
- The generated ARM registry is the explicitly identified GPL/MIT exception.
  Its exact input bytes and reproduction procedure are now available.

The source license scan was inspected for additional license/copyright notices;
its license detections in explanatory documentation are not new code licenses.
Source matching, scanning, tests and Git history are evidence for review, not a
replacement for the submitting human's responsibility. The exact historical AI
models/prompts have not been reconstructed. CONTRIBUTING.md discloses known
usage without inventing historical human approvals. Fedora reviewers and the
human maintainer must make the final acceptance decision for the actual package.
