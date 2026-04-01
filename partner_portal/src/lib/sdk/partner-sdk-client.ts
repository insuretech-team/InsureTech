/**
 * partner-sdk-client.ts
 * ─────────────────────
 * Server-side SDK client factory for Partner portal Next.js API route handlers.
 *
 * The SDK (@lifeplus/insuretech-sdk) is auto-generated from the protobuf
 * service definitions via the API pipeline script.
 *
 * Authentication: The portal uses cookie-based server-side sessions.
 * The gateway validates the session cookie and injects partner_id / user_id
 * from the JWT into every downstream gRPC call — the portal does NOT need to
 * pass those values. apiKey is a required config field but auth is handled by
 * the forwarded session cookie.
 *
 * Usage:
 *   import { makeSdkClient } from "@lib/sdk/partner-sdk-client";
 *   const sdk = makeSdkClient(req);
 *   const { data, error, response } = await sdk.listPartners({ query: { page_size: 50 } });
 */

import {
  // ── Auth ────────────────────────────────────────────────────────────────────
  authServiceLogin,
  authServiceLogout,
  authServiceGetCurrentSession,
  authServiceRefreshToken,
  authServiceChangePassword,
  authServiceGetUserProfile,
  authServiceUpdateUserProfile,
  createInsureTechClient,
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

export function makeSdkClient(
  request: Request,
  sessionOverrides?: {
    portal?: string;
    userId?: string;
    partnerId?: string;
    tenantId?: string;
  }
) {
  const cookieHeader = request.headers.get("cookie") ?? "";
  const csrf = extractCsrf(cookieHeader);

  const extraHeaders: Record<string, string> = {};
  if (cookieHeader) extraHeaders["cookie"] = cookieHeader;
  if (csrf) extraHeaders["X-CSRF-Token"] = csrf;

  // Forward portal + partner-id headers so the backend authz interceptor
  // can correctly resolve the Casbin domain.
  // Super admin: x-portal=PORTAL_SYSTEM (no x-partner-id needed)
  // Partner admin: x-portal=PORTAL_PARTNER + x-partner-id={partner_id}
  // Priority: sessionOverrides (from server session store) > request headers (browser-forwarded)
  const xPortal =
    sessionOverrides?.portal ?? request.headers.get("x-portal") ?? "";
  const xPartnerId =
    sessionOverrides?.partnerId ?? request.headers.get("x-partner-id") ?? "";
  const xUserId =
    sessionOverrides?.userId ?? request.headers.get("x-user-id") ?? "";
  const xTenantId =
    sessionOverrides?.tenantId ?? request.headers.get("x-tenant-id") ?? "";
  if (xPortal) extraHeaders["x-portal"] = xPortal;
  if (xPartnerId) extraHeaders["x-partner-id"] = xPartnerId;
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
    mobileLogin: (
      opts: Omit<Parameters<typeof authServiceLogin>[0], "client">
    ) =>
      authServiceLogin({ client: sdkClient, throwOnError: false, ...opts }),

    logout: (
      opts?: Omit<Parameters<typeof authServiceLogout>[0], "client">
    ) =>
      authServiceLogout({
        client: sdkClient,
        throwOnError: false,
        ...(opts as any),
      }),

    getCurrentSession: (
      opts?: Omit<Parameters<typeof authServiceGetCurrentSession>[0], "client">
    ) =>
      authServiceGetCurrentSession({
        client: sdkClient,
        throwOnError: false,
        ...(opts as any),
      }),

    refreshToken: (
      opts?: Omit<Parameters<typeof authServiceRefreshToken>[0], "client">
    ) =>
      authServiceRefreshToken({
        client: sdkClient,
        throwOnError: false,
        ...(opts as any),
      }),

    changePassword: (
      opts: Omit<Parameters<typeof authServiceChangePassword>[0], "client">
    ) =>
      authServiceChangePassword({
        client: sdkClient,
        throwOnError: false,
        ...(opts as any),
      }),

    getUserProfile: (
      opts?: Omit<Parameters<typeof authServiceGetUserProfile>[0], "client">
    ) =>
      authServiceGetUserProfile({
        client: sdkClient,
        throwOnError: false,
        ...(opts as any),
      }),

    updateUserProfile: (
      opts: Omit<Parameters<typeof authServiceUpdateUserProfile>[0], "client">
    ) =>
      authServiceUpdateUserProfile({
        client: sdkClient,
        throwOnError: false,
        ...(opts as any),
      }),

    // ── Partners ────────────────────────────────────────────────────────────
    // Partner SDK methods will be added here as they become available
    // For now, use makeDirectHttp for partner operations

    // ── Claims ──────────────────────────────────────────────────────────────
    // Claim SDK methods will be added here as they become available

    // ── Agents ──────────────────────────────────────────────────────────────
    // Agent SDK methods will be added here as they become available

    // ── Direct HTTP for SDK-missing operations ─────────────────────────────
    /** GET /v1/partners/me — resolve the caller's own partner organization */
    getMyPartner: () =>
      makeDirectHttp(request, sessionOverrides).get(`/v1/partners/me`),
  };
}

export type PartnerSdkClient = ReturnType<typeof makeSdkClient>;

/**
 * makeDirectHttp — returns typed helpers for direct HTTP calls to the gateway.
 * Use this for endpoints not (yet) exposed as typed SDK methods.
 * Shares the same cookie/CSRF auth headers as makeSdkClient.
 *
 * sessionOverrides are forwarded as x-portal/x-user-id/x-partner-id/x-tenant-id
 * so the backend authz interceptor gets the correct Casbin domain.
 */
export function makeDirectHttp(
  request: Request,
  sessionOverrides?: {
    portal?: string;
    userId?: string;
    partnerId?: string;
    tenantId?: string;
  }
) {
  const cookieHeader = request.headers.get("cookie") ?? "";
  const csrf = cookieHeader.match(/(?:^|;\s*)csrf_token=([^;]*)/)?.[1] ?? "";
  const extraHeaders: Record<string, string> = {
    "Content-Type": "application/json",
  };
  if (cookieHeader) extraHeaders["cookie"] = cookieHeader;
  if (csrf) extraHeaders["X-CSRF-Token"] = decodeURIComponent(csrf);

  // Forward portal context headers — same logic as makeSdkClient.
  // Super admin: x-portal=PORTAL_SYSTEM (no x-partner-id needed).
  // Partner admin: x-portal=PORTAL_PARTNER + x-partner-id={partner_id}.
  const xPortal =
    sessionOverrides?.portal ?? request.headers.get("x-portal") ?? "";
  const xPartnerId =
    sessionOverrides?.partnerId ?? request.headers.get("x-partner-id") ?? "";
  const xUserId =
    sessionOverrides?.userId ?? request.headers.get("x-user-id") ?? "";
  const xTenantId =
    sessionOverrides?.tenantId ?? request.headers.get("x-tenant-id") ?? "";
  if (xPortal) extraHeaders["x-portal"] = xPortal;
  if (xPartnerId) extraHeaders["x-partner-id"] = xPartnerId;
  if (xUserId) extraHeaders["x-user-id"] = xUserId;
  if (xTenantId) extraHeaders["x-tenant-id"] = xTenantId;

  const base =
    process.env.INSURETECH_API_BASE_URL ??
    process.env.NEXT_PUBLIC_INSURETECH_API_BASE_URL ??
    "http://localhost:8080";

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
      envelope = raw
        ? {
            error: {
              message: raw,
              code: "PARSE_ERROR",
              http_status_code: res.status,
              retryable: false,
              field_violations: [],
            },
          }
        : {};
    }
    // Unwrap the unified ApiResponse<T> envelope
    // { success, data, error, meta } — same shape for all endpoints
    const success =
      typeof envelope.success === "boolean" ? envelope.success : res.ok;
    const data = success
      ? ((envelope.data as Record<string, unknown>) ?? {})
      : null;
    const error = (envelope.error as Record<string, unknown> | null) ?? null;
    const message = !success
      ? typeof (error as any)?.message === "string"
        ? (error as any).message
        : raw
      : undefined;
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

export type DirectHttpClient = ReturnType<typeof makeDirectHttp>;
export type DirectHttpResult = Awaited<ReturnType<DirectHttpClient["get"]>>;
