import { NextResponse } from "next/server";

import { directHttp } from "@/lib/server/insuretech";
import { loadContext } from "@/lib/server/portal-data";

export async function GET(request: Request) {
  const searchParams = new URL(request.url).searchParams;
  const insurerId = searchParams.get("insurerId") ?? "";
  const type = searchParams.get("type") ?? "";

  const context = await loadContext(request, insurerId);
  if (!context) {
    return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  }

  const query = new URLSearchParams({ page_size: "100", active_only: "true" });
  if (type) query.set("type", type);

  const result = await directHttp(request, `/v1/documents/templates?${query.toString()}`, {
    session: context.session,
  });

  if (!result.ok) {
    return NextResponse.json({ ok: false, message: result.message ?? "Failed to list templates", data: [] });
  }

  const templates = (result.data.templates as unknown[]) ?? [];

  return NextResponse.json({ ok: true, data: templates });
}
