package structplan_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	validation "github.com/faustbrian/go-validation"
	"github.com/faustbrian/go-validation/structplan"
)

func TestTypedPlanStopsAfterTerminalField(t *testing.T) {
	ctx := contextFor(t)
	builder := structplan.New[account](ctx.Limits())
	warning := validation.ValidatorFunc[string](func(current validation.Context, _ string) validation.Report {
		return validation.NewReport(ctx.Limits()).Add(validation.NewViolation(
			current.Path(), "before", validation.Warning, nil, nil,
		))
	})
	if err := structplan.Add(builder, "name", func(value account) string { return value.Name }, warning); err != nil {
		t.Fatal(err)
	}
	terminal := validation.ValidatorFunc[string](func(current validation.Context, _ string) validation.Report {
		caller, cancel := context.WithCancel(context.Background())
		cancel()
		return validation.AsyncAll[string](caller, ctx, "").Add(validation.NewViolation(
			current.Path(), "terminal", validation.Warning, nil, nil,
		))
	})
	if err := structplan.Add(builder, "email", func(value account) string { return value.Email }, terminal); err != nil {
		t.Fatal(err)
	}
	var laterCalls atomic.Int32
	later := validation.ValidatorFunc[string](func(validation.Context, string) validation.Report {
		laterCalls.Add(1)
		return validation.NewReport(ctx.Limits())
	})
	if err := structplan.Add(builder, "later", func(value account) string { return value.Email }, later); err != nil {
		t.Fatal(err)
	}
	plan, err := builder.Compile()
	if err != nil {
		t.Fatal(err)
	}
	report := plan.Validate(ctx, account{})
	if !errors.Is(report.Err(), context.Canceled) || laterCalls.Load() != 0 {
		t.Fatalf("calls=%d report=%v err=%v", laterCalls.Load(), report, report.Err())
	}
	violations := report.Violations()
	if len(violations) != 2 || violations[0].Path().String() != "name" ||
		violations[1].Path().String() != "email" {
		t.Fatalf("violations = %#v", violations)
	}
}
