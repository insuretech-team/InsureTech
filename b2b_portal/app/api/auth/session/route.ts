import { NextResponse } from "next/server";

import {
  getCurrentSession,
  getErrorMessage,
  toPortalSessionFromCurrentSession,
} from "@lib/auth/backend-auth";
import { SESSION_COOKIE_NAME } from "@lib/auth/session";

export async function GET(request: Request) {
  try {
    const cookieHeader = request.headers.get("cookie") ?? "";
    const hasSessionCookie = cookieHeader.includes(`${SESSION_COOKIE_NAME}=`);
    if (!hasSessionCookie) {
      return NextResponse.json({ ok: false, message: "No active session" }, { status: 401 });
    }

    let result: Awaited<ReturnType<typeof getCurrentSession>>;
    try {
      result = await getCurrentSession(cookieHeader);
    } catch (error) {
      return NextResponse.json(
        { ok: false, message: getErrorMessage(error, "Session service unavailable") },
        { status: 502 }
      );
    }

    if (result.error) {
      const status = result.response?.status ?? 401;
      return NextResponse.json(
        { ok: false, message: getErrorMessage(result.error, "No active session") },
        { status }
      );
    }

    const session = await toPortalSessionFromCurrentSession(result.data ?? {}, cookieHeader);
    if (!session) {
      return NextResponse.json({ ok: false, message: "No active session" }, { status: 401 });
    }

    // Re-mint metadata cookies on every session refresh so they stay in sync
    // with the backend session. This is critical: if the lightweight cookies
    // (portal_role, portal_user_id, portal_biz_id) expire or are cleared while
    // the backend session_token is still valid, session-headers.ts would fall
    // back to PORTAL_B2B with no org context, causing 403 for superadmin on
    // org/dept tabs.
    const response = NextResponse.json({ ok: true, session }, { status: 200 });
    const portalRole = session.principal.role ?? "BUSINESS_ADMIN";
    const portalUserId = session.principal.user?.userId ?? "";
    const portalBizId = session.principal.businessId ?? "";
    const cookieOpts = {
      path: "/",
      httpOnly: false, // must be readable by edge middleware + session-headers helper
      sameSite: "strict" as const,
      secure: process.env.NODE_ENV === "production",
      maxAge: 60 * 60 * 12,
    };
    response.cookies.set({ name: "portal_role",    value: portalRole,    ...cookieOpts });
    response.cookies.set({ name: "portal_user_id", value: portalUserId,  ...cookieOpts });
    response.cookies.set({ name: "portal_biz_id",  value: portalBizId,   ...cookieOpts });

    // Re-mint contact info cookies on every session refresh so they stay in sync.
    // These are sourced from existing cookies (set at login) — if empty, preserve
    // whatever was already in the browser (no-op by writing empty string is safe).
    const existingCookieHeader = request.headers.get("cookie") ?? "";
    const extractCk = (name: string) => {
      const m = existingCookieHeader.match(new RegExp(`(?:^|;\\s*)${name}=([^;]*)`));
      return m ? decodeURIComponent(m[1]) : "";
    };
    const portalMobile = extractCk("portal_mobile");
    const portalEmail  = extractCk("portal_email");
    response.cookies.set({ name: "portal_mobile", value: portalMobile, ...cookieOpts });
    response.cookies.set({ name: "portal_email",  value: portalEmail,  ...cookieOpts });

    return response;
  } catch (error) {
    return NextResponse.json(
      { ok: false, message: getErrorMessage(error, "Session endpoint failed") },
      { status: 502 }
    );
  }
}
