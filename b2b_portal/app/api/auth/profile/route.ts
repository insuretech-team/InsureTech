import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, badRequest } from "@lib/sdk/api-helpers";
import { resolveUserIdFromSession } from "@lib/auth/resolve-user-id";

// Reads a cookie value from a raw Cookie header string.
function extractCookieValue(cookieHeader: string, name: string): string {
  const m = cookieHeader.match(new RegExp(`(?:^|;\\s*)${name}=([^;]*)`));
  return m ? decodeURIComponent(m[1]) : "";
}

/** GET /api/auth/profile -- get current user profile */
export async function GET(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  const userId = await resolveUserIdFromSession(request, hdrs);
  if (!userId) return NextResponse.json({ ok: false, message: "Cannot resolve user identity" }, { status: 401 });
  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.getUserProfile({ path: { user_id: userId } });

  // Read user identity fields from lightweight cookies set at login.
  // mobile_number and email live on the User record, not UserProfile — they are
  // auth credentials and are read-only in the My Profile form.
  const cookieHeader = request.headers.get("cookie") ?? "";
  const mobile_number = extractCookieValue(cookieHeader, "portal_mobile");
  const email         = extractCookieValue(cookieHeader, "portal_email");

  // 404 means the user exists but has no profile row yet (new user).
  // Return identity fields so the form can still display mobile/email.
  if (result.response.status === 404) {
    return NextResponse.json({ ok: true, profile: { mobile_number, email } }, { status: 200 });
  }
  if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });

  // result.data is UserProfileRetrievalResponse = { profile?: UserProfile, error?: Error }
  // Unwrap the nested profile object — do NOT spread result.data directly or the
  // form state will contain a "profile" key which gets sent back on PATCH.
  const responseData = result.data as Record<string, unknown> ?? {};
  const raw = (responseData.profile as Record<string, unknown>) ?? {};

  // date_of_birth: the DB stores it as a DATE, the gateway serialises it via
  // protojson as an RFC3339 timestamp string (e.g. "1990-06-15T00:00:00Z").
  // Convert to YYYY-MM-DD for the HTML date input.
  // Treat zero-epoch / sentinel values (year ≤ 1970 or year === 1900) as empty
  // so new users see a blank date field rather than a confusing placeholder date.
  let dateOfBirth = "";
  const rawDOB = raw.date_of_birth as string | undefined;
  if (rawDOB) {
    const d = new Date(rawDOB);
    if (!isNaN(d.getTime()) && d.getFullYear() > 1970 && d.getFullYear() !== 1900) {
      dateOfBirth = d.toISOString().slice(0, 10); // "YYYY-MM-DD"
    }
  }

  // address_line1 is the proto field name — expose it as both address_line1 and
  // address so the form's single "Address" input maps correctly.
  const address_line1 = (raw.address_line1 as string) ?? "";

  const profile = {
    ...raw,
    date_of_birth: dateOfBirth,
    address_line1,
    // Convenience alias used by the form's "Address" field
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
  const userId = await resolveUserIdFromSession(request, hdrs);
  if (!userId) return NextResponse.json({ ok: false, message: "Cannot resolve user identity" }, { status: 401 });

  // Build a clean payload for protojson.Unmarshal in the gateway:
  // 1. user_id is required by UpdateUserProfileRequest.
  // 2. date_of_birth must be RFC3339 (gateway uses protojson) — HTML date gives "YYYY-MM-DD".
  // 3. address (form convenience alias) maps to address_line1 (proto field name).
  // 4. Strip read-only identity fields (email, mobile_number) — User record only.
  // 5. Strip the "address" alias and any stale "profile" nesting from old form state.
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

  // Convert YYYY-MM-DD → RFC3339 for protojson Timestamp parsing.
  // Validate the 18-year age requirement (DB CHECK constraint) before hitting gateway.
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
    transformed.date_of_birth = dob; // already RFC3339
  }
  // If no date_of_birth provided, omit it so the gateway skips the field.

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.updateUserProfile({
    path: { user_id: userId },
    body: transformed as Parameters<typeof sdk.updateUserProfile>[0]['body'],
  });
  if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
  return NextResponse.json({ ok: true, profile: result.data }, { status: 200 });
}
