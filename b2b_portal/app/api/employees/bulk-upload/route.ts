/**
 * /api/employees/bulk-upload  POST
 *
 * Accepts multipart/form-data:
 *   - file        : Excel (.xlsx) or CSV file
 *   - business_id : organisation UUID
 *
 * Proxies directly to the gateway POST /v1/b2b/employees/bulk-upload
 * preserving the multipart body (no re-parsing needed).
 */
import { NextResponse } from "next/server";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";

function getBaseUrl(): string {
  return (
    process.env.INSURETECH_API_BASE_URL ??
    process.env.NEXT_PUBLIC_INSURETECH_API_BASE_URL ??
    "http://localhost:8080"
  );
}

export async function POST(request: Request) {
  try {
    const hdrs = await resolvePortalHeaders(request);

    const cookieHeader = request.headers.get("cookie") ?? "";
    const csrf =
      cookieHeader.match(/(?:^|;\s*)csrf_token=([^;]*)/)?.[1] ?? "";

    // Forward the raw multipart body to the gateway unchanged.
    // We copy all headers except host so the gateway can validate CSRF
    // and parse the multipart boundary correctly.
    const forwardHeaders: Record<string, string> = {};
    if (cookieHeader) forwardHeaders["cookie"] = cookieHeader;
    if (csrf) forwardHeaders["X-CSRF-Token"] = decodeURIComponent(csrf);

    // Forward portal context
    const xPortal = hdrs?.portal ?? request.headers.get("x-portal") ?? "";
    const xBusinessId =
      hdrs?.businessId ?? request.headers.get("x-business-id") ?? "";
    const xUserId =
      hdrs?.userId ?? request.headers.get("x-user-id") ?? "";
    const xTenantId =
      hdrs?.tenantId ?? request.headers.get("x-tenant-id") ?? "";
    if (xPortal) forwardHeaders["x-portal"] = xPortal;
    if (xBusinessId) forwardHeaders["x-business-id"] = xBusinessId;
    if (xUserId) forwardHeaders["x-user-id"] = xUserId;
    if (xTenantId) forwardHeaders["x-tenant-id"] = xTenantId;

    // Forward content-type (multipart boundary must be preserved)
    const ct = request.headers.get("content-type");
    if (ct) forwardHeaders["content-type"] = ct;

    const body = await request.arrayBuffer();

    const gatewayRes = await fetch(
      `${getBaseUrl()}/v1/b2b/employees/bulk-upload`,
      {
        method: "POST",
        headers: forwardHeaders,
        body,
        cache: "no-store",
      }
    );

    const raw = await gatewayRes.text();
    let data: Record<string, unknown>;
    try {
      data = raw ? (JSON.parse(raw) as Record<string, unknown>) : {};
    } catch {
      data = { ok: false, message: raw || "Gateway error" };
    }

    return NextResponse.json(data, { status: gatewayRes.status });
  } catch (err) {
    return NextResponse.json(
      { ok: false, message: err instanceof Error ? err.message : "Error" },
      { status: 502 }
    );
  }
}
