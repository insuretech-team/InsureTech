import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/partner-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, badRequest } from "@lib/sdk/api-helpers";

const VALID_STATUS_TRANSITIONS: Record<string, string[]> = {
  SUBMITTED: ["UNDER_REVIEW", "REJECTED"],
  UNDER_REVIEW: ["APPROVED", "REJECTED"],
  APPROVED: ["PAID"],
  REJECTED: [],
  PAID: [],
};

/** GET /api/claims/[id] - Get claim details */
export async function GET(
  request: Request,
  { params }: { params: { id: string } }
) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.getClaim({ path: { claim_id: params.id } });

  if (!result.response.ok) {
    return NextResponse.json(
      { ok: false, message: sdkErrorMessage(result) },
      { status: result.response.status }
    );
  }

  return NextResponse.json({ ok: true, data: result.data }, { status: 200 });
}

/** PATCH /api/claims/[id] - Update claim details */
export async function PATCH(
  request: Request,
  { params }: { params: { id: string } }
) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

  let body: Record<string, unknown>;
  try {
    body = await request.json();
  } catch {
    return badRequest("Invalid request body");
  }

  // Validate status transitions if status is being updated
  if (body.status && typeof body.status === "string") {
    const sdk = makeSdkClient(request, hdrs);
    const currentResult = await sdk.getClaim({ path: { claim_id: params.id } });
    
    if (currentResult.response.ok) {
      const currentData = currentResult.data as Record<string, unknown>;
      const currentStatus = currentData.status as string;
      const newStatus = body.status;

      const allowedTransitions = VALID_STATUS_TRANSITIONS[currentStatus] ?? [];
      if (!allowedTransitions.includes(newStatus)) {
        return NextResponse.json(
          {
            ok: false,
            message: `Invalid status transition from ${currentStatus} to ${newStatus}. Allowed: ${allowedTransitions.join(", ") || "none"}`,
          },
          { status: 400 }
        );
      }
    }
  }

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.updateClaim({
    path: { claim_id: params.id },
    body,
  });

  if (!result.response.ok) {
    return NextResponse.json(
      { ok: false, message: sdkErrorMessage(result) },
      { status: result.response.status }
    );
  }

  return NextResponse.json({ ok: true, data: result.data }, { status: 200 });
}
