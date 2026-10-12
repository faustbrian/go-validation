package structplan

import (
	"fmt"
	"reflect"
	"sync"

	validation "github.com/faustbrian/go-validation/v2"
)

// Cache is an instance-owned bounded cache of immutable compiled tag plans.
type Cache struct {
	mu     sync.RWMutex
	limits validation.Limits
	plans  map[reflect.Type]any
}

// NewCache creates an owned cache without upfront limits validation. Use
// DefaultLimits/positive fields; CompileCached checks capacity before tag
// compilation validates remaining limits. Zero Cache is not initialized.
func NewCache(limits validation.Limits) *Cache {
	return &Cache{limits: limits, plans: make(map[reflect.Type]any)}
}

// CompileCached uses reflect.Type identity and holds the cache lock through
// trusted CompileTags; no application callbacks or IO execute under that lock.
// Failed compilation consumes no entry. The cache owns definitions, not input
// instances or validators, and nil cache rejects with ErrInvalidPlan.
func CompileCached[T any](cache *Cache) (*TagPlan[T], error) {
	if cache == nil {
		return nil, ErrInvalidPlan
	}
	typeOf := reflect.TypeFor[T]()
	cache.mu.Lock()
	defer cache.mu.Unlock()
	if existing, ok := cache.plans[typeOf]; ok {
		return existing.(*TagPlan[T]), nil
	}
	if len(cache.plans) >= cache.limits.MaxCompiledPlans {
		return nil, fmt.Errorf("%w: compiled plans", validation.ErrLimitExceeded)
	}
	plan, err := CompileTags[T](cache.limits)
	if err != nil {
		return nil, err
	}
	cache.plans[typeOf] = plan
	return plan, nil
}

// Len returns the current number of cached plans.
func (cache *Cache) Len() int {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	return len(cache.plans)
}

// Clear removes cached entries under the instance lock. Returned immutable
// plans remain valid; later calls may compile a new plan for the same type.
// Len/Clear on nil receivers panic; no closed state or provider shutdown exists.
func (cache *Cache) Clear() {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.plans = make(map[reflect.Type]any)
}
