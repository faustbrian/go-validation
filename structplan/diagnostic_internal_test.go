package structplan

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestConstructionFormatterPreservesPrivateUnknownCause(t *testing.T) {
	const marker = "privateConstructionMarker"
	cause := errors.New(marker)
	err := privateConstructionError(cause)
	if err.Error() != "validation plan construction failed" ||
		strings.Contains(fmt.Sprint(err), marker) {
		t.Fatalf("formatter exposed unknown cause: %v", err)
	}
	if !errors.Is(err, cause) || !errors.Is(errors.Unwrap(err), cause) {
		t.Fatal("formatter discarded the inspectable original cause")
	}
}
