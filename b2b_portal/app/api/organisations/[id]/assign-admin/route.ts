/**
 * POST /api/organisations/[id]/assign-admin
 * Assigns an existing user as org admin by calling b2bServiceAssignOrgAdmin.
 * Body: { userId: string }
 */
import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { sdkErrorMessage } from "@lib/sdk/api-helpers";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";

type RouteContext = { params: Promise<{ id: string }> };

export async function POST(request: Request, { params }: RouteContext) {
  const { id } = await params;
  if (!id) {
    return NextResponse.json({ ok: false, message: "organisation_id required" }, { status: 400 });
  }
  try {
    const body = (await request.json()) as Record<string, unknown>;
    // The SDK uses member_id (not user_id) for assignOrgAdmin.
    // The caller may pass either userId or memberId — accept both.
    const memberId = String(body.memberId ?? body.userId ?? "").trim();
    if (!memberId) {
      return NextResponse.json({ ok: false, message: "memberId (or userId) is required" }, { status: 400 });
    }
    const hdrs = await resolvePortalHeaders(request);
    const result = await makeSdkClient(request, hdrs ?? undefined).assignOrgAdmin({
      path: { organisation_id: id },
      body: { organisation_id: id, member_id: memberId },
    });
    if (!result.response.ok) {
      return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
    }
    return NextResponse.json({ ok: true, message: result.data?.message ?? "Admin assigned successfully" });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}
