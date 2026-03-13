/**
 * /api/organisations/[id]  GET | PATCH | DELETE | POST
 */
import { NextResponse } from "next/server";
import { makeSdkClient, makeDirectHttp } from "@lib/sdk/b2b-sdk-client";
import { sdkErrorMessage } from "@lib/sdk/api-helpers";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import type { Organisation } from "@lifeplus/insuretech-sdk";
import type { Organisation as UiOrg } from "@lib/types/b2b";

type RouteContext = { params: Promise<{ id: string }> };
function mapOrg(org: Organisation): UiOrg {
  return {
    id: org.organisation_id ?? "", name: org.name ?? "", code: org.code ?? "",
    industry: org.industry ?? "", contactEmail: org.contact_email ?? "",
    contactPhone: org.contact_phone ?? "", address: org.address ?? "",
    status: org.status ?? "ORGANISATION_STATUS_ACTIVE",
    totalEmployees: org.total_employees ?? 0, createdAt: org.created_at ?? "",
  };
}

export async function GET(request: Request, { params }: RouteContext) {
  const { id } = await params;
  if (!id) return NextResponse.json({ ok: false, message: "organisation_id required" }, { status: 400 });
  try {
    const hdrs = await resolvePortalHeaders(request);
    const result = await makeSdkClient(request, hdrs ?? undefined).getOrganisation({ path: { organisation_id: id } });
    if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
    return NextResponse.json({ ok: true, organisation: result.data?.organisation ? mapOrg(result.data.organisation) : null });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}

export async function PATCH(request: Request, { params }: RouteContext) {
  const { id } = await params;
  if (!id) return NextResponse.json({ ok: false, message: "organisation_id required" }, { status: 400 });
  try {
    const hdrs = await resolvePortalHeaders(request);
    const body = (await request.json()) as Record<string, unknown>;
    // Only include non-empty fields to avoid sending empty strings to the backend
    const payload: Record<string, unknown> = {};
    if (body.name && String(body.name).trim()) payload.name = String(body.name).trim();
    if (body.industry) payload.industry = String(body.industry);
    if (body.contactEmail) payload.contact_email = String(body.contactEmail);
    if (body.contactPhone) payload.contact_phone = String(body.contactPhone);
    if (body.address) payload.address = String(body.address);
    if (body.status) payload.status = body.status;

    const http = makeDirectHttp(request, hdrs ?? undefined);
    const result = await http.patch(`/v1/b2b/organisations/${id}`, payload);
    if (!result.ok) return NextResponse.json({ ok: false, message: result.data?.error ?? "Update failed" }, { status: result.status });
    const org = (result.data as Record<string, unknown>)?.organisation as Organisation | undefined;
    return NextResponse.json({ ok: true, message: (result.data as Record<string, unknown>)?.message ?? "Updated", organisation: org ? mapOrg(org) : null });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}

export async function DELETE(request: Request, { params }: RouteContext) {
  const { id } = await params;
  if (!id) return NextResponse.json({ ok: false, message: "organisation_id required" }, { status: 400 });
  try {
    const hdrs = await resolvePortalHeaders(request);
    const result = await makeSdkClient(request, hdrs ?? undefined).deleteOrganisation({ path: { organisation_id: id } });
    if (!result.response.ok) {
      return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
    }
    return NextResponse.json({ ok: true, message: "Organisation deleted" });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}

export async function POST(request: Request, { params }: RouteContext) {
  const { id } = await params;
  if (!id) return NextResponse.json({ ok: false, message: "organisation_id required" }, { status: 400 });
  try {
    const hdrs = await resolvePortalHeaders(request);
    const http = makeDirectHttp(request, hdrs ?? undefined);
    const result = await http.patch(`/v1/b2b/organisations/${id}`, { status: "ORGANISATION_STATUS_ACTIVE" });
    if (!result.ok) return NextResponse.json({ ok: false, message: result.data?.error ?? "Approve failed" }, { status: result.status });
    const org = (result.data as Record<string, unknown>)?.organisation as Organisation | undefined;
    return NextResponse.json({ ok: true, message: "Organisation approved", organisation: org ? mapOrg(org) : null });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}
