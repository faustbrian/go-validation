# Migration and versioning

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
