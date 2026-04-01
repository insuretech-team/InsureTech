import { NextResponse } from "next/server";
import { makeSdkClient, makeDirectHttp } from "@lib/sdk/partner-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, badRequest } from "@lib/sdk/api-helpers";

const VALID_COMMISSION_TYPES = ["PERCENTAGE", "FIXED", "TIERED"] as const;

/** GET /api/partners/[id]/commission - Get partner commission structure */
export async function GET(
  request: Request,
  { params }: { params: { id: string } }
) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.getPartnerCommission({ path: { partner_id: params.id } });

  if (!result.response.ok) {
    return NextResponse.json(
      { ok: false, message: sdkErrorMessage(result) },
      { status: result.response.status }
    );
  }

  return NextResponse.json({ ok: true, data: result.data }, { status: 200 });
}

/** PATCH /api/partners/[id]/commission - Update partner commission structure */
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

  // Validate commission type
  if (body.commission_type && !VALID_COMMISSION_TYPES.includes(body.commission_type as typeof VALID_COMMISSION_TYPES[number])) {
    return badRequest(
      `Invalid commission_type. Must be one of: ${VALID_COMMISSION_TYPES.join(", ")}`
    );
  }

  // Use makeDirectHttp since updatePartnerCommission might not be in SDK yet
  const http = makeDirectHttp(request);
  const result = await http.patch(`/v1/partners/${params.id}/commission`, body);

  if (!result.ok) {
    return NextResponse.json(
      { ok: false, message: result.error || "Failed to update commission structure" },
      { status: result.status }
    );
  }

  return NextResponse.json(
    { ok: true, message: "Commission structure updated successfully", data: result.data },
    { status: 200 }
  );
}
