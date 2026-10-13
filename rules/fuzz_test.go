package rules_test

import (
	"encoding/binary"
	"math"
	"testing"

	validation "github.com/faustbrian/go-validation/v3"
	"github.com/faustbrian/go-validation/v3/rules"
)

// FuzzUnicodeAndMalformedPrimitives exercises malformed strings through the
// package-defined syntactic profiles and byte/rune rules; their results are
// ignored. The ordered-map seam asserts NaN admission and actual-value delivery,
// not external standards conformance.
func FuzzUnicodeAndMalformedPrimitives(f *testing.F) {
	f.Add("\x00\x00\x00\x00\x00\x00\xf8\x7f", "https://example.com", "a@example.com")
	f.Add("é", "https://example.com", "a@example.com")
	f.Add(string([]byte{0xff, 0xfe}), "://", "not-mail")
	f.Fuzz(func(t *testing.T, text, rawURL, email string) {
		ctx := contextFor(t)
		_ = rules.ByteLength(0, 1024).Validate(ctx, text)
		_ = rules.RuneLength(0, 1024).Validate(ctx, text)
		_ = rules.URL().Validate(ctx, rawURL)
		_ = rules.Hostname().Validate(ctx, text)
		_ = rules.IP().Validate(ctx, text)
		_ = rules.CIDR().Validate(ctx, text)
		_ = rules.Email().Validate(ctx, email)
		_ = rules.UUID().Validate(ctx, text)
		_ = rules.Identifier().Validate(ctx, text)
		var encoded [8]byte
		copy(encoded[:], text)
		key := math.Float64frombits(binary.LittleEndian.Uint64(encoded[:]))
		calls := 0
		child := validation.ValidatorFunc[int](func(ctx validation.Context, value int) validation.Report {
			calls++
			if value != 7 {
				t.Fatalf("map entry became %d", value)
			}
			return validation.NewReport(ctx.Limits())
		})
		report := rules.Values[float64, int](validation.CollectAll, child).Validate(ctx, map[float64]int{key: 7})
		if math.IsNaN(key) {
			if calls != 0 || !report.HasCode("key") {
				t.Fatalf("NaN admission = %d, %v", calls, report)
			}
		} else if calls != 1 || !report.Empty() {
			t.Fatalf("ordinary entry = %d, %v", calls, report)
		}
	})
}
