#!/usr/bin/env python3
"""Check local audited bytes, optionally against separately fetched references.

No network access is performed. --fetch prints pinned URLs for the reviewer;
place each downloaded file under its reference ID in --upstream-dir.
"""
import argparse
import hashlib
import json
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent


def digest(data):
    return hashlib.sha256(data).hexdigest()


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--fetch', action='store_true', help='print reference IDs and URLs only')
    parser.add_argument('--upstream-dir', type=Path, help='verify pinned upstream files in this directory')
    args = parser.parse_args()
    audit = json.loads((ROOT / 'docs/provenance.json').read_text())
    if args.fetch:
        for key, ref in audit['references'].items():
            print(key, ref['url'])
        return
    files = audit['files']
    expected = {str(p.relative_to(ROOT)) for p in (ROOT / 'assets').iterdir()}
    expected |= {str(p.relative_to(ROOT)) for p in (ROOT / 'internal/hyfetch/data').iterdir()}
    if set(files) != expected:
        raise SystemExit('audit inventory differs from embedded assets/data')
    manifest = dict(line.split('  ', 1)[::-1] for line in (ROOT / 'docs/assets.sha256').read_text().splitlines())
    if manifest != {name: entry['sha256'] for name, entry in files.items() if name.startswith('assets/')}:
        raise SystemExit('asset checksum manifest differs from provenance audit')
    references = {}
    if args.upstream_dir:
        for key, ref in audit['references'].items():
            data = (args.upstream_dir / key).read_bytes()
            if digest(data) != ref['sha256']:
                raise SystemExit(f'upstream hash differs: {key}')
            references[key] = data
    for name, entry in files.items():
        data = (ROOT / name).read_bytes()
        if digest(data) != entry['sha256']:
            raise SystemExit(f'local file changed without provenance review: {name}')
        if entry['reference'] not in audit['references']:
            raise SystemExit(f'unknown upstream reference: {name}')
        if entry['comparison'] not in ('text', 'bytes', 'json-subset'):
            raise SystemExit(f'unknown comparison: {name}')
        if not references:
            continue
        ref = references[entry['reference']]
        if entry['comparison'] == 'text':
            matches = bool(data.strip(b'\n')) and data.strip(b'\n') in ref
        elif entry['comparison'] == 'bytes':
            matches = data == ref
        else:
            local, upstream = json.loads(data), json.loads(ref)
            matches = all(key in upstream and value == upstream[key] for key, value in local.items())
        if not matches:
            raise SystemExit(f'local file does not match upstream: {name}')
    scope = 'pinned upstream references' if references else 'audited local checksums (upstream comparison not run)'
    print(f'Verified {len(files)} files against {scope}')


if __name__ == '__main__':
    main()
