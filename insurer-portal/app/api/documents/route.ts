import { NextResponse } from "next/server";

import { directHttp } from "@/lib/server/insuretech";
import { mapLiveDocument } from "@/lib/server/mappers";
import { loadContext } from "@/lib/server/portal-data";

export async function GET(request: Request) {
  const searchParams = new URL(request.url).searchParams;
  const insurerId = searchParams.get("insurerId") ?? "";
  const entityType = searchParams.get("entityType") ?? "";
  const entityId = searchParams.get("entityId") ?? "";

  const context = await loadContext(request, insurerId);
  if (!context) {
    return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  }

  const query = new URLSearchParams({ page_size: "100" });
  if (entityType) query.set("entity_type", entityType);
  if (entityId) query.set("entity_id", entityId);

  const result = await directHttp(request, `/v1/documents?${query.toString()}`, {
    session: context.session,
  });

  if (!result.ok) {
    return NextResponse.json({ ok: false, message: result.message ?? "Failed to list documents", data: [] });
  }

  const docs = ((result.data.documents as unknown[]) ?? []).map((entry) => mapLiveDocument(entry));

  return NextResponse.json({ ok: true, data: docs });
}
