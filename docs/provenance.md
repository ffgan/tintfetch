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
binary packages; the application source remains MIT licensed. There are no
third-party Go modules. A distribution reviewer may still request individual
artwork attribution or further translation history before acceptance.
