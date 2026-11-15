#!/usr/bin/env python3
"""Create the portable, deterministic corresponding-source archive for a build."""
import argparse
import gzip
import importlib.util
import os
from pathlib import Path
import re
import subprocess
import tarfile

ROOT = Path(__file__).resolve().parent.parent


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('version')
    parser.add_argument('output', type=Path)
    args = parser.parse_args()
    if not re.fullmatch(r'[A-Za-z0-9][A-Za-z0-9.+~_-]*', args.version):
        raise SystemExit('version must be an archive-safe version string')
    spec = importlib.util.spec_from_file_location('release_source', ROOT / 'scripts/check-release-source.py')
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    epoch = os.environ.get('SOURCE_DATE_EPOCH')
    if epoch is None:
        try:
            result = subprocess.run(['git', 'log', '-1', '--format=%ct'], cwd=ROOT, text=True, capture_output=True)
            recorded = result.stdout.strip() if result.returncode == 0 else '0'
        except FileNotFoundError:
            recorded = '0'
        # Exported source trees have no Git metadata. Zero remains deterministic;
        # packagers can supply the actual release timestamp explicitly.
        epoch = recorded
    timestamp = int(epoch)
    if timestamp < 0:
        raise SystemExit('SOURCE_DATE_EPOCH must be nonnegative')
    args.output.parent.mkdir(parents=True, exist_ok=True)
    with args.output.open('wb') as output, gzip.GzipFile(filename='', mode='wb', fileobj=output, mtime=0) as compressed, tarfile.open(fileobj=compressed, mode='w', format=tarfile.PAX_FORMAT) as archive:
        prefix = f'tintfetch-{args.version}/'
        for name in sorted(module.required_source()):
            path = ROOT / name
            if path.is_symlink() or not path.is_file():
                raise SystemExit(f'source must be a regular file: {name}')
            info = tarfile.TarInfo(prefix + name)
            info.size = path.stat().st_size
            info.mode = 0o755 if path.stat().st_mode & 0o111 else 0o644
            info.mtime = timestamp
            with path.open('rb') as source:
                archive.addfile(info, source)
    print(f'Created corresponding source: {args.output}')


if __name__ == '__main__':
    main()
