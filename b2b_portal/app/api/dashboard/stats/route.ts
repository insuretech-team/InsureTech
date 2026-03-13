/**
 * GET /api/dashboard/stats
 *
 * Returns KPI counts for the dashboard stats cards.
 * Super Admin:  total orgs, total employees across all orgs, pending orgs, active POs
 * B2B Admin:    own org member count, own org dept count, own employee count, active POs
 *
 * We parallelise all backend calls with Promise.allSettled so a single
 * failing RPC doesn't blank the whole dashboard.
 */
import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage } from "@lib/sdk/api-helpers";

export type DashboardStats = {
  // Super Admin cards
  totalOrganisations?: number;
  pendingOrganisations?: number;
  // Shared cards
  totalEmployees:    number;
  totalDepartments:  number;
  activePurchaseOrders: number;
  // B2B Admin cards
  totalMembers?:     number;
};

function settled<T>(result: PromiseSettledResult<T>, fallback: T): T {
  return result.status === "fulfilled" ? result.value : fallback;
}

export async function GET(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) {
    return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  }

  const isSuperAdmin = hdrs.portal === "PORTAL_SYSTEM";
  const sdk = makeSdkClient(request, hdrs);

  if (isSuperAdmin) {
    // Super Admin — fetch all-platform counts in parallel
    const [orgsRes, empsRes, posRes, deptsRes] = await Promise.allSettled([
      sdk.listOrganisations({ query: { page_size: 1000 } }),
      sdk.listEmployees({ query: { page_size: 1 } }),
      sdk.listPurchaseOrders({ query: { page_size: 1000 } }),
      sdk.listDepartments({ query: { page_size: 1 } }),
    ]);

    const orgs      = settled(orgsRes, null);
    const emps      = settled(empsRes, null);
    const pos       = settled(posRes, null);
    const depts     = settled(deptsRes, null);

    const orgList   = orgs?.data?.organisations ?? [];
    const pendingOrgs = orgList.filter(
      (o) => (o.status ?? "").includes("PENDING")
    ).length;

    const stats: DashboardStats = {
      totalOrganisations:   orgList.length,
      pendingOrganisations: pendingOrgs,
      totalEmployees:       emps?.data?.total_count ?? 0,
      totalDepartments:     depts?.data?.total_count ?? 0,
      activePurchaseOrders: (pos?.data?.purchase_orders ?? []).filter(
        (p) => {
          const s = p.purchase_order?.status ?? "";
          return s.includes("ACTIVE") || !s.includes("CANCELLED");
        }
      ).length,
    };

    return NextResponse.json({ ok: true, stats, role: "SYSTEM_ADMIN" });

  } else {
    // B2B Admin / HR Manager — scoped to their org via x-business-id header
    const orgId = hdrs.businessId;
    if (!orgId) {
      return NextResponse.json({ ok: false, message: "No organisation context" }, { status: 400 });
    }

    const [empsRes, deptsRes, posRes, membersRes] = await Promise.allSettled([
      sdk.listEmployees({ query: { page_size: 1, business_id: orgId } }),
      sdk.listDepartments({ query: { page_size: 1, business_id: orgId } }),
      sdk.listPurchaseOrders({ query: { page_size: 1000, business_id: orgId } }),
      sdk.listOrgMembers({ path: { organisation_id: orgId } }),
    ]);

    const emps    = settled(empsRes, null);
    const depts   = settled(deptsRes, null);
    const pos     = settled(posRes, null);
    const members = settled(membersRes, null);

    const stats: DashboardStats = {
      totalEmployees:       emps?.data?.total_count ?? 0,
      totalDepartments:     depts?.data?.total_count ?? 0,
      activePurchaseOrders: (pos?.data?.purchase_orders ?? []).filter(
        (p) => {
          const s = p.purchase_order?.status ?? "";
          return s.includes("ACTIVE") || !s.includes("CANCELLED");
        }
      ).length,
      totalMembers:         (members?.data?.members ?? []).length,
    };

    return NextResponse.json({ ok: true, stats, role: "B2B_ORG_ADMIN" });
  }
}
