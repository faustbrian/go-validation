package validation

import "context"

// AsyncValidator is the separate contract for cancellation-aware I/O
// validation. Implementations are not deterministic Validator values.
type AsyncValidator[T any] interface {
	ValidateAsync(context.Context, Context, T) Report
}

// AsyncValidatorFunc adapts a context-aware function to AsyncValidator.
type AsyncValidatorFunc[T any] func(context.Context, Context, T) Report

// ValidateAsync calls the context-aware function.
func (f AsyncValidatorFunc[T]) ValidateAsync(
	ctx context.Context, validationContext Context, value T,
) (report Report) {
	defer func() {
		if recover() != nil {
			report = panicReport(validationContext)
		}
	}()
	return f(ctx, validationContext, value)
}

// IsolateAsyncPanics wraps an arbitrary asynchronous validator with the same
// secret-safe panic containment provided by AsyncValidatorFunc.
func IsolateAsyncPanics[T any](validator AsyncValidator[T]) AsyncValidator[T] {
	return AsyncValidatorFunc[T](func(ctx context.Context,
		validationContext Context, value T,
	) Report {
		return validator.ValidateAsync(ctx, validationContext, value)
	})
}

// AsyncAll executes context-aware validators with bounded concurrency and
// merges their reports in declaration order. Cancellation stops unscheduled
// work; validators already running remain responsible for honoring ctx.
func AsyncAll[T any](ctx context.Context, validationContext Context, value T,
	validators ...AsyncValidator[T],
) Report {
	if terminal := ContextReport(validationContext, ctx); terminal.ContextError() != nil {
		return terminal
	}
	if len(validators) == 0 {
		return NewReport(validationContext.Limits())
	}
	workerCount := min(validationContext.Limits().MaxCustomConcurrency,
		len(validators))
	type job struct {
		index     int
		validator AsyncValidator[T]
	}
	jobs := make(chan job)
	results := make([]Report, len(validators))
	done := make(chan struct{}, workerCount)
	for range workerCount {
		go func() {
			defer func() { done <- struct{}{} }()
			for current := range jobs {
				if current.validator == nil {
					results[current.index] = NewReport(validationContext.Limits())
				} else {
					results[current.index] = IsolateAsyncPanics(current.validator).
						ValidateAsync(ctx, validationContext, value)
				}
			}
		}()
	}
	callerTerminal := func() Report {
		for index, validator := range validators {
			select {
			case jobs <- job{index: index, validator: validator}:
			case <-ctx.Done():
				return ContextReport(validationContext, ctx)
			}
		}
		return NewReport(validationContext.Limits())
	}()
	close(jobs)
	for range workerCount {
		<-done
	}
	merged := mergeAsyncReports(validationContext, results)
	callerTerminal = callerTerminal.Merge(ContextReport(validationContext, ctx))
	if callerTerminal.ContextError() != nil {
		return callerTerminal.Merge(merged)
	}
	return merged
}

func mergeAsyncReports(validationContext Context, reports []Report) Report {
	result := NewReport(validationContext.Limits())
	for _, report := range reports {
		result = result.Merge(report)
	}
	return result
}
