/**
 * GET /api/billing/payments
 * Lists payments for the current organisation via the gateway billing service.
 */
import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { sdkErrorMessage } from "@lib/sdk/api-helpers";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";

export async function GET(request: Request) {
  try {
    const hdrs = await resolvePortalHeaders(request);
    if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

    const url = new URL(request.url);
    const pageSize = Number(url.searchParams.get("page_size") ?? 20);
    const page = Number(url.searchParams.get("page") ?? 1);
    const status = url.searchParams.get("status") ?? undefined;

    const sdk = makeSdkClient(request, hdrs);
    const result = await sdk.listPayments({
      query: {
        page_size: pageSize,
        page,
        ...(status ? { status } : {}),
        ...(hdrs.businessId ? { business_id: hdrs.businessId } : {}),
      },
    });

    if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result), payments: [] }, { status: result.response.status });
    const payload = result.data as Record<string, unknown> | null;
    return NextResponse.json({ ok: true, payments: payload?.payments ?? [], total: payload?.total ?? 0 });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error", payments: [] }, { status: 502 });
  }
}
