/**
 * src/lib/sdk/index.ts — Customer Portal
 * ───────────────────────────────────────
 * Single unified entry point for all SDK client utilities.
 */

// ─── Shared primitives ────────────────────────────────────────────────────────
export type { ApiResult, JsonMap } from "./shared";
export { parseJson } from "./shared";

// ─── Browser-side clients ─────────────────────────────────────────────────────
export { authClient } from "./auth-client";
export type { AuthOkResponse, ProfileResponse, SessionsResponse, OtpResponse } from "./auth-client";

export { policyClient } from "./policy-client";
export type { PolicyListResult, PolicySingleResult } from "./policy-client";

export { claimClient } from "./claim-client";
export type { ClaimListResult, ClaimSingleResult, ClaimCreatePayload } from "./claim-client";

export { paymentClient } from "./payment-client";
export type { PaymentListResult, PaymentSingleResult } from "./payment-client";

export { quotationClient } from "./quotation-client";
export type { QuotationListResult, QuotationSingleResult, QuotationCreatePayload } from "./quotation-client";
