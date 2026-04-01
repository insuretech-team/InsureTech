# B2B Backend Service — Central Reference

> **Single source of truth** for the B2B microservice backend: architecture, request flows, database schema, authorization middleware, gRPC endpoints, gateway wiring, role mapping, events, deployment, and test coverage.
>
> **Test Status:** ✅ 37 unit tests PASS · ✅ 33 live DB tests PASS · Total 70 tests
> **Auth Status:** ✅ Fixed — all org-management methods go through Casbin

---

## 1. Critical Bug Fixes

### Nil Publisher Panic in B2B Service (FIXED)

**File:** `backend/inscore/microservices/b2b/cmd/server/main.go`

**Root Cause:**
When Kafka was unavailable at startup, the event publisher was initialized as a typed nil pointer:

```go
var publisher *events.Publisher  // typed nil
if kafkaProducer != nil {
    publisher = events.NewPublisher(kafkaProducer)
}
// publisher stays as (*events.Publisher)(nil)
```

This typed nil was then passed to `NewB2BService()` as the `EventPublisher` interface. In Go, a typed nil pointer assigned to an interface creates a **non-nil interface value** (the interface holds a non-nil type descriptor but nil data pointer). So `if s.publisher != nil` passed, but calling any method on it caused:

```
panic: runtime error: invalid memory address or nil pointer dereference
    events.(*Publisher).PublishOrgMemberAdded(0x0, ...)
```

This panic was triggered every time `AddOrgMember` was called — which happens during B2B admin creation. This is why B2B admin creation was completely broken.

**The Fix:**

```go
// Always create a Publisher — even when Kafka is unavailable.
// events.NewPublisher accepts nil producer and the internal publish()
// method already no-ops gracefully when producer == nil.
publisher := events.NewPublisher(kafkaProducer)
```

The publisher's internal `publish()` method already has:

```go
func (p *Publisher) publish(topic string, key string, value []byte) error {
    if p.producer == nil {
        return nil  // no-op when Kafka unavailable
    }
    // ... actual kafka send
}
```

So passing nil `kafkaProducer` to `NewPublisher` is safe and intentional.

**Impact:**
- ✅ B2B admin creation (`AddOrgMember` + `PublishOrgMemberAdded`) now works without Kafka
- ✅ Member removal (`RemoveOrgMember` + `PublishOrgMemberRemoved`) now works without Kafka
- ✅ All org member events work correctly with or without Kafka
- ✅ Service degrades gracefully: events are skipped when Kafka unavailable, core functionality continues

---

## 2. Architecture

### System Overview

```
┌──────────────────────────────────────────────────────────┐
│                    REST GATEWAY (HTTP :8080)               │
│  1. Validate session cookie → user_id                      │
│  2. Call B2B.ResolveMyOrganisation(user_id) gRPC           │
│  3. Inject x-business-id metadata with org_id              │
│  4. Route request to B2B service                           │
└─────────────────────────┬──────────────────────────────────┘
                          │ gRPC (Port 50112)
                          ▼
┌──────────────────────────────────────────────────────────┐
│                B2B gRPC SERVER (:50112)                    │
│  ┌─────────────────────────────────────────────────────┐  │
│  │ AuthZ Interceptor → CheckAccess before every handler│  │
│  ├─────────────────────────────────────────────────────┤  │
│  │ B2BHandler (gRPC) — 21 RPC methods                  │  │
│  │   → Error mapping (InvalidArg → NotFound → Internal)│  │
│  ├─────────────────────────────────────────────────────┤  │
│  │ B2BService (business logic / validation / events)   │  │
│  ├─────────────────────────────────────────────────────┤  │
│  │ PortalRepository (GORM Raw SQL)                     │  │
│  │   org_repository · department_repository            │  │
│  │   employee_repository · purchase_order_repository   │  │
│  │   catalog_repository                                │  │
│  └────────────────────────┬────────────────────────────┘  │
└───────────────────────────┼──────────────────────────────┘
                            │ SQL
                            ▼
              ┌──────────────────────────────┐
              │    PostgreSQL (b2b_schema)    │
              └──────────────────────────────┘
```

### Folder Structure

```
microservices/b2b/
├── cmd/server/main.go                    # Entry point: init service with AuthZ middleware
├── internal/
│   ├── config/config.go                  # Service configuration
│   ├── consumers/handlers.go             # Kafka event handlers
│   ├── domain/interfaces.go              # Service interfaces (B2BService, B2BRepository)
│   ├── events/publisher.go               # Event publishing (Kafka)
│   ├── grpc/
│   │   ├── b2b_handler.go                # gRPC handler (all RPC methods)
│   │   └── server.go                     # gRPC server setup + service registration
│   ├── middleware/
│   │   ├── authz_interceptor.go          # Authorization interceptor (Casbin)
│   │   └── authz_interceptor_test.go
│   ├── repository/
│   │   ├── org_repository.go             # Organisation CRUD (6 methods)
│   │   ├── department_repository.go      # Department CRUD (5 methods)
│   │   ├── employee_repository.go        # Employee CRUD (5 methods + batch enrichment)
│   │   ├── catalog_repository.go         # Plan catalog (2 methods)
│   │   └── purchase_order_repository.go  # PO CRUD (3 methods)
│   └── service/
│       ├── b2b_service.go                # Business logic (~1160 lines)
│       └── errors.go                     # ErrInvalidArgument, ErrNotFound
```

**Key characteristics:** Clean Architecture (handler → service → repo) · gRPC-first (Protocol Buffers) · Event-driven (Kafka) · Middleware-based AuthZ (Casbin) · Soft-delete across all entities

---

## 3. Authorization Middleware

### AuthZ Interceptor Flow (`authz_interceptor.go`)

```
1. Extract metadata: x-user-id, x-portal, x-business-id, x-tenant-id
2. Is method ResolveMyOrganisation?
   YES → pass through (bootstrap call)
   NO  → continue
3. Is portal == PORTAL_SYSTEM (super_admin)?
   YES → skip "missing org" check, domain = "system:root"
   NO  → if x-business-id empty → PermissionDenied
          if x-business-id present → domain = "b2b:{org_id}"
4. Map method → resource + action (see table below)
5. Call AuthZ.CheckAccess(user_id, domain, resource, action)
   ALLOW → call handler
   DENY  → PermissionDenied
```

### Request Headers

| Header | Purpose | Required? |
|--------|---------|-----------| 
| `x-user-id` | User identity | **Yes** |
| `x-portal` | Portal type (`PORTAL_SYSTEM` or `PORTAL_B2B`) | For authorization |
| `x-business-id` | Organisation ID | Yes (except bootstrap) |
| `x-tenant-id` | Tenant ID | Fallback for org |

### Method-to-Action Mapping

| Method Pattern | Action | Resource |
|---|---|---|
| `Get*`, `List*` | `GET` | `svc:b2b/*` |
| `Create*`, `Add*`, `Assign*` | `POST` | `svc:b2b/*` |
| `Update*` | `PATCH` | `svc:b2b/*` |
| `Delete*`, `Remove*` | `DELETE` | `svc:b2b/*` |
| `ResolveMyOrganisation` | _(none)_ | _(none)_ — skips Casbin |

### Domain Resolution

| User Type | Portal Header | Casbin Domain |
|-----------|---------------|---------------|
| Super Admin | `PORTAL_SYSTEM` | `system:root` |
| B2B Org Admin | `PORTAL_B2B` + org_id | `b2b:{org_id}` |
| HR Manager | `PORTAL_B2B` + org_id | `b2b:{org_id}` |
| Viewer | `PORTAL_B2B` + org_id | `b2b:{org_id}` |

---

## 4. gRPC Endpoints (21 Methods)

### Organisation CRUD (5)

| RPC | Action | Description |
|-----|--------|-------------|
| `CreateOrganisation` | POST | Create org with name, code, industry, contact info |
| `GetOrganisation` | GET | Get single org by ID |
| `ListOrganisations` | GET | List orgs (paginated) |
| `UpdateOrganisation` | PATCH | Update org fields |
| `DeleteOrganisation` | DELETE | Soft-delete org |

### Organisation Members (4)

| RPC | Action | Description |
|-----|--------|-------------|
| `ListOrgMembers` | GET | List org members |
| `AddOrgMember` | POST | Add user to org with role |
| `AssignOrgAdmin` | POST | Promote member to admin |
| `RemoveOrgMember` | DELETE | Remove member from org |

### Department CRUD (5)

| RPC | Action | Description |
|-----|--------|-------------|
| `ListDepartments` | GET | List departments for org (paginated, sorted by name) |
| `GetDepartment` | GET | Get single department |
| `CreateDepartment` | POST | Create department (generates UUID) |
| `UpdateDepartment` | PATCH | Partial update (dynamic SET clause) |
| `DeleteDepartment` | DELETE | Soft-delete (safety: refuses if active employees exist) |

### Employee CRUD (5)

| RPC | Action | Description |
|-----|--------|-------------|
| `ListEmployees` | GET | List employees for org (paginated, enriched with dept/plan names) |
| `GetEmployee` | GET | Get single employee |
| `CreateEmployee` | POST | Create employee (UUID generated, assigns to dept) |
| `UpdateEmployee` | PATCH | Partial update |
| `DeleteEmployee` | DELETE | Soft-delete |

### Purchase Order (3 + Catalog)

| RPC | Action | Description |
|-----|--------|-------------|
| `ListPurchaseOrderCatalog` | GET | Merged catalog (DB plans + seeded fallback plans) |
| `ListPurchaseOrders` | GET | List POs for org |
| `GetPurchaseOrder` | GET | Get single PO |
| `CreatePurchaseOrder` | POST | Create PO (auto-generates `PO-YYYYMMDD-XXXXXXXX` number, calculates premium) |

### Bootstrap (1)

| RPC | Action | Description |
|-----|--------|-------------|
| `ResolveMyOrganisation` | _(skip)_ | Discover user's org (no auth check) — returns org_id, role, name |

---

## 5. Database Schema

```sql
-- ORGANISATIONS
CREATE TABLE b2b_schema.organisations (
  organisation_id UUID PRIMARY KEY,
  tenant_id UUID NOT NULL,
  name VARCHAR NOT NULL,
  code VARCHAR UNIQUE NOT NULL,
  industry VARCHAR,
  contact_email VARCHAR,
  contact_phone VARCHAR,
  address VARCHAR,
  status VARCHAR DEFAULT 'ORGANISATION_STATUS_ACTIVE',
  total_employees INT DEFAULT 0,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- ORG MEMBERS
CREATE TABLE b2b_schema.org_members (
  member_id UUID PRIMARY KEY,
  organisation_id UUID NOT NULL REFERENCES organisations,
  user_id UUID NOT NULL,
  role VARCHAR DEFAULT 'ORG_MEMBER_ROLE_HR_MANAGER',
  status VARCHAR DEFAULT 'ORG_MEMBER_STATUS_ACTIVE',
  joined_at TIMESTAMP DEFAULT NOW(),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(organisation_id, user_id)
);

-- DEPARTMENTS
CREATE TABLE b2b_schema.departments (
  department_id UUID PRIMARY KEY,
  name VARCHAR NOT NULL,
  business_id UUID NOT NULL REFERENCES organisations,
  employee_no INT DEFAULT 0,
  total_premium JSONB,     -- {"amount": 0, "currency": "BDT", "decimal_amount": 0}
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP     -- Soft delete marker
);

-- EMPLOYEES
CREATE TABLE b2b_schema.employees (
  employee_uuid UUID PRIMARY KEY,
  name VARCHAR NOT NULL,
  employee_id VARCHAR NOT NULL,
  department_id UUID NOT NULL REFERENCES departments,
  business_id UUID NOT NULL REFERENCES organisations,
  insurance_category VARCHAR,
  assigned_plan_id UUID,
  coverage_amount JSONB,
  premium_amount JSONB,
  status VARCHAR DEFAULT 'EMPLOYEE_STATUS_ACTIVE',
  number_of_dependent INT DEFAULT 0,
  email VARCHAR, mobile_number VARCHAR,
  date_of_birth DATE, date_of_joining DATE,
  gender VARCHAR, user_id UUID,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP
);

-- PURCHASE ORDERS
CREATE TABLE b2b_schema.purchase_orders (
  purchase_order_id UUID PRIMARY KEY,
  purchase_order_number VARCHAR NOT NULL UNIQUE,
  business_id UUID NOT NULL REFERENCES organisations,
  department_id UUID NOT NULL REFERENCES departments,
  product_id UUID, plan_id UUID,
  insurance_category VARCHAR,
  employee_count INT NOT NULL,
  number_of_dependents INT DEFAULT 0,
  coverage_amount JSONB, estimated_premium JSONB,
  status VARCHAR DEFAULT 'PURCHASE_ORDER_STATUS_SUBMITTED',
  requested_by VARCHAR, notes VARCHAR,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP
);

-- INSURANCE CATALOG (cross-schema)
-- insurance_schema.products (product_id, product_name, category, status)
-- insurance_schema.product_plans (plan_id, product_id, plan_name, premium_amount JSONB, status)
```

### Data Patterns

| Pattern | Implementation |
|---------|---------------|
| **Soft delete** | `deleted_at TIMESTAMP` — all queries filter `WHERE deleted_at IS NULL` |
| **Money fields** | JSONB `{"amount": 50000, "currency": "BDT", "decimal_amount": 500}` |
| **Enums** | VARCHAR strings: `EMPLOYEE_STATUS_ACTIVE`, `ORGANISATION_STATUS_PENDING`, etc. |
| **UUIDs** | All primary keys are `uuid.NewString()` |

---

## 6. Request Flow Example: Create Employee

```
1. REST GATEWAY
   • User authenticated, session → user_id
   • Call B2B.ResolveMyOrganisation(user_id) → org_id, role
   • Inject x-business-id=org_id into gRPC metadata

2. AuthZ INTERCEPTOR
   • Map "CreateEmployee" → resource=svc:b2b/*, action=POST
   • CheckAccess(user_id, b2b:{org_id}, svc:b2b/*, POST) → ALLOW

3. SERVICE LAYER
   • Validate: name, employee_id, department_id, business_id required
   • Generate UUID → repo.CreateEmployee(input)
   • Enrich: batch-fetch dept names + catalog plan names
   • Build EmployeeView with department_name & assigned_plan_name

4. REPOSITORY
   • INSERT INTO b2b_schema.employees (...)
   • Call GetEmployee(uuid) to return fresh data

5. RESPONSE
   → CreateEmployeeResponse { employee: EmployeeView, message }
```

### Pagination (List Endpoints)

```
page_token="" → offset=0 → SELECT ... LIMIT 10 OFFSET 0  → next_page_token="10"
page_token="10" → offset=10 → SELECT ... LIMIT 10 OFFSET 10 → next_page_token="20"
page_token="40" → offset=40 → 5 items returned → next_page_token="" (no more)
```

### Department Delete Safety

```
1. Check: SELECT COUNT(*) FROM employees
     WHERE department_id=$1 AND status='EMPLOYEE_STATUS_ACTIVE' AND deleted_at IS NULL
2. If count > 0 → Error "department has N active employee(s): reassign before deleting"
3. If count = 0 → UPDATE departments SET deleted_at=NOW() WHERE department_id=$1
```

---

## 7. Role Mapping & Event-Driven Assignment

### Org Member → Casbin Role

| Org Member Role | Casbin Role | Trigger |
|-----------------|-------------|---------|
| `ORG_MEMBER_ROLE_BUSINESS_ADMIN` | `b2b_org_admin` | `AssignOrgAdmin` / `AddOrgMember(BUSINESS_ADMIN)` |
| `ORG_MEMBER_ROLE_HR_MANAGER` | `partner_user` | `AddOrgMember(HR_MANAGER)` |
| `ORG_MEMBER_ROLE_VIEWER` | `partner_user` | `AddOrgMember(VIEWER)` |

### AssignOrgAdmin Flow

```
SuperAdmin calls AssignOrgAdmin(org_id, user_id)
  → b2b_service publishes B2BAdminAssigned event (x-user-id = real caller)
  → consumer HandleB2BAdminAssigned fires
  → lookupRoleID("b2b_org_admin") → authz.AssignRole(user_id, role_id, b2b:{org_id})
  → ensureScopedRolePolicies copies b2b:root policies → b2b:{org_id}
  → user now has b2b_org_admin role ✅
```

### AddOrgMember Flow (HR_MANAGER/VIEWER)

```
→ consumer HandleOrgMemberAdded → assignPartnerUserRole
→ authz.AssignRole(user_id, partner_user_role_id, b2b:{org_id})
```

---

## 8. Event Publishing (Kafka)

### Kafka Configuration

**Broker:** `localhost:9092`  
**Client ID:** `b2b-service`  
**Status:** Connected when available; service degrades gracefully when unavailable

### Published by B2B

| Event | Topic | Key Data |
|-------|-------|----------|
| `OrganisationCreated` | `b2b.organisation.created` | org_id, tenant_id, name, code, caller_id |
| `OrganisationUpdated` | `b2b.organisation.updated` | org_id, name, status, caller_id |
| `OrganisationApproved` | `b2b.organisation.approved` | org_id, approved_by |
| `OrgMemberAdded` | `b2b.org_member.added` | member_id, org_id, user_id, role, added_by |
| `OrgMemberRemoved` | `b2b.org_member.removed` | member_id, org_id, user_id, removed_by |
| `B2BAdminAssigned` | `b2b.admin.assigned` | org_id, user_id, assigned_by |

**Published Topics:**
- `b2b.org_member.added` — triggered by `AddOrgMember` RPC
- `b2b.org_member.removed` — triggered by `RemoveOrgMember` RPC
- `b2b.organisation.created` — triggered by `CreateOrganisation` RPC
- `b2b.admin.assigned` — triggered by `AssignOrgAdmin` RPC
- `b2b.organisation.approved` — triggered by org approval workflow

### Consumed by B2B

| Topic | Handler | Action |
|-------|---------|--------|
| `b2b.organisation.created` | `HandleOrganisationCreated` | Seed Casbin policies for new domain (creates p-rules) |
| `b2b.org_member.added` | `HandleOrgMemberAdded` | Assign `partner_user` role for HR_MANAGER/VIEWER (creates g-rule) |
| `b2b.admin.assigned` | `HandleB2BAdminAssigned` | Assign `b2b_org_admin` role (creates g-rule) |
| `b2b.organisation.approved` | `HandleOrganisationApproved` | Update org status and sync with AuthZ |
| `authn.user.registered` | `HandleUserRegistered` | (from AuthN) Sync user data |
| `authz.role.assigned` | `HandleRoleAssigned` | (from AuthZ) Acknowledgement |

### Graceful Degradation

When Kafka broker is unavailable:
- All event publish operations silently no-op (return nil)
- Core business operations continue normally (CRUD on DB)
- No panics; no service downtime
- When Kafka reconnects, new events are published normally

---

## 9. Business Logic Details

### PO Number Generation

```
Format: PO-YYYYMMDD-XXXXXXXX
Example: PO-20240115-A1B2C3D4
Algorithm: suffix = uppercase(uuid[:8])
```

### Premium Calculation

```
estimated_premium = plan.premium_amount × employee_count
Example: Seba plan (50,000 BDT) × 5 employees = 250,000 BDT
```

### Catalog Merging (Seeded Fallback)

```
1. Query DB: SELECT ... FROM insurance_schema.product_plans JOIN products
2. Merge with seeded plans (hardcoded fallback):
   - Seba (Health, 50k BDT), Surokkha (Health, 43k BDT), Verosa (Life, 85k BDT)
3. DB plans take priority; seeded plans fill gaps → catalog never empty
```

### Enrichment Batch Queries

```
Employee list returns include:
  • department_name: batch query SELECT department_id, name WHERE dept_id = ANY($1)
  • assigned_plan_name: catalog lookup by plan_id
```

---

## 10. Error Handling

```
Service Layer:
  ErrInvalidArgument  →  codes.InvalidArgument (400)
  ErrNotFound         →  codes.NotFound (404)
  Other errors        →  codes.Internal (500)

Middleware Layer:
  Missing metadata    →  codes.Unauthenticated (401)
  Missing user_id     →  codes.Unauthenticated (401)
  Missing org context →  codes.PermissionDenied (403)
  Casbin denies       →  codes.PermissionDenied (403)
  AuthZ service error →  codes.Internal (500)
```

---

## 11. Auth Flow Fixes Applied

| Bug | Root Cause | Fix |
|-----|-----------|-----|
| Super admin got `PermissionDenied` on org management | `requiresOrgContext()` bypassed Casbin for 10 methods | Only `ResolveMyOrganisation` bypasses auth |
| B2B admin got `PermissionDenied` on portal | `b2b_org_admin` policies had no wildcard `*` action | Added `svc:b2b/* → *` alongside verb-specific policies |
| Events published with hardcoded `"superadmin"` caller | Hardcoded strings in event publishers | `resolveCallerID(ctx)` reads real `x-user-id` from gRPC metadata |
| HR_MANAGER and VIEWER got no Casbin role | `HandleOrgMemberAdded` only handled BUSINESS_ADMIN | Added `assignPartnerUserRole()` for HR_MANAGER and VIEWER |
| `partner_user` fallback failed | No wildcard action in `b2b:root` policies | Added `svc:b2b/* → *` wildcard to `partner_user` in seeder |

---

## 12. Deployment & Configuration

### Environment Variables

```bash
# B2B Service
GRPC_PORT=50112
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=insuretech
KAFKA_BROKERS=localhost:9092
AUTHZ_SERVICE_URL=localhost:50060

# Gateway
AUTHN_SERVICE_URL=localhost:50051
AUTHZ_SERVICE_URL=localhost:50060
B2B_SERVICE_URL=localhost:50112
```

### Startup Sequence (`main.go`)

```
1. Initialize Logger → 2. Load Environment → 3. Parse services.yaml
→ 4. Validate B2BConfig → 5. Initialize Database → 6. Create Repository
→ 7. Create B2BService → 8. Create gRPC Server (port 50112)
→ 9. Health Check → 10. Start Server (goroutine)
→ 11. Wait for SIGTERM → 12. Graceful Shutdown
```

---

## 13. Test Coverage

```
Unit tests (no DB):
  middleware   — 12 tests PASS
  consumers    — 16 tests PASS
  service      —  9 tests PASS

Live DB tests (INSURETECH_LIVE_DB_TESTS=1, run with -p 1):
  repository   — 24 tests PASS
  service      — 18 tests PASS

TOTAL: 37 unit + 33 live DB = 70 tests ALL PASS
```

### Key Test Scenarios

- ✅ Super admin can call all org management RPCs without org ID
- ✅ B2B admin assigned via AssignOrgAdmin gets correct b2b_org_admin role
- ✅ B2B org HR managers/viewers get partner_user role
- ✅ ResolveMyOrganisation always passes through (no auth)
- ✅ Non-system users without org context get PermissionDenied
- ✅ Duplicate role assignments silently ignored (idempotent)
- ✅ Department delete refused if active employees exist
- ✅ Soft delete sets deleted_at, not hard delete

---

## 14. Known Issues & TODOs

| Status | Item |
|--------|------|
| ⚠️ | Role IDs are looked up at runtime (not hardcoded), but no retry on lookup failure |
| ⚠️ | No circuit breaker for AuthZ service calls |
| ⚠️ | No caching for AuthZ CheckAccess calls (every request hits AuthZ gRPC) |
| 📋 | Event replay mechanism for DLQ messages |
| 📋 | Rate limiting at gateway level |
| 📋 | Monitoring dashboards for Kafka lag + AuthZ latency |
