# B2B Portal (Frontend) — Central Reference

> **Single source of truth** for the B2B Portal Next.js frontend: architecture, authentication, SDK/header injection, API routes, route guards, state management, org scoping, form validation, and UI components.

---

## 1. Architecture

### Stack
Next.js (App Router) · TypeScript · Radix UI · Tailwind CSS · Server-side API routes as BFF (Backend-For-Frontend) · `@lifeplus/insuretech-sdk` (auto-generated from proto)

### Folder Structure

```
b2b_portal/
├── app/
│   ├── api/
│   │   ├── auth/login|logout|session      # Server-side auth routes
│   │   ├── organisations/                  # Org CRUD + members + admins
│   │   │   ├── route.ts                   # GET list, POST create
│   │   │   ├── [id]/route.ts              # GET, PATCH, DELETE, POST approve
│   │   │   ├── [id]/admins/route.ts       # POST create admin
│   │   │   ├── [id]/members/route.ts      # GET list members
│   │   │   ├── [id]/members/[memberId]/   # DELETE remove member
│   │   │   └── me/route.ts               # GET current org context
│   │   ├── departments/                    # Dept CRUD
│   │   │   ├── route.ts                   # GET list, POST create
│   │   │   └── [id]/route.ts              # GET, PATCH, DELETE
│   │   ├── employees/                      # Employee CRUD
│   │   │   ├── route.ts                   # GET list, POST create
│   │   │   └── [id]/route.ts              # GET, PATCH, DELETE
│   │   └── purchase-orders/
│   │       ├── route.ts, [id]/, catalog/
│   ├── organisations|departments|employees|purchase-orders/  # Pages
│   ├── login/page.tsx
│   └── layout.tsx, page.tsx (dashboard), globals.css
├── components/
│   ├── auth/login-form.tsx
│   ├── dashboard/
│   │   ├── organisations/Organisations.tsx + data-table/
│   │   ├── departments/Departments.tsx + data-table/
│   │   ├── employees/employees-table.tsx + data-table/
│   │   └── overview-activity/, stats-cards/
│   ├── organisations/org-detail-panel, org-member-panel
│   ├── modals/add-organisation-modal, add-department-modal, add-employee-modal
│   └── ui/ (Radix UI primitives)
├── src/
│   ├── lib/
│   │   ├── sdk/b2b-sdk-client.ts          # SDK client factory (makeSdkClient)
│   │   ├── sdk/session-headers.ts          # resolvePortalHeaders
│   │   ├── clients/
│   │   │   ├── auth-client.ts, organisation-client.ts
│   │   │   ├── department-client.ts, employee-client.ts
│   │   │   ├── api-client.ts (shared parseJson, ApiResult)
│   │   │   └── b2b-dashboard-client.ts
│   │   ├── auth/backend-auth.ts, session.ts, session-store.ts
│   │   ├── types/auth.ts, b2b.ts, employee-form.ts
│   │   └── proto-generated/ (protobuf types)
│   └── hooks/
│       ├── useCrudList.ts                  # Generic CRUD list hook
│       ├── useOrganisationForm.ts          # Org form validation
│       ├── useEmployeeForm.ts              # Employee form validation
│       └── useToast.ts
├── middleware.ts                            # Edge middleware (auth guards)
├── next.config.ts, package.json
```

**Key dependencies:** `@bufbuild/protobuf`, `@lifeplus/insuretech-sdk` (local tgz), `@radix-ui/*`, `@tanstack/react-table`, `recharts`, `lucide-react`

---

## 2. Authentication Flow

### Login → Cookie Set → Session

```
1. User submits mobile + password → POST /api/auth/login
2. API route normalizes BD mobile number to E.164 (+880XXXXXXXXXX)
3. Calls backend loginWithMobile() → gets session_token in Set-Cookie
4. Portal sets cookies:
   - session_token   (httpOnly, secure, sameSite=strict, 12h) — real auth
   - csrf_token      (httpOnly, secure, sameSite=lax, 12h) — CSRF protection
   - portal_role     (readable) — UX-level role guard
   - portal_user_id  (readable) — x-user-id header
   - portal_biz_id   (readable) — x-business-id header
```

### Cookie Details

| Cookie | httpOnly | Purpose | Values |
|--------|----------|---------|--------|
| `session_token` | ✅ | Backend auth — real security boundary | Encrypted session token |
| `csrf_token` | ✅ | CSRF protection for state-changing requests | `crypto.randomBytes(16)` |
| `portal_role` | ❌ | Edge middleware role guards (UX only) | `SYSTEM_ADMIN`, `B2B_ORG_ADMIN`, `BUSINESS_ADMIN`, `HR_MANAGER`, `VIEWER` |
| `portal_user_id` | ❌ | Inject `x-user-id` header | User UUID |
| `portal_biz_id` | ❌ | Inject `x-business-id` header | Org UUID (empty for super-admin) |

### Role Mapping (Login)

```typescript
function mapUserTypeToRole(userType):
  UserType.SYSTEM_USER     → "SYSTEM_ADMIN"    → businessId: ""
  UserType.B2B_ORG_ADMIN   → "B2B_ORG_ADMIN"   → businessId: from /organisations/me
  Default                  → "BUSINESS_ADMIN"   → businessId: from /organisations/me
```

### Mobile Number Normalization

Bangladesh E.164 format validation (server + client side):
```
01712345678 → +8801712345678       +880 171-234-5678 → +8801712345678
8801712345678 → +8801712345678     008801712345678 → +8801712345678
Regex: /^880(13|14|15|16|17|18|19)\d{8}$/
```

### Login Error Mapping

Backend gRPC/HTTP errors are mapped to user-friendly messages:
- 401 / "unauthenticated" → "Mobile number or password is incorrect."
- 429 / "locked" → "Account temporarily locked."
- "inactive" / "suspended" → "Account not active."
- 400 / "invalid mobile" → "Invalid mobile number."
- 500+ → "Service temporarily unavailable."

---

## 3. Portal Headers & SDK

### Header Resolution (`resolvePortalHeaders`)

All API routes call `resolvePortalHeaders(request)` which:
1. Requires `session_token` cookie (returns null if missing)
2. Reads `portal_role`, `portal_user_id`, `portal_biz_id` cookies
3. Maps role to portal: `SYSTEM_ADMIN` → `PORTAL_SYSTEM`, all others → `PORTAL_B2B`
4. Returns `{ portal, userId, businessId, tenantId }`

### SDK Client (`makeSdkClient`)

Builds SDK client with all auth headers:

```
cookie:          session_token=<token>; csrf_token=<csrf>
X-CSRF-Token:    <csrf_token> (state-changing requests)
x-portal:        PORTAL_SYSTEM | PORTAL_B2B
x-business-id:   <org_id> (empty for super-admin)
x-user-id:       <user_id>
x-tenant-id:     <tenant_id> (default: 00000000-0000-0000-0000-000000000001)
```

### Headers per Role

| Role | x-portal | x-business-id | Casbin Domain |
|------|----------|---------------|---------------|
| `SYSTEM_ADMIN` | `PORTAL_SYSTEM` | _(empty)_ | `system:root` |
| `B2B_ORG_ADMIN` | `PORTAL_B2B` | org UUID | `b2b:{orgId}` |
| `BUSINESS_ADMIN` | `PORTAL_B2B` | org UUID | `b2b:{orgId}` |
| `HR_MANAGER` / `VIEWER` | `PORTAL_B2B` | org UUID | `b2b:{orgId}` |

---

## 4. API Routes

### Auth (`/api/auth`)

| Method | Endpoint | Purpose |
|--------|----------|---------|
| POST | `/api/auth/login` | Normalize BD mobile to E.164, validate credentials, set 5 cookies (session_token, csrf_token, portal_role, portal_user_id, portal_biz_id) |
| POST | `/api/auth/logout` | Clear all cookies, call backend AuthN.Logout |
| GET | `/api/auth/session` | Validate session_token, return principal object |

### Organisations (`/api/organisations`)

| Method | Endpoint | Details |
|--------|----------|---------|
| GET | `/api/organisations` | Super Admin: list all organisations |
| POST | `/api/organisations` | Super Admin: create org (optional: include admin bootstrap object) |
| GET | `/api/organisations/{id}` | Get single organisation |
| PATCH | `/api/organisations/{id}` | Update org fields |
| DELETE | `/api/organisations/{id}` | Soft delete organisation |
| POST | `/api/organisations/{id}/approve` | **NEW:** Approve pending org → ACTIVE (Super Admin only) |
| GET | `/api/organisations/me` | Resolve current user's org context (SYSTEM_ADMIN → null; B2B roles → org object) |

### Org Members & Admins

| Method | Endpoint | Details |
|--------|----------|---------|
| GET | `/api/organisations/{id}/members` | List org members |
| POST | `/api/organisations/{id}/members` | Add member with role |
| DELETE | `/api/organisations/{id}/members/{memberId}` | Remove member |
| POST | `/api/organisations/{id}/admins` | **TWO MODES:** Mode A: body `{email, password, mobileNumber, fullName}` creates new B2B admin. Mode B: body `{memberId}` promotes existing member. ⚠️ Backend JSON tags: camelCase (mobileNumber, fullName) NOT snake_case |
| POST | `/api/organisations/{id}/assign-admin` | Promote existing member to B2B_ORG_ADMIN |

### Departments (`/api/departments`)

| Method | Endpoint | Details |
|--------|----------|-------|
| GET | `/api/departments` | Scoped by `business_id` query param (super-admin) or from session (B2B admin); falls back to `/organisations/me` |
| POST | `/api/departments` | Requires `{name, businessId}` in request body |
| GET | `/api/departments/{id}` | Get single department |
| PATCH | `/api/departments/{id}` | Update department |
| DELETE | `/api/departments/{id}` | Soft delete (backend refuses if active employees exist) |

### Employees (`/api/employees`)

| Method | Endpoint | Details |
|--------|----------|-------|
| GET | `/api/employees?page_size=&business_id=&department_id=` | Filtered by org + optional dept |
| POST | `/api/employees` | Create with businessId context |
| GET | `/api/employees/{id}` | **FULL RECORD:** Returns email, mobileNumber, gender, dateOfBirth, dateOfJoining, departmentId, businessId, insuranceCategory, coverageAmount, assignedPlanId |
| PATCH | `/api/employees/{id}` | Partial update |
| DELETE | `/api/employees/{id}` | Soft delete |

### Purchase Orders (`/api/purchase-orders`)

| Method | Endpoint | Details |
|--------|----------|-------|
| GET | `/api/purchase-orders/catalog` | Requires `resolvePortalHeaders` (was broken without auth) |
| GET | `/api/purchase-orders` | List purchase orders |
| POST | `/api/purchase-orders` | Create purchase order |
| GET | `/api/purchase-orders/{id}` | Get single PO |
| DELETE | `/api/purchase-orders/{id}` | Delete PO |

### Dashboard (`/api/dashboard`)

| Method | Endpoint | Details |
|--------|----------|-------|
| GET | `/api/dashboard/stats` | **NEW:** Role-aware KPI stats (Super Admin vs B2B Admin views) |
| GET | `/api/dashboard/activity` | **NEW:** Recent activity feed (top 10 items, sorted by createdAt DESC) |

---

## 5. Route Guards

### Edge Middleware (`middleware.ts`)

**Public paths:** `/login`, `/api/auth/login`
**Bypassed paths:** `/_next`, `/public`, `/logos`, `/favicon.ico`, `/api/auth/session`, `/api/auth/logout`

**Middleware flow:**
1. Public/static route → allow
2. No `session_token` cookie → redirect to `/login?next={path}`
3. Logged in + on `/login` → redirect to `/`
4. Protected routes → require session

> **Important:** Middleware does NOT validate the token — only checks cookie presence. Real auth is backend-side.

---

## 6. Organisation Scoping (Key Logic)

### The `/organisations/me` Pattern

The critical org-scoping logic differentiates system admin vs org-bound users:

```typescript
// In employees-table.tsx
const meResult = await organisationClient.getMe();

if (meResult.organisation?.id) {
  // USER IS B2B ADMIN → LOCKED to their org
  setResolvedOrg(meResult.organisation);       // Show readonly label
  setOrganisations([meResult.organisation]);   // Only their org
  setSelectedOrgId(meResult.organisation.id);  // Auto-select
} else {
  // USER IS SYSTEM ADMIN → Can see ALL orgs
  const listResult = await organisationClient.list();
  setOrganisations(listResult.organisations);   // Show dropdown
}
```

**UI Behavior:**
- **B2B admin:** Readonly label "This B2B admin is locked to their organisation context."
- **System admin:** Dropdown selector "Employees are loaded by organisation id."

### How `/api/organisations/me` Works

```typescript
// Server-side route handler
if (session?.principal.role === "SYSTEM_ADMIN") {
  return { ok: true, organisation: null };     // System admin → no org lock
}
// Otherwise call backend /v1/b2b/organisations/me → resolve actual org
```

---

## 6.5 POST /api/organisations/{id}/admins — Two Modes

This endpoint supports two distinct workflows for managing B2B admins:

### Mode A: Create New Admin User

**Request body:**
```json
{
  "email": "admin@company.com",
  "password": "SecurePass123!",
  "mobileNumber": "+8801712345678",
  "fullName": "John Admin"
}
```

- Creates a brand-new B2B admin user with these credentials
- Used during org creation (admin bootstrap)
- ⚠️ **IMPORTANT:** Backend JSON tags are **camelCase** (mobileNumber, fullName) NOT snake_case

### Mode B: Promote Existing Member

**Request body:**
```json
{
  "memberId": "uuid-of-existing-member"
}
```

- Promotes an existing organisation member (who already has an account) to B2B_ORG_ADMIN role
- Used to grant admin privileges to existing users

**Backend implementation:** Detects which mode based on request body shape (if memberId exists, use Mode B; else Mode A).

---

## 6.6 Organisation Approval Flow

### Org Status: PENDING → ACTIVE

1. **Approve button visibility:** Only appears in org detail panel when `status === "PENDING"`
2. **User action:** Super Admin clicks "Approve" button
3. **API call:** `POST /api/organisations/{id}/approve` (Super Admin only)
4. **Backend:** Updates org status from PENDING to ACTIVE
5. **Frontend response:** Status badge in detail panel updates immediately (e.g., "PENDING" → "ACTIVE" badge)
6. **UI sync:** Approve button disappears; org can now be used

This workflow allows Super Admins to onboard organisations with a two-step process: create (PENDING) → review → approve (ACTIVE).

---

## 6.7 Employee Edit Form Pattern

### Full Record Fetch on Edit

When `useEmployeeForm` is in edit mode:

1. Fetches `GET /api/employees/{uuid}` before rendering the form
2. Backend returns **mapViewFull()** with ALL fields:
   - Personal: email, mobileNumber, gender, dateOfBirth
   - Employment: dateOfJoining, departmentId, businessId
   - Insurance: insuranceCategory, coverageAmount, assignedPlanId
3. Form populates all fields from the fetched record
4. Department dropdown is reloaded using `businessId` from the fetched employee data

This ensures that edit forms always have the complete employee profile, even if the list view only shows partial data.

---

## 6.8 Purchase Order Form Behavior

### Insurance Category Auto-Derivation

- **insuranceCategory** is automatically derived from the selected plan on the frontend
- **NOT sent to the backend** — the backend derives it from `plan_id`
- Frontend only computes it for UI display

### Form UX Improvements

- **All field labels visible** (were `sr-only` before — now accessible to all users)
- **Inline field error validation** — errors appear below fields during validation
- **Modal state on error:** `onSubmit()` returns `Promise<boolean>`
  - Returns `false` on validation/API error → modal stays open (user can correct and retry)
  - Returns `true` on success → modal closes

---

## 6.9 Dashboard API Routes

### GET /api/dashboard/stats

**Purpose:** Fetch role-aware KPI statistics for the dashboard.

**Role-aware behavior:**
- **SYSTEM_ADMIN:** System-wide KPIs (e.g., total orgs, total employees, total POs, system health)
- **B2B_ORG_ADMIN / BUSINESS_ADMIN / HR_MANAGER / VIEWER:** Organisation-scoped KPIs (e.g., org employees, active POs, pending approvals)

**Frontend component:** `StatsCards`
- Calls `dashboardClient.getStats()`
- Displays loading state while fetching
- Shows error state if API fails
- Renders stats in card grid with icons/trends

### GET /api/dashboard/activity

**Purpose:** Fetch recent activity feed for the dashboard.

**Return format:** List of activity items (top 10), sorted by `createdAt` DESC (newest first).

**Activity types may include:**
- Organisation created/approved/updated
- Employee added/updated/deleted
- Purchase order created/processed
- Department created/updated
- Member added/removed

**Frontend component:** `OverviewActivity`
- Calls `dashboardClient.getActivity()`
- Displays loading state while fetching
- Shows error state if API fails
- Renders recent items in chronological list with timestamps

---

## 7. State Management & Hooks

**Technology:** React hooks + in-memory state (NO Redux or Zustand)

### `useCrudList<T>(fetcher, dataKey)` — Generic CRUD

```typescript
// Used by Organisations, Departments, Employees pages
const { data, loading, error, reload } = useCrudList<Department>(
  () => departmentClient.list(), "departments"
);
// reload() is stable (refs prevent stale closures), safe to pass to children
```

### Session Store (`session-store.ts`)

- In-memory `Map` for server-side sessions
- TTL: 12 hours (43,200,000ms)
- CSRF token: `crypto.randomBytes(16)` per session

---

## 8. Key Components

### Organisation Management
- **`Organisations.tsx`** — Main list page (super-admin only), uses `useCrudList()`
- **`org-detail-panel.tsx`** — Side panel: Info tab, Members tab, Departments tab. Actions: Approve (pending; visible only when status=PENDING), Delete. Status badge updates immediately on approval.
- **`org-member-panel.tsx`** — Member list. Add/remove/assign-admin
- **`add-organisation-modal.tsx`** — Create mode: org fields + admin fields. Edit mode: 2 tabs (Organisation, B2B Admins)

### Employee Management
- **`employees-table.tsx`** — Org dropdown/lock logic (see Section 6), loads employees by selected org
- **`add-employee-modal.tsx`** — 3 sections: Personal Info, Employment, Insurance. Loads departments dynamically
- **Employee edit flow (useEmployeeForm):** In edit mode, fetches `GET /api/employees/{uuid}` which returns full record with all fields. Department dropdown reloads using `businessId` from the fetched record.
  - Fields populated: email, mobileNumber, gender, dateOfBirth, dateOfJoining, departmentId, businessId, insuranceCategory, coverageAmount, assignedPlanId

### Department Management
- **`Departments.tsx`** — List page, uses `useCrudList()` with business_id resolution
- **Delete flow:** Confirmation dialog → `departmentClient.delete(id)` → reload list. Backend refuses if active employees exist

### Purchase Orders
- **Purchase Order Form:** insuranceCategory is auto-derived from selected plan (NOT sent to backend — backend derives from plan_id). All field labels visible (were sr-only before). Inline field error validation. onSubmit returns `Promise<boolean>` — modal stays open on error.

### Dashboard
- **`StatsCards`** — Fetches `GET /api/dashboard/stats`. Role-aware stats: Super Admin vs B2B Admin views. Has loading/error states.
- **`OverviewActivity`** — Fetches `GET /api/dashboard/activity`. Renders real recent items from backend (top 10, sorted by createdAt DESC). Has loading/error states.

### Form Validation
- **Organisation create:** name (required), code (auto-generated if blank), admin fields: email, password (8+ chars, upper, lower, digit, symbol), mobile (BD format)
- **Employee create/edit:** name, employeeId, departmentId, dateOfJoining required. Coverage amount supports currency conversion. Edit mode: all fields pre-populated from full record fetch.
- **Member roles:** `ORG_MEMBER_ROLE_BUSINESS_ADMIN`/`_ADMIN` → "B2B Admin", `_HR_STAFF` → "HR Staff", `_EMPLOYEE` → "Employee"

### Data Types

```typescript
interface Organisation {
  id: string; name: string; code?: string; industry: string;
  contactEmail: string; contactPhone: string; address: string;
  status: "Active"|"Inactive"|"Suspended"|"Pending"; totalEmployees?: number;
}
interface Department { id: string; name: string; employeeNo: number; totalPremium: string; }
interface EmployeeListItem {
  id: string; name: string; employeeID: string; department: string;
  insuranceCategory?: string; assignedPlan?: string; coverage: string;
  premiumAmount: string; status: "Active"|"Inactive"; numberOfDependent: number;
}
```

---

## 9. Browser-Side API Clients

All clients follow the same pattern: `fetch('/api/...')` → `parseJson<T>(response)` with content-type validation.

| Client | Methods |
|--------|---------|
| `authClient` | `login(payload)`, `logout()`, `getSession()` |
| `organisationClient` | `list()`, `get(id)`, `getMe()`, `create()`, `update()`, `delete()`, `listMembers()`, `assignAdmin()`, `createAdmin()`, `removeMember()`, `approve(id)` |
| `departmentClient` | `list(pageSize, offset)`, `create({name, businessId})`, `update(id, {name})`, `delete(id)` |
| `employeeClient` | `list({pageSize, offset, businessId, departmentId, status})`, `get(id)` ← returns full record, `create(payload)`, `update(id, payload)`, `delete(id)` |
| `dashboardClient` | **NEW:** `getStats()` → role-aware KPI stats, `getActivity()` → recent activity feed (top 10, DESC) |
| `purchaseOrderClient` | `list()`, `get(id)`, `create(payload)`, `delete(id)`, `getCatalog()` |

---

## 10. Security Model — Defense-in-Depth

```
Layer 1: httpOnly session cookie (session_token)
  → Only backend can read/validate · Protects from XSS

Layer 2: CSRF protection
  → X-CSRF-Token header on mutations · csrf_token httpOnly cookie · sameSite=strict

Layer 3: Edge middleware UX guards (middleware.ts)
  → Redirects unauthenticated before app logic · Checks cookie presence only

Layer 4: API route auth checks (resolvePortalHeaders)
  → Verify session_token exists · Extract portal metadata · Inject into SDK

Layer 5: Backend Casbin policy enforcement
  → Real security boundary · PORTAL_SYSTEM → system:root · PORTAL_B2B → b2b:{orgId}
```

**Key insight:** The portal is stateless. Session data and authorization decisions live entirely in the backend. The portal reads cookies, forwards them to the gateway, and renders responses. No sensitive data stored in the browser.

**No explicit frontend permission checks** — all CRUD operations are visible to authenticated users. Backend enforces permissions via 403 responses. The org-locking via `getMe()` is UX-only; the backend controls actual access.
