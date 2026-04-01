import { cookies } from "next/headers";
import { redirect } from "next/navigation";

import { getSession, createSession } from "./session-store";
import { toPortalSessionFromCurrentSession } from "./backend-auth";
import type { PortalPrincipal } from "@lib/types/auth";

export const SESSION_COOKIE_NAME = "session_token";

/**
 * Returns the portal session for the current request.
 *
 * Strategy (two-tier):
 *  1. Fast path — check the in-memory session store (valid within the same
 *     Node.js process lifetime, e.g. after a fresh login in the same process).
 *  2. Fallback — if the store is empty (e.g. after a dev HMR reload or a cold
 *     server start), validate the session_token cookie against the backend
 *     gateway directly and re-hydrate the store so subsequent calls are fast.
 *
 * This prevents the infinite redirect loop that occurred in dev mode:
 *   HMR wipes in-memory store → requireServerSession() always returns null
 *   → redirect("/login") → middleware sees cookie → redirect("/") →
 *   portal_kyc_verified=false → redirect("/kyc") → repeat.
 */
export async function getServerSession() {
  const cookieStore = await cookies();
  const sessionToken = cookieStore.get(SESSION_COOKIE_NAME)?.value;
  if (!sessionToken) return null;
  const cookieHeader = cookieStore
    .getAll()
    .map((cookie) => `${cookie.name}=${encodeURIComponent(cookie.value)}`)
    .join("; ");

  // 1. Fast path — in-memory store hit
  const stored = getSession(sessionToken);
  if (stored) return stored;

  // 2. Fallback — validate against backend gateway and re-hydrate store
  try {
    const { makeSdkClient } = await import("@lib/sdk/b2b-sdk-client");
    const fakeReq = new Request("http://localhost/api/session-hydrate", {
      headers: { cookie: cookieHeader },
    });
    const sdk = makeSdkClient(fakeReq);
    const result = await sdk.getCurrentSession();
    if (!result.response.ok || !result.data) return null;

    const portalSession = await toPortalSessionFromCurrentSession(result.data, cookieHeader);
    if (!portalSession) return null;

    // Re-hydrate the store using the backend session token as the key
    // so that subsequent in-process calls hit the fast path.
    const rehydrated = createSession(portalSession.principal, sessionToken);
    return rehydrated;
  } catch {
    return null;
  }
}

export async function requireServerSession() {
  const session = await getServerSession();
  if (!session) {
    redirect("/login");
  }
  return session;
}

export async function requireServerSessionRole(
  allowedRoles: PortalPrincipal["role"][],
  fallbackPath = "/"
) {
  const session = await requireServerSession();
  if (!allowedRoles.includes(session.principal.role)) {
    redirect(fallbackPath);
  }
  return session;
}
