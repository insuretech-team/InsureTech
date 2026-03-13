/**
 * /api/employees  GET | POST
 * SDK: b2bServiceListEmployees / b2bServiceCreateEmployee
 */
import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { parseMoneyDecimal, sdkErrorMessage } from "@lib/sdk/api-helpers";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import type { EmployeeView, InsuranceType, EmployeeGender, EmployeeStatus, Money } from "@lifeplus/insuretech-sdk";
import type { Employee as UiEmployee } from "@lib/types/b2b";

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

    if (!businessId) {
      const cookieHeader = request.headers.get("cookie") ?? "";
      if (cookieHeader) {
        try {
          const base =
            process.env.INSURETECH_API_BASE_URL ??
            process.env.NEXT_PUBLIC_INSURETECH_API_BASE_URL ??
            "http://localhost:8080";
          const meRes = await fetch(`${base}/v1/b2b/organisations/me`, {
            method: "GET",
            headers: { cookie: cookieHeader },
            cache: "no-store",
          });
          if (meRes.ok) {
            const meData = (await meRes.json()) as Record<string, unknown>;
            if (typeof meData.organisation_id === "string" && meData.organisation_id) {
              businessId = meData.organisation_id;
            }
          }
        } catch {
          // proceed without — backend will enforce based on session cookie anyway
        }
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
    return NextResponse.json({ ok: true, employees: (result.data?.employees ?? []).map(mapView) });
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
    const result = await sdk.createEmployee({
      body: {
        name: String(body.name ?? ""), employee_id: String(body.employeeId ?? ""),
        business_id: businessId, department_id: String(body.departmentId ?? ""),
        insurance_category: body.insuranceCategory as InsuranceType | undefined,
        assigned_plan_id: body.assignedPlanId ? String(body.assignedPlanId) : "",
        coverage_amount: safeCov > 0
          ? { amount: Math.round(safeCov * 100), currency: "BDT", decimal_amount: safeCov } as unknown as Money
          : undefined,
        number_of_dependent: Number(body.numberOfDependent ?? 0),
        email: String(body.email ?? ""),
        mobile_number: body.mobileNumber ? String(body.mobileNumber) : undefined,
        date_of_birth: body.dateOfBirth ? String(body.dateOfBirth) : undefined,
        date_of_joining: body.dateOfJoining ? String(body.dateOfJoining) : undefined,
        gender: body.gender as EmployeeGender | undefined,
      },
    });
    if (!result.response.ok) {
      return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
    }
    return NextResponse.json(
      { ok: true, message: result.data?.message ?? "Employee created", employee: result.data?.employee ? mapView(result.data.employee) : null },
      { status: 201 }
    );
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}
