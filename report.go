package validation

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Report is an immutable, ordered, deduplicated collection of violations with
// an optional context-terminal state.
type Report struct {
	limits     Limits
	violations []Violation
	keys       map[string]struct{}
	truncated  bool
	hasErrors  bool
	terminal   contextTerminal
}

// ContextReport returns an empty report that snapshots ctx.Err().
// An active context produces no terminal state. Require a nonnil context;
// nil panics on direct use. Only exact standard cancellation/deadline
// identities are retained, never a request context or custom cancellation cause.
//
//nolint:revive // The frozen additive API preserves validation context before caller context.
func ContextReport(validationContext Context, ctx context.Context) Report {
	report := NewReport(validationContext.Limits())
	switch ctx.Err() {
	case context.Canceled:
		report.terminal = contextCanceled
	case context.DeadlineExceeded:
		report.terminal = contextDeadlineExceeded
	}
	return report
}

// NewReport creates an empty report governed by the supplied limits without
// validating/defaulting them. Use Context.Limits or positive DefaultLimits; a
// zero Report can merge into a configured report but is not an admission
// constructor with default retention bounds.
func NewReport(limits Limits) Report {
	return Report{limits: limits, keys: make(map[string]struct{})}
}

// Add returns a copy-on-write report with a unique violation retained within
// limits. Diagnostic validation precedes display-path admission; refusal
// replaces input with fixed invalid_violation or rooted path_limit findings.
// Blocking state is recorded before dedup/retention, so dropped errors still
// block Err. Identity includes typed path/code/severity/sorted parameters, not
// cause identity; causes are borrowed interfaces and are never rendered here.
func (r Report) Add(violation Violation) Report {
	if !validDiagnostic(violation, r.limits) {
		violation = invalidViolation()
	}
	if violation.path.exceedsRenderedLength(r.limits.MaxPathLength) {
		violation = NewViolation(RootPath(), "path_limit", Error, nil,
			ErrLimitExceeded)
	}
	key := violationKey(violation)
	if violation.severity == Error {
		r.hasErrors = true
	}
	if _, exists := r.keys[key]; exists {
		return r
	}
	if len(r.violations) >= r.limits.MaxViolations {
		r.truncated = true
		return r
	}
	r.violations = append(append([]Violation(nil), r.violations...), violation)
	r.keys = cloneSetWith(r.keys, key)
	return r
}

// Merge appends other's retained violations under receiver limits, preserving
// hidden blocking/truncation flags. The receiver's first terminal wins. No
// dropped finding is reconstructed, so Len/HasCode can omit a blocking reason.
func (r Report) Merge(other Report) Report {
	for _, violation := range other.violations {
		r = r.Add(violation)
	}
	if other.truncated {
		r.truncated = true
	}
	if other.hasErrors {
		r.hasErrors = true
	}
	if r.terminal == contextActive {
		r.terminal = other.terminal
	}
	return r
}

// Len returns the number of retained violations.
func (r Report) Len() int { return len(r.violations) }

// Empty reports whether no violations were retained, not complete success.
// A terminal or truncated report can be empty while Err is nonnil.
func (r Report) Empty() bool { return len(r.violations) == 0 }

// ContextError reports the captured context terminal identity, if any.
func (r Report) ContextError() error { return r.terminal.err() }

// Truncated reports whether a violation was omitted by the configured limit.
func (r Report) Truncated() bool { return r.truncated }

// HasErrors reports whether any blocking violation was observed, including one
// omitted because MaxViolations was reached.
func (r Report) HasErrors() bool { return r.hasErrors }

// Violations returns a defensive copy preserving validation order.
func (r Report) Violations() []Violation {
	return append([]Violation(nil), r.violations...)
}

// HasCode reports whether a retained violation has code.
func (r Report) HasCode(code string) bool {
	for _, violation := range r.violations {
		if violation.code == code {
			return true
		}
	}
	return false
}

// Err returns a typed error if validation did not complete or contains a
// blocking violation.
func (r Report) Err() error {
	if r.terminal != contextActive {
		return &ContextError{report: r, terminal: r.terminal}
	}
	if r.hasErrors {
		return &InvalidError{report: r}
	}
	return nil
}

// String returns a cause/parameter-free summary. Active reports use failure-
// oriented prose even when empty or warnings only; it is not a success test.
// Use Err for acceptance and restrict explicit path/parameter/cause projection
// according to application confidentiality policy.
func (r Report) String() string {
	if terminal := r.ContextError(); terminal != nil {
		result := "validation canceled"
		if r.terminal == contextDeadlineExceeded {
			result = "validation deadline exceeded"
		}
		if len(r.violations) > 0 {
			noun := "violations"
			if len(r.violations) == 1 {
				noun = "violation"
			}
			result += fmt.Sprintf(" with %d %s", len(r.violations), noun)
		}
		if r.truncated {
			result += " (truncated)"
		}
		return result
	}
	noun := "violations"
	if len(r.violations) == 1 {
		noun = "violation"
	}
	result := fmt.Sprintf("validation failed with %d %s", len(r.violations), noun)
	if r.truncated {
		result += " (truncated)"
	}
	return result
}

// Keep identity independent of human rendering and map iteration order.
// Length prefixes plus typed segments avoid display collisions and ambiguous
// concatenation; cause identity is deliberately excluded from deduplication.
func violationKey(violation Violation) string {
	keys := make([]string, 0, len(violation.parameters))
	for key := range violation.parameters {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var result strings.Builder
	writeIdentityPart(&result, violation.path.identity())
	writeIdentityPart(&result, violation.code)
	writeIdentityPart(&result, strconv.Itoa(int(violation.severity)))
	for _, key := range keys {
		writeIdentityPart(&result, key)
		writeIdentityPart(&result, violation.parameters[key])
	}
	return result.String()
}

func writeIdentityPart(result *strings.Builder, value string) {
	result.WriteString(strconv.Itoa(len(value)))
	result.WriteByte(':')
	result.WriteString(value)
}

func cloneSetWith(source map[string]struct{}, key string) map[string]struct{} {
	clone := make(map[string]struct{}, len(source))
	for key := range source {
		clone[key] = struct{}{}
	}
	clone[key] = struct{}{}
	return clone
}
