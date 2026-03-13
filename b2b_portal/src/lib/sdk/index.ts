/**
 * src/lib/sdk/index.ts
 * ─────────────────────
 * Single unified entry point for all SDK, client, and helper utilities.
 *
 * Import boundary rule:
 *   - Components / hooks / pages  → import from "@lib/sdk"  (this file)
 *   - API route handlers          → import from "@lib/sdk"  (this file)
 *   - Files inside src/lib/sdk/   → import directly from sibling (e.g. "./shared")
 *
 * Do NOT import from @lib/clients — that folder is deprecated and will be deleted.
 */

// ─── Shared primitives (parseJson, ApiResult, JsonMap) ────────────────────────
export type { ApiResult, JsonMap } from "./shared";
export { parseJson } from "./shared";

// ─── Server-side helpers (Next.js API routes only) ────────────────────────────
// NOTE: intentionally NOT re-exported here — import directly by path in API routes:
//   import { sdkErrorMessage, badRequest, gatewayError, ... } from "@lib/sdk/api-helpers"

// ─── Server-side SDK clients ──────────────────────────────────────────────────
// NOTE: intentionally NOT re-exported here — import directly by path in API routes:
//   import { makeSdkClient } from "@lib/sdk/b2b-sdk-client"
//   import { makeDocgenSdkClient } from "@lib/sdk/docgen-sdk-client"
//   import { resolvePortalHeaders } from "@lib/sdk/session-headers"
//   import { sdkErrorMessage, ... } from "@lib/sdk/api-helpers"

// ─── Browser-side clients (used in components & hooks) ────────────────────────
export { authClient } from "./auth-client";
export type {
  AuthOkResponse,
  ProfileResponse,
  SessionsResponse,
  TotpResponse,
  OtpResponse,
  ProfilePhotoUrlResponse,
} from "./auth-client";

export { departmentClient } from "./department-client";
export type { DepartmentListResult, DepartmentSingleResult } from "./department-client";

export { employeeClient } from "./employee-client";
export type {
  EmployeeCreatePayload,
  EmployeeUpdatePayload,
  EmployeeListResult,
  EmployeeFullRecord,
  EmployeeSingleResult,
} from "./employee-client";

export { organisationClient } from "./organisation-client";
export type {
  OrgCreatePayload,
  OrgUpdatePayload,
  OrgAdminCreatePayload,
  OrgListResult,
  OrgSingleResult,
  OrgMembersResult,
  OrgMemberResult,
  OrgMember,
  OrgMemberRole,
  OrgMemberStatus,
} from "./organisation-client";

export { purchaseOrderClient } from "./purchase-order-client";
export type {
  CatalogItem,
  PurchaseOrderCreatePayload,
  PurchaseOrderUpdatePayload,
  POListResult,
  POSingleResult,
  POCatalogResult,
} from "./purchase-order-client";

export { docgenClient } from "./docgen-client";
export type {
  DocumentStatus,
  DocumentRecord,
  GenerateDocumentPayload,
  DocumentListResult,
  DocumentSingleResult,
  DocumentDownloadResult,
} from "./docgen-client";

export { b2bDashboardClient } from "./dashboard-config";
