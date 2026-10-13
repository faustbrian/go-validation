// Package validationjsonrpc projects reports into JSON-RPC invalid-params
// errors. The data member is a package-owned JSON-RPC 2.0 extension.
package validationjsonrpc

import validation "github.com/faustbrian/go-validation/v3"

// Error is a JSON-RPC error object.
type Error struct {
	// Code is -32602 from InvalidParams, including warning-only projected reports.
	Code int `json:"code"`
	// Message is fixed Invalid params from the projector, not localized rule prose.
	Message string `json:"message"`
	// Data is package-defined extension data, not a complete JSON-RPC envelope.
	Data Data `json:"data"`
}

// Data is the stable machine-readable invalid-params payload.
type Data struct {
	// Violations is a detached ordered finding slice, excluding causes/terminal context.
	Violations []Violation `json:"violations"`
	// Truncated marks report retention loss; omitted when false in JSON.
	Truncated bool `json:"truncated,omitempty"`
	// HasErrors preserves blocking state even when retained findings are warnings only.
	HasErrors bool `json:"has_errors"`
}

// Violation is a safe JSON-RPC validation finding.
type Violation struct {
	// Path is the human display path rather than JSON Pointer; application names can be sensitive.
	Path string `json:"path"`
	// Code is the machine rule identity.
	Code string `json:"code"`
	// Parameters is copied application-asserted public-safe data, not a redaction guarantee.
	Parameters map[string]string `json:"parameters,omitempty"`
	// Severity is error or warning from the source finding.
	Severity string `json:"severity"`
}

// InvalidParams projects findings using the standard -32602 error code.
// Applications must route context termination before calling it.
func InvalidParams(report validation.Report) Error {
	violations := report.Violations()
	projected := make([]Violation, 0, len(violations))
	for _, violation := range violations {
		projected = append(projected, Violation{
			Path: violation.Path().String(), Code: violation.Code(),
			Parameters: violation.Parameters(), Severity: severity(violation.Severity()),
		})
	}
	return Error{Code: -32602, Message: "Invalid params",
		Data: Data{Violations: projected, Truncated: report.Truncated(),
			HasErrors: report.HasErrors()}}
}

func severity(value validation.Severity) string {
	if value == validation.Warning {
		return "warning"
	}
	return "error"
}
