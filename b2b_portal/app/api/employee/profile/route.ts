import { NextResponse } from "next/server";

import { makeDirectHttp } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";

export async function GET(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) {
    return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  }

  const result = await makeDirectHttp(request, hdrs).get("/v1/b2b-self/profile");
  if (!result.ok) {
    return NextResponse.json(
      { ok: false, message: result.message ?? "Unable to load employee profile" },
      { status: result.status || 500 }
    );
  }

  const payload = result.data as Record<string, unknown>;
  const employeeView = (payload.employee ?? payload) as Record<string, unknown>;
  const employee = (employeeView.employee ?? employeeView) as Record<string, unknown>;

  return NextResponse.json({
    ok: true,
    profile: {
      ...employee,
      department_name:
        typeof employeeView.department_name === "string" ? employeeView.department_name : "",
      assigned_plan_name:
        typeof employeeView.assigned_plan_name === "string" ? employeeView.assigned_plan_name : "",
    },
  }, { status: 200 });
}
