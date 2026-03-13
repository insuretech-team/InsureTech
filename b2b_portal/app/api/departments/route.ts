/**
 * /api/departments  GET | POST
 */
import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { parseMoneyDecimal, sdkErrorMessage } from "@lib/sdk/api-helpers";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import type { Department } from "@lifeplus/insuretech-sdk";
import type { Department as UiDepartment } from "@lib/types/b2b";

function fmt(d: number) { return d > 0 ? `BDT ${Math.round(d).toLocaleString("en-US", { maximumFractionDigits: 0 })}` : "—"; }
function mapDept(d: Department): UiDepartment {
  return { id: d.department_id ?? "", name: d.name ?? "", employeeNo: d.employee_no ?? 0, totalPremium: fmt(parseMoneyDecimal(d.total_premium)) };
}

export async function GET(request: Request) {
  try {
    const hdrs = await resolvePortalHeaders(request);
    const sdk = makeSdkClient(request, hdrs ?? undefined);
    const url = new URL(request.url);

    // Resolve business_id: explicit query param (super_admin) or from session (b2b_admin).
    // For super_admin (PORTAL_SYSTEM) business_id is not needed — backend uses system:root domain.
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
          // proceed without — backend enforces via session cookie
        }
      }
    }

    const result = await sdk.listDepartments({
      query: {
        page_size: Number(url.searchParams.get("page_size") ?? 50),
        business_id: businessId,
      },
    });
    if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result), departments: [] }, { status: result.response.status });
    return NextResponse.json({ ok: true, departments: (result.data?.departments ?? []).map(mapDept) });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error", departments: [] }, { status: 502 });
  }
}

export async function POST(request: Request) {
  try {
    const hdrs = await resolvePortalHeaders(request);
    const sdk = makeSdkClient(request, hdrs ?? undefined);
    const body = (await request.json()) as Record<string, unknown>;
    const businessId = String(body.businessId ?? "").trim();
    const result = await sdk.createDepartment({ body: { name: String(body.name ?? ""), business_id: businessId } });
    if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
    return NextResponse.json(
      { ok: true, message: result.data?.message ?? "Department created", department: result.data?.department ? mapDept(result.data.department) : null },
      { status: 201 }
    );
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}
