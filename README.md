# validation

[![CI](https://github.com/faustbrian/go-validation/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-validation/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-validation/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Mutation](https://img.shields.io/badge/mutation-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-checked_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-validation.svg)](https://pkg.go.dev/github.com/faustbrian/go-validation)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-validation?sort=semver)](https://github.com/faustbrian/go-validation/releases)
[![Go](https://img.shields.io/badge/go-1.26.6-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`validation` is a typed, transport-neutral validation package for Go 1.26
and later. Ordinary functions and `Validator[T]` are the primary API. Reports
retain stable paths and rule codes without retaining rejected values.

## Five-minute quickstart

```go
package main

import (
	"fmt"

	validation "github.com/faustbrian/go-validation"
	"github.com/faustbrian/go-validation/rules"
)

func main() {
	ctx, _ := validation.NewContext(validation.DefaultLimits())
	validator := validation.All(validation.CollectAll,
		rules.RuneLength(3, 40),
		rules.Prefix("usr_"),
	)
	report := validator.Validate(ctx.WithPath(validation.Field("username")), "x")
	for _, violation := range report.Violations() {
		fmt.Println(violation.Path(), violation.Code())
	}
	// Output:
	// username rune_length
	// username prefix
}
```

Use `validation.Value[T]` when input presence matters:

```go
missing := validation.Missing[string]()
null := validation.Null[string]()
empty := validation.Present("")
_ = []validation.Value[string]{missing, null, empty}
```

Core validators never perform I/O. Use `AsyncValidator[T]` and `AsyncAll` for
context-aware external checks. Reflection is optional and isolated in
`structplan`; typed plans require no tags or registry.

## Documentation

- [Documentation index](docs/README.md)
- [API and packages](docs/api.md)
- [Rule catalog](docs/rules.md)
- [Normative semantics](docs/semantics.md)
- [Specification decisions](docs/specification-decisions.md)
- [Specification conformance](specification/README.md)
- [Guides](docs/guides.md)
- [Security model](docs/security.md)
- [Laravel and cline/struct adoption](docs/adoption.md)
- [Compatibility](docs/compatibility.md)

## Local verification

```sh
make check
```

The Makefile delegates to the same released `golib` contract used by CI.
Hosted CI is a release integrator's final external verification step, not a
prerequisite for local development.

For ecosystem-wide selection and ownership guidance, see the versioned
[Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.3.0/docs/ecosystem/README.md)
and its [Foundations family](https://github.com/faustbrian/go-library-tools/blob/v1.3.0/docs/ecosystem/design-language.md#package-families-and-selection).

## License

MIT. See [LICENSE](LICENSE).
