# Insurance Engine — Technical Documentation
**LabAid InsureTech Platform · SRS v3.11**
*Service: Insurance Engine | Language: C# .NET 8 | Port: 5001*

---

## Table of Contents

1. [System Overview](#1-system-overview)
2. [System Architecture](#2-system-architecture)
3. [Technology Stack](#3-technology-stack)
4. [Module Breakdown](#4-module-breakdown)
   - 4.1 [Product Management & Catalog (FG-003)](#41-product-management--catalog-fg-003)
   - 4.2 [Policy Lifecycle Management (FG-004)](#42-policy-lifecycle-management-fg-004)
   - 4.3 [Policy Management & Renewals (FG-005)](#43-policy-management--renewals-fg-005)
   - 4.4 [Policy Cancellation & Refund (FG-005.1)](#44-policy-cancellation--refund-fg-0051)
   - 4.5 [Policy Endorsement & Amendment (FG-005.2)](#45-policy-endorsement--amendment-fg-0052)
   - 4.6 [Business Rules & Workflows (FG-006)](#46-business-rules--workflows-fg-006)
   - 4.7 [Claims Management (FG-008)](#47-claims-management-fg-008)
   - 4.8 [Fraud Detection & Risk Controls (FG-016)](#48-fraud-detection--risk-controls-fg-016)
5. [State Machines](#5-state-machines)
6. [Claims Approval Matrix](#6-claims-approval-matrix)
7. [Database Schema](#7-database-schema)
8. [Kafka Event Topology](#8-kafka-event-topology)
9. [gRPC Service Contracts](#9-grpc-service-contracts)
10. [Internal Integration Points](#10-internal-integration-points)
11. [FR ID Reference Index](#11-fr-id-reference-index)

---

## 1. System Overview

The **Insurance Engine** is the core business-logic microservice of the LabAid InsureTech Platform. It owns the complete insurance value chain — from product catalog management through to claim settlement — and is the single source of truth for all policy and claims state.

### Responsibility Summary

| Concern | Owned by Insurance Engine |
|---|---|
| Product catalog & pricing | ✅ Yes |
| Policy issuance & lifecycle | ✅ Yes |
| Policy renewal, lapse, reinstatement | ✅ Yes |
| Policy cancellation & endorsements | ✅ Yes |
| Claims submission & approval workflows | ✅ Yes |
| Fraud detection rules | ✅ Yes |
| Premium calculation & business rules | ✅ Yes |
| Payment processing | ❌ No — Payment Service (Node.js, Port 3001) |
| Partner & agent management | ❌ No — Partner Management (C# .NET, Port 5002) |
| Analytics dashboards | ❌ No — Analytics & Reporting (C# .NET, Port 5003) |
| Notifications dispatch | ❌ No — Kafka → Kafka Service (Go, Port 8086) |

### Business Context

The Insurance Engine supports Bangladesh's micro-insurance market (200–2,000 BDT premiums) and manages:

- **Product categories:** Health, Life, Motor, Travel, Micro-insurance (Phase M1); Fire, Device, Livestock, Fisheries, Crop, Goods in Transit, Pet, Home Appliance (Phase M2)
- **Policy number format:** `LBT-YYYY-XXXX-NNNNNN` (e.g., `LBT-2025-0001-000123`)
- **Claim number format:** `CLM-YYYY-XXXX-NNNNNN` (e.g., `CLM-2025-0001-000045`)
- **Compliance:** IDRA & BFIU regulatory alignment

---

## 2. System Architecture

### 2.1 Position in Microservices Ecosystem

```
┌─────────────────────────────────────────────────────────────┐
│                    CLOUDFLARE PROXY                         │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│              API GATEWAY (Go, Port 8080)                    │
│         OAuth2 + JWT · Rate Limiting · Routing              │
└─────────────────────┬───────────────────────────────────────┘
                      │  gRPC (Category 1 — internal)
                      │
          ┌───────────▼──────────────┐
          │   INSURANCE ENGINE       │  ◄── THIS SERVICE
          │   C# .NET 8 · Port 5001  │
          │   VSA + CQRS/MediatR     │
          └───────┬──────────┬───────┘
                  │          │
         Publishes│          │Calls via gRPC
         Kafka    │          │
         Events   │     ┌────▼──────────────────┐
                  │     │  Payment Service       │
                  │     │  Node.js · Port 3001   │
                  │     └────────────────────────┘
                  │
         ┌────────▼────────────────────────────────┐
         │           Apache Kafka                  │
         │  (Domain Events → Notification, etc.)   │
         └─────────────────────────────────────────┘
```

### 2.2 Internal Architecture Pattern: VSA + CQRS/MediatR

The Insurance Engine uses **Vertical Slice Architecture (VSA)** internally, where every feature (e.g., "Issue Policy", "Submit Claim") is a self-contained vertical slice owning all its layers.

```
InsuranceEngine/
│
├── Features/
│   ├── Products/
│   │   ├── CreateProduct/
│   │   │   ├── CreateProductCommand.cs
│   │   │   ├── CreateProductCommandHandler.cs
│   │   │   └── CreateProductValidator.cs
│   │   ├── GetProduct/
│   │   │   ├── GetProductQuery.cs
│   │   │   └── GetProductQueryHandler.cs
│   │   └── ...
│   │
│   ├── Policies/
│   │   ├── IssuePolicy/
│   │   ├── RenewPolicy/
│   │   ├── CancelPolicy/
│   │   ├── LapsePolicy/
│   │   ├── EndorsePolicy/
│   │   └── ...
│   │
│   ├── Claims/
│   │   ├── SubmitClaim/
│   │   ├── ReviewClaim/
│   │   ├── ApproveClaim/
│   │   ├── RejectClaim/
│   │   ├── SettleClaim/
│   │   └── ...
│   │
│   └── BusinessRules/
│       ├── CalculatePremium/
│       ├── ValidateEligibility/
│       └── DetectFraud/
│
├── Domain/
│   ├── Entities/           ← Policy, Claim, Product, Applicant …
│   ├── Events/             ← Domain events (PolicyIssued, ClaimSubmitted …)
│   ├── Enums/              ← PolicyStatus, ClaimStatus, ClaimType …
│   └── ValueObjects/       ← Money, PolicyNumber, ClaimNumber …
│
├── Infrastructure/
│   ├── Persistence/        ← EF Core DbContext, Migrations, Repositories
│   ├── Kafka/              ← Event producers
│   ├── gRPC/               ← Proto-generated stubs + service impl
│   └── ExternalClients/    ← Payment Service client, AI Engine client
│
└── Common/
    ├── Behaviours/         ← MediatR pipeline (validation, logging, etc.)
    └── Exceptions/
```

**CQRS separation:**

| Side | Role | Storage |
|---|---|---|
| **Command** (Write) | Mutates state via MediatR commands | PostgreSQL (primary) → Kafka event |
| **Query** (Read) | Reads from optimized projections via MediatR queries | PostgreSQL read replica / Redis cache |

---

## 3. Technology Stack

| Layer | Technology | Purpose |
|---|---|---|
| **Runtime** | C# .NET 8 | Core application runtime |
| **Architecture pattern** | Vertical Slice Architecture (VSA) | Feature-oriented code organization |
| **Mediator** | MediatR | In-process CQRS messaging (commands, queries, notifications) |
| **Validation** | FluentValidation | Request/command validation pipeline |
| **ORM** | Entity Framework Core (EF Core) | Database access and migrations |
| **Primary DB** | PostgreSQL 17 (JSONB support) | Transactional data — policies, claims, products |
| **Cache** | Redis 7.0+ | Product catalog cache (5-min TTL), session data |
| **Financial ledger** | TigerBeetle | Double-entry bookkeeping for premium/claim transactions |
| **Message broker** | Apache Kafka | Async domain event publishing |
| **Service contract** | Protocol Buffers (proto3) | All data models and gRPC service definitions |
| **Inter-service comms** | gRPC (Category 1) | < 100 ms gateway ↔ microservice communication |
| **Object storage** | S3-compatible (AWS/DigitalOcean) | Policy PDFs, claim documents |
| **Serialization** | System.Text.Json + Protobuf | REST responses + gRPC payloads |
| **Logging** | Structured logging → ELK Stack | Centralized log aggregation |
| **Tracing** | Jaeger (OpenTelemetry) | Distributed tracing with trace ID propagation |
| **Metrics** | Prometheus + Grafana | Service health, business KPIs |

### Performance Targets (Insurance Engine)

| Metric | Target |
|---|---|
| gRPC API response (p95) | < 100 ms |
| Policy issuance end-to-end | < 30 s (including PDF generation) |
| Premium calculation | < 2 s |
| Product catalog load (cached) | < 500 ms |
| Report generation | < 30 s |
| DB query average | < 100 ms |

---

## 4. Module Breakdown

### 4.1 Product Management & Catalog (FG-003)

**Priority:** M1 (core) / M2 (extended categories)
**FR IDs:** FR-021 through FR-029

This module manages the insurance product catalog — creation, configuration, versioning, and display.

#### Supported Product Categories

**Phase M1:**
- Health Insurance
- Auto (Motor) Insurance
- Travel Insurance
- Life Insurance

**Phase M2 (additional):**
- Fire Insurance
- Device Insurance
- Livestock Insurance
- Fisheries Insurance
- Crop Insurance
- Goods in Transit
- Pet Insurance
- Home Appliance Insurance

#### Key Functional Requirements

| FR ID | Requirement | Priority |
|---|---|---|
| FR-021 | Product catalog with M1 categories (Health → Auto → Travel → Life) | M1 |
| FR-021-A | Extended M2 categories (Fire, Device, Livestock, Fisheries, etc.) | M2 |
| FR-022 | Product search by name, category, coverage type, premium range; < 500 ms; fuzzy Bengali text matching | M1 |
| FR-023 | Display product details: coverage, premium, tenure, exclusions, T&Cs; PDF download | M2 |
| FR-023-A | Unit-wise plan purchase; user can increase/decrease coverage amount | M2 |
| FR-023-B | Risk assessment questionnaire for every plan (multi-question flow) | M2 |
| FR-024 | Premium calculator with dynamic inputs (age, sum assured, tenure, riders); real-time < 2 s | M3 |
| FR-025 | Side-by-side product comparison (up to 2 products) | M3 |
| FR-026 | Business Admin can CRUD products; version history maintained | M1 |
| FR-027 | Product variants with configurable riders and add-ons; dynamic pricing recalculation | D |
| FR-028 | Redis cache for product catalog; 5-minute TTL; auto-invalidation on update | M3 |
| FR-029 | Multi-language product descriptions (Bengali primary, English secondary) | M3 |

#### Domain Behaviour Notes

- Products carry **base premium** + optional **rider premiums**; total is recalculated dynamically when riders are toggled (FR-027).
- Coverage amount is **unit-based** — the user adjusts units to increase/decrease the sum insured (FR-023-A).
- Every product has a **risk assessment questionnaire** attached (FR-023-B); answers feed into underwriting logic.
- Product catalog is cached in **Redis (TTL = 5 min)** and invalidated on any product update event.
- Business Admin controls product lifecycle: Draft → Active → Deactivated. Version history is immutable.

---

### 4.2 Policy Lifecycle Management (FG-004)

**Priority:** M1/M2
**FR IDs:** FR-030 through FR-041

This is the most critical module — it handles the complete policy purchase journey from product selection to policy issuance, plus the ongoing policy dashboard.

#### Policy Purchase Flow

```
Product Selection
      │
      ▼
Applicant Details
(name, DOB, NID*, address, occupation, income, health declaration)
      │
      ▼
Nominee Details
(1 nominee required — FR-032)
      │
      ▼
Risk Assessment Questions answered
      │
      ▼
Premium Calculation
      │
      ▼
Payment (delegated to Payment Service)
      │
      ▼
Payment Confirmation received (Kafka event)
      │
      ▼
Policy Issuance
  ├── Policy number generated: LBT-YYYY-XXXX-NNNNNN
  ├── Policy status → ACTIVE (Non-Life: immediate upon payment — FR-037)
  ├── PDF policy document generated (< 30 s) with QR code
  └── Events published → Notification service (SMS + email delivery < 5 min)
```

#### Key Functional Requirements

| FR ID | Requirement | Priority |
|---|---|---|
| FR-030 | End-to-end purchase flow; < 10 min; progress saved at each step | M1 |
| FR-031 | Applicant info: full name, DOB, NID (optional), address, occupation, income, health declaration | M1 |
| FR-032 | Single nominee/beneficiary required | M1 |
| FR-032-A | Nominee income range is optional | M1 |
| FR-033 | NID/Mobile uniqueness validation; duplicate policy prevention | M1 |
| FR-034 | Unique policy number: `LBT-YYYY-XXXX-NNNNNN`; sequential; collision-free | M1 |
| FR-035 | Digital policy PDF with QR code; generated < 30 s of payment confirmation | M2 |
| FR-036 | Policy document delivery via SMS link + email attachment; < 5 min; retry on failure | M2 |
| FR-037 | Immediate activation on payment confirmation (Non-Life); real-time status update | M2 |
| FR-038 | Cooling-off period: 5 days from issuance for full refund | M3 |
| FR-039 | Policy statuses: `Pending Payment`, `Active`, `Suspended`, `Cancelled`, `Lapsed`, `Expired`; transitions logged | M1 |
| FR-040 | Customer policy dashboard: all active/past policies, renewal prompts, payment history; < 3 s load | M1 |
| FR-041 | Order history with: Coverage details, Refer, Active Plans, Claimed Plans, Expired Plans; referral (max 1) | D |

#### Policy Number Generation Logic

```
Format: LBT-{YEAR}-{PRODUCT_CODE}-{SEQUENCE}

Example: LBT-2025-0001-000123
          │    │    │    └── Sequence (zero-padded 6 digits, per product per year)
          │    │    └─────── Product code (4 digits)
          │    └──────────── Year (4 digits)
          └───────────────── Company prefix
```

Collision prevention: DB-level `UNIQUE` constraint + sequence generator (never rely on application-level only).

---

### 4.3 Policy Management & Renewals (FG-005)

**Priority:** M2/M3
**FR IDs:** FR-042 through FR-050

Handles everything that happens to a policy after issuance: renewals (manual and automated), grace periods, lapses, and reinstatements.

#### Key Functional Requirements

| FR ID | Requirement | Priority |
|---|---|---|
| FR-042 | Family Insurance Wallet — group policies for up to 6 family members under one account | D |
| FR-043 | Renewal reminders: 30 days, 7 days, 1 day before expiry via SMS, email, push | M2 |
| FR-044 | Manual renewal: one-click, reuses existing policy data; < 3 min; new policy document issued | M2 |
| FR-045 | Auto-repurchase with stored payment method (opt-in); charged 7 days before expiry | M3 |
| FR-046 | Editable fields during renewal: current address, nominee info; verification for major changes | M3 |
| FR-047 | Grace period: 30 days post-expiry; coverage continues; daily reminders; status = `Grace Period` | M2 |
| FR-048 | Auto-lapse after grace period if unpaid; reinstatement within 90 days | M2 |
| FR-049 | Policy document PDF download with version history for all renewals | M1 |
| FR-050 | Lifecycle event audit trail: issuance, renewal, lapse, reinstatement, cancellation | M1 |

#### Renewal State Flow

```
ACTIVE
  │  (30 days before expiry — reminder sent)
  │  (7 days before expiry — reminder sent)
  │  (1 day before expiry — reminder sent)
  │
  ▼ (on expiry date — payment not received)
GRACE PERIOD  ──── (payment received within 30 days) ────► ACTIVE (renewed)
  │
  │ (30 days passed, still no payment)
  ▼
LAPSED
  │
  │ (within 90 days of lapse — reinstatement request + medical underwriting)
  ▼
ACTIVE (reinstated)  OR  remains LAPSED forever after 90-day window
```

---

### 4.4 Policy Cancellation & Refund (FG-005.1)

**Priority:** M1
**FR IDs:** FR-051 through FR-055

Handles voluntary cancellation requests by customers, agents, or admins, including pro-rata refund calculation.

#### Key Functional Requirements

| FR ID | Requirement | Priority |
|---|---|---|
| FR-051 | Cancellation workflow: request submission with reason dropdown + attachment support | M1 |
| FR-052 | Approval workflow: Business Admin + Focal Person approval required for policies > 30 days old; 48 hr SLA | M1 |
| FR-053 | Pro-rata refund formula: `(Premium Paid − Days Covered − Admin Fee − Cancellation Charge)`; transparent breakdown | M1 |
| FR-054 | Refund processing within 7 working days via MFS or bank transfer | M1 |
| FR-055 | Policy status → `CANCELLED`; all stakeholders notified; IDRA reporting triggered | M1 |

#### Cancellation Approval Rules

```
Request submitted by: Customer / Agent / Admin
         │
         ▼
Policy age ≤ 30 days?
    ├── YES → Cooling-off period — auto-approve, full refund (FR-038)
    └── NO  → Requires: Business Admin approval + Focal Person approval (48 hr SLA)
                    │
                    ▼
              Pro-rata refund calculated
                    │
                    ▼
              Refund initiated (7 working days) via original payment channel
                    │
                    ▼
              Policy status → CANCELLED
              IDRA report event triggered
```

---

### 4.5 Policy Endorsement & Amendment (FG-005.2)

**Priority:** M1/M2/D
**FR IDs:** FR-056 through FR-060

Handles mid-term changes (endorsements) to an active policy.

#### Key Functional Requirements

| FR ID | Requirement | Priority |
|---|---|---|
| FR-056 | Endorsement types: Address, Sum Insured, Nominee, Contact changes | M1 |
| FR-057 | Additional premium calculated for mid-term sum insured increase | D |
| FR-058 | Pro-rata refund/credit calculated for sum insured decrease | M2 |
| FR-059 | Endorsement document generated with suffix: `PLN-001/END-01`; PDF; version tracked | M1 |
| FR-060 | Approval required for sum insured changes > 10% of original | M1 |

#### Endorsement Document Naming

```
Original Policy:   LBT-2025-0001-000123
First Endorsement: LBT-2025-0001-000123/END-01
Second Endorsement:LBT-2025-0001-000123/END-02
```

---

### 4.6 Business Rules & Workflows (FG-006)

**Priority:** M1/M2/M3/D
**FR IDs:** FR-061 through FR-069

Contains critical business logic: premium calculation fallbacks, duplicate detection, claim state machine definition, and grace period logic.

#### Key Functional Requirements

| FR ID | Requirement | Priority |
|---|---|---|
| FR-061 | Premium fallback: if insurer API fails → use cached rates (max 24 hrs old); if unavailable → queue + notify customer within 2 hrs | M1 |
| FR-062 | Premium edge cases: age-based loading, occupation risk factors, pre-existing conditions | M2 |
| FR-063 | Duplicate policy detection: block same product + same insured within 30 days; allow cross-product | M1 |
| FR-064 | Policy merge workflow: Focal Person merges duplicate accounts after NID verification; transfers policies; consolidates claims history | M3 |
| FR-065 | Claim status state machine: `Submitted → Under Review → Documents Requested → Approved/Rejected → Payment Initiated → Settled/Closed` | M1 |
| FR-066 | Claim state transition rules: auto-move to `Documents Requested` if incomplete; Business Admin + Focal Person approval for > BDT 50K | M1 |
| FR-067 | Gamified renewal rewards: discounts or gift vouchers for early renewals | D |
| FR-068 | Grace period logic: 30-day post-expiry; coverage continues; auto-lapse after grace | M3 |
| FR-069 | Lapsed policy reinstatement: within 90 days; medical underwriting required; Focal Person approval | D |

#### Premium Calculation Logic

```
Base Premium
    + Age Loading (if applicable)
    + Occupation Risk Factor
    + Pre-existing Condition Loading
    + Rider Premiums (each configured separately)
    ─────────────────────────────
    = Total Premium

Fallback chain:
  1. Live insurer API (preferred)
  2. Cached rates (max 24 hrs stale) ← used if API fails
  3. Queue + notify customer within 2 hrs ← used if cache unavailable
```

---

### 4.7 Claims Management (FG-008)

**Priority:** M1/M2/M3/D
**FR IDs:** FR-081 through FR-101

This is the second most critical module. It owns the complete claims lifecycle from submission through settlement.

#### Claim Submission Flow

```
Customer selects active policy
        │
        ▼
Claim eligibility validation (< 3 s)
  ├── Policy active? ✓
  ├── Within coverage period? ✓
  ├── Claim type covered by product? ✓
  └── No duplicate submission for same incident? ✓
        │
        ▼
Fixed-step claim form:
  Step 1: Policy selection
  Step 2: Incident details (date, type, description)
  Step 3: Claim reason selection
  Step 4: Document upload (images, bills, reports)
        │
        ▼
Claim number generated: CLM-YYYY-XXXX-NNNNNN
SHA-256 hash generated for document integrity
        │
        ▼
Auto-routed to correct approval level (see Claims Approval Matrix)
        │
        ▼
Partner/insurer notified within 60 s
        │
        ▼
Real-time status tracking visible to customer
```

#### Key Functional Requirements

| FR ID | Requirement | Priority |
|---|---|---|
| FR-081 | Fixed-step claim form; draft saving; claim tracker shown; < 5 min completion | M1 |
| FR-082 | Eligibility validation: policy active, coverage period, claim type, no duplicate | M1 |
| FR-083 | Unique claim number `CLM-YYYY-XXXX-NNNNNN`; SHA-256 document hash | M1 |
| FR-084 | Auto-notify partner/insurer on submission; shared status dashboard; < 60 s | M2 |
| FR-085 | Real-time status tracking; customer push + SMS notifications on status change | M3 |
| FR-086 | Tiered approval workflow per Claims Approval Matrix | M3 |
| FR-087 | Document verification: image quality check, OCR extraction, fraud detection; < 10 s | M3 |
| FR-088 | Chat interface: customer ↔ partner agent ↔ focal person; real-time + file attachment | M3 |
| FR-089 | WebRTC video call for claim verification; call recording for audit | D |
| FR-090 | Partner can add verification notes; approve/reject with mandatory reason | M2 |
| FR-091 | Joint approval: Business Admin + Focal Person for claims BDT 50K–2L | M3 |
| FR-092 | Auto-payment upon approval via customer's selected channel; < 24 hrs | M3 |
| FR-093 | Zero Human Touch Claims (ZHTC): auto-verify + pay claims < BDT 10K (partner pre-agreement); 95% automation | D |
| FR-094 | Fraud detection: frequent claims > 3 in 6 months, duplicate docs, policy-to-claim < 48 hrs | M3 |
| FR-095 | Auto-revoke customer access for confirmed fraud; appeal process available | M3 |
| FR-096 | Balance sheet on Customer/Partner/Agent/InsureTech level for selected periods | M3 |
| FR-097 | TAT tracking per approval level; SLA breach alerts | M3 |
| FR-098 | Claim history and analytics for risk assessment and premium adjustment | M3 |

#### Document Requirements (FR-099)

| Constraint | Value |
|---|---|
| Accepted formats | PDF, JPG, PNG |
| Max file size | 5 MB per file |
| Max total per claim | 25 MB |
| Minimum image resolution | 300 DPI |
| Upload compression | Client-side: 5 MB → 1–2 MB (JPEG 80%, 1920×1080 max) |
| Upload method | Chunked (1 MB chunks, tus.io resume protocol) |
| Direct upload | Presigned S3 URLs (30-min expiry) |

#### Co-payment & Deductible Calculation (FR-100)

```
Net Reimbursement = (Claim Amount − Deductible) × Co-payment Percentage

Where:
  Deductible = Product-level config (tracked annually)
  Co-payment % = Product-level config

Example:
  Claim Amount:   BDT 30,000
  Deductible:     BDT 2,000   (already exhausted this year: 0)
  Co-payment:     80%
  Net:            (30,000 − 0) × 0.80 = BDT 24,000
```

#### Claim Reimbursement Timeline (FR-101)

Insurance company processes reimbursement within **7–15 working days** via document review and bank/MFS transfer.

---

### 4.8 Fraud Detection & Risk Controls (FG-016)

**Priority:** M2/M3
**FR IDs:** FR-175 through FR-188

The Insurance Engine implements real-time fraud detection rules that auto-flag suspicious claims and policy activities.

#### Key Functional Requirements

| FR ID | Requirement | Priority |
|---|---|---|
| FR-175 | Auto-flag claims submitted within 48 hrs of policy purchase; route to manual review | M2 |
| FR-176 | Detect same claim type > 2 times in 12 months; flag for pattern analysis | M2 |
| FR-177 | Flag claims where amount exactly matches 100% of policy limit | M2 |
| FR-178 | Validate medical provider against approved network; flag non-network claims | M2 |
| FR-179 | Device fingerprinting: flag > 3 accounts from same device | M3 |
| FR-180 | Fraud dashboard for Business Admin + Focal Person; drill-down; real-time alerts | M2 |
| FR-181 | RACI matrix enforced for monitoring and incident escalation | M1 |

#### Fraud Detection Rules

| Rule ID | Rule Description | Threshold | Action |
|---|---|---|---|
| FR-182 | Rapid Policy-to-Claim | < 48 hours | Auto-flag + manual review |
| FR-183 | Frequent Claims | > 2 same claim type in 12 months | Flag + pattern analysis |
| FR-184 | Amount Matching | 100% of coverage | Flag + enhanced verification |
| FR-185 | Non-Network Provider | Medical provider not in approved list | Flag + provider verification |
| FD-186 | Geographic Anomaly | Claim location > 100 km from registered address | Flag + location verification |
| FD-187 | Device Fingerprinting | > 3 accounts from same device | Flag + identity verification |
| FD-188 | Behavioral Pattern | ML-based anomaly scoring | Risk score + monitoring queue |

#### Customer Risk Scoring Matrix

| Risk Factor | Points |
|---|---|
| Transaction Frequency (> 3 claims in 6 months) | +20 |
| Transaction Amount (single > BDT 50K) | +15 |
| Geographic Anomaly (claim location far from address) | +10 |
| KYC Completeness (missing NID verification) | +25 |
| Device Fingerprinting (multiple accounts, same device) | +15 |
| Behavioral Anomaly (unusual activity patterns) | +10 |

| Score Range | Risk Category | Review Frequency |
|---|---|---|
| 0–30 | Low Risk | Annual review |
| 31–60 | Medium Risk | Semi-annual review |
| > 60 | High Risk | Quarterly review + enhanced monitoring |

---

## 5. State Machines

### 5.1 Policy Status State Machine

```
                    ┌─────────────────┐
                    │  PENDING_PAYMENT │
                    └────────┬────────┘
                             │ Payment confirmed
                             ▼
                    ┌────────────────┐
               ┌───►    ACTIVE       ◄───────────────┐
               │    └───────┬────────┘               │
               │            │                        │ Reinstatement
               │            │ Expiry date reached    │ (within 90 days)
               │            │ + No payment           │
               │            ▼                        │
               │    ┌──────────────────┐             │
               │    │   GRACE_PERIOD   │             │
               │    └───────┬──────────┘             │
               │            │ 30 days pass, no pay   │
               │            ▼                        │
               │    ┌────────────────┐               │
               │    │     LAPSED     ├───────────────┘
               │    └────────────────┘
               │
               │ Suspension trigger (admin)
               ▼
      ┌──────────────────┐
      │    SUSPENDED     │
      └──────────────────┘

      ┌──────────────────┐
      │    CANCELLED     │  ← Cancellation workflow (FG-005.1)
      └──────────────────┘

      ┌──────────────────┐
      │     EXPIRED      │  ← End of tenure, policy not renewed
      └──────────────────┘
```

**Valid transitions:**

| From | To | Trigger |
|---|---|---|
| PENDING_PAYMENT | ACTIVE | Payment confirmed |
| ACTIVE | GRACE_PERIOD | Expiry date reached, no payment |
| ACTIVE | CANCELLED | Cancellation workflow approved |
| ACTIVE | SUSPENDED | Admin suspension |
| ACTIVE | EXPIRED | Tenure ended, not renewed |
| GRACE_PERIOD | ACTIVE | Payment received within 30 days |
| GRACE_PERIOD | LAPSED | 30 days passed, still no payment |
| LAPSED | ACTIVE | Reinstatement within 90 days + medical underwriting |
| SUSPENDED | ACTIVE | Admin re-activation |

### 5.2 Claim Status State Machine

```
Customer submits claim
         │
         ▼
    SUBMITTED
         │
         │ Auto-validation passes
         ▼
    UNDER_REVIEW
         │
    ┌────┴──────────────┐
    │                   │
    │ Docs complete      │ Docs incomplete
    │                   ▼
    │           PENDING_DOCUMENTS
    │                   │
    │           (docs uploaded)
    │                   │
    └────────►──────────┘
         │
    ┌────┴─────┐
    │          │
    ▼          ▼
APPROVED    REJECTED
    │
    │ Payment initiated
    ▼
PAYMENT_INITIATED
    │
    │ Payment confirmed
    ▼
  SETTLED ──────────────────────────────► CLOSED
               (all activities done)
```

**Transition rules (FR-065, FR-066):**

| Transition | Rule |
|---|---|
| SUBMITTED → UNDER_REVIEW | Automatic on submission |
| UNDER_REVIEW → PENDING_DOCUMENTS | Triggered if documents are incomplete |
| UNDER_REVIEW → APPROVED | Requires approval per Claims Approval Matrix |
| UNDER_REVIEW → REJECTED | Requires rejection reason (mandatory field) |
| APPROVED → PAYMENT_INITIATED | Automatic within 24 hrs of approval |
| PAYMENT_INITIATED → SETTLED | Payment gateway confirms settlement |
| Any → DISPUTED | Customer disputes the decision |

---

## 6. Claims Approval Matrix

*(Only the Insurance Company approves — FR-086, FR-091)*

| Claimed Amount | Approval Level | Approver(s) | Maximum TAT |
|---|---|---|---|
| BDT 0 – 10,000 | L1 — Auto / Officer | System Auto-Approval **OR** Claims Officer | 24 hours |
| BDT 10,001 – 50,000 | L2 — Manager | Claims Manager | 3 days |
| BDT 50,001 – 2,00,000 | L3 — Head | Business Admin + Focal Person (Joint — **both required**) | 7 days |
| BDT 2,00,001+ | Board | Board + Insurer Approval | 15 days |

**Rules:**
- Joint approval (L3): Both Business Admin AND Focal Person must approve; either alone is insufficient.
- Timeout escalation: If approver does not act within TAT, auto-escalate to next level.
- TAT tracked per approval level (FR-097); SLA breach triggers email alert.

---

## 7. Database Schema

All entities are derived from proto3 definitions (proto-first strategy).

### 7.1 Core Tables

#### `products`
```sql
CREATE TABLE products (
    product_id          UUID PRIMARY KEY,
    name                VARCHAR(255) NOT NULL,
    name_bn             VARCHAR(255),                 -- Bengali name
    category            VARCHAR(100) NOT NULL,        -- Health, Motor, Travel, Life …
    status              VARCHAR(50)  NOT NULL,        -- Draft, Active, Deactivated
    base_premium        DECIMAL(12,2) NOT NULL,
    sum_insured_min     DECIMAL(12,2),
    sum_insured_max     DECIMAL(12,2),
    tenure_months_min   INT,
    tenure_months_max   INT,
    description         TEXT,
    description_bn      TEXT,
    terms_url           VARCHAR(500),
    version             INT NOT NULL DEFAULT 1,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_status ON products(status);
```

#### `product_riders`
```sql
CREATE TABLE product_riders (
    rider_id        UUID PRIMARY KEY,
    product_id      UUID NOT NULL REFERENCES products(product_id),
    rider_name      VARCHAR(255) NOT NULL,
    premium_amount  DECIMAL(12,2) NOT NULL,
    coverage_amount DECIMAL(12,2) NOT NULL,
    is_mandatory    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

#### `risk_questions`
```sql
CREATE TABLE risk_questions (
    question_id     UUID PRIMARY KEY,
    product_id      UUID NOT NULL REFERENCES products(product_id),
    question_text   TEXT NOT NULL,
    question_text_bn TEXT,
    question_type   VARCHAR(50) NOT NULL,  -- boolean, multiple_choice, numeric
    options         JSONB,                 -- for multiple_choice
    sort_order      INT NOT NULL DEFAULT 0,
    is_required     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

#### `policies`
```sql
CREATE TABLE policies (
    policy_id           UUID PRIMARY KEY,
    policy_number       VARCHAR(50) UNIQUE NOT NULL,  -- LBT-YYYY-XXXX-NNNNNN
    product_id          UUID NOT NULL REFERENCES products(product_id),
    customer_id         UUID NOT NULL,
    partner_id          UUID,
    agent_id            UUID,
    status              VARCHAR(50) NOT NULL,          -- PolicyStatus enum
    premium_amount      DECIMAL(12,2) NOT NULL CHECK (premium_amount > 0),
    sum_insured         DECIMAL(12,2) NOT NULL,
    tenure_months       INT NOT NULL,
    start_date          DATE NOT NULL,
    end_date            DATE NOT NULL,
    issued_at           TIMESTAMPTZ,
    policy_document_url VARCHAR(1000),
    qr_code_data        TEXT,
    cooling_off_ends_at TIMESTAMPTZ,                  -- issued_at + 5 days (FR-038)
    grace_period_ends_at TIMESTAMPTZ,                 -- end_date + 30 days (FR-047)
    lapse_reinstate_deadline TIMESTAMPTZ,             -- lapse date + 90 days (FR-048)
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);  -- Partitioned by month (FR-239)

CREATE INDEX idx_policies_customer_id ON policies(customer_id);
CREATE INDEX idx_policies_status ON policies(status);
CREATE INDEX idx_policies_product_id ON policies(product_id);
CREATE INDEX idx_policies_end_date ON policies(end_date);  -- For renewal reminders
```

#### `policy_nominees`
```sql
CREATE TABLE policy_nominees (
    nominee_id      UUID PRIMARY KEY,
    policy_id       UUID NOT NULL REFERENCES policies(policy_id),
    full_name       VARCHAR(255) NOT NULL,
    relationship    VARCHAR(100) NOT NULL,
    share_percent   DECIMAL(5,2) NOT NULL CHECK (share_percent > 0 AND share_percent <= 100),
    date_of_birth   DATE,
    nid_number      VARCHAR(50),
    phone_number    VARCHAR(20),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_single_nominee CHECK (
        (SELECT COUNT(*) FROM policy_nominees pn WHERE pn.policy_id = policy_id) <= 1
    )  -- FR-032: single nominee
);
```

#### `policy_applicants`
```sql
CREATE TABLE policy_applicants (
    applicant_id            UUID PRIMARY KEY,
    policy_id               UUID NOT NULL REFERENCES policies(policy_id),
    full_name               VARCHAR(255) NOT NULL,
    date_of_birth           DATE NOT NULL,
    nid_number              VARCHAR(50),              -- Optional (FR-031)
    occupation              VARCHAR(255),
    annual_income           DECIMAL(12,2),
    address                 JSONB,                    -- Structured address
    health_declaration      JSONB,                    -- Pre-existing conditions, smoker, blood group
    risk_answers            JSONB,                    -- Risk assessment Q&A (FR-023-B)
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

#### `policy_lifecycle_events`
```sql
CREATE TABLE policy_lifecycle_events (
    event_id        UUID PRIMARY KEY,
    policy_id       UUID NOT NULL REFERENCES policies(policy_id),
    event_type      VARCHAR(100) NOT NULL, -- Issued, Renewed, Lapsed, Reinstated, Cancelled, Suspended, Expired
    from_status     VARCHAR(50),
    to_status       VARCHAR(50),
    performed_by    UUID,
    notes           TEXT,
    metadata        JSONB,
    occurred_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- Append-only; no UPDATE or DELETE (immutable audit — FR-050, FR-206)
CREATE INDEX idx_policy_events_policy_id ON policy_lifecycle_events(policy_id);
```

#### `policy_endorsements`
```sql
CREATE TABLE policy_endorsements (
    endorsement_id      UUID PRIMARY KEY,
    policy_id           UUID NOT NULL REFERENCES policies(policy_id),
    endorsement_number  VARCHAR(100) NOT NULL,  -- PLN-001/END-01
    endorsement_type    VARCHAR(100) NOT NULL,   -- Address, SumInsured, Nominee, Contact
    changes             JSONB NOT NULL,           -- Before/after snapshot
    additional_premium  DECIMAL(12,2) DEFAULT 0, -- For sum insured increase
    refund_amount       DECIMAL(12,2) DEFAULT 0, -- For sum insured decrease
    status              VARCHAR(50) NOT NULL,     -- Pending, Approved, Rejected
    approved_by         UUID,
    document_url        VARCHAR(1000),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    approved_at         TIMESTAMPTZ
);
```

#### `claims`
```sql
CREATE TABLE claims (
    claim_id            UUID PRIMARY KEY,
    claim_number        VARCHAR(50) UNIQUE NOT NULL,  -- CLM-YYYY-XXXX-NNNNNN
    policy_id           UUID NOT NULL REFERENCES policies(policy_id),
    customer_id         UUID NOT NULL,
    status              VARCHAR(50) NOT NULL,           -- ClaimStatus enum
    claim_type          VARCHAR(100) NOT NULL,          -- ClaimType enum
    claimed_amount      DECIMAL(12,2) NOT NULL,
    approved_amount     DECIMAL(12,2),
    settled_amount      DECIMAL(12,2),
    deductible_amount   DECIMAL(12,2) DEFAULT 0,
    copay_percentage    DECIMAL(5,2) DEFAULT 100,
    incident_date       DATE NOT NULL,
    incident_description TEXT NOT NULL,
    rejection_reason    TEXT,
    document_hash       VARCHAR(64),                   -- SHA-256 of submitted docs
    fraud_score         DECIMAL(5,2),                  -- 0–100
    fraud_flagged       BOOLEAN DEFAULT FALSE,
    submitted_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    approved_at         TIMESTAMPTZ,
    settled_at          TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);  -- Partitioned by month

CREATE INDEX idx_claims_policy_id ON claims(policy_id);
CREATE INDEX idx_claims_customer_id ON claims(customer_id);
CREATE INDEX idx_claims_status ON claims(status);
CREATE INDEX idx_claims_fraud_flagged ON claims(fraud_flagged) WHERE fraud_flagged = TRUE;
```

#### `claim_documents`
```sql
CREATE TABLE claim_documents (
    document_id     UUID PRIMARY KEY,
    claim_id        UUID NOT NULL REFERENCES claims(claim_id),
    document_type   VARCHAR(100) NOT NULL,  -- bill, prescription, police_report …
    file_url        VARCHAR(1000) NOT NULL,
    file_hash       VARCHAR(64) NOT NULL,   -- SHA-256 (FR-083)
    file_size_bytes INT,
    mime_type       VARCHAR(100),
    verified        BOOLEAN DEFAULT FALSE,
    verified_by     UUID,
    ocr_extracted   JSONB,                  -- OCR output (FR-087)
    uploaded_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    verified_at     TIMESTAMPTZ
);
```

#### `claim_approvals`
```sql
CREATE TABLE claim_approvals (
    approval_id     UUID PRIMARY KEY,
    claim_id        UUID NOT NULL REFERENCES claims(claim_id),
    approval_level  VARCHAR(20) NOT NULL,    -- L1, L2, L3, Board
    approver_id     UUID NOT NULL,
    approver_role   VARCHAR(100) NOT NULL,   -- Claims Officer, Manager, Business Admin, Focal Person, Board
    decision        VARCHAR(50) NOT NULL,    -- Pending, Approved, Rejected, NeedsMoreInfo
    approved_amount DECIMAL(12,2),
    notes           TEXT NOT NULL,           -- Mandatory (FR-090)
    decided_at      TIMESTAMPTZ,
    deadline_at     TIMESTAMPTZ NOT NULL,    -- TAT deadline (FR-097)
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

#### `fraud_checks`
```sql
CREATE TABLE fraud_checks (
    check_id        UUID PRIMARY KEY,
    claim_id        UUID NOT NULL REFERENCES claims(claim_id),
    rule_triggered  VARCHAR(100) NOT NULL,   -- e.g., FR-182, FR-183 …
    fraud_score     DECIMAL(5,2) NOT NULL,
    risk_factors    JSONB,
    flagged         BOOLEAN DEFAULT FALSE,
    reviewed_by     UUID,
    reviewed_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

#### `insured_assets` (polymorphic bridge table)
```sql
-- Bridge table for product-type-specific asset details
CREATE TABLE insured_assets (
    asset_id        UUID PRIMARY KEY,
    policy_id       UUID NOT NULL REFERENCES policies(policy_id),
    asset_type      VARCHAR(50) NOT NULL,  -- vehicle, health, travel, device …
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Motor / Vehicle
CREATE TABLE vehicle_details (
    asset_id            UUID PRIMARY KEY REFERENCES insured_assets(asset_id),
    registration_number VARCHAR(50) NOT NULL,
    make                VARCHAR(100),
    model               VARCHAR(100),
    year                INT,
    engine_number       VARCHAR(100),
    chassis_number      VARCHAR(100)
);

-- Health
CREATE TABLE health_declarations (
    asset_id                UUID PRIMARY KEY REFERENCES insured_assets(asset_id),
    has_pre_existing        BOOLEAN DEFAULT FALSE,
    conditions              JSONB,     -- Array of condition strings
    is_smoker               BOOLEAN DEFAULT FALSE,
    blood_group             VARCHAR(10),
    height_cm               DECIMAL(5,2),
    weight_kg               DECIMAL(5,2)
);

-- Travel
CREATE TABLE travel_details (
    asset_id            UUID PRIMARY KEY REFERENCES insured_assets(asset_id),
    destination_country VARCHAR(100) NOT NULL,
    travel_start_date   DATE NOT NULL,
    travel_end_date     DATE NOT NULL,
    traveller_count     INT DEFAULT 1
);
```

### 7.2 Data Retention Policy

| Data Type | Hot (PostgreSQL) | Warm (S3) | Cold (Glacier) | Total Retention |
|---|---|---|---|---|
| Active Policies | Lifetime | — | — | Policy lifetime |
| Expired Policies | 1 year | 5 years | 14 years | 20 years |
| Claims Data | Until settlement | After settlement | 20 years | 20 years |
| Audit Logs | 90 days | 1 year | 6 years | 7 years |
| Policy Documents (PDF) | S3 always (active) | S3 (expired) | Glacier | 20 years |

---

## 8. Kafka Event Topology

### 8.1 Events Published by Insurance Engine

#### Policy Events

| Topic | Event Message | Trigger |
|---|---|---|
| `policy.issued` | `PolicyIssuedEvent` | Policy activated after payment |
| `policy.renewed` | `PolicyRenewedEvent` | Manual or auto-renewal completed |
| `policy.cancelled` | `PolicyCancelledEvent` | Cancellation workflow approved |
| `policy.lapsed` | `PolicyLapsedEvent` | Grace period expired, no payment |
| `policy.reinstated` | `PolicyReinstatedEvent` | Lapsed policy reinstated |
| `policy.suspended` | `PolicySuspendedEvent` | Admin suspension |
| `policy.expired` | `PolicyExpiredEvent` | Tenure ended, not renewed |
| `policy.endorsed` | `PolicyEndorsedEvent` | Endorsement/amendment approved |
| `policy.grace_period_started` | `PolicyGracePeriodStartedEvent` | Expiry reached, entered grace |

#### Claim Events

| Topic | Event Message | Trigger |
|---|---|---|
| `claim.submitted` | `ClaimSubmittedEvent` | Customer submits claim |
| `claim.under_review` | `ClaimUnderReviewEvent` | Validation passed, under review |
| `claim.documents_requested` | `ClaimDocumentsRequestedEvent` | Documents incomplete |
| `claim.approved` | `ClaimApprovedEvent` | All required approvals received |
| `claim.rejected` | `ClaimRejectedEvent` | Claim rejected with reason |
| `claim.payment_initiated` | `ClaimPaymentInitiatedEvent` | Triggers Payment Service |
| `claim.settled` | `ClaimSettledEvent` | Payment confirmed settled |
| `claim.fraud_flagged` | `ClaimFraudFlaggedEvent` | Fraud rule triggered |

#### Product Events

| Topic | Event Message | Trigger |
|---|---|---|
| `product.created` | `ProductCreatedEvent` | Business Admin creates product |
| `product.updated` | `ProductUpdatedEvent` | Product details modified → Redis invalidated |
| `product.deactivated` | `ProductDeactivatedEvent` | Product deactivated |

### 8.2 Events Consumed by Insurance Engine

| Topic | Source Service | Purpose |
|---|---|---|
| `payment.confirmed` | Payment Service | Activate policy / settle claim |
| `payment.failed` | Payment Service | Update policy to PENDING_PAYMENT |
| `payment.refunded` | Payment Service | Update cancelled policy refund status |

### 8.3 Event Schema Example

```json
// PolicyIssuedEvent (published to topic: policy.issued)
{
  "event_id": "evt_abc123",
  "event_type": "PolicyIssued",
  "schema_version": "1.0",
  "occurred_at": "2025-03-01T10:00:00Z",
  "payload": {
    "policy_id": "pol_xyz789",
    "policy_number": "LBT-2025-0001-000123",
    "customer_id": "usr_cde456",
    "product_id": "prd_fgh012",
    "partner_id": null,
    "premium_amount": 1700.00,
    "sum_insured": 100000.00,
    "start_date": "2025-03-01",
    "end_date": "2026-03-01",
    "policy_document_url": "https://storage.labaid.com/policies/LBT-2025-0001-000123.pdf"
  }
}
```

---

## 9. gRPC Service Contracts

Defined in `proto/insuretech/policy/services/v1/insurance_service.proto`

```protobuf
syntax = "proto3";

package insuretech.policy.services.v1;

option csharp_namespace = "Insuretech.Policy.Services.V1";

service InsuranceEngineService {

  // ── Product Operations ──────────────────────────────────
  rpc CreateProduct       (CreateProductRequest)        returns (CreateProductResponse);
  rpc GetProduct          (GetProductRequest)           returns (GetProductResponse);
  rpc ListProducts        (ListProductsRequest)         returns (ListProductsResponse);
  rpc UpdateProduct       (UpdateProductRequest)        returns (UpdateProductResponse);
  rpc DeactivateProduct   (DeactivateProductRequest)    returns (DeactivateProductResponse);
  rpc CalculatePremium    (CalculatePremiumRequest)     returns (CalculatePremiumResponse);

  // ── Policy Operations ───────────────────────────────────
  rpc IssuePolicy         (IssuePolicyRequest)          returns (IssuePolicyResponse);
  rpc GetPolicy           (GetPolicyRequest)            returns (GetPolicyResponse);
  rpc ListPolicies        (ListPoliciesRequest)         returns (ListPoliciesResponse);
  rpc RenewPolicy         (RenewPolicyRequest)          returns (RenewPolicyResponse);
  rpc CancelPolicy        (CancelPolicyRequest)         returns (CancelPolicyResponse);
  rpc SuspendPolicy       (SuspendPolicyRequest)        returns (SuspendPolicyResponse);
  rpc ReinstatePolicy     (ReinstatePolicyRequest)      returns (ReinstatePolicyResponse);
  rpc EndorsePolicy       (EndorsePolicyRequest)        returns (EndorsePolicyResponse);
  rpc GetPolicyDocument   (GetPolicyDocumentRequest)    returns (GetPolicyDocumentResponse);

  // ── Claim Operations ────────────────────────────────────
  rpc SubmitClaim         (SubmitClaimRequest)          returns (SubmitClaimResponse);
  rpc GetClaim            (GetClaimRequest)             returns (GetClaimResponse);
  rpc ListClaims          (ListClaimsRequest)           returns (ListClaimsResponse);
  rpc ReviewClaim         (ReviewClaimRequest)          returns (ReviewClaimResponse);
  rpc ApproveClaim        (ApproveClaimRequest)         returns (ApproveClaimResponse);
  rpc RejectClaim         (RejectClaimRequest)          returns (RejectClaimResponse);
  rpc RequestClaimDocuments (RequestDocumentsRequest)   returns (RequestDocumentsResponse);
  rpc UploadClaimDocument (UploadClaimDocumentRequest)  returns (UploadClaimDocumentResponse);

  // ── Fraud / Risk ────────────────────────────────────────
  rpc RunFraudCheck       (FraudCheckRequest)           returns (FraudCheckResponse);
  rpc GetCustomerRiskScore(CustomerRiskScoreRequest)    returns (CustomerRiskScoreResponse);
}
```

---

## 10. Internal Integration Points

| Integrated Service | Protocol | Direction | Purpose |
|---|---|---|---|
| **API Gateway** (Go, 8080) | gRPC | Inbound | Receives all client requests |
| **Auth Service** (Go, 8081) | gRPC | Outbound (call) | Validate JWT on every command |
| **Authorization** (Go, 8082) | gRPC | Outbound (call) | RBAC permission check |
| **Payment Service** (Node.js, 3001) | Kafka (consume) | Inbound event | `payment.confirmed` → activate policy |
| **Payment Service** (Node.js, 3001) | Kafka (publish) | Outbound event | `claim.payment_initiated` → trigger payout |
| **Partner Management** (C# .NET, 5002) | gRPC | Outbound (call) | Validate partner_id, get commission rate |
| **AI Engine** (Python, 4001) | gRPC | Outbound (call) | Fraud scoring, document OCR validation |
| **Storage Service** (Go, 8084) | gRPC | Outbound (call) | Upload policy PDFs, claim documents to S3 |
| **Kafka Service** (Go, 8086) | Kafka | Outbound (publish) | All domain events → Notification, Analytics |
| **DBManager** (Go, 8083) | gRPC | Outbound (call) | Schema migrations coordination |
| **Analytics & Reporting** (C# .NET, 5003) | Kafka (consume) | Inbound event | Reporting consumes all Insurance Engine events |

---

## 11. FR ID Reference Index

Complete mapping of all Functional Requirements owned by the Insurance Engine:

| FR ID | Module | Short Description | Priority |
|---|---|---|---|
| FR-021 | FG-003 Product | Product catalog — M1 categories | M1 |
| FR-021-A | FG-003 Product | Product catalog — M2 extended categories | M2 |
| FR-022 | FG-003 Product | Product search (name, category, premium range) | M1 |
| FR-023 | FG-003 Product | Product detail display + PDF download | M2 |
| FR-023-A | FG-003 Product | Unit-wise plan + adjustable coverage amount | M2 |
| FR-023-B | FG-003 Product | Risk assessment questionnaire per plan | M2 |
| FR-024 | FG-003 Product | Dynamic premium calculator | M3 |
| FR-025 | FG-003 Product | Side-by-side product comparison (≤ 2) | M3 |
| FR-026 | FG-003 Product | Business Admin product CRUD + version history | M1 |
| FR-027 | FG-003 Product | Product variants with riders and add-ons | D |
| FR-028 | FG-003 Product | Redis product cache (5-min TTL) | M3 |
| FR-029 | FG-003 Product | Multi-language product descriptions (BN + EN) | M3 |
| FR-030 | FG-004 Policy | End-to-end purchase flow (< 10 min) | M1 |
| FR-031 | FG-004 Policy | Applicant info collection | M1 |
| FR-032 | FG-004 Policy | Single nominee required | M1 |
| FR-032-A | FG-004 Policy | Nominee income range optional | M1 |
| FR-033 | FG-004 Policy | NID/Mobile uniqueness, duplicate policy detection | M1 |
| FR-034 | FG-004 Policy | Policy number generation: LBT-YYYY-XXXX-NNNNNN | M1 |
| FR-035 | FG-004 Policy | Digital PDF policy with QR code (< 30 s) | M2 |
| FR-036 | FG-004 Policy | Policy document delivery: SMS + email | M2 |
| FR-037 | FG-004 Policy | Immediate activation on payment (Non-Life) | M2 |
| FR-038 | FG-004 Policy | Cooling-off period: 5 days, full refund | M3 |
| FR-039 | FG-004 Policy | Policy status management + transition logging | M1 |
| FR-040 | FG-004 Policy | Customer policy dashboard (< 3 s) | M1 |
| FR-041 | FG-004 Policy | Order history + referral (max 1) | D |
| FR-042 | FG-005 Renewal | Family Insurance Wallet (≤ 6 members) | D |
| FR-043 | FG-005 Renewal | Renewal reminders: 30d / 7d / 1d | M2 |
| FR-044 | FG-005 Renewal | Manual one-click renewal (< 3 min) | M2 |
| FR-045 | FG-005 Renewal | Auto-repurchase (opt-in, stored payment) | M3 |
| FR-046 | FG-005 Renewal | Editable fields during renewal | M3 |
| FR-047 | FG-005 Renewal | Grace period: 30 days, coverage continues | M2 |
| FR-048 | FG-005 Renewal | Auto-lapse after grace; reinstatement ≤ 90 days | M2 |
| FR-049 | FG-005 Renewal | Policy PDF download + version history | M1 |
| FR-050 | FG-005 Renewal | Lifecycle event audit trail | M1 |
| FR-051 | FG-005.1 Cancel | Cancellation workflow + reason dropdown | M1 |
| FR-052 | FG-005.1 Cancel | Approval workflow for policies > 30 days old | M1 |
| FR-053 | FG-005.1 Cancel | Pro-rata refund calculation | M1 |
| FR-054 | FG-005.1 Cancel | Refund within 7 working days | M1 |
| FR-055 | FG-005.1 Cancel | Status → CANCELLED + stakeholder notification | M1 |
| FR-056 | FG-005.2 Endorse | Endorsement types: Address/SumInsured/Nominee/Contact | M1 |
| FR-057 | FG-005.2 Endorse | Additional premium for sum insured increase | D |
| FR-058 | FG-005.2 Endorse | Pro-rata credit for sum insured decrease | M2 |
| FR-059 | FG-005.2 Endorse | Endorsement doc with suffix PLN-001/END-01 | M1 |
| FR-060 | FG-005.2 Endorse | Approval for sum insured change > 10% | M1 |
| FR-061 | FG-006 Rules | Premium calculation fallback chain | M1 |
| FR-062 | FG-006 Rules | Premium edge cases: age, occupation, conditions | M2 |
| FR-063 | FG-006 Rules | Duplicate policy detection (same product + 30 days) | M1 |
| FR-064 | FG-006 Rules | Policy merge workflow (Focal Person) | M3 |
| FR-065 | FG-006 Rules | Claim status state machine definition | M1 |
| FR-066 | FG-006 Rules | Claim state transition rules + approval routing | M1 |
| FR-067 | FG-006 Rules | Gamified renewal rewards | D |
| FR-068 | FG-006 Rules | Grace period logic (30 days) | M3 |
| FR-069 | FG-006 Rules | Lapsed policy reinstatement (90 days) | D |
| FR-081 | FG-008 Claims | Fixed-step claim form; claim tracker | M1 |
| FR-082 | FG-008 Claims | Claim eligibility validation (< 3 s) | M1 |
| FR-083 | FG-008 Claims | Claim number generation + SHA-256 hash | M1 |
| FR-084 | FG-008 Claims | Partner/insurer notification (< 60 s) | M2 |
| FR-085 | FG-008 Claims | Real-time status tracking + push/SMS | M3 |
| FR-086 | FG-008 Claims | Tiered approval per Claims Approval Matrix | M3 |
| FR-087 | FG-008 Claims | Document verification: OCR + quality check | M3 |
| FR-088 | FG-008 Claims | In-claim chat interface | M3 |
| FR-089 | FG-008 Claims | WebRTC video call for claim inspection | D |
| FR-090 | FG-008 Claims | Partner approval notes (mandatory reason) | M2 |
| FR-091 | FG-008 Claims | Joint approval for claims BDT 50K–2L | M3 |
| FR-092 | FG-008 Claims | Auto-payment within 24 hrs of approval | M3 |
| FR-093 | FG-008 Claims | ZHTC: auto-settle claims < BDT 10K | D |
| FR-094 | FG-008 Claims | Fraud detection rules (frequency, duplicates) | M3 |
| FR-095 | FG-008 Claims | Auto-revoke for confirmed fraud | M3 |
| FR-096 | FG-008 Claims | Balance sheet per stakeholder level | M3 |
| FR-097 | FG-008 Claims | TAT tracking + SLA breach alerts | M3 |
| FR-098 | FG-008 Claims | Claim history analytics | M3 |
| FR-099 | FG-008.1 Docs | Document requirements (PDF/JPG/PNG, 5 MB, 300 DPI) | M1 |
| FR-100 | FG-008.1 Docs | Co-payment + deductible calculation | M1 |
| FR-101 | FG-008.1 Docs | Claim reimbursement workflow (7–15 working days) | M1 |
| FR-175 | FG-016 Fraud | Flag claims submitted < 48 hrs of policy purchase | M2 |
| FR-176 | FG-016 Fraud | Flag same claim type > 2 in 12 months | M2 |
| FR-177 | FG-016 Fraud | Flag claims at 100% of coverage | M2 |
| FR-178 | FG-016 Fraud | Validate medical provider against approved network | M2 |
| FR-179 | FG-016 Fraud | Device fingerprinting: > 3 accounts from same device | M3 |
| FR-180 | FG-016 Fraud | Fraud dashboard (Business Admin + Focal Person) | M2 |
| FR-181 | FG-016 Fraud | RACI enforcement for monitoring + escalation | M1 |
| FR-182 | FG-016 Rules | Fraud rule: Rapid Policy-to-Claim < 48 hrs | M1 |
| FR-183 | FG-016 Rules | Fraud rule: Frequent Claims > 2 in 12 months | M1 |
| FR-184 | FG-016 Rules | Fraud rule: Amount exactly 100% of coverage | M1 |
| FR-185 | FG-016 Rules | Fraud rule: Non-network medical provider | M1 |
| FD-186 | FG-016 Rules | Fraud rule: Geographic anomaly > 100 km | M1 |
| FD-187 | FG-016 Rules | Fraud rule: > 3 accounts from same device | M1 |
| FD-188 | FG-016 Rules | Fraud rule: ML-based behavioral scoring | M1 |

---

*Document generated from SRS v3.11 (Feb 2026) · LabAid InsureTech Platform*
*Last updated: February 2026 · Service port: 5001 · Language: C# .NET 8*
