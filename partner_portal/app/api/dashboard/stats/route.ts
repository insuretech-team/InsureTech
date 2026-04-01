import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/partner-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage } from "@lib/sdk/api-helpers";

/** GET /api/dashboard/stats - Get partner-specific dashboard statistics */
export async function GET(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

  const { searchParams } = new URL(request.url);
  const partnerId = searchParams.get("partner_id") ?? hdrs.partnerId;
  const startDate = searchParams.get("start_date") ?? undefined;
  const endDate = searchParams.get("end_date") ?? undefined;

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.getDashboardStats({
    query: {
      partner_id: partnerId,
      start_date: startDate,
      end_date: endDate,
    },
  });

  if (!result.response.ok) {
    return NextResponse.json(
      { ok: false, message: sdkErrorMessage(result) },
      { status: result.response.status }
    );
  }

  const data = result.data as Record<string, unknown>;

  // Ensure all expected metrics are present
  const stats = {
    total_claims: data.total_claims ?? 0,
    approved_claims: data.approved_claims ?? 0,
    rejected_claims: data.rejected_claims ?? 0,
    pending_claims: data.pending_claims ?? 0,
    commission_earned: data.commission_earned ?? 0,
    commission_pending: data.commission_pending ?? 0,
    active_policies: data.active_policies ?? 0,
    total_agents: data.total_agents ?? 0,
    active_agents: data.active_agents ?? 0,
    ...data,
  };

  return NextResponse.json({ ok: true, data: stats }, { status: 200 });
}
