# InsureTech — Full Technical Audit Report

**Date:** March 13, 2025  
**Scope:** Repository-wide backend analysis with focus on **Insurance Engine** implementation and alignment to SRS, API docs, and database.

**Sources used:**
- Repository: `d:\InsureTech`
- Documentation: `documentation/SRS_v3/` (SRS_V3.11.rtf, SPECS_V3.7 sections)
- API docs: `docs/`, `api/docs/` (HTML)
- Proto: `proto/insuretech/`, `gen/csharp/`
- Migrations: `backend/Inscore/db/migrations/`, `backend/insurance_engine/migration_script.sql`

---

## 1. Full Backend Architecture Overview

### 1.1 Repository layout

| Area | Path | Description |
|------|------|-------------|
| **Backend .NET** | `backend/insurance_engine/`, `backend/polisync/` | C# .NET 8 services |
| **Backend Go** | `backend/Inscore/` | Go microservices (gateway, auth, insurance, etc.) |
| **Proto / Gen** | `proto/insuretech/`, `gen/csharp/` | gRPC contracts and generated code |
| **Documentation** | `documentation/SRS_v3/`, `documentation/Planning/` | SRS, BRD, planning |
| **API docs** | `docs/`, `api/docs/` | Generated HTML (endpoints, schemas) |
| **Infra** | `backend/infra/` | nginx, docker, scripts |
| **SDKs** | `sdks/` | Go, TypeScript, sdk-generator |

### 1.2 Backend projects and solutions

**Solution files:**
- `backend/insurance_engine/InsuranceEngine.sln`
- `backend/polisync/PoliSync.sln`

**Insurance Engine (C#):**

| Project | Path | Role |
|---------|------|------|
| InsuranceEngine.ApiHost | `backend/insurance_engine/src/InsuranceEngine.ApiHost/` | HTTP + gRPC host |
| InsuranceEngine.Policy | `backend/insurance_engine/src/InsuranceEngine.Policy/` | Policy, quotes, claims, beneficiaries, underwriting |
| InsuranceEngine.Products | `backend/insurance_engine/src/InsuranceEngine.Products/` | Product catalog, pricing |
| InsuranceEngine.SharedKernel | `backend/insurance_engine/src/InsuranceEngine.SharedKernel/` | Behaviors, interfaces, CQRS base |

**PoliSync (C#) — vertical slices:**

| Project | Path | Role |
|---------|------|------|
| PoliSync.ApiHost | `backend/polisync/src/PoliSync.ApiHost/` | gRPC host, DI, Kafka consumers |
| PoliSync.Claims | `backend/polisync/src/PoliSync.Claims/` | Claims slice |
| PoliSync.Commission | `backend/polisync/src/PoliSync.Commission/` | Commission |
| PoliSync.Endorsement | `backend/polisync/src/PoliSync.Endorsement/` | Endorsements |
| PoliSync.Orders | `backend/polisync/src/PoliSync.Orders/` | Orders |
| PoliSync.Policy | `backend/polisync/src/PoliSync.Policy/` | Policy |
| PoliSync.Products | `backend/polisync/src/PoliSync.Products/` | Products |
| PoliSync.Quotes | `backend/polisync/src/PoliSync.Quotes/` | Quotes/Quotations |
| PoliSync.Refund | `backend/polisync/src/PoliSync.Refund/` | Refunds |
| PoliSync.Renewal | `backend/polisync/src/PoliSync.Renewal/` | Renewals |
| PoliSync.Underwriting | `backend/polisync/src/PoliSync.Underwriting/` | Underwriting |
| PoliSync.Infrastructure | `backend/polisync/src/PoliSync.Infrastructure/` | DbContext, repos, InsuranceServiceClient, gRPC clients |
| PoliSync.SharedKernel | `backend/polisync/src/PoliSync.SharedKernel/` | CQRS, domain base |
| PoliSync.Proto | `backend/polisync/src/PoliSync.Proto/` | Proto references |

**Inscore (Go):**
- `backend/Inscore/microservices/insurance/` — Insurance gRPC service (full CRUD)
- `backend/Inscore/cmd/` — gateways, dbmanager, authn, etc.
- `backend/Inscore/db/migrations/` — schema-based SQL (e.g. `insurance_schema/`)

### 1.3 Architecture map (high level)

```
                    ┌─────────────────────────────────────────────────────────┐
                    │                    Clients / Gateways                     │
                    │   (Web, Mobile, Partner APIs, API Gateway Go :8080)       │
                    └─────────────────────────────┬───────────────────────────┘
                                                  │
         ┌────────────────────────────────────────┼────────────────────────────────────────┐
         │                                        │                                        │
         ▼                                        ▼                                        ▼
┌─────────────────────┐              ┌─────────────────────┐              ┌─────────────────────┐
│  PoliSync (C#)      │              │ Insurance Engine    │              │ Inscore (Go)        │
│  Insurance Commerce │              │ (C#)                │              │ microservices       │
│  & Policy Engine    │──────────────▶│ Port: (default)     │              │ insurance :50115    │
│  gRPC API           │  InsuranceServiceClient   │              │ (gRPC)   │  authn, authz, etc. │
│                     │  default URL: :50115      │              │          │                     │
│  • QuotesGrpcService│  → Go Insurance           │  REST + gRPC  │          │  Full InsuranceService│
│  • PolicyGrpcService│  • Policy, Claim, Quote   │  (4 RPCs only)│          │  CRUD (90+ RPCs)    │
│  • ClaimGrpcService │  • Endorsement, Renewal   │              │          │                     │
│  • etc.             │  • Underwriting           │  DbContext:   │          │  DB: same schema   │
│                     │                            │  insurance_  │          │  insurance_schema  │
│  PoliSyncDbContext  │                            │  schema      │          │                     │
│  (InsuranceDb)      │                            │              │          │                     │
└─────────────────────┘                            └─────────────────────┘  └─────────────────────┘
         │                                        │                                        │
         └────────────────────────────────────────┼────────────────────────────────────────┘
                                                  ▼
                                    ┌─────────────────────────┐
                                    │  PostgreSQL             │
                                    │  insuretech_primary /   │
                                    │  insurance_schema       │
                                    │  (policies, claims,     │
                                    │   products, quotes, …)  │
                                    └─────────────────────────┘
```

**Finding:** PoliSync’s `InsuranceServiceClient` is configured to **Go Insurance (Inscore)** at **localhost:50115** (`backend/Inscore/configs/services.yaml`: `insurance.ports.grpc: 50115`). The **C# Insurance Engine** is a separate process (no port in appsettings; SRS mentions port 5001) and implements only **4** of the shared `InsuranceService` gRPC methods. It is **not** the backend used by PoliSync for Policy/Claims/Quotes CRUD.

---

## 2. Documentation Sources Summary

### 2.1 SRS (Source of Truth)

| Item | Path | Content |
|------|------|--------|
| SRS V3.11 | `documentation/SRS_v3/SRS_V3.11.rtf` | Main SRS (RTF) |
| SPECS V3.7 | `documentation/SRS_v3/SPECS_V3.7/sections/*.md` | Structured sections (markdown) |

**Sections used for this audit:**
- `02_system_overview.md` — Business context, system boundaries
- `03_architecture.md` — VSA, microservices, tech stack (Insurance Engine C# @ 5001)
- `04_functional_requirements.md` — FRs by feature group (FG-001–FG-023)
- `06_data_model.md` — Proto-first strategy, DB strategy, entities
- `08_integration.md` — External systems, internal gRPC (InsuranceEngineService: IssuePolicy, CalculatePremium, ProcessRenewal, SubmitClaim)

**Relevant feature groups (SRS):**
- FG-003 Product Management (FR-021–FR-029)
- FG-004 Policy Lifecycle (FR-030–FR-040)
- FG-005 Policy Management & Renewals (FR-084–FR-102)
- FG-008 Claims (FR-041–FR-058, FR-103–FR-105)
- FG-006 Business Rules (FR-214–FR-222)

### 2.2 API documentation (HTML)

- **Location:** `docs/`, `api/docs/`
- **Content:** Per-endpoint and per-schema HTML (e.g. `endpoint_quotes_post.html`, `schema_*`, `enum_*`). Same structure in both folders (~1900+ files each).
- **Use:** API contracts, DTO/schema names, request/response shapes, workflows.

---

## 3. Database Structure Analysis

### 3.1 PostgreSQL sources

**Insurance Engine (C#):**
- `backend/insurance_engine/migration_script.sql` — Only creates `__EFMigrationsHistory`. No DDL for business tables.
- **Schema:** EF Core `PolicyDbContext` and `ProductsDbContext` use:
  - Default schema: `insurance_schema`
  - Tables: `policies`, `policy_riders`, `quotes`, `nominees`, `beneficiaries`, `individual_beneficiaries`, `business_beneficiaries`, `health_declarations`, `underwriting_decisions`, `claims`, `claim_approvals`, `claim_documents` (Policy); Products uses its own context and product tables.

**Inscore (Go) — insurance_schema:**
- **Path:** `backend/Inscore/db/migrations/insurance_schema/*.up.sql`
- **Nature:** “Enhance” migrations (indexes, triggers, comments). They assume tables like `insurance_schema.policies`, `insurance_schema.claims`, etc. already exist.
- **No base DDL in repo:** Initial `CREATE TABLE` for `insurance_schema` is not in the scanned paths; likely from another migration source or manual/out-of-repo scripts.

**Tables implied by migrations and DbContext:**

| Table (insurance_schema) | Evidence |
|--------------------------|----------|
| policies | PolicyDbContext, 20250129_001_enhance_policies.up.sql |
| policy_riders, policy_nominees | PolicyDbContext, enhance_policy_riders, enhance_policy_nominees |
| quotes | PolicyDbContext, 20250129_032_enhance_quotes.up.sql |
| claims, claim_documents, claim_approvals | PolicyDbContext, enhance_claims, enhance_claim_documents, enhance_claim_approvals |
| beneficiaries, individual_beneficiaries, business_beneficiaries | PolicyDbContext, enhance_beneficiaries, etc. |
| health_declarations, underwriting_decisions | PolicyDbContext, enhance_health_declarations, enhance_underwriting_decisions |
| products, product_riders, pricing_configs | ProductsDbContext / Inscore migrations |
| renewal_schedules, renewal_reminders, grace_periods | Inscore migrations |
| endorsements, insurers, insurer_configs, insurer_products | Inscore migrations |
| fraud_rules, fraud_cases, fraud_alerts | Inscore migrations |
| quotations (policy) | proto policy.entity.v1.Quotation; PoliSync/Go use “quotation” |

### 3.2 Documentation vs database

- **SRS (06_data_model.md):** Describes Proto → PostgreSQL mapping and tables (e.g. policies, claims, claim_documents). No conflict with the tables above.
- **Gap:** Base `CREATE TABLE` scripts for `insurance_schema` are not present in the repo; only EF model and “enhance” migrations are. Verification of exact column list and constraints would require the missing initial DDL or a live DB.

### 3.3 Insurance Engine entities vs schema

- **PolicyDbContext** (`backend/insurance_engine/src/InsuranceEngine.Policy/Infrastructure/Persistence/PolicyDbContext.cs`): Uses `insurance_schema` and table names (e.g. `policies`, `quotes`, `beneficiaries`) consistent with migration names.
- **Products:** Separate `ProductsDbContext`; product entities align with product-related migrations in Inscore.
- **Conclusion:** Entity and table naming are aligned; full column-level alignment cannot be confirmed without base schema DDL.

---

## 4. Insurance Engine Architecture (C#)

### 4.1 Components

| Layer | Location | Description |
|-------|----------|-------------|
| **ApiHost** | `InsuranceEngine.ApiHost/` | Program.cs, DI, Controllers, gRPC mapping |
| **Policy** | `InsuranceEngine.Policy/` | Policies, quotes, claims, beneficiaries, underwriting, nominees |
| **Products** | `InsuranceEngine.Products/` | Products, pricing, product questions |
| **SharedKernel** | `InsuranceEngine.SharedKernel/` | Behaviors (validation, logging, transaction), interfaces |

### 4.2 gRPC exposure

- **Contract:** `proto/insuretech/insurance/services/v1/insurance_service.proto` → generated `Insuretech.Insurance.Services.V1.InsuranceService`.
- **Implemented in C# Engine:** Only in `InsuranceEngine.ApiHost/GrpcServices/InsuranceGrpcService.cs`:
  - `GetProduct`
  - `CreatePolicy`
  - `UpdatePolicy` (partial: maps to IssuePolicy for status Active)
  - `CreateClaim`
- **Total RPCs in proto:** 90+ (full CRUD for Product, ProductPlan, Rider, PricingConfig, Policy, Claim, Quote, UnderwritingDecision, HealthDeclaration, RenewalSchedule, RenewalReminder, GracePeriod, Insurer, InsurerConfig, InsurerProduct, FraudRule, FraudCase, FraudAlert, Beneficiary, Individual/Business Beneficiary, Endorsement, Quotation, PolicyServiceRequest, ServiceProvider).
- **Conclusion:** The C# Insurance Engine implements **4 RPCs**; the rest are either unimplemented (Unimplemented) or implemented elsewhere (PoliSync + Go Insurance).

### 4.3 REST API (Controllers)

**PoliciesController** (`api/policies`):
- GET, GET {id}, POST, POST {id}/issue, POST {id}/cancel, POST {id}/renew, GET {id}/grace-period, GET {id}/renewal-schedule, GET/POST/PUT/DELETE nominees.

**ProductsController** (`api/products`):
- GET, GET {id}, GET search, POST, PUT {id}, POST {id}/activate, deactivate, discontinue, calculate-premium.

**ClaimsController** (`api/claims`):
- POST, GET {id}, GET customer/{customerId}.

**BeneficiariesController** (`api/v1/beneficiaries`):
- POST individual, POST business, GET {id}, PATCH {id}, POST {id}/kyc, risk-score, GET, GET {id}/quotes, audit-trail, documents, media, workflow-history, commission-statement.

**UnderwritingController** (`api/underwriting`):
- POST quotes, GET quotes/{id}, GET quotes, PATCH quotes/{id}/status, POST/GET health-declarations, POST/GET decisions.

---

## 5. Implemented Modules (Insurance Engine C#)

### 5.1 Policy

- **Commands:** CreatePolicy, IssuePolicy, CancelPolicy, RenewPolicy, ApplyForQuote, UpdateQuoteStatus, RecordUnderwritingDecision, SubmitClaim; Beneficiaries (Create Individual/Business, Update, CompleteKYC, UpdateRiskScore); Nominees (Add, Update, Delete); HealthDeclarations (Submit).
- **Queries:** GetPolicy, ListPolicies, GetQuote, ListQuotes, GetGracePeriod, GetRenewalSchedule, Claim (by id, by customer), GetUnderwritingHistory, Beneficiaries (Get, List, GetBeneficiaryQuotes), HealthDeclarations (Get, GetByQuote), Nominees (List).
- **Controllers:** PoliciesController, ClaimsController, BeneficiariesController, UnderwritingController.
- **Persistence:** PolicyDbContext (insurance_schema).

### 5.2 Products

- **Commands:** CreateProduct, UpdateProduct, ActivateProduct, DeactivateProduct, DiscontinueProduct, CalculatePremium.
- **Queries:** GetProduct, ListProducts, SearchProducts, GetProductCode.
- **Controller:** ProductsController.
- **Extra:** GetProductQuestions (local proto `insurance_engine.proto`) in `InsuranceEngine.Products/GrpcServices/InsuranceGrpcService.cs` (separate from main InsuranceService).

### 5.3 Shared

- MediatR, ValidationBehavior, LoggingBehavior, TransactionBehavior, Kafka event bus, health checks, Swagger in Development.

---

## 6. Partially Implemented Modules

### 6.1 gRPC (InsuranceService contract)

- **Implemented:** GetProduct, CreatePolicy, UpdatePolicy (Issue only), CreateClaim.
- **Not implemented in C# Engine:** All other InsuranceService RPCs (e.g. ListProducts, UpdateProduct, DeleteProduct; GetPolicy, UpdatePolicy full, DeletePolicy, ListPolicies; GetClaim, UpdateClaim, DeleteClaim, ListClaims; full Quote/Quotation CRUD; UnderwritingDecision, HealthDeclaration, RenewalSchedule, RenewalReminder, GracePeriod; Insurer, Fraud, Beneficiary, Endorsement, ServiceProvider; etc.). Callers using the shared contract against the C# host would get Unimplemented for these.

### 6.2 UpdatePolicy (business transitions)

- **File:** `InsuranceEngine.ApiHost/GrpcServices/InsuranceGrpcService.cs`
- Only “Issue” (status Active) is implemented; other transitions throw `Unimplemented`.

### 6.3 REST vs SRS

- REST covers policy lifecycle, quotes, claims, beneficiaries, underwriting, products. Gaps vs SRS (e.g. full renewal reminders, grace period behaviour, cancellation/refund workflow details) would need a feature-by-feature check against FRs.

---

## 7. Missing Modules (vs SRS / Proto)

### 7.1 In C# Insurance Engine (vs SRS “Insurance Engine”)

- **Renewal:** RenewalSchedule, RenewalReminder CRUD and automation (FR-085–FR-092) — logic partially in commands (e.g. RenewPolicy) but no dedicated renewal gRPC CRUD.
- **Endorsement:** No endorsement module or gRPC (FR-098–FR-102).
- **Refund:** No refund module in Engine (refund lives in PoliSync + payment).
- **Commission:** No commission in Engine (in PoliSync).
- **Insurer / InsurerConfig / InsurerProduct:** No insurer management in Engine (in proto and Go).
- **Fraud:** FraudRule, FraudCase, FraudAlert — not in Engine (in proto and Go).
- **Service provider:** Not in Engine.
- **Policy service request:** Not in Engine.
- **Full ProductPlan / Rider / PricingConfig gRPC:** Only product-level and premium calc; no Plan/Rider/Config CRUD in Engine.

### 7.2 Proto vs implementation

- **proto/insuretech/insurance/services/v1/insurance_service.proto:** Defines the full InsuranceService (90+ RPCs).
- **C# Insurance Engine:** Implements 4 RPCs.
- **Go Inscore insurance:** Implements full InsuranceService (see `backend/Inscore/microservices/insurance/service/insurance_service.go`).
- **PoliSync:** Implements the same interface (e.g. QuotesGrpcService, PolicyGrpcService) and delegates persistence to **Go Insurance** via `InsuranceServiceClient` (URL 50115).

---

## 8. API Documentation vs Implementation

### 8.1 HTML docs vs actual endpoints

- **Docs:** `docs/`, `api/docs/` contain endpoint and schema HTML (e.g. `endpoint_quotes_post.html` for POST /v1/quotes).
- **Insurance Engine REST:** Controllers use routes like `api/policies`, `api/claims`, `api/v1/beneficiaries`, `api/underwriting`, `api/products`. A direct path comparison (e.g. `/v1/quotes` vs `api/underwriting` or PoliSync routes) was not fully enumerated; recommend a script to list all `endpoint_*.html` paths and compare to registered routes in both Engine and PoliSync.

### 8.2 DTO and schema names

- REST DTOs live in `InsuranceEngine.Policy/Application/DTOs/` and `InsuranceEngine.Products/Application/DTOs/`.
- API docs reference schema names (e.g. in HTML); these are likely generated from OpenAPI or proto. No systematic DTO-name mismatch was verified; a name-by-name check against docs would be needed for a full sign-off.

### 8.3 SRS FR vs implementation

- **FR-030 (policy purchase flow):** Partially covered by CreatePolicy, ApplyForQuote, IssuePolicy, nominees, payment (payment is external).
- **FR-034 (policy number format LBT-YYYY-XXXX-NNNNNN):** Policy number generation exists (e.g. PolicyNumberGenerator); format not verified in code.
- **FR-043 (claim number CLM-YYYY-XXXX-NNNNNN):** Claim submission exists; claim number format not verified.
- **FR-218–FR-219 (claim status state machine, rules):** Claim entities and approval exist; full state machine and rules need to be checked against SRS.

---

## 9. Inter-Module Communication

### 9.1 PoliSync → “Insurance”

- **PoliSync** uses `InsuranceServiceClient` (gRPC) to call the service at `InsuranceService:Url` (default **http://localhost:50115**).
- **Port 50115** is the **Go Insurance** microservice (`backend/Inscore/configs/services.yaml`), not the C# Insurance Engine.
- So: **PoliSync talks to Go Insurance for Policy/Claim/Quote/Endorsement/Renewal/Underwriting persistence**, not to the C# Engine.

### 9.2 API contract alignment

- **Contract:** Same `proto/insuretech/insurance/services/v1/insurance_service.proto` for PoliSync, Go Insurance, and C# Engine.
- **Go:** Implements all RPCs.
- **PoliSync:** Implements the same interface and delegates to Go.
- **C# Engine:** Implements 4 RPCs; if something (e.g. gateway) pointed at the C# host instead of Go, all other RPCs would be Unimplemented.

### 9.3 DTO and proto types

- PoliSync uses `Insuretech.Policy.Entity.V1.Policy`, `Insuretech.Claims.Entity.V1.Claim`, etc. from gen/csharp.
- Insurance Engine’s gRPC layer uses the same generated types for GetProduct, CreatePolicy, UpdatePolicy, CreateClaim. So for those 4 RPCs, DTO/proto types are aligned.

### 9.4 Mismatches

- **Port/role:** SRS says Insurance Engine is C# (5001). PoliSync’s “Insurance” backend is Go (50115). So “Insurance Engine” in SRS is not the same as the service PoliSync currently uses for insurance CRUD.
- **Coverage:** C# Engine does not implement the full InsuranceService; it cannot replace Go for PoliSync without implementing the remaining RPCs (or PoliSync changing provider).

---

## 10. Documentation vs Implementation Issues

| Issue | Location / Evidence |
|-------|----------------------|
| **Insurance Engine port vs SRS** | SRS: Insurance Engine C# @ 5001. PoliSync: InsuranceServiceClient @ 50115 (Go). C# Engine port not set in appsettings. |
| **Two implementations of InsuranceService** | Go (full) and C# (4 RPCs). Not clearly documented which is “the” Insurance Engine for which consumers. |
| **Base DB schema missing in repo** | insurance_schema base CREATE TABLE not in repo; only EF model and “enhance” migrations. |
| **REST vs doc paths** | Engine uses `api/[controller]`; docs use patterns like `/v1/quotes`. Need mapping between doc paths and actual routes (Engine + PoliSync + Gateway). |
| **UpdatePolicy gRPC** | Comment in code: “Generic UpdatePolicy transitions not yet fully implemented”. |
| **Local proto unused for main host** | `InsuranceEngine.Products/Protos/insurance_engine.proto` (GetProductQuestions, IssuePolicy, CalculatePremium, etc.) is not the contract used by ApiHost; ApiHost uses gen from `proto/insuretech/insurance/...`. |

---

## 11. Implementation Coverage

### 11.1 Documentation vs code (high level)

- **SRS modules (e.g. FG-003–FG-008):** Product, Policy, Claims, Underwriting, Beneficiaries are largely implemented in the C# Engine at REST and in part at gRPC. Renewal, Endorsement, Refund, Commission, Insurer, Fraud are in PoliSync/Go, not in the C# Engine.
- **Insurance Engine (C#) completion (by gRPC contract):** 4 / 90+ ≈ **&lt;5%** of InsuranceService RPCs.
- **Insurance Engine (C#) completion (by REST features):** Policy lifecycle, quotes, claims, beneficiaries, underwriting, products are present; exact FR-level percentage would require a full FR-to-code matrix.

### 11.2 Insurance Engine completion (summary)

| Dimension | Status |
|-----------|--------|
| **REST API** | Policy, Products, Claims, Beneficiaries, Underwriting implemented. |
| **gRPC (InsuranceService)** | 4 RPCs implemented; 86+ RPCs not implemented. |
| **SRS “Insurance Engine” (C# @ 5001)** | Partially implemented; renewal/endorsement/insurer/fraud/service provider not in Engine. |
| **PoliSync’s insurance backend** | Fully served by **Go Insurance (50115)**, not by C# Engine. |

---

## 12. Recommendations

### 12.1 Next modules / features (priority)

1. **Clarify architecture:** Document that PoliSync uses **Go Insurance (50115)** for insurance CRUD and that the **C# Insurance Engine** is a separate service (e.g. for REST and future business logic). Decide whether the C# Engine should ever replace Go for gRPC and, if so, for which callers.
2. **If C# Engine should implement InsuranceService:** Add implementations for at least: ListProducts, GetPolicy, ListPolicies, GetClaim, ListClaims, CreateQuote, GetQuote, UpdateQuote, ListQuotes (and other Quote/Quotation RPCs used by PoliSync). Then either switch PoliSync’s `InsuranceService:Url` to the C# host or run both and route by capability.
3. **Renewal and endorsement:** Implement renewal schedule/reminder and endorsement logic in the C# Engine if SRS assigns them to “Insurance Engine,” and expose via gRPC if PoliSync or others will call it.
4. **Refund and commission:** Keep in PoliSync/payment as today, unless SRS explicitly assigns them to the Engine.

### 12.2 Architecture improvements

1. **Single source of truth for “Insurance”:** Either (a) make C# Engine the only InsuranceService implementation and move PoliSync to it, or (b) formally define Go as the insurance CRUD service and C# as a separate “engine” (e.g. rules, pricing, underwriting) and document the split.
2. **Port and discovery:** Set and document the C# Insurance Engine port (e.g. 5001) in appsettings and deployment config; document which services call which host (Go vs C#).
3. **Database:** Add or reference base `insurance_schema` DDL (CREATE TABLE) in the repo so schema is reproducible and auditable; keep EF and “enhance” migrations in sync with it.

### 12.3 Documentation improvements

1. **SRS / deployment:** State explicitly: “Insurance CRUD for PoliSync is provided by Go Insurance (50115); C# Insurance Engine (5001) provides REST and the following gRPC methods: …”
2. **API docs:** Align HTML endpoint paths with actual routes (Engine, PoliSync, Gateway) and note which service owns which path.
3. **Contract:** Either retire `InsuranceEngine.Products/Protos/insurance_engine.proto` or document that it is only for GetProductQuestions and not the main InsuranceService.

### 12.4 Testing and compatibility

1. **Integration tests:** PoliSync.InsuranceTest and PoliSync.DbTest call InsuranceService at 50115; add tests that run against the C# Engine when it implements more RPCs.
2. **Contract tests:** Add contract tests (e.g. for proto) so Go and C# implementations stay aligned with the same `.proto` and with docs.

---

## Appendix A: File References

| Topic | Path |
|-------|------|
| SRS system overview | `documentation/SRS_v3/SPECS_V3.7/sections/02_system_overview.md` |
| SRS architecture | `documentation/SRS_v3/SPECS_V3.7/sections/03_architecture.md` |
| SRS functional requirements | `documentation/SRS_v3/SPECS_V3.7/sections/04_functional_requirements.md` |
| SRS data model | `documentation/SRS_v3/SPECS_V3.7/sections/06_data_model.md` |
| SRS integration | `documentation/SRS_v3/SPECS_V3.7/sections/08_integration.md` |
| InsuranceService proto | `proto/insuretech/insurance/services/v1/insurance_service.proto` |
| Go Insurance service | `backend/Inscore/microservices/insurance/service/insurance_service.go` |
| Go Insurance main | `backend/Inscore/microservices/insurance/cmd/server/main.go` |
| C# Engine gRPC | `backend/insurance_engine/src/InsuranceEngine.ApiHost/GrpcServices/InsuranceGrpcService.cs` |
| PoliSync Insurance client | `backend/polisync/src/PoliSync.Infrastructure/Clients/InsuranceServiceClient.cs` |
| PoliSync Program (DI, URL) | `backend/polisync/src/PoliSync.ApiHost/Program.cs` |
| Services config (ports) | `backend/Inscore/configs/services.yaml` |
| PolicyDbContext | `backend/insurance_engine/src/InsuranceEngine.Policy/Infrastructure/Persistence/PolicyDbContext.cs` |
| Insurance Engine migration | `backend/insurance_engine/migration_script.sql` |
| Inscore insurance migrations | `backend/Inscore/db/migrations/insurance_schema/*.up.sql` |

---

*End of Technical Audit Report.*
