/**
 * department-client.ts
 * ─────────────────────
 * Browser-side client for /api/departments.
 */
import { parseJson, type ApiResult } from "./api-client";
import type { Department } from "@lib/types/b2b";

export type DepartmentListResult = ApiResult<{ departments: Department[]; total?: number }>;
export type DepartmentSingleResult = ApiResult<{ department?: Record<string, unknown> }>;

export const departmentClient = {
  async list(pageSize = 50, offset = 0, businessId?: string): Promise<DepartmentListResult> {
    const params = new URLSearchParams({ page_size: String(pageSize), offset: String(offset) });
    if (businessId) params.set("business_id", businessId);
    const res = await fetch(`/api/departments?${params}`, { method: "GET", cache: "no-store" });
    return parseJson<DepartmentListResult>(res);
  },

  async create(payload: { name: string; businessId: string }): Promise<DepartmentSingleResult> {
    const res = await fetch("/api/departments", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    return parseJson<DepartmentSingleResult>(res);
  },

  async update(id: string, payload: { name: string }): Promise<DepartmentSingleResult> {
    const res = await fetch(`/api/departments/${id}`, {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    return parseJson<DepartmentSingleResult>(res);
  },

  async delete(id: string): Promise<ApiResult> {
    const res = await fetch(`/api/departments/${id}`, { method: "DELETE" });
    return parseJson<ApiResult>(res);
  },
};
