/**
 * api-client.ts
 * ─────────────
 * Shared browser-side API client utilities.
 * All domain-specific clients (employee, department, purchase-order,
 * organisation) import from here to avoid duplicating parseJson.
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
