package validation_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	validation "github.com/faustbrian/go-validation"
)

func fanoutContext(t *testing.T) validation.Context {
	t.Helper()
	limits := validation.DefaultLimits()
	limits.MaxCollectionSize = 1
	ctx, err := validation.NewContext(limits)
	if err != nil {
		t.Fatal(err)
	}
	return ctx
}

func requireFanoutRefusal(t *testing.T, report validation.Report) {
	t.Helper()
	violations := report.Violations()
	if len(violations) != 1 || violations[0].Code() != "collection_limit" ||
		violations[0].Severity() != validation.Error ||
		!errors.Is(violations[0].Cause(), validation.ErrLimitExceeded) ||
		!errors.Is(report.Err(), validation.ErrInvalid) || report.ContextError() != nil ||
		report.Truncated() {
		t.Fatalf("atomic fanout refusal: %v %#v", report.Err(), violations)
	}
}

func TestCompositionFanoutAdmission(t *testing.T) {
	ctx := fanoutContext(t)
	for _, mode := range []validation.Mode{validation.ShortCircuit, validation.CollectAll} {
		for name, compose := range map[string]func(...validation.Validator[int]) validation.Validator[int]{
			"all": func(values ...validation.Validator[int]) validation.Validator[int] {
				return validation.All(mode, values...)
			},
			"any": func(values ...validation.Validator[int]) validation.Validator[int] {
				return validation.Any(mode, values...)
			},
		} {
			t.Run(name, func(t *testing.T) {
				calls := 0
				validator := validation.ValidatorFunc[int](func(validation.Context, int) validation.Report {
					calls++
					return validation.NewReport(ctx.Limits())
				})
				if report := compose(validator).Validate(ctx, 1); report.Err() != nil || !report.Empty() || calls != 1 {
					t.Fatalf("inclusive one-validator control: calls=%d report=%v", calls, report)
				}
				for _, values := range [][]validation.Validator[int]{{validator, validator}, {nil, validator}} {
					calls = 0
					report := compose(values...).Validate(ctx, 1)
					if calls != 0 {
						t.Fatalf("mode %v dispatched refused fanout: %d", mode, calls)
					}
					requireFanoutRefusal(t, report)
				}
			})
		}
	}
}

func TestAsyncFanoutAdmission(t *testing.T) {
	ctx := fanoutContext(t)
	var calls atomic.Int32
	validator := validation.AsyncValidatorFunc[int](func(context.Context, validation.Context, int) validation.Report {
		calls.Add(1)
		return validation.NewReport(ctx.Limits())
	})
	if report := validation.AsyncAll(context.Background(), ctx, 1, validator); report.Err() != nil || !report.Empty() || calls.Load() != 1 {
		t.Fatalf("inclusive one-validator control: calls=%d report=%v", calls.Load(), report)
	}
	for _, values := range [][]validation.AsyncValidator[int]{{validator, validator}, {nil, validator}} {
		calls.Store(0)
		report := validation.AsyncAll(context.Background(), ctx, 1, values...)
		if calls.Load() != 0 {
			t.Fatalf("dispatched refused async fanout: %d", calls.Load())
		}
		requireFanoutRefusal(t, report)
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	calls.Store(0)
	report := validation.AsyncAll(canceled, ctx, 1, validator, validator)
	if calls.Load() != 0 || !errors.Is(report.Err(), context.Canceled) || report.Len() != 0 || report.HasCode("collection_limit") {
		t.Fatalf("cancellation must precede fanout admission: calls=%d report=%v", calls.Load(), report.Err())
	}
}

func TestFanoutRefusalPreservesReportAdmission(t *testing.T) {
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
			var calls atomic.Int32
			validator := validation.ValidatorFunc[int](func(validation.Context, int) validation.Report {
				calls.Add(1)
				return validation.NewReport(limits)
			})
			async := validation.AsyncValidatorFunc[int](func(context.Context, validation.Context, int) validation.Report {
				calls.Add(1)
				return validation.NewReport(limits)
			})
			for _, mode := range []validation.Mode{validation.ShortCircuit, validation.CollectAll} {
				for _, report := range []validation.Report{
					validation.All(mode, validator, validator).Validate(ctx, 1),
					validation.Any(mode, validator, validator).Validate(ctx, 1),
					validation.AsyncAll(context.Background(), ctx, 1, async, async),
				} {
					violations := report.Violations()
					if calls.Load() != 0 || len(violations) != 1 || violations[0].Code() != test.code ||
						violations[0].Severity() != validation.Error || violations[0].Path().String() != "" ||
						!errors.Is(violations[0].Cause(), test.cause) || !errors.Is(report.Err(), validation.ErrInvalid) ||
						report.ContextError() != nil || report.Truncated() {
						t.Fatalf("report admission precedence: calls=%d report=%v violations=%#v", calls.Load(), report, violations)
					}
				}
			}
		})
	}
}
