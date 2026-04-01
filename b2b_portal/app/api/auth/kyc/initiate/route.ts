import { NextResponse } from "next/server";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";

/**
 * POST /api/auth/kyc/initiate
 *
 * Creates a kyc_verifications record in InsureTech DB before the FLVE session
 * starts. This gives us the canonical InsureTech KYC UUID (kyc_verification_id)
 * that FLVE stores in provider_reference.
 *
 * Body: { user_id: string }
 * Response: { ok: true, kyc_verification_id: string } | { ok: false, message: string }
 */
export async function POST(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) {
    return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  }

  let body: { user_id?: string };
  try {
    body = await request.json() as { user_id?: string };
  } catch {
    return NextResponse.json({ ok: false, message: "Invalid request body" }, { status: 400 });
  }

  const userId = body.user_id;
  if (!userId) {
    return NextResponse.json({ ok: false, message: "user_id is required" }, { status: 400 });
  }

  // Call InsureTech gateway InitiateKYC directly (not in SDK yet)
  const gatewayUrl = process.env.INSURETECH_GATEWAY_URL ?? "http://localhost:8080";
  const resp = await fetch(`${gatewayUrl}/v1/auth/users/${userId}/kyc`, {
    method:  "POST",
    headers: {
      "Content-Type":  "application/json",
      "cookie":        request.headers.get("cookie") ?? "",
      "x-portal":      hdrs.portal ?? "b2b",
      "x-business-id": hdrs.businessId ?? "",
    },
    body: JSON.stringify({ verification_method: "FLVE_EKYC" }),
  });

  if (!resp.ok) {
    // Non-fatal — eKYC will still work; FLVE session just won't have a pre-created record
    return NextResponse.json({ ok: true, kyc_verification_id: null }, { status: 200 });
  }

  const data = await resp.json() as Record<string, unknown>;
  // Gateway wraps InitiateKYCResponse in ApiResponse envelope: { success, data, ... }
  // Unwrap if needed
  const inner = (data.success !== undefined && data.data !== null && typeof data.data === "object")
    ? data.data as Record<string, unknown>
    : data;

  // Proto JSON marshaling uses camelCase: provider_reference → providerReference
  const kycId = (inner.kyc_id ?? inner.kycId ?? inner.kyc_verification_id ?? null) as string | null;
  // provider_reference (snake) or providerReference (camel) — FLVE session_id
  const sessionId = (inner.provider_reference ?? inner.providerReference ?? inner.session_id ?? null) as string | null;
  const steps = (inner.steps ?? []) as unknown[];
  const sessionState = (inner.session_state ?? inner.sessionState ?? inner.status ?? "") as string;

  return NextResponse.json({
    ok: true,
    kyc_verification_id: kycId,
    session_id: sessionId,
    steps,
    session_state: sessionState,
  }, { status: 200 });
}
