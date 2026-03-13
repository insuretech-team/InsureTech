import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, badRequest } from "@lib/sdk/api-helpers";
import { resolveUserIdFromSession } from "@lib/auth/resolve-user-id";

/** POST /api/auth/change-password */
export async function POST(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  let body: { old_password?: string; new_password?: string };
  try { body = await request.json(); } catch { return badRequest("Invalid request body"); }
  if (!body.old_password?.trim()) return badRequest("old_password is required");
  if (!body.new_password?.trim()) return badRequest("new_password is required");
  const userId = await resolveUserIdFromSession(request, hdrs);
  if (!userId) return NextResponse.json({ ok: false, message: "Cannot resolve user identity" }, { status: 401 });
  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.changePassword({ body: { user_id: userId, old_password: body.old_password, new_password: body.new_password } });
  if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
  return NextResponse.json({ ok: true, message: "Password changed successfully" }, { status: 200 });
}
