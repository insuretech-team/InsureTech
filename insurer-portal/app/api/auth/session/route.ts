import { NextResponse } from "next/server";

import { applyPortalCookies, fetchCurrentSession } from "@/lib/server/insuretech";

export async function GET(request: Request) {
  const session = await fetchCurrentSession(request);
  if (!session) {
    return NextResponse.json({ ok: false, message: "No active session" }, { status: 401 });
  }

  const response = NextResponse.json({ ok: true, data: { session } });
  applyPortalCookies(response, session);
  return response;
}
