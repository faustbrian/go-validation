package rules_test

import (
	"math"
	"testing"

	validation "github.com/faustbrian/go-validation/v3"
	"github.com/faustbrian/go-validation/v3/rules"
)

func TestUniqueRetainsGoEqualityForNaN(t *testing.T) {
	ctx := contextFor(t)
	nan := math.NaN()
	if report := rules.Unique[float64]().Validate(ctx, []float64{nan, nan}); !report.Empty() {
		t.Fatalf("NaNs compare unequal: %v", report)
	}
	if report := rules.Unique[float64]().Validate(ctx, []float64{1, 1}); !report.HasCode("unique") {
		t.Fatalf("finite duplicate = %v", report)
	}
	if report := rules.Unique[float64]().Validate(ctx, []float64{0, math.Copysign(0, -1)}); !report.HasCode("unique") {
		t.Fatalf("signed zeros compare equal: %v", report)
	}
}

func TestUniqueCompositeNaNUsesGoEquality(t *testing.T) {
	ctx := contextFor(t)
	type key struct{ Number float64 }
	value := key{math.NaN()}
	if report := rules.Unique[key]().Validate(ctx, []key{value, value}); !report.Empty() {
		t.Fatalf("composite NaNs remain unequal: %v", report)
	}
	array := [1]float64{math.NaN()}
	if report := rules.Unique[[1]float64]().Validate(ctx, [][1]float64{array, array}); !report.Empty() {
		t.Fatalf("array NaNs remain unequal: %v", report)
	}
}

func TestOrderedMapRulesRejectNonreflexiveKeysBeforeCallbacks(t *testing.T) {
	ctx := contextFor(t)
	nan := math.NaN()
	values := map[float64]int{nan: 7}
	values[nan] = 9
	values[1] = 11
	calls := 0
	keys := validation.ValidatorFunc[float64](func(ctx validation.Context, _ float64) validation.Report {
		calls++
		return validation.NewReport(ctx.Limits())
	})
	child := validation.ValidatorFunc[int](func(ctx validation.Context, value int) validation.Report {
		calls++
		if value == 0 {
			t.Error("fabricated zero value")
		}
		return validation.NewReport(ctx.Limits())
	})
	for _, mode := range []validation.Mode{validation.ShortCircuit, validation.CollectAll} {
		for name, report := range map[string]validation.Report{
			"keys":   rules.Keys[float64, int](mode, keys).Validate(ctx, values),
			"values": rules.Values[float64, int](mode, child).Validate(ctx, values),
		} {
			if report.Err() == nil || report.Len() != 1 || !report.HasCode("key") || report.Violations()[0].Path().String() != "" {
				t.Errorf("%s admission = %v", name, report)
			}
		}
	}
	if calls != 0 {
		t.Fatalf("callbacks before key admission: %d", calls)
	}
	type floatingKey float32
	alias := map[floatingKey]int{floatingKey(float32(nan)): 7}
	if report := rules.Values[floatingKey, int](validation.CollectAll, child).Validate(ctx, alias); !report.HasCode("key") {
		t.Fatalf("aliased NaN admission = %v", report)
	}
}

func TestOrderedMapValuesPreserveActualValuesAndFloatOrder(t *testing.T) {
	ctx := contextFor(t)
	var seen []int
	child := validation.ValidatorFunc[int](func(ctx validation.Context, value int) validation.Report {
		seen = append(seen, value)
		return validation.NewReport(ctx.Limits())
	})
	values := map[float64]int{math.Inf(-1): 7, 0: 9, math.Inf(1): 11}
	report := rules.Values[float64, int](validation.CollectAll, child).Validate(ctx, values)
	if !report.Empty() || len(seen) != 3 || seen[0] != 7 || seen[1] != 9 || seen[2] != 11 {
		t.Fatalf("actual values in key order = %v, %v", seen, report)
	}
}

func TestFieldsEqualRetainsGoNaNInequality(t *testing.T) {
	ctx := contextFor(t)
	type pair struct{ Left, Right float64 }
	rule := rules.FieldsEqual("Left", "Right", func(value pair) float64 { return value.Left }, func(value pair) float64 { return value.Right })
	if report := rule.Validate(ctx, pair{math.NaN(), math.NaN()}); !report.HasCode("equal") {
		t.Fatalf("NaN field equality = %v", report)
	}
	if report := rule.Validate(ctx, pair{1, 1}); !report.Empty() {
		t.Fatalf("finite field equality = %v", report)
	}
}

// Typed presence does not traverse object graphs; only reflective normalization
// rejects cyclic dereference chains. A present non-nil pointer remains present.
func TestTypedRequiredDoesNotInventGraphTraversal(t *testing.T) {
	var cycle any
	cycle = &cycle
	report := rules.Required[any]().Validate(contextFor(t), validation.Present(cycle))
	if !report.Empty() {
		t.Fatalf("typed pointer presence = %v", report)
	}
}
