import { NextResponse } from "next/server";

import { directHttp, fetchCurrentSession } from "@/lib/server/insuretech";

export async function POST(
  request: Request,
  { params }: { params: Promise<{ proposalId: string }> },
) {
  const session = await fetchCurrentSession(request);
  if (!session) {
    return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  }

  const { proposalId } = await params;
  const payload = (await request.json().catch(() => null)) as
    | { action?: "approve" | "reject"; reason?: string }
    | null;

  if (!proposalId || !payload?.action) {
    return NextResponse.json({ ok: false, message: "Missing proposal action." }, { status: 400 });
  }

  const result = await directHttp(
    request,
    `/v1/insurance-proposals/${proposalId}/${payload.action}`,
    {
      method: "POST",
      session,
      body:
        payload.action === "reject"
          ? { reason: payload.reason ?? "", reviewed_by_user_id: session.userId }
          : { reviewed_by_user_id: session.userId },
    },
  );

  if (!result.ok) {
    return NextResponse.json(
      { ok: false, message: result.message ?? `Unable to ${payload.action} the proposal.` },
      { status: result.status || 502 },
    );
  }

  return NextResponse.json({ ok: true, data: { updated: true } });
}
