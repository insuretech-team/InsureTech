import { NextResponse } from "next/server";

import {
  applyPortalCookies,
  buildSessionFromLogin,
  loginWithMobile,
  normalizeMobileNumber,
} from "@/lib/server/insuretech";

export async function POST(request: Request) {
  const payload = (await request.json().catch(() => null)) as
    | { mobileNumber?: string; password?: string }
    | null;

  const normalized = normalizeMobileNumber(payload?.mobileNumber ?? "");
  if (!normalized || !payload?.password?.trim()) {
    return NextResponse.json(
      { ok: false, message: "A valid Bangladesh mobile number and password are required." },
      { status: 400 },
    );
  }

  const result = await loginWithMobile(request, {
    mobileNumber: normalized,
    password: payload.password,
  });

  if (!result.response.ok) {
    const status = result.response.status || 401;
    const errorRecord = (result.error as Record<string, unknown> | undefined) ?? {};
    const message =
      (typeof errorRecord.message === "string" && errorRecord.message) ||
      "Sign in failed. Please verify your credentials.";
    return NextResponse.json({ ok: false, message }, { status });
  }

  const { session, sessionToken, csrfToken } = await buildSessionFromLogin(
    request,
    result.data ?? {},
    result.response.headers,
  );

  const loginRecord = (result.data as Record<string, unknown> | undefined) ?? {};
  const userRecord =
    loginRecord.user && typeof loginRecord.user === "object"
      ? (loginRecord.user as Record<string, unknown>)
      : {};

  const response = NextResponse.json({ ok: true, data: { session } });
  applyPortalCookies(response, session, {
    sessionToken,
    csrfToken,
    email: (userRecord.email as string | undefined) ?? "",
    mobile: (userRecord.mobile_number as string | undefined) ?? "",
  });

  return response;
}
