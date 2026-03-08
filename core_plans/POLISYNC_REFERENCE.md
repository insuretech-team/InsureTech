# PoliSync — Complete Reference & Implementation Guide

> **C# .NET 8 Insurance Commerce & Policy Engine**
> **Proto-first · CQRS · MediatR · gRPC · Kafka · PostgreSQL**
> **Last Updated:** March 6, 2026 | **Overall Completion: ~55%**
>
> **Proto Source of Truth:** `E:\Projects\InsureTech\proto\insuretech\`
> **C# Generated Stubs:** `E:\Projects\InsureTech\gen\csharp\`

---

## 1. What Is PoliSync?

PoliSync is the C# .NET 8 insurance commerce and policy engine for LabAid InsureTech. It owns business logic, domain rules, and gRPC APIs for the full insurance lifecycle — product discovery, premium calculation, policy issuance, underwriting, claims, renewals, endorsements, and commission.

**All database persistence is delegated** to the Go `insurance-service` (port 50115). Go owns schema migrations and GORM mappings; C# owns domain logic.

Proto-generated C# stubs: `E:\Projects\InsureTech\gen\csharp\` (no drift).

---

## 2. Platform Architecture

```
Client → Gateway (Go :8080)
  ├── InScore Go Services (:50050–50113, :50180–50321)
  │     authn, authz, kyc, partner, b2b, fraud, payment,
  │     notification, workflow, docgen, storage, etc.
  └── PoliSync C# Services (:50120–50211)
        product, quote, order, commission, policy,
        underwriting, claim (+ endorsement, renewal, refund)
              │ gRPC (all DB ops)
              ▼
        insurance-service (Go :50115) → GORM → PostgreSQL
```

### Data Gateway Pattern

Each C# bounded context delegates persistence via a gateway interface:

| Interface | Implementation | Go RPC Target |
|-----------|---------------|--------------|
| `IPolicyDataGateway` | `GoPolicyDataGateway` | `InsuranceService.Create/Get/Update/List/DeletePolicy` |
| `IClaimDataGateway` | `GoClaimDataGateway` | `InsuranceService.Create/Get/Update/List/DeleteClaim` |
| `IUnderwritingDataGateway` | `GoUnderwritingDataGateway` | `InsuranceService.*UnderwritingDecision*, *HealthDeclaration*, *Quote*` |
| `IEndorsementDataGateway` | `GoEndorsementDataGateway` | `InsuranceService.Create/Get/Update/Delete/ListEndorsement*` |
| `IRenewalDataGateway` | `GoRenewalDataGateway` | `InsuranceService.*RenewalSchedule*, *RenewalReminder*, *GracePeriod*` |
| `IRefundPaymentGateway` | `PaymentServiceGrpcClient` | `PaymentService.InitiateRefund` |
| `IProductDataGateway` | `GoProductDataGateway` | `InsuranceService.*Product*, *ProductPlan*, *Rider*, *PricingConfig*` |
| `IQuotationDataGateway` | `GoQuotationDataGateway` | `InsuranceService.Create/Get/Update/Delete/ListQuotation` |
| `IOrderDataGateway` | `GoOrderDataGateway` | `OrderService.Create/Get/List/InitiatePayment/ConfirmPayment/Cancel/GetOrderStatus` |
| `IBeneficiaryDataGateway` | `GoBeneficiaryDataGateway` | `InsuranceService.*Beneficiary*, *Individual*, *Business*` |

---

## 3. Technology Stack

| Concern | Choice |
|---|---|
| Runtime | .NET 8 LTS (C# 12) |
| gRPC | `Grpc.AspNetCore 2.60` |
| DB persistence | Delegated to Go `insurance-service` :50115 |
| CQRS | MediatR 12 |
| Validation | FluentValidation 11 |
| Messaging | Confluent.Kafka (idempotent, Acks=All) |
| Caching | Redis 7+ via IDistributedCache (5-min TTL) |
| PII | AES-256-GCM (`AesGcmPiiEncryptor`) |
| Money type | `common.v1.Money` (int64 paisa + currency companion) |
| Observability | OpenTelemetry 1.7 + Serilog 3.1 + Prometheus |
| Auth | server side session (web portals) | JWT for app(ios/android)
| Testing | xUnit + Moq + Testcontainers |

---

## 4. Service Registry

*(Source: `backend/inscore/configs/services.yaml`)*

### Go / InScore

| Service | gRPC | HTTP | Description |
|---------|------|------|-------------|
| tenant | 50050 | 50051 | Org & config |
| authn | 50060 | 50061 | Identity, session, OTP |
| authz | 50070 | 50071 | RBAC/ABAC |
| audit | 50080 | 50081 | Audit logging |
| kyc | 50090 | 50091 | NID verification |
| partner | 50100 | 50101 | Agency/broker |
| beneficiary | 50110 | 50111 | Nominee management |
| b2b | 50112 | 50113 | B2B employee/dept |
| **insurance** ⭐ | **50115** | **50116** | **PoliSync data layer** |
| workflow | 50180 | 50181 | Process orchestration |
| payment | 50190 | 50191 | Payment gateway |
| ledger | 50200 | 50201 | Double-entry (LedgerClaw) |
| fraud | 50220 | 50221 | Fraud detection |
| notification | 50230 | 50231 | SMS/email/push |
| docgen | 50280 | 50281 | PDF generation |
| storage | 50290 | 50291 | S3/blob storage |
| gateway | — | 8080 | HTTP API gateway |

### C# / PoliSync

| Service | gRPC | HTTP | Domain |
|---------|------|------|--------|
| product | 50120 | 50121 | Product catalog, plans, riders, pricing |
| quote | 50130 | 50131 | Quotation lifecycle, premium calc |
| order | 50140 | 50141 | Checkout, 5-step purchase flow |
| commission | 50150 | 50151 | Agent commission, revenue share |
| policy | 50160 | 50161 | Policy lifecycle (+ endorsement, renewal, refund) |
| underwriting | 50170 | 50171 | Risk scoring, health declaration |
| claim | 50210 | 50211 | FNOL, 4-tier approval, settlement |

### Product Categories (SRS FR-021)

**Life** · **NonLife:** Health, Auto, Travel, Fire, Marine, Property, Liability, Loss-of-Income, Goods-in-Transit, Crops, Livestock, Credit, Surety, Engineering, Aviation, Pet, Cyber, Terrorism, Natural-Disaster, Device

---

## 5. Implementation Status

### Phase Summary

| Phase | gRPC Service | Domain Logic | Tests | Remaining |
|-------|-------------|-------------|-------|-----------|
| 1 Infrastructure | ✅ | ✅ | ✅ | Done |
| 2 Products & Pricing | ✅ Full | ✅ ~80% | ❌ | ~3 days |
| 3 Quotes & Underwriting | ✅ UW; 🟡 Quote CRUD | 🟡 30% | Contract ✅ | ~1 week |
| 4 Orders & Policy | ✅ Policy; 🟡 Order orchestration | 🟡 35% | Contract ✅ | ~1 week |
| 5 Endorsement & Renewal | ✅ Both | ❌ 10% | Contract ✅ | ~1.5 weeks |
| 6 Claims & Fraud | ✅ Claim | ❌ 10% | Contract ✅ | ~1.5 weeks |
| 7 Commission & Refund | ✅ Both | ❌ 10% | Contract ✅ | ~1 week |
| 8 Hardening | ❌ | ❌ | ❌ | ~1.5 weeks |

### What's Done Per Phase

**Phase 1 – Infrastructure (✅ 100%):**
SharedKernel (13 files), Infrastructure (50+ files incl. 21 repos, 11 gRPC client stubs, KafkaEventBus, AesGcmPiiEncryptor, InsuranceServiceClient), ApiHost (interceptors, health checks, OTEL, Serilog)

**Phase 2 – Products (~70%):**
Product aggregate with state machine, ProductPlan, Rider, PricingConfig, PricingEngine (12 commands, 12 queries, ProductGrpcService, 4 EF configs). Missing: validators, mappers, Redis cache wiring, tests.

**Phase 3 – Underwriting (✅), Quotes (🟡):**
UnderwritingGrpcService + GoUnderwritingDataGateway + contract tests. Underwriting approval applies loading-adjusted premium back to linked quotation. Underwriting rejection projects quotation status/reason (`REJECTED`). Kafka flow wired: `insuretech.quotation.submitted.v1` consumer seeds underwriting quote and `insuretech.underwriting.decision_made.v1` is published on approve/reject. Risk scoring service and 20+ case score matrix tests are now implemented. Quote side now has `QuotesGrpcService` plus `IQuotationDataGateway` / `GoQuotationDataGateway` over InsuranceService quotation CRUD. Remaining gaps: dedicated public quote proto, premium calculation query, expiry job, and tests.

**Phase 4 – Policy (✅ gRPC), Orders (🟡):**
PolicyGrpcService (Create/Get/List/Update/Cancel/Renew/GenerateDocument/Issue) + GoPolicyDataGateway + contract tests + 4 commands + 3 queries. Orders now delegate to Go `orders-service` via `OrderServiceGrpcClient` plus `IOrderDataGateway` / `GoOrderDataGateway`; C# validates quotation and order transitions before delegating the lifecycle RPCs. `OrderPaymentConfirmedConsumer` now issues the policy inside PoliSync and publishes a `policy.issued` projection so Go `orders-service` can link `policy_id` and move the order to `POLICY_ISSUED`.

**Phase 5 – Endorsement & Renewal (✅ gRPC):**
Both gRPC services + data gateways + contract tests.

**Phase 6 – Claims (✅ gRPC):**
ClaimGrpcService (Submit/Get/List/Upload/Approve/Reject/Settle/RequestMoreDocuments/Dispute) + GoClaimDataGateway + contract tests + 2 commands.

**Phase 7 – Commission & Refund (✅ gRPC):**
Both gRPC services + contract tests. RefundGrpcService calls Go `PaymentService.InitiateRefund` via `IRefundPaymentGateway`.

---

## 6. State Machine Specifications & Implementation Guide

This section specifies every state machine, its transition rules, and C# implementation instructions for the domain aggregate layer.

---

### 6.1 Product State Machine

```
         Create()
           │
           ▼
        ┌──────┐
        │ DRAFT │
        └───┬───┘
            │ Activate()
            ▼
       ┌────────┐
    ┌──│ ACTIVE │──┐
    │  └────────┘  │
    │ Deactivate() │ Discontinue()
    ▼              ▼
┌──────────┐  ┌──────────┐
│ INACTIVE │  │ ARCHIVED │ (terminal)
└────┬─────┘  └──────────┘
     │ Activate()    ▲
     └───► ACTIVE    │
     │ Discontinue() │
     └───────────────┘
```

**Status:** ✅ Implemented in `src/PoliSync.Products/Domain/Product.cs`

| Transition | From | To | Rule |
|-----------|------|-----|------|
| `Create()` | — | DRAFT | Validates: `ProductCode` non-empty, `BasePremium > 0`, `SumInsuredMin < SumInsuredMax`, `CoveragePeriodDays ≥ 1` |
| `Activate()` | DRAFT, INACTIVE | ACTIVE | Raises `ProductActivatedEvent` → invalidate Redis cache |
| `Deactivate()` | ACTIVE | INACTIVE | Raises `ProductDeactivatedEvent` → invalidate Redis cache |
| `Discontinue()` | ACTIVE, INACTIVE | ARCHIVED | Terminal state. Raises `ProductDiscontinuedEvent` |
| `Update()` | DRAFT only | DRAFT | Increments `Version`; rejected if not DRAFT |

**Remaining implementation:**
- [ ] `CreateProductCommandValidator` — FluentValidation: `ProductCode` format `^[A-Z]{2,10}-\d{3,6}$`, `BasePremiumBdt > 0`, `SumInsuredMinBdt > 0 && < SumInsuredMaxBdt`
- [ ] Redis cache invalidation in `ProductActivatedEventHandler` and `ProductDeactivatedEventHandler`
- [ ] Unit test: all 5 transitions + rejected transitions (e.g., `Activate()` from `ARCHIVED` → fail)

---

### 6.2 Quotation State Machine

```
  CreateQuotation()
        │
        ▼
    ┌───────┐
    │ DRAFT │
    └───┬───┘
        │ SubmitQuotation()
        ▼
  ┌───────────┐
  │ SUBMITTED │
  └─────┬─────┘
        │ (partner/insurer action)
        ├──► RECEIVED ──► APPROVED ──► (convert to Order)
        │                     │
        │                     └──► REJECTED
        └──► EXPIRED (auto, after 30 days)
```

**Status:** ❌ Not implemented

#### Implementation Steps

**Step 1 — Create `src/PoliSync.Quotes/Domain/Quotation.cs`**

```csharp
public sealed class Quotation : Entity
{
    public Guid TenantId { get; private set; }
    public string QuotationNumber { get; private set; } = string.Empty; // QT-{auto}
    public Guid ProductId { get; private set; }
    public Guid PlanId { get; private set; }
    public Guid CustomerId { get; private set; }
    public QuotationStatus Status { get; private set; }
    public DateTime ExpiryDate { get; private set; }

    // Premium breakdown — all int64 paisa
    public long BasePremium { get; private set; }
    public long RiderPremium { get; private set; }
    public long LoadingAmount { get; private set; }
    public long DiscountAmount { get; private set; }
    public long VatTax { get; private set; }
    public long ServiceFee { get; private set; }
    public long TotalPayable { get; private set; }

    public string? RejectionReason { get; private set; }

    public static Result<Quotation> Create(Guid tenantId, Guid productId, Guid planId,
        Guid customerId, long basePremium, long riderPremium, int expiryDays = 30)
    {
        if (basePremium <= 0) return Result.Fail<Quotation>("INVALID_PREMIUM", "Base premium must be positive");

        var q = new Quotation
        {
            TenantId = tenantId,
            QuotationNumber = $"QT-{Guid.NewGuid().ToString()[..8].ToUpper()}",
            ProductId = productId,
            PlanId = planId,
            CustomerId = customerId,
            Status = QuotationStatus.Draft,
            BasePremium = basePremium,
            RiderPremium = riderPremium,
            ExpiryDate = DateTime.UtcNow.AddDays(expiryDays)
        };
        q.RecalcTotal();
        q.RaiseDomainEvent(new QuotationCreatedEvent(q.Id));
        return Result.Ok(q);
    }

    public Result Submit()
    {
        if (Status != QuotationStatus.Draft)
            return Result.Fail("INVALID_TRANSITION", $"Cannot submit from {Status}");
        if (DateTime.UtcNow > ExpiryDate)
            return Result.Fail("EXPIRED", "Quotation has expired");

        Status = QuotationStatus.Submitted;
        RaiseDomainEvent(new QuotationSubmittedEvent(Id));
        return Result.Ok();
    }

    public Result Approve()
    {
        if (Status is not (QuotationStatus.Submitted or QuotationStatus.Received))
            return Result.Fail("INVALID_TRANSITION", $"Cannot approve from {Status}");
        if (DateTime.UtcNow > ExpiryDate)
            return Result.Fail("EXPIRED", "Quotation has expired");

        Status = QuotationStatus.Approved;
        RaiseDomainEvent(new QuotationApprovedEvent(Id));
        return Result.Ok();
    }

    public Result Reject(string reason)
    {
        if (Status is not (QuotationStatus.Submitted or QuotationStatus.Received))
            return Result.Fail("INVALID_TRANSITION", $"Cannot reject from {Status}");

        Status = QuotationStatus.Rejected;
        RejectionReason = reason;
        RaiseDomainEvent(new QuotationRejectedEvent(Id, reason));
        return Result.Ok();
    }

    public Result Expire()
    {
        if (Status is QuotationStatus.Approved or QuotationStatus.Rejected)
            return Result.Fail("TERMINAL", "Already in terminal state");

        Status = QuotationStatus.Expired;
        return Result.Ok();
    }

    public void ApplyLoading(long loadingAmount)
    {
        LoadingAmount = loadingAmount;
        RecalcTotal();
    }

    private void RecalcTotal()
    {
        var subtotal = BasePremium + RiderPremium + LoadingAmount - DiscountAmount;
        VatTax = (long)(subtotal * 0.15m);  // 15% VAT
        TotalPayable = subtotal + VatTax + ServiceFee;
    }
}

public enum QuotationStatus
{
    Draft = 1, Submitted = 2, Received = 3,
    Approved = 4, Rejected = 5, Expired = 6
}
```

**Step 2 — `CalculatePremiumQuery` (stateless, no DB write)**

```csharp
// Load product + plan + riders from cache/insurance-service
// Apply pricing rules via PricingEngine.Evaluate()
// Apply underwriting loading factor
// Return breakdown: { base, rider, loading, discount, vat, service_fee, total }
```

**Step 3 — Expiry background job**: `IHostedService` that runs hourly, queries quotations where `expiry_date < now AND status NOT IN (APPROVED, REJECTED, EXPIRED)`, calls `Expire()` on each.

---

### 6.3 Order State Machine

```
  CreateOrder(quotationId)
        │
        ▼
   ┌─────────┐
   │ PENDING  │
   └────┬─────┘
        │ InitiatePayment()
        ▼
┌──────────────────┐
│ PAYMENT_INITIATED │
└────────┬─────────┘
         │ ConfirmOrder() (payment callback)
         ▼
      ┌──────┐
      │ PAID │──── IssuePolicyCommand dispatched
      └──┬───┘
         │ Policy created
         ▼
  ┌──────────────┐
  │ POLICY_ISSUED │ (terminal success)
  └──────────────┘

  Any non-terminal → CancelOrder() → CANCELLED
  Payment timeout → FAILED
```

**Status:** ❌ Not implemented

#### Implementation Steps

**Step 1 — Create `src/PoliSync.Orders/Domain/Order.cs`**

```csharp
public sealed class Order : Entity
{
    public Guid QuotationId { get; private set; }
    public Guid CustomerId { get; private set; }
    public long TotalPayable { get; private set; }  // paisa
    public OrderStatus Status { get; private set; }
    public Guid? PaymentId { get; private set; }
    public string? PaymentGatewayRef { get; private set; }
    public Guid? PolicyId { get; private set; }

    public static Result<Order> Create(Guid quotationId, Guid customerId, long totalPayable)
    {
        if (totalPayable <= 0) return Result.Fail<Order>("INVALID_AMOUNT", "Amount must be positive");
        var o = new Order { QuotationId = quotationId, CustomerId = customerId,
                            TotalPayable = totalPayable, Status = OrderStatus.Pending };
        o.RaiseDomainEvent(new OrderCreatedEvent(o.Id));
        return Result.Ok(o);
    }

    public Result InitiatePayment(Guid paymentId, string gatewayRef)
    {
        if (Status != OrderStatus.Pending)
            return Result.Fail("INVALID_TRANSITION", $"Cannot initiate payment from {Status}");
        PaymentId = paymentId;
        PaymentGatewayRef = gatewayRef;
        Status = OrderStatus.PaymentInitiated;
        return Result.Ok();
    }

    public Result ConfirmPayment()
    {
        if (Status != OrderStatus.PaymentInitiated)
            return Result.Fail("INVALID_TRANSITION", $"Cannot confirm from {Status}");
        Status = OrderStatus.Paid;
        RaiseDomainEvent(new OrderPaymentConfirmedEvent(Id, QuotationId, CustomerId));
        return Result.Ok();
    }

    public Result LinkPolicy(Guid policyId)
    {
        if (Status != OrderStatus.Paid)
            return Result.Fail("INVALID_TRANSITION", "Order must be PAID to link policy");
        PolicyId = policyId;
        Status = OrderStatus.PolicyIssued;
        return Result.Ok();
    }

    public Result Cancel(string reason)
    {
        if (Status is OrderStatus.PolicyIssued or OrderStatus.Cancelled)
            return Result.Fail("INVALID_TRANSITION", $"Cannot cancel from {Status}");
        Status = OrderStatus.Cancelled;
        RaiseDomainEvent(new OrderCancelledEvent(Id, reason));
        return Result.Ok();
    }
}

public enum OrderStatus
{
    Pending = 1, PaymentInitiated = 2, Paid = 3,
    PolicyIssued = 4, Cancelled = 5, Failed = 6
}
```

**Step 2 — `InitiatePaymentCommandHandler`:**
Call `PaymentGrpcClient.CreatePaymentIntentAsync(order.TotalPayable, "BDT")` → store `paymentId` + `gatewayRef`.

**Step 3 — `ConfirmOrderCommandHandler`:**
On payment gateway callback → `order.ConfirmPayment()` → dispatch `IssuePolicyCommand` via MediatR.

---

### 6.4 Policy State Machine ★ Core

```
                  IssuePolicy()
                      │
                      ▼
            ┌──────────────────┐
            │ PENDING_PAYMENT  │
            └────────┬─────────┘
                     │ (payment confirmed, FR-037)
                     ▼
                ┌────────┐
         ┌──────│ ACTIVE │──────┬──────────┐
         │      └────────┘      │          │
    Cancel()   Suspend()    (coverage   (non-renewal
    │  ↑15-day    │        end date)   past grace)
    │  cooling    ▼           ▼          ▼
    │  off     ┌───────────┐ ┌────────┐ ┌────────┐
    │  period  │ SUSPENDED │ │EXPIRED/│ │ LAPSED │
    ▼  FR-038  └─────┬─────┘ │MATURED │ └───┬────┘
   ┌───────────┐     │       └────────┘     │
   │ CANCELLED │  Reinstate()  (terminal)  Reinstate()
   │ (terminal)│     │                  (within 90 days,
   └───────────┘     ▼                   medical UW req'd
   triggers   ┌────────────┐             FR-222)
   refund     │REINSTATED  │◄────────────┘
   calc       │  → ACTIVE  │
              └────────────┘
```

**Status:** ✅ gRPC contract done, ✅ `src/PoliSync.Policy/Domain/PolicyAggregate.cs` exists. Remaining gaps are SRS parity: strict policy-number sequencing, richer nominee/approval rules, and fuller event/data-gateway coverage.

**SRS Rules (FR-030 to FR-040, FR-084 to FR-097):**
- **Cooling-off:** 15-day window from issuance for full refund cancellation (FR-038)
- **Reinstatement:** Within 90 days of lapse, requires medical underwriting + Focal Person approval (FR-090/FR-222)
- **Duplicate detection:** Block same product + same NID within 30 days (FR-216)
- **Policy number:** `LBT-YYYY-XXXX-NNNNNN` (FR-034) — sequential, year-prefixed
- **Cancellation approval:** Policies >30 days old require Business Admin + Focal Person approval (FR-094)
- **NID uniqueness:** Validate across policies to prevent duplicate insurance (FR-033)

#### Implementation Steps

**Step 1 — Audit and extend the existing `src/PoliSync.Policy/Domain/PolicyAggregate.cs`**

Current branch reality:
- `PolicyAggregate.Create(...)` already exists and creates a proto-backed aggregate with `PendingPayment`
- `IssuePolicy()`, `CancelPolicy(reason)`, `SuspendPolicy()`, `ReinstatePolicy()`, `MarkAsLapsed()`, and `SetDocumentUrl()` already exist
- `IssuePolicyCommandHandler` already creates and issues policies through the aggregate

Remaining policy-domain gaps to close:
- current policy-number generation is randomised, not the strict sequential `LBT-YYYY-XXXX-NNNNNN` flow required by SRS/proto commentary
- cooling-off, duplicate-detection, older-policy cancellation approval, and nominee invariants are not yet enforced at the aggregate level
- document generation, notification, renewal, and commission follow-through still need tighter event/data-gateway orchestration

**Step 2 — Policy number generation** (in `IssuePolicyCommandHandler`):

```csharp
// SRS FR-034: format LBT-YYYY-XXXX-NNNNNN
// Proto check_constraint: policy_number ~ '^LBT-[0-9]{4}-[A-Z0-9]{4}-[0-9]{6}$'
var seq = await _gateway.GetNextSequenceAsync("insurance_schema.policy_number_seq");
var policyNumber = $"LBT-{DateTime.UtcNow.Year}-{seq:D6}";
```

**Step 3 — Pre-issuance checks** (in `IssuePolicyCommandHandler`):

```csharp
// 1. Verify quotation is APPROVED
// 2. Verify order is PAID
// 3. KYC check
var kycResult = await _kycClient.VerifyNID(nominee.NidNumber);
if (!kycResult.IsVerified) return Result.Fail("KYC_FAILED", "NID verification failed");
// 4. Fraud check (proto contract uses int32 0-100)
var fraudResult = await _fraudClient.CheckFraud(request);
if (fraudResult.FraudScore > 75) return Result.Fail("FRAUD_FLAGGED", "Policy flagged for investigation");
// 5. Issue
var policy = PolicyAggregate.Issue(...);
// 6. Persist via data gateway
await _gateway.CreatePolicyAsync(policy);
// 7. Post-issuance (async via domain events):
//    - DocgenGrpcClient.GeneratePolicyDocument()
//    - StorageGrpcClient.Upload()
//    - NotificationService → SMS/email
//    - RenewalScheduler → create schedule
//    - CommissionHandler → calculate payout
```

**Step 4 — Nominee with PII:**

```csharp
public sealed class Nominee : ValueObject
{
    public string Name { get; }
    public string Relationship { get; }
    public DateTime DateOfBirth { get; }
    public string NidNumber { get; }     // AES-256-GCM encrypted at rest
    public string PhoneNumber { get; }   // AES-256-GCM encrypted at rest
    public int Percentage { get; }       // Must sum to 100 across all nominees
    public bool IsPrimary { get; }
}
```

---

### 6.5 Endorsement State Machine

```
  CreateEndorsement(policyId, type, changes)
        │
        ▼
  ┌───────────┐
  │ REQUESTED │
  └─────┬─────┘
        │ (auto or manual routing)
        ▼
┌──────────────────┐
│ PENDING_APPROVAL │
└───────┬──────────┘
        ├──► APPROVED ──► ApplyEndorsement() ──► APPLIED (terminal)
        │                  │
        │                  └─ atomic: update policy + increment version
        │                             + generate new document
        └──► REJECTED (terminal)
```

**Status:** ✅ gRPC done, ❌ endorsement domain logic is still largely missing on the current branch

**Types:** `ADDRESS_CHANGE`, `NOMINEE_CHANGE`, `SUM_CHANGE`, `PLAN_CHANGE`, `RIDER_ADD`, `RIDER_REMOVE`, `COVERAGE_DATE_CHANGE`

#### Key Business Rules (SRS FR-098 to FR-102)

- Changes stored as JSONB diff (`{ "old": {...}, "new": {...} }`)
- `SUM_CHANGE` increasing **>10%** → requires approval (FR-102) + re-underwriting if health-related
- Sum increase → additional premium calculated mid-term (FR-099)
- Sum decrease → pro-rata refund credited to premium account (FR-100)
- Endorsement doc suffix: `POL-001/END-01`, `POL-001/END-02`, etc. (FR-101)
- `ApplyEndorsement()` must be **atomic** with policy update (single transaction via data gateway)
- Policy `version` incremented on apply; new document generated via DocgenGrpcClient

#### Implementation Steps

1. Create `src/PoliSync.Endorsement/Domain/Endorsement.cs` with state machine methods: `Create()`, `Approve()`, `Reject()`, `Apply()`
2. Create `ApplyEndorsementCommandHandler` — call data gateway in single transaction: update endorsement status + update policy fields + increment version
3. Wire `DocgenGrpcClient.GeneratePolicyDocument()` in `EndorsementAppliedEventHandler`
4. Kafka: publish `insuretech.endorsement.applied.v1` after successful apply

---

### 6.6 Renewal & Grace Period State Machine

```
  PolicyIssuedEvent
        │
        ▼
  Create RenewalSchedule (status=PENDING)
        │
   ┌────┴─────────────────────────────────┐
   │ RenewalSchedulerService (daily cron) │
   └─────┬──────────────┬────────────┬────┘
    T-60 days       T-30 days    T-7 days
    1st reminder    2nd reminder  final reminder
         │               │            │
         ▼               ▼            ▼
    Send SMS+Email  Send SMS+Email  Send SMS+Email+Push
    Status=NOTIFIED
         │
         │ (customer pays renewal)
         ├──► Status = RENEWED ──► new Policy issued
         │
         │ (T=0, coverage ends, no payment)
         ▼
    ActivateGracePeriod (30 days)
    GracePeriod status = ACTIVE
         │
         │ (customer pays during grace)
         ├──► GracePeriod = CLEARED → Policy renewed
         │
         │ (T+30, grace expires)
         ▼
    AutoLapsePolicy()
    Policy.Lapse() → LAPSED
    GracePeriod = EXPIRED
```

**Status:** ✅ gRPC done, ❌ renewal scheduler/domain logic still not implemented on the current branch

#### Implementation Steps

1. `RenewalSchedule`, `RenewalReminder`, `GracePeriod` aggregates
2. `RenewalSchedulerService` : `IHostedService`
   ```csharp
   // Cron: daily at 00:05 UTC
   // Query: SELECT * FROM renewal_schedules WHERE renewal_date IN (now+60, now+30, now+7)
   // For each: dispatch SendRenewalReminderCommand
   // Query: SELECT * FROM grace_periods WHERE end_date < now AND status = 'ACTIVE'
   // For each: dispatch AutoLapsePolicyCommand
   ```
3. `SendRenewalReminderCommandHandler` → `NotificationGrpcClient.SendSMS/Email/Push()`
4. `ActivateGracePeriodCommand` — creates `GracePeriod(startDate=coverage_end, endDate=coverage_end+30, status=ACTIVE)`
5. `AutoLapsePolicyCommandHandler` — calls `PolicyAggregate.Lapse()` via data gateway; publishes `insuretech.policy.lapsed.v1`

**Configuration:**
```json
{ "ReminderDays": [60, 30, 7], "GracePeriodDays": 30, "SchedulerCronUtc": "5 0 * * *" }
```

---

### 6.7 Underwriting Decision Engine

```
  SubmitHealthDeclaration
        │
        ▼
  Calculate RiskScore
  ┌─────────────────────────────────────────────────────┐
  │ Base: 50                                            │
  │ + Age:   18-35(+0), 36-50(+10), 51-65(+20), 66+(+30)│
  │ + BMI:   18.5-24.9(+0), 25-29.9(+5), 30+(+15)     │
  │ + Smoker: Yes(+15), No(+0)                         │
  │ + Pre-existing: each +10 (cap +30)                  │
  │ + Family history: each +8 (cap +16)                 │
  └─────────────────────────────────────────────────────┘
        │
   Score 0–40  → APPROVED (loading = 0%)
   Score 41–60 → APPROVED_WITH_LOADING (loading = 10–25%)
   Score 61–75 → REFERRED (manual underwriter review)
   Score 76+   → DECLINED
        │
        ▼
  UnderwritingDecision persisted
  Loading factor applied back to Quotation premium
```

**Status:** ✅ gRPC + data gateway done. Approval applies loading to quotation. Rejection updates quotation status. `HealthDeclarationAggregate` and `UnderwritingDecisionAggregate` are implemented and wired into `UnderwritingGrpcService`.

**Remaining:**
- [x] `RiskScorer` domain service — scoring algorithm implemented per formula above (`UnderwritingRiskScorer`)
- [x] `HealthDeclaration` aggregate with validation (`HealthDeclarationAggregate`)
- [x] `UnderwritingDecision` aggregate storing decision + score + loading (`UnderwritingDecisionAggregate`)
- [x] Unit tests: all score combinations × decision thresholds (20+ test cases) (`UnderwritingRiskScorerTests`)

---

### 6.8 Claim Approval Matrix

```
  FileClaim (FNOL)
        │
        ▼
     ┌───────┐
     │ FILED │
     └───┬───┘
         │ FraudGrpcClient.CheckClaim()
         ▼
   ┌────────────┐
   │ FRAUD_CHECK│
   └──────┬─────┘
     fraud_score > 75? → FLAGGED_FOR_INVESTIGATION
     fraud_score ≤ 75:
          │
          │ Route by claimed_amount (paisa):
          │
    ┌─────┴──────────────────────────────────────────┐
    │ BDT 0–10K:  L1 Auto/Officer (24hr TAT)        │
    │   ZHTC auto-approve IF fraud_score < 30       │
    ├────────────────────────────────────────────────┤
    │ BDT 10K–50K: L2 Claims Manager (3-day TAT)    │
    ├────────────────────────────────────────────────┤
    │ BDT 50K–2L:  L3 Business Admin + Focal Person │
    │              JOINT approval (7-day TAT)        │
    ├────────────────────────────────────────────────┤
    │ BDT 2L+:     Board + Insurer Approval          │
    │              (15-day TAT)                       │
    └────────────────────────────────────────────────┘
    Claim types: CASHLESS (hospital network) / REIMBURSEMENT
    Co-pay/Deductible: (Claim - Deductible) × Co-pay% (FR-104)
    Claim number format: CLM-YYYY-XXXX-NNNNNN (FR-043)
         │
         ▼
  APPROVED → Settlement → PaymentGrpcClient.InitiateSettlement()
      or
  REJECTED → NotificationGrpcClient
```

**Status:** ✅ gRPC service exists and `src/PoliSync.Claims/Domain/ClaimAggregate.cs` exists. Remaining gaps are consistent aggregate usage, proto-aligned fraud-score handling, and full approval/workflow coverage.

**Configuration:** (from `appsettings.json`)
```json
{
  "AutoApproveThresholdPaisa": 1000000,
  "AutoApproveFraudScoreMax": 30,
  "FraudFlagThreshold": 75,
  "L1ThresholdPaisa": 5000000,
  "L2ThresholdPaisa": 20000000,
  "L3ThresholdPaisa": 50000000
}
```

> ⚠️ **Proto note:** `FraudService.CheckFraud` returns `fraud_score` as **int32 (0-100)**, not float. The current branch still uses `double` `0.30` / `0.75` style thresholds in `ClaimAggregate` and `appsettings.json`; that is an implementation drift that should be normalized to integer `30` / `75`.
```

#### Implementation Steps

1. Use the existing `ClaimAggregate` as the canonical domain object, then extend it where behavior is still bypassed or proto-misaligned.

Current branch reality:
- `ClaimAggregate.FileClaim(...)` exists
- `ApplyFraudCheck(...)` exists
- `GetRequiredApprovalLevel()` exists
- `AddApproval(...)` exists
- `Settle(...)` exists
- `AddDocument(...)` exists

Remaining gaps:
- normalize fraud-score scale to match fraud proto (`int32 0-100`)
- ensure gRPC handlers consistently use aggregate methods instead of directly mutating proto entities where business rules should live
- finish approval-matrix, escalation, and settlement orchestration around the aggregate rather than partial inline handler logic

**Proto enums (source of truth):**
```
ClaimStatus: SUBMITTED, UNDER_REVIEW, PENDING_DOCUMENTS, APPROVED, REJECTED, SETTLED, DISPUTED
ClaimType: HEALTH_HOSPITALIZATION, HEALTH_SURGERY, MOTOR_ACCIDENT, MOTOR_THEFT, TRAVEL_MEDICAL, TRAVEL_BAGGAGE_LOSS, DEVICE_DAMAGE, DEVICE_THEFT, DEATH
ClaimProcessingType: MANUAL, AUTO_ADJUDICATED, AI_ASSISTED
ApprovalDecision: PENDING, APPROVED, REJECTED, NEEDS_MORE_INFO
```

2. `ApprovalRouter` domain service — determines approval level from `claimed_amount`:
   ```csharp
   public ApprovalLevel Route(long claimedAmountPaisa, int fraudScore, ClaimsConfig config)
   {
       if (fraudScore > config.FraudFlagThreshold) return ApprovalLevel.FlaggedForInvestigation;
       if (claimedAmountPaisa <= config.AutoApproveThresholdPaisa && fraudScore < config.AutoApproveFraudScoreMax)
           return ApprovalLevel.AutoApprove;
       if (claimedAmountPaisa <= config.L1ThresholdPaisa) return ApprovalLevel.L1;
       if (claimedAmountPaisa <= config.L2ThresholdPaisa) return ApprovalLevel.L2;
       if (claimedAmountPaisa <= config.L3ThresholdPaisa) return ApprovalLevel.L3;
       return ApprovalLevel.Board;
   }
   ```
3. `FileClaimCommandHandler`: create claim → call `FraudGrpcClient.CheckClaim()` → set `fraud_score` + `fraud_flags` → route to approver
4. For L2+ claims: `WorkflowGrpcClient.CreateTask()` to assign to specific reviewer
5. `SettleClaimCommandHandler`: `PaymentGrpcClient.InitiateSettlement(amount, bankAccount/bKashNumber)` — within 7-15 working days (FR-105)

**Fraud Detection Rules (SRS FR-186 to FR-192):**

| Rule | Trigger | Action |
|------|---------|--------|
| FD-001 Rapid Claim | Claim < 48hrs after purchase | Auto-flag + manual review |
| FD-002 Frequent | Same claim type >2x in 12 months | Flag + pattern analysis |
| FD-003 Max Amount | Claim = 100% of coverage | Flag + enhanced verification |
| FD-004 Non-Network | Medical provider not in approved list | Flag + provider verification |
| FD-005 Geographic | Claim >100km from registered address | Flag + location verification |
| FD-006 Device | Multiple accounts from same device (>3) | Flag + identity verification |
| FD-007 Behavioral | ML-based unusual activity patterns | Risk scoring + monitoring |

---

### 6.9 Commission Calculation

```
  PolicyIssuedEvent (Kafka consumer)
        │
        ▼
  Load CommissionConfig for (partner_id, product_id, plan_id)
        │
        ├── FLAT:       commission = flat_amount
        ├── PERCENTAGE: commission = total_payable × rate / 10000 (basis points)
        └── TIERED:     find tier by total_payable range → apply tier rate
        │
        ▼
  tax_deduction = commission × 0.10  (10% BD withholding tax)
  net_payout = commission - tax_deduction
        │
        ▼
  Create CommissionPayout (status=PENDING)
  Publish insuretech.commission.payout_created.v1
```

**Monthly revenue share cron:**
```
  SUM all total_payable per partner for period
  platform_share = total × 0.20 (20%)
  partner_share = total - platform_share - tax
  Create RevenueShare (status=CALCULATED)
```

---

### 6.10 Refund Calculation

```
  PolicyCancelledEvent (Kafka consumer)
        │
        ▼
  days_remaining = coverage_end - cancellation_date
  days_total     = coverage_end - coverage_start
  pro_rata       = total_paid × days_remaining / days_total
  penalty        = pro_rata × 0.10  (10% cancellation penalty)
  net_refund     = pro_rata - penalty
  // Within 15-day cooling-off (FR-038): full refund, no penalty
        │
        ▼
  Create Refund (status=CALCULATED)
  → ApproveRefund → PaymentGrpcClient.InitiateRefund()
  → Refund within 7 working days via MFS or bank (FR-096)
  → Publish insuretech.refund.initiated.v1
```

---

## 8. Go Services Consumed by C# (with exact RPCs)

These Go services are **fully implemented or will be** — C# just calls them via gRPC:

### InsuranceService (:50115) — Data Layer [75+ RPCs]

| CRUD Group | RPCs |
|-----------|------|
| Product | Create, Get, Update, Delete, List |
| ProductPlan | Create, Get, List |
| Rider | Create, Get, List |
| PricingConfig | Create, Get |
| Policy | Create, Get, Update, Delete, List |
| Claim | Create, Get, Update, Delete, List |
| Quote | Create, Get, Update, Delete, List |
| UnderwritingDecision | Create, Get, Update, Delete, List |
| HealthDeclaration | Create, Get, Update, Delete, GetByQuote |
| RenewalSchedule | Create, Get, Update, Delete, List |
| RenewalReminder | Create, Get, Update, Delete, List |
| GracePeriod | Create, Get, Update, Delete, GetByPolicy, ListActive |
| Insurer | Create, Get, Update, Delete, List |
| InsurerConfig | Create, Get, Update, Delete, GetByInsurer |
| InsurerProduct | Create, Get, Update, Delete, List |
| FraudRule/Case/Alert | Full CRUD + ListActive |
| Beneficiary | Create, Get, Update, Delete, List (Individual + Business) |
| Endorsement | Create, Get, Update, Delete, ListByPolicy |
| Quotation | Create, Get, Update, Delete, List |
| PolicyServiceRequest | Create, Get, Update, Delete, ListByPolicy |
| ServiceProvider | Create, Get, Update, Delete, List |

### FraudService (:50220) — Scoring [10 RPCs]
`CheckFraud`(entity_type, entity_id, data → fraud_score **int32 0-100**, risk_level, triggered_rules), GetAlert, ListAlerts, CreateCase, GetCase, UpdateCase, ListRules, CreateRule, ActivateRule, DeactivateRule

### PaymentService (:50190) — Payments [7 RPCs]
`InitiatePayment`(user_id, policy_id, amount, currency, payment_method [BKASH/NAGAD/ROCKET/CARD], callback_url, **idempotency_key**), `VerifyPayment`, `GetPayment`, `ListPayments`, `InitiateRefund`(payment_id, refund_amount, reason), `GetRefundStatus`, `ReconcilePayments`

### OrdersService (:50142) — Order Lifecycle Persistence [7 RPCs]
`CreateOrder`, `GetOrder`, `ListOrders`, `InitiatePayment`, `ConfirmPayment`, `CancelOrder`, `GetOrderStatus`

### NotificationService (:50230) — Messaging [9 RPCs]
`SendNotification`(recipient_id, type, channel, subject, message, template_id, template_data, priority, schedule_after_seconds), `SendBulkNotifications`, `GetNotificationStatus`, `GetUserNotifications`, `MarkAsRead`, `UpdatePreferences`, CreateTemplate, UpdateTemplate, DeactivateTemplate

### WorkflowService (:50180) — Approval Routing [7 RPCs]
`StartWorkflow`(definition_id, entity_type, entity_id, context), `CompleteTask`(task_id, decision [APPROVED/REJECTED/RETURNED], comments), `GetMyTasks`, GetWorkflowInstance, GetWorkflowHistory, CreateDefinition, GetDefinition

### PartnerService (:50100) — Partner/Agent [11 RPCs]
Create, Get, Update, Delete, ListPartners, VerifyPartner, UpdatePartnerStatus, `GetPartnerCommission`(partner_id, date_range), `UpdateCommissionStructure`, GetAPICredentials, RotateAPIKey

### Other Go Services (C# consumes as-is)
| Service | Port | Used For |
|---------|------|----------|
| DocgenService | :50280 | `GenerateDocument`(template, data) → PDF URL |
| StorageService | :50290 | Upload/Download/Delete blobs |
| AuditService | :50080 | `CreateAuditEvent`(entity, action, details) |
| KycService | :50090 | `VerifyIdentity`(nid_number) → verified/rejected |
| SupportService | :50320 | Ticket management |

---

## 9. Kafka Topics

Current implementation note: the Go `orders-service` still uses non-versioned topics for the order→policy projection path. PoliSync now consumes `orders.order.payment_confirmed` and publishes `policy.issued` for the Go order projector.

| Topic | Producer | Consumers | Status |
|-------|----------|-----------|--------|
| `insuretech.product.created.v1` | product-svc | analytics | 🟡 Wired |
| `insuretech.product.activated.v1` | product-svc | cache-invalidate | 🟡 Wired |
| `insuretech.product.deactivated.v1` | product-svc | cache-invalidate | 🟡 Wired |
| `insuretech.product.pricing_updated.v1` | product-svc | quote-svc | 🟡 Wired |
| `insuretech.quotation.submitted.v1` | quote-svc | underwriting-svc | 🟡 Wired (consumer added) |
| `insuretech.quotation.approved.v1` | quote-svc | notification | ❌ |
| `orders.order.payment_confirmed` | go order-svc | polisync policy consumer | ✅ Wired |
| `policy.issued` | polisync policy consumer | go order-svc | ✅ Wired |
| `insuretech.policy.cancelled.v1` | policy-svc | refund, notification | ❌ |
| `insuretech.policy.lapsed.v1` | renewal-svc | notification | ❌ |
| `insuretech.endorsement.applied.v1` | endorsement-svc | notification, docgen | ❌ |
| `insuretech.renewal.reminder_sent.v1` | renewal-svc | notification | ❌ |
| `insuretech.underwriting.decision_made.v1` | underwriting-svc | quote-svc | 🟡 Wired (publisher added) |
| `insuretech.claim.filed.v1` | claim-svc | fraud-svc, notification | ❌ |
| `insuretech.claim.approved.v1` | claim-svc | notification, settlement | ❌ |
| `insuretech.claim.settled.v1` | claim-svc | analytics, audit | ❌ |
| `insuretech.commission.payout_created.v1` | commission-svc | notification, payment | ❌ |
| `insuretech.refund.initiated.v1` | refund-svc | payment, notification | ❌ |

---

## 10. Database Schema

Schema DDL managed by Go migrations. C# uses Go insurance-service for DB ops.

### `insurance_schema`
`products`, `product_plans`, `pricing_configs`, `riders`, `quotations`, `orders`, `policies`, `nominees`*(PII)*, `policy_riders`, `endorsements`, `renewal_schedules`, `renewal_reminders`, `grace_periods`, `health_declarations`, `underwriting_decisions`, `claims`, `claim_documents`, `claim_approvals`, `settlements`, `refunds`

### `commission_schema`
`commission_configs`, `commission_payouts`, `revenue_shares`

---

## 11. Configuration Reference

```json
{
  "Cache": { "ProductTtlSeconds": 300, "QuotationTtlSeconds": 3600 },
  "Commission": { "WithholdingTaxRate": 0.10, "CancellationPenaltyRate": 0.10, "PlatformRevenueShareRate": 0.20 },
  "Renewal": { "ReminderDays": [60, 30, 7], "GracePeriodDays": 30, "SchedulerCronUtc": "5 0 * * *" },
  "Claims": {
    "AutoApproveThresholdPaisa": 1000000, "AutoApproveFraudScoreMax": 30,
    "FraudFlagThreshold": 75, "L1ThresholdPaisa": 5000000,
    "L2ThresholdPaisa": 20000000, "L3ThresholdPaisa": 50000000
  },
  "Quotation": { "ExpiryDays": 30 },
  "GrpcClients": {
    "InsuranceService": "http://insurance-service:50115",
    "OrdersService": "http://orders-service:50142",
    "AuthzService": "http://authz-service:50070",
    "FraudService": "http://fraud-service:50220",
    "PaymentService": "http://payment-service:50190",
    "NotificationService": "http://notification-service:50230",
    "DocgenService": "http://docgen-service:50280",
    "StorageService": "http://storage-service:50290",
    "WorkflowService": "http://workflow-service:50180",
    "AuditService": "http://audit-service:50080",
    "KycService": "http://kyc-service:50090"
  },
  "Jwt": { "Issuer": "insuretech-authn", "Audience": "insuretech-platform" }
}
```

---

## 12. Security & Compliance (SRS SEC + IDRA + BFIU)

| Req | Implementation | Status |
|-----|---------------|--------|
| JWT verification | `JwtAuthInterceptor` | ✅ |
| RBAC | `AuthzGrpcClient` | 🔴 Stub |
| ABAC partner isolation | `ICurrentUser.partner_id` | ❌ |
| Tenant isolation | `tenant_id` scoped queries | ❌ |
| PII encryption | AES-256-GCM (NID, phone) | ✅ Infra ready |
| PII masking | NID last-3, phone mid-mask, email user-mask (SEC-002) | ❌ |
| Audit trail | `AuditGrpcClient.Log()` | 🔴 Stub |
| KYC on issuance | `KycGrpcClient.VerifyNID()` | 🔴 Stub |
| Fraud screening | `FraudGrpcClient.Check*()` | 🔴 Stub |
| File virus scan | ClamAV on uploads (SEC-010) | ❌ |
| API rate limiting | 1000/hr auth, 100/hr anon (SEC-021) | ❌ |
| Data retention | 20 years policies/claims, 7 years audit | ❌ |

### AML/CFT Compliance (BFIU)

20+ automated monitoring rules including:
- **TM-001** Structuring: 3+ transactions of 9K-10K BDT in 7 days
- **TM-002** Rapid claim within 7 days of purchase
- **TM-008** >3 policies in 7 days → Enhanced Due Diligence
- **TM-009** Premium >BDT 5 lakh → management approval
- **TM-010** >2 cancellations in 3 months → investigation
- **SAR workflow:** Auto-flag → CO review 24hrs → escalate → BFIU within 3 business days → no tipping off

### IDRA Reporting

| Report | Frequency | Due By |
|--------|-----------|--------|
| IC-1 (Premium Collection) | Monthly | 10th of month |
| IC-2 (Claims Intimation) | Monthly | 10th of month |
| IC-3 (Claims Settlement) | Quarterly | Q+15 days |
| IC-4 (Financial Performance) | Quarterly | Q+20 days |
| FCR + CARAMELS | Annual | Y+90 days |
| Incident >BDT 1L, breach, outage >4hr | Event-based | 48 hrs |

### Data Retention (SRS)

| Data | Hot (Postgres) | Warm (S3) | Cold (Glacier) | Total |
|------|---------------|-----------|----------------|-------|
| Active Policies | Lifetime | — | — | Lifetime |
| Expired Policies | 1 year | 5 years | 20 years | 20 years |
| Claims | 2 years | 5 years | 20 years | 20 years |
| Audit Logs | 90 days | 1 year | 7 years | 7 years |
| Sessions | Redis 7 days | — | — | 7 days |
| STR/SAR docs | 10+ years secured offline | — | — | 10+ years |

---

## 13. Performance Targets

| Endpoint | p95 | Strategy |
|----------|-----|----------|
| `CalculatePremium` | < 200ms | Stateless, cache-first |
| `GetProduct` | < 100ms | Redis 300s TTL |
| `GetPolicy` | < 150ms | Via insurance-service |
| `IssuePolicy` (sync) | < 500ms | Async PDF gen |
| `FileClaim` | < 300ms | Sync FNOL + fraud |
| `ListPolicies` | < 400ms | Keyset pagination |

**Throughput:** 1,000 concurrent quotes · 500 policies/min · 200 FNOL/hr

---

## 14. Build & Run

```bash
# Start Go insurance-service first
cd backend/inscore && go run cmd/service-manager/main.go start insurance

# Build & run PoliSync
cd backend/polisync
dotnet restore && dotnet build
dotnet watch run --project src/PoliSync.ApiHost

# Verify
curl http://localhost:50121/health
grpcurl -plaintext localhost:50120 list

# Env vars
$env:DB_PASSWORD="..."
$env:REDIS_PASSWORD="..."
$env:PII_ENCRYPTION_KEY=$(openssl rand -base64 32)
```

---

## 15. Proto Changes Required for SRS Compliance

> Gap analysis: current proto definitions vs SRS v3.7 + InsuranceKnowledgeBank requirements.
> Proto files are at `E:\Projects\InsureTech\proto\insuretech\`

### P0 — Blocking (must fix before C# domain logic implementation)

| # | Proto File | Change | SRS Ref | Rationale |
|---|-----------|--------|---------|-----------|
| 1 | `policy/entity/v1/policy.proto` → `PolicyStatus` enum | Add `POLICY_STATUS_REINSTATED = 8` and `POLICY_STATUS_MATURED = 9` | FR-090, FR-222, FR-084 | SRS reinstatement flow sets status to REINSTATED (not ACTIVE) until underwriting review completes. MATURED needed for life policies reaching term. C# state machine needs these as valid targets. |
| 2 | `policy/entity/v1/quotation.proto` → `QuotationStatus` enum | Add `QUOTATION_STATUS_EXPIRED = 6` | FR-031 | Quotations expire after 30 days. Currently no EXPIRED state — C# expiry job can't set valid status. |
| 3 | `claims/entity/v1/claim.proto` → `ClaimType` enum | Add `CLAIM_TYPE_CASHLESS = 10` and `CLAIM_TYPE_REIMBURSEMENT = 11` — **OR** add separate `ClaimMode` enum | FR-105 | SRS distinguishes cashless (provider pre-auth) vs reimbursement. Current proto only tracks domain types (HEALTH_HOSPITALIZATION etc.), not the processing/settlement mode. Suggest: add a new `ClaimMode` enum + field `claim_mode = 32` rather than polluting `ClaimType`. |
| 4 | `policy/entity/v1/policy.proto` → `Policy` message | Add fields: `cooling_off_end_date` (DATE), `cancellation_reason` (TEXT), `cancellation_approved_by` (UUID), `reinstated_at` (TIMESTAMPTZ), `reinstatement_reason` (TEXT), `duplicate_check_hash` (VARCHAR(64) — SHA-256 of product_id+customer_nid) | FR-038, FR-094, FR-222, FR-216 | **cooling_off_end_date**: 15-day window start_date+15d, enables C# to check full-refund eligibility. **cancellation fields**: SRS requires reason + approver for policies >30 days. **reinstatement fields**: audit trail for reinstatement. **duplicate_check_hash**: unique index prevents same customer NID + product combo within 30 days (FR-216). |

### P1 — Early Phase (needed before Phase 3-5 implementation)

| # | Proto File | Change | SRS Ref | Rationale |
|---|-----------|--------|---------|-----------|
| 5 | `products/entity/v1/product.proto` → `ProductCategory` enum | Add 14 missing categories to match `InsuranceType` enum in `common/v1/types.proto`: FIRE, MARINE, PROPERTY, LIABILITY, LOSS_OF_INCOME, GOODS_IN_TRANSIT, CROPS, LIVESTOCK, CREDIT, SURETY, ENGINEERING, AVIATION, PET, CYBER | SRS Non-Life categories, IDRA classification | Current `ProductCategory` has 7 values (MOTOR, HEALTH, TRAVEL, HOME, DEVICE, AGRICULTURAL, LIFE). But `common.v1.InsuranceType` already has 21 values. **Option A:** Deprecate `ProductCategory`, use `InsuranceType` directly on Product (preferred — single source of truth). **Option B:** Mirror all 21 values into `ProductCategory`. |
| 6 | `endorsement/entity/v1/endorsement.proto` → `Endorsement` message | Add check_constraint on `endorsement_number`: `'^END-[0-9]{4}-[A-Z0-9]{4}-[0-9]{6}-[0-9]{2}$'` matching policy doc suffix `END-01` format | FR-099 | Policy number `LBT-YYYY-XXXX-NNNNNN` has check_constraint; endorsement number should match the SRS document suffix convention. |
| 7 | `claims/entity/v1/claim.proto` → `Claim` message | Add fields: `co_pay_percentage` (DECIMAL 5,2), `network_provider_id` (UUID FK to service_providers), `claim_mode` (new enum CASHLESS/REIMBURSEMENT — see P0 #3) | FR-104, FR-105 | `co_pay_percentage` stores the plan-level co-pay rate (separate from calculated `co_pay_amount`). `network_provider_id` links to approved provider list for FD-004 geographic/network validation. |
| 8 | `insurance/services/v1/insurance_service.proto` | Add composite RPCs: `FindDuplicatePolicy(product_id, customer_nid_hash)`, `GetPoliciesByCustomerNID(nid_hash)`, `ListExpiringPolicies(days_ahead, page)` | FR-216, FR-033, FR-087 | Currently only generic CRUD. C# needs duplicate detection, NID-based lookup for uniqueness validation, and batch expiry queries for the renewal cron. Without these, C# must do multiple CRUD calls + client-side filtering. |

### P2 — Compliance & Reporting (needed before Phase 7-8)

| # | Proto File | Change | SRS Ref | Rationale |
|---|-----------|--------|---------|-----------|
| 9 | `renewal/entity/v1/renewal_schedule.proto` → `RenewalSchedule` | Add `reminder_days_before` (INT[] default `{30,15,7,1}`) | FR-087 | SRS specifies 4-tier reminder schedule. Currently no field to configure per-schedule reminder cadence — hardcoding in C# is brittle. |
| 10 | `policy/entity/v1/policy.proto` → `PolicyStatus` | Consider adding `POLICY_STATUS_PENDING_REINSTATEMENT = 10` | FR-222 | Reinstatement requires medical UW + Focal Person approval. An intermediate status prevents premature claims on policies under review. Without it, the policy goes directly from LAPSED to REINSTATED with no auditable in-between. |
| 11 | NEW: `kyc/services/v1/kyc_service.proto` | Add `BatchVerifyNID` RPC and `GetVerificationStatus` | FR-033, SEC-005 | Current KYC proto (if exists) likely only has single-NID verification. B2B bulk employee upload needs batch NID verification. |
| 12 | NEW: `report/entity/v1/idra_report.proto` | Define `IDRAReport` entity (report_type [IC-1..IC-4, FCR], period, status, generated_at, file_url, submitted_at) + `GenerateIDRAReport` RPC | SEC-011 to SEC-014 | IDRA monthly/quarterly/annual reporting has structured output requirements. A proto entity ensures consistent report generation and audit trail. |

### Summary of Impact

```
Proto changes needed:  12 items
  P0 (blocking):       4  → Must fix before C# state machine coding
  P1 (early phase):    4  → Before Phase 3-5
  P2 (compliance):     4  → Before Phase 7-8

Affected proto files:  6 existing + 2 new
  policy.proto         → 2 enum additions + 6 new fields
  quotation.proto      → 1 enum addition
  claim.proto          → 1 new enum + 3 new fields
  product.proto        → Align with InsuranceType or add 14 values
  endorsement.proto    → 1 check_constraint
  insurance_service.proto → 3 new composite RPCs
  kyc_service.proto    → 1 new RPC (BatchVerifyNID)
  idra_report.proto    → NEW entity + RPC

Migration impact: All changes are additive (new enum values, new fields,
  new messages). No breaking changes to existing deployed services.
  Go services need regenerated stubs + migration for new columns.
```
