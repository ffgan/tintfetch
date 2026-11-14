#!/usr/bin/env python3
"""Verify the audited asset snapshot against the pinned upstream script."""
import hashlib
from pathlib import Path
import sys

ROOT = Path(__file__).resolve().parent.parent
REFERENCE_SHA256 = "129a266b0f4fa8bbf34f45f135d2e5f02eb392a169a58c2dd85b570b4754aee0"

def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: check-provenance.py /path/to/pinned/neofetch")
    upstream = Path(sys.argv[1]).read_bytes()
    if hashlib.sha256(upstream).hexdigest() != REFERENCE_SHA256:
        raise SystemExit("upstream hash differs from the pinned audit reference")
    manifest = {}
    for line in (ROOT / "docs/assets.sha256").read_text().splitlines():
        digest, name = line.split("  ", 1)
        manifest[name] = digest
    assets = sorted((ROOT / "assets").iterdir())
    if set(manifest) != {str(p.relative_to(ROOT)) for p in assets}:
        raise SystemExit("asset manifest does not match the asset collection")
    text = upstream.decode()
    for asset in assets:
        data = asset.read_bytes()
        if hashlib.sha256(data).hexdigest() != manifest[str(asset.relative_to(ROOT))]:
            raise SystemExit(f"local asset changed: {asset.name}")
        if data.decode().strip("\n") not in text:
            raise SystemExit(f"asset not found in reference: {asset.name}")
    print(f"Verified {len(assets)} assets against pinned MIT-licensed reference")

if __name__ == "__main__":
    main()
