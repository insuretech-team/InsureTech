import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/partner-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, badRequest } from "@lib/sdk/api-helpers";

/** GET /api/partners - List partners with pagination */
export async function GET(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

  const { searchParams } = new URL(request.url);
  const page = parseInt(searchParams.get("page") ?? "1", 10);
  const pageSize = parseInt(searchParams.get("page_size") ?? "20", 10);
  const status = searchParams.get("status") ?? undefined;
  const partnerType = searchParams.get("partner_type") ?? undefined;

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.listPartners({
    query: {
      page,
      page_size: pageSize,
      status,
      partner_type: partnerType,
    },
  });

  if (!result.response.ok) {
    return NextResponse.json(
      { ok: false, message: sdkErrorMessage(result) },
      { status: result.response.status }
    );
  }

  return NextResponse.json({ ok: true, data: result.data }, { status: 200 });
}

/** POST /api/partners - Create a new partner */
export async function POST(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

  let body: Record<string, unknown>;
  try {
    body = await request.json();
  } catch {
    return badRequest("Invalid request body");
  }

  // Validate required fields
  if (!body.partner_name || typeof body.partner_name !== "string") {
    return badRequest("partner_name is required");
  }
  if (!body.partner_type || typeof body.partner_type !== "string") {
    return badRequest("partner_type is required");
  }

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.createPartner({ body });

  if (!result.response.ok) {
    return NextResponse.json(
      { ok: false, message: sdkErrorMessage(result) },
      { status: result.response.status }
    );
  }

  return NextResponse.json(
    { ok: true, data: result.data },
    { status: result.response.status || 201 }
  );
}
