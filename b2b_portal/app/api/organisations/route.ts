/**
 * /api/organisations  GET | POST
 */
import { NextResponse } from "next/server";
import { makeDirectHttp, makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { sdkErrorMessage } from "@lib/sdk/api-helpers";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import type { Organisation } from "@lifeplus/insuretech-sdk";
import type { Organisation as UiOrg } from "@lib/types/b2b";

function mapOrg(org: Organisation): UiOrg {
  return {
    id: org.organisation_id ?? "", name: org.name ?? "", code: org.code ?? "",
    industry: org.industry ?? "", contactEmail: org.contact_email ?? "",
    contactPhone: org.contact_phone ?? "", address: org.address ?? "",
    status: org.status ?? "ORGANISATION_STATUS_ACTIVE",
    totalEmployees: org.total_employees ?? 0, createdAt: org.created_at ?? "",
  };
}

function getTenantIdFallback(): string {
  return (
    process.env.DEFAULT_TENANT_ID?.trim() ||
    process.env.NEXT_PUBLIC_DEFAULT_TENANT_ID?.trim() ||
    "00000000-0000-0000-0000-000000000001"
  );
}

function normalizeOrganisationCode(name: string, rawCode: unknown): string {
  const provided = typeof rawCode === "string" ? rawCode : "";
  const sanitized = provided.toUpperCase().replace(/[^A-Z0-9]+/g, "-").replace(/^-+|-+$/g, "");
  if (sanitized) {
    return sanitized;
  }

  const fallbackBase = name.toUpperCase().replace(/[^A-Z0-9]+/g, "-").replace(/^-+|-+$/g, "").slice(0, 12) || "ORG";
  const suffix = Math.random().toString(36).slice(2, 6).toUpperCase();
  return `${fallbackBase}-${suffix}`;
}

export async function GET(request: Request) {
  try {
    const tenantId = getTenantIdFallback();
    const url = new URL(request.url);
    const hdrs = await resolvePortalHeaders(request);
    const result = await makeSdkClient(request, hdrs ?? undefined).listOrganisations({
      query: { tenant_id: tenantId, page_size: Number(url.searchParams.get("page_size") ?? 50) },
    });
    if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result), organisations: [] }, { status: result.response.status });
    return NextResponse.json({ ok: true, organisations: (result.data?.organisations ?? []).map(mapOrg) });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error", organisations: [] }, { status: 502 });
  }
}

export async function POST(request: Request) {
  try {
    const hdrs = await resolvePortalHeaders(request);
    const sdk = makeSdkClient(request, hdrs ?? undefined);
    const http = makeDirectHttp(request, hdrs ?? undefined);
    const body = (await request.json()) as Record<string, unknown>;
    const name = String(body.name ?? "").trim();
    const tenantId = getTenantIdFallback();

    const result = await sdk.createOrganisation({
      body: {
        tenant_id: tenantId,
        name,
        code: normalizeOrganisationCode(name, body.code),
        industry: body.industry ? String(body.industry).trim() : undefined,
        contact_email: body.contactEmail ? String(body.contactEmail).trim() : undefined,
        contact_phone: body.contactPhone ? String(body.contactPhone).trim() : undefined,
        address: body.address ? String(body.address).trim() : undefined,
      },
    });
    if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
    const organisation = result.data?.organisation ? mapOrg(result.data.organisation) : null;
    const organisationID = result.data?.organisation?.organisation_id ?? "";
    if (body.admin && organisationID) {
      // Pass hdrs so x-portal/x-user-id are forwarded — without them the backend
      // authz interceptor can't resolve the Casbin domain and returns 403.
      const adminResult = await http.post(`/v1/b2b/organisations/${organisationID}/admins`, body.admin);
      if (!adminResult.ok) {
        return NextResponse.json(
          {
            ok: false,
            message: `Organisation created but admin bootstrap failed: ${String(adminResult.data.message ?? "Unknown error")}`,
            organisation,
          },
          { status: adminResult.status }
        );
      }
    }
    return NextResponse.json(
      { ok: true, message: result.data?.message ?? "Organisation created", organisation },
      { status: 201 }
    );
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}
