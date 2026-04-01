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

export type DashboardStatsResponse = {
  ok: boolean;
  stats?: DashboardStats;
  role?: string;
  message?: string;
  needsOrganisation?: boolean;
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

    // SDK interceptor unwraps the envelope, result.data is the inner payload directly.
    function unwrap(r: typeof orgs): Record<string, unknown> {
      if (!r?.data) return {};
      return r.data as Record<string, unknown>;
    }

    const orgsData  = unwrap(orgs);
    const empsData  = unwrap(emps);
    const posData   = unwrap(pos);
    const deptsData = unwrap(depts);

    const orgList   = (orgsData?.organisations ?? []) as Record<string, unknown>[];
    const pendingOrgs = orgList.filter(
      (o) => ((o.status as string) ?? "").includes("PENDING")
    ).length;

    const stats: DashboardStats = {
      totalOrganisations:   orgList.length,
      pendingOrganisations: pendingOrgs,
      totalEmployees:       (empsData?.total_count as number) ?? 0,
      totalDepartments:     (deptsData?.total_count as number) ?? 0,
      activePurchaseOrders: ((posData?.purchase_orders ?? []) as Record<string, unknown>[]).filter(
        (p) => {
          const po = p.purchase_order as Record<string, unknown> | undefined;
          const s = (po?.status as string) ?? "";
          return s.includes("ACTIVE") || !s.includes("CANCELLED");
        }
      ).length,
    };

    return NextResponse.json({ ok: true, stats, role: "SYSTEM_ADMIN" });

  } else {
    // B2B Admin / HR Manager — scoped to their org via x-business-id header.
    // If portal_biz_id cookie is missing (race on first load after login),
    // fall back to resolving the org from the session via /organisations/me.
    let orgId = hdrs.businessId;
    if (!orgId) {
      try {
        const meResult = await makeSdkClient(request, hdrs).getMyOrganisation();
        if (meResult.ok && typeof meResult.data.organisation_id === "string") {
          orgId = meResult.data.organisation_id;
        }
      } catch { /* ignore — will 400 below */ }
    }
    if (!orgId) {
      const emptyStats: DashboardStats = {
        totalEmployees: 0,
        totalDepartments: 0,
        activePurchaseOrders: 0,
        totalMembers: 0,
      };
      return NextResponse.json({
        ok: true,
        stats: emptyStats,
        role: "B2B_ORG_ADMIN",
        needsOrganisation: true,
        message: "Your account is active, but no organisation is linked yet.",
      } satisfies DashboardStatsResponse);
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

    // Unwrap ApiResponse<T> envelope
    function unwrapB2B(r: { data?: unknown } | null): Record<string, unknown> {
      if (!r?.data) return {};
      return r.data as Record<string, unknown>;
    }
    const empsData    = unwrapB2B(emps);
    const deptsData   = unwrapB2B(depts);
    const posData     = unwrapB2B(pos);
    const membersData = unwrapB2B(members);

    const stats: DashboardStats = {
      totalEmployees:       (empsData?.total_count as number) ?? 0,
      totalDepartments:     (deptsData?.total_count as number) ?? 0,
      activePurchaseOrders: ((posData?.purchase_orders ?? []) as Record<string, unknown>[]).filter(
        (p) => {
          const po = p.purchase_order as Record<string, unknown> | undefined;
          const s = (po?.status as string) ?? "";
          return s.includes("ACTIVE") || !s.includes("CANCELLED");
        }
      ).length,
      totalMembers:         ((membersData?.members ?? []) as unknown[]).length,
    };

    return NextResponse.json({ ok: true, stats, role: "B2B_ORG_ADMIN" } satisfies DashboardStatsResponse);
  }
}
