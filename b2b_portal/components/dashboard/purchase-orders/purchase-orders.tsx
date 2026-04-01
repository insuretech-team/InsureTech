"use client";

import { useCallback, useEffect, useMemo, useState } from "react";

import DashboardLayout from "../dashboard-layout";

import PurchaseOrderCard from "./card";
import { buildPurchaseOrderColumns, type PurchaseOrder } from "./data-table/columns";
import { DataTable } from "./data-table/data-table";
import { bffClient, type CatalogItem, type PurchaseOrderCreatePayload } from "@lib/sdk/b2b-sdk-client";
import { usePortalPrincipal } from "@lib/auth/portal-session-context";

type DepartmentOption = { id: string; name: string };

const PurchaseOrders = () => {
  const principal = usePortalPrincipal();
  const [data, setData] = useState<PurchaseOrder[]>([]);
  const [departments, setDepartments] = useState<DepartmentOption[]>([]);
  const [catalog, setCatalog] = useState<CatalogItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  // Resolved org id — for B2B admin this is their org, for super admin it stays empty
  // (super admin can select dept across all orgs but PO creation is per-dept so dept dropdown suffices)
  const [resolvedOrgId, setResolvedOrgId] = useState("");

  // Step 1: resolve org context from the server-hydrated portal session
  useEffect(() => {
    let cancelled = false;
    const role = principal?.role ?? "";
    if (!role) {
      setResolvedOrgId("");
      return () => { cancelled = true; };
    }

    const isSuperAdmin = role === "SYSTEM_ADMIN";
    if (isSuperAdmin) {
      setResolvedOrgId("");
      return () => { cancelled = true; };
    }

    void bffClient.organisations.getMe().then((res) => {
      if (cancelled || !res.ok || !res.organisation?.id) return;
      setResolvedOrgId(res.organisation.id);
    });
    return () => { cancelled = true; };
  }, [principal?.role]);

  const loadAll = useCallback(async () => {
    setLoading(true);
    try {
      const [ordersResult, deptResult, catalogResult] = await Promise.all([
        bffClient.purchaseOrders.list(),
        bffClient.departments.list(200, 0, resolvedOrgId || undefined),
        bffClient.purchaseOrders.getCatalog(),
      ]);
      setData(ordersResult.ok && Array.isArray(ordersResult.purchaseOrders) ? ordersResult.purchaseOrders : []);
      setDepartments(deptResult.ok && Array.isArray(deptResult.departments) ? deptResult.departments : []);
      setCatalog(catalogResult.ok && Array.isArray(catalogResult.items) ? catalogResult.items : []);
    } finally {
      setLoading(false);
    }
  }, [resolvedOrgId]);

  useEffect(() => { void loadAll(); }, [loadAll]);

  const summary = useMemo(() => {
    const statuses = {
      total: data.length,
      draft: 0,
      submitted: 0,
      approved: 0,
      fulfilled: 0,
      rejected: 0,
    };

    for (const item of data) {
      switch (item.status) {
        case "In Draft":
          statuses.draft += 1;
          break;
        case "Submitted":
          statuses.submitted += 1;
          break;
        case "Approved":
          statuses.approved += 1;
          break;
        case "Fulfilled":
          statuses.fulfilled += 1;
          break;
        case "Rejected":
          statuses.rejected += 1;
          break;
        default:
          break;
      }
    }

    return [
      { id: 1, title: "Total Orders", value: statuses.total, icon: "./quotations/comment-quote.svg", bgColor: "var(--brand-surface-1)" },
      { id: 2, title: "In Draft", value: statuses.draft, icon: "./quotations/form.svg", bgColor: "var(--brand-surface-2)" },
      { id: 3, title: "Submitted", value: statuses.submitted, icon: "./quotations/paper-plane.svg", bgColor: "var(--brand-surface-5)" },
      { id: 4, title: "Approved", value: statuses.approved, icon: "./quotations/check-circle.svg", bgColor: "var(--brand-surface-3)" },
      { id: 5, title: "Fulfilled", value: statuses.fulfilled, icon: "./quotations/inbox-in.svg", bgColor: "var(--brand-surface-neutral)" },
      { id: 6, title: "Rejected", value: statuses.rejected, icon: "./quotations/cross-circle.svg", bgColor: "var(--brand-surface-4)" },
    ];
  }, [data]);

  async function handleCreatePurchaseOrder(payload: PurchaseOrderCreatePayload): Promise<boolean> {
    setSubmitting(true);
    try {
      const result = await bffClient.purchaseOrders.create(payload);
      if (!result.ok) {
        window.alert(result.message ?? "Failed to submit purchase order");
        return false;
      }
      await loadAll();
      return true;
    } catch (error) {
      window.alert(error instanceof Error ? error.message : "Failed to submit purchase order");
      return false;
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <DashboardLayout>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-6">
        {summary.map((item) => (
          <PurchaseOrderCard key={item.id} {...item} />
        ))}
      </div>

      <div className="space-y-4 py-4">
        <DataTable
          columns={buildPurchaseOrderColumns(loadAll)}
          data={data}
          loading={loading}
          departments={departments}
          catalog={catalog}
          submitting={submitting}
          onCreatePurchaseOrder={handleCreatePurchaseOrder}
        />
      </div>
    </DashboardLayout>
  );
};

export default PurchaseOrders;
