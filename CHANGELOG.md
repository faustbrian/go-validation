# Changelog

All notable changes follow Keep a Changelog. The public API follows semantic
versioning.

## [Unreleased]

### Changed

- Replace copied repository verification tooling with the checksum-pinned
  `go-library-tools` v1.0.7 contract while preserving package behavior,
  public APIs, dependency checksums, and mutation evidence.

### Documentation

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
