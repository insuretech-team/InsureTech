import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, badRequest } from "@lib/sdk/api-helpers";

/** POST /api/auth/verify-otp — verify OTP code */
export async function POST(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  let body: { otp_id?: string; code?: string };
  try { body = await request.json(); } catch { return badRequest("Invalid request body"); }
  if (!body.otp_id?.trim()) return badRequest("otp_id is required");
  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.verifyOtp({ body: { otp_id: body.otp_id, code: body.code } });
  if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
  return NextResponse.json({ ok: true, message: "OTP verified", data: result.data }, { status: 200 });
}
