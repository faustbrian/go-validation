package validation

import (
	"context"
	"errors"
)

var (
	// ErrInvalid marks a non-empty validation result containing errors.
	ErrInvalid = errors.New("validation failed")
	// ErrLimitExceeded marks rejected work that exceeded a configured bound.
	ErrLimitExceeded = errors.New("validation limit exceeded")
	// ErrInvalidLimit marks an invalid Limits configuration.
	ErrInvalidLimit = errors.New("invalid validation limit")
	// ErrValidatorPanic is the safe cause used for an isolated custom panic.
	ErrValidatorPanic = errors.New("validator panicked")
	// ErrInvalidViolation marks an unsafe or malformed custom diagnostic.
	ErrInvalidViolation = errors.New("invalid validation diagnostic")
)

// InvalidError exposes a validation report through errors.As.
type InvalidError struct{ report Report }

// Error returns a value-safe summary.
func (e *InvalidError) Error() string { return e.report.String() }

// Unwrap makes InvalidError compatible with errors.Is and ErrInvalid.
func (e *InvalidError) Unwrap() error { return ErrInvalid }

// Report returns the immutable validation report.
func (e *InvalidError) Report() Report { return e.report }

// ContextError is the structured terminal error returned by Report.Err.
type ContextError struct {
	report   Report
	terminal contextTerminal
}

// Error returns a bounded context-terminal summary.
func (e *ContextError) Error() string {
	return [...]string{
		"validation context not terminated",
		"validation canceled",
		"validation deadline exceeded",
	}[e.terminal]
}

// Unwrap exposes stable context and validation error identities.
func (e *ContextError) Unwrap() []error {
	terminal := e.terminal.err()
	if terminal == nil {
		return nil
	}
	unwrapped := make([]error, 0, 2)
	unwrapped = append(unwrapped, terminal)
	if e.report.hasErrors {
		unwrapped = append(unwrapped, &InvalidError{report: e.report})
	}
	return unwrapped
}

// Report returns the immutable partial validation report.
func (e *ContextError) Report() Report { return e.report }

type contextTerminal uint8

const (
	contextActive contextTerminal = iota
	contextCanceled
	contextDeadlineExceeded
)

func (terminal contextTerminal) err() error {
	return [...]error{nil, context.Canceled, context.DeadlineExceeded}[terminal]
}
