import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/partner-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, badRequest } from "@lib/sdk/api-helpers";

// Reads a cookie value from a raw Cookie header string.
function extractCookieValue(cookieHeader: string, name: string): string {
  const m = cookieHeader.match(new RegExp(`(?:^|;\\s*)${name}=([^;]*)`));
  return m ? decodeURIComponent(m[1]) : "";
}

async function resolveUserIdFromSession(request: Request): Promise<string> {
  const cookieHeader = request.headers.get("cookie") ?? "";
  return extractCookieValue(cookieHeader, "portal_user_id");
}

/** GET /api/auth/profile -- get current user profile */
export async function GET(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  const userId = await resolveUserIdFromSession(request);
  if (!userId) return NextResponse.json({ ok: false, message: "Cannot resolve user identity" }, { status: 401 });
  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.getUserProfile({ path: { user_id: userId } });

  const cookieHeader = request.headers.get("cookie") ?? "";
  const mobile_number = extractCookieValue(cookieHeader, "portal_mobile");
  const email         = extractCookieValue(cookieHeader, "portal_email");

  if (result.response.status === 404) {
    return NextResponse.json({ ok: true, profile: { mobile_number, email } }, { status: 200 });
  }
  if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });

  const responseData = result.data as Record<string, unknown> ?? {};
  const raw = (responseData.profile as Record<string, unknown>) ?? {};

  let dateOfBirth = "";
  const rawDOB = raw.date_of_birth as string | undefined;
  if (rawDOB) {
    const d = new Date(rawDOB);
    if (!isNaN(d.getTime()) && d.getFullYear() > 1970 && d.getFullYear() !== 1900) {
      dateOfBirth = d.toISOString().slice(0, 10);
    }
  }

  const address_line1 = (raw.address_line1 as string) ?? "";

  const profile = {
    ...raw,
    date_of_birth: dateOfBirth,
    address_line1,
    address: address_line1,
    mobile_number,
    email,
  };
  return NextResponse.json({ ok: true, profile }, { status: 200 });
}

/** PATCH /api/auth/profile -- update current user profile */
export async function PATCH(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  let body: Record<string, unknown>;
  try { body = await request.json() as Record<string, unknown>; } catch { return badRequest("Invalid request body"); }
  const userId = await resolveUserIdFromSession(request);
  if (!userId) return NextResponse.json({ ok: false, message: "Cannot resolve user identity" }, { status: 401 });

  const transformed: Record<string, unknown> = {
    user_id: userId,
    full_name:      body.full_name      ?? "",
    occupation:     body.occupation     ?? "",
    employer:       body.employer       ?? "",
    address_line1:  body.address_line1  ?? body.address ?? "",
    address_line2:  body.address_line2  ?? "",
    city:           body.city           ?? "",
    district:       body.district       ?? "",
    division:       body.division       ?? "",
    country:        body.country        ?? "",
    postal_code:    body.postal_code    ?? "",
    nid_number:     body.nid_number     ?? "",
    marital_status: body.marital_status ?? "",
    gender:         body.gender         ?? "",
  };

  const dob = body.date_of_birth as string | undefined;
  if (dob && /^\d{4}-\d{2}-\d{2}$/.test(dob)) {
    const dobDate = new Date(dob);
    const minAge = new Date();
    minAge.setFullYear(minAge.getFullYear() - 18);
    if (dobDate > minAge) {
      return NextResponse.json({ ok: false, message: "Date of birth must be at least 18 years in the past." }, { status: 400 });
    }
    transformed.date_of_birth = `${dob}T00:00:00Z`;
  } else if (dob && dob.includes("T")) {
    transformed.date_of_birth = dob;
  }

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.updateUserProfile({
    path: { user_id: userId },
    body: transformed as Parameters<typeof sdk.updateUserProfile>[0]['body'],
  });
  if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
  return NextResponse.json({ ok: true, profile: result.data }, { status: 200 });
}
