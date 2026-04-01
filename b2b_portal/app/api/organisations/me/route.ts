/**
 * /api/organisations/me  GET
 */
import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
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

function mapResolvedOrg(organisationId: string, organisationName: string): UiOrg {
  return {
    id: organisationId,
    name: organisationName,
    code: "",
    industry: "",
    contactEmail: "",
    contactPhone: "",
    address: "",
    status: "ORGANISATION_STATUS_ACTIVE",
    totalEmployees: 0,
    createdAt: "",
  };
}

export async function GET(request: Request) {
  try {
    const hdrs = await resolvePortalHeaders(request);
    if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

    const sdk = makeSdkClient(request, hdrs);

    // Super admins have no own org — skip the lookup entirely.
    if (hdrs.portal === "PORTAL_SYSTEM") {
      return NextResponse.json({ ok: true, organisation: null });
    }

    // Use the SDK's getMyOrganisation helper (routes via makeDirectHttp, no hardcoded fetch).
    const result = await sdk.getMyOrganisation();
    if (result.ok) {
      const organisationId = String(result.data.organisation_id ?? "");
      if (!organisationId) {
        return NextResponse.json({ ok: true, organisation: null });
      }
      const organisationResult = await sdk.getOrganisation({ path: { organisation_id: organisationId } });
      if (organisationResult.response.ok && organisationResult.data?.organisation) {
        return NextResponse.json({ ok: true, organisation: mapOrg(organisationResult.data.organisation as Organisation) });
      }

      if (organisationResult.response.status === 403 || organisationResult.response.status === 404) {
        return NextResponse.json({
          ok: true,
          organisation: mapResolvedOrg(organisationId, String(result.data.organisation_name ?? "")),
        });
      }

      if (!organisationResult.response.ok || !organisationResult.data?.organisation) {
        return NextResponse.json(
          { ok: false, message: "Failed to load resolved organisation" },
          { status: organisationResult.response.status }
        );
      }
    }

    if (result.status === 403 || result.status === 404) {
      return NextResponse.json({ ok: true, organisation: null });
    }

    return NextResponse.json(
      { ok: false, message: String(result.data?.message ?? "Failed to resolve organisation context") },
      { status: result.status }
    );
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}
