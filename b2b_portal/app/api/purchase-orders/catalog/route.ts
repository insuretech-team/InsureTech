/**
 * /api/purchase-orders/catalog  GET
 */
import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { parseMoneyDecimal, sdkErrorMessage } from "@lib/sdk/api-helpers";

const INS: Record<string, string> = {
  INSURANCE_TYPE_UNSPECIFIED: "Unspecified", INSURANCE_TYPE_LIFE: "Life",
  INSURANCE_TYPE_HEALTH: "Health", INSURANCE_TYPE_AUTO: "Auto", INSURANCE_TYPE_TRAVEL: "Travel",
};

function fmt(d: number) { return d > 0 ? `BDT ${Math.round(d).toLocaleString("en-US", { maximumFractionDigits: 0 })}` : "—"; }

export async function GET(request: Request) {
  try {
    // resolvePortalHeaders is required — catalog endpoint needs auth context
    const hdrs = await resolvePortalHeaders(request);
    const result = await makeSdkClient(request, hdrs ?? undefined).listPurchaseOrderCatalog();
    if (!result.response.ok) {
      return NextResponse.json({ ok: false, items: [], message: sdkErrorMessage(result) }, { status: result.response.status });
    }
    const items = (result.data?.items ?? []).map((item) => ({
      planId:            item.plan_id      ?? "",
      productId:         item.product_id   ?? "",
      productName:       item.product_name ?? "",
      planName:          item.plan_name    ?? "",
      insuranceCategory: INS[item.insurance_category ?? ""] ?? "Unspecified",
      premiumAmount:     fmt(parseMoneyDecimal(item.premium_amount)),
    }));
    return NextResponse.json({ ok: true, items });
  } catch (err) {
    return NextResponse.json({ ok: false, items: [], message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}
