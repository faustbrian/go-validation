// Package validationhttp provides router-neutral HTTP report projection.
//
// Deprecated: use github.com/faustbrian/go-validation/v2/adapters/http. This
// path remains supported for the longer of 180 days after successor public
// availability and two published stable minor releases.
package validationhttp

import (
	"encoding/json"
	"fmt"
	"net/http"

	validation "github.com/faustbrian/go-validation/v2"
)

// Problem is RFC 9457-inspired, with package extensions; it is not complete.
type Problem struct {
	// Type is an application-overridable problem URI; FromReport uses the package fixed type.
	Type string `json:"type"`
	// Title is mutable public prose; FromReport selects warnings or failure text.
	Title string `json:"title"`
	// Status is the projected 200/422 status; WriteProblem accepts manually supplied statuses only within 100..999.
	Status int `json:"status"`
	// Errors is a detached ordered finding slice; causes/terminal context are not included.
	Errors []Error `json:"errors"`
	// Truncated reports whether findings were omitted by MaxViolations.
	Truncated bool `json:"truncated,omitempty"`
}

// Error is one safe HTTP validation finding.
type Error struct {
	// Path is the human display path, which may contain application names.
	Path string `json:"path"`
	// Code is the stable machine identity.
	Code string `json:"code"`
	// Parameters is a defensive copy of application-asserted public-safe rule parameters.
	Parameters map[string]string `json:"parameters,omitempty"`
	// Severity is error or warning from the source finding.
	Severity string `json:"severity"`
}

// FromReport projects validation findings into an unprocessable-content
// problem. Applications must route context termination before calling it.
func FromReport(report validation.Report) Problem {
	violations := report.Violations()
	errors := make([]Error, 0, len(violations))
	for _, violation := range violations {
		errors = append(errors, Error{Path: violation.Path().String(),
			Code: violation.Code(), Parameters: violation.Parameters(),
			Severity: severity(violation.Severity())})
	}
	status := http.StatusOK
	title := "Validation warnings"
	if report.HasErrors() {
		status = http.StatusUnprocessableEntity
		title = "Validation failed"
	}
	return Problem{Type: "https://validation.invalid/problem", Title: title,
		Status: status, Errors: errors, Truncated: report.Truncated()}
}

func severity(value validation.Severity) string {
	if value == validation.Warning {
		return "warning"
	}
	return "error"
}

// WriteProblem writes caller-owned problem data without routing or binding.
// Invalid status rejects before touching writer. For valid status it sets the
// content type and writes headers before encoding; failure can leave a partial
// response. Writer/encoder panics are not contained; encoder errors are returned
// unchanged.
// Nil writer rejects only through ordinary Go panic, not a typed check.
func WriteProblem(writer http.ResponseWriter, problem Problem) error {
	if problem.Status < 100 || problem.Status > 999 {
		return fmt.Errorf("%w: HTTP problem status", validation.ErrInvalidViolation)
	}
	writer.Header().Set("Content-Type", "application/problem+json")
	writer.WriteHeader(problem.Status)
	return json.NewEncoder(writer).Encode(problem)
}

// Hook is a transport integration seam that does not own routing or binding.
type Hook[T any] func(*http.Request, T) validation.Report

// Validate borrows request/value and invokes the hook synchronously with core
// panic containment/default validation limits. Nil hooks produce panic findings;
// nil requests are forwarded. Mutation, request-body ownership, cancellation
// handling and latency remain application obligations; no body is read/closed.
func (hook Hook[T]) Validate(request *http.Request, value T) validation.Report {
	return validation.ValidatorFunc[T](func(_ validation.Context, value T) validation.Report {
		return hook(request, value)
	}).Validate(validation.Context{}, value)
}
