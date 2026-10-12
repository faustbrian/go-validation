package validation

import "fmt"

// Limits bounds validation work performed on untrusted input.
type Limits struct {
	// MaxDepth caps compile-time reflected struct depth (root zero), not runtime pointer/interface cycle traversal; default 32.
	MaxDepth int
	// MaxCollectionSize caps admitted item/key/unique work and supplied compositor/service positions, including nil; default 10000.
	MaxCollectionSize int
	// MaxStringLength caps checked input/catalog bytes before parsing or escaping; default 65536, not all escaped-output allocation.
	MaxStringLength int
	// MaxViolations caps retained findings, default 100; hidden blocking and terminal state survive truncation.
	MaxViolations int
	// MaxPathLength caps human-rendered path bytes at report admission, default 1024; not total path construction or escaped-pointer bytes.
	MaxPathLength int
	// MaxMetadataEntries caps Context metadata and violation parameter entries, default 16.
	MaxMetadataEntries int
	// MaxMetadataKeyLength caps Context keys and machine diagnostic code/parameter names, default 64.
	MaxMetadataKeyLength int
	// MaxMetadataValueLength caps Context locale/operation/metadata and diagnostic parameter bytes, default 256; not confidentiality.
	MaxMetadataValueLength int
	// MaxRegexPatternLength caps Pattern constructor expression bytes, default 1024; not a caller-context timeout.
	MaxRegexPatternLength int
	// MaxCustomConcurrency caps AsyncAll workers per invocation, default 8; concurrent calls multiply the bound.
	MaxCustomConcurrency int
	// MaxStructFields caps typed definitions and total reflected fields visited during compilation, default 256.
	MaxStructFields int
	// MaxTagLength caps each visited validate tag before grammar splitting, default 1024.
	MaxTagLength int
	// MaxCompiledPlans caps distinct reflect.Type entries in an instance Cache, default 256; Clear permits later recompilation.
	MaxCompiledPlans int
}

// DefaultLimits returns conservative limits suitable for application input.
func DefaultLimits() Limits {
	return Limits{
		MaxDepth: 32, MaxCollectionSize: 10_000, MaxStringLength: 65_536,
		MaxViolations: 100,
		MaxPathLength: 1_024, MaxMetadataEntries: 16,
		MaxMetadataKeyLength: 64, MaxMetadataValueLength: 256,
		MaxRegexPatternLength: 1_024, MaxCustomConcurrency: 8,
		MaxStructFields: 256, MaxTagLength: 1_024, MaxCompiledPlans: 256,
	}
}

func (l Limits) validate() error {
	values := []int{l.MaxDepth, l.MaxCollectionSize, l.MaxStringLength,
		l.MaxViolations,
		l.MaxPathLength, l.MaxMetadataEntries, l.MaxMetadataKeyLength,
		l.MaxMetadataValueLength, l.MaxRegexPatternLength,
		l.MaxCustomConcurrency, l.MaxStructFields, l.MaxTagLength,
		l.MaxCompiledPlans}
	for _, value := range values {
		if value <= 0 {
			return fmt.Errorf("%w: limits must be positive", ErrInvalidLimit)
		}
	}
	return nil
}
