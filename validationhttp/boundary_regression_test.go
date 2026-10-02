package validationhttp_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	validation "github.com/faustbrian/go-validation"
	validationhttp "github.com/faustbrian/go-validation/validationhttp"
)

type boundaryWriter struct {
	header                               http.Header
	headerCalls, statusCalls, writeCalls int
	status                               int
	body                                 []byte
}

func (w *boundaryWriter) Header() http.Header    { w.headerCalls++; return w.header }
func (w *boundaryWriter) WriteHeader(status int) { w.statusCalls++; w.status = status }
func (w *boundaryWriter) Write(data []byte) (int, error) {
	w.writeCalls++
	w.body = append(w.body, data...)
	return len(data), nil
}

func TestWriteProblemRejectsStatusBeforeResponseEffects(t *testing.T) {
	for _, status := range []int{0, 99, 1000} {
		writer := &boundaryWriter{header: make(http.Header)}
		err := validationhttp.WriteProblem(writer, validationhttp.Problem{Status: status})
		if !errors.Is(err, validation.ErrInvalidViolation) || writer.headerCalls != 0 ||
			writer.statusCalls != 0 || writer.writeCalls != 0 || len(writer.header) != 0 || len(writer.body) != 0 {
			t.Errorf("status %d: error=%v response effects=%d/%d/%d", status, err, writer.headerCalls, writer.statusCalls, writer.writeCalls)
		}
	}
	for _, status := range []int{100, 999, 200, 422} {
		writer := &boundaryWriter{header: make(http.Header)}
		if err := validationhttp.WriteProblem(writer, validationhttp.Problem{Status: status}); err != nil {
			t.Fatal(err)
		}
		var decoded validationhttp.Problem
		if err := json.Unmarshal(writer.body, &decoded); err != nil {
			t.Fatal(err)
		}
		if writer.status != status || writer.statusCalls != 1 || writer.header.Get("Content-Type") != "application/problem+json" || decoded.Status != status {
			t.Fatalf("valid status %d: response=%d decoded=%d", status, writer.status, decoded.Status)
		}
	}
}

func TestHTTPHookContainsPrivateApplicationPanic(t *testing.T) {
	const marker = "privateHookMarker"
	defer func() {
		if recover() != nil {
			t.Error("application panic escaped function adapter")
		}
	}()
	hook := validationhttp.Hook[int](func(*http.Request, int) validation.Report { panic(marker) })
	report := hook.Validate(nil, 1)
	violations := report.Violations()
	if len(violations) != 1 || violations[0].Code() != "validator_panic" || violations[0].Severity() != validation.Error ||
		!errors.Is(violations[0].Cause(), validation.ErrValidatorPanic) || !errors.Is(report.Err(), validation.ErrInvalid) ||
		strings.Contains(fmt.Sprint(report, report.Err()), marker) {
		t.Fatalf("private panic outcome: report=%v", report)
	}
	calls := 0
	control := validationhttp.Hook[int](func(request *http.Request, value int) validation.Report {
		calls++
		if request != nil || value != 1 {
			t.Error("hook arguments changed")
		}
		return validation.NewReport(validation.DefaultLimits()).Add(validation.NewViolation(validation.RootPath(), "kept", validation.Warning, nil, nil))
	})
	if output := control.Validate(nil, 1); calls != 1 || output.Len() != 1 || !output.HasCode("kept") || output.Err() != nil {
		t.Fatalf("normal hook result changed: calls=%d report=%v", calls, output)
	}
}
