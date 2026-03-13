import { NextResponse } from "next/server";
import { makeDirectHttp } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { resolveUserIdFromSession } from "@lib/auth/resolve-user-id";

/** DELETE /api/auth/sessions/[sessionId] — revoke a specific session */
export async function DELETE(
  request: Request,
  { params }: { params: Promise<{ sessionId: string }> }
) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  const { sessionId } = await params;
  if (!sessionId) return NextResponse.json({ ok: false, message: "sessionId is required" }, { status: 400 });
  const userId = await resolveUserIdFromSession(request, hdrs);
  if (!userId) return NextResponse.json({ ok: false, message: "Cannot resolve user identity" }, { status: 401 });
  // SDK generates DELETE /v1/auth/users/{user_id}/sessions/{session_id} but the gateway
  // route is DELETE /v1/auth/sessions/{session_id} (no user_id in path) — use direct HTTP.
  const http = makeDirectHttp(request, { ...hdrs, userId });
  const result = await http.delete(`/v1/auth/sessions/${sessionId}`);
  if (!result.ok) return NextResponse.json({ ok: false, message: (result.data?.message as string) ?? "Failed to revoke session" }, { status: result.status });
  return NextResponse.json({ ok: true, message: "Session revoked" }, { status: 200 });
}
