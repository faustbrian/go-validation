package validationservice_test

import (
	"context"
	"errors"
	"testing"
	"time"

	validation "github.com/faustbrian/go-validation"
	validationservice "github.com/faustbrian/go-validation/adapters/service"
)

func TestChainPreservesModesNilHooksAndTerminalOutcomes(t *testing.T) {
	vctx, err := validation.NewContext(validation.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	fail := validationservice.Hook[int](func(context.Context, validation.Context, int) validation.Report {
		calls++
		return validation.NewReport(vctx.Limits()).Add(
			validation.NewViolation(vctx.Path(), "stop", validation.Error, nil, nil),
		)
	})
	after := validationservice.Hook[int](func(context.Context, validation.Context, int) validation.Report {
		calls++
		return validation.NewReport(vctx.Limits())
	})
	var nilHook validationservice.Validator[int]
	short := validationservice.Chain(validation.ShortCircuit, nilHook, fail, after).
		Validate(context.Background(), vctx, 1)
	if calls != 1 || !short.HasCode("stop") {
		t.Fatalf("short calls=%d report=%v", calls, short)
	}
	calls = 0
	collect := validationservice.Chain(validation.CollectAll, fail, after).
		Validate(context.Background(), vctx, 1)
	if calls != 2 || !collect.HasCode("stop") {
		t.Fatalf("collect calls=%d report=%v", calls, collect)
	}

	caller, cancel := context.WithCancel(context.Background())
	cancel()
	calls = 0
	terminal := validationservice.Chain(validation.CollectAll, after).
		Validate(caller, vctx, 1)
	if calls != 0 || !errors.Is(terminal.Err(), context.Canceled) {
		t.Fatalf("terminal calls=%d report=%v err=%v", calls, terminal, terminal.Err())
	}
}

func TestChainCallerTerminalPrecedesCarriedTerminalAndPartialFindings(t *testing.T) {
	vctx, err := validation.NewContext(validation.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	caller, cancelCaller := context.WithCancel(context.Background())
	derived, cancelDerived := context.WithDeadline(context.Background(), time.Unix(1, 0))
	defer cancelDerived()
	hook := validationservice.Hook[int](func(context.Context, validation.Context, int) validation.Report {
		cancelCaller()
		return validation.ContextReport(vctx, derived).Add(
			validation.NewViolation(vctx.Path(), "partial", validation.Warning, nil, nil),
		)
	})
	report := validationservice.Chain(validation.CollectAll, hook).Validate(caller, vctx, 1)
	if !errors.Is(report.Err(), context.Canceled) ||
		errors.Is(report.Err(), context.DeadlineExceeded) || !report.HasCode("partial") {
		t.Fatalf("report=%v err=%v", report, report.Err())
	}
}
