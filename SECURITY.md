# Security Policy

## Reporting

Report suspected vulnerabilities privately through the GitHub security
advisory for `faustbrian/go-validation`. Do not open a public issue containing exploit
details, credentials, private fixtures, or affected deployment information.

Include the affected module and version, impact, reproduction, preconditions,
and any suggested mitigation. Maintainers follow the pinned ecosystem
[vulnerability-management process](https://github.com/faustbrian/go-library-tools/blob/77bfd78c12a853f0d490bb27a3fbcb5f34330772/docs/ecosystem/security/vulnerability-management.md)
for severity assignment, acknowledgement and remediation targets, private
embargo handling, advisory ranges, and coordinated affected-module releases.
Targets begin when sufficient private evidence exists to reproduce or bound
the report; the owner explains any changed target to the reporter.

## Supported Versions

The published `github.com/faustbrian/go-validation/v2` line, starting at
`v2.0.0`, is the current security-maintained baseline. Its safeguards and
residual ownership are documented in [the security model](docs/security.md).

The earlier `github.com/faustbrian/go-validation` identity remains publicly
available, including `v1.1.0` consumed by OpeningHours and Temporal and the
toolchain-only `v1.2.0` release. Those versions do not acquire the v2 safeguards
or qualification from this policy or from current-main CI. Their continued
consumer adoption and security-support disposition remain unresolved; this
document does not waive that boundary or certify retained v1 consumers.
Cross-major named types are not interchangeable, so migration must be assessed
at each consumer's public contract. See [`COMPATIBILITY.md`](COMPATIBILITY.md).

## Security Gates

Checks follow the changed trust boundary and applicable release risk in
[`AGENTS.md`](AGENTS.md#proportional-assurance). Ordinary required CI remains
separate from the opt-in native security diagnostic, which qualifies owned
analysis and current-tree/history secret policy for the exact checked-out
source. A diagnostic cannot substitute for ordinary CI or a release rehearsal,
and enabling it does not claim its gates have passed.

The published `v2.0.0` release is a source-only Go module release: it does not
claim uploaded SBOM, provenance, or binary assets. Required supply-chain
evidence is selected and verified at the applicable release boundary; absence
of an asset that was not promised is not evidence that a required gate passed.
A missing selected scanner or unavailable required service is a failed gate,
not a warning.

Security fixes MUST include a regression test that does not publish weaponized
details or real secrets. Credentials MUST be redacted from logs and evidence.

## Repository Assurance

The repository [safety and concurrency policy](AGENTS.md#safety-and-concurrency)
and [supply-chain policy](AGENTS.md#dependencies-and-supply-chain) define shared
trust boundaries and release requirements. Package-specific security guidance
refines those rules for its owned boundary.
