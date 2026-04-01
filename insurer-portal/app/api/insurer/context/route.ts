import { NextResponse } from "next/server";

import { loadContext } from "@/lib/server/portal-data";

export async function GET(request: Request) {
  const insurerId = new URL(request.url).searchParams.get("insurerId") ?? "";
  const context = await loadContext(request, insurerId);

  if (!context) {
    return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  }

  return NextResponse.json({
    ok: true,
    data: {
      session: context.session,
      insurers: context.insurers,
      currentInsurer: context.currentInsurer,
      config: context.config,
      products: context.products,
      source: context.source,
    },
  });
}
