export { authClient } from "./auth-client";
export { b2bDashboardClient } from "./b2b-dashboard-client";

// ─── B2B domain clients (browser-side, use typed SDK routes) ──────────────────
export { departmentClient } from "./department-client";
export type { DepartmentListResult, DepartmentSingleResult } from "./department-client";

export { employeeClient } from "./employee-client";
export type { EmployeeCreatePayload, EmployeeUpdatePayload, EmployeeListResult, EmployeeSingleResult } from "./employee-client";

export { purchaseOrderClient } from "./purchase-order-client";
export type { CatalogItem, PurchaseOrderCreatePayload, PurchaseOrderUpdatePayload, POListResult, POSingleResult, POCatalogResult } from "./purchase-order-client";

export { organisationClient } from "./organisation-client";
export type { OrgCreatePayload, OrgUpdatePayload, OrgListResult, OrgSingleResult } from "./organisation-client";

// ─── DocGen client (browser-side, calls /api/documents/*) ─────────────────────
export { docgenClient } from "./docgen-client";
export type {
  DocumentRecord,
  DocumentStatus,
  GenerateDocumentPayload,
  DocumentListResult,
  DocumentSingleResult,
  DocumentDownloadResult,
} from "./docgen-client";

// ─── Shared API client utilities ──────────────────────────────────────────────
export type { JsonMap, ApiResult } from "./api-client";
