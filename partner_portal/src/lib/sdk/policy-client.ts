/**
 * policy-client.ts — Partner Portal
 * ───────────────────────────────────
 * Browser-side client for /api/policies/* BFF routes.
 */
import { parseJson } from "./shared";
import type { ApiResult } from "./shared";

export type PolicyListResult = ApiResult<{ data?: Record<string, unknown>[]; total?: number }>;
export type PolicySingleResult = ApiResult<{ data?: Record<string, unknown> }>;

export const policyClient = {
  async list(params?: Record<string, string | number>): Promise<PolicyListResult> {
    const qs = params ? "?" + new URLSearchParams(params as Record<string, string>).toString() : "";
    const response = await fetch(`/api/policies${qs}`, { method: "GET", cache: "no-store" });
    return parseJson<PolicyListResult>(response);
  },

  async get(id: string): Promise<PolicySingleResult> {
    const response = await fetch(`/api/policies/${id}`, { method: "GET", cache: "no-store" });
    return parseJson<PolicySingleResult>(response);
  },
};
