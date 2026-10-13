// Package structplan provides reflection-free typed plans and optional
// startup-compiled validate-tag plans. New/Add/Compile or CompileTags are entry
// points; these plans validate without binding, coercion or provider IO.
// Builders are mutable and not concurrent; compiled definitions are immutable
// but borrow callbacks/input. Cache serializes compilation and owns no shutdown.
// Tag plans reject cyclic runtime pointer/interface dereferences independently
// of compile-time type/depth checks; finite acyclic chains have no hop cap.
package structplan

import (
	"errors"
	"fmt"

	validation "github.com/faustbrian/go-validation/v3"
)

var (
	// ErrDuplicateField marks a repeated typed-plan field name.
	ErrDuplicateField = errors.New("duplicate validation field")
	// ErrUnknownRule marks an unsupported strict tag rule.
	ErrUnknownRule = errors.New("unknown validation rule")
	// ErrDuplicateRule marks a repeated rule in one tag.
	ErrDuplicateRule = errors.New("duplicate validation rule")
	// ErrInvalidTag marks malformed tag grammar or parameters.
	ErrInvalidTag = errors.New("invalid validation tag")
	// ErrCycle marks a recursive reflected struct or pointer type traversal.
	ErrCycle = errors.New("cyclic validation type")
	// ErrUnsupportedKind marks a type a tag rule cannot evaluate.
	ErrUnsupportedKind = errors.New("unsupported validation kind")
	// ErrInvalidPlan marks malformed typed plan construction.
	ErrInvalidPlan = errors.New("invalid validation plan")
)

type fieldPlan[T any] func(validation.Context, T) validation.Report

// Builder incrementally defines a reflection-free typed struct plan.
type Builder[T any] struct {
	limits validation.Limits
	fields []fieldPlan[T]
	names  map[string]struct{}
}

// New creates an owned mutable builder without validating limits until
// Compile. Use positive/default limits; zero builder is not an initialized
// builder. Caller owns serialized Add/Compile and borrowed callback lifetime.
func New[T any](limits validation.Limits) *Builder[T] {
	return &Builder[T]{limits: limits, names: make(map[string]struct{})}
}

// Add appends a named typed accessor and validator to builder.
func Add[T, V any](builder *Builder[T], name string, accessor func(T) V,
	validator validation.Validator[V],
) error {
	if builder == nil || name == "" || accessor == nil || validator == nil {
		return ErrInvalidPlan
	}
	if len(name) > builder.limits.MaxPathLength {
		return fmt.Errorf("%w: field path", validation.ErrLimitExceeded)
	}
	if _, exists := builder.names[name]; exists {
		return privateConstructionError(fmt.Errorf("%w: %s", ErrDuplicateField, name))
	}
	if len(builder.fields) >= builder.limits.MaxStructFields {
		return fmt.Errorf("%w: struct fields", validation.ErrLimitExceeded)
	}
	builder.names[name] = struct{}{}
	builder.fields = append(builder.fields, func(ctx validation.Context,
		value T,
	) validation.Report {
		safeField := validation.ValidatorFunc[T](func(fieldContext validation.Context,
			_ T,
		) validation.Report {
			return validator.Validate(fieldContext, accessor(value))
		})
		return safeField.Validate(ctx.WithPath(validation.Field(name)), value)
	})
	return nil
}

// Plan construction diagnostics can contain caller-defined fields, types and
// rules. Preserve their original cause chain behind a fixed public category.
type constructionError struct {
	category string
	cause    error
}

func (err *constructionError) Error() string { return err.category }

func (err *constructionError) Unwrap() error { return err.cause }

func privateConstructionError(cause error) error {
	for _, kind := range []error{
		ErrDuplicateField, ErrUnknownRule, ErrDuplicateRule, ErrInvalidTag,
		ErrCycle, ErrUnsupportedKind, ErrInvalidPlan,
		validation.ErrLimitExceeded, validation.ErrInvalidLimit,
	} {
		if errors.Is(cause, kind) {
			return &constructionError{category: kind.Error(), cause: cause}
		}
	}
	return &constructionError{category: "validation plan construction failed", cause: cause}
}

// Plan is an immutable reflection-free typed struct plan.
type Plan[T any] struct {
	limits validation.Limits
	fields []fieldPlan[T]
}

// Compile validates limits and copies field definitions, not accessor/validator
// state. Later Add calls do not change the plan. Nil receivers panic rather
// than return ErrInvalidPlan; instances need serialized builder access.
func (builder *Builder[T]) Compile() (*Plan[T], error) {
	if _, err := validation.NewContext(builder.limits); err != nil {
		return nil, fmt.Errorf("compile typed plan: %w", err)
	}
	return &Plan[T]{limits: builder.limits,
		fields: append([]fieldPlan[T](nil), builder.fields...)}, nil
}

// Validate evaluates captured fields in declaration order under caller Context
// report limits, without recapping compiled definition count. Field callbacks
// are panic-contained individually; terminal state stops later fields.
// Generic inputs are borrowed and callbacks must support any concurrent reuse.
func (plan *Plan[T]) Validate(ctx validation.Context, value T) validation.Report {
	report := validation.NewReport(ctx.Limits())
	for _, field := range plan.fields {
		current := field(ctx, value)
		report = report.Merge(current)
		if current.ContextError() != nil {
			break
		}
	}
	return report
}
