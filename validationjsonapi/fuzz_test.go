package validationjsonapi_test

import (
	"encoding/json"
	"testing"

	validation "github.com/faustbrian/go-validation/v3"
	"github.com/faustbrian/go-validation/v3/validationjsonapi"
)

// FuzzProjectionPaths requires JSON marshaling of admitted projected paths/
// codes; it is wire-safety evidence, not a complete JSON:API conformance oracle.
func FuzzProjectionPaths(f *testing.F) {
	f.Add("a/b~c", "required")
	f.Fuzz(func(t *testing.T, field, code string) {
		report := validation.NewReport(validation.DefaultLimits()).Add(
			validation.NewViolation(validation.RootPath().Append(validation.Field(field)),
				code, validation.Error, nil, nil))
		if _, err := json.Marshal(validationjsonapi.Errors(report)); err != nil {
			t.Fatal(err)
		}
	})
}
