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
    /** PATCH /v1/b2b/purchase-orders/{id} — not yet a generated SDK method */
    updatePurchaseOrderHttp: async (id: string, body: Record<string, unknown>) => {
      const res = await fetch(`${getBaseUrl()}/v1/b2b/purchase-orders/${id}`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
          ...extraHeaders,
        },
        body: JSON.stringify(body),
        cache: "no-store",
      });
      const data = (await res.json()) as Record<string, unknown>;
      return { ok: res.ok, status: res.status, data };
    },

    /** DELETE /v1/b2b/purchase-orders/{id} — not yet a generated SDK method */
    deletePurchaseOrderHttp: async (id: string) => {
      const res = await fetch(`${getBaseUrl()}/v1/b2b/purchase-orders/${id}`, {
        method: "DELETE",
        headers: extraHeaders,
        cache: "no-store",
      });
      const data = (await res.json()) as Record<string, unknown>;
      return { ok: res.ok, status: res.status, data };
    },
  };
}

export type B2bSdkClient = ReturnType<typeof makeSdkClient>;

/**
 * makeDirectHttp — returns typed helpers for direct HTTP calls to the gateway.
 * Use this for endpoints not (yet) exposed as typed SDK methods.
 * Shares the same cookie/CSRF auth headers as makeSdkClient.
 *
 * sessionOverrides are forwarded as x-portal/x-user-id/x-business-id/x-tenant-id
 * so the backend authz interceptor gets the correct Casbin domain. Without these,
 * superadmin direct HTTP calls (e.g. assign-admin) would 403.
 */
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
    let data: Record<string, unknown>;
    try {
      data = raw ? (JSON.parse(raw) as Record<string, unknown>) : {};
    } catch {
      data = raw ? { message: raw } : {};
    }
    return { ok: res.ok, status: res.status, data };
  };

  return {
    get: (path: string) => doFetch("GET", path),
    post: (path: string, body?: unknown) => doFetch("POST", path, body),
    patch: (path: string, body?: unknown) => doFetch("PATCH", path, body),
    put: (path: string, body?: unknown) => doFetch("PUT", path, body),
    delete: (path: string) => doFetch("DELETE", path),
  };
}
