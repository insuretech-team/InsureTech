import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/partner-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, badRequest } from "@lib/sdk/api-helpers";

/** POST /api/policies/verify - Verify policy eligibility */
export async function POST(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

  let body: { policy_number?: string; customer_nid?: string };
  try {
    body = await request.json();
  } catch {
    return badRequest("Invalid request body");
  }

  // At least one identifier is required
  if (!body.policy_number && !body.customer_nid) {
    return badRequest("Either policy_number or customer_nid is required");
  }

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.verifyPolicy({ body });

  if (!result.response.ok) {
    return NextResponse.json(
      { ok: false, message: sdkErrorMessage(result) },
      { status: result.response.status }
    );
  }

  // Log verification request for audit
  console.log(
    `[AUDIT] Policy verification: ` +
    `policy_number=${body.policy_number ?? "N/A"}, ` +
    `customer_nid=${body.customer_nid ?? "N/A"}, ` +
    `user_id=${hdrs.userId}, ` +
    `timestamp=${new Date().toISOString()}`
  );

  const data = result.data as Record<string, unknown>;
  
  // Ensure response includes all required fields
  const response = {
    ok: true,
    data: {
      policy_number: data.policy_number ?? body.policy_number,
      status: data.status ?? "UNKNOWN",
      coverage_details: data.coverage_details ?? {},
      sum_assured: data.sum_assured ?? 0,
      eligibility_status: data.eligibility_status ?? "UNKNOWN",
      customer_name: data.customer_name,
      policy_start_date: data.policy_start_date,
      policy_end_date: data.policy_end_date,
      ...data,
    },
  };

  return NextResponse.json(response, { status: 200 });
}
