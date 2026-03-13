import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, badRequest } from "@lib/sdk/api-helpers";

/** POST /api/auth/send-email-otp — send OTP to user's email */
export async function POST(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  let body: { email?: string; type?: string };
  try { body = await request.json(); } catch { return badRequest("Invalid request body"); }
  if (!body.email?.trim()) return badRequest("email is required");
  if (!body.type?.trim()) return badRequest("type is required (e.g. email_verification, password_reset_email)");
  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.sendEmailOtp({ body: { email: body.email, type: body.type } });
  if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
  return NextResponse.json({ ok: true, message: "Email OTP sent", data: result.data }, { status: 200 });
}
