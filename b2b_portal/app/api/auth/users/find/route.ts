import { NextResponse } from "next/server";

import { makeDirectHttp } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";

function readCookie(cookieHeader: string, name: string): string {
  const escaped = name.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  const match = cookieHeader.match(new RegExp(`(?:^|;\\s*)${escaped}=([^;]*)`));
  return match ? decodeURIComponent(match[1]) : "";
}

export async function POST(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) {
    return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  }

  const role = readCookie(request.headers.get("cookie") ?? "", "portal_role");
  if (role !== "SYSTEM_ADMIN") {
    return NextResponse.json({ ok: false, message: "Forbidden" }, { status: 403 });
  }

  let body: { identifier?: string };
  try {
    body = await request.json();
  } catch {
    return NextResponse.json({ ok: false, message: "Invalid request body" }, { status: 400 });
  }

  const identifier = body.identifier?.trim();
  if (!identifier) {
    return NextResponse.json({ ok: false, message: "identifier is required" }, { status: 400 });
  }

  const result = await makeDirectHttp(request, hdrs).post("/v1/auth/users:find", { identifier });
  if (!result.ok) {
    return NextResponse.json(
      { ok: false, message: result.data?.message ?? "User not found" },
      { status: result.status }
    );
  }

  const payload = (result.data ?? {}) as Record<string, unknown>;
  const user = (payload.user ?? payload) as Record<string, unknown>;

  return NextResponse.json({
    ok: true,
    user: {
      userId: typeof user.user_id === "string" ? user.user_id : "",
      fullName: typeof user.full_name === "string" ? user.full_name : "",
      email: typeof user.email === "string" ? user.email : "",
      mobileNumber: typeof user.mobile_number === "string" ? user.mobile_number : "",
      userType: typeof user.user_type === "string" ? user.user_type : "",
      emailVerified: Boolean(user.email_verified),
      kycVerified: Boolean(user.kyc_verified),
      passwordChangeRequired: Boolean(user.password_change_required),
    },
  });
}
