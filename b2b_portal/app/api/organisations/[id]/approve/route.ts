/**
 * POST /api/organisations/[id]/approve
 * Approves (activates) a pending organisation by setting its status to ACTIVE.
 * Only Super Admin can call this — enforced by the backend AuthZ interceptor.
 */
import { NextResponse } from "next/server";
import { makeDirectHttp } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import type { Organisation } from "@lifeplus/insuretech-sdk";
import type { Organisation as UiOrg } from "@lib/types/b2b";

type RouteContext = { params: Promise<{ id: string }> };

function mapOrg(org: Organisation): UiOrg {
  return {
    id: org.organisation_id ?? "",
    name: org.name ?? "",
    code: org.code ?? "",
    industry: org.industry ?? "",
    contactEmail: org.contact_email ?? "",
    contactPhone: org.contact_phone ?? "",
    address: org.address ?? "",
    status: org.status ?? "ORGANISATION_STATUS_ACTIVE",
    totalEmployees: org.total_employees ?? 0,
    createdAt: org.created_at ?? "",
  };
}

export async function POST(request: Request, { params }: RouteContext) {
  const { id } = await params;
  if (!id) {
    return NextResponse.json({ ok: false, message: "organisation_id required" }, { status: 400 });
  }
  try {
    const hdrs = await resolvePortalHeaders(request);
    if (!hdrs) {
      return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
    }
    const http = makeDirectHttp(request, hdrs);
    // Approve = patch status to ACTIVE
    const result = await http.patch(`/v1/b2b/organisations/${id}`, {
      status: "ORGANISATION_STATUS_ACTIVE",
    });
    if (!result.ok) {
      return NextResponse.json(
        { ok: false, message: (result.data as Record<string, unknown>)?.error ?? "Approve failed" },
        { status: result.status }
      );
    }
    const org = (result.data as Record<string, unknown>)?.organisation as Organisation | undefined;
    return NextResponse.json({
      ok: true,
      message: "Organisation approved and activated",
      organisation: org ? mapOrg(org) : null,
    });
  } catch (err) {
    return NextResponse.json(
      { ok: false, message: err instanceof Error ? err.message : "Error" },
      { status: 502 }
    );
  }
}
