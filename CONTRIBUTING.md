# Contributing

## Before Editing

1. Read [`AGENTS.md`](AGENTS.md) and the affected package's docs.
2. Run `go test ./...` or the affected package's narrower tests.
3. Identify owned dependencies and reverse dependants in `modules.json`.
4. Preserve unrelated work and generated/corpus provenance.

## Changes

Keep commits focused and conventional. Update every affected changelog with
the behavior and migration impact. Public API changes require compatibility
evidence and documentation. Specification behavior requires a decision record,
fixture coverage, and interoperability evidence.

New direct dependencies and dependency updates must follow the
[dependency governance policy](AGENTS.md#dependencies-and-supply-chain). Package-local
update bots are forbidden; the root policy owns every module and action update.

Specification-backed changes must follow the
[specification governance contract](AGENTS.md#design), update
the affected entries in the
[specification decision register](docs/specification-decisions.md), and
complete the Specification Decisions section of the pull request template.
Record the decision identifier, compatibility and wire impact, changelog
digest, and any superseded identifier with its replacement. An unresolved
interpretation or stale source pin is release-blocking; peer behavior cannot
silently select policy.

Mutation is optional local or scheduled assurance, not a required PR status.
Confirmed behavioral and security defects still require a focused regression
and repair; equivalent/invalid mutants retain reviewed narrow dispositions.

Do not add package-local workflows, permanent replacements, machine-specific
paths, bypass flags, broad mutation exclusions, or aggregate quality metrics
that hide a failing package.

## Verification

Standalone `go test ./...` runs deterministic tests, fuzz seeds and saved corpus
regressions without sibling checkouts or ecosystem ledgers. The manifest has no
required services, build tags or test tags.

With the checksum-pinned public verifier selected by `.golib.yaml`, use:

```bash
make inventory
make check
make ci
```

`check` and `ci` use bounded local checks. Three named root race assertions
exercise `AsyncAll` concurrency/order bounds, cancellation stopping admission
while joining callbacks and retaining partial findings, and deadline cleanup
with retained ordered warnings. They run once with a two-minute timeout, not a
broad stress campaign. Ordinary CI consumes the public verifier without a source
build. The separate manual security diagnostic retains its selected source
qualification, release-substitution refusal and cleanup.

`make specification-check` validates the offline register and history;
`make specification-check-online` additionally re-fetches reviewed authority
digests. Ordinary `ci` uses the offline boundary, while `make check-full` retains
the online refresh. Existing hosted Required and security controls remain intact.

Select one target for short active fuzz exploration. Optional local mutation
selects one package; neither command is an all-change prerequisite:

```bash
GOMAXPROCS=2 go test -run '^$' -fuzz '^FuzzPathAndReportSafety$' -fuzztime=2s -parallel=1 -timeout=2m .
golib mutation --module . --package ./rules
```

Path, transport, async lifecycle, API and specification changes retain affected
conformance, compatibility, race and direct-adapter/consumer checks. Preserve
the minimized fuzz corpus and source authority decisions.

`make check-full` retains full assurance, all seven native 10000x fuzz budgets,
conformance and benchmarks. Relevant release qualification remains
`golib release dry-run --module .`. Select full/release checks for their material
risk or release boundary, not every commit. Report unavailable, skipped and
failed boundaries honestly; a local pass is not online authority, hosted-CI or
public-release qualification.

## Adding A Module

Follow [repository structure policy](AGENTS.md#repository-structure). New modules
require an explicit purpose, ownership boundary, dependency review, package
catalog entry, full quality gates, documentation, changelog, license, security
policy, compatibility plan, and release dry-run.
