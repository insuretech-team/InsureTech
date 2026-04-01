import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/partner-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, badRequest } from "@lib/sdk/api-helpers";

/** GET /api/partners/[id]/agents - List partner agents with pagination */
export async function GET(
  request: Request,
  { params }: { params: { id: string } }
) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

  const { searchParams } = new URL(request.url);
  const page = parseInt(searchParams.get("page") ?? "1", 10);
  const pageSize = parseInt(searchParams.get("page_size") ?? "20", 10);
  const status = searchParams.get("status") ?? undefined;

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.listPartnerAgents({
    path: { partner_id: params.id },
    query: {
      page,
      page_size: pageSize,
      status,
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

/** POST /api/partners/[id]/agents - Create a new partner agent */
export async function POST(
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

  // Validate required fields
  if (!body.full_name || typeof body.full_name !== "string") {
    return badRequest("full_name is required");
  }
  if (!body.mobile_number || typeof body.mobile_number !== "string") {
    return badRequest("mobile_number is required");
  }
  if (!body.email || typeof body.email !== "string") {
    return badRequest("email is required");
  }

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.createPartnerAgent({
    path: { partner_id: params.id },
    body,
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
      message: "Agent created successfully. Verification credentials sent via email/SMS.",
      data: result.data,
    },
    { status: result.response.status || 201 }
  );
}
