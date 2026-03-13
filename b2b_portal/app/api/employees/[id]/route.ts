/**
 * /api/employees/[id]  GET | PATCH | DELETE
 */
import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { parseMoneyDecimal, sdkErrorMessage } from "@lib/sdk/api-helpers";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import type { EmployeeView, InsuranceType, EmployeeGender, EmployeeStatus, Money } from "@lifeplus/insuretech-sdk";
import type { Employee as UiEmployee } from "@lib/types/b2b";

type RouteContext = { params: Promise<{ id: string }> };
const INS: Record<string, string> = {
  INSURANCE_TYPE_UNSPECIFIED: "Unspecified", INSURANCE_TYPE_LIFE: "Life",
  INSURANCE_TYPE_HEALTH: "Health", INSURANCE_TYPE_AUTO: "Auto", INSURANCE_TYPE_TRAVEL: "Travel",
};
function fmt(d: number) { return d > 0 ? `BDT ${Math.round(d).toLocaleString("en-US", { maximumFractionDigits: 0 })}` : "—"; }
function toStatus(v?: EmployeeStatus): "Active" | "Inactive" { return v === "EMPLOYEE_STATUS_INACTIVE" ? "Inactive" : "Active"; }
/** Minimal UI shape returned for list views */
function mapView(v: EmployeeView): UiEmployee {
  const e = v.employee;
  return {
    id: e?.employee_uuid ?? "", name: e?.name ?? "", employeeID: e?.employee_id ?? "",
    department: v.department_name ?? "Unassigned",
    insuranceCategory: INS[e?.insurance_category ?? ""] ?? "Unspecified",
    assignedPlan: v.assigned_plan_name ?? e?.assigned_plan_id ?? "N/A",
    coverage: fmt(parseMoneyDecimal(e?.coverage_amount)),
    premiumAmount: fmt(parseMoneyDecimal(e?.premium_amount)),
    status: toStatus(e?.status), numberOfDependent: e?.number_of_dependent ?? 0,
  };
}

/**
 * Full employee record shape returned by GET /api/employees/[id]
 * Includes ALL fields needed to fully populate the edit form —
 * email, mobileNumber, gender, dateOfBirth, dateOfJoining, departmentId, etc.
 */
function mapViewFull(v: EmployeeView) {
  const e = v.employee;
  const covDecimal = parseMoneyDecimal(e?.coverage_amount);
  const premDecimal = parseMoneyDecimal(e?.premium_amount);
  return {
    // List-view fields
    id: e?.employee_uuid ?? "",
    name: e?.name ?? "",
    employeeID: e?.employee_id ?? "",
    department: v.department_name ?? "Unassigned",
    insuranceCategory: e?.insurance_category ? (
      // Map proto enum string → numeric form value
      e.insurance_category === "INSURANCE_TYPE_HEALTH" ? 1 :
        e.insurance_category === "INSURANCE_TYPE_LIFE" ? 2 :
          e.insurance_category === "INSURANCE_TYPE_AUTO" ? 3 :
            e.insurance_category === "INSURANCE_TYPE_TRAVEL" ? 4 : 0
    ) : 0,
    assignedPlan: v.assigned_plan_name ?? e?.assigned_plan_id ?? "N/A",
    coverage: fmt(covDecimal),
    premiumAmount: fmt(premDecimal),
    status: toStatus(e?.status),
    numberOfDependent: e?.number_of_dependent ?? 0,
    // Full form fields (missing from list view)
    email: e?.email ?? "",
    mobileNumber: e?.mobile_number ?? "",
    gender: e?.gender ?? "",
    dateOfBirth: e?.date_of_birth ?? "",
    dateOfJoining: e?.date_of_joining ?? "",
    departmentId: e?.department_id ?? "",
    businessId: e?.business_id ?? "",
    assignedPlanId: e?.assigned_plan_id ?? "",
    coverageAmount: covDecimal > 0 ? String(covDecimal) : "",
  };
}

export async function GET(request: Request, { params }: RouteContext) {
  const { id } = await params;
  if (!id) return NextResponse.json({ ok: false, message: "employee_uuid required" }, { status: 400 });
  try {
    const hdrs = await resolvePortalHeaders(request);
    const result = await makeSdkClient(request, hdrs ?? undefined).getEmployee({ path: { employee_uuid: id } });
    if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
    return NextResponse.json({ ok: true, employee: result.data?.employee ? mapViewFull(result.data.employee) : null });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}

export async function PATCH(request: Request, { params }: RouteContext) {
  const { id } = await params;
  if (!id) return NextResponse.json({ ok: false, message: "employee_uuid required" }, { status: 400 });
  const empId = id as string;
  try {
    const hdrs = await resolvePortalHeaders(request);
    const sdk = makeSdkClient(request, hdrs ?? undefined);
    const body = (await request.json()) as Record<string, unknown>;
    const cov = body.coverageAmount != null
      ? (typeof body.coverageAmount === "number" ? body.coverageAmount : Number.parseFloat(String(body.coverageAmount)))
      : undefined;
    const result = await sdk.updateEmployee({
      path: { employee_uuid: empId },
      body: {
        employee_uuid: empId,
        name: String(body.name ?? ""),
        department_id: String(body.departmentId ?? ""),
        email: String(body.email ?? ""),
        assigned_plan_id: String(body.assignedPlanId ?? ""),
        mobile_number: body.mobileNumber ? String(body.mobileNumber) : undefined,
        date_of_birth: body.dateOfBirth ? String(body.dateOfBirth) : undefined,
        date_of_joining: body.dateOfJoining ? String(body.dateOfJoining) : undefined,
        gender: body.gender as EmployeeGender | undefined,
        insurance_category: body.insuranceCategory as InsuranceType | undefined,
        coverage_amount: cov !== undefined && !Number.isNaN(cov)
          ? { amount: Math.round(cov * 100), currency: "BDT", decimal_amount: cov } as unknown as Money
          : undefined,
        number_of_dependent: body.numberOfDependent != null ? Number(body.numberOfDependent) : undefined,
        status: body.status as EmployeeStatus | undefined,
      },
    });
    if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
    return NextResponse.json({ ok: true, message: result.data?.message ?? "Updated", employee: result.data?.employee ? mapView(result.data.employee) : null });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}

export async function DELETE(request: Request, { params }: RouteContext) {
  const { id } = await params;
  if (!id) return NextResponse.json({ ok: false, message: "employee_uuid required" }, { status: 400 });
  try {
    const hdrs = await resolvePortalHeaders(request);
    const result = await makeSdkClient(request, hdrs ?? undefined).deleteEmployee({ path: { employee_uuid: id } });
    if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
    return NextResponse.json({ ok: true, message: "Employee deleted" });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}
