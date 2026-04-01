# Implementation Plan: AuthN / AuthZ / B2B Completion

**Branch**: `main` | **Date**: 2026-03-12 | **Spec**: [specs/main/spec.md](spec.md)  
**Research**: [specs/main/research.md](research.md)  
**Data Model**: [specs/main/data-model.md](data-model.md)

---

## Summary

Close the remaining implementation gaps in the three foundational services —
**AuthN** (identity), **AuthZ** (authorisation), and **B2B** (organisation management) —
so that the full B2B onboarding, group-insurance purchase, and KYC flows can run
end-to-end in production.

**AuthN** (~90% complete): Add `GetMe` unified endpoint; harden voice-biometric and KYC
integrations; verify no stub remains in FLVE external client.

**AuthZ** (~95% complete): Switch `GetJWKS` to DB-backed `token_configs` with file
fallback; add `RotateTokenConfig` RPC; promote `acl/` placeholder.

**B2B** (~70% complete): Implement full purchase-order lifecycle (Approve/Reject/Fulfill
RPCs + 3 Kafka events); wire `UpdateDepartmentTotalPremium`; replace hardcoded catalog
fallback with DB seed data; add employee bulk import; wire all stub consumers; add
Prometheus metrics.

All changes are additive (new proto fields, new RPCs, new migrations using ALTER only).
No breaking changes to any existing RPC or event schema.

---

## Technical Context

**Language/Version**: Go 1.22+  
**Primary Dependencies**:
  - gRPC / protobuf (`google.golang.org/grpc`, `google.golang.org/protobuf`)
  - GORM v2 on PostgreSQL 17 (GORM-tagged proto-gen models only)
  - Kafka (Confluent-compatible)
  - Redis 7+ (session, JTI blocklist, permission cache)
  - Casbin v2 (AuthZ)
  - Argon2id (AuthN passwords)
  - OpenTelemetry Go SDK (tracing + metrics)

**Storage**:
  - `authn_schema` — users, sessions, user_profiles, otps, api_keys, kyc_verification, documents, voice_sessions
  - `authz_schema` — casbin_rules, roles, user_roles, policies, portal_configs, access_decision_audits, token_configs
  - `b2b_schema` — organisations, org_members, departments, employees, purchase_orders
  - `insurance_schema` — product_plans, products (read-only cross-schema query from B2B catalog)

**Testing**: Go standard `testing`, testify/require, testcontainers-go (PostgreSQL + Kafka + Redis), gomock  
**Target Platform**: Linux container (Docker / Kubernetes)  
**Project Type**: gRPC microservices with generated OpenAPI/REST gateway  
**Performance Goals**: p95 < 100 ms for `GetMe` and `CheckAccess`; p95 < 200 ms for `CreatePurchaseOrder`  
**Constraints**: Zero-downtime deployment; RS256 key rotation without service restart; PO lifecycle transitions are idempotent  
**Scale/Scope**: ~50 000 B2B employee records at launch; ~200 req/s sustained AuthN load

---

## Constitution Check

*Re-evaluated after Phase 1 design (data-model.md complete).*

| # | Principle | Gate Question | Status |
|---|-----------|---------------|--------|
| I | Proto-First | All new RPCs (`GetMe`, `RotateTokenConfig`, PO lifecycle, `BulkImportEmployees`) have proto messages committed before service code. New event messages added to events proto. | ✅ PASS (pre-condition: write proto first) |
| II | Polyglot Ownership | All changes are in Go services (AuthN/AuthZ/B2B). No TS/C# layer involved. | ✅ PASS |
| III | REST API Standard | All new HTTP paths follow Google AIP custom method style (noun resource + `:verb` suffix): `POST /v1/authz/token-configs:rotate`, `POST /v1/b2b/purchase-orders/{id}:approve`, `POST /v1/b2b/purchase-orders/{id}:reject`, `POST /v1/b2b/purchase-orders/{id}:fulfill`, `POST /v1/b2b/employees:bulk-import`. Standard CRUD uses plain nouns (no verbs). OpenAPI 3.x spec auto-generated from proto. | ✅ PASS |
| IV | Event-Driven | Three new Kafka events for PO lifecycle: `PurchaseOrderApproved`, `PurchaseOrderRejected`, `PurchaseOrderFulfilled`. All via publisher interface + outbox pattern. KYC events (`KYCApproved`, `KYCRejected`) added. | ✅ PASS (outbox required in implementation) |
| V | Security & Compliance | `GetMe` requires authenticated JWT (no anonymous). `RotateTokenConfig` gated by `SYSTEM_USER` admin check. `ApprovePurchaseOrder` checks caller role (`B2B_ORG_ADMIN` or `SYSTEM_USER`). PII fields (NID in bulk import) follow existing AES-256 encryption path. | ✅ PASS |
| VI | VSA & Tests | Unit tests required before implementation for each new RPC. Integration tests use testcontainers. Coverage target ≥ 80%. | ✅ PASS (enforced in tasks) |
| VII | Observability | B2B `internal/metrics/` to be added (gap found in research). AuthN and AuthZ already have metrics. All new RPCs export RED metrics via gRPC stats handler. | ✅ PASS (B2B metrics gap closes in this plan) |
| VIII | Multi-Tenancy | Purchase orders carry `business_id` (scoped per org). All new queries filter by org/tenant context from JWT metadata. | ✅ PASS |
| IX | Versioning | All changes additive (new fields = field numbers > existing max). No existing fields removed. Kafka topics versioned at `v1`. | ✅ PASS |
| X | Simplicity | No new service boundaries. No new abstractions beyond domain interface for catalog swap (R-05). Bulk import is a single RPC, not an async job. | ✅ PASS |
| XI | Hybrid SQL Migration | New columns via `ALTER TABLE … ADD COLUMN IF NOT EXISTS`. No `CREATE TABLE`, no `INSERT`. Seed plans go to `db/seeds/`. DB migration file is ALTER-only. | ✅ PASS |
| XII | Platform Surface | `GetMe` and token management: system + all portals (via JWT). PO lifecycle: `business.labaidinsuretech.com` (B2B portal) via custom method RPCs. Bulk import: same. `FulfillPurchaseOrder` gated by `SYSTEM_USER` role. | ✅ PASS |

---

## Project Structure

### Documentation (this feature)

```text
specs/main/
├── plan.md           ← THIS FILE
├── spec.md           ← Feature specification
├── research.md       ← Phase 0 research findings
├── data-model.md     ← Phase 1 data model + proto additions
├── quickstart.md     ← Phase 1 quickstart guide (see below)
├── contracts/        ← OpenAPI diff contracts (generated after proto commit)
└── tasks.md          ← Phase 2 task breakdown (/speckit.tasks)
```

### Source code (affected paths)

```text
proto/insuretech/
├── authn/services/v1/auth_service.proto     ← GetMe RPC
├── authz/services/v1/authz_service.proto    ← RotateTokenConfig RPC
├── b2b/
│   ├── services/v1/b2b_service.proto        ← ApprovePO, RejectPO, FulfillPO, BulkImport
│   ├── entity/v1/purchase_order.proto       ← new lifecycle fields
│   └── events/v1/b2b_events.proto           ← PO lifecycle events

backend/inscore/microservices/
├── authn/
│   ├── internal/grpc/service_iface.go       ← GetMe added
│   ├── internal/service/auth_service.go     ← GetMe implementation
│   └── internal/service/kyc_orchestrator_service.go ← KYC event publishing
├── authz/
│   ├── internal/service/authz_service.go    ← GetJWKS → DB-backed; RotateTokenConfig
│   ├── internal/grpc/authz_handler.go       ← RotateTokenConfig handler
│   └── internal/acl/acl.go                  ← ACL matrix (promoted from placeholder)
└── b2b/
    ├── internal/grpc/b2b_handler.go         ← PO lifecycle + BulkImport handlers
    ├── internal/service/b2b_service.go      ← PO lifecycle logic; catalog cleanup
    ├── internal/repository/purchase_order_repository.go ← lifecycle update queries
    ├── internal/events/publisher.go         ← PO lifecycle event publishers
    ├── internal/consumers/handlers.go       ← wire HandleOrganisationApproved
    └── internal/metrics/                    ← NEW: Prometheus RED metrics

backend/inscore/db/
├── migrations/                              ← ALTER TABLE for PO new columns + indexes
└── seeds/b2b_catalog_seed.sql               ← moved from hardcoded service map
```

---

## Complexity Tracking

> Violations documented per constitution requirement.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Cross-schema query in B2B catalog repo (`insurance_schema.product_plans`) | Unblocks B2B v1.0 without waiting for products service gRPC contract to stabilise | Pure gRPC call is correct long-term but blocking; interface abstraction + config flag sets migration path |
| AuthN depends on AuthZ loopback gRPC for `GetMe` | Circular-dependency appearance — both are internal services | Gateway aggregation would require BFF logic in proxy (Principle II violation) |

---

## Phase 0 — Research (Complete)

See [research.md](research.md) for full findings. All NEEDS-CLARIFICATION items resolved:

| Research Item | Decision |
|---------------|----------|
| R-01 GetMe aggregation point | AuthN calls AuthZ loopback |
| R-02 JWKS DB vs file | DB first, file fallback |
| R-03 Token rotation RPC scope | New `RotateTokenConfig` RPC |
| R-04 PO approval actors | System/B2B_ORG_ADMIN approves; payment event fulfills |
| R-05 Catalog cross-schema | Keep v1 with interface + config flag |
| R-06 Bulk import batch size | Max 500, per-row errors |
| R-07 HandleUserRegistered no-op | No-op v1.0; invitations deferred |
| R-08 B2B observability gap | Add `internal/metrics/` to B2B |

---

## Phase 1 — Design (Complete)

See [data-model.md](data-model.md) for full proto additions and migration DDL.

**Summary of proto changes:**
- `authn/services/v1`: + `GetMe` RPC + messages
- `authz/services/v1`: + `RotateTokenConfig` RPC + messages
- `b2b/services/v1`: + `ApprovePurchaseOrder`, `RejectPurchaseOrder`, `FulfillPurchaseOrder`, `BulkImportEmployees` RPCs + messages
- `b2b/entity/v1`: + 8 new fields on `PurchaseOrder`
- `b2b/events/v1`: + `PurchaseOrderApprovedEvent`, `PurchaseOrderRejectedEvent`, `PurchaseOrderFulfilledEvent`

**Summary of DB migrations (ALTER only):**
```sql
-- b2b_schema.purchase_orders
ALTER TABLE b2b_schema.purchase_orders ADD COLUMN IF NOT EXISTS approved_by       UUID;
ALTER TABLE b2b_schema.purchase_orders ADD COLUMN IF NOT EXISTS approved_at       TIMESTAMPTZ;
ALTER TABLE b2b_schema.purchase_orders ADD COLUMN IF NOT EXISTS approver_notes    TEXT;
ALTER TABLE b2b_schema.purchase_orders ADD COLUMN IF NOT EXISTS rejected_by       UUID;
ALTER TABLE b2b_schema.purchase_orders ADD COLUMN IF NOT EXISTS rejected_at       TIMESTAMPTZ;
ALTER TABLE b2b_schema.purchase_orders ADD COLUMN IF NOT EXISTS rejection_reason  TEXT;
ALTER TABLE b2b_schema.purchase_orders ADD COLUMN IF NOT EXISTS payment_reference TEXT;
ALTER TABLE b2b_schema.purchase_orders ADD COLUMN IF NOT EXISTS fulfilled_at      TIMESTAMPTZ;

-- Indexes
CREATE INDEX IF NOT EXISTS idx_purchase_orders_business_status
    ON b2b_schema.purchase_orders (business_id, status);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_department_status
    ON b2b_schema.purchase_orders (department_id, status);
CREATE UNIQUE INDEX IF NOT EXISTS idx_token_configs_active
    ON authz_schema.token_configs (is_active) WHERE is_active = true;
```

---

## Phase 2 — Implementation Order

Tasks are ordered by dependency. Each `[TAG]` maps to a constitution principle.

### Service ordering (dependency-safe)

```
AuthZ (independent) → AuthN (depends on AuthZ for GetMe) → B2B (depends on both)
```

---

### TRACK A: AuthZ

#### A-1  Proto: add RotateTokenConfig RPC [PROTO]
- Edit `proto/insuretech/authz/services/v1/authz_service.proto`
- Add `RotateTokenConfigRequest`, `RotateTokenConfigResponse`, `RotateTokenConfig` RPC; HTTP: `POST /v1/authz/token-configs:rotate`
- Run `scripts/generate.ps1` → re-gen `gen/go/`, update `api/openapi.yaml`
- `buf lint` + `buf breaking` must pass
- **Deliverable**: clean `buf lint` run on CI

#### A-2  Service: DB-backed GetJWKS [SEC] [TEST]
- Edit `authz/internal/service/authz_service.go` method `GetJWKS`
- Inject `domain.TokenConfigRepository` dependency (already injected in `main.go`)
- Logic: call `r.tokenConfigRepo.GetActive(ctx)` → build JWK → serve
- In-memory cache with 5-minute TTL (use existing `permission_cache.go` pattern or new simple struct)
- Fallback to `AUTHZ_JWKS_PUBLIC_KEY_PATH` file when `GetActive` returns `ErrRecordNotFound`
- Unit tests (table full → DB; table empty → file; cache hit → no DB call)
- **Deliverable**: `authz/internal/service/authz_service_test.go` green

#### A-3  Service: RotateTokenConfig [PROTO] [SEC] [TEST]
- Implement `(s *AuthZService) RotateTokenConfig(ctx, req)` in `authz_service.go`
- Algorithm:
  1. Validate `new_public_key_pem` parses as RSA key
  2. Check `new_kid` uniqueness in `token_configs`
  3. In transaction: `UPDATE token_configs SET is_active=false WHERE is_active=true`; `INSERT new key`
  4. Call `s.InvalidatePolicyCache(ctx, &InvalidatePolicyCacheRequest{})` (clears JWKS cache too)
  5. Create audit log entry
- Add handler in `authz/internal/grpc/authz_handler.go`
- Unit tests: success path; duplicate KID → error; invalid PEM → error
- **Deliverable**: `RotateTokenConfig` end-to-end test with testcontainers PostgreSQL

#### A-4  ACL sub-module [TEST]
- Promote `authz/internal/acl/acl.go` from placeholder
- Define `ResourceActionMatrix` map: `portal → resource → []allowed_actions`
- Seed SYSTEM portal: `{ "svc:policy/read": ["GET","*"], "svc:claims/create": ["POST"] … }`
- Wire `PortalSeeder` to use matrix when seeding fresh-install `p` rules
- Unit tests: matrix lookup returns correct actions per portal
- **Deliverable**: `authz/internal/acl/acl_test.go` green; seeder integration test passes

---

### TRACK B: AuthN (depends on A-2 complete so AuthZ JWKS is stable)

#### B-1  Proto: add GetMe RPC [PROTO]
- Edit `proto/insuretech/authn/services/v1/auth_service.proto`
- Add `GetMeRequest`, `GetMeResponse` (see data-model.md §1.1)
- Import authz entity proto for `Permission` type
- Run `scripts/generate.ps1`
- `buf lint` + `buf breaking` must pass
- **Deliverable**: generated `auth_service_grpc.pb.go` has `GetMe` method

#### B-2  Service: GetMe implementation [SEC] [TEST]
- Create `authz_client` constructor param in `AuthService` (or new `profileService`)
- `GetMe`:
  1. Extract `user_id` from JWT via metadata (`x-user-id`) or `ValidateToken` call
  2. Load `users` + `user_profiles` from repository
  3. Call `AuthZClient.GetUserPermissions(ctx, userId, domain)` via loopback gRPC
  4. Merge and return `GetMeResponse`
- Wire `AuthZServiceClient` in `authn/cmd/server/main.go` (new gRPC client to authz `:50054`)
- Add to `service_iface.go`
- Unit tests: mock AuthZ client; profile not found → `NotFound`; AuthZ error → `Internal`
- Integration test: full flow login → GetMe returns roles
- **Deliverable**: `GET /v1/auth/me` returns `{ user_id, user_type, portal, tenant_id, roles[], permissions[] }`

#### B-3  KYC events [EVENT] [TEST]
- In `kyc_orchestrator_service.go`, add event publish after `ApproveKYC`/`RejectKYC`
- Define event messages in `proto/insuretech/authn/events/v1/` if not present
- Publish to `insuretech.authn.v1.KYCApproved` / `insuretech.authn.v1.KYCRejected`
- Outbox pattern: write to `outbox_events` table before Kafka produce (or use existing event outbox if present)
- Unit tests: mock Kafka producer; event payload matches proto schema
- **Deliverable**: KYC approval publishes verifiable Kafka event in integration test

#### B-4  Voice biometric depth verification [TEST]
- Audit `voice_service.go`: check each of `InitiateVoiceSession`, `SubmitVoiceSample`, `VerifyVoiceSession`
- For each: is the FLVE client call real or a stub (`return nil, nil`)?
- If stub: replace with FLVE HTTP/gRPC call using existing client pattern
- Map FLVE HTTP error codes to gRPC status codes (400→`InvalidArgument`, 404→`NotFound`, etc.)
- Add timeout context (default 10 s) on FLVE calls
- Unit tests: FLVE returns error → correct gRPC status; FLVE timeout → `DeadlineExceeded`
- **Deliverable**: `voice_service_test.go` green with FLVE error-path tests; no `return nil, nil` stubs

---

### TRACK C: B2B (depends on A-* and B-* tracks)

#### C-1  Proto: PO lifecycle + BulkImport + events [PROTO] [EVENT]
- Edit `proto/insuretech/b2b/services/v1/b2b_service.proto`: add 3 PO RPCs + BulkImport RPC
- Edit `proto/insuretech/b2b/entity/v1/purchase_order.proto`: add 8 new lifecycle fields
- Create/edit `proto/insuretech/b2b/events/v1/b2b_events.proto`: add 3 PO events
- Run `scripts/generate.ps1`
- `buf lint` + `buf breaking` must pass
- **Deliverable**: generated `b2b_service_grpc.pb.go` has all 4 new RPCs

#### C-2  DB migration: PO lifecycle columns + indexes [MIGRATE]
- Create `backend/inscore/db/migrations/<timestamp>_b2b_purchase_order_lifecycle.sql`
- Content: 8 `ALTER TABLE … ADD COLUMN IF NOT EXISTS` + 3 `CREATE INDEX IF NOT EXISTS`
- Run `run_migration.ps1 --Target=primary --dry-run` to validate no forbidden keywords
- **Deliverable**: dry-run passes; migration applied to staging DB

#### C-3  Seed: catalog plans moved to DB seeds [MIGRATE]
- Create `backend/inscore/db/seeds/b2b_catalog_seed.sql` with 3 plans as `INSERT … ON CONFLICT DO NOTHING`
- Product IDs and plan IDs must be updated to real UUIDs that exist / will exist in `insurance_schema.product_plans`
- Remove `seededCatalogPlans` map and `fallbackCatalogPlan` function from `b2b_service.go`
- Remove `mergeCatalogWithSeedFallback` and `mergeCatalogMapWithSeedFallback`
- Update `ListPurchaseOrderCatalog` and related service methods to use DB result directly
- Integration test: seed applied → `ListPurchaseOrderCatalog` returns 3 items; no seed → empty list
- **Deliverable**: `seededCatalogPlans` gone from service; test green with testcontainers

#### C-4  Repository: PO lifecycle update queries [TEST]
- Add to `purchase_order_repository.go`:
  - `UpdatePurchaseOrderStatus(ctx, id, newStatus, callerID, notes string) error`
  - `SetPurchaseOrderApproved(ctx, id, approvedBy, notes string) (*PurchaseOrder, error)`
  - `SetPurchaseOrderRejected(ctx, id, rejectedBy, reason string) (*PurchaseOrder, error)`
  - `SetPurchaseOrderFulfilled(ctx, id, fulfilledBy, paymentRef string) (*PurchaseOrder, error)`
- Each query: `UPDATE … WHERE purchase_order_id = ? AND status = <expected_prev_status>`; rows affected = 0 → `FailedPrecondition`
- Unit tests with testcontainers PostgreSQL

#### C-5  Service: PO lifecycle methods [SEC] [EVENT] [TEST]
- Implement in `b2b_service.go`:
  - `ApprovePurchaseOrder`: status guard `SUBMITTED → APPROVED`; call repo; publish `PurchaseOrderApprovedEvent`; call `UpdateDepartmentTotalPremium` (internal)
  - `RejectPurchaseOrder`: status guard `SUBMITTED → REJECTED`; publish `PurchaseOrderRejectedEvent`
  - `FulfillPurchaseOrder`: status guard `APPROVED → FULFILLED`; publish `PurchaseOrderFulfilledEvent`; call `UpdateDepartmentTotalPremium`
- `UpdateDepartmentTotalPremium`: SUM estimated_premium for all FULFILLED orders in department; call `repo.UpdateDepartmentTotalPremium`
- Publish to `EventPublisher` interface (already exists); add 3 new `Publish*` methods  
- Unit tests: valid transition → event published; invalid transition → `FailedPrecondition`; event publish fail → non-fatal log (do not fail the RPC)

#### C-6  Service: BulkImportEmployees [TEST]
- Implement in `b2b_service.go` method `BulkImportEmployees`
- Algorithm: for each row (batched in DB transaction up to 500):
  1. Validate mobile (`^01[3-9]\d{8}$`)
  2. Validate dept exists
  3. Validate plan exists in catalog
  4. Upsert employee by mobile (ON CONFLICT DO NOTHING for duplicate → error row)
  5. Collect per-row result
- Response carries `results[]`, `total_rows`, `success_count`, `failure_count`
- Unit tests: all success; partial failure; oversized batch (> 500 → reject)

#### C-7  Handler: wire all new RPCs [TEST]
- Add to `b2b/internal/grpc/b2b_handler.go`:
  - `ApprovePurchaseOrder`, `RejectPurchaseOrder`, `FulfillPurchaseOrder`, `BulkImportEmployees`
- Each handler: validate required fields → `codes.InvalidArgument`; delegate to `h.svc.*`; map errors
- Integration tests using test gRPC server

#### C-8  Events publisher: PO lifecycle events [EVENT]
- Add to `b2b/internal/events/publisher.go`:
  - `PublishPurchaseOrderApproved`, `PublishPurchaseOrderRejected`, `PublishPurchaseOrderFulfilled`
- Add topic constants: `TopicPurchaseOrderApproved`, `TopicPurchaseOrderRejected`, `TopicPurchaseOrderFulfilled`
- Follow existing `publish()` helper pattern (nil producer → log + no-op)
- Unit tests: producer called with correct topic and payload

#### C-9  Consumer: wire HandleOrganisationApproved [EVENT] [TEST]
- In `consumers/handlers.go` `HandleOrganisationApproved`:
  - Call `s.repo.UpdateOrganisationStatus(ctx, orgID, ACTIVE)`
  - Log success/failure
- Add `UpdateOrganisationStatus(ctx, orgID, status)` to `B2BRepository` interface and `PortalRepository`
- Unit test: event received → repo called; repo error → error returned

#### C-10  Metrics: add Prometheus RED to B2B [OBS]
- Create `b2b/internal/metrics/metrics.go`
- Register gRPC server metrics with `grpc_prometheus` library (already used in AuthN pattern)
- Wire in `b2b/cmd/server/main.go`
- Expose on `/metrics` endpoint (existing health HTTP server pattern)
- **Deliverable**: `curl http://b2b-svc:2112/metrics` returns `grpc_server_handled_total` counters

---

## Quickstart (for next developer)

See [quickstart.md](quickstart.md) for local setup and test run instructions.

**TL;DR** for this feature:
```powershell
# 1. Apply proto changes and regenerate
cd e:\Projects\InsureTech
.\scripts\generate.ps1

# 2. Apply DB migrations
.\run_migration.ps1 --Target=primary

# 3. Apply seeds
.\run_migration.ps1 --Target=primary --seeds-only

# 4. Run unit tests for changed services
go test ./backend/inscore/microservices/authn/... -short
go test ./backend/inscore/microservices/authz/... -short
go test ./backend/inscore/microservices/b2b/...   -short

# 5. Run integration tests (requires Docker)
go test ./backend/inscore/microservices/authn/... -run Integration -count=1
go test ./backend/inscore/microservices/authz/... -run Integration -count=1
go test ./backend/inscore/microservices/b2b/...   -run Integration -count=1
```

---

## Phase 3 — API Consistency Corrections

Full audit of `api/paths/insuretech/*/services/v1/*.yaml` (34 YAML files, 220+ paths) against
**Principle III** (REST API Standard) and the [Google AIP HTTP annotation guidelines](https://google.aip.dev/127).

### Design rules applied

| # | Rule | Compliant | Violation |
|---|------|-----------|-----------|
| R1 | No verb as a plain path segment | `/v1/policies/{id}:cancel` | `/v1/auth/login`, `/v1/knowledge-base/search` |
| R2 | Custom actions use `:verb` suffix on a noun resource (Google AIP) | `/v1/claims/{id}:approve` | `/v1/auth/login` (no resource noun) |
| R3 | Collection names must be **plural kebab-case** nouns | `/purchase-orders` | `/auth/session/current` (singular) |
| R4 | Literal segments must not shadow parameter segments | `/v1/tasks:pending` (custom method) | `/v1/tasks/my-tasks` vs `/{task_id}` |
| R5 | Filters/views use query params; `:verb` is for state-changing actions only | `GET /tasks?assignee=me` | `GET /kyc-verifications:pending` |
| R6 | One authoritative path per resource — no singular/plural duplicates | `/v1/refunds` | `/v1/policies/{id}/refund` + `/v1/policies/{id}/refunds` |
| R7 | Custom method name = single lowercase word or `kebab-case`; no internal verbs | `:bulk-import`, `:rotate` | `:update-status` |
| R8 | POST on a non-collection path requires `:verb` (no bare POST on `{id}`) | `POST /v1/tickets/{id}:reopen` | `POST /v1/tickets/{ticket_id}` |

---

### CRITICAL — Routing conflicts (fix immediately; breaks gateway routing)

Literal path segments that shadow `{param}` segments cause **non-deterministic routing** in
grpc-gateway and every HTTP router that evaluates concrete-before-param.

| Service | Violating path | Conflicts with | Corrected path | File |
|---------|---------------|----------------|----------------|------|
| B2BService | `GET /v1/b2b/purchase-orders/catalog` | `/v1/b2b/purchase-orders/{purchase_order_id}` | `GET /v1/b2b/plan-catalogs` — promote catalog to its own resource | B2BService.yaml |
| BeneficiaryService | `POST /v1/beneficiaries/individual` | `/v1/beneficiaries/{beneficiary_id}` | `POST /v1/beneficiaries` — add `"type": "individual"` to body | BeneficiaryService.yaml |
| BeneficiaryService | `POST /v1/beneficiaries/business` | `/v1/beneficiaries/{beneficiary_id}` | `POST /v1/beneficiaries` — add `"type": "business"` to body | BeneficiaryService.yaml |
| TaskService | `GET /v1/tasks/my-tasks` | `/v1/tasks/{task_id}` | `GET /v1/tasks?assignee=me` — query param filter | TaskService.yaml |
| WorkflowService | `GET /v1/workflow-tasks/my-tasks` | `/v1/workflow-tasks/{task_id}` (implied) | `GET /v1/workflow-tasks?assignee=me` | WorkflowService.yaml |
| SupportService | `GET /v1/knowledge-base/{slug}` + `GET /v1/knowledge-base/{article_id}` | each other — two param segments on same collection | Keep `GET /v1/knowledge-base/{article_id}` only; slug lookup → `?slug=` | SupportService.yaml |
| PaymentService | `POST /v1/payments/webhook/{provider}` | shadows any payment ID equal to "webhook" | `POST /v1/webhooks/payments/{provider}` — or move to gateway-level webhook handler | PaymentService.yaml |
| RefundService | `POST /v1/policies/{policy_id}/refund` (singular) | `POST /v1/policies/{policy_id}/refunds` (plural) in same service | Remove singular `/refund`; keep `/refunds` | RefundService.yaml |

**Proto action required for B2B catalog fix (R4 in this plan):**
The `purchase-orders/catalog` path maps to `ListPurchaseOrderCatalog` RPC which queries
`insurance_schema.product_plans`. Rename the RPC resource to `plan-catalogs`:
```protobuf
// In b2b_service.proto — update HTTP annotation
rpc ListPurchaseOrderCatalog(...) returns (...) {
  option (google.api.http) = { get: "/v1/b2b/plan-catalogs" };  // was: purchase-orders/catalog
}
```
Update `B2BService.yaml`: remove `/v1/b2b/purchase-orders/catalog`; add `/v1/b2b/plan-catalogs`.

---

### HIGH — Verbs as plain path segments (Google AIP R1/R2 violations)

These paths use HTTP verbs or action words as resource segments without the `:verb` colon
syntax. All are existing live endpoints — mark as **v1.x breaking** (require SDK version bump and
client migration notice).

| Service | Violating path | Method | Corrected path | Migration note |
|---------|---------------|--------|----------------|----------------|
| AuthService | `POST /v1/auth/register` | POST | `POST /v1/auth/users` | Create user resource; `register` = create |
| AuthService | `POST /v1/auth/login` | POST | `POST /v1/auth/sessions` | Create session; body: `{credential_type, …}` |
| AuthService | `POST /v1/auth/logout` | POST | `DELETE /v1/auth/sessions/{session_id}` | Destroy session; or `POST /v1/auth/sessions/{id}:invalidate` |
| AuthService | `POST /v1/auth/email/register` | POST | `POST /v1/auth/users:email-register` | Custom method on users collection |
| AuthService | `POST /v1/auth/email/verify` | POST | `POST /v1/auth/users/{user_id}/email:verify` | Custom method on email sub-resource |
| AuthService | `POST /v1/auth/email/login` | POST | `POST /v1/auth/sessions:email` | Custom method on sessions collection |
| DocumentService | `GET /v1/documents/{document_id}/download` | GET | `GET /v1/documents/{document_id}:download` | `/download` is a verb; move to `:download` custom method |
| MediaService | `GET /v1/media/{media_id}/download` | GET | `GET /v1/media/{media_id}:download` | Same pattern — `:download` custom method |
| MediaService | `POST /v1/media/{media_id}/process` | POST | `POST /v1/media/{media_id}:process` | `/process` is a verb; move to `:process` custom method |
| SupportService | `GET /v1/knowledge-base/search` | GET | `GET /v1/knowledge-base:search` | Custom method on collection; `GET …:search?q=` |

**For auth endpoints (`/register`, `/login`, `/logout`, `/email/register`, `/email/login`, `/email/verify`):**
These are already consumed by all portals and SDKs. They are **BREAKING** changes.
Plan of action:
1. Add the corrected paths as **aliases** now (keep old paths alive for 2 minor versions).
2. SDK generator emits deprecation warning on old paths via `x-deprecated: true` in OpenAPI.
3. Remove old paths in `v2` when all portals have migrated.

**Proto HTTP annotation changes (non-breaking; dual-bind using `additional_bindings`):**
```protobuf
rpc Register(RegisterRequest) returns (RegisterResponse) {
  option (google.api.http) = {
    post: "/v1/auth/users"          // NEW canonical
    body: "*"
    additional_bindings { post: "/v1/auth/register" body: "*" }  // DEPRECATED alias
  };
}
rpc Login(LoginRequest) returns (LoginResponse) {
  option (google.api.http) = {
    post: "/v1/auth/sessions"       // NEW canonical
    body: "*"
    additional_bindings { post: "/v1/auth/login" body: "*" }     // DEPRECATED alias
  };
}
```

---

### MEDIUM — Naming inconsistency (fix in current sprint; non-breaking)

| Service | Rule | Current | Corrected | File |
|---------|------|---------|-----------|------|
| AuthService | R3 — singular collection | `GET /v1/auth/session/current` | `GET /v1/auth/sessions/current` | AuthService.yaml |
| AuthZService | R1 — verb resource | `POST /v1/authz/check` | `POST /v1/authz/access-decisions:check` — or keep as `:check` on service root: `POST /v1/authz:check` | AuthZService.yaml |
| AuthZService | R1 — verb resource | `POST /v1/authz/check:batch` | `POST /v1/authz/access-decisions:batch-check` — or `POST /v1/authz:batch-check` | AuthZService.yaml |
| KYCService | R5 — status as `:verb` | `GET /v1/kyc-verifications:pending` | `GET /v1/kyc-verifications?status=pending` | KYCService.yaml |
| RenewalService | R5 — filter as path | `GET /v1/renewals/upcoming` | `GET /v1/renewals?status=upcoming` | RenewalService.yaml |
| BillingService | R3 — singular sub-resource | `GET /v1/orders/{order_id}/invoice` | `GET /v1/orders/{order_id}/invoices` | BillingService.yaml |
| PartnerService | R7 — verb in `:custom` name | `POST /v1/partners/{id}:update-status` | `PATCH /v1/partners/{id}` with `status` field in body; or decompose into `:activate` / `:deactivate` / `:suspend` | PartnerService.yaml |
| TaskService | R8 — POST on `{id}` without `:verb` | `POST /v1/tasks/{task_id}` | Identify the action; add `:assign`, `:complete`, `:reopen` custom methods as needed | TaskService.yaml |

---

### LOW — Style improvements (backlog / v2)

| Service | Issue | Notes |
|---------|-------|-------|
| AuthService | `biometric` and `voice-biometric` are singleton namespaces | `:authenticate`, `:initiate`, `:submit`, `:verify` custom methods are correct AIP style; no change required |
| NotificationService | Compound custom method names `mark-as-read`, `send-bulk` | Hyphenated names are allowed by AIP; consider simplifying to `:read` and `:bulk-send` |
| PaymentService | `GET /v1/payments/provider/{provider}/references/{ref}` | Deep path for provider-reference lookup; consider `GET /v1/payments?provider={p}&provider_ref={r}` |
| AuthZService | `portals/{portal}/config` has duplicate `get:` declarations | Remove one `get:` in AuthZService.yaml (YAML syntax bug) |
| MediaService | `/media/{id}/thumbnail` and `/media/{id}/optimized` | These are computed views; consider `GET /media/{id}?variant=thumbnail\|optimized` or keep as-is |

---

### Missing paths to ADD (planned new endpoints in this feature)

These paths do **not exist** in any YAML yet and must be added as part of the current plan:

| Service | New path | Method | Proto RPC | YAML file |
|---------|----------|--------|-----------|-----------|
| AuthService | `/v1/auth/me` | GET | `GetMe` | AuthService.yaml |
| AuthZService | `/v1/authz/token-configs` | GET | `ListTokenConfigs` (expose existing) | AuthZService.yaml |
| AuthZService | `/v1/authz/token-configs:rotate` | POST | `RotateTokenConfig` (new) | AuthZService.yaml |
| AuthZService | `/v1/authz/cache:invalidate` | POST | `InvalidatePolicyCache` | AuthZService.yaml |
| B2BService | `/v1/b2b/purchase-orders/{purchase_order_id}:approve` | POST | `ApprovePurchaseOrder` | B2BService.yaml |
| B2BService | `/v1/b2b/purchase-orders/{purchase_order_id}:reject` | POST | `RejectPurchaseOrder` | B2BService.yaml |
| B2BService | `/v1/b2b/purchase-orders/{purchase_order_id}:fulfill` | POST | `FulfillPurchaseOrder` | B2BService.yaml |
| B2BService | `/v1/b2b/employees:bulk-import` | POST | `BulkImportEmployees` | B2BService.yaml |
| B2BService | `/v1/b2b/plan-catalogs` | GET | `ListPurchaseOrderCatalog` (rename path) | B2BService.yaml |

---

### Implementation tasks for Phase 3

Order: Critical → Medium → High (High has breaking changes, deferred to v1.x sprint).

#### D-1  Fix routing conflicts [CRITICAL] [YAML]
- `B2BService.yaml`: remove `/v1/b2b/purchase-orders/catalog`; add `/v1/b2b/plan-catalogs`
- Update proto HTTP annotation for `ListPurchaseOrderCatalog` RPC
- `BeneficiaryService.yaml`: remove `/individual` and `/business`; update `CreateBeneficiary` proto to accept `type` field in body
- `TaskService.yaml`: remove `/tasks/my-tasks`; add `assignee` query param to `GET /v1/tasks`
- `WorkflowService.yaml`: remove `/workflow-tasks/my-tasks`; add `assignee` query param
- `SupportService.yaml`: remove duplicate `{slug}` path; add `?slug=` query param to `GetKnowledgeBaseArticle`
- `PaymentService.yaml`: move `webhook/{provider}` out of `/payments/`; new path `/v1/webhooks/payments/{provider}`
- `RefundService.yaml`: remove singular `/v1/policies/{policy_id}/refund`
- **Deliverable**: `buf lint` passes; `scripts/generate.ps1` regenerates clean OpenAPI

#### D-2  Add missing new-feature paths [YAML]
- Add all 9 paths from the "Missing paths" table above to their respective YAML files
- Follow existing path entry format (operationId, summary, parameters, requestBody, responses)
- **Deliverable**: `openapi.yaml` (regenerated) contains all new path entries

#### D-3  Fix MEDIUM inconsistencies [YAML]
- `AuthService.yaml`: rename `/auth/session/current` → `/auth/sessions/current`
- `AuthZService.yaml`: rename `/authz/check` → `/authz/access-decisions:check`; remove duplicate `get:` in portals config
- `KYCService.yaml`: remove `:pending` path; add `status` query param annotation to list endpoint
- `RenewalService.yaml`: remove `/renewals/upcoming`; add `status` query param
- `BillingService.yaml`: rename `/orders/{id}/invoice` → `/orders/{id}/invoices`
- `PartnerService.yaml`: remove `:update-status`; update `UpdatePartner` to handle status via PATCH body
- **Deliverable**: All MEDIUM violations resolved; zero duplicate routes

#### D-4  Auth path aliasing for v1 deprecation [YAML][PROTO]
- Add `additional_bindings` to `Register`, `Login`, `Logout`, `Email*` RPCs (see proto snippet above)
- Add `x-deprecated: true` + `x-sunset-version: "v2"` to old paths in YAML
- SDK generator must honour `x-deprecated` and emit `@Deprecated` / JSDoc `@deprecated`
- **Deliverable**: Old paths still work; OpenAPI renders deprecation notice on each

---

## Risk Register

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| `buf breaking` rejects new proto fields | Low | Medium | All changes additive (new field numbers only) |
| FLVE voice provider API down in CI | Medium | Low | Use mock in unit tests; integration test marked `short`-skip |
| Cross-schema catalog query breaks when schemas split | Medium | High | `domain.CatalogRepository` interface + `B2B_CATALOG_SOURCE` flag already in plan |
| Token rotation applied while in-flight JWTs use old key | Medium | High | Old key stays in DB with `is_active=false`; JWKS endpoint serves all active+recently-rotated keys |
| PO status race condition (concurrent approve + reject) | Low | Medium | Use `UPDATE … WHERE status = SUBMITTED` with rows-affected check |

---

## Completion Criteria

The feature is DONE when:

- [ ] All 10 implementation tasks (A-1 → A-4, B-1 → B-4, C-1 → C-10) are green in CI
- [ ] `buf lint` and `buf breaking` pass on main
- [ ] `go test ./...` -short passes in all three services
- [ ] `run_migration.ps1 --dry-run --strict` passes with no zombie columns
- [ ] `GET /v1/auth/me` returns roles for a logged-in user
- [ ] `GET /.well-known/jwks.json` returns the active key KID from `token_configs`
- [ ] POST approve + fulfill flow transitions a PO from SUBMITTED → APPROVED → FULFILLED
- [ ] Kafka topics `insuretech.b2b.v1.PurchaseOrderApproved` and `insuretech.b2b.v1.PurchaseOrderFulfilled` receive events in integration test
- [ ] B2B `/metrics` endpoint returns gRPC RED metrics
- [ ] `seededCatalogPlans` no longer exists in `b2b_service.go`
