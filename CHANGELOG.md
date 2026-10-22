# Changelog

All notable changes follow Keep a Changelog. The public API follows semantic
versioning.

## 3.0.1 - 2026-10-22

### Changed

- Publish the current source and documentation with a fresh immutable tag.
- Build with Go 1.27.3 to address GO-2026-6809 in `crypto/x509`.
  Consumers should use Go 1.27.3 or newer.

## [Unreleased]

## [3.0.0] - 2026-10-13

### Breaking

- Prepare v3.0.0 under `github.com/faustbrian/go-validation/v3`; migrate every
  root and subpackage import together. All sixteen packages remain at the
  repository root. Published v2 stays available; major-module named types and
  error identities are distinct. See the [v3 migration guide](docs/migration.md#adopting-v3).
- Empty/all-nil `Any` now emits `any` because no alternative succeeded.
  Use `All` for an intentionally empty set of independent constraints.
- `Keys`/`Values` reject NaN floating keys with a map-level `key` finding
  before child invocation; values are taken from actual entries. NaN equality
  in `Unique` and `FieldsEqual` remains ordinary Go equality.

### Fixed

- Tag plans reject cyclic pointer/interface normalization with field-local
  `cycle` findings, preserving finite acyclic chains and pointer presence.
  `MaxDepth` remains compile-time structural depth, not a runtime hop cap.
- Recursive untagged pointer types fail construction with `ErrCycle`; failed
  compilation does not consume a plan-cache entry or prevent later compilation.

### Documentation

- Clarify callback/input ownership, per-call async joining, partial report
  state, projection/helper variants and resource admission. Distinguish
  compile-time type checks from existing runtime dereference-cycle limitations
  and describe package-defined comparison/primitive profiles truthfully.

### Maintenance

- Adopt patched documentation and authenticated shared tooling, bounded
  selected verification, source documentation contracts and coherent
  contribution/adoption navigation. Runtime dependency versions and the
  Go 1.27 minimum remain unchanged by this major-module transition.

### Source accountability

- [46d76915](https://github.com/faustbrian/go-validation/commit/46d7691580d2b242fd4ebe9002cc9e557683ea00) — ci: use patched documentation tooling.
- [f9e4acb3](https://github.com/faustbrian/go-validation/commit/f9e4acb39521ac917fb7f47bca07e006cc7d5c3d) — ci: adopt bounded published verification.
- [618e441c](https://github.com/faustbrian/go-validation/commit/618e441cab0148b44d9fa5d19f4de0e868dfe5c4) — ci: add bounded selected assurance.
- [19d47174](https://github.com/faustbrian/go-validation/commit/19d471746d29dcae72fd0f65e797b90c4d94f240) — docs: standardize repository development policy.
- [2c10a5bd](https://github.com/faustbrian/go-validation/commit/2c10a5bd3ce379098e1d2778bb75f6bce4e3a137) — chore(ci): pin shared tooling to v2.3.0.
- [cb316be7](https://github.com/faustbrian/go-validation/commit/cb316be749e6ed45fd534491b599173961640fbf) — chore(ci): align manual tooling source pins.
- [c4a9fda1](https://github.com/faustbrian/go-validation/commit/c4a9fda162bb302d02291496b35e29c869c2d03c) — docs: clarify validation ownership and limits.
- [1366c3fe](https://github.com/faustbrian/go-validation/commit/1366c3fea0026e5dac0bd7964c6e2fc935bcb0d8) — chore: enforce source documentation contracts.
- [717be798](https://github.com/faustbrian/go-validation/commit/717be798b83ce345acb419c4ae805d7cc3a4171f) — docs: align adoption and documentation navigation.
- [bdaf4e0a](https://github.com/faustbrian/go-validation/commit/bdaf4e0aa94c6b3188bf1cdedaa3612f69fdac4f) — docs: repair assurance policy navigation.
- [7bf88065](https://github.com/faustbrian/go-validation/commit/7bf880659e78e790b5adef901c35fc33e3186024) — fix: correct validation boundaries in v3.

## [2.0.1] - 2026-10-06

### Documentation

- Record owners, rationale, mitigations and review conditions for existing
  callback and diagnostic trust boundaries without expanding guarantees or
  claiming whole-family security completion.

- Identify the diagnostic privacy and validator fanout admission changes in
  the security model as published v2.0.0 behavior, distinct from historical v1
  artifacts and whole-family security completion.

### Maintenance

- Refresh the immutable shared workflow and add opt-in exact-source
security diagnostics without replacing ordinary CI, release
qualification, or downstream adoption evidence.

## [2.0.0]

These historical notes describe safeguards already published in v2.0.0,
not changes introduced by this patch.

### Changed

- Prepare v2.0.0 under `github.com/faustbrian/go-validation/v2`, retaining
  Go 1.27.0 and all sixteen packages at their existing source locations.
  Update every validation import to include `/v2`; named types from different
  major module paths are not assignment-compatible. See the
  [migration guide](docs/migration.md#adopting-v2) for the behavioral changes
  below. Existing v1 consumers remain on their published module identity.

- Make HTTP and service function hooks in canonical and retained packages
  contain panics through the core validator adapter, discarding private panic
  payloads. Arbitrary interface implementations remain application-owned and
  require explicit isolation. HTTP problem writers now reject status values
  outside 100–999 before touching response headers or the body; callers receive
  an error classified as `ErrInvalidViolation` instead of a partial response
  or a downstream invalid-status panic.

- Bound the next major release's `All`, `Any`, `AsyncAll`, and both service
  chains by `MaxCollectionSize`, counting all supplied positions including nil.
  Oversized fanout fails before invocation or async result/worker allocation.
  Set an explicit appropriate limit for larger trusted compositions; caller
  cancellation continues to take precedence for async and service operations.

- Prepare the next major release's private default diagnostics for pattern
  compilation and struct-plan construction. Errors no longer print caller
  expressions, field names, type names, or rule identifiers. Match stable
  classifications with `errors.Is` and inspect original causes explicitly with
  `errors.As` or `errors.Unwrap` instead of parsing default error text.

## [1.2.0] - 2026-10-01

### Changed

- Require Go 1.27.0 instead of Go 1.26.6. Upgrade consumer toolchains
  before adopting v1.2.0; validation APIs, runtime contracts, and wire
  formats remain unchanged. The minor release signals the supported
  toolchain transition without changing the public API.

- Keep reusable CI and its checked-out tooling on the same v1.8.4 source
  while retaining the checksum-verified v1.4.0 CLI bootstrap.

## [1.1.0] - 2026-09-06

### Added

- Add immutable cancellation and deadline outcomes to `Report`, including
  `ContextReport`, `Report.ContextError`, and structured error projections that
  preserve partial validation findings.
- Add target-oriented `adapters/config`, `adapters/http`, `adapters/jsonapi`,
  `adapters/jsonrpc`, and `adapters/service` packages while retaining the v1
  integration paths.

### Changed

- Make `AsyncAll`, validation compositors, collection and struct-plan
  traversal, test helpers, and service chains preserve terminal context
  identity, ordered partial findings, and bounded admitted-work ownership.
- Deprecate the five original integration paths in favor of their additive
  target-oriented successors; retain them for the longer of 180 days after
  successor availability and two stable minor releases.

- Adopt the `go-library-tools` v1.3.0 schema-v2 cohesion contract and local
  `make cohesion` gate without changing validation API or runtime behavior.
- Pin reusable CI to the immutable v1.3.0 workflow and enforce cohesion
  metadata in the repository's required CI contract.

- Upgrade the checksum-pinned CLI and immutable shared workflow to
  `go-library-tools` v1.4.0, add online specification governance to local
  `make ci`, and preserve the validation API and runtime behavior.

- Replace copied repository verification tooling with the checksum-pinned
  `go-library-tools` v1.2.0 contract so CI executes specification governance
  with the released binary while preserving package behavior,
  public APIs, dependency checksums, and mutation evidence.

### Documentation

- Publish the module's family, capabilities, ownership, lifecycle, supported
  environments, package selection, and delivery status, and link the README to
  the immutable v1.3.0 ecosystem index and family guidance.
- Advance ecosystem and Foundations family guidance to the immutable v1.4.0
  documentation release.

- Add the [specification decision register](docs/specification-decisions.md)
  and conformance governance for bounded RFC 6901, RFC 9457, JSON:API 1.1,
  and JSON-RPC 2.0 claims:
  - VALIDATION-DEC-001 sha256:95c377bc8aff8aaba04b1e6add8daaf3f65c8823cb4ecb75ee0d56003220b80c
  - VALIDATION-DEC-002 sha256:16768491f285e40cd325a055fce9de7de4425f5d7e788d68895305c5019d8e66
  - VALIDATION-DEC-003 sha256:adb41196694f54e3eb093a6a94a67ee9d2c26baf2d1708f954fd1a5912e301cc
  - VALIDATION-DEC-004 sha256:ae1b9622b0ed448c4c2b91e2eb016db6882a82ac092d6b08668823a1e1fd6edf
  - VALIDATION-DEC-005 sha256:ad54e8018ffdec79c08070e8dcec8473d95b04d2b518ad7033f40e73f5bc08ca
  - VALIDATION-DEC-006 sha256:721619a7666ab86da9fd252bf8d05beb9b9568f9e25f0ac569e03a01dca14bc8

- Replace archived monorepo links and completed execution artifacts with a
  standalone, human-oriented documentation structure.

## [1.0.0] - 2026-08-25

### Changed

- Bind the reviewed zero-mutant validation-config inventory to its canonical
  standalone source identity.

- Exclude intentional nested modules from root local-proxy archives so local,
  bootstrap, CI, and public module checksums describe the same source
  boundary.

- Track the pinned documentation-tool lockfile so clean CI checkouts install
  the exact validated cspell dependency.

- Reconcile standalone dependency checksums against deterministic current
  module archives so CI, local verification, and release consumers resolve
  identical content.

- Harden standalone documentation validation with deterministic spelling and
  link checks, package-specific documentation gates, and repository-local
  contributor guidance.

### Documentation

- Replace obsolete standalone-repository links and workflow claims with
  monorepo-canonical targets and current release guidance.

- Link the package README to package-owned documentation.

### Changed

- Publish the module from its standalone `github.com/faustbrian/go-validation` identity while preserving its documented API and behavior.
- Preserve precomputed validator configuration through explicit constructors,
  removing analyzer-ambiguous closure assignments without changing behavior.

### Added

- Typed validators, immutable contexts, presence-aware values, typed paths,
  safe violations, bounded reports, and deterministic composition.
- Standard presence, string, numeric, collection, temporal, cross-field,
  network, email, UUID, and identifier rules.
- Bounded async execution, typed and strict-tag struct plans, transport
  projections, integration hooks, localization, observations, and test tools.
- Fail-closed truncation, structural path identity, reflective typed parity,
  severity-aware projections, and sanitized observation labels.
- Collision-safe parameter identity, exact reflective numeric bounds, and
  transport-visible report blocking state.
- Exact coverage, race, fuzz, mutation, benchmark, vulnerability, API, docs,
  lint, and release automation.
- Automatic secret-safe panic containment in synchronous and asynchronous
  custom function adapters.
- A context-owned string-size budget enforced before typed and reflective
  parsing, matching, comparison, sorting, or hashing.
- Strict malformed typed-plan rejection, field-local extension panic
  containment, total reflective field accounting, and plan-construction fuzzing.
- Complete cross-kind presence tables and warning-preserving `Any` and
  `Dependent` composition.
- Fail-closed custom diagnostic severity, code, parameter, UTF-8, control, and
  safe-formatting enforcement.
- Panic-contained arbitrary async validation with deadline evidence, hostile
  translation sanitization, and cross-transport conformance coverage.
- Allocation-free hostile path rejection plus cache teardown races,
  caller-data immutability, and numeric edge evidence.
- A versioned six-target fuzz-corpus inventory and 90/90 mutation evidence
  spanning every standard-rule family and all hardening-critical boundaries.

[Unreleased]: https://github.com/faustbrian/go-validation/compare/v1.2.0...HEAD
[1.2.0]: https://github.com/faustbrian/go-validation/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/faustbrian/go-validation/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/faustbrian/go-validation/releases/tag/v1.0.0
