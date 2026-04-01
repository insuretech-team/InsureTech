# B2B Portal Documentation Index

This directory contains comprehensive documentation of the B2B Portal's architecture, SDK clients, and API routes.

## Documents

### 1. **ARCHITECTURE_OVERVIEW.md** ⭐ START HERE
Complete guide to the portal's architecture and design patterns.

**Covers:**
- Project structure and file organization
- Authentication & session flow (login, validation, logout)
- SDK architecture (browser-side vs server-side clients)
- Authentication headers and portal header resolution
- API response envelope (GatewayResponse)
- Middleware and route guards
- API route patterns
- Key concepts (BFF pattern, Casbin auth, session refresh)
- Import boundaries
- Environment variables
- File manifest with descriptions

**Key Diagrams:**
- Login flow with cookie setup
- Session validation and metadata refresh
- Logout with cookie expiration
- Auth header forwarding (x-portal, x-user-id, x-business-id)

---

### 2. **SDK_CLIENT_REFERENCE.md** - Detailed SDK Usage
Complete reference for all SDK clients with examples.

**Sections:**
- **Browser-Side Clients** (6 files)
  - `auth-client.ts` - Authentication operations (login, logout, session, 2FA, OTP)
  - `employee-client.ts` - Employee CRUD with form fields
  - `department-client.ts` - Department management
  - `organisation-client.ts` - Organisation management with member roles
  - `purchase-order-client.ts` - Purchase order lifecycle + catalog
  - `docgen-client.ts` - Document generation and management

- **Server-Side Clients** (3 files)
  - `b2b-sdk-client.ts` - makeSdkClient() wrapper around @lifeplus/insuretech-sdk
  - `makeDirectHttp()` - Raw HTTP for endpoints not in SDK
  - `docgen-sdk-client.ts` - makeDocgenClient() for document service

- **Shared Utilities**
  - `shared.ts` - Universal types and helpers
  - `api-helpers.ts` - Server-side response builders
  - `session-headers.ts` - Auth context resolution

**For each client:**
- Method signatures and return types
- Query parameters and request/response payloads
- Usage examples in components and API routes
- TypeScript interfaces

---

### 3. **API_ROUTES_SUMMARY.md** - Complete Route Mapping
Comprehensive list of all Next.js API routes.

**Covers 45+ routes organized by domain:**
- **Authentication** (13 routes) - Login, logout, profile, 2FA, OTP, sessions
- **Employees** (5 routes) - List, create, get, update, delete
- **Departments** (5 routes) - List, create, get, update, delete
- **Organisations** (10 routes) - List, create, get, update, delete, members, admins, approval
- **Purchase Orders** (6 routes) - List, create, get, update, delete, catalog
- **Documents** (5 routes) - Generate, list, get, download, delete
- **Dashboard** (2 routes) - Stats, activity

**For each route:**
- HTTP method and path
- Purpose and description
- Corresponding browser client method
- Query parameters / request body format
- Response type
- Authentication requirements
- Special notes

**Summary table** showing all routes with HTTP methods and auth requirements.

---

## Quick Reference

### Import Patterns

**In React Components:**
```typescript
import {
  authClient,
  employeeClient,
  departmentClient,
  organisationClient,
  purchaseOrderClient,
  docgenClient
} from "@lib/sdk";
```

**In API Routes:**
```typescript
import { makeSdkClient, makeDirectHttp } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, unwrapSdkResult } from "@lib/sdk/api-helpers";
```

### API Route Pattern

```typescript
export async function GET(request: Request) {
  // 1. Resolve auth context
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false }, { status: 401 });
  
  // 2. Create SDK client
  const sdk = makeSdkClient(request, hdrs);
  
  // 3. Call backend
  const result = await sdk.listEmployees({ query: { page_size: 50 } });
  
  // 4. Unwrap response
  const unwrapped = unwrapSdkResult(result);
  if (!unwrapped.ok) {
    return NextResponse.json({ ok: false, message: unwrapped.message }, { status: unwrapped.status });
  }
  
  // 5. Return to browser
  return NextResponse.json({ ok: true, data: unwrapped.data });
}
```

### Authentication Flow

1. **Login** → User submits credentials → `/api/auth/login` → Sets cookies (session_token, csrf_token, portal_*)
2. **Session Check** → Component calls `authClient.getSession()` → Refreshes metadata cookies
3. **API Calls** → Component calls client (e.g., `employeeClient.list()`) → Forwards session_token to `/api/employees` → Route calls backend with x-portal + x-business-id headers
4. **Logout** → User clicks logout → `/api/auth/logout` → Expires all cookies

### Key Files by Function

| Function | File | Type |
|----------|------|------|
| List employees | `employee-client.ts` | Browser |
| Create employee | `employee-client.ts` | Browser |
| Call backend for employees | `b2b-sdk-client.ts` → makeSdkClient() | Server |
| Resolve auth headers | `session-headers.ts` | Server |
| Unwrap API response | `shared.ts` → unwrapGateway() | Both |
| Build error response | `api-helpers.ts` → gatewayError() | Server |
| Handle login | `backend-auth.ts` + `/api/auth/login` | Server |
| Create session | `session-store.ts` → createSession() | Server |
| Validate session | `session.ts` → getServerSession() | Server |

---

## Architecture Highlights

### 🏗️ BFF (Backend for Frontend) Pattern
- Browser never calls backend directly
- All requests go through Next.js API routes
- Routes handle auth, validation, transformation, security

### 🔐 Authentication
- **Cookie-based sessions** (session_token validated by backend)
- **Session store** in-memory on portal server (12-hour TTL)
- **Metadata cookies** (portal_role, portal_user_id, portal_biz_id) for middleware & header resolution
- **CSRF protection** via X-CSRF-Token header

### 👤 Authorization
- **Casbin RBAC** on backend
- **Portal context headers** (x-portal, x-business-id) route to correct Casbin domain
- **Roles:** SYSTEM_ADMIN (superadmin), B2B_ORG_ADMIN, BUSINESS_ADMIN, HR_MANAGER, VIEWER
- **Middleware guards** for UX-level route protection

### 📦 SDK Generation
- Types from `@lifeplus/insuretech-sdk` (auto-generated from protobuf)
- Browser clients are thin fetch wrappers calling API routes
- Server clients wrap the SDK with auth headers + cookie handling
- Direct HTTP client for endpoints not in generated SDK

### 📨 Response Envelope
- Unified `GatewayResponse<T>` structure: `{ success, data, error, meta }`
- Pagination metadata in `meta.pagination`
- Error details: code, message, error_id, http_status_code, retryable, field_violations

---

## File Structure Reference

```
src/lib/
├── sdk/                          # SDK Layer (main entry point: index.ts)
│   ├── shared.ts                # Universal types & unwrap helpers
│   ├── api-helpers.ts           # Server response builders
│   ├── auth-client.ts           # Browser: /api/auth/* client
│   ├── employee-client.ts       # Browser: /api/employees client
│   ├── department-client.ts     # Browser: /api/departments client
│   ├── organisation-client.ts   # Browser: /api/organisations client
│   ├── purchase-order-client.ts # Browser: /api/purchase-orders client
│   ├── docgen-client.ts         # Browser: /api/documents client
│   ├── b2b-sdk-client.ts        # Server: makeSdkClient() + makeDirectHttp()
│   ├── docgen-sdk-client.ts     # Server: makeDocgenClient()
│   ├── session-headers.ts       # Server: resolvePortalHeaders()
│   ├── dashboard-config.ts
│   └── index.ts                 # Central export barrel
├── auth/                        # Session & Auth Management
│   ├── backend-auth.ts          # Server: Gateway auth calls
│   ├── session.ts               # Server: Cookie helpers
│   ├── session-store.ts         # Server: In-memory store
│   └── resolve-user-id.ts       # Server: Fallback user ID resolution
├── types/
│   ├── auth.ts                  # PortalPrincipal, PortalSession, etc
│   └── [other types]
└── proto-generated/             # Auto-generated protobuf types (read-only)

app/api/
├── auth/                        # Auth routes
│   ├── login/route.ts
│   ├── logout/route.ts
│   ├── session/route.ts
│   └── [other auth routes]
├── employees/route.ts           # Employee CRUD
├── departments/route.ts         # Department CRUD
├── organisations/route.ts       # Organisation CRUD
├── purchase-orders/route.ts     # PO CRUD
├── documents/route.ts           # Document routes
└── dashboard/                   # Dashboard stats

middleware.ts                    # Edge middleware for auth & role guards
```

---

## Development Tips

### 1. Adding a New API Route

1. Create `app/api/[feature]/route.ts`
2. Import: `resolvePortalHeaders`, `makeSdkClient`, response builders
3. Call `resolvePortalHeaders()` first to get auth context
4. Create SDK client: `makeSdkClient(request, hdrs)`
5. Call backend method and unwrap response
6. Return `NextResponse.json()`

**Example:**
```typescript
export async function GET(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return unauthorized();
  
  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.someMethod({ ... });
  
  if (!result.response.ok) {
    return gatewayError(sdkErrorMessage(result), result.response.status);
  }
  
  return NextResponse.json({ ok: true, data: result.data });
}
```

### 2. Using Clients in Components

```typescript
"use client";

import { employeeClient } from "@lib/sdk";
import { useEffect, useState } from "react";

export function EmployeeList() {
  const [employees, setEmployees] = useState([]);
  
  useEffect(() => {
    employeeClient.list({ pageSize: 50 }).then(res => {
      if (res.ok) setEmployees(res.employees);
    });
  }, []);
  
  return <div>{/* render employees */}</div>;
}
```

### 3. Debugging Auth Issues

**Superadmin getting 403?**
- Check `x-portal=PORTAL_SYSTEM` header is being sent
- Ensure metadata cookies are fresh (call `/api/auth/session`)
- Verify `portal_role` cookie is set to "SYSTEM_ADMIN"

**Session expired?**
- Check `session_token` cookie exists
- Session TTL is 12 hours from last activity
- Call `/api/auth/session` to refresh metadata cookies

**CSRF token missing?**
- Ensure `csrf_token` cookie is set at login
- Routes should forward `X-CSRF-Token` header automatically
- Check middleware isn't clearing cookies

### 4. Environment Setup

```bash
# Required
INSURETECH_API_BASE_URL=http://localhost:8080

# Optional
INSURETECH_API_KEY=<unused for cookie auth>
DEFAULT_TENANT_ID=00000000-0000-0000-0000-000000000001
NODE_ENV=development
```

---

## Common Patterns

### Error Handling
```typescript
// SDK result
const result = await sdk.listEmployees({ ... });
if (!result.response.ok) {
  const message = sdkErrorMessage(result);  // Extract user-facing error
  return NextResponse.json({ ok: false, message }, { status: result.response.status });
}

// Direct HTTP
const http = makeDirectHttp(request, hdrs);
const res = await http.post("/v1/b2b/organisations/123/admins", body);
if (!res.ok) {
  return NextResponse.json({ ok: false, message: res.message }, { status: res.status });
}
```

### Pagination
```typescript
const result = await sdk.listEmployees({
  query: {
    page_size: 20,      // Items per page
    business_id: "..."  // Required for non-superadmin
  }
});

// Check meta for pagination info
const { total_count, total_pages, has_next } = result.data?.meta?.pagination ?? {};
```

### Form Submission
```typescript
const payload = {
  name: "John Doe",
  employeeId: "EMP-001",
  businessId: hdrs?.businessId ?? "",  // From auth context
  departmentId: "...",
  email: "...",
  // ... other fields
};

const result = await employeeClient.create(payload);
if (result.ok) {
  // Success: result.employee contains created record
  showToast("Employee created!");
} else {
  // Error: result.message has user-facing error
  showError(result.message);
}
```

---

## Where to Go Next

1. **Understanding the current flow?** → Read **ARCHITECTURE_OVERVIEW.md**
2. **Using a specific client?** → Check **SDK_CLIENT_REFERENCE.md** for that client
3. **Need to find an endpoint?** → Search **API_ROUTES_SUMMARY.md**
4. **Building a new feature?** → Follow the API Route Pattern above, use the appropriate client

---

## Summary

The B2B Portal is a **Next.js application** with:
- **Frontend:** React components using fetch-based SDK clients
- **Backend For Frontend (BFF):** Next.js API routes handling auth, validation, transformation
- **Backend:** InsureTech gateway with gRPC services (accessed via @lifeplus/insuretech-sdk)
- **Auth:** Cookie-based sessions with Casbin RBAC
- **SDK:** Thin browser clients + wrapped server clients with automatic auth header injection

**All communication is authenticated via session cookies, with role-based authorization enforced by the backend Casbin engine.**

Created: 2024  
Last Updated: Now
