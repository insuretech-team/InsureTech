import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/partner-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, badRequest } from "@lib/sdk/api-helpers";

/** GET /api/partners/[id]/agents/[agentId] - Get agent details */
export async function GET(
  request: Request,
  { params }: { params: { id: string; agentId: string } }
) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.getPartnerAgent({
    path: { partner_id: params.id, agent_id: params.agentId },
  });

  if (!result.response.ok) {
    return NextResponse.json(
      { ok: false, message: sdkErrorMessage(result) },
      { status: result.response.status }
    );
  }

  return NextResponse.json({ ok: true, data: result.data }, { status: 200 });
}

/** PATCH /api/partners/[id]/agents/[agentId] - Update agent details */
export async function PATCH(
  request: Request,
  { params }: { params: { id: string; agentId: string } }
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
  const result = await sdk.updatePartnerAgent({
    path: { partner_id: params.id, agent_id: params.agentId },
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

/** DELETE /api/partners/[id]/agents/[agentId] - Deactivate agent */
export async function DELETE(
  request: Request,
  { params }: { params: { id: string; agentId: string } }
) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.deletePartnerAgent({
    path: { partner_id: params.id, agent_id: params.agentId },
  });

  if (!result.response.ok) {
    return NextResponse.json(
      { ok: false, message: sdkErrorMessage(result) },
      { status: result.response.status }
    );
  }

  return NextResponse.json(
    { ok: true, message: "Agent deactivated successfully" },
    { status: 200 }
  );
}
