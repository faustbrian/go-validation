package validation

// Mode controls whether composition stops after a decisive result.
type Mode uint8

const (
	// ShortCircuit stops after the first decisive failure or success.
	ShortCircuit Mode = iota + 1
	// CollectAll evaluates every relevant validator in declaration order.
	CollectAll
)

// All requires every validator to pass. The supplied positions, including nil
// entries, must fit Context.Limits().MaxCollectionSize before any invocation.
// The supplied slice/input is borrowed across future calls. Empty/all-nil
// inputs pass; zero/unknown Mode collects all. An arbitrary nested panic is
// contained by the outer function adapter and can replace prior findings.
func All[T any](mode Mode, validators ...Validator[T]) Validator[T] {
	return ValidatorFunc[T](func(ctx Context, value T) Report {
		if len(validators) > ctx.Limits().MaxCollectionSize {
			return collectionLimitReport(ctx)
		}
		report := NewReport(ctx.Limits())
		for _, validator := range validators {
			if validator == nil {
				continue
			}
			current := validator.Validate(ctx, value)
			report = report.Merge(current)
			if current.ContextError() != nil ||
				(mode == ShortCircuit && current.Err() != nil) {
				break
			}
		}
		return report
	})
}

// Any requires at least one validator to pass. Failed alternatives are
// returned only when every alternative fails. The supplied positions, including
// nil entries, must fit Context.Limits().MaxCollectionSize before any invocation.
// Empty/all-nil inputs also return nonblocking success. If any alternative
// passes, only successful alternatives' findings (including warnings) remain;
// a terminal alternative returns accumulated partial findings immediately.
// Caller-owned validator slices must stay stable across future calls.
func Any[T any](mode Mode, validators ...Validator[T]) Validator[T] {
	return ValidatorFunc[T](func(ctx Context, value T) Report {
		if len(validators) > ctx.Limits().MaxCollectionSize {
			return collectionLimitReport(ctx)
		}
		failures := NewReport(ctx.Limits())
		successes := NewReport(ctx.Limits())
		partial := NewReport(ctx.Limits())
		passed := false
		for _, validator := range validators {
			if validator == nil {
				continue
			}
			current := validator.Validate(ctx, value)
			partial = partial.Merge(current)
			if current.ContextError() != nil {
				return partial
			}
			if current.Err() == nil {
				passed = true
				successes = successes.Merge(current)
				if mode == ShortCircuit {
					break
				}
				continue
			}
			failures = failures.Merge(current)
		}
		if passed {
			return successes
		}
		return failures
	})
}

func collectionLimitReport(ctx Context) Report {
	return NewReport(ctx.Limits()).Add(NewViolation(
		ctx.Path(), "collection_limit", Error, nil, ErrLimitExceeded,
	))
}

// Not passes only when validator fails.
func Not[T any](validator Validator[T]) Validator[T] {
	return ValidatorFunc[T](func(ctx Context, value T) Report {
		if validator != nil {
			report := validator.Validate(ctx, value)
			if report.ContextError() != nil {
				return report
			}
			if report.Err() != nil {
				return NewReport(ctx.Limits())
			}
		}
		return NewReport(ctx.Limits()).Add(NewViolation(
			ctx.Path(), "not", Error, nil, nil,
		))
	})
}

// When chooses a validator using a deterministic typed predicate.
func When[T any](predicate func(T) bool, then, otherwise Validator[T]) Validator[T] {
	return ValidatorFunc[T](func(ctx Context, value T) Report {
		selected := otherwise
		if predicate != nil && predicate(value) {
			selected = then
		}
		if selected == nil {
			return NewReport(ctx.Limits())
		}
		return selected.Validate(ctx, value)
	})
}

// Dependent runs dependent only when prerequisite has no blocking errors.
func Dependent[T any](prerequisite, dependent Validator[T]) Validator[T] {
	return ValidatorFunc[T](func(ctx Context, value T) Report {
		report := NewReport(ctx.Limits())
		if prerequisite != nil {
			report = prerequisite.Validate(ctx, value)
			if report.Err() != nil {
				return report
			}
		}
		if dependent == nil {
			return report
		}
		return report.Merge(dependent.Validate(ctx, value))
	})
}
