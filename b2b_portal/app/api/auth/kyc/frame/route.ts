import { NextResponse } from "next/server";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";

/**
 * POST /api/auth/kyc/frame
 *
 * Submits a KYC frame (JPEG image) to the InsureTech gateway for processing.
 * This is a BFF proxy that forwards multipart form data to the gateway.
 *
 * Request body (multipart/form-data):
 *   - user_id: string (user identifier)
 *   - session_id: string (KYC session identifier)
 *   - image_data: binary (JPEG image file)
 *
 * Response: Gateway response JSON { ok: true, ... } | { ok: false, message: string }
 */
export async function POST(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) {
    return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  }

  // Parse multipart form data
  let formData: FormData;
  try {
    formData = await request.formData();
  } catch {
    return NextResponse.json({ ok: false, message: "Invalid multipart form data" }, { status: 400 });
  }

  // Extract required fields
  const userId = formData.get("user_id");
  const sessionId = formData.get("session_id");
  // Accept "file" (sent by EKYCFlow.tsx) or "image_data" (legacy)
  const imageData = formData.get("file") ?? formData.get("image_data");

  if (!userId || typeof userId !== "string") {
    return NextResponse.json({ ok: false, message: "user_id is required" }, { status: 400 });
  }

  if (!sessionId || typeof sessionId !== "string") {
    return NextResponse.json({ ok: false, message: "session_id is required" }, { status: 400 });
  }

  if (!imageData || !(imageData instanceof File)) {
    return NextResponse.json({ ok: false, message: "image file is required" }, { status: 400 });
  }

  // Build the gateway URL
  const gatewayUrl = process.env.INSURETECH_GATEWAY_URL ?? "http://localhost:8080";

  // Create a new FormData to forward to the gateway
  // Gateway expects field name "image_data" for the binary frame
  const forwardFormData = new FormData();
  forwardFormData.append("user_id", userId);
  forwardFormData.append("session_id", sessionId);
  forwardFormData.append("image_data", imageData);

  // Forward to the gateway with session cookie and portal headers
  const resp = await fetch(`${gatewayUrl}/v1/auth/users/${userId}/kyc:submit-frame`, {
    method: "POST",
    headers: {
      "cookie": request.headers.get("cookie") ?? "",
      "x-portal": hdrs.portal ?? "b2b",
      "x-user-id": hdrs.userId ?? "",
      "x-business-id": hdrs.businessId ?? "",
      "x-tenant-id": hdrs.tenantId ?? "",
    },
    body: forwardFormData,
  });

  // Return the gateway response
  if (!resp.ok) {
    let errorMessage = "Gateway error";
    try {
      const errorData = await resp.json() as Record<string, unknown>;
      errorMessage = (errorData.message ?? errorData.error ?? errorMessage) as string;
    } catch {
      // If response is not JSON, use status text
      errorMessage = resp.statusText || "Gateway error";
    }
    return NextResponse.json({ ok: false, message: errorMessage }, { status: resp.status });
  }

  // Gateway returns SubmitKYCFrameResponse with UseProtoNames (snake_case):
  // detection, current_step_detail, liveness_confidence, guidance_messages,
  // overall_progress, session_state — exactly what EKYCFlow.tsx expects.
  // Unwrap ApiResponse envelope if present.
  const raw = await resp.json() as Record<string, unknown>;
  const inner = (raw.success !== undefined && raw.data !== null && typeof raw.data === "object")
    ? raw.data as Record<string, unknown>
    : raw;

  return NextResponse.json(inner, { status: 200 });
}
