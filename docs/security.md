# Security model

## Published v2 diagnostic change

The published v2.0.0 release changes the default error text returned by
`rules.Pattern` and public `structplan` construction/tag compilation. These
errors expose fixed categories rather than caller expressions, field names,
type names, or rule identifiers. Valid matching and plan behavior are unchanged.
Existing sentinel classification through `errors.Is` remains available, and
`errors.As`/`errors.Unwrap` retain the original diagnostic causes for explicit
trusted inspection. Those causes may contain private application configuration;
do not render them into public logs, traces, or responses. Consumers must stop
parsing default diagnostic text and use the structured classification instead.

This describes the published v2 API, not a change to an already published v1
artifact or a claim that the whole package family has completed security work.

## Trust boundaries

Input values, object graphs, map keys, tags, custom validators, translation
catalogs, and transport consumers may be hostile. Applications choose which
parameters and causes are safe; core rules add bounds only and never rejected
values. Paths are locations and can contain caller field/key text, so
`validationobserve` deliberately excludes them from labels.

## Threats and controls

| Threat | Control |
| --- | --- |
| Secret disclosure | values are absent from violations/default formatting; projections omit causes |
| Validation bypass | explicit typed rules, strict tags, mutation-tested branches |
| Path confusion | typed segments, deterministic rendering, RFC 6901 escaping |
| Rule-code collision | bounded machine-safe codes participate in structural dedup identity; invalid custom diagnostics fail closed |
| Log/metric injection | observation labels exclude paths, values, parameters, and causes; invalid custom labels are replaced |
| CPU or memory denial | string/collection/depth/field/tag/path/violation/regex/cache/concurrency limits, including total supplied compositor/service-chain positions before invocation or async result/worker allocation |
| Regex denial | startup compilation with Go RE2 and pattern-length bound |
| Reflection panic/recursion | startup kind checks, inaccessible-field errors, cycle/depth detection |
| Custom panic | function adapters contain panics; sync/async wrappers protect arbitrary implementations; payloads are discarded |
| Hidden I/O/deadlock | I/O uses separate context-aware async contract |
| Cancellation disclosure | reports retain only the standard terminal identity, never a request context or custom cancellation cause |

## Resource budgets

| Limit | Default | Enforcement |
| --- | ---: | --- |
| Depth | 32 | reflective compilation |
| Collection size | 10,000 | item/key/unique traversal and total compositor/service-chain positions (including nil) before work |
| String size | 65,536 bytes | typed and reflective string rules before parsing or comparison |
| Violations | 100 | report add/merge |
| Path length | 1,024 bytes | every report addition |
| Metadata entries | 16 | context construction |
| Metadata key/value | 64/256 bytes | context construction |
| Violation parameters | 16 entries, 64/256-byte key/value | construction and report insertion |
| Locale and operation | 256 bytes each | context construction |
| Regex pattern | 1,024 bytes | pattern construction |
| Custom concurrency | 8 | `AsyncAll` worker pool |
| Struct fields | 256 | typed fields and every reflectively visited field |
| Tag length | 1,024 bytes | tag compilation |
| Cached plans | 256 | instance cache |

Applications should lower limits for smaller protocols. Limits must be
positive. The zero-value `Context` safely uses default limits. Caller maps,
slices, pointers, and objects are read but not mutated.
Custom validators remain application code and can violate the non-mutation
rule; isolate them by ownership and review, not by assuming Go can enforce
purity. `ValidatorFunc` and `AsyncValidatorFunc` contain panics automatically.
Wrap other interface implementations with `IsolatePanics` before direct use.
HTTP and service `Hook` function adapters, in both canonical and retained
packages, use the same core containment owner. Their panic payloads are
discarded; blocking `validator_panic` findings retain `ErrValidatorPanic`
classification, subject to ordinary report diagnostic/path budgets. HTTP
hooks use default validation limits; service hooks preserve the supplied
validation context. Config checks, arbitrary service interfaces, and observers
remain trusted application collaborators; this does not promise their recovery,
preemption, or mutation safety.

Both HTTP `WriteProblem` variants reject statuses outside 100–999 with a fixed
`ErrInvalidViolation`-classified error before accessing the response writer.
Valid statuses preserve their existing JSON and header behavior. Writer and
encoder failures remain caller-owned errors, not recovered panics.

`AsyncAll` joins every admitted callback before return; a callback that ignores
the caller context can still delay the caller and remains application-owned.

The fanout admission published in v2.0.0 applies at execution, using the supplied
validation `Context` (or its safe default limits). `All`, `Any`, `AsyncAll`,
and both service `Chain` variants reject a list exceeding `MaxCollectionSize`
without invoking any validator or retaining partial findings. The refusal is
one blocking violation, normally `collection_limit` with `ErrLimitExceeded`.
Existing report admission still applies: an insufficient diagnostic-code budget
replaces it with rooted `invalid_violation`/`ErrInvalidViolation`; otherwise an
oversized path becomes rooted `path_limit`/`ErrLimitExceeded`. Diagnostic
validation precedes path validation, so `invalid_violation` wins when both
budgets are exceeded. The report's typed error remains `ErrInvalid`.
Async/service caller cancellation
is checked first and remains the terminal outcome. This bounds package-owned
invocation and async state, not the caller's allocation of validator definitions
or the work performed by trusted callback code.

## Residual ownership and review

These records make the existing trust boundaries accountable. They do not
authorize new behavior, waive findings or establish whole-family security
completion. Package maintainers own preserving the documented contracts;
application integrators own the collaborators and disclosure choices below.

### Application callback execution

- Owner: application integrators for callback behavior; package maintainers
  for admission, bounded concurrency and joining admitted work.
- Rationale: callback implementations are application code. Cancellation
  cannot forcibly stop a running callback or enforce its input ownership.
- Mitigation: use reviewed context-aware callbacks, caller deadlines and
  appropriate collection/concurrency limits. Use the documented panic
  wrappers for arbitrary validator implementations. Wrapping does not make
  callbacks preemptible or prevent their mutation of application objects.
- Review condition: callback implementation, ownership, deadline handling,
  concurrency limits or lifecycle requirements change.

### Explicit diagnostic introspection

- Owner: application integrators for diagnostic sinks; package maintainers
  for private default formatting and stable error classification.
- Rationale: explicit error introspection retains original compilation
  causes for trusted diagnosis; these can contain application configuration.
- Mitigation: use default formatting and structured classifications for
  public output. Restrict cause inspection to trusted diagnosis and redact
  application configuration before storing or exporting it.
- Review condition: cause handling, diagnostic destinations, logging policy
  or public error projections change.

### Application-supplied diagnostic content

- Owner: application integrators for parameter, path and catalog content;
  package maintainers for documented bounds, escaping and projections.
- Rationale: bounded and escaped application text can still be confidential.
  The package cannot infer which application strings are safe to disclose.
- Mitigation: supply only public-safe parameters and catalog messages for
  public projections; avoid confidential field/key names in exposed paths.
  Keep private causes out of projections and observation labels. Encoding
  and size limits are not a confidentiality policy.
- Review condition: input schema, parameters, catalog messages, diagnostic
  consumers or application disclosure policy change.

## Reporting vulnerabilities

Do not include a real secret or customer payload in a report. Provide a minimal
synthetic reproducer describing the affected rule, path, and version.
Custom parameter values are bounded safe text but remain application-asserted
message data; the package cannot infer whether arbitrary text is confidential.
Default error formatting omits them and escapes control-bearing paths, while
transport encoders perform their normal JSON escaping.
Translation lookup is panic-safe and cannot alter machine path, code,
severity, ordering, or blocking state. Catalog text exceeding the string
budget, containing invalid UTF-8 or controls, is omitted; accepted text is
HTML-escaped before return.
