# B2B Portal — Complete Org Management Reference
## InsureTech: Wiring, Components, UI/UX & Flows for Super Admin & B2B Admin

> Based on: `AUTHN_AUTHZ_REFERENCE.md` · `B2B_PORTAL_REFERENCE.md` · `B2B_SERVICE_REFERENCE.md`

---

## Table of Contents

1. [System Overview](#1-system-overview)
2. [Full Stack Architecture & Wiring](#2-full-stack-architecture--wiring)
3. [Authentication Flow & Cookie Wiring](#3-authentication-flow--cookie-wiring)
4. [Portal Headers & SDK Client](#4-portal-headers--sdk-client)
5. [Route Guards & Middleware](#5-route-guards--middleware)
6. [API Routes — BFF Layer](#6-api-routes--bff-layer)
7. [Frontend Components & Pages](#7-frontend-components--pages)
8. [State Management & Hooks](#8-state-management--hooks)
9. [AuthZ Wiring & Permission Model](#9-authz-wiring--permission-model)
10. [Super Admin — Org Management UI/UX & Flows](#10-super-admin--org-management-uiux--flows)
11. [B2B Admin — Org Management UI/UX & Flows](#11-b2b-admin--org-management-uiux--flows)
12. [Role Mapping & Event-Driven Assignment](#12-role-mapping--event-driven-assignment)
13. [Database Schema](#13-database-schema)
14. [Changes & Bug Fixes Log](#14-changes--bug-fixes-log)
15. [Error Handling & Security Model](#15-error-handling--security-model)

---

## 1. System Overview

### 1.1 What Is This System?

The InsureTech B2B Portal is a **Next.js (App Router)** frontend that serves as a management portal for:
- **Super Admins** (`SYSTEM_USER`) — platform-wide org management across ALL organisations
- **B2B Org Admins** (`B2B_ORG_ADMIN`) — management scoped to their single organisation
- **HR Managers / Viewers** — limited access to their organisation data

The portal is a **stateless BFF (Backend-for-Frontend)**. It forwards cookies + headers to backend gRPC services and renders results. No sensitive state lives in the browser.

### 1.2 Key Actors

| Actor | UserType (AuthN) | Portal Role Cookie | Casbin Domain | Scope |
|-------|-------------------|-------------------|---------------|-------|
| Super Admin | `SYSTEM_USER` | `SYSTEM_ADMIN` | `system:root` | ALL organisations |
| B2B Org Admin | `B2B_ORG_ADMIN` | `B2B_ORG_ADMIN` | `b2b:{org_id}` | Own organisation only |
| Business Admin | Default | `BUSINESS_ADMIN` | `b2b:{org_id}` | Own organisation only |
| HR Manager | — | `HR_MANAGER` | `b2b:{org_id}` | Own organisation only |
| Viewer | — | `VIEWER` | `b2b:{org_id}` | Own organisation only |

### 1.3 Core Entities

| Entity | Table | Key Fields |
|--------|-------|------------|
| Organisation | `b2b_schema.organisations` | organisation_id, name, code, status, total_employees |
| Org Member | `b2b_schema.org_members` | member_id, organisation_id, user_id, role |
| Department | `b2b_schema.departments` | department_id, name, business_id, employee_no |
| Employee | `b2b_schema.employees` | employee_uuid, name, department_id, business_id, status |
| Purchase Order | `b2b_schema.purchase_orders` | purchase_order_id, PO number, department_id, plan_id |

### 1.4 Capability Matrix

| Capability | Super Admin | B2B Org Admin | HR Manager | Viewer |
|-----------|-------------|---------------|------------|--------|
| List ALL organisations | YES | NO | NO | NO |
| Create organisation | YES | NO | NO | NO |
| Approve organisation | YES | NO | NO | NO |
| Delete organisation | YES | NO | NO | NO |
| View own organisation | YES | YES | YES | YES |
| Update own organisation | YES | YES | NO | NO |
| Assign org admin | YES | NO | NO | NO |
| List members | YES | YES | NO | NO |
| Add member | YES | YES | NO | NO |
| Remove member | YES | YES | NO | NO |
| Manage departments | YES | YES | YES | NO |
| Manage employees | YES | YES | YES | NO |
| View employees | YES | YES | YES | YES |
| Create purchase orders | YES | YES | YES | NO |

---

## 2. Full Stack Architecture & Wiring

### 2.1 End-to-End Architecture

```
Browser
  |
  | HTTPS (httpOnly cookie: session_token, csrf_token)
  v
Next.js App Router (b2b_portal) — Port 3000
  |
  |  app/api/* — BFF Server-side API routes
  |  resolvePortalHeaders() → injects x-portal, x-user-id, x-business-id
  |
  | HTTP → REST Gateway (Port 8080)
  v
REST Gateway
  |  1. Validates session cookie → gets user_id via AuthN.ValidateToken
  |  2. Calls B2B.ResolveMyOrganisation(user_id) → gets org_id
  |  3. Injects x-business-id metadata
  |
  | gRPC (Port 50112)
  v
B2B gRPC Server
  |  AuthZ Interceptor → CheckAccess before EVERY handler
  |  B2BHandler → B2BService → PortalRepository
  |
  | gRPC (Port 50054)       gRPC (Port 50053)
  v                         v
AuthZ Service              AuthN Service
  Casbin PERM model          JWT / Session validation
  PostgreSQL authz_schema    PostgreSQL authn_schema
  |
  v
PostgreSQL b2b_schema
  organisations, org_members, departments, employees, purchase_orders
```

### 2.2 Technology Stack

| Layer | Technology |
|-------|-----------|
| Frontend | Next.js (App Router) · TypeScript · Radix UI · Tailwind CSS |
| SDK | `@lifeplus/insuretech-sdk` (auto-generated from proto, local .tgz) |
| HTTP Client | `@bufbuild/protobuf` connect transport |
| Tables | `@tanstack/react-table` |
| Charts | `recharts` |
| Icons | `lucide-react` |
| Backend | Go microservices (gRPC + Protocol Buffers) |
| Auth | AuthN gRPC :50053 · AuthZ (Casbin) gRPC :50054 |
| B2B Service | gRPC :50112 |
| Database | PostgreSQL (`b2b_schema`, `authn_schema`, `authz_schema`) |
| Messaging | Kafka (event-driven role assignment) |

### 2.3 Folder Structure

```
b2b_portal/
├── app/
│   ├── api/
│   │   ├── auth/
│   │   │   ├── login/route.ts          # POST — normalize mobile, set cookies
│   │   │   ├── logout/route.ts         # POST — invalidate session, clear cookies
│   │   │   └── session/route.ts        # GET — validate session, return principal
│   │   ├── organisations/
│   │   │   ├── route.ts                # GET list, POST create
│   │   │   ├── [id]/route.ts           # GET, PATCH, DELETE, POST approve
│   │   │   ├── [id]/admins/route.ts    # POST create/bootstrap admin
│   │   │   ├── [id]/members/route.ts   # GET list, POST add member
│   │   │   ├── [id]/members/[memberId]/route.ts  # DELETE remove member
│   │   │   └── me/route.ts             # GET — resolve current org context
│   │   ├── departments/
│   │   │   ├── route.ts                # GET list, POST create
│   │   │   └── [id]/route.ts           # GET, PATCH, DELETE
│   │   ├── employees/
│   │   │   ├── route.ts                # GET list, POST create
│   │   │   └── [id]/route.ts           # GET, PATCH, DELETE
│   │   └── purchase-orders/
│   │       ├── route.ts, [id]/, catalog/
│   ├── login/page.tsx
│   ├── organisations/page.tsx
│   ├── departments/page.tsx
│   ├── employees/page.tsx
│   ├── purchase-orders/page.tsx
│   ├── layout.tsx
│   └── page.tsx                        # Dashboard
├── components/
│   ├── auth/login-form.tsx
│   ├── dashboard/
│   │   ├── organisations/
│   │   │   ├── Organisations.tsx       # Main list (super-admin only)
│   │   │   └── data-table/             # TanStack table columns/filters
│   │   ├── departments/
│   │   │   ├── Departments.tsx
│   │   │   └── data-table/
│   │   ├── employees/
│   │   │   ├── employees-table.tsx     # Org lock/dropdown logic
│   │   │   └── data-table/
│   │   ├── overview-activity/          # Activity feed
│   │   └── stats-cards/               # Dashboard KPI cards
│   ├── organisations/
│   │   ├── org-detail-panel.tsx        # Side panel: Info / Members / Departments tabs
│   │   └── org-member-panel.tsx        # Member list + add/remove/assign-admin
│   ├── modals/
│   │   ├── add-organisation-modal.tsx  # Create (org fields + admin) / Edit (2 tabs)
│   │   ├── add-department-modal.tsx
│   │   └── add-employee-modal.tsx      # 3 sections: Personal / Employment / Insurance
│   └── ui/                             # Radix UI primitives
├── src/
│   ├── lib/
│   │   ├── sdk/
│   │   │   ├── b2b-sdk-client.ts       # makeSdkClient() — builds SDK with all headers
│   │   │   └── session-headers.ts      # resolvePortalHeaders()
│   │   ├── clients/
│   │   │   ├── auth-client.ts
│   │   │   ├── organisation-client.ts
│   │   │   ├── department-client.ts
│   │   │   ├── employee-client.ts
│   │   │   ├── api-client.ts           # parseJson, ApiResult pattern
│   │   │   └── b2b-dashboard-client.ts
│   │   ├── auth/
│   │   │   ├── backend-auth.ts
│   │   │   ├── session.ts
│   │   │   └── session-store.ts        # In-memory Map, TTL 12h
│   │   └── types/
│   │       ├── auth.ts
│   │       ├── b2b.ts
│   │       └── employee-form.ts
│   └── hooks/
│       ├── useCrudList.ts              # Generic CRUD list hook
│       ├── useOrganisationForm.ts      # Org form validation
│       ├── useEmployeeForm.ts          # Employee form validation
│       └── useToast.ts
├── middleware.ts                        # Edge middleware — auth cookie guards
└── next.config.ts
```

---

## 3. Authentication Flow & Cookie Wiring

### 3.1 Login Flow (Step by Step)

```
1. User opens /login — enters mobile number + password

2. Client POSTs to /api/auth/login
   { mobile: "01712345678", password: "..." }

3. Server-side route normalizes mobile to BD E.164:
   "01712345678" → "+8801712345678"
   Regex: /^880(13|14|15|16|17|18|19)\d{8}$/

4. Calls backend loginWithMobile(mobile, password)
   → AuthN gRPC :50053 — Login RPC
   → Validates credentials (Argon2id password hash)
   → Creates SERVER_SIDE session → returns session_token

5. Portal reads UserType from response:
   SYSTEM_USER     → portal_role = "SYSTEM_ADMIN",   businessId = ""
   B2B_ORG_ADMIN   → portal_role = "B2B_ORG_ADMIN",  businessId = /organisations/me
   Default         → portal_role = "BUSINESS_ADMIN",  businessId = /organisations/me

6. For B2B roles, calls /api/organisations/me to get org_id
   → Backend: B2B.ResolveMyOrganisation(user_id)

7. Sets 5 cookies on response:
   session_token   httpOnly=true  secure=true  sameSite=strict  maxAge=12h
   csrf_token      httpOnly=true  secure=true  sameSite=lax     maxAge=12h
   portal_role     httpOnly=false (edge middleware can read)
   portal_user_id  httpOnly=false (inject x-user-id header)
   portal_biz_id   httpOnly=false (inject x-business-id header)

8. Redirects to / (Dashboard)
```

### 3.2 Cookie Reference

| Cookie | httpOnly | sameSite | Purpose | Value |
|--------|----------|----------|---------|-------|
| `session_token` | YES | strict | Real auth boundary — backend validates | Encrypted session token from AuthN |
| `csrf_token` | YES | lax | CSRF protection for mutations | `crypto.randomBytes(16).toString('hex')` |
| `portal_role` | NO | — | Edge middleware UX routing | `SYSTEM_ADMIN` / `B2B_ORG_ADMIN` / `BUSINESS_ADMIN` / `HR_MANAGER` / `VIEWER` |
| `portal_user_id` | NO | — | Inject `x-user-id` into SDK headers | User UUID |
| `portal_biz_id` | NO | — | Inject `x-business-id` into SDK headers | Org UUID (empty for SYSTEM_ADMIN) |

### 3.3 Mobile Number Normalization

```
Accepted inputs → Normalized E.164:
  01712345678         → +8801712345678
  +8801712345678      → +8801712345678
  8801712345678       → +8801712345678
  008801712345678     → +8801712345678
  +880 171-234-5678   → +8801712345678

Valid prefixes: 013, 014, 015, 016, 017, 018, 019
```

### 3.4 Login Error Mapping

| Backend Error | User-Facing Message |
|--------------|---------------------|
| 401 / "unauthenticated" | "Mobile number or password is incorrect." |
| 429 / "locked" | "Account temporarily locked." |
| "inactive" / "suspended" | "Account not active." |
| 400 / "invalid mobile" | "Invalid mobile number." |
| 500+ | "Service temporarily unavailable." |

### 3.5 Logout Flow

```
1. Client POSTs /api/auth/logout
2. Server calls AuthN.Logout(session_token) — invalidates backend session
3. Clears all 5 cookies (maxAge=0)
4. Redirects to /login
```

### 3.6 Session Validation

```
GET /api/auth/session
  → Reads session_token cookie
  → Calls AuthN.ValidateToken(session_token)
  → Returns principal: { userId, role, portal, businessId }
  → Used by middleware and client-side session checks
```

---

## 4. Portal Headers & SDK Client

### 4.1 resolvePortalHeaders() — The Critical Wiring Function

Every BFF API route calls this first. It bridges cookies → SDK headers:

```typescript
// src/lib/sdk/session-headers.ts
async function resolvePortalHeaders(request: NextRequest) {
  const sessionToken = request.cookies.get('session_token')?.value
  if (!sessionToken) return null  // → 401

  const role       = request.cookies.get('portal_role')?.value
  const userId     = request.cookies.get('portal_user_id')?.value
  const bizId      = request.cookies.get('portal_biz_id')?.value

  // Map role → portal enum
  const portal = role === 'SYSTEM_ADMIN' ? PORTAL_SYSTEM : PORTAL_B2B

  return { portal, userId, businessId: bizId, tenantId: DEFAULT_TENANT_ID }
}
```

### 4.2 makeSdkClient() — Header Injection

```typescript
// src/lib/sdk/b2b-sdk-client.ts
function makeSdkClient(sessionToken, csrfToken, headers) {
  return new InsureTechSDK({
    headers: {
      'cookie':         `session_token=${sessionToken}; csrf_token=${csrfToken}`,
      'X-CSRF-Token':   csrfToken,           // on state-changing requests
      'x-portal':       headers.portal,      // PORTAL_SYSTEM | PORTAL_B2B
      'x-business-id':  headers.businessId,  // org UUID (empty for super-admin)
      'x-user-id':      headers.userId,      // user UUID
      'x-tenant-id':    headers.tenantId,    // default: 00000000-0000-0000-0000-000000000001
    }
  })
}
```

### 4.3 Headers per Role

| Cookie `portal_role` | `x-portal` sent | `x-business-id` | Casbin Domain resolved |
|----------------------|-----------------|-----------------|------------------------|
| `SYSTEM_ADMIN` | `PORTAL_SYSTEM` | _(empty)_ | `system:root` |
| `B2B_ORG_ADMIN` | `PORTAL_B2B` | org UUID | `b2b:{org_id}` |
| `BUSINESS_ADMIN` | `PORTAL_B2B` | org UUID | `b2b:{org_id}` |
| `HR_MANAGER` | `PORTAL_B2B` | org UUID | `b2b:{org_id}` |
| `VIEWER` | `PORTAL_B2B` | org UUID | `b2b:{org_id}` |

### 4.4 AuthZ Interceptor Header Processing (Backend)

```
B2B gRPC Server receives gRPC call with metadata:
  x-user-id       → who is calling
  x-portal        → PORTAL_SYSTEM or PORTAL_B2B
  x-business-id   → org_id (empty for super-admin)
  x-tenant-id     → tenant scope

Interceptor logic:
  1. Is method ResolveMyOrganisation?
     YES → pass through (bootstrap, no auth check)
  2. Is x-portal == PORTAL_SYSTEM?
     YES → domain = "system:root" (skip x-business-id check)
  3. Is x-portal == PORTAL_B2B?
     x-business-id empty? → PermissionDenied
     x-business-id present? → domain = "b2b:{org_id}"
  4. Map method to resource + action:
     Get*/List*     → GET    svc:b2b/*
     Create*/Add*/Assign* → POST   svc:b2b/*
     Update*        → PATCH  svc:b2b/*
     Delete*/Remove* → DELETE svc:b2b/*
  5. Call AuthZ.CheckAccess(user_id, domain, svc:b2b/*, action)
     ALLOW → proceed to handler
     DENY  → PermissionDenied
```

---

## 5. Route Guards & Middleware

### 5.1 Edge Middleware (middleware.ts)

Runs at the CDN edge before any page/API route renders.

```
Public paths (always allow):
  /login
  /api/auth/login

Bypassed paths (no auth check):
  /_next/*
  /public/*
  /logos/*
  /favicon.ico
  /api/auth/session
  /api/auth/logout

Protected path logic:
  1. No session_token cookie → redirect to /login?next={originalPath}
  2. Has session_token + on /login → redirect to /
  3. Has session_token + on protected route → allow (cookie presence only)

IMPORTANT: Middleware does NOT validate the token.
Real authentication happens backend-side on every API call.
```

### 5.2 API Route Auth Pattern

Every BFF API route follows this pattern:

```typescript
// Example: app/api/organisations/route.ts
export async function GET(request: NextRequest) {
  // Step 1: Resolve portal headers from cookies
  const headers = await resolvePortalHeaders(request)
  if (!headers) return Response.json({ error: 'Unauthorized' }, { status: 401 })

  // Step 2: Build SDK client with all headers
  const sdk = makeSdkClient(
    request.cookies.get('session_token')!.value,
    request.cookies.get('csrf_token')!.value,
    headers
  )

  // Step 3: Call backend SDK method
  const result = await sdk.listOrganisations({ ... })

  // Step 4: Return response
  return Response.json(result)
}
```

### 5.3 Organisation Scoping — The /organisations/me Pattern

The most critical UI logic. Determines what a user can see:

```typescript
// In employees-table.tsx, departments, etc.
const meResult = await organisationClient.getMe()

if (meResult.organisation?.id) {
  // B2B_ORG_ADMIN / BUSINESS_ADMIN / HR_MANAGER / VIEWER
  // LOCKED to their own organisation
  setResolvedOrg(meResult.organisation)        // readonly label
  setOrganisations([meResult.organisation])    // only their org in dropdown
  setSelectedOrgId(meResult.organisation.id)  // auto-selected, not changeable

} else {
  // SYSTEM_ADMIN — can see all organisations
  const listResult = await organisationClient.list()
  setOrganisations(listResult.organisations)   // full dropdown of all orgs
}
```

```typescript
// Server-side: app/api/organisations/me/route.ts
if (session?.principal.role === 'SYSTEM_ADMIN') {
  return Response.json({ ok: true, organisation: null })
  // null = super admin, no org lock
}
// Otherwise: call B2B.ResolveMyOrganisation(user_id) → return actual org
```

**UI Result:**
- B2B Admin sees: readonly label `"This B2B admin is locked to their organisation context."`
- Super Admin sees: dropdown `"Select organisation"` listing all orgs

---

## 6. API Routes — BFF Layer

### 6.1 Auth Routes

| Method | Path | What It Does |
|--------|------|-------------|
| POST | `/api/auth/login` | Normalize mobile → AuthN.Login → set 5 cookies → redirect |
| POST | `/api/auth/logout` | AuthN.Logout → clear all cookies → redirect to /login |
| GET | `/api/auth/session` | AuthN.ValidateToken → return principal (userId, role, portal) |

### 6.2 Organisation Routes

| Method | Path | Who Can Call | Backend SDK Call |
|--------|------|-------------|-----------------|
| GET | `/api/organisations` | Super Admin | `sdk.listOrganisations()` |
| POST | `/api/organisations` | Super Admin | `sdk.createOrganisation()` + optional admin bootstrap |
| GET | `/api/organisations/{id}` | Super Admin, own B2B Admin | `sdk.getOrganisation(id)` |
| PATCH | `/api/organisations/{id}` | Super Admin, own B2B Admin | Direct HTTP PATCH |
| DELETE | `/api/organisations/{id}` | Super Admin | `sdk.deleteOrganisation(id)` |
| POST | `/api/organisations/{id}` | Super Admin | Approve/activate org |
| GET | `/api/organisations/me` | All authenticated | `B2B.ResolveMyOrganisation(user_id)` |

### 6.3 Org Member & Admin Routes

| Method | Path | Who Can Call | Purpose |
|--------|------|-------------|---------|
| GET | `/api/organisations/{id}/members` | Super Admin, B2B Admin | List all org members |
| POST | `/api/organisations/{id}/members` | Super Admin, B2B Admin | Add member (userId + role) |
| DELETE | `/api/organisations/{id}/members/{memberId}` | Super Admin, B2B Admin | Remove member |
| POST | `/api/organisations/{id}/admins` | Super Admin | Mode A: Create new user + assign B2B_ORG_ADMIN (email, password, mobileNumber, fullName); Mode B: Promote existing member (memberId) |
| POST | `/api/organisations/{id}/assign-admin` | Super Admin | Promote existing member to B2B Admin (calls `assignAdmin(orgId, memberId)` correctly) |
| POST | `/api/organisations/{id}/approve` | Super Admin | Approve pending org → PATCH status to ACTIVE |

### 6.4 Department Routes

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/departments` | Super Admin: `business_id` from query param. B2B Admin: from session cookie / `/organisations/me` fallback. Scoped by organisationId always. |
| POST | `/api/departments` | Body must include `{ name, businessId }` — both required. Backend rejects if businessId missing. |
| GET | `/api/departments/{id}` | Standard CRUD |
| PATCH | `/api/departments/{id}` | Partial update (dynamic SET clause in SQL) |
| DELETE | `/api/departments/{id}` | Backend refuses if dept has active employees |

### 6.5 Employee Routes

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/employees?page_size=&business_id=&department_id=` | Filtered by org + optional department |
| POST | `/api/employees` | Create with business_id context |
| GET | `/api/employees/{id}` | **mapViewFull()** returns ALL form fields: email, mobile_number, gender, date_of_birth, date_of_joining, department_id, business_id, insurance_category, assigned_plan_id, coverage_amount. insuranceCategory mapped: proto enum string → numeric form value (1=Health, 2=Life, 3=Auto, 4=Travel) |
| PATCH | `/api/employees/{id}` | Partial update |
| DELETE | `/api/employees/{id}` | Soft delete (sets deleted_at) |

### 6.6 Purchase Order Routes

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/purchase-orders/catalog` | Merged: DB plans + seeded fallback (Seba/Surokkha/Verosa). **NEW:** resolvePortalHeaders() added — was calling makeSdkClient(request) with NO headers → auth failed. Fix: makeSdkClient(request, hdrs) with resolved headers |
| GET | `/api/purchase-orders` | List POs for org |
| GET | `/api/purchase-orders/{id}` | Single PO |
| POST | `/api/purchase-orders` | Auto-generates `PO-YYYYMMDD-XXXXXXXX`, calculates premium |
| DELETE | `/api/purchase-orders/{id}` | **NEW:** Delete PO with confirm dialog |

### 6.7 Dashboard Routes (NEW)

| Method | Path | Purpose |
|--------|------|---------|
| GET | `/api/dashboard/stats` | Fetches KPI stats. Super Admin: totalOrganisations, pendingOrganisations, totalEmployees, totalDepartments, activePurchaseOrders. B2B Admin: totalMembers, totalDepartments, totalEmployees, activePurchaseOrders. Returns: `{ ok, stats, role }` |
| GET | `/api/dashboard/activity` | Fetches last 5 from each source, merges, sorts by createdAt DESC, returns top 10. Super Admin: recent orgs + employees + departments + POs. B2B Admin: recent employees + departments + POs + member joins. timeAgo() humanised timestamps. Returns: `{ ok, activities: [{id, type, title, subtitle, createdAt}] }` |

---

## 7. Frontend Components & Pages

### 7.1 Pages (Updated with New Routes)

| Page | Path | Role Access | Description |
|------|------|-------------|-------------|
| Login | `/login` | Public | Mobile + password form |
| Dashboard | `/` | All authenticated | KPI cards + activity feed |
| Organisations | `/organisations` | Super Admin only | Full org list + detail panel |
| Departments | `/departments` | All | Org-scoped dept management |
| Employees | `/employees` | All | Org-scoped employee management |
| Purchase Orders | `/purchase-orders` | B2B Admin + Super Admin | PO list + create |

### 7.2 Key Components — Detail

#### LoginForm (`components/auth/login-form.tsx`)
```
Fields:
  - Mobile Number (text, normalized on submit)
  - Password (password, toggle visibility)

Behaviour:
  - On submit: POST /api/auth/login
  - Error states: "Invalid credentials", "Account locked", etc.
  - Redirect: reads ?next= query param, or defaults to /
  - Loading state: spinner on button
```

#### Organisations Page (`components/dashboard/organisations/`)
```
OrganisationsList.tsx
  ├── DataTable (TanStack)
  │   Columns: Name | Code | Status | Employees | Actions
  │   Filter: search by name/code, filter by status
  │   Sort: name, created date
  │   Actions per row:
  │     [Eye icon] → opens OrgDetailPanel as side drawer
  │     [Edit icon] → opens AddOrganisationModal in edit mode
  │     [Approve] → POST /api/organisations/{id}/approve (if status=PENDING)
  │     [Delete] → confirm dialog → DELETE /api/organisations/{id}
  │
  ├── [+ Add Organisation] button (Super Admin only)
  │     → opens AddOrganisationModal in create mode
  │
  └── OrgDetailPanel (side drawer/sheet)
        Tabs:
          [Info]    → org fields: name, code, status, contact, address, employee count
          [Members] → OrgMemberPanel (list + invite + remove)
          [Departments] → inline dept list for this org

AddOrganisationModal.tsx — FIXED
  Two modes:
  ┌─── CREATE MODE ────────────────────────────────────────┐
  │ Tab 1: Organisation Info                               │
  │   Fields: Name, Code, Status, Total Employees         │
  │   Address: Street, City, State, Zip, Country          │
  │   Contact: Email, Phone                               │
  │                                                        │
  │ Tab 2: Admin Setup (optional at creation)              │
  │   Mode A: Create new user → email, password, mobile, fullName (camelCase) │
  │   Mode B: Promote existing → memberId only            │
  │   Bug fixed: was sending mobile_number (snake_case), backend expects mobileNumber (camelCase) │
  │   [Skip] or [Create with Admin]                        │
  └────────────────────────────────────────────────────────┘
  ┌─── EDIT MODE ──────────────────────────────────────────┐
  │ Same two tabs pre-populated with existing data         │
  │ Admin tab shows current admin + option to change       │
  │ Edit fields changed from uncontrolled defaultValue → controlled value │
  │ useEffect syncs all edit state when org.id or open changes (prevents stale data when switching orgs) │
  └────────────────────────────────────────────────────────┘

OrgMemberPanel (`components/organisations/org-member-panel.tsx`) — FIXED
  ├── Member list (name, role, status, joined date)
  │   UserIdCell: shows truncated UUID (a1b2c3d4…) with one-click copy button
  │   Role badges: colour-coded (purple=Admin, blue=HR, grey=Viewer)
  ├── [Add Member] → opens inline form
  │     Fields: User ID or email, Role dropdown
  │     Roles: B2B_ORG_ADMIN (Super Admin only), BUSINESS_ADMIN, HR_MANAGER, VIEWER
  ├── [Remove] per member → confirm → DELETE /api/organisations/{id}/members/{memberId}
  ├── Promote button (🛡️, Super Admin only): calls assignAdmin(orgId, memberId) via correct /assign-admin endpoint
  ├── Error shown as inline red banner, not alert()
  └── Removed duplicate assign-existing-user form, consolidated into promote icon
```

#### Departments Page (`components/dashboard/departments/`)
```
DepartmentsTable.tsx (TanStack DataTable) — REBUILT with org-lock pattern
  Columns: Name | Employee No | Business (Super Admin only) | Actions
  Filters: search by name, filter by org (Super Admin only)
  Actions: 👁 View · ✏️ Edit · 🗑 Delete (replaced <details> dropdown)
  
  Super Admin: Select dropdown of all orgs → filter departments
  B2B Admin: Readonly label showing org name, auto-locked
  
  departmentClient.list(50, 0, selectedOrgId) — always scoped by org
  Passes organisationId to AddDepartmentModal and DataTable

AddDepartmentModal.tsx — FIXED
  Added organisationId prop
  Fields:
    - Name (required)
    - Employee Number (readonly, auto-counted)
    - Business / Organisation (Super Admin: dropdown, B2B Admin: locked)

  POST body now sends `{ name, businessId }` — was only sending `{ name }` (backend rejected)

  Super Admin: can select any org from dropdown
  B2B Admin:   org auto-locked from /organisations/me
```

#### Employees Page (`components/dashboard/employees/`)
```
EmployeesTable.tsx (TanStack DataTable)
  Columns: Name | NID | DOB | Gender | Department | Plan | Status | Actions
  Filters:
    - Org selector (Super Admin: dropdown of all orgs)
                   (B2B Admin: locked, shows org name label)
    - Department filter (populates from selected org)
    - Status filter: Active / Inactive
  Actions: 👁 View · ✏️ Edit · 🗑 Delete (replaced <details> dropdown)
  
  View button: opens AddEmployeeModal with employeeUuid (fetches full record via GET /api/employees/{uuid})

AddEmployeeModal.tsx — COMPLETE FIX
  useEmployeeForm hook declared BEFORE departments useEffect (was after — caused TS2448)
  Department fetch now uses `organisationId || values.businessId` — resolves org in edit mode
  values.businessId added to dependency array → depts reload when edit record loads
  
  Section 1: Personal Information
    - First Name, Last Name (required)
    - Email (now populated from GET /api/employees/{id})
    - National ID (required, unique per org)
    - Date of Birth (date picker)
    - Gender (Male / Female / Other)
    - Contact Number, Emergency Contact, Address

  Section 2: Employment Details
    - Employee Code (auto-generated or manual)
    - Department (dropdown filtered by selected org)
    - Designation / Job Title
    - Join Date

  Section 3: Insurance Information
    - Plan (dropdown from /api/purchase-orders/catalog)
    - Sum Insured / Coverage Amount
    - Policy Start Date, End Date
    - Beneficiary Name, Beneficiary Relationship

useEmployeeForm (src/hooks/useEmployeeForm.ts) — FIXED
  In edit mode: now fetches full employee record from GET /api/employees/{uuid} on mount
  mapApiToForm(EmployeeFullRecord) maps ALL fields: email, mobileNumber, gender, dateOfBirth, dateOfJoining, departmentId, businessId, insuranceCategory, coverageAmount, assignedPlanId
  loadingRecord state exposed → modal shows spinner while fetching
  Previously only used sparse initialValues from table row (only name + employeeId)
```

#### Purchase Orders Page (`components/dashboard/purchase-orders/`)
```
PurchaseOrdersTable.tsx — FIXED: Actions Implemented
  Columns: PO Number | Department | Plan | Premium | Status | Created
  Actions: 👁 View · 🗑 Delete (replaced hardcoded `—` stub with TODO comment)
  
  View → PODetailSheet right side drawer (all PO fields)
  Delete → confirm dialog → DELETE /api/purchase-orders/{id} → refresh
  
  Org context resolved on mount via organisationClient.getMe()
  departmentClient.list(200, 0, resolvedOrgId) — B2B Admin sees only their depts

AddPurchaseOrderModal — COMPLETE REWRITE
  All labels changed from sr-only → visible block labels
  insuranceCategory: auto-derived from selected plan via useEffect, shown as readonly field
  Department dropdown: visible label + error message
  Product Plan dropdown: visible label + error message + plan summary card
  Number of Employees / Dependents: visible labels + 2-column grid
  Coverage Amount (BDT): visible label + required validation
  Notes: changed from Input to textarea (3 rows)
  Validation: inline field errors, modal stays open on failure
  Cancel button added
  PurchaseOrderFormInput type now includes insuranceCategory: string
  handleCreatePurchaseOrder returns Promise<boolean> — modal stays open on error
```

#### Dashboard Page (`/`)
```
StatsCards — FIXED: Live Data
  File: components/dashboard/stats-cards/stats-cards.tsx
  Was: imported static mock from lib/stats-cards.ts
  Fix: fetches from /api/dashboard/stats, role-aware cards
  Loading skeleton (4 spinners), error state
  StatsCard props: { title, value, icon (svg path), bgColor, bgIcon (svg path) }
  
  Super Admin view:
    - Total Organisations
    - Pending Organisations
    - Total Employees across all orgs
    - Active Purchase Orders

  B2B Admin view:
    - Total Members
    - Total Departments
    - Total Employees in org
    - Active Purchase Orders for org

OverviewActivity — FIXED: Live Data
  File: components/dashboard/overview-activity/overview-activity.tsx
  Was: imported static mock from lib/overview-activity
  Fix: fetches from /api/dashboard/activity
  Type-based icons: 🏢 org · 👥 member · 👤 employee · 📋 dept · 🗒 PO
  timeAgo() relative timestamps
  Empty/error/loading states
  
  Super Admin:
    - Recent org registrations
    - Recent employee additions
    - Recent department creations
    - Recent purchase orders
  
  B2B Admin:
    - Recent member additions
    - Recent employee additions
    - Recent department creations
    - Recent purchase orders
```

---

## 8. State Management & Hooks

### 8.1 useCrudList — Generic CRUD Hook

```typescript
// src/hooks/useCrudList.ts
// Used by ALL list pages (orgs, departments, employees, POs)

const {
  items,          // the list data
  loading,        // boolean
  error,          // string | null
  pagination,     // { page, pageSize, total }
  create,         // (data) => Promise<void>  — calls POST
  update,         // (id, data) => Promise<void>  — calls PATCH
  remove,         // (id) => Promise<void>  — calls DELETE
  refresh,        // () => void  — re-fetches
  setFilters,     // (filters) => void
} = useCrudList({
  endpoint: '/api/organisations',
  defaultPageSize: 20,
})
```

### 8.2 Session / Auth State

```typescript
// No Redux/Zustand. Auth state is server-side via cookies.
// Client-side reads:
const { data: session } = useSession()
// Calls GET /api/auth/session
// Returns: { userId, role, portal, businessId, ok }

// Role-conditional rendering pattern:
const isSuperAdmin = session?.principal.role === 'SYSTEM_ADMIN'
const isB2BAdmin   = session?.principal.role === 'B2B_ORG_ADMIN'
```

### 8.3 Organisation Context (Client-Side)

```typescript
// useOrganisationContext.ts
// Called in every page that needs org context (departments, employees, POs)

async function resolveOrganisationContext() {
  const me = await fetch('/api/organisations/me').then(r => r.json())

  if (me.organisation) {
    // B2B role — locked to own org
    return {
      isSuperAdmin: false,
      lockedOrgId: me.organisation.id,
      lockedOrgName: me.organisation.name,
      availableOrgs: [me.organisation],
    }
  } else {
    // Super Admin — fetch all orgs for dropdown
    const list = await fetch('/api/organisations').then(r => r.json())
    return {
      isSuperAdmin: true,
      lockedOrgId: null,
      availableOrgs: list.organisations,
    }
  }
}
```

### 8.4 useOrganisationForm

```typescript
// Validates org creation / edit form
const {
  form,       // react-hook-form FormInstance
  errors,     // validation errors
  onSubmit,   // handles create or update
  tab,        // 'info' | 'admin'
  setTab,
} = useOrganisationForm({
  mode: 'create' | 'edit',
  initialData?: OrganisationFormData,
  onSuccess: () => void,
})
```

### 8.5 ApiResult Pattern

```typescript
// All client calls return ApiResult<T>
type ApiResult<T> =
  | { ok: true;  data: T }
  | { ok: false; error: string; status: number }

// Usage:
const result = await organisationClient.list()
if (!result.ok) {
  toast.error(result.error)
  return
}
setOrganisations(result.data.organisations)
```

---

## 9. AuthZ Wiring & Permission Model

### 9.1 Casbin PERM Model (Backend)

```
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act, eft

[role_definition]
g = _, _, _             // g(user, role, domain)

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom)
    && keyMatch2(r.dom, p.dom)
    && keyMatch2(r.obj, p.obj)
    && regexMatch(r.act, p.act)
```

### 9.2 Policy Table (Seeded Defaults)

| Subject | Domain | Object | Action | Effect |
|---------|--------|--------|--------|--------|
| `SYSTEM_ADMIN` | `system:root` | `svc:b2b/*` | `(GET\|POST\|PATCH\|DELETE)` | allow |
| `B2B_ORG_ADMIN` | `b2b:*` | `svc:b2b/organisations/me` | `GET` | allow |
| `B2B_ORG_ADMIN` | `b2b:*` | `svc:b2b/organisations` | `GET\|PATCH` | allow |
| `B2B_ORG_ADMIN` | `b2b:*` | `svc:b2b/members/*` | `GET\|POST\|DELETE` | allow |
| `B2B_ORG_ADMIN` | `b2b:*` | `svc:b2b/departments/*` | `GET\|POST\|PATCH\|DELETE` | allow |
| `B2B_ORG_ADMIN` | `b2b:*` | `svc:b2b/employees/*` | `GET\|POST\|PATCH\|DELETE` | allow |
| `BUSINESS_ADMIN` | `b2b:*` | `svc:b2b/departments/*` | `GET\|POST\|PATCH\|DELETE` | allow |
| `BUSINESS_ADMIN` | `b2b:*` | `svc:b2b/employees/*` | `GET\|POST\|PATCH\|DELETE` | allow |
| `HR_MANAGER` | `b2b:*` | `svc:b2b/employees/*` | `GET\|POST\|PATCH` | allow |
| `VIEWER` | `b2b:*` | `svc:b2b/*` | `GET` | allow |

### 9.3 Role Assignment Flows

#### At Organisation Creation (Super Admin)
```
1. Super Admin fills "Admin Setup" tab in AddOrganisationModal
2. POST /api/organisations/{id}/admins
   body: { mobile, firstName, lastName }
3. Portal calls AuthN.CreateUser(mobile) → gets userId
4. Portal calls AuthZ.AssignRole(userId, B2B_ORG_ADMIN, b2b:{org_id})
5. Kafka event: "org.admin.assigned" → AuthN sends welcome SMS/email
```

#### Existing User Promoted to B2B Admin
```
1. Super Admin opens OrgDetailPanel → Members tab
2. Clicks [Assign Admin] next to a member
3. POST /api/organisations/{id}/assign-admin
   body: { userId }
4. Calls AuthZ.UpdateRole(userId, from: BUSINESS_ADMIN, to: B2B_ORG_ADMIN, domain: b2b:{org_id})
5. Cookie portal_role remains as-is until next login
```

#### When Casbin Policy Is NOT Enough
```
Additional backend enforcement in B2BService:
  - For B2B_ORG_ADMIN: ALWAYS filter by organisation_id from x-business-id header
    Even if policy says "allow", if org_id doesn't match the resource → 403
  - This prevents cross-org data leakage even if policy has a mistake
```

### 9.4 AuthZ Event Flow (Kafka)

```
B2B Service (Producer)
  ├── org.created          → AuthZ seeds domain b2b:{org_id} with default policies
  ├── org.admin.assigned   → AuthZ assigns B2B_ORG_ADMIN role in b2b:{org_id} domain
  ├── org.member.added     → AuthZ assigns BUSINESS_ADMIN role in b2b:{org_id} domain
  ├── org.member.removed   → AuthZ revokes all roles in b2b:{org_id} domain for user
  └── org.deleted          → AuthZ deletes entire b2b:{org_id} domain policies

AuthZ Service (Consumer)
  └── Processes events → updates casbin_rule table in authz_schema
```

---

## 10. Super Admin — Org Management UI/UX & Flows

### 10.1 Navigation Available to Super Admin

```
Sidebar:
  🏠 Dashboard
  🏢 Organisations     ← ONLY visible to Super Admin
  👥 Departments
  👤 Employees
  📋 Purchase Orders
  ⚙️  Settings
```

### 10.2 Page: Organisations (`/organisations`)

#### Layout Wireframe
```
┌─────────────────────────────────────────────────────────────────┐
│  Organisations                              [+ Add Organisation] │
├─────────────────────────────────────────────────────────────────┤
│  [🔍 Search by name or code...]  [Status ▾]  [Sort ▾]           │
├──────────┬──────────────┬────────┬───────────┬──────────────────┤
│ Name     │ Code         │ Status │ Employees │ Actions          │
├──────────┼──────────────┼────────┼───────────┼──────────────────┤
│ Acme Corp│ ACME-001     │ Active │ 245       │ 👁 ✏️  ✅ 🗑       │
│ TechStart│ TECH-002     │ Active │ 52        │ 👁 ✏️  ✅ 🗑       │
│ GlobalCo │ GLOB-003     │Pending │ 0         │ 👁 ✏️  ✅ 🗑       │
│ HealthCo │ HLTH-004     │Inactive│ 18        │ 👁 ✏️  ✅ 🗑       │
├──────────┴──────────────┴────────┴───────────┴──────────────────┤
│  Showing 1–4 of 245            [< Prev]  Page 1 of 62  [Next >] │
└─────────────────────────────────────────────────────────────────┘
```

**Action Icons per row (replaced <details> dropdown with inline buttons):**
- 👁 **View** → opens `OrgDetailPanel` as right side drawer
- ✏️ **Edit** → opens `AddOrganisationModal` in Edit mode
- ✅ **Approve** → visible only when status contains PENDING → calls organisationClient.approve() → if org is open in detail panel, immediately updates status badge to ACTIVE
- 🗑 **Delete** → confirm dialog: "Type org name to confirm" → DELETE

buildOrganisationColumns() now called with 3 args: (onRefresh, onRowClick, onApprove)

**Status badges:**
- `Active` → green pill
- `Pending` → amber pill
- `Inactive` → grey pill
- `Suspended` → red pill

### 10.3 Flow: Create New Organisation

```
Step 1: Click [+ Add Organisation]
  → AddOrganisationModal opens (Tab 1 active)

Step 2: Tab 1 — Organisation Info
  Required:
    ✦ Organisation Name       (text, max 100)
    ✦ Organisation Code       (text, uppercase, unique e.g. ACME-001)
    ✦ Status                  (dropdown: Active | Pending | Inactive)
    ✦ Total Employees         (number, estimated headcount)
  Optional:
    ✦ Contact Email           (email validation)
    ✦ Contact Phone           (BD mobile format)
    ✦ Address Line 1 & 2
    ✦ City, State, ZIP, Country

Step 3: Tab 2 — Admin Setup (optional)
  Option A: Skip — create org without admin now
  Option B: Fill admin:
    ✦ Admin Mobile            (BD mobile, will be normalised)
    ✦ Admin First Name
    ✦ Admin Last Name

Step 4: Click [Create Organisation]
  → POST /api/organisations
    body: { name, code, status, totalEmployees, contactEmail, ... }
  → On success (if admin provided):
    POST /api/organisations/{newOrgId}/admins
    body: { mobile, firstName, lastName }
    → Backend: AuthN.CreateUser → AuthZ.AssignRole(B2B_ORG_ADMIN)
  → Toast: "Organisation created successfully"
  → Table refreshes with new row

Error cases:
  - Code already exists → "Organisation code already taken"
  - Admin mobile already registered → "User already exists, use Assign Admin instead"
  - Validation errors → inline field errors on Tab 1 / Tab 2
```

### 10.4 Flow: View Organisation Detail (OrgDetailPanel)

```
Right side drawer opens with 3 tabs:

TAB 1: Info
  ┌──────────────────────────────────────────────┐
  │ Name:          Acme Corporation               │
  │ Code:          ACME-001                       │
  │ Status:        🟢 Active                      │
  │ Total Employees: 245                          │
  │ Contact Email: admin@acme.com                 │
  │ Phone:         +8801712345678                 │
  │ Address:       123 Main St, Dhaka, BD         │
  │ Created:       2026-01-15                     │
  └──────────────────────────────────────────────┘
  [Edit Organisation]

TAB 2: Members
  ┌──────────────────────────────────────────────┐
  │ [+ Add Member]                               │
  │ Name          Role          Status  Action   │
  │ John Doe      B2B_ORG_ADMIN Active  [Remove] │
  │ Jane Smith    BUSINESS_ADMIN Active  [Remove] │
  │ Bob Jones     HR_MANAGER    Active  [Remove] │
  └──────────────────────────────────────────────┘
  Add Member Form (inline expand):
    - User ID (text input)
    - Role (dropdown: B2B_ORG_ADMIN | BUSINESS_ADMIN | HR_MANAGER | VIEWER)
    [Add]

TAB 3: Departments
  ┌──────────────────────────────────────────────┐
  │ Department Name        Employee No           │
  │ Engineering            42                    │
  │ Operations             18                    │
  │ HR                     8                     │
  └──────────────────────────────────────────────┘
```

### 10.5 Flow: Edit Organisation

```
Click [✏️] on org row OR [Edit Organisation] in detail panel
  → AddOrganisationModal opens in EDIT mode (pre-populated)

Tab 1: Organisation Info (editable)
  All fields pre-filled from existing org data

Tab 2: Admin Management
  Shows current admin(s)
  [Change Admin] → same mobile + name fields
  [Remove Admin] → revoke B2B_ORG_ADMIN role

Click [Save Changes]
  → PATCH /api/organisations/{id}
    body: changed fields only (partial update)
  → Toast: "Organisation updated"
  → Panel refreshes
```

### 10.6 Flow: Approve a Pending Organisation

```
Status = PENDING (amber badge in table)

Option A: Via table action icon [✅ Approve]
Option B: Inside OrgDetailPanel Info tab → [Approve Organisation] button

Confirmation dialog:
  "Are you sure you want to approve Acme Corporation?
   This will activate their account and notify their admin."
  [Cancel] [Approve]

On confirm:
  → POST /api/organisations/{id}/approve
  → Backend: sets status = ACTIVE, sends notification
  → Toast: "Organisation approved and activated"
  → Table row status badge changes to Active
```

### 10.7 Flow: Delete Organisation

```
Click [🗑] on row
  → Confirm dialog:
    "This action is irreversible. Type the organisation name to confirm."
    Input: [____________]
    [Cancel] [Delete]

On match:
  → DELETE /api/organisations/{id}
  → Backend:
    1. Soft-deletes org (sets deleted_at)
    2. Publishes org.deleted event to Kafka
    3. AuthZ consumer deletes all b2b:{org_id} domain policies
    4. All members lose B2B access
  → Toast: "Organisation deleted"
  → Row removed from table
```

### 10.8 Flow: Manage Members (from OrgDetailPanel)

```
Add Member:
  1. Click [+ Add Member] in Members tab
  2. Inline form expands:
     - User ID (UUID of existing user)
     - Role dropdown: B2B_ORG_ADMIN | BUSINESS_ADMIN | HR_MANAGER | VIEWER
  3. [Add] → POST /api/organisations/{id}/members
     body: { userId, role }
  4. AuthZ.AssignRole(userId, role, domain: b2b:{org_id})
  5. Member appears in list

Remove Member:
  1. Click [Remove] next to member
  2. Confirm: "Remove Jane Smith from Acme Corp?"
  3. DELETE /api/organisations/{id}/members/{memberId}
  4. AuthZ.RevokeRole(userId, domain: b2b:{org_id})
  5. Member removed from list

Assign Admin:
  1. Click [Assign Admin] next to a BUSINESS_ADMIN member
  2. Confirm dialog
  3. POST /api/organisations/{id}/assign-admin { userId }
  4. AuthZ: removes BUSINESS_ADMIN, assigns B2B_ORG_ADMIN
  5. Role badge updates to B2B_ORG_ADMIN
```

### 10.9 Super Admin — Department Management

```
/departments page with org dropdown

Org selector: full dropdown of all organisations
  → Changing org selection reloads department list

[+ Add Department]
  → AddDepartmentModal:
    Fields:
      - Department Name (required)
      - Organisation (dropdown — Super Admin can pick any)
    [Create]
    → POST /api/departments { name, businessId }

Edit:  PATCH /api/departments/{id} { name }
Delete: DELETE /api/departments/{id}
  → Backend rejects if department has active employees
```

### 10.10 Super Admin — Employee Management

```
/employees page

Org selector: dropdown of ALL organisations (not locked)
  → On org change: reload dept filter + employee list

Dept filter: dropdown of departments within selected org
  (or "All Departments")

[+ Add Employee]
  → AddEmployeeModal (3 sections — see Section 7.2)
  → businessId = selectedOrgId from dropdown

Employee table actions:
  [👁 View]  → employee detail sheet (read-only all fields)
  [✏️ Edit]  → pre-filled AddEmployeeModal
  [🗑 Delete] → soft delete confirm
```

---

## 11. B2B Admin — Org Management UI/UX & Flows

### 11.1 Navigation Available to B2B Admin

```
Sidebar:
  🏠 Dashboard
  👥 Departments
  👤 Employees
  📋 Purchase Orders
  ⚙️  Organisation Settings   ← limited to own org
  (No "Organisations" menu — cannot see other orgs)
```

### 11.2 Dashboard (B2B Admin View)

```
┌─────────────────────────────────────────────────────────────────┐
│  Welcome, John Doe                          Acme Corporation     │
├──────────────┬──────────────┬──────────────┬────────────────────┤
│ Departments  │  Employees   │ Active POs   │  Pending Actions   │
│     12       │    245       │     8        │        3           │
├──────────────┴──────────────┴──────────────┴────────────────────┤
│  Quick Actions:  [+ Add Employee]  [+ Add Department]  [+ New PO]│
├─────────────────────────────────────────────────────────────────┤
│  Recent Activity                                                 │
│  • 2h ago:  Employee John added to Engineering dept             │
│  • 5h ago:  Purchase Order PO-20260305-A1B2C3D4 created         │
│  • 1d ago:  Department "Operations" updated                     │
└─────────────────────────────────────────────────────────────────┘
```

### 11.3 Organisation Settings Page (`/organisations/settings`)

B2B Admin can view and edit ONLY their own organisation's details.

```
┌─────────────────────────────────────────────────────────────────┐
│  Organisation Settings                         [Save Changes]   │
├─────────────────────────────────────────────────────────────────┤
│  Organisation Name:   [Acme Corporation          ]              │
│  Code:                ACME-001  (read-only)                     │
│  Status:              🟢 Active  (read-only — set by Super Admin)│
│  Total Employees:     [245                       ]              │
├─────────────────────────────────────────────────────────────────┤
│  Contact Information                                            │
│  Email:   [admin@acme.com                        ]              │
│  Phone:   [+8801712345678                        ]              │
├─────────────────────────────────────────────────────────────────┤
│  Address                                                        │
│  Street:  [123 Main Street                       ]              │
│  City:    [Dhaka         ]  State: [Dhaka    ]                  │
│  ZIP:     [1000          ]  Country: [Bangladesh ]              │
└─────────────────────────────────────────────────────────────────┘
```

**What B2B Admin CANNOT change:**
- Organisation Code (readonly)
- Organisation Status (set by Super Admin only)
- Cannot delete their own organisation

**Flow:**
```
1. Navigate to Organisation Settings
2. Portal calls GET /api/organisations/me → populates form
3. Admin edits fields
4. Clicks [Save Changes]
5. PATCH /api/organisations/{ownOrgId}
   body: { name, totalEmployees, contactEmail, contactPhone, address... }
6. Toast: "Organisation updated successfully"
```

### 11.4 Member Management (B2B Admin)

B2B Admin accesses members from Organisation Settings page or Members tab.

```
┌─────────────────────────────────────────────────────────────────┐
│  Team Members                              [+ Invite Member]    │
├──────────────┬──────────────┬─────────────┬─────────────────────┤
│ Name         │ Role         │ Status      │ Actions             │
├──────────────┼──────────────┼─────────────┼─────────────────────┤
│ John Doe     │ B2B_ORG_ADMIN│ Active      │ (You — cannot remove)│
│ Jane Smith   │ BUSINESS_ADMIN│ Active     │ [Change Role] [Remove]│
│ Bob Jones    │ HR_MANAGER   │ Active      │ [Change Role] [Remove]│
│ Alice Brown  │ VIEWER       │ Active      │ [Change Role] [Remove]│
└──────────────┴──────────────┴─────────────┴─────────────────────┘
```

**B2B Admin member constraints:**
- Cannot remove themselves
- Cannot assign `B2B_ORG_ADMIN` role (only Super Admin can)
- Can assign: `BUSINESS_ADMIN`, `HR_MANAGER`, `VIEWER`
- Can remove any member except themselves

**Invite Member Flow:**
```
1. Click [+ Invite Member]
2. Modal opens:
   Fields:
     - User ID (UUID of existing portal user)
     - Role: [BUSINESS_ADMIN ▾]  (B2B_ORG_ADMIN option hidden)
   [Send Invite]

3. POST /api/organisations/{ownOrgId}/members
   body: { userId, role: "BUSINESS_ADMIN" | "HR_MANAGER" | "VIEWER" }

4. Backend:
   - Validates user exists in AuthN
   - Checks user not already member
   - Inserts into org_members
   - Publishes org.member.added Kafka event
   - AuthZ assigns role in b2b:{org_id} domain

5. Toast: "Member added successfully"
6. Member appears in list
```

**Remove Member Flow:**
```
1. Click [Remove] next to member
2. Confirm: "Remove Bob Jones from Acme Corporation?
            They will lose all access immediately."
   [Cancel] [Remove]

3. DELETE /api/organisations/{ownOrgId}/members/{memberId}

4. Backend:
   - Removes from org_members
   - Publishes org.member.removed Kafka event
   - AuthZ revokes all roles in b2b:{org_id} domain

5. Member removed from list immediately
```

### 11.5 Department Management (B2B Admin)

```
/departments page — ORG LOCKED

Org display:
  ┌──────────────────────────────────────────────────────┐
  │ 🔒 Acme Corporation                                  │
  │    This view is locked to your organisation context. │
  └──────────────────────────────────────────────────────┘

Department Table:
┌─────────────────────────────────────────────────────────┐
│ Departments — Acme Corporation         [+ Add Dept]     │
├──────────────────┬──────────────────┬───────────────────┤
│ Department Name  │ Employee Count   │ Actions           │
├──────────────────┼──────────────────┼───────────────────┤
│ Engineering      │ 42               │ [✏️ Edit] [🗑 Del] │
│ Operations       │ 18               │ [✏️ Edit] [🗑 Del] │
│ HR               │ 8                │ [✏️ Edit] [🗑 Del] │
│ Finance          │ 15               │ [✏️ Edit] [🗑 Del] │
└──────────────────┴──────────────────┴───────────────────┘
```

**Add Department Flow:**
```
1. Click [+ Add Department]
2. Modal:
   - Name (required)
   - Organisation: "Acme Corporation" (read-only, auto-filled)
3. [Create] → POST /api/departments { name, businessId: ownOrgId }
4. Toast: "Department created"
5. Table refreshes
```

**Delete Department:**
```
1. [🗑 Delete] → Confirm dialog
2. DELETE /api/departments/{id}
3. Backend check: if any active employees → return 409
   "Cannot delete department with active employees.
    Please reassign or remove employees first."
4. On success: row removed, toast shown
```

### 11.6 Employee Management (B2B Admin)

```
/employees page — ORG LOCKED

Org display: readonly label (no dropdown)
Dept filter: dropdown of OWN org's departments only

Employee Table:
┌────────────────────────────────────────────────────────────────┐
│ Employees — Acme Corporation                  [+ Add Employee] │
├──────────┬────────┬────────────┬─────────────┬────────────────┤
│ Name     │ Dept   │ Plan       │ Status      │ Actions        │
├──────────┼────────┼────────────┼─────────────┼────────────────┤
│ John Doe │ Eng    │ Surokkha   │ 🟢 Active  │ [👁][✏️][🗑]   │
│ Jane S   │ Ops    │ Seba       │ 🟢 Active  │ [👁][✏️][🗑]   │
│ Bob J    │ HR     │ Verosa     │ 🔴 Inactive│ [👁][✏️][🗑]   │
└──────────┴────────┴────────────┴─────────────┴────────────────┘
```

**Add Employee Flow (3-Section Modal):**
```
Section 1 — Personal Information:
  First Name*, Last Name*, National ID*, DOB*, Gender*
  Contact Number, Emergency Contact, Address

Section 2 — Employment Details:
  Employee Code (auto or manual), Department* (own org depts only)
  Designation, Join Date*

Section 3 — Insurance Information:
  Plan* (Seba / Surokkha / Verosa — from catalog)
  Sum Insured, Policy Start Date*, End Date*
  Beneficiary Name, Beneficiary Relationship

[Create Employee] → POST /api/employees
  body: all fields + businessId: ownOrgId
Toast: "Employee created successfully"
```

### 11.7 Purchase Orders (B2B Admin)

```
/purchase-orders page

PO Table:
┌───────────────────────────────────────────────────────────────┐
│ Purchase Orders                              [+ Create PO]    │
├──────────────────────┬──────────┬──────────┬─────────────────┤
│ PO Number            │ Plan     │ Premium  │ Status          │
├──────────────────────┼──────────┼──────────┼─────────────────┤
│ PO-20260305-A1B2C3D4 │ Surokkha │ ৳45,000  │ 🟢 Active       │
│ PO-20260201-B2C3D4E5 │ Seba     │ ৳12,500  │ 🟢 Active       │
└──────────────────────┴──────────┴──────────┴─────────────────┘

Create PO Flow:
  1. Click [+ Create PO]
  2. Modal:
     - Organisation: "Acme Corporation" (read-only)
     - Department: dropdown of own depts
     - Plan: catalog dropdown (Seba | Surokkha | Verosa)
     - No. of Employees: number input
     - Start Date / End Date
     Auto-calculated:
       PO Number: PO-YYYYMMDD-XXXXXXXX
       Premium: plan_rate × employee_count
  3. [Create] → POST /api/purchase-orders
  4. Toast: "Purchase Order created"
```

---

## 12. Role Mapping & Event-Driven Assignment

### 12.1 UserType to Portal Role Mapping

```typescript
// Executed at login in /api/auth/login/route.ts
switch (userType) {
  case UserType.SYSTEM_USER:
    portalRole  = 'SYSTEM_ADMIN'
    businessId  = ''              // Super admin — no org lock
    break

  case UserType.B2B_ORG_ADMIN:
    portalRole  = 'B2B_ORG_ADMIN'
    businessId  = await resolveMyOrganisation(userId)  // call /organisations/me
    break

  default:
    // Resolve from AuthZ — check what role user has in their org
    const orgRole = await authzService.getUserRoleInDomain(userId)
    portalRole  = orgRole   // BUSINESS_ADMIN | HR_MANAGER | VIEWER
    businessId  = orgRole.orgId
}
```

### 12.2 Casbin Domain Lifecycle

```
Organisation created (org_id = "abc-123"):
  → Event: org.created
  → AuthZ creates domain: b2b:abc-123
  → Seeds default policies for B2B_ORG_ADMIN, BUSINESS_ADMIN, HR_MANAGER, VIEWER

Admin assigned (user_id = "usr-456", org_id = "abc-123"):
  → Event: org.admin.assigned
  → AuthZ: g(usr-456, B2B_ORG_ADMIN, b2b:abc-123)
  → User can now call gRPC with x-portal=PORTAL_B2B, x-business-id=abc-123

Member added (user_id = "usr-789", role = BUSINESS_ADMIN):
  → Event: org.member.added
  → AuthZ: g(usr-789, BUSINESS_ADMIN, b2b:abc-123)

Member removed (user_id = "usr-789"):
  → Event: org.member.removed
  → AuthZ: delete g(usr-789, *, b2b:abc-123)

Organisation deleted:
  → Event: org.deleted
  → AuthZ: delete all policies where dom = b2b:abc-123
  → All users lose access immediately
```

### 12.3 Multi-Org Considerations

- A user can belong to MULTIPLE organisations (multiple rows in org_members)
- Their portal_biz_id cookie is set to the PRIMARY org at login
- To switch org context: must re-login or implement org switcher (future feature)
- AuthZ checks the SPECIFIC domain passed via x-business-id each call

---

## 13. Database Schema

### 13.1 b2b_schema Tables

```sql
-- Organisations
CREATE TABLE b2b_schema.organisations (
    organisation_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name              VARCHAR(200) NOT NULL,
    code              VARCHAR(50)  NOT NULL UNIQUE,
    status            VARCHAR(20)  NOT NULL DEFAULT 'PENDING',  -- PENDING|ACTIVE|INACTIVE|SUSPENDED
    total_employees   INT          NOT NULL DEFAULT 0,
    contact_email     VARCHAR(200),
    contact_phone     VARCHAR(20),
    address_line1     VARCHAR(300),
    address_line2     VARCHAR(300),
    city              VARCHAR(100),
    state             VARCHAR(100),
    zip_code          VARCHAR(20),
    country           VARCHAR(100),
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ                              -- soft delete
);

-- Org Members
CREATE TABLE b2b_schema.org_members (
    member_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organisation_id   UUID NOT NULL REFERENCES b2b_schema.organisations(organisation_id),
    user_id           UUID NOT NULL,                          -- FK to authn_schema.users
    role              VARCHAR(50) NOT NULL,                   -- B2B_ORG_ADMIN|BUSINESS_ADMIN|HR_MANAGER|VIEWER
    status            VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    joined_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE(organisation_id, user_id)
);

-- Departments
CREATE TABLE b2b_schema.departments (
    department_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name              VARCHAR(200) NOT NULL,
    business_id       UUID NOT NULL REFERENCES b2b_schema.organisations(organisation_id),
    employee_no       INT  NOT NULL DEFAULT 0,                -- computed / updated via trigger
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);

-- Employees
CREATE TABLE b2b_schema.employees (
    employee_uuid     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_code     VARCHAR(50),
    first_name        VARCHAR(100) NOT NULL,
    last_name         VARCHAR(100) NOT NULL,
    national_id       VARCHAR(20)  NOT NULL,
    date_of_birth     DATE,
    gender            VARCHAR(20),
    contact_number    VARCHAR(20),
    emergency_contact VARCHAR(20),
    address           TEXT,
    department_id     UUID REFERENCES b2b_schema.departments(department_id),
    business_id       UUID NOT NULL REFERENCES b2b_schema.organisations(organisation_id),
    designation       VARCHAR(100),
    join_date         DATE,
    plan_id           UUID,
    sum_insured       NUMERIC(12,2),
    policy_start_date DATE,
    policy_end_date   DATE,
    beneficiary_name  VARCHAR(200),
    beneficiary_rel   VARCHAR(100),
    status            VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ,
    UNIQUE(national_id, business_id)
);

-- Purchase Orders
CREATE TABLE b2b_schema.purchase_orders (
    purchase_order_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    po_number         VARCHAR(50)  NOT NULL UNIQUE,           -- PO-YYYYMMDD-XXXXXXXX
    department_id     UUID REFERENCES b2b_schema.departments(department_id),
    business_id       UUID NOT NULL REFERENCES b2b_schema.organisations(organisation_id),
    plan_id           UUID,
    plan_name         VARCHAR(200),
    employee_count    INT  NOT NULL DEFAULT 0,
    premium_amount    NUMERIC(14,2),
    policy_start_date DATE,
    policy_end_date   DATE,
    status            VARCHAR(20)  NOT NULL DEFAULT 'ACTIVE',
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
```

### 13.2 authz_schema (Casbin)

```sql
CREATE TABLE authz_schema.casbin_rule (
    id     BIGSERIAL PRIMARY KEY,
    ptype  VARCHAR(10),   -- 'p' (policy) or 'g' (group/role)
    v0     VARCHAR(256),  -- sub or user
    v1     VARCHAR(256),  -- dom
    v2     VARCHAR(256),  -- obj
    v3     VARCHAR(256),  -- act
    v4     VARCHAR(256),  -- eft (allow/deny)
    v5     VARCHAR(256)
);

-- Example seeded rows:
-- ptype | v0              | v1          | v2            | v3                      | v4
-- p     | SYSTEM_ADMIN    | system:root | svc:b2b/*     | (GET|POST|PATCH|DELETE) | allow
-- p     | B2B_ORG_ADMIN   | b2b:*       | svc:b2b/org*  | (GET|PATCH)             | allow
-- g     | usr-456         | B2B_ORG_ADMIN | b2b:abc-123 |                         |
```

---

## 14. Changes & Bug Fixes Log

Comprehensive documentation of all recent changes and bug fixes across the InsureTech B2B Portal codebase.

### 14.1 Backend Fix: Nil Publisher Panic

**File:** `backend/inscore/microservices/b2b/cmd/server/main.go`

**Root Cause:**
When Kafka was unavailable, `var publisher *events.Publisher` was typed-nil. Assigning this to the `EventPublisher` interface created a non-nil interface with a nil receiver. Calling methods like `AddOrgMember()` or `PublishOrgMemberAdded()` on this interface would panic because the receiver was nil.

**Symptoms:**
- B2B admin creation failed silently or crashed when Kafka was down
- System relied on Kafka being always available for basic org operations

**Fix:**
Always call `publisher := events.NewPublisher(kafkaProducer)`. The publisher handles nil producer gracefully internally:
- If Kafka is unavailable, the producer is nil
- Publisher methods check for nil producer and return early (graceful degradation)
- B2B admin creation now works without Kafka

**Impact:**
- B2B admin creation is now resilient to Kafka failures
- Org operations continue even when Kafka is temporarily down
- Events are simply not published in offline scenarios

---

### 14.2 BFF Route Fix: /api/organisations/[id]/admins

**File:** `app/api/organisations/[id]/admins/route.ts`

**Root Cause:**
Route was not handling both modes of admin assignment. Additionally, the body format sent from frontend used snake_case (`mobile_number`, `firstName`, `lastName`) but the backend's `assignOrgAdminPayload` struct defined json tags in camelCase (`mobileNumber`, `fullName`), causing fields to arrive empty in the backend.

**Symptoms:**
- Creating new admin with email/password failed (400 error)
- Fields appeared empty in backend request
- Only one mode of admin assignment was supported

**Fix:**
Now handles TWO modes based on body shape:
- **Mode A:** `{ email, password, mobileNumber, fullName }` → creates new user + assigns B2B_ORG_ADMIN role (POST /v1/b2b/organisations/{id}/admins)
- **Mode B:** `{ memberId }` only → promotes existing member via SDK assignOrgAdmin

Body now sends exactly matching backend camelCase json tags:
```json
{ email, password, mobileNumber, fullName }
```

**Impact:**
- Both admin creation and promotion workflows now work
- No more 400 errors due to field mapping
- Clean separation of two distinct flows

---

### 14.3 BFF Route Fix: /api/organisations/[id]/assign-admin

**File:** `app/api/organisations/[id]/assign-admin/route.ts`

**Root Cause:**
`organisation-client.ts assignAdmin()` was incorrectly posting `{ memberId }` to `/admins` (the create-new-user endpoint) instead of the correct `/assign-admin` endpoint.

**Symptoms:**
- Promoting existing members to admin failed
- Wrong endpoint called

**Fix:**
Split into correct endpoints:
- **Create new admin:** POST `/api/organisations/{id}/admins` with full user data
- **Promote existing:** POST `/api/organisations/{id}/assign-admin` with memberId only

**Impact:**
- Member promotion now works correctly
- Clear API contract separation

---

### 14.4 BFF Route Fix: /api/organisations/[id]/approve (NEW)

**File:** `app/api/organisations/[id]/approve/route.ts` (NEW)

**Implementation:**
- POST endpoint approves pending org
- PATCHes status to ORGANISATION_STATUS_ACTIVE
- Super Admin only (enforced by backend AuthZ)
- Returns: `{ ok, message, organisation }`

**Integration Points:**
- Called from organisations table approve button
- Called from OrgDetailPanel info tab approve button
- Updates status badge immediately if detail panel is open

---

### 14.5 Frontend Fix: OrgDetailPanel (org-detail-panel.tsx)

**File:** `components/organisations/org-detail-panel.tsx`

**Root Cause:**
Edit mode used uncontrolled inputs with `defaultValue`. When switching between orgs in the detail panel, stale data persisted because React didn't re-render controlled fields.

**Symptoms:**
- Data from previous org visible when opening different org
- Editing wrong org's data

**Fix:**
1. Changed edit fields from uncontrolled `defaultValue` → controlled `value`
2. Added useEffect that syncs all edit state when `org.id` or `open` changes
3. Added missing React import for useEffect

**Impact:**
- No more stale data when switching orgs
- Clean state sync on panel open/close

---

### 14.6 Frontend Fix: Organisation Columns (data-table/columns.tsx)

**File:** `components/dashboard/organisations/data-table/columns.tsx`

**Root Cause:**
Used `<details>` dropdown menu for actions, was inefficient and inconsistent with other tables.

**Fix:**
- Replaced `<details>` dropdown with inline icon buttons: 👁 View · ✏️ Edit · ✅ Approve · 🗑 Delete
- `onView` prop wired: `buildOrganisationColumns(onRefresh, onRowClick, onApprove)`
- Approve button: only shown when status contains PENDING, calls `organisationClient.approve()`
- Uses `LuCheck` icon (note: `LuCheckCircle` does not exist in react-icons/lu)

**Impact:**
- Consistent UI across all tables
- Better mobile UX with inline buttons
- Faster access to actions

---

### 14.7 Frontend Fix: Department Columns (data-table/columns.tsx)

**File:** `components/dashboard/departments/data-table/columns.tsx`

**Fix:**
- Replaced `<details>` dropdown with inline ✏️ Edit · 🗑 Delete icon buttons

---

### 14.8 Frontend Fix: Employee Columns (data-table/columns.tsx)

**File:** `components/dashboard/employees/data-table/columns.tsx`

**Fix:**
- View button now opens AddEmployeeModal with `employeeUuid` (fetches full record via GET /api/employees/{uuid})
- Replaced `<details>` dropdown with inline 👁 View · ✏️ Edit · 🗑 Delete buttons

---

### 14.9 Frontend Fix: OrgMemberPanel (org-member-panel.tsx)

**File:** `components/organisations/org-member-panel.tsx`

**Fixes:**

1. **UserIdCell:** Shows truncated UUID (a1b2c3d4…) with one-click copy button
2. **Role badges:** Colour-coded:
   - Purple = Admin (B2B_ORG_ADMIN)
   - Blue = HR Manager
   - Grey = Viewer
3. **Promote button** (🛡️, Super Admin only):
   - Calls `assignAdmin(orgId, memberId)` via correct `/assign-admin` endpoint
   - Previously was calling wrong endpoint
4. **Error handling:**
   - Errors shown as inline red banner, not `alert()`
5. **Code cleanup:**
   - Removed duplicate assign-existing-user form
   - Consolidated into promote icon

**Impact:**
- Better UX with inline errors
- Correct API endpoint called
- Cleaner component code

---

### 14.10 Frontend Fix: useEmployeeForm (src/hooks/useEmployeeForm.ts)

**File:** `src/hooks/useEmployeeForm.ts`

**Root Cause:**
In edit mode, hook only used sparse initialValues from table row (only name + employeeId). Missing fields like email, mobile, gender, dates, and insurance info were not populated, causing empty form fields in edit modal.

**Fix:**
1. **In edit mode:** Now fetches full employee record from `GET /api/employees/{uuid}` on mount
2. **mapApiToForm(EmployeeFullRecord)** maps ALL fields:
   - email, mobileNumber, gender
   - dateOfBirth, dateOfJoining
   - departmentId, businessId
   - insuranceCategory, coverageAmount, assignedPlanId
3. **loadingRecord state** exposed → modal shows spinner while fetching
4. **Previous behavior:** Only used sparse initialValues from table row

**Impact:**
- Edit mode now shows all employee data
- No more missing form fields
- Better UX with loading state

---

### 14.11 API Fix: GET /api/employees/[id]

**File:** `app/api/employees/[id]/route.ts`

**Root Cause:**
Endpoint only returned list-view fields (no email, mobile, etc.). When edit mode tried to populate form, critical fields were missing.

**Fix:**
1. Added `mapViewFull()` function
2. Returns ALL form fields: email, mobile_number, gender, date_of_birth, date_of_joining, department_id, business_id, insurance_category, assigned_plan_id, coverage_amount
3. **insuranceCategory mapping:** proto enum string → numeric form value:
   - 1 = Health
   - 2 = Life
   - 3 = Auto
   - 4 = Travel

**Impact:**
- Edit modal can now pre-fill complete employee record
- Insurance category properly converted for form display

---

### 14.12 Frontend Fix: AddEmployeeModal (components/modals/add-employee-modal.tsx)

**File:** `components/modals/add-employee-modal.tsx`

**Root Cause:**
Two issues:
1. `useEmployeeForm` hook declared AFTER departments useEffect, causing TypeScript TS2448 error
2. Department fetch didn't resolve org correctly in edit mode

**Fixes:**
1. **Hook order:** useEmployeeForm declared BEFORE departments useEffect
2. **Department fetch:** Now uses `organisationId || values.businessId` → resolves org in edit mode
3. **Dependency array:** Added `values.businessId` → depts reload when edit record loads

**Impact:**
- No more TypeScript errors
- Departments properly populate based on org context
- Edit mode works correctly

---

### 14.13 Frontend Fix: Departments Page (Departments.tsx)

**File:** `components/dashboard/departments/Departments.tsx`

**Root Cause:**
Department management lacked the same org-lock pattern as employees page. B2B admins couldn't filter properly by their own org.

**Fix:**
Rebuilt with same org-lock pattern as employees-table:
- **Super Admin:** Select dropdown of all orgs → filter departments
- **B2B Admin:** Readonly label showing org name, auto-locked
- `departmentClient.list(50, 0, selectedOrgId)` — always scoped by org
- Passes `organisationId` to AddDepartmentModal and DataTable

**Impact:**
- B2B admins now see only their org's departments
- Consistent scoping across all pages
- No cross-org data leakage

---

### 14.14 Frontend Fix: AddDepartmentModal

**File:** `components/modals/add-department-modal.tsx`

**Root Cause:**
POST body only sent `{ name }`, but backend's department creation endpoint required both name and businessId.

**Fix:**
1. Added `organisationId` prop
2. POST body now sends `{ name, businessId }` — both required
3. Backend no longer rejects due to missing businessId

**Impact:**
- Department creation now works correctly
- Clear API contract

---

### 14.15 API Fix: catalog/route.ts

**File:** `app/api/purchase-orders/catalog/route.ts`

**Root Cause:**
Route was calling `makeSdkClient(request)` with NO headers argument, causing authentication to fail. Headers are needed to validate the user making the request.

**Fix:**
Added `resolvePortalHeaders()`:
```typescript
const headers = await resolvePortalHeaders(request)
const sdk = makeSdkClient(request, headers)  // was: makeSdkClient(request)
```

**Impact:**
- Catalog endpoint now properly authenticated
- Users can fetch plan catalog successfully

---

### 14.16 Frontend Fix: Purchase Orders Page (purchase-orders.tsx)

**File:** `components/dashboard/purchase-orders/purchase-orders.tsx`

**Fixes:**
1. **Org context:** Resolves on mount via `organisationClient.getMe()`
2. **Department filtering:** `departmentClient.list(200, 0, resolvedOrgId)` — B2B Admin sees only their depts
3. **Error handling:** `handleCreatePurchaseOrder` returns `Promise<boolean>` — modal stays open on error

**Impact:**
- Proper org scoping for PO creation
- B2B admins see only relevant departments
- Better error UX

---

### 14.17 Frontend Fix: AddPurchaseOrderModal — Complete Rewrite

**File:** `components/modals/add-purchase-order-modal.tsx`

**Root Cause:**
Modal had poor UX with sr-only labels, unclear field requirements, and validation did not prevent submission on error.

**Complete Rewrite:**

1. **Labels:** All changed from sr-only → visible block labels
2. **Insurance Category:** Auto-derived from selected plan via useEffect, shown as readonly field
3. **Department dropdown:** Visible label + error message
4. **Product Plan dropdown:** Visible label + error message + plan summary card
5. **Number of Employees / Dependents:** Visible labels + 2-column grid
6. **Coverage Amount (BDT):** Visible label + required validation
7. **Notes:** Changed from Input to textarea (3 rows)
8. **Validation:** Inline field errors, modal stays open on failure
9. **Cancel button:** Added
10. **Type update:** `PurchaseOrderFormInput` now includes `insuranceCategory: string`
11. **Error handling:** `handleCreatePurchaseOrder` returns `Promise<boolean>` → modal stays open on error

**Impact:**
- Clearer UX for form completion
- Better error feedback
- Insurance category no longer requires manual selection

---

### 14.18 Frontend Fix: PO Columns — Actions Implemented

**File:** `components/dashboard/purchase-orders/data-table/columns.tsx`

**Root Cause:**
Was showing hardcoded `—` stub with TODO comment saying "backend returns 501". Backend routes were fully working, just not implemented in frontend.

**Fix:**
Backend routes fully working, actions implemented:
- 👁 **View** → PODetailSheet right side drawer (all PO fields)
- 🗑 **Delete** → confirm dialog → DELETE `/api/purchase-orders/{id}` → refresh

**Impact:**
- Full PO management workflow now functional
- Users can view and delete purchase orders

---

### 14.19 New API Routes: /api/dashboard/stats (NEW)

**File:** `app/api/dashboard/stats/route.ts` (NEW)

**Implementation:**

Fetches KPI stats with parallel `Promise.allSettled`:

**Super Admin returns:**
- totalOrganisations
- pendingOrganisations
- totalEmployees
- totalDepartments
- activePurchaseOrders

**B2B Admin returns:**
- totalMembers
- totalDepartments
- totalEmployees
- activePurchaseOrders

**Response:** `{ ok, stats, role }`

---

### 14.20 New API Routes: /api/dashboard/activity (NEW)

**File:** `app/api/dashboard/activity/route.ts` (NEW)

**Implementation:**

Fetches last 5 from each source, merges, sorts by createdAt DESC, returns top 10.

**Super Admin includes:**
- Recent orgs
- Recent employees
- Recent departments
- Recent purchase orders

**B2B Admin includes:**
- Recent employees
- Recent departments
- Recent purchase orders
- Recent member joins

**Features:**
- `timeAgo()` humanised timestamps (e.g., "2 hours ago")
- Type-based icons

**Response:** `{ ok, activities: [{id, type, title, subtitle, createdAt}] }`

---

### 14.21 Frontend Fix: StatsCards — Live Data

**File:** `components/dashboard/stats-cards/stats-cards.tsx`

**Root Cause:**
Component imported static mock from `lib/stats-cards.ts`. Dashboard always showed hardcoded numbers, not real data.

**Fix:**
1. **Fetches from:** `/api/dashboard/stats`, role-aware cards
2. **Loading state:** Skeleton (4 spinners)
3. **Error state:** Error message display
4. **Props:** `{ title, value, icon (svg path), bgColor, bgIcon (svg path) }`

**Impact:**
- Dashboard now shows real statistics
- Live data updates on page load
- Proper loading/error states

---

### 14.22 Frontend Fix: OverviewActivity — Live Data

**File:** `components/dashboard/overview-activity/overview-activity.tsx`

**Root Cause:**
Component imported static mock from `lib/overview-activity`. Activity feed always showed sample data, not real recent events.

**Fix:**
1. **Fetches from:** `/api/dashboard/activity`
2. **Type-based icons:**
   - 🏢 Org
   - 👥 Member
   - 👤 Employee
   - 📋 Department
   - 🗒 Purchase Order
3. **Timestamps:** `timeAgo()` relative format
4. **States:** Empty, error, loading handled

**Impact:**
- Activity feed now shows real recent events
- Users see actual system activity
- Proper state handling

---

### 14.23 Frontend Fix: Organisations.tsx — Approve Wiring

**File:** `components/dashboard/organisations/Organisations.tsx`

**Fix:**
1. `buildOrganisationColumns` now called with 3rd arg `onApprove`
2. On approve: if org is open in detail panel, immediately updates status badge to ACTIVE
3. Removes need for manual panel refresh

**Impact:**
- Seamless UX when approving orgs
- Status updates in real-time

---

## 15. Error Handling & Security Model

### 14.1 HTTP Error Codes & UI Responses

| HTTP Code | Scenario | UI Behaviour |
|-----------|----------|-------------|
| 400 | Validation error | Inline field errors on form |
| 401 | No session / expired | Redirect to /login?next={currentPath} |
| 403 | Insufficient permission | Toast: "You don't have permission to perform this action" |
| 404 | Resource not found | Toast: "Resource not found" / 404 page |
| 409 | Conflict (duplicate code, active employees) | Toast: specific conflict message |
| 429 | Rate limited / locked | Toast: "Too many requests. Try again later." |
| 500 | Backend error | Toast: "Something went wrong. Please try again." |

### 14.2 Security Guardrails

```
1. httpOnly cookies — session_token never accessible to JS
2. CSRF token — sent in X-CSRF-Token header for all mutations
3. Backend org isolation — x-business-id enforced server-side
   Even if client sends wrong org_id, AuthZ intercept blocks it
4. Casbin domain isolation — every policy scoped to a domain
   Cross-org access impossible even with valid JWT
5. Soft deletes — data retained for audit, never truly deleted
6. Mobile normalization — prevents duplicate accounts via format variants
7. Edge middleware — unauthenticated requests never reach Next.js route handlers
8. Token validation — EVERY API call validates session_token against AuthN
   No local JWT verification — prevents token confusion attacks
```

### 14.3 Role Boundary Enforcement (Dual-Layer)

```
Layer 1 — Frontend (UX):
  - "Organisations" menu hidden from B2B Admin
  - Org dropdown locked for B2B Admin (shows label, not dropdown)
  - "Create Org" button hidden from B2B Admin
  - "Assign B2B Admin role" option hidden from B2B Admin
  - "Approve", "Delete" org actions hidden from B2B Admin

Layer 2 — Backend (Casbin + Service):
  - AuthZ intercept checks EVERY gRPC method
  - Service layer re-validates org ownership via x-business-id
  - Database queries always include business_id filter
  → Even if frontend is bypassed, backend enforces all boundaries

Both layers MUST agree. Frontend is UX. Backend is the real security boundary.
```

### 14.4 Audit Trail

```
Every state-changing operation writes to audit log:
  - Who: user_id + role
  - What: action (create/update/delete/approve)
  - Which: resource type + resource_id
  - When: timestamp
  - Context: organisation_id, IP address, user_agent

Kafka events serve dual purpose:
  1. AuthZ synchronization (primary)
  2. Audit log population (secondary)
```

---

## Quick Reference — API + Permission Cheat Sheet

| Action | API Endpoint | Super Admin | B2B Admin | HR Manager | Viewer |
|--------|-------------|-------------|-----------|------------|--------|
| List all orgs | GET /api/organisations | YES | NO | NO | NO |
| Create org | POST /api/organisations | YES | NO | NO | NO |
| View own org | GET /api/organisations/me | YES | YES | YES | YES |
| Update org | PATCH /api/organisations/{id} | YES | Own only | NO | NO |
| Approve org | POST /api/organisations/{id}/approve | YES | NO | NO | NO |
| Delete org | DELETE /api/organisations/{id} | YES | NO | NO | NO |
| List members | GET /api/organisations/{id}/members | YES | Own only | NO | NO |
| Add member | POST /api/organisations/{id}/members | YES | Own only | NO | NO |
| Remove member | DELETE /api/organisations/{id}/members/{mid} | YES | Own only | NO | NO |
| Bootstrap admin | POST /api/organisations/{id}/admins | YES | NO | NO | NO |
| List departments | GET /api/departments | YES | Own only | Own only | Own only |
| Create department | POST /api/departments | YES | Own only | YES | NO |
| Update department | PATCH /api/departments/{id} | YES | Own only | YES | NO |
| Delete department | DELETE /api/departments/{id} | YES | Own only | NO | NO |
| List employees | GET /api/employees | YES | Own only | Own only | Own only |
| Create employee | POST /api/employees | YES | Own only | YES | NO |
| Update employee | PATCH /api/employees/{id} | YES | Own only | YES | NO |
| Delete employee | DELETE /api/employees/{id} | YES | Own only | NO | NO |
| List POs | GET /api/purchase-orders | YES | Own only | Own only | Own only |
| Create PO | POST /api/purchase-orders | YES | YES | YES | NO |

---

*End of B2B Portal Complete Org Management Reference*
*Generated: March 2026 | InsureTech Platform*
