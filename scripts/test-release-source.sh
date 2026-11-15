#!/bin/sh
# Exercise real tar/ZIP releases and the corresponding-source publication gate.
set -eu
root=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
cd "$root"
stage=$(mktemp -d)
trap 'rm -rf "$stage"' EXIT HUP INT TERM
version=0.0.0-source-test
VERSION="$version" DIST_DIR="$stage/dist" bash scripts/build-release.sh linux/amd64 windows/amd64
make dist VERSION="$version"
cmp "build/tintfetch-$version.tar.gz" "$stage/dist/tintfetch_${version}_source.tar.gz"
python3 scripts/check-release-source.py "$stage/dist" "$version"
python3 - "$root" "$stage/dist" "$version" <<'PY'
from pathlib import Path
import io
import subprocess
import sys
import tarfile
import zipfile
root, directory, version = Path(sys.argv[1]), Path(sys.argv[2]), sys.argv[3]
cmd = [sys.executable, str(root / 'scripts/check-release-source.py'), str(directory), version]

def rejects(label, expected):
    result = subprocess.run(cmd, text=True, capture_output=True)
    assert result.returncode and expected in result.stderr, (label, result.stdout, result.stderr)
    print(label + ': correctly rejected')

source = directory / f'tintfetch_{version}_source.tar.gz'
original_source = source.read_bytes()
source.unlink()
rejects('Missing corresponding source', source.name)
source.write_bytes(original_source)
with tarfile.open(fileobj=io.BytesIO(original_source), mode='r:gz') as original, tarfile.open(source, 'w:gz') as changed:
    for member in original.getmembers():
        if member.name.endswith('/go.mod'):
            continue
        changed.addfile(member, original.extractfile(member) if member.isfile() else None)
rejects('Incomplete corresponding source', 'corresponding source missing or differs')
source.write_bytes(original_source)
archive = directory / f'tintfetch_{version}_windows_amd64.zip'
original_zip = archive.read_bytes()
with zipfile.ZipFile(io.BytesIO(original_zip)) as original, zipfile.ZipFile(archive, 'w') as changed:
    for item in original.infolist():
        if item.filename != 'COPYING':
            changed.writestr(item, original.read(item))
rejects('Missing GPL license in binary ZIP', 'missing or outdated COPYING')
archive.write_bytes(original_zip)
subprocess.run(cmd, check=True)
PY
