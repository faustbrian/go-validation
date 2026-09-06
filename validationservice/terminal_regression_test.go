package validationservice_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	validation "github.com/faustbrian/go-validation"
	"github.com/faustbrian/go-validation/validationservice"
)

func TestChainReturnsTerminalWithoutStartingHooks(t *testing.T) {
	vctx, err := validation.NewContext(validation.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var calls atomic.Int32
	hook := validationservice.Hook[int](func(context.Context, validation.Context, int) validation.Report {
		calls.Add(1)
		return validation.NewReport(vctx.Limits())
	})
	for name, chain := range map[string]validationservice.Validator[int]{
		"nonempty": validationservice.Chain(validation.CollectAll, hook),
		"empty":    validationservice.Chain[int](validation.CollectAll),
		"all nil":  validationservice.Chain[int](validation.CollectAll, nil, nil),
	} {
		report := chain.Validate(ctx, vctx, 1)
		if !errors.Is(report.Err(), context.Canceled) {
			t.Fatalf("%s report=%v err=%v", name, report, report.Err())
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("hooks started = %d", calls.Load())
	}
}

func TestChainStopsAfterCallerOrCallbackTerminalAndKeepsPartialReport(t *testing.T) {
	vctx, err := validation.NewContext(validation.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	for _, mode := range []validation.Mode{validation.ShortCircuit, validation.CollectAll} {
		t.Run(map[validation.Mode]string{validation.ShortCircuit: "short", validation.CollectAll: "collect"}[mode], func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			var laterCalls atomic.Int32
			first := validationservice.Hook[int](func(context.Context, validation.Context, int) validation.Report {
				cancel()
				return validation.NewReport(vctx.Limits()).Add(validation.NewViolation(
					vctx.Path(), "partial", validation.Error, nil, nil,
				))
			})
			later := validationservice.Hook[int](func(context.Context, validation.Context, int) validation.Report {
				laterCalls.Add(1)
				return validation.NewReport(vctx.Limits())
			})
			report := validationservice.Chain(mode, first, nil, later).Validate(ctx, vctx, 1)
			if laterCalls.Load() != 0 || !errors.Is(report.Err(), context.Canceled) ||
				!errors.Is(report.Err(), validation.ErrInvalid) || !report.HasCode("partial") {
				t.Fatalf("calls=%d report=%v err=%v", laterCalls.Load(), report, report.Err())
			}
		})
	}
}

func TestChainStopsOnCarriedTerminalWhileCallerRemainsActive(t *testing.T) {
	vctx, err := validation.NewContext(validation.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	var laterCalls atomic.Int32
	first := validationservice.Hook[int](func(context.Context, validation.Context, int) validation.Report {
		derived, cancel := context.WithCancel(context.Background())
		cancel()
		return validation.AsyncAll[int](derived, vctx, 1)
	})
	later := validationservice.Hook[int](func(context.Context, validation.Context, int) validation.Report {
		laterCalls.Add(1)
		return validation.NewReport(vctx.Limits())
	})
	report := validationservice.Chain(validation.CollectAll, first, later).
		Validate(context.Background(), vctx, 1)
	if laterCalls.Load() != 0 || !errors.Is(report.Err(), context.Canceled) {
		t.Fatalf("calls=%d report=%v err=%v", laterCalls.Load(), report, report.Err())
	}
}
