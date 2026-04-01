/**
 * src/lib/sdk/index.ts — Partner Portal
 * ──────────────────────────────────────
 * Single unified entry point for all SDK client utilities.
 */

export type { ApiResult, JsonMap } from "./shared";
export { parseJson } from "./shared";

export { authClient } from "./auth-client";
export type { AuthOkResponse, ProfileResponse, OtpResponse } from "./auth-client";

export { commissionClient } from "./commission-client";
export type { CommissionListResult, CommissionSingleResult } from "./commission-client";

export { policyClient } from "./policy-client";
export type { PolicyListResult, PolicySingleResult } from "./policy-client";
