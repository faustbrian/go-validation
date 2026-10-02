package validationservice_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	validation "github.com/faustbrian/go-validation/v2"
	validationservice "github.com/faustbrian/go-validation/v2/validationservice"
)

func TestServiceHookContainsPrivateApplicationPanic(t *testing.T) {
	const marker = "privateHookMarker"
	defer func() {
		if recover() != nil {
			t.Error("application panic escaped function adapter")
		}
	}()
	ctx, err := validation.NewContext(validation.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	ctx = ctx.WithPath(validation.Field("field"))
	hook := validationservice.Hook[int](func(context.Context, validation.Context, int) validation.Report { panic(marker) })
	report := hook.Validate(context.Background(), ctx, 1)
	violations := report.Violations()
	if len(violations) != 1 || violations[0].Code() != "validator_panic" || violations[0].Severity() != validation.Error ||
		violations[0].Path().String() != "field" || !errors.Is(violations[0].Cause(), validation.ErrValidatorPanic) ||
		!errors.Is(report.Err(), validation.ErrInvalid) || strings.Contains(fmt.Sprint(report, report.Err()), marker) {
		t.Fatalf("private panic outcome: report=%v", report)
	}
	calls := 0
	caller, cancel := context.WithCancel(context.Background())
	cancel()
	control := validationservice.Hook[int](func(got context.Context, vctx validation.Context, value int) validation.Report {
		calls++
		if got != caller || vctx.Path().String() != "field" || value != 1 {
			t.Error("hook arguments changed")
		}
		return validation.NewReport(vctx.Limits()).Add(validation.NewViolation(vctx.Path(), "kept", validation.Warning, nil, nil))
	})
	if output := control.Validate(caller, ctx, 1); calls != 1 || output.Len() != 1 || !output.HasCode("kept") || output.Err() != nil {
		t.Fatalf("normal hook result changed: calls=%d report=%v", calls, output)
	}
}
