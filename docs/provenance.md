# Source and license provenance

The maintainer identified `ffgan/hyfetch` as the migration source. On 2026-11-15,
its HEAD was `2252fc36e76eebc42030c383784f63f9489f5630`; this is an audit
reference, **not a claim that the original migration used that exact commit**.
The reference's `neofetch` script identifies itself as version 8.0.4.

Reference files:

- https://github.com/ffgan/hyfetch/blob/2252fc36e76eebc42030c383784f63f9489f5630/neofetch
- https://github.com/ffgan/hyfetch/blob/2252fc36e76eebc42030c383784f63f9489f5630/LICENSE.md

The original 384 files in `assets/`, with only leading/trailing newlines removed for
comparison, occur verbatim in the reference script. The script includes Dylan
Araps's MIT notice (2015–2021); the reference repository's MIT license extends
the year range to 2015–2023. Both the project license and third-party notices
retain the upstream attribution. `cmd/README` and `configs/config` are legacy
reference material and also retain their existing MIT notices.

To reproduce the asset check, download the pinned `neofetch` reference and run:

```sh
python3 scripts/check-provenance.py /path/to/neofetch
```

The script checks the downloaded reference's SHA-256, the local asset manifest,
and the text matches. It does not download anything during builds.

The Go implementation contains manually translated and modified functions.
The maintainer's source identification and existing upstream references establish
the migration lineage; this audit does not reconstruct each translated line's
historical source revision. New contributors should identify provenance in their
commits and preserve any additional applicable notices. A distribution reviewer
may require further history or individual artwork attribution before acceptance.

Compiled executables contain the Go runtime and standard library. Their BSD
redistribution notice is included in THIRD_PARTY_NOTICES and installed with
binary packages; the Tintfetch application source remains MIT licensed.

The merge of origin/master at de3e5aa added 17 further artwork files. They
are inherited from the remote implementation and do not match the pinned
reference verbatim. Their checksums are recorded, and the paths are listed in
assets-pending-provenance.txt. Their upstream attribution must be reviewed
before a formal distribution submission. The source audit script distinguishes
these pending files from the 384 verified files instead of certifying them.
