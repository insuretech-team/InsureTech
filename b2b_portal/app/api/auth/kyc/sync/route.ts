import { NextResponse } from "next/server";

const ALLOWED_STATUSES = new Set(["true", "false", "pending_review"]);

function sanitizeNextPath(value: string | null): string {
  if (!value || !value.startsWith("/")) {
    return "/";
  }
  if (value.startsWith("//")) {
    return "/";
  }
  return value;
}

export async function GET(request: Request) {
  const url = new URL(request.url);
  const status = url.searchParams.get("status") ?? "";
  const nextPath = sanitizeNextPath(url.searchParams.get("next"));

  if (!ALLOWED_STATUSES.has(status)) {
    return NextResponse.json({ ok: false, message: "Invalid KYC status" }, { status: 400 });
  }

  const response = NextResponse.redirect(new URL(nextPath, request.url));
  response.cookies.set({
    name: "portal_kyc_verified",
    value: status,
    path: "/",
    httpOnly: false,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    maxAge: 60 * 60 * 12,
  });
  return response;
}
