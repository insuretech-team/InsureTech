# InsureTech Platform

A production-grade **insurance commerce platform** with a dual-backend architecture, proto-first API design, and multi-portal frontend ecosystem.

---

## Architecture Overview

InsureTech uses a **dual-backend** design where two complementary engines share 39 proto-defined domains:

### InScore (Go Microservices)

Platform infrastructure services in `backend/inscore/` — 28 `cmd/` packages:

| Service | gRPC | HTTP | Domain |
|---------|------|------|--------|
| authn | 50060 | 50061 | Authentication (JWT, OTP, device binding) |
| authz | 50070 | 50071 | Authorization (Casbin RBAC) |
| audit | 50080 | 50081 | Audit logging |
| kyc | 50090 | 50091 | eKYC / face liveness verification |
| partner | 50100 | 50101 | Distribution partner management |
| insurance | 50115 | — | Core insurance service |
| sync_order | 50142 | 50143 | Order persistence & Kafka events |
| payment | 50190 | 50191 | Payment processing & reconciliation |
| fraud | 50220 | 50221 | ML-based fraud detection |
| notification | 50230 | 50231 | SMS, email, push notifications |
| docgen | 50280 | 50281 | Document generation (policies, certs) |
| storage | 50290 | 50291 | File / object storage |

Additional services: `b2b`, `beneficiary`, `media`, `tenant`, `workflow`, `support`, `ai`, `analytics`, `iot`, `webrtc`, `gateway`, database CLIs (`dbx`, `dbsql`, `dbops`).

### PoliSync (C# .NET 8 Insurance Engine)

Insurance business logic in `backend/polisync/` — 15 `src/` projects:

| Service | gRPC | HTTP | Domain |
|---------|------|------|--------|
| Product | 50120 | 50121 | Product catalog & configuration |
| Quote | 50130 | 50131 | Premium calculation & quotation |
| Order | 50140 | 50141 | Purchase order business logic (CQRS) |
| Commission | 50150 | 50151 | Agent/partner commission rules |
| Policy | 50160 | 50161 | Policy lifecycle & issuance |
| Underwriting | 50170 | 50171 | Risk assessment & approval |
| Claims | 50210 | 50211 | Claims submission & adjudication |

Additional projects: `Endorsement`, `Renewal`, `Refund`, `Workflow`, `SharedKernel`, `Infrastructure`, `Proto`, `ApiHost`.

> **⚠️ Critical — Order Domain Split:**  
> `order` (PoliSync, port 50140) = business logic, validations, CQRS commands  
> `sync_order` (InScore, port 50142) = persistence layer, Kafka events, database  
> They are **two halves of one domain** — changes often require both backends.

### Support Services

| Service | Port |
|---------|------|
| PostgreSQL | 5432 |
| Redis | 6379 |
| Kafka | 9092 |
| Zookeeper | 2181 |
| Gateway / API | 8080 |
| Nginx | 80 / 443 |

---

## Proto-First API Design

All schema definitions live in **Protocol Buffers** — 39 domain packages in `proto/insuretech/`:

```
Core:           insurance, policy, claims, billing, payment, products, orders, renewal, refund, endorsement
User:           authn, authz, kyc, beneficiary, partner, b2b, tenant, insurer
Operational:    commission, underwriting, fraud, workflow, task, audit, analytics, report
Infrastructure: common, services, notification, document, media, storage, apikey, support
Advanced:       ai, iot, voice, webrtc, mfs
```

### Code Generation Pipeline

```
proto/insuretech/        →  buf generate  →  gen/go/      (Go + gRPC)
                                          →  gen/ts/      (TypeScript)
                                          →  gen/csharp/  (C# + gRPC)
                             ↓
                    inject_gorm_tags.go   →  GORM struct tags on Go types
                    organize-csharp       →  C# namespace organization
                    TS sync               →  b2b_portal/src/lib/proto-generated/
```

Run: `scripts/generate.ps1` (Windows) or `scripts/generate.sh` (Linux)

### API Generation Pipeline

Proto definitions feed into a **58-file Python generator** (`api/generator/`) that produces:

- `api/openapi.yaml` — 3.3 MB OpenAPI 3.0 specification
- `api/schemas/` — 672 schema files
- `api/events/` — 170 event schema files
- `api/enums/` — 134 enum files
- `api/paths/` — 30 path files
- `api/ENDPOINT_MAP.md` — All REST routes documented
- `api/docs/` — HTML documentation + Schema Visualizer

Run: `run_api_pipeline.ps1` (Windows) or `run_api_pipeline.sh` (Linux)

---

## Frontend Portals

Five portals consume the backend through a **BFF (Backend-for-Frontend) pattern** — browsers never call the gateway directly.

| Portal | Framework | Purpose |
|--------|-----------|---------|
| `customer_portal` | SvelteKit | End-user policy management |
| `b2b_portal` | Next.js 13+ App Router | Corporate group insurance |
| `insurer-portal` | SvelteKit | Insurance company admin |
| `partner_portal` | SvelteKit | Distribution agent/broker |
| `system_portal` | SvelteKit | Platform super-admin |

**Shared resources:**
- `web_shared/` — Cross-portal components
- `gen/ts/` — Generated TypeScript proto types (read-only)
- Auth: Cookie-based sessions with portal headers (`x-portal`, `x-user-id`, `x-business-id`, `x-tenant-id`)

---

## REST API Standards

All endpoints comply with **9 normative rules** documented in `rules/`:

| # | Rule | Summary |
|---|------|---------|
| 01 | Response Envelope | `{ success, data, error, meta }` on every response |
| 02 | HTTP Status Codes | 201 create, 200 read, 204 delete, 400–500 errors |
| 03 | Error Handling | Structured `{ code, message, error_id, retryable, field_violations }` |
| 04 | Security | Every endpoint declares auth explicitly |
| 05 | Pagination | `{ page, page_size, total_count, has_next, has_prev }` |
| 06 | DI & Testing | Client-consumable contracts, one `ApiResponse<T>` decoder |
| 07 | URL Naming | Resource-based, kebab-case paths |
| 08 | Null / Optional | Required vs nullable explicitly declared |
| 09 | Generator Compliance | OpenAPI generators emit rule-compliant output |

---

## Database

### Proto-First Schema Model

**NEVER write `CREATE TABLE`, `ALTER TABLE`, or `ADD COLUMN` in SQL migrations.**

```
Proto → buf generate → GORM tags injected → Engine auto-derives schema → Tables auto-created/altered
```

SQL migrations are for **data operations only** (seeds, backfills, transforms).

### Naming Conventions

| Element | Pattern | Example |
|---------|---------|---------|
| Primary Key | `{entity}_id` | `policy_id` (never generic `id`) |
| Foreign Key | `{referenced}_id` | `customer_id` |
| Table | snake_case, plural | `policy_riders` |
| Money | field + `_currency` | `premium_amount` + `premium_currency` |
| Audit | JSONB column | `audit_info` |

### Database CLI Tools

```powershell
cd backend/inscore
go run ./cmd/dbx <operation>     # Specialized database operations
go run ./cmd/dbsql <query>       # Lightweight SQL queries
go run ./cmd/dbops <command>     # Bulk operations management
```

---

## SDKs

Generated client SDKs from OpenAPI + proto definitions:

| SDK | Status | Location |
|-----|--------|----------|
| Go | ✅ Available | `sdks/insuretech-go-sdk/` |
| TypeScript | ✅ Available | `sdks/insuretech-typescript-sdk/` |
| Python | 🚧 Planned | — |
| Java | 🚧 Planned | — |

SDK generators: `sdks/sdk-generator/go/` and `sdks/sdk-generator/typescript/`

---

## Quick Start

### Prerequisites

- **Go 1.22+** — InScore backend
- **.NET 8 SDK** — PoliSync engine
- **Node.js 20+** — Frontend portals
- **Docker & Docker Compose** — Infrastructure services
- **Buf CLI** — Proto code generation
- **Python 3.8+** — API generator

### Bootstrap

```powershell
# Full environment setup (installs tools, generates code, starts services)
.\scripts\bootstrap.ps1        # Windows
./scripts/bootstrap.sh         # Linux
```

### Code Generation

```powershell
.\scripts\generate.ps1         # Proto → Go/TS/C# code generation
.\run_api_pipeline.ps1         # Proto → OpenAPI → SDK → docs
```

### Run Services

```powershell
# Start all infrastructure
docker-compose up -d

# Check service health
.\scripts\check_services.ps1

# Run database migrations (data-only)
.\run_migration.ps1

# Start mock server for frontend dev
.\scripts\start_mock_server.ps1
```

### Development Reset

```powershell
.\safe_reset_db.ps1            # Reset dev database (DESTRUCTIVE)
```

---

## Project Structure

```
InsureTech/
├── proto/insuretech/           # 39 proto domain packages (source of truth)
├── gen/                        # Generated code (go/, ts/, csharp/)
├── backend/
│   ├── inscore/                # Go microservices (28 cmd/ packages)
│   └── polisync/               # C# .NET 8 insurance engine (15 src/ projects)
├── api/
│   ├── openapi.yaml            # Master OpenAPI 3.0 spec (3.3 MB)
│   ├── generator/              # 58-file Python proto→OpenAPI generator
│   ├── schemas/                # 672 generated schema files
│   ├── ENDPOINT_MAP.md         # All REST routes documented
│   └── postman/                # Postman collections
├── customer_portal/            # SvelteKit — end-user portal
├── b2b_portal/                 # Next.js 13+ — B2B corporate portal
├── insurer-portal/             # SvelteKit — insurer admin
├── partner_portal/             # SvelteKit — distribution partner
├── system_portal/              # SvelteKit — platform super-admin
├── web_shared/                 # Shared frontend code
├── sdks/                       # SDK generators + generated SDKs
├── rules/                      # 9 normative API rules + DB conventions
├── scripts/                    # 27 utility scripts (bootstrap, generate, deploy)
├── documentation/              # Architecture docs, BRD, SRS, core plans
├── specs/                      # Feature specifications
├── seek/                       # SQLite codebase search index CLI
├── .resources/                 # 14 reference insurance projects
├── .opencode/                  # OpenCode AI agent config (skills, agents, prompts)
├── buf.yaml                    # Buf lint + breaking change config
├── buf.gen.yaml                # Buf code generation plugins
├── docker-compose.yml          # Development services (21 KB)
├── docker-compose-prod.yml     # Production services (14 KB)
└── opencode.json               # Kimi K2.5 agent configuration
```

---

## Documentation

| Document | Location |
|----------|----------|
| **Start Here** | `documentation/About/START_HERE.md` |
| **Architecture Overview** | `documentation/About/ARCHITECTURE_OVERVIEW.md` |
| **PoliSync Reference** | `documentation/About/POLISYNC_REFERENCE.md` |
| **Active Workstreams** | `documentation/core_plans/ACTIVE_WORKSTREAMS.md` |
| **SRS v3.7** | `documentation/SRS_v3/SPECS_V3.7/sections/` |
| **Business Requirements** | `documentation/BRD/BRDV3.7.md` |
| **API Rules** | `rules/00-index.md` |
| **Database Rules** | `rules/dbrules.md` |
| **SDK Documentation** | `sdks/README.md` |
| **API Generator** | `api/generator/README.md` |

---

## Environment Configuration

| File | Purpose |
|------|---------|
| `.env` | Default settings (comprehensive) |
| `.env.dev` | Development overrides |
| `.env.prod` | Production settings |
| `.env.example` | Template for new setups |

**Naming:** Service-specific (`POLISYNC_DB_HOST`, `INSCORE_PORT`) • Shared (`INSURETECH_` prefix) • Secrets never committed.

---

## Git Conventions

- **Branch:** `feature/<domain>/<description>`, `fix/<domain>/<description>`
- **Commit:** `<type>(<scope>): <description>` (e.g., `feat(policy): add renewal flow`)
- **Pre-push:** Always run `buf lint` and `buf breaking` before pushing proto changes

---

## Insurance Product Lifecycle

```
 1. Product Definition     →  PoliSync: product catalog, coverage rules, pricing
 2. Quote Generation       →  PoliSync: premium calculation based on risk factors
 3. KYC / eKYC             →  InScore: identity verification, face liveness
 4. Underwriting           →  PoliSync: risk assessment, approval/rejection
 5. Order Placement        →  PoliSync (order) + InScore (sync_order)
 6. Payment Processing     →  InScore: payment gateway, reconciliation
 7. Policy Issuance        →  PoliSync: policy document, certificate generation
 8. Endorsements           →  PoliSync: policy modifications mid-term
 9. Renewals               →  PoliSync: renewal processing, premium recalculation
10. Claims Processing      →  PoliSync: claim submission, adjudication, payout
11. Fraud Detection        →  InScore: ML-based fraud scoring
12. Commissions            →  PoliSync: agent/partner commission calculation
13. Refunds                →  PoliSync: cancellation and refund processing
```

---

## Multi-Tenant & Multi-Portal

- Each insurer is a **tenant** with isolated data
- `x-tenant-id` header required on all API calls
- Products, policies, and claims are tenant-scoped
- Five portals serve different user personas with role-based access (Casbin RBAC)
- Authorization domains: `system:root` (superadmin), `org:{business_id}` (B2B org scope)

---

## License

Proprietary — All Rights Reserved.
