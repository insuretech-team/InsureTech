# Insurance Engine — Complete Technical Development Blueprint
> **LabAid InsureTech Platform | SRS v3.11 | C# .NET 8 | Port 5001**
> Generated from: SRS v3.11 + Phase 1/2/3 Analysis Scripts
> Architecture: Vertical Slice Architecture (VSA) + CQRS/MediatR
> Communication: gRPC only (no REST)

---

# 1. System Understanding

## 1.1 Business Objective

LabAid InsureTech Platform হলো Bangladesh-focused micro-insurance ecosystem যা digital channel-এ end-to-end insurance lifecycle manage করে। Target: BDT 200-2,000 premium range-এর micro-insurance products, 165 million mobile-using Bangladeshis-এর কাছে পৌঁছানো।

## 1.2 Insurance Engine-এর Exact Role

Insurance Engine (C# .NET 8, Port 5001) হলো platform-এর **central business logic hub**। এটি সমস্ত insurance-specific domain logic contain করে:

| দায়িত্ব | FG Group |
|---|---|
| Product Catalog Management | FG-003 |
| Policy Lifecycle (Purchase → Expiry) | FG-004 |
| Renewals, Grace Period, Lapse | FG-005 |
| Cancellation & Pro-rata Refund | FG-005.1 |
| Endorsement & Amendment | FG-005.2 |
| Business Rules & Workflows | FG-006 |
| Claims Management | FG-008 |
| Fraud Detection Rules | FG-016 |

**Scope boundary (critical):** Reporting/Analytics → Port 5003 (Analytics Service)। Partner Management → Port 5002। Payment → Node.js Port 3001।

## 1.3 System-এ Insurance Engine-এর Position

```
[Clients: Mobile/Web/Portal]
         ↓ REST/GraphQL
[API Gateway (Go:8080)]  ←→  [Auth Service (Go:8081/8082)]
         ↓ gRPC (Category 1)
[Insurance Engine (C#:5001)]  ←→  [Kafka (Go:8086)]
         ↓ PostgreSQL         ↓ Redis         ↓ S3
[PostgreSQL 17]          [Redis 7.0]     [Object Storage]
```

## 1.4 Upstream Dependencies (যেখান থেকে data আসে)

| Service | কোথায় | কীভাবে | কী আসে |
|---|---|---|---|
| API Gateway | Go:8080 | gRPC call + JWT validation metadata | User identity, tenant context |
| Auth Service | Go:8081 | gRPC interceptor token validation | user_id, roles, permissions |
| Authorization | Go:8082 | gRPC policy check | RBAC/ABAC decision |
| Payment Service | Node.js:3001 | Kafka event: `payment.completed` | Payment confirmation |
| Partner Service | C#:5002 | gRPC call | Partner/agent info |
| Storage Service | Go:8084 | gRPC | Document URLs after upload |
| AI Engine | Python:4001 | gRPC | Fraud risk score |
| OCR Service | Python:4002 | gRPC | Extracted document fields |

## 1.5 Downstream Dependencies (যেখানে data যায়)

| Event/Call | Destination | Format |
|---|---|---|
| `policy.issued` | Kafka → Notification Service | Protobuf event |
| `policy.issued` | Kafka → Analytics Service (5003) | Protobuf event |
| `policy.issued` | Kafka → Partner Service (5002, commission) | Protobuf event |
| `claim.submitted` | Kafka → AI Engine (fraud check) | Protobuf event |
| `claim.approved` | Kafka → Payment Service | Protobuf event |
| `document.generate` | OCR/PDF Service (4002) via gRPC | PDF generation request |
| `renewal.reminder` | Kafka → Notification Service | Protobuf event |

---

# 2. Full Insurance Engine Architecture

## 2.1 Project Structure (Vertical Slice Architecture)

```
InsurancEngine/                          ← Solution Root
├── InsuranceEngine.API/                 ← gRPC Server entry point
│   ├── Program.cs
│   ├── appsettings.json
│   ├── GrpcServices/                    ← gRPC service registrations
│   │   ├── ProductGrpcService.cs
│   │   ├── PolicyGrpcService.cs
│   │   ├── ClaimsGrpcService.cs
│   │   └── EndorsementGrpcService.cs
│   └── Interceptors/
│       ├── AuthInterceptor.cs
│       ├── LoggingInterceptor.cs
│       └── ExceptionInterceptor.cs
│
├── InsuranceEngine.Features/            ← ALL business logic (VSA)
│   ├── Products/
│   │   ├── Commands/
│   │   │   ├── CreateProduct/
│   │   │   │   ├── CreateProductCommand.cs
│   │   │   │   ├── CreateProductHandler.cs
│   │   │   │   └── CreateProductValidator.cs
│   │   │   └── UpdateProduct/
│   │   ├── Queries/
│   │   │   ├── GetProduct/
│   │   │   │   ├── GetProductQuery.cs
│   │   │   │   └── GetProductHandler.cs
│   │   │   └── SearchProducts/
│   │   ├── Domain/
│   │   │   ├── Product.cs               ← Entity
│   │   │   ├── ProductCategory.cs       ← Value Object/Enum
│   │   │   └── ProductRules.cs          ← Domain Rules
│   │   └── Infrastructure/
│   │       ├── ProductRepository.cs
│   │       └── ProductCache.cs
│   │
│   ├── Policies/
│   │   ├── Commands/
│   │   │   ├── CreatePolicy/
│   │   │   │   ├── CreatePolicyCommand.cs
│   │   │   │   ├── CreatePolicyHandler.cs   ← orchestrator
│   │   │   │   └── CreatePolicyValidator.cs
│   │   │   ├── ActivatePolicy/
│   │   │   ├── CancelPolicy/
│   │   │   ├── RenewPolicy/
│   │   │   ├── LapsePolicy/
│   │   │   └── UpdatePolicyEndorsement/
│   │   ├── Queries/
│   │   │   ├── GetPolicy/
│   │   │   ├── ListUserPolicies/
│   │   │   └── GetPolicyHistory/
│   │   ├── Domain/
│   │   │   ├── Policy.cs                ← Aggregate Root
│   │   │   ├── PolicyStatus.cs          ← Enum with state machine
│   │   │   ├── PolicyNumber.cs          ← Value Object
│   │   │   ├── Nominee.cs               ← Entity
│   │   │   ├── InsuredAsset.cs          ← Polymorphic bridge entity
│   │   │   ├── Rider.cs
│   │   │   └── PolicyDomainRules.cs     ← Business invariants
│   │   ├── StateMachine/
│   │   │   └── PolicyStateMachine.cs    ← Stateless/custom FSM
│   │   └── Infrastructure/
│   │       ├── PolicyRepository.cs
│   │       └── PolicyEventPublisher.cs
│   │
│   ├── Claims/
│   │   ├── Commands/
│   │   │   ├── SubmitClaim/
│   │   │   ├── ApproveClaim/
│   │   │   ├── RejectClaim/
│   │   │   ├── RequestDocuments/
│   │   │   └── SettleClaim/
│   │   ├── Queries/
│   │   │   ├── GetClaim/
│   │   │   ├── ListClaims/
│   │   │   └── GetClaimApprovalMatrix/
│   │   ├── Domain/
│   │   │   ├── Claim.cs                 ← Aggregate Root
│   │   │   ├── ClaimStatus.cs
│   │   │   ├── ClaimNumber.cs           ← Value Object
│   │   │   ├── ClaimDocument.cs
│   │   │   ├── ClaimApproval.cs
│   │   │   ├── ApprovalLevel.cs         ← Enum (L1/L2/L3/Board)
│   │   │   └── ClaimDomainRules.cs
│   │   ├── StateMachine/
│   │   │   └── ClaimStateMachine.cs
│   │   └── Infrastructure/
│   │       ├── ClaimRepository.cs
│   │       └── ClaimEventPublisher.cs
│   │
│   ├── Endorsements/
│   │   ├── Commands/
│   │   │   ├── CreateEndorsement/
│   │   │   └── ApproveEndorsement/
│   │   └── Domain/
│   │       ├── Endorsement.cs
│   │       └── EndorsementType.cs
│   │
│   ├── Renewals/
│   │   ├── Commands/
│   │   │   └── RenewPolicy/
│   │   ├── BackgroundJobs/
│   │   │   └── RenewalReminderJob.cs    ← Hangfire/Quartz daily job
│   │   └── Domain/
│   │       └── RenewalRules.cs
│   │
│   └── Fraud/
│       ├── Commands/
│       │   └── EvaluateFraudRisk/
│       ├── Rules/
│       │   ├── RapidClaimRule.cs        ← FR-182: <48hr
│       │   ├── FrequentClaimRule.cs     ← FR-183: >2/12mo
│       │   ├── AmountMatchingRule.cs    ← FR-184: 100% coverage
│       │   └── NetworkViolationRule.cs  ← FR-185
│       └── Domain/
│           └── FraudEvaluation.cs
│
├── InsuranceEngine.Domain/              ← Shared domain primitives
│   ├── Common/
│   │   ├── AggregateRoot.cs
│   │   ├── Entity.cs
│   │   ├── ValueObject.cs
│   │   └── IDomainEvent.cs
│   ├── Exceptions/
│   │   ├── DomainException.cs
│   │   ├── PolicyNotFoundException.cs
│   │   └── DuplicatePolicyException.cs
│   └── Events/
│       ├── PolicyIssuedEvent.cs
│       ├── ClaimSubmittedEvent.cs
│       └── PolicyCancelledEvent.cs
│
├── InsuranceEngine.Infrastructure/      ← Cross-cutting infrastructure
│   ├── Persistence/
│   │   ├── InsuranceDbContext.cs
│   │   ├── Configurations/              ← EF Core IEntityTypeConfiguration
│   │   │   ├── PolicyConfiguration.cs
│   │   │   ├── ClaimConfiguration.cs
│   │   │   └── ProductConfiguration.cs
│   │   ├── Migrations/
│   │   ├── UnitOfWork.cs
│   │   └── TransactionHelper.cs
│   ├── Messaging/
│   │   ├── KafkaProducer.cs
│   │   ├── KafkaConsumer.cs
│   │   └── KafkaTopics.cs              ← Topic constants
│   ├── Caching/
│   │   └── DistributedCacheService.cs  ← Redis wrapper (FR-028)
│   ├── ExternalServices/
│   │   ├── AIEngineClient.cs           ← gRPC client
│   │   ├── StorageServiceClient.cs     ← gRPC client
│   │   └── PartnerServiceClient.cs     ← gRPC client
│   └── BackgroundServices/
│       └── PolicyExpiryMonitorService.cs
│
├── InsuranceEngine.Contracts/           ← Proto-generated code (readonly)
│   └── Generated/                       ← Output of protoc
│       ├── insurance_engine_pb.cs
│       └── insurance_engine_grpc.cs
│
└── InsuranceEngine.Tests/
    ├── Unit/
    ├── Integration/
    └── Contract/
```

## 2.2 Layer Responsibilities

| Layer | দায়িত্ব | Pattern |
|---|---|---|
| **API** | gRPC endpoint mapping, interceptors | Decorator, Interceptor |
| **Features (Application)** | Commands, Queries, Handlers, Validators | CQRS, MediatR, FluentValidation |
| **Domain** | Entities, Value Objects, Aggregates, Domain Rules | DDD, Rich Domain Model |
| **Infrastructure** | DB, Cache, Messaging, External clients | Repository, UoW, Adapter |
| **Contracts** | Proto-generated C# classes | Generated (do not edit) |

---

# 3. Source-Based Technical Mapping

## 3.1 SRS v3.11 — Insurance Engine Mapping

| SRS Section | Insurance Engine-এর Relation | Usage Type |
|---|---|---|
| FG-003 (FR-021 to FR-029) | Products module — পুরো slice | Direct implementation |
| FG-004 (FR-030 to FR-041) | Policies module — core lifecycle | Direct implementation |
| FG-005 (FR-042 to FR-060) | Renewals + Cancellation + Endorsement | Direct implementation |
| FG-006 (FR-061 to FR-069) | Business Rules module | Implementation constraint |
| FG-008 (FR-081 to FR-101) | Claims module + Approval Matrix | Direct implementation |
| FG-016 (FR-175 to FR-188) | Fraud detection rules | Direct implementation |
| Section 6.4 (Proto-First) | Contracts layer — code generation | Architecture constraint |
| Section 6.6 (CQRS) | Application layer pattern | Implementation pattern |
| Appendix A (Proto defs) | Contracts/Generated/ — source | Direct usable |
| Claims Approval Matrix | ClaimStateMachine + ApprovalLevel enum | Direct implementation |
| Phase Scripts (1/2/3) | Implementation gap analysis | Reference + roadmap |

## 3.2 Phase Scripts Mapping

| Script | Purpose | Insurance Engine-এর জন্য |
|---|---|---|
| `update_srs_phase1.py` | M1 items status — কোনটা implemented, কোনটা gap | Critical gaps: Kafka (mock), Payment, PDF, Notifications |
| `update_srs_phase2.py` | M2 items — partial/not implemented list | Priority: Pro-rata refund, Grace period, Fraud patterns |
| `update_srs_phase3.py` | M3/D/F roadmap | Redis caching (FR-028) implemented; IoT/AI — future |

**Critical Gaps (Phase 1 script থেকে):**
1. Kafka integration — currently mock, real broker needed
2. PDF generation — mock only, needs OCR Service gRPC call
3. Refund processing — pro-rata calculation needs C# implementation
4. Notification trigger — SMS/Email not integrated (only Kafka event published)

---

# 4. Development Sequence

## Phase 1 — Foundation (Day 1-3)

### 1.1 Project Bootstrap

```bash
# Solution তৈরি
dotnet new sln -n InsuranceEngine
dotnet new webapi -n InsuranceEngine.API --no-https false
dotnet new classlib -n InsuranceEngine.Features
dotnet new classlib -n InsuranceEngine.Domain
dotnet new classlib -n InsuranceEngine.Infrastructure
dotnet new classlib -n InsuranceEngine.Contracts
dotnet new xunit -n InsuranceEngine.Tests
```

### 1.2 Core NuGet Packages

```xml
<!-- InsuranceEngine.API -->
<PackageReference Include="Grpc.AspNetCore" Version="2.62.0" />
<PackageReference Include="Microsoft.AspNetCore.Authentication.JwtBearer" Version="8.0.0" />

<!-- InsuranceEngine.Features -->
<PackageReference Include="MediatR" Version="12.3.0" />
<PackageReference Include="FluentValidation.DependencyInjectionExtensions" Version="11.9.0" />

<!-- InsuranceEngine.Infrastructure -->
<PackageReference Include="Microsoft.EntityFrameworkCore.Design" Version="8.0.0" />
<PackageReference Include="Npgsql.EntityFrameworkCore.PostgreSQL" Version="8.0.0" />
<PackageReference Include="Confluent.Kafka" Version="2.4.0" />
<PackageReference Include="StackExchange.Redis" Version="2.7.0" />
<PackageReference Include="Polly" Version="8.3.0" />
<PackageReference Include="Serilog.AspNetCore" Version="8.0.0" />
<PackageReference Include="OpenTelemetry.Exporter.Jaeger" Version="1.6.0" />
```

### 1.3 Base Abstractions (Domain Layer)

```csharp
// InsuranceEngine.Domain/Common/AggregateRoot.cs
public abstract class AggregateRoot<TId> : Entity<TId>
{
    private readonly List<IDomainEvent> _domainEvents = [];
    public IReadOnlyList<IDomainEvent> DomainEvents => _domainEvents.AsReadOnly();

    protected void AddDomainEvent(IDomainEvent domainEvent)
        => _domainEvents.Add(domainEvent);

    public void ClearDomainEvents() => _domainEvents.Clear();
}

// InsuranceEngine.Domain/Common/ValueObject.cs
public abstract class ValueObject
{
    protected abstract IEnumerable<object> GetEqualityComponents();
    public override bool Equals(object? obj) { /* structural equality */ }
    public override int GetHashCode() { /* hash from components */ }
}
```

### 1.4 Environment Config (appsettings.json)

```json
{
  "ConnectionStrings": {
    "DefaultConnection": "Host=localhost;Database=insurance_engine;Username=postgres;Password=..."
  },
  "Kafka": {
    "BootstrapServers": "localhost:9092",
    "ConsumerGroupId": "insurance-engine"
  },
  "Redis": {
    "ConnectionString": "localhost:6379"
  },
  "Grpc": {
    "Port": 5001,
    "AIEngineUrl": "https://localhost:4001",
    "StorageServiceUrl": "https://localhost:8084",
    "PartnerServiceUrl": "https://localhost:5002"
  },
  "ProductCacheTtl": 300
}
```

### 1.5 DI Registration (Program.cs Pattern)

```csharp
// InsuranceEngine.API/Program.cs
var builder = WebApplication.CreateBuilder(args);

builder.Services.AddGrpc(o => o.Interceptors.Add<AuthInterceptor>());
builder.Services.AddMediatR(cfg => cfg.RegisterServicesFromAssembly(
    typeof(InsuranceEngine.Features.AssemblyReference).Assembly));
builder.Services.AddValidatorsFromAssembly(FeaturesAssembly);
builder.Services.AddDbContext<InsuranceDbContext>(o =>
    o.UseNpgsql(builder.Configuration.GetConnectionString("DefaultConnection")));
builder.Services.AddStackExchangeRedisCache(o =>
    o.Configuration = builder.Configuration["Redis:ConnectionString"]);
builder.Services.AddScoped<IUnitOfWork, UnitOfWork>();

// Feature-specific registrations
builder.Services.AddProductFeature();
builder.Services.AddPolicyFeature();
builder.Services.AddClaimsFeature();
```

---

## Phase 2 — Domain Layer (Day 4-8)

### 2.1 Policy Aggregate (Most Critical)

```csharp
// InsuranceEngine.Features/Policies/Domain/Policy.cs
public sealed class Policy : AggregateRoot<Guid>
{
    // Core Fields
    public PolicyNumber PolicyNumber { get; private set; }
    public Guid CustomerId { get; private set; }
    public Guid ProductId { get; private set; }
    public PolicyStatus Status { get; private set; }
    public decimal PremiumAmount { get; private set; }
    public DateOnly StartDate { get; private set; }
    public DateOnly EndDate { get; private set; }
    public DateOnly? GracePeriodEndDate { get; private set; }

    // Nominees (FR-032: Single nominee required)
    private readonly List<Nominee> _nominees = [];
    public IReadOnlyList<Nominee> Nominees => _nominees.AsReadOnly();

    // Riders
    private readonly List<Rider> _riders = [];

    // Factory method (DDD pattern)
    public static Policy Create(
        Guid customerId, Guid productId, decimal premiumAmount,
        DateOnly startDate, int tenureDays, Nominee nominee)
    {
        var policy = new Policy
        {
            Id = Guid.NewGuid(),
            CustomerId = customerId,
            ProductId = productId,
            Status = PolicyStatus.PendingPayment,
            PremiumAmount = premiumAmount,
            StartDate = startDate,
            EndDate = startDate.AddDays(tenureDays)
        };

        policy.AddNominee(nominee);
        policy.AddDomainEvent(new PolicyCreatedEvent(policy.Id, customerId, productId));
        return policy;
    }

    public void Activate()
    {
        PolicyDomainRules.EnsureCanActivate(this);
        Status = PolicyStatus.Active;
        AddDomainEvent(new PolicyActivatedEvent(Id, CustomerId));
    }

    public void EnterGracePeriod()
    {
        PolicyDomainRules.EnsureCanEnterGracePeriod(this);
        Status = PolicyStatus.GracePeriod;
        GracePeriodEndDate = EndDate.AddDays(30); // FR-047
    }

    public void Lapse()
    {
        Status = PolicyStatus.Lapsed;
        AddDomainEvent(new PolicyLapsedEvent(Id, CustomerId));
    }

    public void Cancel(string reason)
    {
        PolicyDomainRules.EnsureCanCancel(this);
        Status = PolicyStatus.Cancelled;
        AddDomainEvent(new PolicyCancelledEvent(Id, CustomerId, reason));
    }

    // Pro-rata refund calculation (FR-053)
    public decimal CalculateProRataRefund(decimal adminFee, decimal cancellationCharge)
    {
        var totalDays = (EndDate.ToDateTime(TimeOnly.MinValue)
                        - StartDate.ToDateTime(TimeOnly.MinValue)).Days;
        var daysCovered = (DateOnly.FromDateTime(DateTime.UtcNow.Date)
                          - StartDate).Days;
        var unusedDays = totalDays - daysCovered;
        var unusedRatio = (decimal)unusedDays / totalDays;
        return (PremiumAmount * unusedRatio) - adminFee - cancellationCharge;
    }
}
```

### 2.2 Policy Status State Machine

```csharp
// InsuranceEngine.Features/Policies/Domain/PolicyStatus.cs
public enum PolicyStatus
{
    PendingPayment,
    Active,
    GracePeriod,
    Suspended,
    Lapsed,
    Cancelled,
    Expired
}

// InsuranceEngine.Features/Policies/StateMachine/PolicyStateMachine.cs
public static class PolicyStateMachine
{
    private static readonly Dictionary<PolicyStatus, HashSet<PolicyStatus>> _validTransitions = new()
    {
        [PolicyStatus.PendingPayment] = [PolicyStatus.Active, PolicyStatus.Cancelled],
        [PolicyStatus.Active]         = [PolicyStatus.GracePeriod, PolicyStatus.Suspended, PolicyStatus.Cancelled, PolicyStatus.Expired],
        [PolicyStatus.GracePeriod]    = [PolicyStatus.Active, PolicyStatus.Lapsed],
        [PolicyStatus.Lapsed]         = [PolicyStatus.Active], // Reinstatement within 90 days
        [PolicyStatus.Suspended]      = [PolicyStatus.Active, PolicyStatus.Cancelled],
        [PolicyStatus.Cancelled]      = [], // Terminal
        [PolicyStatus.Expired]        = []  // Terminal
    };

    public static void EnsureTransitionAllowed(PolicyStatus from, PolicyStatus to)
    {
        if (!_validTransitions[from].Contains(to))
            throw new InvalidPolicyTransitionException(from, to);
    }
}
```

### 2.3 Claims Aggregate

```csharp
// InsuranceEngine.Features/Claims/Domain/Claim.cs
public sealed class Claim : AggregateRoot<Guid>
{
    public ClaimNumber ClaimNumber { get; private set; }
    public Guid PolicyId { get; private set; }
    public Guid CustomerId { get; private set; }
    public ClaimStatus Status { get; private set; }
    public decimal ClaimAmount { get; private set; }
    public string ClaimReason { get; private set; }
    public DateTime SubmittedAt { get; private set; }
    public ApprovalLevel RequiredApprovalLevel { get; private set; }

    private readonly List<ClaimDocument> _documents = [];
    private readonly List<ClaimApproval> _approvals = [];

    public static Claim Submit(Guid policyId, Guid customerId,
        decimal amount, string reason, List<string> documentUrls)
    {
        // FR-082: Eligibility validation done in handler before calling this
        var claim = new Claim
        {
            Id = Guid.NewGuid(),
            PolicyId = policyId,
            CustomerId = customerId,
            ClaimAmount = amount,
            ClaimReason = reason,
            Status = ClaimStatus.Submitted,
            SubmittedAt = DateTime.UtcNow,
            RequiredApprovalLevel = DetermineApprovalLevel(amount)
        };

        foreach (var url in documentUrls)
            claim._documents.Add(new ClaimDocument(url));

        claim.AddDomainEvent(new ClaimSubmittedEvent(claim.Id, policyId, customerId, amount));
        return claim;
    }

    // Claims Approval Matrix (FR-086, SRS Table)
    private static ApprovalLevel DetermineApprovalLevel(decimal amount) => amount switch
    {
        <= 10_000m  => ApprovalLevel.L1,          // BDT 0-10K: Auto/Officer
        <= 50_000m  => ApprovalLevel.L2,          // BDT 10K-50K: Manager
        <= 200_000m => ApprovalLevel.L3,          // BDT 50K-2L: Business Admin + Focal Person
        _           => ApprovalLevel.Board        // BDT 2L+: Board + Insurer
    };
}
```

### 2.4 Value Objects

```csharp
// PolicyNumber Value Object (FR-034 format: LBT-YYYY-XXXX-NNNNNN)
public sealed class PolicyNumber : ValueObject
{
    public string Value { get; }

    private PolicyNumber(string value) => Value = value;

    public static PolicyNumber Generate(string insuranceTypeCode,
        string insurerCode, int productId, int sequenceNumber)
    {
        var year = DateTime.UtcNow.Year;
        var value = $"{insuranceTypeCode}-{year}-{productId:D4}-{sequenceNumber:D6}";
        return new PolicyNumber(value);
    }

    protected override IEnumerable<object> GetEqualityComponents()
    {
        yield return Value;
    }
}

// ClaimNumber Value Object (FR-083 format: CLM-YYYY-XXXX-NNNNNN)
public sealed class ClaimNumber : ValueObject
{
    public string Value { get; }
    public string Hash { get; } // SHA-256 for document integrity

    public static ClaimNumber Generate(int year, int productId, int seq, string documentPayload)
    {
        var value = $"CLM-{year}-{productId:D4}-{seq:D6}";
        var hash = ComputeSha256(documentPayload);
        return new ClaimNumber { Value = value, Hash = hash };
    }
    // ...
}
```

### 2.5 Product Entity

```csharp
// InsuranceEngine.Features/Products/Domain/Product.cs
public sealed class Product : AggregateRoot<Guid>
{
    public string Name { get; private set; }
    public ProductCategory Category { get; private set; } // Health/Life/Motor/Travel/Micro
    public decimal MinPremium { get; private set; }
    public decimal MaxPremium { get; private set; }
    public decimal MaxCoverage { get; private set; }
    public bool IsActive { get; private set; }
    public List<RiskAssessmentQuestion> AssessmentQuestions { get; private set; } = []; // FR-023-B
    public List<ProductRider> AvailableRiders { get; private set; } = [];
    public PremiumCalculationConfig CalculationConfig { get; private set; }
}

public enum ProductCategory
{
    Health, Life, Motor, Travel, MicroInsurance,
    Fire, Device, Livestock, Fisheries, Crop, GoodsInTransit, Pet, HomeAppliance
}
```

### 2.6 Insured Asset (Polymorphic Bridge)

```csharp
// InsuranceEngine.Features/Policies/Domain/InsuredAsset.cs
// Polymorphic bridge — different child tables per product type
public abstract class InsuredAsset : Entity<Guid>
{
    public Guid PolicyId { get; protected set; }
    public string AssetType { get; protected set; }
}

public sealed class HealthInsuredAsset : InsuredAsset
{
    public string HealthCondition { get; set; }
    public bool HasPreExistingCondition { get; set; }
    public List<string> MedicalHistory { get; set; } = [];
}

public sealed class MotorInsuredAsset : InsuredAsset
{
    public string VehicleRegNumber { get; set; }
    public string Make { get; set; }
    public int Year { get; set; }
    public decimal MarketValue { get; set; }
}
```

---

## Phase 3 — Contracts (gRPC/Proto) (Day 9-12)

### 3.1 Proto File Structure

```
proto/insuretech/
├── insurance_engine/
│   ├── entity/v1/
│   │   ├── product.proto
│   │   ├── policy.proto
│   │   ├── claim.proto
│   │   └── endorsement.proto
│   ├── events/v1/
│   │   ├── policy_events.proto
│   │   └── claim_events.proto
│   └── services/v1/
│       ├── product_service.proto
│       ├── policy_service.proto
│       └── claims_service.proto
```

### 3.2 Core Service Definitions

```protobuf
// proto/insuretech/insurance_engine/services/v1/policy_service.proto
syntax = "proto3";
package insuretech.insurance_engine.services.v1;

import "google/protobuf/timestamp.proto";
import "insuretech/insurance_engine/entity/v1/policy.proto";

service PolicyService {
  // M1 RPCs
  rpc CreatePolicy (CreatePolicyRequest) returns (CreatePolicyResponse);
  rpc GetPolicy (GetPolicyRequest) returns (GetPolicyResponse);
  rpc ListUserPolicies (ListUserPoliciesRequest) returns (ListUserPoliciesResponse);
  rpc ActivatePolicy (ActivatePolicyRequest) returns (ActivatePolicyResponse);
  rpc CancelPolicy (CancelPolicyRequest) returns (CancelPolicyResponse);

  // M2 RPCs
  rpc RenewPolicy (RenewPolicyRequest) returns (RenewPolicyResponse);
  rpc UpdateEndorsement (UpdateEndorsementRequest) returns (UpdateEndorsementResponse);
  rpc GetPolicyHistory (GetPolicyHistoryRequest) returns (GetPolicyHistoryResponse);
}

message CreatePolicyRequest {
  string customer_id = 1;
  string product_id = 2;
  double premium_amount = 3;
  NomineeDto nominee = 4;
  repeated string rider_ids = 5;
  string idempotency_key = 6; // FR-230
}

message NomineeDto {
  string full_name = 1;
  string relationship = 2;
  string phone_number = 3;
  optional string nid_number = 4;
  optional string income_range = 5; // FR-032-A: optional
}
```

```protobuf
// proto/.../claims_service.proto
service ClaimsService {
  rpc SubmitClaim (SubmitClaimRequest) returns (SubmitClaimResponse);
  rpc GetClaim (GetClaimRequest) returns (GetClaimResponse);
  rpc ApproveClaim (ApproveClaimRequest) returns (ApproveClaimResponse);
  rpc RejectClaim (RejectClaimRequest) returns (RejectClaimResponse);
  rpc RequestAdditionalDocuments (RequestDocsRequest) returns (RequestDocsResponse);
  rpc ListClaims (ListClaimsRequest) returns (ListClaimsResponse);
}

message SubmitClaimRequest {
  string policy_id = 1;
  double claim_amount = 2;
  string claim_reason = 3;
  repeated string document_urls = 4;
  string incident_date = 5; // ISO 8601
  string idempotency_key = 6;
}
```

### 3.3 Proto Evolution Strategy

**Recommended Engineering Decision:**
- Field number 1-15: Core fields (1 byte wire format — performance-critical fields এখানে)
- Field number 16+: Extended fields
- `optional` keyword সব non-mandatory fields-এ
- Deprecated fields: `[deprecated = true]` tag + keep জন্য backward compatibility
- Never reuse field numbers
- Version via package: `services/v1/` → `services/v2/` (breaking changes only)

---

## Phase 4 — Application Layer (CQRS) (Day 13-20)

### 4.1 Command Handler Pattern

```csharp
// InsuranceEngine.Features/Policies/Commands/CreatePolicy/CreatePolicyHandler.cs
public sealed class CreatePolicyHandler : IRequestHandler<CreatePolicyCommand, CreatePolicyResponse>
{
    private readonly IUnitOfWork _uow;
    private readonly IPolicyRepository _policyRepo;
    private readonly IProductRepository _productRepo;
    private readonly IKafkaProducer _kafka;
    private readonly IDistributedCache _cache;

    public async Task<CreatePolicyResponse> Handle(
        CreatePolicyCommand cmd, CancellationToken ct)
    {
        // 1. Idempotency check (FR-230)
        var existingByKey = await _policyRepo.FindByIdempotencyKeyAsync(cmd.IdempotencyKey, ct);
        if (existingByKey is not null)
            return MapToResponse(existingByKey); // Return cached response

        // 2. Product validation
        var product = await _productRepo.GetByIdAsync(cmd.ProductId, ct)
            ?? throw new ProductNotFoundException(cmd.ProductId);

        if (!product.IsActive)
            throw new ProductNotAvailableException(cmd.ProductId);

        // 3. Duplicate policy check (FR-063)
        var duplicateExists = await _policyRepo.ExistsWithinDaysAsync(
            cmd.CustomerId, cmd.ProductId, days: 30, ct);
        if (duplicateExists)
            throw new DuplicatePolicyException(cmd.CustomerId, cmd.ProductId);

        // 4. NID uniqueness check (FR-033)
        if (cmd.NidNumber is not null)
        {
            var nidExists = await _policyRepo.ExistsWithNidAsync(cmd.NidNumber, ct);
            if (nidExists)
                throw new DuplicateNidException(cmd.NidNumber);
        }

        // 5. Policy number generation (FR-034)
        var sequenceNumber = await _policyRepo.GetNextSequenceAsync(ct);
        var policyNumber = PolicyNumber.Generate(
            product.InsuranceTypeCode, "LBT", product.NumericId, sequenceNumber);

        // 6. Create aggregate
        var nominee = new Nominee(cmd.NomineeName, cmd.NomineeRelationship,
            cmd.NomineePhone, cmd.NomineeNid);
        var policy = Policy.Create(cmd.CustomerId, cmd.ProductId,
            cmd.PremiumAmount, DateOnly.FromDateTime(DateTime.UtcNow),
            product.TenureDays, nominee);

        // 7. Persist
        await _policyRepo.AddAsync(policy, ct);
        await _uow.CommitAsync(ct);

        // 8. Publish domain events via Kafka
        foreach (var evt in policy.DomainEvents)
            await _kafka.PublishAsync(KafkaTopics.PolicyEvents, evt, ct);

        policy.ClearDomainEvents();

        return new CreatePolicyResponse { PolicyId = policy.Id.ToString(),
            PolicyNumber = policy.PolicyNumber.Value };
    }
}
```

### 4.2 Validator Pattern

```csharp
// InsuranceEngine.Features/Policies/Commands/CreatePolicy/CreatePolicyValidator.cs
public sealed class CreatePolicyValidator : AbstractValidator<CreatePolicyCommand>
{
    public CreatePolicyValidator()
    {
        RuleFor(x => x.CustomerId).NotEmpty();
        RuleFor(x => x.ProductId).NotEmpty();
        RuleFor(x => x.PremiumAmount).GreaterThan(0);

        // Nominee validation (FR-032)
        RuleFor(x => x.NomineeName).NotEmpty().MaximumLength(200);
        RuleFor(x => x.NomineeRelationship).NotEmpty();
        RuleFor(x => x.NomineePhone)
            .Matches(@"^\+880\s?1[3-9]\d{8}$") // Bangladesh mobile format
            .WithMessage("Invalid Bangladesh mobile number");

        // Beneficiary income optional (FR-032-A)
        // income_range: no validation rule — it's optional

        // Idempotency key
        RuleFor(x => x.IdempotencyKey).NotEmpty().Must(k => Guid.TryParse(k, out _))
            .WithMessage("IdempotencyKey must be a valid UUID");
    }
}
```

### 4.3 Fraud Evaluation — Strategy Pattern

```csharp
// InsuranceEngine.Features/Fraud/Rules/ — Strategy Pattern
public interface IFraudRule
{
    string RuleId { get; }
    Task<FraudRuleResult> EvaluateAsync(FraudEvaluationContext ctx, CancellationToken ct);
}

// FR-182: Rapid Policy-Claim (<48hr)
public sealed class RapidClaimRule : IFraudRule
{
    public string RuleId => "FR-182";

    public async Task<FraudRuleResult> EvaluateAsync(FraudEvaluationContext ctx, CancellationToken ct)
    {
        var hoursSincePurchase = (ctx.ClaimSubmissionTime - ctx.PolicyStartDateTime).TotalHours;
        if (hoursSincePurchase < 48)
            return FraudRuleResult.Flagged(RuleId, $"Claim within {hoursSincePurchase:F1}hr of purchase");

        return FraudRuleResult.Clear(RuleId);
    }
}

// FR-183: Frequent Claims (>2 same type in 12 months)
public sealed class FrequentClaimRule : IFraudRule { /* ... */ }

// Orchestrator
public sealed class FraudEvaluationService
{
    private readonly IEnumerable<IFraudRule> _rules;

    public async Task<FraudEvaluation> EvaluateAsync(Claim claim, CancellationToken ct)
    {
        var ctx = await BuildContextAsync(claim, ct);
        var results = await Task.WhenAll(_rules.Select(r => r.EvaluateAsync(ctx, ct)));
        return FraudEvaluation.From(claim.Id, results);
    }
}
```

### 4.4 Premium Calculator — Strategy Pattern

```csharp
// InsuranceEngine.Features/Products/Domain/PremiumCalculation/
public interface IPremiumCalculationStrategy
{
    decimal Calculate(PremiumCalculationInput input);
}

public sealed class StandardPremiumStrategy : IPremiumCalculationStrategy
{
    public decimal Calculate(PremiumCalculationInput input)
    {
        var basePremium = input.SumAssured * input.BaseRate;

        // Age loading (FR-062)
        var ageFactor = input.Age switch
        {
            < 18 => 0.8m,
            < 30 => 1.0m,
            < 45 => 1.2m,
            < 60 => 1.5m,
            _    => 2.0m
        };

        // Occupation risk loading
        var occupationFactor = input.OccupationRiskLevel switch
        {
            OccupationRisk.Low    => 1.0m,
            OccupationRisk.Medium => 1.15m,
            OccupationRisk.High   => 1.35m,
            _                    => 1.0m
        };

        // Rider additions
        var riderPremium = input.Riders.Sum(r => r.AdditionalPremium);

        return (basePremium * ageFactor * occupationFactor) + riderPremium;
    }
}
```

### 4.5 Claims Approval Handler

```csharp
// InsuranceEngine.Features/Claims/Commands/ApproveClaim/ApproveClaimHandler.cs
public sealed class ApproveClaimHandler : IRequestHandler<ApproveClaimCommand, ApproveClaimResponse>
{
    public async Task<ApproveClaimResponse> Handle(ApproveClaimCommand cmd, CancellationToken ct)
    {
        var claim = await _claimRepo.GetByIdAsync(cmd.ClaimId, ct)
            ?? throw new ClaimNotFoundException(cmd.ClaimId);

        // Approval level authorization check
        var approverRole = _currentUser.GetRole();
        ClaimDomainRules.EnsureApproverAuthorized(claim.RequiredApprovalLevel, approverRole);

        // Joint approval check for L3 (FR-091: BDT 50K-2L needs BA + FP)
        if (claim.RequiredApprovalLevel == ApprovalLevel.L3)
        {
            ClaimDomainRules.EnsureJointApprovalComplete(claim, cmd.ApproverId);
        }

        claim.Approve(cmd.ApproverId, cmd.Notes);

        await _claimRepo.UpdateAsync(claim, ct);
        await _uow.CommitAsync(ct);

        // Trigger payment (FR-092: auto-payment upon approval)
        await _kafka.PublishAsync(KafkaTopics.ClaimApproved, new ClaimApprovedEvent
        {
            ClaimId = claim.Id.ToString(),
            PolicyId = claim.PolicyId.ToString(),
            Amount = (double)claim.ClaimAmount,
            CustomerId = claim.CustomerId.ToString()
        }, ct);

        return new ApproveClaimResponse { ClaimId = claim.Id.ToString(), Status = "APPROVED" };
    }
}
```

---

## Phase 5 — Infrastructure Layer (Day 21-26)

### 5.1 EF Core DbContext

```csharp
// InsuranceEngine.Infrastructure/Persistence/InsuranceDbContext.cs
public sealed class InsuranceDbContext : DbContext
{
    public DbSet<Policy> Policies => Set<Policy>();
    public DbSet<Claim> Claims => Set<Claim>();
    public DbSet<Product> Products => Set<Product>();
    public DbSet<Nominee> Nominees => Set<Nominee>();
    public DbSet<ClaimDocument> ClaimDocuments => Set<ClaimDocument>();
    public DbSet<ClaimApproval> ClaimApprovals => Set<ClaimApproval>();
    public DbSet<InsuredAsset> InsuredAssets => Set<InsuredAsset>();

    protected override void OnModelCreating(ModelBuilder mb)
    {
        mb.ApplyConfigurationsFromAssembly(typeof(InsuranceDbContext).Assembly);

        // Polymorphic insured assets — TPH (Table Per Hierarchy)
        mb.Entity<InsuredAsset>()
            .HasDiscriminator<string>("asset_type")
            .HasValue<HealthInsuredAsset>("HEALTH")
            .HasValue<MotorInsuredAsset>("MOTOR");
    }
}
```

### 5.2 Policy Configuration (EF Core)

```csharp
// InsuranceEngine.Infrastructure/Persistence/Configurations/PolicyConfiguration.cs
public sealed class PolicyConfiguration : IEntityTypeConfiguration<Policy>
{
    public void Configure(EntityTypeBuilder<Policy> b)
    {
        b.ToTable("policies");
        b.HasKey(p => p.Id);
        b.Property(p => p.Id).HasColumnName("policy_id");

        // Value Object mapping
        b.OwnsOne(p => p.PolicyNumber, pn =>
        {
            pn.Property(x => x.Value)
              .HasColumnName("policy_number")
              .HasMaxLength(50)
              .IsRequired();
            pn.HasIndex(x => x.Value).IsUnique();
        });

        b.Property(p => p.Status)
            .HasConversion<string>()
            .HasMaxLength(30);

        b.Property(p => p.PremiumAmount)
            .HasPrecision(12, 2)
            .IsRequired();

        // Indexes
        b.HasIndex(p => p.CustomerId);
        b.HasIndex(p => new { p.CustomerId, p.ProductId });
        b.HasIndex(p => p.Status);

        // Audit timestamps
        b.Property<DateTime>("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP");
        b.Property<DateTime>("updated_at");
    }
}
```

### 5.3 Unit of Work

```csharp
// InsuranceEngine.Infrastructure/Persistence/UnitOfWork.cs
public sealed class UnitOfWork : IUnitOfWork
{
    private readonly InsuranceDbContext _ctx;
    private IDbContextTransaction? _transaction;

    public async Task BeginTransactionAsync(CancellationToken ct = default)
        => _transaction = await _ctx.Database.BeginTransactionAsync(ct);

    public async Task CommitAsync(CancellationToken ct = default)
    {
        // Dispatch domain events before commit
        var aggregates = _ctx.ChangeTracker.Entries<AggregateRoot<Guid>>()
            .Where(e => e.Entity.DomainEvents.Any())
            .Select(e => e.Entity);

        // Save DB changes
        await _ctx.SaveChangesAsync(ct);

        if (_transaction is not null)
            await _transaction.CommitAsync(ct);
    }

    public async Task RollbackAsync(CancellationToken ct = default)
    {
        if (_transaction is not null)
            await _transaction.RollbackAsync(ct);
    }
}
```

### 5.4 Policy Repository

```csharp
// InsuranceEngine.Features/Policies/Infrastructure/PolicyRepository.cs
public sealed class PolicyRepository : IPolicyRepository
{
    private readonly InsuranceDbContext _ctx;

    public async Task<Policy?> GetByIdAsync(Guid id, CancellationToken ct)
        => await _ctx.Policies
            .Include(p => p.Nominees)
            .Include(p => p.Riders)
            .FirstOrDefaultAsync(p => p.Id == id, ct);

    public async Task<bool> ExistsWithinDaysAsync(
        Guid customerId, Guid productId, int days, CancellationToken ct)
    {
        var cutoff = DateTime.UtcNow.AddDays(-days);
        return await _ctx.Policies
            .AnyAsync(p => p.CustomerId == customerId
                        && p.ProductId == productId
                        && p.CreatedAt >= cutoff
                        && p.Status != PolicyStatus.Cancelled, ct);
    }

    // Specification Pattern for complex queries
    public async Task<PagedResult<Policy>> ListAsync(
        PolicySpecification spec, CancellationToken ct)
    {
        var query = _ctx.Policies.AsQueryable();
        query = spec.Apply(query);
        var total = await query.CountAsync(ct);
        var items = await query.Skip(spec.Skip).Take(spec.Take).ToListAsync(ct);
        return new PagedResult<Policy>(items, total);
    }
}
```

### 5.5 Kafka Integration

```csharp
// InsuranceEngine.Infrastructure/Messaging/KafkaProducer.cs
public sealed class KafkaProducer : IKafkaProducer
{
    private readonly IProducer<string, byte[]> _producer;

    public async Task PublishAsync<TEvent>(string topic, TEvent @event, CancellationToken ct)
        where TEvent : IMessage<TEvent>
    {
        var bytes = @event.ToByteArray(); // Protobuf serialization
        var message = new Message<string, byte[]>
        {
            Key = Guid.NewGuid().ToString(),
            Value = bytes,
            Headers = new Headers
            {
                { "event-type", Encoding.UTF8.GetBytes(typeof(TEvent).Name) },
                { "timestamp", BitConverter.GetBytes(DateTimeOffset.UtcNow.ToUnixTimeMilliseconds()) }
            }
        };

        await _producer.ProduceAsync(topic, message, ct);
    }
}

// KafkaTopics constants
public static class KafkaTopics
{
    public const string PolicyEvents   = "insurance.policy.events";
    public const string ClaimEvents    = "insurance.claim.events";
    public const string PaymentEvents  = "payment.events";      // consumed
    public const string ClaimApproved  = "insurance.claim.approved";
    public const string RenewalEvents  = "insurance.renewal.events";
}
```

### 5.6 Redis Caching (FR-028)

```csharp
// InsuranceEngine.Infrastructure/Caching/DistributedCacheService.cs
public sealed class DistributedCacheService : IDistributedCacheService
{
    private readonly IDistributedCache _cache;

    // Product cache: 5-minute TTL (FR-028)
    public async Task<Product?> GetProductAsync(Guid productId, CancellationToken ct)
    {
        var key = $"product:{productId}";
        var cached = await _cache.GetAsync(key, ct);
        if (cached is null) return null;
        return MessagePackSerializer.Deserialize<Product>(cached);
    }

    public async Task SetProductAsync(Product product, CancellationToken ct)
    {
        var key = $"product:{product.Id}";
        var bytes = MessagePackSerializer.Serialize(product);
        await _cache.SetAsync(key, bytes, new DistributedCacheEntryOptions
        {
            AbsoluteExpirationRelativeToNow = TimeSpan.FromMinutes(5) // FR-028
        }, ct);
    }

    // Invalidate on product update
    public async Task InvalidateProductAsync(Guid productId, CancellationToken ct)
        => await _cache.RemoveAsync($"product:{productId}", ct);
}
```

---

## Phase 6 — Rule Engine (Day 27-31)

### 6.1 Business Rules Registry

```csharp
// InsuranceEngine.Features/Policies/Domain/PolicyDomainRules.cs
public static class PolicyDomainRules
{
    public static void EnsureCanActivate(Policy policy)
    {
        PolicyStateMachine.EnsureTransitionAllowed(policy.Status, PolicyStatus.Active);
    }

    public static void EnsureCanCancel(Policy policy)
    {
        if (policy.Status is PolicyStatus.Cancelled or PolicyStatus.Expired)
            throw new DomainException("Policy already in terminal state");

        PolicyStateMachine.EnsureTransitionAllowed(policy.Status, PolicyStatus.Cancelled);
    }

    // FR-060: Endorsement approval if sum change >10%
    public static void EnsureSumChangeApprovalRequired(
        decimal currentSum, decimal newSum, out bool approvalRequired)
    {
        var changePercent = Math.Abs((newSum - currentSum) / currentSum * 100);
        approvalRequired = changePercent > 10;
    }

    // FR-063: Duplicate detection
    public static void EnsureNoDuplicateWithinDays(bool duplicateExists, Guid customerId)
    {
        if (duplicateExists)
            throw new DuplicatePolicyException(customerId);
    }
}
```

### 6.2 Claims Domain Rules

```csharp
// InsuranceEngine.Features/Claims/Domain/ClaimDomainRules.cs
public static class ClaimDomainRules
{
    // FR-082: Eligibility
    public static void EnsureClaimEligible(Policy policy, ClaimType claimType)
    {
        if (policy.Status is not PolicyStatus.Active and not PolicyStatus.GracePeriod)
            throw new ClaimIneligibleException("Policy is not active");

        if (DateTime.UtcNow.Date > policy.EndDate.ToDateTime(TimeOnly.MinValue).AddDays(30))
            throw new ClaimIneligibleException("Claim period expired");
    }

    // FR-091: Joint approval for L3
    public static void EnsureJointApprovalComplete(Claim claim, Guid currentApproverId)
    {
        if (claim.RequiredApprovalLevel != ApprovalLevel.L3) return;

        var hasBA = claim.Approvals.Any(a => a.ApproverRole == "BusinessAdmin");
        var hasFP = claim.Approvals.Any(a => a.ApproverRole == "FocalPerson");

        // Current approver adds one; check if both roles covered
        // (validation happens after this approval is added)
    }

    // Approval Matrix (SRS Claims Approval Matrix table)
    public static void EnsureApproverAuthorized(ApprovalLevel required, string approverRole)
    {
        var authorized = required switch
        {
            ApprovalLevel.L1   => approverRole is "System" or "ClaimsOfficer",
            ApprovalLevel.L2   => approverRole is "ClaimsManager",
            ApprovalLevel.L3   => approverRole is "BusinessAdmin" or "FocalPerson",
            ApprovalLevel.Board => approverRole is "BoardMember" or "InsurerApprover",
            _ => false
        };

        if (!authorized)
            throw new UnauthorizedApprovalException(required, approverRole);
    }

    // FR-100: Co-payment and deductible calculation
    public static decimal CalculatePayableAmount(
        decimal claimAmount, decimal deductible, decimal copaymentPercent)
    {
        var afterDeductible = Math.Max(0, claimAmount - deductible);
        return afterDeductible * (1 - copaymentPercent / 100);
    }
}
```

---

## Phase 7 — Integration (Day 32-38)

### 7.1 Payment Event Consumer (Kafka)

```csharp
// InsuranceEngine.Infrastructure/Messaging/PaymentEventConsumer.cs
// Listens to: payment.events topic
// On payment.completed → activate policy
public sealed class PaymentEventConsumer : BackgroundService
{
    protected override async Task ExecuteAsync(CancellationToken ct)
    {
        using var consumer = new ConsumerBuilder<string, byte[]>(_config).Build();
        consumer.Subscribe(KafkaTopics.PaymentEvents);

        while (!ct.IsCancellationRequested)
        {
            var result = consumer.Consume(ct);
            var evt = PaymentCompletedEvent.Parser.ParseFrom(result.Message.Value);

            // Dispatch to MediatR
            await _mediator.Send(new ActivatePolicyCommand { PolicyId = evt.PolicyId }, ct);
            consumer.Commit(result);
        }
    }
}
```

### 7.2 AI Engine Integration (Fraud Score)

```csharp
// InsuranceEngine.Infrastructure/ExternalServices/AIEngineClient.cs
public sealed class AIEngineClient : IAIEngineClient
{
    private readonly AIEngine.AIEngineClient _grpcClient;

    public async Task<FraudRiskScore> GetFraudRiskScoreAsync(
        string claimId, string customerId, decimal amount, CancellationToken ct)
    {
        var request = new FraudRiskRequest
        {
            ClaimId = claimId,
            CustomerId = customerId,
            Amount = (double)amount
        };

        // Retry with Polly (FR-229: exponential backoff)
        var response = await _retryPolicy.ExecuteAsync(
            () => _grpcClient.GetFraudRiskScoreAsync(request, cancellationToken: ct));

        return new FraudRiskScore(response.Score, response.RiskLevel);
    }
}
```

### 7.3 Renewal Background Job

```csharp
// InsuranceEngine.Features/Renewals/BackgroundJobs/RenewalReminderJob.cs
// FR-043: 30 days, 7 days, 1 day reminder
public sealed class RenewalReminderJob : IHostedService
{
    private Timer? _timer;

    public Task StartAsync(CancellationToken ct)
    {
        _timer = new Timer(Execute, null, TimeSpan.Zero, TimeSpan.FromHours(24));
        return Task.CompletedTask;
    }

    private async void Execute(object? state)
    {
        var reminderDays = new[] { 30, 7, 1 };
        foreach (var days in reminderDays)
        {
            var targetDate = DateOnly.FromDateTime(DateTime.UtcNow.AddDays(days));
            var policies = await _policyRepo.GetExpiringOnDateAsync(targetDate);

            foreach (var policy in policies)
            {
                await _kafka.PublishAsync(KafkaTopics.RenewalEvents,
                    new RenewalReminderEvent
                    {
                        PolicyId = policy.Id.ToString(),
                        CustomerId = policy.CustomerId.ToString(),
                        ExpiryDate = policy.EndDate.ToString("yyyy-MM-dd"),
                        DaysRemaining = days
                    });
            }
        }
    }
}
```

---

## Phase 8 — Security + Reliability (Day 39-43)

### 8.1 gRPC Auth Interceptor

```csharp
// InsuranceEngine.API/Interceptors/AuthInterceptor.cs
public sealed class AuthInterceptor : Interceptor
{
    public override async Task<TResponse> UnaryServerHandler<TRequest, TResponse>(
        TRequest request, ServerCallContext context,
        UnaryServerMethod<TRequest, TResponse> continuation)
    {
        var token = context.RequestHeaders
            .FirstOrDefault(h => h.Key == "authorization")?.Value;

        if (string.IsNullOrEmpty(token))
            throw new RpcException(new Status(StatusCode.Unauthenticated, "Token required"));

        // Validate with Auth Service
        var claims = await _authClient.ValidateTokenAsync(token);
        if (claims is null)
            throw new RpcException(new Status(StatusCode.Unauthenticated, "Invalid token"));

        // Inject user context
        context.UserState["user_id"] = claims.UserId;
        context.UserState["roles"] = claims.Roles;

        return await continuation(request, context);
    }
}
```

### 8.2 Exception Interceptor

```csharp
// InsuranceEngine.API/Interceptors/ExceptionInterceptor.cs
public sealed class ExceptionInterceptor : Interceptor
{
    public override async Task<TResponse> UnaryServerHandler<TRequest, TResponse>(
        TRequest request, ServerCallContext context,
        UnaryServerMethod<TRequest, TResponse> continuation)
    {
        try
        {
            return await continuation(request, context);
        }
        catch (DomainException ex)
        {
            _logger.LogWarning(ex, "Domain rule violation");
            throw new RpcException(new Status(StatusCode.InvalidArgument, ex.Message));
        }
        catch (NotFoundException ex)
        {
            throw new RpcException(new Status(StatusCode.NotFound, ex.Message));
        }
        catch (UnauthorizedApprovalException ex)
        {
            throw new RpcException(new Status(StatusCode.PermissionDenied, ex.Message));
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Unhandled exception in {Method}", context.Method);
            throw new RpcException(new Status(StatusCode.Internal, "Internal error"));
        }
    }
}
```

### 8.3 Resilience (Polly)

```csharp
// Program.cs — Polly policies
builder.Services.AddHttpClient<IAIEngineClient>()
    .AddPolicyHandler(Policy.WrapAsync(
        // Retry: 1s, 2s, 4s, 8s, 16s — FR-229
        HttpPolicyExtensions.HandleTransientHttpError()
            .WaitAndRetryAsync(5, attempt => TimeSpan.FromSeconds(Math.Pow(2, attempt - 1))),
        // Circuit Breaker
        HttpPolicyExtensions.HandleTransientHttpError()
            .CircuitBreakerAsync(3, TimeSpan.FromSeconds(30))
    ));
```

### 8.4 Audit Logging

```csharp
// InsuranceEngine.Infrastructure/Audit/AuditLogger.cs
// FR-206: Immutable audit logs — append-only PostgreSQL table
public sealed class AuditLogger : IAuditLogger
{
    public async Task LogAsync(AuditEntry entry, CancellationToken ct)
    {
        // PostgreSQL append-only (no UPDATE/DELETE on audit table)
        await _ctx.Database.ExecuteSqlRawAsync(
            "INSERT INTO audit_logs (id, action, entity_type, entity_id, user_id, timestamp, payload) " +
            "VALUES (@id, @action, @entityType, @entityId, @userId, NOW(), @payload::jsonb)",
            new NpgsqlParameter("@id", Guid.NewGuid()),
            new NpgsqlParameter("@action", entry.Action),
            // ...
        );
    }
}
```

---

## Phase 9 — Testing (Day 44-50)

### 9.1 Unit Test — Domain Rules

```csharp
// InsuranceEngine.Tests/Unit/Policies/PolicyStateMachineTests.cs
public sealed class PolicyStateMachineTests
{
    [Fact]
    public void ActivePolicy_CanTransitionTo_GracePeriod()
    {
        // Arrange
        var policy = Policy.Create(/* ... */);
        policy.Activate();

        // Act & Assert
        var act = () => PolicyStateMachine.EnsureTransitionAllowed(
            PolicyStatus.Active, PolicyStatus.GracePeriod);
        act.Should().NotThrow();
    }

    [Fact]
    public void CancelledPolicy_CannotTransition_ToAnyState()
    {
        var act = () => PolicyStateMachine.EnsureTransitionAllowed(
            PolicyStatus.Cancelled, PolicyStatus.Active);
        act.Should().Throw<InvalidPolicyTransitionException>();
    }

    [Theory]
    [InlineData(5_000, ApprovalLevel.L1)]
    [InlineData(30_000, ApprovalLevel.L2)]
    [InlineData(100_000, ApprovalLevel.L3)]
    [InlineData(500_000, ApprovalLevel.Board)]
    public void ClaimAmount_MapsTo_CorrectApprovalLevel(
        decimal amount, ApprovalLevel expected)
    {
        var claim = Claim.Submit(/* ... amount ... */);
        claim.RequiredApprovalLevel.Should().Be(expected);
    }
}
```

### 9.2 Integration Test

```csharp
// InsuranceEngine.Tests/Integration/CreatePolicyTests.cs
public sealed class CreatePolicyTests : IClassFixture<IntegrationTestFactory>
{
    [Fact]
    public async Task CreatePolicy_WithValidData_ReturnsSuccess()
    {
        // Uses Testcontainers for PostgreSQL + Redis
        var cmd = new CreatePolicyCommand
        {
            CustomerId = Guid.NewGuid().ToString(),
            ProductId = _testProductId.ToString(),
            PremiumAmount = 500,
            NomineeName = "Test Nominee",
            NomineeRelationship = "Spouse",
            NomineePhone = "+8801711000000",
            IdempotencyKey = Guid.NewGuid().ToString()
        };

        var result = await _mediator.Send(cmd);

        result.PolicyNumber.Should().StartWith("LBT-");
        result.PolicyNumber.Should().MatchRegex(@"^[A-Z]+-\d{4}-\d{4}-\d{6}$");
    }

    [Fact]
    public async Task CreatePolicy_WithDuplicateIdempotencyKey_ReturnsSameResult()
    {
        var key = Guid.NewGuid().ToString();
        var result1 = await _mediator.Send(new CreatePolicyCommand { IdempotencyKey = key, /* ... */ });
        var result2 = await _mediator.Send(new CreatePolicyCommand { IdempotencyKey = key, /* ... */ });

        result1.PolicyNumber.Should().Be(result2.PolicyNumber); // Idempotent
    }
}
```

### 9.3 Fraud Rule Tests

```csharp
[Theory]
[InlineData(24, true)]   // 24hr — flagged
[InlineData(48, false)]  // exactly 48hr — boundary
[InlineData(72, false)]  // 72hr — clear
public async Task RapidClaimRule_FlagsCorrectly(double hoursSincePurchase, bool shouldFlag)
{
    var ctx = new FraudEvaluationContext
    {
        PolicyStartDateTime = DateTime.UtcNow.AddHours(-hoursSincePurchase),
        ClaimSubmissionTime = DateTime.UtcNow
    };

    var result = await new RapidClaimRule().EvaluateAsync(ctx, CancellationToken.None);

    result.IsFlagged.Should().Be(shouldFlag);
}
```

---

## Phase 10 — Production Readiness (Day 51-55)

### 10.1 Health Checks

```csharp
// Program.cs
builder.Services.AddHealthChecks()
    .AddNpgsql(connectionString, name: "postgresql")
    .AddRedis(redisConnection, name: "redis")
    .AddKafka(kafkaConfig, name: "kafka")
    .AddGrpcCheck<AIEngineClient>(name: "ai-engine");
```

### 10.2 OpenTelemetry

```csharp
builder.Services.AddOpenTelemetry()
    .WithTracing(tracing => tracing
        .AddAspNetCoreInstrumentation()
        .AddGrpcClientInstrumentation()
        .AddEntityFrameworkCoreInstrumentation()
        .AddJaegerExporter(o => o.AgentHost = "jaeger"))
    .WithMetrics(metrics => metrics
        .AddAspNetCoreInstrumentation()
        .AddRuntimeInstrumentation()
        .AddPrometheusExporter());
```

### 10.3 Database Migrations Ordering

```
migrations/
├── 001_create_products_table.sql
├── 002_create_policies_table.sql
├── 003_create_nominees_table.sql
├── 004_create_riders_table.sql
├── 005_create_insured_assets_table.sql
├── 006_create_claims_table.sql
├── 007_create_claim_documents_table.sql
├── 008_create_claim_approvals_table.sql
├── 009_create_endorsements_table.sql
├── 010_create_fraud_evaluations_table.sql
├── 011_create_audit_logs_table.sql    ← append-only constraint
├── 012_add_indexes.sql
├── 013_add_foreign_keys.sql
└── 014_add_partitioning.sql           ← policies/claims by month
```

---

# 5. Technical Standards Summary

| Concern | Pattern | কোথায় |
|---|---|---|
| Feature organization | Vertical Slice | Features/[Domain]/Commands,Queries |
| Write operations | CQRS Command + MediatR Handler | Features/*/Commands/ |
| Read operations | CQRS Query + MediatR Handler | Features/*/Queries/ |
| Input validation | FluentValidation (per command) | Features/*/Commands/*/Validator |
| DB access | Repository Pattern per slice | Features/*/Infrastructure/ |
| Complex queries | Specification Pattern | PolicySpecification, ClaimSpecification |
| Business object creation | Factory Method (static .Create()) | Domain entities |
| Premium calculation | Strategy Pattern | Products/Domain/PremiumCalculation/ |
| Fraud detection | Strategy Pattern + Chain | Fraud/Rules/*.cs |
| Policy state | Finite State Machine | Policies/StateMachine/ |
| Cross-cutting concern | MediatR Pipeline Behavior | Logging, Validation, Transaction |
| Resilience | Polly (retry + circuit breaker) | Infrastructure/ExternalServices/ |
| Caching | Decorator via IDistributedCache | Infrastructure/Caching/ |
| Event publishing | Adapter (Kafka) | Infrastructure/Messaging/ |

---

# 6. Exact File/Folder Paths

```
InsuranceEngine.Features/
├── Products/
│   ├── Commands/CreateProduct/CreateProductCommand.cs
│   ├── Commands/CreateProduct/CreateProductHandler.cs
│   ├── Commands/CreateProduct/CreateProductValidator.cs
│   ├── Commands/UpdateProduct/...
│   ├── Commands/DeactivateProduct/...
│   ├── Queries/GetProduct/GetProductQuery.cs
│   ├── Queries/GetProduct/GetProductHandler.cs
│   ├── Queries/SearchProducts/SearchProductsQuery.cs
│   ├── Queries/SearchProducts/SearchProductsHandler.cs
│   ├── Domain/Product.cs
│   ├── Domain/ProductCategory.cs
│   ├── Domain/ProductRules.cs
│   ├── Domain/PremiumCalculation/IPremiumCalculationStrategy.cs
│   ├── Domain/PremiumCalculation/StandardPremiumStrategy.cs
│   └── Infrastructure/ProductRepository.cs
│
├── Policies/
│   ├── Commands/CreatePolicy/CreatePolicyCommand.cs
│   ├── Commands/CreatePolicy/CreatePolicyHandler.cs
│   ├── Commands/CreatePolicy/CreatePolicyValidator.cs
│   ├── Commands/ActivatePolicy/ActivatePolicyCommand.cs
│   ├── Commands/ActivatePolicy/ActivatePolicyHandler.cs
│   ├── Commands/CancelPolicy/CancelPolicyCommand.cs
│   ├── Commands/CancelPolicy/CancelPolicyHandler.cs
│   ├── Commands/RenewPolicy/...
│   ├── Commands/UpdatePolicyEndorsement/...
│   ├── Queries/GetPolicy/...
│   ├── Queries/ListUserPolicies/...
│   ├── Queries/GetPolicyHistory/...
│   ├── Domain/Policy.cs
│   ├── Domain/PolicyStatus.cs
│   ├── Domain/PolicyNumber.cs
│   ├── Domain/Nominee.cs
│   ├── Domain/Rider.cs
│   ├── Domain/InsuredAsset.cs
│   ├── Domain/HealthInsuredAsset.cs
│   ├── Domain/MotorInsuredAsset.cs
│   ├── Domain/PolicyDomainRules.cs
│   ├── StateMachine/PolicyStateMachine.cs
│   └── Infrastructure/PolicyRepository.cs
│
├── Claims/
│   ├── Commands/SubmitClaim/...
│   ├── Commands/ApproveClaim/...
│   ├── Commands/RejectClaim/...
│   ├── Commands/RequestDocuments/...
│   ├── Commands/SettleClaim/...
│   ├── Queries/GetClaim/...
│   ├── Queries/ListClaims/...
│   ├── Domain/Claim.cs
│   ├── Domain/ClaimStatus.cs
│   ├── Domain/ClaimNumber.cs
│   ├── Domain/ClaimDocument.cs
│   ├── Domain/ClaimApproval.cs
│   ├── Domain/ApprovalLevel.cs
│   ├── Domain/ClaimDomainRules.cs
│   ├── StateMachine/ClaimStateMachine.cs
│   └── Infrastructure/ClaimRepository.cs
│
├── Endorsements/
│   ├── Commands/CreateEndorsement/...
│   ├── Commands/ApproveEndorsement/...
│   ├── Domain/Endorsement.cs
│   └── Domain/EndorsementType.cs
│
├── Renewals/
│   ├── Commands/RenewPolicy/...
│   ├── BackgroundJobs/RenewalReminderJob.cs
│   └── Domain/RenewalRules.cs
│
└── Fraud/
    ├── Commands/EvaluateFraudRisk/...
    ├── Rules/RapidClaimRule.cs
    ├── Rules/FrequentClaimRule.cs
    ├── Rules/AmountMatchingRule.cs
    ├── Rules/NetworkViolationRule.cs
    ├── Domain/FraudEvaluation.cs
    ├── Domain/FraudEvaluationContext.cs
    └── Services/FraudEvaluationService.cs
```

---

# 7. Hidden Critical Decisions

## 7.1 Transaction Boundary

**Recommended Engineering Decision:**

Policy creation-এ একটি database transaction-এ সব লিখতে হবে:
`Policy + Nominees + Riders` — atomic। Kafka publish হবে **after** commit, not inside transaction। Domain events `UnitOfWork.CommitAsync()` পরে publish হবে।

```
Begin TX
  → INSERT policies
  → INSERT nominees
  → INSERT riders
Commit TX
  → Publish Kafka events (async, outside TX)
```

## 7.2 Idempotency (FR-230)

Policy creation এবং payment API-তে `idempotency_key` (UUID) mandatory। Storage: `idempotency_keys` table-এ 24 ঘন্টা retain করতে হবে। Same key-এ second request → cached response return, no side effects।

## 7.3 Concurrency Handling

Policy sequence number generation-এ race condition prevent করতে PostgreSQL sequence ব্যবহার:
```sql
CREATE SEQUENCE policy_sequence START 1;
SELECT nextval('policy_sequence');
```
Entity update-এ optimistic concurrency: EF Core `[ConcurrencyToken]` + `rowversion` column।

## 7.4 Cache Invalidation (FR-028)

Product cache (5-min TTL) invalidation trigger: `UpdateProduct` command handler-এ `IDistributedCacheService.InvalidateProductAsync()` call করতে হবে। Redis key pattern: `product:{id}`।

## 7.5 Migration Ordering (Critical)

Foreign key dependency order অনুযায়ী migrate করতে হবে:
1. `products` (no dependency)
2. `policies` (depends on products)
3. `nominees` (depends on policies)
4. `claims` (depends on policies)
5. `claim_documents` (depends on claims)
6. `claim_approvals` (depends on claims)
7. `audit_logs` (no FK, append-only)

## 7.6 Proto Evolution

**Recommended Engineering Decision:**
- কোনো field number কখনো reuse করা যাবে না
- Breaking change মানে নতুন `v2/` package
- `optional` সব non-mandatory fields-এ লাগাতে হবে (proto3)
- Generated code `Contracts/Generated/` folder-এ manually edit করা নিষেধ

## 7.7 Claim Status Machine Terminal States

`Cancelled` এবং `Expired` policy-তে claim submit করা যাবে না। `ClaimDomainRules.EnsureClaimEligible()` এই check করবে। Grace period (30 days) চলাকালীন claim eligible।

## 7.8 Pro-rata Refund Formula (FR-053)

```
Refund = (PremiumPaid - (PremiumPaid × DaysCovered/TotalDays)) - AdminFee - CancellationCharge
```
এই calculation `Policy.CalculateProRataRefund()` domain method-এ থাকবে, handler-এ নয়।

---

# 8. AI Execution Script Mode

## Step-by-Step Implementation Instructions

### Step 1: Foundation Bootstrap

**Objective:** Runnable gRPC server with DB connection  
**Input:** appsettings.json, Program.cs template  
**Output:** Compilable project, migrations applied, health check passing  
**Dependency:** PostgreSQL running, Redis running  
**Validation Criteria:**
- `dotnet build` → 0 errors
- `dotnet ef database update` → migrations applied
- gRPC server starts on port 5001
- Health check endpoint responds 200

---

### Step 2: Domain Layer

**Objective:** All domain entities, value objects, state machines  
**Input:** SRS FG-003, FG-004, FG-005, FG-008 FR IDs  
**Output:** `Policy.cs`, `Claim.cs`, `Product.cs`, all Value Objects, State Machines  
**Dependency:** Step 1  
**Validation Criteria:**
- All domain unit tests pass
- Policy state machine rejects invalid transitions
- Claim approval level correctly derived from amount
- Pro-rata refund calculation correct for edge cases

---

### Step 3: Proto Contracts

**Objective:** All service definitions compiled  
**Input:** SRS Appendix A proto definitions  
**Output:** `Contracts/Generated/*.cs` files via `protoc`  
**Dependency:** Step 2  
**Validation Criteria:**
- `protoc` compile without errors
- `PolicyService`, `ClaimsService`, `ProductService` RPCs all defined
- `NomineeDto.income_range` is `optional`

---

### Step 4: CreatePolicy Feature (First Slice)

**Objective:** Complete CreatePolicy vertical slice working end-to-end  
**Input:** Command + Handler + Validator + Repository + gRPC mapping  
**Output:** `CreatePolicy` RPC functional, integration test passing  
**Dependency:** Steps 1-3  
**Validation Criteria:**
- Idempotency: same key → same response
- Duplicate detection (30-day window): returns error
- Policy number format: `LBT-2025-XXXX-NNNNNN`
- Kafka event published after commit
- No Kafka publish if DB transaction fails

---

### Step 5: SubmitClaim Feature

**Objective:** Claim submission with fraud evaluation  
**Input:** Claims domain + Fraud rules + AI Engine gRPC client  
**Output:** `SubmitClaim` RPC functional, fraud flags working  
**Dependency:** Step 4 (needs active policy)  
**Validation Criteria:**
- FR-082: Inactive policy → error
- FR-182: Claim <48hr → fraud flagged
- FR-083: CLM-YYYY-XXXX-NNNNNN format + SHA-256 hash
- Approval level correctly determined by amount

---

### Step 6: ApproveClaim Feature

**Objective:** Tiered approval with joint approval for L3  
**Input:** Approval matrix (SRS table), roles from JWT  
**Output:** `ApproveClaim` RPC with role-based routing  
**Dependency:** Step 5  
**Validation Criteria:**
- L1 BDT 0-10K: Auto/Officer can approve
- L3 BDT 50K-2L: Both BusinessAdmin AND FocalPerson required
- Wrong role → `PermissionDenied` gRPC status
- On approval → `ClaimApprovedEvent` published to Kafka

---

### Step 7: Renewals + Cancellation

**Objective:** Policy lifecycle completion  
**Input:** FR-043, FR-047, FR-048, FR-051, FR-052, FR-053  
**Output:** Renewal reminders background job, cancellation with pro-rata refund  
**Dependency:** Step 4  
**Validation Criteria:**
- Renewal reminder fires at D-30, D-7, D-1
- Grace period (30 days) state correctly set
- Auto-lapse after grace period
- Pro-rata refund amount calculation correct

---

### Step 8: Fraud Detection Module

**Objective:** All 4 primary fraud rules implemented  
**Input:** FR-182 to FR-185 (SRS Fraud Detection Rules table)  
**Output:** `FraudEvaluationService` with all rules  
**Dependency:** Step 5  
**Validation Criteria:**
- FR-182: <48hr → flagged
- FR-183: >2 same type in 12mo → flagged
- FR-184: 100% coverage match → flagged
- FR-185: Non-network provider → flagged
- Rules evaluate independently (no short-circuit without reason)

---

### Step 9: Product Catalog + Caching

**Objective:** Product CRUD with Redis caching  
**Input:** FR-021 to FR-029  
**Output:** All ProductService RPCs, Redis cache with 5-min TTL  
**Dependency:** Step 1  
**Validation Criteria:**
- Second product read → cache hit (Redis)
- Product update → cache invalidated
- Product search <500ms
- Assessment questions returned with product (FR-023-B)

---

### Step 10: Observability + Production Hardening

**Objective:** Tracing, metrics, error handling complete  
**Input:** Jaeger, Prometheus config  
**Output:** All interceptors, health checks, circuit breakers in place  
**Dependency:** All previous steps  
**Validation Criteria:**
- Every gRPC call traceable in Jaeger
- Prometheus metrics exported
- AI Engine circuit breaks after 3 failures
- Unhandled exceptions → `Internal` gRPC status, never expose stack trace
- Audit log written for policy issue, claim approval, cancellation

---

# 9. Implementation Gap Register (From Phase Scripts)

| Gap | Priority | SRS FR | Current Status | Required Action |
|---|---|---|---|---|
| Real Kafka integration | P1-Critical | FR-136 | Mock only | Replace KafkaMock with real Confluent.Kafka |
| PDF Generation | P1-Critical | FR-035 | Mock | Implement gRPC call to OCR Service (4002) |
| Pro-rata refund (C#) | P1-Critical | FR-053 | Go backend only | Implement in `Policy.CalculateProRataRefund()` |
| Grace period workflow | P1-High | FR-047 | Partial entity | Complete state transition + daily job |
| Fraud pattern analysis | P2 | FR-176 | Not implemented | `FrequentClaimRule.cs` |
| Endorsement doc generation | P2 | FR-059 | Not implemented | gRPC call to OCR/PDF Service |
| Document OCR verification | P3 | FR-087 | Not implemented | AI Engine integration |
| Zero Human Touch Claims | Future | FR-093 | Not implemented | ML integration required |

---

*Blueprint Version: 1.0 | Source: SRS v3.11 + Phase 1/2/3 Scripts | Architecture: VSA + CQRS/MediatR | C# .NET 8*
