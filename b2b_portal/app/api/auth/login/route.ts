import { NextResponse } from "next/server";

import {
  getSetCookieHeaders,
  loginWithMobile,
  toPortalSessionFromLogin,
} from "@lib/auth/backend-auth";
import { SESSION_COOKIE_NAME } from "@lib/auth/session";
import type { LoginRequest } from "@lib/types/auth";

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

  // Account locked / too many attempts
  if (
    lower.includes("locked") ||
    lower.includes("too many") ||
    lower.includes("rate limit") ||
    lower.includes("max attempt") ||
    lower.includes("blocked") ||
    httpStatus === 429
  ) {
    return "Your account has been temporarily locked due to too many failed attempts. Please try again later.";
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

// Valid Bangladesh operator prefixes: 013,014,015,016,017,018,019
const BD_PHONE_RE = /^880(13|14|15|16|17|18|19)\d{8}$/;

/**
 * Normalizes a Bangladesh mobile number to canonical E.164 form (+880XXXXXXXXXX).
 *
 * Accepted input variants (spaces, dashes, dots freely ignored):
 *   01712345678          → +8801712345678
 *   1712345678           → +8801712345678   (10 digits, no leading 0)
 *   8801712345678        → +8801712345678
 *   00 8801712345678     → +8801712345678
 *   +880 171-234-5678    → +8801712345678
 *   +88 01712345678      → +8801712345678   (typo with 88 instead of 880)
 *
 * Returns null when the number cannot be recognized as a valid BD number.
 */
function normalizeMobileNumber(value: string): string | null {
  // Strip everything except digits and a leading +
  const stripped = value.trim().replace(/[^\d+]/g, "");

  // Drop the leading + so we work purely with digits from here
  const digits = stripped.startsWith("+") ? stripped.slice(1) : stripped;

  let e164Digits: string; // will hold 880XXXXXXXXXX (13 digits)

  if (digits.startsWith("00880")) {
    // 008801712345678 → 8801712345678
    e164Digits = digits.slice(2);
  } else if (digits.startsWith("880")) {
    // 8801712345678
    e164Digits = digits;
  } else if (digits.startsWith("0088")) {
    // 00881712345678 — uncommon but handle gracefully
    e164Digits = "880" + digits.slice(4);
  } else if (digits.startsWith("88") && digits.length === 13) {
    // 88 followed by 01XXXXXXXXX — missing a zero: treat as typo
    e164Digits = "880" + digits.slice(2);
  } else if (digits.startsWith("0")) {
    // 01712345678 (11 digits local)
    e164Digits = "880" + digits.slice(1);
  } else if (digits.length === 10) {
    // 1712345678 — 10 digits without leading 0
    e164Digits = "880" + digits;
  } else {
    return null;
  }

  if (!BD_PHONE_RE.test(e164Digits)) {
    return null;
  }

  return `+${e164Digits}`;
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

  const normalizedMobile = normalizeMobileNumber(mobileRaw);
  if (!normalizedMobile) {
    return NextResponse.json(
      {
        ok: false,
        message:
          "Invalid mobile number. Please enter a valid Bangladesh number " +
          "(e.g. 01712345678, +8801712345678 or 008801712345678).",
      },
      { status: 400 }
    );
  }

  const result = await loginWithMobile({
    mobileNumber: normalizedMobile,
    password: payload.password,
    deviceId: payload.deviceId,
  });

  if (result.error) {
    const httpStatus = result.response?.status || 500;
    return NextResponse.json(
      { ok: false, message: toUserFriendlyLoginError(result.error, httpStatus) },
      { status: httpStatus }
    );
  }

  const response = NextResponse.json({ ok: true }, { status: result.response.status || 200 });
  const setCookieHeaders = getSetCookieHeaders(result.response.headers);
  const backendSessionCookie = setCookieHeaders.find((value) =>
    value.startsWith(`${SESSION_COOKIE_NAME}=`)
  );
  const sessionToken = backendSessionCookie
    ? extractCookieValue(backendSessionCookie, SESSION_COOKIE_NAME)
    : undefined;
  if (sessionToken) {
    response.cookies.set({
      name: SESSION_COOKIE_NAME,
      value: sessionToken,
      path: "/",
      httpOnly: true,
      sameSite: "strict",
      secure: process.env.NODE_ENV === "production",
      maxAge: 60 * 60 * 12,
    });
  }

  const csrfToken = result.response.headers.get("x-csrf-token") ?? result.data?.csrf_token;
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
  const session = await toPortalSessionFromLogin(result.data ?? {}, sessionCookieHeader);
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
  let portalUserId = session.principal.user?.userId ?? (result.data as Record<string, unknown>)?.user_id as string ?? "";

  // If user_id is still empty, call getCurrentSession using the new session token to resolve it.
  // This handles gateways that don't return user_id in the login JSON body.
  if (!portalUserId && sessionToken) {
    try {
      const { authServiceGetCurrentSession, createInsureTechClient } = await import("@lifeplus/insuretech-sdk");
      const { getApiBaseUrl } = await import("@lib/sdk/api-helpers");
      const { toPortalSessionFromCurrentSession } = await import("@lib/auth/backend-auth");
      const tempClient = createInsureTechClient({ baseUrl: getApiBaseUrl(), apiKey: process.env.INSURETECH_API_KEY ?? "" });
      const cookieStr = `${SESSION_COOKIE_NAME}=${sessionToken}`;
      const sessionRes = await authServiceGetCurrentSession({
        client: tempClient,
        headers: { Cookie: cookieStr },
        throwOnError: false,
      });
      if (sessionRes.response.ok && sessionRes.data) {
        const portalSession = await toPortalSessionFromCurrentSession(sessionRes.data, cookieStr);
        portalUserId = portalSession?.principal?.user?.userId ?? "";
      }
    } catch { /* ignore — userId will remain empty */ }
  }
  const cookieOpts = {
    path: "/",
    httpOnly: false, // must be readable by edge middleware + session-headers helper
    sameSite: "strict" as const,
    secure: process.env.NODE_ENV === "production",
    maxAge: 60 * 60 * 12,
  };

  finalResponse.cookies.set({ name: "portal_role", value: portalRole, ...cookieOpts });
  finalResponse.cookies.set({ name: "portal_user_id", value: portalUserId, ...cookieOpts });
  finalResponse.cookies.set({ name: "portal_biz_id", value: portalBizId, ...cookieOpts });

  // Store user contact info cookies so the My Profile page can display
  // mobile_number and email — these live on the User record, not UserProfile.
  // They are read-only identity fields (auth credentials), not profile fields.
  const portalMobile = (result.data as Record<string, unknown>)?.user
    ? ((result.data as Record<string, unknown>).user as Record<string, unknown>)?.mobile_number as string ?? ""
    : "";
  const portalEmail = (result.data as Record<string, unknown>)?.user
    ? ((result.data as Record<string, unknown>).user as Record<string, unknown>)?.email as string ?? ""
    : "";
  finalResponse.cookies.set({ name: "portal_mobile", value: portalMobile, ...cookieOpts });
  finalResponse.cookies.set({ name: "portal_email", value: portalEmail, ...cookieOpts });

  return finalResponse;
}
