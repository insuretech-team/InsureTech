# Research: AuthN / AuthZ / B2B Completion

**Phase**: 0 — Research  
**Branch**: main  
**Generated**: 2026-03-12  
**Feeds**: plan.md Phase 1 design

---

## Research Objective

Resolve all NEEDS-CLARIFICATION items from the Technical Context and document technology
decisions before Phase 1 design begins. Each finding is expressed as Decision / Rationale /
Alternatives.

---

## R-01  AuthN GetMe — where does AuthZ role-fetch live?

**Question**: Should AuthN call AuthZ internally to enrich the `GetMe` response, or should
the gateway aggregate the two calls client-side?

**Findings**:
- AuthN already holds the session, user type, and portal context. Clients currently make
  two calls (AuthN `ValidateToken` + AuthZ `GetUserPermissions`) to build a dashboard.
- The gateway (`backend/inscore/...`) routes gRPC; adding aggregation there violates
  Principle II (gateway has no business logic).
- Both AuthN and AuthZ are Go services — an internal gRPC loopback is safe.

**Decision**: AuthN `GetMe` calls AuthZ `GetUserPermissions` over loopback gRPC. AuthN
service receives an `authzClient AuthZServiceClient` dependency via constructor injection.

**Rationale**: Single RTT from client; client SDK surface shrinks (one call not two);
AuthN already owns the token context so it is the natural aggregation point.

**Alternatives considered**:
- Gateway aggregation: rejected — gateway is a proxy, not a BFF.
- Client-side aggregation: rejected — doubles mobile app RTT on every page load.
- New BFF microservice: rejected — YAGNI (Principle X).

---

## R-02  AuthZ JWKS — DB table vs file path strategy

**Question**: `GetJWKS` reads from `AUTHZ_JWKS_PUBLIC_KEY_PATH` env var. The
`token_configs` DB table is seeded but never read. What is the intended design?

**Findings**:
- The seeder (`portal_seeder.go`) writes the public key PEM and KID into
  `authz_schema.token_configs` during service startup.
- `GetJWKS` ignores that table and reads the same PEM from file.
- This means the DB table is write-only — a waste and a future footgun (key rotation
  writes to DB but JWKS response doesn't reflect it).
- RS256 JWKS endpoints are typically served from a key store, not a file. The pattern
  "active row in token_configs" already exists — it just needs to be read.

**Decision**: `GetJWKS` prefers `token_configs.GetActive()`. Falls back to file only when
table is empty (first-boot before seeder completes).

**Rationale**: Enables zero-downtime key rotation: rotate key in `token_configs` →
`GetJWKS` reflects it within cache TTL (5 min) without service restart.

**Alternatives considered**:
- Keep file-only: rejected — cannot rotate without restart; violates operational readiness
  (Principle VII).
- External secret manager only (Vault/KMS): deferred — overkill for v1.0; `token_configs`
  is sufficient, and `private_key_ref` already holds the Vault path for future use.

---

## R-03  AuthZ token key rotation — RPC scope

**Question**: Should key rotation be a new proto RPC or a seeder-only operation?

**Findings**:
- Kubernetes restarts can trigger the seeder at any time; seeder idempotently writes a
  key only when `GetActive` returns empty.
- For rolling key rotation (e.g. 90-day PCI cycle), an operator RPC is needed to:
  1. Generate new RSA key pair.
  2. Mark old key `is_active = false`.
  3. Insert new key as active.
  4. Invalidate JWKS cache.
- The AuthN private-key file is separate; AuthZ serves only the verification public key.

**Decision**: Add `RotateTokenConfig` RPC to `authz_service.proto`. Private key generation
happens outside the service (KMS/Vault); the RPC accepts a new public key PEM + private
key Vault reference.

**Rationale**: Keeps private key material out of the RPC wire (constitution Principle V).
Only the public key PEM is written to `token_configs`; private key ref is a Vault path string.

**Alternatives considered**:
- CLI script (not an RPC): rejected — can't audit trail; no gRPC authz check.
- Auto-rotation on TTL: deferred — requires cronjob + Vault integration (v2 feature).

---

## R-04  B2B purchase order approval — who approves?

**Question**: The PO lifecycle (SUBMITTED → APPROVED → FULFILLED → REJECTED) has no RPCs.
Who calls `ApprovePurchaseOrder`, and how is the fulfillment triggered?

**Findings**:
- Current PO status enum: DRAFT, SUBMITTED, APPROVED, FULFILLED, REJECTED.
- `CreatePurchaseOrder` writes `SUBMITTED` immediately (skips DRAFT for now).
- In the Bangladeshi group insurance model, a B2B org admin submits; an insurer or
  system admin approves.
- `FulfillPurchaseOrder` maps to "policy batch issuance" — triggered when payment intent
  is confirmed by the payment service.

**Decision**:
- `ApprovePurchaseOrder`: callable by `SYSTEM_USER` or `B2B_ORG_ADMIN` with `approver_notes`.
- `RejectPurchaseOrder`: callable by same roles; requires `reason`.
- `FulfillPurchaseOrder`: internal-only RPC (no HTTP annotation), called by payment
  event consumer when `PaymentCompleted` event is received for the PO.

**Rationale**: Separating approve from fulfill keeps the two concerns (business approval
vs. financial settlement) clean. The payment consumer triggers fulfillment so the
insurance-domain services (PoliSync) can start policy issuance.

**Alternatives considered**:
- PATCH `/v1/b2b/purchase-orders/{id}` with `{"status":"APPROVED"}`: rejected — status
  transitions are domain verbs; Principle III forbids action verbs in paths.
- Auto-approve on create: rejected — violates Bangladeshi insurance compliance.

---

## R-05  B2B catalog — cross-schema query vs gRPC call

**Question**: `catalog_repository.go` queries `insurance_schema.product_plans` directly
across DB schemas. Is this a constitution violation?

**Findings**:
- Principle II states "Cross-runtime logic sharing happens only through proto contracts
  and Kafka domain events — never through shared libraries or **direct DB cross-access**."
- Go services (not C#) own `insurance_schema`; cross-Go-service DB reads are more
  acceptable than cross-runtime reads but still couple services to each other's schema.
- In the current monorepo Docker Compose, all Go services share one PostgreSQL instance.
- Principle XI confirms that B2B may not read from `insurance_schema` without a contract.
- The proper contract is a gRPC call to the products/plans service (Go `insurance-service`
  or PoliSync bridge).

**Decision for v1.0**: Keep the cross-schema read **but** gate it behind a
`domain.CatalogRepository` interface with a swappable gRPC implementation.
Config flag `B2B_CATALOG_SOURCE=db|grpc` controls which is used.
Remove hardcoded `seededCatalogPlans`; replace with proper DB seed data.

**Decision for v2.0**: Switch `B2B_CATALOG_SOURCE=grpc` pointing to products service
once that service stabilises its RPC contract.

**Rationale**: Avoids a blocking dependency on the products service for B2B v1.0.
Interface boundary is set now so the switch is a config change, not a refactor.

**Alternatives considered**:
- Force gRPC immediately: rejected — products service RPC contract is not stable.
- Keep hardcoded fallback: rejected — breaks production when DB is populated with real data.

---

## R-06  B2B employee bulk import — batch size and error handling

**Question**: What is a safe batch size and how should partial failures be reported?

**Findings**:
- Average B2B org in Bangladesh: 50–2000 employees.
- PostgreSQL can handle batch inserts of 500 rows trivially.
- Partial failure handling: per-row errors with row index in response (not fail-all).
- NID is the primary dedup key; mobile number is secondary.

**Decision**: Max 500 rows per request. Response includes `[]ImportResult{ row_index, 
success, error_code, error_detail }`. Mandatory fields: `full_name`, `mobile`,
`department_id`, `plan_id`. `nid` is optional (encouraged, not required).

**Alternatives considered**:
- Async job (upload CSV, poll status): deferred — overkill for ≤500 rows.
- Fail-all-or-nothing: rejected — HR use case requires per-row feedback.

---

## R-07  B2B stub consumers — pending_org_invitations table

**Question**: `HandleUserRegistered` should check pending org invitations but no such
table exists in `b2b_schema`.

**Findings**:
- No `pending_org_invitations` or `invitations` table exists in any migration file.
- The ACTIVE_WORKSTREAMS.md mentions "invite link flow" as planned.
- For v1.0, `HandleUserRegistered` cannot auto-link without an invitations table.

**Decision**: `HandleUserRegistered` logs and is a no-op for v1.0. The invitation table
and auto-link logic are deferred to the "B2B Invite Flow" feature.

**Alternatives considered**:
- Phone-match employee record: could work, but phone matching without invite confirmation
  would violate zero-trust (Principle V).

---

## R-08  Observability gaps

**Question**: Do all three services expose Prometheus RED metrics and OTel traces?

**Findings**:
- AuthN: `internal/metrics/` exists — has Prometheus metrics.
- AuthZ: `internal/metrics/` exists — has Prometheus metrics.
- B2B: no `internal/metrics/` directory found; no metrics registration visible.

**Decision**: Add `internal/metrics/` to B2B with counter/histogram for each gRPC method.
Wire via gRPC `StatsHandler` (same pattern as AuthN/AuthZ).

---

## Summary of Decisions

| ID | Decision | Impact area |
|----|----------|-------------|
| R-01 | `GetMe` aggregates in AuthN via AuthZ loopback gRPC | AuthN proto + service |
| R-02 | `GetJWKS` reads `token_configs` first, file as fallback | AuthZ service |
| R-03 | `RotateTokenConfig` RPC added to authz proto | AuthZ proto + service |
| R-04 | PO approval is multi-RPC; fulfill triggered by payment event | B2B proto + service + events |
| R-05 | Keep cross-schema query v1; add interface + config flag for v2 gRPC switch | B2B domain + catalog repo |
| R-06 | Bulk import: max 500 rows, per-row error response | B2B proto + service |
| R-07 | `HandleUserRegistered` is no-op v1.0; invitations deferred | B2B consumers |
| R-08 | Add Prometheus metrics to B2B service | B2B internal/metrics |
