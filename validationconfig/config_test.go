package validationconfig_test

import (
	"errors"
	"testing"

	validation "github.com/faustbrian/go-validation/v2"
	"github.com/faustbrian/go-validation/v2/rules"
	"github.com/faustbrian/go-validation/v2/validationconfig"
)

func TestCheckImplementsSmallConfigContract(t *testing.T) {
	ctx, err := validation.NewContext(validation.DefaultLimits(), validation.WithOperation("startup"))
	if err != nil {
		t.Fatal(err)
	}
	check := validationconfig.CheckValue("bad", ctx, rules.Email())
	var contract validationconfig.Validator = check
	if err := contract.Validate(); !errors.Is(err, validation.ErrInvalid) {
		t.Fatalf("Validate() error = %v", err)
	}
}
