# Specification decisions

This register bounds the package's standards claims and records every observable
interpretation owned by the validation paths and transport projections.

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and
"OPTIONAL" in this document are to be interpreted as described in BCP 14
[RFC2119] [RFC8174] when, and only when, they appear in all capitals, as
shown here.

Source pins, update monitoring, conformance bindings, and append-only decision
history live under the [specification directory](../specification/README.md).

The package does not claim external-specification conformance for the generic
`URL`, `Email`, `UUID`, hostname, or IP rules. Those names identify
package-defined bounded syntactic profiles documented in the
[rule catalog](rules.md).

## VALIDATION-DEC-001: RFC 6901 pointer serialization

Status `resolved`; owner `validation maintainers`; classification `interoperability policy`;
decision scope `normative`; specification `RFC 6901 JSON Pointer`;
version `RFC 6901`; source authority `rfc6901-source`; authority URL
https://www.rfc-editor.org/rfc/rfc6901.txt; section `Sections 3 and 4`; requirement strength
`not specified`.

| Field | Decision |
| --- | --- |
| Issue | Typed validation paths need one stable JSON Pointer serialization without losing segment boundaries or mis-escaping reference tokens. |
| Credible interpretations | `Render each typed segment as one reference token and represent the root as the empty string.`; `Flatten the human-readable path and treat that rendering as one pointer token.` |
| Known peer behavior | RFC 6901 implementations serialize the root as an empty pointer and escape each reference token; no maintained peer harness is present in this repository. |
| Selected behavior | JSONPointer returns the empty string for RootPath and emits one slash-prefixed reference token per typed segment, replacing each tilde with ~0 before replacing each slash with ~1. |
| Rationale | This preserves segment identity and follows the RFC 6901 string syntax without coupling pointer serialization to the human-readable path. |
| Security consequences | Hostile field and key text cannot inject additional pointer segments through slash or tilde characters. |
| Resource consequences | Serialization is linear in the bounded path content and performs no I/O. |
| Compatibility consequences | Root representation, segment boundaries, and escape order are part of the v1 compatibility contract. |
| Wire consequences | JSON:API source pointers use the RFC 6901 string form; human-readable paths remain a separate representation. |
| Executable evidence | `TestPathPreservesTypedSegmentsAndEscapesPointers`; `TestErrorsUseJSONPointersAndStableCodes`; `TestTransportProjectionsPreserveConformanceAndEscapeLocations` |
| Fixture evidence | None. |
| Fuzz evidence | `FuzzPathAndReportSafety`; `FuzzProjectionPaths` |
| Interoperability evidence | None. |
| Differential evidence | None. |
| Public APIs | `RootPath`; `Path.JSONPointer`; `validationjsonapi.Errors` |
| Documentation | `docs/specification-decisions.md`; `docs/semantics.md`; `docs/api.md` |
| Upstream status | RFC 6901 errata are monitored independently from the pinned normative text. |
| Reconsider when | RFC 6901 is superseded or a supported transport requires a different pointer representation. |

## VALIDATION-DEC-002: Typed item token without append semantics

Status `resolved`; owner `validation maintainers`; classification `omission`;
decision scope `application-policy`; specification `RFC 6901 JSON Pointer`;
version `RFC 6901`; source authority `rfc6901-source`; authority URL
https://www.rfc-editor.org/rfc/rfc6901.txt; section `Sections 3 and 4`; requirement strength
`not specified`.

| Field | Decision |
| --- | --- |
| Issue | The package-owned Item segment has no concrete object member or array index, while RFC 6901 assigns the literal - token an array-evaluation result that does not identify an existing value. |
| Credible interpretations | `Reject Item paths as non-serializable.`; `Serialize Item as a literal - reference token while leaving evaluation and mutation semantics to consumers.` |
| Known peer behavior | JSON Pointer and JSON Patch implementations interpret - only in their own evaluation or mutation context; this package has no maintained peer harness for its typed Item abstraction. |
| Selected behavior | JSONPointer serializes Item as the literal - reference token but neither evaluates that token nor claims JSON Patch append semantics. |
| Rationale | The stable token preserves a generic collection-item location without making the validation package a pointer evaluator or patch processor. |
| Security consequences | Consumers cannot infer that an Item pointer authorizes an append operation; that interpretation remains outside this API. |
| Resource consequences | Item serialization writes one bounded token and performs no collection lookup. |
| Compatibility consequences | The literal /- rendering remains stable, but consumers must choose any evaluation or mutation policy independently. |
| Wire consequences | An Item segment appears as /- in pointer-bearing projections and as [] in the human-readable path. |
| Executable evidence | `TestPathPreservesTypedSegmentsAndEscapesPointers` |
| Fixture evidence | None. |
| Fuzz evidence | `FuzzPathAndReportSafety` |
| Interoperability evidence | None. |
| Differential evidence | None. |
| Public APIs | `Item`; `Path.JSONPointer`; `Path.String` |
| Documentation | `docs/specification-decisions.md`; `docs/semantics.md`; `docs/api.md` |
| Upstream status | No upstream issue is needed because the mapping from the package-owned Item type is application policy. |
| Reconsider when | The package adds JSON Pointer evaluation, JSON Patch operations, or a transport that forbids a non-resolving pointer. |

## VALIDATION-DEC-003: JSON:API validation error projection

Status `resolved`; owner `validation maintainers`; classification `optional behavior`;
decision scope `transport-specific`; specification `JSON:API 1.1`;
version `JSON:API 1.1`; source authority `jsonapi-source`; authority URL
https://jsonapi.org/format/1.1/; section `Document Structure, Errors, Error Objects, and Meta Information`; requirement strength
`MAY`.

| Field | Decision |
| --- | --- |
| Issue | JSON:API permits several error members and arbitrary meta members but does not define this package's severity, truncation, safe-parameter, or report-level blocking representation. |
| Credible interpretations | `Emit only standard error members and discard package aggregation state.`; `Use JSON:API error objects and source pointers while carrying package-owned state in meta objects.` |
| Known peer behavior | JSON:API implementations expose different application-specific error metadata; no maintained peer implements this package's Report extensions. |
| Selected behavior | Errors emits an errors array whose objects contain string status, stable code, RFC 6901 source pointer, and package-owned parameters and severity meta; top-level meta contains truncated and has_errors. Callers remain responsible for ensuring each pointer identifies an existing value in the request document. |
| Rationale | The standard members preserve JSON:API interoperability while meta keeps bounded validation state machine-readable without inventing standard members. |
| Security consequences | Only safe parameters and escaped pointer text are projected; causes and rejected values remain excluded. |
| Resource consequences | Projection is linear in the report's bounded retained violations and performs no I/O. |
| Compatibility consequences | Error object field names, status strings, meta field names, severity values, ordering, and truncation semantics are part of the v1 transport contract. |
| Wire consequences | Errors use status 422 and warnings use status 200; source.pointer is RFC 6901 syntax, while meta members are package extensions rather than JSON:API-defined semantics. An Item segment serializes as /- and does not identify an existing array value. |
| Executable evidence | `TestErrorsUseJSONPointersAndStableCodes`; `TestWarningsAndTruncationRemainMachineReadable`; `TestTruncatedBlockingStateSurvivesWarningRetention`; `TestTransportProjectionsPreserveConformanceAndEscapeLocations` |
| Fixture evidence | None. |
| Fuzz evidence | `FuzzProjectionPaths` |
| Interoperability evidence | None. |
| Differential evidence | None. |
| Public APIs | `validationjsonapi.Document`; `validationjsonapi.Error`; `validationjsonapi.Errors` |
| Documentation | `docs/specification-decisions.md`; `docs/guides.md`; `docs/api.md` |
| Upstream status | The mutable official JSON:API 1.1 publication page is monitored separately as a release authority. |
| Reconsider when | JSON:API publishes a new stable version, changes error-object requirements, or defines conflicting metadata semantics. |

## VALIDATION-DEC-004: JSON-RPC invalid-params projection

Status `resolved`; owner `validation maintainers`; classification `optional behavior`;
decision scope `transport-specific`; specification `JSON-RPC 2.0`;
version `JSON-RPC 2.0`; source authority `jsonrpc-source`; authority URL
https://www.jsonrpc.org/specification; section `Section 5.1 Error object and reserved error codes`; requirement strength
`not specified`.

| Field | Decision |
| --- | --- |
| Issue | JSON-RPC 2.0 defines -32602 Invalid params and allows server-defined data, but it does not define a validation finding schema. |
| Credible interpretations | `Return only the reserved code and message.`; `Attach a package-owned bounded violations payload in the optional data member.` |
| Known peer behavior | JSON-RPC servers use application-specific data shapes for Invalid params; no maintained peer implements this package's Report projection. |
| Selected behavior | InvalidParams returns code -32602, message Invalid params, and package-owned data containing ordered violations, truncation, report-level has_errors, safe parameters, and severity. |
| Rationale | The reserved code preserves protocol meaning while server-defined data retains stable validation details without claiming those fields are standardized. |
| Security consequences | The projection excludes causes and rejected values and retains only bounded safe report data. |
| Resource consequences | Projection is linear in the report's bounded retained violations and performs no transport or request processing. |
| Compatibility consequences | The code, message, data field names, finding order, and severity values are part of the v1 transport contract. |
| Wire consequences | The error object uses the JSON-RPC 2.0 Invalid params code; the data object is package policy and callers remain responsible for the surrounding response envelope and request id. |
| Executable evidence | `TestInvalidParamsPreservesSafeStableData`; `TestInvalidParamsProjectsWarningsAndTruncation`; `TestTransportProjectionsPreserveConformanceAndEscapeLocations` |
| Fixture evidence | None. |
| Fuzz evidence | None. |
| Interoperability evidence | None. |
| Differential evidence | None. |
| Public APIs | `validationrpc.Error`; `validationrpc.Data`; `validationrpc.InvalidParams` |
| Documentation | `docs/specification-decisions.md`; `docs/guides.md`; `docs/api.md` |
| Upstream status | The mutable official JSON-RPC 2.0 publication page is monitored separately as a release authority. |
| Reconsider when | JSON-RPC publishes a successor or an adopted profile defines an incompatible Invalid params data schema. |

## VALIDATION-DEC-005: RFC 9457-inspired HTTP problem projection

Status `resolved`; owner `validation maintainers`; classification `optional behavior`;
decision scope `transport-specific`; specification `RFC 9457 Problem Details for HTTP APIs`;
version `RFC 9457`; source authority `rfc9457-source`; authority URL
https://www.rfc-editor.org/rfc/rfc9457.txt; section `Sections 3, 3.1, and 3.2`; requirement strength
`MAY`.

| Field | Decision |
| --- | --- |
| Issue | RFC 9457 permits extension members but does not define this package's warning success response, validation error array, truncation member, or fixed problem type. |
| Credible interpretations | `Claim complete RFC 9457 problem-details conformance.`; `Expose a deliberately RFC 9457-inspired projection with synchronized status and package-owned extension fields.` |
| Known peer behavior | Problem Details implementations define application-specific validation extensions; no maintained peer implements this package's warning and truncation policy. |
| Selected behavior | FromReport and WriteProblem expose an RFC 9457-inspired shape with type, title, synchronized numeric status, errors, and optional truncated; the package does not claim complete RFC 9457 compliance. |
| Rationale | The familiar media type and core member names provide a router-neutral integration seam while explicit claim limits prevent extensions and warning behavior from being mistaken for a full Problem Details implementation. |
| Security consequences | Only safe parameters and bounded human-readable paths are emitted; causes and rejected values remain excluded. |
| Resource consequences | Projection and encoding are linear in the bounded report and use the caller-owned response writer. |
| Compatibility consequences | The fixed type URI, titles, status mapping, media type, and extension field names are part of the v1 transport contract but not a general RFC 9457 implementation promise. |
| Wire consequences | WriteProblem emits application/problem+json and the same HTTP status as Problem.status; warnings produce 200 with Validation warnings, while blocking reports produce 422 with Validation failed. |
| Executable evidence | `TestProblemAndWriterAreRouterNeutralAndEscaped`; `TestWarningProblemPreservesSeverityAndTruncation`; `TestTransportProjectionsPreserveConformanceAndEscapeLocations` |
| Fixture evidence | None. |
| Fuzz evidence | None. |
| Interoperability evidence | None. |
| Differential evidence | None. |
| Public APIs | `validationhttp.Problem`; `validationhttp.FromReport`; `validationhttp.WriteProblem` |
| Documentation | `docs/specification-decisions.md`; `docs/guides.md`; `docs/api.md` |
| Upstream status | RFC 9457 errata are monitored independently from the pinned normative text. |
| Reconsider when | RFC 9457 is superseded, the package claims complete Problem Details conformance, or the HTTP projection changes its status or extension policy. |

## VALIDATION-DEC-006: Shared projection state and escaping policy

Status `resolved`; owner `validation maintainers`; classification `omission`;
decision scope `application-policy`; specification `JSON:API 1.1`;
version `JSON:API 1.1`; source authority `jsonapi-source`; authority URL
https://jsonapi.org/format/1.1/; section `Errors, Error Objects, and Meta Information`; requirement strength
`not specified`.

Additional authoritative source: `{"id":"rfc6901-source","version":"RFC 6901","url":"https://www.rfc-editor.org/rfc/rfc6901.txt","specifications":["RFC 6901 JSON Pointer"]}`

Additional authoritative source: `{"id":"rfc9457-source","version":"RFC 9457","url":"https://www.rfc-editor.org/rfc/rfc9457.txt","specifications":["RFC 9457 Problem Details for HTTP APIs"]}`

Additional authoritative source: `{"id":"jsonrpc-source","version":"JSON-RPC 2.0","url":"https://www.jsonrpc.org/specification","specifications":["JSON-RPC 2.0"]}`

| Field | Decision |
| --- | --- |
| Issue | The referenced transport specifications do not jointly define how one bounded Report preserves ordering, severity, truncation, blocking state, safe parameters, and hostile location text across projections. |
| Credible interpretations | `Let each transport independently select and potentially lose report state.`; `Apply one package-owned projection policy and vary only the transport representation.` |
| Known peer behavior | The cross-projection tests compare package-owned transports, not independent peers; no practical maintained peer exposes the same package-specific Report state. |
| Selected behavior | All projections preserve retained violation order, stable codes, safe parameters, severity, truncation, and report-level blocking state where their public shape supports it, and JSON encoding escapes hostile location text. |
| Rationale | One application policy prevents transport adapters from silently changing validation meaning while keeping specification-defined and package-defined fields distinct. |
| Security consequences | Hostile location text is JSON-escaped and rejected values, causes, and unsafe data remain outside every projection. |
| Resource consequences | Every projection processes only bounded retained violations and performs no hidden network or storage work. |
| Compatibility consequences | Changing ordering, severity, truncation, blocking state, safe parameters, or escaping requires compatibility and specification-decision review across every transport. |
| Wire consequences | Transport field names differ intentionally, but retained findings and aggregation state remain semantically aligned; HTTP does not expose has_errors as a separate member because status carries the blocking state. |
| Executable evidence | `TestTransportProjectionsPreserveConformanceAndEscapeLocations`; `TestWarningProblemPreservesSeverityAndTruncation`; `TestWarningsAndTruncationRemainMachineReadable`; `TestInvalidParamsProjectsWarningsAndTruncation` |
| Fixture evidence | None. |
| Fuzz evidence | `FuzzPathAndReportSafety`; `FuzzProjectionPaths` |
| Interoperability evidence | None. |
| Differential evidence | None. |
| Public APIs | `validationhttp.FromReport`; `validationjsonapi.Errors`; `validationrpc.InvalidParams` |
| Documentation | `docs/specification-decisions.md`; `docs/guides.md`; `docs/compatibility.md` |
| Upstream status | No upstream source defines the combined package policy; each pinned specification and its change authority is monitored independently. |
| Reconsider when | A supported transport specification defines conflicting aggregation semantics or a maintained peer with the same report model becomes available. |

[RFC2119]: https://www.rfc-editor.org/rfc/rfc2119
[RFC8174]: https://www.rfc-editor.org/rfc/rfc8174
