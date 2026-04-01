import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/partner-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage } from "@lib/sdk/api-helpers";

/**
 * Masks an API key, showing only the last 4 characters
 * Example: "sk_live_1234567890abcdef" -> "sk_live_************cdef"
 */
function maskApiKey(key: string): string {
  if (!key || key.length <= 4) return key;
  const visiblePart = key.slice(-4);
  const prefix = key.includes("_") ? key.split("_").slice(0, 2).join("_") + "_" : "";
  const maskedLength = key.length - prefix.length - 4;
  return prefix + "*".repeat(maskedLength) + visiblePart;
}

/** GET /api/partners/[id]/credentials - Get partner API credentials (masked) */
export async function GET(
  request: Request,
  { params }: { params: { id: string } }
) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.getPartnerCredentials({ path: { partner_id: params.id } });

  if (!result.response.ok) {
    return NextResponse.json(
      { ok: false, message: sdkErrorMessage(result) },
      { status: result.response.status }
    );
  }

  // Mask API keys in the response
  const data = result.data as Record<string, unknown>;
  if (data && typeof data === "object") {
    if (typeof data.api_key === "string") {
      data.api_key = maskApiKey(data.api_key);
    }
    if (typeof data.api_secret === "string") {
      data.api_secret = maskApiKey(data.api_secret);
    }
  }

  return NextResponse.json({ ok: true, data }, { status: 200 });
}
