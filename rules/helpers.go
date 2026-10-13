// Package rules provides typed validators without coercion/defaulting. Use
// constructors to compose core Validator values; generic inputs/callbacks are
// borrowed, not deep copied. These are package-defined bounded syntactic and
// Go comparison profiles, not full URL/email/UUID or numerical specification
// conformance. Core function adapters contain panics, not arbitrary blocking.
package rules

import (
	"reflect"

	validation "github.com/faustbrian/go-validation/v3"
)

func pass(ctx validation.Context) validation.Report {
	return validation.NewReport(ctx.Limits())
}

func fail(ctx validation.Context, code string,
	parameters map[string]string,
) validation.Report {
	return pass(ctx).Add(validation.NewViolation(
		ctx.Path(), code, validation.Error, parameters, nil,
	))
}

func rejectOversizedString(ctx validation.Context,
	value string,
) (validation.Report, bool) {
	if len(value) <= ctx.Limits().MaxStringLength {
		return validation.Report{}, false
	}
	return fail(ctx, "string_limit", nil), true
}

func valueExceedsStringLimit(ctx validation.Context, value any) bool {
	reflected := reflect.ValueOf(value)
	return reflected.Len() > ctx.Limits().MaxStringLength
}

func isStringType[T any]() bool { return reflect.TypeFor[T]().Kind() == reflect.String }
