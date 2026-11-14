#!/bin/sh
set -eu
stage=$(mktemp -d)
trap 'rm -rf "$stage"' EXIT HUP INT TERM
mkdir -p "$stage/one" "$stage/two"
# Build the same source from different absolute paths, without VCS metadata.
tar --exclude=./.git --exclude=./build -cf "$stage/source.tar" .
tar -xf "$stage/source.tar" -C "$stage/one"
tar -xf "$stage/source.tar" -C "$stage/two"
(cd "$stage/one" && GOFLAGS=-buildvcs=false make build VERSION=0.1.0-test)
(cd "$stage/two" && GOFLAGS=-buildvcs=false make build VERSION=0.1.0-test)
cmp "$stage/one/build/tintfetch" "$stage/two/build/tintfetch"
echo 'Reproducible binary verified across two source paths'
