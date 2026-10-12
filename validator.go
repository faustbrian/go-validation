package validation

// Validator is the deterministic, side-effect-free validation contract.
type Validator[T any] interface {
	// Validate must return ordered findings without mutating borrowed input or
	// hiding IO. Determinism/purity are implementation obligations, not enforced
	// by Context. Instances reused concurrently must support that concurrency.
	// A Report may contain warnings, truncation or terminal state; use Err for
	// complete success. Arbitrary interface panics propagate on direct calls;
	// use IsolatePanics where containment is required.
	Validate(Context, T) Report
}

// ValidatorFunc adapts an ordinary function to Validator.
type ValidatorFunc[T any] func(Context, T) Report

// Validate calls the borrowed function synchronously and contains panics,
// including a nil function, as a blocking validator_panic finding. Report
// diagnostic/path budgets still apply. It cannot prevent mutation/IO/blocking
// or force callback concurrency safety; no input deep copy is made.
func (f ValidatorFunc[T]) Validate(ctx Context, value T) (report Report) {
	defer func() {
		if recover() != nil {
			report = panicReport(ctx)
		}
	}()
	return f(ctx, value)
}

// IsolatePanics explicitly wraps a custom validator with panic containment.
// The panic payload and rejected value are deliberately discarded.
func IsolatePanics[T any](validator Validator[T]) Validator[T] {
	return ValidatorFunc[T](func(ctx Context, value T) (report Report) {
		defer func() {
			if recover() != nil {
				report = panicReport(ctx)
			}
		}()
		return validator.Validate(ctx, value)
	})
}

func panicReport(ctx Context) Report {
	return NewReport(ctx.Limits()).Add(NewViolation(
		ctx.Path(), "validator_panic", Error, nil, ErrValidatorPanic,
	))
}
