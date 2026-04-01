import { NextResponse } from "next/server";

import { directHttp } from "@/lib/server/insuretech";
import { mapProposal } from "@/lib/server/mappers";
import { loadContext } from "@/lib/server/portal-data";
import { fallbackProposals } from "@/lib/mock-data";

export async function GET(request: Request) {
  const searchParams = new URL(request.url).searchParams;
  const insurerId = searchParams.get("insurerId") ?? "";
  const status = searchParams.get("status") ?? "";
  const context = await loadContext(request, insurerId);

  if (!context) {
    return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  }

  const query = new URLSearchParams({
    insurer_id: context.currentInsurer.id,
    page_size: "100",
  });
  if (status) query.set("status", status);

  const result = await directHttp(
    request,
    `/v1/insurance-proposals?${query.toString()}`,
    { session: context.session },
  );

  const proposals = result.ok
    ? (((result.data.proposals as unknown[]) ?? []).map((entry) => mapProposal(entry, "live")))
    : fallbackProposals;

  return NextResponse.json({
    ok: true,
    data: proposals,
    message: result.ok ? undefined : "Showing fallback proposal data while the live list is unavailable.",
  });
}
