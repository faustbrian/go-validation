// Package validationtext applies borrowed application catalogs through
// Messages without changing machine identity. It owns no provider IO, worker
// or shutdown; callback privacy/concurrency/latency remain caller obligations.
// Lookup panics/unsafe text yield empty prose, while accepted text is HTML
// escaped. Caller-owned result slices do not retain mutable parameter maps.
package validationtext

import (
	"html"
	"unicode/utf8"

	validation "github.com/faustbrian/go-validation/v2"
)

// Catalog looks up application-facing prose without changing rule semantics.
type Catalog interface {
	// Lookup receives copied parameters but borrowed provider state. Return
	// promptly without IO hidden behind an unbounded sync call; concurrent
	// reuse requires safety. Panic/missing/unsafe text becomes empty Text,
	// without changing machine identity or exposing the panic payload.
	Lookup(locale, code string, parameters map[string]string) (string, bool)
}

// Message pairs immutable machine identity with optional localized prose.
type Message struct {
	// Path is the raw human display path; only translated Text receives HTML escaping.
	Path string
	// Code is the unchanged machine rule identity.
	Code string
	// Text is HTML-escaped accepted catalog prose or empty on refusal; escaping can expand beyond the pre-escape byte bound.
	Text string
}

// Messages translates findings in report order. Missing entries leave Text
// empty rather than embedding core fallback prose. Lookup is synchronous and
// panic-contained, not preemptible. Accepted text is checked before HTML
// escaping, which can expand output. Path/code are unchanged and unescaped;
// captured terminal/truncation state is not included in this finding-only list.
func Messages(ctx validation.Context, report validation.Report,
	catalog Catalog,
) []Message {
	violations := report.Violations()
	messages := make([]Message, 0, len(violations))
	for _, violation := range violations {
		text := safeLookup(catalog, ctx.Locale(), violation.Code(),
			violation.Parameters(), ctx.Limits().MaxStringLength)
		messages = append(messages, Message{Path: violation.Path().String(),
			Code: violation.Code(), Text: text})
	}
	return messages
}

func safeLookup(catalog Catalog, locale, code string,
	parameters map[string]string, maximum int,
) (text string) {
	if catalog == nil {
		return ""
	}
	defer func() {
		if recover() != nil {
			text = ""
		}
	}()
	text, found := catalog.Lookup(locale, code, parameters)
	if !found || len(text) > maximum || !utf8.ValidString(text) {
		return ""
	}
	for _, character := range text {
		if character < ' ' || character == 0x7f {
			return ""
		}
	}
	return html.EscapeString(text)
}
