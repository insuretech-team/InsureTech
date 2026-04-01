import { NextResponse } from "next/server";

import { unwrapGateway, type GatewayResponse } from "@lib/sdk/shared";
import { getPasswordValidationMessage } from "@/src/lib/utils/password";

function getApiBaseUrl(): string {
  return (
    process.env.INSURETECH_API_BASE_URL ??
    process.env.NEXT_PUBLIC_INSURETECH_API_BASE_URL ??
    "http://localhost:8080"
  );
}

export async function POST(request: Request) {
  let payload: { email?: string; otpId?: string; otpCode?: string; newPassword?: string };
  try {
    payload = (await request.json()) as { email?: string; otpId?: string; otpCode?: string; newPassword?: string };
  } catch {
    return NextResponse.json({ ok: false, message: "Invalid request body" }, { status: 400 });
  }

  const email = payload.email?.trim().toLowerCase();
  const otpId = payload.otpId?.trim();
  const otpCode = payload.otpCode?.trim();
  const newPassword = payload.newPassword ?? "";

  if (!email) {
    return NextResponse.json({ ok: false, message: "Email is required" }, { status: 400 });
  }
  if (!otpId) {
    return NextResponse.json({ ok: false, message: "OTP session is required" }, { status: 400 });
  }
  if (!otpCode) {
    return NextResponse.json({ ok: false, message: "OTP code is required" }, { status: 400 });
  }

  const passwordError = getPasswordValidationMessage(newPassword, "New password");
  if (passwordError) {
    return NextResponse.json({ ok: false, message: passwordError }, { status: 400 });
  }

  try {
    const response = await fetch(`${getApiBaseUrl()}/v1/auth/email/password:reset`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        email,
        otp_id: otpId,
        otp_code: otpCode,
        new_password: newPassword,
      }),
      cache: "no-store",
    });
    const body = (await response.json().catch(() => null)) as GatewayResponse<{
      message?: string;
    }> | null;

    if (!body) {
      return NextResponse.json({ ok: false, message: "Unable to set password" }, { status: 502 });
    }

    const unwrapped = unwrapGateway(body, response.status);
    if (!unwrapped.ok) {
      return NextResponse.json({ ok: false, message: unwrapped.message }, { status: unwrapped.status });
    }

    return NextResponse.json({
      ok: true,
      message: unwrapped.data.message ?? "Password set successfully. You can now sign in.",
    });
  } catch (error) {
    return NextResponse.json(
      { ok: false, message: error instanceof Error ? error.message : "Unable to set password" },
      { status: 502 }
    );
  }
}
