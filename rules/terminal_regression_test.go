package rules_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	validation "github.com/faustbrian/go-validation"
	"github.com/faustbrian/go-validation/rules"
)

func terminalValidator[T any](ctx validation.Context, calls *atomic.Int32) validation.Validator[T] {
	return validation.ValidatorFunc[T](func(current validation.Context, _ T) validation.Report {
		if calls.Add(1) == 1 {
			caller, cancel := context.WithCancel(context.Background())
			cancel()
			return validation.AsyncAll[T](caller, ctx, *new(T)).Add(
				validation.NewViolation(current.Path(), "partial", validation.Warning, nil, nil),
			)
		}
		return validation.NewReport(ctx.Limits()).Add(validation.NewViolation(
			ctx.Path(), "later", validation.Error, nil, nil,
		))
	})
}

func TestCollectionCompositorsStopOnTerminalReports(t *testing.T) {
	ctx := contextFor(t)
	var itemCalls atomic.Int32
	items := rules.Items(validation.CollectAll, terminalValidator[int](ctx, &itemCalls)).
		Validate(ctx, []int{1, 2})
	if !errors.Is(items.Err(), context.Canceled) || itemCalls.Load() != 1 ||
		!items.HasCode("partial") || items.Violations()[0].Path().String() != "[0]" {
		t.Fatalf("Items calls=%d report=%v err=%v", itemCalls.Load(), items, items.Err())
	}

	var keyCalls atomic.Int32
	keys := rules.Keys[string, int](validation.CollectAll,
		terminalValidator[string](ctx, &keyCalls)).Validate(ctx, map[string]int{"a": 1, "b": 2})
	if !errors.Is(keys.Err(), context.Canceled) || keyCalls.Load() != 1 ||
		keys.Violations()[0].Path().String() != "[a]" {
		t.Fatalf("Keys calls=%d report=%v err=%v", keyCalls.Load(), keys, keys.Err())
	}

	var valueCalls atomic.Int32
	values := rules.Values[string, int](validation.CollectAll,
		terminalValidator[int](ctx, &valueCalls)).Validate(ctx, map[string]int{"a": 1, "b": 2})
	if !errors.Is(values.Err(), context.Canceled) || valueCalls.Load() != 1 ||
		values.Violations()[0].Path().String() != "[a]" {
		t.Fatalf("Values calls=%d report=%v err=%v", valueCalls.Load(), values, values.Err())
	}
}
