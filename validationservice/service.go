// Package validationservice provides transport-neutral service hook contracts.
//
// Deprecated: use github.com/faustbrian/go-validation/adapters/service. This
// path remains supported for the longer of 180 days after successor public
// availability and two published stable minor releases.
package validationservice

import (
	"context"

	validation "github.com/faustbrian/go-validation"
)

// Validator is a cancellation-aware service-boundary validation contract.
type Validator[T any] interface {
	Validate(context.Context, validation.Context, T) validation.Report
}

// Hook adapts a service-boundary function to Validator.
type Hook[T any] func(context.Context, validation.Context, T) validation.Report

// Validate invokes the service hook.
func (hook Hook[T]) Validate(ctx context.Context,
	validationContext validation.Context, value T,
) validation.Report {
	return hook(ctx, validationContext, value)
}

// Chain evaluates service hooks in declaration order and preserves caller
// cancellation or deadline as a terminal validation outcome.
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
