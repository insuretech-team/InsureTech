# Insurance Engine - Enterprise Technical Audit Prompt

You must perform a **deep enterprise technical audit** of only the **Insurance Engine module** inside my repository.

This project follows a **protobuf / gRPC based multi-microservice architecture**, but the most critical runtime truth is:

⚠️ **Database schema created through migrations**

Because runtime stability depends on migration-created tables.

If code, proto, generated classes, and database differ: **runtime errors will happen**.

Therefore migration-created schema must be treated as highest runtime validation layer.

---

# 🏢 Project Context

## Project Information

- **Project:** LabAid InsureTech Platform
- **Service:** Insurance Engine (Port 5001) ← **YOUR SCOPE**
- **Technology:** .NET 10, PostgreSQL, Kafka, Redis, gRPC
- **Architecture:** Vertical Slice Architecture (VSA) + CQRS (MediatR)
- **Partner:** LifePlus

## ⚠️ AUDIT SCOPE: Insurance Engine ONLY

This audit covers **ONLY the Insurance Engine microservice**, not the entire platform.

### Insurance Engine's 8 Modules (YOUR RESPONSIBILITY):

1. **Product Catalog** (FG-001 to FG-008)
2. **Underwriting** (FG-009 to FG-012)
3. **Policy Lifecycle** (FG-013 to FG-016)
4. **Renewal** (FG-019 to FG-022)
5. **Endorsement** (FG-023 to FG-026)
6. **Cancellation** (FG-027 to FG-030)
7. **Claims Management** (FG-031 to FG-040)
8. **Fraud Detection** (FG-041 to FG-044)

### Other Platform Services (NOT YOUR SCOPE):

- ❌ Payment Service (Port 5002)
- ❌ Analytics & Reporting Service (Port 5003) - includes FG-017, FG-018
- ❌ Document Service (Port 5004)
- ❌ Notification Service (Port 5005)
- ❌ Customer Service (Port 5006)

**Insurance Engine interacts with these services via:**

- Kafka events (publish/subscribe)
- gRPC calls (when needed)

But does NOT implement their core business logic.

## Repository Structure

```
D:\InsureTech\
├── proto\                                          (Root-level shared proto repository)
│   ├── insurance_engine\
│   │   ├── products.proto
│   │   ├── policies.proto
│   │   ├── claims.proto
│   │   ├── underwriting.proto
│   │   └── [other modules].proto
│   └── common\
│       └── shared_types.proto
│
├── gen\                                            (Root-level generated code)
│   └── csharp\
│       ├── InsuranceEngine\
│       │   ├── Products\
│       │   ├── Policies\
│       │   └── [other modules]\
│       └── Common\
│
└── backend\
    └── insurance_engine\
        └── src\
            ├── InsuranceEngine.SharedKernel\       (Shared infrastructure)
            │   └── Persistence\
            │       ├── Entities\                   (Core shared entities)
            │       │   ├── Product.cs
            │       │   ├── Policy.cs
            │       │   ├── Claim.cs
            │       │   └── [other core entities]
            │       └── Migrations\                 (EF Core migrations - HIGHEST PRIORITY)
            │           ├── 20250101_InitialCreate.cs
            │           └── [subsequent migrations]
            │
            └── InsuranceEngine.[Module]\           (Module-specific projects)
                ├── InsuranceEngine.Products\
                ├── InsuranceEngine.Policies\
                ├── InsuranceEngine.Claims\
                ├── InsuranceEngine.Underwriting\
                ├── InsuranceEngine.Renewals\
                ├── InsuranceEngine.Endorsements\
                ├── InsuranceEngine.Cancellations\
                └── InsuranceEngine.FraudDetection\

Each Module Project Structure:
InsuranceEngine.[Module]\
├── GrpcServices\                                   (gRPC service implementations)
│   └── [Module]ServiceImpl.cs
├── Application\
│   ├── Commands\                                   (MediatR command handlers)
│   │   ├── Create[Entity]\
│   │   │   ├── Create[Entity]Command.cs
│   │   │   ├── Create[Entity]Handler.cs
│   │   │   └── Create[Entity]Validator.cs        (May be here or in Validators folder)
│   │   └── [other commands]
│   ├── Queries\                                    (MediatR query handlers)
│   │   ├── Get[Entity]ById\
│   │   └── List[Entities]\
│   └── Validators\                                 (If separate from Commands)
│       └── [validators if not in command folders]
└── Domain\                                         (Module-specific domain entities if any)
    └── [module-specific entities]
```

---

# 🎯 Critical Source Priority (Strict Order)

## 1️⃣ Migration-Created Database Schema (Highest Runtime Truth) ⚠️

**Location:** `D:\InsureTech\backend\insurance_engine\src\InsuranceEngine.SharedKernel\Persistence\Migrations\`

The database has already been created using project migrations.

This is the most important runtime truth.

Audit must first verify:

- Which tables were actually created
- Exact columns created
- Data types
- Nullable rules
- Foreign keys
- Constraints
- Indexes
- Default values

**Why Critical:** If code differs from migration, runtime errors will occur.

---

## 2️⃣ Proto Files (Communication Contract)

**Location:** `D:\InsureTech\proto\insurance_engine\*.proto`

Proto files define cross-service contract.

Treat proto as:

- Service contract (gRPC service definitions)
- Field contract (message structure)
- Cross-team compatibility contract
- Version control contract

---

## 3️⃣ Generated C# Classes (`gen` folder)

**Location:** `D:\InsureTech\gen\csharp\InsuranceEngine\`

Generated classes define implementation-side contract from proto.

Verify:

- Proper namespace usage
- No manual modifications to generated code
- Version consistency with proto

---

## 4️⃣ Backend Implementation

**Locations:**

- **Core Entities:** `D:\InsureTech\backend\insurance_engine\src\InsuranceEngine.SharedKernel\Persistence\Entities\`
- **gRPC Services:** `D:\InsureTech\backend\insurance_engine\src\InsuranceEngine.[Module]\GrpcServices\`
- **Commands:** `D:\InsureTech\backend\insurance_engine\src\InsuranceEngine.[Module]\Application\Commands\`
- **Queries:** `D:\InsureTech\backend\insurance_engine\src\InsuranceEngine.[Module]\Application\Queries\`
- **Validators:** `D:\InsureTech\backend\insurance_engine\src\InsuranceEngine.[Module]\Application\Validators\` or inside `Commands\` folder

---

## 5️⃣ HTML Documentation (`doc/*.html`)

**Location:** `D:\InsureTech\backend\insurance_engine\doc\*.html`

Treat HTML as exact API contract:

- Endpoint request structure
- Endpoint response structure
- Validation rules
- Status codes
- Error responses

---

## 6️⃣ SRS Documentation

**Location:** `D:\InsureTech\backend\insurance_engine\doc\SRS_V3\LabAid_InsureTech_SRS_v3.11.md`

Treat SRS as business definition only.

⚠️ **CRITICAL SCOPE LIMITATION:**

The SRS document covers the entire LabAid InsureTech Platform, which includes multiple microservices:

- Insurance Engine (Port 5001) ← **YOUR RESPONSIBILITY**
- Payment Service (Port 5002)
- Analytics & Reporting Service (Port 5003)
- Document Service (Port 5004)
- Notification Service (Port 5005)
- Customer Service (Port 5006)
- Other services...

**AUDIT MUST ONLY VERIFY INSURANCE ENGINE SCOPE:**

### Insurance Engine Functional Requirements (FR):

Audit ONLY these FR IDs from SRS:

- **FG-001 to FG-008:** Product Catalog Module
- **FG-009 to FG-012:** Underwriting Module
- **FG-013 to FG-016:** Policy Lifecycle Module
- **FG-019 to FG-022:** Renewal Module
- **FG-023 to FG-026:** Endorsement Module
- **FG-027 to FG-030:** Cancellation Module
- **FG-031 to FG-040:** Claims Management Module
- **FG-041 to FG-044:** Fraud Detection Module

### OUT OF SCOPE (DO NOT AUDIT):

- ❌ **FG-017, FG-018:** Reporting (belongs to Analytics Service - Port 5003)
- ❌ Payment gateway integration (belongs to Payment Service - Port 5002)
- ❌ Document storage/retrieval (belongs to Document Service - Port 5004)
- ❌ Email/SMS sending (belongs to Notification Service - Port 5005)
- ❌ Customer profile management (belongs to Customer Service - Port 5006)

### Insurance Engine Business Logic ONLY:

Verify implementation of:
✅ Product definition and management
✅ Risk assessment and underwriting decisions
✅ Policy issuance and lifecycle management
✅ Renewal calculation and processing
✅ Endorsement/modification handling
✅ Cancellation and refund calculation
✅ Claims processing workflow and approval matrix
✅ Fraud detection rules and scoring

### Cross-Service Boundaries:

Insurance Engine should:

- **Publish events** to Kafka (for other services to consume)
- **Call other services** via gRPC when needed
- **NOT implement** other services' core business logic

Example:

- ✅ Insurance Engine calculates refund amount (its logic)
- ❌ Insurance Engine processes payment refund (Payment Service's job)
- ✅ Insurance Engine publishes `insurance.policy.cancelled` event with refund amount
- ✅ Payment Service consumes event and processes refund

---

# 🔍 Primary Audit Objective

Determine whether these 5 layers are fully aligned:

```
Migration Schema ↔ Proto ↔ Generated C# ↔ Backend Code ↔ HTML Contract
```

Any mismatch must be reported with:

1. Exact file paths
2. Severity level
3. Runtime impact
4. Recommended fix

---

# 📋 Audit Steps

## Step 1 — Migration First Audit (Most Critical) ⚠️

**Priority:** P0 (Highest)

**Location:** `D:\InsureTech\backend\insurance_engine\src\InsuranceEngine.SharedKernel\Persistence\Migrations\`

Analyze migrations first. Extract exact schema:

### What to Extract:

- All tables created
- All columns with exact names
- Data types (int, varchar, jsonb, etc.)
- Nullable vs NOT NULL
- Primary keys
- Foreign keys
- Unique constraints
- Check constraints
- Default values
- Indexes

### Then Verify:

#### Migration vs Entity Classes

**Compare:** Migration schema vs `InsuranceEngine.SharedKernel\Persistence\Entities\*.cs`

Check:

- Entity property matches actual table column কিনা
- Property type matches column type কিনা (e.g., `string` vs `varchar(255)`)
- Nullable property vs nullable column mismatch আছে কিনা
- Navigation properties match foreign keys কিনা
- Missing properties আছে কিনা
- Extra properties আছে কিনা (যা table এ নেই)

#### Migration vs Proto

**Compare:** Migration schema vs `D:\InsureTech\proto\insurance_engine\*.proto`

Check:

- Proto message fields match table columns কিনা
- Field types compatible কিনা (proto int32 vs SQL integer)
- Optional/required fields match nullable columns কিনা
- Repeated fields properly handled কিনা

#### Migration vs Generated C#

**Compare:** Migration schema vs `D:\InsureTech\gen\csharp\InsuranceEngine\**\*.cs`

Check:

- Generated class properties match table structure কিনা
- Enum types match database enum/varchar কিনা
- Optional field handling correct কিনা

⚠️ **This is highest priority because migration mismatch causes direct runtime failure.**

---

## Step 2 — Proto Contract Audit

**Location:** `D:\InsureTech\proto\insurance_engine\*.proto`

Analyze all proto files. Identify:

### Proto Service Definitions:

- Service names
- RPC methods
- Request message types
- Response message types
- Stream vs unary

### Proto Messages:

- Message definitions
- Field names and types
- Optional vs required fields
- Repeated fields
- Nested messages
- Enums
- Oneof fields

### Verify Backend Follows Proto Contract:

#### Proto Service → gRPC Service Implementation

**Compare:** Proto services vs `InsuranceEngine.[Module]\GrpcServices\*ServiceImpl.cs`

Check:

- All RPC methods implemented কিনা
- Method signatures match কিনা
- Request/Response types correct কিনা
- Return types match কিনা
- Missing methods আছে কিনা

#### Proto Messages → Usage in Code

**Compare:** Proto messages vs Command/Query DTOs

Check:

- Proto messages directly used কিনা (preferred)
- Manual DTOs created unnecessarily কিনা (duplication)
- Mapping logic exists if custom DTOs used

---

## Step 3 — Generated C# Integrity Audit

**Location:** `D:\InsureTech\gen\csharp\InsuranceEngine\`

Verify generated code usage:

### Check for:

- ✅ Generated classes reused correctly in gRPC services
- ❌ Duplicate manual models created (anti-pattern)
- ❌ Namespace conflicts with generated code
- ❌ Manual modifications to generated files
- ✅ Proper dependency injection of generated stubs

### Specific Checks:

```csharp
// ✅ GOOD: Using generated code
using InsuranceEngine.Products;

public class ProductServiceImpl : ProductService.ProductServiceBase
{
    public override Task<CreateProductResponse> CreateProduct(
        CreateProductRequest request, ...)
}

// ❌ BAD: Manual DTO duplication
public class CreateProductRequestDto  // Generated version exists!
{
    public string Name { get; set; }
    ...
}
```

---

## Step 4 — Backend Implementation Audit

**Locations:** Multiple directories per module

Analyze:

### gRPC Service Implementations

**Location:** `InsuranceEngine.[Module]\GrpcServices\*ServiceImpl.cs`

Check:

- Inherits from generated `*.ProductServiceBase`
- Implements all proto-defined RPC methods
- Delegates to MediatR handlers
- Proper error handling
- ServerCallContext usage

### MediatR Command/Query Handlers

**Locations:**

- `InsuranceEngine.[Module]\Application\Commands\`
- `InsuranceEngine.[Module]\Application\Queries\`

Check:

- Command/Query → Handler mapping
- Validation before handler execution
- Entity retrieval from repository
- Business logic implementation
- Event publishing (Kafka)

### Repositories

Check:

- Entity CRUD operations
- EF Core DbContext usage
- Query optimization
- Transaction handling

### Validators

**Location:** `InsuranceEngine.[Module]\Application\Validators\` or inside `Commands\`

Check:

- FluentValidation rules exist
- Validation matches proto constraints
- Validation matches database constraints
- Business rule validation

### DI Registrations

Check:

- gRPC services registered
- MediatR registered
- Validators registered
- Repositories registered

Build complete architecture map.

---

## Step 5 — HTML API Contract Audit

**Location:** `D:\InsureTech\doc\*.html`

Verify for each documented endpoint:

### Request Contract:

- Endpoint path matches gRPC method কিনা
- Request parameters documented correctly কিনা
- Request body structure matches proto message কিনা
- Required vs optional fields match কিনা

### Response Contract:

- Response structure matches proto message কিনা
- Success response documented correctly কিনা
- Error responses documented কিনা
- Status codes correct কিনা

### Validation Rules:

- Documented validation matches FluentValidation rules কিনা
- Field constraints match database constraints কিনা

---

## Step 5.5 — SRS Business Requirements Audit

**Location:** `D:\InsureTech\documentation\SRS_V3\LabAid_InsureTech_SRS_v3.11.md`

**⚠️ CRITICAL: Audit ONLY Insurance Engine FR IDs**

### Module-wise FR Mapping:

#### 1. Product Catalog Module

**SRS FR IDs:** FG-001 to FG-008
Verify implementation of:

- FG-001: Product definition capability
- FG-002: Product configuration management
- FG-003: Pricing rules management
- FG-004: Coverage limits definition
- FG-005: Exclusions management
- FG-006: Product activation/deactivation
- FG-007: Product versioning
- FG-008: Product catalog querying

#### 2. Underwriting Module

**SRS FR IDs:** FG-009 to FG-012
Verify implementation of:

- FG-009: Application submission
- FG-010: Risk assessment rules
- FG-011: Premium calculation logic
- FG-012: Approval/Rejection decision workflow

#### 3. Policy Lifecycle Module

**SRS FR IDs:** FG-013 to FG-016
Verify implementation of:

- FG-013: Policy issuance
- FG-014: Policy information management
- FG-015: Policy status tracking
- FG-016: Policy document generation (Insurance Engine creates request, Document Service generates)

#### 4. Renewal Module

**SRS FR IDs:** FG-019 to FG-022
Verify implementation of:

- FG-019: Renewal notice generation
- FG-020: Renewal premium calculation
- FG-021: Renewal processing
- FG-022: Lapsed policy handling

#### 5. Endorsement Module

**SRS FR IDs:** FG-023 to FG-026
Verify implementation of:

- FG-023: Endorsement request handling
- FG-024: Coverage modification
- FG-025: Premium adjustment calculation
- FG-026: Endorsement approval workflow

#### 6. Cancellation Module

**SRS FR IDs:** FG-027 to FG-030
Verify implementation of:

- FG-027: Cancellation request processing
- FG-028: Refund calculation logic
- FG-029: Cancellation reason tracking
- FG-030: Policy termination handling

#### 7. Claims Management Module

**SRS FR IDs:** FG-031 to FG-040
Verify implementation of:

- FG-031: Claim registration
- FG-032: Document verification requirements
- FG-033: Claim investigation workflow
- FG-034: Loss assessment
- FG-035: Approval matrix (L1/L2/L3/Board)
- FG-036: Claim approval decision
- FG-037: Claim rejection with reasons
- FG-038: Settlement calculation (Insurance Engine calculates, Payment Service disburses)
- FG-039: Claim status tracking
- FG-040: Claim history management

#### 8. Fraud Detection Module

**SRS FR IDs:** FG-041 to FG-044
Verify implementation of:

- FG-041: Fraud pattern detection rules
- FG-042: Risk scoring for claims
- FG-043: Alert generation
- FG-044: Investigation workflow

### ❌ OUT OF SCOPE (Do NOT Audit These):

- **FG-017, FG-018:** Report generation (Analytics Service responsibility)
- **FG-045+:** Payment processing (Payment Service responsibility)
- **FG-050+:** Document storage (Document Service responsibility)
- **FG-055+:** Notification sending (Notification Service responsibility)

### For Each FR, Verify:

1. ✅ Business rule implemented in code
2. ✅ Validation rules match SRS specification
3. ✅ Database schema supports the requirement
4. ✅ Proto contract includes necessary fields
5. ✅ API endpoint exists (if applicable)

### Cross-Service Interaction Verification:

When Insurance Engine needs other services:

- ✅ Events published to Kafka (don't implement the consuming service's logic)
- ✅ gRPC calls made to other services (when synchronous needed)
- ❌ Should NOT duplicate other services' business logic

Example Correct Pattern:

```
Insurance Engine (Cancellation Module):
1. Calculates refund amount ✅ (FG-028 - its logic)
2. Publishes "insurance.policy.cancelled" event with refund_amount ✅
3. Payment Service consumes event and processes refund ✅ (not Insurance Engine's job)
```

---

## Step 6 — Kafka Event Contract Audit

Verify Kafka event schemas:

### Event Topics to Check:

- `insurance.policy.issued`
- `insurance.policy.renewed`
- `insurance.policy.cancelled`
- `insurance.policy.endorsed`
- `insurance.claim.filed`
- `insurance.claim.approved`
- `insurance.claim.settled`
- `insurance.fraud.detected`
- `insurance.premium.due`

### For Each Event:

Check:

- Event payload structure defined in proto কিনা
- Event publishing code exists কিনা
- Event payload matches current entity structure কিনা
- Event versioning strategy আছে কিনা
- Consumer contract breach risk আছে কিনা

---

## Step 7 — State Machine Alignment Audit

Verify state transition logic:

### Policy States:

```
Draft → Active → Renewed
              ↓
         Suspended → Cancelled
              ↓
           Lapsed
```

Check:

- Enum values in code match migration-created column type কিনা
- State transition validation exists কিনা
- Invalid state transition prevented কিনা
- State change events published correctly কিনা

### Claim States:

```
Submitted → UnderReview → Investigating
                              ↓
                      Approved/Rejected
                              ↓
                         Settled/Closed
```

Check same as above.

---

## Step 8 — Business Rules Engine Audit

Verify business logic implementation:

### Premium Calculation Rules:

Check:

- Hardcoded vs database-driven rules
- Rule versioning strategy
- Base premium + risk factors calculation
- Rule conflict detection

### Risk Scoring Rules:

Check:

- Risk factor weights
- Score calculation logic
- Auto-approval thresholds

### Fraud Detection Rules:

Check:

- Detection logic implementation
- Rule configuration
- Alert generation
- False positive handling

### Approval Matrix Rules (L1/L2/L3/Board):

Check:

- Claim amount thresholds (50k, 200k, 500k BDT)
- Routing logic
- Escalation workflow
- Approval delegation

---

## Step 9 — Runtime Error Risk Detection

Detect mismatches that may cause:

### Insert Failures:

- Missing columns in entity but exist in table
- NOT NULL column without default value
- Foreign key constraint violations

### Update Failures:

- Attempting to update non-existent columns
- Type mismatch during update

### Get/Query Failures:

- Entity property not in table
- Navigation property misconfiguration
- Incorrect joins

### Mapping Exceptions:

- Proto → Entity mapping errors
- Entity → Proto response mapping errors
- Enum conversion failures

### Null Reference Exceptions:

- Non-nullable database field mapped to nullable property
- Missing null checks in handlers

### Serialization Failures:

- Complex types not properly handled in proto
- JSONB column deserialization errors

### Foreign Key Conflicts:

- Deleting referenced entities
- Inserting with invalid foreign keys

### Enum Conversion Issues:

- Proto enum → Database enum mismatch
- Missing enum values
- Invalid enum string conversion

---

## Step 10 — Beneficiary Critical Stability Check

**Context:** Beneficiary functionality currently partially works.

Must remain stable during any refactoring.

### Individual Beneficiary:

Verify these working flows remain safe:

- ✅ Create individual beneficiary
- ✅ Get all individual beneficiaries
- ✅ Get individual beneficiary by ID

### Business Beneficiary:

Verify:

- ✅ Create business beneficiary
- ✅ Get all business beneficiaries

### Hidden Risks to Detect:

- Migration change that might break these
- Proto change without code update
- Entity property change without migration
- Validation rule that might reject valid requests

---

## Step 11 — Caching Layer Alignment Audit

**Technology:** Redis (StackExchange.Redis)

### Verify:

- Cached entity structure matches migration schema কিনা
- Cache key naming conventions consistent কিনা
- Cache invalidation on entity updates করা হচ্ছে কিনা
- Stale data risk assessment
- Cache serialization format matches entity structure কিনা

### Specific Checks:

- Product catalog caching
- Premium calculation rules caching
- Fraud detection rules caching

---

## Step 12 — Dead Code Detection

Detect orphaned/unused code:

### Unused DTOs:

- Manual DTOs when generated proto messages exist
- DTOs defined but never referenced

### Unreferenced Entities:

- Entity classes not mapped in DbContext
- Entity classes not used in any query/command

### Unused Proto Messages:

- Proto messages defined but never used in any RPC

### Orphaned Migrations:

- Migrations that were partially applied
- Conflicting migrations

### Commented-Out Critical Code:

- Important validation commented out
- Business rules temporarily disabled

---

## Step 13 — Cross-Service Contract Audit

Verify gRPC contract with other services:

### Expected Integration Points:

- **Payment Service:** Premium collection, refund processing
- **Document Service:** Policy documents, claim documents storage
- **Notification Service:** Email, SMS alerts
- **Customer Service:** Customer information retrieval

### For Each Integration:

Check:

- Insurance Engine's proto matches expected contract কিনা
- Request/Response messages compatible কিনা
- Error handling for service unavailability
- Retry/resilience patterns (Polly)

---

## Step 14 — Module Completion Audit

Determine implementation status of 8 modules:

**⚠️ Map each module to its SRS FR IDs:**

### 1. Product Catalog Module

**SRS FR:** FG-001 to FG-008
Check:

- Migration: tables for products, product_configs, pricing_rules
- Proto: ProductService, messages (CreateProductRequest, Product, etc.)
- Implementation: Command/Query handlers for all CRUD operations
- SRS Compliance: All FR-001 to FR-008 business rules implemented
- Status: Fully/Partially/Missing (X% complete)

### 2. Underwriting Module

**SRS FR:** FG-009 to FG-012
Check:

- Migration: tables for underwriting_applications, risk_factors
- Proto: UnderwritingService
- Implementation: Risk assessment, premium calculation handlers
- SRS Compliance: All FR-009 to FR-012 implemented
- Status: Fully/Partially/Missing (X% complete)

### 3. Policy Lifecycle Module

**SRS FR:** FG-013 to FG-016
Check:

- Migration: tables for policies, insured_assets, vehicle_details, etc.
- Proto: PolicyService
- Implementation: Policy CRUD, status management handlers
- SRS Compliance: All FR-013 to FR-016 implemented
- Status: Fully/Partially/Missing (X% complete)

### 4. Renewal Module

**SRS FR:** FG-019 to FG-022
Check:

- Migration: renewals table
- Proto: RenewalService
- Implementation: Renewal calculation, processing handlers
- SRS Compliance: All FR-019 to FR-022 implemented
- Status: Fully/Partially/Missing (X% complete)

### 5. Endorsement Module

**SRS FR:** FG-023 to FG-026
Check:

- Migration: endorsements table
- Proto: EndorsementService
- Implementation: Modification, adjustment handlers
- SRS Compliance: All FR-023 to FR-026 implemented
- Status: Fully/Partially/Missing (X% complete)

### 6. Cancellation Module

**SRS FR:** FG-027 to FG-030
Check:

- Migration: cancellations table
- Proto: CancellationService
- Implementation: Refund calculation, termination handlers
- SRS Compliance: All FR-027 to FR-030 implemented
- Status: Fully/Partially/Missing (X% complete)

### 7. Claims Management Module

**SRS FR:** FG-031 to FG-040
Check:

- Migration: claims, claim_documents, claim_assessments tables
- Proto: ClaimsService
- Implementation: Claim processing, approval matrix handlers
- SRS Compliance: All FR-031 to FR-040 implemented
- Status: Fully/Partially/Missing (X% complete)

### 8. Fraud Detection Module

**SRS FR:** FG-041 to FG-044
Check:

- Migration: fraud_alerts, fraud_rules tables
- Proto: FraudDetectionService
- Implementation: Pattern detection, scoring handlers
- SRS Compliance: All FR-041 to FR-044 implemented
- Status: Fully/Partially/Missing (X% complete)

### For Each Module Report:

- ✅ Tables created in migration
- ✅ Proto service defined
- ✅ Generated C# classes exist
- ✅ gRPC service implemented
- ✅ Command handlers implemented (for each SRS FR)
- ✅ Query handlers implemented
- ✅ Validators implemented
- ✅ SRS business rules implemented
- ✅ HTML docs exist

**Completion Percentage:** (Implemented FRs / Total FRs) × 100

### ❌ Explicitly Exclude from Audit:

- FG-017, FG-018 (Reporting - Analytics Service)
- Payment processing FRs (Payment Service)
- Document generation FRs (Document Service)
- Notification FRs (Notification Service)

---

# 📊 Final Enterprise Audit Report Structure

Generate comprehensive markdown report with following sections:

## 1️⃣ Executive Summary

```markdown
### Overall System Health: X%

### Critical Metrics:

- Migration Alignment: X%
- Proto Compliance: X%
- Generated Code Usage: X%
- Backend Implementation: X%
- Documentation Coverage: X%
- Module Maturity: X%

### Issue Breakdown:

- P0 Critical: X issues (blocks production)
- P1 High: X issues (causes bugs)
- P2 Medium: X issues (technical debt)
- P3 Low: X issues (optimization)

### Risk Assessment:

- Runtime Error Risk: High/Medium/Low
- Regression Risk: High/Medium/Low
- Integration Risk: High/Medium/Low
```

---

## 2️⃣ Migration Schema Audit (Highest Priority)

### Tables Created:

List all tables from migrations with:

- Table name
- Columns
- Types
- Constraints
- Indexes

### Migration vs Entity Mismatches:

| Table | Entity | Issue | Impact | Fix |
| ----- | ------ | ----- | ------ | --- |

### Migration vs Proto Mismatches:

| Table | Proto Message | Issue | Impact | Fix |
| ----- | ------------- | ----- | ------ | --- |

---

## 3️⃣ Proto Contract Audit

### Proto Services Defined:

List all services with RPC methods

### Proto Implementation Coverage:

| Proto Service | gRPC Implementation | Status | Issues |
| ------------- | ------------------- | ------ | ------ |

### Proto Message Usage:

| Proto Message | Used in Code | Duplicate DTO? | Issue |
| ------------- | ------------ | -------------- | ----- |

---

## 4️⃣ Generated C# Audit

### Generated Classes Found:

List all generated classes

### Usage Analysis:

| Generated Class | Used Correctly | Manual Duplicate | Issue |
| --------------- | -------------- | ---------------- | ----- |

### Namespace Conflicts:

List any conflicts

---

## 5️⃣ Backend Architecture Audit

### gRPC Services:

| Module | Service Class | Proto Base | Implemented Methods | Missing Methods |
| ------ | ------------- | ---------- | ------------------- | --------------- |

### Command/Query Handlers:

| Module | Commands | Queries | Validators | Issues |
| ------ | -------- | ------- | ---------- | ------ |

### Entity Repository Pattern:

| Entity | Repository | CRUD Complete | Issues |
| ------ | ---------- | ------------- | ------ |

---

## 6️⃣ HTML Contract Audit

### Documented Endpoints:

| Endpoint | Implementation | Request Match | Response Match | Issues |
| -------- | -------------- | ------------- | -------------- | ------ |

---

## 7️⃣ Event Contract Audit

### Kafka Topics:

| Topic | Proto Defined | Publisher | Payload Match | Issues |
| ----- | ------------- | --------- | ------------- | ------ |

---

## 8️⃣ State Machine Audit

### Policy States:

| State | Enum Value | Migration Column | Transitions Valid | Issues |
| ----- | ---------- | ---------------- | ----------------- | ------ |

### Claim States:

Similar table...

---

## 9️⃣ Module Completion Audit

### Per-Module Breakdown:

#### Product Catalog Module

- **Completion:** X%
- **Migration:** ✅/❌ Tables created
- **Proto:** ✅/❌ Service defined
- **Implementation:** ✅/❌ Handlers complete
- **Critical Issues:** X
- **Recommendations:** ...

Repeat for all 8 modules...

---

## 🔟 Runtime Risk Audit

### P0 Critical Risks (Blocks Production):

| Risk | Location | Issue | Impact | Fix |
| ---- | -------- | ----- | ------ | --- |

### P1 High Risks (Causes Bugs):

Similar table...

### P2 Medium Risks (Technical Debt):

Similar table...

### P3 Low Risks (Optimization):

Similar table...

---

## 1️⃣1️⃣ Schema Conflict Audit

### Migration vs Code Conflicts:

Detailed list of all mismatches

---

## 1️⃣2️⃣ Integration Risk Audit

### Cross-Service Compatibility:

| External Service | Contract Status | Risk Level | Issues |
| ---------------- | --------------- | ---------- | ------ |

---

## 1️⃣3️⃣ Regression Risk Audit

### Beneficiary Stability:

| Feature | Current Status | Regression Risk | Protection Needed |
| ------- | -------------- | --------------- | ----------------- |

### Other Working Features:

List features that must remain stable

---

## 1️⃣4️⃣ Dead Code Analysis

### Unused Components:

| Type | Location | Reason | Action |
| ---- | -------- | ------ | ------ |

---

## 1️⃣5️⃣ Action Plan

### Immediate Actions (P0):

1. [File path] - Issue - Fix
2. ...

### Short-term Actions (P1):

1. [File path] - Issue - Fix
2. ...

### Long-term Actions (P2/P3):

1. [File path] - Issue - Fix
2. ...

---

# 🎯 Completion Metrics Required

Provide detailed percentage calculations:

```
Migration Alignment = (Matching Columns / Total Columns) × 100
Proto Alignment = (Implemented RPCs / Defined RPCs) × 100
Generated C# Usage = (Used Generated / Total Generated) × 100
Endpoint Compliance = (Correct Endpoints / Total Documented) × 100
Documentation Coverage = (Documented Features / Total Features) × 100
Module Maturity = Average of all module completion %
```

---

# ⚠️ Critical Rules

## Every Finding Must Include:

1. **Category:** [Migration/Proto/Code/Contract/Event/State/Business Rule]
2. **Severity:** [P0 Critical / P1 High / P2 Medium / P3 Low]
3. **Module:** [Products/Policies/Claims/Underwriting/Renewals/Endorsements/Cancellations/FraudDetection]
4. **File Path:** [Exact full path starting from D:\InsureTech\]
5. **Issue:** [Specific detailed mismatch description]
6. **Impact:** [Exact runtime error or risk description]
7. **Fix:** [Exact code change or migration change needed]
8. **Priority:** [P0/P1/P2/P3]

## No Assumptions Allowed:

- ❌ Do not assume table structure
- ❌ Do not assume proto messages exist
- ❌ Do not assume code is correct
- ✅ Verify everything against actual files

## Migration Mismatch = P0 Critical:

- Any mismatch between migration and code is **Critical** severity
- Must be fixed before any deployment

---

# 📁 File Locations Summary

```
PROTO REPOSITORY:
D:\InsureTech\proto\insurance_engine\*.proto

GENERATED CODE:
D:\InsureTech\gen\csharp\InsuranceEngine\**\*.cs

MIGRATIONS (HIGHEST PRIORITY):
D:\InsureTech\backend\insurance_engine\src\InsuranceEngine.SharedKernel\Persistence\Migrations\*.cs

CORE ENTITIES:
D:\InsureTech\backend\insurance_engine\src\InsuranceEngine.SharedKernel\Persistence\Entities\*.cs

MODULE PROJECTS:
D:\InsureTech\backend\insurance_engine\src\InsuranceEngine.[Module]\
├── GrpcServices\*ServiceImpl.cs
├── Application\
│   ├── Commands\**\*.cs
│   ├── Queries\**\*.cs
│   └── Validators\*.cs (or in Commands folder)
└── Domain\*.cs (if module-specific entities exist)

DOCUMENTATION:
D:\InsureTech\backend\insurance_engine\doc\*.html
D:\InsureTech\backend\insurance_engine\doc\SRS_V3\LabAid_InsureTech_SRS_v3.11.md
```

---

**END OF AUDIT PROMPT**

---

## Expected Output

A single comprehensive markdown file with:

- 15 major audit sections
- Detailed tables for findings
- Exact file paths
- Priority-based action plan
- Completion metrics
- Risk assessment

Minimum 50+ pages of detailed analysis expected for enterprise-grade audit.
