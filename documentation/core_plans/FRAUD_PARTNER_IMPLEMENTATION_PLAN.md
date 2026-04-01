# Fraud & Partner Microservices - Revised Implementation Plan (Proto-First)

Generated: February 27, 2026  
Scope Inputs Reviewed:
- `documentation/SRS_v3/SPECS_V3.7/sections/*` (especially 4.9, 4.10, 4.16, 5, 6, 7, 8)
- `proto/insuretech/fraud/*` and `proto/insuretech/partner/*`
- Existing service implementations (`backend/inscore/microservices/authn`, `authz`, `partner`, `fraud`)
- Runtime and infra wiring (`backend/inscore/configs/services.yaml`, gateway wiring, db migrations)

---

## 1. Current State Summary

### Fraud service (implemented baseline, not greenfield)
- `backend/inscore/microservices/fraud` is implemented with `internal/config`, `internal/repository`, `internal/service`, `internal/grpc`, and `internal/events` packages plus live tests.
- `backend/inscore/microservices/fraud/cmd/server/main.go` registers `FraudService`, health, reflection, and an AuthZ unary interceptor.
- Fraud startup is already aligned to `backend/inscore/configs/services.yaml`; `FRAUD_PORT` / `FRAUD_GRPC_PORT` env overrides are explicitly ignored with a warning.
- Runtime methods exist for `CheckFraud`, alert/case CRUD, rule CRUD, activation/deactivation, publisher hooks, and consumer-triggered checks.
- Remaining fraud work is hardening: richer rule semantics, broader async consumers, stronger metrics packaging, and full SRS coverage.

### Partner service (implemented baseline, needs hardening)
- `backend/inscore/microservices/partner` already has domain interfaces, service modules, gRPC handlers/server, repositories, config loading, and live tests.
- The previously flagged RPCs now have concrete service implementations:
  - `UpdatePartner`
  - `DeletePartner`
  - `GetPartnerCommission`
  - `UpdateCommissionStructure`
  - `GetPartnerAPICredentials`
  - `RotatePartnerAPIKey`
- Partner gRPC server already uses the shared AuthZ interceptor pattern and loads ports from `services.yaml`, ignoring legacy env overrides.
- Remaining partner work is around deeper KYB/verification workflow rigor, policy-lifecycle commission automation, broader metrics, and stronger contract/E2E testing.

### Gateway + routing status
- Gateway wiring already includes both `partner` and `fraud` service blocks.
- Partner HTTP routes exist for CRUD, verification, commission, and credential rotation.
- Fraud HTTP routes exist for checks, alerts, cases, and rule lifecycle operations.
- CSRF is already applied to state-changing partner/fraud routes.

### Data/migrations
- Proto entities for fraud and partner are present and already reflected in runtime code paths.
- Enhancement migrations for partner/fraud tables exist and should still be treated as proto-first rollout artifacts.
- The remaining gap is not schema absence; it is keeping migration validation, backup verification, and runtime behavior aligned on every contract change.

---

## 2. Critical Gaps (Priority Ordered)

1. The largest gap is no longer service existence; it is incomplete business-rule coverage versus SRS requirements for fraud and partner flows.
2. Fraud runtime exists, but advanced rule depth, investigative workflows, and broader event consumption still need completion and validation.
3. Partner runtime exists, but commission automation from policy lifecycle events and stricter approval/audit semantics remain incomplete.
4. AuthN/AuthZ-style structure is mostly present now, but metrics/consumer/test maturity is still uneven compared with the stronger auth services.
5. Proto-first discipline still needs to be enforced continuously so schema changes, generated code, and rollout steps do not drift.
6. Gateway integration is already present; remaining work is end-to-end authorization, portal behavior validation, and regression testing.
7. Compliance, audit, and security mappings from SRS section 7 still need explicit verification against service behavior and emitted events.

---

## 3. Non-Negotiable Implementation Rules

1. Proto is the single source of truth.
2. No duplicated DTO/model types; CRUD uses generated proto structs with injected gorm tags.
3. Migration flow is always:
   - proto update/freeze
   - codegen
   - db migrate
   - schema check
   - only then service CRUD implementation
4. Follow AuthN/AuthZ structure:
   - `cmd/server`
   - `internal/config`
   - `internal/domain`
   - `internal/repository`
   - `internal/service`
   - `internal/grpc` (handler/interceptors)
   - `internal/events` (publisher/consumer)
   - `internal/metrics`
5. Service ports come from `backend/inscore/configs/services.yaml`; env port overrides are warnings-only or ignored.
6. Gateway never performs service-level CRUD logic.

---

## 4. Service Structures On Current Branch

## 4.1 Fraud (`backend/inscore/microservices/fraud`)
- `cmd/server/main.go`
- `internal/config/config.go`
- `internal/domain/*`
- `internal/repository/fraud_rule_repository.go`
- `internal/repository/fraud_alert_repository.go`
- `internal/repository/fraud_case_repository.go`
- `internal/service/fraud_engine_service.go`
- `internal/service/fraud_case_service.go`
- `internal/grpc/fraud_handler.go`
- `internal/events/publisher.go`
- `internal/events/consumer.go`
- tests: repository/service live tests + publisher/consumer tests

Current nuance:
- most of this structure already exists in runtime code
- metrics are currently lighter-weight (`MetricsSnapshot()` in service) than a dedicated `internal/metrics` package
- this plan section should now be read as hardening and normalization guidance, not as a missing-service blueprint

## 4.2 Partner (`backend/inscore/microservices/partner`)
- Existing structure already includes `internal/domain/interfaces.go`, repositories, service split, gRPC server/handler, and live tests.
- Keep refactoring toward authz-style boundaries where helpful, but do not treat partner as a stubbed service anymore.
- Add consumer and metrics modules where policy-lifecycle automation requires them.
- Expand tests from current live-service coverage to broader unit + gRPC integration + eventing coverage.

---

## 5. SRS -> Implementation Mapping

### Partner
- Core M1/M2/M3 requirements:
  - FR-059..069, FR-070..072, FR-141, FR-145
- Deferred/phase-gated:
  - FR-142, FR-144, FR-146, FR-206 (D)

### Fraud
- Core M1/M2/M3 requirements:
  - FR-186..192
  - FR-054/055 integration behavior
- Deferred:
  - FR-166 (advanced AI model maturity)
  - FR-234 graph-db visualization (Neo4j/Neptune, D)

---

## 6. Proto-First Delivery Sequence

### Phase 0 - Contract Freeze
1. Validate `partner_service.proto`, `partner.proto`, `partner_events.proto`.
2. Validate `fraud_service.proto`, `fraud_rule.proto`, `fraud_alert.proto`, `fraud_case.proto`, `fraud_events.proto`.
3. Run proto lint/breaking checks and regenerate code.

### Phase 1 - Schema Readiness
1. Apply migrations from proto pipeline + enhancement migrations.
2. Verify table/constraint/index existence for:
   - `partner_schema.partners`, `agents`, `commissions`
   - `insurance_schema.fraud_rules`, `fraud_alerts`, `fraud_cases`
3. Run schema consistency checks on primary/backup.

### Phase 2 - Repository Completeness
1. Partner repos:
   - Add update/delete/list-filter methods required by all RPCs.
   - Support page token/filter/order semantics.
   - Guarantee PII encryption/decryption path coverage.
2. Fraud repos:
   - Full CRUD/list for rules/alerts/cases.
   - JSONB-safe read/write for rule conditions, alert details, case evidence.

### Phase 3 - Service Logic Completeness
1. Partner:
   - Implement real update/delete.
   - Implement commission aggregation and structure update logic.
   - Implement API credential retrieval/rotation via authn integration client.
   - Implement focal-person verification and state transitions with audit semantics.
2. Fraud:
   - Implement rule evaluation engine for FD-001..FD-007.
   - Implement score aggregation/risk-level mapping.
   - Implement alert creation and case lifecycle transitions.
   - Implement confirmed-fraud action workflow for downstream account restrictions.

### Phase 4 - Eventing
1. Partner publisher:
   - `PartnerOnboarded`, `PartnerVerified`, `AgentRegistered`, `CommissionCalculated`.
2. Partner consumer:
   - consume `PolicyIssued`, `PolicyRenewed` for commission generation.
3. Fraud publisher:
   - `FraudAlertTriggered`, `FraudCaseCreated`, `FraudConfirmed`.
4. Fraud consumer:
   - consume claim/policy/customer risk events for async fraud analysis.

### Phase 5 - Transport + Gateway
1. Build grpc handlers (validation, error-to-status mapping) for partner/fraud.
2. Ensure auth interceptors align with AuthN/AuthZ style.
3. Add gateway service registration for fraud.
4. Add gateway handlers/routes for partner and fraud REST endpoints mapped from proto HTTP annotations.

### Phase 6 - Observability + Compliance
1. Structured logging via pkg logger + zap in all layers.
2. Prometheus metrics:
   - fraud check latency, hit rate by rule, alert volume
   - partner onboarding SLA, commission calc latency/failures
3. Audit/security requirements mapping:
   - partner approval trails (FR-066/070/072)
   - fraud incident and escalation traces (FR-191/192, SEC/BFIU alignment)

### Phase 7 - Test & Release Gates
1. Unit tests for service logic and rule engines.
2. Live repository tests on cloud postgres.
3. gRPC integration tests for all RPCs.
4. Kafka integration tests for produced/consumed events.
5. Gateway E2E tests for REST paths.
6. Performance checks vs NFR targets for p95 latency.

---

## 7. RPC Completion Matrix (Current vs Target)

## 7.1 PartnerService RPCs
- `CreatePartner`: implemented baseline -> continue hardening onboarding validation and transactional semantics.
- `GetPartner`: implemented baseline -> continue tightening authz/error mapping consistency.
- `UpdatePartner`: implemented persistence path -> verify update-mask/field-coverage semantics against proto.
- `ListPartners`: implemented baseline pagination/filtering -> verify ordering/page-token/total-count behavior against contract.
- `DeletePartner`: implemented soft-delete path -> ensure audit/event semantics are complete.
- `VerifyPartner`: implemented baseline -> deepen KYB and focal-person approval workflow.
- `UpdatePartnerStatus`: implemented baseline -> tighten state machine and reason logging rules.
- `GetPartnerCommission`: implemented baseline aggregate response -> extend as needed for richer reconciliation/reporting.
- `UpdateCommissionStructure`: implemented baseline -> harden validation and approval workflow.
- `GetPartnerAPICredentials`: implemented via integration/authn-backed path -> verify rotation/audit guarantees.
- `RotatePartnerAPIKey`: implemented via integration/authn-backed path -> verify lifecycle and traceability.

## 7.2 FraudService RPCs
- The runtime already implements:
  - `CheckFraud`
  - `GetFraudAlert`, `ListFraudAlerts`
  - `CreateFraudCase`, `GetFraudCase`, `UpdateFraudCase`
  - `ListFraudRules`, `CreateFraudRule`, `UpdateFraudRule`
  - `ActivateFraudRule`, `DeactivateFraudRule`
- Remaining target work is not transport coverage; it is deeper rule sophistication, broader SRS mapping, stronger eventing, and stronger end-to-end validation.

---

## 8. Concrete Refactor Actions Required Immediately

1. Preserve the existing `services.yaml`-driven startup behavior for fraud/partner; do not regress to env-driven port selection.
2. Add stronger metrics packages and dashboards so fraud/partner observability reaches authn/authz maturity.
3. Expand partner event consumption for policy-issued / policy-renewed commission automation.
4. Expand fraud async intake beyond the current baseline so claim/policy/customer risk signals are covered consistently.
5. Add contract, gateway, and portal E2E tests around the already-wired partner/fraud routes.

---

## 9. Definition of Done

1. All partner/fraud proto RPCs implemented and passing integration tests.
2. Fraud and partner service startup follow centralized config path + service port conventions.
3. No duplicate model/types outside generated proto entities.
4. Event publish/consume flows verified with Kafka integration tests.
5. Gateway exposes and secures partner/fraud endpoints.
6. SRS M1/M2/M3 fraud + partner requirements mapped to tested behavior.

---

## 10. Execution Order Recommendation

1. Fraud hardening first: deepen rule coverage, event intake, and runtime verification against proto/SRS.
2. Partner hardening second: tighten commission automation, verification workflow, and audit guarantees.
3. Gateway and cross-service E2E validation third, because the transport wiring already exists.
4. Compliance/analytics/deferred D/F items remain a separate roadmap after the current runtime is behaviorally complete.
