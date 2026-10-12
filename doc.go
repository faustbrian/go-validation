// Package validation provides typed, deterministic, bounded application
// validation without transport, binding, or persistence dependencies. Build a
// Context and use Validator/rules or an optional structplan; Report.Err is the
// complete success predicate, not Empty or HasErrors alone. Value preserves
// missing/null/present states without coercion, defaulting or input mutation.
// Generic input objects and callback implementations are borrowed, not deep
// copied; purity and concurrent safety remain application obligations. AsyncAll
// owns bounded per-call workers and joins admitted callbacks, but cannot
// preempt noncooperative work. Transport, IO policy and provider shutdown
// remain outside core; limits are selected budgets, not total CPU/heap caps.
package validation
