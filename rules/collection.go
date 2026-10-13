package rules

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"

	validation "github.com/faustbrian/go-validation/v3"
)

// SliceSize requires an inclusive slice-size range.
func SliceSize[T any](minimum, maximum int) validation.Validator[[]T] {
	return validation.ValidatorFunc[[]T](func(ctx validation.Context, values []T) validation.Report {
		if len(values) >= minimum && len(values) <= maximum {
			return pass(ctx)
		}
		return fail(ctx, "size", sizeParameters(minimum, maximum))
	})
}

// MapSize requires an inclusive map-size range.
func MapSize[K comparable, V any](minimum, maximum int) validation.Validator[map[K]V] {
	return validation.ValidatorFunc[map[K]V](func(ctx validation.Context, values map[K]V) validation.Report {
		size := len(values)
		if size >= minimum && size <= maximum {
			return pass(ctx)
		}
		return fail(ctx, "size", sizeParameters(minimum, maximum))
	})
}

func sizeParameters(minimum, maximum int) map[string]string {
	return map[string]string{"minimum": strconv.Itoa(minimum), "maximum": strconv.Itoa(maximum)}
}

// Unique uses Go comparable equality after collection admission. NaNs do not
// compare equal and therefore can repeat; dynamic noncomparable any values
// become validator_panic findings through the function adapter. It does not
// canonicalize values or implement JSON structural equality.
func Unique[T comparable]() validation.Validator[[]T] {
	return uniqueValidator[T](isStringType[T]())
}

func uniqueValidator[T comparable](checkStringLimit bool) validation.Validator[[]T] {
	return validation.ValidatorFunc[[]T](func(ctx validation.Context, values []T) validation.Report {
		if len(values) > ctx.Limits().MaxCollectionSize {
			return fail(ctx, "collection_limit", nil)
		}
		seen := make(map[T]struct{}, len(values))
		for index, value := range values {
			if checkStringLimit && valueExceedsStringLimit(ctx, value) {
				return fail(ctx.WithPath(validation.Index(index)),
					"string_limit", nil)
			}
			if _, exists := seen[value]; exists {
				return fail(ctx.WithPath(validation.Index(index)), "unique", nil)
			}
			seen[value] = struct{}{}
		}
		return pass(ctx)
	})
}

// Items applies a validator to slice items in index order.
func Items[T any](mode validation.Mode, validator validation.Validator[T]) validation.Validator[[]T] {
	return validation.ValidatorFunc[[]T](func(ctx validation.Context, values []T) validation.Report {
		if len(values) > ctx.Limits().MaxCollectionSize {
			return fail(ctx, "collection_limit", nil)
		}
		report := pass(ctx)
		for index, value := range values {
			current := validator.Validate(ctx.WithPath(validation.Index(index)), value)
			report = report.Merge(current)
			if current.ContextError() != nil ||
				(mode == validation.ShortCircuit && current.Err() != nil) {
				break
			}
		}
		return report
	})
}

// Keys validates ordered map keys in sorted order for deterministic reports.
// Non-reflexive floating keys (NaNs) fail key admission before child invocation.
func Keys[K cmp.Ordered, V any](mode validation.Mode,
	validator validation.Validator[K],
) validation.Validator[map[K]V] {
	return keysValidator[K, V](mode, validator, isStringType[K]())
}

func keysValidator[K cmp.Ordered, V any](
	mode validation.Mode,
	validator validation.Validator[K],
	checkStringLimit bool,
) validation.Validator[map[K]V] {
	return validation.ValidatorFunc[map[K]V](func(ctx validation.Context, values map[K]V) validation.Report {
		if len(values) > ctx.Limits().MaxCollectionSize {
			return fail(ctx, "collection_limit", nil)
		}
		keys := make([]K, 0, len(values))
		oversizedKey := false
		for key := range values {
			if key != key {
				return fail(ctx, "key", nil)
			}
			if checkStringLimit && valueExceedsStringLimit(ctx, key) {
				oversizedKey = true
			}
			keys = append(keys, key)
		}
		if oversizedKey {
			return fail(ctx, "string_limit", nil)
		}
		slices.Sort(keys)
		report := pass(ctx)
		for _, key := range keys {
			current := validator.Validate(ctx.WithPath(validation.Key(fmt.Sprint(key))), key)
			report = report.Merge(current)
			if current.ContextError() != nil ||
				(mode == validation.ShortCircuit && current.Err() != nil) {
				break
			}
		}
		return report
	})
}

// Values validates map values in sorted key order and locates each finding at
// its typed key segment. Non-reflexive floating keys (NaNs) fail key admission
// before child invocation because they cannot identify a deterministic keyed
// location. Values come from actual map entries, not a second key lookup.
// Inputs must remain unchanged.
func Values[K cmp.Ordered, V any](mode validation.Mode,
	validator validation.Validator[V],
) validation.Validator[map[K]V] {
	return valuesValidator[K, V](mode, validator, isStringType[K]())
}

func valuesValidator[K cmp.Ordered, V any](
	mode validation.Mode,
	validator validation.Validator[V],
	checkStringLimit bool,
) validation.Validator[map[K]V] {
	return validation.ValidatorFunc[map[K]V](func(ctx validation.Context, values map[K]V) validation.Report {
		if len(values) > ctx.Limits().MaxCollectionSize {
			return fail(ctx, "collection_limit", nil)
		}
		type entry struct {
			key   K
			value V
		}
		entries := make([]entry, 0, len(values))
		oversizedKey := false
		for key, value := range values {
			if key != key {
				return fail(ctx, "key", nil)
			}
			if checkStringLimit && valueExceedsStringLimit(ctx, key) {
				oversizedKey = true
			}
			entries = append(entries, entry{key: key, value: value})
		}
		if oversizedKey {
			return fail(ctx, "string_limit", nil)
		}
		slices.SortFunc(entries, func(left, right entry) int { return cmp.Compare(left.key, right.key) })
		report := pass(ctx)
		for _, entry := range entries {
			current := validator.Validate(
				ctx.WithPath(validation.Key(fmt.Sprint(entry.key))), entry.value,
			)
			report = report.Merge(current)
			if current.ContextError() != nil ||
				(mode == validation.ShortCircuit && current.Err() != nil) {
				break
			}
		}
		return report
	})
}
