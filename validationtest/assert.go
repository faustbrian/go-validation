// Package validationtest provides reusable report assertions and conformance
// helpers for consumers.
package validationtest

import validation "github.com/faustbrian/go-validation/v2"

// TestingT is the subset of testing.TB used by these helpers.
type TestingT interface {
	// Helper marks the helper call site without changing the validation result.
	Helper()
	// Fatalf reports a failed assertion; conventional implementations stop the
	// test, but helpers also return after captured terminal-state refusal.
	Fatalf(format string, arguments ...any)
}

// Case is a reusable typed validation fixture.
type Case[T any] struct {
	// Name is a test diagnostic label, not validation configuration.
	Name string
	// Value is a shallow borrowed fixture value.
	Value T
	// Valid requires Empty for positive conformance cases, including rejection of warnings; this differs from RequireValid.
	Valid bool
	// Codes lists required retained codes for negative cases; it is not exact finding equality.
	Codes []string
}

// RequireValid requires complete nonblocking success; warning-only reports
// pass. It refuses captured cancellation/deadline before checking Err.
func RequireValid(t TestingT, report validation.Report) {
	t.Helper()
	if !requireComplete(t, report) {
		return
	}
	if report.Err() != nil {
		t.Fatalf("expected valid report, got %v", report)
	}
}

// RequireCode fails the test unless report contains code.
func RequireCode(t TestingT, report validation.Report, code string) {
	t.Helper()
	if !requireComplete(t, report) {
		return
	}
	if !report.HasCode(code) {
		t.Fatalf("expected violation code %q, got %v", code, report)
	}
}

// RejectMutations requires each supplied value to complete with blocking
// findings. It is a consumer helper, not proof of an exhaustive mutation
// campaign; inputs/validator are borrowed without fanout limits or containment.
func RejectMutations[T any](t TestingT, ctx validation.Context,
	validator validation.Validator[T], mutations []T,
) {
	t.Helper()
	for index, mutation := range mutations {
		report := validator.Validate(ctx, mutation)
		if !requireComplete(t, report) {
			return
		}
		if report.Err() == nil {
			t.Fatalf("mutation %d was accepted", index)
		}
	}
}

// Conformance runs reusable fixtures against a validator.
func Conformance[T any](t TestingT, ctx validation.Context,
	validator validation.Validator[T], cases []Case[T],
) {
	t.Helper()
	for _, current := range cases {
		report := validator.Validate(ctx, current.Value)
		if !requireComplete(t, report) {
			return
		}
		if current.Valid {
			if !report.Empty() {
				t.Fatalf("case %q: expected valid report, got %v",
					current.Name, report)
			}
			continue
		}
		if report.Err() == nil {
			t.Fatalf("case %q: invalid value was accepted", current.Name)
		}
		for _, code := range current.Codes {
			if !report.HasCode(code) {
				t.Fatalf("case %q: expected code %q, got %v",
					current.Name, code, report)
			}
		}
	}
}

func requireComplete(t TestingT, report validation.Report) bool {
	if report.ContextError() == nil {
		return true
	}
	t.Fatalf("validation did not complete: %v", report.Err())
	return false
}
