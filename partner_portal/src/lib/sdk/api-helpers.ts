/**
 * api-helpers.ts
 * Shared utilities for Next.js API route handlers.
 * Aligned with the new unified ApiResponse<T> gateway envelope.
 */

import { NextResponse } from "next/server";
import type { JsonMap, GatewayResponse } from "./shared";
import { unwrapGateway, extractGatewayError } from "./shared";
export type { JsonMap };
export { unwrapGateway, extractGatewayError };

export function getApiBaseUrl(): string {
  return (
    process.env.INSURETECH_API_BASE_URL ??
    process.env.NEXT_PUBLIC_INSURETECH_API_BASE_URL ??
    "http://localhost:8080"
  );
}

export function getCookieValue(cookieHeader: string, name: string): string {
  const match = cookieHeader.match(new RegExp(`(?:^|;\\s*)${name}=([^;]*)`));
  return match ? decodeURIComponent(match[1]) : "";
}

export function getCsrfToken(cookieHeader: string): string {
  return getCookieValue(cookieHeader, "csrf_token");
}

export function getRecord(value: unknown): JsonMap {
  if (value && typeof value === "object" && !Array.isArray(value))
    return value as JsonMap;
  return {};
}

export function getStringField(source: JsonMap, ...keys: string[]): string {
  for (const key of keys) {
    const value = source[key];
    if (typeof value === "string" && value.trim()) return value;
  }
  return "";
}

export function getNumberField(source: JsonMap, ...keys: string[]): number {
  for (const key of keys) {
    const value = source[key];
    if (typeof value === "number") return value;
    if (typeof value === "string" && value.trim()) {
      const parsed = Number.parseInt(value, 10);
      if (!Number.isNaN(parsed)) return parsed;
    }
  }
  return 0;
}

export function parseMoneyDecimal(value: unknown): number {
  if (value == null) return 0;
  if (typeof value === "bigint") return Number(value) / 100;
  if (typeof value === "number") return value;
  if (typeof value === "string") {
    const p = Number.parseFloat(value);
    return Number.isNaN(p) ? 0 : p;
  }
  if (typeof value === "object") {
    const bag = value as JsonMap;
    const decimal = bag.decimal_amount ?? bag.decimalAmount;
    if (typeof decimal === "number") return decimal;
    if (typeof decimal === "string") {
      const p = Number.parseFloat(decimal);
      if (!Number.isNaN(p)) return p;
    }
    const amount = bag.amount;
    if (typeof amount === "number") return amount / 100;
    if (typeof amount === "string") {
      const p = Number.parseFloat(amount);
      if (!Number.isNaN(p)) return p;
    }
  }
  return 0;
}

// ─── Standard BFF response builders ──────────────────────────────────────────

/**
 * sdkErrorMessage — extracts a user-facing error message from an SDK result.
 * The SDK result now carries a GatewayResponse envelope in result.data.
 */
export function sdkErrorMessage(result: unknown): string {
  return extractGatewayError(result);
}

/**
 * unwrapSdkResult — extracts the payload from an SDK call result.
 * The SDK client interceptor already strips the ApiResponse envelope,
 * so result.data is the inner T directly.
 */
export function unwrapSdkResult<T>(result: {
  data?: T | null;
  error?: unknown;
  response?: Response;
}) {
  if (result.error || !result.response?.ok) {
    const msg =
      typeof result.error === "string"
        ? result.error
        : result.error &&
            typeof result.error === "object" &&
            "message" in result.error
          ? String((result.error as Record<string, unknown>).message)
          : "An unexpected error occurred";
    return {
      ok: false as const,
      message: msg,
      code: "ERROR",
      status: result.response?.status ?? 500,
    };
  }
  return { ok: true as const, data: result.data as T, meta: undefined };
}

export function badRequest(message: string): NextResponse {
  return NextResponse.json({ ok: false, message }, { status: 400 });
}

export function gatewayError(message: string, status = 502): NextResponse {
  return NextResponse.json({ ok: false, message }, { status });
}

export function notFound(message = "Not found"): NextResponse {
  return NextResponse.json({ ok: false, message }, { status: 404 });
}

export function unauthorized(message = "Unauthorized"): NextResponse {
  return NextResponse.json({ ok: false, message }, { status: 401 });
}

export function forbidden(message = "Forbidden"): NextResponse {
  return NextResponse.json({ ok: false, message }, { status: 403 });
}

export function internalError(message = "Internal error"): NextResponse {
  return NextResponse.json({ ok: false, message }, { status: 500 });
}
