// Package validationjsonapi projects reports into JSON:API error objects.
// Severity and aggregation metadata are package-owned extensions.
package validationjsonapi

import validation "github.com/faustbrian/go-validation/v3"

// Document is a JSON:API error document with report-level metadata.
type Document struct {
	// Errors holds detached ordered retained findings; it is empty, not omitted, when there are none.
	Errors []Error `json:"errors"`
	// Meta preserves hidden blocking/truncation state independently from retained errors.
	Meta DocumentMeta `json:"meta"`
}

// DocumentMeta preserves aggregation state.
type DocumentMeta struct {
	// Truncated reports omitted findings, not total omitted count.
	Truncated bool `json:"truncated"`
	// HasErrors remains true when blocking findings were dropped by retention limits.
	HasErrors bool `json:"has_errors"`
}

// Error is a JSON:API validation error object.
type Error struct {
	// Status is the per-finding string 422 for blocking errors or 200 for warnings, not a response-writer decision.
	Status string `json:"status"`
	// Code is the stable machine rule identity.
	Code string `json:"code"`
	// Source serializes a caller-owned typed path as a pointer, not proof that the request location exists.
	Source Source `json:"source"`
	// Meta contains copied parameters and package-defined severity extensions.
	Meta ErrorMeta `json:"meta"`
}

// ErrorMeta preserves safe rule parameters and severity.
type ErrorMeta struct {
	// Parameters is a defensive copy, with confidentiality asserted by the application.
	Parameters map[string]string `json:"parameters,omitempty"`
	// Severity is error or warning; this is package-owned JSON:API metadata.
	Severity string `json:"severity"`
}

// Source identifies the invalid document location.
type Source struct {
	// Pointer is RFC6901 token serialization; generic Item emits -, without JSON Patch or pointer-evaluation semantics.
	Pointer string `json:"pointer"`
}

// Errors projects report findings in their deterministic order. Applications
// must route context termination before calling it.
func Errors(report validation.Report) Document {
	violations := report.Violations()
	projected := make([]Error, 0, len(violations))
	for _, violation := range violations {
		status := "422"
		severity := "error"
		if violation.Severity() == validation.Warning {
			status = "200"
			severity = "warning"
		}
		projected = append(projected, Error{
			Status: status, Code: violation.Code(),
			Source: Source{Pointer: violation.Path().JSONPointer()},
			Meta: ErrorMeta{Parameters: violation.Parameters(),
				Severity: severity},
		})
	}
	return Document{Errors: projected,
		Meta: DocumentMeta{Truncated: report.Truncated(),
			HasErrors: report.HasErrors()}}
}
