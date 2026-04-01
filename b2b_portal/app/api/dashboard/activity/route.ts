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

type ActivityResponse = {
  ok: boolean;
  activities?: ActivityItem[];
  message?: string;
  needsOrganisation?: boolean;
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

    function unwrapAct(r: Awaited<ReturnType<typeof sdk.listOrganisations>> | null): Record<string, unknown> {
      if (!r?.data) return {};
      const d = r.data as Record<string, unknown>;
      return (d.data && typeof d.data === "object" ? d.data : d) as Record<string, unknown>;
    }
    const orgsData  = unwrapAct(settled(orgsRes, null));
    const empsData  = unwrapAct(settled(empsRes, null) as Awaited<ReturnType<typeof sdk.listOrganisations>> | null);
    const deptsData = unwrapAct(settled(deptsRes, null) as Awaited<ReturnType<typeof sdk.listOrganisations>> | null);
    const posData   = unwrapAct(settled(posRes, null) as Awaited<ReturnType<typeof sdk.listOrganisations>> | null);

    for (const org of (orgsData?.organisations ?? []) as Record<string, unknown>[]) {
      activities.push({
        id:        `org-${org.organisation_id as string}`,
        type:      "org",
        title:     `Organisation registered: ${(org.name as string) ?? "Unknown"}`,
        subtitle:  `Code: ${(org.code as string) ?? "—"} · Status: ${((org.status as string) ?? "").replace("ORGANISATION_STATUS_", "")}`,
        createdAt: (org.created_at as string) ?? "",
      });
    }
    for (const emp of (empsData?.employees ?? []) as Record<string, unknown>[]) {
      const e = emp.employee as Record<string, unknown> | undefined;
      activities.push({
        id:        `emp-${e?.employee_uuid}`,
        type:      "employee",
        title:     `Employee added: ${(e?.name as string) ?? "Unknown"}`,
        subtitle:  `ID: ${(e?.employee_id as string) ?? "—"}`,
        createdAt: (e?.created_at as string) ?? "",
      });
    }
    for (const dept of (deptsData?.departments ?? []) as Record<string, unknown>[]) {
      activities.push({
        id:        `dept-${dept.department_id}`,
        type:      "department",
        title:     `Department created: ${(dept.name as string) ?? "Unknown"}`,
        subtitle:  `Employees: ${(dept.employee_no as number) ?? 0}`,
        createdAt: (dept.created_at as string) ?? "",
      });
    }
    for (const poView of (posData?.purchase_orders ?? []) as Record<string, unknown>[]) {
      const po = poView.purchase_order as Record<string, unknown> | undefined;
      activities.push({
        id:        `po-${po?.purchase_order_id ?? Math.random()}`,
        type:      "po",
        title:     `Purchase order: ${(po?.purchase_order_number as string) ?? "—"}`,
        subtitle:  `Plan: ${(poView.plan_name as string) ?? (po?.plan_id as string) ?? "—"} · Status: ${((po?.status as string) ?? "").replace("PURCHASE_ORDER_STATUS_", "")}`,
        createdAt: (po?.created_at as string) ?? "",
      });
    }
  } else {
    // Fall back to resolving org from session if portal_biz_id cookie is missing
    // (race condition on first dashboard load after login).
    let resolvedOrgId = orgId;
    if (!resolvedOrgId) {
      try {
        const meResult = await sdk.getMyOrganisation();
        if (meResult.ok && typeof meResult.data.organisation_id === "string") {
          resolvedOrgId = meResult.data.organisation_id;
        }
      } catch { /* ignore */ }
    }
    if (!resolvedOrgId) {
      return NextResponse.json({
        ok: true,
        activities: [],
        needsOrganisation: true,
        message: "Your account is active, but no organisation activity is available yet.",
      } satisfies ActivityResponse);
    }

    const [empsRes, deptsRes, posRes, membersRes] = await Promise.allSettled([
      sdk.listEmployees({ query: { page_size: 5, business_id: resolvedOrgId } }),
      sdk.listDepartments({ query: { page_size: 5, business_id: resolvedOrgId } }),
      sdk.listPurchaseOrders({ query: { page_size: 5, business_id: resolvedOrgId } }),
      sdk.listOrgMembers({ path: { organisation_id: resolvedOrgId } }),
    ]);

    function unwrapActB2B(r: Awaited<ReturnType<typeof sdk.listEmployees>> | null): Record<string, unknown> {
      if (!r?.data) return {};
      const d = r.data as Record<string, unknown>;
      return (d.data && typeof d.data === "object" ? d.data : d) as Record<string, unknown>;
    }
    const b2bEmps    = unwrapActB2B(settled(empsRes, null));
    const b2bDepts   = unwrapActB2B(settled(deptsRes, null) as Awaited<ReturnType<typeof sdk.listEmployees>> | null);
    const b2bPos     = unwrapActB2B(settled(posRes, null) as Awaited<ReturnType<typeof sdk.listEmployees>> | null);
    const b2bMembers = unwrapActB2B(settled(membersRes, null) as Awaited<ReturnType<typeof sdk.listEmployees>> | null);

    for (const emp of (b2bEmps?.employees ?? []) as Record<string, unknown>[]) {
      const e = emp.employee as Record<string, unknown> | undefined;
      activities.push({
        id:        `emp-${e?.employee_uuid}`,
        type:      "employee",
        title:     `Employee added: ${(e?.name as string) ?? "Unknown"}`,
        subtitle:  `ID: ${(e?.employee_id as string) ?? "—"}`,
        createdAt: (e?.created_at as string) ?? "",
      });
    }
    for (const dept of (b2bDepts?.departments ?? []) as Record<string, unknown>[]) {
      activities.push({
        id:        `dept-${dept.department_id}`,
        type:      "department",
        title:     `Department created: ${(dept.name as string) ?? "Unknown"}`,
        subtitle:  `Employees: ${(dept.employee_no as number) ?? 0}`,
        createdAt: (dept.created_at as string) ?? "",
      });
    }
    for (const poView of (b2bPos?.purchase_orders ?? []) as Record<string, unknown>[]) {
      const po = poView.purchase_order as Record<string, unknown> | undefined;
      activities.push({
        id:        `po-${po?.purchase_order_id ?? Math.random()}`,
        type:      "po",
        title:     `Purchase order: ${(po?.purchase_order_number as string) ?? "—"}`,
        subtitle:  `Plan: ${(poView.plan_name as string) ?? (po?.plan_id as string) ?? "—"} · Status: ${((po?.status as string) ?? "").replace("PURCHASE_ORDER_STATUS_", "")}`,
        createdAt: (po?.created_at as string) ?? "",
      });
    }
    for (const m of (b2bMembers?.members ?? []) as Record<string, unknown>[]) {
      activities.push({
        id:        `member-${m.member_id as string}`,
        type:      "member",
        title:     `Member joined organisation`,
        subtitle:  `Role: ${((m.role as string) ?? "").replace("ORG_MEMBER_ROLE_", "")}`,
        createdAt: (m.joined_at as string) ?? "",
      });
    }
  }

  // Sort by most recent first, limit to 10
  const sorted = activities
    .filter((a) => a.createdAt)
    .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
    .slice(0, 10);

  return NextResponse.json({ ok: true, activities: sorted } satisfies ActivityResponse);
}
