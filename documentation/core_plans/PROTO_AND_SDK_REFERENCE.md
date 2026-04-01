# Proto And SDK Reference

This is the consolidated contract deep-dive. It replaces the split between:

- `PROTO_INDEX.md`
- `PROTO_COMMON_TYPES.md`
- `PROTO_CORE_MODULES.md`
- `PROTO_FILES_SUMMARY.md`
- `COMPLETE_SDK_CONTENT_SUMMARY.md`
- `SDK_DOCUMENTATION.md`
- `PROTO_API_AND_SDK_REFERENCE.md`

## 1. Source Of Truth Model

The project is proto-first.

The correct dependency direction is:

1. `proto/`
2. generated gateway/OpenAPI view in `api/`
3. generated SDKs in `sdks/`
4. service and portal consumption

If a field or RPC definition matters, verify it in `proto/`.

## 2. Current Contract Estate

The repository currently contains about `180` proto files.

The module set currently spans `39` top-level areas, including:

- AI
- analytics
- API key
- audit
- authn
- authz
- B2B
- beneficiary
- billing
- claims
- commission
- common
- document
- endorsement
- fraud
- insurance
- insurer
- IoT
- KYC
- media
- MFS
- notification
- orders
- partner
- payment
- policy
- products
- refund
- renewal
- report
- services
- storage
- support
- task
- tenant
- underwriting
- voice
- webRTC
- workflow

## 3. Common Contract Building Blocks

The old proto docs correctly identified a few cross-cutting contract layers that still matter.

### 3.1 Common types

Shared common entities include:

- money
- address
- phone
- email
- document
- audit metadata

### 3.2 Error model

The contract set uses a shared error vocabulary with:

- standard error message structure
- field violations
- shared error codes

### 3.3 Security annotations

The common security proto layer captures field-level intent around:

- classification
- data category
- privacy handling
- security event semantics

### 3.4 Database annotations

Proto files also carry DB-shaping metadata such as:

- table options
- column options
- foreign keys
- indexes

This is why the database model is treated as proto-first rather than SQL-first.

## 3.5 Payment and billing proto field reference

> Updated 2026-03-06 after Phase 2 proto extension. All fields are live in `gen/` after `buf generate`.

### payment/entity/v1/payment.proto — Payment message (fields 1–53)

| # | Field | Type | DB Column | Notes |
|---|---|---|---|---|
| 1 | payment_id | string | UUID PK | |
| 2 | transaction_id | string | VARCHAR(255) | Internal reference |
| 3 | tigerbeetle_transfer_id | string | UUID | Ledger transfer ref |
| 4 | policy_id | string | UUID | FK to policies |
| 5 | claim_id | string | UUID | FK to claims |
| 6 | type | PaymentType | VARCHAR(50) | PREMIUM, CLAIM, etc. |
| 7 | method | PaymentMethod | VARCHAR(50) | CARD, BANK_TRANSFER |
| 8 | status | PaymentStatus | VARCHAR(50) | See enum below |
| 9 | amount | Money | BIGINT | Paisa |
| 10 | currency | string | VARCHAR(3) | ISO 4217 |
| 11 | payer_id | string | UUID | Customer user_id |
| 12 | payee_id | string | UUID | Platform account |
| 13 | initiated_at | Timestamp | TIMESTAMPTZ | |
| 14 | completed_at | Timestamp | TIMESTAMPTZ | |
| 15 | created_at | Timestamp | TIMESTAMPTZ | |
| 16 | updated_at | Timestamp | TIMESTAMPTZ | |
| 17 | gateway | string | VARCHAR(50) | SSLCOMMERZ, BKASH |
| 18 | gateway_response | string | JSONB | Raw provider response |
| 19 | receipt_url | string | TEXT | |
| 20 | retry_count | int32 | INT | |
| 21 | failure_reason | string | TEXT | |
| 22 | idempotency_key | string | VARCHAR(255) UNIQUE | CG-3 |
| 23 | order_id | string | UUID | FK to orders |
| 24 | invoice_id | string | UUID | FK to billing.invoices |
| 25 | tenant_id | string | UUID | Multi-tenant isolation |
| 26 | customer_id | string | UUID | Payer (may differ from payer_id for agent flows) |
| 27 | organisation_id | string | UUID | B2B context |
| 28 | purchase_order_id | string | UUID | B2B PO reference |
| 29 | provider | string | VARCHAR(50) | sslcommerz, bkash, nagad, manual |
| 30 | provider_reference | string | VARCHAR(255) | Canonical provider ref |
| 31 | tran_id | string | VARCHAR(255) UNIQUE | Merchant tran_id sent to SSLCommerz |
| 32 | val_id | string | VARCHAR(255) | SSLCommerz validation ID |
| 33 | session_key | string | VARCHAR(255) | SSLCommerz session key |
| 34 | bank_tran_id | string | VARCHAR(255) | Bank transaction ID |
| 35 | card_type | string | VARCHAR(50) | VISA, MASTERCARD, bKash |
| 36 | card_brand | string | VARCHAR(50) | |
| 37 | card_issuer | string | VARCHAR(100) | Issuing bank |
| 38 | card_issuer_country | string | VARCHAR(50) | |
| 39 | validated_at | Timestamp | TIMESTAMPTZ | When provider validated |
| 40 | validation_status | string | VARCHAR(50) | VALID, VALIDATED, INVALID |
| 41 | risk_level | string | VARCHAR(20) | 0=safe, 1=moderate, 2=high |
| 42 | risk_title | string | VARCHAR(100) | |
| 43 | callback_received_at | Timestamp | TIMESTAMPTZ | Browser callback |
| 44 | ipn_received_at | Timestamp | TIMESTAMPTZ | IPN callback |
| 45 | manual_review_status | ManualReviewStatus | VARCHAR(50) | NOT_REQUIRED/PENDING/APPROVED/REJECTED |
| 46 | manual_proof_file_id | string | UUID | Media service file ID |
| 47 | verified_by | string | UUID | Admin who approved |
| 48 | verified_at | Timestamp | TIMESTAMPTZ | |
| 49 | rejection_reason | string | TEXT | |
| 50 | receipt_number | string | VARCHAR(100) UNIQUE | RCP-YYYYMMDD-XXXXXXXX |
| 51 | receipt_document_id | string | UUID | Document-service doc ID |
| 52 | receipt_file_id | string | UUID | Media/storage file ID |
| 53 | ledger_transaction_id | string | VARCHAR(255) | TigerBeetle ref |

**PaymentStatus enum:**
`UNSPECIFIED(0)` `INITIATED(1)` `PENDING(2)` `PROCESSING(3)` `SUCCESS(4)` `FAILED(5)` `REFUNDED(6)` `CANCELLED(7)` `VERIFIED(8)` `MANUAL_REVIEW_REQUIRED(9)` `RECEIPT_PENDING(10)`

**ManualReviewStatus enum:**
`UNSPECIFIED(0)` `NOT_REQUIRED(1)` `PENDING(2)` `APPROVED(3)` `REJECTED(4)`

---

### payment/services/v1/payment_service.proto — PaymentService RPCs

| RPC | HTTP | Auth | Notes |
|---|---|---|---|
| InitiatePayment | POST /v1/payments | customer | Returns gateway_page_url for redirect |
| VerifyPayment | POST /v1/payments/{id}/verify | customer | SSLCommerz val_id server-side check |
| GetPayment | GET /v1/payments/{id} | any | |
| ListPayments | GET /v1/payments | any | |
| InitiateRefund | POST /v1/payments/{id}/refunds | agent/admin | |
| GetRefundStatus | GET /v1/refunds/{id} | any | |
| ListPaymentMethods | GET /v1/users/{id}/payment-methods | customer | |
| AddPaymentMethod | POST /v1/users/{id}/payment-methods | customer | |
| ReconcilePayments | POST /v1/payments/reconcile | admin | |
| HandleGatewayWebhook | POST /v1/payments/webhook/{provider} | **public** | SSLCommerz IPN — no auth |
| GetPaymentByProviderReference | GET /v1/payments/provider/{p}/references/{ref} | admin | |
| SubmitManualPaymentProof | POST /v1/payments/{id}/submit-proof | customer | Bank transfer proof upload |
| ReviewManualPayment | POST /v1/payments/{id}/review | agent/admin | Approve or reject proof |
| GenerateReceipt | POST /v1/payments/{id}/generate-receipt | any | Async PDF trigger |
| GetPaymentReceipt | GET /v1/payments/{id}/receipt | any | |

**Key request extensions (Phase 2 typed fields in InitiatePaymentRequest):**
`order_id` `invoice_id` `tenant_id` `customer_id` `organisation_id` `purchase_order_id`
`customer_name` `customer_email` `customer_phone` `customer_address_line1` `customer_city` `customer_postcode` `customer_country`

---

### payment/events/v1/payment_events.proto — Kafka events

| Topic | Event message | Key fields |
|---|---|---|
| `insuretech.payment.v1.payment.initiated` | PaymentInitiatedEvent | order_id, provider, tran_id, occurred_at |
| `insuretech.payment.v1.payment.completed` | PaymentCompletedEvent | order_id, val_id, receipt_number, occurred_at |
| `insuretech.payment.v1.payment.failed` | PaymentFailedEvent | order_id, provider, error_code, occurred_at |
| `insuretech.payment.v1.refund.processed` | RefundProcessedEvent | order_id, invoice_id, occurred_at |
| `insuretech.payment.v1.payment.verified` | PaymentVerifiedEvent | order_id, val_id, verified_at |
| `insuretech.payment.v1.payment.manual_review_requested` | ManualPaymentProofSubmittedEvent | order_id, manual_proof_file_id |
| `insuretech.payment.v1.payment.manual_review_completed` | ManualPaymentReviewedEvent | order_id, approved, rejection_reason |
| `insuretech.payment.v1.payment.receipt_generated` | ReceiptGeneratedEvent | order_id, receipt_number, receipt_file_id |
| `insuretech.payment.v1.payment.reconciliation_mismatch` | PaymentReconciliationMismatchEvent | order_id, expected_amount, actual_amount |

---

### billing/entity/v1/invoice.proto — Invoice message (fields 1–27)

| # | Field | Type | DB Column | Notes |
|---|---|---|---|---|
| 1 | invoice_id | string | UUID PK | |
| 2 | invoice_number | string | VARCHAR(100) UNIQUE | INV-2026-000001 |
| 3 | business_id | string | UUID | B2B: FK to organisations |
| 4 | amount | Money | BIGINT | Base amount (excl. tax) |
| 5 | due_date | Timestamp | TIMESTAMPTZ | |
| 6 | status | InvoiceStatus | VARCHAR(50) | See enum |
| 7 | issued_at | Timestamp | TIMESTAMPTZ | |
| 8 | paid_at | Timestamp | TIMESTAMPTZ | |
| 9 | payment_id | string | UUID | FK to payments |
| 10 | invoice_pdf_url | string | TEXT | S3 URL |
| 11 | policy_ids | []string | UUID[] | Policies on invoice |
| 12 | created_at | Timestamp | TIMESTAMPTZ | |
| 13 | updated_at | Timestamp | TIMESTAMPTZ | |
| 14 | order_id | string | UUID | FK to orders (B2C) |
| 15 | customer_id | string | UUID | FK to users |
| 16 | organisation_id | string | UUID | B2B FK |
| 17 | purchase_order_id | string | UUID | B2B PO FK |
| 18 | tenant_id | string | UUID NOT NULL | Multi-tenant key |
| 19 | tax_amount | Money | BIGINT | VAT/tax |
| 20 | total_amount | Money | BIGINT | amount + tax |
| 21 | currency | string | VARCHAR(3) | BDT default |
| 22 | notes | string | TEXT | Memo/billing notes |
| 23 | issued_by | string | UUID | Admin who issued |
| 24 | cancelled_at | Timestamp | TIMESTAMPTZ | |
| 25 | overdue_at | Timestamp | TIMESTAMPTZ | |
| 26 | metadata | map<string,string> | JSONB | Extension KVs |
| 27 | credit_note_id | string | UUID | Credit note ref |

**InvoiceStatus enum:**
`UNSPECIFIED(0)` `DRAFT(1)` `ISSUED(2)` `PENDING(3, legacy)` `APPROVED(4)` `PAID(5)` `OVERDUE(6)` `CANCELLED(7)` `CREDIT_NOTE_ISSUED(8)`

---

### billing/services/v1/billing_service.proto — BillingService RPCs

| RPC | HTTP | Notes |
|---|---|---|
| CreateInvoice | POST /v1/invoices | B2C or B2B |
| GetInvoice | GET /v1/invoices/{id} | |
| ListInvoices | GET /v1/invoices | Filter by customer/org/order/status |
| MarkInvoicePaid | POST /v1/invoices/{id}/mark-paid | Called by payment-service |
| CancelInvoice | POST /v1/invoices/{id}/cancel | Optional credit note |
| IssueInvoice | POST /v1/invoices/{id}/issue | DRAFT → ISSUED |
| GetInvoicePDF | GET /v1/invoices/{id}/pdf | Pre-signed URL |
| GenerateInvoicePDF | POST /v1/invoices/{id}/generate-pdf | Async trigger |
| GetInvoiceByOrderId | GET /v1/orders/{id}/invoice | Link lookup |

---

### orders/entity/v1/order.proto — Order message

| Field | Type | Notes |
|---|---|---|
| order_id | string UUID PK | |
| order_number | string | ORD-YYYYMMDD-XXXXXX |
| tenant_id | string | Multi-tenant |
| quotation_id | string | FK to quotations |
| customer_id | string | FK to users |
| product_id | string | FK to products |
| plan_id | string | FK to product_plans |
| status | OrderStatus | See enum |
| total_payable | Money | Paisa |
| currency | string | BDT default |
| payment_id | string | FK to payments |
| payment_gateway_ref | string | Provider ref |
| policy_id | string | FK to policies |
| cancellation_reason | string | |
| failure_reason | string | |
| created_at / updated_at / paid_at | Timestamp | |

**OrderStatus enum:**
`UNSPECIFIED(0)` `PENDING(1)` `PAYMENT_INITIATED(2)` `PAID(3)` `POLICY_ISSUED(4)` `CANCELLED(5)` `FAILED(6)`

---

## 4. Business-Critical Contract Areas

## 4.1 Identity and access

Primary modules:

- `insuretech/authn/`
- `insuretech/authz/`

Core concepts:

- users
- sessions
- OTP
- profiles
- API keys
- roles
- user-role assignments
- policy rules
- Casbin rules
- access audits

## 4.2 B2B

Primary modules:

- `insuretech/b2b/`

Core concepts:

- organizations
- organization members
- departments
- employees
- purchase orders

This is one of the most visibly consumed contract areas in the current portal code.

## 4.3 Commerce and policy lifecycle

Primary modules:

- `insuretech/products/`
- `insuretech/orders/`
- `insuretech/payment/`
- `insuretech/policy/`
- `insuretech/claims/`
- `insuretech/underwriting/`
- `insuretech/renewal/`
- `insuretech/refund/`
- `insuretech/endorsement/`
- `insuretech/commission/`

Important implementation note:

These modules are consumed across both Go services and PoliSync, so contract clarity matters more here than in any single-runtime explanation.

## 4.4 Support and document operations

Primary modules:

- `insuretech/document/`
- `insuretech/storage/`
- `insuretech/support/`
- `insuretech/report/`
- `insuretech/workflow/`
- `insuretech/task/`

## 4.5 Risk and ecosystem modules

Primary modules:

- `insuretech/fraud/`
- `insuretech/partner/`
- `insuretech/beneficiary/`
- `insuretech/insurer/`
- `insuretech/insurance/`

## 5. API Projection

`api/` is the HTTP/OpenAPI projection of the contract set.

Important files:

- `api/openapi.yaml`
- `api/ENDPOINT_MAP.md`
- `api/docs/`

Current reality:

- the generated API docs are broad and active
- the B2B, payment, policy, underwriting, partner, fraud, workflow, support, insurer, and many other domains are already represented there
- these generated docs should be treated as evidence of the live contract surface, not as the primary design source

## 6. SDK Output

Current generated SDK locations:

- `sdks/insuretech-typescript-sdk/`
- `sdks/insuretech-go-sdk/`

## 6.1 TypeScript SDK shape

The TypeScript SDK includes:

- public entrypoints
- client wrapper
- generated client modules
- generated type modules
- core serialization/auth utilities

This is why the older SDK summary docs spent time on:

- `src/index.ts`
- `src/client-wrapper.ts`
- `src/client.gen.ts`
- `src/sdk.gen.ts`
- `src/types.gen.ts`
- generated core helpers

Those details are still useful, but no longer need separate top-level docs.

## 6.2 Portal SDK consumption

### B2B portal

The B2B portal consumes the TypeScript SDK through:

- local package dependency
- BFF route wrappers
- `src/lib/sdk/*` convenience clients

### System portal

The system portal also depends on the TypeScript SDK, though its UI still often renders demo data instead of live API data.

## 7. Generation Model

The intended workflow is:

1. update proto contracts
2. run code generation
3. regenerate OpenAPI outputs
4. regenerate SDK artifacts
5. update consuming portal or service code

This keeps contracts, gateway docs, and SDKs aligned.

## 8. Developer Rules

## 8.1 Do first

- edit `proto/`

## 8.2 Do not do first

- do not hand-edit generated API docs as if they were the source
- do not hand-edit generated SDK output unless you are intentionally fixing the generator pipeline

## 8.3 When trying to understand a domain

Use this order:

1. the relevant proto package
2. generated API map
3. generated SDK/client surface
4. consuming code in portals or services

## 9. Why These Docs Were Consolidated

The old proto and SDK doc cluster had good information, but it was fragmented by perspective:

- index perspective
- common-type perspective
- module-summary perspective
- SDK package perspective

That made maintenance harder and increased the risk that summaries would drift from the real proto tree.

This document keeps the useful structure but collapses it into one contract-focused reference.

## 10. Canonical Use

Use this document when you need:

- the big picture of the contract estate
- the relationship between proto, API, and SDK generation
- a starting point before drilling into specific proto packages
