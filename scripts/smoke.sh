#!/bin/sh
set -eu
binary=${1:-build/tintfetch}
stage=$(mktemp -d)
trap 'rm -rf "$stage"' EXIT HUP INT TERM
export HOME="$stage/home" XDG_CONFIG_HOME="$stage/home/.config"
mkdir -p "$XDG_CONFIG_HOME"
"$binary" --version
"$binary" --print-config > "$stage/config.conf"
"$binary" --config "$stage/config.conf" --plain > "$stage/plain"
"$binary" --config "$stage/config.conf" --json > "$stage/info.json"
python3 - "$stage" <<'PY'
import json, pathlib, sys
root = pathlib.Path(sys.argv[1])
data = json.loads((root / 'info.json').read_text())
assert isinstance(data['result'], list)
assert '\x1b' not in (root / 'plain').read_text()
assert '\x1b' not in json.dumps(data)
assert data['result'], 'JSON output must contain configured information'
PY
if "$binary" --no-config --unknown-option > "$stage/out" 2> "$stage/err"; then
    echo 'unknown option unexpectedly accepted' >&2
    exit 1
fi
