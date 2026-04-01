/**
 * payment-client.ts
 * ─────────────────
 * Browser-side client for /api/payments/* BFF routes (customer portal).
 */
import { parseJson } from "./shared";
import type { ApiResult } from "./shared";

export type PaymentListResult = ApiResult<{ data?: Record<string, unknown>[]; total?: number }>;
export type PaymentSingleResult = ApiResult<{ data?: Record<string, unknown> }>;

export const paymentClient = {
  async list(params?: Record<string, string | number>): Promise<PaymentListResult> {
    const qs = params ? "?" + new URLSearchParams(params as Record<string, string>).toString() : "";
    const response = await fetch(`/api/payments${qs}`, { method: "GET", cache: "no-store" });
    return parseJson<PaymentListResult>(response);
  },

  async get(id: string): Promise<PaymentSingleResult> {
    const response = await fetch(`/api/payments/${id}`, { method: "GET", cache: "no-store" });
    return parseJson<PaymentSingleResult>(response);
  },
};
