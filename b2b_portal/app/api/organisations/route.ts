/**
 * /api/organisations  GET | POST
 */
import { NextResponse } from "next/server";
import { makeDirectHttp, makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { unwrapSdkResult } from "@lib/sdk/api-helpers";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import type { Organisation } from "@lifeplus/insuretech-sdk";
import type { Organisation as UiOrg } from "@lib/types/b2b";
import { getBangladeshMobileValidationMessage, normalizeBangladeshMobile } from "@lib/utils/bd-mobile";
import { getPasswordValidationMessage } from "@/src/lib/utils/password";

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

function normalizeOrganisationCode(name: string): string {
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
    const unwrapped = unwrapSdkResult(result);
    if (!unwrapped.ok) return NextResponse.json({ ok: false, message: unwrapped.message, organisations: [] }, { status: unwrapped.status });
    const d = unwrapped.data as Record<string, unknown>;
    return NextResponse.json({ ok: true, organisations: ((d?.organisations ?? []) as Organisation[]).map(mapOrg) });
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
    if (!name) {
      return NextResponse.json({ ok: false, message: "Organisation name is required" }, { status: 400 });
    }
    const contactPhone = body.contactPhone ? String(body.contactPhone).trim() : "";
    const normalizedContactPhone = contactPhone ? normalizeBangladeshMobile(contactPhone) : null;
    if (contactPhone && !normalizedContactPhone) {
      return NextResponse.json({ ok: false, message: getBangladeshMobileValidationMessage("Contact phone") }, { status: 400 });
    }
    const adminAssignment = body.adminAssignment as Record<string, unknown> | undefined;
    const assignmentUserId = typeof adminAssignment?.userId === "string" ? adminAssignment.userId.trim() : "";
    const assignmentTemporaryPassword =
      typeof adminAssignment?.temporaryPassword === "string" ? adminAssignment.temporaryPassword : "";
    if (adminAssignment) {
      if (!assignmentUserId) {
        return NextResponse.json({ ok: false, message: "Assigned admin userId is required" }, { status: 400 });
      }
      const passwordError = getPasswordValidationMessage(assignmentTemporaryPassword, "Temporary password");
      if (passwordError) {
        return NextResponse.json({ ok: false, message: passwordError }, { status: 400 });
      }
    }

    const result = await sdk.createOrganisation({
      body: {
        tenant_id: tenantId,
        name,
        code: normalizeOrganisationCode(name),
        industry: body.industry ? String(body.industry).trim() : undefined,
        contact_email: body.contactEmail ? String(body.contactEmail).trim() : undefined,
        contact_phone: normalizedContactPhone ?? undefined,
        address: body.address ? String(body.address).trim() : undefined,
      },
    });
    const unwrapped = unwrapSdkResult(result);
    if (!unwrapped.ok) return NextResponse.json({ ok: false, message: unwrapped.message }, { status: unwrapped.status });
    const organisation = unwrapped.data?.organisation ? mapOrg(unwrapped.data.organisation) : null;
    const organisationID = unwrapped.data?.organisation?.organisation_id ?? "";
    if (adminAssignment && organisationID) {
      const assignResult = await http.post(`/v1/b2b/organisations/${organisationID}/admins:assign`, {
        userId: assignmentUserId,
        temporaryPassword: assignmentTemporaryPassword,
      });
      if (!assignResult.ok) {
        return NextResponse.json(
          {
            ok: false,
            message: `Organisation created but admin assignment failed: ${String(assignResult.data.message ?? "Unknown error")}`,
            organisation,
          },
          { status: assignResult.status }
        );
      }
    } else if (body.admin && organisationID) {
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
      { ok: true, message: "Organisation created", organisation },
      { status: 201 }
    );
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}
