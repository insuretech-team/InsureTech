import { NextResponse } from "next/server";

import { makeSdkClient } from "@lib/sdk/partner-sdk-client";
import { getSetCookieHeaders } from "@lib/auth/backend-auth";
import { SESSION_COOKIE_NAME } from "@lib/auth/session";

const CSRF_COOKIE_NAME = "csrf_token";

function getCookieValue(cookieHeader: string, cookieName: string): string | undefined {
  const target = `${cookieName}=`;
  for (const rawPart of cookieHeader.split(";")) {
    const part = rawPart.trim();
    if (part.startsWith(target)) {
      return decodeURIComponent(part.slice(target.length));
    }
  }
  return undefined;
}

function expireSessionCookie(response: NextResponse) {
  response.cookies.set({
    name: SESSION_COOKIE_NAME,
    value: "",
    path: "/",
    httpOnly: true,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    expires: new Date(0),
  });
}

function expireCsrfCookie(response: NextResponse) {
  response.cookies.set({
    name: CSRF_COOKIE_NAME,
    value: "",
    path: "/",
    httpOnly: true,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    expires: new Date(0),
  });
}

function expirePortalCookies(response: NextResponse) {
  for (const name of ["portal_role", "portal_user_id", "portal_partner_id", "portal_mobile", "portal_email"]) {
    response.cookies.set({
      name,
      value: "",
      path: "/",
      httpOnly: false,
      sameSite: "lax",
      secure: process.env.NODE_ENV === "production",
      expires: new Date(0),
    });
  }
}

export async function POST(request: Request) {
  const cookieHeader = request.headers.get("cookie") ?? "";
  const csrfToken = getCookieValue(cookieHeader, CSRF_COOKIE_NAME);
  const sdk = makeSdkClient(request);
  let sessionId = "";

  try {
    const currentSessionResult = await sdk.getCurrentSession();
    if (!currentSessionResult.error) {
      sessionId = currentSessionResult.data?.session?.session_id ?? "";
    }
  } catch {
    // Ignore session lookup failures; we can still clear local cookies.
  }

  if (!sessionId) {
    const response = NextResponse.json({ ok: true, message: "No active session" }, { status: 200 });
    expireSessionCookie(response);
    expireCsrfCookie(response);
    expirePortalCookies(response);
    return response;
  }

  let result: Awaited<ReturnType<typeof sdk.logout>>;
  try {
    result = await sdk.logout({
      body: {
        session_id: sessionId,
        logout_reason: "user_initiated",
      },
    });
  } catch (error) {
    const msg = error instanceof Error ? error.message : "Logout failed";
    const response = NextResponse.json({ ok: false, message: msg }, { status: 502 });
    expireSessionCookie(response);
    expireCsrfCookie(response);
    return response;
  }

  if (!result.response.ok) {
    const status = result.response.status ?? 500;
    const errPayload = "error" in result ? result.error as Record<string, unknown> | undefined : undefined;
    const errMsg = typeof errPayload?.message === "string" ? errPayload.message : "Logout failed";
    const response = NextResponse.json({ ok: false, message: errMsg }, { status });
    for (const setCookie of getSetCookieHeaders(result.response.headers)) {
      response.headers.append("set-cookie", setCookie);
    }
    expireSessionCookie(response);
    expireCsrfCookie(response);
    return response;
  }

  const response = NextResponse.json({ ok: true }, { status: result.response.status || 200 });
  for (const setCookie of getSetCookieHeaders(result.response.headers)) {
    response.headers.append("set-cookie", setCookie);
  }
  expireSessionCookie(response);
  expireCsrfCookie(response);
  expirePortalCookies(response);
  return response;
}
