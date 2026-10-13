# validation

[![CI](https://github.com/faustbrian/go-validation/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-validation/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-validation/actions/workflows/ci.yml)
[![Mutation](https://img.shields.io/badge/mutation-local_and_scheduled-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-structural_checks_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-validation/v3.svg)](https://pkg.go.dev/github.com/faustbrian/go-validation/v3)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-validation?sort=semver)](https://github.com/faustbrian/go-validation/releases)
[![Go](https://img.shields.io/badge/go-1.27.0-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`validation` is a typed, transport-neutral validation package for Go 1.27
and later. Ordinary functions and `Validator[T]` are the primary API. Reports
retain stable paths and rule codes without retaining rejected values.

## Five-minute quickstart

The v3 major-version module uses the `/v3` import identity. Require the
public v3.0.0 tag to resolve before installing this release:

```sh
go get github.com/faustbrian/go-validation/v3@v3.0.0
```

See the [v3 migration guide](docs/migration.md#adopting-v3) before upgrading
existing validation boundaries.

```go
package main

import (
	"fmt"

	validation "github.com/faustbrian/go-validation/v3"
	"github.com/faustbrian/go-validation/v3/rules"
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

For context-aware validation, use `report.Err() == nil` as the complete
success predicate. Cancellation and deadlines are terminal outcomes rather
than validation findings, so a terminal report can still be `Empty()` and can
retain completed warnings or errors.

Target-oriented integrations live under `adapters/config`, `adapters/http`,
`adapters/jsonapi`, `adapters/jsonrpc`, and `adapters/service`. The original
`validationconfig`, `validationhttp`, `validationjsonapi`, `validationrpc`, and
`validationservice` paths remain compatible during their deprecation period.

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
go test ./...
make check
```

With the checksum-pinned public verifier, `make check` and `make ci` select
bounded local checks; `make check-full` preserves explicit full assurance.
See [contributor verification](CONTRIBUTING.md#verification) for selected 2s
fuzzing, optional single-package mutation, asynchronous lifecycle race
assertions and relevant specification, compatibility and release boundaries.
Hosted CI is a release integrator's final external verification step, not a
prerequisite for local development.

For ecosystem-wide selection and ownership guidance, see the versioned
[Golib ecosystem index](https://github.com/faustbrian/go-library-tools/tree/v1.4.0/docs/ecosystem)
and its [Foundations family](https://github.com/faustbrian/go-library-tools/blob/v1.8.4/docs/ecosystem/design-language.md#package-families-and-selection).

## License

MIT. See [LICENSE](LICENSE).
