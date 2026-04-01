/**
 * shared.ts
 * ─────────
 * Universal browser/server shared primitives for customer portal.
 * Single source of truth for parseJson, ApiResult, and JsonMap.
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

/** Standard shape returned by every /api/* route handler. */
export type ApiResult<T extends object = object> = {
  ok: boolean;
  message?: string;
} & T;
