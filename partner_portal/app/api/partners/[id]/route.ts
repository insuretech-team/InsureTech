import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/partner-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, badRequest } from "@lib/sdk/api-helpers";

/** GET /api/partners/[id] - Get partner details */
export async function GET(
  request: Request,
  { params }: { params: { id: string } }
) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.getPartner({ path: { partner_id: params.id } });

  if (!result.response.ok) {
    return NextResponse.json(
      { ok: false, message: sdkErrorMessage(result) },
      { status: result.response.status }
    );
  }

  return NextResponse.json({ ok: true, data: result.data }, { status: 200 });
}

/** PATCH /api/partners/[id] - Update partner details */
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

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.updatePartner({
    path: { partner_id: params.id },
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

/** DELETE /api/partners/[id] - Soft delete partner */
export async function DELETE(
  request: Request,
  { params }: { params: { id: string } }
) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.deletePartner({ path: { partner_id: params.id } });

  if (!result.response.ok) {
    return NextResponse.json(
      { ok: false, message: sdkErrorMessage(result) },
      { status: result.response.status }
    );
  }

  return NextResponse.json({ ok: true, message: "Partner deleted successfully" }, { status: 200 });
}
