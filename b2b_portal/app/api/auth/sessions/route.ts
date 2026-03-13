import { NextResponse } from "next/server";
import { makeDirectHttp } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { resolveUserIdFromSession } from "@lib/auth/resolve-user-id";

/** GET /api/auth/sessions -- list all active sessions for current user */
export async function GET(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  const userId = await resolveUserIdFromSession(request, hdrs);
  if (!userId) return NextResponse.json({ ok: false, message: "Cannot resolve user identity" }, { status: 401 });
  // SDK query type doesn't include active_only — use direct HTTP to pass it correctly.
  // active_only=true filters revoked/expired sessions at the DB level.
  const http = makeDirectHttp(request, { ...hdrs, userId });
  const result = await http.get(`/v1/auth/users/${userId}/sessions?active_only=true`);
  if (!result.ok) return NextResponse.json({ ok: false, message: (result.data?.message as string) ?? "Failed to list sessions" }, { status: result.status });
  return NextResponse.json({ ok: true, sessions: result.data }, { status: 200 });
}

/** DELETE /api/auth/sessions -- revoke ALL sessions (sign out everywhere) */
export async function DELETE(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  const userId = await resolveUserIdFromSession(request, hdrs);
  if (!userId) return NextResponse.json({ ok: false, message: "Cannot resolve user identity" }, { status: 401 });
  // SDK generates POST /v1/auth/users/{id}/sessions:revoke-all (kebab) but the gateway
  // route is POST /v1/auth/users/{id}/sessions:revokeAll (camelCase) — use direct HTTP.
  const http = makeDirectHttp(request, { ...hdrs, userId });
  const result = await http.post(`/v1/auth/users/${userId}/sessions:revokeAll`, {
    user_id: userId,
    exclude_current_session: false,
  });
  if (!result.ok) return NextResponse.json({ ok: false, message: (result.data?.message as string) ?? "Failed to revoke all sessions" }, { status: result.status });
  return NextResponse.json({ ok: true, message: "All sessions revoked" }, { status: 200 });
}
