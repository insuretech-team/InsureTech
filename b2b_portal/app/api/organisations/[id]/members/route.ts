import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { sdkErrorMessage } from "@lib/sdk/api-helpers";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import type { OrgMember, OrgMemberRole } from "@lifeplus/insuretech-sdk";

type RouteContext = { params: Promise<{ id: string }> };

export async function GET(request: Request, { params }: RouteContext) {
  const { id } = await params;
  if (!id) {
    return NextResponse.json({ ok: false, message: "organisation_id required", members: [] }, { status: 400 });
  }

  try {
    const hdrs = await resolvePortalHeaders(request);
    const result = await makeSdkClient(request, hdrs ?? undefined).listOrgMembers({
      path: { organisation_id: id },
    });


    if (!result.response.ok) {
      return NextResponse.json(
        { ok: false, message: sdkErrorMessage(result), members: [] },
        { status: result.response.status }
      );
    }

    return NextResponse.json({
      ok: true,
      members: (result.data?.members ?? []) as OrgMember[],
    });
  } catch (err) {
    return NextResponse.json(
      { ok: false, message: err instanceof Error ? err.message : "Error", members: [] },
      { status: 502 }
    );
  }
}

export async function POST(request: Request, { params }: RouteContext) {
  const { id } = await params;
  if (!id) {
    return NextResponse.json({ ok: false, message: "organisation_id required" }, { status: 400 });
  }
  try {
    const body = (await request.json()) as Record<string, unknown>;
    const userId = String(body.userId ?? "").trim();
    const role = String(body.role ?? "ORG_MEMBER_ROLE_HR_MANAGER") as OrgMemberRole;
    if (!userId) {
      return NextResponse.json({ ok: false, message: "userId is required" }, { status: 400 });
    }
    const hdrs = await resolvePortalHeaders(request);
    const result = await makeSdkClient(request, hdrs ?? undefined).addOrgMember({
      path: { organisation_id: id },
      body: { organisation_id: id, user_id: userId, role },
    });
    if (!result.response.ok) {
      return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
    }
    return NextResponse.json({ ok: true, message: result.data?.message ?? "Member added", member: result.data?.member });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}
