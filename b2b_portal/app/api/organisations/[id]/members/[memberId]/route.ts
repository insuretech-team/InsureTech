import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { sdkErrorMessage } from "@lib/sdk/api-helpers";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";

type RouteContext = { params: Promise<{ id: string; memberId: string }> };

export async function DELETE(request: Request, { params }: RouteContext) {
  const { id, memberId } = await params;
  if (!id || !memberId) {
    return NextResponse.json({ ok: false, message: "organisation_id and member_id are required" }, { status: 400 });
  }

  try {
    const hdrs = await resolvePortalHeaders(request);
    const result = await makeSdkClient(request, hdrs ?? undefined).removeOrgMember({
      path: {
        organisation_id: id,
        member_id: memberId,
      },
    });

    if (!result.response.ok) {
      return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
    }

    return NextResponse.json({
      ok: true,
      message: "Member removed",
    });
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}
