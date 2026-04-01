/**
 * src/lib/sdk/index.ts
 * ─────────────────────
 * Single unified entry point for all SDK, client, and helper utilities.
 *
 * Import boundary rule:
 *   - Components / hooks / pages  → import from "@lib/sdk/b2b-sdk-client" (bffClient)
 *   - API route handlers          → import from "@lib/sdk/b2b-sdk-client" (makeSdkClient / makeDirectHttp)
 *   - Files inside src/lib/sdk/   → import directly from sibling (e.g. "./shared")
 *
 * All browser-side clients are consolidated into bffClient in b2b-sdk-client.ts.
 * All server-side SDK calls go through makeSdkClient / makeDirectHttp.
 */

// ─── Shared primitives ────────────────────────────────────────────────────────
export type { ApiResult, JsonMap } from "./shared";
export { parseJson } from "./shared";

// ─── Single unified client (browser + server) ─────────────────────────────────
// Browser-side: use bffClient (calls /api/* BFF routes)
// Server-side:  use makeSdkClient / makeDirectHttp (calls gateway directly)
export { bffClient } from "./b2b-sdk-client";
export type {
  // Auth
  AuthOkResponse, ProfileResponse, SessionsResponse,
  TotpResponse, OtpResponse, ProfilePhotoUrlResponse,
  // Employees
  EmployeeCreatePayload, EmployeeUpdatePayload,
  EmployeeListResult, EmployeeFullRecord, EmployeeSingleResult,
  // Departments
  DepartmentListResult, DepartmentSingleResult,
  // Organisations
  OrgCreatePayload, OrgUpdatePayload, OrgAdminCreatePayload,
  OrgListResult, OrgSingleResult, OrgMembersResult, OrgMemberResult,
  OrgMember, OrgMemberRole, OrgMemberStatus,
  // Purchase Orders
  CatalogItem, PurchaseOrderCreatePayload, PurchaseOrderUpdatePayload,
  POListResult, POSingleResult, POCatalogResult,
  // Documents
  DocumentStatus, DocumentRecord, GenerateDocumentPayload,
  DocumentListResult, DocumentSingleResult, DocumentDownloadResult,
  // SDK client type
  B2bSdkClient,
} from "./b2b-sdk-client";

// ─── Dashboard config (static mock data — no network calls) ───────────────────
export { b2bDashboardClient } from "./dashboard-config";

// ─── Server-side helpers (Next.js API routes only — import by path) ───────────
// import { makeSdkClient, makeDirectHttp, makeDocgenClient } from "@lib/sdk/b2b-sdk-client"
// import { resolvePortalHeaders } from "@lib/sdk/session-headers"
// import { sdkErrorMessage, badRequest, gatewayError, ... } from "@lib/sdk/api-helpers"
export type { DocgenClient } from "./b2b-sdk-client";
