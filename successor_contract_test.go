//lint:file-ignore SA1019 Compatibility coverage requires the deprecated v1 paths.
//nolint:staticcheck // Compatibility coverage requires the deprecated v1 paths.
package validation_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	validation "github.com/faustbrian/go-validation"
	validationconfig "github.com/faustbrian/go-validation/adapters/config"
	validationhttp "github.com/faustbrian/go-validation/adapters/http"
	validationjsonapi "github.com/faustbrian/go-validation/adapters/jsonapi"
	validationjsonrpc "github.com/faustbrian/go-validation/adapters/jsonrpc"
	validationservice "github.com/faustbrian/go-validation/adapters/service"
	legacyconfig "github.com/faustbrian/go-validation/validationconfig"
	legacyhttp "github.com/faustbrian/go-validation/validationhttp"
	legacyjsonapi "github.com/faustbrian/go-validation/validationjsonapi"
	legacyjsonrpc "github.com/faustbrian/go-validation/validationrpc"
	legacyservice "github.com/faustbrian/go-validation/validationservice"
)

type chainSamplingContext struct {
	context.Context
	cancelAt int32
	terminal error
	calls    atomic.Int32
}

func (*chainSamplingContext) Done() <-chan struct{} { return nil }

func (ctx *chainSamplingContext) Err() error {
	if ctx.calls.Add(1) >= ctx.cancelAt {
		return ctx.terminal
	}
	return nil
}

func TestSuccessorPackagesOwnExportedTypeIdentity(t *testing.T) {
	types := map[string]reflect.Type{
		"config Validator":  reflect.TypeOf((*validationconfig.Validator)(nil)).Elem(),
		"config Check":      reflect.TypeOf(validationconfig.Check[int]{}),
		"http Problem":      reflect.TypeOf(validationhttp.Problem{}),
		"http Error":        reflect.TypeOf(validationhttp.Error{}),
		"jsonapi Document":  reflect.TypeOf(validationjsonapi.Document{}),
		"jsonapi Meta":      reflect.TypeOf(validationjsonapi.DocumentMeta{}),
		"jsonapi Error":     reflect.TypeOf(validationjsonapi.Error{}),
		"jsonapi ErrorMeta": reflect.TypeOf(validationjsonapi.ErrorMeta{}),
		"jsonapi Source":    reflect.TypeOf(validationjsonapi.Source{}),
		"jsonrpc Error":     reflect.TypeOf(validationjsonrpc.Error{}),
		"jsonrpc Data":      reflect.TypeOf(validationjsonrpc.Data{}),
		"jsonrpc Violation": reflect.TypeOf(validationjsonrpc.Violation{}),
		"service Validator": reflect.TypeOf((*validationservice.Validator[int])(nil)).Elem(),
		"service Hook":      reflect.TypeOf(validationservice.Hook[int](nil)),
	}
	wants := map[string]string{
		"config Validator":  "github.com/faustbrian/go-validation/adapters/config",
		"config Check":      "github.com/faustbrian/go-validation/adapters/config",
		"http Problem":      "github.com/faustbrian/go-validation/adapters/http",
		"http Error":        "github.com/faustbrian/go-validation/adapters/http",
		"jsonapi Document":  "github.com/faustbrian/go-validation/adapters/jsonapi",
		"jsonapi Meta":      "github.com/faustbrian/go-validation/adapters/jsonapi",
		"jsonapi Error":     "github.com/faustbrian/go-validation/adapters/jsonapi",
		"jsonapi ErrorMeta": "github.com/faustbrian/go-validation/adapters/jsonapi",
		"jsonapi Source":    "github.com/faustbrian/go-validation/adapters/jsonapi",
		"jsonrpc Error":     "github.com/faustbrian/go-validation/adapters/jsonrpc",
		"jsonrpc Data":      "github.com/faustbrian/go-validation/adapters/jsonrpc",
		"jsonrpc Violation": "github.com/faustbrian/go-validation/adapters/jsonrpc",
		"service Validator": "github.com/faustbrian/go-validation/adapters/service",
		"service Hook":      "github.com/faustbrian/go-validation/adapters/service",
	}
	for name, typ := range types {
		if typ.PkgPath() != wants[name] {
			t.Errorf("%s package = %q, want %q", name, typ.PkgPath(), wants[name])
		}
	}
}

func TestSuccessorsMatchLegacyBehaviorWithoutSharingNamedTypes(t *testing.T) {
	vctx := testContext(t)
	report := finding(vctx, "warning", validation.Warning).
		Merge(finding(vctx, "invalid", validation.Error))

	assertJSONEqual(t, "http", legacyhttp.FromReport(report), validationhttp.FromReport(report))
	assertJSONEqual(t, "jsonapi", legacyjsonapi.Errors(report), validationjsonapi.Errors(report))
	assertJSONEqual(t, "jsonrpc", legacyjsonrpc.InvalidParams(report), validationjsonrpc.InvalidParams(report))

	legacyWriter := httptest.NewRecorder()
	successorWriter := httptest.NewRecorder()
	if err := legacyhttp.WriteProblem(legacyWriter, legacyhttp.FromReport(report)); err != nil {
		t.Fatal(err)
	}
	if err := validationhttp.WriteProblem(successorWriter, validationhttp.FromReport(report)); err != nil {
		t.Fatal(err)
	}
	if legacyWriter.Code != successorWriter.Code ||
		legacyWriter.Header().Get("Content-Type") != successorWriter.Header().Get("Content-Type") ||
		!bytes.Equal(legacyWriter.Body.Bytes(), successorWriter.Body.Bytes()) {
		t.Fatalf("HTTP writes differ: legacy=%#v successor=%#v", legacyWriter, successorWriter)
	}

	rootValidator := validation.ValidatorFunc[int](func(validation.Context, int) validation.Report {
		return report
	})
	legacyCheck := legacyconfig.CheckValue(1, vctx, rootValidator)
	successorCheck := validationconfig.CheckValue(1, vctx, rootValidator)
	if !errors.Is(legacyCheck.Validate(), validation.ErrInvalid) ||
		!errors.Is(successorCheck.Validate(), validation.ErrInvalid) {
		t.Fatalf("config errors = %v %v", legacyCheck.Validate(), successorCheck.Validate())
	}

	legacyHook := legacyservice.Hook[int](func(context.Context, validation.Context, int) validation.Report {
		return report
	})
	successorHook := validationservice.Hook[int](func(context.Context, validation.Context, int) validation.Report {
		return report
	})
	legacyReport := legacyservice.Chain(validation.CollectAll, legacyHook).
		Validate(context.Background(), vctx, 1)
	successorReport := validationservice.Chain(validation.CollectAll, successorHook).
		Validate(context.Background(), vctx, 1)
	if legacyReport.String() != successorReport.String() ||
		!reflect.DeepEqual(legacyReport.Violations(), successorReport.Violations()) {
		t.Fatalf("service reports differ: %v %v", legacyReport, successorReport)
	}

	legacyRequestHook := legacyhttp.Hook[int](func(_ *http.Request, _ int) validation.Report { return report })
	successorRequestHook := validationhttp.Hook[int](func(_ *http.Request, _ int) validation.Report { return report })
	request := httptest.NewRequest("GET", "/", nil)
	if legacyRequestHook.Validate(request, 1).String() !=
		successorRequestHook.Validate(request, 1).String() {
		t.Fatal("HTTP hook behavior differs")
	}
}

func TestSuccessorConfigAndServicePreserveTerminalOutcomes(t *testing.T) {
	vctx := testContext(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	terminal := validation.ContextReport(vctx, ctx)
	validator := validation.ValidatorFunc[int](func(validation.Context, int) validation.Report {
		return terminal
	})
	if err := validationconfig.CheckValue(1, vctx, validator).Validate(); !errors.Is(err, context.Canceled) {
		t.Fatalf("config error = %v", err)
	}
	var calls int
	hook := validationservice.Hook[int](func(context.Context, validation.Context, int) validation.Report {
		calls++
		return validation.NewReport(vctx.Limits())
	})
	report := validationservice.Chain(validation.CollectAll, hook).Validate(ctx, vctx, 1)
	if calls != 0 || !errors.Is(report.Err(), context.Canceled) {
		t.Fatalf("service calls=%d report=%v err=%v", calls, report, report.Err())
	}

	derived, cancelDerived := context.WithDeadline(context.Background(), time.Unix(1, 0))
	defer cancelDerived()
	first := validationservice.Hook[int](func(context.Context, validation.Context, int) validation.Report {
		calls++
		return validation.ContextReport(vctx, derived)
	})
	later := validationservice.Hook[int](func(context.Context, validation.Context, int) validation.Report {
		calls++
		return validation.NewReport(vctx.Limits())
	})
	report = validationservice.Chain(validation.CollectAll, first, later).
		Validate(context.Background(), vctx, 1)
	if calls != 1 || !errors.Is(report.Err(), context.DeadlineExceeded) {
		t.Fatalf("derived terminal calls=%d report=%v err=%v", calls, report, report.Err())
	}
}

func TestLegacyAndSuccessorChainsUseEveryCallerContextBoundary(t *testing.T) {
	vctx := testContext(t)
	type runner func(context.Context, validation.Mode, *int) validation.Report
	runners := map[string]runner{
		"legacy": func(ctx context.Context, mode validation.Mode, calls *int) validation.Report {
			first := legacyservice.Hook[int](func(context.Context, validation.Context, int) validation.Report {
				*calls++
				return finding(vctx, "first", validation.Error)
			})
			second := legacyservice.Hook[int](func(context.Context, validation.Context, int) validation.Report {
				*calls++
				return finding(vctx, "second", validation.Warning)
			})
			return legacyservice.Chain(mode, first, nil, second).Validate(ctx, vctx, 1)
		},
		"successor": func(ctx context.Context, mode validation.Mode, calls *int) validation.Report {
			first := validationservice.Hook[int](func(context.Context, validation.Context, int) validation.Report {
				*calls++
				return finding(vctx, "first", validation.Error)
			})
			second := validationservice.Hook[int](func(context.Context, validation.Context, int) validation.Report {
				*calls++
				return finding(vctx, "second", validation.Warning)
			})
			return validationservice.Chain(mode, first, nil, second).Validate(ctx, vctx, 1)
		},
	}
	cases := []struct {
		name        string
		cancelAt    int32
		mode        validation.Mode
		wantCalls   int
		wantSamples int32
		wantCode    string
	}{
		{name: "initial sample", cancelAt: 1, mode: validation.CollectAll, wantSamples: 2},
		{name: "before first hook", cancelAt: 2, mode: validation.CollectAll, wantSamples: 3},
		{name: "after first hook", cancelAt: 3, mode: validation.CollectAll, wantCalls: 1, wantSamples: 4, wantCode: "first"},
		{name: "before second hook", cancelAt: 4, mode: validation.CollectAll, wantCalls: 1, wantSamples: 5, wantCode: "first"},
		{name: "short-circuit final sample", cancelAt: 4, mode: validation.ShortCircuit, wantCalls: 1, wantSamples: 4, wantCode: "first"},
	}
	for runnerName, run := range runners {
		for _, test := range cases {
			t.Run(runnerName+"/"+test.name, func(t *testing.T) {
				ctx := &chainSamplingContext{Context: context.Background(), cancelAt: test.cancelAt, terminal: context.Canceled}
				calls := 0
				report := run(ctx, test.mode, &calls)
				if calls != test.wantCalls || ctx.calls.Load() != test.wantSamples ||
					!errors.Is(report.Err(), context.Canceled) ||
					(test.wantCode != "" && !report.HasCode(test.wantCode)) {
					t.Fatalf("calls=%d samples=%d report=%v err=%v", calls, ctx.calls.Load(), report, report.Err())
				}
				if test.mode == validation.ShortCircuit && !errors.Is(report.Err(), validation.ErrInvalid) {
					t.Fatalf("short-circuit error lost invalid identity: %v", report.Err())
				}
			})
		}
	}
}

func TestLegacyAndSuccessorChainsSampleEmptyAndAllNilReturns(t *testing.T) {
	vctx := testContext(t)
	for _, test := range []struct {
		name string
		run  func(context.Context) validation.Report
	}{
		{"legacy empty", func(ctx context.Context) validation.Report {
			return legacyservice.Chain[int](validation.CollectAll).Validate(ctx, vctx, 1)
		}},
		{"legacy nil", func(ctx context.Context) validation.Report {
			return legacyservice.Chain[int](validation.CollectAll, nil).Validate(ctx, vctx, 1)
		}},
		{"successor empty", func(ctx context.Context) validation.Report {
			return validationservice.Chain[int](validation.CollectAll).Validate(ctx, vctx, 1)
		}},
		{"successor nil", func(ctx context.Context) validation.Report {
			return validationservice.Chain[int](validation.CollectAll, nil).Validate(ctx, vctx, 1)
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx := &chainSamplingContext{Context: context.Background(), cancelAt: 2, terminal: context.DeadlineExceeded}
			report := test.run(ctx)
			if ctx.calls.Load() != 2 || !errors.Is(report.Err(), context.DeadlineExceeded) {
				t.Fatalf("samples=%d report=%v err=%v", ctx.calls.Load(), report, report.Err())
			}
		})
	}
}

func assertJSONEqual(t *testing.T, name string, legacy, successor any) {
	t.Helper()
	legacyJSON, err := json.Marshal(legacy)
	if err != nil {
		t.Fatal(err)
	}
	successorJSON, err := json.Marshal(successor)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(legacyJSON, successorJSON) {
		t.Fatalf("%s JSON differs: %s != %s", name, legacyJSON, successorJSON)
	}
}
