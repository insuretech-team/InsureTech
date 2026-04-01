# Insurance Engine — Complete Module Specification
## LabAid InsureTech Platform | SRS v3.11

> **Service:** Insurance Engine
> **Language:** C# .NET 8
> **Port:** 5001
> **Architecture:** Vertical Slice Architecture (VSA) + CQRS/MediatR
> **Database:** PostgreSQL 17 (EF Core)
> **Cache:** StackExchange.Redis
> **Messaging:** Apache Kafka (outbound domain events)
> **Inter-service:** gRPC + Protocol Buffers (Category 1 API — <100ms)
> **Resilience:** Polly (retry, circuit breaker)
> **Validation:** FluentValidation
> **Source of Truth:** SRS v3.11 | FINAL_DRAFT | Feb 2026

---

## Index

1. [Service Identity & Technical Foundation](#1-service-identity--technical-foundation)
2. [Module 01 — Product Catalog](#2-module-01--product-catalog-fg-003)
3. [Module 02 — Policy Lifecycle](#3-module-02--policy-lifecycle-fg-004)
4. [Module 03 — Policy Renewals](#4-module-03--policy-renewals-fg-005)
5. [Module 04 — Policy Cancellation & Refund](#5-module-04--policy-cancellation--refund-fg-0051)
6. [Module 05 — Policy Endorsement & Amendment](#6-module-05--policy-endorsement--amendment-fg-0052)
7. [Module 06 — Business Rules & Workflows](#7-module-06--business-rules--workflows-fg-006)
8. [Module 07 — Claims Management](#8-module-07--claims-management-fg-008)
9. [Module 08 — Claims Document Processing](#9-module-08--claims-document-processing-fg-0081)
10. [Module 09 — Fraud Detection & Risk Controls](#10-module-09--fraud-detection--risk-controls-fg-016)
11. [Module 10 — Payment Processing (Engine-side)](#11-module-10--payment-processing-engine-side-fg-007)
12. [Module 11 — Audit & Logging](#12-module-11--audit--logging-fg-019)
13. [Module 12 — Notifications (Outbound Events)](#13-module-12--notifications-outbound-events-fg-012)
14. [Proto / gRPC Service Definitions](#14-proto--grpc-service-definitions)
15. [PostgreSQL Database Schema](#15-postgresql-database-schema)
16. [Kafka Event Topology](#16-kafka-event-topology)
17. [State Machines](#17-state-machines)
18. [Claims Approval Matrix](#18-claims-approval-matrix)
19. [Fraud Detection Rules Reference](#19-fraud-detection-rules-reference)
20. [NFR — Performance & Compliance Constraints](#20-nfr--performance--compliance-constraints)
21. [Integration Dependencies](#21-integration-dependencies)
22. [FR → Acceptance Criteria Master Table](#22-fr--acceptance-criteria-master-table)

---

## 1. Service Identity & Technical Foundation

### 1.1 Microservice Profile

```
Service Name    : Insurance Engine
Language        : C# .NET 8
Port            : 5001
Protocol        : gRPC (internal) + Kafka (events)
Architecture    : VSA (Vertical Slice Architecture)
Pattern         : CQRS via MediatR
ORM             : EF Core → PostgreSQL 17
Cache           : StackExchange.Redis
Resilience      : Polly (exponential backoff, circuit breaker)
Validation      : FluentValidation (per-command/query)
Source of Truth : Proto3 definitions
```

### 1.2 Architectural Rules (VSA + CQRS)

- **Vertical Slice**: কোনো shared horizontal layer নেই। প্রতিটা feature তার নিজের handler, validator, mapper, repository নিয়ে একটা self-contained slice।
- **CQRS separation**:
  - **Command** → PostgreSQL write (strong consistency) → Kafka event publish
  - **Query** → Redis cache → Read replica (eventual consistency)
- **MediatR pipeline**: `Request → FluentValidation → Handler → Domain Event`
- **No cross-slice dependency**: Slice A, Slice B এর কোনো internal class use করতে পারবে না।
- **Polly policy**: সব external gRPC calls এ retry (exponential: 1s→2s→4s→8s→16s, max 5) + circuit breaker।
- **Idempotency**: সব policy issuance এবং payment API তে `Idempotency-Key` (UUID) header required। 24 ঘণ্টা key store করতে হবে।

### 1.3 Proto-First Strategy

```
Proto Definitions (canonical)
    ├── EF Core Entity → PostgreSQL table
    ├── gRPC Request/Response
    ├── Kafka event schema
    └── C# generated classes (gen/csharp/)
```

Proto namespace: `insuretech.policy.entity.v1`, `insuretech.claims.entity.v1`, `insuretech.payment.entity.v1`

### 1.4 External Dependencies

| Dependency | Purpose | Timeout |
|-----------|---------|---------|
| Auth Service (Go:8081) | JWT validation | 5s |
| Partner Management (C#:5002) | Commission calc, partner validation | 10s |
| Payment Service (Node.js:3001) | Payment initiation, status | 15s |
| AI Engine (Python:4001) | Fraud scoring, document OCR result | 30s |
| Storage Service (Go:8084) | Document URL, presigned S3 | 10s |
| Kafka Service (Go:8086) | Event publish | async |
| bKash/Nagad API | MFS payment gateway | 15s read, 5s connect |
| NID API | Identity verification | 10s |

---

## 2. Module 01 — Product Catalog (FG-003)

### 2.1 Scope

SRS Section: 4.3 | FR-021 থেকে FR-029

Insurance Engine এর Product Catalog module সমস্ত insurance product এর definition, categorization, search, pricing এবং lifecycle manage করে।

### 2.2 Functional Requirements

| FR ID | বিবরণ | Priority | Acceptance Criteria |
|-------|-------|----------|---------------------|
| FR-021 | Product catalog with categories: Health, Life, Motor, Travel, Micro-insurance। Display preference: Health → Auto → Travel → Life | M1 | Products by category, search + filter enabled |
| FR-022 | Product search by name, category, coverage type, premium range | M1 | Search results <500ms, fuzzy matching Bengali text |
| FR-023 | Product details: coverage, premium, tenure, exclusions, T&C | M2 | All info visible before purchase, PDF download |
| FR-023-A | Unit-wise plan purchase। User coverage amount increase/decrease করতে পারবে | M2 | Coverage amount adjustable |
| FR-023-B | প্রতিটা plan এ risk assessment questions থাকবে (multiple) | M2 | Every plan must have assessment questions |
| FR-024 | Premium calculator: dynamic inputs (age, sum assured, tenure, riders) | M3 | Real-time calc <2s, premium breakdown shown |
| FR-025 | Product comparison side-by-side (max 2 products) | M3 | Comparison table: features, coverage, pricing |
| FR-026 | Business Admin: product create, update, deactivate | M1 | CRUD operations, version history maintained |
| FR-027 | Product variants: configurable riders + add-ons। B2B product এ existing addon এর coverage increase করা যাবে | D | Base + optional riders, dynamic pricing recalc |
| FR-028 | Redis cache product catalog, 5-minute TTL | M3 | Cache hit rate >80%, auto-invalidation on update |
| FR-029 | Multi-language product description: Bengali + English | M3 | Language toggle, i18n format |

### 2.3 Product Categories & Display Order

```
Priority Display Order:
1. Health Insurance    (সর্বোচ্চ priority)
2. Auto/Motor Insurance
3. Travel Insurance
4. Life Insurance
5. Micro-Insurance
```

### 2.4 Product Data Model (Proto → DB)

```protobuf
// proto/insuretech/products/entity/v1/product.proto
message Product {
  string product_id = 1;           // UUID
  string product_code = 2;         // e.g. "HEALTH-001"
  string name_en = 3;
  string name_bn = 4;              // Bengali name
  ProductCategory category = 5;
  ProductStatus status = 6;
  double base_premium = 7;
  double min_sum_insured = 8;
  double max_sum_insured = 9;
  int32 min_tenure_months = 10;
  int32 max_tenure_months = 11;
  int32 min_age = 12;
  int32 max_age = 13;
  string description_en = 14;
  string description_bn = 15;
  string terms_url = 16;
  repeated RiskAssessmentQuestion questions = 17;
  repeated Rider available_riders = 18;
  repeated string exclusions = 19;
  google.protobuf.Timestamp effective_date = 20;
  google.protobuf.Timestamp created_at = 21;
  google.protobuf.Timestamp updated_at = 22;
  string created_by = 23;          // admin user_id
  int32 version = 24;
}

enum ProductCategory {
  PRODUCT_CATEGORY_UNSPECIFIED = 0;
  PRODUCT_CATEGORY_HEALTH = 1;
  PRODUCT_CATEGORY_AUTO = 2;
  PRODUCT_CATEGORY_TRAVEL = 3;
  PRODUCT_CATEGORY_LIFE = 4;
  PRODUCT_CATEGORY_MICRO = 5;
}

enum ProductStatus {
  PRODUCT_STATUS_UNSPECIFIED = 0;
  PRODUCT_STATUS_DRAFT = 1;
  PRODUCT_STATUS_ACTIVE = 2;
  PRODUCT_STATUS_INACTIVE = 3;
  PRODUCT_STATUS_ARCHIVED = 4;
}

message RiskAssessmentQuestion {
  string question_id = 1;
  string question_en = 2;
  string question_bn = 3;
  QuestionType type = 4;           // BOOLEAN, MULTIPLE_CHOICE, TEXT
  repeated string options = 5;
  bool is_required = 6;
}

message Rider {
  string rider_id = 1;
  string name_en = 2;
  string name_bn = 3;
  double additional_premium = 4;
  double additional_coverage = 5;
}
```

**PostgreSQL Table:**
```sql
CREATE TABLE products (
    product_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_code    VARCHAR(50) UNIQUE NOT NULL,
    name_en         VARCHAR(255) NOT NULL,
    name_bn         VARCHAR(255) NOT NULL,
    category        VARCHAR(50) NOT NULL,
    status          VARCHAR(50) NOT NULL DEFAULT 'DRAFT',
    base_premium    DECIMAL(12,2) NOT NULL CHECK (base_premium > 0),
    min_sum_insured DECIMAL(15,2) NOT NULL,
    max_sum_insured DECIMAL(15,2) NOT NULL,
    min_tenure_months INT NOT NULL,
    max_tenure_months INT NOT NULL,
    min_age         INT NOT NULL,
    max_age         INT NOT NULL,
    description_en  TEXT,
    description_bn  TEXT,
    terms_url       VARCHAR(500),
    questions       JSONB DEFAULT '[]',
    exclusions      JSONB DEFAULT '[]',
    effective_date  TIMESTAMP NOT NULL,
    version         INT NOT NULL DEFAULT 1,
    created_by      UUID NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE product_riders (
    rider_id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id            UUID NOT NULL REFERENCES products(product_id),
    name_en               VARCHAR(255) NOT NULL,
    name_bn               VARCHAR(255) NOT NULL,
    additional_premium    DECIMAL(12,2) NOT NULL,
    additional_coverage   DECIMAL(15,2) NOT NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_effective_date ON products(effective_date DESC);
-- Full-text search (Bengali + English)
CREATE INDEX idx_products_fts ON products USING GIN(
    to_tsvector('simple', name_en || ' ' || name_bn)
);
```

### 2.5 Cache Strategy

```
Redis Key Pattern: product:catalog:{category}
TTL: 300 seconds (5 minutes)
Invalidation: On any product CREATE/UPDATE/DEACTIVATE → delete all product:catalog:* keys
Premium calculator cache: product:premium:{product_id}:{age}:{sum_assured}:{tenure} → TTL 3600s
```

### 2.6 Business Rules

- Product version history maintain করতে হবে (audit purposes)।
- Deactivated product এ নতুন policy নেওয়া যাবে না, কিন্তু existing active policies চলবে।
- Sum insured range এর বাইরে purchase block।
- Age range এর বাইরে purchase block।
- Risk assessment questions: সব `is_required=true` questions এর উত্তর mandatory।

### 2.7 CQRS Slices

```
Commands:
  CreateProductCommand       → CreateProductHandler
  UpdateProductCommand       → UpdateProductHandler
  DeactivateProductCommand   → DeactivateProductHandler

Queries:
  GetProductByIdQuery        → GetProductByIdHandler (Redis → DB)
  ListProductsQuery          → ListProductsHandler (Redis → DB)
  SearchProductsQuery        → SearchProductsHandler (PostgreSQL FTS)
  CalculatePremiumQuery      → CalculatePremiumHandler (Redis → calculation)
  CompareProductsQuery       → CompareProductsHandler
```

---

## 3. Module 02 — Policy Lifecycle (FG-004)

### 3.1 Scope

SRS Section: 4.4 | FR-030 থেকে FR-041

Policy purchase থেকে issuance পর্যন্ত সম্পূর্ণ lifecycle। এটা Insurance Engine এর central module।

### 3.2 Functional Requirements

| FR ID | বিবরণ | Priority | Acceptance Criteria |
|-------|-------|----------|---------------------|
| FR-030 | End-to-end purchase flow: product selection → applicant details → nominee → payment → issuance | M1 | Complete flow <10 min, progress saved each step |
| FR-031 | Applicant info: full name, DOB, NID (optional), address, occupation, income, health declaration | M1 | Mandatory fields validated, conditional by product type |
| FR-032 | Single nominee/beneficiary support | M1 | Only 1 nominee required |
| FR-032-A | Beneficiary income range optional | M1 | Beneficiary submittable without income range |
| FR-033 | NID/Mobile uniqueness across policies — duplicate prevention | M1 | DB constraint enforced, user notified of existing policies |
| FR-034 | Policy number format: `LBT-YYYY-XXXX-NNNNNN` | M1 | Sequential, year-based prefix, collision-free |
| FR-035 | Digital policy document (PDF) with QR code | M2 | PDF within 30s of payment confirmation, QR scannable |
| FR-036 | Policy doc delivery via SMS link + email attachment | M2 | Delivery within 5 min, retry on failure |
| FR-037 | Non-Life policy activates immediately upon payment confirmation | M2 | Status real-time update, customer notified |
| FR-038 | Cooling-off period: 5 days from issuance for full refund | M3 | Cancellation within 24hr, refund initiated |
| FR-039 | Policy statuses: Pending Payment, Active, Suspended, Cancelled, Lapsed, Expired | M1 | Status transitions logged + timestamped, notifications triggered |
| FR-040 | Customer policy dashboard: active + past policies, renewal prompts, payment history | M1 | Load <3s, real-time status |
| FR-041 | Order history: Coverage details, Refer, Active Plans, Claimed Plans, Expired Plans। Max referral limit 1। Certificate download। | D | — |

### 3.3 Policy Purchase Flow (Step-by-Step)

```
Step 1: Product Selection
  → User selects product from catalog
  → System validates product is ACTIVE
  → System validates user age fits product range

Step 2: Risk Assessment
  → System presents all required questions for selected product (FR-023-B)
  → User answers all is_required=true questions
  → Responses stored for underwriting

Step 3: Applicant Details
  → full_name (required)
  → date_of_birth (required)
  → nid_number (optional)
  → mobile_number (required — uniqueness check)
  → occupation (required)
  → annual_income (required)
  → address (required)
  → health_declaration (conditional by product)

Step 4: Nominee Details
  → full_name (required)
  → relationship (required)
  → date_of_birth (required)
  → nid_number (optional)
  → phone_number (optional)
  → income_range (OPTIONAL — FR-032-A)
  → share_percentage = 100 (single nominee — FR-032)

Step 5: Coverage Configuration
  → base coverage selection
  → riders selection (optional)
  → sum insured adjustment (unit-wise — FR-023-A)
  → Premium recalculation

Step 6: Payment
  → Payment method selection (bKash/Nagad/Rocket/Manual)
  → Payment initiation → Payment Service (Node.js:3001)
  → Idempotency-Key header required

Step 7: Policy Issuance
  → Payment confirmed → status: ACTIVE
  → Policy number generated: LBT-YYYY-XXXX-NNNNNN
  → PDF generated (within 30s)
  → SMS + Email sent
  → Kafka event: PolicyIssued published
```

### 3.4 Policy Number Format

```
LBT - YYYY - XXXX - NNNNNN
 │     │       │       │
 │     │       │       └── Sequential 6-digit number (padded, collision-free)
 │     │       └────────── Product code (4-digit numeric ID)
 │     └────────────────── Year (4-digit)
 └──────────────────────── Company prefix "LBT" (LabAid InsureTech)

Example: LBT-2025-0012-000001
```

**Generation Strategy:**
- PostgreSQL sequence per `(year, product_code)` combination।
- Race condition prevent করতে `SELECT ... FOR UPDATE` বা advisory lock।
- Year rollover এ sequence reset।

### 3.5 Policy Data Model (Proto → DB)

```protobuf
// proto/insuretech/policy/entity/v1/policy.proto
message Policy {
  string policy_id = 1;                   // UUID
  string policy_number = 2;               // LBT-YYYY-XXXX-NNNNNN
  string product_id = 3;
  string customer_id = 4;
  string partner_id = 5;                  // optional
  string agent_id = 6;                    // optional
  PolicyStatus status = 7;
  double premium_amount = 8;
  double sum_insured = 9;
  int32 tenure_months = 10;
  google.protobuf.Timestamp start_date = 11;
  google.protobuf.Timestamp end_date = 12;
  google.protobuf.Timestamp issued_at = 13;
  google.protobuf.Timestamp created_at = 14;
  google.protobuf.Timestamp updated_at = 15;
  string policy_document_url = 16;
  repeated Nominee nominees = 17;
  repeated Rider selected_riders = 18;
  Applicant applicant = 19;
  repeated RiskAssessmentAnswer risk_answers = 20;
  string idempotency_key = 21;
}

enum PolicyStatus {
  POLICY_STATUS_UNSPECIFIED = 0;
  POLICY_STATUS_PENDING_PAYMENT = 1;
  POLICY_STATUS_ACTIVE = 2;
  POLICY_STATUS_GRACE_PERIOD = 3;
  POLICY_STATUS_LAPSED = 4;
  POLICY_STATUS_SUSPENDED = 5;
  POLICY_STATUS_CANCELLED = 6;
  POLICY_STATUS_EXPIRED = 7;
}

message Applicant {
  string full_name = 1;
  google.protobuf.Timestamp date_of_birth = 2;
  string nid_number = 3;                  // optional
  string occupation = 4;
  double annual_income = 5;
  string address = 6;
  HealthDeclaration health_declaration = 7;
}

message HealthDeclaration {
  bool has_pre_existing_conditions = 1;
  repeated string conditions = 2;
  bool is_smoker = 3;
  string blood_group = 4;
}

message Nominee {
  string nominee_id = 1;
  string full_name = 2;
  string relationship = 3;
  double share_percentage = 4;            // must be 100 (single nominee)
  google.protobuf.Timestamp date_of_birth = 5;
  string nid_number = 6;                 // optional
  string phone_number = 7;
  // income_range: intentionally omitted — FR-032-A (optional)
}

message RiskAssessmentAnswer {
  string question_id = 1;
  string answer = 2;
}
```

**PostgreSQL Tables:**

```sql
CREATE TABLE policies (
    policy_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_number       VARCHAR(30) UNIQUE NOT NULL,
    product_id          UUID NOT NULL REFERENCES products(product_id),
    customer_id         UUID NOT NULL,
    partner_id          UUID,
    agent_id            UUID,
    status              VARCHAR(50) NOT NULL DEFAULT 'PENDING_PAYMENT',
    premium_amount      DECIMAL(12,2) NOT NULL CHECK (premium_amount > 0),
    sum_insured         DECIMAL(15,2) NOT NULL CHECK (sum_insured > 0),
    tenure_months       INT NOT NULL CHECK (tenure_months > 0),
    start_date          DATE,
    end_date            DATE,
    issued_at           TIMESTAMP,
    policy_document_url VARCHAR(1000),
    idempotency_key     UUID UNIQUE,
    created_at          TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

-- Yearly partitions (FR-239)
CREATE TABLE policies_2025 PARTITION OF policies
    FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');
CREATE TABLE policies_2026 PARTITION OF policies
    FOR VALUES FROM ('2026-01-01') TO ('2027-01-01');

CREATE TABLE policy_applicants (
    applicant_id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_id           UUID NOT NULL REFERENCES policies(policy_id),
    full_name           VARCHAR(255) NOT NULL,
    date_of_birth       DATE NOT NULL,
    nid_number          VARCHAR(50),
    occupation          VARCHAR(255) NOT NULL,
    annual_income       DECIMAL(15,2) NOT NULL,
    address             TEXT NOT NULL,
    health_declaration  JSONB DEFAULT '{}',
    risk_answers        JSONB DEFAULT '[]',
    created_at          TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE policy_nominees (
    nominee_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_id           UUID NOT NULL REFERENCES policies(policy_id),
    full_name           VARCHAR(255) NOT NULL,
    relationship        VARCHAR(100) NOT NULL,
    share_percentage    DECIMAL(5,2) NOT NULL DEFAULT 100.00,
    date_of_birth       DATE,
    nid_number          VARCHAR(50),
    phone_number        VARCHAR(20),
    -- income_range intentionally excluded (FR-032-A)
    created_at          TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE policy_riders (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_id   UUID NOT NULL REFERENCES policies(policy_id),
    rider_id    UUID NOT NULL REFERENCES product_riders(rider_id),
    premium     DECIMAL(12,2) NOT NULL,
    coverage    DECIMAL(15,2) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE policy_status_history (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_id   UUID NOT NULL REFERENCES policies(policy_id),
    old_status  VARCHAR(50),
    new_status  VARCHAR(50) NOT NULL,
    reason      TEXT,
    changed_by  UUID,
    changed_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_policies_customer_id ON policies(customer_id);
CREATE INDEX idx_policies_status ON policies(status);
CREATE INDEX idx_policies_product_id ON policies(product_id);
CREATE INDEX idx_policies_end_date ON policies(end_date);
CREATE INDEX idx_policies_policy_number ON policies(policy_number);
```

### 3.6 Duplicate Policy Prevention (FR-033 & FR-063)

```
Rule: Same product + Same insured (NID or mobile) within 30 days → BLOCK
Cross-product: ALLOWED (different product type হলে চলবে)

Check sequence:
1. Extract customer's NID from applicant data
2. Query: SELECT COUNT(*) FROM policies p
          JOIN policy_applicants pa ON p.policy_id = pa.policy_id
          WHERE p.product_id = :product_id
            AND (pa.nid_number = :nid OR p.customer_id = :customer_id)
            AND p.created_at > NOW() - INTERVAL '30 days'
            AND p.status NOT IN ('CANCELLED', 'EXPIRED')
3. If COUNT > 0 → reject with clear error message
```

### 3.7 CQRS Slices

```
Commands:
  InitiatePolicyPurchaseCommand  → InitiatePolicyPurchaseHandler
  ConfirmPolicyPaymentCommand    → ConfirmPolicyPaymentHandler
  IssuePolicyCommand             → IssuePolicyHandler
  SuspendPolicyCommand           → SuspendPolicyHandler

Queries:
  GetPolicyByIdQuery             → GetPolicyByIdHandler
  GetPolicyByNumberQuery         → GetPolicyByNumberHandler
  ListCustomerPoliciesQuery      → ListCustomerPoliciesHandler
  GetPolicyDashboardQuery        → GetPolicyDashboardHandler
```

---

## 4. Module 03 — Policy Renewals (FG-005)

### 4.1 Scope

SRS Section: 4.5 | FR-042 থেকে FR-050

Policy expiry এর আগে ও পরে renewal management।

### 4.2 Functional Requirements

| FR ID | বিবরণ | Priority | Acceptance Criteria |
|-------|-------|----------|---------------------|
| FR-042 | Family Insurance Wallet: একাউন্টের অধীনে পরিবারের সব policies একসাথে manage | D | Unified dashboard, bulk payment, relationship mgmt |
| FR-043 | Renewal reminders: 30 days, 7 days, 1 day before expiry — SMS + email + push | M2 | On schedule, delivery tracked |
| FR-044 | Manual renewal: one-click, existing policy data reuse | M2 | <3 min, updated policy doc issued |
| FR-045 | Auto repurchase: stored payment method, opt-in | M3 | Consent recorded, auto-charge 7 days before expiry |
| FR-046 | Address + nominee update allowed during renewal | M3 | Limited editable fields, major change → verification |
| FR-047 | Grace period: 30 days post-expiry, coverage continues | M2 | Status "GRACE_PERIOD", daily reminders |
| FR-048 | Auto-lapse after grace period, reinstatement within 90 days with penalty | M2 | Status "LAPSED", reinstatement workflow |
| FR-049 | PDF download with version history for all renewals | M1 | All versions accessible, issue date marked |
| FR-050 | Policy lifecycle event audit trail: issuance, renewal, lapse, reinstatement, cancellation | M1 | Immutable event log, queryable |

### 4.3 Renewal Business Rules

```
Renewal Timeline:
  T-30 days  → Reminder SMS + Email + Push (FR-043)
  T-7 days   → Reminder SMS + Email + Push
  T-1 day    → Final reminder SMS + Email + Push
  T+0        → Policy EXPIRES (status → EXPIRED)
  T+1 to T+30 → GRACE_PERIOD (coverage continues, daily reminders)
  T+31       → Auto-LAPSE if payment not received

Auto-repurchase (FR-045):
  T-7 days   → Charge stored payment method
  On success → New policy issued, linked to old policy
  On failure → Notify customer, manual intervention required

Reinstatement (FR-048):
  Window: 90 days from lapse date
  Requires: Outstanding premium + reinstatement penalty (configurable)
  Medical underwriting: May be required per product type
  Approval: Focal Person approval required
```

### 4.4 Grace Period Logic (FR-047 + FR-048 + FR-068)

```
State: EXPIRED → GRACE_PERIOD
  Trigger: end_date passed + payment not received
  Duration: 30 days
  Coverage: CONTINUES during grace period
  Reminders: Daily SMS/push

State: GRACE_PERIOD → LAPSED
  Trigger: grace_period_end passed + still no payment
  Auto-process: Background job runs daily at 2:00 AM BST
  Kafka event: PolicyLapsed published

State: LAPSED → ACTIVE (Reinstatement)
  Window: 90 days from lapse_date
  Process: Customer pays outstanding + penalty → Focal Person approves
  After 90 days: Cannot reinstate, must buy new policy
```

**DB columns (policies table extension):**
```sql
ALTER TABLE policies ADD COLUMN grace_period_end DATE;
ALTER TABLE policies ADD COLUMN lapsed_at TIMESTAMP;
ALTER TABLE policies ADD COLUMN reinstatement_fee DECIMAL(12,2);
ALTER TABLE policies ADD COLUMN renewal_of_policy_id UUID REFERENCES policies(policy_id);
ALTER TABLE policies ADD COLUMN auto_renew_enabled BOOLEAN DEFAULT FALSE;
ALTER TABLE policies ADD COLUMN auto_renew_payment_token VARCHAR(500);
```

### 4.5 CQRS Slices

```
Commands:
  InitiateRenewalCommand           → InitiateRenewalHandler
  ConfirmRenewalPaymentCommand     → ConfirmRenewalPaymentHandler
  EnableAutoRenewalCommand         → EnableAutoRenewalHandler
  DisableAutoRenewalCommand        → DisableAutoRenewalHandler
  ReinstatePolicyCommand           → ReinstatePolicyHandler
  ProcessGracePeriodExpiryCommand  → ProcessGracePeriodExpiryHandler (scheduled)
  ProcessAutoLapseCommand          → ProcessAutoLapseHandler (scheduled)

Queries:
  GetRenewalEligibilityQuery       → GetRenewalEligibilityHandler
  ListUpcomingRenewalsQuery        → ListUpcomingRenewalsHandler
```

---

## 5. Module 04 — Policy Cancellation & Refund (FG-005.1)

### 5.1 Scope

SRS Section: 4.5.1 | FR-051 থেকে FR-055

### 5.2 Functional Requirements

| FR ID | বিবরণ | Priority | Acceptance Criteria |
|-------|-------|----------|---------------------|
| FR-051 | Cancellation workflow: customer/agent/admin request + reason | M1 | Request form with reason dropdown, attachment support |
| FR-052 | Approval workflow: policies >30 days old → Business Admin + Focal Person approval, 48hr SLA | M1 | Approval routing, 48hr SLA enforced |
| FR-053 | Pro-rata refund: `(Premium Paid - Days Covered - Admin Fee - Cancellation Charge)` | M1 | Refund calculator, configurable fees |
| FR-054 | Refund within 7 working days via MFS or bank transfer | M1 | Payment gateway integration, notification |
| FR-055 | Status → CANCELLED, all stakeholders notified, IDRA reporting | M1 | Multi-channel notification |

### 5.3 Cancellation Business Rules

```
Policy age ≤ 5 days (Cooling-off period — FR-038):
  → Full refund, no admin fee, no cancellation charge
  → No approval required
  → Automatic

Policy age ≤ 30 days:
  → Pro-rata refund calculation
  → No dual approval required
  → Business Admin or Agent can approve

Policy age > 30 days:
  → Dual approval: Business Admin + Focal Person (both must approve)
  → 48hr SLA from submission
  → Escalation if not approved within 48hr
```

### 5.4 Pro-Rata Refund Formula

```
Variables:
  P   = Total Premium Paid (BDT)
  D   = Days covered so far
  T   = Total tenure in days
  AF  = Admin Fee (configurable, e.g. 5% of P)
  CC  = Cancellation Charge (configurable, e.g. 2% of P)

Formula:
  Premium used = (P / T) × D
  Refund = P - Premium_used - AF - CC
  Refund = P - ((P/T) × D) - AF - CC

Example:
  P = 1200 BDT, T = 365 days, D = 60 days, AF = 60 BDT (5%), CC = 24 BDT (2%)
  Premium_used = (1200/365) × 60 = 197.26 BDT
  Refund = 1200 - 197.26 - 60 - 24 = 918.74 BDT
```

**DB:**
```sql
CREATE TABLE policy_cancellations (
    cancellation_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_id           UUID NOT NULL REFERENCES policies(policy_id),
    requested_by        UUID NOT NULL,
    requester_role      VARCHAR(50) NOT NULL,
    reason              VARCHAR(255) NOT NULL,
    reason_detail       TEXT,
    attachment_url      VARCHAR(1000),
    refund_amount       DECIMAL(12,2),
    admin_fee           DECIMAL(12,2),
    cancellation_charge DECIMAL(12,2),
    days_covered        INT,
    status              VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    -- approval fields
    bizadmin_approved_by  UUID,
    bizadmin_approved_at  TIMESTAMP,
    focal_approved_by     UUID,
    focal_approved_at     TIMESTAMP,
    rejection_reason      TEXT,
    -- refund
    refund_initiated_at   TIMESTAMP,
    refund_completed_at   TIMESTAMP,
    refund_payment_ref    VARCHAR(255),
    created_at          TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 5.5 CQRS Slices

```
Commands:
  RequestCancellationCommand      → RequestCancellationHandler
  ApproveCancellationCommand      → ApproveCancellationHandler
  RejectCancellationCommand       → RejectCancellationHandler
  ProcessCancellationRefundCommand → ProcessCancellationRefundHandler

Queries:
  GetCancellationStatusQuery      → GetCancellationStatusHandler
  ListPendingCancellationsQuery   → ListPendingCancellationsHandler
  CalculateCancellationRefundQuery → CalculateCancellationRefundHandler
```

---

## 6. Module 05 — Policy Endorsement & Amendment (FG-005.2)

### 6.1 Scope

SRS Section: 4.5.2 | FR-056 থেকে FR-060

Mid-term policy changes।

### 6.2 Functional Requirements

| FR ID | বিবরণ | Priority | Acceptance Criteria |
|-------|-------|----------|---------------------|
| FR-056 | Endorsement support: Address, Sum Insured, Nominee, Contact changes | M1 | Amendment forms, validation |
| FR-057 | Additional premium for mid-term sum insured increase | D | Premium calculator, payment integration |
| FR-058 | Pro-rata refund for sum insured decrease | M2 | Credit to premium account |
| FR-059 | Endorsement document: suffix format `PLN-001/END-01` | M1 | PDF generation, version tracking |
| FR-060 | Approval required for sum insured changes >10% | M1 | Approval workflow, threshold configurable |

### 6.3 Endorsement Types

```
Type 1: Address Change
  → No approval required
  → No premium impact
  → Endorsement doc generated

Type 2: Contact Change (phone/email)
  → No approval required
  → No premium impact

Type 3: Nominee Change
  → No approval required (single nominee, FR-032)
  → New nominee replaces old
  → Endorsement doc generated

Type 4: Sum Insured Change (Increase)
  → Change ≤ 10%: Auto-approved
  → Change > 10%: Business Admin approval required (FR-060)
  → Additional premium calculated (pro-rata remaining days)
  → Payment required before endorsement effective

Type 5: Sum Insured Change (Decrease)
  → Pro-rata refund for overpaid premium (FR-058)
  → Credit issued, no payment required
  → Endorsement doc generated
```

### 6.4 Endorsement Document Numbering

```
Format: {policy_number}/END-{sequence}
Example: LBT-2025-0012-000001/END-01
         LBT-2025-0012-000001/END-02

Sequence: per-policy, starting at 01, incrementing
```

**DB:**
```sql
CREATE TABLE policy_endorsements (
    endorsement_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_id           UUID NOT NULL REFERENCES policies(policy_id),
    endorsement_number  VARCHAR(50) UNIQUE NOT NULL,  -- PLN-001/END-01
    endorsement_type    VARCHAR(50) NOT NULL,
    change_description  TEXT NOT NULL,
    old_values          JSONB NOT NULL,
    new_values          JSONB NOT NULL,
    additional_premium  DECIMAL(12,2) DEFAULT 0,
    refund_amount       DECIMAL(12,2) DEFAULT 0,
    requires_approval   BOOLEAN DEFAULT FALSE,
    approved_by         UUID,
    approved_at         TIMESTAMP,
    rejection_reason    TEXT,
    status              VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    endorsement_doc_url VARCHAR(1000),
    effective_date      DATE,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 6.5 CQRS Slices

```
Commands:
  CreateEndorsementCommand      → CreateEndorsementHandler
  ApproveEndorsementCommand     → ApproveEndorsementHandler
  RejectEndorsementCommand      → RejectEndorsementHandler

Queries:
  GetEndorsementByIdQuery       → GetEndorsementByIdHandler
  ListPolicyEndorsementsQuery   → ListPolicyEndorsementsHandler
  CalculateEndorsementPremiumQuery → CalculateEndorsementPremiumHandler
```

---

## 7. Module 06 — Business Rules & Workflows (FG-006)

### 7.1 Scope

SRS Section: 4.6 | FR-061 থেকে FR-069

Cross-cutting business rules যা পুরো Insurance Engine জুড়ে apply হয়।

### 7.2 Functional Requirements

| FR ID | বিবরণ | Priority | Acceptance Criteria |
|-------|-------|----------|---------------------|
| FR-061 | Premium calculation fallback: Insurer API fail → cached rates (max 24hr) → queue + notify customer 2hr | M1 | Fallback tested, cache validation, queue notification |
| FR-062 | Premium edge cases: age-based loading, occupation risk, pre-existing conditions — clear messaging | M2 | All edge cases handled, actuarial validation |
| FR-063 | Duplicate policy detection: same product + person within 30 days → block; cross-product → allow | M1 | Accurate detection, clear error |
| FR-064 | Policy merge workflow: Focal Person merges duplicate accounts (NID verified), transfers policies | M3 | Data integrity, audit logged |
| FR-065 | Claim state machine: `Submitted → Under Review → Docs Requested → Approved/Rejected → Payment Initiated → Settled/Closed` | M1 | Invalid transitions blocked |
| FR-066 | Auto-transition to "Docs Requested" if incomplete; >BDT 50K → BizAdmin + Focal approval | M1 | Routing correct, notifications sent |
| FR-067 | Gamified renewal rewards: discounts/vouchers for early renewal | D | Points engine, voucher integration |
| FR-068 | Grace period logic (30 days) → auto-lapse → 90-day reinstatement | M3 | Enforced, customer notified |
| FR-069 | Lapsed policy reinstatement: within 90 days, medical underwriting, Focal approval | D | Workflow, approval required |

### 7.3 Premium Calculation with Fallback (FR-061)

```
Primary Flow:
  1. Call Insurer gRPC API for live rate
  2. On success → calculate premium → cache result (24hr TTL)
  3. Return premium to user

Fallback Flow:
  1. Insurer API fails / timeout
  2. Check Redis: product:rates:{product_id}:{age_band} → max 24hr old
  3. Cache hit → use cached rate → flag response as "cached_rate"
  4. Cache miss → queue calculation request to Kafka
     → Notify customer: "Quote being prepared, you'll be notified within 2 hours"
  5. Background worker resolves when API recovers → notify customer

Polly Config:
  Retry: 5 attempts, exponential: 1s, 2s, 4s, 8s, 16s
  Circuit Breaker: Open after 3 consecutive failures, 30s window
```

### 7.4 Claim State Machine (FR-065 + FR-066)

```
                    ┌─────────────┐
                    │  SUBMITTED  │
                    └──────┬──────┘
                           │ Auto-validate eligibility
                    ┌──────▼──────────────┐
                    │    UNDER_REVIEW     │
                    └──────┬──────────────┘
              ┌────────────┼────────────────┐
              │ Incomplete │                │ Complete docs
     ┌────────▼──────┐     │       ┌────────▼────────┐
     │ DOCS_REQUESTED│     │       │    APPROVED     │
     └────────┬──────┘     │       └────────┬────────┘
              │ Docs        │                │
              │ received    │       ┌────────▼────────┐
              └────────────►│       │PAYMENT_INITIATED │
                            │       └────────┬────────┘
                    ┌───────▼───┐            │
                    │ REJECTED  │   ┌────────▼────────┐
                    └───────────┘   │    SETTLED      │
                                    └─────────────────┘

Invalid transitions are BLOCKED at the application layer.
Every transition is logged in claim_status_history.
```

### 7.5 Duplicate Policy Detection Rules

```sql
-- FR-033 + FR-063: Check before any new policy issuance
SELECT COUNT(*)
FROM policies p
JOIN policy_applicants pa ON p.policy_id = pa.policy_id
WHERE p.product_id = :new_product_id
  AND (
    pa.nid_number = :applicant_nid
    OR p.customer_id = :customer_id
  )
  AND p.created_at > NOW() - INTERVAL '30 days'
  AND p.status NOT IN ('CANCELLED', 'EXPIRED', 'LAPSED');

-- If COUNT > 0 → BLOCK with error: "Duplicate policy detected"
-- Cross-product (different product_id) → ALLOW
```

---

## 8. Module 07 — Claims Management (FG-008)

### 8.1 Scope

SRS Section: 4.8 | FR-081 থেকে FR-098

Claim submission থেকে settlement পর্যন্ত সম্পূর্ণ workflow।

### 8.2 Functional Requirements

| FR ID | বিবরণ | Priority | Acceptance Criteria |
|-------|-------|----------|---------------------|
| FR-081 | Claim form: policy selection → incident details → claim reason → doc upload। Claim tracker must show। | M1 | Form <5 min, draft saving each step |
| FR-082 | Eligibility validation: policy active, within coverage, claim type covered, no duplicate | M1 | Validation <3s, clear errors |
| FR-083 | Claim number: `CLM-YYYY-XXXX-NNNNNN` + SHA-256 hash | M1 | Collision-free, integrity hash |
| FR-084 | Auto-notify partner/insurer on submission, shared dashboard | M2 | Notification <60s |
| FR-085 | Real-time status: Submitted, Under Review, Approved, Rejected, Settled। Push + SMS | M3 | Updates <5s |
| FR-086 | Tiered approval per Approval Matrix | M3 | Auto-routing, escalation on timeout |
| FR-087 | Document verification: image quality, OCR extraction, fraud detection | M3 | Validation <10s, OCR >85% |
| FR-088 | Chat between customer, partner agent, focal person | M3 | Real-time, file attachment, history |
| FR-089 | WebRTC video call for claim verification | D | HD, screen share, recording |
| FR-090 | Partner: add notes, approve/reject with mandatory reason | M2 | Timestamped, reason required |
| FR-091 | Joint approval BizAdmin + Focal for BDT 50K–2L | M3 | Both required, 5-day timeout escalation |
| FR-092 | Auto-payment on approval via customer's selected channel | M3 | Within 24hr, confirmation sent |
| FR-093 | Zero Human Touch Claims (ZHTC): <BDT 10K with partner pre-agreement | D | 95% automation, ML fraud check, instant |
| FR-094 | Fraud detection: frequent claims, duplicate docs, rapid policy-to-claim | M3 | Auto-flag, risk score, review queue |
| FR-095 | Auto-revoke access for confirmed fraud | M3 | Account suspension, appeal process |
| FR-096 | Balance sheet: Customer, Partner, Agent, InsureTech level | M3 | Daily/monthly/quarterly, export |
| FR-097 | TAT tracking per approval level, SLA breach alert | M3 | Real-time monitoring, email alert |
| FR-098 | Claim history + analytics for risk assessment | M3 | Frequency report, avg amount, settlement ratio |

### 8.3 Claim Number Format

```
CLM - YYYY - XXXX - NNNNNN
 │     │       │       │
 │     │       │       └── Sequential 6-digit (padded)
 │     │       └────────── Product code (4-digit)
 │     └────────────────── Year
 └──────────────────────── "CLM" prefix

Example: CLM-2025-0012-000045

Document Integrity Hash:
  SHA-256 of: claim_id + customer_id + submitted_at + total_claimed_amount
  Stored: claim.submission_hash
  Verified: On every claim access
```

### 8.4 Claim Data Model (Proto → DB)

```protobuf
// proto/insuretech/claims/entity/v1/claim.proto
message Claim {
  string claim_id = 1;
  string claim_number = 2;
  string policy_id = 3;
  string customer_id = 4;
  ClaimStatus status = 5;
  ClaimType type = 6;
  double claimed_amount = 7;
  double approved_amount = 8;
  double settled_amount = 9;
  google.protobuf.Timestamp incident_date = 10;
  string incident_description = 11;
  repeated ClaimDocument documents = 12;
  repeated ClaimApproval approvals = 13;
  google.protobuf.Timestamp submitted_at = 14;
  google.protobuf.Timestamp approved_at = 15;
  google.protobuf.Timestamp settled_at = 16;
  google.protobuf.Timestamp created_at = 17;
  google.protobuf.Timestamp updated_at = 18;
  string rejection_reason = 19;
  FraudCheckResult fraud_check = 20;
  string submission_hash = 21;           // SHA-256
  string payment_channel = 22;           // bKash/Nagad/Rocket/bank
  double deductible_amount = 23;
  double copayment_percentage = 24;
  double final_payable_amount = 25;
}

enum ClaimStatus {
  CLAIM_STATUS_UNSPECIFIED = 0;
  CLAIM_STATUS_SUBMITTED = 1;
  CLAIM_STATUS_UNDER_REVIEW = 2;
  CLAIM_STATUS_PENDING_DOCUMENTS = 3;
  CLAIM_STATUS_APPROVED = 4;
  CLAIM_STATUS_REJECTED = 5;
  CLAIM_STATUS_PAYMENT_INITIATED = 6;
  CLAIM_STATUS_SETTLED = 7;
  CLAIM_STATUS_DISPUTED = 8;
  CLAIM_STATUS_CLOSED = 9;
}

enum ClaimType {
  CLAIM_TYPE_UNSPECIFIED = 0;
  CLAIM_TYPE_HEALTH_HOSPITALIZATION = 1;
  CLAIM_TYPE_HEALTH_SURGERY = 2;
  CLAIM_TYPE_MOTOR_ACCIDENT = 3;
  CLAIM_TYPE_MOTOR_THEFT = 4;
  CLAIM_TYPE_TRAVEL_MEDICAL = 5;
  CLAIM_TYPE_TRAVEL_BAGGAGE_LOSS = 6;
  CLAIM_TYPE_DEVICE_DAMAGE = 7;
  CLAIM_TYPE_DEVICE_THEFT = 8;
  CLAIM_TYPE_DEATH = 9;
}

message ClaimDocument {
  string document_id = 1;
  string document_type = 2;
  string file_url = 3;
  string file_hash = 4;                  // SHA-256
  int64 file_size_bytes = 5;
  string mime_type = 6;
  google.protobuf.Timestamp uploaded_at = 7;
  bool verified = 8;
  string verified_by = 9;
  string ocr_extracted_text = 10;
}

message ClaimApproval {
  string approval_id = 1;
  string approver_id = 2;
  string approver_role = 3;
  ApprovalDecision decision = 4;
  double approved_amount = 5;
  string notes = 6;
  google.protobuf.Timestamp decided_at = 7;
  int32 approval_level = 8;              // L1, L2, L3, Board
}

enum ApprovalDecision {
  APPROVAL_DECISION_UNSPECIFIED = 0;
  APPROVAL_DECISION_PENDING = 1;
  APPROVAL_DECISION_APPROVED = 2;
  APPROVAL_DECISION_REJECTED = 3;
  APPROVAL_DECISION_NEEDS_MORE_INFO = 4;
}

message FraudCheckResult {
  double fraud_score = 1;                // 0–100
  repeated string risk_factors = 2;
  bool flagged = 3;
  string reviewed_by = 4;
  google.protobuf.Timestamp reviewed_at = 5;
}
```

**PostgreSQL Tables:**

```sql
CREATE TABLE claims (
    claim_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    claim_number        VARCHAR(30) UNIQUE NOT NULL,
    policy_id           UUID NOT NULL REFERENCES policies(policy_id),
    customer_id         UUID NOT NULL,
    status              VARCHAR(50) NOT NULL DEFAULT 'SUBMITTED',
    claim_type          VARCHAR(100) NOT NULL,
    claimed_amount      DECIMAL(12,2) NOT NULL CHECK (claimed_amount > 0),
    approved_amount     DECIMAL(12,2),
    settled_amount      DECIMAL(12,2),
    incident_date       DATE NOT NULL,
    incident_description TEXT NOT NULL,
    submission_hash     VARCHAR(64) NOT NULL,
    payment_channel     VARCHAR(50),
    deductible_amount   DECIMAL(12,2) DEFAULT 0,
    copayment_pct       DECIMAL(5,2) DEFAULT 0,
    final_payable       DECIMAL(12,2),
    rejection_reason    TEXT,
    submitted_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    approved_at         TIMESTAMP,
    settled_at          TIMESTAMP,
    created_at          TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (submitted_at);

CREATE TABLE claim_documents (
    document_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    claim_id            UUID NOT NULL REFERENCES claims(claim_id),
    document_type       VARCHAR(100) NOT NULL,
    file_url            VARCHAR(1000) NOT NULL,
    file_hash           VARCHAR(64) NOT NULL,
    file_size_bytes     BIGINT NOT NULL,
    mime_type           VARCHAR(100) NOT NULL,
    verified            BOOLEAN DEFAULT FALSE,
    verified_by         UUID,
    ocr_extracted_text  TEXT,
    uploaded_at         TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE claim_approvals (
    approval_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    claim_id            UUID NOT NULL REFERENCES claims(claim_id),
    approver_id         UUID NOT NULL,
    approver_role       VARCHAR(100) NOT NULL,
    approval_level      INT NOT NULL,
    decision            VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    approved_amount     DECIMAL(12,2),
    notes               TEXT,
    decided_at          TIMESTAMP,
    due_at              TIMESTAMP,
    created_at          TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE fraud_checks (
    check_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    claim_id            UUID NOT NULL REFERENCES claims(claim_id),
    fraud_score         DECIMAL(5,2) NOT NULL,
    risk_factors        JSONB DEFAULT '[]',
    flagged             BOOLEAN DEFAULT FALSE,
    reviewed_by         UUID,
    reviewed_at         TIMESTAMP,
    created_at          TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE claim_status_history (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    claim_id    UUID NOT NULL REFERENCES claims(claim_id),
    old_status  VARCHAR(50),
    new_status  VARCHAR(50) NOT NULL,
    reason      TEXT,
    changed_by  UUID,
    changed_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_claims_policy_id ON claims(policy_id);
CREATE INDEX idx_claims_customer_id ON claims(customer_id);
CREATE INDEX idx_claims_status ON claims(status);
CREATE INDEX idx_claims_submitted_at ON claims(submitted_at DESC);
CREATE INDEX idx_claims_fraud_flagged ON fraud_checks(flagged) WHERE flagged = TRUE;
```

### 8.5 Co-payment & Deductible Calculation (FR-100)

```
Formula:
  Net Payable = (Claimed Amount - Annual Deductible Remaining) × (1 - Copayment%)

Variables:
  Claimed Amount      : Amount customer is claiming
  Annual Deductible   : Fixed annual deductible from product config
  Deductible Remaining: Annual deductible - already used this year
  Copayment %         : Customer's share percentage (from product config)

Example:
  Claimed = 50,000 BDT
  Annual Deductible = 5,000 BDT
  Already used deductible = 2,000 BDT
  Remaining deductible = 3,000 BDT
  Copayment = 10%

  After deductible = 50,000 - 3,000 = 47,000 BDT
  Copayment amount = 47,000 × 10% = 4,700 BDT (customer pays)
  Net Payable = 47,000 × 90% = 42,300 BDT
```

### 8.6 CQRS Slices

```
Commands:
  SubmitClaimCommand              → SubmitClaimHandler
  UploadClaimDocumentCommand      → UploadClaimDocumentHandler
  RequestAdditionalDocsCommand    → RequestAdditionalDocsHandler
  ApproveClaim_L1Command          → ApproveClaim_L1Handler
  ApproveClaim_L2Command          → ApproveClaim_L2Handler
  ApproveClaim_L3Command          → ApproveClaim_L3Handler
  ApproveClaim_BoardCommand       → ApproveClaim_BoardHandler
  RejectClaimCommand              → RejectClaimHandler
  InitiateClaimPaymentCommand     → InitiateClaimPaymentHandler
  SettleClaimCommand              → SettleClaimHandler
  DisputeClaimCommand             → DisputeClaimHandler

Queries:
  GetClaimByIdQuery               → GetClaimByIdHandler
  GetClaimByNumberQuery           → GetClaimByNumberHandler
  ListCustomerClaimsQuery         → ListCustomerClaimsHandler
  GetClaimDashboardQuery          → GetClaimDashboardHandler (Admin)
  GetPendingApprovalsQuery        → GetPendingApprovalsHandler
  GetTATReportQuery               → GetTATReportHandler
```

---

## 9. Module 08 — Claims Document Processing (FG-008.1)

### 9.1 Scope

SRS Section: 4.8.1 | FR-099 থেকে FR-101

### 9.2 Functional Requirements

| FR ID | বিবরণ | Priority | Acceptance Criteria |
|-------|-------|----------|---------------------|
| FR-099 | Document specs: PDF/JPG/PNG, max 5MB/file, 25MB total/claim, 300 DPI minimum | M1 | Client-side validation, OCR quality check |
| FR-100 | Co-payment + deductible calculation | M1 | Product-level config, breakdown display |
| FR-101 | Reimbursement workflow: doc review → bank/MFS transfer 7–15 working days (insurance company) | M1 | Document verification, payment processing, notifications |

### 9.3 Document Upload Rules

```
Per-file limits:
  Max size: 5 MB
  Formats: PDF, JPG, JPEG, PNG
  Min resolution: 300 DPI
  Max total per claim: 25 MB

Upload flow (FR-242):
  1. Client-side validation (size, format)
  2. Client-side compression: JPEG 80% quality, max 1920×1080
  3. Presigned S3 URL (30-min expiry) from Storage Service
  4. Direct upload to S3 (tus.io chunked, 1MB chunks)
  5. Insurance Engine receives S3 key → stores in claim_documents

OCR (via AI Engine Python:4002):
  Trigger: on document upload
  Accuracy target: >85%
  Fields extracted: date, amount, provider name, diagnosis codes
  Result stored: claim_documents.ocr_extracted_text (JSONB)
  Timeout: 30s → manual review queue
```

### 9.4 Document Type Taxonomy

```
Health Claims:
  - hospital_bill
  - prescription
  - discharge_summary
  - investigation_report
  - doctor_certificate

Motor Claims:
  - police_report
  - repair_estimate
  - photos_damage
  - driving_license
  - vehicle_registration

Travel Claims:
  - flight_itinerary
  - medical_report
  - baggage_loss_report
  - police_report
  - receipts

Death Claims:
  - death_certificate
  - hospital_report
  - nominee_nid
  - nominee_bank_details
```

---

## 10. Module 09 — Fraud Detection & Risk Controls (FG-016)

### 10.1 Scope

SRS Section: 4.16 | FR-175 থেকে FR-188

### 10.2 Functional Requirements

| FR ID | বিবরণ | Priority | Acceptance Criteria |
|-------|-------|----------|---------------------|
| FR-175 | Claims <48hr of policy purchase → auto-flag for manual review | M2 | Auto-flag + notify Claims Officer |
| FR-176 | Same claim type >2x in 12 months → flag + pattern analysis | M2 | Historical analysis, risk scoring |
| FR-177 | Claim amount = 100% of coverage → flag + enhanced verification | M2 | Suspicious pattern, additional docs |
| FR-178 | Non-approved network medical provider → flag | M2 | Provider DB, real-time validation |
| FR-179 | >3 accounts from same device → device fingerprinting alert | M3 | Browser/mobile ID tracking |
| FR-180 | Fraud dashboard: BizAdmin + Focal Person, drill-down | M2 | Real-time alerts, risk scores |
| FR-181 | RACI for monitoring + incident escalation per roles | M1 | Responsibility matrix enforced |

### 10.3 Fraud Detection Rules Reference

| Rule ID | Description | Threshold | Action |
|---------|-------------|-----------|--------|
| FR-182 | Rapid Policy-Claim | < 48 hours from policy issued_at | Auto-flag + manual review |
| FR-183 | Frequent Claims | >2 same type in 12 months | Flag + pattern analysis |
| FR-184 | Full Coverage Claim | Claimed = 100% of sum_insured | Flag + enhanced verification |
| FR-185 | Non-Network Provider | Medical provider not in approved list | Flag + provider verification |
| FD-186 | Geographic Anomaly | Claim location vs registered address > 100km | Flag + location verification |
| FD-187 | Device Fingerprinting | >3 accounts from same device_id | Flag + identity verification |
| FD-188 | Behavioral Pattern | ML-based anomaly score | Risk score + monitoring |

### 10.4 Fraud Risk Scoring

```
Risk Score (0–100):
  +20 pts  : >3 claims in 6 months
  +15 pts  : Single transaction >50K BDT
  +10 pts  : Claim location >100km from registered address
  +25 pts  : Missing NID verification
  +15 pts  : >3 accounts from same device
  +10 pts  : Unusual activity pattern (ML-detected)
  +15 pts  : Claim within 48hr of policy
  +10 pts  : Claim type repeated >2 in 12 months
  +5 pts   : Non-network provider

Risk Levels:
  0–30     : Low Risk → normal processing
  31–60    : Medium Risk → enhanced review
  61–79    : High Risk → mandatory manual review
  80–100   : Critical → auto-flag + senior approval required
```

### 10.5 Fraud Detection Flow

```
On ClaimSubmitted event:
  1. Rapid claim check (policy_issued_at vs claim_submitted_at)
  2. Frequency check (same claim_type, same customer, last 12 months)
  3. Amount check (claimed_amount vs policy.sum_insured)
  4. Provider validation (medical_provider_name vs approved_network list)
  5. Geographic check (claim location vs customer.address)
  6. Device check (customer device_id — from Auth service)
  7. ML score (call AI Engine: POST /fraud/score)
  8. Aggregate risk score
  9. If flagged → insert fraud_checks record, update claim.fraud_score
  10. If score ≥ 80 → notify Claims Officer immediately via Kafka
```

**DB:**
```sql
CREATE TABLE approved_medical_providers (
    provider_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    registration_no VARCHAR(100) UNIQUE NOT NULL,
    location        VARCHAR(255),
    district        VARCHAR(100),
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE fraud_alerts (
    alert_id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    claim_id        UUID REFERENCES claims(claim_id),
    policy_id       UUID REFERENCES policies(policy_id),
    customer_id     UUID NOT NULL,
    alert_type      VARCHAR(100) NOT NULL,
    rule_id         VARCHAR(20) NOT NULL,
    risk_score      DECIMAL(5,2) NOT NULL,
    details         JSONB NOT NULL,
    status          VARCHAR(50) DEFAULT 'OPEN',
    resolved_by     UUID,
    resolved_at     TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 10.6 CQRS Slices

```
Commands:
  RunFraudCheckCommand         → RunFraudCheckHandler
  ResolveFraudAlertCommand     → ResolveFraudAlertHandler
  FlagCustomerForFraudCommand  → FlagCustomerForFraudHandler
  RevokeCustomerAccessCommand  → RevokeCustomerAccessHandler

Queries:
  GetFraudDashboardQuery       → GetFraudDashboardHandler
  GetCustomerRiskScoreQuery    → GetCustomerRiskScoreHandler
  ListOpenFraudAlertsQuery     → ListOpenFraudAlertsHandler
```

---

## 11. Module 10 — Payment Processing (Engine-side) (FG-007)

### 11.1 Scope

SRS Section: 4.7 | FR-070 থেকে FR-080

Insurance Engine এর payment-related responsibilities (actual payment processing Payment Service Node.js:3001 করে, কিন্তু Insurance Engine payment lifecycle track করে)।

### 11.2 Functional Requirements

| FR ID | বিবরণ | Priority | Acceptance Criteria |
|-------|-------|----------|---------------------|
| FR-070 | Payment methods: bKash, Nagad, Rocket, Bank Transfer (MFS/Banking only) | M1 | All MFS integrated, manual verification |
| FR-071 | bKash: production credentials + sandbox | M1 | Success rate >99%, fallback to manual |
| FR-072 | Nagad + Rocket: tokenization for recurring | M3 | PCI-DSS SAQ-A, secure token storage |
| FR-073 | Manual payment: proof upload (bank receipt, bKash screenshot) | M1 | Image <5MB, admin verify within 24hr |
| FR-074 | Payment verification: pending → verified → policy activated OR rejected → refund | M2 | Admin approval manual, automated MFS |
| FR-075 | Payment receipt: transaction ID, amount, date, policy number | M2 | PDF within 5 min via SMS/email |
| FR-076 | Installment plans: Monthly or Yearly (M2 scope) | M3 | Auto-reminders, 15-day grace |
| FR-077 | Retry: exponential backoff, max 3 retries | M2 | Customer notified each attempt |
| FR-078 | Refund: configurable rules, 7-day processing | M2 | Credited to original payment method |
| FR-079 | TigerBeetle: double-entry bookkeeping | M2 | All transactions recorded, real-time reconciliation |
| FR-080 | Payment audit trail: immutable logs, 20-year retention | M1 | PostgreSQL + S3, 20-year retention |

### 11.3 Payment Data Model (Proto → DB)

```protobuf
// proto/insuretech/payment/entity/v1/payment.proto
message Payment {
  string payment_id = 1;
  string transaction_id = 2;            // External gateway TX ID
  string policy_id = 3;
  string claim_id = 4;                  // For claim settlement
  PaymentType type = 5;
  PaymentMethod method = 6;
  PaymentStatus status = 7;
  double amount = 8;
  string currency = 9;                  // "BDT"
  string payer_id = 10;
  string payee_id = 11;
  google.protobuf.Timestamp initiated_at = 12;
  google.protobuf.Timestamp completed_at = 13;
  google.protobuf.Timestamp created_at = 14;
  string gateway = 15;                  // "bKash", "Nagad", "Rocket", "manual"
  string gateway_response = 16;         // JSON string
  string receipt_url = 17;
  int32 retry_count = 18;
  string idempotency_key = 19;
  string proof_url = 20;                // For manual payments
  bool verified_manually = 21;
  string verified_by = 22;
}

enum PaymentType {
  PAYMENT_TYPE_UNSPECIFIED = 0;
  PAYMENT_TYPE_PREMIUM = 1;
  PAYMENT_TYPE_CLAIM_SETTLEMENT = 2;
  PAYMENT_TYPE_REFUND = 3;
  PAYMENT_TYPE_COMMISSION = 4;
  PAYMENT_TYPE_REINSTATEMENT_FEE = 5;
}

enum PaymentMethod {
  PAYMENT_METHOD_UNSPECIFIED = 0;
  PAYMENT_METHOD_BKASH = 1;
  PAYMENT_METHOD_NAGAD = 2;
  PAYMENT_METHOD_ROCKET = 3;
  PAYMENT_METHOD_CARD = 4;
  PAYMENT_METHOD_BANK_TRANSFER = 5;
  PAYMENT_METHOD_CASH = 6;
  PAYMENT_METHOD_CHEQUE = 7;
}

enum PaymentStatus {
  PAYMENT_STATUS_UNSPECIFIED = 0;
  PAYMENT_STATUS_INITIATED = 1;
  PAYMENT_STATUS_PENDING = 2;
  PAYMENT_STATUS_PENDING_MANUAL_VERIFY = 3;
  PAYMENT_STATUS_PROCESSING = 4;
  PAYMENT_STATUS_SUCCESS = 5;
  PAYMENT_STATUS_FAILED = 6;
  PAYMENT_STATUS_REFUNDED = 7;
  PAYMENT_STATUS_CANCELLED = 8;
}
```

**PostgreSQL Table:**
```sql
CREATE TABLE payments (
    payment_id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id          VARCHAR(255) UNIQUE,
    policy_id               UUID REFERENCES policies(policy_id),
    claim_id                UUID REFERENCES claims(claim_id),
    payment_type            VARCHAR(50) NOT NULL,
    payment_method          VARCHAR(50) NOT NULL,
    status                  VARCHAR(50) NOT NULL DEFAULT 'INITIATED',
    amount                  DECIMAL(12,2) NOT NULL CHECK (amount > 0),
    currency                VARCHAR(10) NOT NULL DEFAULT 'BDT',
    payer_id                UUID,
    payee_id                UUID,
    gateway                 VARCHAR(50),
    gateway_response        JSONB,
    receipt_url             VARCHAR(1000),
    retry_count             INT DEFAULT 0,
    idempotency_key         UUID UNIQUE NOT NULL,
    proof_url               VARCHAR(1000),
    verified_manually       BOOLEAN DEFAULT FALSE,
    verified_by             UUID,
    initiated_at            TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at            TIMESTAMP,
    created_at              TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Idempotency table (24hr TTL)
CREATE TABLE payment_idempotency_keys (
    key             UUID PRIMARY KEY,
    payment_id      UUID NOT NULL REFERENCES payments(payment_id),
    response_cached JSONB,
    expires_at      TIMESTAMP NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_payments_policy_id ON payments(policy_id);
CREATE INDEX idx_payments_claim_id ON payments(claim_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_payer_id ON payments(payer_id);
```

### 11.4 Payment Gateway Retry Logic

```
bKash / Nagad Timeout Config:
  connect_timeout = 5s
  read_timeout = 15s

Polly Retry Policy:
  Max retries: 3 (for gateway failures)
  Delays: 1s → 3s → 9s (exponential)
  On each retry: Kafka event PaymentRetryAttempted → customer notified

Circuit Breaker:
  Open after 3 consecutive gateway failures
  Half-open after 30s → test with single request
  Full open after success

Fallback:
  Gateway down >5 min → queue for manual processing → notify admin
```

### 11.5 CQRS Slices

```
Commands:
  InitiatePaymentCommand          → InitiatePaymentHandler
  ConfirmPaymentCommand           → ConfirmPaymentHandler (webhook)
  VerifyManualPaymentCommand      → VerifyManualPaymentHandler
  RejectManualPaymentCommand      → RejectManualPaymentHandler
  InitiateRefundCommand           → InitiateRefundHandler
  ProcessRefundCommand            → ProcessRefundHandler

Queries:
  GetPaymentByIdQuery             → GetPaymentByIdHandler
  GetPaymentStatusQuery           → GetPaymentStatusHandler
  ListPolicyPaymentsQuery         → ListPolicyPaymentsHandler
  GetPaymentAuditTrailQuery       → GetPaymentAuditTrailHandler
```

---

## 12. Module 11 — Audit & Logging (FG-019)

### 12.1 Scope

SRS Section: 4.19 | FR-206 থেকে FR-211

### 12.2 Functional Requirements

| FR ID | বিবরণ | Priority | Acceptance Criteria |
|-------|-------|----------|---------------------|
| FR-206 | Immutable audit logs: policy issue, claim approval, payment, dispute resolution | M1 | Append-only PostgreSQL, tamper detection |
| FR-207 | Data retention: 20-year minimum for regulatory compliance | M2 | Tiered storage hot/warm/cold, automated archival |
| FR-208 | User action tracking: IP, device, timestamp, action type | M3 | Comprehensive, GDPR compliant |
| FR-209 | Partner additional logs per MOU | F | Partner-specific tables, isolated |
| FR-210 | Regulatory portal for IDRA/BFIU data access | M2 | Secure portal, access audit trail |
| FR-211 | Log aggregation + alerting on suspicious patterns | M2 | ELK/CloudWatch, anomaly detection |

### 12.3 Audit Log Data Model

```sql
CREATE TABLE audit_logs (
    log_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type     VARCHAR(100) NOT NULL,  -- 'policy', 'claim', 'payment', etc.
    entity_id       UUID NOT NULL,
    action          VARCHAR(100) NOT NULL,  -- 'POLICY_ISSUED', 'CLAIM_APPROVED', etc.
    actor_id        UUID,                   -- user who performed action
    actor_role      VARCHAR(100),
    old_state       JSONB,
    new_state       JSONB,
    metadata        JSONB DEFAULT '{}',     -- IP, device, user-agent
    ip_address      INET,
    device_info     VARCHAR(500),
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Append-only enforcement (no UPDATE/DELETE allowed)
-- TimescaleDB hypertable for time-series performance
SELECT create_hypertable('audit_logs', 'created_at');

-- Retention: 90 days hot in TimescaleDB, then S3, then Glacier (20yr)
SELECT add_retention_policy('audit_logs', INTERVAL '90 days');
```

### 12.4 Critical Actions That Must Be Audit-Logged

```
Policy Events:
  POLICY_PURCHASE_INITIATED
  POLICY_ISSUED
  POLICY_ACTIVATED
  POLICY_SUSPENDED
  POLICY_CANCELLED
  POLICY_EXPIRED
  POLICY_LAPSED
  POLICY_REINSTATED
  POLICY_RENEWED
  POLICY_ENDORSED

Claim Events:
  CLAIM_SUBMITTED
  CLAIM_DOCS_REQUESTED
  CLAIM_APPROVED_L1
  CLAIM_APPROVED_L2
  CLAIM_APPROVED_L3
  CLAIM_APPROVED_BOARD
  CLAIM_REJECTED
  CLAIM_SETTLED
  FRAUD_FLAGGED
  FRAUD_CONFIRMED
  FRAUD_CLEARED

Payment Events:
  PAYMENT_INITIATED
  PAYMENT_SUCCESS
  PAYMENT_FAILED
  PAYMENT_REFUNDED
  MANUAL_PAYMENT_VERIFIED
  MANUAL_PAYMENT_REJECTED

Admin Events:
  PRODUCT_CREATED
  PRODUCT_UPDATED
  PRODUCT_DEACTIVATED
  CANCELLATION_APPROVED
  CANCELLATION_REJECTED
  ENDORSEMENT_APPROVED
```

---

## 13. Module 12 — Notifications (Outbound Events) (FG-012)

### 13.1 Scope

SRS Section: 4.12 | FR-136 থেকে FR-145

Insurance Engine Kafka-তে domain events publish করে। Kafka Service (Go:8086) সেগুলো consume করে SMS/Email/Push পাঠায়।

### 13.2 Functional Requirements

| FR ID | বিবরণ | Priority | Acceptance Criteria |
|-------|-------|----------|---------------------|
| FR-136 | Kafka event-driven notifications: in-app push, SMS, email | M1 | Event published <100ms |
| FR-137 | Notification triggers: OTP, verification, purchase confirmation, claim updates, renewal reminders, payment confirmations | M1 | Template-based, personalized |
| FR-138 | Opt-in/opt-out for marketing | M2 | GDPR-compliant consent |
| FR-141 | Delivery status tracking: queued, sent, delivered, failed, bounced | M2 | Max 3 retries, exponential backoff |
| FR-142 | Message templates with Bengali/English placeholders | M2 | Variable substitution |
| FR-143 | Rate limiting: max 5 notifications/hour/user | M3 | Redis-based, critical alerts exempt |
| FR-144 | Notification history: 90 days in dashboard | M3 | Read/unread status |

### 13.3 Notification Events Published by Insurance Engine

```
Policy Events → Topic: insurance.policy.events
  PolicyIssued          (trigger: policy activated after payment)
  PolicyRenewed         (trigger: renewal payment confirmed)
  PolicyCancelled       (trigger: cancellation approved)
  PolicyExpiringSoon    (trigger: scheduled — 30d/7d/1d before end_date)
  PolicyExpired         (trigger: end_date passed)
  PolicyLapsed          (trigger: grace period expired)
  PolicyReinstated      (trigger: reinstatement approved)

Claim Events → Topic: insurance.claim.events
  ClaimSubmitted        (trigger: claim successfully submitted)
  ClaimDocsRequested    (trigger: docs requested by reviewer)
  ClaimApproved         (trigger: final approval granted)
  ClaimRejected         (trigger: claim rejected with reason)
  ClaimSettled          (trigger: payment to customer confirmed)
  FraudFlagged          (trigger: fraud score ≥ threshold)

Payment Events → Topic: insurance.payment.events
  PremiumPaymentReceived     (trigger: payment success)
  PremiumPaymentFailed       (trigger: all retries exhausted)
  RefundInitiated            (trigger: refund process started)
  RefundCompleted            (trigger: refund credited)
  ManualPaymentPendingReview (trigger: manual proof uploaded)
```

---

## 14. Proto / gRPC Service Definitions

### 14.1 InsuranceEngineService

```protobuf
// proto/insuretech/policy/services/v1/insurance_service.proto
syntax = "proto3";
package insuretech.policy.services.v1;

option csharp_namespace = "Insuretech.Policy.Services.V1";
option go_package = "github.com/labaid/insuretech/proto/policy/services/v1;policyservicev1";

service InsuranceEngineService {

  // Product Catalog
  rpc CreateProduct(CreateProductRequest) returns (CreateProductResponse);
  rpc UpdateProduct(UpdateProductRequest) returns (UpdateProductResponse);
  rpc DeactivateProduct(DeactivateProductRequest) returns (DeactivateProductResponse);
  rpc GetProduct(GetProductRequest) returns (GetProductResponse);
  rpc ListProducts(ListProductsRequest) returns (ListProductsResponse);
  rpc SearchProducts(SearchProductsRequest) returns (SearchProductsResponse);
  rpc CalculatePremium(CalculatePremiumRequest) returns (CalculatePremiumResponse);

  // Policy Lifecycle
  rpc InitiatePolicyPurchase(InitiatePolicyPurchaseRequest) returns (InitiatePolicyPurchaseResponse);
  rpc ConfirmPolicyPayment(ConfirmPolicyPaymentRequest) returns (ConfirmPolicyPaymentResponse);
  rpc IssuePolicy(IssuePolicyRequest) returns (IssuePolicyResponse);
  rpc GetPolicy(GetPolicyRequest) returns (GetPolicyResponse);
  rpc ListCustomerPolicies(ListCustomerPoliciesRequest) returns (ListCustomerPoliciesResponse);
  rpc GetPolicyDashboard(GetPolicyDashboardRequest) returns (GetPolicyDashboardResponse);

  // Renewals
  rpc InitiateRenewal(InitiateRenewalRequest) returns (InitiateRenewalResponse);
  rpc ConfirmRenewalPayment(ConfirmRenewalPaymentRequest) returns (ConfirmRenewalPaymentResponse);
  rpc EnableAutoRenewal(EnableAutoRenewalRequest) returns (EnableAutoRenewalResponse);
  rpc ReinstatePolicy(ReinstatePolicyRequest) returns (ReinstatePolicyResponse);

  // Cancellation
  rpc RequestCancellation(RequestCancellationRequest) returns (RequestCancellationResponse);
  rpc ApproveCancellation(ApproveCancellationRequest) returns (ApproveCancellationResponse);
  rpc RejectCancellation(RejectCancellationRequest) returns (RejectCancellationResponse);
  rpc CalculateCancellationRefund(CalculateCancellationRefundRequest) returns (CalculateCancellationRefundResponse);

  // Endorsement
  rpc CreateEndorsement(CreateEndorsementRequest) returns (CreateEndorsementResponse);
  rpc ApproveEndorsement(ApproveEndorsementRequest) returns (ApproveEndorsementResponse);
  rpc RejectEndorsement(RejectEndorsementRequest) returns (RejectEndorsementResponse);

  // Claims
  rpc SubmitClaim(SubmitClaimRequest) returns (SubmitClaimResponse);
  rpc UploadClaimDocument(UploadClaimDocumentRequest) returns (UploadClaimDocumentResponse);
  rpc ApproveClaim(ApproveClaimRequest) returns (ApproveClaimResponse);
  rpc RejectClaim(RejectClaimRequest) returns (RejectClaimResponse);
  rpc SettleClaim(SettleClaimRequest) returns (SettleClaimResponse);
  rpc GetClaim(GetClaimRequest) returns (GetClaimResponse);
  rpc ListCustomerClaims(ListCustomerClaimsRequest) returns (ListCustomerClaimsResponse);

  // Fraud
  rpc RunFraudCheck(RunFraudCheckRequest) returns (RunFraudCheckResponse);
  rpc GetFraudDashboard(GetFraudDashboardRequest) returns (GetFraudDashboardResponse);
  rpc ResolveFraudAlert(ResolveFraudAlertRequest) returns (ResolveFraudAlertResponse);

  // Payments (engine-side tracking)
  rpc InitiatePayment(InitiatePaymentRequest) returns (InitiatePaymentResponse);
  rpc ConfirmPayment(ConfirmPaymentRequest) returns (ConfirmPaymentResponse);
  rpc InitiateRefund(InitiateRefundRequest) returns (InitiateRefundResponse);

  // NOTE: Reporting (FG-017, FG-018) belongs to Analytics & Reporting Service (Port 5003)
  // Insurance Engine publishes domain events to Kafka → Port 5003 consumes and generates reports
}
```

### 14.2 Key Request/Response Messages

```protobuf
// IssuePolicyRequest
message IssuePolicyRequest {
  string idempotency_key = 1;          // UUID — required
  string product_id = 2;
  string customer_id = 3;
  Applicant applicant = 4;
  Nominee nominee = 5;
  repeated string selected_rider_ids = 6;
  double sum_insured = 7;
  int32 tenure_months = 8;
  string payment_id = 9;
  string partner_id = 10;
  string agent_id = 11;
  repeated RiskAssessmentAnswer risk_answers = 12;
}

message IssuePolicyResponse {
  string policy_id = 1;
  string policy_number = 2;
  PolicyStatus status = 3;
  double premium_amount = 4;
  google.protobuf.Timestamp start_date = 5;
  google.protobuf.Timestamp end_date = 6;
  string policy_document_url = 7;
}

// SubmitClaimRequest
message SubmitClaimRequest {
  string policy_id = 1;
  string customer_id = 2;
  ClaimType claim_type = 3;
  double claimed_amount = 4;
  google.protobuf.Timestamp incident_date = 5;
  string incident_description = 6;
  repeated string document_ids = 7;
  string payment_channel = 8;
}

message SubmitClaimResponse {
  string claim_id = 1;
  string claim_number = 2;
  ClaimStatus status = 3;
  string submission_hash = 4;
  double fraud_score = 5;
  bool fraud_flagged = 6;
}

// CalculatePremiumRequest
message CalculatePremiumRequest {
  string product_id = 1;
  int32 age = 2;
  double sum_insured = 3;
  int32 tenure_months = 4;
  repeated string rider_ids = 5;
  string occupation = 6;
  bool is_smoker = 7;
  bool has_pre_existing = 8;
}

message CalculatePremiumResponse {
  double base_premium = 1;
  double riders_premium = 2;
  double total_premium = 3;
  repeated PremiumBreakdownItem breakdown = 4;
  bool is_cached_rate = 5;
  google.protobuf.Timestamp rate_valid_until = 6;
}
```

---

## 15. PostgreSQL Database Schema

### 15.1 Complete Table Inventory

```
Insurance Engine Owns These Tables:
├── products
├── product_riders
├── policies                    (partitioned by year)
│   ├── policies_2025
│   └── policies_2026
├── policy_applicants
├── policy_nominees
├── policy_riders               (selected riders per policy)
├── policy_status_history
├── policy_cancellations
├── policy_endorsements
├── claims                      (partitioned by month)
├── claim_documents
├── claim_approvals
├── fraud_checks
├── fraud_alerts
├── approved_medical_providers
├── claim_status_history
├── payments
├── payment_idempotency_keys
└── audit_logs                  (TimescaleDB hypertable)
-- NOTE: idra_reports table belongs to Analytics & Reporting Service (Port 5003)
```

### 15.2 Foreign Key Relationships

```
products (product_id)
  └── product_riders (product_id → products)
  └── policies (product_id → products)
      └── policy_applicants (policy_id → policies)
      └── policy_nominees (policy_id → policies)
      └── policy_riders (policy_id → policies, rider_id → product_riders)
      └── policy_status_history (policy_id → policies)
      └── policy_cancellations (policy_id → policies)
      └── policy_endorsements (policy_id → policies)
      └── claims (policy_id → policies)
          └── claim_documents (claim_id → claims)
          └── claim_approvals (claim_id → claims)
          └── fraud_checks (claim_id → claims)
          └── claim_status_history (claim_id → claims)
      └── payments (policy_id → policies)
          └── payment_idempotency_keys (payment_id → payments)
```

### 15.3 Key Constraints

```sql
-- Policy: premium must be positive
ALTER TABLE policies ADD CONSTRAINT chk_premium_positive CHECK (premium_amount > 0);
-- Policy: sum_insured must be positive
ALTER TABLE policies ADD CONSTRAINT chk_sum_insured_positive CHECK (sum_insured > 0);
-- Nominee: share must be 100 (single nominee rule)
ALTER TABLE policy_nominees ADD CONSTRAINT chk_single_nominee_share CHECK (share_percentage = 100.00);
-- Claim: claimed amount must be positive
ALTER TABLE claims ADD CONSTRAINT chk_claimed_positive CHECK (claimed_amount > 0);
-- Payment: amount must be positive
ALTER TABLE payments ADD CONSTRAINT chk_payment_positive CHECK (amount > 0);
-- Product: min_sum <= max_sum
ALTER TABLE products ADD CONSTRAINT chk_sum_range CHECK (min_sum_insured <= max_sum_insured);
```

---

## 16. Kafka Event Topology

### 16.1 Topics & Partitioning

| Topic | Partitions | Key | Retention | Publisher |
|-------|-----------|-----|-----------|-----------|
| `insurance.policy.events` | 6 | `policy_id` | 30 days | Insurance Engine |
| `insurance.claim.events` | 6 | `claim_id` | 30 days | Insurance Engine |
| `insurance.payment.events` | 6 | `payment_id` | 30 days | Insurance Engine |
| `insurance.fraud.events` | 3 | `customer_id` | 90 days | Insurance Engine |
| `insurance.audit.events` | 12 | `entity_id` | 365 days | Insurance Engine |

### 16.2 Complete Event Inventory

```
Topic: insurance.policy.events
  PolicyIssued
  PolicyRenewed
  PolicyCancelled
  PolicySuspended
  PolicyActivated
  PolicyExpiringSoon       (scheduled — 30d/7d/1d)
  PolicyExpired
  PolicyLapsed
  PolicyReinstated
  PolicyEndorsed
  PolicyPurchaseInitiated
  PolicyPurchaseFailed

Topic: insurance.claim.events
  ClaimSubmitted
  ClaimDocsRequested
  ClaimUnderReview
  ClaimApproved_L1
  ClaimApproved_L2
  ClaimApproved_L3
  ClaimApproved_Board
  ClaimRejected
  ClaimPaymentInitiated
  ClaimSettled
  ClaimDisputed
  ClaimClosed

Topic: insurance.payment.events
  PremiumPaymentInitiated
  PremiumPaymentSuccess
  PremiumPaymentFailed
  PremiumPaymentRetry
  RefundInitiated
  RefundSuccess
  ManualPaymentUploaded
  ManualPaymentVerified
  ManualPaymentRejected
  ClaimSettlementInitiated
  ClaimSettlementSuccess

Topic: insurance.fraud.events
  FraudCheckCompleted
  FraudFlagged
  FraudConfirmed
  FraudCleared
  CustomerAccessRevoked
```

### 16.3 Event Schema Pattern

```json
// PolicyIssued event example
{
  "event_id": "uuid",
  "event_type": "PolicyIssued",
  "event_version": "1.0",
  "timestamp": "2025-03-01T10:00:00Z",
  "aggregate_id": "policy_id",
  "aggregate_type": "Policy",
  "payload": {
    "policy_id": "uuid",
    "policy_number": "LBT-2025-0012-000001",
    "customer_id": "uuid",
    "product_id": "uuid",
    "premium_amount": 1500.00,
    "sum_insured": 50000.00,
    "start_date": "2025-03-01",
    "end_date": "2026-02-28",
    "partner_id": "uuid",
    "agent_id": "uuid"
  },
  "metadata": {
    "correlation_id": "uuid",
    "source_service": "insurance-engine",
    "version": "1.0"
  }
}
```

---

## 17. State Machines

### 17.1 Policy Status State Machine

```
PENDING_PAYMENT
  → ACTIVE            : Payment confirmed (immediate for Non-Life)
  → CANCELLED         : Payment not received + customer cancels

ACTIVE
  → GRACE_PERIOD      : End date passed, payment not received
  → SUSPENDED         : Admin suspension
  → CANCELLED         : Cancellation approved
  → EXPIRED           : End date passed, no renewal (only if grace period skipped)

GRACE_PERIOD (30 days post-expiry)
  → ACTIVE            : Payment received during grace period
  → LAPSED            : 30 days grace period expired, no payment

LAPSED
  → ACTIVE            : Reinstatement within 90 days + penalty paid + Focal approval
  → (Terminal)        : After 90 days — cannot reinstate

SUSPENDED
  → ACTIVE            : Admin lifts suspension
  → CANCELLED         : Admin or customer cancels

CANCELLED   [Terminal]
EXPIRED     [Terminal]
```

### 17.2 Claim Status State Machine

```
SUBMITTED
  → UNDER_REVIEW       : Auto-transition after eligibility validation passes
  → REJECTED           : Eligibility validation fails (policy inactive, not covered, duplicate)

UNDER_REVIEW
  → PENDING_DOCUMENTS  : Reviewer requests more docs (auto if incomplete)
  → APPROVED           : L1 approval (claim < 10K BDT)
  → REJECTED           : Reviewer rejects with reason

PENDING_DOCUMENTS
  → UNDER_REVIEW       : Customer uploads required documents

APPROVED
  → PAYMENT_INITIATED  : Claim payment triggered
  → DISPUTED           : Customer disputes approved amount

PAYMENT_INITIATED
  → SETTLED            : Payment confirmed to customer
  → UNDER_REVIEW       : Payment failed → back for review

REJECTED   [Terminal — unless disputed]
SETTLED    [Terminal]
CLOSED     [Terminal — admin closes]
DISPUTED   → UNDER_REVIEW (re-review cycle)
```

### 17.3 Payment Status State Machine

```
INITIATED
  → PENDING             : Gateway accepted, awaiting callback
  → PENDING_MANUAL_VERIFY: Manual proof uploaded
  → FAILED              : Gateway immediate rejection

PENDING
  → SUCCESS             : Gateway webhook confirms
  → FAILED              : Timeout or gateway failure
  → PROCESSING          : Gateway processing (intermediate)

PENDING_MANUAL_VERIFY
  → SUCCESS             : Admin verifies proof
  → FAILED              : Admin rejects proof

FAILED
  → INITIATED           : Retry (max 3 retries — Polly)
  → CANCELLED           : Customer cancels after failures

SUCCESS   [Terminal]
REFUNDED  [Terminal]
CANCELLED [Terminal]
```

---

## 18. Claims Approval Matrix

### 18.1 Tiered Approval (Insurance Company Only)

| Amount Range | Level | Approver | Max TAT | Auto-escalation |
|-------------|-------|---------|---------|----------------|
| BDT 0–10,000 | L1 Auto/Officer | System Auto OR Claims Officer | 24 hours | Yes, after 24hr |
| BDT 10,001–50,000 | L2 Manager | Claims Manager | 3 days | Yes, after 72hr |
| BDT 50,001–2,00,000 | L3 Head (Joint) | Business Admin + Focal Person (BOTH required) | 7 days | Yes, after 5 days |
| BDT 2,00,001+ | Board | Board + Insurer Approval | 15 days | Yes, after 10 days |

### 18.2 Routing Logic

```csharp
// Approval routing based on claim amount
ApprovalLevel GetApprovalLevel(decimal claimedAmount) {
    return claimedAmount switch {
        <= 10_000     => ApprovalLevel.L1_Auto,
        <= 50_000     => ApprovalLevel.L2_Manager,
        <= 2_00_000   => ApprovalLevel.L3_Joint,
        _             => ApprovalLevel.Board
    };
}

// L3 Joint Approval: BOTH Business Admin and Focal Person must approve
bool IsL3Approved(Claim claim) {
    var approvals = claim.Approvals.Where(a => a.ApprovalLevel == 3);
    bool bizAdminApproved = approvals.Any(a => a.ApproverRole == "BusinessAdmin" && a.Decision == "APPROVED");
    bool focalApproved = approvals.Any(a => a.ApproverRole == "FocalPerson" && a.Decision == "APPROVED");
    return bizAdminApproved && focalApproved;
}
```

### 18.3 Auto-escalation Rules

```
L1 (24hr): System → Claims Officer → if timeout → escalate to L2 Manager + alert
L2 (3 days): if not decided in 72hr → alert Business Admin, escalate
L3 (7 days): if not decided in 5 days → alert Board
Board (15 days): if not decided in 10 days → alert Compliance Team
```

---

## 19. Fraud Detection Rules Reference

### 19.1 Engine-level Rules (FR-182 to FD-188)

| Rule ID | Detection Rule | Threshold | Triggered By | Score Impact |
|---------|---------------|-----------|-------------|-------------|
| FR-182 | Rapid policy-to-claim | policy issued_at to claim submitted_at < 48hr | ClaimSubmitted event | +35 pts |
| FR-183 | Frequent same-type claims | >2 same ClaimType in 12 months for same customer | ClaimSubmitted event | +20 pts |
| FR-184 | Full coverage claim | claimed_amount / sum_insured ≥ 0.98 | ClaimSubmitted event | +25 pts |
| FR-185 | Non-network provider | medical provider not in approved_medical_providers | ClaimSubmitted event | +20 pts |
| FD-186 | Geographic anomaly | claim district vs customer district > 100km | ClaimSubmitted event | +15 pts |
| FD-187 | Device anomaly | same device_id registered for >3 accounts | AccountCreated event | +15 pts |
| FD-188 | ML behavioral | AI Engine fraud score > threshold | Continuous | ML score |

### 19.2 AML Rules (from Security Section)

| TM Rule | Description | Threshold | Action |
|---------|-------------|-----------|--------|
| TM-002 | Rapid policy + claim | Claim < 7 days of purchase | Flag + manual review |
| TM-004 | Frequent claims | >3 claims in 6 months | Flag + pattern analysis |
| TM-005 | Amount near coverage | >90% of coverage claimed | Flag + doc verification |
| TM-008 | Rapid purchases | >3 policies in 7 days | Flag + EDD |
| TM-009 | High-value premium | >BDT 5 lakh | Enhanced due diligence |
| TM-010 | Frequent cancellations | >2 cancellations in 3 months | Flag + investigation |

---

## 20. NFR — Performance & Compliance Constraints

### 20.1 Performance Targets

| Metric | Target | Notes |
|--------|--------|-------|
| gRPC API response (p95) | < 100ms | Category 1 internal API |
| Policy issuance end-to-end | < 500ms | After payment confirmation |
| Premium calculation | < 2s | Including cache check |
| Claim submission | < 3s | Eligibility validation included |
| Fraud check | < 5s | Including ML Engine call |
| Document upload presigned URL | < 500ms | S3 presigned URL generation |
| PDF generation | < 30s | Policy document after payment |
| Database query (p95) | < 100ms | All CRUD operations |
| Redis cache hit | < 5ms | Product catalog, session data |
| Search results | < 500ms | Product search with FTS |

### 20.2 Reliability

| NFR | Target |
|-----|--------|
| System availability | 99.5% (M1) → 99.9% (M2) |
| RTO | < 4 hours |
| RPO | < 1 hour |
| Idempotency key TTL | 24 hours |
| Payment retry max | 3 attempts |
| gRPC circuit breaker | Open after 3 failures in 30s |

### 20.3 Compliance Requirements

| Requirement | Details |
|-------------|---------|
| Audit log retention | 20 years minimum (IDRA requirement) |
| Payment audit trail | PostgreSQL + S3 (20 years) |
| Data encryption at rest | AES-256 |
| Data encryption in transit | TLS 1.3 |
| NID data retention | 20 years (IDRA) |
| Financial records | 7 years minimum (BFIU) |
| Policy document | 20 years in S3/Glacier |
| Idempotency enforcement | All payment + policy issuance APIs |
| IDRA monthly data feed | Insurance Engine publishes events → Port 5003 generates Form IC-1 (Premium) |
| IDRA monthly data feed | Insurance Engine publishes events → Port 5003 generates Form IC-2 (Claims) |
| Fraud event reporting | FraudFlagged event published to Kafka → Port 5003 files IDRA report within 48hr |

### 20.4 Data Retention Tiers

| Data | Hot (PostgreSQL) | Warm (S3) | Cold (Glacier) | Total |
|------|-----------------|-----------|----------------|-------|
| Active Policies | Lifetime | — | — | Lifetime |
| Expired Policies | 1 year | 5 years | 14 years | 20 years |
| Claims Data | Until settled | Post-settlement 5yr | 14 years | 20 years |
| Audit Logs | 90 days (TimescaleDB) | 1 year | 18+ years | 20 years |
| IoT Telemetry | 90 days | 9 months | — | 1 year |
| Payments | 2 years | 5 years | 13 years | 20 years |

---

## 21. Integration Dependencies

### 21.1 Inbound (Other services call Insurance Engine)

| Caller | Method | Purpose |
|--------|--------|---------|
| API Gateway (Go:8080) | gRPC | All client-facing requests routed here |
| Payment Service (Node.js:3001) | gRPC webhook equivalent | Payment confirmation events |
| AI Engine (Python:4001) | gRPC callback | Fraud score, OCR result |

### 21.2 Outbound (Insurance Engine calls other services)

| Target Service | Protocol | Purpose | Timeout |
|---------------|---------|---------|---------|
| Auth Service (Go:8081) | gRPC | Token validation | 5s |
| Partner Management (C#:5002) | gRPC | Partner validation, commission | 10s |
| Payment Service (Node.js:3001) | gRPC | Initiate payment, get status | 15s |
| AI Engine (Python:4001) | gRPC | Fraud check, document OCR | 30s |
| Storage Service (Go:8084) | gRPC | Presigned S3 URL, doc storage | 10s |
| Kafka Service (Go:8086) | Kafka Producer | Domain event publishing | async |
| NID API (External) | REST | NID verification | 10s |
| TigerBeetle | TigerBeetle protocol | Financial ledger entries | 10s |

### 21.3 Payment Gateway Integration (via Payment Service)

| Gateway | Phase | Timeout | Fallback |
|---------|-------|---------|---------|
| bKash | M1 | connect:5s, read:15s | Queue manual |
| Nagad | M2 | connect:5s, read:15s | Queue manual |
| Rocket | M3 | connect:5s, read:15s | Queue manual |
| Manual (proof upload) | M1 | — | Admin verify 24hr |

### 21.4 Retry & Circuit Breaker Config (Polly)

```csharp
// Standard Polly policy for all external gRPC calls
var retryPolicy = Policy
    .Handle<RpcException>()
    .Or<TimeoutException>()
    .WaitAndRetryAsync(5, retryAttempt =>
        TimeSpan.FromSeconds(Math.Pow(2, retryAttempt)),  // 2s, 4s, 8s, 16s, 32s
        onRetry: (exception, timespan, retryCount, context) => {
            // Log retry attempt
            // Publish PaymentRetryAttempted event to Kafka
        });

var circuitBreaker = Policy
    .Handle<RpcException>()
    .CircuitBreakerAsync(
        exceptionsAllowedBeforeBreaking: 3,
        durationOfBreak: TimeSpan.FromSeconds(30));

var combinedPolicy = Policy.WrapAsync(retryPolicy, circuitBreaker);
```

---

## 22. FR → Acceptance Criteria Master Table

### 22.1 M1 Priority FRs (Must-have: March 2025)

| FR ID | Module | Requirement Summary | Acceptance Criteria |
|-------|--------|--------------------|--------------------|
| FR-021 | Product Catalog | Product categories: Health, Auto, Travel, Life, Micro | By category, search + filter |
| FR-022 | Product Catalog | Product search | <500ms, Bengali fuzzy match |
| FR-026 | Product Catalog | Business Admin CRUD | Version history maintained |
| FR-030 | Policy Lifecycle | End-to-end purchase flow | <10 min, progress saved |
| FR-031 | Policy Lifecycle | Applicant info collection | All mandatory fields validated |
| FR-032 | Policy Lifecycle | Single nominee | Only 1 nominee required |
| FR-032-A | Policy Lifecycle | Beneficiary income optional | Submit without income range |
| FR-033 | Policy Lifecycle | NID/Mobile uniqueness | DB constraint, user notified |
| FR-034 | Policy Lifecycle | Policy number format | LBT-YYYY-XXXX-NNNNNN, collision-free |
| FR-039 | Policy Lifecycle | Policy status machine | Transitions logged + timestamped |
| FR-040 | Policy Lifecycle | Customer dashboard | Load <3s, real-time status |
| FR-049 | Renewals | PDF with version history | All versions accessible |
| FR-050 | Renewals | Lifecycle event audit | Immutable log, queryable |
| FR-051 | Cancellation | Cancellation workflow | Reason dropdown, attachment support |
| FR-052 | Cancellation | Dual approval >30 days | 48hr SLA routing |
| FR-053 | Cancellation | Pro-rata refund calc | Configurable fees, breakdown |
| FR-054 | Cancellation | 7-day refund processing | Gateway integration, notification |
| FR-055 | Cancellation | CANCELLED status + notify | Multi-channel, IDRA reporting |
| FR-056 | Endorsement | Endorsement types | Amendment forms, validation |
| FR-059 | Endorsement | Endorsement document | PLN-001/END-01 format, PDF |
| FR-060 | Endorsement | >10% change approval | Approval workflow, threshold config |
| FR-061 | Business Rules | Premium fallback | Fallback tested, cache, queue notify |
| FR-063 | Business Rules | Duplicate policy block | 30-day window, cross-product allowed |
| FR-065 | Business Rules | Claim state machine | Invalid transitions blocked |
| FR-066 | Business Rules | Claim transition rules | Auto-routing, notifications |
| FR-081 | Claims | Claim submission form | <5 min, draft saving |
| FR-082 | Claims | Claim eligibility check | <3s, clear errors |
| FR-083 | Claims | Claim number + hash | CLM format, SHA-256 integrity |
| FR-099 | Doc Processing | Document size limits | 5MB/file, 25MB total, 300 DPI |
| FR-100 | Doc Processing | Co-payment deductible | Product config, breakdown |
| FR-101 | Doc Processing | Reimbursement workflow | Doc review, 7–15 day payment |
| FR-070 | Payment | MFS payment methods | bKash/Nagad/Rocket/manual |
| FR-071 | Payment | bKash integration | >99% success rate |
| FR-073 | Payment | Manual payment proof | <5MB upload, 24hr verify |
| FR-080 | Payment | Payment audit trail | PostgreSQL + S3, 20-year |
| FR-181 | Fraud | RACI for monitoring | Responsibility matrix enforced |
| FR-206 | Audit | Immutable audit logs | Append-only, tamper detection |

### 22.2 M2 Priority FRs (Grand Launch: April 2025)

| FR ID | Module | Requirement Summary | Acceptance Criteria |
|-------|--------|--------------------|--------------------|
| FR-023 | Product Catalog | Product details display | All info before purchase, PDF download |
| FR-023-A | Product Catalog | Unit-wise plan, coverage adjust | Coverage amount adjustable |
| FR-023-B | Product Catalog | Risk assessment questions | Every plan has questions |
| FR-035 | Policy Lifecycle | PDF + QR code | Within 30s of payment |
| FR-036 | Policy Lifecycle | Doc delivery SMS + email | Within 5 min, retry |
| FR-037 | Policy Lifecycle | Immediate activation Non-Life | Real-time status update |
| FR-043 | Renewals | Renewal reminders 30d/7d/1d | On schedule, tracked |
| FR-044 | Renewals | Manual one-click renewal | <3 min, updated doc |
| FR-047 | Renewals | Grace period 30 days | Status GRACE_PERIOD, daily reminders |
| FR-048 | Renewals | Auto-lapse + reinstatement | 90-day window, penalty |
| FR-058 | Endorsement | Sum decrease refund | Credit to premium account |
| FR-062 | Business Rules | Premium edge cases | All cases handled, clear messaging |
| FR-074 | Payment | Verification workflow | Admin approval manual, automated MFS |
| FR-075 | Payment | Payment receipt PDF | Within 5 min SMS/email |
| FR-077 | Payment | Retry exponential backoff | Max 3, customer notified |
| FR-078 | Payment | Refund processing | 7-day, original method |
| FR-079 | Payment | TigerBeetle double-entry | All transactions, real-time reconcile |
| FR-084 | Claims | Auto-notify partner on submit | <60s notification |
| FR-090 | Claims | Partner notes + approve/reject | Timestamped, reason required |
| FR-175 | Fraud | Rapid claim flag <48hr | Auto-flag, Claims Officer notified |
| FR-176 | Fraud | Frequent claim type flag | Historical analysis |
| FR-177 | Fraud | Full coverage claim flag | Enhanced verification |
| FR-178 | Fraud | Non-network provider flag | Provider DB, real-time |
| FR-180 | Fraud | Fraud dashboard | Real-time alerts, drill-down |
| FR-207 | Audit | 20-year retention | Tiered storage, automated archival |
| FR-210 | Audit | IDRA/BFIU portal access | Secure portal, audit trail |
| FR-211 | Audit | Log aggregation + alerts | ELK/CloudWatch, anomaly detection |

### 22.3 M3 Priority FRs (Enhancement: August 2025)

| FR ID | Module | Requirement Summary | Acceptance Criteria |
|-------|--------|--------------------|--------------------|
| FR-024 | Product Catalog | Premium calculator | <2s, breakdown shown |
| FR-025 | Product Catalog | Product comparison 2x | Comparison table |
| FR-028 | Product Catalog | Redis cache 5-min TTL | >80% hit rate |
| FR-029 | Product Catalog | Bengali + English i18n | Language toggle |
| FR-038 | Policy Lifecycle | 5-day cooling-off refund | 24hr cancellation, refund |
| FR-045 | Renewals | Auto-repurchase opt-in | Consent recorded, 7d before |
| FR-046 | Renewals | Update during renewal | Limited fields, verification major |
| FR-064 | Business Rules | Policy merge workflow | Data integrity, audit |
| FR-068 | Business Rules | Grace period full logic | Enforced, customer notified |
| FR-085 | Claims | Real-time status + push/SMS | <5s updates |
| FR-086 | Claims | Tiered approval routing | Auto-routing, escalation |
| FR-087 | Claims | Doc verification OCR | <10s, >85% accuracy |
| FR-088 | Claims | Chat interface | Real-time, file attach |
| FR-091 | Claims | Joint L3 approval 50K–2L | Both required, 5-day escalation |
| FR-092 | Claims | Auto-payment on approval | 24hr, confirmation sent |
| FR-094 | Claims | Fraud detection full | Auto-flag, risk score |
| FR-095 | Claims | Auto-revoke on fraud | Suspension, appeal |
| FR-096 | Claims | Balance sheet all levels | Daily/monthly/quarterly |
| FR-097 | Claims | TAT tracking + SLA alert | Real-time monitoring |
| FR-098 | Claims | Claim history analytics | Frequency, avg amount |
| FR-072 | Payment | Nagad + Rocket tokenized | PCI-DSS SAQ-A |
| FR-179 | Fraud | Device fingerprinting | Browser/mobile ID |
| FR-208 | Audit | User action full tracking | IP, device, GDPR |

---

## Appendix A — VSA Folder Structure

```
InsuranceEngine/
├── src/
│   ├── InsuranceEngine.API/              ← gRPC service entry point
│   │   └── Program.cs
│   │
│   ├── InsuranceEngine.Application/      ← All CQRS slices
│   │   ├── Products/
│   │   │   ├── CreateProduct/
│   │   │   │   ├── CreateProductCommand.cs
│   │   │   │   ├── CreateProductHandler.cs
│   │   │   │   └── CreateProductValidator.cs
│   │   │   ├── UpdateProduct/
│   │   │   ├── GetProduct/
│   │   │   ├── ListProducts/
│   │   │   ├── SearchProducts/
│   │   │   └── CalculatePremium/
│   │   │
│   │   ├── Policies/
│   │   │   ├── InitiatePurchase/
│   │   │   ├── ConfirmPayment/
│   │   │   ├── IssuePolicy/
│   │   │   ├── GetPolicy/
│   │   │   ├── ListPolicies/
│   │   │   └── GetDashboard/
│   │   │
│   │   ├── Renewals/
│   │   │   ├── InitiateRenewal/
│   │   │   ├── ConfirmRenewal/
│   │   │   ├── EnableAutoRenewal/
│   │   │   ├── ReinstatePolicy/
│   │   │   └── ProcessGracePeriodExpiry/  ← Background job
│   │   │
│   │   ├── Cancellations/
│   │   │   ├── RequestCancellation/
│   │   │   ├── ApproveCancellation/
│   │   │   ├── RejectCancellation/
│   │   │   └── CalculateRefund/
│   │   │
│   │   ├── Endorsements/
│   │   │   ├── CreateEndorsement/
│   │   │   ├── ApproveEndorsement/
│   │   │   └── RejectEndorsement/
│   │   │
│   │   ├── Claims/
│   │   │   ├── SubmitClaim/
│   │   │   ├── UploadDocument/
│   │   │   ├── ApproveClaim_L1/
│   │   │   ├── ApproveClaim_L2/
│   │   │   ├── ApproveClaim_L3/
│   │   │   ├── ApproveClaim_Board/
│   │   │   ├── RejectClaim/
│   │   │   ├── SettleClaim/
│   │   │   └── GetClaim/
│   │   │
│   │   ├── Fraud/
│   │   │   ├── RunFraudCheck/
│   │   │   ├── ResolveFraudAlert/
│   │   │   └── GetFraudDashboard/
│   │   │
│   │   ├── Payments/
│   │   │   ├── InitiatePayment/
│   │   │   ├── ConfirmPayment/
│   │   │   └── InitiateRefund/
│   │   │
│   │   └── Shared/
│   │       ├── Behaviors/              ← MediatR pipeline behaviors
│   │       │   ├── ValidationBehavior.cs
│   │       │   ├── LoggingBehavior.cs
│   │       │   └── AuditBehavior.cs
│   │       └── Interfaces/
│   │
│   ├── InsuranceEngine.Domain/           ← Domain entities, value objects
│   │   ├── Policy/
│   │   ├── Claim/
│   │   ├── Product/
│   │   └── Payment/
│   │
│   ├── InsuranceEngine.Infrastructure/   ← EF Core, Redis, Kafka, external clients
│   │   ├── Persistence/
│   │   │   ├── InsuranceDbContext.cs
│   │   │   └── Migrations/
│   │   ├── Cache/
│   │   │   └── RedisCache.cs
│   │   ├── Messaging/
│   │   │   └── KafkaEventPublisher.cs
│   │   └── ExternalServices/
│   │       ├── PaymentServiceClient.cs
│   │       ├── AIEngineClient.cs
│   │       └── StorageServiceClient.cs
│   │
│   └── InsuranceEngine.Contracts/        ← Proto-generated C# classes
│       └── gen/csharp/                   ← from root/gen/csharp/
│
└── tests/
    ├── InsuranceEngine.UnitTests/
    ├── InsuranceEngine.IntegrationTests/
    └── InsuranceEngine.E2ETests/
```

---

## Appendix B — Environment Configuration

```
INSURANCE_ENGINE_PORT=5001
POSTGRES_CONNECTION_STRING=Host=...;Database=insuretech;Username=...;Password=...
REDIS_CONNECTION_STRING=localhost:6379
KAFKA_BOOTSTRAP_SERVERS=localhost:9092
AUTH_SERVICE_GRPC=http://auth-service:8081
PARTNER_SERVICE_GRPC=http://partner-service:5002
PAYMENT_SERVICE_GRPC=http://payment-service:3001
AI_ENGINE_GRPC=http://ai-engine:4001
STORAGE_SERVICE_GRPC=http://storage-service:8084
TIGERBEETLE_ADDRESS=localhost:3000
S3_BUCKET=insuretech-documents
POLICY_NUMBER_SEQUENCE_PREFIX=LBT
CLAIM_NUMBER_SEQUENCE_PREFIX=CLM
PREMIUM_CACHE_TTL_SECONDS=300
FRAUD_SCORE_FLAG_THRESHOLD=60
FRAUD_SCORE_CRITICAL_THRESHOLD=80
GRACE_PERIOD_DAYS=30
REINSTATEMENT_WINDOW_DAYS=90
COOLING_OFF_DAYS=5
CANCELLATION_ADMIN_FEE_PERCENT=5.0
CANCELLATION_CHARGE_PERCENT=2.0
DUPLICATE_POLICY_WINDOW_DAYS=30
MAX_IDEMPOTENCY_KEY_TTL_HOURS=24
```

---

## Appendix C — Audit Categories for Antigravity Verification

```
AUDIT CATEGORY 1: Proto Compliance
  → All C# classes in gen/csharp/ match proto definitions?
  → All enum values mapped correctly?
  → All required fields (non-null) enforced in handlers?

AUDIT CATEGORY 2: API Contract Compliance
  → All gRPC methods in InsuranceEngineService implemented?
  → Request/Response message shapes match HTML API docs?
  → Error codes consistent with OpenAPI docs?

AUDIT CATEGORY 3: SRS Business Compliance
  → Policy number format: LBT-YYYY-XXXX-NNNNNN?
  → Claim number format: CLM-YYYY-XXXX-NNNNNN?
  → Single nominee only (FR-032)?
  → Beneficiary income optional (FR-032-A)?
  → Cooling-off 5 days (FR-038)?
  → Grace period 30 days (FR-047)?
  → Reinstatement window 90 days (FR-048)?
  → Endorsement document suffix correct (FR-059)?
  → Sum insured change >10% needs approval (FR-060)?

AUDIT CATEGORY 4: CRUD Completeness
  → Products: Create, Read, Update, Deactivate (no Delete)?
  → Policies: full lifecycle CRUD?
  → Claims: full submission + approval workflow?
  → Endorsements: Create + Approve + Reject?

AUDIT CATEGORY 5: State Machine Correctness
  → Policy status transitions: all valid paths implemented?
  → Invalid transitions: blocked at handler level?
  → Claim status transitions: all valid paths?
  → History logged on every status change?

AUDIT CATEGORY 6: Fraud Rules
  → FR-182 (48hr rapid claim) implemented?
  → FR-183 (>2 same type/12mo) implemented?
  → FR-184 (100% coverage) implemented?
  → FR-185 (non-network provider) implemented?

AUDIT CATEGORY 7: Financial Rules
  → Pro-rata refund formula correct?
  → Co-payment + deductible formula correct?
  → TigerBeetle entries: every payment creates double-entry?
  → Idempotency enforced on all payment + issuance APIs?

AUDIT CATEGORY 8: Approval Matrix
  → L1 auto ≤10K BDT?
  → L2 Manager 10K–50K?
  → L3 JOINT (both BizAdmin + Focal) 50K–2L?
  → Board 2L+?
  → TAT enforcement + auto-escalation?

AUDIT CATEGORY 9: Architecture & Pattern
  → VSA: no cross-slice dependencies?
  → CQRS: commands → PostgreSQL write, queries → Redis → read replica?
  → Polly: retry + circuit breaker on all external calls?
  → FluentValidation: every command has a validator?
  → Audit log: every critical action logged?
  → Kafka: domain event published after every state change?
```

---

*Document generated from SRS v3.11 (FINAL_DRAFT, Feb 2026)*
*Target: AI Tool feed for Insurance Engine audit verification*
*Coverage: 11 domain modules + 1 notification module (total 12), ~120 FRs, complete proto/DB/Kafka/state machine specs*
*Note: Reporting (FG-017, FG-018, FR-189–205) belongs to Analytics & Reporting Service (Port 5003) — NOT Insurance Engine*