/**
 * /api/employees  GET | POST
 * SDK: b2bServiceListEmployees / b2bServiceCreateEmployee
 */
import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { parseMoneyDecimal, sdkErrorMessage, unwrapSdkResult } from "@lib/sdk/api-helpers";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import type { EmployeeView, InsuranceType, EmployeeGender, EmployeeStatus, Money } from "@lifeplus/insuretech-sdk";
import type { Employee as UiEmployee } from "@lib/types/b2b";
import {
  getBangladeshMobileValidationMessage,
  normalizeBangladeshMobile,
} from "@/src/lib/utils/bd-mobile";

const INS: Record<string, string> = {
  INSURANCE_TYPE_UNSPECIFIED: "Unspecified", INSURANCE_TYPE_LIFE: "Life",
  INSURANCE_TYPE_HEALTH: "Health", INSURANCE_TYPE_AUTO: "Auto", INSURANCE_TYPE_TRAVEL: "Travel",
};
function fmt(d: number) { return d > 0 ? `BDT ${Math.round(d).toLocaleString("en-US", { maximumFractionDigits: 0 })}` : "—"; }
function toStatus(v?: EmployeeStatus): "Active" | "Inactive" { return v === "EMPLOYEE_STATUS_INACTIVE" ? "Inactive" : "Active"; }
function mapView(v: EmployeeView): UiEmployee {
  const e = v.employee;
  return {
    id: e?.employee_uuid ?? "", name: e?.name ?? "", employeeID: e?.employee_id ?? "",
    department: v.department_name ?? "Unassigned",
    insuranceCategory: INS[e?.insurance_category ?? ""] ?? "Unspecified",
    assignedPlan: v.assigned_plan_name ?? e?.assigned_plan_id ?? "N/A",
    coverage: fmt(parseMoneyDecimal(e?.coverage_amount)),
    premiumAmount: fmt(parseMoneyDecimal(e?.premium_amount)),
    status: toStatus(e?.status),
    numberOfDependent: e?.number_of_dependent ?? 0,
  };
}

export async function GET(request: Request) {
  try {
    const hdrs = await resolvePortalHeaders(request);
    const sdk = makeSdkClient(request, hdrs ?? undefined);
    const url = new URL(request.url);

    // business_id may be supplied explicitly (super_admin selects an org from the dropdown).
    // For b2b_admin it is never in the query — we resolve it from the session.
    // For super_admin (PORTAL_SYSTEM) business_id is not required by backend.
    let businessId = url.searchParams.get("business_id") ?? hdrs?.businessId ?? undefined;

    if (!businessId && hdrs?.portal !== "PORTAL_SYSTEM") {
      try {
        const meResult = await sdk.getMyOrganisation();
        if (meResult.ok && typeof meResult.data.organisation_id === "string" && meResult.data.organisation_id) {
          businessId = meResult.data.organisation_id;
        }
      } catch {
        // proceed without — backend will enforce based on session cookie anyway
      }
    }

    const result = await sdk.listEmployees({
      query: {
        page_size: Number(url.searchParams.get("page_size") ?? 50),
        business_id: businessId,
        department_id: url.searchParams.get("department_id") ?? undefined,
      },
    });
    if (!result.response.ok) {
      return NextResponse.json({ ok: false, message: sdkErrorMessage(result), employees: [] }, { status: result.response.status });
    }
    // SDK interceptor unwraps the envelope; result.data is the payload directly.
    const payload = result.data as Record<string, unknown> | null;
    return NextResponse.json({ ok: true, employees: ((payload?.employees ?? []) as unknown[]).map((v) => mapView(v as EmployeeView)) });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error", employees: [] }, { status: 502 });
  }
}

export async function POST(request: Request) {
  try {
    const hdrs = await resolvePortalHeaders(request);
    const sdk = makeSdkClient(request, hdrs ?? undefined);
    const body = (await request.json()) as Record<string, unknown>;
    const businessId = String(body.businessId ?? "").trim();
    const cov = typeof body.coverageAmount === "number" ? body.coverageAmount : Number.parseFloat(String(body.coverageAmount ?? "0"));
    const safeCov = Number.isNaN(cov) ? 0 : cov;
    const mobileNumber = body.mobileNumber ? String(body.mobileNumber).trim() : "";
    const normalizedMobileNumber = mobileNumber ? normalizeBangladeshMobile(mobileNumber) : null;
    if (mobileNumber && !normalizedMobileNumber) {
      return NextResponse.json(
        { ok: false, message: getBangladeshMobileValidationMessage("Employee mobile number") },
        { status: 400 }
      );
    }
    const result = await sdk.createEmployee({
      body: {
        user_id: hdrs?.userId ?? "",
        name: String(body.name ?? ""), employee_id: String(body.employeeId ?? ""),
        business_id: businessId, department_id: String(body.departmentId ?? ""),
        insurance_category: body.insuranceCategory as InsuranceType | undefined,
        assigned_plan_id: body.assignedPlanId ? String(body.assignedPlanId) : "",
        coverage_amount: safeCov > 0
          ? { amount: Math.round(safeCov * 100), currency: "BDT", decimal_amount: safeCov } as unknown as Money
          : undefined,
        number_of_dependent: Number(body.numberOfDependent ?? 0),
        email: String(body.email ?? ""),
        mobile_number: normalizedMobileNumber ?? undefined,
        date_of_birth: body.dateOfBirth ? String(body.dateOfBirth) : undefined,
        date_of_joining: body.dateOfJoining ? String(body.dateOfJoining) : undefined,
        gender: body.gender as EmployeeGender | undefined,
      },
    });
    const unwrapped = unwrapSdkResult(result);
    if (!unwrapped.ok) return NextResponse.json({ ok: false, message: unwrapped.message }, { status: unwrapped.status });
    const d = unwrapped.data as Record<string, unknown>;
    return NextResponse.json(
      { ok: true, message: (d?.message as string) ?? "Employee created", employee: d?.employee ? mapView(d.employee as EmployeeView) : null },
      { status: 201 }
    );
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}
