package structplan_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	validation "github.com/faustbrian/go-validation"
	"github.com/faustbrian/go-validation/rules"
	"github.com/faustbrian/go-validation/structplan"
)

func requirePrivateConstructionError(t *testing.T, err, category error, marker string) {
	t.Helper()
	if !errors.Is(err, category) {
		t.Fatalf("construction classification: %v, want %v", err, category)
	}
	for _, message := range []string{err.Error(), fmt.Sprint(err), fmt.Sprintf("%+v", err)} {
		if strings.Contains(message, marker) || !strings.Contains(message, category.Error()) {
			t.Fatalf("default construction diagnostic is not categorical/private: %q", message)
		}
	}
	// Original diagnostic details remain available only after explicit unwrap.
	for cause := errors.Unwrap(err); cause != nil; cause = errors.Unwrap(cause) {
		if strings.Contains(cause.Error(), marker) {
			return
		}
	}
	t.Fatal("original construction diagnostic cause was discarded")
}

func TestTypedPlanConstructionDiagnosticIsPrivate(t *testing.T) {
	const field = "privateTypedFieldMarker"
	builder := structplan.New[string](validation.DefaultLimits())
	accessor := func(value string) string { return value }
	if err := structplan.Add(builder, field, accessor, rules.Prefix("a")); err != nil {
		t.Fatal(err)
	}
	err := structplan.Add(builder, field, accessor, rules.Prefix("b"))
	requirePrivateConstructionError(t, err, structplan.ErrDuplicateField, field)
	plan, err := builder.Compile()
	if err != nil || plan == nil {
		t.Fatalf("duplicate refusal corrupted builder: %v", err)
	}
	if report := plan.Validate(contextFor(t), "a"); !report.Empty() {
		t.Fatalf("duplicate refusal changed original validator: %v", report)
	}
	if report := plan.Validate(contextFor(t), "b"); report.Len() != 1 || !report.HasCode("prefix") {
		t.Fatalf("original field behavior changed: %v", report)
	}
}

type privateCycleMarker struct {
	Next *privateCycleMarker
}

func TestTagPlanConstructionDiagnosticIsPrivate(t *testing.T) {
	for _, test := range []struct {
		name     string
		compile  func() error
		category error
		marker   string
	}{
		{"unknown rule", func() error {
			type value struct {
				Value string `validate:"privateRuleMarker"`
			}
			plan, err := structplan.CompileTags[value](validation.DefaultLimits())
			if plan != nil {
				t.Fatal("unknown rule returned a partial plan")
			}
			return err
		}, structplan.ErrUnknownRule, "privateRuleMarker"},
		{"duplicate rule", func() error {
			type value struct {
				PrivateFieldMarker string `validate:"required,required"`
			}
			plan, err := structplan.CompileTags[value](validation.DefaultLimits())
			if plan != nil {
				t.Fatal("duplicate rule returned a partial plan")
			}
			return err
		}, structplan.ErrDuplicateRule, "PrivateFieldMarker"},
		{"malformed tag", func() error {
			type value struct {
				PrivateFieldMarker string `validate:"required=1"`
			}
			plan, err := structplan.CompileTags[value](validation.DefaultLimits())
			if plan != nil {
				t.Fatal("malformed tag returned a partial plan")
			}
			return err
		}, structplan.ErrInvalidTag, "PrivateFieldMarker"},
		{"inaccessible field", func() error {
			type value struct {
				privateFieldMarker string `validate:"required"`
			}
			plan, err := structplan.CompileTags[value](validation.DefaultLimits())
			if plan != nil {
				t.Fatal("inaccessible field returned a partial plan")
			}
			return err
		}, structplan.ErrInvalidTag, "privateFieldMarker"},
		{"recursive type", func() error {
			plan, err := structplan.CompileTags[privateCycleMarker](validation.DefaultLimits())
			if plan != nil {
				t.Fatal("recursive type returned a partial plan")
			}
			return err
		}, structplan.ErrCycle, "privateCycleMarker"},
	} {
		t.Run(test.name, func(t *testing.T) {
			requirePrivateConstructionError(t, test.compile(), test.category, test.marker)
		})
	}
	type valid struct {
		Value string `validate:"required"`
	}
	plan, err := structplan.CompileTags[valid](validation.DefaultLimits())
	if err != nil || plan == nil {
		t.Fatalf("valid tag plan rejected: %v", err)
	}
	if report := plan.Validate(contextFor(t), valid{Value: "a"}); !report.Empty() {
		t.Fatalf("valid tag behavior changed: %v", report)
	}
	if report := plan.Validate(contextFor(t), valid{}); !report.HasCode("required") {
		t.Fatalf("invalid value classification changed: %v", report)
	}
}
