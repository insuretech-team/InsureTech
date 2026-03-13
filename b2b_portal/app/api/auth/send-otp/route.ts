import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, badRequest } from "@lib/sdk/api-helpers";

/** POST /api/auth/send-otp — send OTP to a recipient (phone/email) */
export async function POST(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  let body: { recipient?: string; type?: string; channel?: string };
  try { body = await request.json(); } catch { return badRequest("Invalid request body"); }
  if (!body.recipient?.trim()) return badRequest("recipient is required");
  if (!body.type?.trim()) return badRequest("type is required (e.g. registration, login, reset_password)");
  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.sendOtp({ body: { recipient: body.recipient, type: body.type, channel: body.channel } });
  if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
  return NextResponse.json({ ok: true, message: "OTP sent", data: result.data }, { status: 200 });
}
