# API Routes Summary

Complete mapping of all Next.js API routes and their corresponding SDK clients.

---

## Authentication Routes

### `POST /api/auth/login`
**Purpose:** Authenticate user with mobile number + password  
**Handler:** `app/api/auth/login/route.ts`  
**Browser Client:** `authClient.login(payload)`  
**Returns:** `PortalAuthResponse` with session + cookies  
**Sets Cookies:**
- `session_token` (HttpOnly)
- `csrf_token` (HttpOnly)
- `portal_role`, `portal_user_id`, `portal_biz_id` (plain text for middleware)

---

### `POST /api/auth/logout`
**Purpose:** End current session  
**Handler:** `app/api/auth/logout/route.ts`  
**Browser Client:** `authClient.logout()`  
**Returns:** `AuthOkResponse`  
**Clears Cookies:** All auth-related cookies

---

### `GET /api/auth/session`
**Purpose:** Get current session + refresh metadata cookies  
**Handler:** `app/api/auth/session/route.ts`  
**Browser Client:** `authClient.getSession()`  
**Returns:** `PortalAuthResponse` with session  
**Re-mints Cookies:** `portal_role`, `portal_user_id`, `portal_biz_id`, `portal_email`, `portal_mobile`

**Critical:** This endpoint must be called regularly to keep metadata cookies in sync with backend session. If metadata cookies expire while session_token is valid, subsequent API calls will use PORTAL_B2B instead of PORTAL_SYSTEM for superadmin.

---

### `POST /api/auth/refresh`
**Purpose:** Refresh auth token  
**Handler:** `app/api/auth/refresh/route.ts`  
**Browser Client:** `authClient.refreshToken()`  
**Returns:** `AuthOkResponse`

---

### `GET /api/auth/profile`
**Purpose:** Get user profile  
**Handler:** `app/api/auth/profile/route.ts`  
**Browser Client:** `authClient.getProfile()`  
**Returns:** `ProfileResponse`

---

### `PATCH /api/auth/profile`
**Purpose:** Update user profile fields  
**Handler:** `app/api/auth/profile/route.ts`  
**Browser Client:** `authClient.updateProfile(payload)`  
**Returns:** `ProfileResponse`

---

### `GET /api/auth/profile-photo-url`
**Purpose:** Get S3 upload URL for profile photo  
**Handler:** `app/api/auth/profile-photo-url/route.ts`  
**Browser Client:** `authClient.getProfilePhotoUploadUrl()`  
**Returns:** `ProfilePhotoUrlResponse`

---

### `POST /api/auth/change-password`
**Purpose:** Change user password  
**Handler:** `app/api/auth/change-password/route.ts`  
**Browser Client:** `authClient.changePassword({ old_password, new_password })`  
**Returns:** `AuthOkResponse`

---

### `GET /api/auth/sessions`
**Purpose:** List active sessions for current user  
**Handler:** `app/api/auth/sessions/route.ts` (GET)  
**Browser Client:** `authClient.listSessions()`  
**Returns:** `SessionsResponse`

---

### `DELETE /api/auth/sessions`
**Purpose:** Revoke all sessions  
**Handler:** `app/api/auth/sessions/route.ts` (DELETE)  
**Browser Client:** `authClient.revokeAllSessions()`  
**Returns:** `AuthOkResponse`

---

### `DELETE /api/auth/sessions/[sessionId]`
**Purpose:** Revoke specific session  
**Handler:** `app/api/auth/sessions/[sessionId]/route.ts`  
**Browser Client:** `authClient.revokeSession(sessionId)`  
**Returns:** `AuthOkResponse`

---

### `POST /api/auth/totp`
**Purpose:** Enable two-factor authentication  
**Handler:** `app/api/auth/totp/route.ts` (POST)  
**Browser Client:** `authClient.enableTotp()`  
**Returns:** `TotpResponse`

---

### `DELETE /api/auth/totp`
**Purpose:** Disable two-factor authentication  
**Handler:** `app/api/auth/totp/route.ts` (DELETE)  
**Browser Client:** `authClient.disableTotp(totpCode)`  
**Returns:** `AuthOkResponse`

---

### `POST /api/auth/send-otp`
**Purpose:** Request SMS OTP  
**Handler:** `app/api/auth/send-otp/route.ts`  
**Browser Client:** `authClient.sendOtp(purpose?)`  
**Returns:** `AuthOkResponse`

---

### `POST /api/auth/verify-otp`
**Purpose:** Verify SMS OTP  
**Handler:** `app/api/auth/verify-otp/route.ts`  
**Browser Client:** `authClient.verifyOtp(otp, purpose?)`  
**Returns:** `OtpResponse`

---

### `POST /api/auth/send-email-otp`
**Purpose:** Request email OTP  
**Handler:** `app/api/auth/send-email-otp/route.ts`  
**Browser Client:** `authClient.sendEmailOtp(purpose?)`  
**Returns:** `AuthOkResponse`

---

### `POST /api/auth/verify-email`
**Purpose:** Verify email (via token or OTP)  
**Handler:** `app/api/auth/verify-email/route.ts`  
**Browser Client:** `authClient.verifyEmail({ token?, otp? })`  
**Returns:** `AuthOkResponse`

---

## Employee Routes

### `GET /api/employees`
**Purpose:** List employees  
**Handler:** `app/api/employees/route.ts` (GET)  
**Browser Client:** `employeeClient.list(options?)`  
**Query Params:**
```
page_size=50
offset=0
business_id=<optional>
department_id=<optional>
status=<optional>
```
**Returns:** `EmployeeListResult` with array of employees

**Server-Side SDK:** `sdk.listEmployees({ query: { ... } })`

---

### `POST /api/employees`
**Purpose:** Create employee  
**Handler:** `app/api/employees/route.ts` (POST)  
**Browser Client:** `employeeClient.create(payload)`  
**Returns:** `EmployeeSingleResult` with created employee

---

### `GET /api/employees/[id]`
**Purpose:** Get single employee with all form fields  
**Handler:** `app/api/employees/[id]/route.ts` (GET)  
**Browser Client:** `employeeClient.get(id)`  
**Returns:** `EmployeeSingleResult` with `EmployeeFullRecord`

---

### `PATCH /api/employees/[id]`
**Purpose:** Update employee  
**Handler:** `app/api/employees/[id]/route.ts` (PATCH)  
**Browser Client:** `employeeClient.update(id, payload)`  
**Returns:** `EmployeeSingleResult`

---

### `DELETE /api/employees/[id]`
**Purpose:** Delete employee  
**Handler:** `app/api/employees/[id]/route.ts` (DELETE)  
**Browser Client:** `employeeClient.delete(id)`  
**Returns:** `ApiResult`

---

## Department Routes

### `GET /api/departments`
**Purpose:** List departments  
**Handler:** `app/api/departments/route.ts` (GET)  
**Browser Client:** `departmentClient.list(pageSize?, offset?, businessId?)`  
**Query Params:**
```
page_size=50
offset=0
business_id=<optional>
```
**Returns:** `DepartmentListResult`

---

### `POST /api/departments`
**Purpose:** Create department  
**Handler:** `app/api/departments/route.ts` (POST)  
**Browser Client:** `departmentClient.create({ name, businessId })`  
**Returns:** `DepartmentSingleResult`

---

### `GET /api/departments/[id]`
**Purpose:** Get single department  
**Handler:** `app/api/departments/[id]/route.ts` (GET)  
**Browser Client:** `departmentClient.get(id)`  
**Returns:** `DepartmentSingleResult`

---

### `PATCH /api/departments/[id]`
**Purpose:** Update department  
**Handler:** `app/api/departments/[id]/route.ts` (PATCH)  
**Browser Client:** `departmentClient.update(id, { name })`  
**Returns:** `DepartmentSingleResult`

---

### `DELETE /api/departments/[id]`
**Purpose:** Delete department  
**Handler:** `app/api/departments/[id]/route.ts` (DELETE)  
**Browser Client:** `departmentClient.delete(id)`  
**Returns:** `ApiResult`

---

## Organisation Routes

### `GET /api/organisations`
**Purpose:** List organisations (superadmin only)  
**Handler:** `app/api/organisations/route.ts` (GET)  
**Browser Client:** `organisationClient.list()`  
**Auth:** Requires `SYSTEM_ADMIN` role  
**Returns:** `OrgListResult`

---

### `POST /api/organisations`
**Purpose:** Create organisation with optional admin  
**Handler:** `app/api/organisations/route.ts` (POST)  
**Browser Client:** `organisationClient.create(payload)`  
**Returns:** `OrgSingleResult`

---

### `GET /api/organisations/me`
**Purpose:** Get current user's organisation  
**Handler:** `app/api/organisations/me/route.ts`  
**Browser Client:** `organisationClient.getMe()`  
**Returns:** `OrgSingleResult`

---

### `GET /api/organisations/[id]`
**Purpose:** Get organisation details  
**Handler:** `app/api/organisations/[id]/route.ts` (GET)  
**Browser Client:** `organisationClient.get(id)`  
**Returns:** `OrgSingleResult`

---

### `PATCH /api/organisations/[id]`
**Purpose:** Update organisation  
**Handler:** `app/api/organisations/[id]/route.ts` (PATCH)  
**Browser Client:** `organisationClient.update(id, payload)`  
**Returns:** `OrgSingleResult`

---

### `DELETE /api/organisations/[id]`
**Purpose:** Delete organisation  
**Handler:** `app/api/organisations/[id]/route.ts` (DELETE)  
**Browser Client:** `organisationClient.delete(id)`  
**Returns:** `ApiResult`

---

### `GET /api/organisations/[id]/members`
**Purpose:** List organisation members  
**Handler:** `app/api/organisations/[id]/members/route.ts` (GET)  
**Browser Client:** `organisationClient.listMembers(id)`  
**Returns:** `OrgMembersResult`

---

### `POST /api/organisations/[id]/members`
**Purpose:** Add existing user as member  
**Handler:** `app/api/organisations/[id]/members/route.ts` (POST)  
**Browser Client:** `organisationClient.addMember(id, userId, role)`  
**Returns:** `OrgMemberResult`

---

### `DELETE /api/organisations/[id]/members/[memberId]`
**Purpose:** Remove member from organisation  
**Handler:** `app/api/organisations/[id]/members/[memberId]/route.ts`  
**Browser Client:** `organisationClient.removeMember(id, memberId)`  
**Returns:** `ApiResult`

---

### `POST /api/organisations/[id]/admins`
**Purpose:** Create NEW user and add as organisation admin  
**Handler:** `app/api/organisations/[id]/admins/route.ts`  
**Browser Client:** `organisationClient.createAdmin(id, payload)`  
**Returns:** `OrgMemberResult`

---

### `POST /api/organisations/[id]/assign-admin`
**Purpose:** Promote EXISTING member to admin role  
**Handler:** `app/api/organisations/[id]/assign-admin/route.ts`  
**Browser Client:** `organisationClient.assignAdmin(id, memberId)` OR `organisationClient.assignExistingAdmin(id, userId)`  
**Returns:** `OrgMemberResult`

---

### `POST /api/organisations/[id]/approve`
**Purpose:** Approve pending organisation (sets status to ACTIVE)  
**Handler:** `app/api/organisations/[id]/approve/route.ts`  
**Browser Client:** `organisationClient.approve(id)`  
**Returns:** `OrgSingleResult`

---

## Purchase Order Routes

### `GET /api/purchase-orders`
**Purpose:** List purchase orders  
**Handler:** `app/api/purchase-orders/route.ts` (GET)  
**Browser Client:** `purchaseOrderClient.list(options?)`  
**Query Params:**
```
page_size=50
offset=0
status=<optional>
```
**Returns:** `POListResult`

---

### `POST /api/purchase-orders`
**Purpose:** Create purchase order  
**Handler:** `app/api/purchase-orders/route.ts` (POST)  
**Browser Client:** `purchaseOrderClient.create(payload)`  
**Returns:** `POSingleResult`

---

### `GET /api/purchase-orders/catalog`
**Purpose:** Get available plans catalog  
**Handler:** `app/api/purchase-orders/catalog/route.ts`  
**Browser Client:** `purchaseOrderClient.getCatalog()`  
**Returns:** `POCatalogResult` with array of `CatalogItem`

---

### `GET /api/purchase-orders/[id]`
**Purpose:** Get single purchase order  
**Handler:** `app/api/purchase-orders/[id]/route.ts` (GET)  
**Browser Client:** `purchaseOrderClient.get(id)`  
**Returns:** `POSingleResult`

---

### `PATCH /api/purchase-orders/[id]`
**Purpose:** Update purchase order  
**Handler:** `app/api/purchase-orders/[id]/route.ts` (PATCH)  
**Browser Client:** `purchaseOrderClient.update(id, payload)`  
**Returns:** `POSingleResult`

---

### `DELETE /api/purchase-orders/[id]`
**Purpose:** Delete purchase order  
**Handler:** `app/api/purchase-orders/[id]/route.ts` (DELETE)  
**Browser Client:** `purchaseOrderClient.delete(id)`  
**Returns:** `ApiResult`

---

## Document Routes

### `POST /api/documents`
**Purpose:** Generate document from template  
**Handler:** `app/api/documents/route.ts` (POST)  
**Browser Client:** `docgenClient.generate(payload)`  
**Returns:** `DocumentSingleResult`

---

### `GET /api/documents`
**Purpose:** List documents for entity  
**Handler:** `app/api/documents/route.ts` (GET)  
**Browser Client:** `docgenClient.list(options)`  
**Returns:** `DocumentListResult`

---

### `GET /api/documents/[id]`
**Purpose:** Get single document  
**Handler:** `app/api/documents/[id]/route.ts` (GET)  
**Browser Client:** `docgenClient.get(documentId)`  
**Returns:** `DocumentSingleResult`

---

### `GET /api/documents/[id]/download`
**Purpose:** Download document content (base64)  
**Handler:** `app/api/documents/[id]/download/route.ts`  
**Browser Client:** `docgenClient.download(documentId)`  
**Returns:** `DocumentDownloadResult` with base64 content

---

### `DELETE /api/documents/[id]`
**Purpose:** Delete document  
**Handler:** `app/api/documents/[id]/route.ts` (DELETE)  
**Browser Client:** `docgenClient.delete(documentId)`  
**Returns:** `ApiResult`

---

## Dashboard Routes

### `GET /api/dashboard/stats`
**Purpose:** Get dashboard statistics  
**Handler:** `app/api/dashboard/stats/route.ts`  
**Returns:** Statistics for dashboard cards

---

### `GET /api/dashboard/activity`
**Purpose:** Get activity log  
**Handler:** `app/api/dashboard/activity/route.ts`  
**Returns:** Recent activity events

---

## Summary

**Total API Routes:** 45+

**Authentication:** 13 routes  
**Employees:** 5 routes  
**Departments:** 5 routes  
**Organisations:** 10 routes  
**Purchase Orders:** 6 routes  
**Documents:** 5 routes  
**Dashboard:** 2 routes  

All routes except `/api/auth/login` require a valid `session_token` cookie. Authorization is enforced by the backend using Casbin RBAC.
