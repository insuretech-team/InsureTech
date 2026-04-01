import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/partner-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage } from "@lib/sdk/api-helpers";

/** POST /api/partners/[id]/credentials/rotate - Rotate partner API credentials */
export async function POST(
  request: Request,
  { params }: { params: { id: string } }
) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.rotatePartnerCredentials({ path: { partner_id: params.id } });

  if (!result.response.ok) {
    return NextResponse.json(
      { ok: false, message: sdkErrorMessage(result) },
      { status: result.response.status }
    );
  }

  // Log credential rotation event for audit
  // TODO: Implement audit logging when audit system is ready
  console.log(`[AUDIT] Partner credentials rotated: partner_id=${params.id}, user_id=${hdrs.userId}, timestamp=${new Date().toISOString()}`);

  return NextResponse.json(
    {
      ok: true,
      message: "API credentials rotated successfully. Old credentials are now invalid.",
      data: result.data,
    },
    { status: 200 }
  );
}
