/**
 * quotation-client.ts
 * ───────────────────
 * Browser-side client for /api/quotations/* BFF routes (customer portal).
 */
import { parseJson } from "./shared";
import type { ApiResult } from "./shared";

export type QuotationListResult = ApiResult<{ data?: Record<string, unknown>[]; total?: number }>;
export type QuotationSingleResult = ApiResult<{ data?: Record<string, unknown> }>;
export type QuotationCreatePayload = Record<string, unknown>;

export const quotationClient = {
  async list(params?: Record<string, string | number>): Promise<QuotationListResult> {
    const qs = params ? "?" + new URLSearchParams(params as Record<string, string>).toString() : "";
    const response = await fetch(`/api/quotations${qs}`, { method: "GET", cache: "no-store" });
    return parseJson<QuotationListResult>(response);
  },

  async get(id: string): Promise<QuotationSingleResult> {
    const response = await fetch(`/api/quotations/${id}`, { method: "GET", cache: "no-store" });
    return parseJson<QuotationSingleResult>(response);
  },

  async create(payload: QuotationCreatePayload): Promise<QuotationSingleResult> {
    const response = await fetch("/api/quotations", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    return parseJson<QuotationSingleResult>(response);
  },
};
