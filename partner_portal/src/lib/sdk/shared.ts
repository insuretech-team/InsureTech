/**
 * shared.ts
 * Universal browser/server primitives aligned with the new ApiResponse<T> envelope.
 */

export type JsonMap = Record<string, unknown>;

/** Parse JSON from a fetch Response, throwing on non-JSON content-types. */
export async function parseJson<T>(response: Response): Promise<T> {
  const ct = response.headers.get("content-type") ?? "";
  if (!ct.includes("application/json")) {
    throw new Error(`Unexpected response type (status ${response.status})`);
  }
  return (await response.json()) as T;
}

/** Standard shape returned by every /api/* BFF route handler to the browser. */
export type ApiResult<T extends object = object> = {
  ok: boolean;
  message?: string;
} & T;

// ─── Gateway envelope types (mirrors respond package in Go gateway) ────────────

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
 * Maps to ApiResponse<T> in openapi.yaml and the Go respond package.
 */
export interface GatewayResponse<T> {
  success: boolean;
  data: T | null;
  error: GatewayError | null;
  meta: GatewayMeta;
}

/**
 * unwrapGateway — extracts data from a GatewayResponse.
 * Returns { ok: true, data } on success, { ok: false, message, code } on error.
 * Use in BFF API route handlers to safely extract typed data.
 */
export function unwrapGateway<T>(
  body: GatewayResponse<T>,
  httpStatus?: number
):
  | { ok: true; data: T; meta: GatewayMeta }
  | {
      ok: false;
      message: string;
      code: string;
      status: number;
      retryable: boolean;
    } {
  if (body.success && body.data !== null && body.data !== undefined) {
    return { ok: true, data: body.data, meta: body.meta };
  }
  const err = body.error;
  return {
    ok: false,
    message: err?.message ?? "An unexpected error occurred",
    code: err?.code ?? "UNKNOWN_ERROR",
    status: err?.http_status_code ?? httpStatus ?? 500,
    retryable: err?.retryable ?? false,
  };
}

/**
 * extractGatewayError — extracts a user-facing error message from an SDK result
 * or any SDK result object. Safe to call on any shape.
 *
 * With the SDK interceptor, error responses have the gateway error object
 * directly in result.error (unwrapped from envelope.error).
 */
export function extractGatewayError(result: unknown): string {
  if (!result || typeof result !== "object")
    return "An unexpected error occurred";
  const r = result as Record<string, unknown>;

  // SDK result shape: { error: GatewayError, response }
  // The interceptor unwraps envelope.error into result.error directly.
  if ("error" in r && r.error && typeof r.error === "object") {
    const e = r.error as Record<string, unknown>;
    if (typeof e.message === "string") return e.message;
  }

  // Legacy: data.success === false (pre-interceptor shape)
  if ("data" in r && r.data && typeof r.data === "object") {
    const d = r.data as Record<string, unknown>;
    if (
      "success" in d &&
      d.success === false &&
      "error" in d &&
      d.error &&
      typeof d.error === "object"
    ) {
      const e = d.error as Record<string, unknown>;
      return typeof e.message === "string"
        ? e.message
        : "An unexpected error occurred";
    }
    if (typeof d.message === "string") return d.message;
  }

  // Fallback
  if (typeof r.message === "string") return r.message;

  return "An unexpected error occurred";
}
