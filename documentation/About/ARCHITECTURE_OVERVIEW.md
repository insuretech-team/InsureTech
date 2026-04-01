# B2B Portal Architecture Overview

## Project Structure

```
E:\Projects\InsureTech\b2b_portal\
├── app/                           # Next.js App Router (13+)
│   ├── api/auth/                  # Auth API routes
│   │   ├── login/
│   │   ├── logout/
│   │   ├── session/
│   │   ├── refresh/
│   │   ├── profile/
│   │   ├── change-password/
│   │   ├── totp/
│   │   ├── send-otp/
│   │   ├── verify-otp/
│   │   ├── send-email-otp/
│   │   ├── verify-email/
│   │   ├── sessions/                # Multi-session management
│   │   ├── sessions/[sessionId]/
│   │   └── profile-photo-url/
│   ├── api/employees/             # B2B CRUD API routes
│   │   ├── route.ts               # GET list, POST create
│   │   ├── [id]/route.ts          # GET single, PATCH update, DELETE
│   │   ├── bulk-upload/
│   │   └── template/
│   ├── api/organisations/         # Organisation management
│   │   ├── route.ts
│   │   ├── [id]/route.ts
│   │   ├── [id]/members/
│   │   ├── [id]/admins/
│   │   ├── [id]/assign-admin/
│   │   ├── [id]/approve/
│   │   └── me/
│   ├── api/departments/
│   ├── api/purchase-orders/
│   ├── api/documents/
│   ├── api/document-templates/
│   ├── api/dashboard/
│   ├── login/page.tsx
│   ├── employees/page.tsx
│   ├── departments/page.tsx
│   ├── organisations/page.tsx
│   └── [other routes]/
├── src/lib/
│   ├── sdk/                       # **Main SDK Layer - Browser & Server**
│   │   ├── shared.ts              # Universal types (ApiResult, GatewayResponse)
│   │   ├── api-helpers.ts         # Server-only: response builders, unwrappers
│   │   ├── auth-client.ts         # Browser: /api/auth/* client
│   │   ├── employee-client.ts     # Browser: /api/employees client
│   │   ├── department-client.ts   # Browser: /api/departments client
│   │   ├── organisation-client.ts # Browser: /api/organisations client
│   │   ├── purchase-order-client.ts
│   │   ├── docgen-client.ts       # Browser: /api/documents client
│   │   ├── b2b-sdk-client.ts      # Server: @lifeplus/insuretech-sdk wrapper
│   │   ├── docgen-sdk-client.ts   # Server: Direct HTTP for docgen endpoints
│   │   ├── session-headers.ts     # Server: Resolves auth headers (x-portal, x-user-id, etc)
│   │   ├── dashboard-config.ts
│   │   └── index.ts               # Central export barrel
│   ├── auth/                      # **Session & Auth Management**
│   │   ├── backend-auth.ts        # Server: Gateway auth service calls
│   │   ├── session.ts             # Server: Next.js session cookie helper
│   │   ├── session-store.ts       # Server: In-memory session store
│   │   └── resolve-user-id.ts     # Server: Fallback to resolve user_id
│   ├── types/
│   │   ├── auth.ts                # PortalPrincipal, PortalSession, PortalAuthResponse
│   │   ├── b2b.ts
│   │   └── [other types]/
│   └── proto-generated/           # Generated Protobuf types (read-only)
├── components/                    # React components (client-side)
├── middleware.ts                  # Edge middleware: auth & role guards
├── next.config.ts
└── package.json
```

---

## Authentication & Session Flow

### 1. **Login Flow** (`/api/auth/login`)

```
Browser Login Form
    ↓
POST /api/auth/login
    ↓ (server-side)
backend-auth.ts → authServiceLogin() [calls gateway]
    ↓
Gateway validates credentials → Returns SessionResponse + JWT
    ↓
backend-auth.ts → Creates PortalSession (local in-memory store)
    ↓ (Set cookies on response)
Cookies:
  - session_token (HttpOnly): sessionId from store
  - csrf_token (HttpOnly): generated token
  - portal_role: User's role (SYSTEM_ADMIN, BUSINESS_ADMIN, etc)
  - portal_user_id: User ID
  - portal_biz_id: Organisation/Business ID (B2B users only)
  - portal_email: Contact email (optional)
  - portal_mobile: Mobile number (optional)
    ↓
Browser receives response → Stores session_token cookie
```

**Key Points:**
- Session is stored in **in-memory Map** on the server (session-store.ts)
- session_token is a UUID, hashed and stored with session data
- Lightweight metadata cookies (portal_role, portal_user_id, portal_biz_id) are set for middleware to read
- These metadata cookies are NOT HttpOnly so middleware can check them
- All auth requests forward session_token cookie to the gateway for validation

### 2. **Session Validation** (`/api/auth/session`)

```
GET /api/auth/session (includes session_token cookie)
    ↓
backend-auth.ts → getCurrentSession() [calls gateway]
    ↓
Gateway validates session_token cookie
    ↓
Returns current session + user profile
    ↓
Re-mint metadata cookies to keep in sync with backend
    ↓
Response includes:
  - portal_role (re-minted)
  - portal_user_id (re-minted)
  - portal_biz_id (re-minted)
  - portal_email (preserved from existing cookie)
  - portal_mobile (preserved from existing cookie)
```

**Critical:** Metadata cookies must be re-minted on every session refresh. If they expire while session_token is still valid, resolvePortalHeaders() falls back to PORTAL_B2B with no org context, causing 403 errors for superadmin.

### 3. **Logout Flow** (`/api/auth/logout`)

```
POST /api/auth/logout (with session_token cookie)
    ↓
backend-auth.ts → logoutCurrentSession() [calls gateway]
    ↓
Gateway invalidates session_token
    ↓
Expire all auth cookies:
  - session_token (set expires=Date(0))
  - csrf_token
  - portal_role
  - portal_user_id
  - portal_biz_id
    ↓
Browser loses session → Middleware redirects to /login
```

---

## SDK Architecture

### Browser-Side Clients (Fetch-based)

Used by **React components and hooks**. They call Next.js API routes (BFF pattern).

```typescript
// src/lib/sdk/auth-client.ts
export const authClient = {
  login(payload: PortalLoginRequest): Promise<PortalAuthResponse>
  logout(): Promise<AuthOkResponse>
  getSession(): Promise<PortalAuthResponse>
  refreshToken(): Promise<AuthOkResponse>
  getProfile(): Promise<ProfileResponse>
  updateProfile(payload): Promise<ProfileResponse>
  changePassword(payload): Promise<AuthOkResponse>
  enableTotp(): Promise<TotpResponse>
  disableTotp(totpCode): Promise<AuthOkResponse>
  sendOtp(purpose): Promise<AuthOkResponse>
  verifyOtp(otp, purpose): Promise<OtpResponse>
  sendEmailOtp(purpose): Promise<AuthOkResponse>
  verifyEmail(payload): Promise<AuthOkResponse>
  listSessions(): Promise<SessionsResponse>
  revokeSession(sessionId): Promise<AuthOkResponse>
  revokeAllSessions(): Promise<AuthOkResponse>
  getProfilePhotoUploadUrl(): Promise<ProfilePhotoUrlResponse>
};
```

```typescript
// src/lib/sdk/employee-client.ts
export const employeeClient = {
  list(options?: { pageSize?, offset?, businessId?, departmentId?, status? }): Promise<EmployeeListResult>
  get(id: string): Promise<EmployeeSingleResult>
  create(payload: EmployeeCreatePayload): Promise<EmployeeSingleResult>
  update(id: string, payload: EmployeeUpdatePayload): Promise<EmployeeSingleResult>
  delete(id: string): Promise<ApiResult>
};
```

```typescript
// src/lib/sdk/organisation-client.ts
export const organisationClient = {
  list(): Promise<OrgListResult>
  get(id: string): Promise<OrgSingleResult>
  getMe(): Promise<OrgSingleResult>
  create(payload: OrgCreatePayload): Promise<OrgSingleResult>
  update(id: string, payload: OrgUpdatePayload): Promise<OrgSingleResult>
  delete(id: string): Promise<ApiResult>
  listMembers(id: string): Promise<OrgMembersResult>
  addMember(id: string, userId: string, role: string): Promise<OrgMemberResult>
  assignAdmin(id: string, memberId: string): Promise<OrgMemberResult>
  createAdmin(id: string, payload: OrgAdminCreatePayload): Promise<OrgMemberResult>
  removeMember(id: string, memberId: string): Promise<ApiResult>
  assignExistingAdmin(id: string, userId: string): Promise<ApiResult>
  approve(id: string): Promise<OrgSingleResult>
};
```

**Import Pattern:**
```typescript
import { authClient, employeeClient, organisationClient } from "@lib/sdk";
```

### Server-Side SDK Clients

Used by **Next.js API route handlers** to call the backend gateway.

#### `makeSdkClient()` - Wrapped @lifeplus/insuretech-sdk

```typescript
// src/lib/sdk/b2b-sdk-client.ts
export function makeSdkClient(request: Request, sessionOverrides?: PortalHeaders) {
  // Extracts cookies + resolves x-portal, x-user-id, x-business-id headers
  // Returns typed wrapper around @lifeplus/insuretech-sdk
  
  return {
    // ── Auth ─────────────────────────────────────────
    emailLogin(opts)
    logout(opts)
    validateToken(opts)
    registerEmailUser(opts)
    getCurrentSession(opts)
    refreshToken(opts)
    changePassword(opts)
    getUserProfile(opts)
    updateUserProfile(opts)
    getProfilePhotoUploadUrl(opts)
    listSessions(opts)
    revokeSession(opts)
    revokeAllSessions(opts)
    enableTotp(opts)
    disableTotp(opts)
    sendOtp(opts)
    verifyOtp(opts)
    sendEmailOtp(opts)
    verifyEmail(opts)
    
    // ── Employees ────────────────────────────────────
    listEmployees(opts)
    createEmployee(opts)
    getEmployee(opts)
    updateEmployee(opts)
    deleteEmployee(opts)
    
    // ── Departments ──────────────────────────────────
    listDepartments(opts)
    createDepartment(opts)
    getDepartment(opts)
    updateDepartment(opts)
    deleteDepartment(opts)
    
    // ── Purchase Orders ──────────────────────────────
    listPurchaseOrders(opts)
    createPurchaseOrder(opts)
    getPurchaseOrder(opts)
    listPurchaseOrderCatalog(opts)
    updatePurchaseOrderHttp(id, body)  // Direct HTTP
    deletePurchaseOrderHttp(id)        // Direct HTTP
    
    // ── Organisations ────────────────────────────────
    listOrganisations(opts)
    createOrganisation(opts)
    getOrganisation(opts)
    updateOrganisation(opts)
    deleteOrganisation(opts)
    listOrgMembers(opts)
    addOrgMember(opts)
    assignOrgAdmin(opts)
    removeOrgMember(opts)
  };
}
```

**Usage in API Routes:**
```typescript
export async function GET(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  const sdk = makeSdkClient(request, hdrs ?? undefined);
  
  const result = await sdk.listEmployees({
    query: { page_size: 50, business_id: "..." }
  });
  
  if (!result.response.ok) {
    return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
  }
  
  return NextResponse.json({ ok: true, data: result.data });
}
```

#### `makeDirectHttp()` - Raw HTTP with auth headers

```typescript
// src/lib/sdk/b2b-sdk-client.ts
export function makeDirectHttp(request: Request, sessionOverrides?: PortalHeaders) {
  // Returns untyped HTTP helpers for endpoints not in generated SDK
  
  return {
    get(path: string): Promise<HttpResult>
    post(path: string, body?: unknown): Promise<HttpResult>
    patch(path: string, body?: unknown): Promise<HttpResult>
    put(path: string, body?: unknown): Promise<HttpResult>
    delete(path: string): Promise<HttpResult>
  };
}
```

**Usage:**
```typescript
const http = makeDirectHttp(request, hdrs ?? undefined);
const res = await http.post(`/v1/b2b/organisations/${id}/admins`, adminPayload);
```

#### `makeDocgenClient()` - Document generation

```typescript
// src/lib/sdk/docgen-sdk-client.ts
export function makeDocgenClient(request: Request, sessionOverrides?: PortalHeaders) {
  return {
    generate(payload: GenerateDocumentPayload): Promise<DocumentSingleResult>
    list(path: string, query): Promise<DocumentListResult>
    get(documentId: string): Promise<DocumentSingleResult>
    download(documentId: string): Promise<DocumentDownloadResult>
    delete(documentId: string): Promise<ApiResult>
    createTemplate(payload: CreateDocumentTemplatePayload): Promise<DocumentSingleResult>
    getTemplate(templateId: string): Promise<DocumentSingleResult>
    listTemplates(query): Promise<DocumentListResult>
    updateTemplate(templateId: string, payload): Promise<DocumentSingleResult>
    deactivateTemplate(templateId: string): Promise<ApiResult>
    deleteTemplate(templateId: string): Promise<ApiResult>
  };
}
```

---

## Authentication Headers

### Session Resolution (`resolvePortalHeaders()`)

Located in `src/lib/sdk/session-headers.ts`. Called by every API route to resolve auth context.

```typescript
export async function resolvePortalHeaders(request: Request): Promise<PortalHeaders | null> {
  const cookieHeader = request.headers.get("cookie") ?? "";
  
  // Require backend session cookie — if absent, request is unauthenticated
  const sessionToken = extractCookie(cookieHeader, "session_token");
  if (!sessionToken) return null;
  
  // Read lightweight metadata cookies (set at login, re-minted on session refresh)
  const role = extractCookie(cookieHeader, "portal_role") || "BUSINESS_ADMIN";
  const userId = extractCookie(cookieHeader, "portal_user_id");
  const businessId = extractCookie(cookieHeader, "portal_biz_id");
  
  const portal = roleToPortal(role);  // Maps role to x-portal header value
  const tenantId = process.env.DEFAULT_TENANT_ID ?? "00000000-0000-0000-0000-000000000001";
  
  return { portal, userId, businessId, tenantId };
}
```

### Header Forwarding

The `makeSdkClient()` and `makeDirectHttp()` functions forward these headers to the gateway:

| Header | Value | Purpose |
|--------|-------|---------|
| `x-portal` | `PORTAL_SYSTEM` (superadmin) or `PORTAL_B2B` (org users) | Route to correct Casbin domain |
| `x-user-id` | User UUID | Identifies requesting user |
| `x-business-id` | Organisation/Business UUID | Identifies org context (B2B only) |
| `x-tenant-id` | Tenant UUID | Multi-tenant identifier |
| `X-CSRF-Token` | CSRF token from cookie | CSRF protection |
| `cookie` | Raw cookie header | Session validation |

**Critical:** Without `x-portal` and `x-business-id`, the backend authz interceptor returns 403 (Casbin domain resolution fails).

---

## API Response Envelope

All API endpoints (both BFF and direct HTTP) return a unified envelope:

```typescript
// src/lib/sdk/shared.ts

export interface GatewayResponse<T> {
  success: boolean;
  data: T | null;
  error: GatewayError | null;
  meta: GatewayMeta;
}

export interface GatewayError {
  code: string;
  message: string;
  error_id: string;
  http_status_code: number;
  retryable: boolean;
  field_violations: Array<{ field: string; description: string }>;
}

export interface GatewayMeta {
  request_id: string;
  timestamp: string;
  pagination?: {
    page: number;
    page_size: number;
    total_count: number;
    total_pages: number;
    has_next: boolean;
    has_prev: boolean;
  } | null;
}
```

### Unwrapping Responses

```typescript
// In server-side API routes
function unwrapGateway<T>(body: GatewayResponse<T>, httpStatus?: number) {
  if (body.success && body.data !== null) {
    return { ok: true, data: body.data, meta: body.meta };
  }
  
  const err = body.error;
  return {
    ok: false,
    message: err?.message ?? 'An unexpected error occurred',
    code: err?.code ?? 'UNKNOWN_ERROR',
    status: err?.http_status_code ?? httpStatus ?? 500,
    retryable: err?.retryable ?? false,
  };
}
```

---

## Middleware & Route Guards

### Edge Middleware (`middleware.ts`)

Runs at the edge (Cloudflare Workers / Vercel Edge) before the request reaches Next.js.

```typescript
export function middleware(request: NextRequest) {
  // Public paths (no auth required)
  if (pathname === "/login" || pathname.startsWith("/api/auth/login")) {
    return NextResponse.next();
  }
  
  // Check session_token cookie
  const hasSessionCookie = Boolean(request.cookies.get("session_token")?.value);
  
  if (!hasSessionCookie && !isPublic) {
    // Redirect to login with ?next parameter
    return NextResponse.redirect(new URL("/login?next=" + pathname, request.url));
  }
  
  // Role-based route guards (uses portal_role cookie)
  const role = request.cookies.get("portal_role")?.value ?? "";
  const guard = ROLE_GUARDS.find(g => pathname.startsWith(g.prefix));
  if (guard && !guard.allowedRoles.includes(role)) {
    // Redirect to appropriate default page for their role
    return NextResponse.redirect(new URL(fallback, request.url));
  }
  
  return NextResponse.next();
}
```

**Role Guard Configuration:**
```typescript
const ROLE_GUARDS = [
  { prefix: "/organisations", allowedRoles: ["SYSTEM_ADMIN"] },
  { prefix: "/team", allowedRoles: ["SYSTEM_ADMIN", "B2B_ORG_ADMIN", "BUSINESS_ADMIN"] },
  { prefix: "/departments", allowedRoles: ["SYSTEM_ADMIN", "B2B_ORG_ADMIN", "BUSINESS_ADMIN", "HR_MANAGER", "VIEWER"] },
  { prefix: "/employees", allowedRoles: ["SYSTEM_ADMIN", "B2B_ORG_ADMIN", "BUSINESS_ADMIN", "HR_MANAGER", "VIEWER"] },
  { prefix: "/purchase-orders", allowedRoles: ["SYSTEM_ADMIN", "B2B_ORG_ADMIN", "BUSINESS_ADMIN", "HR_MANAGER", "VIEWER"] },
];
```

**Important:** Middleware does NOT check API routes (`/api/*`). They handle auth via session cookie forwarding to the backend.

---

## API Route Pattern

All API routes follow this pattern:

```typescript
import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, unwrapSdkResult } from "@lib/sdk/api-helpers";

export async function GET(request: Request) {
  try {
    // 1. Resolve auth context from cookies
    const hdrs = await resolvePortalHeaders(request);
    if (!hdrs) {
      return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
    }
    
    // 2. Create SDK client (includes auth headers)
    const sdk = makeSdkClient(request, hdrs);
    
    // 3. Call backend service
    const result = await sdk.listEmployees({ query: { page_size: 50 } });
    
    // 4. Unwrap response
    const unwrapped = unwrapSdkResult(result);
    if (!unwrapped.ok) {
      return NextResponse.json({ ok: false, message: unwrapped.message }, { status: unwrapped.status });
    }
    
    // 5. Return to browser
    return NextResponse.json({ ok: true, employees: unwrapped.data?.employees ?? [] });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}

export async function POST(request: Request) {
  try {
    const hdrs = await resolvePortalHeaders(request);
    if (!hdrs) {
      return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
    }
    
    const sdk = makeSdkClient(request, hdrs);
    const body = await request.json();
    
    const result = await sdk.createEmployee({ body: { ... } });
    
    if (!result.response.ok) {
      return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
    }
    
    return NextResponse.json({ ok: true, employee: result.data?.employee ?? null }, { status: 201 });
  } catch (err) {
    return NextResponse.json({ ok: false, message: "Error" }, { status: 502 });
  }
}
```

---

## Key Concepts

### 1. **BFF (Backend for Frontend) Pattern**
- Browser never calls gateway directly
- All requests go through Next.js API routes
- API routes add headers, validate, transform, and secure data

### 2. **Cookie-Based Sessions**
- Session validation happens on the gateway (backend)
- Portal stores session metadata in cookies for middleware to check
- Session TTL: 12 hours

### 3. **Multi-Role Access Control**
- **SYSTEM_ADMIN**: Super admin, can manage all orgs (x-portal=PORTAL_SYSTEM)
- **B2B_ORG_ADMIN**: Organisation admin (x-portal=PORTAL_B2B, has x-business-id)
- **BUSINESS_ADMIN**: Old role (treated as B2B admin)
- **HR_MANAGER**: Can manage employees, departments
- **VIEWER**: Read-only access

### 4. **Casbin Authorization**
- Backend uses Casbin for RBAC
- Routes requests to domain: `system:root` (superadmin) or `org:{business_id}` (B2B)
- Portal must send correct `x-portal` + `x-business-id` headers

### 5. **Session Refresh Strategy**
- Session is validated on every `/api/auth/session` call
- Metadata cookies are re-minted to stay in sync with backend
- If metadata cookies expire while session_token is valid, fallback is PORTAL_B2B (not PORTAL_SYSTEM)

### 6. **Error Handling**
- Gateway returns structured error with: code, message, error_id, http_status_code, retryable flag
- Portal unwraps envelope using `unwrapGateway()` / `unwrapSdkResult()`
- User-facing errors extracted via `extractGatewayError()`

---

## Import Boundaries

### ✅ Correct Imports

**In components/hooks:**
```typescript
import { authClient, employeeClient } from "@lib/sdk";
```

**In API routes:**
```typescript
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, badRequest, gatewayError } from "@lib/sdk/api-helpers";
```

### ❌ Avoid

```typescript
// Don't import from deprecated folder
import { ... } from "@lib/clients";

// Don't use makeDirectHttp in components
import { makeDirectHttp } from "@lib/sdk/b2b-sdk-client";  // Server-only!
```

---

## Environment Variables

```bash
# Required
INSURETECH_API_BASE_URL=http://localhost:8080
# OR
NEXT_PUBLIC_INSURETECH_API_BASE_URL=http://localhost:8080

# Optional
INSURETECH_API_KEY=<optional, unused for cookie auth>
DEFAULT_TENANT_ID=00000000-0000-0000-0000-000000000001
```

---

## File Manifest

| File | Purpose |
|------|---------|
| `shared.ts` | Universal types: ApiResult, GatewayResponse, unwrapGateway, extractGatewayError |
| `api-helpers.ts` | Server helpers: badRequest, gatewayError, unauthorized, forbidden, sdkErrorMessage |
| `auth-client.ts` | Browser: fetch-based auth client for components |
| `employee-client.ts` | Browser: fetch-based employee CRUD client |
| `department-client.ts` | Browser: fetch-based department client |
| `organisation-client.ts` | Browser: fetch-based organisation client |
| `purchase-order-client.ts` | Browser: fetch-based purchase order client |
| `docgen-client.ts` | Browser: fetch-based document generation client |
| `b2b-sdk-client.ts` | Server: makeSdkClient + makeDirectHttp factories |
| `docgen-sdk-client.ts` | Server: makeDocgenClient for document service |
| `session-headers.ts` | Server: resolvePortalHeaders for auth context extraction |
| `index.ts` | Central export barrel |
| `backend-auth.ts` | Server: Gateway auth service calls (login, logout, getCurrentSession) |
| `session.ts` | Server: Next.js session cookie helpers (getServerSession, requireServerSession) |
| `session-store.ts` | Server: In-memory session storage (createSession, getSession, clearSession) |
| `resolve-user-id.ts` | Server: Fallback to resolve user_id from gateway |

---

## Summary

The B2B Portal uses a **BFF architecture** with:

1. **Browser-side clients** (fetch-based) → Call Next.js API routes
2. **API route handlers** → Validate session, add headers, call gateway
3. **Server-side SDK wrapper** → makeSdkClient wraps @lifeplus/insuretech-sdk
4. **Session management** → Cookie-based, validated by gateway
5. **Auth headers** → x-portal, x-user-id, x-business-id for Casbin routing
6. **Role guards** → Middleware + server-side validation
7. **Unified response envelope** → GatewayResponse with success/data/error/meta

All communication with the gateway is authenticated via session cookies + CSRF tokens. The portal is stateless — session data lives in the backend, with lightweight metadata cookies for edge middleware to check.
