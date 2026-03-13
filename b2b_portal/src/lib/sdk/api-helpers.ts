/**
 * api-helpers.ts
 * ──────────────
 * Shared utilities for Next.js API route handlers.
 */

import { NextResponse } from "next/server";

import type { JsonMap } from "./shared";
export type { JsonMap };

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
  if (value && typeof value === "object" && !Array.isArray(value)) return value as JsonMap;
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
    if (typeof decimal === "string") { const p = Number.parseFloat(decimal); if (!Number.isNaN(p)) return p; }
    const raw = bag.amount;
    if (typeof raw === "bigint") return Number(raw) / 100;
    if (typeof raw === "number") return raw / 100;
    if (typeof raw === "string") { const p = Number.parseFloat(raw); return Number.isNaN(p) ? 0 : p / 100; }
  }
  return 0;
}

/**
 * Safely extract error message from a @hey-api/client-fetch SDK result.
 *
 * The SDK result union is:
 *   { data: T; error: undefined; request: Request; response: Response }
 * | { data: undefined; error: Error; request: Request; response: Response }
 *
 * So `error` only exists on the second branch. Use `result.response.ok` to
 * discriminate, then cast to extract the message.
 */
export function sdkErrorMessage(result: { response: Response; error?: unknown }): string {
  const err = result.error;
  if (!err) return "Request failed";
  if (err instanceof Error) return err.message;
  if (err && typeof err === "object") {
    const m = (err as JsonMap).message;
    if (typeof m === "string") return m;
  }
  return "Request failed";
}

export function badRequest(message: string) {
  return NextResponse.json({ ok: false, message }, { status: 400 });
}

export function gatewayError(message: string) {
  return NextResponse.json({ ok: false, message }, { status: 502 });
}
