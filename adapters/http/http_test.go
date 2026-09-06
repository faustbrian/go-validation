package validationhttp_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	validation "github.com/faustbrian/go-validation"
	validationhttp "github.com/faustbrian/go-validation/adapters/http"
)

func TestProblemWriterAndHookPreserveLegacyBehavior(t *testing.T) {
	limits := validation.DefaultLimits()
	limits.MaxViolations = 1
	report := validation.NewReport(limits).
		Add(validation.NewViolation(validation.RootPath().Append(validation.Field("<token>")),
			"warning", validation.Warning, nil, nil)).
		Add(validation.NewViolation(validation.RootPath(), "blocked", validation.Error, nil, nil))
	problem := validationhttp.FromReport(report)
	if problem.Status != http.StatusUnprocessableEntity || problem.Title != "Validation failed" ||
		!problem.Truncated || len(problem.Errors) != 1 || problem.Errors[0].Severity != "warning" {
		t.Fatalf("problem = %#v", problem)
	}
	recorder := httptest.NewRecorder()
	if err := validationhttp.WriteProblem(recorder, problem); err != nil {
		t.Fatal(err)
	}
	if recorder.Code != http.StatusUnprocessableEntity ||
		recorder.Header().Get("Content-Type") != "application/problem+json" ||
		strings.Contains(recorder.Body.String(), "<token>") {
		t.Fatalf("response = %d %#v %q", recorder.Code, recorder.Header(), recorder.Body.String())
	}
	called := false
	hook := validationhttp.Hook[string](func(_ *http.Request, value string) validation.Report {
		called = value == "input"
		return validation.NewReport(limits)
	})
	if output := hook.Validate(httptest.NewRequest(http.MethodPost, "/", nil), "input"); !called || !output.Empty() {
		t.Fatalf("hook called=%v output=%v", called, output)
	}
}

func TestWarningProblemIsSuccessful(t *testing.T) {
	report := validation.NewReport(validation.DefaultLimits()).Add(
		validation.NewViolation(validation.RootPath(), "warning", validation.Warning, nil, nil),
	)
	problem := validationhttp.FromReport(report)
	if problem.Status != http.StatusOK || problem.Title != "Validation warnings" ||
		problem.Errors[0].Severity != "warning" {
		t.Fatalf("problem = %#v", problem)
	}
}
