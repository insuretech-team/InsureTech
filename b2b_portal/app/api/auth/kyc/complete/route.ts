import { NextResponse } from "next/server";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { resolveUserIdFromSession } from "@lib/auth/resolve-user-id";

/**
 * POST /api/auth/kyc/complete
 *
 * Called by EKYCFlow after FLVE completes the session.
 * Calls InsureTech CompleteKYCSession → sets status to PENDING_REVIEW.
 *
 * Body: { session_id: string, profile_image_url?: string }
 * Response: { ok: true } | { ok: false, message: string }
 */
export async function POST(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) {
    return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  }

  let body: { session_id?: string; profile_image_url?: string };
  try {
    body = await request.json() as { session_id?: string; profile_image_url?: string };
  } catch {
    return NextResponse.json({ ok: false, message: "Invalid request body" }, { status: 400 });
  }

  if (!body.session_id) {
    return NextResponse.json({ ok: false, message: "session_id is required" }, { status: 400 });
  }

  const userId = await resolveUserIdFromSession(request, hdrs);
  if (!userId) {
    return NextResponse.json({ ok: false, message: "Cannot resolve user identity" }, { status: 401 });
  }

  // Call InsureTech gateway CompleteKYCSession directly (not in SDK yet)
  const gatewayUrl = process.env.INSURETECH_GATEWAY_URL ?? "http://localhost:8080";
  const resp = await fetch(`${gatewayUrl}/v1/auth/users/${userId}/kyc:complete`, {
    method:  "POST",
    headers: {
      "Content-Type":  "application/json",
      "cookie":        request.headers.get("cookie") ?? "",
      "x-portal":      hdrs.portal ?? "b2b",
      "x-business-id": hdrs.businessId ?? "",
    },
    body: JSON.stringify({
      session_id:        body.session_id,
      profile_image_url: body.profile_image_url ?? "",
    }),
  });

  if (!resp.ok) {
    const errData = await resp.json().catch(() => ({})) as { message?: string };
    return NextResponse.json(
      { ok: false, message: errData.message ?? `Gateway error ${resp.status}` },
      { status: resp.status }
    );
  }

  // Pass through useful fields from the gateway response to the client
  // (liveness_confidence, captured_image_base64, profile_image_url) for the
  // done-screen score grid and image preview — matches Svelte client behaviour.
  const gwData = await resp.json().catch(() => ({})) as Record<string, unknown>;
  const payload = {
    ok: true,
    profile_image_url:    gwData.profile_image_url    ?? gwData.profileImageUrl    ?? "",
    liveness_confidence:  gwData.liveness_confidence  ?? gwData.livenessConfidence ?? 0,
    captured_image_base64: gwData.captured_image_base64 ?? gwData.capturedImageBase64 ?? "",
  };

  const response = NextResponse.json(payload, { status: 200 });
  response.cookies.set({
    name: "portal_kyc_verified",
    value: "pending_review",
    path: "/",
    httpOnly: false,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    maxAge: 60 * 60 * 12,
  });
  return response;
}
