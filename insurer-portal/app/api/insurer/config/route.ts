import { NextResponse } from "next/server";

import { directHttp, fetchCurrentSession } from "@/lib/server/insuretech";

export async function PATCH(request: Request) {
  const session = await fetchCurrentSession(request);
  if (!session) {
    return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  }

  const payload = (await request.json().catch(() => null)) as
    | {
        insurerId?: string;
        apiBaseUrl?: string;
        authType?: string;
        authCredentials?: string;
        businessModel?: string;
        autoUnderwritingEnabled?: boolean;
        realTimeClaimNotification?: boolean;
        paymentTerms?: string;
        claimSettlementDays?: number;
      }
    | null;

  if (!payload?.insurerId) {
    return NextResponse.json({ ok: false, message: "insurerId is required" }, { status: 400 });
  }

  const result = await directHttp(
    request,
    `/v1/insurers/${payload.insurerId}/config`,
    {
      method: "PUT",
      session,
      body: {
        insurer_id: payload.insurerId,
        api_base_url: payload.apiBaseUrl ?? "",
        auth_type: payload.authType ?? "",
        auth_credentials: payload.authCredentials ?? "",
        business_model: payload.businessModel ?? "",
        auto_underwriting_enabled: Boolean(payload.autoUnderwritingEnabled),
        real_time_claim_notification: Boolean(payload.realTimeClaimNotification),
        payment_terms: payload.paymentTerms ?? "",
        claim_settlement_days: payload.claimSettlementDays ?? 0,
      },
    },
  );

  if (!result.ok) {
    return NextResponse.json(
      { ok: false, message: result.message ?? "Unable to update insurer configuration." },
      { status: result.status || 502 },
    );
  }

  return NextResponse.json({ ok: true, data: { saved: true } });
}
