# B2B Portal — Org & User Management Upgrade Plan

## Context Overview (Read This First)

### Backend Stack
- **B2B gRPC service** on port `50112` handles all business logic
- **AuthZ service** (Casbin + Kafka) handles role-based access control
- **Gateway** (REST→gRPC) exposes HTTP endpoints consumed by the portal
- **Three Casbin roles** relevant to portal:
  - `super_admin` → domain `system:root` → full CRUD on all orgs
  - `b2b_org_admin` → domain `b2b:{org_id}` → full CRUD within their org
  - `partner_user` → domain `b2b:{org_id}` → CRUD on departments/employees/POs

### Role Assignment Flow (Kafka Event Chain)
```
super_admin calls AssignOrgAdmin(org_id, user_id)
  → B2BAdminAssigned event published to Kafka
  → authz consumer assigns b2b_org_admin role in b2b:{org_id}
  → ensureScopedRolePolicies copies b2b:root policies to b2b:{org_id}

super_admin calls AddOrgMember(org_id, user_id, role=HR_MANAGER|VIEWER)
  → OrgMemberAdded event published
  → authz consumer assigns partner_user role in b2b:{org_id}
```

### Auth Fix Applied (Previous Session)
- `authz_interceptor.go`: Only `ResolveMyOrganisation` bypasses Casbin; ALL other methods go through Casbin
- `portal_seeder.go`: `b2b_org_admin` and `partner_user` have wildcard `*` + verb policies in `b2b:root`
- `b2b_service.go`: real `x-user-id` used in all published events
- `consumers/handlers.go`: `HR_MANAGER`+`VIEWER` → `partner_user`; `BUSINESS_ADMIN` → `b2b_org_admin`

### SDK
- Package: `@lifeplus/insuretech-sdk` built from `E:\Projects\InsureTech\sdks\insuretech-typescript-sdk`
- Build: `tsup` → CJS + ESM + `.d.ts` → tar → copied to `b2b_portal` as local dep
- Pipeline: `run_api_pipeline.ps1` (proto → OpenAPI → SDK gen → build → tar → install)
- B2B SDK functions: `b2bServiceListOrganisations`, `b2bServiceCreateOrganisation`, `b2bServiceGetOrganisation`, `b2bServiceUpdateOrganisation`, `b2bServiceDeleteOrganisation`, `b2bServiceListOrgMembers`, `b2bServiceAssignOrgAdmin`, `b2bServiceAddOrgMember`, `b2bServiceRemoveOrgMember`, `b2bServiceListDepartments`, `b2bServiceCreateDepartment`, `b2bServiceUpdateDepartment`, `b2bServiceDeleteDepartment`, `b2bServiceListEmployees`, `b2bServiceCreateEmployee`, `b2bServiceUpdateEmployee`, `b2bServiceDeleteEmployee`, `b2bServiceListPurchaseOrders`, `b2bServiceCreatePurchaseOrder`, `b2bServiceListPurchaseOrderCatalog`

### Portal Stack
- **Next.js 16** (App Router) + React 19 + TypeScript
- **Tailwind CSS** + Radix UI + shadcn/ui components
- **TanStack Table** for data grids
- Pattern: Next.js API Routes (`app/api/*/route.ts`) proxy to backend → browser clients (`src/lib/clients/*.ts`) call API routes → React components use `useCrudList` hook
- Auth: mobile OTP → session cookie + CSRF token → `x-business-id` resolved from session

### What Already Exists
| Feature | Files | Status |
|---------|-------|--------|
| Login page | `app/(auth)/login/page.tsx` + `api/auth/login/route.ts` | ✅ Done |
| Session check | `api/auth/session/route.ts` | ✅ Done |
| Org list/create | `api/organisations/route.ts` + `components/dashboard/organisations/Organisations.tsx` | ✅ Done |
| Org get/update/delete | `api/organisations/[id]/route.ts` | ✅ Done |
| Org admin create | `api/organisations/[id]/admins/route.ts` | ✅ Done |
| Org member list | `api/organisations/[id]/members/route.ts` | ✅ Done |
| Employee CRUD | `api/employees/route.ts` + `components/employees/` | ✅ Done |
| Department CRUD | `api/departments/route.ts` + `components/departments/` | ✅ Done |
| Dashboard layout | `app/(dashboard)/layout.tsx` | ✅ Done |
| useCrudList hook | `src/hooks/use-crud-list.ts` | ✅ Done |

### What Is Missing / Broken
| Feature | Gap |
|---------|-----|
| `run_api_pipeline.ps1` SDK tar+copy step | Missing — SDK not auto-installed to portal |
| Assign existing user as org admin | No UI (only create-new-admin exists) |
| Remove org admin role | No API route or UI |
| Add org member (HR_MANAGER/VIEWER role) | No API route or UI |
| Remove org member | No UI for `RemoveOrgMember` |
| HR user management page | Missing — HR managers need their own filtered view |
| Role-based nav filtering | Partial — needs b2b_org_admin vs partner_user split |
| Org approval flow | UI exists in modal but no dedicated approval page |
| Purchase order management | Missing full UI |
| Super admin org management | Full page missing (currently shared with b2b admin) |

---

## Phase 0 — Fix SDK Pipeline (run_api_pipeline.ps1)

### Problem
`run_api_pipeline.ps1` generates and builds the TypeScript SDK but does NOT copy the built tarball into `b2b_portal` and reinstall it as a local dependency. Every time the API changes, the portal uses a stale SDK.

### Steps

**Step 0.1 — Add SDK tar+copy+install block to `run_api_pipeline.ps1`**

Add at the end of the existing pipeline, after `npm run build` in the SDK directory:

```powershell
# ── STEP: Package SDK as tarball and install into b2b_portal ─────────────────
Write-Host "`n[SDK] Packaging TypeScript SDK..." -ForegroundColor Cyan

$sdkDir   = "$PSScriptRoot\sdks\insuretech-typescript-sdk"
$portalDir = "$PSScriptRoot\b2b_portal"
$sdkDepsDir = "$portalDir\sdk"

# 1. Build the SDK (tsup → dist/)
Set-Location $sdkDir
npm run build
if ($LASTEXITCODE -ne 0) { throw "[SDK] Build failed" }

# 2. Pack to tarball (produces lifeplus-insuretech-sdk-*.tgz)
$packOutput = npm pack --json | ConvertFrom-Json
$tarball    = $packOutput[0].filename
Write-Host "[SDK] Packed: $tarball" -ForegroundColor Green

# 3. Copy tarball to portal/sdk/ directory
if (-not (Test-Path $sdkDepsDir)) { New-Item -ItemType Directory -Path $sdkDepsDir | Out-Null }
Copy-Item -Path "$sdkDir\$tarball" -Destination "$sdkDepsDir\insuretech-sdk.tgz" -Force
Write-Host "[SDK] Copied tarball to $sdkDepsDir\insuretech-sdk.tgz" -ForegroundColor Green

# 4. Update package.json dep to point to local tarball
$pkgJson = Get-Content "$portalDir\package.json" | ConvertFrom-Json
$pkgJson.dependencies.'@lifeplus/insuretech-sdk' = "file:./sdk/insuretech-sdk.tgz"
$pkgJson | ConvertTo-Json -Depth 10 | Set-Content "$portalDir\package.json"

# 5. Reinstall in portal
Set-Location $portalDir
npm install
if ($LASTEXITCODE -ne 0) { throw "[SDK] Portal npm install failed" }
Write-Host "[SDK] SDK installed in portal successfully." -ForegroundColor Green
Set-Location $PSScriptRoot
```

**File to edit**: `E:\Projects\InsureTech\run_api_pipeline.ps1`

**Portal `sdk/` directory**: `E:\Projects\InsureTech\b2b_portal\sdk\` (gitignored)

**Add to `.gitignore`**:
```
b2b_portal/sdk/
```

---

## Phase 1 — New API Routes (Backend Proxy Layer)

All routes follow the same pattern:
1. Extract session from cookie → get `x-user-id`, `x-business-id`, `x-portal`, `x-tenant-id`
2. Call SDK function with auth headers
3. Return JSON response

### 1.1 — `POST /api/organisations/[id]/members` (Add org member)

**File**: `app/api/organisations/[id]/members/route.ts`

Currently only `GET` exists. Add `POST` to call `b2bServiceAddOrgMember`.

```typescript
// POST: Add a member to the org
export async function POST(req: Request, { params }: { params: { id: string } }) {
  const session = await getServerSession(cookies());
  const body = await req.json(); // { userId, role: 'HR_MANAGER' | 'VIEWER' }
  const result = await b2bServiceAddOrgMember({
    headers: buildAuthHeaders(session),
    body: {
      organisationId: params.id,
      userId: body.userId,
      role: body.role,
    }
  });
  return NextResponse.json(result.data);
}
```

### 1.2 — `DELETE /api/organisations/[id]/members/[memberId]`

**File**: `app/api/organisations/[id]/members/[memberId]/route.ts` (new file)

Calls `b2bServiceRemoveOrgMember`.

### 1.3 — `POST /api/organisations/[id]/assign-admin` (Assign existing user as admin)

**File**: `app/api/organisations/[id]/assign-admin/route.ts` (new file)

Calls `b2bServiceAssignOrgAdmin` with `{ organisationId, userId }`.
Distinct from `/admins` (which creates a NEW user AND assigns admin role).

### 1.4 — `DELETE /api/organisations/[id]/members/[memberId]/admin`

**File**: `app/api/organisations/[id]/members/[memberId]/admin/route.ts` (new file)

Removes the admin role. Calls `b2bServiceRemoveOrgMember` then re-adds as `VIEWER`.

### 1.5 — `GET /api/employees` (already exists — add `departmentId` filter)

**File**: `app/api/employees/route.ts`

Add `departmentId` query param → pass to `b2bServiceListEmployees`.

### 1.6 — `GET/POST/PATCH/DELETE /api/purchase-orders`

**File**: `app/api/purchase-orders/route.ts` (new file)

Wraps: `b2bServiceListPurchaseOrders`, `b2bServiceCreatePurchaseOrder`.

### 1.7 — `GET /api/purchase-orders/catalog`

**File**: `app/api/purchase-orders/catalog/route.ts` (new file)

Wraps: `b2bServiceListPurchaseOrderCatalog`.

---

## Phase 2 — Browser Clients (`src/lib/clients/`)

Each client is a thin fetch wrapper — mirrors the API routes.

### 2.1 — Update `organisation-client.ts`

Add missing methods:
```typescript
// Add org member (HR_MANAGER | VIEWER | BUSINESS_ADMIN)
addMember(orgId: string, userId: string, role: OrgMemberRole): Promise<OrgMember>

// Remove org member
removeMember(orgId: string, memberId: string): Promise<void>

// Assign existing user as admin (no new user creation)
assignAdmin(orgId: string, userId: string): Promise<void>
```

### 2.2 — New `purchase-order-client.ts`

```typescript
list(orgId: string, opts?: PaginationOpts): Promise<PaginatedResult<PurchaseOrder>>
getCatalog(): Promise<CatalogPlan[]>
create(payload: CreatePOPayload): Promise<PurchaseOrder>
```

### 2.3 — Update `employee-client.ts`

Add `departmentId?: string` filter to `list()`.

---

## Phase 3 — State Management (Zustand Stores)

**Directory**: `src/lib/stores/`

### 3.1 — `org-store.ts`

```typescript
interface OrgStore {
  // Current org context (b2b admin's own org, or super admin's selected org)
  currentOrg: Organisation | null;
  setCurrentOrg: (org: Organisation) => void;

  // All orgs (super_admin only)
  orgs: Organisation[];
  setOrgs: (orgs: Organisation[]) => void;

  // Members of currentOrg
  members: OrgMember[];
  setMembers: (members: OrgMember[]) => void;

  // Loading/error state
  isLoading: boolean;
  error: string | null;
}
```

### 3.2 — `user-store.ts` (already partially exists — extend)

Add:
```typescript
portal: 'PORTAL_SYSTEM' | 'PORTAL_B2B';
casbinRole: 'super_admin' | 'b2b_org_admin' | 'partner_user';
orgId: string | null;
```

### 3.3 — `purchase-order-store.ts`

```typescript
interface POStore {
  orders: PurchaseOrder[];
  catalog: CatalogPlan[];
  setOrders: (orders: PurchaseOrder[]) => void;
  setCatalog: (plans: CatalogPlan[]) => void;
}
```

---

## Phase 4 — UI Components

### 4.1 — Org Member Management Panel
**File**: `components/organisations/org-member-panel.tsx`

Renders inside the org detail view (modal or page).

```
┌─────────────────────────────────────────────────────┐
│  Organisation Members            [+ Add Member]      │
├────────────────┬───────────┬───────────┬────────────┤
│ Name           │ Role      │ Status    │ Actions    │
├────────────────┼───────────┼───────────┼────────────┤
│ Rahim Ahmed    │ Admin     │ Active    │ [Remove]   │
│ Karim Uddin    │ HR Manager│ Active    │ [Remove]   │
│ Fatema Begum   │ Viewer    │ Inactive  │ [Remove]   │
├────────────────┴───────────┴───────────┴────────────┤
│ Showing 3 of 3                                       │
└─────────────────────────────────────────────────────┘
```

**Props**: `orgId`, `currentUserRole: 'super_admin' | 'b2b_org_admin'`, `onMemberAdded`, `onMemberRemoved`

**Sub-components**: `OrgMemberRow`, `AddMemberModal`

### 4.2 — Add Member Modal
**File**: `components/organisations/add-member-modal.tsx`

Two tabs:
1. **Create New User** — name, mobile, email → `POST /api/organisations/[id]/admins`
2. **Assign Existing User** — user ID / mobile search → role dropdown → `POST /api/organisations/[id]/assign-admin` or `/members`

### 4.3 — Org Management Table (Super Admin)
**File**: `components/organisations/org-management-table.tsx`

TanStack Table with columns: Name, Code, Industry, Contact Email, Status, Created At, Actions (View/Edit/Delete/Approve).
Top bar: Search + Status filter + `[+ Create Organisation]`. Pagination: 10/25/50.

### 4.4 — Org Detail Side Panel
**File**: `components/organisations/org-detail-panel.tsx`

shadcn `Sheet` with tabs:
- **Info** — org fields + inline edit form + delete button
- **Members** — `<OrgMemberPanel>` component
- **Departments** — count + link to dept management

### 4.5 — Portal Header
**File**: `components/layout/portal-header.tsx`

Shows org name (from session), user name, role badge, logout button.

### 4.6 — Role-Based Navigation Sidebar
**File**: `components/layout/sidebar.tsx`

```typescript
const navItems = {
  super_admin:  ['Organisations /organisations', 'Users /users', 'Audit /audit', 'Settings /settings'],
  b2b_org_admin:['Dashboard /dashboard', 'Departments /departments', 'Employees /employees',
                 'Purchase Orders /purchase-orders', 'Team /team', 'Settings /settings'],
  partner_user: ['Dashboard /dashboard', 'Departments /departments',
                 'Employees /employees', 'Purchase Orders /purchase-orders'],
};
```

Nav items rendered based on `session.casbinRole` from session store.

### 4.7 — Purchase Order Components
**File**: `components/purchase-orders/purchase-order-list.tsx`

Columns: PO Number, Plan, Product, Employees, Total Premium, Status, Date.
Filter by status. `[+ New PO]` → `CreatePOModal`.

**File**: `components/purchase-orders/create-po-modal.tsx`

2-step wizard: (1) Select catalog plan → view details. (2) Select employees → confirm premium → submit.

### 4.8 — Team Management
**File**: `components/team/team-management.tsx`

B2B admin manages team (HR_MANAGER | VIEWER only — no admin creation here).
Columns: Name, Role badge, Status, `[Remove]` button.
`[+ Add Team Member]` → `AddMemberModal` restricted to HR_MANAGER|VIEWER roles.

---

## Phase 5 — Pages (App Router)

### 5.1 — Upgrade `/organisations` page (Super Admin)
**File**: `app/(dashboard)/organisations/page.tsx`

Add org detail side panel (`OrgDetailPanel`) opening on row click.
Panel includes member management and approve/reject actions for PENDING orgs.

### 5.2 — New `/dashboard` page (B2B Admin)
**File**: `app/(dashboard)/dashboard/page.tsx`

Summary cards: Total Employees, Active Departments, Active POs, Org Status.
Links to respective pages. Only for `b2b_org_admin` and `partner_user`.

### 5.3 — New `/team` page (B2B Admin Only)
**File**: `app/(dashboard)/team/page.tsx`

```typescript
export default function TeamPage() {
  const { session } = useSession();
  return <TeamManagement orgId={session.orgId} />;
}
```

Middleware blocks `partner_user` from this route.

### 5.4 — New `/purchase-orders` page
**File**: `app/(dashboard)/purchase-orders/page.tsx`

Accessible to `b2b_org_admin` and `partner_user`.

### 5.5 — Upgrade `/departments` page
Add employee count per dept. Row click filters employees by `departmentId`.

### 5.6 — Upgrade `/employees` page
Add `departmentId` query filter, department name column, status badge.

---

## Phase 6 — Route Guards & Session

### 6.1 — Middleware
**File**: `middleware.ts`

```typescript
const roleGuards: Record<string, string[]> = {
  '/organisations':   ['super_admin'],
  '/team':            ['b2b_org_admin'],
  '/purchase-orders': ['b2b_org_admin', 'partner_user'],
  '/employees':       ['b2b_org_admin', 'partner_user'],
  '/departments':     ['b2b_org_admin', 'partner_user'],
  '/dashboard':       ['b2b_org_admin', 'partner_user'],
};
```

Unauthenticated → redirect `/login`. Wrong role → redirect `/dashboard` (or `/organisations` for super_admin).

### 6.2 — Session casbinRole Mapping
**File**: `src/lib/backend-auth.ts`

```typescript
function mapUserTypeToCasbinRole(userType: string, portal: string): CasbinRole {
  if (portal === 'PORTAL_SYSTEM') return 'super_admin';
  switch (userType) {
    case 'BUSINESS_ADMIN': return 'b2b_org_admin';
    case 'HR_MANAGER':
    case 'VIEWER':
    default:               return 'partner_user';
  }
}
```

Session must carry: `casbinRole`, `orgId`, `portal`, `userId`, `tenantId`.

---

## Phase 7 — Auth Headers & SDK Integration

### 7.1 — Centralised `buildAuthHeaders`
**File**: `src/lib/api-helpers.ts`

```typescript
export function buildAuthHeaders(session: PortalSession) {
  return {
    'x-user-id':    session.userId,
    'x-business-id': session.orgId ?? '',
    'x-portal':     session.portal,       // PORTAL_SYSTEM | PORTAL_B2B
    'x-tenant-id':  session.tenantId ?? 'root',
    'Cookie':       session.rawCookieHeader,
    'X-CSRF-Token': session.csrfToken,
  };
}
```

### 7.2 — SDK Import Pattern
All API routes use named imports from `@lifeplus/insuretech-sdk` (local tarball installed by pipeline):

```typescript
import {
  b2bServiceListOrganisations, b2bServiceCreateOrganisation,
  b2bServiceGetOrganisation, b2bServiceUpdateOrganisation,
  b2bServiceDeleteOrganisation, b2bServiceListOrgMembers,
  b2bServiceAddOrgMember, b2bServiceRemoveOrgMember,
  b2bServiceAssignOrgAdmin, b2bServiceListDepartments,
  b2bServiceCreateDepartment, b2bServiceUpdateDepartment,
  b2bServiceDeleteDepartment, b2bServiceListEmployees,
  b2bServiceCreateEmployee, b2bServiceUpdateEmployee,
  b2bServiceDeleteEmployee, b2bServiceListPurchaseOrders,
  b2bServiceCreatePurchaseOrder, b2bServiceListPurchaseOrderCatalog,
} from '@lifeplus/insuretech-sdk';
```

---

## Complete File Change Matrix

| File | Action | Phase |
|------|--------|-------|
| `run_api_pipeline.ps1` | Add SDK tar+copy+install block | 0 |
| `b2b_portal/.gitignore` | Add `sdk/` to gitignore | 0 |
| `app/api/organisations/[id]/members/route.ts` | Add POST handler | 1.1 |
| `app/api/organisations/[id]/members/[memberId]/route.ts` | New — DELETE | 1.2 |
| `app/api/organisations/[id]/assign-admin/route.ts` | New — POST | 1.3 |
| `app/api/organisations/[id]/members/[memberId]/admin/route.ts` | New — DELETE | 1.4 |
| `app/api/employees/route.ts` | Add departmentId filter | 1.5 |
| `app/api/purchase-orders/route.ts` | New — GET/POST | 1.6 |
| `app/api/purchase-orders/catalog/route.ts` | New — GET | 1.7 |
| `src/lib/clients/organisation-client.ts` | Add addMember/removeMember | 2.1 |
| `src/lib/clients/purchase-order-client.ts` | New | 2.2 |
| `src/lib/clients/employee-client.ts` | Add departmentId filter | 2.3 |
| `src/lib/stores/org-store.ts` | New | 3.1 |
| `src/lib/stores/user-store.ts` | Extend with casbinRole/portal/orgId | 3.2 |
| `src/lib/stores/purchase-order-store.ts` | New | 3.3 |
| `components/organisations/org-member-panel.tsx` | New | 4.1 |
| `components/organisations/add-member-modal.tsx` | New | 4.2 |
| `components/organisations/org-management-table.tsx` | New | 4.3 |
| `components/organisations/org-detail-panel.tsx` | New | 4.4 |
| `components/layout/portal-header.tsx` | New | 4.5 |
| `components/layout/sidebar.tsx` | Upgrade role-based nav | 4.6 |
| `components/purchase-orders/purchase-order-list.tsx` | New | 4.7 |
| `components/purchase-orders/create-po-modal.tsx` | New | 4.7 |
| `components/team/team-management.tsx` | New | 4.8 |
| `app/(dashboard)/organisations/page.tsx` | Upgrade with OrgDetailPanel | 5.1 |
| `app/(dashboard)/dashboard/page.tsx` | New | 5.2 |
| `app/(dashboard)/team/page.tsx` | New | 5.3 |
| `app/(dashboard)/purchase-orders/page.tsx` | New | 5.4 |
| `app/(dashboard)/departments/page.tsx` | Add dept→employee link | 5.5 |
| `app/(dashboard)/employees/page.tsx` | Add departmentId filter + status badge | 5.6 |
| `middleware.ts` | Role-based route guards | 6.1 |
| `src/lib/backend-auth.ts` | Add casbinRole mapping | 6.2 |
| `src/lib/api-helpers.ts` | Centralise buildAuthHeaders | 7.1 |

---

## Implementation Order for Next Agent

**Execute in this exact order:**

1. **Phase 0** — Fix `run_api_pipeline.ps1` (SDK tarball pipeline)
2. **Phase 7.1** — Centralise `buildAuthHeaders` in `api-helpers.ts`
3. **Phase 6.2** — Fix `backend-auth.ts` `casbinRole` mapping in session
4. **Phase 1** — All missing API routes (proxy layer)
5. **Phase 2** — Update/create browser clients
6. **Phase 3** — Zustand stores
7. **Phase 4** — UI components (member panel → add member modal → org table → detail panel → sidebar → PO components → team)
8. **Phase 5** — Wire all pages
9. **Phase 6.1** — Middleware route guards

---

## Testing Checklist

- [ ] `run_api_pipeline.ps1` builds SDK, tars, copies to `b2b_portal/sdk/`, runs `npm install`
- [ ] Super admin logs in → `/organisations` → full CRUD
- [ ] Super admin creates org → assigns admin → Kafka event → b2b_org_admin role assigned
- [ ] Super admin assigns existing user as admin via panel → Kafka event fires
- [ ] B2B admin logs in → `/dashboard` → cannot reach `/organisations` (middleware blocks)
- [ ] B2B admin manages departments, employees, purchase orders
- [ ] B2B admin adds HR manager → HR manager logs in → `partner_user` role → sees depts/employees/POs
- [ ] B2B admin adds viewer → viewer gets read-only access
- [ ] B2B admin removes team member → role revoked
- [ ] HR manager blocked from `/team` → redirected to `/dashboard`
- [ ] All API routes use centralised `buildAuthHeaders`
- [ ] SDK imported from local tarball `@lifeplus/insuretech-sdk` works correctly
