package validationconfig_test

import (
	"context"
	"errors"
	"testing"

	validation "github.com/faustbrian/go-validation"
	validationconfig "github.com/faustbrian/go-validation/adapters/config"
	"github.com/faustbrian/go-validation/rules"
)

func TestCheckImplementsConfigContractAndPreservesTerminalErrors(t *testing.T) {
	ctx, err := validation.NewContext(validation.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	check := validationconfig.CheckValue("bad", ctx, rules.Email())
	var contract validationconfig.Validator = check
	if err := contract.Validate(); !errors.Is(err, validation.ErrInvalid) {
		t.Fatalf("invalid error = %v", err)
	}
	caller, cancel := context.WithCancel(context.Background())
	cancel()
	terminal := validation.ValidatorFunc[string](func(validation.Context, string) validation.Report {
		return validation.ContextReport(ctx, caller)
	})
	if err := validationconfig.CheckValue("value", ctx, terminal).Validate(); !errors.Is(err, context.Canceled) {
		t.Fatalf("terminal error = %v", err)
	}
}
