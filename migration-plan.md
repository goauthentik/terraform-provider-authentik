# Migrate terraform-provider-authentik from SDKv2 to terraform-plugin-framework

## Status

| Phase | Status | Notes |
| --- | --- | --- |
| 0a — mux-compatible provider schema | Done (`d1c083b`) | `url`/`token` now plain `Optional` with explicit env fallback in `providerConfigure`; `docs/index.md` regenerated. |
| 0b — guardrails | Done (`fe80fd8`) | Test imports rewritten to `terraform-plugin-testing`; `resource.UnitTest`→`resource.Test`; docs-drift CI job added; `CheckDestroy`/`ImportState` added to `authentik_group`, `authentik_application`, `authentik_user`, `authentik_stage_captcha`, `authentik_provider_oauth2`. Also fixed the leaked-machine-path bug in `docs/resources/system_settings.md` permanently (per-resource template override, no longer regenerates differently per contributor). **Not yet run against a live authentik** — no instance available in this environment; `go build`/`go vet`/`gofmt` and the non-network unit tests pass. |
| 0c — rename `pkg/provider`→`pkg/sdkprovider` | Done (`1bc98f32`) | Pure `git mv` + package rename; `main.go` updated. `pkg/provider` is now free. |
| 1a — framework provider | Done (`b306f8c`) | `pkg/provider/provider.go` added, schema byte-identical to SDKv2 post-0a. Shared memoised `helpers.NewAPIClient`/`ClientOptions`/`APIURL` in `pkg/helpers/client.go`; `pkg/sdkprovider/provider.go` now delegates to it. Registers zero resources/data sources so far. |
| 1b — mux server + registry manifest | Done (`148354d`) | `main.go` now serves `tf6muxserver` (SDKv2 upgraded tf5→tf6, muxed with the framework provider). `terraform-registry-manifest.json` + goreleaser wiring added. `TestProviderSchemaMatchesSDKv2` + `TestRegistrySchemasValidate` added. `go tool tfplugindocs` rerun end-to-end against the muxed binary — zero doc diff. Kept `TestProvider`/`InternalValidate` rather than removing it (still-useful orthogonal coverage). |
| 1c — shared framework helpers | Done (`d41e18a`) | `values.go` (StringOrNull/Int32OrNull/BoolOrNull/ListOrNull/StringPtrOrNull, StringPtr/StringPtrEmpty), `httperror.go`, `list.go` (MergeList/MergeStringList/MergeInt32List), `expression.go` (ExpressionType/Value), `validators.go` (OneOf/RelativeDuration), `span.go`. All framework-only, unit tested. `defaults.go`/`describe.go` deferred to 1e. |
| 1d — resource/datasource base + tracing | Done (`804f513`) | `resourceBase`/`dataSourceBase` in `pkg/provider/{resource,datasource}.go`: `Configure` (nil-safe), `ImportState` via passthrough ID, `span()` helper. No consumers yet — own unit tests exercise them directly. `golangci-lint run ./...` clean. |
| 1e — documentation parity | Done (`b880bb3`) | `describe.go` (`Desc`/`WithDefault`/`Generated`/`MarkResourceDeprecated`), `defaults.go` (`StringDefault`/`BoolDefault`/`Int32Default`/`Float64Default` value-exposing wrappers). `golangci-lint run ./...` clean, full non-network test suite green. |
| 2 — pattern-setter resource | Done | `authentik_group` resource (`dc35615`), `authentik_group` data source (`404adb5`), `authentik_application` resource (`34dd87c`) — all three migrated, all verified byte-identical `tfplugindocs` output. See "Discoveries during Phase 2" below for what the plan didn't anticipate. |
| 3 — resource batches | In progress — batch 1 (stages) 20/27 | Parts 1-5 of 7 done (`a0baac9`, `3936085`, `4a391dc`, `954b9a1`, `85656c6`), plus two fix-up commits for defects the parts uncovered in already-migrated code: `cd132d0` (discovery #5, twelve resources) and `71f71e0` (discovery #6, four resources). Migrated so far: dummy, user_delete, user_logout, deny, invitation, endpoints, source, prompt, authenticator_endpoint_gdtc, consent, authenticator_totp, authenticator_static, redirect, mutual_tls, password, account_lockdown, identification, user_write, authenticator_webauthn, authenticator_sms. Remaining 7: captcha, email, authenticator_duo, authenticator_email, authenticator_validate, user_login, prompt_field — i.e. the four H4 secret-carriers, the two `Float64` resources, and `authenticator_validate`'s #935 case, so parts 6-7 are the hard ones. Zero docs drift throughout. Acceptance tests run in CI, not in this environment. |
| 4 — data sources | Not started | |
| 5 — cleanup | Not started | |

## Discoveries during Phase 2 (read before starting Phase 3)

Seven things the original plan didn't spell out. The first four and #7 were found while
porting the three pattern-setter objects; #5 and #6 were found later, in Phase 3 batch 1
parts 4 and 5, and are the two that had already been copied into a dozen resources before
anyone noticed - both are write/read-mapping rules that no docs diff can detect. All are
now established pattern — follow them in every Phase 3 batch rather than rediscovering them
per-resource.

1. **Moving a resource's tests to `pkg/provider` creates an import cycle with
   `pkg/sdkprovider`'s own tests**, because both need the same muxed-provider test
   scaffolding (`ProviderFactories` etc.), and that scaffolding must import both
   provider packages to build the mux server. Fix, now in place: a new `pkg/acctest`
   package holds `MuxServer`/`ProviderFactories`/`ProviderTestFactories`/`PreCheck`/
   `APIClientFromEnv`. Every acceptance test file in `pkg/sdkprovider` is an **external**
   test package (`package sdkprovider_test`) so it can import `pkg/acctest` freely; files
   that also poke package-private internals for white-box unit tests split into an
   internal half (`..._internal_test.go`, `package sdkprovider`, no `pkg/acctest`
   import) and an external half for the acceptance test. The same split applies in
   `pkg/provider` once a resource has both an acceptance test and an internal unit test
   (see `resource_group_internal_test.go` / `resource_group_test.go`).
2. **A resource can't always be deleted wholesale from `pkg/sdkprovider`.** `data_source_
   groups.go` (plural, still SDKv2, migrates in Phase 4) clones `dataSourceGroup()`'s
   schema map at runtime and calls `mapFromGroup()` on every result — the exact
   "clone-and-mutate trick" the Phase 4 section already calls out for `data_source_users.go`
   too. When migrating a singular data source whose plural sibling hasn't moved yet,
   keep the schema-constructor and mapping functions in `pkg/sdkprovider`, and only
   delete the now-dead `*Read`/`ReadContext`-wiring functions plus the registry map
   entry. Check for this cross-file dependency *before* deleting a data source file —
   `data_source_users.go` will hit the same thing when its turn comes.
3. **`StringPtrOrNull`/`Int32PtrOrNull` are only correct for genuinely-nullable API
   fields** (openapi `Nullable*` wrapper types, where `nil`/unset is structurally
   distinct from an empty value). For the far more common case of a plain `*string`
   field that the authentik API always populates with `""` rather than leaving
   structurally absent (H1's "146 TypeString attributes" case), use the prior-aware
   `StringOrNull(prior, res.GetX())` instead — `StringPtrOrNull` will silently collapse
   "prior was null" into `""` and reintroduce the exact bug H1 exists to prevent. This
   was caught by an internal unit test on `authentik_application`, not by the docs-diff
   check — write one for every resource with H1-affected string fields, don't rely on
   docs drift to catch it (docs render the same whether the value is null or `""`).
4. **H3's mutable-id resources show `id = (known after apply)` on every update once
   migrated**, not just ones that change the slug/identifier — `id` is Computed-only
   with no `UseStateForUnknown` (correctly, since it can't safely predict a changed
   slug), and the framework marks *any* unmodified Computed attribute unknown on update
   absent a plan modifier. This is the accepted tradeoff already documented for
   sibling-derived attributes like `grant_types`; it now also applies to the ten H3
   resources's `id`. `Update()` must read the prior id from `req.State`, not
   `req.Plan` (which holds the unknown placeholder) — `req.State.GetAttribute(ctx,
   path.Root("id"), &priorID)`.
5. **An attribute with a schema `Default` must be mapped from the API verbatim, never
   through the prior-aware H1 helpers.** Found in Phase 3 batch 1 part 4, after parts 1-3
   had already shipped the bug. H1 says "use `StringOrNull(prior, v)`", and it is right —
   but only for attributes that *can* be null. H2 makes 318 attributes
   `Optional+Computed+Default`, and those can never legitimately be null in state, so the
   two rules collide. The failure path is `terraform import`: `ImportState` writes only
   the id, so Read seeds its model with every other attribute null, and
   `BoolOrNull(null, false)` / `Int32OrNull(null, 0)` / `StringOrNull(null, "")` all
   return null. The next plan then applies the default and reports a spurious
   `null -> false` update, where SDKv2 wrote the concrete value and planned clean. It
   bites exactly when the API value equals the type's zero value, so it is invisible on
   attributes whose default is a non-zero value (`mode = "flow"`, `token_count = 6`) and
   live on the ones whose default is `false` / `""` / `0` — or on a non-zero default the
   user has set *to* the zero value. The rule:

   > `Default != nil` → `types.BoolValue(res.GetX())`. No `Default` → the prior-aware
   > `*OrNull` helper. Required → concrete, or `MergeList` for lists.

   Two things worth knowing: **the docs-drift gate cannot see this** (docs render
   identically whether the value is null or `false`), and neither can any test the repo
   had before part 4 — it needs an `ImportStateVerify` step or a direct `fromAPI` unit
   test with a null-seeded model. `resource_stage_defaults_internal_test.go` is the
   executable spec; the fix across the twelve already-migrated resources is `cd132d0`.
6. **Optional list attributes must be written as a non-nil empty slice, never as a nil
   one — use `helpers.SliceOrEmpty`, not bare `ElementsAs`.** Found in Phase 3 batch 1
   part 5, alongside #5, and it had also already been copied into four resources. Two
   independent problems, one fix:

   - Every generated API request model gates its list fields on `if !IsNil(o.X)` in
     `ToMap()`. `ElementsAs` writes a **nil** slice for a null list, so the field is
     omitted from the request body entirely; the API leaves absent fields unchanged on
     `PUT`, so removing an optional list from config no longer clears it server-side.
     SDKv2's `CastSlice` always returned a non-nil slice and so always sent `[]`. The
     failure is loud rather than silent — the plan says null, the API keeps the old
     contents, and the framework raises `Provider produced inconsistent result after
     apply` — but it breaks a completely ordinary operation.
   - `ElementsAs` **cannot write an unknown value into a `[]T` target at all**; it
     returns `Received unknown value, however the target type cannot handle unknown
     values`. An `Optional+Computed` list is unknown in the plan whenever config omits
     it, so `authentik_group.users` failed `Create` outright. Treating unknown as "send
     nothing" is right for those: the point of `Optional+Computed` is that the server
     decides when config is silent.

   `SliceOrEmpty` is the list analogue of `StringPtrEmpty` and the counterpart of
   `MergeList` on the read side. Note the asymmetry: **null lists stay null on the way
   in (`MergeList`) but go out as `[]` (`SliceOrEmpty`)** — H1 governs state, the API's
   `!IsNil` gate governs the wire, and they genuinely disagree. Fixed across
   `authentik_group`, `authentik_application`, `authentik_stage_prompt` and
   `authentik_stage_mutual_tls` in `71f71e0`; the group unknown case is pinned by
   `TestResourceGroupToRequest_ListEncodings`, which asserts against `ToMap()` rather
   than the struct field, since `ToMap` is what decides what reaches the wire.
7. **Data source `Default` doesn't exist in the framework at all** — `datasource/schema`
   attribute types have no `Default` field (data sources aren't planned, so there's
   nothing to default at plan time). SDKv2's data-source-level `Default: true` (e.g.
   `authentik_group`'s `include_users`) must be applied by hand in `Read` when the
   config value is null; keep `helpers.WithDefault(...)` in the description text so the
   generated docs still say "Defaults to `true`." even though nothing enforces it
   structurally.

## Context

`pkg/provider` is built on `terraform-plugin-sdk/v2`, which HashiCorp has put into
maintenance mode — it gets no new protocol features and its data model (untyped
`map[string]any`, `d.GetOk` conflating zero with unset, `DiffSuppressFunc`) is the root
cause of several open bugs in this provider. The framework gives us explicit
null/unknown handling, real nested attributes, semantic equality instead of diff
suppression, and protocol 6.

The provider is large (96 resources, 25 data sources, ~17k non-test LOC, ~1036 schema
attributes) but extremely uniform: one resource template repeated 96 times, one data
source template repeated 25 times, and only five custom schema functions in total (three
validators, two diff suppressors), all living in one 114-line file. It also
avoids almost everything hard — zero `CustomizeDiff`, `TypeSet`, `Timeouts`,
`SchemaVersion`, `StateUpgraders`, `StateFunc`, `d.HasChange`, `d.GetRawConfig`, and all
96 importers are plain `ImportStatePassthroughContext`.

Outcome: identical Terraform configuration and state for users, a provider served over
protocol 6, and SDKv2 + `terraform-plugin-sdk` v1 + `hashicorp/go-cty` dropped from
`go.mod`.

## Decisions

| Question | Decision |
| --- | --- |
| Strategy | `terraform-plugin-mux`, incremental — batches of resources over several releases |
| Collections | Keep `ListAttribute`; keep `helpers.ListConsistentMerge`. No state upgraders. |
| Null handling | Real null semantics: null → omit, explicit `""`/`false`/`0` → send |
| Layout | Rename SDKv2 to `pkg/sdkprovider`; framework code goes in `pkg/provider` |

## New dependencies

```
github.com/hashicorp/terraform-plugin-framework            v1.19.0
github.com/hashicorp/terraform-plugin-framework-validators  v0.19.0
github.com/hashicorp/terraform-plugin-framework-jsontypes   v0.2.0
github.com/hashicorp/terraform-plugin-mux                   v0.23.1
github.com/hashicorp/terraform-plugin-testing               v1.16.0
```

`terraform-plugin-go` is already an indirect dependency at v0.31.0 and becomes direct.

---

## Hazards that shape the design

These five are the whole difficulty of the migration. Everything else is translation.

### H1. Null vs empty string, on ~216 attributes

`internal/fwschemadata/value_semantic_equality.go` returns early when the prior value is
null, so **semantic equality cannot bridge null ↔ `""`** — no custom type fixes this.

Today 146 `TypeString`, 51 `TypeList`, 10 `TypeBool` and 9 `TypeInt` attributes are
`Optional` with no `Default`, and the authentik API returns `""` / `[]` for the unset ones.
SDKv2's legacy type system treats `""` and null as the same thing; the framework does not.
A naive port writes `""` into state while the plan says null, and Terraform core raises
`Provider produced inconsistent result after apply` — a hard error, hit on Create, because
every Create tail-calls Read.

The write direction has the same problem in reverse. Three distinct unset encodings exist
today and `ValueStringPointer()` collapses all three to `nil`:

| helper | unset → | effect on a nullable API field |
| --- | --- | --- |
| `GetP[string]` | `nil` | leave unchanged |
| `GetStringP` | `&""` | **clear the field** (this is the #865 fix) |
| `CastSlice` | `[]T{}` | send empty array |

**Fix:** prior-value-aware mapping, in `pkg/helpers`. `StringOrNull(prior types.String, v string) types.String`
returns null when `v == ""` and `prior` is null, and `types.StringValue(v)` otherwise;
likewise `Int32OrNull`, `BoolOrNull`, `ListOrNull[T]`. For the write direction keep the two
behaviours *named*: `StringPtr` (null → nil) and `StringPtrEmpty` (null → `&""`), so the
choice is visible at each of the ~210 call sites instead of implied.

### H2. `Default` requires `Computed: true`

Every attribute type's `ValidateImplementation` raises `nonComputedAttributeWithDefaultDiag`
when a default is set without `Computed`, and that runs during the `GetProviderSchema` RPC —
a fatal error at `terraform plan`, not a lint. So all 318 defaulted attributes become
`Optional + Computed + Default`, flipping the `Computed` bit in the protocol schema for 318
attributes. SDKv2 forbids that same combination, which is why the docs rule in 1e has to be
`Computed && Default == nil`.

Verified safe, so don't work around them: `stringdefault.StaticString` composes fine with
`CustomType: jsontypes.NormalizedType{}` (defaults are applied at the `tftypes` level, via
`resp.PlanValue.ToTerraformValue(ctx)`), and `resource_system_settings.go`'s
`init()`-computed `defaultFlags` is fine because `Schema()` runs long after package init.

### H3. `id` is not injected, and `UseStateForUnknown` on it is wrong for 10 resources

SDKv2 auto-injects a root `id`; the framework does not, and `terraform-plugin-testing`
errors with `no 'id' found in attributes`. All 96 resources and 25 data sources need an
explicit `"id": schema.StringAttribute{Computed: true}`.

`UseStateForUnknown()` on `id` is only safe when the ID is a stable key. Ten resources
derive it from a **mutable** attribute, and pinning the old value there makes apply fail:

- `res.Slug` — `resource_application.go`, `resource_flow.go`, and the six
  `resource_source_{kerberos,ldap,oauth,plex,saml,scim,telegram}.go` (9 files)
- `res.Identifier` — `resource_token.go`

The UUID-keyed ones (`res.Pk`, `res.PbmUuid`, `res.BrandUuid`, `res.TokenUuid`,
`res.LicenseUuid`, `res.Id`) are stable and may use it.

### H4. `resp.State.Set` is wholesale, `d.Set` is incremental

`d.Set` mutates one key on top of the existing `ResourceData`. `resp.State.Set(ctx, &model)`
writes the entire object, so anything the model does not carry becomes **null**. Eighteen
resources never write some of their own attributes during Read — mostly write-only secrets —
and would silently lose them:

```
resource_endpoints_enrollment_token.go        expires, expiring, retrieve_key
resource_token.go                             expires, expiring, retrieve_key
resource_provider_oauth2.go                   jwks_sources
resource_provider_proxy.go                    jwks_sources
resource_provider_ssf.go                      event_retention, jwt_federation_providers, signing_key
resource_rbac_permission_role.go              model, role
resource_rbac_permission_user.go              model, object_id, permission, user
resource_source_kerberos.go                   spnego_keytab, sync_keytab, sync_password
resource_source_ldap.go                       bind_password
resource_source_oauth.go                      consumer_secret
resource_source_telegram.go                   bot_token
resource_stage_authenticator_duo.go           admin_secret_key, client_secret
resource_stage_authenticator_email.go         password
resource_stage_authenticator_endpoint_gdtc.go configure_flow, friendly_name
resource_stage_captcha.go                     private_key
resource_stage_email.go                       password
resource_stage_user_login.go                  geoip_binding, network_binding
resource_task_schedule.go                     app_model
```

plus `resource_certificate_key_pair.go`, where `key_data` is written only when
`CryptoCertificatekeypairsViewPrivateKeyRetrieve` succeeds.

**Fix — and it is the same fix as H1:** seed the model from the prior value before
overwriting it, so unwritten attributes and null-vs-empty both fall out for free.

```go
func (r *xResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data xModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)   // seed
	// fetch, then overwrite only what the API actually returns
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
```

Create and Update seed from `req.Plan` instead. Do **not** reach for `WriteOnly` for the
secrets here — it needs Terraform 1.11+ and is a config-breaking change.

### H5. Both muxed servers receive `ConfigureProvider`

`tf6muxserver` calls `ConfigureProvider` on every underlying server. `providerConfigure`
makes a live `RootAPI.RootConfigRetrieve` call and calls the global `sentry.Init`, so
naive muxing gives two round trips per plan/apply, two `api.APIClient` connection pools,
and a global Sentry client re-initialised mid-configure. The shared, memoised client
constructor in 1a is what prevents this — it is not just tidiness.

---

## Phase 0 — Prep (three separate PRs, no behaviour change)

Order matters: 0b lands **first**. The `ImportState`/`CheckDestroy` steps it adds are the
only tests that detect H1, H3 and H4, and the docs-drift job is the guardrail for
everything after it.

### 0a. Make the SDKv2 provider config schema mux-compatible

This is mandatory, not stylistic. SDKv2's `core_schema.go` sniffs the environment when an
attribute is `Required` with a `DefaultFunc` — it *runs* the func and, if it yields a value,
emits the attribute as `Optional` instead:

```go
if reqd && s.DefaultFunc != nil {
	v, err := s.DefaultFunc()
	if err != nil || (err == nil && v != nil) {
		reqd = false
		opt = true
	}
}
```

So `url` and `token` are currently emitted as **Optional when `AUTHENTIK_URL` is set and
Required when it isn't**. `tf6muxserver` compares provider schemas with a full `cmp.Equal`
on `tfprotov6.Schema` (ignoring only `MinItems`/`MaxItems`), so leaving this in place
produces intermittent `Invalid Provider Server Combination` errors depending on the
developer's shell. This is also the reason `make gen` prefixes `AUTHENTIK_URL=""` — that
hack becomes dead once the `DefaultFunc`s are gone.

In `pkg/provider/provider.go`:

- Drop all three `schema.EnvDefaultFunc` calls.
- `url` and `token` become `Optional: true` (they are already effectively optional —
  the env var satisfies `Required` today).
- Read `AUTHENTIK_URL` / `AUTHENTIK_TOKEN` / `AUTHENTIK_INSECURE` inside
  `providerConfigure`, falling back to the env var when the config value is empty, and
  raise an explicit diagnostic when `url` or `token` resolve to empty.

Also drop the now-dead `AUTHENTIK_URL=""` prefix from the Makefile `gen` target.

This changes `docs/index.md` (`url`/`token` move from Required to Optional) and replaces
Terraform's built-in "argument required" error with our own message. Ship it on its own
so the diff is legible.

### 0b. Guardrails: test harness, docs-drift CI, and the missing test steps

**Land this first.** Three things, no provider code:

1. A pure import-path rewrite across ~104 files:
   - `terraform-plugin-sdk/v2/helper/resource` → `terraform-plugin-testing/helper/resource`
   - `terraform-plugin-sdk/helper/acctest` (v1!) → `terraform-plugin-testing/helper/acctest`
   - `resource.UnitTest` → `resource.Test` (`UnitTest` is deprecated in plugin-testing)

   `terraform-plugin-testing` v1.16.0 still imports SDKv2 and still declares
   `ProviderFactories map[string]func() (*schema.Provider, error)`, so this is genuinely
   zero-edit beyond the imports — `ProviderFactories` keeps working and is swapped in 1b.
   `PreCheck`, `Steps`, `Config`, `Check`, `ComposeTestCheckFunc`, `TestCheckResourceAttr`
   and `ExpectError` are unchanged, so the 343 `TestCheckResourceAttr` calls and every
   `testAcc*` config function are untouched. Dropping v1 `acctest` removes
   `terraform-plugin-sdk v1.17.2` from `go.mod`.
2. The docs-drift CI job (`make gen && git diff --exit-code docs/`). It does not exist today
   and it is the single highest-value guardrail for this migration.

   `docs/` is currently **out of sync with the code**, so committing a regeneration is a
   prerequisite for the gate to pass. Running `tfplugindocs` on an unmodified tree produces
   three diffs unrelated to anything here: attribute reordering in
   `docs/data-sources/provider_oauth2_config.md` and `docs/resources/stage_captcha.md`, and
   — worth fixing on its own — `docs/resources/system_settings.md` has a leaked absolute
   path from another contributor's machine
   (`{{codefile "shell" "/Users/connorpeshek/projects/..."}}`), meaning that page's import
   block has been broken in the published docs. Land the regeneration in this PR so the
   gate starts from a clean tree.
3. `ImportState: true, ImportStateVerify: true` plus `CheckDestroy` on ~5 representative
   resources. Both are currently zero across the whole repo. `ImportStateVerify` asserts
   that Read reproduces post-apply state — i.e. it is a direct H4 detector — and these
   steps must exist *before* anything is migrated, or the first batch has nothing to fail
   against.

### 0c. Rename `pkg/provider` → `pkg/sdkprovider`

Pure `git mv` + package rename + import fixups. Do it before any framework code exists so
that later PRs have small diffs. `pkg/provider` is then free for the framework provider.

---

## Phase 1 — Framework scaffolding

### 1a. `pkg/provider/provider.go` — the framework provider

Implements `provider.Provider` with `Metadata`, `Schema`, `Configure`, `Resources`,
`DataSources`. Mirror the SDKv2 schema exactly as it stands after 0a: `url` (Optional),
`insecure` (Optional bool), `token` (Optional, Sensitive), `headers` (Optional, Sensitive
`MapAttribute` of `types.StringType`).

Byte-for-byte, because mux compares the full `tfprotov6.Schema`:

- Use `MarkdownDescription` with the identical strings. The SDKv2 `init()` sets the global
  `schema.DescriptionKind = StringMarkdown`, so its attributes emit `MARKDOWN`; using
  `Description` on the framework side would emit `PLAIN` and mismatch.
- Leave the framework `schema.Schema`'s own description **empty**. SDKv2's provider block
  emits `""`/`PLAIN`. This costs nothing because `templates/index.md.tmpl` hardcodes the
  provider prose.
- Env fallback must be spelled the same way as 0a — `if data.URL.ValueString() == ""` →
  `os.Getenv`, not `IsNull()` — so the two `Configure` implementations agree.

Factor the URL/client construction out of `pkg/sdkprovider/provider.go` into a shared
function so both providers build byte-identical clients and there is one place that
appends `/api/v3`. Suggested home: `pkg/helpers/client.go`, exposing

```go
func NewAPIClient(ctx context.Context, opts ClientOptions) (*api.APIClient, error)
```

where `ClientOptions` carries `URL`, `Token`, `Insecure`, `Headers`, `Version`, `Testing`.
Move `TestingTransport`, `GetTLSTransport` and `tracingTransport` there too (the
abandoned `tf-framework` branch already has this file as `pkg/helpers/transport.go` —
reuse it, but keep the subpath `/api/v3` logic, which that branch dropped).

**Memoise it on `ClientOptions`** — per H5, both muxed servers call `ConfigureProvider`,
so an unmemoised constructor doubles the `RootConfigRetrieve` round trip, doubles the
connection pool, and re-runs the global `sentry.Init` on every plan and apply.

Split the URL logic out as a pure `func APIURL(raw string) (string, error)` so the seven
`TestProviderConfigure_PathBasedURL` cases become seven table rows against it. That removes
the need for any `schema.TestResourceDataRaw` replacement — hand-building a `tfsdk.Config`
from `tftypes.NewValue` for this would be ~40 lines of noise.

`Configure` sets both `resp.ResourceData` and `resp.DataSourceData` to the same
`*APIClient`. Keep `APIClient` in `pkg/helpers` with an exported accessor so both
provider packages and their in-package tests can reach the client — today's
`APIClient.client` field is unexported, which is why in-package tests construct it
directly, and the replacement must stay constructible from `pkg/helpers`.

### 1b. `main.go` — mux server, and the registry manifest

**Blocker to fix in this PR:** there is no `terraform-registry-manifest.json` in the repo,
and `.goreleaser.yml` does not package one. The registry defaults to
`"protocol_versions": ["5.0"]` when the file is absent, so shipping a protocol-6 binary
without it publishes wrong metadata. Add at the repo root:

```json
{
    "version": 1,
    "metadata": {
        "protocol_versions": ["6.0"]
    }
}
```

and add the glob to **both** `checksum.extra_files` and `release.extra_files` in
`.goreleaser.yml`, matching HashiCorp's framework scaffolding:

```yaml
    - glob: 'terraform-registry-manifest.json'
      name_template: '{{ .ProjectName }}_{{ .Version }}_manifest.json'
```


Per the mux docs: `tf5to6server.UpgradeServer(ctx, sdkprovider.Provider(version, false).GRPCProvider)`,
then `tf6muxserver.NewMuxServer(ctx, providerserver.NewProtocol6(provider.New(version, false)), ...)`,
served with `tf6server.Serve("registry.terraform.io/goauthentik/authentik", ...)` and
`tf6server.WithManagedDebug()` when `-debug` is passed.

Keep the `-version` flag and the `var version string = "dev"` ldflag target —
`.github/workflows/release.yml` shells out to `./binary -version` to name the Sentry
release, and `.goreleaser.yml` sets `-X main.version=`.

In `pkg/sdkprovider/provider_test.go`, replace `providerFactories` with
`ProtoV6ProviderFactories` pointing at **the muxed server**, so acceptance tests exercise
the same composition users get. Two factories, mirroring today's pair:

```go
var providerFactories = map[string]func() (tfprotov6.ProviderServer, error){...}
var providerTestFactories = ... // testing == true, keeps the mock-failed-request path
```

`TestProvider`'s `InternalValidate()` is replaced by two tests, and both are strict
upgrades on it:

- `TestProviderSchemaMatchesSDKv2` — fetch each server's `tfprotov6.Schema` and compare
  them with mux's own comparison options. This is what keeps 0a honest; without it the
  environment-dependent schema problem resurfaces silently.
- A registry-wide sweep calling `schemaResponse.Schema.ValidateImplementation(ctx)` for
  every framework resource and data source, table-driven over the registry rather than 121
  separate tests. This is what catches H2 at build time instead of at `terraform plan`.

The two `schema.TestResourceDataRaw` unit tests
(`TestResourceGroupReadRolesPreserveConfiguredOrder`,
`TestResourceStageAuthenticatorValidateNotConfiguredAction`) exist only to exercise a
mapping function. They get simpler, not harder: build the model struct literally and call
the resource's `toRequest`/`fromResponse` directly, with no Terraform machinery at all.

### 1c. Shared helpers

The two H1 cases that must not be flattened into "null → omit", and which are the reason
the write-direction helpers stay separately named:

- `resource_application.go:96-98` uses `GetStringP` (which returns `&""`, never nil) for
  `group`, `meta_icon` and `meta_launch_url` precisely so that clearing the attribute
  sends `""` and the API clears the field. That is the fix from "Handle launch url and
  icon removal" (#865); sending nil instead silently regresses it. These three use
  `StringPtrEmpty`, and `GetStringP`'s `strings.TrimSpace` must be kept too.
- `resource_stage_authenticator_validate.go` `not_configured_action` must never be omitted
  from the request body (#935, and there is a unit test asserting exactly that).

New framework-side files in `pkg/helpers`. The SDKv2-coupled `errors.go`, `resource.go` and
`wrapper.go` stay untouched until Phase 5, so the new functions get distinct names rather
than a parallel package:

- **`values.go`** — the H1 core. Read direction: `StringOrNull(prior types.String, v string)`,
  `StringPtrOrNull`, `Int32OrNull`, `BoolOrNull`, `ListOrNull[T]`. Write direction:
  `StringPtr` (null → nil) and `StringPtrEmpty` (null → `&""`).
- **`httperror.go`** — `IsNotFound(hr *http.Response) bool` and
  `HTTPError(hr *http.Response, err error) diag.Diagnostics`. Deliberately *not* a
  `HTTPToDiag` port: only `Read` acts on `IsNotFound` by calling
  `resp.State.RemoveResource(ctx)`. A 404 during Update becomes an error (today it silently
  drops the resource and reports success — changelog it); a 404 during Delete stays a
  success.
- **`list.go` additions** — `MergeStringList(ctx, prior types.List, api []string) (types.List, diag.Diagnostics)`
  and friends, delegating to the existing SDK-free `ListConsistentMerge`. The algorithm is
  unchanged; the `old` argument now comes from `req.State`/`req.Plan` and must be null-aware.
- **`expression.go`** — `ExpressionType`/`ExpressionValue` implementing
  `basetypes.StringTypable` and `StringValuableWithSemanticEquals`, replacing
  `DiffSuppressExpression` on 18 attributes. The comparison must be **symmetric**:
  `strings.TrimRight(a, "\n") == strings.TrimRight(b, "\n")`. The existing
  `DiffSuppressExpression` trims only the config side, and the framework calls
  `proposedNewValuable.StringSemanticEquals(ctx, priorValuable)` — a direct port compares
  the wrong sides.
- **`validators.go`** — `OneOf[T ~string](allowed []T) validator.String` wrapping
  `stringvalidator.OneOf` (85 call sites), and `RelativeDuration() validator.String`.
  `ValidateJSON` is not needed — `jsontypes.Normalized` implements `ValidateAttribute` —
  but keep appending the `JSONDescription` prose or 15 doc files change.
- **`defaults.go` + `describe.go`** — see "Documentation parity" below.
- **`span.go`** — `Span(ctx, name) func()` for the `defer` tracing pattern.

`pkg/helpers/paginate.go` and `ListConsistentMerge` carry over untouched. `wrapper.go`
(`SetWrapper`, `GetP`, `GetIntP`, `GetStringP`, `GetJSON`, `SetJSON`, the `ResourceData`
interface) is deleted in Phase 5 once nothing calls it. Use `types.Int32` /
`Int32Attribute` — both `TypeInt` and `Int32Attribute` serialise to protocol `number`, so
it is a free change — but note `ValueInt32Pointer()` is *not* equivalent to `GetIntP`
(H1 again).

### 1d. Resource base + Sentry tracing

`pkg/provider/resource.go` and `datasource.go`:

- A `resourceBase` struct embedded by every resource, providing `Configure` (with the
  mandatory `req.ProviderData == nil` check, per the framework's configure docs) and
  `ImportState` via `resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)`.
  A `dataSourceBase` does the same minus `ImportState`.
- Sentry spans. Do **not** port `tracing.go`'s wrapper approach. `resource.Resource` has
  nine optional companion interfaces (`ResourceWithConfigure`, `ResourceWithImportState`,
  `ResourceWithConfigValidators`, `ResourceWithValidateConfig`, `ResourceWithModifyPlan`,
  `ResourceWithUpgradeState`, `ResourceWithMoveState`, `ResourceWithIdentity`,
  `ResourceWithUpgradeIdentity`), and a wrapper type either implements one or it does not —
  Go has no way to make that conditional on the wrapped value. A wrapper that implements
  `ResourceWithUpgradeState` for a resource with no upgraders makes the framework call it
  and fail; a wrapper that omits `ResourceWithImportState` silently breaks import.
  Instead put a one-line span helper on `resourceBase` and open each CRUD method with
  `defer r.span(ctx, "create")()`. This is exactly equivalent to today's behaviour —
  `tr()` does not propagate the span's context into the wrapped call either.

### 1e. Documentation parity

`init()` in the SDKv2 provider installs a `schema.SchemaDescriptionBuilder` that appends
`" Defaults to \`x\`."` when a default exists and `" Generated."` when an attribute is
computed. Those strings are in all 122 generated doc files today and the framework has no
equivalent global hook — without a replacement, 318 defaults lose their documentation.

Compose the description at **construction** time rather than post-processing the schema:

```go
MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedPolicyEngineModeEnumValues),
                                  helpers.WithDefault(api.POLICYENGINEMODE_ANY)),
```

`defaults.go` still provides value-exposing wrappers (`helpers.StringDefault(v)` etc.,
each embedding the corresponding `stringdefault`/`booldefault`/… value and exposing
`.Value()`), because `defaults.String` has no value accessor of its own.

A whole-schema `Describe(schema.Schema)` walker was the first idea and is worse:
`schema.Attribute` is an interface over value-type structs, so changing a description means
type-switching across six concrete types and *rebuilding* each one — ~80 lines of
reflection-flavoured code re-running on every `GetProviderSchema`, and it still needs the
wrapper types. Construction-time composition is type-checked, has no runtime cost, and lets
the handful of attributes that should *not* say "Generated" say so explicitly.

The predicate matters either way. Per H2, all 318 defaulted attributes are now `Computed`,
so `" Generated."` must be gated on `Computed && Default == nil`. Keying it on `Computed`
alone — which is what the SDKv2 builder does — would newly stamp it onto all 318 and
produce a large, wrong docs diff. tfplugindocs still lists `Optional + Computed` under
"Optional", so no attribute changes section.

Also carry over the `"Category --- Body"` convention in every resource's
`MarkdownDescription` — `templates/resources.md.tmpl` splits on `" --- "` to derive the
registry subcategory, so dropping it silently removes the subcategory from a doc page.

The docs-drift job lands in 0b, before any of this.

Two other doc-visible details: `helpers.MarkDeprecated` maps to `Schema.DeprecationMessage`,
but the SDKv2 version *also* appends a `~>` callout to the description — replicate both or
the page changes. And keep the `JSONDescription` prose on the 15 JSON attributes even though
`jsontypes.Normalized` supersedes `ValidateJSON`, or those 15 pages change.

---

## Phase 2 — Pattern-setter resource

**The deliverable of this PR is the helper package; the objects are its proof.** Pick
objects that force every helper to be designed against a real case, because whatever shape
lands here gets copied 95 times. Three objects, one PR:

- **`authentik_group` (resource).** In one small file: `jsontypes.Normalized` with a
  `Default` on a custom type (`attributes`, `"{}"`), `ListConsistentMerge` ×3 with prior-state
  plumbing, Optional-only lists (H1's list variant), an `Optional+Computed` `list(int32)`
  (`users`), and a `Default: false` bool. It also owns
  `TestResourceGroupReadRolesPreserveConfiguredOrder`, the executable spec for list
  ordering — port that first, as the regression guard.
- **`authentik_group` (data source).** Adds the data-source-only mechanics:
  `ExactlyOneOf` → `datasourcevalidator.ExactlyOneOf`, a `Default: true` bool
  (`include_users`), and the `users_obj` `list(object)`.
- **`authentik_application` (resource).** The cheapest object that forces H1's *string*
  variant — five Optional-only strings the API returns `""` for — plus the nullable-int
  variant (`protocol_provider` via `api.NewNullableInt32`), an `Optional+Computed`
  server-generated `uuid`, an enum `Default`, and uniquely **H3's mutable-`id` trap**
  (`d.SetId(res.Slug)` with `slug` Required and mutable).

Setting the pattern with something tiny like `authentik_policy_expression` (3 attributes)
would mean designing the helpers blind and then copying the wrong ones 95 times.

Sketch of the target shape:

```go
type groupResource struct{ resourceBase }

type groupModel struct {
    ID          types.String         `tfsdk:"id"`
    Name        types.String         `tfsdk:"name"`
    IsSuperuser types.Bool           `tfsdk:"is_superuser"`
    Parents     types.List           `tfsdk:"parents"`
    Users       types.List           `tfsdk:"users"`
    Roles       types.List           `tfsdk:"roles"`
    Attributes  jsontypes.Normalized `tfsdk:"attributes"`
}
```

Points the template must nail down, because all 96 resources repeat them:

- **`id` must be declared explicitly**, and `UseStateForUnknown()` on it only for
  stable-key resources — see H3 for the ten that must not have it. For the 15 resources
  with `int32` PKs, keep `id` a **string** in state (they currently do `strconv.Itoa`) and
  parse it in each CRUD method — changing it to `Int32` would break existing state.
- **Do not reuse one `read()` helper across Create/Read/Update.** All 96 resources
  currently end Create and Update with `return resourceXRead(...)`. If that shared helper
  calls `resp.State.RemoveResource(ctx)` on a 404, Create fails with
  `Missing Resource State After Create`. Split it into
  `fetch() (obj, found bool, diags)` and let only `Read` act on `!found`.
- **`UseStateForUnknown()` on the ~49 `Optional+Computed` attributes that have no default**,
  or users see `(known after apply)` on every plan where SDKv2 showed the prior value.
  On the 318 that are `Optional+Computed` only because of H2 it is a harmless no-op —
  a default already makes the plan value known — so leave it off there rather than adding
  noise. The rule:

  > `UseStateForUnknown()` is safe iff the server value is a pure function of this
  > resource's own config plus immutable creation-time state. Any attribute whose value
  > depends on a **sibling** attribute must not have it.

  `authentik_provider_oauth2.grant_types` is server-derived from `client_type`, so it must
  **not** have it: the plan would assert the stale list, the API would change it, and apply
  would error. Accept `(known after apply)` there, or clear it in `ModifyPlan` when
  `client_type` changes. `client_secret` (generated once, then stable) and
  `authentik_application.uuid` are both fine.
- **Defaults**: ~60 are typed API enum constants (`api.CLIENTTYPEENUM_CONFIDENTIAL`) and
  need an explicit `string(...)` conversion. One is computed at package init from a
  marshalled struct (`defaultFlags` in `resource_system_settings.go`) — that still works,
  it is just a `string` variable passed to `helpers.StringDefault`.
- **404**: only `Read` calls `resp.State.RemoveResource(ctx)`.
- **Create/Update tail-call Read** today. In the framework, set state from the
  create/update response directly and only re-read where the API genuinely returns less
  than the retrieve endpoint — a tail-call Read that overwrites the plan value will trip
  the framework's data-consistency errors, which are hard errors now rather than the
  warnings SDKv2 demoted them to.

Verify the semantic-equality path here too: with `jsontypes.Normalized`, a config value
of `{"a": 1}` against an API response of `{"a":1}` must produce an empty plan, and with
`helpers.Expression`, a heredoc value ending in `\n` against the API's stripped value
must produce an empty plan.

One behaviour change to note in the changelog: semantic equality is a *state-write-time*
mechanism, not a plan-time one. Terraform core computes the proposed new state from config
against prior state and knows nothing about it. So after `terraform import` — state holding
the API's stripped value, config holding the heredoc — the first plan shows a diff and one
no-op Update, then converges. `DiffSuppressFunc` showed no diff at all. Low impact today
(there are no `ImportState` test steps yet), but real.

---

## Phase 3 — Resource batches

Group the remaining resources by family so each PR shares an API surface and a reviewer can
diff them against each other. Two deviations from "easiest first":

- **`authentik_provider_oauth2` + `authentik_system_settings` get their own PR, early.**
  They are the two hardest objects — `list(object)` with hand-rolled converters and
  `ListConsistentMerge[CustomRedirectURI]`, four relative-duration validators,
  sibling-derived `Optional+Computed` (`grant_types`), sensitive computed `client_secret`,
  plus `footer_links` untyped passthrough, the `init()`-computed default, and
  `DeleteContext: schema.NoopContext`. If `allowed_redirect_uris` cannot be represented
  without a breaking change, that must surface in the third PR, not the twelfth.
- **`authentik_user` goes early regardless of family** — 20+ test configs depend on it.
- **The 18 H4 resources get one dedicated PR** so the seed-from-prior-state pattern gets a
  single focused review rather than being spot-checked across six batches.

Otherwise, batch by the `Description: "X --- "` subcategory (Directory / Applications /
Customization / Flows & Stages / System / RBAC / Sources / Endpoints), ~8–12 objects per PR,
which keeps each PR's docs diff and test surface coherent. Rough order:

1. **Stages** (27 files, `resource_stage_*.go`) — mostly flat strings/bools/enums, and the
   only `TypeFloat` attributes (`resource_stage_authenticator_validate.go`,
   `resource_stage_captcha.go` → `Float64Attribute`).
2. **Property mappings** (14 files) — nearly identical; all carry the `expression`
   attribute, so this batch validates `helpers.Expression` at scale.
3. **Policies** (9 files) — `resource_policy_geoip.go` needs care: it is the main user
   of `GetP[bool]`, so it is where the null-semantics change is user-visible.
4. **Sources** (7 files) — all declare a `uuid` Computed attribute alongside the slug `id`.
5. **Providers** (11 files) — the int32-PK cluster. `resource_provider_oauth2.go` is the
   single largest file and owns `allowed_redirect_uris` (`TypeList` of untyped `TypeMap`).
   Keep it as `ListAttribute{ElementType: types.MapType{ElemType: types.StringType}}` to
   preserve config and state exactly; converting it to a real nested attribute changes
   the state type and needs a state upgrader, so treat it as a separate follow-up.
   `resource_system_settings.go` `footer_links` has the same shape and the same rule.
6. **RBAC / outposts / endpoints / everything else** (27 files: 4 RBAC, 4 outpost,
   4 endpoints, and 15 singletons — application, application_entitlement, blueprint,
   brand, certificate_key_pair, enterprise_license, event_rule, event_transport, flow,
   flow_stage_binding, rac_endpoint, system_settings, task_schedule, token, user) — includes the only
   three `ForceNew` files (`resource_outpost_provider_attachment.go`,
   `resource_rbac_permission_role.go`, `resource_rbac_permission_user.go` →
   `RequiresReplace()` plan modifiers), the `MarkDeprecated` resource (→ the
   `DeprecationMessage` field on `schema.Schema`, plus the same `~>` callout appended to
   the description so it still reaches the docs), and `resource_user.go`, whose
   create-only `password` (`d.IsNewResource()`) keeps its current write-once behaviour —
   per H4, `WriteOnly` is off the table here. This batch also holds the only two resources with no update path
   (`resource_rbac_permission_user.go` declares none, `resource_outpost_provider_attachment.go`
   sets `UpdateContext: nil`). `Update` is a required method on `resource.Resource`, so
   implement it as an error diagnostic — every attribute on both is already `ForceNew`, so
   with `RequiresReplace()` in place it is unreachable. `resource_rbac_permission_user.go`
   also has `CreateContext`/`ReadContext` set to `schema.NoopContext`, and
   `resource_system_settings.go` has `DeleteContext: schema.NoopContext` — all become
   empty method bodies (the framework removes state itself on a clean `Delete` return).

Each PR moves the resource files, moves their tests, deletes the SDKv2 originals, and
removes the entries from `pkg/sdkprovider/provider.go`'s `ResourcesMap` while adding them
to `Resources()`. A resource must be served by exactly one of the two servers — the mux
server errors if both claim it, which is a useful safety net.

---

## Phase 4 — Data sources

25 data sources, three of which need real design work rather than translation:

**Use `ListAttribute` with an `ObjectType`, not `ListNestedAttribute`.** SDKv2's
`CoreConfigSchema` routes Computed-only nested lists to *attributes*, not blocks:

```go
default: // SchemaConfigModeAuto
	if schema.Computed && !schema.Optional {
		ret.Attributes[name] = schema.coreConfigSchemaAttribute()
		continue
	}
```

All three nested sites are Computed-only, so they are already protocol attributes of type
`list(object(...))`. `schema.ListAttribute{Computed: true, ElementType: types.ObjectType{AttrTypes: ...}}`
reproduces that byte-for-byte — no state change, no docs change. `ListNestedAttribute`
would emit a protocol `NestedType` instead and *change* the schema.

- **`data_source_users.go` / `data_source_groups.go`** clone and mutate another data
  source's schema map at runtime, then iterate their own schema map inside `Read` to build
  the API query. Framework schemas are static values, so both need their attribute maps
  written out, and the `Read` query-building loop replaced with straight field checks.
  Factor the shared per-user / per-group attribute map into a function returning
  `map[string]schema.Attribute` (and a matching `map[string]attr.Type` for the object
  element) so the singular and plural data sources stay in sync — that is what the
  clone-and-mutate trick was approximating. `data_source_groups.go`'s `Default != nil`
  branch exists purely to work around `GetOk` treating `false` as unset and can be deleted.
  This also fixes a latent panic for free: `data_source_users.go:107,109` and
  `data_source_groups.go:110` do `v.([]string)` on a `TypeList` value from `d.GetOk`, which
  is `[]any` — so setting `groups_by_name`, `groups_by_pk` or `members_by_username` panics
  today. `ElementsAs` makes it compile-time correct; add a test.
- **`data_source_group.go`** `users_obj` already uses a shared `dataSourceGroupMember()`
  schema function; it becomes the same shared attribute/type map pair.
- The `ExactlyOneOf` / `ConflictsWith` / `RequiredWith` sets (data sources only) become
  `datasourcevalidator.ExactlyOneOf` / `Conflicting` / `RequiredTogether` returned from
  `DataSourceWithConfigValidators.ConfigValidators`. Note
  `data_source_application_entitlement.go:28,35` has self-referential `RequiredWith`
  (`app` requires `app`, `name` requires `name`) — almost certainly meant to be
  `RequiredTogether(app, name)` or a no-op; decide explicitly rather than porting the bug.
- **Every data source must declare a Computed `id`.** Only 6 of the 25 declare it today
  (`data_source_application_entitlement.go`, `certificate_key_pair`, `outpost`,
  `policy_binding`, `policy_expression`, `stage_prompt_field`); the other 19 rely on
  SDKv2's implicit `id`, which is in state and listed under Read-Only in every generated
  doc page. Omitting it in the framework silently removes an attribute users may reference.
  Keep the same values, including the synthetic ones — `d.SetId("0")` in the plural data
  sources and `d.SetId("-1")` in the `managed_list` branch of the six
  `data_source_property_mapping_*.go` files.

---

## Phase 5 — Cleanup

- Delete `pkg/sdkprovider` entirely and drop the mux wiring; `main.go` becomes a plain
  `providerserver.Serve(ctx, provider.New(version, false), providerserver.ServeOpts{
  Address: "registry.terraform.io/goauthentik/authentik", Debug: debugMode})`.
- Delete `pkg/helpers/wrapper.go`, and the SDKv2 halves of `errors.go` and `resource.go`.
- Switch the test factories to `providerserver.NewProtocol6WithError`.
- Give the framework provider schema a real description now that mux equality no longer
  constrains it (see 1a).
- Drop `terraform-plugin-sdk/v2`, `terraform-plugin-sdk` (v1), `hashicorp/go-cty` and
  `terraform-plugin-mux` from `go.mod`.
- Regenerate docs and commit; the docs-drift job from 0b keeps them honest afterwards.
- Add a `.golangci.yml`. There is none today, so CI runs golangci-lint's defaults
  (`errcheck`, `govet`, `ineffassign`, `staticcheck`, `unused`) — `errcheck` will start
  flagging unchecked `resp.State.Set` results and `unused` will flag helpers that lose
  their last caller mid-migration.

## Incidental bugs to fix in passing

Cheap while the files are already open, and each is a real defect:

- `pkg/helpers/paginate.go`: on a non-first-page error the loop `continue`s without
  incrementing `page` — an infinite loop on a transient failure.
- `resource_outpost.go:70` uses `context.Background()` instead of the passed `ctx`.
- `data_source_rbac_permission.go:48` and `data_source_certificate_key_pair.go:82` have
  copy-pasted wrong error messages ("No matching flows/groups found").
- `resource_endpoints_connector_agent.go:145-146` sets `name` twice.
- `resource_token.go:88` shadows the `int` builtin with a local variable.

## Verification

Per PR:

1. `make build` and `golangci-lint run`.
2. `make test` against a live authentik (`AUTHENTIK_URL` + `AUTHENTIK_TOKEN`; the
   devcontainer or the docker-compose flow in `README.md:34-52` provides one). This runs
   the real acceptance tests — `TF_ACC=1` is set for the whole suite.
3. `go tool github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs` (the
   `AUTHENTIK_URL=""` prefix goes away in 0a) and inspect `git diff docs/`.
   tfplugindocs builds the binary and runs
   `terraform providers schema -json`, so under mux it sees both halves and unmigrated
   resources render exactly as before. **For a faithful batch the diff should be empty**;
   any change to a resource in the batch is a schema difference you introduced, and any
   change to a resource *outside* the batch is a mux or shared-helper regression. Treat
   this as the gate on every PR. Expected exceptions: `docs/index.md` in 0a (`url`/`token`
   Required → Optional) and nested-attribute rendering for the three data sources in
   Phase 4.

   Set `MarkdownDescription`, not `Description`. The framework emits markdown only when
   `MarkdownDescription` is non-empty; falling back to `Description` flips
   `description_kind` to plain text and reflows every doc page.

For each migrated resource, add a state-compatibility test following the pattern in the
framework's migration/testing docs: step 1 pins the last released provider via
`ExternalProviders` (`source: goauthentik/authentik`, `VersionConstraint` set to the
current release) and applies a config; step 2 switches to `ProtoV6ProviderFactories` with
the same config and asserts `plancheck.ExpectEmptyPlan()` in
`ConfigPlanChecks.PreApply`. That is the only mechanism that actually proves state and
plan parity across the migration, and the repo has nothing like it today.

Beyond the ~5 seeded in 0b, add `ImportState: true, ImportStateVerify: true` to every
migrated resource, or at minimum one per batch. It is the cheapest H4 detector there is —
it asserts that Read reproduces post-apply state exactly.
