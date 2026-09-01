# Compatibility

The minimum supported toolchain is Go 1.26.6. The module has no runtime
third-party dependencies. Linux, macOS, and Windows are supported where the Go
standard library is supported; CI's primary environment is Linux.

The v1 compatibility contract includes exported Go API, sentinel relationships,
stable standard-rule codes, path rendering and JSON-pointer escaping,
declaration-order aggregation, deduplication identity, and transport field
names. Application translations, custom codes, benchmark timings, and
observation backend adapters are outside that contract.

Standards-backed compatibility is governed by the
[specification decision register](specification-decisions.md). A change to
pointer serialization, transport field names, status mapping, package-owned
extension data, or shared projection state requires a decision digest and
compatibility review. Generic URL, email, UUID, hostname, and IP rules are
package-defined syntactic profiles rather than external conformance claims.

`api/baseline.txt` is the mechanical exported-API snapshot. It complements
behavior tests; it does not prove semantic compatibility by itself.
