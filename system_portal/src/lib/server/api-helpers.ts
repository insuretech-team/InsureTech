/**
 * api-helpers.ts
 * ─────────────────
 * Shared utilities for SvelteKit server routes.
 * Aligned with the unified ApiResponse<T> gateway envelope.
 */

import { error, json, type RequestEvent } from '@sveltejs/kit';

export type JsonMap = Record<string, unknown>;

// ─── Standard Response Types ──────────────────────────────────────────────────

export interface GatewayError {
  code: string;
  message: string;
  error_id: string;
  http_status_code: number;
  retryable: boolean;
  field_violations: Array<{ field: string; description: string }>;
}

export interface GatewayMeta {
  request_id: string;
  timestamp: string;
  pagination?: {
    page: number;
    page_size: number;
    total_count: number;
    total_pages: number;
    has_next: boolean;
    has_prev: boolean;
  } | null;
}

/**
 * GatewayResponse<T> — the unified envelope returned by every gateway endpoint.
 */
export interface GatewayResponse<T> {
  success: boolean;
  data: T | null;
  error: GatewayError | null;
  meta: GatewayMeta;
}

/**
 * Standard API result shape for SvelteKit routes.
 */
export type ApiResult<T extends object = object> = {
  ok: boolean;
  message?: string;
} & T;

// ─── Cookie Helpers ───────────────────────────────────────────────────────────

export function extractCookie(cookieHeader: string, name: string): string {
  const escaped = name.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
  const m = cookieHeader.match(new RegExp(`(?:^|;\\s*)${escaped}=([^;]*)`));
  return m ? decodeURIComponent(m[1]) : '';
}

export function getCsrfToken(cookieHeader: string): string {
  return extractCookie(cookieHeader, 'csrf_token');
}

// ─── Data Extraction Helpers ──────────────────────────────────────────────────

export function getRecord(value: unknown): JsonMap {
  if (value && typeof value === 'object' && !Array.isArray(value)) return value as JsonMap;
  return {};
}

export function getStringField(source: JsonMap, ...keys: string[]): string {
  for (const key of keys) {
    const value = source[key];
    if (typeof value === 'string' && value.trim()) return value;
  }
  return '';
}

export function getNumberField(source: JsonMap, ...keys: string[]): number {
  for (const key of keys) {
    const value = source[key];
    if (typeof value === 'number') return value;
    if (typeof value === 'string' && value.trim()) {
      const parsed = Number.parseInt(value, 10);
      if (!Number.isNaN(parsed)) return parsed;
    }
  }
  return 0;
}

export function parseMoneyDecimal(value: unknown): number {
  if (value == null) return 0;
  if (typeof value === 'bigint') return Number(value) / 100;
  if (typeof value === 'number') return value;
  if (typeof value === 'string') {
    const p = Number.parseFloat(value);
    return Number.isNaN(p) ? 0 : p;
  }
  if (typeof value === 'object') {
    const bag = value as JsonMap;
    const decimal = bag.decimal_amount ?? bag.decimalAmount;
    if (typeof decimal === 'number') return decimal;
    if (typeof decimal === 'string') {
      const p = Number.parseFloat(decimal);
      if (!Number.isNaN(p)) return p;
    }
    const amount = bag.amount;
    if (typeof amount === 'number') return amount / 100;
    if (typeof amount === 'string') {
      const p = Number.parseFloat(amount);
      if (!Number.isNaN(p)) return p;
    }
  }
  return 0;
}

// ─── SDK Result Helpers ───────────────────────────────────────────────────────

/**
 * Extracts a user-facing error message from an SDK result.
 */
export function sdkErrorMessage(result: unknown): string {
  return extractGatewayError(result);
}

/**
 * Unwraps the payload from an SDK call result.
 * The SDK client interceptor already strips the ApiResponse envelope,
 * so result.data is the inner T directly.
 */
export function unwrapSdkResult<T>(result: { data?: T | null; error?: unknown; response?: Response }) {
  if (result.error || !result.response?.ok) {
    const msg =
      typeof result.error === 'string'
        ? result.error
        : result.error && typeof result.error === 'object' && 'message' in result.error
          ? String((result.error as Record<string, unknown>).message)
          : 'An unexpected error occurred';
    return { ok: false as const, message: msg, code: 'ERROR', status: result.response?.status ?? 500 };
  }
  return { ok: true as const, data: result.data as T, meta: undefined };
}

/**
 * Extracts a user-facing error message from any SDK result or envelope.
 */
export function extractGatewayError(result: unknown): string {
  if (!result || typeof result !== 'object') return 'An unexpected error occurred';
  const r = result as Record<string, unknown>;

  // SDK result shape: { error: GatewayError, response }
  if ('error' in r && r.error && typeof r.error === 'object') {
    const e = r.error as Record<string, unknown>;
    if (typeof e.message === 'string') return e.message;
  }

  // Legacy: data.success === false (pre-interceptor shape)
  if ('data' in r && r.data && typeof r.data === 'object') {
    const d = r.data as Record<string, unknown>;
    if (
      'success' in d &&
      d.success === false &&
      'error' in d &&
      d.error &&
      typeof d.error === 'object'
    ) {
      const e = d.error as Record<string, unknown>;
      return typeof e.message === 'string' ? e.message : 'An unexpected error occurred';
    }
    if (typeof d.message === 'string') return d.message;
  }

  // Fallback
  if (typeof r.message === 'string') return r.message;

  return 'An unexpected error occurred';
}

/**
 * Unwraps data from a GatewayResponse envelope.
 */
export function unwrapGateway<T>(
  body: GatewayResponse<T>,
  httpStatus?: number
):
  | { ok: true; data: T; meta: GatewayMeta }
  | { ok: false; message: string; code: string; status: number; retryable: boolean } {
  if (body.success && body.data !== null && body.data !== undefined) {
    return { ok: true, data: body.data, meta: body.meta };
  }
  const err = body.error;
  return {
    ok: false,
    message: err?.message ?? 'An unexpected error occurred',
    code: err?.code ?? 'UNKNOWN_ERROR',
    status: err?.http_status_code ?? httpStatus ?? 500,
    retryable: err?.retryable ?? false,
  };
}

// ─── Response Builders ────────────────────────────────────────────────────────

export function okResponse<T extends object>(data: T, status = 200): Response {
  return json({ ok: true, ...data } satisfies ApiResult<T>, { status });
}

export function errorResponse(message: string, status = 400): Response {
  return json({ ok: false, message } satisfies ApiResult, { status });
}

export function badRequest(message: string): Response {
  return errorResponse(message, 400);
}

export function unauthorized(message = 'Unauthorized'): Response {
  return errorResponse(message, 401);
}

export function forbidden(message = 'Forbidden'): Response {
  return errorResponse(message, 403);
}

export function notFound(message = 'Not found'): Response {
  return errorResponse(message, 404);
}

export function gatewayError(message: string, status = 502): Response {
  return errorResponse(message, status);
}

export function internalError(message = 'Internal error'): Response {
  return errorResponse(message, 500);
}

// ─── SvelteKit Error Helpers ────────────────────────────────────────────────────

export function throwUnauthorized(message = 'Unauthorized') {
  throw error(401, message);
}

export function throwForbidden(message = 'Forbidden') {
  throw error(403, message);
}

export function throwNotFound(message = 'Not found') {
  throw error(404, message);
}

// ─── Auth Check Helpers ─────────────────────────────────────────────────────────

export function requireUser(event: RequestEvent): NonNullable<App.Locals['user']> {
  if (!event.locals.user) {
    throwUnauthorized('Authentication required');
  }
  return event.locals.user as NonNullable<App.Locals['user']>;
}

export function requireRole(
  event: RequestEvent,
  allowedRoles: string[],
  currentRole?: string
): string {
  const role = currentRole ?? event.locals.user?.role;
  if (!role || !allowedRoles.includes(role)) {
    throwForbidden('Insufficient permissions');
  }
  return String(role);
}

// ─── Request Parsing Helpers ────────────────────────────────────────────────────

export async function parseJsonBody<T = JsonMap>(event: RequestEvent): Promise<T> {
  try {
    const body = await event.request.json();
    return body as T;
  } catch {
    throw error(400, 'Invalid JSON body');
  }
}

export function getQueryParam(event: RequestEvent, key: string): string | null {
  const url = new URL(event.request.url);
  return url.searchParams.get(key);
}

export function getQueryParams(event: RequestEvent): URLSearchParams {
  const url = new URL(event.request.url);
  return url.searchParams;
}
