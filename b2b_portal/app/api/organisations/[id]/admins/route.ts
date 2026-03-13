/**
 * POST /api/organisations/[id]/admins
 *
 * Two modes — determined by the request body:
 *
 * Mode A — Create new B2B admin user (email + password + mobile):
 *   body: { email, password, mobileNumber, fullName? }
 *   → forwards to backend POST /v1/b2b/organisations/{id}/admins
 *     which calls AuthN.RegisterEmailUser + AssignRole(B2B_ORG_ADMIN)
 *
 * Mode B — Promote an existing org member to admin (memberId only):
 *   body: { memberId }
 *   → forwards to backend PUT /v1/b2b/organisations/{id}/assign-admin
 *     via the SDK assignOrgAdmin call (member_id based)
 *
 * This split prevents accidentally registering a new user when the caller
 * only intends to promote an existing member, and vice-versa.
 */
import { NextResponse } from "next/server";
import { makeDirectHttp, makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { sdkErrorMessage } from "@lib/sdk/api-helpers";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";

type RouteContext = { params: Promise<{ id: string }> };

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
    const body = (await request.json()) as Record<string, unknown>;

    // Mode B: promote existing member — body has memberId but no email/password
    const memberId = typeof body.memberId === "string" ? body.memberId.trim() : "";
    if (memberId && !body.email && !body.password) {
      const result = await makeSdkClient(request, hdrs).assignOrgAdmin({
        path: { organisation_id: id },
        body: { organisation_id: id, member_id: memberId },
      });
      if (!result.response.ok) {
        return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
      }
      return NextResponse.json({ ok: true, message: result.data?.message ?? "Admin assigned successfully" });
    }

    // Mode A: create new admin user — requires email + password + mobileNumber
    const email = typeof body.email === "string" ? body.email.trim() : "";
    const password = typeof body.password === "string" ? body.password : "";
    const mobileNumber = typeof body.mobileNumber === "string" ? body.mobileNumber.trim() : "";

    if (!email || !password || !mobileNumber) {
      return NextResponse.json(
        { ok: false, message: "email, password, and mobileNumber are required to create a new admin" },
        { status: 400 }
      );
    }

    // Backend assignOrgAdminPayload struct uses camelCase JSON tags:
    // json:"email", json:"password", json:"mobileNumber", json:"fullName"
    const result = await makeDirectHttp(request, hdrs).post(`/v1/b2b/organisations/${id}/admins`, {
      email,
      password,
      mobileNumber,
      fullName: typeof body.fullName === "string" ? body.fullName.trim() : undefined,
    });

    return NextResponse.json(
      { ok: result.ok, message: result.data?.message ?? (result.ok ? "Admin created" : "Failed to create admin") },
      { status: result.ok ? 201 : result.status }
    );
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}
