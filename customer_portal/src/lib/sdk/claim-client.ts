/**
 * claim-client.ts
 * ───────────────
 * Browser-side client for /api/claims/* BFF routes (customer portal).
 */
import { parseJson } from "./shared";
import type { ApiResult } from "./shared";

export type ClaimListResult = ApiResult<{ data?: Record<string, unknown>[]; total?: number }>;
export type ClaimSingleResult = ApiResult<{ data?: Record<string, unknown> }>;
export type ClaimCreatePayload = Record<string, unknown>;

export const claimClient = {
  async list(params?: Record<string, string | number>): Promise<ClaimListResult> {
    const qs = params ? "?" + new URLSearchParams(params as Record<string, string>).toString() : "";
    const response = await fetch(`/api/claims${qs}`, { method: "GET", cache: "no-store" });
    return parseJson<ClaimListResult>(response);
  },

  async get(id: string): Promise<ClaimSingleResult> {
    const response = await fetch(`/api/claims/${id}`, { method: "GET", cache: "no-store" });
    return parseJson<ClaimSingleResult>(response);
  },

  async create(payload: ClaimCreatePayload): Promise<ClaimSingleResult> {
    const response = await fetch("/api/claims", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    return parseJson<ClaimSingleResult>(response);
  },
};
