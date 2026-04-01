import { NextResponse } from "next/server";

import { directHttp, fetchCurrentSession } from "@/lib/server/insuretech";

function amountToMoney(value?: number) {
  if (!value || Number.isNaN(value)) {
    return { amount: 0, currency: "BDT" };
  }
  return { amount: Math.round(value * 100), currency: "BDT" };
}

export async function POST(
  request: Request,
  { params }: { params: Promise<{ claimId: string }> },
) {
  const session = await fetchCurrentSession(request);
  if (!session) {
    return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  }

  const { claimId } = await params;
  const payload = (await request.json().catch(() => null)) as
    | {
        action?: "approve" | "reject" | "settle";
        amount?: number;
        reason?: string;
        paymentReference?: string;
      }
    | null;

  if (!claimId || !payload?.action) {
    return NextResponse.json({ ok: false, message: "Missing claim action." }, { status: 400 });
  }

  const path = `/v1/claims/${claimId}/${payload.action}`;
  const body =
    payload.action === "approve"
      ? {
          claim_id: claimId,
          approver_id: session.userId,
          approved_amount: amountToMoney(payload.amount),
          notes: payload.reason ?? "",
        }
      : payload.action === "reject"
        ? {
            claim_id: claimId,
            approver_id: session.userId,
            reason: payload.reason ?? "Rejected from insurer portal",
          }
        : {
            claim_id: claimId,
            payment_method: "BANK_TRANSFER",
            payment_reference: payload.paymentReference ?? "",
          };

  const result = await directHttp(request, path, {
    method: "POST",
    session,
    body,
  });

  if (!result.ok) {
    return NextResponse.json(
      { ok: false, message: result.message ?? `Unable to ${payload.action} this claim.` },
      { status: result.status || 502 },
    );
  }

  return NextResponse.json({ ok: true, data: { updated: true } });
}
