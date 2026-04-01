import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/partner-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage } from "@lib/sdk/api-helpers";

/** GET /api/dashboard/activity - Get recent activity feed */
export async function GET(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

  const { searchParams } = new URL(request.url);
  const page = parseInt(searchParams.get("page") ?? "1", 10);
  const pageSize = parseInt(searchParams.get("page_size") ?? "20", 10);
  const partnerId = searchParams.get("partner_id") ?? hdrs.partnerId;

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.getDashboardActivity({
    query: {
      partner_id: partnerId,
      page,
      page_size: pageSize,
    },
  });

  if (!result.response.ok) {
    return NextResponse.json(
      { ok: false, message: sdkErrorMessage(result) },
      { status: result.response.status }
    );
  }

  return NextResponse.json({ ok: true, data: result.data }, { status: 200 });
}
