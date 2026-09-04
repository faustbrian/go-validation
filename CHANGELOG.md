# Changelog

All notable changes follow Keep a Changelog. The public API follows semantic
versioning.

## [Unreleased]

### Changed

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

[Unreleased]: https://github.com/faustbrian/go-validation/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/faustbrian/go-validation/releases/tag/v1.0.0
