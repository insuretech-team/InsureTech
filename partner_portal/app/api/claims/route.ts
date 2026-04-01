import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/partner-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, badRequest } from "@lib/sdk/api-helpers";

/** GET /api/claims - List claims with filtering */
export async function GET(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

  const { searchParams } = new URL(request.url);
  const page = parseInt(searchParams.get("page") ?? "1", 10);
  const pageSize = parseInt(searchParams.get("page_size") ?? "20", 10);
  const partnerId = searchParams.get("partner_id") ?? undefined;
  const status = searchParams.get("status") ?? undefined;
  const startDate = searchParams.get("start_date") ?? undefined;
  const endDate = searchParams.get("end_date") ?? undefined;

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.listClaims({
    query: {
      page,
      page_size: pageSize,
      partner_id: partnerId,
      status,
      start_date: startDate,
      end_date: endDate,
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

/** POST /api/claims - Create a new claim */
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
  if (!body.policy_number || typeof body.policy_number !== "string") {
    return badRequest("policy_number is required");
  }
  if (!body.claim_amount || typeof body.claim_amount !== "number") {
    return badRequest("claim_amount is required and must be a number");
  }
  if (!body.claim_type || typeof body.claim_type !== "string") {
    return badRequest("claim_type is required");
  }

  const sdk = makeSdkClient(request, hdrs);
  
  // First validate policy coverage and eligibility
  const policyResult = await sdk.verifyPolicy({
    body: { policy_number: body.policy_number as string },
  });

  if (!policyResult.response.ok) {
    return NextResponse.json(
      { ok: false, message: "Policy verification failed: " + sdkErrorMessage(policyResult) },
      { status: policyResult.response.status }
    );
  }

  const policyData = policyResult.data as Record<string, unknown>;
  if (policyData.status !== "ACTIVE") {
    return NextResponse.json(
      { ok: false, message: "Policy is not active and cannot be used for claims" },
      { status: 400 }
    );
  }

  // Create the claim
  const result = await sdk.createClaim({ body });

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
