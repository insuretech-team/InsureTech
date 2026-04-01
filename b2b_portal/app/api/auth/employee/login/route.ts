import { NextResponse } from "next/server";

import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { SESSION_COOKIE_NAME } from "@lib/auth/session";
import {
  getSetCookieHeaders,
  toPortalSessionFromCurrentSession,
  toPortalSessionFromLogin,
} from "@lib/auth/backend-auth";
import type { EmployeePortalLoginRequest } from "@lib/types/auth";

const CSRF_COOKIE_NAME = "csrf_token";

function getApiBaseUrl(): string {
  return (
    process.env.INSURETECH_API_BASE_URL ??
    process.env.NEXT_PUBLIC_INSURETECH_API_BASE_URL ??
    "http://localhost:8080"
  );
}

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

function toUserFriendlyEmployeeLoginError(error: unknown, httpStatus: number): string {
  let raw = "";
  if (typeof error === "string") {
    raw = error;
  } else if (error && typeof error === "object") {
    for (const key of ["message", "error", "detail", "description"] as const) {
      const value = (error as Record<string, unknown>)[key];
      if (typeof value === "string" && value.trim()) {
        raw = value;
        break;
      }
    }
  }
  const lower = raw.toLowerCase();

  if (
    lower.includes("locked") ||
    lower.includes("too many") ||
    lower.includes("rate limit") ||
    httpStatus === 422 ||
    httpStatus === 429
  ) {
    return raw.trim() || "Your account is temporarily locked. Please try again later.";
  }

  if (httpStatus === 401 || lower.includes("invalid credentials") || lower.includes("invalid password")) {
    return "Email or password is incorrect. Please try again.";
  }

  if (lower.includes("not active") || lower.includes("inactive")) {
    return "Your employee access is not active yet. Complete activation first.";
  }

  if (httpStatus >= 500 || lower.includes("internal") || lower.includes("unavailable")) {
    return "The sign-in service is temporarily unavailable. Please try again in a moment.";
  }

  return raw.trim() || "Employee sign-in failed.";
}

export async function POST(request: Request) {
  let payload: EmployeePortalLoginRequest;
  try {
    payload = (await request.json()) as EmployeePortalLoginRequest;
  } catch {
    return NextResponse.json({ ok: false, message: "Invalid login payload" }, { status: 400 });
  }

  const email = payload.email?.trim().toLowerCase();
  if (!email) {
    return NextResponse.json({ ok: false, message: "Email is required" }, { status: 400 });
  }
  if (!payload.password?.trim()) {
    return NextResponse.json({ ok: false, message: "Password is required" }, { status: 400 });
  }

  let backendResponse: Response;
  let rawBackendBody = "";
  try {
    backendResponse = await fetch(`${getApiBaseUrl()}/v1/auth/email-password:login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        email,
        password: payload.password,
        device_id: payload.deviceId ?? "b2b-employee-web",
        device_type: "WEB",
        device_name: "B2B Employee Portal Web",
      }),
      cache: "no-store",
    });
    rawBackendBody = await backendResponse.text();
  } catch (error) {
    return NextResponse.json(
      { ok: false, message: toUserFriendlyEmployeeLoginError(error, 502) },
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
      backendResponse.status ??
      500;
    return NextResponse.json(
      {
        ok: false,
        message: toUserFriendlyEmployeeLoginError(backendError ?? rawBackendBody, httpStatus),
      },
      { status: httpStatus }
    );
  }

  const loginData =
    backendPayload && typeof backendPayload.data === "object"
      ? (backendPayload.data as Record<string, unknown>)
      : backendPayload;

  let sessionToken =
    typeof loginData.session_token === "string" && loginData.session_token
      ? loginData.session_token
      : undefined;

  if (!sessionToken) {
    const setCookieHeaders = getSetCookieHeaders(backendResponse.headers);
    const backendSessionCookie = setCookieHeaders.find((value) =>
      value.startsWith(`${SESSION_COOKIE_NAME}=`)
    );
    sessionToken = backendSessionCookie
      ? extractCookieValue(backendSessionCookie, SESSION_COOKIE_NAME)
      : undefined;
  }

  const response = NextResponse.json({ ok: true }, { status: backendResponse.status || 200 });
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
    (typeof loginData.csrf_token === "string" ? loginData.csrf_token : undefined);
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
  const session = await toPortalSessionFromLogin(loginData ?? {}, sessionCookieHeader);

  if (sessionToken && session?.principal) {
    const { createSession } = await import("@lib/auth/session-store");
    createSession(session.principal, sessionToken);
  }

  let portalUserId = session.principal.user?.userId ?? "";
  if (!portalUserId && sessionToken) {
    try {
      const cookieStr = `${SESSION_COOKIE_NAME}=${sessionToken}`;
      const tempReq = new Request(request.url, { headers: { cookie: cookieStr } });
      const tempSdk = makeSdkClient(tempReq);
      const currentSessionResult = await tempSdk.getCurrentSession();
      if (currentSessionResult.response?.ok && currentSessionResult.data) {
        const portalSession = await toPortalSessionFromCurrentSession(currentSessionResult.data, cookieStr);
        portalUserId = portalSession?.principal.user?.userId ?? "";
      }
    } catch {
      portalUserId = session.principal.user?.userId ?? "";
    }
  }

  const cookieOpts = {
    path: "/",
    httpOnly: false,
    sameSite: "lax" as const,
    secure: process.env.NODE_ENV === "production",
    maxAge: 60 * 60 * 12,
  };

  const finalResponse = NextResponse.json({ ok: true, session }, { status: response.status });
  for (const cookie of response.cookies.getAll()) {
    finalResponse.cookies.set(cookie);
  }

  finalResponse.cookies.set({ name: "portal_role", value: "B2B_BENEFICIARY", ...cookieOpts });
  finalResponse.cookies.set({ name: "portal_user_id", value: portalUserId, ...cookieOpts });
  finalResponse.cookies.set({ name: "portal_biz_id", value: session.principal.businessId ?? "", ...cookieOpts });
  finalResponse.cookies.set({ name: "portal_password_change_required", value: "false", ...cookieOpts });
  finalResponse.cookies.set({ name: "portal_kyc_verified", value: "true", ...cookieOpts });
  finalResponse.cookies.set({ name: "portal_mobile", value: "", ...cookieOpts });
  finalResponse.cookies.set({
    name: "portal_email",
    value: session.principal.user?.email ?? email,
    ...cookieOpts,
  });

  return finalResponse;
}
