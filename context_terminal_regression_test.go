package validation_test

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	validation "github.com/faustbrian/go-validation"
)

type observableTerminalContext struct {
	context.Context
	done        chan struct{}
	errObserved chan struct{}
	terminal    error
	once        sync.Once
}

type finalSampleContext struct {
	context.Context
	calls atomic.Int32
}

//nolint:errorlint // Exact identity is the public contract.
func exactErrorIdentity(actual, expected error) bool {
	return actual == expected
}

func (*finalSampleContext) Done() <-chan struct{} { return nil }

func (ctx *finalSampleContext) Err() error {
	if ctx.calls.Add(1) >= 2 {
		return context.Canceled
	}
	return nil
}

func (ctx *observableTerminalContext) Done() <-chan struct{} { return ctx.done }

func (ctx *observableTerminalContext) Err() error {
	select {
	case <-ctx.done:
		ctx.once.Do(func() { close(ctx.errObserved) })
		return ctx.terminal
	default:
		return nil
	}
}

func canceledReport(vctx validation.Context) validation.Report {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return validation.AsyncAll[int](ctx, vctx, 0)
}

func deadlineReport(vctx validation.Context) validation.Report {
	ctx, cancel := context.WithDeadline(context.Background(), time.Unix(1, 0))
	defer cancel()
	return validation.AsyncAll[int](ctx, vctx, 0)
}

func finding(vctx validation.Context, code string, severity validation.Severity) validation.Report {
	return validation.NewReport(vctx.Limits()).Add(validation.NewViolation(
		vctx.Path(), code, severity, nil, nil,
	))
}

func TestContextReportAndContextErrorExposeBoundedImmutableState(t *testing.T) {
	vctx := testContext(t)
	sampled := &finalSampleContext{Context: context.Background()}
	if report := validation.ContextReport(vctx, sampled); report.Err() != nil || sampled.calls.Load() != 1 {
		t.Fatalf("active sample count=%d report=%v", sampled.calls.Load(), report)
	}
	if report := validation.ContextReport(vctx, context.Background()); report.Err() != nil || report.ContextError() != nil || !report.Empty() {
		t.Fatalf("active report=%v err=%v context=%v", report, report.Err(), report.ContextError())
	}

	cause := errors.New("token=secret")
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(cause)
	report := validation.ContextReport(vctx, ctx)
	if !exactErrorIdentity(report.ContextError(), context.Canceled) ||
		!errors.Is(report.Err(), context.Canceled) || errors.Is(report.Err(), cause) ||
		report.String() != "validation canceled" || !report.Empty() || report.HasErrors() {
		t.Fatalf("canceled report=%v err=%v context=%v", report, report.Err(), report.ContextError())
	}
	var terminal *validation.ContextError
	if !errors.As(report.Err(), &terminal) || terminal == nil {
		t.Fatalf("structured terminal = %#v", terminal)
	}
	if !exactErrorIdentity(terminal.Report().ContextError(), context.Canceled) {
		t.Fatalf("structured terminal report = %#v", terminal.Report())
	}
	first := terminal.Unwrap()
	second := terminal.Unwrap()
	third := terminal.Unwrap()
	if len(first) != 1 || len(second) != 1 || len(third) != 1 || &first[0] == &second[0] {
		t.Fatalf("unwrap slices = %#v %#v", first, second)
	}
	first[0] = errors.New("mutated")
	if !exactErrorIdentity(second[0], context.Canceled) ||
		!exactErrorIdentity(third[0], context.Canceled) {
		t.Fatalf("unwrap mutation escaped: %#v %#v", second, third)
	}
	mutated := terminal.Report().Add(validation.NewViolation(
		vctx.Path(), "new", validation.Error, nil, nil,
	))
	if !mutated.HasCode("new") || terminal.Report().HasCode("new") {
		t.Fatal("ContextError.Report did not preserve report immutability")
	}

	deadline := validation.ContextReport(vctx, deadlineContext(t))
	deadlineErr := deadline.Err()
	if deadlineErr == nil {
		t.Fatalf("deadline report unexpectedly valid: %v", deadline)
	}
	if !exactErrorIdentity(deadline.ContextError(), context.DeadlineExceeded) ||
		!errors.Is(deadlineErr, context.DeadlineExceeded) ||
		deadlineErr.Error() != "validation deadline exceeded" ||
		deadline.String() != "validation deadline exceeded" {
		t.Fatalf("deadline report=%v err=%v context=%v", deadline, deadlineErr, deadline.ContextError())
	}

	zero := &validation.ContextError{}
	if zero.Error() != "validation context not terminated" || zero.Unwrap() != nil ||
		zero.Report().Err() != nil || errors.Is(zero, context.Canceled) ||
		errors.Is(zero, validation.ErrInvalid) {
		t.Fatalf("zero ContextError = %q %#v %#v", zero.Error(), zero.Unwrap(), zero.Report())
	}
}

func TestContextErrorPreservesPartialInvalidReportAndTerminalFormatting(t *testing.T) {
	limits := validation.DefaultLimits()
	limits.MaxViolations = 1
	vctx, err := validation.NewContext(limits)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	report := validation.ContextReport(vctx, ctx).
		Add(validation.NewViolation(vctx.Path(), "warning", validation.Warning, nil, nil)).
		Add(validation.NewViolation(vctx.Path(), "blocked", validation.Error, nil, nil))
	if report.String() != "validation canceled with 1 violation (truncated)" ||
		!report.Truncated() || !report.HasErrors() || report.Len() != 1 ||
		!errors.Is(report.Err(), context.Canceled) ||
		!errors.Is(report.Err(), validation.ErrInvalid) {
		t.Fatalf("partial report=%v err=%v", report, report.Err())
	}
	var invalid *validation.InvalidError
	if !errors.As(report.Err(), &invalid) ||
		!exactErrorIdentity(invalid.Report().ContextError(), context.Canceled) {
		t.Fatalf("invalid projection = %#v", invalid)
	}
	var terminal *validation.ContextError
	if !errors.As(report.Err(), &terminal) || len(terminal.Unwrap()) != 2 {
		t.Fatalf("terminal projection = %#v", terminal)
	}
	unwrapped := terminal.Unwrap()
	if !exactErrorIdentity(unwrapped[0], context.Canceled) ||
		!errors.Is(unwrapped[1], validation.ErrInvalid) {
		t.Fatalf("unwrap order = %#v", unwrapped)
	}
	wide := testContext(t)
	plural := validation.ContextReport(wide, ctx).
		Add(validation.NewViolation(wide.Path(), "first", validation.Warning, nil, nil)).
		Add(validation.NewViolation(wide.Path(), "second", validation.Warning, nil, nil))
	if plural.String() != "validation canceled with 2 violations" {
		t.Fatalf("bounded plural report = %v", plural)
	}
}

func deadlineContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithDeadline(context.Background(), time.Unix(1, 0))
	t.Cleanup(cancel)
	return ctx
}

func TestAsyncAllReturnsTerminalForPreTerminatedContexts(t *testing.T) {
	vctx := testContext(t)
	var calls atomic.Int32
	validator := validation.AsyncValidatorFunc[int](func(
		context.Context, validation.Context, int,
	) validation.Report {
		calls.Add(1)
		return finding(vctx, "unexpected", validation.Error)
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for _, validators := range [][]validation.AsyncValidator[int]{nil, {validator}} {
		report := validation.AsyncAll(ctx, vctx, 1, validators...)
		if !errors.Is(report.Err(), context.Canceled) || !report.Empty() {
			t.Fatalf("canceled report = %v, err = %v", report, report.Err())
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("callbacks started = %d", calls.Load())
	}

	expired, stop := context.WithDeadline(context.Background(), time.Unix(1, 0))
	defer stop()
	report := validation.AsyncAll(expired, vctx, 1, validator)
	if !errors.Is(report.Err(), context.DeadlineExceeded) || !report.Empty() {
		t.Fatalf("deadline report = %v, err = %v", report, report.Err())
	}
	if calls.Load() != 0 {
		t.Fatalf("deadline callback started = %d", calls.Load())
	}
}

func TestAsyncAllCancellationStopsAdmissionJoinsAndRetainsPartialResults(t *testing.T) {
	limits := validation.DefaultLimits()
	limits.MaxCustomConcurrency = 1
	vctx, err := validation.NewContext(limits)
	if err != nil {
		t.Fatal(err)
	}
	ctx := &observableTerminalContext{
		Context: context.Background(), done: make(chan struct{}),
		errObserved: make(chan struct{}),
		terminal:    context.Canceled,
	}
	started := make(chan struct{})
	release := make(chan struct{})
	returned := make(chan struct{})
	first := validation.AsyncValidatorFunc[int](func(
		context.Context, validation.Context, int,
	) validation.Report {
		close(started)
		<-release
		close(returned)
		return finding(vctx, "partial", validation.Warning)
	})
	var laterCalls atomic.Int32
	later := validation.AsyncValidatorFunc[int](func(
		context.Context, validation.Context, int,
	) validation.Report {
		laterCalls.Add(1)
		return finding(vctx, "later", validation.Error)
	})
	done := make(chan validation.Report, 1)
	go func() { done <- validation.AsyncAll(ctx, vctx, 1, first, later) }()
	<-started
	close(ctx.done)
	<-ctx.errObserved
	close(release)
	report := <-done
	select {
	case <-returned:
	default:
		t.Fatal("AsyncAll returned before admitted callback")
	}
	if laterCalls.Load() != 0 || !errors.Is(report.Err(), context.Canceled) ||
		!report.HasCode("partial") || report.HasErrors() {
		t.Fatalf("calls=%d report=%v err=%v", laterCalls.Load(), report, report.Err())
	}
}

func TestAsyncAllDeadlineJoinsBoundedWorkAndKeepsWarnings(t *testing.T) {
	limits := validation.DefaultLimits()
	limits.MaxCustomConcurrency = 2
	vctx, err := validation.NewContext(limits)
	if err != nil {
		t.Fatal(err)
	}
	ctx := &observableTerminalContext{
		Context: context.Background(), done: make(chan struct{}),
		errObserved: make(chan struct{}),
		terminal:    context.DeadlineExceeded,
	}
	var active atomic.Int32
	var maximum atomic.Int32
	started := make(chan struct{}, 2)
	returned := make(chan struct{}, 2)
	validators := make([]validation.AsyncValidator[int], 2)
	for index := range validators {
		code := []string{"first", "second"}[index]
		validators[index] = validation.AsyncValidatorFunc[int](func(
			ctx context.Context, _ validation.Context, _ int,
		) validation.Report {
			current := active.Add(1)
			for previous := maximum.Load(); current > previous; previous = maximum.Load() {
				if maximum.CompareAndSwap(previous, current) {
					break
				}
			}
			started <- struct{}{}
			<-ctx.Done()
			active.Add(-1)
			returned <- struct{}{}
			return finding(vctx, code, validation.Warning)
		})
	}
	done := make(chan validation.Report, 1)
	go func() { done <- validation.AsyncAll(ctx, vctx, 1, validators...) }()
	<-started
	<-started
	close(ctx.done)
	<-ctx.errObserved
	report := <-done
	if maximum.Load() != 2 || active.Load() != 0 || len(returned) != 2 ||
		!errors.Is(report.Err(), context.DeadlineExceeded) || report.HasErrors() {
		t.Fatalf("maximum=%d active=%d returned=%d report=%v err=%v",
			maximum.Load(), active.Load(), len(returned), report, report.Err())
	}
	violations := report.Violations()
	if len(violations) != 2 || violations[0].Code() != "first" ||
		violations[1].Code() != "second" {
		t.Fatalf("violations = %#v", violations)
	}
}

func TestAsyncAllFinalContextSampleIsTheReturnBoundary(t *testing.T) {
	vctx := testContext(t)
	ctx := &finalSampleContext{Context: context.Background()}
	report := validation.AsyncAll(ctx, vctx, 1,
		validation.AsyncValidatorFunc[int](func(received context.Context,
			_ validation.Context, _ int,
		) validation.Report {
			if received != ctx {
				t.Fatal("callback did not receive original caller context")
			}
			return finding(vctx, "completed", validation.Warning)
		}),
	)
	if ctx.calls.Load() != 2 || !errors.Is(report.Err(), context.Canceled) ||
		!report.HasCode("completed") {
		t.Fatalf("samples=%d report=%v err=%v", ctx.calls.Load(), report, report.Err())
	}
}

func TestAsyncAllPreservesPanicAndBlockingFindingsWithCancellation(t *testing.T) {
	vctx := testContext(t)
	t.Run("panic", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		report := validation.AsyncAll(ctx, vctx, 1,
			validation.AsyncValidatorFunc[int](func(
				context.Context, validation.Context, int,
			) validation.Report {
				cancel()
				panic("token=secret")
			}),
		)
		if !errors.Is(report.Err(), context.Canceled) ||
			!errors.Is(report.Err(), validation.ErrInvalid) ||
			!report.HasCode("validator_panic") ||
			report.String() != "validation canceled with 1 violation" ||
			strings.Contains(report.String(), "secret") {
			t.Fatalf("report=%v err=%v", report, report.Err())
		}
	})

	t.Run("blocking partial", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		report := validation.AsyncAll(ctx, vctx, 1,
			validation.AsyncValidatorFunc[int](func(
				context.Context, validation.Context, int,
			) validation.Report {
				cancel()
				return finding(vctx, "invalid", validation.Error)
			}),
		)
		if !errors.Is(report.Err(), context.Canceled) ||
			!errors.Is(report.Err(), validation.ErrInvalid) ||
			!report.HasCode("invalid") {
			t.Fatalf("report=%v err=%v", report, report.Err())
		}
	})
}

func TestAsyncAllCarriedTerminalDoesNotStopActiveCallerAdmission(t *testing.T) {
	limits := validation.DefaultLimits()
	limits.MaxCustomConcurrency = 1
	vctx, err := validation.NewContext(limits)
	if err != nil {
		t.Fatal(err)
	}
	derived, cancel := context.WithCancel(context.Background())
	cancel()
	var calls atomic.Int32
	validators := []validation.AsyncValidator[int]{
		validation.AsyncValidatorFunc[int](func(
			context.Context, validation.Context, int,
		) validation.Report {
			calls.Add(1)
			return validation.ContextReport(vctx, derived).
				Merge(finding(vctx, "first", validation.Warning))
		}),
		validation.AsyncValidatorFunc[int](func(
			context.Context, validation.Context, int,
		) validation.Report {
			calls.Add(1)
			return deadlineReport(vctx).
				Merge(finding(vctx, "second", validation.Error))
		}),
		validation.AsyncValidatorFunc[int](func(
			context.Context, validation.Context, int,
		) validation.Report {
			calls.Add(1)
			return finding(vctx, "third", validation.Warning)
		}),
	}
	report := validation.AsyncAll(context.Background(), vctx, 1, validators...)
	if calls.Load() != 3 || !errors.Is(report.Err(), context.Canceled) ||
		errors.Is(report.Err(), context.DeadlineExceeded) ||
		!errors.Is(report.Err(), validation.ErrInvalid) {
		t.Fatalf("calls=%d report=%v err=%v", calls.Load(), report, report.Err())
	}
	violations := report.Violations()
	if len(violations) != 3 || violations[0].Code() != "first" ||
		violations[1].Code() != "second" || violations[2].Code() != "third" {
		t.Fatalf("violations = %#v", violations)
	}
}

func TestReportMergeUsesReceiverTerminalPrecedence(t *testing.T) {
	vctx := testContext(t)
	canceled := canceledReport(vctx).Merge(finding(vctx, "first", validation.Warning))
	deadline := deadlineReport(vctx).Merge(finding(vctx, "second", validation.Error))

	report := canceled.Merge(deadline)
	if !errors.Is(report.Err(), context.Canceled) ||
		errors.Is(report.Err(), context.DeadlineExceeded) ||
		!errors.Is(report.Err(), validation.ErrInvalid) {
		t.Fatalf("receiver-first err = %v", report.Err())
	}
	violations := report.Violations()
	if len(violations) != 2 || violations[0].Code() != "first" ||
		violations[1].Code() != "second" {
		t.Fatalf("violations = %#v", violations)
	}

	reversed := deadline.Merge(canceled)
	if !errors.Is(reversed.Err(), context.DeadlineExceeded) ||
		errors.Is(reversed.Err(), context.Canceled) {
		t.Fatalf("reversed err = %v", reversed.Err())
	}
}

func TestCompositorsPreserveAndStopOnTerminalReports(t *testing.T) {
	vctx := testContext(t)
	terminal := validation.ValidatorFunc[int](func(validation.Context, int) validation.Report {
		return canceledReport(vctx)
	})
	var laterCalls atomic.Int32
	later := validation.ValidatorFunc[int](func(validation.Context, int) validation.Report {
		laterCalls.Add(1)
		return finding(vctx, "later", validation.Warning)
	})

	for name, validator := range map[string]validation.Validator[int]{
		"all":       validation.All(validation.CollectAll, terminal, later),
		"any":       validation.Any(validation.CollectAll, terminal, later),
		"when":      validation.When(func(int) bool { return true }, terminal, later),
		"dependent": validation.Dependent(terminal, later),
	} {
		laterCalls.Store(0)
		report := validator.Validate(vctx, 1)
		if !errors.Is(report.Err(), context.Canceled) || laterCalls.Load() != 0 {
			t.Fatalf("%s calls=%d report=%v err=%v", name, laterCalls.Load(), report, report.Err())
		}
	}

	report := validation.Not(terminal).Validate(vctx, 1)
	if !errors.Is(report.Err(), context.Canceled) || report.HasCode("not") {
		t.Fatalf("Not report=%v err=%v", report, report.Err())
	}
}

func TestCompositorsRetainOrderedPartialReportsAroundTerminal(t *testing.T) {
	vctx := testContext(t)
	warning := validation.ValidatorFunc[int](func(validation.Context, int) validation.Report {
		return finding(vctx, "before", validation.Warning)
	})
	terminal := validation.ValidatorFunc[int](func(validation.Context, int) validation.Report {
		return canceledReport(vctx).Merge(finding(vctx, "terminal", validation.Warning))
	})
	var laterCalls atomic.Int32
	later := validation.ValidatorFunc[int](func(validation.Context, int) validation.Report {
		laterCalls.Add(1)
		return finding(vctx, "later", validation.Error)
	})
	for name, validator := range map[string]validation.Validator[int]{
		"all": validation.All(validation.CollectAll, warning, terminal, later),
		"any": validation.Any(validation.CollectAll,
			validation.ValidatorFunc[int](func(validation.Context, int) validation.Report {
				return finding(vctx, "before", validation.Error)
			}), terminal, later),
		"dependent": validation.Dependent(warning, terminal),
	} {
		laterCalls.Store(0)
		report := validator.Validate(vctx, 1)
		if !errors.Is(report.Err(), context.Canceled) || laterCalls.Load() != 0 {
			t.Fatalf("%s calls=%d report=%v err=%v", name, laterCalls.Load(), report, report.Err())
		}
		violations := report.Violations()
		if len(violations) != 2 || violations[0].Code() != "before" ||
			violations[1].Code() != "terminal" {
			t.Fatalf("%s violations = %#v", name, violations)
		}
	}
	not := validation.Not(terminal).Validate(vctx, 1)
	if !errors.Is(not.Err(), context.Canceled) || !not.HasCode("terminal") || not.HasCode("not") {
		t.Fatalf("Not report=%v err=%v", not, not.Err())
	}
}
