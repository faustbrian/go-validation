# Migration and versioning

## Adopting v2

After v2.0.0 is publicly available, require
`github.com/faustbrian/go-validation/v2@v2.0.0` and insert `/v2` after
`go-validation` in every root and subpackage import. Go 1.27.0 remains the
minimum toolchain. The module stays at the repository root; there is no `v2`
source directory. Do not bridge v1 and v2 named types by assignment.

All sixteen packages remain available, including the five deprecated
integration paths under the v2 module. Prefer `adapters/config`,
`adapters/http`, `adapters/jsonapi`, `adapters/jsonrpc`, and `adapters/service`
for new integrations. Their existing deprecation policy is unchanged; the
retained v2 paths do not preserve cross-major Go type identity. Existing
consumers may continue to pin published v1 versions independently.

Review these intentional behavioral changes at each application boundary:

- Pattern and struct-plan construction errors no longer print expressions,
  field names, type names, or rule identifiers by default. Use `errors.Is`
  for stable classifications. Explicit `errors.As` or `errors.Unwrap` access
  to original causes is a trusted diagnostic boundary, not safe public prose.
- `All`, `Any`, `AsyncAll`, and both service chains count every supplied
  position, including nil, against `MaxCollectionSize`. Oversized fanout is
  refused before invocation or async allocation; set an appropriate explicit
  limit for larger trusted compositions. Cancellation still takes precedence
  for async and service operations.
- Canonical and retained HTTP/service function hooks contain panics and
  discard payloads. Arbitrary interface implementations remain caller-owned
  and require explicit isolation. HTTP problem writers reject statuses outside
  100–999 as `ErrInvalidViolation` before changing headers or writing a body.

Re-run application response fixtures and success/cancellation checks rather
than comparing private diagnostic text. The released v1.2.0 API snapshot is
preserved in `api/v1.2.0.txt`; `api/baseline.txt` tracks the current v2 API.

## Earlier v1 migrations

Before adopting v1, pin the module version and record existing response payload
fixtures. Migrate one boundary at a time using the [adoption guide](adoption.md).

Public exported symbols, stable rule codes, path rendering, report ordering,
and projection fields follow semantic versioning after v1. Adding a new rule is
minor; changing an existing pass/fail boundary, code, or path is breaking.
Application prose is not part of the semantic contract.

Run `golib api check` during upgrades and compare `CHANGELOG.md`. Re-run local
truth tables for application-specific optional/null decoding because Go
decoders can collapse states before this package sees them.

When upgrading from an earlier pre-v1 snapshot, initialize limits from
`DefaultLimits` instead of a positional literal. `MaxStringLength` now rejects
oversized typed, reflective, collection-key, and translation input with
`string_limit` before parsing or hashing. Custom diagnostics that violate
severity, code, metadata, UTF-8, or control-character constraints now fail
closed as `invalid_violation`. Custom validator panics become
`validator_panic`, without retaining the panic payload. Application message
catalog output is bounded, control-free, valid UTF-8, and HTML-escaped; compare
machine code and path rather than translated prose in compatibility fixtures.

## Adopting terminal-aware validation and target adapters

Upgrade to v1.1.0 before changing imports. Replace complete-success checks
based on `Report.Empty()` or `!Report.HasErrors()` with `Report.Err() == nil`
where validation can observe a caller context. Classify cancellation and
deadlines with `errors.Is`; use `errors.As` with `*validation.ContextError` to
inspect an immutable partial report. A partial terminal report can also match
`validation.ErrInvalid` and expose `*validation.InvalidError`.

After the public version resolves without `replace` or `go.work`, imports may
move independently:

| Legacy path | Successor path |
| --- | --- |
| `validationconfig` | `adapters/config` |
| `validationhttp` | `adapters/http` |
| `validationjsonapi` | `adapters/jsonapi` |
| `validationrpc` | `adapters/jsonrpc` |
| `validationservice` | `adapters/service` |

Successor named types deliberately have their own package identity. Do not
rely on assignment compatibility between legacy and successor named types;
migrate one integration boundary at a time. Legacy paths remain supported for
the longer of 180 days after successor public availability and two published
stable minor releases. Rolling back an import migration does not restore the
old success-equivalent cancellation behavior.
