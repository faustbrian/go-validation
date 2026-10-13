package validation_test

import (
	"testing"

	validation "github.com/faustbrian/go-validation/v3"
)

func TestAnyRequiresAnEvaluatedSuccessfulAlternative(t *testing.T) {
	ctx := testContext(t)
	for _, mode := range []validation.Mode{validation.ShortCircuit, validation.CollectAll} {
		for _, children := range [][]validation.Validator[int]{nil, {}, {nil}, {nil, nil}} {
			report := validation.Any(mode, children...).Validate(ctx, 1)
			if !report.HasCode("any") || report.Err() == nil {
				t.Fatalf("empty/nil Any must fail without a successful alternative: %v", report)
			}
		}
		failure := validation.ValidatorFunc[int](func(ctx validation.Context, _ int) validation.Report {
			return validation.NewReport(ctx.Limits()).Add(validation.NewViolation(ctx.Path(), "failed", validation.Error, nil, nil))
		})
		if report := validation.Any(mode, nil, failure, nil).Validate(ctx, 1); !report.HasCode("failed") {
			t.Fatalf("nil does not satisfy an evaluated alternative: %v", report)
		}
		limits := validation.DefaultLimits()
		limits.MaxCollectionSize = 1
		bounded, err := validation.NewContext(limits)
		if err != nil {
			t.Fatal(err)
		}
		if report := validation.Any[int](mode, nil, nil).Validate(bounded, 1); !report.HasCode("collection_limit") {
			t.Fatalf("nil positions bypass admission: %v", report)
		}
	}
}
