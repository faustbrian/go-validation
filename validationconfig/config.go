// Package validationconfig adapts typed validators to a small config contract.
//
// Deprecated: use github.com/faustbrian/go-validation/v2/adapters/config. This
// path remains supported for the longer of 180 days after successor public
// availability and two published stable minor releases.
package validationconfig

import validation "github.com/faustbrian/go-validation/v2"

// Validator is the minimal validation contract used by configuration loaders.
type Validator interface {
	// Validate returns nil for complete nonblocking success, including warnings;
	// return operational/terminal failures rather than treating them as success.
	Validate() error
}

// Check binds a value, immutable context, and typed validator.
type Check[T any] struct {
	value     T
	context   validation.Context
	validator validation.Validator[T]
}

// CheckValue captures shallow borrowed value/validator state without validation
// or nil refusal. Reuse concurrently only when validator/input allow it. It
// owns no IO, close or mutation lifecycle; use a core wrapper for panic safety.
func CheckValue[T any](value T, ctx validation.Context,
	validator validation.Validator[T],
) Check[T] {
	return Check[T]{value: value, context: ctx, validator: validator}
}

// Validate returns the report error when validation does not complete or the
// value is invalid. It directly calls the captured interface; a nil validator
// or arbitrary implementation panic propagates. No input copy is made.
func (check Check[T]) Validate() error {
	return check.validator.Validate(check.context, check.value).Err()
}
