import { NextResponse } from "next/server";

import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { SESSION_COOKIE_NAME } from "@lib/auth/session";
import { toPortalSessionFromLogin, getSetCookieHeaders } from "@lib/auth/backend-auth";
import type { LoginRequest } from "@lib/types/auth";
import {
  getBangladeshMobileValidationMessage,
  normalizeBangladeshMobile,
} from "@/src/lib/utils/bd-mobile";

/**
 * Maps backend errors (gRPC codes, HTTP status codes, raw strings) to clean,
 * user-facing messages. Never leaks internal RPC details to the UI.
 */
function toUserFriendlyLoginError(error: unknown, httpStatus: number): string {
  // Extract the raw error string from whatever shape the SDK returns
  let raw = "";
  if (typeof error === "string") {
    raw = error;
  } else if (error && typeof error === "object") {
    for (const key of ["message", "error", "detail", "description"] as const) {
      const v = (error as Record<string, unknown>)[key];
      if (typeof v === "string" && v.trim()) { raw = v; break; }
    }
  }
  const lower = raw.toLowerCase();

  // ── gRPC / HTTP status → friendly message map ──────────────────────────────

  // Account locked / too many attempts
  if (
    lower.includes("locked") ||
    lower.includes("too many") ||
    lower.includes("rate limit") ||
    lower.includes("max attempt") ||
    lower.includes("blocked") ||
    httpStatus === 422 ||
    httpStatus === 429
  ) {
    if (raw.trim()) {
      return raw.trim().replace(/\.$/, "");
    }
    return "Your account has been temporarily locked due to too many failed attempts. Please try again later.";
  }

  // Wrong password / user not found
  if (
    httpStatus === 401 ||
    lower.includes("unauthenticated") ||
    lower.includes("invalid password") ||
    lower.includes("invalid credentials") ||
    lower.includes("wrong password") ||
    lower.includes("password") ||
    lower.includes("not found") ||
    lower.includes("no user") ||
    lower.includes("user not found")
  ) {
    return "Mobile number or password is incorrect. Please try again.";
  }

  // Account not active / disabled
  if (
    lower.includes("inactive") ||
    lower.includes("disabled") ||
    lower.includes("suspended") ||
    lower.includes("banned") ||
    lower.includes("not active")
  ) {
    return "Your account is not active. Please contact your administrator.";
  }

  // Invalid mobile number (should be caught client-side, but just in case authn rejects it)
  if (
    lower.includes("invalid mobile") ||
    lower.includes("invalid_argument") ||
    lower.includes("invalidargument") ||
    lower.includes("mobile_number") ||
    lower.includes("phone") ||
    httpStatus === 400
  ) {
    return "Invalid mobile number. Please enter a valid Bangladesh number (e.g. 01712345678).";
  }

  // Server / network errors
  if (httpStatus >= 500 || lower.includes("unavailable") || lower.includes("internal")) {
    return "The service is temporarily unavailable. Please try again in a moment.";
  }

  // Generic fallback — never show raw RPC text
  return "Login failed. Please check your mobile number and password and try again.";
}

const CSRF_COOKIE_NAME = "csrf_token";

function extractCookieValue(setCookieHeader: string, cookieName: string): string | undefined {
  const [nameValue] = setCookieHeader.split(";", 1);
  if (!nameValue) {
    return undefined;
  }
  const separatorIndex = nameValue.indexOf("=");
  if (separatorIndex <= 0) {
    return undefined;
  }
  const name = nameValue.slice(0, separatorIndex).trim();
  if (name !== cookieName) {
    return undefined;
  }
  return nameValue.slice(separatorIndex + 1);
}

function getApiBaseUrl(): string {
  return (
    process.env.INSURETECH_API_BASE_URL ??
    process.env.NEXT_PUBLIC_INSURETECH_API_BASE_URL ??
    "http://localhost:8080"
  );
}

export async function POST(request: Request) {
  let payload: LoginRequest;
  try {
    payload = (await request.json()) as LoginRequest;
  } catch {
    return NextResponse.json({ ok: false, message: "Invalid login payload" }, { status: 400 });
  }

  const mobileRaw = payload.mobileNumber?.trim();
  if (!mobileRaw) {
    return NextResponse.json(
      { ok: false, message: "Mobile number is required" },
      { status: 400 }
    );
  }
  if (!payload.password?.trim()) {
    return NextResponse.json(
      { ok: false, message: "Password is required" },
      { status: 400 }
    );
  }

  const normalizedMobile = normalizeBangladeshMobile(mobileRaw);
  if (!normalizedMobile) {
    return NextResponse.json(
      {
        ok: false,
        message: getBangladeshMobileValidationMessage(),
      },
      { status: 400 }
    );
  }

  const requestCookie = request.headers.get("cookie");
  let backendResponse: Response;
  let rawBackendBody = "";
  try {
    backendResponse = await fetch(`${getApiBaseUrl()}/v1/auth/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...(requestCookie ? { cookie: requestCookie } : {}),
      },
      body: JSON.stringify({
        mobile_number: normalizedMobile,
        password: payload.password,
        device_id: payload.deviceId ?? "b2b-portal-web",
        device_type: "WEB",
        device_name: "B2B Portal Web",
      }),
      cache: "no-store",
    });
    rawBackendBody = await backendResponse.text();
  } catch (error) {
    return NextResponse.json(
      { ok: false, message: toUserFriendlyLoginError(error, 502) },
      { status: 502 }
    );
  }

  let backendPayload: Record<string, unknown> = {};
  if (rawBackendBody) {
    try {
      backendPayload = JSON.parse(rawBackendBody) as Record<string, unknown>;
    } catch {
      backendPayload = {};
    }
  }

  const backendError =
    backendPayload && typeof backendPayload.error === "object"
      ? (backendPayload.error as Record<string, unknown>)
      : undefined;

  if (!backendResponse.ok || backendPayload.success === false) {
    const httpStatus =
      (typeof backendError?.http_status_code === "number" ? backendError.http_status_code : undefined) ??
      (backendResponse.status || 500);
    return NextResponse.json(
      {
        ok: false,
        message: toUserFriendlyLoginError(backendError ?? rawBackendBody, httpStatus),
      },
      { status: httpStatus }
    );
  }

  const loginData =
    backendPayload && typeof backendPayload.data === "object"
      ? (backendPayload.data as Record<string, unknown>)
      : backendPayload;

  const response = NextResponse.json({ ok: true }, { status: backendResponse.status || 200 });

  // Primary: read session_token from the JSON response body.
  // The gateway Login handler keeps session_token in the proto JSON body (in addition
  // to setting the HttpOnly cookie) specifically so the Next.js BFF can read it here
  // without relying on Set-Cookie header forwarding.
  // Set-Cookie is a forbidden header in the Fetch API Headers constructor for the
  // generated SDK path, so we read the backend response headers directly here.
  let sessionToken: string | undefined =
    typeof loginData?.session_token === "string" && loginData.session_token
      ? loginData.session_token
      : undefined;

  // Fallback: try reading from Set-Cookie headers (works when no SDK interceptor is involved).
  if (!sessionToken) {
    const setCookieHeaders = getSetCookieHeaders(backendResponse.headers);
    const backendSessionCookie = setCookieHeaders.find((value) =>
      value.startsWith(`${SESSION_COOKIE_NAME}=`)
    );
    sessionToken = backendSessionCookie
      ? extractCookieValue(backendSessionCookie, SESSION_COOKIE_NAME)
      : undefined;
  }
  if (sessionToken) {
    response.cookies.set({
      name: SESSION_COOKIE_NAME,
      value: sessionToken,
      path: "/",
      httpOnly: true,
      sameSite: "lax",
      secure: process.env.NODE_ENV === "production",
      maxAge: 60 * 60 * 12,
    });
  }

  const csrfToken =
    backendResponse.headers.get("x-csrf-token") ??
    (typeof loginData?.csrf_token === "string" ? loginData.csrf_token : undefined);
  if (csrfToken) {
    response.cookies.set({
      name: CSRF_COOKIE_NAME,
      value: csrfToken,
      path: "/",
      httpOnly: true,
      sameSite: "lax",
      secure: process.env.NODE_ENV === "production",
      maxAge: 60 * 60 * 12,
    });
  }

  const sessionCookieHeader = sessionToken ? `${SESSION_COOKIE_NAME}=${sessionToken}` : undefined;
  // loginData is the unwrapped LoginResponse from the gateway ApiResponse envelope.
  const session = await toPortalSessionFromLogin(loginData ?? {}, sessionCookieHeader);

  // Store the session in the in-memory session store using the backend token as the key.
  // This allows requireServerSession() to find it on the fast path without calling the gateway.
  if (sessionToken && session?.principal) {
    const { createSession } = await import("@lib/auth/session-store");
    createSession(session.principal, sessionToken);
  }

  const finalResponse = NextResponse.json({ ok: true, session }, { status: response.status });
  for (const cookie of response.cookies.getAll()) {
    finalResponse.cookies.set(cookie);
  }

  // Set lightweight metadata cookies used by:
  //   1. Edge middleware (portal_role) — for role-based page routing without hitting DB
  //   2. API route session-headers helper — to inject x-portal/x-user-id/x-business-id
  //      into backend SDK calls so the Casbin authz interceptor gets the right domain.
  //
  // These are NOT security boundaries — the backend session cookie enforces real authz.
  const portalRole = session.principal.role ?? "BUSINESS_ADMIN";
  const portalBizId = session.principal.businessId ?? "";

  // user_id may not be in the login JSON response body — it lives in the HttpOnly session cookie.
  // Prefer: from toPortalSessionFromLogin, then from result.data directly, then from getCurrentSession.
  let portalUserId = session.principal.user?.userId ?? (loginData as Record<string, unknown>)?.user_id as string ?? "";

  // If user_id is still empty, call getCurrentSession using the new session token to resolve it.
  // This handles gateways that don't return user_id in the login JSON body.
  if (!portalUserId && sessionToken) {
    try {
      const { toPortalSessionFromCurrentSession } = await import("@lib/auth/backend-auth");
      const cookieStr = `${SESSION_COOKIE_NAME}=${sessionToken}`;
      // Build a fake Request carrying the new session cookie so makeSdkClient can forward it.
      const tempReq = new Request(request.url, { headers: { cookie: cookieStr } });
      const tempSdk = makeSdkClient(tempReq);
      const sessionRes = await tempSdk.getCurrentSession();
      if (sessionRes.response?.ok && sessionRes.data) {
        const portalSession = await toPortalSessionFromCurrentSession(sessionRes.data, cookieStr);
        portalUserId = portalSession?.principal?.user?.userId ?? "";
      }
    } catch { /* ignore — userId will remain empty */ }
  }

  let resolvedKycVerified: boolean | undefined = loginData?.user
    ? Boolean((loginData.user as Record<string, unknown>)?.kyc_verified)
    : undefined;

  // "pending_review" means the user completed eKYC but manual approval is still pending.
  // We track this separately so the login cookie can be set to "pending_review" (no gate)
  // rather than "false" (gate re-triggers). Starts as false; set to true if KYC status
  // is PENDING_REVIEW in the backend record.
  let kycIsPendingReview = false;

  if (sessionToken && portalUserId) {
    try {
      const cookieStr = `${SESSION_COOKIE_NAME}=${sessionToken}`;
      const tempReq = new Request(request.url, { headers: { cookie: cookieStr } });
      const tempSdk = makeSdkClient(tempReq);

      // 1. Fetch user profile to get kyc_verified boolean
      if (resolvedKycVerified === undefined) {
        const profileRes = await tempSdk.getUserProfile({ path: { user_id: portalUserId } });
        if (profileRes.response?.ok && profileRes.data) {
          const rawProfile = profileRes.data as Record<string, unknown>;
          const profile = (rawProfile.profile ?? rawProfile) as Record<string, unknown>;
          resolvedKycVerified = Boolean(profile.kyc_verified);
        }
      }

      // 2. If profile says not verified, check the KYC record status.
      //    A user who finished eKYC (PENDING_REVIEW or VERIFIED in the KYC record)
      //    must NOT be re-gated to /kyc — the profile update may have raced or the
      //    admin approval step is still pending.
      if (!resolvedKycVerified) {
        const gatewayUrl = process.env.INSURETECH_GATEWAY_URL ?? process.env.INSURETECH_API_BASE_URL ?? "http://localhost:8080";
        const kycStatusRes = await fetch(`${gatewayUrl}/v1/auth/users/${portalUserId}/kyc`, {
          headers: {
            cookie: `${SESSION_COOKIE_NAME}=${sessionToken}`,
            "x-portal": "b2b",
          },
          cache: "no-store",
        }).catch(() => null);
        if (kycStatusRes?.ok) {
          const kycBody = await kycStatusRes.json().catch(() => ({})) as Record<string, unknown>;
          const kycData = (kycBody.data ?? kycBody) as Record<string, unknown>;
          const kycStatus = ((kycData.status ?? "") as string).toUpperCase();
          if (kycStatus === "PENDING_REVIEW") {
            kycIsPendingReview = true;
          } else if (kycStatus === "VERIFIED") {
            // KYC record says VERIFIED but profile not yet updated — treat as verified
            resolvedKycVerified = true;
          }
        }
      }
    } catch {
      // Ignore lookup errors and fall back to the safe default below.
    }
  }

  const cookieOpts = {
    path: "/",
    httpOnly: false, // must be readable by edge middleware + session-headers helper
    sameSite: "lax" as const,
    secure: process.env.NODE_ENV === "production",
    maxAge: 60 * 60 * 12,
  };

  const passwordChangeRequired = Boolean(
    (loginData as Record<string, unknown>)?.password_change_required ??
    session?.passwordChangeRequired
  );

  finalResponse.cookies.set({ name: "portal_role",    value: portalRole,    ...cookieOpts });
  finalResponse.cookies.set({ name: "portal_user_id", value: portalUserId,  ...cookieOpts });
  finalResponse.cookies.set({ name: "portal_biz_id",  value: portalBizId,   ...cookieOpts });
  finalResponse.cookies.set({
    name: "portal_password_change_required",
    value: passwordChangeRequired ? "true" : "false",
    ...cookieOpts,
  });

  // KYC verification status cookie — read by edge middleware to gate unverified admins.
  // "false"         → redirect to /kyc on every page navigation until verified.
  // "true"          → fully verified, no gate.
  // "pending_review"→ submitted but awaiting manual approval, no gate.
  // Absent/unknown  → no gate (backward compat for existing sessions).
  //
  // IMPORTANT: Only set "false" when we have a *confirmed* false from the backend
  // AND the KYC record is not already PENDING_REVIEW or VERIFIED.
  // If the profile fetch failed or returned nothing (resolvedKycVerified === undefined),
  // we must NOT default to "false" — that would trap the user in an infinite /kyc
  // redirect loop. Omitting the cookie causes the middleware to skip the gate entirely
  // (safe backward-compat path).
  if (resolvedKycVerified === true) {
    finalResponse.cookies.set({ name: "portal_kyc_verified", value: "true", ...cookieOpts });
  } else if (kycIsPendingReview) {
    // User completed eKYC; awaiting admin approval — allow access, show pending UI.
    finalResponse.cookies.set({ name: "portal_kyc_verified", value: "pending_review", ...cookieOpts });
  } else if (resolvedKycVerified === false) {
    finalResponse.cookies.set({ name: "portal_kyc_verified", value: "false", ...cookieOpts });
  }
  // If resolvedKycVerified is undefined (profile unavailable), we intentionally skip
  // setting portal_kyc_verified. The /kyc page itself will re-check the DB on load.

  // Store user contact info cookies so the My Profile page can display
  // mobile_number and email — these live on the User record, not UserProfile.
  // They are read-only identity fields (auth credentials), not profile fields.
  const portalMobile = loginData?.user
    ? ((loginData.user as Record<string, unknown>)?.mobile_number as string) ?? ""
    : "";
  const portalEmail = loginData?.user
    ? ((loginData.user as Record<string, unknown>)?.email as string) ?? ""
    : "";
  finalResponse.cookies.set({ name: "portal_mobile", value: portalMobile, ...cookieOpts });
  finalResponse.cookies.set({ name: "portal_email", value: portalEmail, ...cookieOpts });

  return finalResponse;
}
