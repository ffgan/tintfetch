# Contributing

Use Go 1.26 or newer. Keep changes focused, describe the problem and resulting
behavior, and add regression tests for changes to parsing, command execution,
or probe logic. Preserve source and license attribution when importing code or
artwork. Contributions are licensed under the repository's MIT license; no CLA
is currently required.

Before opening a pull request:

```sh
gofmt -w cmd configs internal e2e main.go
make check
go test -race ./...
scripts/smoke.sh build/tintfetch
```

Run `make build` before the smoke test. Tests must not require a real desktop,
network access, elevated privileges, or a particular distribution. Use fixtures
for probe behavior. Do not turn an unavailable optional probe into a process
exit; return no value or a useful error instead. Keep process termination at the CLI boundary; probe libraries must return errors. Keep external processes bounded and pass arguments
without a shell unless a feature explicitly opts into Bash execution.

Add native feature support and tests before marking an option implemented in
`docs/compatibility.md`. Keep user documentation and the manual consistent.
Retain existing native CLI and JSON contracts; exercise the real binary through e2e tests.

To report a bug, include the build version, OS/architecture, command, expected
behavior and diagnostic output. Remove hostnames, usernames, IP addresses and
other personal information before sharing output. For vulnerabilities, follow
SECURITY.md.
