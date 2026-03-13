import { NextResponse } from "next/server";

import {
  getCurrentSession,
  getErrorMessage,
  getSetCookieHeaders,
  logoutCurrentSession,
} from "@lib/auth/backend-auth";
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
  for (const name of ["portal_role", "portal_user_id", "portal_biz_id"]) {
    response.cookies.set({
      name,
      value: "",
      path: "/",
      httpOnly: false,
      sameSite: "strict",
      secure: process.env.NODE_ENV === "production",
      expires: new Date(0),
    });
  }
}

export async function POST(request: Request) {
  const cookieHeader = request.headers.get("cookie") ?? "";
  const csrfToken = getCookieValue(cookieHeader, CSRF_COOKIE_NAME);
  let sessionId = "";

  try {
    const currentSessionResult = await getCurrentSession(cookieHeader);
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

  let result: Awaited<ReturnType<typeof logoutCurrentSession>>;
  try {
    result = await logoutCurrentSession(cookieHeader, csrfToken, sessionId);
  } catch (error) {
    const response = NextResponse.json(
      { ok: false, message: getErrorMessage(error, "Logout failed") },
      { status: 502 }
    );
    expireSessionCookie(response);
    expireCsrfCookie(response);
    return response;
  }

  if (result.error) {
    const status = result.response?.status ?? 500;
    const response = NextResponse.json(
      { ok: false, message: getErrorMessage(result.error, "Logout failed") },
      { status }
    );
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
