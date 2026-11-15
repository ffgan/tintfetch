#!/usr/bin/env python3
"""Check that a release includes matching complete source and license notices.

Run from the same source checkout used to compile the release. Archive contents
are read in memory; no archive paths are extracted. Builds/publish must finish
this check before sharing binary archives.
"""
import argparse
from pathlib import Path
import tarfile
import zipfile

ROOT = Path(__file__).resolve().parent.parent
SOURCE_DIRS = ('assets', 'cmd', 'configs', 'internal', 'scripts', 'docs', 'packaging', 'debian', 'e2e', '.github')
SOURCE_FILES = ('main.go', 'go.mod', '.go-version', 'Makefile', 'LICENSE', 'COPYING', 'THIRD_PARTY_NOTICES', 'README.md', 'CONTRIBUTING.md', 'SECURITY.md', 'CHANGELOG.md', '.gitignore', '.editorconfig')
NOTICES = ('LICENSE', 'COPYING', 'THIRD_PARTY_NOTICES')


def required_source():
    result = set(SOURCE_FILES)
    for name in SOURCE_DIRS:
        for path in (ROOT / name).rglob('*'):
            if path.is_file() and '__pycache__' not in path.parts and path.suffix != '.pyc':
                # Debian and Go build artifacts are not corresponding source.
                rel = path.relative_to(ROOT).as_posix()
                if rel.startswith(('debian/.debhelper/', 'debian/tintfetch/')) or path.name in ('debhelper-build-stamp', 'files') or path.name.endswith(('.substvars', '.debhelper')):
                    continue
                result.add(rel)
    return result


def tar_contents(path):
    with tarfile.open(path, 'r:gz') as archive:
        result = {}
        for member in archive.getmembers():
            if member.isdir():
                continue
            if not member.isfile() or member.name in result:
                raise ValueError(f'non-regular or repeated archive entry: {member.name}')
            result[member.name] = archive.extractfile(member).read()
        return result


def check(directory, version):
    source_name = f'tintfetch_{version}_source.tar.gz'
    source = tar_contents(directory / source_name)
    prefix = f'tintfetch-{version}/'
    for name in required_source():
        if source.get(prefix + name) != (ROOT / name).read_bytes():
            raise ValueError(f'corresponding source missing or differs from build checkout: {name}')
    binaries = sorted(p for p in directory.iterdir() if p.name.startswith(f'tintfetch_{version}_') and p.name != source_name and (p.name.endswith('.tar.gz') or p.suffix == '.zip'))
    if not binaries:
        raise ValueError('release contains no binary archives')
    for path in binaries:
        if path.suffix == '.zip':
            with zipfile.ZipFile(path) as archive:
                if len(set(archive.namelist())) != len(archive.namelist()):
                    raise ValueError(f'duplicate ZIP entries: {path.name}')
                contents = {name: archive.read(name) for name in archive.namelist()}
        else:
            contents = tar_contents(path)
        for notice in NOTICES:
            if contents.get(notice) != (ROOT / notice).read_bytes():
                raise ValueError(f'{path.name}: missing or outdated {notice}')
        if source_name.encode() not in contents.get('SOURCE.txt', b''):
            raise ValueError(f'{path.name}: corresponding-source notice missing or incorrect')
    print(f'Checked matching corresponding source and notices for {len(binaries)} binary archives')


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('directory', type=Path)
    parser.add_argument('version')
    args = parser.parse_args()
    try:
        check(args.directory, args.version)
    except (OSError, ValueError, tarfile.TarError, zipfile.BadZipFile) as exc:
        raise SystemExit(str(exc)) from exc


if __name__ == '__main__':
    main()
