# Implementation Baseline

This document explains the current repository in implementation terms, not aspirational terms.

## 1. Repo Topology

The codebase has four main layers.

### 1.1 Contracts

- `proto/` is the primary contract layer.
- It contains about `180` proto files across auth, authz, B2B, products, quotations, orders, policies, claims, payment, underwriting, renewal, refund, document, storage, partner, fraud, workflow, analytics, support, report, voice, and tenant domains.

### 1.2 Go platform and infrastructure

- Root: `backend/inscore/`
- Contains:
  - gateway
  - identity and authorization services
  - B2B service
  - partner and fraud services
  - document/media/storage services
  - payment, orders, insurance, workflow, notification, support, analytics, AI, and related services

Current implementation footprint:

- `29` entrypoints under `backend/inscore/cmd`
- `24` service folders under `backend/inscore/microservices`

### 1.3 C# domain host

- Root: `backend/polisync/`
- Contains `14` projects under `backend/polisync/src`
- Hosts business-domain gRPC services in a single .NET application

### 1.4 Web portals

- `b2b_portal/` - most integrated
- `customer_portal/` - mostly presentation/static data
- `system_portal/` - admin shell with mixed live/demo state

## 2. Runtime Architecture

## 2.1 Service registry

The authoritative runtime map is `backend/inscore/configs/services.yaml`.

The registry groups the platform into:

- Core infrastructure:
  - tenant
  - authn
  - authz
  - audit
- User and partner management:
  - kyc
  - partner
  - beneficiary
  - b2b
- Insurance data layer:
  - insurance
- Commerce and production:
  - product
  - quote
  - order
  - orders
  - commission
- Policy lifecycle:
  - policy
  - underwriting
  - workflow
- Financials:
  - payment
  - ledger
- Claims and fraud:
  - claim
  - fraud
- Communications and media:
  - notification
  - support
  - webrtc
  - media
  - ocr
- Documents and storage:
  - docgen
  - storage
- Intelligence:
  - iot
  - analytics
  - ai
- Edge/frontends:
  - gateway
  - b2b_portal
  - redis

This is important because several old plan docs described services as missing even though the runtime registry already includes them.

## 2.2 Gateway

The gateway entrypoint is `backend/inscore/cmd/gateway/main.go`.

Current facts:

- default HTTP port is `8080`
- exposes `healthz` and `readyz`
- uses resilient downstream client management
- validates auth with AuthN and authorization with AuthZ
- fronts both Go services and PoliSync services

The router in `backend/inscore/cmd/gateway/internal/routes/router.go` currently contains `202` registrations.

### Gateway domain coverage already present

- Auth and profile management
- B2B organizations, members, departments, employees, bulk upload, purchase orders
- Media APIs
- Document and template APIs
- Storage APIs
- Partner APIs
- Fraud APIs
- Product APIs
- Quotation APIs
- Order APIs
- Policy APIs
- Underwriting APIs
- Claim APIs
- Commission APIs

That means the HTTP surface is already broad. The main challenge is not absence of routing, but depth/completeness by domain and portal adoption.

## 2.3 Go services that are clearly real now

The following areas have active code, entrypoints, and service logic, not just placeholders:

- AuthN
  - `backend/inscore/cmd/authn`
  - `backend/inscore/microservices/authn`
- AuthZ
  - `backend/inscore/cmd/authz`
  - `backend/inscore/microservices/authz`
- B2B
  - `backend/inscore/cmd/b2b`
  - `backend/inscore/microservices/b2b`
- Partner
  - `backend/inscore/cmd/partner`
  - `backend/inscore/microservices/partner`
- Fraud
  - `backend/inscore/cmd/fraud`
  - `backend/inscore/microservices/fraud`
- Payment
  - `backend/inscore/cmd/payment`
  - `backend/inscore/microservices/payment`
- Orders data/service layer
  - `backend/inscore/cmd/orders`
  - `backend/inscore/microservices/orders`
- Document generation
  - `backend/inscore/cmd/docgen`
  - `backend/inscore/microservices/docgen`
- Media processing
  - `backend/inscore/cmd/media`
  - `backend/inscore/microservices/media`
- Storage
  - `backend/inscore/cmd/storage`
  - `backend/inscore/microservices/storage`
- Workflow
  - `backend/inscore/cmd/workflow`
  - `backend/inscore/microservices/workflow`

## 2.4 PoliSync is active infrastructure, not theory

`backend/polisync/src/PoliSync.ApiHost/Program.cs` already maps:

- `ProductGrpcService`
- `QuotesGrpcService`
- `OrderGrpcService`
- `PolicyGrpcService`
- `ClaimGrpcService`
- `CommissionGrpcService`
- `UnderwritingGrpcService`
- `EndorsementGrpcService`
- `RenewalGrpcService`
- `RefundGrpcService`

PoliSync also wires:

- PostgreSQL
- Redis
- Kafka
- JWT authentication
- health checks
- gRPC reflection
- Go data gateways and gRPC clients

This matters because earlier docs often treated product/order/policy/claim work as mostly planned. The infrastructure is already there; the question is domain completeness and contract alignment.

## 3. Domain Reality By Area

## 3.1 AuthN and AuthZ

Current reality:

- both services exist and run
- gateway routes are already extensive
- JWT, sessions, profile, document, KYC, API key, TOTP, and voice-related flows are present in the gateway surface
- B2B and downstream services rely on them synchronously

Interpretation:

Auth is a live platform dependency, not a stubbed subsystem.

## 3.2 B2B

Current reality:

- B2B service exists with service entrypoint, repository, service, middleware, event publisher, consumer startup, and gRPC handler layers
- gateway routes already expose:
  - organizations
  - org members
  - admin assignment
  - departments
  - employees
  - employee bulk upload
  - purchase-order catalog
  - purchase-order list/get/create
- purchase-order update/delete are still explicitly `501` at the gateway because the proto does not define those RPCs yet

Interpretation:

B2B is one of the most concrete end-to-end areas in the repository.

## 3.3 Payment and Orders

> **Status updated 2026-03-06** — Phase 2 payment and billing work complete.

**What is real and working (Go — inscore):**

- `backend/inscore/microservices/payment/` — full Go payment service
  - gRPC port 50190. Pure gRPC (no HTTP server of its own).
  - `PaymentService` has 15 RPCs, all implemented: `InitiatePayment`, `VerifyPayment`, `GetPayment`, `ListPayments`, `InitiateRefund`, `GetRefundStatus`, `ListPaymentMethods`, `AddPaymentMethod`, `ReconcilePayments`, `HandleGatewayWebhook`, `GetPaymentByProviderReference`, `SubmitManualPaymentProof`, `ReviewManualPayment`, `GenerateReceipt`, `GetPaymentReceipt`
  - `PaymentRepository` has 10 methods including `GetPaymentByTranID`, `GetPaymentByOrderID`, `GetPaymentByProviderReference`
  - SSLCommerz provider client at `internal/providers/sslcommerz/client.go` — `InitSession`, `ValidatePayment`, `QueryPayment`, `InitiateRefund`
  - All 9 Kafka event topics canonical: `insuretech.payment.v1.payment.initiated`, `.completed`, `.failed`, `.verified`, `.manual_review_requested`, `.manual_review_completed`, `.receipt_generated`, `.reconciliation_mismatch`, and refund topic
  - Payment entity has 53 proto fields covering all SSLCommerz correlation, manual review, and receipt fields
  - Gateway routes wired for all 15 RPCs + 4 SSLCommerz HTTP callback routes (IPN/success/fail/cancel)

- `backend/inscore/microservices/orders/` — Go orders data layer
  - gRPC port 50142. Kafka consumer subscribes to payment events.
  - Consumer uses typed `evt.GetOrderId()` (no longer the `correlation_id` workaround)
  - `InitiatePayment` in order service passes all typed fields to payment service (no more untyped metadata map)

**What is real and working (C# — PoliSync):**

- `backend/polisync/src/PoliSync.Orders/` — C# order host (port 50140/50141)
  - Owns order creation, cart, quotation-to-order lifecycle
  - Publishes order events consumed by Go payment and orders services

**Billing (proto + gateway — service not yet implemented):**

- `billing/entity/v1/invoice.proto` — 27-field Invoice message in `billing_schema`
- `billing/services/v1/billing_service.proto` — 9 RPCs (CreateInvoice, GetInvoice, ListInvoices, MarkInvoicePaid, CancelInvoice, IssueInvoice, GetInvoicePDF, GenerateInvoicePDF, GetInvoiceByOrderId)
- `billing/events/v1/billing_events.proto` — 6 event types
- Gateway routes wired for all billing RPCs (port 50195 expected)
- **Billing Go microservice not yet created** — next implementation task

**What still needs attention:**

- billing microservice Go implementation (server + handler + service + repository)
- migration for `billing_schema.invoices` table
- `payment_schema.payments` migration to add the Phase 2 columns (23–53) if not run yet
- end-to-end SSLCommerz sandbox test with ngrok IPN
- receipt PDF generation (docgen-service integration)

**Original section:**

Current reality:

- Go payment service exists with repository, service, gRPC handler, health checks, and Kafka publisher support
- Go orders service exists and references payment integration
- PoliSync order host exists and is exposed through the gateway
- gateway already exposes `/v1/orders`

Interpretation:

Older docs that described payment as missing are outdated. The remaining work is ownership clarity, workflow completion, and end-to-end issuance/reconciliation alignment.

## 3.4 Fraud and Partner

Current reality:

- partner service contains repository, service, grpc, events, and startup structure
- fraud service contains repository, service, grpc, metrics, and startup structure
- gateway now includes partner and fraud routes

Interpretation:

The older fraud/partner plan describing fraud as absent and partner as mostly non-functional is no longer accurate for this branch. These domains are implemented enough to be described as partial-to-substantial, not blank.

## 3.5 DocGen and Media

Current reality:

- docgen service contains server, handler, repositories, service, worker, and template/document CRUD support
- media service contains server, handler, repositories, service, processing jobs, processors, worker, and Kafka publisher support
- gateway exposes document/media/storage routes already

Interpretation:

These are not greenfield workstreams. They are existing services needing hardening, infrastructure integration, and operational maturity.

## 3.6 Insurance service and mixed ownership

Current reality:

- there is a Go `insurance` service
- there is also PoliSync coverage for products, orders, policy, underwriting, claims, renewal, refund, and commission
- some domains therefore have both Go-side and C#-side presence

Interpretation:

This repo already has overlapping ownership boundaries in several business domains. Documentation must reflect that honestly instead of pretending a single runtime owns everything cleanly today.

## 4. Portal Reality

## 4.1 B2B portal

Current footprint:

- `15` page routes under `b2b_portal/app`
- `38` API routes under `b2b_portal/app/api`

The B2B portal is the only portal with a real backend-for-frontend layer that consistently proxies the gateway and uses the SDK.

## 4.2 Customer portal

Current footprint:

- `10` page routes
- data is mostly local constants and component-level arrays

Examples of static/demo evidence:

- `customer_portal/components/dashboard/employees/employees-table.tsx`
- `customer_portal/components/claims/claims-page.tsx`
- `customer_portal/components/payments/payments-page.tsx`
- `customer_portal/components/policies/policies-page.tsx`
- `customer_portal/lib/quotations.ts`
- `customer_portal/lib/payments.ts`

## 4.3 System portal

Current footprint:

- `11` `+page.svelte` routes
- has generated types and a server API helper
- many dashboard pages still read from demo datasets

Examples of demo-data evidence:

- `system_portal/src/routes/dashboard/products/+page.svelte`
- `system_portal/src/routes/dashboard/policies/+page.svelte`
- `system_portal/src/routes/dashboard/claims/+page.svelte`
- `system_portal/src/routes/dashboard/analytics/+page.svelte`
- `system_portal/src/lib/data_detailed/products_demo.ts`
- `system_portal/src/lib/data_detailed/policies_demo.ts`
- `system_portal/src/lib/data_detailed/claims_demo.ts`
- `system_portal/src/lib/data_detailed/analyticsData.ts`

## 5. Contract, API, and SDK Reality

The correct dependency direction is:

1. `proto/`
2. generated API/OpenAPI outputs in `api/`
3. generated SDKs in `sdks/`
4. portal and service consumption

Current implementation evidence:

- `api/ENDPOINT_MAP.md` is large and current
- `api/docs/` has generated HTML docs
- `sdks/insuretech-typescript-sdk/` is consumed by B2B and system portals
- generated proto/type output is checked into portals

Interpretation:

The repo is already operating as a contract-first system, but the documentation had become fragmented across too many derivative summaries.

## 6. Database and Migration Reality

The current rules are:

- proto-first schema definition
- SQL migrations for enhancements
- `backend/inscore/cmd/dbmanager` for migration and direct DB operations

Important implication:

The `dbmanager` note should not remain a standalone orphaned file. It is part of the platform baseline, not an isolated plan.

## 7. Key Corrections To Earlier Plan Narratives

These statements are now false or materially incomplete:

- "payment is missing"
- "storage is missing"
- "notification/audit/analytics/kyc need to be implemented from scratch"
- "fraud is not implemented"
- "partner has no real gateway integration"
- "document generation is not implemented"
- "bulk employee upload is missing"
- "purchase orders are mostly roadmap-only"

The more accurate statement is:

Many of these domains already exist in code, but vary in maturity, completeness, and portal adoption.

## 8. Architectural Risks That The Docs Must Not Hide

### 8.1 Mixed ownership across Go and C#

There is real overlap between:

- Go insurance/payment/orders layers
- PoliSync products/orders/policy/claim/underwriting/renewal/refund layers

This requires explicit ownership decisions in future work.

### 8.2 Uneven portal maturity

- B2B is integration-heavy
- customer portal is mostly UI shell
- system portal is mixed live/demo

Any document that speaks about "the portals" as if all are equally real is misleading.

### 8.3 Generated artifacts versus design source

OpenAPI docs and SDKs are downstream artifacts. They should not be treated as the design source in future documentation.

## 9. Canonical Reading Rule

When documentation conflicts with code, trust this order:

1. proto contracts
2. service registry in `services.yaml`
3. actual service entrypoints and handlers
4. gateway router
5. portal API routes and components
6. older planning documents
