import { NextResponse } from "next/server";

import { makeSdkClient } from "@lib/sdk/partner-sdk-client";
import { SESSION_COOKIE_NAME } from "@lib/auth/session";
import { toPortalSessionFromLogin, getSetCookieHeaders } from "@lib/auth/backend-auth";
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

  // Use makeSdkClient with no session overrides (public endpoint — no session cookie yet).
  // sdk.login() maps to authServiceLogin → POST /v1/auth/login (mobile_number + password).
  // Do NOT use emailLogin — that is OTP-based (requires otp_id + code, not password).
  const sdk = makeSdkClient(request);
  const result = await sdk.mobileLogin({
    body: {
      mobile_number: normalizedMobile,
      password: payload.password,
      device_id: payload.deviceId ?? "partner-portal-web",
      device_type: "WEB",
      device_name: "Partner Portal Web",
    },
  });

  if (!result.response.ok) {
    const httpStatus = result.response.status || 500;
    const errPayload = "error" in result ? result.error : undefined;
    return NextResponse.json(
      { ok: false, message: toUserFriendlyLoginError(errPayload, httpStatus) },
      { status: httpStatus }
    );
  }

  const response = NextResponse.json({ ok: true }, { status: result.response.status || 200 });

  // Primary: read session_token from the JSON response body.
  // The gateway Login handler keeps session_token in the proto JSON body (in addition
  // to setting the HttpOnly cookie) specifically so the Next.js BFF can read it here
  // without relying on Set-Cookie header forwarding.
  // Set-Cookie is a forbidden header in the Fetch API Headers constructor — the SDK
  // interceptor (which rewrites the response body) silently drops it, so we cannot
  // rely on result.response.headers.get('set-cookie') or getSetCookieHeaders().
  const loginData = result.data as Record<string, unknown> | undefined;
  let sessionToken: string | undefined =
    typeof loginData?.session_token === "string" && loginData.session_token
      ? loginData.session_token
      : undefined;

  // Fallback: try reading from Set-Cookie headers (works when no SDK interceptor is involved).
  if (!sessionToken) {
    const setCookieHeaders = getSetCookieHeaders(result.response.headers);
    const backendSessionCookie = setCookieHeaders.find((value: string) =>
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
  // loginData is the unwrapped LoginResponse (SDK interceptor strips the ApiResponse envelope).
  const session = await toPortalSessionFromLogin(loginData ?? {}, sessionCookieHeader);
  const finalResponse = NextResponse.json({ ok: true, session }, { status: response.status });
  for (const cookie of response.cookies.getAll()) {
    finalResponse.cookies.set(cookie);
  }

  // Set lightweight metadata cookies used by:
  //   1. Edge middleware (portal_role) — for role-based page routing without hitting DB
  //   2. API route session-headers helper — to inject x-portal/x-user-id/x-partner-id
  //      into backend SDK calls so the Casbin authz interceptor gets the right domain.
  //
  // These are NOT security boundaries — the backend session cookie enforces real authz.
  const portalRole = session.principal.role ?? "PARTNER_ADMIN";
  const portalPartnerId = session.principal.partnerId ?? "";

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
      if (sessionRes.response.ok && sessionRes.data) {
        const portalSession = await toPortalSessionFromCurrentSession(sessionRes.data, cookieStr);
        portalUserId = portalSession?.principal?.user?.userId ?? "";
      }
    } catch { /* ignore — userId will remain empty */ }
  }
  const cookieOpts = {
    path: "/",
    httpOnly: false, // must be readable by edge middleware + session-headers helper
    sameSite: "lax" as const,
    secure: process.env.NODE_ENV === "production",
    maxAge: 60 * 60 * 12,
  };

  finalResponse.cookies.set({ name: "portal_role", value: portalRole, ...cookieOpts });
  finalResponse.cookies.set({ name: "portal_user_id", value: portalUserId, ...cookieOpts });
  finalResponse.cookies.set({ name: "portal_partner_id", value: portalPartnerId, ...cookieOpts });

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
