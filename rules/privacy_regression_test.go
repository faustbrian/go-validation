package rules_test

import (
	"errors"
	"fmt"
	"regexp/syntax"
	"strings"
	"testing"

	validation "github.com/faustbrian/go-validation/v2"
	"github.com/faustbrian/go-validation/v2/rules"
)

func TestPatternCompilationDiagnosticIsPrivate(t *testing.T) {
	const expression = "(?P<privatePatternMarker"
	validator, err := rules.Pattern(expression, validation.DefaultLimits())
	if validator != nil || err == nil {
		t.Fatalf("invalid expression accepted: %v %v", validator, err)
	}
	var cause *syntax.Error
	if !errors.As(err, &cause) || cause.Expr != expression || cause.Code != syntax.ErrInvalidNamedCapture {
		t.Fatalf("original syntax cause unavailable: %v", err)
	}
	for _, message := range []string{err.Error(), fmt.Sprint(err), fmt.Sprintf("%+v", err)} {
		if strings.Contains(message, "privatePatternMarker") || !strings.Contains(message, "compile validation pattern") {
			t.Fatalf("default pattern diagnostic is not categorical/private: %q", message)
		}
	}
	validator, err = rules.Pattern("^a$", validation.DefaultLimits())
	if err != nil || validator == nil {
		t.Fatalf("valid expression rejected: %v", err)
	}
	if report := validator.Validate(contextFor(t), "a"); !report.Empty() {
		t.Fatalf("valid matching changed: %v", report)
	}
	if report := validator.Validate(contextFor(t), "b"); !report.HasCode("pattern") {
		t.Fatalf("nonmatch classification changed: %v", report)
	}
}
