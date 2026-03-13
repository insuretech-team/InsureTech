/**
 * /api/organisations/me  GET
 */
import { NextResponse } from "next/server";
import { getCurrentSession, toPortalSessionFromCurrentSession } from "@lib/auth/backend-auth";
import { makeDirectHttp, makeSdkClient } from "@lib/sdk/b2b-sdk-client";
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

export async function GET(request: Request) {
  try {
    const cookieHeader = request.headers.get("cookie") ?? "";
    if (cookieHeader) {
      const sessionResult = await getCurrentSession(cookieHeader);
      if (!sessionResult.error && sessionResult.data) {
        const session = await toPortalSessionFromCurrentSession(sessionResult.data, cookieHeader);
        if (session?.principal.role === "SYSTEM_ADMIN") {
          return NextResponse.json({ ok: true, organisation: null });
        }
      }
    }

    const http = makeDirectHttp(request);
    const result = await http.get("/v1/b2b/organisations/me");
    if (result.ok) {
      const organisationId = String(result.data.organisation_id ?? "");
      if (!organisationId) {
        return NextResponse.json({ ok: true, organisation: null });
      }

      const organisationResult = await makeSdkClient(request).getOrganisation({
        path: { organisation_id: organisationId },
      });
      if (!organisationResult.response.ok || !organisationResult.data?.organisation) {
        return NextResponse.json({ ok: false, message: "Failed to load resolved organisation" }, { status: organisationResult.response.status });
      }
      return NextResponse.json({ ok: true, organisation: mapOrg(organisationResult.data.organisation as Organisation) });
    }

    if (result.status === 403 || result.status === 404) {
      return NextResponse.json({ ok: true, organisation: null });
    }

    return NextResponse.json(
      { ok: false, message: String(result.data.message ?? "Failed to resolve organisation context") },
      { status: result.status }
    );
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}
