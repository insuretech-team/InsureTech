import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, badRequest } from "@lib/sdk/api-helpers";
import { resolveUserIdFromSession } from "@lib/auth/resolve-user-id";

/** POST /api/auth/totp -- enable TOTP (2FA). Returns QR code / secret. */
export async function POST(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  const userId = await resolveUserIdFromSession(request, hdrs);
  if (!userId) return NextResponse.json({ ok: false, message: "Cannot resolve user identity" }, { status: 401 });
  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.enableTotp({ path: { user_id: userId }, body: { user_id: userId } });
  if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
  return NextResponse.json({ ok: true, totp: result.data }, { status: 200 });
}

/** DELETE /api/auth/totp -- disable TOTP */
export async function DELETE(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  let body: { totp_code?: string };
  try { body = await request.json(); } catch { return badRequest("Invalid request body"); }
  if (!body.totp_code?.trim()) return badRequest("totp_code is required");
  const userId = await resolveUserIdFromSession(request, hdrs);
  if (!userId) return NextResponse.json({ ok: false, message: "Cannot resolve user identity" }, { status: 401 });
  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.disableTotp({ path: { user_id: userId }, body: { user_id: userId, totp_code: body.totp_code } });
  if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
  return NextResponse.json({ ok: true, message: "TOTP disabled" }, { status: 200 });
}
