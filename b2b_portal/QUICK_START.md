# Quick Start Guide

A fast reference for common tasks in the B2B Portal.

---

## 🚀 5-Minute Overview

### What is this?
A **Next.js B2B Portal** for managing employees, departments, organisations, and insurance policies. Uses a **BFF (Backend for Frontend) pattern** with cookie-based authentication.

### How does auth work?
1. User logs in → Browser gets `session_token` + metadata cookies
2. Every API call forwards `session_token` to the backend
3. Backend validates it and injects `x-portal` + `x-business-id` headers
4. Routes are protected by role-based access control (Casbin)

### What's the project structure?

```
Browser Components
       ↓
authClient, employeeClient, etc. (fetch-based)
       ↓
/api/auth/*, /api/employees/*, etc. (Next.js routes)
       ↓
makeSdkClient(), makeDirectHttp() (SDK wrappers)
       ↓
Backend Gateway (gRPC services)
```

---

## 📖 Documentation Map

| Document | Read If... |
|----------|-----------|
| **ARCHITECTURE_OVERVIEW.md** | You want to understand how everything fits together |
| **SDK_CLIENT_REFERENCE.md** | You need to know how to use a specific client |
| **API_ROUTES_SUMMARY.md** | You're looking for a specific API endpoint |
| **QUICK_START.md** (this file) | You want quick answers to common questions |

---

## ⚡ Common Tasks

### Task: List Employees in a Component

```typescript
"use client";

import { employeeClient } from "@lib/sdk";
import { useEffect, useState } from "react";

export function EmployeeList() {
  const [employees, setEmployees] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    employeeClient
      .list({ pageSize: 50, offset: 0 })
      .then((response) => {
        if (response.ok) {
          setEmployees(response.employees ?? []);
        } else {
          setError(response.message ?? "Failed to load employees");
        }
      })
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <table>
      <tbody>
        {employees.map((emp) => (
          <tr key={emp.id}>
            <td>{emp.name}</td>
            <td>{emp.department}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
```

### Task: Create an Employee

```typescript
const response = await employeeClient.create({
  name: "John Doe",
  employeeId: "EMP-001",
  businessId: "org-uuid", // Get from auth context
  departmentId: "dept-uuid",
  email: "john@example.com",
  insuranceCategory: 1, // Health
  coverageAmount: 100000,
});

if (response.ok) {
  console.log("Created:", response.employee);
  showToast("Employee created!");
} else {
  showError(response.message);
}
```

### Task: Build a New API Route

**File:** `app/api/my-feature/route.ts`

```typescript
import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, unwrapSdkResult, badRequest } from "@lib/sdk/api-helpers";

export async function GET(request: Request) {
  try {
    // 1. Resolve auth (required for all routes except /api/auth/login)
    const hdrs = await resolvePortalHeaders(request);
    if (!hdrs) {
      return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
    }

    // 2. Create SDK client with auth headers
    const sdk = makeSdkClient(request, hdrs);

    // 3. Call backend service
    const result = await sdk.listEmployees({
      query: { page_size: 50, business_id: hdrs.businessId },
    });

    // 4. Unwrap response
    const unwrapped = unwrapSdkResult(result);
    if (!unwrapped.ok) {
      return NextResponse.json(
        { ok: false, message: unwrapped.message },
        { status: unwrapped.status }
      );
    }

    // 5. Return to browser
    return NextResponse.json({
      ok: true,
      employees: unwrapped.data?.employees ?? [],
    });
  } catch (error) {
    return NextResponse.json(
      { ok: false, message: "Internal error" },
      { status: 500 }
    );
  }
}

export async function POST(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false }, { status: 401 });

  const body = await request.json();
  const sdk = makeSdkClient(request, hdrs);

  // ... create logic
}
```

### Task: Handle Superadmin Routes

**Problem:** Superadmin getting 403 on `/departments` or `/organisations`?

**Solution:** Ensure `x-portal=PORTAL_SYSTEM` is being sent. This happens automatically when:
1. `portal_role` cookie is "SYSTEM_ADMIN"
2. `resolvePortalHeaders()` maps it to `x-portal=PORTAL_SYSTEM`

**Debug checklist:**
```typescript
const hdrs = await resolvePortalHeaders(request);
console.log("Portal header:", hdrs?.portal); // Should be "PORTAL_SYSTEM"
console.log("Business ID:", hdrs?.businessId); // Should be empty for superadmin

const sdk = makeSdkClient(request, hdrs);
// SDK automatically forwards x-portal=PORTAL_SYSTEM
```

### Task: Login Flow

**Already implemented!** See `app/api/auth/login/route.ts` for reference.

**In components:**
```typescript
const response = await authClient.login({
  mobileNumber: "+8801700000000",
  password: "password123",
});

if (response.ok) {
  // Session stored in cookie, user is logged in
  router.push("/");
} else {
  showError(response.message);
}
```

### Task: Get Current Session

**Why?** To refresh metadata cookies (portal_role, portal_user_id, portal_biz_id).

**When?** Call periodically or after any role/org change.

```typescript
const response = await authClient.getSession();

if (response.ok) {
  const session = response.session;
  console.log("User:", session?.principal.displayName);
  console.log("Role:", session?.principal.role);
  console.log("Org:", session?.principal.organisationName);
} else {
  // Session expired, redirect to login
  router.push("/login");
}
```

### Task: Call an Endpoint NOT in the Generated SDK

**Use `makeDirectHttp()`:**

```typescript
export async function POST(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return unauthorized();

  const http = makeDirectHttp(request, hdrs); // Raw HTTP with auth
  const body = await request.json();

  // Call endpoint not in @lifeplus/insuretech-sdk
  const res = await http.post(
    `/v1/b2b/organisations/${id}/admins:assign`,
    { member_id: memberId }
  );

  if (!res.ok) {
    return NextResponse.json({ ok: false, message: res.message }, { status: res.status });
  }

  return NextResponse.json({ ok: true, data: res.data });
}
```

---

## 🔑 Key Concepts

### Session Token vs Metadata Cookies

| Cookie | Purpose | HttpOnly | Lifetime |
|--------|---------|----------|----------|
| `session_token` | Validate with backend | ✅ Yes | 12 hours |
| `csrf_token` | CSRF protection | ✅ Yes | 12 hours |
| `portal_role` | Middleware checks | ❌ No | 12 hours |
| `portal_user_id` | Resolve headers | ❌ No | 12 hours |
| `portal_biz_id` | Resolve headers | ❌ No | 12 hours |

**Important:** Metadata cookies must stay in sync with backend session! Call `/api/auth/session` periodically to refresh them.

### Portal Header Values

| Header | Value | When Used |
|--------|-------|-----------|
| `x-portal` | `PORTAL_SYSTEM` | `portal_role === "SYSTEM_ADMIN"` |
| `x-portal` | `PORTAL_B2B` | Any B2B role (admin, HR manager, viewer) |
| `x-business-id` | `{organisation_id}` | B2B users (empty for superadmin) |
| `x-user-id` | `{user_id}` | Always included |
| `x-tenant-id` | Default tenant | Always included |

**Casbin routing:**
- Superadmin: `system:root` domain (no org context)
- B2B users: `org:{business_id}` domain (specific org context)

### API Response Shape

Every API endpoint returns:

```typescript
{
  ok: boolean;
  message?: string;
  [data_field]?: T;  // e.g., employee, employees, organisation, etc
}
```

**Server-side:** Use `unwrapSdkResult()` or `unwrapGateway()` to unwrap the gateway envelope.

**Browser-side:** Response is already unwrapped by the API route.

---

## 🐛 Debugging Tips

### Check Session is Valid

```typescript
// In API route
const hdrs = await resolvePortalHeaders(request);
if (!hdrs) {
  console.error("No session found - user is not authenticated");
  return unauthorized();
}

console.log("Authenticated as:", hdrs.userId);
console.log("Portal context:", hdrs.portal);
console.log("Business ID:", hdrs.businessId);
```

### Trace an API Call

```typescript
// 1. Check browser cookies
console.log(document.cookie); // Should contain session_token

// 2. Check Network tab
// Request should have Cookie header with session_token
// Response should include Set-Cookie headers to refresh metadata cookies

// 3. Check Redux/state
// Session should be stored somewhere (usually context or store)

// 4. Check backend logs
// Gateway should log the session validation + auth context resolution
```

### Handle Errors Properly

```typescript
const result = await employeeClient.list();

if (!result.ok) {
  // result.message is user-facing error (from backend)
  showToast(result.message);
  
  // Log details for debugging
  console.error("API Error:", result);
}
```

### Test Superadmin Routes

1. Login as superadmin
2. Navigate to `/organisations` (superadmin-only route)
3. Check `portal_role` cookie is "SYSTEM_ADMIN"
4. Call `/api/organisations` — should succeed
5. Check Network tab — `x-portal: PORTAL_SYSTEM` header is sent

If getting 403:
- Refresh page → calls `/api/auth/session` → refreshes metadata cookies
- Check backend logs for Casbin domain resolution errors

---

## 📚 Imports Cheat Sheet

**Browser (Components/Hooks):**
```typescript
import {
  authClient,
  employeeClient,
  departmentClient,
  organisationClient,
  purchaseOrderClient,
  docgenClient,
  type ApiResult,
  type AuthResponse,
} from "@lib/sdk";
```

**Server (API Routes):**
```typescript
import { makeSdkClient, makeDirectHttp } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import {
  sdkErrorMessage,
  unwrapSdkResult,
  badRequest,
  unauthorized,
  forbidden,
  notFound,
  gatewayError,
  internalError,
} from "@lib/sdk/api-helpers";
import { parseJson, type ApiResult } from "@lib/sdk";
```

---

## ❌ Common Mistakes

### ❌ Calling makeDirectHttp in a Component
```typescript
// DON'T DO THIS
"use client";
import { makeDirectHttp } from "@lib/sdk/b2b-sdk-client"; // Server-only!
```

**Fix:** Create an API route and call it from the component.

### ❌ Forgetting to Check `hdrs`
```typescript
// DON'T DO THIS
const hdrs = await resolvePortalHeaders(request);
const sdk = makeSdkClient(request, hdrs); // Could be null!
```

**Fix:** Always check hdrs first:
```typescript
const hdrs = await resolvePortalHeaders(request);
if (!hdrs) return unauthorized();
const sdk = makeSdkClient(request, hdrs);
```

### ❌ Not Unwrapping SDK Result
```typescript
// DON'T DO THIS
const result = await sdk.listEmployees({ ... });
return NextResponse.json({ ok: true, employees: result.data });
```

**Fix:** Unwrap properly:
```typescript
const result = await sdk.listEmployees({ ... });
const unwrapped = unwrapSdkResult(result);
if (!unwrapped.ok) return gatewayError(unwrapped.message, unwrapped.status);
return NextResponse.json({ ok: true, employees: unwrapped.data?.employees });
```

### ❌ Forgetting to Forward Headers
```typescript
// DON'T DO THIS
const http = makeDirectHttp(request); // Missing hdrs!
const res = await http.post("/v1/b2b/organisations/123/admins", body);
// Will get 403 because x-portal header is not sent
```

**Fix:** Always pass hdrs:
```typescript
const hdrs = await resolvePortalHeaders(request);
if (!hdrs) return unauthorized();
const http = makeDirectHttp(request, hdrs); // Forwards x-portal, x-business-id
```

---

## 🔗 Links

- **Full Architecture:** See `ARCHITECTURE_OVERVIEW.md`
- **All Clients:** See `SDK_CLIENT_REFERENCE.md`
- **All Routes:** See `API_ROUTES_SUMMARY.md`
- **Gateway Docs:** Check @lifeplus/insuretech-sdk package

---

## Summary

1. **Components** use thin fetch clients (`authClient`, `employeeClient`, etc.)
2. **API Routes** call backend via SDK wrappers (`makeSdkClient()`)
3. **Auth** via `session_token` cookie + metadata cookies + x-portal headers
4. **Always** resolve headers first, unwrap responses properly
5. **Debug** by checking cookies, headers, and backend logs

**Got a question?** Check the docs in this order:
1. Quick Start (this file) for common patterns
2. SDK_CLIENT_REFERENCE.md for client methods
3. API_ROUTES_SUMMARY.md for endpoint details
4. ARCHITECTURE_OVERVIEW.md for deep dives

---

Created: 2024  
Last Updated: Now
