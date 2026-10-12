// Package validationservice provides transport-neutral service hook contracts.
package validationservice

import (
	"context"

	validation "github.com/faustbrian/go-validation/v2"
)

// Validator is a cancellation-aware service-boundary validation contract.
type Validator[T any] interface {
	// Validate borrows caller context/input and must return cooperatively.
	// Instances reused concurrently must support it; direct interface panics
	// propagate unless invoked through a containing Hook/Chain boundary.
	Validate(context.Context, validation.Context, T) validation.Report
}

// Hook adapts a service-boundary function to Validator.
type Hook[T any] func(context.Context, validation.Context, T) validation.Report

// Validate invokes the borrowed hook synchronously with core panic containment
// and supplied validation limits. Nil hook panics become findings; cancellation
// is forwarded rather than captured here. No callback preemption/input copy
// or provider lifecycle is owned.
func (hook Hook[T]) Validate(ctx context.Context,
	validationContext validation.Context, value T,
) validation.Report {
	return validation.ValidatorFunc[T](func(validationContext validation.Context, value T) validation.Report {
		return hook(ctx, validationContext, value)
	}).Validate(validationContext, value)
}

// Chain evaluates service hooks in declaration order and preserves caller
// cancellation or deadline as a terminal validation outcome.
// Total supplied positions, including nil, must fit MaxCollectionSize before
// invocation; caller cancellation takes precedence over this admission.
// The validator slice is borrowed until future calls, not copied at creation.
// Unknown/zero Mode behaves as CollectAll. Chain owns no workers or shutdown;
// each direct provider must return before terminal state can be captured.
// A provider panic aborts this invocation through the Hook boundary and can
// replace already accumulated findings with the core panic report.
func Chain[T any](mode validation.Mode, validators ...Validator[T]) Validator[T] {
	return Hook[T](func(ctx context.Context,
		validationContext validation.Context, value T,
	) validation.Report {
		finish := func(report validation.Report) validation.Report {
			terminal := validation.ContextReport(validationContext, ctx)
			if terminal.ContextError() != nil {
				return terminal.Merge(report)
			}
			return report
		}
		if terminal := validation.ContextReport(validationContext, ctx); terminal.ContextError() != nil {
			return finish(terminal)
		}
		if len(validators) > validationContext.Limits().MaxCollectionSize {
			return finish(collectionLimitReport(validationContext))
		}
		report := validation.NewReport(validationContext.Limits())
		for _, validator := range validators {
			if validator == nil {
				continue
			}
			if terminal := validation.ContextReport(validationContext, ctx); terminal.ContextError() != nil {
				return finish(terminal.Merge(report))
			}
			current := validator.Validate(ctx, validationContext, value)
			report = report.Merge(current)
			if terminal := validation.ContextReport(validationContext, ctx); terminal.ContextError() != nil {
				return finish(terminal.Merge(report))
			}
			if current.ContextError() != nil {
				break
			}
			if mode == validation.ShortCircuit && current.Err() != nil {
				break
			}
		}
		return finish(report)
	})
}

func collectionLimitReport(ctx validation.Context) validation.Report {
	return validation.NewReport(ctx.Limits()).Add(validation.NewViolation(
		ctx.Path(), "collection_limit", validation.Error, nil, validation.ErrLimitExceeded,
	))
}
