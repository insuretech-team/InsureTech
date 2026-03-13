import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, badRequest } from "@lib/sdk/api-helpers";

/** POST /api/auth/refresh — refresh the current session token */
export async function POST(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  let body: { refresh_token?: string; device_id?: string; ip_address?: string };
  try { body = await request.json(); } catch { body = {}; }
  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.refreshToken({
    body: {
      refresh_token: body.refresh_token,
      device_id: body.device_id ?? "web",
      ip_address: body.ip_address,
    },
  });
  if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
  return NextResponse.json({ ok: true, data: result.data }, { status: 200 });
}
