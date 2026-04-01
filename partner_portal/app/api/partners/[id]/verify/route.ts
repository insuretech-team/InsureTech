import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/partner-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, badRequest } from "@lib/sdk/api-helpers";

/** POST /api/partners/[id]/verify - Verify partner (approve/reject) */
export async function POST(
  request: Request,
  { params }: { params: { id: string } }
) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

  let body: { approved?: boolean; rejection_reason?: string };
  try {
    body = await request.json();
  } catch {
    return badRequest("Invalid request body");
  }

  if (typeof body.approved !== "boolean") {
    return badRequest("approved field is required and must be a boolean");
  }

  if (!body.approved && !body.rejection_reason?.trim()) {
    return badRequest("rejection_reason is required when rejecting a partner");
  }

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.verifyPartner({
    path: { partner_id: params.id },
    body: {
      approved: body.approved,
      rejection_reason: body.rejection_reason,
    },
  });

  if (!result.response.ok) {
    return NextResponse.json(
      { ok: false, message: sdkErrorMessage(result) },
      { status: result.response.status }
    );
  }

  return NextResponse.json(
    {
      ok: true,
      message: body.approved ? "Partner approved successfully" : "Partner rejected",
      data: result.data,
    },
    { status: 200 }
  );
}
