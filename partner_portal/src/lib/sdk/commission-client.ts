/**
 * commission-client.ts — Partner Portal
 * ──────────────────────────────────────
 * Browser-side client for /api/commissions/* BFF routes.
 */
import { parseJson } from "./shared";
import type { ApiResult } from "./shared";

export type CommissionListResult = ApiResult<{ data?: Record<string, unknown>[]; total?: number }>;
export type CommissionSingleResult = ApiResult<{ data?: Record<string, unknown> }>;

export const commissionClient = {
  async list(params?: Record<string, string | number>): Promise<CommissionListResult> {
    const qs = params ? "?" + new URLSearchParams(params as Record<string, string>).toString() : "";
    const response = await fetch(`/api/commissions${qs}`, { method: "GET", cache: "no-store" });
    return parseJson<CommissionListResult>(response);
  },

  async get(id: string): Promise<CommissionSingleResult> {
    const response = await fetch(`/api/commissions/${id}`, { method: "GET", cache: "no-store" });
    return parseJson<CommissionSingleResult>(response);
  },
};
