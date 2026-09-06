// Package validationconfig adapts typed validators to a small config contract.
//
// Deprecated: use github.com/faustbrian/go-validation/adapters/config. This
// path remains supported for the longer of 180 days after successor public
// availability and two published stable minor releases.
package validationconfig

import validation "github.com/faustbrian/go-validation"

// Validator is the minimal validation contract used by configuration loaders.
type Validator interface {
	Validate() error
}

// Check binds a value, immutable context, and typed validator.
type Check[T any] struct {
	value     T
	context   validation.Context
	validator validation.Validator[T]
}

// CheckValue constructs a reusable configuration validation contract.
func CheckValue[T any](value T, ctx validation.Context,
	validator validation.Validator[T],
) Check[T] {
	return Check[T]{value: value, context: ctx, validator: validator}
}

// Validate returns the report error when validation does not complete or the
// value is invalid.
func (check Check[T]) Validate() error {
	return check.validator.Validate(check.context, check.value).Err()
}
