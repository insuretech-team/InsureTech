/**
 * /api/purchase-orders/[id]  GET | PATCH | DELETE
 * SDK: getPurchaseOrder. Update/Delete → direct HTTP (not in SDK).
 */
import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { parseMoneyDecimal, sdkErrorMessage } from "@lib/sdk/api-helpers";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import type { PurchaseOrder } from "@lifeplus/insuretech-sdk";
import type { PurchaseOrder as UiPO } from "@lib/types/b2b";

type RouteContext = { params: Promise<{ id: string }> };
const STATUS: Record<string, string> = {
  PURCHASE_ORDER_STATUS_UNSPECIFIED: "Pending", PURCHASE_ORDER_STATUS_DRAFT: "In Draft",
  PURCHASE_ORDER_STATUS_SUBMITTED: "Submitted", PURCHASE_ORDER_STATUS_APPROVED: "Approved",
  PURCHASE_ORDER_STATUS_FULFILLED: "Fulfilled", PURCHASE_ORDER_STATUS_REJECTED: "Rejected",
};
function fmt(d: number) { return d > 0 ? `BDT ${Math.round(d).toLocaleString("en-US", { maximumFractionDigits: 0 })}` : "—"; }
type POWithMeta = PurchaseOrder & { department_name?: string; product_name?: string; plan_name?: string };
function mapPO(po: POWithMeta): UiPO {
  return {
    id: po.purchase_order_id ?? "", purchaseOrderNumber: po.purchase_order_number ?? "",
    productName: po.product_name ?? "Unknown Product", planName: po.plan_name ?? "Unknown Plan",
    insuranceCategory: po.insurance_category ?? "Unspecified",
    department: po.department_name ?? "Unassigned",
    employeeCount: po.employee_count ?? 0, numberOfDependents: po.number_of_dependents ?? 0,
    coverageAmount: fmt(parseMoneyDecimal(po.coverage_amount)),
    estimatedPremium: fmt(parseMoneyDecimal(po.estimated_premium)),
    status: STATUS[po.status ?? ""] ?? "Pending",
    submittedAt: po.created_at ? new Date(po.created_at).toLocaleDateString("en-GB") : "-",
    notes: po.notes ?? "",
  };
}

export async function GET(request: Request, { params }: RouteContext) {
  const { id } = await params;
  if (!id) return NextResponse.json({ ok: false, message: "purchase_order_id required" }, { status: 400 });
  try {
    const hdrs = await resolvePortalHeaders(request);
    const result = await makeSdkClient(request, hdrs ?? undefined).getPurchaseOrder({ path: { purchase_order_id: id } });
    if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
    const po = result.data?.purchase_order as POWithMeta | undefined;
    return NextResponse.json({ ok: true, purchaseOrder: po ? mapPO(po) : null });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}

export async function PATCH(request: Request, { params }: RouteContext) {
  const { id } = await params;
  if (!id) return NextResponse.json({ ok: false, message: "purchase_order_id required" }, { status: 400 });
  try {
    const hdrs = await resolvePortalHeaders(request);
    const body = (await request.json()) as Record<string, unknown>;
    const sdk = makeSdkClient(request, hdrs ?? undefined);
    const patch: Record<string, unknown> = {};
    if (body.status != null) patch.status = body.status;
    if (body.notes != null) patch.notes = String(body.notes);
    if (body.employeeCount != null) patch.employee_count = Number(body.employeeCount);
    if (body.numberOfDependents != null) patch.number_of_dependents = Number(body.numberOfDependents);
    const result = await sdk.updatePurchaseOrderHttp(id, patch);
    if (!result.ok) return NextResponse.json({ ok: false, message: (result.data?.message as string) ?? "Failed to update" }, { status: result.status });
    return NextResponse.json({ ok: true, message: (result.data?.message as string) ?? "Updated" });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}

export async function DELETE(request: Request, { params }: RouteContext) {
  const { id } = await params;
  if (!id) return NextResponse.json({ ok: false, message: "purchase_order_id required" }, { status: 400 });
  try {
    const hdrs = await resolvePortalHeaders(request);
    const result = await makeSdkClient(request, hdrs ?? undefined).deletePurchaseOrderHttp(id);
    if (!result.ok) return NextResponse.json({ ok: false, message: (result.data?.message as string) ?? "Failed to delete" }, { status: result.status });
    return NextResponse.json({ ok: true, message: (result.data?.message as string) ?? "Deleted" });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}
