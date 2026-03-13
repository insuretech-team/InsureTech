/**
 * shared.ts
 * ─────────
 * Universal browser/server shared primitives.
 * Single source of truth for parseJson, ApiResult, and JsonMap.
 *
 * Safe to import anywhere — no runtime dependencies, no Next.js server APIs.
 * All browser-side clients import from here instead of duplicating these utilities.
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
