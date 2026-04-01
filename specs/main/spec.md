# Feature Specification: AuthN / AuthZ / B2B Completion

**Feature ID**: main  
**Branch**: main  
**Status**: In Progress  
**Author**: speckit.specify  
**Version**: 1.0.0  
**Created**: 2026-03-12  

---

## Overview

Complete the three core identity and B2B microservices of the LabAid InsureTech platform so
that they satisfy all production acceptance criteria. The three services share a dependency
chain: AuthN (**identity**) is consumed by AuthZ (**authorisation**), which is in turn
consumed by B2B (**B2B org management**). All three must be production-ready before any
downstream insurance workflows can go live.

### Services in scope

| Service | gRPC port | Schema | Status |
|---------|-----------|--------|--------|
| `authn` | `:50053` | `authn_schema` | ~90 % complete |
| `authz` | `:50054` | `authz_schema` | ~95 % complete |
| `b2b`   | `:50055` | `b2b_schema`   | ~70 % complete |

---

## Goals

1. **AuthN** — Close all unimplemented RPCs; harden voice-biometric, KYC, and FLVE
   integration paths; add `GetMe` end-point (unified user+profile+roles response).
2. **AuthZ** — Implement DB-backed `GetJWKS` (read active key from `token_configs` instead
   of file), promote `acl/` placeholder to working ACL sub-module, add token-config rotation
   RPC.
3. **B2B** — Implement full purchase-order lifecycle (Approve / Reject / Fulfill RPCs +
   events), wire `UpdateDepartmentTotalPremium`, replace hardcoded catalog fallback with
   proper DB seed data, implement employee bulk-import, and wire all stub event consumers.

---

## Non-Goals

- New insurance-domain features (policy, claims, payment) — out of scope.
- Portal UI for these features — tracked separately in B2B portal workstream.
- SDK regeneration — handled by the `scripts/generate.ps1` pipeline after proto changes.

---

## Requirements

### REQ-AUTHN-01 — GetMe unified endpoint

> **Priority**: HIGH  
> A single gRPC RPC `GetMe` returns the calling user's identity (from AuthN) merged with
> their active roles and tenant assignments (from AuthZ). Required for portal onboarding.

**Acceptance criteria:**
- [ ] Proto message `GetMeRequest` / `GetMeResponse` added to `auth_service.proto`.
- [ ] Response includes: `user_id`, `user_type`, `email`, `phone`, `profile`, `portal`, `tenant_id`, `roles[]`, `permissions[]`.
- [ ] AuthZ gRPC client is injected into AuthN service for role/permission fetch.
- [ ] Validated with unit test (mock AuthZ client) and integration test.

### REQ-AUTHN-02 — Voice biometric integration depth

> **Priority**: MEDIUM  
> `InitiateVoiceSession`, `SubmitVoiceSample`, `VerifyVoiceSession` are in proto and
> service_iface.go. The external voice biometric provider (FLVE) integration must be
> verified complete and non-stub.

**Acceptance criteria:**
- [ ] `voice_service.go` calls the external FLVE client (no no-op stubs).
- [ ] FLVE client errors map to gRPC status codes (not wrapped `errors.New`).
- [ ] Unit tests cover FLVE error paths.
- [ ] Voice session state machine: INITIATED → SAMPLES_COLLECTED → VERIFIED / FAILED.

### REQ-AUTHN-03 — KYC orchestrator integration

> **Priority**: MEDIUM  
> `InitiateKYC`, `SubmitKYCFrame`, `CompleteKYCSession` are implemented; the KYC
> microservice call chain must be verified end-to-end, and `ApproveKYC`/`RejectKYC`
> must publish events to Kafka.

**Acceptance criteria:**
- [ ] `kyc_orchestrator_service.go` calls the KYC microservice via gRPC (no stub).  
- [ ] `ApproveKYC` publishes `KYCApprovedEvent` to `insuretech.authn.v1.KYCApproved`.
- [ ] `RejectKYC` publishes `KYCRejectedEvent` to `insuretech.authn.v1.KYCRejected`.
- [ ] Events conform to outbox pattern (write event row before Kafka produce step).

### REQ-AUTHZ-01 — Token config DB-backed JWKS

> **Priority**: HIGH  
> `GetJWKS` currently reads the RSA public key from the file path
> `AUTHZ_JWKS_PUBLIC_KEY_PATH`. It should prefer the active record from
> `authz_schema.token_configs` so that key rotation is reflected without service restart.

**Acceptance criteria:**
- [ ] `GetJWKS` queries `TokenConfigRepository.GetActive(ctx)` first; falls back to env-var
  file only when table is empty (bootstrap scenario).
- [ ] In-memory cache (TTL = 5 min) protects the DB lookup on hot paths.
- [ ] Unit test: table has active record → served from DB.
- [ ] Unit test: table empty → served from file-path fallback.
- [ ] Integration test: rotate key in DB → next JWKS response reflects new KID.

### REQ-AUTHZ-02 — Token rotation RPC

> **Priority**: MEDIUM  
> Add `RotateTokenConfig` RPC to proto and service to allow online key rotation without
> downtime. New key is added as active; old key is marked inactive (not deleted, for
> in-flight JWT validation).

**Acceptance criteria:**
- [ ] Proto: `RotateTokenConfig(RotateTokenConfigRequest) returns (RotateTokenConfigResponse)`.
- [ ] Old active key's `is_active` set to `false`; new key inserted as active.
- [ ] `InvalidatePolicyCache` called after rotation (cache depends on JWKS keys).
- [ ] Old key remains in table with `rotated_at` timestamp (for in-flight token validation).
- [ ] Audit log entry created.
- [ ] Unit + integration tests.

### REQ-AUTHZ-03 — ACL sub-module implementation

> **Priority**: LOW  
> `internal/acl/` is an empty placeholder. It should contain the resource-action ACL
> matrix used by `CheckAccess` for human-readable documentation and admin UI tooling.

**Acceptance criteria:**
- [ ] `acl/acl.go` defines `ResourceActionMatrix` map keyed by `portal → resource → []allowed_actions`.
- [ ] `CheckAccess` delegates to this matrix for unknown `object` patterns (not just Casbin).
- [ ] Seeder uses this matrix to seed `p`-type Casbin rules on fresh installs.

### REQ-B2B-01 — Purchase order approval lifecycle

> **Priority**: CRITICAL  
> `CreatePurchaseOrder` sets status `SUBMITTED`. The SUBMITTED → APPROVED → FULFILLED →
> REJECTED lifecycle transitions have no RPC, no event, and no service logic.

**Acceptance criteria:**
- [ ] Proto: add `ApprovePurchaseOrder`, `RejectPurchaseOrder`, `FulfillPurchaseOrder` RPCs.
- [ ] `ApprovePurchaseOrder`: transitions `SUBMITTED → APPROVED`; only callable by
  `B2B_ORG_ADMIN` or `SYSTEM_USER`.
- [ ] `RejectPurchaseOrder`: transitions `SUBMITTED → REJECTED` with required `reason` field.
- [ ] `FulfillPurchaseOrder`: transitions `APPROVED → FULFILLED`; triggered after payment
  confirmation (can be called externally or via event consumer).
- [ ] `PurchaseOrderApprovedEvent` published to `insuretech.b2b.v1.PurchaseOrderApproved`.
- [ ] `PurchaseOrderFulfilledEvent` published to `insuretech.b2b.v1.PurchaseOrderFulfilled`.
- [ ] `PurchaseOrderRejectedEvent` published to `insuretech.b2b.v1.PurchaseOrderRejected`.
- [ ] UpdateDepartmentTotalPremium called when PO is Fulfilled.
- [ ] Unit tests for each transition; invalid transitions return `FailedPrecondition`.

### REQ-B2B-02 — UpdateDepartmentTotalPremium service exposure

> **Priority**: HIGH  
> `UpdateDepartmentTotalPremium` exists in the repository interface but has no service
> method or handler RPC. It is called internally after purchase order fulfillment.

**Acceptance criteria:**
- [ ] Called internally from `FulfillPurchaseOrder` (not a public RPC — internal only).
- [ ] Recalculates total premium as SUM(employee_count × plan_premium) for all FULFILLED
  purchase orders in the department.
- [ ] Unit test verifies total is updated correctly.

### REQ-B2B-03 — Catalog seed data cleanup

> **Priority**: HIGH  
> `seededCatalogPlans` in `b2b_service.go` contains hardcoded fixed UUIDs matching no
> real product records. This will corrupt production catalog if `insurance_schema.product_plans`
> is empty.

**Acceptance criteria:**
- [ ] Move the 3 seeded plans to `backend/inscore/db/seeds/b2b_catalog_seed.sql` as
  `INSERT INTO insurance_schema.product_plans … ON CONFLICT DO NOTHING`.
- [ ] Remove `seededCatalogPlans` map and `fallbackCatalogPlan` function from service.
- [ ] `mergeCatalogWithSeedFallback` removed; `ListCatalogPlans` returns DB result directly.
- [ ] Seed file applied as part of `run_migration.ps1` seed step.
- [ ] Integration test: empty product_plans → returns empty list (no silent fallback).

### REQ-B2B-04 — Employee bulk import

> **Priority**: MEDIUM  
> ACTIVE_WORKSTREAMS.md references employee bulk import. It is absent from proto, handler,
> and service.

**Acceptance criteria:**
- [ ] Proto: `BulkImportEmployees(BulkImportEmployeesRequest) returns (BulkImportEmployeesResponse)`.
- [ ] Request accepts `repeated EmployeeImportRow` (max 500 rows per call).
- [ ] Partial failure mode: each row has a per-row `success/error` in response.
- [ ] Duplicate NID/mobile detected and returned as error row (not hard fail).
- [ ] Unit test for partial failure scenario.

### REQ-B2B-05 — Wire stub event consumers

> **Priority**: HIGH  
> `HandleOrganisationApproved`, `HandleUserRegistered`, `HandleRoleAssigned` in
> `consumers/handlers.go` are no-ops (log + return nil). They must perform real work.

**Acceptance criteria:**
- [ ] `HandleOrganisationApproved`: update `organisations.status = ACTIVE` for the org.
- [ ] `HandleUserRegistered`: check `pending_org_invitations` table; if match found, auto-add
  user as org member with `partner_user` role.
- [ ] `HandleRoleAssigned`: emit notification event to notification service (informational).
- [ ] Unit tests for each consumer with mocked repo/client.

### REQ-B2B-06 — Employee ↔ B2C user account linking

> **Priority**: MEDIUM  
> Employees have a `user_id` field in proto but there is no mechanism to link an Employee
> record to an AuthN B2C_CUSTOMER user account after enrollment.

**Acceptance criteria:**
- [ ] `CreateEmployee` or a dedicated `LinkEmployeeUser` RPC accepts `user_id` and writes
  it to the employee record.
- [ ] `HandleUserRegistered` consumer links if phone/email matches a pending employee record.
- [ ] Integration test: create employee → register user with matching mobile → employee.user_id populated.

---

## Out-of-Scope / Deferred

| Item | Reason |
|------|--------|
| Multi-key concurrent JWKS validation | Complex — deferred to post v1.0 |
| Webhook delivery for PO events | Notification service owns this |
| Portal SSO / SAML federation | Not in sprint |
| Regulatory reporting from B2B data | Belongs to report pipeline |

---

## Dependencies

| Service | Protocol | Used By |
|---------|----------|---------|
| AuthZ gRPC (`:50054`) | Proto | AuthN `GetMe`, B2B consumers |
| AuthN gRPC (`:50053`) | Proto | B2B employee enrollment |
| KYC microservice | gRPC | AuthN `kyc_orchestrator_service.go` |
| FLVE (external) | REST/gRPC | AuthN `voice_service.go` |
| Notification service | Kafka | B2B `HandleRoleAssigned` consumer |
| Payment service | Kafka | B2B `FulfillPurchaseOrder` consumer |
