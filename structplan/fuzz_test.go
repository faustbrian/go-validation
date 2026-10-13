package structplan

import (
	"testing"

	validation "github.com/faustbrian/go-validation/v3"
)

type fuzzAlias string

type FuzzEmbedded struct {
	Alias fuzzAlias `validate:"min=1,max=8"`
}

type fuzzShape struct {
	FuzzEmbedded
	Pointer *string        `validate:"email"`
	Dynamic any            `validate:"required"`
	Mapping map[string]int `validate:"max=8"`
	Slice   []int          `validate:"max=8"`
	Array   [2]int         `validate:"min=1,max=2"`
}

type fuzzGeneric[T any] struct {
	Value T `validate:"required"`
}

// FuzzTagGrammarNeverPanics probes bounded raw grammar strings; it does not
// assert acceptance or returned error classifications.
func FuzzTagGrammarNeverPanics(f *testing.F) {
	f.Add("required,min=1,max=5")
	f.Add("unknown,,min=-1")
	f.Fuzz(func(t *testing.T, tag string) {
		if len(tag) > 4096 {
			t.Skip()
		}
		_, _ = parseTag(tag)
	})
}

// FuzzCompiledPlansAcrossSupportedShapes exercises aliases, nil/non-nil
// pointers, interfaces and collection admission on fixed type shapes,
// including finite constructed pointer/interface cycles. Cyclic
// dereference chains must produce field-local cycle findings.
func FuzzCompiledPlansAcrossSupportedShapes(f *testing.F) {
	f.Add("name", "a@example.com", 1, false)
	f.Add("name", "a@example.com", 2, false)
	f.Add("", "bad", 20, true)
	f.Fuzz(func(t *testing.T, text, email string, count int, nilPointer bool) {
		if len(text) > 256 || len(email) > 256 {
			t.Skip()
		}
		count %= 32
		if count < 0 {
			count = -count
		}
		values := make([]int, count)
		mapping := make(map[string]int, count)
		for index := range count {
			mapping[string(rune('a'+index%26))] = index
		}
		var pointer *string
		if !nilPointer {
			pointer = &email
		}
		limits := validation.DefaultLimits()
		limits.MaxCollectionSize = 16
		ctx, err := validation.NewContext(limits)
		if err != nil {
			t.Fatal(err)
		}
		plan, err := CompileTags[fuzzShape](limits)
		if err != nil {
			t.Fatal(err)
		}
		dynamic := any(text)
		cyclic := count%3 == 1
		if cyclic {
			var cycle any
			cycle = &cycle
			dynamic = cycle
		} else if count%3 == 2 {
			for range count {
				next := dynamic
				dynamic = &next
			}
		}
		report := plan.Validate(ctx, fuzzShape{
			FuzzEmbedded: FuzzEmbedded{Alias: fuzzAlias(text)},
			Pointer:      pointer, Dynamic: dynamic, Mapping: mapping, Slice: values,
		})
		if cyclic != report.HasCode("cycle") {
			t.Fatalf("cyclic graph admission = %v", report)
		}
		generic, err := CompileTags[fuzzGeneric[string]](limits)
		if err != nil {
			t.Fatal(err)
		}
		_ = generic.Validate(ctx, fuzzGeneric[string]{Value: text})
	})
}

// FuzzTypedPlanConstructionNeverPanics probes named definitions, nil providers
// and panic accessors under fixed limits without changing persistent builders.
func FuzzTypedPlanConstructionNeverPanics(f *testing.F) {
	f.Add("field", false, false, false)
	f.Add("", true, false, true)
	f.Fuzz(func(t *testing.T, name string, nilAccessor, nilValidator,
		panicAccessor bool,
	) {
		if len(name) > 4096 {
			t.Skip()
		}
		limits := validation.DefaultLimits()
		limits.MaxPathLength = 64
		builder := New[fuzzGeneric[string]](limits)
		var accessor func(fuzzGeneric[string]) string
		if !nilAccessor {
			accessor = func(value fuzzGeneric[string]) string {
				if panicAccessor {
					panic("secret-value")
				}
				return value.Value
			}
		}
		var validator validation.Validator[string]
		if !nilValidator {
			validator = validation.ValidatorFunc[string](func(ctx validation.Context,
				_ string,
			) validation.Report {
				return validation.NewReport(ctx.Limits())
			})
		}
		if err := Add(builder, name, accessor, validator); err != nil {
			return
		}
		plan, err := builder.Compile()
		if err != nil {
			t.Fatal(err)
		}
		_ = plan.Validate(validation.Context{}, fuzzGeneric[string]{Value: "value"})
	})
}
