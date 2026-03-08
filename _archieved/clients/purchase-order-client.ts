/**
 * purchase-order-client.ts
 * ─────────────────────────
 * Browser-side client for /api/purchase-orders.
 */
import { parseJson, type ApiResult } from "./api-client";
import type { PurchaseOrder } from "@lib/types/b2b";

export type CatalogItem = {
  planId: string;
  productId: string;
  productName: string;
  planName: string;
  insuranceCategory: string;
  premiumAmount: string;
};

export type PurchaseOrderCreatePayload = {
  departmentId:        string;
  planId:              string;
  insuranceCategory?:  string;   // human label e.g. "Health" — API route maps to proto enum
  employeeCount:       number;
  numberOfDependents?: number;
  coverageAmount?:     number;
  notes?:              string;
};

export type PurchaseOrderUpdatePayload = {
  status?: number;
  notes?: string;
  employeeCount?: number;
  numberOfDependents?: number;
  coverageAmount?: number;
};

export type POListResult = ApiResult<{ purchaseOrders: PurchaseOrder[]; total?: number }>;
export type POSingleResult = ApiResult<{ purchaseOrder?: PurchaseOrder | null }>;
export type POCatalogResult = ApiResult<{ items: CatalogItem[] }>;

export const purchaseOrderClient = {
  async list(options?: { pageSize?: number; offset?: number; status?: number }): Promise<POListResult> {
    const params = new URLSearchParams({ page_size: String(options?.pageSize ?? 50), offset: String(options?.offset ?? 0) });
    if (options?.status != null) params.set("status", String(options.status));
    const res = await fetch(`/api/purchase-orders?${params}`, { method: "GET", cache: "no-store" });
    return parseJson<POListResult>(res);
  },

  async get(id: string): Promise<POSingleResult> {
    const res = await fetch(`/api/purchase-orders/${id}`, { method: "GET", cache: "no-store" });
    return parseJson<POSingleResult>(res);
  },

  async create(payload: PurchaseOrderCreatePayload): Promise<POSingleResult> {
    const res = await fetch("/api/purchase-orders", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    return parseJson<POSingleResult>(res);
  },

  async update(id: string, payload: PurchaseOrderUpdatePayload): Promise<POSingleResult> {
    const res = await fetch(`/api/purchase-orders/${id}`, {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    return parseJson<POSingleResult>(res);
  },

  async delete(id: string): Promise<ApiResult> {
    const res = await fetch(`/api/purchase-orders/${id}`, { method: "DELETE" });
    return parseJson<ApiResult>(res);
  },

  async getCatalog(): Promise<POCatalogResult> {
    const res = await fetch("/api/purchase-orders/catalog", { method: "GET", cache: "no-store" });
    return parseJson<POCatalogResult>(res);
  },
};
