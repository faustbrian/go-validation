# Source documentation audit

## Pre-edit scope and plan

Base `cb316be749e6ed45fd534491b599173961640fbf`, clean detached main-source
worktree; this worker owns Tier A comments/docs only. The unrelated original
chore branch, dirty runtime/adapters/security files and mutation evidence remain
untouched. Outcome: truthful contracts and adjacent rationale, not runtime,
dependency, fixture, policy or standards repair.

Reproduce through `git ls-files '*.go'`, module/package catalogs and native
`source-docs inventory/check --json`: 91 Go files, one v2 module, sixteen
packages including five canonical/retained adapter pairs. Initial 94 exported
comment omissions are structural input, not prose quality or correctness.
Direct consumers are adapters, structplan/rules, observability/text/test helpers
and the published external contracts recorded by the catalog. Named risks:
false purity/IO prevention, shallow generic-value ownership, context cooperation
and admitted-work joining, per-call versus aggregate concurrency, partial/
truncated blocking state, typed-path identity versus rendering, diagnostic
confidentiality versus escaping, reflection/cache/tag bounds and variant drift.

## Complete tracked Go inventory

- `.`: `async.go`, `async_all_test.go`, `benchmark_test.go`, `composition.go`, `composition_test.go`, `context.go`, `context_terminal_regression_test.go`, `core_test.go`, `doc.go`, `errors.go`, `example_test.go`, `fanout_admission_test.go`, `fuzz_test.go`, `hardening_test.go`, `limits.go`, `mutation_contract_test.go`, `path.go`, `projection_conformance_test.go`, `report.go`, `resource_test.go`, `successor_contract_test.go`, `validator.go`, `value.go`, `violation.go`.
- `adapters/config`: `config.go`, `config_test.go`.
- `adapters/http`: `boundary_regression_test.go`, `http.go`, `http_test.go`.
- `adapters/jsonapi`: `fuzz_test.go`, `jsonapi.go`, `jsonapi_test.go`.
- `adapters/jsonrpc`: `rpc.go`, `rpc_test.go`.
- `adapters/service`: `fanout_admission_test.go`, `hook_privacy_regression_test.go`, `service.go`, `service_test.go`.
- `rules`: `collection.go`, `collection_test.go`, `crossfield.go`, `crossfield_test.go`, `fuzz_test.go`, `hardening_test.go`, `helpers.go`, `mutation_contract_test.go`, `numeric.go`, `numeric_test.go`, `presence.go`, `presence_test.go`, `primitives.go`, `primitives_test.go`, `privacy_regression_test.go`, `string.go`, `string_test.go`, `temporal.go`, `temporal_test.go`, `terminal_regression_test.go`.
- `structplan`: `cache.go`, `diagnostic_internal_test.go`, `fuzz_test.go`, `hardening_test.go`, `internal_test.go`, `mutation_contract_test.go`, `plan.go`, `plan_test.go`, `privacy_regression_test.go`, `tags.go`, `terminal_regression_test.go`.
- `validationconfig`: `config.go`, `config_test.go`.
- `validationhttp`: `boundary_regression_test.go`, `http.go`, `http_test.go`.
- `validationjsonapi`: `fuzz_test.go`, `jsonapi.go`, `jsonapi_test.go`.
- `validationobserve`: `observe.go`, `observe_test.go`.
- `validationrpc`: `rpc.go`, `rpc_test.go`.
- `validationservice`: `fanout_admission_test.go`, `hook_privacy_regression_test.go`, `service.go`, `service_test.go`, `terminal_regression_test.go`.
- `validationtest`: `assert.go`, `assert_test.go`, `terminal_regression_test.go`.
- `validationtext`: `text.go`, `text_test.go`.

## Missing comments versus truthful-rationale questions

Repair interface obligations and limits/projection/options/test-helper fields
from source, not generic boilerplate. Audit all public methods and private
composition/async worker/path/report/state/cache/tag/rule/transport/callback
constraints alongside their tests. Compare each canonical adapter against its
retained deprecated variant: distinct nominal/sentinel types are not aliases.
Package comments must state real responsibility and ownership without inventing
complete HTTP/JSON:API/JSON-RPC/RFC primitive conformance.

Questions requiring source investigation: constructors' zero/nil behavior and
all defaults; callback purity/latency/concurrency cannot be enforced by Context;
metadata/paths/causes remain application-asserted content even when bounded.
Check report retention/dedup/terminal precedence, callback-under-lock behavior,
cache ownership, admission before allocation and reflective coercion profiles
before making stronger statements. Preserve pinned authorities/append-only
decision history and historical published evidence.

## Generation, variants and smallest sufficient proof

No generated Go/platform/tag split was found; retained package deprecations
and the frozen ContextReport analyzer directive remain intact. Inventory tests,
fuzz targets, benchmarks, examples and helpers by actual source/input boundary,
not test counts as acceptance. Native docs operation runs existing Examples
across `./...`; no separate docs shell script is declared.

Configured native root API baseline is go-doc: `.`, adapters/config,http,
jsonapi,jsonrpc,service, rules, structplan, validationconfig,http,jsonapi,observe,
rpc,service,test,text in that exact order. Concatenate unmodified native
`go doc -all` stdout and normalize final trailing newlines to one. Producer
Go 1.27.2; public minimum 1.27.0 remains unchanged. Frozen v1 snapshot stays exact.

Only after explicit serialized Go grant: formatting, all 91 native comment-free
AST/signature/literal/directive proof and raw significant-token comparison
(report any raw differences; normalize only optional implicit semicolons before
closing braces when required), parsed/rendered comments, sixteen package renders,
stable configured full API snapshot generation, native API equality and existing
Examples. Source/offline-spec/local-link/marker and pinned changed-added/staged
secret checks plus one ordinary complete Tier A truth/usefulness review precede
an explicit-path Intent/ACK safe-wrapper commit. No race/leak, active fuzz,
mutation, conformance, service, benchmark or external-consumer campaign for
unchanged behavior. Three owned caches/helpers/renders/snapshots are cleaned
on success/failure; parent owns publication, CI and aggregate goal custody.

## Source-grounded boundary coverage

All 91 tracked sources are audited by their inventory families, including
existing admission/terminal/privacy/mutation-contract/successor/transport
regressions and rule/tag/cache helpers. Seven fuzz entrypoints (including two
projection variants), seven root Examples and one multi-workload benchmark
retain their executable bodies and expected outputs. Canonical/retained
adapter bodies are exact after imports at the base; deprecations and distinct
nominal exported identities remain intact. No variants were treated as aliases.

Root Context has no caller cancellation capability and cannot enforce purity
or forbid IO; values are shallow generic copies. AsyncAll owns queue closure,
unique per-job result slots and joins before merging declaration order.
Admission/concurrency bounds are per invocation; caller cancellation stops
admission when selected, child terminal state does not cancel peers, and
noncooperative callbacks can prevent completion. Composition/service wrappers
can replace partial findings on an arbitrary nested panic; typed-plan field
wrappers isolate each field separately.

Reports copy on write and retain hidden blocking/truncation/first-terminal
state. Dedup uses typed path identity plus code/severity/ordered parameters,
not display text or cause identity. Diagnostic validation precedes path
admission; fixed replacement findings remain distinct from input metadata.
NewReport does not validate/default limits; zero Context supplies defaults.
Path length bounds display bytes rather than escaped pointer bytes or prior
construction allocations. Causes, raw paths, parameters, operation labels and
catalog prose remain application-asserted disclosure choices, not inferred
safe data. Catalog text admission precedes HTML expansion; only Text is escaped.

Tag plans check compile-time struct type cycles/depth/visited fields and grammar,
not every runtime rule/kind compatibility. Runtime pointer/interface cycles
can loop in indirect/required traversal before admission; this is an existing
source-audited termination defect, not an intentional safety interpretation.
No hanging probe, process control or runtime repair is authorized in this pass.
Builder/plan/cached plan ownership is explicit: mutable serialized builders,
immutable copied definitions with borrowed callbacks, lock-held trusted tag
compilation (no application callback/IO), and old plans surviving Cache.Clear.

Rule profiles use Go equality/comparison and approximate floating tolerances.
NaN values can repeat under Unique; NaN map-key lookup yields zero values for
Values. Primitive names are package-defined syntactic profiles, not full
external specification conformance. UUID's overbroad RFC9562 comment is
corrected without changing matching or frozen expectations. Size-checking
constructors/callback return boundaries are distinguished from total heap/time
budgets. Transport projectors preserve ordered retained findings and hidden
blocking/truncation, omit terminal/cause state, and rely on caller routing.
HTTP writes headers before encoding, so error can leave a partial response.
RequireValid accepts warnings, while positive conformance Case requires Empty.

Six existing RFC6901/JSON:API1.1/JSON-RPC2.0/RFC9457 decisions and source pins
remain immutable, including append-only history. Open rationale: runtime
dereference graphs are not bounded even though compile-time type recursion is;
NaN keyed-value behavior and empty Any success are existing Go/package behavior
without a stronger domain guarantee. These are documented, not repaired.

## Verification and delivery boundary

The pinned Go 1.27.2 bounded session formatted all 91 sources and proved exact
comment-free AST/signature/literal/directive identity against the base. Raw
significant tokens differ only in `rules/temporal.go`: the canonical multiline
Clock method adds one grammar-optional implicit semicolon immediately before
`}`. Normalizing only that optional token gives exact identity. All 360 comment
groups parsed/rendered; sixteen configured native `go doc -all` outputs formed
a 58,393-byte snapshot stable on two generations. Native `golib api check`
passed. Native `golib docs check` used the configured `^Example` selector with
the declared timeout: seven existing Examples passed, other packages correctly
reported no matching tests. Its link result had zero errors and two declared
exclusions, not universal remote availability.

All 56 tracked test/helper files were reviewed by contract family: controlled
context samplers/admission counters/join barriers, typed-path collision and
safe-text boundaries, retained/dropped severity and terminal precedence,
transport writers and successor nominal identities, grammar/type/value
admission, privacy wrappers, clocks/catalogs/observers and test doubles. Fuzz
bodies keep bounded fixed shapes; they are not cyclic-object or standards
oracles. Existing source/fixtures are retained, not represented as fresh broad
runtime, race, fuzz, mutation, conformance or performance qualification.

The catalog names direct published consumers API Query, International,
Localized, Opening Hours, Service's reference HTTP integration and Temporal.
Their consumed signatures remain exact; no consumer implementation, published
pin, fixture or dependency changed and no external-consumer campaign ran.
Sixty-five tracked dependency/catalog/API-v1/specification/verification inputs
remain byte-exact. Offline native specification checking accepts all six
existing decisions. Native source checking now reports zero structural
findings; no debt/generation/build markers exist in tracked Go sources.

The final review added only NewViolation's DefaultLimits/copy-versus-borrow
clarification; its affected root render/baseline refresh is recorded below.
The original full-source proof still covers unchanged runtime, directives,
Examples and fifteen unaffected rendered packages. No cyclic-graph probe ran:
its nontermination is source-audited, not a runtime observation. Three exact
owned Go caches, parser/helper and rendered outputs were removed after the
initial session; local commit/publication evidence is kept separate.

Final affected proof: `violation.go` has exact raw tokens/AST/literals/directives
and eleven parsed/rendered comment groups; root native GoDoc is stable twice.
The assembled configured sixteen-package baseline is 58,576 bytes and native
API equality passes. Fifteen unchanged package renders and seven unchanged
Examples retain the initial session result, not a new execution claim. The
final helper/render/three-cache session was also removed with absence verified.
Ordinary complete final-diff review covers all source/prose and the native API
snapshot, including both canonical/retained adapter surfaces. No runtime repair,
source lock/dependency/fixture changes or broadened test campaign is included.

Affected local Markdown targets and staged whitespace are valid. Pinned
Gitleaks v8.30.1 with the retained native policy found zero findings on staged
added/new-file content; no policy, ignore entries or findings were suppressed.
Its owned snapshots/private report were removed immediately. Source checking
is zero findings and the six offline decisions remain valid. Publication/CI
and aggregate source-goal completion remain coordinator-owned boundaries.

## Objective Enforcement

The checksum-pinned public Tools v2.4.0 verifier now supplies
`make source-docs-check` as a prerequisite of the existing local and CI Make
routes. Ordinary shared CI selects strict source-documentation checking and
requires its result. Native prerequisites, operation budgets, security and
release controls are retained. Hosted exact-head evidence is recorded by the
coordination ledger; configuration alone does not establish green delivery.
