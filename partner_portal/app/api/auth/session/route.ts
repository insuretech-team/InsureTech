import { NextResponse } from "next/server";

import { makeSdkClient } from "@lib/sdk/partner-sdk-client";
import { toPortalSessionFromCurrentSession } from "@lib/auth/backend-auth";
import { SESSION_COOKIE_NAME } from "@lib/auth/session";

export async function GET(request: Request) {
  try {
    const cookieHeader = request.headers.get("cookie") ?? "";
    const hasSessionCookie = cookieHeader.includes(`${SESSION_COOKIE_NAME}=`);
    if (!hasSessionCookie) {
      return NextResponse.json({ ok: false, message: "No active session" }, { status: 401 });
    }

    const sdk = makeSdkClient(request);
    let result: Awaited<ReturnType<typeof sdk.getCurrentSession>>;
    try {
      result = await sdk.getCurrentSession();
    } catch (error) {
      const msg = error instanceof Error ? error.message : "Session service unavailable";
      return NextResponse.json({ ok: false, message: msg }, { status: 502 });
    }

    if (!result.response.ok) {
      const status = result.response.status ?? 401;
      const errPayload = "error" in result ? result.error as Record<string, unknown> | undefined : undefined;
      const errMsg = typeof errPayload?.message === "string" ? errPayload.message : "No active session";
      return NextResponse.json({ ok: false, message: errMsg }, { status });
    }

    const session = await toPortalSessionFromCurrentSession(result.data ?? {}, cookieHeader);
    if (!session) {
      return NextResponse.json({ ok: false, message: "No active session" }, { status: 401 });
    }

    // Re-mint metadata cookies on every session refresh
    const response = NextResponse.json({ ok: true, session }, { status: 200 });
    const portalRole = session.principal.role ?? "PARTNER_ADMIN";
    const portalUserId = session.principal.user?.userId ?? "";
    const portalPartnerId = session.principal.partnerId ?? "";
    const cookieOpts = {
      path: "/",
      httpOnly: false,
      sameSite: "lax" as const,
      secure: process.env.NODE_ENV === "production",
      maxAge: 60 * 60 * 12,
    };
    response.cookies.set({ name: "portal_role",       value: portalRole,       ...cookieOpts });
    response.cookies.set({ name: "portal_user_id",    value: portalUserId,     ...cookieOpts });
    response.cookies.set({ name: "portal_partner_id", value: portalPartnerId, ...cookieOpts });

    // Re-mint contact info cookies
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
    const msg = error instanceof Error ? error.message : "Session endpoint failed";
    return NextResponse.json({ ok: false, message: msg }, { status: 502 });
  }
}
