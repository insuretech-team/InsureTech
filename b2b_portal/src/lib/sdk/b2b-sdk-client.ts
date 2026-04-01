/**
 * b2b-sdk-client.ts
 * ─────────────────
 * Server-side SDK client factory for B2B portal Next.js API route handlers.
 *
 * The SDK (@lifeplus/insuretech-sdk) is auto-generated from the protobuf
 * service definitions via the API pipeline script.
 *
 * Authentication: The portal uses cookie-based server-side sessions.
 * The gateway validates the session cookie and injects business_id / user_id
 * from the JWT into every downstream gRPC call — the portal does NOT need to
 * pass those values. apiKey is a required config field but auth is handled by
 * the forwarded session cookie.
 *
 * SDK methods NOT generated (PO update/delete, resolveMyOrg):
 * Those RPCs are gateway-only and not exposed as REST endpoints, so they are
 * handled by direct authenticated HTTP inside their route files.
 *
 * Usage:
 *   import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
 *   const sdk = makeSdkClient(req);
 *   const { data, error, response } = await sdk.listEmployees({ query: { page_size: 50 } });
 */

import {
  // ── Auth ────────────────────────────────────────────────────────────────────
  authServiceLogin,
  authServiceEmailLogin,
  authServiceLogout,
  authServiceRegisterEmailUser,
  authServiceValidateToken,
  authServiceGetCurrentSession,
  authServiceRefreshToken,
  authServiceChangePassword,
  authServiceGetUserProfile,
  authServiceUpdateUserProfile,
  authServiceGetProfilePhotoUploadUrl,
  authServiceListSessions,
  authServiceRevokeSession,
  authServiceRevokeAllSessions,
  authServiceEnableTotp,
  authServiceDisableTotp,
  authServiceSendOtp,
  authServiceVerifyOtp,
  authServiceSendEmailOtp,
  authServiceVerifyEmail,
  createInsureTechClient,
  b2bServiceListEmployees,
  b2bServiceCreateEmployee,
  b2bServiceGetEmployee,
  b2bServiceUpdateEmployee,
  b2bServiceDeleteEmployee,
  b2bServiceListDepartments,
  b2bServiceCreateDepartment,
  b2bServiceGetDepartment,
  b2bServiceUpdateDepartment,
  b2bServiceDeleteDepartment,
  b2bServiceListPurchaseOrders,
  b2bServiceCreatePurchaseOrder,
  b2bServiceGetPurchaseOrder,
  b2bServiceListPurchaseOrderCatalog,
  b2bServiceListOrganisations,
  b2bServiceCreateOrganisation,
  b2bServiceDeleteOrganisation,
  b2bServiceGetOrganisation,
  b2bServiceUpdateOrganisation,
  b2bServiceListOrgMembers,
  b2bServiceAddOrgMember,
  b2bServiceAssignOrgAdmin,
  b2bServiceRemoveOrgMember,
  paymentServiceListPayments,
  productServiceListProducts,
  billingServiceListInvoices,
  

} from "@lifeplus/insuretech-sdk";

// ─── Helpers ──────────────────────────────────────────────────────────────────

function getBaseUrl(): string {
  return (
    process.env.INSURETECH_API_BASE_URL ??
    process.env.NEXT_PUBLIC_INSURETECH_API_BASE_URL ??
    "http://localhost:8080"
  );
}

function extractCsrf(cookieHeader: string): string {
  const m = cookieHeader.match(/(?:^|;\s*)csrf_token=([^;]*)/);
  return m ? decodeURIComponent(m[1]) : "";
}

// ─── Factory ──────────────────────────────────────────────────────────────────

export function makeSdkClient(request: Request, sessionOverrides?: { portal?: string; userId?: string; businessId?: string; tenantId?: string }) {
  const cookieHeader = request.headers.get("cookie") ?? "";
  const csrf = extractCsrf(cookieHeader);

  const extraHeaders: Record<string, string> = {};
  if (cookieHeader) extraHeaders["cookie"] = cookieHeader;
  if (csrf) extraHeaders["X-CSRF-Token"] = csrf;

  // Forward portal + business-id headers so the backend authz interceptor
  // can correctly resolve the Casbin domain.
  // Super admin: x-portal=PORTAL_SYSTEM (no x-business-id needed)
  // B2B admin:   x-portal=PORTAL_B2B + x-business-id={org_id}
  // Priority: sessionOverrides (from server session store) > request headers (browser-forwarded)
  const xPortal = sessionOverrides?.portal ?? request.headers.get("x-portal") ?? "";
  const xBusinessId = sessionOverrides?.businessId ?? request.headers.get("x-business-id") ?? "";
  const xUserId = sessionOverrides?.userId ?? request.headers.get("x-user-id") ?? "";
  const xTenantId = sessionOverrides?.tenantId ?? request.headers.get("x-tenant-id") ?? "";
  if (xPortal) extraHeaders["x-portal"] = xPortal;
  if (xBusinessId) extraHeaders["x-business-id"] = xBusinessId;
  if (xUserId) extraHeaders["x-user-id"] = xUserId;
  if (xTenantId) extraHeaders["x-tenant-id"] = xTenantId;

  // apiKey is required by InsureTechClientConfig but auth is done via cookie.
  // The gateway validates the session cookie — apiKey is unused by the backend.
  const sdkClient = createInsureTechClient({
    apiKey: process.env.INSURETECH_API_KEY ?? "",
    baseUrl: getBaseUrl(),
    headers: extraHeaders,
  });

  return {
    // ── Auth ────────────────────────────────────────────────────────────────
    mobileLogin: (opts: Omit<Parameters<typeof authServiceLogin>[0], "client">) =>
      authServiceLogin({ client: sdkClient, throwOnError: false, ...opts }),

    emailLogin: (opts: Omit<Parameters<typeof authServiceEmailLogin>[0], "client">) =>
      authServiceEmailLogin({ client: sdkClient, throwOnError: false, ...opts }),

    logout: (opts?: Omit<Parameters<typeof authServiceLogout>[0], "client">) =>
      authServiceLogout({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    validateToken: (opts: Omit<Parameters<typeof authServiceValidateToken>[0], "client">) =>
      authServiceValidateToken({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    registerEmailUser: (opts: Omit<Parameters<typeof authServiceRegisterEmailUser>[0], "client">) =>
      authServiceRegisterEmailUser({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    getCurrentSession: (opts?: Omit<Parameters<typeof authServiceGetCurrentSession>[0], "client">) =>
      authServiceGetCurrentSession({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    refreshToken: (opts?: Omit<Parameters<typeof authServiceRefreshToken>[0], "client">) =>
      authServiceRefreshToken({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    changePassword: (opts: Omit<Parameters<typeof authServiceChangePassword>[0], "client">) =>
      authServiceChangePassword({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    getUserProfile: (opts?: Omit<Parameters<typeof authServiceGetUserProfile>[0], "client">) =>
      authServiceGetUserProfile({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    updateUserProfile: (opts: Omit<Parameters<typeof authServiceUpdateUserProfile>[0], "client">) =>
      authServiceUpdateUserProfile({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    getProfilePhotoUploadUrl: (opts?: Omit<Parameters<typeof authServiceGetProfilePhotoUploadUrl>[0], "client">) =>
      authServiceGetProfilePhotoUploadUrl({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    listSessions: (opts?: Omit<Parameters<typeof authServiceListSessions>[0], "client">) =>
      authServiceListSessions({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    revokeSession: (opts: Omit<Parameters<typeof authServiceRevokeSession>[0], "client">) =>
      authServiceRevokeSession({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    revokeAllSessions: (opts?: Omit<Parameters<typeof authServiceRevokeAllSessions>[0], "client">) =>
      authServiceRevokeAllSessions({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    enableTotp: (opts?: Omit<Parameters<typeof authServiceEnableTotp>[0], "client">) =>
      authServiceEnableTotp({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    disableTotp: (opts: Omit<Parameters<typeof authServiceDisableTotp>[0], "client">) =>
      authServiceDisableTotp({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    sendOtp: (opts: Omit<Parameters<typeof authServiceSendOtp>[0], "client">) =>
      authServiceSendOtp({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    verifyOtp: (opts: Omit<Parameters<typeof authServiceVerifyOtp>[0], "client">) =>
      authServiceVerifyOtp({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    sendEmailOtp: (opts: Omit<Parameters<typeof authServiceSendEmailOtp>[0], "client">) =>
      authServiceSendEmailOtp({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    verifyEmail: (opts: Omit<Parameters<typeof authServiceVerifyEmail>[0], "client">) =>
      authServiceVerifyEmail({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    // ── Employees ──────────────────────────────────────────────────────────
    listEmployees: (opts?: Omit<Parameters<typeof b2bServiceListEmployees>[0], "client">) =>
      b2bServiceListEmployees({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    createEmployee: (opts: Omit<Parameters<typeof b2bServiceCreateEmployee>[0], "client">) =>
      b2bServiceCreateEmployee({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    getEmployee: (opts: Omit<Parameters<typeof b2bServiceGetEmployee>[0], "client">) =>
      b2bServiceGetEmployee({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    updateEmployee: (opts: Omit<Parameters<typeof b2bServiceUpdateEmployee>[0], "client">) =>
      b2bServiceUpdateEmployee({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    deleteEmployee: (opts: Omit<Parameters<typeof b2bServiceDeleteEmployee>[0], "client">) =>
      b2bServiceDeleteEmployee({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    // ── Departments ────────────────────────────────────────────────────────
    listDepartments: (opts?: Omit<Parameters<typeof b2bServiceListDepartments>[0], "client">) =>
      b2bServiceListDepartments({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    createDepartment: (opts: Omit<Parameters<typeof b2bServiceCreateDepartment>[0], "client">) =>
      b2bServiceCreateDepartment({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    getDepartment: (opts: Omit<Parameters<typeof b2bServiceGetDepartment>[0], "client">) =>
      b2bServiceGetDepartment({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    updateDepartment: (opts: Omit<Parameters<typeof b2bServiceUpdateDepartment>[0], "client">) =>
      b2bServiceUpdateDepartment({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    deleteDepartment: (opts: Omit<Parameters<typeof b2bServiceDeleteDepartment>[0], "client">) =>
      b2bServiceDeleteDepartment({ client: sdkClient, throwOnError: false, ...(opts as any) }),
    // ── Payments ───────────────────────────────────────────────────────────
    listPayments: (opts?: Omit<Parameters<typeof paymentServiceListPayments>[0], "client">) =>
      paymentServiceListPayments({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    // ── Billing ────────────────────────────────────────────────────────────
    listInvoices: (opts?: Omit<Parameters<typeof billingServiceListInvoices>[0], "client">) =>
      billingServiceListInvoices({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    // ── Products ───────────────────────────────────────────────────────────
    listProducts: (opts?: Omit<Parameters<typeof productServiceListProducts>[0], "client">) =>
      productServiceListProducts({ client: sdkClient, throwOnError: false, ...(opts as any) }),
    
    
         // ── Purchase Orders ────────────────────────────────────────────────────
    // NOTE: UpdatePurchaseOrder / DeletePurchaseOrder are not exposed as REST
    // endpoints in the generated SDK. Those operations fall back to direct HTTP.
    listPurchaseOrders: (opts?: Omit<Parameters<typeof b2bServiceListPurchaseOrders>[0], "client">) =>
      b2bServiceListPurchaseOrders({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    createPurchaseOrder: (opts: Omit<Parameters<typeof b2bServiceCreatePurchaseOrder>[0], "client">) =>
      b2bServiceCreatePurchaseOrder({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    getPurchaseOrder: (opts: Omit<Parameters<typeof b2bServiceGetPurchaseOrder>[0], "client">) =>
      b2bServiceGetPurchaseOrder({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    listPurchaseOrderCatalog: (opts?: Omit<Parameters<typeof b2bServiceListPurchaseOrderCatalog>[0], "client">) =>
      b2bServiceListPurchaseOrderCatalog({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    // ── Organisations ──────────────────────────────────────────────────────
    listOrganisations: (opts?: Omit<Parameters<typeof b2bServiceListOrganisations>[0], "client">) =>
      b2bServiceListOrganisations({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    createOrganisation: (opts: Omit<Parameters<typeof b2bServiceCreateOrganisation>[0], "client">) =>
      b2bServiceCreateOrganisation({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    deleteOrganisation: (opts: Omit<Parameters<typeof b2bServiceDeleteOrganisation>[0], "client">) =>
      b2bServiceDeleteOrganisation({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    getOrganisation: (opts: Omit<Parameters<typeof b2bServiceGetOrganisation>[0], "client">) =>
      b2bServiceGetOrganisation({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    updateOrganisation: (opts: Omit<Parameters<typeof b2bServiceUpdateOrganisation>[0], "client">) =>
      b2bServiceUpdateOrganisation({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    listOrgMembers: (opts: Omit<Parameters<typeof b2bServiceListOrgMembers>[0], "client">) =>
      b2bServiceListOrgMembers({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    addOrgMember: (opts: Omit<Parameters<typeof b2bServiceAddOrgMember>[0], "client">) =>
      b2bServiceAddOrgMember({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    assignOrgAdmin: (opts: Omit<Parameters<typeof b2bServiceAssignOrgAdmin>[0], "client">) =>
      b2bServiceAssignOrgAdmin({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    removeOrgMember: (opts: Omit<Parameters<typeof b2bServiceRemoveOrgMember>[0], "client">) =>
      b2bServiceRemoveOrgMember({ client: sdkClient, throwOnError: false, ...(opts as any) }),

    // ── Direct HTTP for SDK-missing operations ─────────────────────────────
    // Delegate to makeDirectHttp so there are zero hardcoded fetch() calls here.
    /** PATCH /v1/b2b/purchase-orders/{id} — not yet a generated SDK method */
    updatePurchaseOrderHttp: (id: string, body: Record<string, unknown>) =>
      makeDirectHttp(request, sessionOverrides).patch(`/v1/b2b/purchase-orders/${id}`, body),

    /** DELETE /v1/b2b/purchase-orders/{id} — not yet a generated SDK method */
    deletePurchaseOrderHttp: (id: string) =>
      makeDirectHttp(request, sessionOverrides).delete(`/v1/b2b/purchase-orders/${id}`),

    /** GET /v1/b2b/organisations/me — resolve the caller's own organisation */
    getMyOrganisation: () =>
      makeDirectHttp(request, sessionOverrides).get(`/v1/b2b/organisations/me`),
  };
}

export type B2bSdkClient = ReturnType<typeof makeSdkClient>;

// ─────────────────────────────────────────────────────────────────────────────
// BROWSER-SIDE BFF CLIENT
// Thin fetch wrappers that call Next.js /api/* BFF routes from the browser.
// Components and hooks MUST use these — never call the gateway directly.
// ─────────────────────────────────────────────────────────────────────────────

import { parseJson, type ApiResult } from "./shared";
import type { Employee, Department, Organisation, PurchaseOrder } from "@lib/types/b2b";

// ─── Employee types ───────────────────────────────────────────────────────────
export type EmployeeCreatePayload = {
  name: string; employeeId: string; businessId: string; departmentId: string;
  email?: string; mobileNumber?: string; dateOfBirth?: string; dateOfJoining?: string;
  gender?: string; insuranceCategory?: number; assignedPlanId?: string;
  coverageAmount?: number; numberOfDependent?: number;
};
export type EmployeeUpdatePayload = Partial<EmployeeCreatePayload> & { status?: number };
export type EmployeeFullRecord = {
  id: string; name: string; employeeID: string; department: string;
  insuranceCategory: number; assignedPlan: string; coverage: string;
  premiumAmount: string; status: "Active" | "Inactive"; numberOfDependent: number;
  email: string; mobileNumber: string; gender: string;
  dateOfBirth: string; dateOfJoining: string;
  departmentId: string; businessId: string;
  assignedPlanId: string; coverageAmount: string;
};
export type EmployeeListResult = ApiResult<{ employees: Employee[]; total?: number }>;
export type EmployeeSingleResult = ApiResult<{ employee?: EmployeeFullRecord }>;

// ─── Department types ─────────────────────────────────────────────────────────
export type DepartmentListResult = ApiResult<{ departments: Department[]; total?: number }>;
export type DepartmentSingleResult = ApiResult<{ department?: Record<string, unknown> }>;

// ─── Organisation types ───────────────────────────────────────────────────────
export type OrgMemberRole =
  | "ORG_MEMBER_ROLE_UNSPECIFIED" | "ORG_MEMBER_ROLE_BUSINESS_ADMIN"
  | "ORG_MEMBER_ROLE_HR_MANAGER" | "ORG_MEMBER_ROLE_VIEWER";
export type OrgMemberStatus =
  | "ORG_MEMBER_STATUS_UNSPECIFIED" | "ORG_MEMBER_STATUS_ACTIVE" | "ORG_MEMBER_STATUS_INACTIVE";
export type OrgMember = {
  member_id?: string; organisation_id?: string; user_id?: string;
  role?: OrgMemberRole; status?: OrgMemberStatus; joined_at?: string;
  email?: string; full_name?: string; mobile_number?: string;
};
export type OrgCreatePayload = {
  name: string; code?: string; industry?: string;
  contactEmail?: string; contactPhone?: string; address?: string;
  admin?: OrgAdminCreatePayload;
  adminAssignment?: OrgAdminAssignmentPayload;
};
export type OrgUpdatePayload = Partial<OrgCreatePayload>;
export type OrgAdminCreatePayload = {
  email: string; password: string; fullName?: string; mobileNumber?: string;
};
export type OrgAdminAssignmentPayload = {
  userId: string; temporaryPassword: string;
};
export type OrgListResult = ApiResult<{ organisations: Organisation[] }>;
export type OrgSingleResult = ApiResult<{ organisation?: Organisation }>;
export type OrgMembersResult = ApiResult<{ members: OrgMember[] }>;
export type OrgMemberResult = ApiResult<{ member?: OrgMember }>;
export type PortalUserLookup = {
  userId: string;
  fullName?: string;
  email?: string;
  mobileNumber?: string;
  userType?: string;
  emailVerified?: boolean;
  kycVerified?: boolean;
  passwordChangeRequired?: boolean;
};
export type PortalUserLookupResult = ApiResult<{ user?: PortalUserLookup }>;

// ─── Purchase Order types ─────────────────────────────────────────────────────
export type CatalogItem = {
  planId: string; productId: string; productName: string;
  planName: string; insuranceCategory: string; premiumAmount: string;
};
export type PurchaseOrderCreatePayload = {
  departmentId: string; planId: string; insuranceCategory?: string;
  employeeCount: number; numberOfDependents?: number;
  coverageAmount?: number; notes?: string;
};
export type PurchaseOrderUpdatePayload = {
  status?: number; notes?: string; employeeCount?: number;
  numberOfDependents?: number; coverageAmount?: number;
};
export type POListResult = ApiResult<{ purchaseOrders: PurchaseOrder[]; total?: number }>;
export type POSingleResult = ApiResult<{ purchaseOrder?: PurchaseOrder | null }>;
export type POCatalogResult = ApiResult<{ items: CatalogItem[] }>;

// ─── Docgen types (browser-facing) ───────────────────────────────────────────
export type DocumentStatus =
  | "DOCUMENT_STATUS_PENDING" | "DOCUMENT_STATUS_PROCESSING"
  | "DOCUMENT_STATUS_COMPLETED" | "DOCUMENT_STATUS_FAILED";
export interface DocumentRecord {
  document_id: string; template_id: string; entity_type: string;
  entity_id: string; document_type: string; status: DocumentStatus;
  file_url?: string; download_url?: string; created_at?: string; updated_at?: string;
}
export interface GenerateDocumentPayload {
  template_id: string; entity_type: string; entity_id: string;
  data?: Record<string, unknown>; include_qr_code?: boolean;
}
export type DocumentListResult = ApiResult<{ documents?: DocumentRecord[]; total?: number }>;
export type DocumentSingleResult = ApiResult<{ document?: DocumentRecord }>;
export type DocumentDownloadResult = ApiResult<{
  content?: string; content_type?: string; file_name?: string;
}>;

// ─── Auth BFF types ───────────────────────────────────────────────────────────
export type AuthOkResponse = { ok: boolean; message?: string };
export type ProfileResponse = { ok: boolean; message?: string; profile?: Record<string, unknown> };
export type SessionsResponse = { ok: boolean; message?: string; sessions?: Record<string, unknown> };
export type TotpResponse = { ok: boolean; message?: string; totp?: Record<string, unknown> };
export type OtpResponse = { ok: boolean; message?: string; data?: Record<string, unknown> };
export type ProfilePhotoUrlResponse = { ok: boolean; message?: string; uploadUrl?: Record<string, unknown> };
export type EmployeeLoginOrganisation = {
  organisationId: string;
  organisationName: string;
  organisationCode: string;
};
export type EmployeeLoginOrganisationsResult = ApiResult<{ organisations: EmployeeLoginOrganisation[] }>;
export type EmployeeActivationResult = ApiResult<{ otpId?: string; expiresInSeconds?: number; userId?: string }>;
export type EmployeeCoverageResult = ApiResult<{ coverage?: Record<string, unknown> }>;

// ─────────────────────────────────────────────────────────────────────────────
// bffClient — the single browser-side client. Import this in components/hooks.
// All calls go to /api/* BFF routes. Never to the gateway directly.
// ─────────────────────────────────────────────────────────────────────────────
import type {
  EmployeePortalLoginRequest,
  PortalAuthResponse,
  PortalLoginRequest,
} from "@lib/types/auth";

export const bffClient = {
  // ── Auth ──────────────────────────────────────────────────────────────────
  auth: {
    async login(payload: PortalLoginRequest): Promise<PortalAuthResponse> {
      return parseJson<PortalAuthResponse>(await fetch("/api/auth/login", {
        method: "POST", headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      }));
    },
    async loginEmployee(payload: EmployeePortalLoginRequest): Promise<PortalAuthResponse> {
      return parseJson<PortalAuthResponse>(await fetch("/api/auth/employee/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      }));
    },
    async logout(): Promise<AuthOkResponse> {
      return parseJson<AuthOkResponse>(await fetch("/api/auth/logout", { method: "POST", keepalive: true }));
    },
    async getSession(): Promise<PortalAuthResponse> {
      return parseJson<PortalAuthResponse>(await fetch("/api/auth/session", { method: "GET", cache: "no-store" }));
    },
    async refreshToken(): Promise<AuthOkResponse> {
      return parseJson<AuthOkResponse>(await fetch("/api/auth/refresh", { method: "POST" }));
    },
    async getProfile(): Promise<ProfileResponse> {
      return parseJson<ProfileResponse>(await fetch("/api/auth/profile", { method: "GET", cache: "no-store" }));
    },
    async updateProfile(payload: Record<string, unknown>): Promise<ProfileResponse> {
      return parseJson<ProfileResponse>(await fetch("/api/auth/profile", {
        method: "PATCH", headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload),
      }));
    },
    async getProfilePhotoUploadUrl(): Promise<ProfilePhotoUrlResponse> {
      return parseJson<ProfilePhotoUrlResponse>(await fetch("/api/auth/profile-photo-url", { method: "GET", cache: "no-store" }));
    },
    async changePassword(payload: { old_password: string; new_password: string }): Promise<AuthOkResponse> {
      return parseJson<AuthOkResponse>(await fetch("/api/auth/change-password", {
        method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload),
      }));
    },
    async listSessions(): Promise<SessionsResponse> {
      return parseJson<SessionsResponse>(await fetch("/api/auth/sessions", { method: "GET", cache: "no-store" }));
    },
    async revokeSession(sessionId: string): Promise<AuthOkResponse> {
      return parseJson<AuthOkResponse>(await fetch(`/api/auth/sessions/${sessionId}`, { method: "DELETE" }));
    },
    async revokeAllSessions(): Promise<AuthOkResponse> {
      return parseJson<AuthOkResponse>(await fetch("/api/auth/sessions", { method: "DELETE" }));
    },
    async enableTotp(): Promise<TotpResponse> {
      return parseJson<TotpResponse>(await fetch("/api/auth/totp", { method: "POST" }));
    },
    async disableTotp(totpCode: string): Promise<AuthOkResponse> {
      return parseJson<AuthOkResponse>(await fetch("/api/auth/totp", {
        method: "DELETE", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ totp_code: totpCode }),
      }));
    },
    async sendOtp(purpose?: string): Promise<AuthOkResponse> {
      return parseJson<AuthOkResponse>(await fetch("/api/auth/send-otp", {
        method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ purpose }),
      }));
    },
    async verifyOtp(otp: string, purpose?: string): Promise<OtpResponse> {
      return parseJson<OtpResponse>(await fetch("/api/auth/verify-otp", {
        method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ otp, purpose }),
      }));
    },
    async sendEmailOtp(purpose?: string): Promise<AuthOkResponse> {
      return parseJson<AuthOkResponse>(await fetch("/api/auth/send-email-otp", {
        method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ purpose }),
      }));
    },
    async verifyEmail(payload: { token?: string; otp?: string }): Promise<AuthOkResponse> {
      return parseJson<AuthOkResponse>(await fetch("/api/auth/verify-email", {
        method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload),
      }));
    },
    async findPortalUser(identifier: string): Promise<PortalUserLookupResult> {
      return parseJson<PortalUserLookupResult>(await fetch("/api/auth/users/find", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ identifier }),
      }));
    },
    async searchEmployeeOrganisations(query: string): Promise<EmployeeLoginOrganisationsResult> {
      const params = new URLSearchParams({ q: query, page_size: "8" });
      return parseJson<EmployeeLoginOrganisationsResult>(await fetch(`/api/auth/employee/organisations?${params}`, {
        method: "GET",
        cache: "no-store",
      }));
    },
    async activateEmployee(payload: {
      organisationId?: string;
      organisationCode: string;
      employeeId: string;
      email: string;
    }): Promise<EmployeeActivationResult> {
      return parseJson<EmployeeActivationResult>(await fetch("/api/auth/employee/activate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      }));
    },
    async completeEmployeeActivation(payload: {
      email: string;
      otpId: string;
      otpCode: string;
      newPassword: string;
    }): Promise<AuthOkResponse> {
      return parseJson<AuthOkResponse>(await fetch("/api/auth/employee/complete-activation", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      }));
    },
  },

  employeeSelf: {
    async getProfile(): Promise<ProfileResponse> {
      return parseJson<ProfileResponse>(await fetch("/api/employee/profile", {
        method: "GET",
        cache: "no-store",
      }));
    },
    async getCoverage(): Promise<EmployeeCoverageResult> {
      return parseJson<EmployeeCoverageResult>(await fetch("/api/employee/coverage", {
        method: "GET",
        cache: "no-store",
      }));
    },
  },

  // ── Employees ─────────────────────────────────────────────────────────────
  employees: {
    async list(options?: { pageSize?: number; offset?: number; businessId?: string; departmentId?: string; status?: number }): Promise<EmployeeListResult> {
      const p = new URLSearchParams({ page_size: String(options?.pageSize ?? 50), offset: String(options?.offset ?? 0) });
      if (options?.businessId) p.set("business_id", options.businessId);
      if (options?.departmentId) p.set("department_id", options.departmentId);
      if (options?.status != null) p.set("status", String(options.status));
      return parseJson<EmployeeListResult>(await fetch(`/api/employees?${p}`, { method: "GET", cache: "no-store" }));
    },
    async get(id: string): Promise<EmployeeSingleResult> {
      return parseJson<EmployeeSingleResult>(await fetch(`/api/employees/${id}`, { method: "GET", cache: "no-store" }));
    },
    async create(payload: EmployeeCreatePayload): Promise<EmployeeSingleResult> {
      return parseJson<EmployeeSingleResult>(await fetch("/api/employees", {
        method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload),
      }));
    },
    async update(id: string, payload: EmployeeUpdatePayload): Promise<EmployeeSingleResult> {
      return parseJson<EmployeeSingleResult>(await fetch(`/api/employees/${id}`, {
        method: "PATCH", headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload),
      }));
    },
    async delete(id: string): Promise<ApiResult> {
      return parseJson<ApiResult>(await fetch(`/api/employees/${id}`, { method: "DELETE" }));
    },
  },

  // ── Departments ───────────────────────────────────────────────────────────
  departments: {
    async list(pageSize = 50, offset = 0, businessId?: string): Promise<DepartmentListResult> {
      const p = new URLSearchParams({ page_size: String(pageSize), offset: String(offset) });
      if (businessId) p.set("business_id", businessId);
      return parseJson<DepartmentListResult>(await fetch(`/api/departments?${p}`, { method: "GET", cache: "no-store" }));
    },
    async create(payload: { name: string; businessId: string }): Promise<DepartmentSingleResult> {
      return parseJson<DepartmentSingleResult>(await fetch("/api/departments", {
        method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload),
      }));
    },
    async update(id: string, payload: { name: string }): Promise<DepartmentSingleResult> {
      return parseJson<DepartmentSingleResult>(await fetch(`/api/departments/${id}`, {
        method: "PATCH", headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload),
      }));
    },
    async delete(id: string): Promise<ApiResult> {
      return parseJson<ApiResult>(await fetch(`/api/departments/${id}`, { method: "DELETE" }));
    },
  },

  // ── Organisations ─────────────────────────────────────────────────────────
  organisations: {
    async list(): Promise<OrgListResult> {
      return parseJson<OrgListResult>(await fetch("/api/organisations", { method: "GET", cache: "no-store" }));
    },
    async get(id: string): Promise<OrgSingleResult> {
      return parseJson<OrgSingleResult>(await fetch(`/api/organisations/${id}`, { method: "GET", cache: "no-store" }));
    },
    async getMe(): Promise<OrgSingleResult> {
      return parseJson<OrgSingleResult>(await fetch("/api/organisations/me", { method: "GET", cache: "no-store" }));
    },
    async create(payload: OrgCreatePayload): Promise<OrgSingleResult> {
      return parseJson<OrgSingleResult>(await fetch("/api/organisations", {
        method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload),
      }));
    },
    async update(id: string, payload: OrgUpdatePayload): Promise<OrgSingleResult> {
      return parseJson<OrgSingleResult>(await fetch(`/api/organisations/${id}`, {
        method: "PATCH", headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload),
      }));
    },
    async delete(id: string): Promise<ApiResult> {
      return parseJson<ApiResult>(await fetch(`/api/organisations/${id}`, { method: "DELETE" }));
    },
    async listMembers(id: string): Promise<OrgMembersResult> {
      return parseJson<OrgMembersResult>(await fetch(`/api/organisations/${id}/members`, { method: "GET", cache: "no-store" }));
    },
    async assignAdmin(id: string, memberId: string): Promise<OrgMemberResult> {
      return parseJson<OrgMemberResult>(await fetch(`/api/organisations/${id}/admins?action=assign`, {
        method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ memberId }),
      }));
    },
    async createAdmin(id: string, payload: OrgAdminCreatePayload): Promise<OrgMemberResult> {
      return parseJson<OrgMemberResult>(await fetch(`/api/organisations/${id}/admins`, {
        method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload),
      }));
    },
    async removeMember(id: string, memberId: string): Promise<ApiResult> {
      return parseJson<ApiResult>(await fetch(`/api/organisations/${id}/members/${memberId}`, { method: "DELETE" }));
    },
    async addMember(id: string, userId: string, role: string): Promise<OrgMemberResult> {
      return parseJson<OrgMemberResult>(await fetch(`/api/organisations/${id}/members`, {
        method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ userId, role }),
      }));
    },
    async assignExistingAdmin(id: string, userId: string, temporaryPassword?: string): Promise<ApiResult> {
      return parseJson<ApiResult>(await fetch(`/api/organisations/${id}/admins?action=assign`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ userId, temporaryPassword }),
      }));
    },
    async approve(id: string): Promise<OrgSingleResult> {
      return parseJson<OrgSingleResult>(await fetch(`/api/organisations/${id}/approve`, {
        method: "POST", headers: { "Content-Type": "application/json" },
      }));
    },
  },

  // ── Purchase Orders ───────────────────────────────────────────────────────
  purchaseOrders: {
    async list(options?: { pageSize?: number; offset?: number; status?: number }): Promise<POListResult> {
      const p = new URLSearchParams({ page_size: String(options?.pageSize ?? 50), offset: String(options?.offset ?? 0) });
      if (options?.status != null) p.set("status", String(options.status));
      return parseJson<POListResult>(await fetch(`/api/purchase-orders?${p}`, { method: "GET", cache: "no-store" }));
    },
    async get(id: string): Promise<POSingleResult> {
      return parseJson<POSingleResult>(await fetch(`/api/purchase-orders/${id}`, { method: "GET", cache: "no-store" }));
    },
    async create(payload: PurchaseOrderCreatePayload): Promise<POSingleResult> {
      return parseJson<POSingleResult>(await fetch("/api/purchase-orders", {
        method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload),
      }));
    },
    async update(id: string, payload: PurchaseOrderUpdatePayload): Promise<POSingleResult> {
      return parseJson<POSingleResult>(await fetch(`/api/purchase-orders/${id}`, {
        method: "PATCH", headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload),
      }));
    },
    async delete(id: string): Promise<ApiResult> {
      return parseJson<ApiResult>(await fetch(`/api/purchase-orders/${id}`, { method: "DELETE" }));
    },
    async getCatalog(): Promise<POCatalogResult> {
      return parseJson<POCatalogResult>(await fetch("/api/purchase-orders/catalog", { method: "GET", cache: "no-store" }));
    },
  },

  // ── Billing & Invoices ────────────────────────────────────────────────────
  billing: {
    async listInvoices(options?: { pageSize?: number; page?: number; status?: string }): Promise<ApiResult<{ invoices?: unknown[]; total?: number }>> {
      const p = new URLSearchParams({ page_size: String(options?.pageSize ?? 20) });
      if (options?.page) p.set("page", String(options.page));
      if (options?.status) p.set("status", options.status);
      return parseJson<ApiResult<{ invoices?: unknown[]; total?: number }>>(await fetch(`/api/billing/invoices?${p}`, { cache: "no-store" }));
    },
    async listPayments(options?: { pageSize?: number; page?: number; status?: string }): Promise<ApiResult<{ payments?: unknown[]; total?: number }>> {
      const p = new URLSearchParams({ page_size: String(options?.pageSize ?? 20) });
      if (options?.page) p.set("page", String(options.page));
      if (options?.status) p.set("status", options.status);
      return parseJson<ApiResult<{ payments?: unknown[]; total?: number }>>(await fetch(`/api/billing/payments?${p}`, { cache: "no-store" }));
    },
  },

  // ── Insurance Plans ───────────────────────────────────────────────────────
  insurancePlans: {
    async list(options?: { category?: string; pageSize?: number }): Promise<ApiResult<{ items?: unknown[] }>> {
      const p = new URLSearchParams({ page_size: String(options?.pageSize ?? 50) });
      if (options?.category) p.set("category", options.category);
      return parseJson<ApiResult<{ items?: unknown[] }>>(await fetch(`/api/insurance-plans?${p}`, { cache: "no-store" }));
    },
  },

  // ── Documents ─────────────────────────────────────────────────────────────
  documents: {
    async generate(payload: GenerateDocumentPayload): Promise<DocumentSingleResult> {
      return parseJson<DocumentSingleResult>(await fetch("/api/documents", {
        method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload),
      }));
    },
    async list(options: { entityType: string; entityId: string; status?: string; page?: number; pageSize?: number }): Promise<DocumentListResult> {
      const p = new URLSearchParams({ entity_type: options.entityType, entity_id: options.entityId });
      if (options.status) p.set("status", options.status);
      if (options.page) p.set("page", String(options.page));
      if (options.pageSize) p.set("page_size", String(options.pageSize));
      return parseJson<DocumentListResult>(await fetch(`/api/documents?${p}`, { cache: "no-store" }));
    },
    async get(documentId: string): Promise<DocumentSingleResult> {
      return parseJson<DocumentSingleResult>(await fetch(`/api/documents/${documentId}`, { cache: "no-store" }));
    },
    async download(documentId: string): Promise<DocumentDownloadResult> {
      return parseJson<DocumentDownloadResult>(await fetch(`/api/documents/${documentId}/download`, { cache: "no-store" }));
    },
    async delete(documentId: string): Promise<ApiResult> {
      return parseJson<ApiResult>(await fetch(`/api/documents/${documentId}`, { method: "DELETE" }));
    },
  },
};

/**
 * makeDirectHttp — returns typed helpers for direct HTTP calls to the gateway.
 * Use this for endpoints not (yet) exposed as typed SDK methods.
 * Shares the same cookie/CSRF auth headers as makeSdkClient.
 *
 * sessionOverrides are forwarded as x-portal/x-user-id/x-business-id/x-tenant-id
 * so the backend authz interceptor gets the correct Casbin domain. Without these,
 * superadmin direct HTTP calls (e.g. admins:assign) would 403.
 */
// ─────────────────────────────────────────────────────────────────────────────
// makeDocgenClient — server-side docgen helpers for API route handlers.
// The docgen service is NOT in the generated SDK. It is exposed by the gateway
// under /v1/documents and /v1/document-templates. Uses makeDirectHttp internally.
// ─────────────────────────────────────────────────────────────────────────────

export function makeDocgenClient(
  request: Request,
  sessionOverrides?: { portal?: string; userId?: string; businessId?: string; tenantId?: string }
) {
  const http = makeDirectHttp(request, sessionOverrides);

  return {
    // ── Documents ──────────────────────────────────────────────────────────
    generate: (payload: {
      template_id: string; entity_type: string; entity_id: string;
      data?: Record<string, unknown>; include_qr_code?: boolean;
    }) => http.post("/v1/documents", payload),

    getDocument: (documentId: string) =>
      http.get(`/v1/documents/${documentId}`),

    listDocuments: (options: {
      entityType: string; entityId: string;
      status?: string; page?: number; pageSize?: number;
    }) => {
      const params = new URLSearchParams({
        entity_type: options.entityType,
        entity_id: options.entityId,
      });
      if (options.status) params.set("status", options.status);
      if (options.page) params.set("page", String(options.page));
      if (options.pageSize) params.set("page_size", String(options.pageSize));
      return http.get(`/v1/entities/${options.entityType}/${options.entityId}/documents?${params}`);
    },

    downloadDocument: (documentId: string) =>
      http.get(`/v1/documents/${documentId}/download`),

    deleteDocument: (documentId: string) =>
      http.delete(`/v1/documents/${documentId}`),

    // ── Templates ──────────────────────────────────────────────────────────
    createTemplate: (payload: {
      name: string; type: string; description?: string;
      template_content: string; output_format: string; variables?: string[];
    }) => http.post("/v1/document-templates", payload),

    getTemplate: (templateId: string) =>
      http.get(`/v1/document-templates/${templateId}`),

    listTemplates: (options?: {
      type?: string; activeOnly?: boolean; pageSize?: number; pageToken?: string;
    }) => {
      const params = new URLSearchParams();
      if (options?.type) params.set("type", options.type);
      if (options?.activeOnly) params.set("active_only", "true");
      if (options?.pageSize) params.set("page_size", String(options.pageSize));
      if (options?.pageToken) params.set("page_token", options.pageToken);
      const qs = params.toString() ? `?${params}` : "";
      return http.get(`/v1/document-templates${qs}`);
    },

    updateTemplate: (templateId: string, payload: {
      template?: {
        name?: string; description?: string;
        template_content?: string; output_format?: string; is_active?: boolean;
      };
    }) => http.patch(`/v1/document-templates/${templateId}`, payload),

    deactivateTemplate: (templateId: string, reason?: string) =>
      http.post(`/v1/document-templates/${templateId}/deactivate`, { reason: reason ?? "" }),

    deleteTemplate: (templateId: string) =>
      http.delete(`/v1/document-templates/${templateId}`),
  };
}

export type DocgenClient = ReturnType<typeof makeDocgenClient>;

export function makeDirectHttp(request: Request, sessionOverrides?: { portal?: string; userId?: string; businessId?: string; tenantId?: string }) {
  const cookieHeader = request.headers.get("cookie") ?? "";
  const csrf = cookieHeader.match(/(?:^|;\s*)csrf_token=([^;]*)/)?.[1] ?? "";
  const extraHeaders: Record<string, string> = { "Content-Type": "application/json" };
  if (cookieHeader) extraHeaders["cookie"] = cookieHeader;
  if (csrf) extraHeaders["X-CSRF-Token"] = decodeURIComponent(csrf);

  // Forward portal context headers — same logic as makeSdkClient.
  // Super admin: x-portal=PORTAL_SYSTEM (no x-business-id needed).
  // B2B admin:   x-portal=PORTAL_B2B + x-business-id={org_id}.
  const xPortal = sessionOverrides?.portal ?? request.headers.get("x-portal") ?? "";
  const xBusinessId = sessionOverrides?.businessId ?? request.headers.get("x-business-id") ?? "";
  const xUserId = sessionOverrides?.userId ?? request.headers.get("x-user-id") ?? "";
  const xTenantId = sessionOverrides?.tenantId ?? request.headers.get("x-tenant-id") ?? "";
  if (xPortal) extraHeaders["x-portal"] = xPortal;
  if (xBusinessId) extraHeaders["x-business-id"] = xBusinessId;
  if (xUserId) extraHeaders["x-user-id"] = xUserId;
  if (xTenantId) extraHeaders["x-tenant-id"] = xTenantId;

  const base = process.env.INSURETECH_API_BASE_URL ?? process.env.NEXT_PUBLIC_INSURETECH_API_BASE_URL ?? "http://localhost:8080";

  const doFetch = async (method: string, path: string, body?: unknown) => {
    const res = await fetch(`${base}${path}`, {
      method,
      headers: extraHeaders,
      body: body !== undefined ? JSON.stringify(body) : undefined,
      cache: "no-store",
    });
    const raw = await res.text();
    let envelope: Record<string, unknown>;
    try {
      envelope = raw ? (JSON.parse(raw) as Record<string, unknown>) : {};
    } catch {
      envelope = raw ? { error: { message: raw, code: 'PARSE_ERROR', http_status_code: res.status, retryable: false, field_violations: [] } } : {};
    }
    // Unwrap the unified ApiResponse<T> envelope
    // { success, data, error, meta } — same shape for all endpoints
    const success = typeof envelope.success === 'boolean' ? envelope.success : res.ok;
    const data = success ? (envelope.data as Record<string, unknown> ?? {}) : null;
    const error = envelope.error as Record<string, unknown> | null ?? null;
    const errorMessage = error && typeof error.message === "string" ? error.message : "";
    const message = !success ? (errorMessage || raw) : undefined;
    return { ok: success, status: res.status, data: data ?? {}, error, message };
  };

  return {
    get: (path: string) => doFetch("GET", path),
    post: (path: string, body?: unknown) => doFetch("POST", path, body),
    patch: (path: string, body?: unknown) => doFetch("PATCH", path, body),
    put: (path: string, body?: unknown) => doFetch("PUT", path, body),
    delete: (path: string) => doFetch("DELETE", path),
  };
}
