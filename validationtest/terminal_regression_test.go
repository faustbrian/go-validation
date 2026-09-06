package validationtest_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	validation "github.com/faustbrian/go-validation"
	"github.com/faustbrian/go-validation/validationtest"
)

type terminalRecordingT struct{ messages []string }

func (*terminalRecordingT) Helper() {}

func (test *terminalRecordingT) Fatalf(format string, arguments ...any) {
	test.messages = append(test.messages, fmt.Sprintf(format, arguments...))
}

func TestHelpersRejectTerminalReportsBeforeInterpretingFindings(t *testing.T) {
	vctx, err := validation.NewContext(validation.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	terminal := validation.AsyncAll[int](ctx, vctx, 1)
	validator := validation.ValidatorFunc[int](func(validation.Context, int) validation.Report {
		return terminal
	})

	checks := []func(*terminalRecordingT){
		func(test *terminalRecordingT) { validationtest.RequireValid(test, terminal) },
		func(test *terminalRecordingT) { validationtest.RequireCode(test, terminal, "missing") },
		func(test *terminalRecordingT) { validationtest.RejectMutations(test, vctx, validator, []int{1}) },
		func(test *terminalRecordingT) {
			validationtest.Conformance(test, vctx, validator, []validationtest.Case[int]{{Name: "terminal", Valid: true}})
		},
	}
	for index, check := range checks {
		recorder := &terminalRecordingT{}
		check(recorder)
		if len(recorder.messages) != 1 || !strings.Contains(recorder.messages[0], "did not complete") {
			t.Fatalf("check %d messages = %#v", index, recorder.messages)
		}
	}
}

func TestRequireValidAcceptsCompletedWarnings(t *testing.T) {
	recorder := &terminalRecordingT{}
	report := validation.NewReport(validation.DefaultLimits()).Add(
		validation.NewViolation(validation.RootPath(), "warning", validation.Warning, nil, nil),
	)
	validationtest.RequireValid(recorder, report)
	if len(recorder.messages) != 0 {
		t.Fatalf("messages = %#v", recorder.messages)
	}
}
