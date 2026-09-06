package validationjsonrpc_test

import (
	"encoding/json"
	"strings"
	"testing"

	validation "github.com/faustbrian/go-validation"
	validationjsonrpc "github.com/faustbrian/go-validation/adapters/jsonrpc"
)

type safeCause string

func (cause safeCause) Error() string { return string(cause) }

func TestInvalidParamsPreservesSafeStableData(t *testing.T) {
	limits := validation.DefaultLimits()
	limits.MaxViolations = 2
	report := validation.NewReport(limits).
		Add(validation.NewViolation(validation.RootPath().Append(validation.Field("items")).
			Append(validation.Index(0)), "required", validation.Error,
			map[string]string{"minimum": "1"}, safeCause("safe cause"))).
		Add(validation.NewViolation(validation.RootPath(), "warning", validation.Warning, nil, nil)).
		Add(validation.NewViolation(validation.RootPath(), "hidden", validation.Error, nil, nil))
	projected := validationjsonrpc.InvalidParams(report)
	if projected.Code != -32602 || projected.Message != "Invalid params" ||
		len(projected.Data.Violations) != 2 || !projected.Data.Truncated ||
		!projected.Data.HasErrors || projected.Data.Violations[0].Path != "items[0]" ||
		projected.Data.Violations[0].Severity != "error" ||
		projected.Data.Violations[1].Severity != "warning" {
		t.Fatalf("projection = %#v", projected)
	}
	encoded, err := json.Marshal(projected)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "safe cause") {
		t.Fatalf("projection leaked cause: %s", encoded)
	}
}
