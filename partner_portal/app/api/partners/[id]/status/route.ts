import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/partner-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, badRequest } from "@lib/sdk/api-helpers";

const VALID_STATUSES = [
  "PENDING_VERIFICATION",
  "ACTIVE",
  "INACTIVE",
  "SUSPENDED",
  "REJECTED",
] as const;

/** PATCH /api/partners/[id]/status - Update partner status */
export async function PATCH(
  request: Request,
  { params }: { params: { id: string } }
) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

  let body: { status?: string; reason?: string };
  try {
    body = await request.json();
  } catch {
    return badRequest("Invalid request body");
  }

  if (!body.status || typeof body.status !== "string") {
    return badRequest("status field is required");
  }

  if (!VALID_STATUSES.includes(body.status as typeof VALID_STATUSES[number])) {
    return badRequest(
      `Invalid status. Must be one of: ${VALID_STATUSES.join(", ")}`
    );
  }

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.updatePartnerStatus({
    path: { partner_id: params.id },
    body: {
      status: body.status,
      reason: body.reason,
    },
  });

  if (!result.response.ok) {
    return NextResponse.json(
      { ok: false, message: sdkErrorMessage(result) },
      { status: result.response.status }
    );
  }

  return NextResponse.json(
    { ok: true, message: "Partner status updated successfully", data: result.data },
    { status: 200 }
  );
}
