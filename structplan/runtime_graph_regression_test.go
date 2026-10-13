package structplan_test

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"testing"
	"time"

	validation "github.com/faustbrian/go-validation/v3"
	"github.com/faustbrian/go-validation/v3/structplan"
)

type runtimeRequiredGraph struct {
	Dynamic any    `validate:"required"`
	Next    string `validate:"required"`
}
type runtimeMeasuredGraph struct {
	Dynamic any `validate:"min=1,max=8"`
}
type recursiveRuntimePointer *recursiveRuntimePointer
type runtimePointerGraph struct {
	Dynamic recursiveRuntimePointer `validate:"required"`
}

// Each historical hang is isolated in a killable subprocess; compilation and
// synchronous validation deliberately expose no context cancellation API.
func TestRuntimeDereferenceGraphsAreBounded(t *testing.T) {
	if probe := os.Getenv("VALIDATION_RUNTIME_GRAPH_PROBE"); probe != "" {
		limits := validation.DefaultLimits()
		limits.MaxDepth = 4
		ctx, err := validation.NewContext(limits)
		if err != nil {
			t.Fatal(err)
		}
		var first, second any
		switch probe {
		case "self-required", "self-measured":
			first = &first
		case "two-node":
			first = &second
			second = &first
		case "deep-acyclic":
			first = "ok"
			for range 5 {
				next := first
				first = &next
			}
		case "compile-pointer-only":
			type untagged struct{ Dynamic recursiveRuntimePointer }
			_, err := structplan.CompileTags[untagged](limits)
			if !errors.Is(err, structplan.ErrCycle) {
				t.Fatalf("recursive pointer compilation = %v", err)
			}
			cache := structplan.NewCache(limits)
			if _, err := structplan.CompileCached[untagged](cache); !errors.Is(err, structplan.ErrCycle) {
				t.Fatalf("cached recursive pointer compilation = %v", err)
			}
			if cache.Len() != 0 {
				t.Fatal("failed compile consumed cache entry")
			}
			if _, err := structplan.CompileCached[runtimeRequiredGraph](cache); err != nil {
				t.Fatalf("next cache compile = %v", err)
			}
			if cache.Len() != 1 {
				t.Fatal("valid compile lost cache entry")
			}
			return
		case "pointer-only":
			var pointer recursiveRuntimePointer
			pointer = recursiveRuntimePointer(&pointer)
			plan, err := structplan.CompileTags[runtimePointerGraph](limits)
			if err != nil {
				t.Fatal(err)
			}
			report := plan.Validate(ctx, runtimePointerGraph{pointer})
			if !report.HasCode("cycle") || report.Err() == nil {
				t.Fatalf("pointer-cycle report = %v", report)
			}
			return
		default:
			t.Fatalf("unknown bounded probe %q", probe)
		}
		if probe == "self-measured" {
			plan, err := structplan.CompileTags[runtimeMeasuredGraph](limits)
			if err != nil {
				t.Fatal(err)
			}
			report := plan.Validate(ctx, runtimeMeasuredGraph{first})
			if report.Len() != 1 || !report.HasCode("cycle") || report.Violations()[0].Path().String() != "Dynamic" {
				t.Fatalf("measured cycle report = %v", report)
			}
		} else {
			plan, err := structplan.CompileTags[runtimeRequiredGraph](limits)
			if err != nil {
				t.Fatal(err)
			}
			report := plan.Validate(ctx, runtimeRequiredGraph{Dynamic: first})
			findings := report.Violations()
			if probe == "deep-acyclic" {
				if len(findings) != 1 || findings[0].Code() != "required" || findings[0].Path().String() != "Next" {
					t.Fatalf("finite acyclic graph was not preserved: %v", report)
				}
				return
			}
			if len(findings) != 2 || findings[0].Code() != "cycle" || findings[0].Path().String() != "Dynamic" || findings[1].Code() != "required" || findings[1].Path().String() != "Next" {
				t.Fatalf("graph admission report = %v", report)
			}
		}
		return
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for _, probe := range []string{"self-required", "self-measured", "two-node", "deep-acyclic", "pointer-only", "compile-pointer-only"} {
		t.Run(probe, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			command := exec.CommandContext(ctx, executable, "-test.run=^TestRuntimeDereferenceGraphsAreBounded$", "-test.timeout=3s")
			command.Env = append(os.Environ(), "VALIDATION_RUNTIME_GRAPH_PROBE="+probe)
			output, err := command.CombinedOutput()
			if err != nil {
				t.Fatalf("bounded runtime graph: %v; deadline=%v; output=%s", err, ctx.Err(), output)
			}
		})
	}
}

func TestRuntimeDereferencePreservesAcyclicAndPresenceSemantics(t *testing.T) {
	limits := validation.DefaultLimits()
	plan, err := structplan.CompileTags[runtimeRequiredGraph](limits)
	if err != nil {
		t.Fatal(err)
	}
	empty := ""
	nonempty := "ok"
	var deep any = nonempty
	for range 64 {
		next := deep
		deep = &next
	}
	for _, test := range []struct {
		name  string
		value any
		limit int
		code  string
	}{
		{"deep acyclic interface/pointer chain", deep, 1, ""},
		{"ordinary chain", &nonempty, 3, ""},
		{"compile depth does not cap pointer hops", &nonempty, 1, ""},
		{"scalar at first hop", nonempty, 1, ""},
		{"present pointer to empty", &empty, 3, ""},
		{"nil interface", nil, 1, "required"},
		{"nil pointer", (*string)(nil), 2, "required"},
	} {
		t.Run(test.name, func(t *testing.T) {
			callerLimits := limits
			callerLimits.MaxDepth = test.limit
			ctx, err := validation.NewContext(callerLimits)
			if err != nil {
				t.Fatal(err)
			}
			report := plan.Validate(ctx, runtimeRequiredGraph{Dynamic: test.value, Next: "ok"})
			if test.code == "" {
				if !report.Empty() {
					t.Fatalf("expected admitted graph, got %v", report)
				}
			} else if report.Len() != 1 || !report.HasCode(test.code) {
				t.Fatalf("graph report = %v, want %s", report, test.code)
			}
		})
	}
}
