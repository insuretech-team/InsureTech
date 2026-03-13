/**
 * /api/departments/[id]  GET | PATCH | DELETE
 */
import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { parseMoneyDecimal, sdkErrorMessage } from "@lib/sdk/api-helpers";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import type { Department } from "@lifeplus/insuretech-sdk";
import type { Department as UiDepartment } from "@lib/types/b2b";

type RouteContext = { params: Promise<{ id: string }> };
function fmt(d: number) { return d > 0 ? `BDT ${Math.round(d).toLocaleString("en-US", { maximumFractionDigits: 0 })}` : "—"; }
function mapDept(d: Department): UiDepartment {
  return { id: d.department_id ?? "", name: d.name ?? "", employeeNo: d.employee_no ?? 0, totalPremium: fmt(parseMoneyDecimal(d.total_premium)) };
}

export async function GET(request: Request, { params }: RouteContext) {
  const { id } = await params;
  if (!id) return NextResponse.json({ ok: false, message: "department_id required" }, { status: 400 });
  try {
    const hdrs = await resolvePortalHeaders(request);
    const result = await makeSdkClient(request, hdrs ?? undefined).getDepartment({ path: { department_id: id } });
    if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
    return NextResponse.json({ ok: true, department: result.data?.department ? mapDept(result.data.department) : null });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}

export async function PATCH(request: Request, { params }: RouteContext) {
  const { id } = await params;
  if (!id) return NextResponse.json({ ok: false, message: "department_id required" }, { status: 400 });
  try {
    const hdrs = await resolvePortalHeaders(request);
    const sdk = makeSdkClient(request, hdrs ?? undefined);
    const body = (await request.json()) as Record<string, unknown>;
    const result = await sdk.updateDepartment({ path: { department_id: id }, body: { department_id: id, name: String(body.name ?? "") } });
    if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
    return NextResponse.json({ ok: true, message: result.data?.message ?? "Updated", department: result.data?.department ? mapDept(result.data.department) : null });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}

export async function DELETE(request: Request, { params }: RouteContext) {
  const { id } = await params;
  if (!id) return NextResponse.json({ ok: false, message: "department_id required" }, { status: 400 });
  try {
    const hdrs = await resolvePortalHeaders(request);
    const result = await makeSdkClient(request, hdrs ?? undefined).deleteDepartment({ path: { department_id: id } });
    if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
    return NextResponse.json({ ok: true, message: "Department deleted" });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}
