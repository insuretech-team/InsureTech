/**
 * GET /api/insurance-plans
 * Lists available insurance plans (products) from the gateway catalog.
 * Reuses the purchase-orders catalog endpoint which already returns plan/product data.
 */
import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";

import { resolvePortalHeaders } from "@lib/sdk/session-headers";

export async function GET(request: Request) {
  try {
    const hdrs = await resolvePortalHeaders(request);
    const url = new URL(request.url);
    const category = url.searchParams.get("category") ?? undefined;
    const pageSize = Number(url.searchParams.get("page_size") ?? 50);

    const sdk = makeSdkClient(request, hdrs ?? undefined);

    // Use the PO catalog as the source of insurance plans — it returns plan/product data
    // including planName, insuranceCategory, premiumAmount, productName.
    let items: unknown[] = [];
    const result = await sdk.listPurchaseOrderCatalog({ query: {} });
    if (result.response.ok) {
      const raw = result.data as Record<string, unknown> | null;
      const payload = (raw?.data && typeof raw.data === "object" ? raw.data : raw) as Record<string, unknown> | null;
      // Filter by category if provided
      const allItems = (payload?.items ?? payload?.products ?? []) as Record<string, unknown>[];
      items = category
        ? allItems.filter((i) => ((i.insuranceCategory as string) ?? "").toLowerCase() === category.toLowerCase())
        : allItems;
    }

    return NextResponse.json({ ok: true, items });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error", items: [] }, { status: 502 });
  }
}
