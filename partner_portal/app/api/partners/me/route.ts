import { NextResponse } from "next/server";
import { makeDirectHttp } from "@lib/sdk/partner-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";

/** GET /api/partners/me - Get current partner's organization data */
export async function GET(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

  // Use makeDirectHttp for endpoints not yet in SDK
  const http = makeDirectHttp(request);
  const result = await http.get("/v1/partners/me");

  if (!result.ok) {
    return NextResponse.json(
      { ok: false, message: result.error || "Failed to fetch partner data" },
      { status: result.status }
    );
  }

  return NextResponse.json({ ok: true, data: result.data }, { status: 200 });
}
