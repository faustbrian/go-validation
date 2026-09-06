package validationjsonapi_test

import (
	"testing"

	validation "github.com/faustbrian/go-validation"
	validationjsonapi "github.com/faustbrian/go-validation/adapters/jsonapi"
)

func TestErrorsPreservePathsSeverityAndHiddenBlockingState(t *testing.T) {
	limits := validation.DefaultLimits()
	limits.MaxViolations = 1
	report := validation.NewReport(limits).
		Add(validation.NewViolation(validation.RootPath().Append(validation.Field("a/b~c")),
			"warning", validation.Warning, map[string]string{"minimum": "1"}, nil)).
		Add(validation.NewViolation(validation.RootPath(), "blocked", validation.Error, nil, nil))
	document := validationjsonapi.Errors(report)
	if len(document.Errors) != 1 || !document.Meta.Truncated || !document.Meta.HasErrors ||
		document.Errors[0].Status != "200" || document.Errors[0].Code != "warning" ||
		document.Errors[0].Source.Pointer != "/a~1b~0c" ||
		document.Errors[0].Meta.Parameters["minimum"] != "1" ||
		document.Errors[0].Meta.Severity != "warning" {
		t.Fatalf("document = %#v", document)
	}
}

func TestErrorsProjectBlockingSeverity(t *testing.T) {
	report := validation.NewReport(validation.DefaultLimits()).Add(
		validation.NewViolation(validation.RootPath(), "required", validation.Error, nil, nil),
	)
	document := validationjsonapi.Errors(report)
	if document.Errors[0].Status != "422" || document.Errors[0].Meta.Severity != "error" {
		t.Fatalf("document = %#v", document)
	}
}
