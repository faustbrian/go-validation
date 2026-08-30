# Specification conformance matrix

The [source manifest](manifest.tsv) pins RFC 6901, RFC 9457, JSON:API 1.1,
JSON-RPC 2.0, and their separate errata or release authorities. The
[decision register](../docs/specification-decisions.md) is the canonical human
record; the JSON files in this directory bind the same decisions to executable
evidence and append-only content digests.

The package claims only the behaviors named below. It serializes JSON Pointer
tokens but does not evaluate pointers or implement JSON Patch. Its HTTP shape
is RFC 9457-inspired rather than a complete implementation. JSON:API and
JSON-RPC extension data remain package policy. Generic primitive validators do
not claim conformance to namesake external specifications.

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and
"OPTIONAL" in this document are to be interpreted as described in BCP 14
[RFC2119] [RFC8174] when, and only when, they appear in all capitals, as
shown here.

A source digest change MUST NOT silently alter behavior. Review the changed
publication, errata or release authority, affected decisions, compatibility,
and executable evidence before updating any pin.

## Decision matrix

| Decision | Source | Executable evidence | Interoperability and differential status |
| --- | --- | --- | --- |
| VALIDATION-DEC-001 | RFC 6901 Sections 3 and 4 | `TestPathPreservesTypedSegmentsAndEscapesPointers`; `TestErrorsUseJSONPointersAndStableCodes`; `TestTransportProjectionsPreserveConformanceAndEscapeLocations`; `FuzzPathAndReportSafety`; `FuzzProjectionPaths` | not assessed; no maintained peer harness |
| VALIDATION-DEC-002 | RFC 6901 Sections 3 and 4 | `TestPathPreservesTypedSegmentsAndEscapesPointers`; `FuzzPathAndReportSafety` | not assessed; typed Item is package policy |
| VALIDATION-DEC-003 | JSON:API 1.1 Error Objects | `TestErrorsUseJSONPointersAndStableCodes`; `TestWarningsAndTruncationRemainMachineReadable`; `TestTruncatedBlockingStateSurvivesWarningRetention`; `TestTransportProjectionsPreserveConformanceAndEscapeLocations`; `FuzzProjectionPaths` | not assessed; package metadata has no maintained peer |
| VALIDATION-DEC-004 | JSON-RPC 2.0 Section 5.1 | `TestInvalidParamsPreservesSafeStableData`; `TestInvalidParamsProjectsWarningsAndTruncation`; `TestTransportProjectionsPreserveConformanceAndEscapeLocations` | not assessed; server-defined data has no maintained peer |
| VALIDATION-DEC-005 | RFC 9457 Sections 3, 3.1, and 3.2 | `TestProblemAndWriterAreRouterNeutralAndEscaped`; `TestWarningProblemPreservesSeverityAndTruncation`; `TestTransportProjectionsPreserveConformanceAndEscapeLocations` | not assessed; inspired shape has no maintained peer |
| VALIDATION-DEC-006 | All four pinned sources | `TestTransportProjectionsPreserveConformanceAndEscapeLocations`; transport warning and truncation tests; projection fuzz targets | not assessed; owned cross-projection evidence is not independent peer evidence |

Run the structural checks with the canonical `golib specification check`.
Use `golib specification check --online` to re-fetch the bounded authorities
and verify every reviewed content digest.

[RFC2119]: https://www.rfc-editor.org/rfc/rfc2119
[RFC8174]: https://www.rfc-editor.org/rfc/rfc8174
