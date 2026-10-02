package validationservice_test

import (
	"context"
	"errors"
	"testing"

	validation "github.com/faustbrian/go-validation"
	validationservice "github.com/faustbrian/go-validation/validationservice"
)

func TestServiceChainFanoutAdmission(t *testing.T) {
	limits := validation.DefaultLimits()
	limits.MaxCollectionSize = 1
	vctx, err := validation.NewContext(limits)
	if err != nil {
		t.Fatal(err)
	}
	for _, mode := range []validation.Mode{validation.ShortCircuit, validation.CollectAll} {
		calls := 0
		hook := validationservice.Hook[int](func(context.Context, validation.Context, int) validation.Report {
			calls++
			return validation.NewReport(limits)
		})
		report := validationservice.Chain(mode, hook).Validate(context.Background(), vctx, 1)
		if calls != 1 || report.Err() != nil || !report.Empty() {
			t.Fatalf("inclusive one-hook control: calls=%d report=%v", calls, report)
		}
		for _, values := range [][]validationservice.Validator[int]{{hook, hook}, {nil, hook}} {
			calls = 0
			report = validationservice.Chain(mode, values...).Validate(context.Background(), vctx, 1)
			if calls != 0 {
				t.Fatalf("mode %v dispatched refused chain: %d", mode, calls)
			}
			violations := report.Violations()
			if len(violations) != 1 || violations[0].Code() != "collection_limit" ||
				violations[0].Severity() != validation.Error ||
				!errors.Is(violations[0].Cause(), validation.ErrLimitExceeded) ||
				!errors.Is(report.Err(), validation.ErrInvalid) ||
				report.ContextError() != nil || report.Truncated() {
				t.Fatalf("atomic chain refusal: %v %#v", report.Err(), violations)
			}
		}
		canceled, cancel := context.WithCancel(context.Background())
		cancel()
		calls = 0
		report = validationservice.Chain(mode, hook, hook).Validate(canceled, vctx, 1)
		if calls != 0 || !errors.Is(report.Err(), context.Canceled) || report.Len() != 0 || report.HasCode("collection_limit") {
			t.Fatalf("cancellation must precede chain admission: calls=%d report=%v", calls, report.Err())
		}
	}
}

func TestServiceChainRefusalPreservesReportAdmission(t *testing.T) {
	for _, test := range []struct {
		name                   string
		codeBudget, pathBudget int
		code                   string
		cause                  error
	}{
		{"diagnostic", 1, 1024, "invalid_violation", validation.ErrInvalidViolation},
		{"path", 64, 1, "path_limit", validation.ErrLimitExceeded},
		{"both", 1, 1, "invalid_violation", validation.ErrInvalidViolation},
	} {
		t.Run(test.name, func(t *testing.T) {
			limits := validation.DefaultLimits()
			limits.MaxCollectionSize = 1
			limits.MaxMetadataKeyLength, limits.MaxPathLength = test.codeBudget, test.pathBudget
			ctx, err := validation.NewContext(limits)
			if err != nil {
				t.Fatal(err)
			}
			ctx = ctx.WithPath(validation.Field("long"))
			calls := 0
			hook := validationservice.Hook[int](func(context.Context, validation.Context, int) validation.Report {
				calls++
				return validation.NewReport(limits)
			})
			for _, mode := range []validation.Mode{validation.ShortCircuit, validation.CollectAll} {
				report := validationservice.Chain(mode, hook, hook).Validate(context.Background(), ctx, 1)
				violations := report.Violations()
				if calls != 0 || len(violations) != 1 || violations[0].Code() != test.code ||
					violations[0].Severity() != validation.Error || violations[0].Path().String() != "" ||
					!errors.Is(violations[0].Cause(), test.cause) || !errors.Is(report.Err(), validation.ErrInvalid) ||
					report.ContextError() != nil || report.Truncated() {
					t.Fatalf("report admission precedence: calls=%d report=%v violations=%#v", calls, report, violations)
				}
			}
		})
	}
}
