# SDK Client Files Reference

## Overview

The portal's SDK layer is split into **browser-side clients** (fetch-based) and **server-side clients** (SDK wrappers). All clients route through Next.js API routes, following the BFF (Backend for Frontend) pattern.

---

## Browser-Side Clients (Used in Components)

These are used by React components and hooks to call Next.js API routes. They perform simple `fetch()` calls with JSON serialization.

### 1. `auth-client.ts` - Authentication

**Import:**
```typescript
import { authClient } from "@lib/sdk";
```

**Methods:**

| Method | Endpoint | Purpose |
|--------|----------|---------|
| `login(payload)` | `POST /api/auth/login` | Authenticate user with mobile + password |
| `logout()` | `POST /api/auth/logout` | End session |
| `getSession()` | `GET /api/auth/session` | Get current session + user profile |
| `refreshToken()` | `POST /api/auth/refresh` | Refresh auth token |
| `getProfile()` | `GET /api/auth/profile` | Get user profile |
| `updateProfile(payload)` | `PATCH /api/auth/profile` | Update profile fields |
| `getProfilePhotoUploadUrl()` | `GET /api/auth/profile-photo-url` | Get S3 upload URL |
| `changePassword(old_password, new_password)` | `POST /api/auth/change-password` | Change password |
| `listSessions()` | `GET /api/auth/sessions` | List active sessions |
| `revokeSession(sessionId)` | `DELETE /api/auth/sessions/{id}` | Revoke a session |
| `revokeAllSessions()` | `DELETE /api/auth/sessions` | Revoke all sessions |
| `enableTotp()` | `POST /api/auth/totp` | Enable 2FA |
| `disableTotp(totpCode)` | `DELETE /api/auth/totp` | Disable 2FA |
| `sendOtp(purpose)` | `POST /api/auth/send-otp` | Request SMS OTP |
| `verifyOtp(otp, purpose)` | `POST /api/auth/verify-otp` | Verify SMS OTP |
| `sendEmailOtp(purpose)` | `POST /api/auth/send-email-otp` | Request email OTP |
| `verifyEmail(token or otp)` | `POST /api/auth/verify-email` | Verify email |

**Types:**
```typescript
export type PortalLoginRequest = { mobileNumber?: string; password: string; deviceId?: string };
export type PortalAuthResponse = { ok: boolean; message?: string; session?: PortalSession };
export type ProfileResponse = { ok: boolean; message?: string; profile?: Record<string, unknown> };
export type SessionsResponse = { ok: boolean; message?: string; sessions?: Record<string, unknown> };
export type TotpResponse = { ok: boolean; message?: string; totp?: Record<string, unknown> };
export type OtpResponse = { ok: boolean; message?: string; data?: Record<string, unknown> };
```

**Usage Example:**
```typescript
const response = await authClient.login({
  mobileNumber: "+8801700000000",
  password: "SecurePass123"
});

if (response.ok) {
  console.log("Logged in:", response.session?.principal.displayName);
} else {
  console.error("Login failed:", response.message);
}
```

---

### 2. `employee-client.ts` - Employee Management

**Import:**
```typescript
import { employeeClient, type EmployeeListResult, type EmployeeSingleResult } from "@lib/sdk";
```

**Methods:**

| Method | Endpoint | Returns |
|--------|----------|---------|
| `list(options?)` | `GET /api/employees?page_size=50&offset=0&...` | `EmployeeListResult` |
| `get(id)` | `GET /api/employees/{id}` | `EmployeeSingleResult` |
| `create(payload)` | `POST /api/employees` | `EmployeeSingleResult` |
| `update(id, payload)` | `PATCH /api/employees/{id}` | `EmployeeSingleResult` |
| `delete(id)` | `DELETE /api/employees/{id}` | `ApiResult` |

**Query Options:**
```typescript
interface ListOptions {
  pageSize?: number;        // Default: 50
  offset?: number;          // Default: 0
  businessId?: string;      // Filter by org/business
  departmentId?: string;    // Filter by department
  status?: number;          // Filter by status
}
```

**Payload Types:**
```typescript
export type EmployeeCreatePayload = {
  name: string;
  employeeId: string;
  businessId: string;
  departmentId: string;
  email?: string;
  mobileNumber?: string;
  dateOfBirth?: string;
  dateOfJoining?: string;
  gender?: string;
  insuranceCategory?: number;
  assignedPlanId?: string;
  coverageAmount?: number;
  numberOfDependent?: number;
};

export type EmployeeUpdatePayload = Partial<EmployeeCreatePayload> & {
  status?: number;  // Employee status code
};

export type EmployeeFullRecord = {
  id: string;
  name: string;
  employeeID: string;
  department: string;
  insuranceCategory: number;
  assignedPlan: string;
  coverage: string;
  premiumAmount: string;
  status: "Active" | "Inactive";
  numberOfDependent: number;
  email: string;
  mobileNumber: string;
  gender: string;
  dateOfBirth: string;
  dateOfJoining: string;
  departmentId: string;
  businessId: string;
  assignedPlanId: string;
  coverageAmount: string;
};
```

**Usage Example:**
```typescript
// List employees
const { ok, employees } = await employeeClient.list({
  businessId: "org-uuid-123",
  pageSize: 20,
  offset: 0
});

// Create employee
const { ok, employee } = await employeeClient.create({
  name: "John Doe",
  employeeId: "EMP-001",
  businessId: "org-uuid-123",
  departmentId: "dept-uuid-456",
  email: "john@example.com",
  insuranceCategory: 1,  // Health
  coverageAmount: 100000,
  numberOfDependent: 2
});

// Update employee
await employeeClient.update("emp-uuid", {
  name: "Jane Doe",
  status: 1  // Active
});

// Delete employee
await employeeClient.delete("emp-uuid");
```

---

### 3. `department-client.ts` - Department Management

**Import:**
```typescript
import { departmentClient, type DepartmentListResult } from "@lib/sdk";
```

**Methods:**

| Method | Endpoint |
|--------|----------|
| `list(pageSize?, offset?, businessId?)` | `GET /api/departments?page_size=50&offset=0&business_id=...` |
| `create(name, businessId)` | `POST /api/departments` |
| `update(id, name)` | `PATCH /api/departments/{id}` |
| `delete(id)` | `DELETE /api/departments/{id}` |

**Usage Example:**
```typescript
// List all departments
const { ok, departments, total } = await departmentClient.list(50, 0, "org-uuid");

// Create department
const { ok, department } = await departmentClient.create({
  name: "Human Resources",
  businessId: "org-uuid"
});

// Update department
await departmentClient.update("dept-uuid", { name: "HR & Admin" });

// Delete department
await departmentClient.delete("dept-uuid");
```

---

### 4. `organisation-client.ts` - Organisation Management

**Import:**
```typescript
import {
  organisationClient,
  type OrgListResult,
  type OrgSingleResult,
  type OrgMembersResult,
  type OrgMember
} from "@lib/sdk";
```

**Methods:**

| Method | Endpoint | Purpose |
|--------|----------|---------|
| `list()` | `GET /api/organisations` | List all organisations |
| `get(id)` | `GET /api/organisations/{id}` | Get organisation details |
| `getMe()` | `GET /api/organisations/me` | Get current user's organisation |
| `create(payload)` | `POST /api/organisations` | Create new organisation |
| `update(id, payload)` | `PATCH /api/organisations/{id}` | Update organisation |
| `delete(id)` | `DELETE /api/organisations/{id}` | Delete organisation |
| `listMembers(id)` | `GET /api/organisations/{id}/members` | List organisation members |
| `addMember(id, userId, role)` | `POST /api/organisations/{id}/members` | Add member to org |
| `assignAdmin(id, memberId)` | `POST /api/organisations/{id}/assign-admin` | Promote existing member to admin |
| `createAdmin(id, payload)` | `POST /api/organisations/{id}/admins` | Create new admin user |
| `removeMember(id, memberId)` | `DELETE /api/organisations/{id}/members/{memberId}` | Remove member |
| `assignExistingAdmin(id, userId)` | `POST /api/organisations/{id}/assign-admin` | Assign existing user as admin |
| `approve(id)` | `POST /api/organisations/{id}/approve` | Approve pending organisation |

**Payload Types:**
```typescript
export type OrgCreatePayload = {
  name: string;
  code?: string;                        // Unique org code
  industry?: string;
  contactEmail?: string;
  contactPhone?: string;
  address?: string;
  admin?: OrgAdminCreatePayload;       // Optional: create admin during org creation
};

export type OrgUpdatePayload = Partial<OrgCreatePayload>;

export type OrgAdminCreatePayload = {
  email: string;
  password: string;
  fullName?: string;
  mobileNumber?: string;
};

export type OrgMember = {
  member_id?: string;
  organisation_id?: string;
  user_id?: string;
  role?: "ORG_MEMBER_ROLE_BUSINESS_ADMIN" | "ORG_MEMBER_ROLE_HR_MANAGER" | "ORG_MEMBER_ROLE_VIEWER";
  status?: "ORG_MEMBER_STATUS_ACTIVE" | "ORG_MEMBER_STATUS_INACTIVE";
  joined_at?: string;
  email?: string;
  full_name?: string;
  mobile_number?: string;
};
```

**Usage Example:**
```typescript
// List organisations (superadmin only)
const { ok, organisations } = await organisationClient.list();

// Get current user's organisation
const { ok, organisation } = await organisationClient.getMe();

// Create organisation with admin
const { ok, organisation } = await organisationClient.create({
  name: "Acme Corp",
  code: "ACME",
  industry: "Technology",
  contactEmail: "admin@acme.com",
  contactPhone: "+8801700000000",
  admin: {
    email: "admin@acme.com",
    password: "SecurePass123!",
    fullName: "Admin User"
  }
});

// List organisation members
const { ok, members } = await organisationClient.listMembers("org-uuid");

// Add member with HR_MANAGER role
const { ok, member } = await organisationClient.addMember(
  "org-uuid",
  "user-uuid",
  "ORG_MEMBER_ROLE_HR_MANAGER"
);

// Promote existing member to admin
const { ok } = await organisationClient.assignAdmin("org-uuid", "member-uuid");

// Remove member
await organisationClient.removeMember("org-uuid", "member-uuid");

// Approve pending organisation
const { ok } = await organisationClient.approve("org-uuid");
```

---

### 5. `purchase-order-client.ts` - Purchase Order Management

**Import:**
```typescript
import {
  purchaseOrderClient,
  type POListResult,
  type POSingleResult,
  type POCatalogResult,
  type CatalogItem
} from "@lib/sdk";
```

**Methods:**

| Method | Endpoint |
|--------|----------|
| `list(options?)` | `GET /api/purchase-orders?page_size=50&offset=0&status=...` |
| `get(id)` | `GET /api/purchase-orders/{id}` |
| `create(payload)` | `POST /api/purchase-orders` |
| `update(id, payload)` | `PATCH /api/purchase-orders/{id}` |
| `delete(id)` | `DELETE /api/purchase-orders/{id}` |
| `getCatalog()` | `GET /api/purchase-orders/catalog` |

**Payload Types:**
```typescript
export type PurchaseOrderCreatePayload = {
  departmentId: string;
  planId: string;
  insuranceCategory?: string;    // e.g. "Health", "Life"
  employeeCount: number;
  numberOfDependents?: number;
  coverageAmount?: number;
  notes?: string;
};

export type PurchaseOrderUpdatePayload = {
  status?: number;
  notes?: string;
  employeeCount?: number;
  numberOfDependents?: number;
  coverageAmount?: number;
};

export type CatalogItem = {
  planId: string;
  productId: string;
  productName: string;
  planName: string;
  insuranceCategory: string;
  premiumAmount: string;
};
```

**Usage Example:**
```typescript
// Get catalog of available plans
const { ok, items } = await purchaseOrderClient.getCatalog();
// items: [{ planId: "...", productName: "Health Plan A", ... }]

// Create purchase order
const { ok, purchaseOrder } = await purchaseOrderClient.create({
  departmentId: "dept-uuid",
  planId: "plan-uuid",
  insuranceCategory: "Health",
  employeeCount: 50,
  numberOfDependents: 75,
  coverageAmount: 100000,
  notes: "Coverage for IT department"
});

// List purchase orders
const { ok, purchaseOrders } = await purchaseOrderClient.list({
  pageSize: 20,
  offset: 0,
  status: 1  // Active
});

// Update purchase order
await purchaseOrderClient.update("po-uuid", {
  employeeCount: 60,
  notes: "Updated headcount"
});

// Delete purchase order
await purchaseOrderClient.delete("po-uuid");
```

---

### 6. `docgen-client.ts` - Document Generation

**Import:**
```typescript
import {
  docgenClient,
  type DocumentRecord,
  type DocumentListResult,
  type GenerateDocumentPayload
} from "@lib/sdk";
```

**Methods:**

| Method | Endpoint |
|--------|----------|
| `generate(payload)` | `POST /api/documents` |
| `list(options)` | `GET /api/documents?entity_type=...&entity_id=...` |
| `get(documentId)` | `GET /api/documents/{id}` |
| `download(documentId)` | `GET /api/documents/{id}/download` |
| `delete(documentId)` | `DELETE /api/documents/{id}` |

**Payload Types:**
```typescript
export type DocumentStatus =
  | "DOCUMENT_STATUS_PENDING"
  | "DOCUMENT_STATUS_PROCESSING"
  | "DOCUMENT_STATUS_COMPLETED"
  | "DOCUMENT_STATUS_FAILED";

export interface DocumentRecord {
  document_id: string;
  template_id: string;
  entity_type: string;
  entity_id: string;
  document_type: string;
  status: DocumentStatus;
  file_url?: string;
  download_url?: string;
  created_at?: string;
  updated_at?: string;
}

export interface GenerateDocumentPayload {
  template_id: string;
  entity_type: string;        // e.g. "employee", "policy"
  entity_id: string;
  data?: Record<string, unknown>;  // Variables to inject into template
  include_qr_code?: boolean;
}
```

**Usage Example:**
```typescript
// Generate document
const { ok, document } = await docgenClient.generate({
  template_id: "tpl-uuid",
  entity_type: "employee",
  entity_id: "emp-uuid",
  data: {
    employee_name: "John Doe",
    policy_number: "POL-123456"
  },
  include_qr_code: true
});

// List documents for an employee
const { ok, documents } = await docgenClient.list({
  entityType: "employee",
  entityId: "emp-uuid",
  status: "DOCUMENT_STATUS_COMPLETED",
  pageSize: 10
});

// Get single document
const { ok, document } = await docgenClient.get("doc-uuid");

// Download document (returns base64 content)
const { ok, content, content_type, file_name } = await docgenClient.download("doc-uuid");

// Delete document
await docgenClient.delete("doc-uuid");
```

---

## Server-Side Clients (Used in API Routes)

These are used only in Next.js API route handlers to call the backend gateway. They wrap the `@lifeplus/insuretech-sdk` library.

### 7. `b2b-sdk-client.ts` - Main SDK Wrapper

**Import:**
```typescript
import { makeSdkClient, makeDirectHttp } from "@lib/sdk/b2b-sdk-client";
import type { B2bSdkClient } from "@lib/sdk/b2b-sdk-client";
```

**Factory Function:**
```typescript
export function makeSdkClient(
  request: Request,
  sessionOverrides?: {
    portal?: string;      // x-portal header
    userId?: string;      // x-user-id header
    businessId?: string;  // x-business-id header
    tenantId?: string;    // x-tenant-id header
  }
): B2bSdkClient
```

**Returned Methods:**

All methods return `{ data: T, error: Error, response: Response }` from the SDK.

#### Auth Methods
```typescript
emailLogin(opts: { body: { email, password, device_id? } })
logout(opts?)
validateToken(opts: { body: { token } })
registerEmailUser(opts: { body: { email, password, full_name? } })
getCurrentSession(opts?)
refreshToken(opts?)
changePassword(opts: { body: { old_password, new_password } })
getUserProfile(opts?)
updateUserProfile(opts: { body: { full_name?, mobile_number?, ... } })
getProfilePhotoUploadUrl(opts?)
listSessions(opts?)
revokeSession(opts: { path_params: { session_id } })
revokeAllSessions(opts?)
enableTotp(opts?)
disableTotp(opts: { body: { totp_code } })
sendOtp(opts: { body: { purpose?, mobile_number? } })
verifyOtp(opts: { body: { otp, purpose? } })
sendEmailOtp(opts: { body: { purpose? } })
verifyEmail(opts: { body: { token?, otp? } })
```

#### Employee Methods
```typescript
listEmployees(opts?: { query: { page_size?, business_id?, department_id? } })
createEmployee(opts: { body: { name, employee_id, business_id, ... } })
getEmployee(opts: { path_params: { id } })
updateEmployee(opts: { path_params: { id }, body: { name?, ... } })
deleteEmployee(opts: { path_params: { id } })
```

#### Department Methods
```typescript
listDepartments(opts?: { query: { page_size?, business_id? } })
createDepartment(opts: { body: { name, business_id } })
getDepartment(opts: { path_params: { id } })
updateDepartment(opts: { path_params: { id }, body: { name } })
deleteDepartment(opts: { path_params: { id } })
```

#### Purchase Order Methods
```typescript
listPurchaseOrders(opts?: { query: { page_size?, business_id? } })
createPurchaseOrder(opts: { body: { department_id, plan_id, employee_count, ... } })
getPurchaseOrder(opts: { path_params: { id } })
listPurchaseOrderCatalog(opts?: { query: { page_size? } })
updatePurchaseOrderHttp(id: string, body: Record<string, unknown>)  // Direct HTTP
deletePurchaseOrderHttp(id: string)                                 // Direct HTTP
```

#### Organisation Methods
```typescript
listOrganisations(opts?: { query: { tenant_id, page_size? } })
createOrganisation(opts: { body: { tenant_id, name, code, industry?, ... } })
getOrganisation(opts: { path_params: { id } })
updateOrganisation(opts: { path_params: { id }, body: { name?, ... } })
deleteOrganisation(opts: { path_params: { id } })
listOrgMembers(opts: { path_params: { id }, query: { page_size? } })
addOrgMember(opts: { path_params: { id }, body: { user_id, role } })
assignOrgAdmin(opts: { path_params: { id }, body: { member_id } })
removeOrgMember(opts: { path_params: { id, member_id } })
```

**Usage in API Route:**
```typescript
export async function GET(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false }, { status: 401 });
  
  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.listEmployees({
    query: { page_size: 50, business_id: hdrs.businessId }
  });
  
  if (!result.response.ok) {
    return NextResponse.json(
      { ok: false, message: sdkErrorMessage(result) },
      { status: result.response.status }
    );
  }
  
  return NextResponse.json({ ok: true, employees: result.data?.employees ?? [] });
}
```

---

### 8. `makeDirectHttp()` - Raw HTTP Client

**Import:**
```typescript
import { makeDirectHttp } from "@lib/sdk/b2b-sdk-client";
```

**Factory Function:**
```typescript
export function makeDirectHttp(
  request: Request,
  sessionOverrides?: PortalHeaders
) {
  return {
    get(path: string): Promise<HttpResult>
    post(path: string, body?: unknown): Promise<HttpResult>
    patch(path: string, body?: unknown): Promise<HttpResult>
    put(path: string, body?: unknown): Promise<HttpResult>
    delete(path: string): Promise<HttpResult>
  };
}
```

**Return Type:**
```typescript
interface HttpResult {
  ok: boolean;
  status: number;
  data: Record<string, unknown>;
  error: GatewayError | null;
  message?: string;
}
```

**Usage for Endpoints NOT in Generated SDK:**
```typescript
export async function POST(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false }, { status: 401 });
  
  const http = makeDirectHttp(request, hdrs);
  const body = await request.json();
  
  // Call gateway endpoint not exposed as SDK method
  const res = await http.post(`/v1/b2b/organisations/${id}/admins`, body);
  
  if (!res.ok) {
    return NextResponse.json({ ok: false, message: res.message }, { status: res.status });
  }
  
  return NextResponse.json({ ok: true, member: res.data });
}
```

---

### 9. `docgen-sdk-client.ts` - Document Service Client

**Import:**
```typescript
import { makeDocgenClient } from "@lib/sdk/docgen-sdk-client";
import type {
  GenerateDocumentPayload,
  GenerateDocumentResponse,
  ListDocumentsResponse
} from "@lib/sdk/docgen-sdk-client";
```

**Factory Function:**
```typescript
export function makeDocgenClient(
  request: Request,
  sessionOverrides?: PortalHeaders
) {
  return {
    generate(payload: GenerateDocumentPayload): Promise<DocumentSingleResult>
    list(path: string, query: Record<string, unknown>): Promise<DocumentListResult>
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

**Usage Example:**
```typescript
export async function POST(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false }, { status: 401 });
  
  const docgen = makeDocgenClient(request, hdrs);
  const body = await request.json();
  
  const result = await docgen.generate(body);
  
  if (!result.ok) {
    return NextResponse.json({ ok: false, message: result.message }, { status: result.status });
  }
  
  return NextResponse.json({ ok: true, document: result.document });
}
```

---

## Shared Utilities

### `shared.ts` - Universal Types & Helpers

**Import:**
```typescript
import { parseJson, type ApiResult, type JsonMap } from "@lib/sdk";
```

**Types:**
```typescript
export type ApiResult<T extends object = object> = {
  ok: boolean;
  message?: string;
} & T;

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
```

**Functions:**
```typescript
export async function parseJson<T>(response: Response): Promise<T>
export function unwrapGateway<T>(body: GatewayResponse<T>, httpStatus?: number)
export function extractGatewayError(result: unknown): string
```

---

### `api-helpers.ts` - Server-Side Response Builders

**Import:**
```typescript
import {
  badRequest,
  gatewayError,
  notFound,
  unauthorized,
  forbidden,
  internalError,
  sdkErrorMessage,
  unwrapSdkResult,
  getApiBaseUrl,
  getCookieValue,
  getCsrfToken,
  parseMoneyDecimal
} from "@lib/sdk/api-helpers";
```

**Response Builders:**
```typescript
badRequest(message: string): NextResponse        // 400
gatewayError(message, status?): NextResponse     // 502 default
notFound(message?): NextResponse                 // 404
unauthorized(message?): NextResponse             // 401
forbidden(message?): NextResponse                // 403
internalError(message?): NextResponse            // 500
```

**Error Extraction:**
```typescript
sdkErrorMessage(result: unknown): string         // Extract error from SDK result
unwrapSdkResult(result): UnwrappedResult         // Unwrap GatewayResponse
extractGatewayError(result: unknown): string     // Extract from any shape
```

**Utilities:**
```typescript
getApiBaseUrl(): string                          // Get gateway base URL
getCookieValue(cookieHeader, name): string       // Extract cookie value
getCsrfToken(cookieHeader): string               // Get CSRF token
parseMoneyDecimal(value): number                 // Parse money to decimal
```

---

## Best Practices

1. **Always resolve headers first in API routes:**
   ```typescript
   const hdrs = await resolvePortalHeaders(request);
   if (!hdrs) return NextResponse.json({ ok: false }, { status: 401 });
   ```

2. **Use unwrapSdkResult for SDK calls:**
   ```typescript
   const result = await sdk.listEmployees({ ... });
   const unwrapped = unwrapSdkResult(result);
   if (!unwrapped.ok) return NextResponse.json({ ok: false, ... });
   ```

3. **Use sdkErrorMessage for error extraction:**
   ```typescript
   const message = sdkErrorMessage(result);
   ```

4. **Import only what you need from shared types:**
   ```typescript
   import type { ApiResult } from "@lib/sdk";
   ```

5. **Never call makeDirectHttp in browser code** — it's server-only!

6. **Always include `cache: "no-store"` in fetch calls** to avoid stale data.

7. **Forward session headers via sessionOverrides:**
   ```typescript
   const sdk = makeSdkClient(request, hdrs);
   const http = makeDirectHttp(request, hdrs);
   ```

---

## Summary

| File | Type | Usage | Export |
|------|------|-------|--------|
| `auth-client.ts` | Browser | Components | `authClient` |
| `employee-client.ts` | Browser | Components | `employeeClient` |
| `department-client.ts` | Browser | Components | `departmentClient` |
| `organisation-client.ts` | Browser | Components | `organisationClient` |
| `purchase-order-client.ts` | Browser | Components | `purchaseOrderClient` |
| `docgen-client.ts` | Browser | Components | `docgenClient` |
| `b2b-sdk-client.ts` | Server | API routes | `makeSdkClient`, `makeDirectHttp` |
| `docgen-sdk-client.ts` | Server | API routes | `makeDocgenClient` |
| `shared.ts` | Both | Everywhere | Types + utilities |
| `api-helpers.ts` | Server | API routes | Response builders |
| `session-headers.ts` | Server | API routes | `resolvePortalHeaders` |

All clients are exported through the barrel `src/lib/sdk/index.ts` for convenient imports.
