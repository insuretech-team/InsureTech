import { NextResponse } from "next/server";

import { loadOverview } from "@/lib/server/portal-data";

export async function GET(request: Request) {
  const insurerId = new URL(request.url).searchParams.get("insurerId") ?? "";
  const overview = await loadOverview(request, insurerId);

  if (!overview) {
    return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  }

  return NextResponse.json({ ok: true, data: overview });
}
