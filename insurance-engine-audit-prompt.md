# Insurance Engine - Enterprise Technical Audit Prompt

You must perform a **deep enterprise technical audit** of only the **Insurance Engine module** inside my repository.

This project follows a **protobuf / gRPC based multi-microservice architecture**, and the absolute source of truth is:

⚠️ **Proto Files (The Communication Contract)**

Proto files define the API, the data structures, and the cross-service boundaries. The backend implementation and database schema must strictly follow the definitions in the Protos.

If code, database, and proto differ: **Proto wins**.

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

7. **Claims Management** (FG-031 to FG-040)
8. **Fraud Detection** (FG-041 to FG-044)
9. **Beneficiary** (FR-032 / FG-045)

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
│   ├── insuretech\                                (Namespace folder)
│   │   ├── products\
│   │   ├── policy\
│   │   ├── claims\
│   │   ├── underwriting\
│   │   ├── renewal\
│   │   ├── endorsement\
│   │   └── fraud\
│   └── common\
│
├── gen\                                            (Root-level generated code)
│   └── csharp\
│       ├── InsuranceEngine\
│       │   ├── Products\
│       │   ├── Policies\
│       │   └── [other modules]\
│
└── backend\
    └── insurance_engine\
        └── src\
            ├── InsuranceEngine.SharedKernel\       (Shared infrastructure)
            │   └── Persistence\
            │       ├── Entities\                   (Core shared entities - derived from Protos)
            │       │   ├── ProductEntity.cs
            │       │   ├── PolicyEntity.cs
            │       │   ├── ClaimEntity.cs
            │       │   └── [other core entities]
            │       └── DbContext\                  (EF Core context)
            │
            └── InsuranceEngine.[Module]\           (Module-specific projects)
                ├── InsuranceEngine.Products\
                ├── InsuranceEngine.Policy\
                ├── InsuranceEngine.Claims\
                ├── InsuranceEngine.Underwriting\
                ├── InsuranceEngine.Renewals\
                ├── InsuranceEngine.Endorsements\
                ├── InsuranceEngine.Cancellations\
                └── InsuranceEngine.FraudDetection\

Each Module Project Structure:
InsuranceEngine.[Module]\
├── GrpcServices\                                   (gRPC service implementations)
│   └── [Module]GrpcService.cs                      (Inherits from generated base)
├── Application\
│   ├── Commands\                                   (MediatR command handlers)
│   │   ├── [Action][Entity]Command.cs
│   │   └── [Action][Entity]CommandHandler.cs       (Returns Proto Response)
│   ├── Queries\                                    (MediatR query handlers)
│   │   ├── Get[Entity]Query.cs
│   │   └── List[Entities]Query.cs
│   └── Validators\                                 (FluentValidation)
└── Domain\                                         (Module-specific domain aggregates)
```

---

# 🎯 Critical Source Priority (Strict Order)

## 1️⃣ Proto Files (Communication Contract - HIGHEST TRUTH) ⚠️

**Location:** `D:\InsureTech\proto\insuretech\`

Proto files define the cross-service contract and the expected state of the system.

Treat proto as:
- Service contract (gRPC service definitions)
- Field contract (message structure)
- The blueprint for both Database Entities and Generated Classes.

---

## 2️⃣ Generated C# Classes (`gen` folder)

**Location:** `D:\InsureTech\gen\csharp\InsuranceEngine\`

Generated classes define implementation-side contract from proto. 
- Must NOT be manually modified.
- MediatR handlers MUST return these types for consistency with the API.

---

## 3️⃣ Backend Implementation (Vertical Slice Pattern)

**Location:** `D:\InsureTech\backend\insurance_engine\src\InsuranceEngine.[Module]\`

Implementation must follow the **Policy Project Pattern**:
- **gRPC Service**: Entry point, maps Proto Request -> MediatR Command.
- **MediatR Handler**: Business logic, interacts with Persistence Entities, returns Proto Response type.
- **Persistence**: SharedKernel entities mapped to database via snake_case convention.

---

### Phase 3: Module-Specific Logic & Validation

This phase ensures that the business logic implemented in the MediatR Handlers matches the **SRS v3.11** requirements and the **Proto Service Contracts**.

#### 3.1 Product Management (FG-003)
- [ ] **FR-023-A**: Verify support for unit-wise plan purchase and coverage increase/decrease.
- [ ] **FR-023-B**: Verify that each plan includes risk assessment questions.
- [ ] **FR-024**: Verify premium calculator handles dynamic inputs (age, sum assured, tenure, riders) with components breakdown.
- [ ] **Proto Match**: `CalculatePremium` RPC must return `base_premium`, `rider_premium`, and `total_premium`.

#### 3.2 Policy Lifecycle (FG-004)
- [ ] **FR-031/032**: Verify mandatory information collection (KYC, Health Declaration) and single nominee support (FR-032).
- [ ] **FR-034**: Verify policy number generation format: `LBT-YYYY-XXXX-NNNNNN`.
- [ ] **FR-039**: Verify state machine transitions: `PENDING_PAYMENT` → `ACTIVE` → `LAPSED` → `EXPIRED`.
- [ ] **Proto Match**: `IssuePolicy` RPC must validate `quote_id` and `payment_id` before transition to `ACTIVE`.

#### 3.3 Renewals & Grace Period (FG-005)
- [ ] **FR-047**: Verify 30-day grace period implementation where status remains "Grace Period" with continued coverage.
- [ ] **FR-048**: Verify auto-lapse after grace period and reinstatement within 90 days.
- [ ] **Proto Match**: `RenewPolicy` RPC must allow updating address and nominee information (FR-046).

#### 3.4 Cancellations & Endorsements (FG-005.1/2)
- [ ] **FR-052**: Verify joint approval (Business Admin + Focal Person) for policies >30 days old.
- [ ] **FR-053**: Verify pro-rata refund formula: `(Premium Paid - Days Covered - Fees)`.
- [ ] **FR-059**: Verify endorsement document suffix: `PLN-001/END-01`.
- [ ] **FR-060**: Verify approval requirement for Sum Insured changes >10%.

#### 3.5 Underwriting & Quoting (FG-004/006)
- [ ] **FR-062**: Verify handling of medical loading, occupation risk, and pre-existing conditions.
- [ ] **FR-069**: Verify reinstatement requires medical underwriting and Focal Person approval.
- [ ] **Proto Match**: `SubmitHealthDeclaration` must be completed before `ApproveUnderwriting`.

#### 3.6 Claims Management (FG-008)
- [ ] **FR-086/542-545**: Verify Tiered Approval Matrix:
    - BDT 0–10K: Officer (24h TAT)
    - BDT 10K–50K: Manager (3d TAT)
    - BDT 50K–2L: Joint (BA + FP) (7d TAT)
    - BDT 2L+: Board (15d TAT)
- [ ] **FR-100**: Verify calculation of deductibles and co-payment.
- [ ] **Proto Match**: `SettleClaim` must record `payment_method` and `payment_reference`.

#### 3.7 Fraud Detection & AML (FG-008/TM Rules)
- [ ] **FR-094**: Verify fraud flags: frequent claims (>3 in 6 months), rapid policy-to-claim (<48hrs).
- [ ] **TM-001/008**: Verify AML triggers: Multiple transactions < threshold, >3 policies in 7 days.
- [ ] **SEC-022**: Verify AML Triggers (e.g., Premium > BDT 5L without income proof).

#### 3.8 Beneficiary & KYC (FR-032)
- [ ] **FR-032-A**: Verify that beneficiary income range is optional.
- [ ] **KYC Sync**: Verify NID verification status is tracked before policy issuance.

---

## 4️⃣ Database Schema (Persistence Layer)

**Location:** `D:\InsureTech\backend\insurance_engine\src\InsuranceEngine.SharedKernel\Persistence\`

Database must be kept in sync with the Protos.
- Table columns must match Proto fields.
- Data types must be compatible (e.g., Money/decimal stored as long in paisa).
- Nullable rules in Protos must reflect in the DB schema.

**Note:** If Protos change, the database must be updated to match.

---

## 5️⃣ API Documentation (If Available)

**Location:** `D:\InsureTech\documentation\` or `D:\InsureTech\backend\insurance_engine\doc\`

Treat any HTML/Markdown API docs as contract:
- Endpoint request structure
- Endpoint response structure
- Validation rules
- Status codes

---

## 6️⃣ SRS Documentation

**Location:** `D:\InsureTech\documentation\SRS_v3\LabAid_InsureTech_SRS_v3.11.md`

Treat SRS as business definition only.

---

# 🔍 Primary Audit Objective

Determine whether these 5 layers are fully aligned:

```
Proto (Truth) ↔ Generated C# ↔ Backend Code ↔ Migration/Database ↔ HTML Contract
```

Any mismatch must be reported with:
1. Exact file paths
2. Severity level
3. Recommended fix

---

# 📋 Audit Steps

## Step 1 — Proto Contract Audit (HIGHEST PRIORITY) ⚠️

**Priority:** P0 (Highest)
**Location:** `D:\InsureTech\proto\insuretech\`

Analyze the absolute source of truth first. Extract the contract:

### What to Extract:
- Every Service definition and RPC method.
- Every Request and Response message structure.
- Field types, repeated fields, and enums.

### Then Verify Sync with Other Layers:

#### Proto vs Backend Implementation
**Compare:** Protos vs `InsuranceEngine.[Module]\GrpcServices\*GrpcService.cs`
Check:
- RPC methods map to MediatR commands correctly.
- MediatR handlers return the generated Proto response types (Mandatory!).

#### Proto vs Persistence Entities
**Compare:** Protos vs `InsuranceEngine.SharedKernel\Persistence\Entities\*.cs`
Check:
- Entity properties exist for every required field in Protos.
- Property types match (e.g., long for currency).
- Nullable properties match optional fields in Protos.

---

## Step 2 — Persistence Layer Audit (Sync Check)

**Location:** `D:\InsureTech\backend\insurance_engine\src\InsuranceEngine.SharedKernel\Persistence\`

Verify that the database schema and EF Core configurations stay in sync with the Proto-derived entities.

### Check for:
- ✅ Tables and columns match entity structure (derived from Proto).
- ✅ Snake_case naming convention applied correctly.
- ✅ Data consistency (e.g., decimal values stored as long in paisa).

⚠️ **Mismatch here causes runtime persistence errors.**

---

## Step 3 — Generated C# Integrity Audit

**Location:** `D:\InsureTech\gen\csharp\InsuranceEngine\`

Verify generated code usage:
- ✅ Generated classes reused correctly in gRPC services.
- ❌ Duplicate manual models created (anti-pattern).
- ❌ Manual modifications to generated files.

---

## Step 4 — Backend Implementation Audit

**Locations:** `InsuranceEngine.[Module]\`

Analyze:
### gRPC Service Implementations
- Inherits from generated `*ServiceBase`.
- Implements all proto-defined RPC methods.
- Maps Proto Requests to MediatR Commands.

### MediatR Command/Query Handlers
- **Crucial:** Returns the Proto Response type.
- Performs business logic using SharedKernel entities.
- Published Kafka events for cross-service workflow.

---

## Step 5 — SRS & FR Audit (Module Specific Details)

Analyze each module against its specific SRS Functional Requirements. Verify that the business logic in MediatR handlers, validation rules, and persistence entities supports these requirements.

### Module-wise FR Mapping:

#### 1. Product Catalog (FG-001 to FG-008)
- FG-001: Product definition capability (Core entity)
- FG-002: Product configuration management
- FG-003: Pricing rules management (MediatR logic)
- FG-004: Coverage limits definition
- FG-005: Exclusions management
- FG-006: Product activation/deactivation
- FG-007: Product versioning
- FG-008: Product catalog querying (Proto & DB sync)

#### 2. Underwriting (FG-009 to FG-012)
- FG-009: Application submission
- FG-010: Risk assessment rules
- FG-011: Premium calculation logic (MediatR handler)
- FG-012: Approval/Rejection decision workflow

#### 3. Policy Lifecycle (FG-013 to FG-016)
- FG-013: Policy issuance
- FG-014: Policy information management
- FG-015: Policy status tracking (State machine)
- FG-016: Policy document generation request (Kafka event)

#### 4. Renewal (FG-019 to FG-022)
- FG-019: Renewal notice generation (Kafka event)
- FG-020: Renewal premium calculation
- FG-021: Renewal processing
- FG-022: Lapsed policy handling

#### 5. Endorsement (FG-023 to FG-026)
- FG-023: Endorsement request handling
- FG-024: Coverage modification
- FG-025: Premium adjustment calculation
- FG-026: Endorsement approval workflow

#### 6. Cancellation (FG-027 to FG-030)
- FG-027: Cancellation request processing
- FG-028: Refund calculation logic (MediatR handler)
- FG-029: Cancellation reason tracking
- FG-030: Policy termination handling

#### 7. Claims Management (FG-031 to FG-040)
- FG-031: Claim registration
- FG-032: Document verification requirements
- FG-033: Claim investigation workflow
- FG-034: Loss assessment
- FG-035: Approval matrix (L1/L2/L3/Board)
- FG-036: Claim approval decision
- FG-037: Claim rejection with reasons
- FG-038: Settlement calculation
- FG-039: Claim status tracking
- FG-040: Claim history management

#### 8. Fraud Detection (FG-041 to FG-044)
- FG-041: Fraud pattern detection rules
- FG-042: Risk scoring for claims
- FG-043: Alert generation (Kafka event)
- FG-044: Investigation workflow

---

## Step 6 — Kafka Event Contract Audit
- `insurance.policy.issued`, `renewed`, `cancelled`, etc.
- Verify payload structure matches Proto definitions.

---

## Step 7 — State Machine Audit
- Policy States: Draft → Active → Renewed / Lapsed / Cancelled.
- Claim States: Submitted → UnderReview → Approved/Rejected → Settled.

---

## Step 8 — Module Completion Audit
Determine implementation status of all 8 modules based on the Proto+VSA standard.

---

# 📊 Final Enterprise Audit Report Structure

Generate comprehensive markdown report including:
1. Executive Summary (System Health %)
2. Proto vs Code Compliance
3. Backend Architecture Consistency
4. Database/Persistence Sync
5. Multi-Module Completion Status
6. Action Plan (P0 to P3 issues)

---

# 🎯 Completion Metrics Required

Provide detailed percentage calculations in the final report:

```
Proto Alignment = (Implemented RPCs / Defined RPCs) × 100
MediatR Pattern Compliance = (Handlers returning Proto Response / Total Handlers) × 100
Persistence Sync = (Matching Columns in DB vs Proto Fields / Total Proto Fields) × 100
Module Maturity = Average of all module completion %
SRS Compliance = (Implemented FRs / Total FRs) × 100
```

---

# ⚠️ Critical Rules

## Every Finding Must Include:
1. **Category:** [Proto/Code/Persistence/Event/Business Rule]
2. **Severity:** [P0 Critical / P1 High / P2 Medium / P3 Low]
3. **Module:** [Module Name]
4. **File Path:** [Exact full path]
5. **Issue:** [Specific mismatch description]
6. **Impact:** [Exact runtime error risk]
7. **Fix:** [Exact change needed]

## No Assumptions Allowed:
- ✅ Verify everything against actual files.
- ❌ Do not assume code matches proto.

---

# 📁 File Locations Summary

```
PROTO REPOSITORY:
D:\InsureTech\proto\insuretech\

GENERATED CODE:
D:\InsureTech\gen\csharp\InsuranceEngine\

PERSISTENCE (Entities & DB Context):
D:\InsureTech\backend\insurance_engine\src\InsuranceEngine.SharedKernel\Persistence\

MODULE PROJECTS:
D:\InsureTech\backend\insurance_engine\src\InsuranceEngine.[Module]\
├── GrpcServices\*GrpcService.cs
├── Application\
│   ├── Commands\
│   └── Queries\
└── Domain\

DOCUMENTATION:
D:\InsureTech\documentation\SRS_v3\LabAid_InsureTech_SRS_v3.11.md
D:\InsureTech\documentation\
```

---

**END OF AUDIT PROMPT**
