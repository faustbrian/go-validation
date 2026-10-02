# Compatibility

The minimum supported toolchain is Go 1.27.0. The module has no runtime
third-party dependencies. Linux, macOS, and Windows are supported where the Go
standard library is supported; CI's primary environment is Linux.

The v2 compatibility contract includes exported Go API, sentinel relationships,
stable standard-rule codes, path rendering and JSON-pointer escaping,
declaration-order aggregation, deduplication identity, and transport field
names. Application translations, custom codes, benchmark timings, and
observation backend adapters are outside that contract.

The v2 module uses `github.com/faustbrian/go-validation/v2`. Its intentional
changes from v1 are documented in the [migration guide](migration.md#adopting-v2).
The unmodified released v1.2.0 snapshot remains in `api/v1.2.0.txt`; it is
historical evidence, not the current v2 check baseline.

v1.1.0 additively introduces terminal cancellation/deadline state and five
target-oriented adapter paths. Existing declarations and signatures remain.
The intentional behavioral correction is that `AsyncAll` and service `Chain`
no longer return success-equivalent reports after observing caller context
termination. Legacy and successor transport projections retain their existing
validation-finding shapes; applications route terminal outcomes before using
those projections.

Standards-backed compatibility is governed by the
[specification decision register](specification-decisions.md). A change to
pointer serialization, transport field names, status mapping, package-owned
extension data, or shared projection state requires a decision digest and
compatibility review. Generic URL, email, UUID, hostname, and IP rules are
package-defined syntactic profiles rather than external conformance claims.

`api/baseline.txt` is the mechanical exported-API snapshot. It complements
behavior tests; it does not prove semantic compatibility by itself.
