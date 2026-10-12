// Package validationobserve emits findings-only code/severity/operation
// observations through Report. It omits paths, parameters, causes and values,
// but syntactically valid application labels can still be confidential or
// high-cardinality. Observers are borrowed synchronous callbacks with caller-
// owned concurrency, lifetime and latency; no provider goroutines are started.
package validationobserve

import validation "github.com/faustbrian/go-validation/v2"

// Observation contains bounded labels only. It intentionally excludes paths,
// parameters, causes, and rejected values.
type Observation struct {
	// Code is a retained machine code; application-defined identities can still be confidential/high-cardinality.
	Code string
	// Severity is error or warning.
	Severity string
	// Operation is caller Context operation, with non-label characters replaced; syntax checks are not privacy classification.
	Operation string
}

// Observer receives bounded machine labels, not confidentiality-classified data.
type Observer interface {
	// Observe runs synchronously, once per retained finding; no panic/latency
	// containment, nil guard or aggregate truncation/terminal observation is
	// added. Caller owns observer lifetime, privacy and concurrent safety.
	Observe(Observation)
}

// Report emits one observation per retained violation, excluding omitted
// findings/terminal state. Nil observer is not guarded; observer panics/blocking
// propagate. Operation syntax screening is not automatic secrecy protection.
func Report(ctx validation.Context, report validation.Report, observer Observer) {
	for _, violation := range report.Violations() {
		severity := "error"
		if violation.Severity() == validation.Warning {
			severity = "warning"
		}
		operation := safeLabel(ctx.Operation(), "invalid_operation")
		observer.Observe(Observation{Code: violation.Code(), Severity: severity,
			Operation: operation})
	}
}

func safeLabel(value string, replacement string) string {
	if len(value) == 0 {
		return value
	}
	for _, character := range value {
		if (character >= 'a' && character <= 'z') ||
			(character >= 'A' && character <= 'Z') ||
			(character >= '0' && character <= '9') ||
			character == '_' || character == '-' || character == '.' ||
			character == ':' {
			continue
		}
		return replacement
	}
	return value
}
