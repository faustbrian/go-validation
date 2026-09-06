// Package validationrpc projects reports into JSON-RPC invalid-params errors.
// The data member is a package-owned JSON-RPC 2.0 extension.
//
// Deprecated: use github.com/faustbrian/go-validation/adapters/jsonrpc. This
// path remains supported for the longer of 180 days after successor public
// availability and two published stable minor releases.
package validationrpc

import validation "github.com/faustbrian/go-validation"

// Error is a JSON-RPC error object.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
}

// Data is the stable machine-readable invalid-params payload.
type Data struct {
	Violations []Violation `json:"violations"`
	Truncated  bool        `json:"truncated,omitempty"`
	HasErrors  bool        `json:"has_errors"`
}

// Violation is a safe JSON-RPC validation finding.
type Violation struct {
	Path       string            `json:"path"`
	Code       string            `json:"code"`
	Parameters map[string]string `json:"parameters,omitempty"`
	Severity   string            `json:"severity"`
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
