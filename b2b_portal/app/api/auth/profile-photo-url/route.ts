import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage } from "@lib/sdk/api-helpers";
import { resolveUserIdFromSession } from "@lib/auth/resolve-user-id";

/** GET /api/auth/profile-photo-url -- get presigned upload URL for profile photo */
export async function GET(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  const userId = await resolveUserIdFromSession(request, hdrs);
  if (!userId) return NextResponse.json({ ok: false, message: "Cannot resolve user identity" }, { status: 401 });
  const { searchParams } = new URL(request.url);
  const contentType = searchParams.get("content_type") ?? "image/jpeg";
  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.getProfilePhotoUploadUrl({
    path: { user_id: userId },
    body: { user_id: userId, content_type: contentType },
  });
  if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
  return NextResponse.json({ ok: true, uploadUrl: result.data }, { status: 200 });
}
