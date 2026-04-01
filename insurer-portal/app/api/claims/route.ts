import { NextResponse } from "next/server";

import { directHttp } from "@/lib/server/insuretech";
import { mapClaim } from "@/lib/server/mappers";
import { loadContext } from "@/lib/server/portal-data";
import { fallbackClaims } from "@/lib/mock-data";

export async function GET(request: Request) {
  const searchParams = new URL(request.url).searchParams;
  const insurerId = searchParams.get("insurerId") ?? "";
  const context = await loadContext(request, insurerId);

  if (!context) {
    return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  }

  const result = await directHttp(request, "/v1/claims?page_size=100", {
    session: context.session,
  });

  const claims = result.ok
    ? (((result.data.claims as unknown[]) ?? []).map((entry) => mapClaim(entry, "live")))
    : fallbackClaims;

  return NextResponse.json({
    ok: true,
    data: claims,
    message: result.ok ? undefined : "Showing fallback claim data while the live claim queue is unavailable.",
  });
}
