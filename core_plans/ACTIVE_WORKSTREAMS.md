# Active Workstreams

## Rule

Most remaining work in this repo is not "start from zero" work. It is alignment, completion, hardening, and operationalization of code that already exists.

## 1. Payment And Order Alignment

> **Updated 2026-03-06** — Phase 2 largely complete. See `PAYMENT_ORDER_IMPLEMENTATION_PLAN.md` and `PROTO_AND_SDK_REFERENCE.md § 3.5` for detail.

### What Already Exists

- Go payment service in `backend/inscore/microservices/payment/` — 15 RPCs, SSLCommerz, manual review, receipts
- Go orders service/data layer in `backend/inscore/microservices/orders/` — typed event consumers
- PoliSync order host in `backend/polisync/src/PoliSync.Orders/`
- Gateway fully wired: 15 payment RPCs + 4 SSLCommerz HTTP callbacks + 9 billing routes
- Payment proto extended to 53 fields (order_id, tran_id, val_id, manual review, receipt fields)
- Billing proto created: `billing/entity`, `billing/services`, `billing/events`
- Kafka topics renamed to canonical `insuretech.payment.v1.*` format
- `orders/consumers` uses typed `evt.GetOrderId()` — metadata workaround removed
- `InitiatePayment` passes typed fields (order_id, customer_id, tenant_id) — no more untyped metadata map

### What Still Needs Attention

- **billing microservice Go implementation** — proto + gateway exist, service not yet created
- `billing_schema.invoices` DB migration
- `payment_schema.payments` Phase 2 column migration (fields 23–53) if not already run
- end-to-end SSLCommerz sandbox test with public ngrok IPN URL
- receipt PDF generation via docgen-service
- policy issuance trigger after `PaymentCompletedEvent` consumed

### Documentation Position

Implementation plan is now a completion backlog, not a design document. The system is largely built — remaining work is billing service implementation, integration testing, and operational hardening.

## 2. DocGen And Media Hardening

### What Already Exists

- document service server, repositories, handlers, and worker flow
- media service server, repositories, processors, worker flow, and Kafka publisher path
- B2B portal document and template routes already exist

### What Still Needs Attention

- production infrastructure assumptions
- async job reliability
- renderer/storage integration hardening
- performance, monitoring, and operational test coverage

## 3. Fraud And Partner Completion

### What Already Exists

- partner CRUD, verification, commission, and API-key related service logic
- fraud alert, case, and rule service logic
- repositories, gRPC handlers, and event publisher paths

### What Still Needs Attention

- broader gateway and portal adoption
- completeness against the latest business flows
- integration tests and operational readiness checks

## 4. Portal Convergence

### B2B

Continue extending the already-wired portal path rather than redesigning it.

### Customer

Move from hard-coded data to the same contract-driven API pattern used in B2B.

### System

Replace demo datasets with live SDK-backed queries incrementally by domain.

## 5. Auth And Authorization Baseline

AuthN and AuthZ are already central and live in the stack. Future work here should focus on:

- permission coverage consistency across portals and services
- seed-data/bootstrap discipline
- auditability of cross-portal role behavior

## Recommended Execution Order

1. contract and ownership cleanup where Go and PoliSync overlap
2. B2B stabilization on already-wired flows
3. docgen/media and fraud/partner hardening
4. customer and system portal integration work
