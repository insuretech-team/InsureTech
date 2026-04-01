import { NextResponse } from "next/server";

import {
  clearPortalCookies,
  fetchCurrentSession,
  logoutCurrentSession,
} from "@/lib/server/insuretech";

export async function POST(request: Request) {
  const session = await fetchCurrentSession(request);
  await logoutCurrentSession(request, session).catch(() => null);

  const response = NextResponse.json({ ok: true, data: null });
  clearPortalCookies(response);
  return response;
}
