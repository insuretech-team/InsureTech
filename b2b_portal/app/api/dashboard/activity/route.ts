/**
 * GET /api/dashboard/activity
 *
 * Returns recent activity items for the dashboard feed.
 * Assembles activity from recent orgs, employees, departments, and purchase orders.
 * Parallel fetches with Promise.allSettled — any single failure won't blank the feed.
 */
import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";

type ActivityItem = {
  id: string;
  type: "org" | "employee" | "department" | "po" | "member";
  title: string;
  subtitle: string;
  createdAt: string;
};

function settled<T>(r: PromiseSettledResult<T>, fallback: T): T {
  return r.status === "fulfilled" ? r.value : fallback;
}

export async function GET(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

  const isSuperAdmin = hdrs.portal === "PORTAL_SYSTEM";
  const sdk = makeSdkClient(request, hdrs);
  const orgId = hdrs.businessId;

  const activities: ActivityItem[] = [];

  if (isSuperAdmin) {
    const [orgsRes, empsRes, deptsRes, posRes] = await Promise.allSettled([
      sdk.listOrganisations({ query: { page_size: 5 } }),
      sdk.listEmployees({ query: { page_size: 5 } }),
      sdk.listDepartments({ query: { page_size: 5 } }),
      sdk.listPurchaseOrders({ query: { page_size: 5 } }),
    ]);

    for (const org of settled(orgsRes, null)?.data?.organisations ?? []) {
      activities.push({
        id:        `org-${org.organisation_id}`,
        type:      "org",
        title:     `Organisation registered: ${org.name ?? "Unknown"}`,
        subtitle:  `Code: ${org.code ?? "—"} · Status: ${(org.status ?? "").replace("ORGANISATION_STATUS_", "")}`,
        createdAt: org.created_at ?? "",
      });
    }
    for (const emp of settled(empsRes, null)?.data?.employees ?? []) {
      const e = emp.employee;
      activities.push({
        id:        `emp-${e?.employee_uuid}`,
        type:      "employee",
        title:     `Employee added: ${e?.name ?? "Unknown"}`,
        subtitle:  `ID: ${e?.employee_id ?? "—"}`,
        createdAt: e?.created_at ?? "",
      });
    }
    for (const dept of settled(deptsRes, null)?.data?.departments ?? []) {
      activities.push({
        id:        `dept-${dept.department_id}`,
        type:      "department",
        title:     `Department created: ${dept.name ?? "Unknown"}`,
        subtitle:  `Employees: ${dept.employee_no ?? 0}`,
        createdAt: dept.created_at ?? "",
      });
    }
    for (const poView of settled(posRes, null)?.data?.purchase_orders ?? []) {
      const po = poView.purchase_order;
      activities.push({
        id:        `po-${po?.purchase_order_id ?? Math.random()}`,
        type:      "po",
        title:     `Purchase order: ${po?.purchase_order_number ?? "—"}`,
        subtitle:  `Plan: ${poView.plan_name ?? po?.plan_id ?? "—"} · Status: ${(po?.status ?? "").replace("PURCHASE_ORDER_STATUS_", "")}`,
        createdAt: po?.created_at ?? "",
      });
    }
  } else {
    if (!orgId) return NextResponse.json({ ok: false, message: "No organisation context" }, { status: 400 });

    const [empsRes, deptsRes, posRes, membersRes] = await Promise.allSettled([
      sdk.listEmployees({ query: { page_size: 5, business_id: orgId } }),
      sdk.listDepartments({ query: { page_size: 5, business_id: orgId } }),
      sdk.listPurchaseOrders({ query: { page_size: 5, business_id: orgId } }),
      sdk.listOrgMembers({ path: { organisation_id: orgId } }),
    ]);

    for (const emp of settled(empsRes, null)?.data?.employees ?? []) {
      const e = emp.employee;
      activities.push({
        id:        `emp-${e?.employee_uuid}`,
        type:      "employee",
        title:     `Employee added: ${e?.name ?? "Unknown"}`,
        subtitle:  `ID: ${e?.employee_id ?? "—"}`,
        createdAt: e?.created_at ?? "",
      });
    }
    for (const dept of settled(deptsRes, null)?.data?.departments ?? []) {
      activities.push({
        id:        `dept-${dept.department_id}`,
        type:      "department",
        title:     `Department created: ${dept.name ?? "Unknown"}`,
        subtitle:  `Employees: ${dept.employee_no ?? 0}`,
        createdAt: dept.created_at ?? "",
      });
    }
    for (const poView of settled(posRes, null)?.data?.purchase_orders ?? []) {
      const po = poView.purchase_order;
      activities.push({
        id:        `po-${po?.purchase_order_id ?? Math.random()}`,
        type:      "po",
        title:     `Purchase order: ${po?.purchase_order_number ?? "—"}`,
        subtitle:  `Plan: ${poView.plan_name ?? po?.plan_id ?? "—"} · Status: ${(po?.status ?? "").replace("PURCHASE_ORDER_STATUS_", "")}`,
        createdAt: po?.created_at ?? "",
      });
    }
    for (const m of settled(membersRes, null)?.data?.members ?? []) {
      activities.push({
        id:        `member-${m.member_id}`,
        type:      "member",
        title:     `Member joined organisation`,
        subtitle:  `Role: ${(m.role ?? "").replace("ORG_MEMBER_ROLE_", "")}`,
        createdAt: m.joined_at ?? "",
      });
    }
  }

  // Sort by most recent first, limit to 10
  const sorted = activities
    .filter((a) => a.createdAt)
    .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
    .slice(0, 10);

  return NextResponse.json({ ok: true, activities: sorted });
}
