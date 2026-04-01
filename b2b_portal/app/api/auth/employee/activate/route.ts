import { NextResponse } from "next/server";

import { unwrapGateway, type GatewayResponse } from "@lib/sdk/shared";

function getApiBaseUrl(): string {
  return (
    process.env.INSURETECH_API_BASE_URL ??
    process.env.NEXT_PUBLIC_INSURETECH_API_BASE_URL ??
    "http://localhost:8080"
  );
}

export async function POST(request: Request) {
  let payload: { organisationId?: string; organisationCode?: string; employeeId?: string; email?: string };
  try {
    payload = (await request.json()) as { organisationId?: string; organisationCode?: string; employeeId?: string; email?: string };
  } catch {
    return NextResponse.json({ ok: false, message: "Invalid request body" }, { status: 400 });
  }

  const organisationId = payload.organisationId?.trim();
  const organisationCode = payload.organisationCode?.trim().toUpperCase();
  const employeeId = payload.employeeId?.trim();
  const email = payload.email?.trim().toLowerCase();

  if (!organisationCode && !organisationId) {
    return NextResponse.json({ ok: false, message: "Organisation selection is required" }, { status: 400 });
  }
  if (!employeeId) {
    return NextResponse.json({ ok: false, message: "Employee ID is required" }, { status: 400 });
  }
  if (!email) {
    return NextResponse.json({ ok: false, message: "Email is required" }, { status: 400 });
  }

  try {
    const response = await fetch(`${getApiBaseUrl()}/v1/b2b/employees:activate`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        organisation_id: organisationId,
        organisation_code: organisationCode,
        employee_id: employeeId,
        email,
      }),
      cache: "no-store",
    });
    const body = (await response.json().catch(() => null)) as GatewayResponse<{
      otp_id?: string;
      expires_in_seconds?: number;
      user_id?: string;
      message?: string;
    }> | null;

    if (!body) {
      return NextResponse.json({ ok: false, message: "Activation failed" }, { status: 502 });
    }

    const unwrapped = unwrapGateway(body, response.status);
    if (!unwrapped.ok) {
      return NextResponse.json({ ok: false, message: unwrapped.message }, { status: unwrapped.status });
    }

    return NextResponse.json({
      ok: true,
      message: unwrapped.data.message ?? "Verification code sent to your email.",
      otpId: unwrapped.data.otp_id ?? "",
      expiresInSeconds: unwrapped.data.expires_in_seconds ?? 0,
      userId: unwrapped.data.user_id ?? "",
    });
  } catch (error) {
    return NextResponse.json(
      { ok: false, message: error instanceof Error ? error.message : "Activation failed" },
      { status: 502 }
    );
  }
}
