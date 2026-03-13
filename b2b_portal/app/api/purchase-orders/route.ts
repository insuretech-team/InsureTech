/**
 * /api/purchase-orders  GET | POST
 */
import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { parseMoneyDecimal, sdkErrorMessage } from "@lib/sdk/api-helpers";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import type { PurchaseOrder, Money } from "@lifeplus/insuretech-sdk";
import type { PurchaseOrder as UiPO } from "@lib/types/b2b";

const STATUS: Record<string, string> = {
  PURCHASE_ORDER_STATUS_UNSPECIFIED: "Pending", PURCHASE_ORDER_STATUS_DRAFT: "In Draft",
  PURCHASE_ORDER_STATUS_SUBMITTED: "Submitted", PURCHASE_ORDER_STATUS_APPROVED: "Approved",
  PURCHASE_ORDER_STATUS_FULFILLED: "Fulfilled", PURCHASE_ORDER_STATUS_REJECTED: "Rejected",
};
const INS: Record<string, string> = {
  INSURANCE_TYPE_UNSPECIFIED: "Unspecified", INSURANCE_TYPE_LIFE: "Life",
  INSURANCE_TYPE_HEALTH: "Health", INSURANCE_TYPE_AUTO: "Auto", INSURANCE_TYPE_TRAVEL: "Travel",
};

function fmt(d: number) { return d > 0 ? `BDT ${Math.round(d).toLocaleString("en-US", { maximumFractionDigits: 0 })}` : "—"; }
type POWithMeta = PurchaseOrder & { department_name?: string; product_name?: string; plan_name?: string };
function mapPO(po: POWithMeta): UiPO {
  return {
    id: po.purchase_order_id ?? "", purchaseOrderNumber: po.purchase_order_number ?? "",
    productName: po.product_name ?? "Unknown Product", planName: po.plan_name ?? "Unknown Plan",
    insuranceCategory: INS[po.insurance_category ?? ""] ?? "Unspecified",
    department: po.department_name ?? "Unassigned",
    employeeCount: po.employee_count ?? 0, numberOfDependents: po.number_of_dependents ?? 0,
    coverageAmount: fmt(parseMoneyDecimal(po.coverage_amount)),
    estimatedPremium: fmt(parseMoneyDecimal(po.estimated_premium)),
    status: STATUS[po.status ?? ""] ?? "Pending",
    submittedAt: po.created_at ? new Date(po.created_at).toLocaleDateString("en-GB") : "-",
    notes: po.notes ?? "",
  };
}

export async function GET(request: Request) {
  try {
    const hdrs = await resolvePortalHeaders(request);
    const sdk = makeSdkClient(request, hdrs ?? undefined);
    const url = new URL(request.url);
    const result = await sdk.listPurchaseOrders({ query: { page_size: Number(url.searchParams.get("page_size") ?? 50) } });
    if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result), purchaseOrders: [] }, { status: result.response.status });
    const pos = (result.data?.purchase_orders ?? []) as POWithMeta[];
    return NextResponse.json({ ok: true, purchaseOrders: pos.map(mapPO) });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error", purchaseOrders: [] }, { status: 502 });
  }
}

export async function POST(request: Request) {
  try {
    const hdrs = await resolvePortalHeaders(request);
    const sdk = makeSdkClient(request, hdrs ?? undefined);
    const body = (await request.json()) as Record<string, unknown>;
    const cov = typeof body.coverageAmount === "number" ? body.coverageAmount : Number.parseFloat(String(body.coverageAmount ?? "0"));
    const safeCov = Number.isNaN(cov) ? 0 : cov;
    // Note: insurance_category is NOT a field on PurchaseOrderCreationRequest proto —
    // the backend derives it automatically from the plan_id. We only forward the fields
    // the proto actually accepts.
    const result = await sdk.createPurchaseOrder({
      body: {
        department_id: String(body.departmentId ?? ""),
        plan_id: String(body.planId ?? ""),
        employee_count: Number(body.employeeCount ?? 0),
        number_of_dependents: Number(body.numberOfDependents ?? 0),
        coverage_amount: safeCov > 0 ? { amount: Math.round(safeCov * 100), currency: "BDT", decimal_amount: safeCov } as unknown as Money : undefined,
        notes: body.notes ? String(body.notes) : undefined,
      },
    });
    if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
    const poRaw = result.data?.purchase_order as POWithMeta | undefined;
    return NextResponse.json({ ok: true, message: result.data?.message ?? "Purchase order created", purchaseOrder: poRaw ? mapPO(poRaw) : null }, { status: 201 });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}
