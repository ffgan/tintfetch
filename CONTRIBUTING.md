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

## AI assistance and human review

Tintfetch has used substantial AI-assisted development, as confirmed by the
maintainer. The history includes commits attributed to Cursor Agent and
cursor[bot]. Recent source-provenance and distribution-packaging work has also
been assisted by OpenAI Codex. The exact tools, models and prompts for every
historical commit have not been reconstructed. Agent authorship in Git is not
evidence that a human reviewed the resulting code or cleared its provenance.

AI assistance is allowed, but the human contributor remains responsible for
understanding, reviewing and testing the submitted changes and establishing
permission to redistribute any incorporated material. Disclose substantial AI
assistance in the PR and commit message, including the tool name and affected
scope; an `Assisted-by: OpenAI Codex` or `Assisted-by: Cursor` trailer is suitable
when accurate. Do not invent historical tool/model details, human sign-offs,
review approvals or copyright ownership.

Treat generated code, notices, citations and tests as unverified suggestions
until checked. Record external source URLs and revisions, preserve applicable
license notices, and flag uncertain origins instead of assigning the repository's
MIT license by assumption. AI-generated tests and automated license scans do not
establish legal provenance or substitute for human review.

For Fedora submissions, disclose AI assistance in packaging, patches and review
text as applicable. A human maintainer must review and own the submission; an AI
assistant cannot issue the final acceptance decision or claim Fedora approval.
See https://docs.fedoraproject.org/en-US/council/policy/ai-contribution-policy/.
